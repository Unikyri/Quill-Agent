# OpenAI responses

> **Source:** https://docs.qwencloud.com/api-reference/chat/openai-responses

POST/compatible-mode/v1/responses Python

 Basic call

 Copy ```\nimport os
from openai import OpenAI

client = OpenAI(
 # If environment variable is not set, replace with: api_key="sk-xxx"
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)

response = client.responses.create(
 model="qwen3.7-plus",
 input="What can you do?"
)

# Get model response
print(response.output_text)
``` 200

 application/json

 Copy ```\n{
 "created_at": 1771165900,
 "id": "f75c28fb-4064-48ed-90da-4d2cc4362xxx",
 "model": "qwen3.7-plus",
 "object": "response",
 "output": [
 {
 "content": [
 {
 "annotations": [],
 "text": "Hello! I am Qwen3.5, a large language model developed by Alibaba Cloud with knowledge up to 2026, designed to assist you with complex reasoning, creative tasks, and multilingual conversations.",
 "type": "output_text"
 }
 ],
 "id": "msg_89ad23e6-f128-4d4c-b7a1-a786e7880xxx",
 "role": "assistant",
 "status": "completed",
 "type": "message"
 }
 ],
 "parallel_tool_calls": false,
 "status": "completed",
 "tool_choice": "auto",
 "tools": [],
 "usage": {
 "input_tokens": 57,
 "input_tokens_details": {
 "cached_tokens": 0
 },
 "output_tokens": 44,
 "output_tokens_details": {
 "reasoning_tokens": 0
 },
 "total_tokens": 101,
 "x_details": [
 {
 "input_tokens": 57,
 "output_tokens": 44,
 "total_tokens": 101,
 "x_billing_type": "response_api"
 }
 ]
 }
}
``` The legacy URL path `/api/v2/apps/protocols/compatible-mode/v1/responses` will be deprecated soon. Please migrate to the new path `/compatible-mode/v1/responses` as soon as possible. 
## [​ ](#compatibility-with-openai) Compatibility with OpenAI

This API is OpenAI-compatible, but key differences exist in parameters, features, and behaviors.
Requests process only the parameters listed in this document. Any OpenAI parameters not mentioned are ignored.
Key differences:

- 
**Unsupported parameters**: Some parameters are not supported, such as `background` (synchronous only).


- 
**Additional parameters**: Supports extra parameters beyond OpenAI&#x27;s spec, such as `enable_thinking`.


