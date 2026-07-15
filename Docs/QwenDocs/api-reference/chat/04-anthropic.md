# Anthropic Messages API

> **Source:** https://docs.qwencloud.com/api-reference/chat/anthropic

POST/apps/anthropic/v1/messages Python

 Basic call

 Copy ```\nimport anthropic
import os

client = anthropic.Anthropic(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/apps/anthropic",
)

message = client.messages.create(
 model="qwen3.7-plus",
 max_tokens=1024,
 system="You are a helpful assistant",
 messages=[
 {
 "role": "user",
 "content": "Who are you?"
 }
 ],
 thinking={"type": "disabled"},
)

print(message.content[0].text)
``` 200

 application/json

 Copy ```\n{
 "id": "msg_e2898f19-fc0e-4cb3-bd9b-5b7dc4ea3bc9",
 "type": "message",
 "role": "assistant",
 "model": "qwen3.7-plus",
 "content": [
 {
 "type": "thinking",
 "thinking": "Let me analyze this question...",
 "signature": ""
 },
 {
 "type": "text",
 "text": "Hello! I am Qwen..."
 }
 ],
 "stop_reason": "end_turn",
 "stop_sequence": null,
 "usage": {
 "input_tokens": 22,
 "output_tokens": 223,
 "cache_creation_input_tokens": 0,
 "cache_read_input_tokens": 0
 }
}
``` ## [​ ](#faq) FAQ

**After configuring Claude Desktop or Claude Code, the connection test fails with `Model discovery — Gateway /v1/models returned HTTP 404`, or the request URL contains `/v1/v1/models`. How do I fix it?**
The model discovery feature of clients such as Claude Desktop and Claude Code automatically appends `/v1/models` to the configured base URL. Check the following two points:

- 
**Do not end the base URL with `/v1/`**: it should end at `/apps/anthropic` (for example, for Singapore use `https://dashscope-intl.aliyuncs.com/apps/anthropic`). If you mistakenly enter `.../apps/anthropic/v1/`, the client appends `/v1/models` and produces the duplicated path `/v1/v1/models`, which returns HTTP 404.


