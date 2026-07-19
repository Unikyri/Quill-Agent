package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/quill/backend/internal/config"
)

func newDashScopeTestService(t *testing.T, serverURL string) *DashScopeService {
	t.Helper()
	cfg := &config.Config{
		QwenAPIKey:            "test-key",
		QwenNativeBaseURL:     serverURL,
		QwenExtractionModel:   "extract-model",
		QwenReasoningModel:    "reason-model",
		QwenEmbeddingModel:    "embed-model",
		QwenEmbeddingDims:     3,
		QwenMaxConcurrency:    2,
		QwenTurboConcurrency:  2,
		LLMMaxConcurrency:     2,
		LLMTPMTurbo:           1_000_000,
		LLMTPMMax:             1_000_000,
		LLMRPM:                10_000,
		LLMInteractiveReserve: 0.30,
		LLMRampStep:           1,
		QwenRetryMaxAttempts:  2,
		QwenAPITimeout:        5 * time.Second,
	}
	svc := NewDashScopeService(cfg, nil)
	// Retry tests must never wait on exponential backoff.
	svc.retrySleep = func(context.Context, time.Duration) error { return nil }
	svc.jitter = func(delay time.Duration) time.Duration { return delay }
	return svc
}

func writeNativeChatResponse(w http.ResponseWriter, content string, toolCalls []QwenToolCall, input, output, cached int) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status_code": 200,
		"output": map[string]interface{}{
			"choices": []interface{}{map[string]interface{}{
				"finish_reason": "stop",
				"message": map[string]interface{}{
					"role":       "assistant",
					"content":    content,
					"tool_calls": toolCalls,
				},
			}},
		},
		"usage": map[string]interface{}{
			"input_tokens":  input,
			"output_tokens": output,
			"total_tokens":  input + output,
			"prompt_tokens_details": map[string]interface{}{
				"cached_tokens": cached,
			},
		},
	})
}

func TestDashScopeServiceChatNativeRequestAndUsage(t *testing.T) {
	var gotPath, gotAuth string
	var got map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotAuth = r.URL.Path, r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("decode request: %v", err)
			return
		}
		writeNativeChatResponse(w, "native answer", nil, 12, 4, 7)
	}))
	defer server.Close()

	svc := newDashScopeTestService(t, server.URL)
	answer, err := svc.Chat(context.Background(), "extract-model", []QwenMessage{{Role: "user", Content: "hello"}})
	if err != nil {
		t.Fatalf("Chat: %v", err)
	}
	if answer != "native answer" {
		t.Fatalf("answer = %q, want native answer", answer)
	}
	if gotPath != "/api/v1/services/aigc/text-generation/generation" {
		t.Fatalf("path = %q, want native generation endpoint", gotPath)
	}
	if gotAuth != "Bearer test-key" {
		t.Fatalf("authorization = %q, want bearer header", gotAuth)
	}
	if got["model"] != "extract-model" {
		t.Fatalf("model = %#v, want extract-model", got["model"])
	}
	input := got["input"].(map[string]interface{})
	messages := input["messages"].([]interface{})
	if len(messages) != 1 || messages[0].(map[string]interface{})["role"] != "user" {
		t.Fatalf("native input.messages = %#v", messages)
	}
	params := got["parameters"].(map[string]interface{})
	if params["result_format"] != "message" {
		t.Fatalf("result_format = %#v, want message", params["result_format"])
	}
	usage := svc.UsageSnapshot()
	if usage.InputTokens != 12 || usage.OutputTokens != 4 || usage.CachedTokens != 7 || usage.Requests != 1 {
		t.Fatalf("usage = %+v, want input=12 output=4 cached=7 requests=1", usage)
	}
}

