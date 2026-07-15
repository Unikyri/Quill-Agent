# Wan 2.6 — Generate or edit an image

> **Source:** https://docs.qwencloud.com/api-reference/image-generation/wan26-image-gen-edit/create-task

POST/services/aigc/image-generation/generation cURL

 cURL (image editing, async)

 Copy ```\ncurl --location 'https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/image-generation/generation' \
--header 'Content-Type: application/json' \
--header "Authorization: Bearer $DASHSCOPE_API_KEY" \
--header 'X-DashScope-Async: enable' \
--data '{
 "model": "wan2.6-image",
 "input": {
 "messages": [
 {
 "role": "user",
 "content": [
 {
 "text": "Generate a tomato and egg stir-fry based on the style of image 1 and the background of image 2"
 },
 {
 "image": "https://cdn.wanx.aliyuncs.com/tmp/pressure/umbrella1.png"
 },
 {
 "image": "https://img.alicdn.com/imgextra/i3/O1CN01SfG4J41UYn9WNt4X1_!!6000000002530-49-tps-1696-960.webp"
 }
 ]
 }
 ]
 },
 "parameters": {
 "prompt_extend": true,
 "watermark": false,
 "n": 1,
 "enable_interleave": false,
 "size": "1K"
 }
}'
``` 200 400 Copy ```\n{
 "output": {
 "task_status": "PENDING",
 "task_id": "0385dc79-5ff8-4d82-bcb6-xxxxxx"
 },
 "request_id": "4909100c-7b5a-9f92-bfe5-xxxxxx"
}
``` [Get an API key](/api-reference/preparation/api-key) and [set it as an environment variable](/api-reference/preparation/export-api-key-env). If using the SDK, [install it](/api-reference/preparation/install-sdk). 
Image generation tasks take 1 to 2 minutes. Use the async API to avoid request timeouts by splitting the process into two steps:

- **Create a task** (this endpoint) and receive a `task_id`.

- **[Query the result](/api-reference/image-generation/wan26-image-gen-edit/query-result)** by polling with the `task_id`.


