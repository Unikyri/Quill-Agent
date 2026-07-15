# OpenAI chat

> **Source:** https://docs.qwencloud.com/api-reference/chat/openai-chat

POST/compatible-mode/v1/chat/completions Python

 Text input

 Copy ```\nimport os
from openai import OpenAI

client = OpenAI(
 # If the environment variable is not set, replace the following line with: api_key="sk-xxx"
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)

completion = client.chat.completions.create(
 model="qwen3.7-plus"
 messages=[
 {"role": "system", "content": "You are a helpful assistant."},
 {"role": "user", "content": "Who are you?"},
 ],
 # extra_body={"enable_thinking": False},
)
print(completion.model_dump_json())
``` 200

 application/json

 Copy ```\n{
 "choices": [
 {
 "message": {
 "role": "assistant",
 "content": "I am a large-scale language model developed by Alibaba Cloud. My name is Qwen."
 },
 "finish_reason": "stop",
 "index": 0,
 "logprobs": null
 }
 ],
 "object": "chat.completion",
 "usage": {
 "prompt_tokens": 3019,
 "completion_tokens": 104,
 "total_tokens": 3123,
 "prompt_tokens_details": {
 "cached_tokens": 2048
 }
 },
 "created": 1735120033,
 "system_fingerprint": null,
 "model": "qwen3.7-plus",
 "id": "chatcmpl-6ada9ed2-7f33-9de2-8bb0-78bd4035025a"
}
``` ### Authorizations
 [​ ](#authorization) Authorizationstring header required DashScope API key.

 ### Body
application/json [​ ](#model) modelstring required The name of the model to use. Supported models include Qwen large language models (commercial and open source), Qwen-VL, Qwen-Coder, Qwen-Omni, and Qwen-Math. For specific model names and billing information, see [Text Generation - Qwen](#).

 [​ ](#messages) messages(System message · object | User message · object | Assistant message · object | Tool message · object)[] required The conversation history for the model, listed in chronological order.

 A system message that defines the role, tone, task objectives, or constraints for the model. Place it at the beginning of the messages array. Do not set a System Message for the QwQ model. A System Message has no effect on the QVQ model.

 - System message 
- User message 
- Assistant message 
- Tool message 

 Show child attributes

 [​ ](#messagesrole) messages.roleenum<string> required The role for a system message. Fixed as `system`.

 Available options:system [​ ](#messagescontent) messages.contentstring required A system instruction that defines the model&#x27;s role, behavior, response style, and task constraints.

 [​ ](#stream) streamboolean defaultfalse Enables streaming output mode. When `true`, the model generates and sends output incrementally. A data block (chunk) is returned as soon as part of the content is generated. You can read these chunks in real time to assemble the full reply. Set this to `true` to improve the reading experience and reduce the risk of timeouts.

 [​ ](#stream-options) stream_optionsobject Configuration options for streaming output. This parameter is effective only when `stream` is set to `true`.

 Show child attributes

 [​ ](#stream-optionsinclude-usage) stream_options.include_usageboolean defaultfalse Specifies whether to include token consumption information in the last data block of the response.

 [​ ](#modalities) modalitiesstring[] default["text"] Specifies the modalities of the output data. This parameter applies only to Qwen-Omni models. Valid values: `["text","audio"]` or `["text"]`.

 [​ ](#audio) audioobject The voice and format of the output audio. This parameter applies only to Qwen-Omni models, and you must set the `modalities` parameter to `["text","audio"]`.

 Show child attributes

 [​ ](#audiovoice) audio.voicestring required The voice used for the output audio. For more information, see [Voice list](#).

 [​ ](#audioformat) audio.formatstring required The format of the output audio. Only `wav` is supported.

 [​ ](#temperature) temperaturenumber The sampling temperature controls the diversity of the generated text. Higher values increase diversity, while lower values make the output more deterministic. The value must be greater than or equal to 0 and less than 2. Both `temperature` and `top_p` control the diversity of the generated text. Set only one of them. Do not modify the default temperature value for QVQ models.

 [​ ](#top-p) top_pnumber The probability threshold for nucleus sampling. A higher `top_p` value produces more diverse text. A lower `top_p` value produces more deterministic text. Value range: (0, 1.0]. Both `temperature` and `top_p` control the diversity of the generated text. Set only one of them. Do not modify the default `top_p` value for QVQ models.

 [​ ](#top-k) top_kinteger Specifies the number of candidate tokens to use for sampling during generation. A larger value produces more random output, whereas a smaller value produces more deterministic output. If set to `null` or a value greater than 100, the `top_k` strategy is disabled and only `top_p` takes effect. The value must be an integer greater than or equal to 0.


**Default `top_k` values:**


- QVQ series, qwen-vl-plus-2025-07-10, and qwen-vl-plus-2025-08-15: 10

- QwQ series: 40

- Other qwen-vl-plus series, models before qwen-vl-max-2025-08-13, qwen2.5-omni-7b: 1

- Qwen3-Omni-Flash series: 50

- All other models: 20


You must not change the default `top_k` value for QVQ models.


**This is not a standard OpenAI parameter.** When using the Python SDK, include it in `extra_body`: `extra_body={"top_k": xxx}`.

 [​ ](#presence-penalty) presence_penaltynumber Controls how strongly the model avoids repeating content. Valid values: -2.0 to 2.0. Positive values reduce repetition. Negative values increase it. For creative writing or brainstorming, increase this value. For technical documents or formal text, decrease this value.


**Default `presence_penalty` values:**


- Qwen3.5 (non-thinking mode), qwen3-max-preview (thinking mode), Qwen3 (non-thinking mode), Qwen3-Instruct series, qwen3-0.6b/1.7b/4b (thinking mode), QVQ series, qwen-max, qwen-max-latest, qwen2.5-vl series, qwen-vl-max series, qwen-vl-plus, Qwen3-VL (non-thinking): 1.5

- qwen-vl-plus-latest, qwen-vl-plus-2025-08-15: 1.2

- qwen-vl-plus-2025-01-25: 1.0

- qwen3-8b/14b/32b/30b-a3b/235b-a22b (thinking mode), qwen3.6-plus/qwen3.6-plus-2026-04-02, qwen3.5-plus/qwen3.5-plus-latest/2025-04-28 (thinking mode), qwen-turbo/qwen-turbo/2025-04-28 (thinking mode): 0.5

- All other models: 0.0


**How it works:** When the parameter value is positive, the model penalizes tokens that already appear in the generated text. The penalty does not depend on how many times a token appears. This reduces the likelihood of those tokens reappearing, which decreases repetition and increases lexical diversity.


When using the qwen-vl-plus-2025-01-25 model for text extraction, set `presence_penalty` to 1.5.


Do not modify the default `presence_penalty` value for QVQ models.

 [​ ](#response-format) response_formatobject The format of the response content. Defaults to `{"type": "text"}`.


Valid values:


- `{"type": "text"}`: Returns a plain text response.

- `{"type": "json_object"}`: Returns a JSON string that conforms to standard JSON syntax.

- `{"type": "json_schema", "json_schema": {...}}`: Returns a JSON string that conforms to a custom schema.


If you specify `{"type": "json_object"}`, explicitly instruct the model to output JSON in the prompt, such as by adding "Please output in JSON format." Otherwise, an error occurs.


For supported models, see Structured output.

 Show child attributes

 [​ ](#response-formattype) response_format.typeenum<string> required The format type. `text` returns plain text. `json_object` returns a JSON string that conforms to standard JSON syntax. `json_schema` returns a JSON string that conforms to a custom schema.

 Available options:text,json_object,json_schema [​ ](#response-formatjson-schema) response_format.json_schemaobject Defines the configuration for structured output. Required when `type` is `json_schema`.

 Show child attributes

 [​ ](#response-formatjson-schemaname) response_format.json_schema.namestring required A unique name for the schema. Can contain only letters, numbers, underscores, and hyphens. Maximum 64 characters.

 [​ ](#response-formatjson-schemadescription) response_format.json_schema.descriptionstring A description of the schema&#x27;s purpose. This helps the model understand the semantic context of the output.

 [​ ](#response-formatjson-schemaschema) response_format.json_schema.schemaobject An object that conforms to the JSON Schema standard defining the output data structure.

 [​ ](#response-formatjson-schemastrict) response_format.json_schema.strictboolean defaultfalse Specifies whether the model must strictly follow all schema constraints. `true` (recommended) ensures full compliance. `false` may result in output that does not conform to the specification.

 [​ ](#max-tokens) max_tokensinteger deprecated **Deprecated.** Use `max_completion_tokens` for new integrations.


The maximum number of tokens in the response. Generation stops when this limit is reached, and the `finish_reason` field is set to `length`. The default and maximum values correspond to the model&#x27;s maximum output length. `max_tokens` does not limit the length of the chain-of-thought.

 [​ ](#max-completion-tokens) max_completion_tokensinteger The maximum number of tokens in this response, including chain-of-thought tokens. Generation stops when this limit is reached, and `finish_reason` is set to `length`. The default and maximum values correspond to the model&#x27;s maximum output length.


**Difference from `max_tokens`:** `max_completion_tokens` limits the total length of both the chain-of-thought and the final response, while `max_tokens` does not limit the chain-of-thought length. For reasoning models, `max_completion_tokens` is recommended.


**Supported models:**


- Qwen-Max: Qwen3.7-Max and later

- Qwen-Plus: Qwen3.5-Plus and later

- Qwen-Flash: Qwen3.5-Flash and later

- DeepSeek: deepseek-v3, deepseek-r1, deepseek-r1-0528, deepseek-v3.1, deepseek-v3.2, deepseek-v3.2-exp, deepseek-v4-pro, deepseek-v4-flash and later


The actual number of output tokens may differ from the configured value by up to 10 tokens.

 [​ ](#vl-high-resolution-images) vl_high_resolution_imagesboolean defaultfalse Increases the maximum pixel limit for input images to the pixel value corresponding to 16384 tokens. When `true`, a fixed-resolution strategy is used and the `max_pixels` setting is ignored. If an image exceeds the pixel limit, its total pixel count is downscaled to meet the limit.


**Pixel limits when `vl_high_resolution_images` is `true`:**


- Qwen3.5 series, Qwen3-VL series, qwen-vl-max, qwen-vl-max-latest, qwen-vl-max-0813, qwen-vl-plus, qwen-vl-plus-latest, qwen-vl-plus-0815: 16,777,216 (each token corresponds to 32×32 pixels, i.e., 16,384×32×32)

- QVQ series and other Qwen2.5-VL series models: 12,845,056 (each token corresponds to 28×28 pixels, i.e., 16,384×28×28)


If `vl_high_resolution_images` is `false`, the actual pixel limit is determined by `max_pixels`.


**This is not a standard OpenAI parameter.** When using the Python SDK, include it in `extra_body`: `extra_body={"vl_high_resolution_images": xxx}`.

 [​ ](#n) ninteger default1 The number of responses to generate. Must be an integer in the range of 1-4. This is useful for scenarios that require multiple candidate responses, such as creative writing or ad copy. Supported only by Qwen3 (non-thinking mode) models. If you pass the `tools` parameter, set `n` to 1. Increasing `n` increases output token consumption but does not affect input token consumption.

 [​ ](#enable-thinking) enable_thinkingboolean Enables the thinking mode for hybrid thinking models. This mode is available for the Qwen3.5, Qwen3, Qwen3-Omni-Flash, and Qwen3-VL models. When enabled, the thinking content is returned in the `reasoning_content` field.


Default values differ by model. For supported models and their default `enable_thinking` values, see the Model List.


**This is not a standard OpenAI parameter.** When using the Python SDK, place it in `extra_body`: `extra_body={"enable_thinking": xxx}`.

 [​ ](#thinking-budget) thinking_budgetinteger The maximum number of tokens for the thinking process. Applies to Qwen3.5, Qwen3-VL, and the commercial and open source versions of Qwen3 models. The default value is the model&#x27;s maximum chain-of-thought length. For more information, see the Model List.


**This is not a standard OpenAI parameter.** When using the Python SDK, place it in `extra_body`: `extra_body={"thinking_budget": xxx}`.

 [​ ](#tool-stream) tool_streamboolean defaultfalse Only takes effect when `stream=true`. Controls streaming behavior for complex tool arguments. Ordinary tool arguments (all parameter types are string) are streamed as long as `stream=true` is enabled; `tool_stream` has no effect on them. Complex tools are those whose definitions include parameters of type array or object.


**Supported Qwen models:**


- qwen-max series: text modality of the qwen3.7-max series

- qwen-plus series: text modality of the qwen3.7-plus and qwen3.6-plus series, and all modalities of the qwen3.5-plus series

- qwen-flash series: all modalities of the qwen3.6-flash and qwen3.5-flash series


**Qwen usage:**


- `tool_stream=false`: Complex tool arguments are returned all at once (default). More accurate for complex formats.

- `tool_stream=true`: Complex tool arguments are streamed. Avoids timeout risk for complex formats.


**Supported GLM models:** glm-5.1.


**GLM usage:**


- `tool_stream=false`: Tool arguments are returned all at once (default). More accurate for complex formats.

- `tool_stream=true`: Tool arguments are streamed. Avoids timeout risk for complex formats.


**This is not a standard OpenAI parameter.** When using the Python SDK, include it in `extra_body`: `extra_body={"tool_stream": true}`.

 [​ ](#enable-code-interpreter) enable_code_interpreterboolean defaultfalse Specifies whether to enable the code interpreter feature.


**This is not a standard OpenAI parameter.** When using the Python SDK, include it in `extra_body`: `extra_body={"enable_code_interpreter": xxx}`.

 [​ ](#seed) seedinteger The random number seed. Ensures reproducible results. If you use the same seed value and other parameters remain unchanged, the model returns the same result whenever possible. Valid values: [0, 2^31-1].

 [​ ](#logprobs) logprobsboolean defaultfalse Specifies whether to return the log probabilities of the output tokens. Content generated during the thinking phase (`reasoning_content`) does not include log probabilities.


**Supported models:**


- Qwen-plus series snapshots (excluding the stable model)

- Qwen-turbo series snapshots (excluding the stable model)

- Qwen3-vl-plus models (including the stable model)

- Qwen3-vl-flash models (including the stable model)

- Qwen3 open source models


 [​ ](#top-logprobs) top_logprobsinteger default0 The number of most likely candidate tokens to return at each generation step. Valid values: 0 to 5. This parameter applies only if `logprobs` is set to `true`.

 [​ ](#stop) stopstring Stop words. If a string or token specified in `stop` appears in the generated text, generation stops immediately. If `stop` is an array, do not use a `token_id` and a string as elements simultaneously.

 [​ ](#tools) toolsobject[] An array of one or more tool objects that the model can call in function calling. When `tools` is set and the model determines that a tool needs to be called, the response returns tool information in the `tool_calls` field. For usage examples, see the [Function calling guide](/developer-guides/text-generation/function-calling).

 Show child attributes

 [​ ](#toolstype) tools.typeenum<string> required Tool type. Currently supports only `function`.

 Available options:function [​ ](#toolsfunction) tools.functionobject required Show child attributes

 [​ ](#toolsfunctionname) tools.function.namestring required The tool name. Must contain only letters, digits, underscores, and hyphens. Up to 64 tokens.

 [​ ](#toolsfunctiondescription) tools.function.descriptionstring required A description of the tool. Helps the model determine when and how to call the tool.

 [​ ](#toolsfunctionparameters) tools.function.parametersobject default{} The tool&#x27;s parameters described using a valid JSON Schema. If empty, the tool has no input parameters.

 [​ ](#tool-choice) tool_choiceenum<string> default"auto" The tool selection policy. `auto` lets the model select. `none` disables tool calling. An object with `{"type": "function", "function": {"name": "..."}}` forces a specific tool. Models in thinking mode do not support forcing a specific tool.

 [​ ](#parallel-tool-calls) parallel_tool_callsboolean defaultfalse Specifies whether to enable parallel tool calling.

 [​ ](#enable-search) enable_searchboolean defaultfalse Enables web search. Enabling web search may increase token consumption.


**This is not a standard OpenAI parameter.** When using the Python SDK, include it in `extra_body`: `extra_body={"enable_search": True}`.

 [​ ](#search-options) search_optionsobject The web search strategy. Takes effect only when `enable_search` is `true`.


**Properties:**


- `forced_search` (boolean, default: `false`): Forces web search. `true` forcefully enables web search. `false` lets the model decide.

- `search_strategy` (string, default: `turbo`): The search scale strategy. `turbo` balances speed and effectiveness. `max` uses a more comprehensive strategy with multiple search engines. `agent` calls search and LLM multiple times for multi-round retrieval (applicable only to qwen3.7-max, qwen3.7-max-2026-06-08, qwen3.7-max-2026-05-20, qwen3.7-max-preview, qwen3.7-max-2026-05-17, qwen3.5-plus, qwen3.5-plus-2026-02-15, qwen3.5-flash, qwen3.5-flash-2026-02-23, qwen3-max, qwen3-max-2026-01-23, qwen3-max-2025-09-23, qwen3.5-omni-plus, qwen3.5-omni-plus-2026-03-15, qwen3.5-omni-flash, and qwen3.5-omni-flash-2026-03-15). `agent_max` adds web extraction support to the `agent` strategy (applicable only to the thinking mode of qwen3.7-max, qwen3.7-max-2026-06-08, qwen3.7-max-2026-05-20, qwen3.7-max-preview, qwen3.7-max-2026-05-17, qwen3-max, and qwen3-max-2026-01-23). When `agent` or `agent_max` is enabled, only return search sources (`enable_source: true`) is supported. All other web search features are unavailable.

- `enable_search_extension` (boolean, default: `false`): Enables domain-specific search.


**This is not a standard OpenAI parameter.** When using the Python SDK, include it in `extra_body`: `extra_body={"search_options": xxx}`.

 Show child attributes

 [​ ](#search-optionsforced-search) search_options.forced_searchboolean defaultfalse Forces web search. `true` forcefully enables web search. `false` lets the model decide.

 [​ ](#search-optionssearch-strategy) search_options.search_strategyenum<string> default"turbo" The search scale strategy. `turbo` (default) balances speed and effectiveness. `max` uses a more comprehensive strategy with multiple search engines; response time may be longer. `agent` calls web search and LLM multiple times for multi-round retrieval. Applicable only to qwen3.7-max, qwen3.7-max-2026-06-08, qwen3.7-max-2026-05-20, qwen3.7-max-preview, qwen3.7-max-2026-05-17, qwen3.5-plus, qwen3.5-plus-2026-02-15, qwen3.5-flash, qwen3.5-flash-2026-02-23, qwen3-max, qwen3-max-2026-01-23, qwen3-max-2025-09-23, qwen3.5-omni-plus, qwen3.5-omni-plus-2026-03-15, qwen3.5-omni-flash, and qwen3.5-omni-flash-2026-03-15. `agent_max` adds web extraction support to the `agent` strategy. Applicable only to the thinking mode of qwen3.7-max, qwen3.7-max-2026-06-08, qwen3.7-max-2026-05-20, qwen3.7-max-preview, qwen3.7-max-2026-05-17, qwen3-max, and qwen3-max-2026-01-23. When `agent` or `agent_max` is enabled, only return search sources (`enable_source: true`) is supported.

 Available options:turbo,max,agent,agent_max [​ ](#search-optionsenable-search-extension) search_options.enable_search_extensionboolean defaultfalse Enables domain-specific search.

 ### Response
200-application/json [​ ](#id) idstring The unique identifier for this request.

 [​ ](#choices) choicesobject[] An array of generated content from the model.

 Show child attributes

 [​ ](#choicesfinish-reason) choices.finish_reasonenum<string> The reason the model stopped generating. `stop`: stopped naturally or by stop parameter. `length`: reached maximum output length. `tool_calls`: stopped to call a tool.

 Available options:stop,length,tool_calls [​ ](#choicesindex) choices.indexinteger The index of this object in the choices array.

 [​ ](#choiceslogprobs) choices.logprobsobject | null Log probability information for tokens in the model&#x27;s output.

 Show child attributes

 [​ ](#choiceslogprobscontent) choices.logprobs.contentobject[] An array of tokens and their corresponding log probabilities.

 Show child attributes

 [​ ](#choiceslogprobscontenttoken) choices.logprobs.content.tokenstring The text of the current token.

 [​ ](#choiceslogprobscontentbytes) choices.logprobs.content.bytesinteger[] A list of raw UTF-8 bytes for the current token.

 [​ ](#choiceslogprobscontentlogprob) choices.logprobs.content.logprobnumber | null The log probability of the current token. `null` indicates an extremely low probability.

 [​ ](#choiceslogprobscontenttop-logprobs) choices.logprobs.content.top_logprobsobject[] The most likely candidate tokens for the current position.

 Show child attributes

 [​ ](#choiceslogprobscontenttop-logprobstoken) choices.logprobs.content.top_logprobs.tokenstring The candidate token text.

 [​ ](#choiceslogprobscontenttop-logprobsbytes) choices.logprobs.content.top_logprobs.bytesinteger[] Raw UTF-8 bytes for the token.

 [​ ](#choiceslogprobscontenttop-logprobslogprob) choices.logprobs.content.top_logprobs.logprobnumber | null The log probability of this candidate token.

 [​ ](#choicesmessage) choices.messageobject The message generated by the model.

 Show child attributes

 [​ ](#choicesmessagecontent) choices.message.contentstring The content of the model&#x27;s response.

 [​ ](#choicesmessagereasoning-content) choices.message.reasoning_contentstring The content of the model&#x27;s chain-of-thought reasoning.

 [​ ](#choicesmessagerefusal) choices.message.refusalstring | null This field is always `null`.

 [​ ](#choicesmessagerole) choices.message.roleenum<string> Always `assistant`.

 Available options:assistant [​ ](#choicesmessageaudio) choices.message.audioobject | null This field is always `null`.

 [​ ](#choicesmessagefunction-call) choices.message.function_callobject | null (Deprecated) This field is always `null`. Use `tool_calls` instead.

 [​ ](#choicesmessagetool-calls) choices.message.tool_callsobject[] Information about tools and their input parameters that the model generates after initiating a function call.

 Show child attributes

 [​ ](#choicesmessagetool-callsid) choices.message.tool_calls.idstring The unique identifier for this tool response.

 [​ ](#choicesmessagetool-callstype) choices.message.tool_calls.typeenum<string> The type of the tool. Currently, only `function` is supported.

 Available options:function [​ ](#choicesmessagetool-callsfunction) choices.message.tool_calls.functionobject Show child attributes

 [​ ](#choicesmessagetool-callsfunctionname) choices.message.tool_calls.function.namestring The name of the tool.

 [​ ](#choicesmessagetool-callsfunctionarguments) choices.message.tool_calls.function.argumentsstring The input parameters, formatted as a JSON string. Model outputs are non-deterministic. Validate the parameters before calling the function.

 [​ ](#choicesmessagetool-callsindex) choices.message.tool_calls.indexinteger The index of this tool in the `tool_calls` array.

 [​ ](#created) createdinteger The Unix timestamp, in seconds, when the request was created.

 [​ ](#model) modelstring The model used for this request.

 [​ ](#object) objectenum<string> Always `chat.completion`.

 Available options:chat.completion [​ ](#service-tier) service_tierstring | null Currently fixed as `null`.

 [​ ](#system-fingerprint) system_fingerprintstring | null Currently fixed as `null`.

 [​ ](#usage) usageobject Token consumption details for this request.

 Show child attributes

 [​ ](#usagecompletion-tokens) usage.completion_tokensinteger The number of tokens in the model&#x27;s output.

 [​ ](#usageprompt-tokens) usage.prompt_tokensinteger The number of tokens in the input.

 [​ ](#usagetotal-tokens) usage.total_tokensinteger The total number of tokens consumed (prompt_tokens + completion_tokens).

 [​ ](#usagecompletion-tokens-details) usage.completion_tokens_detailsobject A fine-grained breakdown of output tokens.

 Show child attributes

 [​ ](#usagecompletion-tokens-detailsaudio-tokens) usage.completion_tokens_details.audio_tokensinteger | null The number of audio tokens in the output (Qwen-Omni).

 [​ ](#usagecompletion-tokens-detailsreasoning-tokens) usage.completion_tokens_details.reasoning_tokensinteger | null The number of tokens in the thinking process.

 [​ ](#usagecompletion-tokens-detailstext-tokens) usage.completion_tokens_details.text_tokensinteger The number of text tokens in the output.

 [​ ](#usageprompt-tokens-details) usage.prompt_tokens_detailsobject A fine-grained breakdown of input tokens.

 Show child attributes

 [​ ](#usageprompt-tokens-detailsaudio-tokens) usage.prompt_tokens_details.audio_tokensinteger | null The number of tokens in the input audio.

 [​ ](#usageprompt-tokens-detailscached-tokens) usage.prompt_tokens_details.cached_tokensinteger The number of tokens that hit the cache.

 [​ ](#usageprompt-tokens-detailstext-tokens) usage.prompt_tokens_details.text_tokensinteger The number of text tokens in the input.

 [​ ](#usageprompt-tokens-detailsimage-tokens) usage.prompt_tokens_details.image_tokensinteger The number of image tokens in the input.

 [​ ](#usageprompt-tokens-detailsvideo-tokens) usage.prompt_tokens_details.video_tokensinteger The number of tokens for the input video.

 [​ ](#usageprompt-tokens-detailscache-creation) usage.prompt_tokens_details.cache_creationobject Information about explicit cache creation.

 Show child attributes

 [​ ](#usageprompt-tokens-detailscache-creationephemeral-5m-input-tokens) usage.prompt_tokens_details.cache_creation.ephemeral_5m_input_tokensinteger The number of tokens used to create the explicit cache.

 [​ ](#usageprompt-tokens-detailscache-creation-input-tokens) usage.prompt_tokens_details.cache_creation_input_tokensinteger The number of tokens used to create the explicit cache.

 [​ ](#usageprompt-tokens-detailscache-type) usage.prompt_tokens_details.cache_typestring When using an explicit cache, the value is `ephemeral`. Otherwise, this field does not exist.