func TestDashScopeServiceAgentToolRoundTrip(t *testing.T) {
	var round atomic.Int32
	var secondMessages []QwenMessage
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body dashScopeChatRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode native request: %v", err)
			return
		}
		if round.Add(1) == 1 {
			if len(body.Parameters.Tools) != 1 || body.Parameters.ToolChoice != "auto" {
				t.Errorf("tool parameters = %+v, want one tool and auto choice", body.Parameters)
			}
			writeNativeChatResponse(w, "", []QwenToolCall{{ID: "call-1", Type: "function", Function: QwenToolCallFunction{Name: "lookup", Arguments: `{"query":"dragon"}`}}}, 5, 2, 0)
			return
		}
		secondMessages = body.Input.Messages
		writeNativeChatResponse(w, "final native answer", nil, 9, 3, 0)
	}))
	defer server.Close()

	svc := newDashScopeTestService(t, server.URL)
	answer, err := svc.RunAgentLoop(context.Background(), []QwenMessage{{Role: "user", Content: "find dragon"}}, []QwenTool{{Type: "function", Function: QwenToolFunction{Name: "lookup", Description: "look up", Parameters: map[string]interface{}{}}}}, &recordingDashScopeExecutor{}, 3)
	if err != nil {
		t.Fatalf("RunAgentLoop: %v", err)
	}
	if answer != "final native answer" || round.Load() != 2 {
		t.Fatalf("answer=%q rounds=%d, want final answer after two rounds", answer, round.Load())
	}
	if len(secondMessages) != 3 || secondMessages[1].Role != "assistant" || secondMessages[2].Role != "tool" || secondMessages[2].ToolCallID != "call-1" {
		t.Fatalf("second native messages = %+v, want assistant/tool round trip", secondMessages)
	}
}

func TestDashScopeStructuredSchemasAndLocalEntityValidation(t *testing.T) {
	var body map[string]interface{}
	content := `{"characters":[{"name":"Mira","aliases":[],"type":"character","status":"active","description":"pilot","properties":{}},{"name":"bad","aliases":[],"type":"dragon","status":"active","description":"invalid","properties":{}}],"places":[],"objects":[],"events":[],"factions":[],"world_rules":[],"plot_developments":[]}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode native extraction request: %v", err)
		}
		writeNativeChatResponse(w, content, nil, 1, 1, 0)
	}))
	defer server.Close()
	svc := newDashScopeTestService(t, server.URL)
	entities, err := svc.ExtractEntities(context.Background(), "Mira arrives.", "")
	if err != nil {
		t.Fatalf("ExtractEntities: %v", err)
	}
	if len(entities.Characters) != 1 || entities.Characters[0].Type != "character" {
		t.Fatalf("validated characters=%+v, want only canonical character", entities.Characters)
	}
	parameters := body["parameters"].(map[string]interface{})
	format := parameters["response_format"].(map[string]interface{})
	if format["type"] != "json_schema" {
		t.Fatalf("response format=%+v, want strict json_schema", format)
	}
	schemaDescriptor := format["json_schema"].(map[string]interface{})
	if schemaDescriptor["strict"] != true {
		t.Fatalf("schema descriptor=%+v, want strict=true", schemaDescriptor)
	}
	root := schemaDescriptor["schema"].(map[string]interface{})
	properties := root["properties"].(map[string]interface{})
	characterItems := properties["characters"].(map[string]interface{})["items"].(map[string]interface{})
	typeSchema := characterItems["properties"].(map[string]interface{})["type"].(map[string]interface{})
	enum := typeSchema["enum"].([]interface{})
	if len(enum) != 7 {
		t.Fatalf("entity type enum length=%d, want 7", len(enum))
	}
	messages := body["input"].(map[string]interface{})["messages"].([]interface{})
	if len(messages) != 2 {
		t.Fatalf("messages=%d, want stable system + variable paragraph", len(messages))
	}
	system := messages[0].(map[string]interface{})
	blocks := system["content"].([]interface{})
	if len(blocks) != 1 {
		t.Fatalf("stable extraction prefix=%#v, want lore and deterministic cache padding", blocks)
	}
	stableText := blocks[0].(map[string]interface{})["text"].(string)
	if !strings.Contains(stableText, "UNIVERSE LORE") || !strings.Contains(stableText, "stable context") || len(stableText) < 4096 {
		t.Fatalf("stable extraction prefix=%#v, want lore and deterministic cache padding", blocks[0])
	}
	if blocks[0].(map[string]interface{})["cache_control"].(map[string]interface{})["type"] != "ephemeral" {
		t.Fatalf("stable prefix cache marker=%#v", blocks[0])
	}
	if messages[1].(map[string]interface{})["content"] != "Mira arrives." {
		t.Fatalf("variable extraction message=%#v, want paragraph only", messages[1])
	}
}

func TestDashScopeStructuredRelationshipSchemaAndValidation(t *testing.T) {
	var body map[string]interface{}
	content := `{"relationships":[{"source":"Mira","target":"Aurelia","type":"LOCATED_IN","properties":{}},{"source":"Mira","target":"Aurelia","type":"NOT_AN_EDGE","properties":{}}]}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode native relationship request: %v", err)
		}
		writeNativeChatResponse(w, content, nil, 1, 1, 0)
	}))
	defer server.Close()
	svc := newDashScopeTestService(t, server.URL)
	relationships, err := svc.AnalyzeRelationships(context.Background(), "Mira arrives.", []string{"Mira", "Aurelia"})
	if err != nil {
		t.Fatalf("AnalyzeRelationships: %v", err)
	}
	if len(relationships) != 1 || relationships[0]["type"] != "LOCATED_IN" {
		t.Fatalf("validated relationships=%+v, want one migration edge", relationships)
	}
	parameters := body["parameters"].(map[string]interface{})
	format := parameters["response_format"].(map[string]interface{})
	if format["type"] != "json_schema" {
		t.Fatalf("response format=%+v, want strict json_schema", format)
	}
	schemaDescriptor := format["json_schema"].(map[string]interface{})
	if schemaDescriptor["strict"] != true {
		t.Fatalf("schema descriptor=%+v, want strict=true", schemaDescriptor)
	}
	relProps := schemaDescriptor["schema"].(map[string]interface{})["properties"].(map[string]interface{})["relationships"].(map[string]interface{})
	itemProps := relProps["items"].(map[string]interface{})["properties"].(map[string]interface{})
	enum := itemProps["type"].(map[string]interface{})["enum"].([]interface{})
	if len(enum) != 16 {
		t.Fatalf("relationship type enum length=%d, want migration + legacy compatibility set", len(enum))
	}
}