- 
**Effortless context caching**: Enable automatic server-side caching for multi-turn conversations with a single request header. See [Session Cache](/developer-guides/text-generation/context-cache#session-cache).


 ### Authorizations
 [​ ](#authorization) Authorizationstring header required DashScope API key.

 ### Header Parameters
 [​ ](#x-dashscope-session-cache) x-dashscope-session-cacheenum<string> Controls [session cache](/developer-guides/text-generation/context-cache#session-cache) for multi-turn conversations using `previous_response_id`. When enabled, the server automatically caches conversation context to reduce latency and cost.


- `enable`: Enables session cache. Cache creation tokens are billed at 125% of the standard input price; cache hits at 10%. Cache is valid for 5 minutes (resets on hit). Minimum 1024 tokens for cache creation.

- `disable`: Disables session cache. Falls back to implicit cache if supported by the model.


Supported models: `qwen3.7-max`, `qwen3.7-max-2026-06-08`, `qwen3.7-max-2026-05-20`, `qwen3-max`, `qwen3.7-plus`, `qwen3.7-plus-2026-05-26`, `qwen3.6-plus`, `qwen3.5-plus`, `qwen3.5-flash`, `qwen-plus`, `qwen-flash`, `qwen3-coder-plus`, `qwen3-coder-flash`.


Pass via SDK: `default_headers` (Python) or `defaultHeaders` (Node.js).

 Available options:enable,disable ### Body
application/json [​ ](#model) modelstring required The model name. Supported models include qwen3.7-max, qwen3.7-max-2026-06-08, qwen3.7-max-2026-05-20, qwen3.7-max-preview, qwen3.7-max-2026-05-17, qwen3-max, qwen3-max-2026-01-23, qwen3.7-plus, qwen3.7-plus-2026-05-26, qwen3.6-plus, qwen3.6-plus-2026-04-02, qwen3.6-35b-a3b, qwen3.5-plus, qwen3.5-plus-2026-04-20, qwen3.5-plus-2026-02-15, qwen3.6-flash, qwen3.6-flash-2026-04-16, qwen3.5-flash, qwen3.5-flash-2026-02-23, qwen3.5-397b-a17b, qwen3.5-122b-a10b, qwen3.5-27b, qwen3.5-35b-a3b, qwen-plus, qwen-flash, qwen3-coder-plus, qwen3-coder-flash, qwen-plus-character, and qwen-flash-character.

 [​ ](#input) inputstring required The model input. Supports a plain text string or a message array arranged in conversational order.

 [​ ](#instructions) instructionsstring A system instruction inserted at the beginning of the context. When `previous_response_id` is used, the instructions specified in the previous turn are not carried over to the current context.

 [​ ](#previous-response-id) previous_response_idstring The unique ID of the previous response. The current response id is valid for 7 days. You can use this parameter to create a multi-turn conversation. The server automatically retrieves and combines the input and output of that turn as context. If both an input message array and `previous_response_id` are provided, the new messages in input are appended to the historical context. Cannot be used together with `conversation`. For usage examples, see the [Multi-turn conversations guide](/developer-guides/text-generation/multi-turn).

 [​ ](#conversation) conversationstring The conversation to which the current response belongs. Historical items in the conversation are automatically passed as context to the current request. The input and output of the current request are also automatically added to the conversation after the response completes. Cannot be used together with `previous_response_id`.

 [​ ](#stream) streamboolean defaultfalse Specifies whether to enable streaming output. If this parameter is set to `true`, the model response data is streamed back to the client in real time.

 [​ ](#tools) toolsobject[] A list of tools the model can use. Supported tool types: `web_search`, `code_interpreter`, `web_extractor`, `web_search_image`, `image_search`, `file_search`, `mcp`, `function`.


**Built-in tools** use `{"type": "<tool_name>"}` format. For example: `{"type": "web_search"}`.


**MCP tools** use the following format:


Copy ```\n{
 "type": "mcp",
 "server_protocol": "sse",
 "server_label": "amap-maps",
 "server_description": "AMAP MCP Server...",
 "server_url": "https://dashscope-intl.aliyuncs.com/api/v1/mcps/amap-maps/sse",
 "headers": {
 "Authorization": "Bearer $DASHSCOPE_API_KEY"
 }
}
``` 
**Function tools** use the following format:


Copy ```\n[{
 "type": "function",
 "name": "get_weather",
 "description": "Get weather information for a specified city",
 "parameters": {
 "type": "object",
 "properties": {
 "city": {
 "type": "string",
 "description": "The name of the city"
 }
 },
 "required": ["city"]
 }
}]
``` For usage examples, see the [Function calling guide](/developer-guides/text-generation/function-calling) and the [Web search guide](/developer-guides/text-generation/web-search).
``` Show child attributes

 [​ ](#toolstype) tools.typestring required The tool type. Valid values: `web_search`, `code_interpreter`, `web_extractor`, `web_search_image`, `image_search`, `file_search`, `mcp`, `function`.

 [​ ](#tool-choice) tool_choiceenum<string> Controls how the model selects and calls tools. Supports string format and object format.


**String format:**


- `auto`: The model automatically decides whether to call a tool.

- `none`: Prevents the model from calling any tool.

- `required`: Forces the model to call a tool. Available only when there is a single tool in the tools list.


**Object format:** Specifies the range of available tools for the model. The model can select and call tools only from the predefined list.

 [​ ](#temperature) temperaturenumber The sampling temperature that controls the diversity of the generated text. A higher temperature results in more diverse text. A lower temperature results in more deterministic text. Value range: [0, 2). Both `temperature` and `top_p` control the diversity of the generated text. We recommend that you set only one of them.

 [​ ](#top-p) top_pnumber The probability threshold for nucleus sampling that controls the diversity of the generated text. A higher `top_p` value results in more diverse text. A lower `top_p` value results in more deterministic text. Value range: (0, 1.0]. Both `temperature` and `top_p` control the diversity of the generated text. We recommend that you set only one of them.

 [​ ](#enable-thinking) enable_thinkingboolean Specifies whether to enable thinking mode. If set to `true`, the model thinks before replying. The thinking content is returned through an output item of the `reasoning` type. Reasoning tokens are counted in `output_tokens_details.reasoning_tokens` and are billed as reasoning tokens. When thinking mode is enabled, we recommend enabling the built-in tools to achieve the best model performance on complex tasks.


**This is not a standard OpenAI parameter.** The Python SDK passes it using `extra_body={"enable_thinking": True}`. The Node.js SDK and curl use `enable_thinking: true` directly as a top-level parameter.

 ### Response
200-application/json [​ ](#id) idstring The unique ID for this response. It is valid for 7 days. You can use this parameter in the `previous_response_id` parameter to create a multi-turn conversation.

 [​ ](#created-at) created_atnumber The Unix timestamp in seconds for this request.

 [​ ](#object) objectenum<string> The object type. The value is `response`.

 Available options:response [​ ](#status) statusenum<string> The status of the response generation.

 Available options:completed,failed,in_progress,cancelled,queued,incomplete [​ ](#model) modelstring The ID of the model that is used to generate the response.

 [​ ](#output) outputobject[] An array of output items generated by the model. The type and order of elements in the array depend on the model&#x27;s response.

 Show child attributes

 [​ ](#outputtype) output.typeenum<string> The type of the output item.


- `message`: Contains the final reply content generated by the model.

- `reasoning`: Returned when thinking mode (`enable_thinking: true`) is enabled. Reasoning tokens are counted in `output_tokens_details.reasoning_tokens` and are billed as reasoning tokens.

- `function_call`: Returned when a user-defined function tool is used. You need to handle the function call and return the result.

- `web_search_call`: Returned when the `web_search` tool is used.

- `code_interpreter_call`: Returned when the `code_interpreter` tool is used.

- `web_extractor_call`: Returned when the `web_extractor` tool is used. It must be used with the `web_search` tool.

- `web_search_image_call`: Returned when the `web_search_image` tool is used. It contains a list of searched images.

- `image_search_call`: Returned when the `image_search` tool is used. It contains a list of similar images.

- `mcp_call`: Returned when the `mcp` tool is used. It contains the result of the MCP service call.

- `file_search_call`: Returned when the `file_search` tool is used. It contains the search query and results from the knowledge base.


 Available options:message,reasoning,function_call,web_search_call,code_interpreter_call,web_extractor_call,web_search_image_call,image_search_call,mcp_call,file_search_call [​ ](#outputid) output.idstring The unique identifier for the output item. This field is included in all types of output items.

 [​ ](#outputrole) output.roleenum<string> The role of the message. The value is `assistant`. This field exists only when the type is `message`.

 Available options:assistant [​ ](#outputstatus) output.statusenum<string> The status of the output item.

 Available options:completed,in_progress [​ ](#outputname) output.namestring The tool or function name. This field exists when the type is `function_call`, `web_search_image_call`, `image_search_call`, or `mcp_call`. For `web_search_image_call` and `image_search_call`, the values are fixed as `web_search_image` and `image_search`, respectively. For `mcp_call`, the value is the specific function name called in the MCP service, such as `amap-maps-maps_geo`.

 [​ ](#outputarguments) output.argumentsstring The arguments for the tool call, in JSON string format. This field exists when the type is `function_call`, `web_search_image_call`, `image_search_call`, or `mcp_call`. You need to parse it using `JSON.parse()` before use.


**Arguments content by tool type:**


- `web_search_image_call`: `{"queries": ["search term 1", "search term 2"]}`, where queries is a list of search terms auto-generated by the model.

- `image_search_call`: `{"img_idx": 0, "bbox": [0, 0, 1000, 1000]}`, where `img_idx` is the index of the input image (starting from 0), and `bbox` is the bounding box coordinates `[x1, y1, x2, y2]` for the search area, with a range of 0-1000.

- `function_call`: An argument object generated according to the schema of the user-defined function parameters.

- `mcp_call`: An argument object for the function called in the MCP service.


 [​ ](#outputcall-id) output.call_idstring The unique identifier for the function call. This field exists only when the type is `function_call`. When returning the function call result, you must use this ID to associate the request with the response.

 [​ ](#outputcontent) output.contentobject[] An array of message content. This field exists only when the type is `message`.

 Show child attributes

 [​ ](#outputcontenttype) output.content.typeenum<string> The content type. The value is `output_text`.

 Available options:output_text [​ ](#outputcontenttext) output.content.textstring The text content that is generated by the model.

 [​ ](#outputcontentannotations) output.content.annotationsobject[] An array of text annotations. This is usually an empty array.

 [​ ](#outputsummary) output.summaryobject[] An array of reasoning summaries. This field exists only when the type is `reasoning`. Each element contains a `type` field with a value of `summary_text` and a `text` field that contains the summary text.

 Show child attributes

 [​ ](#outputsummarytype) output.summary.typeenum<string> Available options:summary_text [​ ](#outputsummarytext) output.summary.textstring [​ ](#outputaction) output.actionobject The search action information. This field exists only when the type is `web_search_call`.

 Show child attributes

 [​ ](#outputactionquery) output.action.querystring The search query keyword.

 [​ ](#outputactiontype) output.action.typeenum<string> The search type. The value is `search`.

 Available options:search [​ ](#outputactionsources) output.action.sourcesobject[] A list of search sources. Each element contains a `type` field and a `url` field.

 Show child attributes

 [​ ](#outputactionsourcestype) output.action.sources.typestring [​ ](#outputactionsourcesurl) output.action.sources.urlstring [​ ](#outputcode) output.codestring The code that is generated and executed by the model. This field exists only when the type is `code_interpreter_call`.

 [​ ](#outputoutputs) output.outputsobject[] An array of code execution outputs. This field exists only when the type is `code_interpreter_call`. Each element contains a `type` field with a value of `logs` and a `logs` field that contains the code execution log.

 Show child attributes

 [​ ](#outputoutputstype) output.outputs.typeenum<string> Available options:logs [​ ](#outputoutputslogs) output.outputs.logsstring [​ ](#outputcontainer-id) output.container_idstring The identifier for the code interpreter container. This field exists only when the type is `code_interpreter_call`. It is used to associate multiple code executions within the same session.

 [​ ](#outputgoal) output.goalstring A description of the extraction goal that explains what information needs to be extracted from the web page. This field exists only when the type is `web_extractor_call`.

 [​ ](#outputoutput) output.outputstring The output result of the tool call, in string format.


- When type is `web_extractor_call`, this is the summary of the extracted web content.

- When type is `web_search_image_call` or `image_search_call`, this is a JSON string containing an array of image search results. Each element contains a `title` field (image title), a `url` field (image URL), and an `index` field (sequence number).

- When type is `mcp_call`, this is the JSON string result returned by the MCP service.


 [​ ](#outputurls) output.urlsstring[] A list of URLs of the extracted web pages. This field exists only when the type is `web_extractor_call`.

 [​ ](#outputserver-label) output.server_labelstring The MCP service label. This field exists only when the type is `mcp_call`. It identifies the MCP service used for this call.

 [​ ](#outputqueries) output.queriesstring[] The list of queries used for knowledge base retrieval. This field exists only when the type is `file_search_call`. Array elements are strings representing search queries generated by the model.

 [​ ](#outputresults) output.resultsobject[] An array of knowledge base retrieval results. This field exists only when the type is `file_search_call`.

 Show child attributes

 [​ ](#outputresultsfile-id) output.results.file_idstring The file ID of the matching document.

 [​ ](#outputresultsfilename) output.results.filenamestring The filename of the matching document.

 [​ ](#outputresultsscore) output.results.scorenumber The relevance score of the match, ranging from 0 to 1. A higher value indicates greater relevance.

 [​ ](#outputresultstext) output.results.textstring A snippet of the matched document content.

 [​ ](#parallel-tool-calls) parallel_tool_callsboolean Whether parallel tool calls are enabled.

 [​ ](#tool-choice) tool_choicestring The value of the `tool_choice` parameter from the echo request. Valid values are `auto`, `none`, and `required`.

 [​ ](#tools) toolsobject[] The complete content of the tools parameter from the echo request. The structure is the same as the tools parameter in the request body.

 [​ ](#error) errorobject | null The error object that is returned when the model fails to generate a response. This field is `null` on success.

 Show child attributes

 [​ ](#errorcode) error.codestring The error code.

 [​ ](#errormessage) error.messagestring A human-readable error message.

 [​ ](#usage) usageobject The token consumption information for this request.

 Show child attributes

 [​ ](#usageinput-tokens) usage.input_tokensinteger The number of input tokens.

 [​ ](#usageoutput-tokens) usage.output_tokensinteger The number of tokens that are output by the model.

 [​ ](#usagetotal-tokens) usage.total_tokensinteger The total number of tokens consumed. This is the sum of `input_tokens` and `output_tokens`.

 [​ ](#usageinput-tokens-details) usage.input_tokens_detailsobject The fine-grained categorization of input tokens.

 Show child attributes

 [​ ](#usageinput-tokens-detailscached-tokens) usage.input_tokens_details.cached_tokensinteger The number of tokens that hit the cache.

 [​ ](#usageoutput-tokens-details) usage.output_tokens_detailsobject The fine-grained categorization of output tokens.

 Show child attributes

 [​ ](#usageoutput-tokens-detailsreasoning-tokens) usage.output_tokens_details.reasoning_tokensinteger The number of tokens in the thinking process.

 [​ ](#usagex-details) usage.x_detailsobject[] Detailed token breakdown by billing type.

 Show child attributes

 [​ ](#usagex-detailsinput-tokens) usage.x_details.input_tokensinteger The number of input tokens.

 [​ ](#usagex-detailsoutput-tokens) usage.x_details.output_tokensinteger The number of tokens that are output by the model.

 [​ ](#usagex-detailstotal-tokens) usage.x_details.total_tokensinteger The total number of tokens consumed.

 [​ ](#usagex-detailsx-billing-type) usage.x_details.x_billing_typestring The value is `response_api`.

 [​ ](#usagex-tools) usage.x_toolsobject Statistical information about tool usage. If built-in tools are used, this field includes the number of calls for each tool. Example: `{"web_search": {"count": 1}}`
