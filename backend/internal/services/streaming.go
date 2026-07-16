package services

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

// StreamChunk is one incremental event emitted while consuming a Qwen
// streaming chat completion. Type is one of "text", "tool_call", "done", or
// "error"; only the field(s) matching Type are populated.
type StreamChunk struct {
	Type     string
	Text     string
	ToolCall *QwenToolCall
	Finish   string
	Err      error
}

// accToolCall buffers a tool call's fragmented deltas (name arrives once,
// arguments arrive as concatenated string fragments) until the stream signals
// completion, keyed by the delta's index.
type accToolCall struct {
	ID   string
	Name string
	Args strings.Builder
}

// streamResponseFrame mirrors one SSE "data:" frame of a Qwen streaming chat
// completion response.
type streamResponseFrame struct {
	Choices []streamChoice `json:"choices"`
	Output  struct {
		Text         string         `json:"text"`
		Choices      []streamChoice `json:"choices"`
		FinishReason string         `json:"finish_reason"`
	} `json:"output"`
	Usage streamUsage `json:"usage"`
}

type streamChoice struct {
	Delta struct {
		Content   string           `json:"content"`
		ToolCalls []streamToolCall `json:"tool_calls"`
	} `json:"delta"`
	Message struct {
		Content   string           `json:"content"`
		ToolCalls []streamToolCall `json:"tool_calls"`
	} `json:"message"`
	FinishReason string `json:"finish_reason"`
}

type streamToolCall struct {
	Index    int    `json:"index"`
	ID       string `json:"id"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

// streamUsage is provider-neutral enough for both the compatible and native
// DashScope SSE envelopes. DashScope includes cache counts in the final frame.
type streamUsage struct {
	InputTokens         int `json:"input_tokens"`
	OutputTokens        int `json:"output_tokens"`
	TotalTokens         int `json:"total_tokens"`
	PromptTokensDetails struct {
		CachedTokens             int `json:"cached_tokens"`
		CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	} `json:"prompt_tokens_details"`
}

// ChatCompletionStream sends payload with stream:true and returns a channel
// of StreamChunk values read from the SSE response body. The channel is
// closed when the stream ends (successfully or on error).
//
// Verified against a live DashScope streaming + tool-call response (PR3
// spike): the final chunk before [DONE] carries finish_reason:"tool_calls",
// tool-call deltas are keyed by index, function.name arrives once, and
// function.arguments arrives as concatenated string fragments.
func (s *QwenService) ChatCompletionStream(ctx context.Context, payload QwenRequest) (<-chan StreamChunk, error) {
	payload = s.normalizeRequestMessages(payload)
	payload.Stream = true
	resp, release, err := s.sendQwenRequest(ctx, s.tierForModel(payload.Model), payload.Model, http.MethodPost, "/chat/completions", payload, s.estimateChatTokens(payload), true)
	if err != nil {
		return nil, err
	}

	ch := make(chan StreamChunk)
	go func() { release(readStream(resp.Body, ch)) }()
	return ch, nil
}

// readStream parses SSE frames off body, accumulates tool-call deltas by
// index, and emits StreamChunk values on ch. It closes ch and body when done.
func readStream(body io.ReadCloser, ch chan<- StreamChunk) (success bool) {
	return readStreamWithUsage(body, ch, nil)
}

// readStreamWithUsage is the shared SSE parser. The optional callback receives
// the last provider-reported usage frame after a successful stream, avoiding
// double-counting when a provider repeats cumulative usage on every frame.
func readStreamWithUsage(body io.ReadCloser, ch chan<- StreamChunk, onUsage func(streamUsage)) (success bool) {
	defer close(ch)
	defer body.Close()

	acc := map[int]*accToolCall{}
	finished := false
	var lastUsage streamUsage
	hasUsage := false

	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		// SSE streams may include event/id/retry fields and comments between
		// data frames. They are transport metadata, not JSON payloads. Native
		// DashScope responses emit these fields while the compatible endpoint
		// commonly omits them, so ignore them consistently for both transports.
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "" {
			continue
		}
		if data == "[DONE]" {
			break
		}

		var frame streamResponseFrame
		if err := json.Unmarshal([]byte(data), &frame); err != nil {
			ch <- StreamChunk{Type: "error", Err: fmt.Errorf("unmarshal stream chunk: %w", err)}
			return false
		}
		if frame.Usage.InputTokens != 0 || frame.Usage.OutputTokens != 0 || frame.Usage.TotalTokens != 0 || frame.Usage.PromptTokensDetails.CachedTokens != 0 || frame.Usage.PromptTokensDetails.CacheCreationInputTokens != 0 {
			lastUsage = frame.Usage
			hasUsage = true
		}
		choices := frame.Choices
		if len(choices) == 0 {
			choices = frame.Output.Choices
		}
		if len(choices) == 0 && frame.Output.Text != "" {
			var choice streamChoice
			choice.Delta.Content = frame.Output.Text
			choice.FinishReason = frame.Output.FinishReason
			choices = []streamChoice{choice}
		}
		if len(choices) == 0 {
			continue
		}
		choice := choices[0]
		if choice.Delta.Content == "" {
			choice.Delta.Content = choice.Message.Content
		}
		if len(choice.Delta.ToolCalls) == 0 {
			choice.Delta.ToolCalls = choice.Message.ToolCalls
		}
		if choice.FinishReason == "" {
			choice.FinishReason = frame.Output.FinishReason
		}

		if choice.Delta.Content != "" {
			ch <- StreamChunk{Type: "text", Text: choice.Delta.Content}
		}

		for _, tc := range choice.Delta.ToolCalls {
			entry, ok := acc[tc.Index]
			if !ok {
				entry = &accToolCall{}
				acc[tc.Index] = entry
			}
			if tc.ID != "" {
				entry.ID = tc.ID
			}
			if tc.Function.Name != "" {
				entry.Name = tc.Function.Name
			}
			entry.Args.WriteString(tc.Function.Arguments)
		}

		switch choice.FinishReason {
		case "tool_calls":
			finished = true
			flushToolCalls(acc, ch)
		case "stop", "length", "content_filter":
			finished = true
			ch <- StreamChunk{Type: "done", Finish: choice.FinishReason}
		}
	}

	if err := scanner.Err(); err != nil {
		ch <- StreamChunk{Type: "error", Err: fmt.Errorf("read stream: %w", err)}
		return false
	}

	if !finished {
		ch <- StreamChunk{Type: "error", Err: fmt.Errorf("stream ended without a finish_reason completion signal")}
		return false
	}
	if hasUsage && onUsage != nil {
		onUsage(lastUsage)
	}
	return true
}

// flushToolCalls emits one tool_call StreamChunk per accumulated entry, in
// index order, so multi-tool-call responses dispatch deterministically.
func flushToolCalls(acc map[int]*accToolCall, ch chan<- StreamChunk) {
	indices := make([]int, 0, len(acc))
	for idx := range acc {
		indices = append(indices, idx)
	}
	sort.Ints(indices)

	for _, idx := range indices {
		entry := acc[idx]
		ch <- StreamChunk{
			Type: "tool_call",
			ToolCall: &QwenToolCall{
				ID:   entry.ID,
				Type: "function",
				Function: QwenToolCallFunction{
					Name:      entry.Name,
					Arguments: entry.Args.String(),
				},
			},
		}
	}
}