type recordingDashScopeExecutor struct{}

func (*recordingDashScopeExecutor) ExecuteTool(name, args string) (string, error) {
	if name != "lookup" || !strings.Contains(args, "dragon") {
		return "", fmt.Errorf("unexpected tool call %s %s", name, args)
	}
	return "dragon memory", nil
}

func TestDashScopeServiceEmbeddingNativeRequest(t *testing.T) {
	var got dashScopeEmbeddingRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/services/embeddings/text-embedding/text-embedding" {
			t.Errorf("path = %q, want native embedding endpoint", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("decode embedding request: %v", err)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status_code": 200,
			"output": map[string]interface{}{"embeddings": []interface{}{
				map[string]interface{}{"embedding": []float32{0.1, 0.2, 0.3}, "text_index": 0},
			}},
			"usage": map[string]interface{}{"total_tokens": 3},
		})
	}))
	defer server.Close()

	svc := newDashScopeTestService(t, server.URL)
	result, err := svc.GenerateEmbeddingBatch(context.Background(), []string{"a paragraph"})
	if err != nil {
		t.Fatalf("GenerateEmbeddingBatch: %v", err)
	}
	if len(result) != 1 || len(result[0]) != 3 {
		t.Fatalf("embedding result = %#v", result)
	}
	if got.Model != "embed-model" || len(got.Input.Texts) != 1 || got.Parameters.Dimension != 3 || got.Parameters.OutputType != "dense" {
		t.Fatalf("embedding request = %+v", got)
	}
}

func TestDashScopeServiceEmbeddingBatchChunksAtTen(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		var request dashScopeEmbeddingRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("decode embedding request: %v", err)
			return
		}
		embeddings := make([]interface{}, len(request.Input.Texts))
		for i := range request.Input.Texts {
			embeddings[i] = map[string]interface{}{"embedding": []float32{0.1, 0.2, 0.3}, "text_index": i}
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status_code": 200,
			"output":      map[string]interface{}{"embeddings": embeddings},
		})
	}))
	defer server.Close()

	texts := make([]string, 11)
	for i := range texts {
		texts[i] = fmt.Sprintf("paragraph-%d", i)
	}
	svc := newDashScopeTestService(t, server.URL)
	result, err := svc.GenerateEmbeddingBatch(context.Background(), texts)
	if err != nil {
		t.Fatalf("GenerateEmbeddingBatch: %v", err)
	}
	if calls.Load() != 2 || len(result) != len(texts) {
		t.Fatalf("calls=%d result_len=%d, want two requests and eleven embeddings", calls.Load(), len(result))
	}
}

func TestDashScopeResponseRejectsBodyStatusAndCodeErrors(t *testing.T) {
	cases := []string{
		`{"status_code":400,"message":"bad request"}`,
		`{"status_code":401,"code":"Unauthorized","message":"bad key"}`,
		`{"code":"Model.NotFound","message":"missing model"}`,
	}
	for _, body := range cases {
		if _, _, err := parseDashScopeResponse([]byte(body)); err == nil {
			t.Errorf("parseDashScopeResponse(%s) returned nil error", body)
		}
	}
}

