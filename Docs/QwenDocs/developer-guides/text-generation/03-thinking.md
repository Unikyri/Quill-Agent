# Thinking

> **Source:** https://docs.qwencloud.com/developer-guides/text-generation/thinking

Solve complex tasks with step-by-step thinking

 Copy page Thinking (reasoning) models reason before answering — outputting `reasoning_content` (Chat Completions / DashScope) or `reasoning_summary_text` events (Responses API). Models support thinking in one of two modes:

- **Hybrid**: toggle thinking on or off per request with `enable_thinking`. Qwen3.7, Qwen3.6, Qwen3.5 (enabled by default); Qwen3, Qwen3-VL, Qwen3-Omni, DeepSeek-V3/V4 (disabled by default).

- **Thinking-only**: always thinks — cannot be disabled. QwQ, DeepSeek-R1, `-thinking` variants.


## [​ ](#enable-thinking) Enable thinking

- OpenAI Chat Completions 
- OpenAI Responses API 
- DashScope 

 Python Node.js curl Copy ```\nimport os
from openai import OpenAI
client = OpenAI(api_key=os.getenv("DASHSCOPE_API_KEY"), base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1")

completion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=[{"role": "user", "content": "If 3x + 7 = 22, what is x?"}],
 extra_body={"enable_thinking": True}, # ← enable thinking
 stream=True,
)
for chunk in completion:
 if not chunk.choices:
 continue
 delta = chunk.choices[0].delta
 if hasattr(delta, "reasoning_content") and delta.reasoning_content:
 print(delta.reasoning_content, end="", flush=True) # ← phase 1: thinking
 if hasattr(delta, "content") and delta.content:
 print(delta.content, end="", flush=True) # ← phase 2: answer

``` Thinking arrives as `response.reasoning_summary_text.delta` events, followed by `response.output_text.delta` for the answer.Python Node.js curl Copy ```\nimport os
from openai import OpenAI
client = OpenAI(api_key=os.getenv("DASHSCOPE_API_KEY"), base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1")

stream = client.responses.create(
 model="qwen3.7-plus",
 input="If 3x + 7 = 22, what is x?",
 extra_body={"enable_thinking": True}, # ← enable thinking
 stream=True,
)
for chunk in stream:
 if chunk.type == "response.reasoning_summary_text.delta":
 print(chunk.delta, end="", flush=True) # ← phase 1: thinking
 elif chunk.type == "response.output_text.delta":
 print(chunk.delta, end="", flush=True) # ← phase 2: answer

``` Python Java curl Copy ```\nimport dashscope
dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;
from dashscope import MultiModalConversation

responses = MultiModalConversation.call(
 model="qwen3.7-plus",
 messages=[{"role": "user", "content": [{"text": "If 3x + 7 = 22, what is x?"}]}],
 enable_thinking=True, # ← enable thinking
 stream=True,
 incremental_output=True, # ← recommended: new tokens only
)
for chunk in responses:
 msg = chunk.output.choices[0].message
 if msg.reasoning_content:
 print(msg.reasoning_content, end="", flush=True) # ← phase 1: thinking
 if msg.content and msg.content[0].get("text"):
 print(msg.content[0]["text"], end="", flush=True) # ← phase 2: answer

``` 
## [​ ](#control-thinking-depth) Control thinking depth

### [​ ](#token-budget) Token budget

Use `thinking_budget` to cap the number of thinking tokens. If the limit is reached, the model stops thinking and generates its answer immediately. All thinking-capable models from Qwen3 onward support this parameter. Chat Completions and DashScope only — not supported by the Responses API.
- OpenAI Chat Completions 
- DashScope 

 Copy ```\nextra_body={"enable_thinking": True, "thinking_budget": 500}

``` Copy ```\nenable_thinking=True,
thinking_budget=500,

``` 
### [​ ](#prompt-level-control) Prompt-level control

With `enable_thinking: true`, add `/no_think` to skip thinking for one turn. `/think` restores it. Last instruction wins. Supported by open-source Qwen3 hybrid models and `qwen-plus-2025-04-28`.
### [​ ](#preserve-thinking-in-multi-turn) Preserve thinking in multi-turn

By default, models do not read `reasoning_content` from the `messages` array in multi-turn conversations. Set `preserve_thinking` to `true` to append the `reasoning_content` in assistant messages to the next input, allowing the model to reference previous reasoning.
Supported models: `qwen3.7-max`, `qwen3.7-max-2026-06-08`, `qwen3.7-max-2026-05-20`, `qwen3.7-max-preview`, `qwen3.7-max-2026-05-17`, `qwen3.7-plus`, `qwen3.7-plus-2026-05-26`, `qwen3.6-max-preview`, `qwen3.6-plus`, `qwen3.6-plus-2026-04-02`.
- OpenAI Chat Completions 
- DashScope 

 Copy ```\nextra_body={"enable_thinking": True, "preserve_thinking": True}

``` `preserve_thinking` is not an OpenAI standard parameter. When using the Python SDK, pass it via `extra_body`. Copy ```\nenable_thinking=True,
preserve_thinking=True,

``` 
## [​ ](#function-calling-with-thinking-mode) Function calling with thinking mode

When thinking is enabled during [function calling](/developer-guides/text-generation/function-calling), the model reasons about which tools to call and how to use results before responding. The response includes `reasoning_content` before each tool call.
**Key points:**

- Pass `enable_thinking: true` alongside your `tools` array — no other config needed.

- In multi-turn tool-call flows, include the assistant&#x27;s `reasoning_content` when sending tool results back. Omitting it degrades accuracy.

- Streaming delivers thinking tokens first, then tool call deltas. See [Streaming with tool calls](/developer-guides/text-generation/streaming#streaming-with-tool-calls) for the parse pattern.

- `thinking_budget` works the same as in regular thinking mode.


 Thinking mode is most valuable for complex tool orchestration — multi-step reasoning about which tools to call, parameter selection, and result interpretation. For simple single-tool calls, the overhead may not be worth it. 
## [​ ](#notes) Notes


- **Streaming required for some models**: Non-streaming is supported by Qwen3.7 Max, Qwen3.7 Plus, Qwen3.6 Plus, Qwen3.5 Plus/Flash, Qwen3 Max, Qwen Plus/Flash/Turbo (commercial), and Qwen3.5 open-source models. Qwen3 open-source models require streaming. Streaming is always recommended to avoid timeout risks.

- **No audio output in thinking mode** (Qwen3-Omni): Text and image inputs work normally; audio output is not available when thinking is enabled.


 [Previous ](/developer-guides/text-generation/structured-output)[Batch API Process bulk requests asynchronously at 50% off Next ](/developer-guides/text-generation/batch)
