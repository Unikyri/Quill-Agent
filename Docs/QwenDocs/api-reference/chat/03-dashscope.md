# DashScope chat

> **Source:** https://docs.qwencloud.com/api-reference/chat/dashscope

POST/api/v1/services/aigc/text-generation/generation Python

 Text input

 Copy ```\nimport os
import dashscope

dashscope.base_http_api_url = 'https://dashscope-intl.aliyuncs.com/api/v1'
messages = [
 {'role': 'system', 'content': 'You are a helpful assistant.'},
 {'role': 'user', 'content': 'Who are you?'}
]
response = dashscope.Generation.call(
 api_key=os.getenv('DASHSCOPE_API_KEY'),
 model='qwen-plus',
 messages=messages,
 result_format='message'
)
print(response)
``` 200 400 401 429 Copy ```\n{
 "status_code": 200,
 "request_id": "902fee3b-f7f0-9a8c-96a1-6b4ea25af114",
 "code": "",
 "message": "",
 "output": {
 "text": null,
 "finish_reason": null,
 "choices": [
 {
 "finish_reason": "stop",
 "message": {
 "role": "assistant",
 "content": "I am a large-scale language model developed by Alibaba Cloud. My name is Qwen.",
 "tool_calls": null,
 "reasoning_content": null
 }
 }
 ]
 },
 "usage": {
 "input_tokens": 22,
 "output_tokens": 17,
 "total_tokens": 39,
 "image_tokens": null,
 "video_tokens": null,
 "audio_tokens": null
 }
}
``` [Get an API key](/api-reference/preparation/api-key) and [set it as an environment variable](/api-reference/preparation/export-api-key-env). To use the SDK, [install it](/api-reference/preparation/install-sdk). 
## [​ ](#endpoint) Endpoint


- HTTP (text-only, such as `qwen-plus`): `POST https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/text-generation/generation`

- HTTP (multimodal, such as `qwen3.7-plus`, `qwen3-vl-plus`): `POST https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation`

- SDK `base_http_api_url`: `https://dashscope-intl.aliyuncs.com/api/v1`


**Python SDK:**
Copy ```\ndashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

``` 
**Java SDK:**
Copy ```\n// Option 1: Set during instantiation
import com.alibaba.dashscope.protocol.Protocol;
Generation gen = new Generation(Protocol.HTTP.getValue(), "https://dashscope-intl.aliyuncs.com/api/v1");

// Option 2: Set globally
import com.alibaba.dashscope.utils.Constants;
Constants.baseHttpApiUrl = "https://dashscope-intl.aliyuncs.com/api/v1";

``` 