func TestDashScopeHealthCheckRejectsBodyStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"status_code": 401, "message": "invalid key"})
	}))
	defer server.Close()
	svc := newDashScopeTestService(t, server.URL)
	if err := svc.HealthCheck(context.Background()); err == nil {
		t.Fatal("HealthCheck accepted an error status in a successful HTTP body")
	}
}

func TestDashScopeServiceRetries429(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = io.WriteString(w, `{"code":"Throttled","message":"try again"}`)
			return
		}
		writeNativeChatResponse(w, "after retry", nil, 2, 1, 0)
	}))
	defer server.Close()

	svc := newDashScopeTestService(t, server.URL)
	answer, err := svc.Chat(context.Background(), "extract-model", []QwenMessage{{Role: "user", Content: "retry"}})
	if err != nil {
		t.Fatalf("Chat after retry: %v", err)
	}
	if answer != "after retry" || calls.Load() != 2 {
		t.Fatalf("answer=%q calls=%d, want retry then success", answer, calls.Load())
	}
}

func TestDashScopeSendNativeRequestHonorsModelConcurrencyGate(t *testing.T) {
	firstStarted := make(chan struct{})
	allowFirst := make(chan struct{})
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			close(firstStarted)
			<-allowFirst
		}
		writeNativeChatResponse(w, "ok", nil, 1, 1, 0)
	}))
	defer server.Close()
	svc := newDashScopeTestService(t, server.URL)
	svc.maxSem = make(chan struct{}, 1)
	firstDone := make(chan error, 1)
	go func() {
		_, err := svc.Chat(context.Background(), "reason-model", []QwenMessage{{Role: "user", Content: "first"}})
		firstDone <- err
	}()
	select {
	case <-firstStarted:
	case <-time.After(time.Second):
		t.Fatal("first native request did not reach the server")
	}
	secondDone := make(chan error, 1)
	go func() {
		_, err := svc.Chat(context.Background(), "reason-model", []QwenMessage{{Role: "user", Content: "second"}})
		secondDone <- err
	}()
	select {
	case err := <-secondDone:
		t.Fatalf("second request bypassed maxSem gate: %v", err)
	case <-time.After(100 * time.Millisecond):
	}
	if calls.Load() != 1 {
		t.Fatalf("server calls=%d while first held gate, want 1", calls.Load())
	}
	close(allowFirst)
	if err := <-firstDone; err != nil {
		t.Fatalf("first request: %v", err)
	}
	if err := <-secondDone; err != nil {
		t.Fatalf("second request after release: %v", err)
	}
	if calls.Load() != 2 {
		t.Fatalf("server calls=%d, want two serialized requests", calls.Load())
	}
}

func TestDashScopeSendNativeRequestReleasesGateAcrossFallback(t *testing.T) {
	var models []string
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&request)
		models = append(models, request["model"].(string))
		if calls.Add(1) == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = io.WriteString(w, `{"code":"Throttled","message":"try again"}`)
			return
		}
		writeNativeChatResponse(w, "fallback answer", nil, 1, 1, 0)
	}))
	defer server.Close()
	svc := newDashScopeTestService(t, server.URL)
	svc.retryMaxAttempts = 1
	svc.fallbackOn429 = true
	svc.fallbackModel = "fallback-model"
	svc.maxSem = make(chan struct{}, 1)
	answer, err := svc.Chat(context.Background(), "reason-model", []QwenMessage{{Role: "user", Content: "fallback"}})
	if err != nil || answer != "fallback answer" {
		t.Fatalf("fallback answer=%q err=%v", answer, err)
	}
	if len(models) != 2 || models[0] != "reason-model" || models[1] != "fallback-model" {
		t.Fatalf("models=%v, want primary then fallback after gate release", models)
	}
}

func TestDashScopeServiceRerankNativeRequest(t *testing.T) {
	var request dashScopeRerankRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/services/rerank/text-rerank/text-rerank" {
			t.Errorf("path=%q, want native rerank endpoint", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("decode rerank request: %v", err)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status_code": 200,
			"output": map[string]interface{}{"results": []interface{}{
				map[string]interface{}{"index": 1, "relevance_score": 0.9},
				map[string]interface{}{"index": 0, "relevance_score": 0.2},
			}},
		})
	}))
	defer server.Close()
	svc := newDashScopeTestService(t, server.URL)
	svc.rerankModel = "rerank-model"
	results, err := svc.Rerank(context.Background(), "dragon", []string{"quiet lake", "dragon attack"}, 2)
	if err != nil {
		t.Fatalf("Rerank: %v", err)
	}
	if request.Model != "rerank-model" || request.Input.Query != "dragon" || len(request.Input.Documents) != 2 || request.Parameters.TopN != 2 || request.Parameters.ReturnDocuments {
		t.Fatalf("rerank request=%+v", request)
	}
	if len(results) != 2 || results[0].Index != 1 || results[0].Score != 0.9 {
		t.Fatalf("rerank results=%+v", results)
	}
}