The request body uses the same `messages` format and parameters as the [synchronous endpoint](/api-reference/image-generation/wan26-image-gen-edit/synchronous), but requires the `X-DashScope-Async: enable` header. ### Authorizations
 [​ ](#authorization) Authorizationstring header required DashScope API Key. Create one in the [Qwen Cloud console](https://home.qwencloud.com/api-keys).

 ### Header Parameters
 [​ ](#x-dashscope-async) X-DashScope-Asyncenum<string> required Asynchronous processing configuration. **Must be set to `enable`**.

 Available options:enable ### Body
application/json [​ ](#model) modelenum<string> required Model name. Set to `wan2.6-image`.

 Available options:wan2.6-image Example:wan2.6-image [​ ](#input) inputobject required Input data containing the messages array.

 Show child attributes

 [​ ](#inputmessages) input.messagesobject[] required Array of request content. Only single-turn conversations are supported. Provide one message with `role: user`.

 Required range:items: 1–1 Show child attributes

 [​ ](#inputmessagesrole) input.messages.roleenum<string> required Message role. Must be `user`.

 Available options:user [​ ](#inputmessagescontent) input.messages.contentobject[] required Message content array. Must contain exactly one `text` object. Image objects depend on the mode:


- Image editing (`enable_interleave=false`): **1 to 4** image objects required.

- Interleaved text-image (`enable_interleave=true`): **0 to 1** image objects.


When using multiple images, include multiple `image` objects in the array. Image order is determined by array position.

 Show child attributes

 [​ ](#inputmessagescontenttext) input.messages.content.textstring Positive prompt describing the desired image content, style, and composition. Supports Chinese and English. Maximum 2,000 characters (each Chinese character, letter, digit, or symbol counts as one character). Excess is auto-truncated. The `content` array must contain exactly one `text` object.

 Example:Generate a tomato and egg stir-fry based on the style of image 1 and the background of image 2 Required range:length <= 2000 [​ ](#inputmessagescontentimage) input.messages.content.imagestring Input image as a public URL (HTTP/HTTPS) or Base64-encoded string (`data:{mime_type};base64,{data}`).


**Image constraints:**


- Formats: JPEG, JPG, PNG (alpha channel not supported), BMP, WEBP.

- Resolution: Width and height each between 240 and 8,000 pixels.

- File size: Max 10 MB.


**Image quantity limits:**


- When `enable_interleave=false` (image editing): must input 1 to 4 images.

- When `enable_interleave=true` (interleaved text-image): can input 0 to 1 images.


 Example:https://cdn.wanx.aliyuncs.com/tmp/pressure/umbrella1.png [​ ](#parameters) parametersobject Image processing parameters.

 Show child attributes

 [​ ](#parametersnegative-prompt) parameters.negative_promptstring Negative prompt describing content you do NOT want in the image. Supports Chinese and English. Maximum 500 characters. Excess is auto-truncated.


Example: `Low resolution, low quality, deformed limbs, deformed fingers, oversaturated colors, wax-like appearance, no facial details, overly smooth skin, AI-looking artifacts, chaotic composition, blurry or distorted text.`

 Required range:length <= 500 [​ ](#parameterssize) parameters.sizestring Output image resolution. Supports two methods: referencing input image proportions or directly specifying dimensions.


**Image editing mode** (`enable_interleave=false`):


- Method 1 (recommended): `1K` (default) or `2K`. Output total pixels close to 1280*1280 or 2048*2048, maintaining the aspect ratio of the last input image.

- Method 2: Specify `width*height` in pixels. Total pixels must be between [768*768, 2048*2048], aspect ratio [1:4, 4:1]. Actual values are multiples of 16.


**Interleaved text-image mode** (`enable_interleave=true`):


- Method 1 (default): References input image proportions. If total pixels <= 1280*1280, output matches input. If > 1280*1280, output scales to ~1280*1280.

- Method 2: Specify `width*height`. Total pixels must be between [768*768, 1280*1280], aspect ratio [1:4, 4:1].


**Recommended resolutions:** 1280*1280 (1:1), 800*1200 (2:3), 1200*800 (3:2), 960*1280 (3:4), 1280*960 (4:3), 720*1280 (9:16), 1280*720 (16:9), 1344*576 (21:9).

 Example:1K [​ ](#parametersenable-interleave) parameters.enable_interleaveboolean defaultfalse Controls image generation mode:


- `false` (default): Image editing mode. Supports multi-image input (1-4 images) and subject consistency generation. Can generate 1 to 4 result images.

- `true`: Interleaved text-image output mode. Supports 0-1 input images. Generates mixed content containing both text and images. **Requires** `stream=true` and `X-DashScope-Sse: enable` header.


 [​ ](#parametersn) parameters.ninteger default4 Number of images to generate. Behavior depends on mode:


- Image editing (`enable_interleave=false`): Range 1-4. Default: 4.

- Interleaved text-image (`enable_interleave=true`): Must be 1. Use `max_images` to control image count instead.


**Note:** `n` directly affects cost. Cost = unit price x number of successfully generated images.

 Required range:1 <= x <= 4 [​ ](#parametersmax-images) parameters.max_imagesinteger default5 Only effective in interleaved text-image mode (`enable_interleave=true`). Specifies the maximum number of images the model can generate in a single response. Range: 1-5. Default: 5. The actual number of generated images is determined by model inference and may be less than this value.


**Note:** `max_images` affects cost. Cost = unit price x number of successfully generated images.

 Required range:1 <= x <= 5 [​ ](#parametersprompt-extend) parameters.prompt_extendboolean defaulttrue Only effective in image editing mode (`enable_interleave=false`). Enables intelligent prompt rewriting that optimizes and refines the positive prompt. The negative prompt is not affected.

 [​ ](#parametersstream) parameters.streamboolean defaultfalse Controls whether the response uses streaming output. In interleaved text-image mode (`enable_interleave=true`), you **must** set this to `true`.

 [​ ](#parameterswatermark) parameters.watermarkboolean defaultfalse Adds a watermark label in the bottom-right corner of the image with fixed text "AI Generated".

 [​ ](#parametersseed) parameters.seedinteger Random number seed. Range: [0, 2147483647]. Same seed produces more consistent (but not identical) results. If omitted, a random seed is used.

 Required range:0 <= x <= 2147483647 ### Response
200-application/json [​ ](#request-id) request_idstring default"4909100c-7b5a-9f92-bfe5-xxxxxx" Unique request identifier for troubleshooting.

 [​ ](#output) outputobject Show child attributes

 [​ ](#outputtask-id) output.task_idstring The task identifier. Use this to query the task status via `GET /tasks/{task_id}`. Valid for 24 hours.

 [​ ](#outputtask-status) output.task_statusenum<string> The initial task status. Typically `PENDING`.

 Available options:PENDING