- 
**Add models manually to skip discovery**: the Anthropic-compatible endpoint provides only the Messages API (`/v1/messages`) and does not provide a model list endpoint (`/v1/models`), so the model discovery request returns 404 as well. Manually add models (for example, `qwen3.7-plus`) under Models in the client to skip automatic discovery.


 ### Authorizations
 [​ ](#x-api-key) x-api-keystring header required Qwen Cloud API key passed via `x-api-key` header. `Authorization: Bearer` header is also supported.

 ### Body
application/json [​ ](#model) modelstring required Model name. Supported models:


**Qwen Max**: qwen3.7-max, qwen3.7-max-2026-06-08, qwen3.7-max-2026-05-20, qwen3.6-max-preview, qwen3-max, qwen3-max-2026-01-23, qwen3-max-preview


**Qwen Plus**: qwen3.6-plus, qwen3.6-plus-2026-04-02, qwen3.5-plus, qwen3.5-plus-2026-04-20, qwen3.5-plus-2026-02-15, qwen-plus, qwen-plus-latest, qwen-plus-2025-09-11


**Qwen Flash**: qwen3.6-flash, qwen3.6-flash-2026-04-16, qwen3.5-flash, qwen3.5-flash-2026-02-23, qwen-flash, qwen-flash-2025-07-28


**Qwen Turbo**: qwen-turbo, qwen-turbo-latest


**Qwen Coder**: qwen3-coder-next, qwen3-coder-plus, qwen3-coder-plus-2025-09-23, qwen3-coder-flash


**Qwen VL**: qwen3-vl-plus, qwen3-vl-flash, qwen-vl-max, qwen-vl-plus


**Qwen Open-source**: qwen3.6-27b, qwen3.5-397b-a17b, qwen3.5-122b-a10b, qwen3.5-27b, qwen3.5-35b-a3b


**Third-party models**: deepseek-v4-pro, deepseek-v4-flash, deepseek-v3.2

 [​ ](#max-tokens) max_tokensinteger required Maximum number of tokens to generate.

 [​ ](#messages) messagesobject[] required Message array, alternating between `user` and `assistant` turns.

 Show child attributes

 [​ ](#messagesrole) messages.roleenum<string> required Message role.

 Available options:user,assistant [​ ](#messagescontent) messages.contentstring required Message content. Can be a plain text string or a structured content array.

 [​ ](#system) systemstring System prompt to set the model&#x27;s role or behavior. Passed as a top-level parameter; the `messages` array does not accept the `system` role. A string is equivalent to a single `type="text"` content block. To use context caching, pass an array of content blocks with `cache_control`.

 [​ ](#stream) streamboolean defaultfalse Enable streaming output. Default is `false`.

 [​ ](#temperature) temperaturenumber Controls diversity of generated text, range [0, 2). Higher values produce more random output. This range differs from Anthropic&#x27;s native [0.0, 1.0] — verify this parameter when migrating from Anthropic.

 [​ ](#top-p) top_pnumber Probability threshold for nucleus sampling. Both `temperature` and `top_p` control diversity — set only one at a time.

 [​ ](#top-k) top_kinteger Size of the sampling candidate set during generation.

 [​ ](#stop-sequences) stop_sequencesstring[] Text sequences that stop generation. The model stops before outputting the sequence. When hit, `stop_reason` is still `end_turn` and the matched sequence is not included in the response.

 [​ ](#thinking) thinkingobject Extended thinking configuration. When enabled, the model reasons before responding, and the response includes `thinking`-type content blocks. Not all models support thinking mode.

 Show child attributes

 [​ ](#thinkingtype) thinking.typeenum<string> `enabled` (enable thinking) or `disabled` (disable thinking).

 Available options:enabled,disabled [​ ](#thinkingbudget-tokens) thinking.budget_tokensinteger Maximum tokens for the thinking process. Does not overlap with `max_tokens`: this parameter limits thinking; `max_tokens` limits the final reply. Takes effect when `type` is `enabled`.

 [​ ](#reasoning-effort) reasoning_effortenum<string> Controls reasoning intensity. Default is `max`. Supported models: deepseek-v4-pro, deepseek-v4-flash. Values `low` or `medium` are mapped to `high`; `xhigh` is mapped to `max`.

 Available options:high,max [​ ](#tools) toolsobject[] Tool definition array for function calling.

 Show child attributes

 [​ ](#toolsname) tools.namestring required Tool name.

 [​ ](#toolsinput-schema) tools.input_schemaobject required JSON Schema definition for the tool&#x27;s input parameters.

 [​ ](#toolsdescription) tools.descriptionstring Description of the tool&#x27;s functionality.

 [​ ](#tool-choice) tool_choiceobject Tool choice strategy. `{"type": "auto"}`: model decides whether to call tools (default). `{"type": "any"}`: force calling any tool. `{"type": "none"}`: disable tool calling. `{"type": "tool", "name": "tool_name"}`: force calling a specific tool.

 Show child attributes

 [​ ](#tool-choicetype) tool_choice.typeenum<string> Strategy type.

 Available options:auto,any,none,tool [​ ](#tool-choicename) tool_choice.namestring When `type` is `tool`, specifies the name of the tool to call.

 [​ ](#output-config) output_configobject Structured output configuration. When enabled, the model returns a JSON string. Behavior varies by model:


- **Strict structured outputs**: Available for deepseek and glm series models. The model strictly follows the provided JSON Schema, guaranteeing the same field types and hierarchy.

- **Regular structured outputs**: For all other models, schema field constraints are not enforced — the API automatically falls back to a plain JSON mode (only guaranteeing that the output is a valid JSON string). In this fallback mode, the request must satisfy both of the following: (1) the `output_config` parameter is explicitly provided; (2) the `system` or `messages` content contains the keyword "JSON" (case-insensitive). If the keyword "JSON" is missing, the API throws: `&#x27;messages&#x27; must contain the word &#x27;json&#x27; in some form`.


 Show child attributes

 [​ ](#output-configformat) output_config.formatobject required The output format definition.

 Show child attributes

 [​ ](#output-configformattype) output_config.format.typeenum<string> required Fixed value: `json_schema`.

 Available options:json_schema [​ ](#output-configformatschema) output_config.format.schemaobject required JSON Schema object that follows the standard JSON Schema specification. Should include `type` (data type), `properties` (field definitions), `required` (array of required field names), and `additionalProperties` (must be set to `false`).

 ### Response
200-application/json [​ ](#id) idstring Unique message identifier.

 [​ ](#type) typeenum<string> Always `message`.

 Available options:message [​ ](#role) roleenum<string> Always `assistant`.

 Available options:assistant [​ ](#model) modelstring The model used for generation.

 [​ ](#content) contentobject[] Content array. Element types can be `text`, `thinking` (returned when thinking is enabled), or `tool_use` (tool call).

 Show child attributes

 [​ ](#contenttype) content.typeenum<string> Content block type.

 Available options:text,thinking,tool_use [​ ](#contenttext) content.textstring Model-generated text reply when `type` is `text`.

 [​ ](#contentthinking) content.thinkingstring The model&#x27;s reasoning before the final response.

 [​ ](#contentsignature) content.signaturestring Signature when `type` is `thinking`. Currently always an empty string.

 [​ ](#contentid) content.idstring Unique identifier for the tool call when `type` is `tool_use`.

 [​ ](#contentname) content.namestring Name of the tool being called when `type` is `tool_use`.

 [​ ](#contentinput) content.inputobject Tool call input parameters when `type` is `tool_use`.

 [​ ](#stop-reason) stop_reasonenum<string> Stop reason: `end_turn` (normal completion), `max_tokens` (token limit reached), `tool_use` (tool call).

 Available options:end_turn,max_tokens,tool_use [​ ](#stop-sequence) stop_sequencestring | null Always `null`.

 [​ ](#usage) usageobject Token usage statistics. In streaming, the `usage` in the `message_start` event only contains `input_tokens` and `output_tokens`; all 4 fields appear in the `message_delta` event.

 Show child attributes

 [​ ](#usageinput-tokens) usage.input_tokensinteger Number of input tokens.

 [​ ](#usageoutput-tokens) usage.output_tokensinteger Number of output tokens.

 [​ ](#usagecache-creation-input-tokens) usage.cache_creation_input_tokensinteger Number of input tokens consumed for cache creation.

 [​ ](#usagecache-read-input-tokens) usage.cache_read_input_tokensinteger Number of input tokens consumed by cache reads.