func TestDashScopeServiceStreamingSetsNativeSSEHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/services/aigc/text-generation/generation" || r.Header.Get("X-DashScope-SSE") != "enable" {
			t.Errorf("stream request path/header = %q/%q", r.URL.Path, r.Header.Get("X-DashScope-SSE"))
		}
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, _ := w.(http.Flusher)
		_, _ = io.WriteString(w, "data: {\"choices\":[{\"delta\":{\"content\":\"native \"},\"finish_reason\":null}]}\n\n")
		if flusher != nil {
			flusher.Flush()
		}
		_, _ = io.WriteString(w, "data: {\"choices\":[{\"delta\":{\"content\":\"stream\"},\"finish_reason\":\"stop\"}]}\n\n")
		_, _ = io.WriteString(w, "data: [DONE]\n\n")
	}))
	defer server.Close()

	svc := newDashScopeTestService(t, server.URL)
	stream, err := svc.ChatCompletionStream(context.Background(), QwenRequest{Model: "reason-model", Messages: []QwenMessage{{Role: "user", Content: "stream"}}})
	if err != nil {
		t.Fatalf("ChatCompletionStream: %v", err)
	}
	var text strings.Builder
	var sawDone bool
	for chunk := range stream {
		if chunk.Type == "text" {
			text.WriteString(chunk.Text)
		}
		if chunk.Type == "done" {
			sawDone = true
		}
		if chunk.Type == "error" {
			t.Fatalf("stream error: %v", chunk.Err)
		}
	}
	if text.String() != "native stream" || !sawDone {
		t.Fatalf("stream text=%q done=%v", text.String(), sawDone)
	}
}

func TestDashScopeServiceStreamingParsesNativeOutputToolCallsAndUsage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, "data: {\"output\":{\"choices\":[{\"message\":{\"content\":\"native \"}}]}}\n\n")
		_, _ = io.WriteString(w, "data: {\"output\":{\"choices\":[{\"message\":{\"tool_calls\":[{\"id\":\"call-native\",\"function\":{\"name\":\"lookup\",\"arguments\":\"{\\\"q\\\":\\\"x\\\"}\"}}]},\"finish_reason\":\"tool_calls\"}]},\"usage\":{\"input_tokens\":8,\"output_tokens\":2,\"prompt_tokens_details\":{\"cached_tokens\":5,\"cache_creation_input_tokens\":1}}}\n\n")
		_, _ = io.WriteString(w, "data: [DONE]\n\n")
	}))
	defer server.Close()
	svc := newDashScopeTestService(t, server.URL)
	stream, err := svc.ChatCompletionStream(context.Background(), QwenRequest{Model: "reason-model", Messages: []QwenMessage{{Role: "user", Content: "stream"}}})
	if err != nil {
		t.Fatalf("ChatCompletionStream: %v", err)
	}
	var text strings.Builder
	var calls []QwenToolCall
	for chunk := range stream {
		switch chunk.Type {
		case "text":
			text.WriteString(chunk.Text)
		case "tool_call":
			calls = append(calls, *chunk.ToolCall)
		case "error":
			t.Fatalf("stream error: %v", chunk.Err)
		}
	}
	if text.String() != "native " || len(calls) != 1 || calls[0].Function.Name != "lookup" {
		t.Fatalf("text=%q calls=%+v", text.String(), calls)
	}
	usage := svc.UsageSnapshot()
	if usage.InputTokens != 8 || usage.OutputTokens != 2 || usage.CachedTokens != 5 || usage.CacheCreationInputTokens != 1 || usage.Requests != 1 {
		t.Fatalf("stream usage=%+v", usage)
	}
}

