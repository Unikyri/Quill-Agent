# Streaming output

> **Source:** https://docs.qwencloud.com/developer-guides/text-generation/streaming

Receive text output token by token as it is generated.

 Copy page ## [​ ](#enable-streaming) Enable streaming

- OpenAI Chat Completions 
- OpenAI Responses API 
- DashScope 

 No usage by default. Set `stream_options` to get token counts in the **last chunk only**.Python Node.js curl Copy ```\nimport os
from openai import OpenAI
client = OpenAI(api_key=os.getenv("DASHSCOPE_API_KEY"), base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1")

completion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=[{"role": "user", "content": "Hi"}],
 stream=True, # ← enable streaming
 stream_options={"include_usage": True}, # ← usage in last chunk only
)
for chunk in completion:
 if chunk.choices:
 print(chunk.choices[0].delta.content or "", end="", flush=True)
 elif chunk.usage: # ← last chunk: usage only
 print(f"\nTokens: {chunk.usage.total_tokens}")

``` Python Node.js curl Copy ```\nimport os
from openai import OpenAI
client = OpenAI(api_key=os.getenv("DASHSCOPE_API_KEY"), base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1")

stream = client.responses.create(
 model="qwen3.7-plus",
 input="Hi",
 stream=True, # ← enable streaming
)
for event in stream:
 if event.type == "response.output_text.delta":
 print(event.delta, end="", flush=True)
 elif event.type == "response.completed":
 print(f"\nTokens: {event.response.usage.total_tokens}")

``` Unlike OpenAI compatible (last chunk only), DashScope returns real-time token usage in **every chunk** — useful for monitoring costs or stopping early. Qwen3.5 and Qwen3.6 models use the `multimodal-generation` endpoint shown below. Earlier models such as qwen-plus and qwen3-max use the `text-generation` endpoint instead. See [DashScope API reference](/api-reference/chat/dashscope) for details. Python Java curl Copy ```\nfrom http import HTTPStatus
import dashscope
dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;
from dashscope import MultiModalConversation

responses = MultiModalConversation.call(
 model="qwen3.7-plus",
 messages=[{"role": "user", "content": [{"text": "Hi"}]}],
 stream=True, # ← enable streaming
 incremental_output=True, # ← recommended: new tokens only per chunk
)
for resp in responses:
 if resp.status_code == HTTPStatus.OK:
 content = resp.output.choices[0].message.content
 if content:
 print(content[0]["text"], end="", flush=True)

``` 
## [​ ](#event-format) Event format

Each SSE event is a `data:` line containing a JSON chunk. The final `data: [DONE]` signals the end of the stream.
Copy ```\ndata: {"choices":[{"delta":{"content":"I am"},...,"finish_reason":null}],...}
data: {"choices":[{"delta":{"content":" Qwen"},...,"finish_reason":null}],...}
data: {"choices":[{"delta":{"content":""},...,"finish_reason":"stop"}],...}
data: [DONE]

``` 
## [​ ](#streaming-with-thinking-mode) Streaming with thinking mode

Two-phase streaming: thinking first, then the answer.
- OpenAI Chat Completions 
- OpenAI Responses API 
- DashScope 

 Copy ```\nfor chunk in completion:
 delta = chunk.choices[0].delta
 if hasattr(delta, "reasoning_content") and delta.reasoning_content:
 print(delta.reasoning_content, end="", flush=True) # ← phase 1: thinking
 if hasattr(delta, "content") and delta.content:
 print(delta.content, end="", flush=True) # ← phase 2: answer

``` Copy ```\nfor event in stream:
 if event.type == "response.reasoning_summary_text.delta":
 print(event.delta, end="", flush=True) # ← phase 1: thinking
 elif event.type == "response.output_text.delta":
 print(event.delta, end="", flush=True) # ← phase 2: answer

``` Copy ```\nfor chunk in completion:
 msg = chunk.output.choices[0].message
 if msg.reasoning_content: # ← phase 1: thinking
 print(msg.reasoning_content, end="", flush=True)
 if msg.content: # ← phase 2: answer
 print(msg.content, end="", flush=True)

``` 
→ Full config: [Reasoning](/developer-guides/text-generation/thinking) | Qwen3-Omni: [Audio and video](/developer-guides/speech/multimodal-speech)
## [​ ](#streaming-with-tool-calls) Streaming with tool calls

When streaming [function calling](/developer-guides/text-generation/function-calling) responses, tool call arguments arrive as incremental deltas that must be concatenated before JSON parsing.
- Chat Completions 
- Responses API 

 Each chunk&#x27;s `delta` may contain `tool_calls[i].function.arguments` — a partial JSON string. Accumulate all fragments per tool call index, then `JSON.parse()` the complete string.Copy ```\ntool_args = {}
for chunk in completion:
 delta = chunk.choices[0].delta
 if delta.tool_calls:
 for tc in delta.tool_calls:
 tool_args.setdefault(tc.index, "")
 tool_args[tc.index] += tc.function.arguments or ""
# After stream ends:
for idx, args_str in tool_args.items():
 parsed = json.loads(args_str)

``` Tool calls emit `response.function_call_arguments.delta` events. Concatenate deltas until `response.function_call_arguments.done`.Copy ```\nargs_buffer = ""
for event in stream:
 if event.type == "response.function_call_arguments.delta":
 args_buffer += event.delta
 elif event.type == "response.function_call_arguments.done":
 parsed = json.loads(args_buffer)
 args_buffer = ""

``` 
 With [thinking mode](/developer-guides/text-generation/thinking) enabled, the stream delivers three phases: thinking tokens, then tool call deltas, then (after you send tool results) the final answer. 
## [​ ](#notes) Notes


- **Nginx proxy**: set `proxy_buffering off` or SSE events will buffer

- **High concurrency**: size connection pool, monitor file descriptors

- **Web frontend**: use `ReadableStream` + `TextDecoderStream`

- **Quality**: streaming does not affect response quality

- **Streaming-only models**: QwQ and QVQ only support streaming output. Non-streaming calls to these models will fail or return empty content.


## [​ ](#faq) FAQ

**What is the difference between streaming and non-streaming calls?**

- **Timeout**: Non-streaming calls have a fixed maximum timeout of 300 seconds. If the model does not finish generating within 300 seconds, the request will time out and fail.

- **Output structure**: Non-streaming calls return the complete response as a single JSON object. Streaming calls return data chunks incrementally via the SSE protocol, where each chunk contains partial generated content that the client must concatenate.

- **Feature compatibility**: Both support JSON Mode, Function Call, and other features with no functional differences.


We recommend using streaming output to avoid timeout issues and provide a better user experience.
**Does streaming output support JSON Mode (structured output)?**
Yes. Set `stream` to `true` and `response_format` to `{"type": "json_object"}` in the same request. The model will return JSON content incrementally via streaming. The concatenated complete output will be valid JSON. [Previous ](/developer-guides/text-generation/context-cache)[Async task management Two async patterns on Qwen Cloud: task-based for media generation, batch for high-volume text processing Next ](/developer-guides/run-and-scale/async-task-management)
