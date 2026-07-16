package services

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/quill/backend/internal/config"
)

// newSSEServer serves the given raw SSE data lines (already prefixed with
// "data: " where needed) as a chat/completions streaming response, flushing
// after each line so the client-side scanner observes them incrementally.
func newSSEServer(t *testing.T, lines []string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatalf("ResponseWriter does not support flushing")
		}
		for _, line := range lines {
			fmt.Fprintf(w, "data: %s\n\n", line)
			flusher.Flush()
		}
	}))
}

func newStreamTestService(baseURL string) *QwenService {
	cfg := &config.Config{
		QwenBaseURL:          baseURL,
		QwenAPIKey:           "test-key",
		QwenMaxConcurrency:   1,
		QwenTurboConcurrency: 1,
	}
	return NewQwenService(cfg, nil)
}

// collectStream drains a StreamChunk channel with a timeout so a bug that
// forgets to close the channel fails the test instead of hanging forever.
func collectStream(t *testing.T, ch <-chan StreamChunk) []StreamChunk {
	t.Helper()
	var chunks []StreamChunk
	timeout := time.After(5 * time.Second)
	for {
		select {
		case chunk, ok := <-ch:
			if !ok {
				return chunks
			}
			chunks = append(chunks, chunk)
		case <-timeout:
			t.Fatal("timed out waiting for stream channel to close")
		}
	}
}

func TestChatCompletionStream_TextOnly(t *testing.T) {
	server := newSSEServer(t, []string{
		`{"choices":[{"delta":{"content":"Hello "},"finish_reason":null}]}`,
		`{"choices":[{"delta":{"content":"world"},"finish_reason":null}]}`,
		`{"choices":[{"delta":{},"finish_reason":"stop"}]}`,
		`[DONE]`,
	})
	defer server.Close()

	svc := newStreamTestService(server.URL)
	ch, err := svc.ChatCompletionStream(context.Background(), QwenRequest{Model: "qwen-max"})
	if err != nil {
		t.Fatalf("ChatCompletionStream: %v", err)
	}

	chunks := collectStream(t, ch)

	var text string
	var sawDone bool
	for _, c := range chunks {
		switch c.Type {
		case "text":
			text += c.Text
		case "done":
			sawDone = true
			if c.Finish != "stop" {
				t.Errorf("done chunk Finish = %q, want %q", c.Finish, "stop")
			}
		case "tool_call":
			t.Errorf("unexpected tool_call chunk in text-only stream: %+v", c)
		case "error":
			t.Errorf("unexpected error chunk: %v", c.Err)
		}
	}

	if text != "Hello world" {
		t.Errorf("accumulated text = %q, want %q", text, "Hello world")
	}
	if !sawDone {
		t.Error("expected a done chunk with finish_reason stop")
	}
}

func TestChatCompletionStream_IgnoresSSEMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = fmt.Fprint(w, "id: response-1\nevent: result\ndata: {\"choices\":[{\"delta\":{\"content\":\"hello\"},\"finish_reason\":\"stop\"}]}\n\ndata: [DONE]\n\n")
	}))
	defer server.Close()

	stream, err := newStreamTestService(server.URL).ChatCompletionStream(context.Background(), QwenRequest{Model: "qwen-max"})
	if err != nil {
		t.Fatalf("ChatCompletionStream: %v", err)
	}
	chunks := collectStream(t, stream)

	var text string
	for _, chunk := range chunks {
		if chunk.Type == "error" {
			t.Fatalf("stream error: %v", chunk.Err)
		}
		if chunk.Type == "text" {
			text += chunk.Text
		}
	}
	if text != "hello" {
		t.Fatalf("stream text=%q, want %q", text, "hello")
	}
}

func TestChatCompletionStream_LengthFinishIsTerminal(t *testing.T) {
	server := newSSEServer(t, []string{
		`{"choices":[{"delta":{"content":"partial"},"finish_reason":"length"}]}`,
		`[DONE]`,
	})
	defer server.Close()

	stream, err := newStreamTestService(server.URL).ChatCompletionStream(context.Background(), QwenRequest{Model: "qwen-max"})
	if err != nil {
		t.Fatalf("ChatCompletionStream: %v", err)
	}
	chunks := collectStream(t, stream)

	var text string
	var sawDone bool
	for _, chunk := range chunks {
		switch chunk.Type {
		case "text":
			text += chunk.Text
		case "done":
			sawDone = true
			if chunk.Finish != "length" {
				t.Errorf("done finish=%q, want length", chunk.Finish)
			}
		case "error":
			t.Fatalf("stream error: %v", chunk.Err)
		}
	}
	if text != "partial" || !sawDone {
		t.Fatalf("stream text=%q done=%v, want partial/true", text, sawDone)
	}
}