func TestDashScopeCacheMarkerUsesNativeContentBlock(t *testing.T) {
	var body map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request: %v", err)
		}
		writeNativeChatResponse(w, "ok", nil, 1, 1, 1)
	}))
	defer server.Close()
	svc := newDashScopeTestService(t, server.URL)
	_, err := svc.Chat(context.Background(), "reason-model", []QwenMessage{
		{Role: "system", Content: "stable instructions"},
		{Role: "user", Content: "variable prompt"},
	})
	if err != nil {
		t.Fatalf("Chat: %v", err)
	}
	messages := body["input"].(map[string]interface{})["messages"].([]interface{})
	first := messages[0].(map[string]interface{})
	blocks, ok := first["content"].([]interface{})
	if !ok || len(blocks) != 1 {
		t.Fatalf("first native message content=%#v, want one content block", first["content"])
	}
	block := blocks[0].(map[string]interface{})
	if block["cache_control"].(map[string]interface{})["type"] != "ephemeral" {
		t.Fatalf("cache_control=%#v", block["cache_control"])
	}
}

func TestDashScopeCacheBatchGuardRejectsMarkedMessages(t *testing.T) {
	marked := []QwenMessage{{Role: "system", Content: "stable", CacheControl: &CacheControl{Type: "ephemeral"}}}
	if err := ensureDashScopeCacheBatchCompatible(marked, true); err == nil {
		t.Fatal("cache-marked messages must be rejected for Batch requests")
	}
	if err := ensureDashScopeCacheBatchCompatible(marked, false); err != nil {
		t.Fatalf("non-Batch native request should allow cache marker: %v", err)
	}
}

// TestCompressDashScopeToolResultsSummaryIsAssistantRole guards against a
// regression where the compressed context was injected as Role: "tool".
// DashScope's native API rejects a "tool" message that doesn't directly
// follow an assistant message with tool_calls (InvalidParameter: messages
// with role "tool" must be a response to a preceding message with
// "tool_calls") — the compressed head here is just [system, user], so the
// summary must be an assistant message instead.
func TestCompressDashScopeToolResultsSummaryIsAssistantRole(t *testing.T) {
	// Tiny budget forces the 80%-usage threshold to trip immediately.
	budgetMgr := NewContextBudgetManager(NewTokenizer(), 10, 0)

	msgs := []QwenMessage{
		{Role: "system", Content: "sys"},
		{Role: "user", Content: "usr"},
		{Role: "assistant", ToolCalls: []QwenToolCall{{ID: "1", Function: QwenToolCallFunction{Name: "f"}}}},
		{Role: "tool", ToolCallID: "1", Content: "result A"},
		{Role: "assistant", ToolCalls: []QwenToolCall{{ID: "2", Function: QwenToolCallFunction{Name: "f"}}}},
		{Role: "tool", ToolCallID: "2", Content: "result B"},
	}

	summarizeCalled := false
	compressedMsgs, compressed := compressDashScopeToolResults(context.Background(), budgetMgr, msgs, func(_ context.Context, prompt string) (string, error) {
		summarizeCalled = true
		if !strings.Contains(prompt, "result A") {
			t.Errorf("summarize prompt = %q, want it to include the compressed tool result", prompt)
		}
		return "summary text", nil
	})
	if !compressed {
		t.Fatal("compressDashScopeToolResults did not attempt compression — test setup did not exceed threshold")
	}
	if !summarizeCalled {
		t.Fatal("summarize callback was not invoked")
	}

	for i, msg := range compressedMsgs {
		if msg.Role == "tool" {
			if i == 0 || compressedMsgs[i-1].Role != "assistant" || len(compressedMsgs[i-1].ToolCalls) == 0 {
				t.Errorf("compressedMsgs[%d] has role \"tool\" without a preceding assistant tool_calls message", i)
			}
		}
	}

	foundSummary := false
	for _, msg := range compressedMsgs {
		if msg.Role == "assistant" && strings.Contains(msg.Content, "summary text") {
			foundSummary = true
		}
	}
	if !foundSummary {
		t.Error("compressedMsgs does not contain the summary as an assistant message")
	}
}

func TestDashScopeNativeBaseURLNormalizesCompatibleAndAPIV1Forms(t *testing.T) {
	for input, want := range map[string]string{
		"https://dashscope-intl.aliyuncs.com":                    "https://dashscope-intl.aliyuncs.com/api/v1",
		"https://dashscope-intl.aliyuncs.com/api/v1":             "https://dashscope-intl.aliyuncs.com/api/v1",
		"https://dashscope-intl.aliyuncs.com/compatible-mode/v1": "https://dashscope-intl.aliyuncs.com/api/v1",
	} {
		if got := nativeAPIBaseURL(input); got != want {
			t.Errorf("nativeAPIBaseURL(%q) = %q, want %q", input, got, want)
		}
	}
}
