# Wan v2 — Generate an image

> **Source:** https://docs.qwencloud.com/api-reference/image-generation/wan-text-to-image-v2/create-task

POST/services/aigc/image-generation/generation cURL

 cURL - wan2.6-t2i (async)

 Copy ```\ncurl --location 'https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/image-generation/generation' \
--header 'Content-Type: application/json' \
--header "Authorization: Bearer $DASHSCOPE_API_KEY" \
--header 'X-DashScope-Async: enable' \
--data '{
 "model": "wan2.6-t2i",
 "input": {
 "messages": [
 {
 "role": "user",
 "content": [
 {
 "text": "A flower shop with exquisite windows, a beautiful wooden door, and flowers on display"
 }
 ]
 }
 ]
 },
 "parameters": {
 "prompt_extend": true,
 "watermark": false,
 "n": 1,
 "negative_prompt": "",
 "size": "1280*1280"
 }
}'
``` 200 400 Copy ```\n{
 "request_id": "&#x3C;string>",
 "output": {
 "task_id": "&#x3C;string>",
 "task_status": "PENDING"
 }
}
``` [Get an API key](/api-reference/preparation/api-key) and [set it as an environment variable](/api-reference/preparation/export-api-key-env). If using the SDK, [install it](/api-reference/preparation/install-sdk). 
## [​ ](#supported-models) Supported models

The Wan text-to-image family uses different endpoints and request formats by version:
ModelEndpointInput formatResolutionMax prompt length`wan2.6-t2i``/services/aigc/image-generation/generation``messages` array1280*1280 to 1440*1440, ratio 1:4 to 4:12,100 chars`wan2.5-t2i-preview``/services/aigc/text2image/image-synthesis``prompt` string1280*1280 to 1440*1440, ratio 1:4 to 4:12,000 chars`wan2.2-t2i-plus``/services/aigc/text2image/image-synthesis``prompt` string512–1440 per side, max 1440*1440500 chars`wan2.2-t2i-flash``/services/aigc/text2image/image-synthesis``prompt` string512–1440 per side, max 1440*1440500 chars`wan2.1-t2i-plus``/services/aigc/text2image/image-synthesis``prompt` string512–1440 per side, max 1440*1440500 chars`wan2.1-t2i-turbo``/services/aigc/text2image/image-synthesis``prompt` string512–1440 per side, max 1440*1440500 chars`wanx2.0-t2i-turbo``/services/aigc/text2image/image-synthesis``prompt` string512–1440 per side, max 1440*1440800 chars 
 **wan2.6-t2i** also supports a [synchronous endpoint](/api-reference/image-generation/wan-text-to-image-v2/synchronous) (single-request, immediate response). 
## [​ ](#sdk-version-requirements) SDK version requirements


- **wan2.6-t2i**: DashScope Python SDK **1.25.7+**, Java SDK **2.22.6+**

- **wan2.5 and earlier**: DashScope Python SDK **1.25.2+**, Java SDK **2.22.2+**


 ### Authorizations
 [​ ](#authorization) Authorizationstring header required DashScope API Key. Create one in the [Qwen Cloud console](https://home.qwencloud.com/api-keys).

 ### Header Parameters
 [​ ](#x-dashscope-async) X-DashScope-Asyncenum<string> required Must be `enable` to create an asynchronous task.

 Available options:enable ### Body
application/json object [​ ](#model) modelenum<string> required The model name. For the wan2.6-t2i model, use `wan2.6-t2i`.

 Available options:wan2.6-t2i Example:wan2.6-t2i [​ ](#input) inputobject required The input object containing the message array.

 Show child attributes

 [​ ](#inputmessages) input.messagesobject[] required An array for the request content. Only single-turn conversations are supported (one set of role and content).

 Show child attributes

 [​ ](#inputmessagesrole) input.messages.roleenum<string> required The role of the message. Must be `user`.

 Available options:user [​ ](#inputmessagescontent) input.messages.contentobject[] required An array for the message content.

 Show child attributes

 [​ ](#inputmessagescontenttext) input.messages.content.textstring required The positive prompt describing the content, style, and composition of the image to generate. Chinese and English are supported. Maximum 2,100 characters. Only one text is allowed per request.

 Example:A flower shop with exquisite windows, a beautiful wooden door, and flowers on display [​ ](#parameters) parametersobject Parameters for wan2.6-t2i model.

 Show child attributes

 [​ ](#parametersnegative-prompt) parameters.negative_promptstring Optional. Describes content you do not want in the image. Maximum 500 characters. Chinese and English supported. Example: low resolution, low quality, deformed limbs, deformed fingers, oversaturated, waxy, no facial details, overly smooth, AI-like, chaotic composition, blurry text, distorted text.

 [​ ](#parameterssize) parameters.sizestring default"1280*1280" The resolution of the output image in `width*height` format. Default: `1280*1280`. The total number of pixels must be between `1280*1280` and `1440*1440`, and the aspect ratio must be between 1:4 and 4:1. Recommended resolutions: 1:1 (`1280*1280`), 3:4 (`1104*1472`), 4:3 (`1472*1104`), 9:16 (`960*1696`), 16:9 (`1696*960`).

 Example:1280*1280 [​ ](#parametersn) parameters.ninteger default4 The number of images to generate. Value: 1 to 4. Default: `4`. Note: billing is per image (Cost = Unit price x Number of images). Set to 1 for testing.

 Required range:1 <= x <= 4 [​ ](#parametersprompt-extend) parameters.prompt_extendboolean defaulttrue Whether to enable prompt rewriting. When enabled, this feature uses an LLM to optimize the positive prompt, significantly improving results for shorter prompts but adding a few seconds to processing time. Default: `true`.

 [​ ](#parameterswatermark) parameters.watermarkboolean defaultfalse Whether to add a watermark in the lower-right corner of the image with the text "AI-generated". Default: `false`.

 [​ ](#parametersseed) parameters.seedinteger Optional. Random number seed. Range: [0, 2147483647]. Using the same seed produces more consistent results, but identical results are not guaranteed. If not provided, the system uses a random number.

 Required range:0 <= x <= 2147483647 ### Response
200-application/json [​ ](#request-id) request_idstring Unique request identifier for troubleshooting.

 [​ ](#output) outputobject Show child attributes

 [​ ](#outputtask-id) output.task_idstring The task identifier. Use this to query the task status via `GET /tasks/{task_id}`. Valid for 24 hours.

 [​ ](#outputtask-status) output.task_statusenum<string> The initial task status. Typically `PENDING`.

 Available options:PENDING
