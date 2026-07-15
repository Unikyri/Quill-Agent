# OpenAI compatibility

> **Source:** https://docs.qwencloud.com/api-reference/toolkitframework/openai-compatible/overview

Migrate from OpenAI by changing three parameters: base_url, api_key, and model.

 Copy page Qwen Cloud provides OpenAI-compatible APIs. If you have existing code that uses the OpenAI SDK or REST API, you can switch to Qwen models by changing three parameters: `base_url`, `api_key`, and `model`.
## [​ ](#quick-migration) Quick migration

- Python 
- Node.js 
- curl 

 Copy ```\nimport os
from openai import OpenAI

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)

completion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=[{"role": "user", "content": "Hello!"}],
)
print(completion.choices[0].message.content)

``` Copy ```\nimport OpenAI from "openai";

const openai = new OpenAI({
 apiKey: process.env.DASHSCOPE_API_KEY,
 baseURL: "https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
});

async function main() {
 const completion = await openai.chat.completions.create({
 model: "qwen3.7-plus",
 messages: [{ role: "user", content: "Hello!" }],
 });
 console.log(completion.choices[0].message.content);
}

main();

``` Copy ```\ncurl https://dashscope-intl.aliyuncs.com/compatible-mode/v1/chat/completions \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY" \
 -H "Content-Type: application/json" \
 -d &#x27;{
 "model": "qwen3.7-plus",
 "messages": [{"role": "user", "content": "Hello!"}]
 }&#x27;

``` 
 Before you begin, [get an API key](/api-reference/preparation/api-key) and [set it as an environment variable](/api-reference/preparation/export-api-key-env). If you use the OpenAI SDK, [install it](/api-reference/preparation/install-sdk). 
## [​ ](#supported-apis) Supported APIs