---
 ### Authorizations
 [​ ](#authorization) Authorizationstring header required Your DashScope API key. See [Get API key](/api-reference/preparation/api-key) for details.

 ### Body
application/json [​ ](#model) modelstring required The name of the model to call. Supports Qwen large language models (commercial and open-source), Qwen-Coder, and math models. For a list of models, see [Text generation — Qwen](/developer-guides/getting-started/text-generation-models).

 Example:qwen-plus [​ ](#input) inputobject required The input to the model.

 Show child attributes

 [​ ](#inputmessages) input.messages(System message · object | User message · object | Assistant message · object | Tool message · object)[] required The conversation context, provided as an ordered list of messages. Each message is a system, user, assistant, or tool message object.

 Sets the role, tone, task objective, or constraints for the model. Usually placed first in the messages array. Do not set for QwQ models. Does not take effect for QVQ models.

 - System message 
- User message 
- Assistant message 
- Tool message 

 Show child attributes

 [​ ](#inputmessagesrole) input.messages.roleenum<string> required Fixed as `system`.

 Available options:system [​ ](#inputmessagescontent) input.messages.contentstring required The system message content that sets context for the model.

 [​ ](#parameters) parametersobject Optional generation parameters for text models.

 Show child attributes

 [​ ](#parametersresult-format) parameters.result_formatenum<string> default"text" The format of the returned data. Set to `message` for multi-turn conversations.


**Default values:** `text` for most models, except Qwen3-Max, Qwen3-VL, QwQ, and Qwen3 open source models (excluding qwen3-next-80b-a3b-instruct), which default to `message`.


When the model is Qwen-VL/QVQ, setting `text` has no effect. For Qwen3-Max, Qwen3-VL, and Qwen3 models in thinking mode, this can only be set to `message`.

 Available options:message,text [​ ](#parameterstemperature) parameters.temperaturenumber Sampling temperature. Controls output diversity. Higher values produce more diverse output; lower values produce more deterministic output. Range: [0, 2).


Do not modify the default temperature value for QVQ models.

 Required range:0 <= x < 2 [​ ](#parameterstop-p) parameters.top_pnumber Nucleus sampling threshold. Higher values produce more diverse output. Range: (0, 1.0].


**Default values by model:**


- Qwen3.5 (non-thinking), Qwen3 (non-thinking), Qwen3-Instruct series, Qwen3-Coder series, qwen-max series, qwen-plus series (non-thinking), qwen-flash series (non-thinking), qwen-turbo series (non-thinking), qwen open source series, qwen-vl-max-2025-08-13, Qwen3-VL (non-thinking): **0.8**

- qwen-vl-plus series, qwen-vl-max, qwen-vl-max-latest, qwen-vl-max-2025-04-08, qwen2.5-vl-3b/7b/32b/72b-instruct: **0.001**

- QVQ series, qwen-vl-plus-2025-07-10, qwen-vl-plus-2025-08-15: **0.5**

- qwen3-max-preview (thinking mode), Qwen3-Omni-Flash series: **1.0**

- Qwen3.5 (thinking), Qwen3 (thinking), Qwen3-VL (thinking), Qwen3-Thinking, QwQ series, Qwen3-Omni-Captioner: **0.95**


Do not modify the default `top_p` value for QVQ models.

 Required range:0 < x <= 1 [​ ](#parameterstop-k) parameters.top_kinteger The size of the candidate token set for sampling. A larger value increases randomness; a smaller value increases determinism. If `None` or greater than 100, top_k is disabled and only top_p takes effect. Must be >= 0.


**Default values by model:**


- QVQ series, qwen-vl-plus-2025-07-10, qwen-vl-plus-2025-08-15: **10**

- QwQ series: **40**

- Other Qwen-VL-Plus series, Qwen-VL-Max models released before August 13 2025, qwen2.5-omni-7b: **1**

- Qwen3-Omni-Flash series: **50**

- All other models: **20**


Do not modify the default `top_k` value for QVQ models.

 Required range:x >= 0 [​ ](#parametersmax-tokens) parameters.max_tokensinteger deprecated **Deprecated.** Use `max_completion_tokens` for new integrations.


The maximum length of the model&#x27;s answer, which excludes chain-of-thought content. That is: Model answer = Model output - Chain-of-thought (if any). Does not limit thinking chain length. Default is the model&#x27;s maximum output length.

 [​ ](#parametersmax-completion-tokens) parameters.max_completion_tokensinteger The maximum number of tokens in this response, including chain-of-thought tokens. Generation stops when this limit is reached, and `finish_reason` is set to `length`. The default and maximum values correspond to the model&#x27;s maximum output length.


**Difference from `max_tokens`:** `max_completion_tokens` limits the total length of both the chain-of-thought and the final response, while `max_tokens` does not limit the chain-of-thought length. For reasoning models, `max_completion_tokens` is recommended.


**Supported models:**


- Qwen-Max: Qwen3.7-Max and later

- Qwen-Plus: Qwen3.5-Plus and later

- Qwen-Flash: Qwen3.5-Flash and later

- DeepSeek: deepseek-v3, deepseek-r1, deepseek-r1-0528, deepseek-v3.1, deepseek-v3.2, deepseek-v3.2-exp, deepseek-v4-pro, deepseek-v4-flash and later


The actual number of output tokens may differ from the configured value by up to 10 tokens.


For HTTP calls, place `max_completion_tokens` in the `parameters` field.

 [​ ](#parametersstream) parameters.streamboolean defaultfalse Whether to stream the response. For HTTP streaming, also set the `X-DashScope-SSE: enable` header. For Java SDK streaming, use the `streamCall` interface.

 [​ ](#parametersincremental-output) parameters.incremental_outputboolean defaultfalse When streaming, whether to return only the new delta tokens (true) or the full accumulated text so far (false).

 [​ ](#parametersenable-thinking) parameters.enable_thinkingboolean Whether to enable thinking mode. Applies to hybrid thinking models: Qwen3.7, Qwen3.6, Qwen3.5, Qwen3, and Qwen3-VL series, as well as the DeepSeek-V4-Pro/V4-Flash series, DeepSeek-V3.2/V3.2-exp/V3.1 series, Kimi-K2.6/K2.5 series, and GLM series. When enabled, thinking content is returned in the `reasoning_content` field.

 [​ ](#parametersthinking-budget) parameters.thinking_budgetinteger Maximum length of the thinking chain. Applies to commercial and open source versions of Qwen3.5, Qwen3-VL, and Qwen3. Default is the model&#x27;s maximum chain-of-thought length.

 [​ ](#parameterspreserve-thinking) parameters.preserve_thinkingboolean defaultfalse Whether to append the `reasoning_content` from assistant messages in the conversation history to the model input. Useful when the model needs to refer to the historical thinking process.


Currently supported by: qwen3.7-max, qwen3.7-max-2026-05-20 and subsequent snapshots, qwen3.6-max-preview, qwen3.7-plus, qwen3.7-plus-2026-05-26, qwen3.6-plus, qwen3.6-plus-2026-04-02, qwen3.6-flash, kimi-k2.6, kimi-k2.7-code-highspeed, and kimi-k2.7-code.


- If the historical messages do not contain `reasoning_content`, enabling this parameter does not cause an error.

- When enabled, the `reasoning_content` from the conversation history is included in the input token count and is billed.


For HTTP calls, place `preserve_thinking` in the `parameters` object.

 [​ ](#parametersreasoning-effort) parameters.reasoning_effortstring default"high" Controls the inference effort of DeepSeek-V4 and GLM series models. Valid values: `high` (high-effort inference) and `max` (maximum-effort inference). `low` and `medium` are mapped to `high`, and `xhigh` is mapped to `max`.


Applies to: glm-5.2, glm-5.1, glm-5, deepseek-v4-pro, and deepseek-v4-flash.


For HTTP calls, place `reasoning_effort` in the `parameters` object.

 [​ ](#parameterstool-stream) parameters.tool_streamboolean defaultfalse Controls streaming behavior for complex tool arguments. Only takes effect in streaming mode. This parameter only affects the streaming output behavior of complex tool parameters. Simple tool parameters, where all parameter types are strings, can be streamed as long as streaming calls are enabled; `tool_stream` has no effect on them. Complex tools are tools where some parameter types in the tool definition are arrays or objects.


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


For HTTP calls, place `tool_stream` in the `parameters` object.

 [​ ](#parametersenable-code-interpreter) parameters.enable_code_interpreterboolean defaultfalse Whether to enable the code interpreter feature. For more information, see Code interpreter.

 [​ ](#parametersclear-thinking) parameters.clear_thinkingboolean defaultfalse Controls whether to use the `reasoning_content` (thinking process) from previous turns as context input for the model in a multi-turn conversation. Supported only by GLM series models: glm-5.2, glm-5.1, glm-5, and glm-4.7.


- `true`: Ignores the `reasoning_content` from previous turns and uses only visible text, tool calls, and results as context input. This can reduce context length and cost.

- `false` (default): Retains the `reasoning_content` from previous turns and provides it to the model along with the context.


 [​ ](#parametersrepetition-penalty) parameters.repetition_penaltynumber Penalty for token repetition. A value of 1.0 means no penalty. Higher values reduce repetition. Must be a positive number.


When using the `qwen-vl-plus_2025-01-25` model for text extraction, set `repetition_penalty` to 1.0. Do not modify the default `repetition_penalty` value for QVQ models.

 [​ ](#parameterspresence-penalty) parameters.presence_penaltynumber Controls how much the model avoids repeating content already present in the text. Range: [-2.0, 2.0]. Positive values reduce repetition; negative values increase it.


**Default values by model:**


- Qwen3.5 (non-thinking), qwen3-max-preview (thinking), Qwen3 (non-thinking), Qwen3-Instruct series, qwen3-0.6b/1.7b/4b (thinking), QVQ series, qwen-max, qwen-max-latest, qwen2.5-vl series, qwen-vl-max series, qwen-vl-plus, Qwen3-VL (non-thinking): **1.5**

- qwen-vl-plus-latest, qwen-vl-plus-2025-08-15: **1.2**

- qwen-vl-plus-2025-01-25: **1.0**

- qwen3-8b/14b/32b/30b-a3b/235b-a22b (thinking), qwen-plus/qwen-plus-latest/2025-04-28 (thinking), qwen-turbo/qwen-turbo/2025-04-28 (thinking): **0.5**

- All other models: **0.0**


When using `qwen-vl-plus-2025-01-25` for text extraction, set `presence_penalty` to 1.5. Do not modify the default for QVQ models.

 Required range:-2 <= x <= 2 [​ ](#parametersseed) parameters.seedinteger Random seed for reproducible results. Range: [0, 2³¹−1]. With the same seed and parameters, the model returns the same result whenever possible.

 Required range:x >= 0 [​ ](#parametersstop) parameters.stopstring Stop sequences. When the generated text contains a specified string or token ID, generation stops immediately. Do not mix strings and token IDs in the same array. Not supported by all models; check model documentation.

 [​ ](#parameterstools) parameters.toolsobject[] An array of tool objects for function calling. When using tools, you must set `result_format` to `message`. Not supported by qwen-vl series models. For usage examples, see the [Function calling guide](/developer-guides/text-generation/function-calling).

 Show child attributes

 [​ ](#parameterstoolstype) parameters.tools.typeenum<string> required The type of tool. Currently only `function` is supported.

 Available options:function [​ ](#parameterstoolsfunction) parameters.tools.functionobject required Show child attributes

 [​ ](#parameterstoolsfunctionname) parameters.tools.function.namestring required The name of the tool function. Can contain letters, numbers, underscores, and hyphens. Maximum 64 characters.

 [​ ](#parameterstoolsfunctiondescription) parameters.tools.function.descriptionstring required A description of the tool function that helps the model decide when and how to call it.

 [​ ](#parameterstoolsfunctionparameters) parameters.tools.function.parametersobject A JSON Schema object describing the function parameters. Defaults to `{}`.

 [​ ](#parameterstool-choice) parameters.tool_choiceenum<string> default"auto" Defines the tool selection strategy. Thinking mode models do not support forcing a specific tool.

 [​ ](#parametersparallel-tool-calls) parameters.parallel_tool_callsboolean defaultfalse Whether to enable parallel tool calls. Not supported by thinking mode models when forcing a specific tool. See [Parallel tool calls](/developer-guides/text-generation/function-calling).

 [​ ](#parametersresponse-format) parameters.response_formatobject default{"type":"text"} The format of the returned content. If set to `json_object`, you must instruct the model to output JSON in the prompt.

 Show child attributes

 [​ ](#parametersresponse-formattype) parameters.response_format.typeenum<string> The output format type. `text`: plain text. `json_object`: standard JSON string. `json_schema`: JSON matching the provided schema.

 Available options:text,json_object,json_schema [​ ](#parametersresponse-formatjson-schema) parameters.response_format.json_schemaobject Required when `type` is `json_schema`. Defines the JSON Schema for structured output. For supported models, see [Structured output](/developer-guides/text-generation/structured-output).

 Show child attributes

 [​ ](#parametersresponse-formatjson-schemaname) parameters.response_format.json_schema.namestring Unique schema name (letters, numbers, underscores, hyphens; max 64 characters).

 [​ ](#parametersresponse-formatjson-schemadescription) parameters.response_format.json_schema.descriptionstring Description of the schema&#x27;s purpose.

 [​ ](#parametersresponse-formatjson-schemaschema) parameters.response_format.json_schema.schemaobject A JSON Schema object defining the output data structure.

 [​ ](#parametersresponse-formatjson-schemastrict) parameters.response_format.json_schema.strictboolean defaultfalse Whether the model must strictly adhere to all schema constraints. `true` is recommended.

 [​ ](#parameterslogprobs) parameters.logprobsboolean defaultfalse Whether to return log probabilities of the output tokens. Supported by: snapshot models of qwen-plus/qwen-turbo series; qwen3-vl-plus/qwen3-vl-flash series; Qwen3 open source models. See model page for supported models.

 [​ ](#parameterstop-logprobs) parameters.top_logprobsinteger default0 Number of most likely candidate tokens to return at each generation step. Valid values: 0–5. Only takes effect when `logprobs` is `true`. Supported by the same models as `logprobs`.

 Required range:0 <= x <= 5 [​ ](#parametersn) parameters.ninteger default1 The number of responses to generate. Range: 1–4. Currently only non-thinking mode Qwen3 models are supported. Fixed at 1 when `tools` is specified. Increases output token consumption.

 Required range:1 <= x <= 4 [​ ](#parametersvl-high-resolution-images) parameters.vl_high_resolution_imagesboolean defaultfalse Whether to enable high-resolution image processing. When enabled, uses a fixed-resolution strategy where `max_pixels` is ignored. Default: false.


**Pixel limits when enabled (true):**


- Qwen3.5 series, Qwen3-VL series, qwen-vl-max, qwen-vl-max-latest, qwen-vl-max-0813, qwen-vl-plus, qwen-vl-plus-latest, qwen-vl-plus-0815: fixed at **16777216** pixels (16384 tokens × 32×32 pixels)

- QVQ series and other Qwen2.5-VL series: fixed at **12845056** pixels (16384 tokens × 28×28 pixels)


When `false`, the pixel limit is determined by `max_pixels`.

 [​ ](#parametersvl-enable-image-hw-output) parameters.vl_enable_image_hw_outputboolean defaultfalse Whether to return the dimensions of the scaled image in the response (`image_hw` field). When streaming, returned in the last chunk. Applies to Qwen-VL series models.

 ### Response
200-application/json [​ ](#status-code) status_codeinteger The status code of the request. `200` indicates success. The Java SDK does not return this field; if a call fails, an exception is thrown containing the status_code.

 [​ ](#request-id) request_idstring A unique identifier for this request. In the Java SDK, this is `requestId`.

 [​ ](#code) codestring The error code. Empty string if the request was successful. Only the Python SDK returns this field.

 [​ ](#message) messagestring A human-readable error message. Empty string if the request was successful.

 [​ ](#output) outputobject The model&#x27;s output.

 Show child attributes

 [​ ](#outputtext) output.textstring | null The generated text. Returned when `result_format` is `text`.

 [​ ](#outputfinish-reason) output.finish_reasonstring | null The reason generation stopped. Returned when `result_format` is `text`. Values: `null` (still generating), `stop` (natural end or stop condition triggered), `length` (max tokens reached), `tool_calls` (tool call triggered).

 [​ ](#outputchoices) output.choicesobject[] The output choices. Returned when `result_format` is `message`.

 Show child attributes

 [​ ](#outputchoicesfinish-reason) output.choices.finish_reasonstring | null The reason generation stopped. Values: `null` (generating), `stop`, `length`, `tool_calls`.

 [​ ](#outputchoicesmessage) output.choices.messageobject The assistant&#x27;s output message.

 Show child attributes

 [​ ](#outputchoicesmessagerole) output.choices.message.rolestring Always `assistant`.

 [​ ](#outputchoicesmessagecontent) output.choices.message.contentstring The message content. A string for text models; an array for Qwen-VL/Qwen-Audio models. Empty when `tool_calls` is present.

 [​ ](#outputchoicesmessagereasoning-content) output.choices.message.reasoning_contentstring | null The deep thinking content. Returned when thinking mode is enabled.

 [​ ](#outputchoicesmessagetool-calls) output.choices.message.tool_callsobject[] | null Tool calls the model wants to make. Present when the model triggers a function call.

 Show child attributes

 [​ ](#outputchoicesmessagetool-callsid) output.choices.message.tool_calls.idstring The ID of the tool call.

 [​ ](#outputchoicesmessagetool-callstype) output.choices.message.tool_calls.typeenum<string> The type of tool. Currently only `function` is supported.

 Available options:function [​ ](#outputchoicesmessagetool-callsfunction) output.choices.message.tool_calls.functionobject Show child attributes

 [​ ](#outputchoicesmessagetool-callsfunctionname) output.choices.message.tool_calls.function.namestring The name of the tool function.

 [​ ](#outputchoicesmessagetool-callsfunctionarguments) output.choices.message.tool_calls.function.argumentsstring The input parameters for the tool, as a JSON string.

 [​ ](#outputchoicesmessagetool-callsindex) output.choices.message.tool_calls.indexinteger The index of this tool call in the tool_calls array.

 [​ ](#outputchoiceslogprobs) output.choices.logprobsobject | null Log probability information for this choice. Returned when `logprobs` is `true`.

 Show child attributes

 [​ ](#outputchoiceslogprobscontent) output.choices.logprobs.contentobject[] An array of tokens with log probability information.

 Show child attributes

 [​ ](#outputchoiceslogprobscontenttoken) output.choices.logprobs.content.tokenstring [​ ](#outputchoiceslogprobscontentbytes) output.choices.logprobs.content.bytesinteger[] [​ ](#outputchoiceslogprobscontentlogprob) output.choices.logprobs.content.logprobnumber | null [​ ](#outputchoiceslogprobscontenttop-logprobs) output.choices.logprobs.content.top_logprobsobject[] Show child attributes

 [​ ](#outputchoiceslogprobscontenttop-logprobstoken) output.choices.logprobs.content.top_logprobs.tokenstring [​ ](#outputchoiceslogprobscontenttop-logprobsbytes) output.choices.logprobs.content.top_logprobs.bytesinteger[] [​ ](#outputchoiceslogprobscontenttop-logprobslogprob) output.choices.logprobs.content.top_logprobs.logprobnumber | null [​ ](#usage) usageobject Token usage information for this request.

 Show child attributes

 [​ ](#usageinput-tokens) usage.input_tokensinteger Number of tokens in the user input.

 [​ ](#usageoutput-tokens) usage.output_tokensinteger Number of tokens in the model output.

 [​ ](#usagetotal-tokens) usage.total_tokensinteger Total tokens (input + output). Returned for plain text input.

 [​ ](#usageimage-tokens) usage.image_tokensinteger | null Number of tokens in the input image. Returned when the input includes an image.

 [​ ](#usagevideo-tokens) usage.video_tokensinteger | null Number of tokens in the input video. Returned when the input includes a video.

 [​ ](#usageaudio-tokens) usage.audio_tokensinteger | null Number of tokens in the input audio. Returned when the input includes audio.

 [​ ](#usageinput-tokens-details) usage.input_tokens_detailsobject Detailed input token breakdown for Qwen-VL and QVQ models.

 Show child attributes

 [​ ](#usageinput-tokens-detailstext-tokens) usage.input_tokens_details.text_tokensinteger [​ ](#usageinput-tokens-detailsimage-tokens) usage.input_tokens_details.image_tokensinteger [​ ](#usageinput-tokens-detailsvideo-tokens) usage.input_tokens_details.video_tokensinteger [​ ](#usageoutput-tokens-details) usage.output_tokens_detailsobject Detailed output token breakdown.

 Show child attributes

 [​ ](#usageoutput-tokens-detailstext-tokens) usage.output_tokens_details.text_tokensinteger Number of tokens in the output text.

 [​ ](#usageoutput-tokens-detailsreasoning-tokens) usage.output_tokens_details.reasoning_tokensinteger Number of tokens in the thinking process.

 [​ ](#usageprompt-tokens-details) usage.prompt_tokens_detailsobject Fine-grained classification of input tokens.

 Show child attributes

 [​ ](#usageprompt-tokens-detailscached-tokens) usage.prompt_tokens_details.cached_tokensinteger Number of tokens that hit the cache. See [Context cache](/developer-guides/text-generation/context-cache).

 [​ ](#usageprompt-tokens-detailscache-creation-input-tokens) usage.prompt_tokens_details.cache_creation_input_tokensinteger Number of tokens used to create the explicit cache.

 [​ ](#usageprompt-tokens-detailscache-type) usage.prompt_tokens_details.cache_typestring If explicit caching is used, the value is `ephemeral`. Otherwise not returned.

 [​ ](#usageprompt-tokens-detailscache-creation) usage.prompt_tokens_details.cache_creationobject Information about explicit cache creation.

 Show child attributes

 [​ ](#usageprompt-tokens-detailscache-creationephemeral-5m-input-tokens) usage.prompt_tokens_details.cache_creation.ephemeral_5m_input_tokensinteger The number of tokens used to create a 5-minute explicit cache.