func TestChatCompletionStream_SingleToolCallFragmented(t *testing.T) {
	// Mirrors the real DashScope shape observed in the PR3 live spike:
	// name + id arrive once on the first fragment, arguments arrive as
	// concatenated string fragments across subsequent chunks, and the
	// final chunk (empty delta) carries finish_reason:"tool_calls".
	server := newSSEServer(t, []string{
		`{"choices":[{"delta":{"tool_calls":[{"index":0,"id":"call_abc","type":"function","function":{"name":"search_vector_memory","arguments":"{\"query\": \""}}]},"finish_reason":null}]}`,
		`{"choices":[{"delta":{"tool_calls":[{"index":0,"function":{"arguments":"dragon"}}]},"finish_reason":null}]}`,
		`{"choices":[{"delta":{"tool_calls":[{"index":0,"function":{"arguments":"\"}"}}]},"finish_reason":null}]}`,
		`{"choices":[{"delta":{},"finish_reason":"tool_calls"}]}`,
		`[DONE]`,
	})
	defer server.Close()

	svc := newStreamTestService(server.URL)
	ch, err := svc.ChatCompletionStream(context.Background(), QwenRequest{Model: "qwen-max"})
	if err != nil {
		t.Fatalf("ChatCompletionStream: %v", err)
	}

	chunks := collectStream(t, ch)

	var toolCalls []StreamChunk
	for _, c := range chunks {
		if c.Type == "error" {
			t.Fatalf("unexpected error chunk: %v", c.Err)
		}
		if c.Type == "tool_call" {
			toolCalls = append(toolCalls, c)
		}
	}

	if len(toolCalls) != 1 {
		t.Fatalf("got %d tool_call chunks, want exactly 1 (no partial dispatch)", len(toolCalls))
	}
	tc := toolCalls[0].ToolCall
	if tc == nil {
		t.Fatal("tool_call chunk has nil ToolCall")
	}
	if tc.ID != "call_abc" {
		t.Errorf("ToolCall.ID = %q, want %q", tc.ID, "call_abc")
	}
	if tc.Function.Name != "search_vector_memory" {
		t.Errorf("ToolCall.Function.Name = %q, want %q", tc.Function.Name, "search_vector_memory")
	}
	wantArgs := `{"query": "dragon"}`
	if tc.Function.Arguments != wantArgs {
		t.Errorf("ToolCall.Function.Arguments = %q, want %q", tc.Function.Arguments, wantArgs)
	}
}

func TestChatCompletionStream_MultipleInterleavedToolCalls(t *testing.T) {
	server := newSSEServer(t, []string{
		`{"choices":[{"delta":{"tool_calls":[{"index":0,"id":"call_0","type":"function","function":{"name":"search_vector_memory","arguments":"{\"q\":"}}]},"finish_reason":null}]}`,
		`{"choices":[{"delta":{"tool_calls":[{"index":1,"id":"call_1","type":"function","function":{"name":"query_entity_graph","arguments":"{\"entity\":"}}]},"finish_reason":null}]}`,
		`{"choices":[{"delta":{"tool_calls":[{"index":0,"function":{"arguments":"\"dragon\"}"}}]},"finish_reason":null}]}`,
		`{"choices":[{"delta":{"tool_calls":[{"index":1,"function":{"arguments":"\"Aria\"}"}}]},"finish_reason":null}]}`,
		`{"choices":[{"delta":{},"finish_reason":"tool_calls"}]}`,
		`[DONE]`,
	})
	defer server.Close()

	svc := newStreamTestService(server.URL)
	ch, err := svc.ChatCompletionStream(context.Background(), QwenRequest{Model: "qwen-max"})
	if err != nil {
		t.Fatalf("ChatCompletionStream: %v", err)
	}

	chunks := collectStream(t, ch)

	byName := map[string]*QwenToolCall{}
	var toolCallCount int
	for _, c := range chunks {
		if c.Type == "error" {
			t.Fatalf("unexpected error chunk: %v", c.Err)
		}
		if c.Type == "tool_call" {
			toolCallCount++
			byName[c.ToolCall.Function.Name] = c.ToolCall
		}
	}

	if toolCallCount != 2 {
		t.Fatalf("got %d tool_call chunks, want exactly 2", toolCallCount)
	}

	vecCall, ok := byName["search_vector_memory"]
	if !ok {
		t.Fatal("missing search_vector_memory tool call")
	}
	if vecCall.ID != "call_0" || vecCall.Function.Arguments != `{"q":"dragon"}` {
		t.Errorf("search_vector_memory call = %+v, want ID=call_0 Arguments={\"q\":\"dragon\"}", vecCall)
	}

	graphCall, ok := byName["query_entity_graph"]
	if !ok {
		t.Fatal("missing query_entity_graph tool call")
	}
	if graphCall.ID != "call_1" || graphCall.Function.Arguments != `{"entity":"Aria"}` {
		t.Errorf("query_entity_graph call = %+v, want ID=call_1 Arguments={\"entity\":\"Aria\"}", graphCall)
	}
}

func TestChatCompletionStream_NoCompletionSignal(t *testing.T) {
	// Stream is cut off ([DONE] arrives) without any finish_reason ever
	// being set — must surface an error and must NOT dispatch a partial
	// tool call.
	server := newSSEServer(t, []string{
		`{"choices":[{"delta":{"tool_calls":[{"index":0,"id":"call_x","type":"function","function":{"name":"search_vector_memory","arguments":"{\"query\": \"drag"}}]},"finish_reason":null}]}`,
		`[DONE]`,
	})
	defer server.Close()

	svc := newStreamTestService(server.URL)
	ch, err := svc.ChatCompletionStream(context.Background(), QwenRequest{Model: "qwen-max"})
	if err != nil {
		t.Fatalf("ChatCompletionStream: %v", err)
	}

	chunks := collectStream(t, ch)

	var sawError bool
	for _, c := range chunks {
		if c.Type == "tool_call" {
			t.Errorf("must not dispatch a partial tool call, got: %+v", c.ToolCall)
		}
		if c.Type == "error" {
			sawError = true
			if c.Err == nil {
				t.Error("error chunk must carry a non-nil Err")
			}
		}
	}
	if !sawError {
		t.Error("expected an error chunk when stream ends without a finish signal")
	}
}