APIBase URL (for SDK)Description[Chat Completions](#chat-completions)`https://dashscope-intl.aliyuncs.com/compatible-mode/v1`Text generation, vision, function calling[Responses](#responses-api)`https://dashscope-intl.aliyuncs.com/api/v2/apps/protocols/compatible-mode/v1`Built-in tools, simplified multi-turn[Embedding](#embedding)`https://dashscope-intl.aliyuncs.com/compatible-mode/v1`Text embeddings[File](#file-api)`https://dashscope-intl.aliyuncs.com/compatible-mode/v1`File upload and management[Batch](#batch-api)`https://dashscope-intl.aliyuncs.com/compatible-mode/v1`Asynchronous bulk processing at 50% cost[Conversations](#conversations-api)`https://dashscope-intl.aliyuncs.com/api/v2/apps/protocols/compatible-mode/v1`Auto-managed multi-turn context 
 The Responses API and Conversations API use a **different** `base_url` from the other four APIs. Make sure you set the correct `base_url` for the API you are calling. 
## [​ ](#chat-completions) Chat Completions

The Chat Completions API (`/v1/chat/completions`) is largely compatible with [OpenAI&#x27;s Chat API](https://platform.openai.com/docs/api-reference/chat/create). The key differences are listed below.
### [​ ](#qwen-specific-parameters) Qwen-specific parameters

These parameters are not part of the OpenAI standard. In the OpenAI Python SDK, pass them via `extra_body`.
ParameterTypeDescription`enable_thinking`BooleanEnable deep reasoning mode. Some models require streaming. See [Thinking](/developer-guides/text-generation/thinking#notes).`thinking_budget`IntegerMax tokens for the thinking process. Same streaming requirements as `enable_thinking`.`enable_search`BooleanEnable web search. Replaces OpenAI&#x27;s `web_search_options`.`search_options`ObjectConfigure search behavior (strategy, forced search, etc.).`top_k`IntegerSampling candidate set size. Range: (0, 100].`vl_high_resolution_images`BooleanEnable high-resolution mode for vision models.`enable_code_interpreter`BooleanEnable code interpreter. Streaming required (not required for Responses API). 
### [​ ](#behavioral-differences) Behavioral differences


- **`response_format`** supports `json_object` only (no `json_schema`).

- **`tool_choice`** supports `auto`, `none`, and specific function object (`{"type": "function", "function": {"name": "..."}}`). The `required` value is not supported.

- **`tools`** supports `function` type only.

- **`parallel_tool_calls`** defaults to `false` (OpenAI defaults to `true`).

- **`n`** supports 1-4 and is limited to specific models (qwen-plus, qwen-plus-character).

- **`web_search_options`** is not supported. Use `extra_body.enable_search` and `extra_body.search_options` instead.


### [​ ](#unsupported-parameters) Unsupported parameters

The following parameters are silently ignored: `frequency_penalty`, `logit_bias`, `max_completion_tokens`, `metadata`, `prediction`, `prompt_cache_key`, `reasoning_effort`, `service_tier`, `store`, `verbosity`.
For the full API reference and code examples, see [Chat Completions](/api-reference/chat/openai-chat).
## [​ ](#responses-api) Responses API

The Responses API uses a **different base_url**: `https://dashscope-intl.aliyuncs.com/api/v2/apps/protocols/compatible-mode/v1`.
Compared to Chat Completions, the Responses API offers:

- **Built-in tools**: `web_search`, `code_interpreter`, `web_extractor`, and `image_search` -- no external tool setup required.

- **Simplified multi-turn**: Pass `previous_response_id` instead of building a full message history.

- **Conversation integration**: Pair with the [Conversations API](#conversations-api) for automatic context management.

- **Session cache**: Automatically caches context across turns to reduce latency and cost. Enable with the `x-dashscope-session-cache: enable` header. See [Session cache](/developer-guides/text-generation/context-cache#session-cache).


### [​ ](#migrate-from-chat-completions) Migrate from Chat Completions

To switch from Chat Completions to the Responses API:

- Change `base_url` to `https://dashscope-intl.aliyuncs.com/api/v2/apps/protocols/compatible-mode/v1` and the endpoint path from `/v1/chat/completions` to `/v1/responses`.

- Read the response with `output_text` instead of `choices[0].message.content`.

- For multi-turn conversations, pass `previous_response_id` instead of manually appending messages.


For the full reference and code examples, see [Responses](/api-reference/chat/openai-responses).
## [​ ](#embedding) Embedding

The Embedding API (`/v1/embeddings`) is compatible with [OpenAI&#x27;s Embedding API](https://platform.openai.com/docs/api-reference/embeddings). Key differences:

- **`encoding_format`**: Only `float` is supported (default and only option).

- **`user`**: Not supported.

- **`dimensions`**: Available values depend on the model. For example, `text-embedding-v4` supports 2,048, 1,536, 1,024 (default), 768, 512, 256, 128, and 64.


For supported models and code examples, see [Embedding](/api-reference/text-embedding/openai-embedding).
## [​ ](#file-api) File API

The File API (`/v1/files`) is compatible with [OpenAI&#x27;s Files API](https://platform.openai.com/docs/api-reference/files), with these differences:

- **`purpose`** must be `file-extract` (for document analysis with Qwen-Long/Qwen-Doc) or `batch` (for batch processing). OpenAI values like `fine-tune` and `assistants` are not supported.

- **File content retrieval** (`GET /v1/files/{file_id}/content`) is not supported.

- **List filtering**: The `purpose` and `order` parameters on `GET /v1/files` are not supported.

- **Storage limits**: 10,000 files, 100 GB total. Files never expire.


For the full reference, see [File](/api-reference/platform-api/file).
## [​ ](#batch-api) Batch API

The Batch API (`/v1/batches`) is compatible with [OpenAI&#x27;s Batch API](https://platform.openai.com/docs/api-reference/batch), with these differences:

- **50% cost discount** compared to real-time pricing.

- **`completion_window`**: Supports 24h to 336h (14 days). Accepts "h" (hours) and "d" (days) units with integer values. OpenAI is fixed at 24h.

- **Extra metadata**: `metadata.ds_name` (task name) and `metadata.ds_description` (task description).

- **Extra list filters**: `ds_name`, `input_file_ids`, `status`, `create_after`, `create_before`.

- **Input file limits**: Up to 50,000 requests per file, 500 MB total, 6 MB per line. All requests in a file must use the same model.


For the full workflow guide, see [Batch API](/api-reference/platform-api/batch/create-batch).
## [​ ](#conversations-api) Conversations API

The Conversations API is a Qwen-specific feature with no direct OpenAI equivalent. It automatically manages multi-turn context across devices and sessions. It uses the same `base_url` as the Responses API: `https://dashscope-intl.aliyuncs.com/api/v2/apps/protocols/compatible-mode/v1`.
Use it with the Responses API to inject historical context without manual message synchronization.
For the full reference, see [Conversations](/api-reference/platform-api/conversations). [Previous ](/api-reference/platform-api/batch/cancel-batch)[Generate a temporary API key Short-lived access tokens Next ](/api-reference/more/generate-a-temporary-api-key)
