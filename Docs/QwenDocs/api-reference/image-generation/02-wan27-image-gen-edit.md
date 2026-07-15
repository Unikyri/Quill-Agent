# Wan 2.7 — Generate or edit an image

> **Source:** https://docs.qwencloud.com/api-reference/image-generation/wan27-image-gen-edit/create-task

POST/services/aigc/image-generation/generation cURL

 cURL (text-to-image, async)

 Copy curl --location &#x27;https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/image-generation/generation&#x27; \
--header &#x27;Content-Type: application/json&#x27; \
--header "Authorization: Bearer $DASHSCOPE_API_KEY" \
--header &#x27;X-DashScope-Async: enable&#x27; \
--data &#x27;{
 "model": "wan2.7-image-pro",
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
 "n": 1,
 "size": "2K",
 "watermark": false,
 "thinking_mode": true
 }
}&#x27; 200 400 Copy {
 "request_id": "ccf4b2f4-bf30-9e13-9461-3a28c6a7bxxx",
 "output": {
 "task_id": "8811b4a4-00ac-4aa2-a2fd-017d3b90cxxx",
 "task_status": "PENDING"
 }
} [Get an API key](/api-reference/preparation/api-key) and [set it as an environment variable](/api-reference/preparation/export-api-key-env). If using the SDK, [install it](/api-reference/preparation/install-sdk). 
Image generation tasks take 1 to 2 minutes. Use the async API to avoid request timeouts by splitting the process into two steps:

- **Create a task** (this endpoint) and receive a `task_id`.

- **[Query the result](/api-reference/image-generation/wan27-image-gen-edit/query-result)** by polling with the `task_id`.


The request body uses the same `messages` format and parameters as the [synchronous endpoint](/api-reference/image-generation/wan27-image-gen-edit/synchronous), but requires the `X-DashScope-Async: enable` header. ### Authorizations
 [​ ](#authorization) Authorizationstring header required DashScope API Key. Create one in the [Qwen Cloud console](https://home.qwencloud.com/api-keys).

 ### Header Parameters
 [​ ](#x-dashscope-async) X-DashScope-Asyncenum<string> required Asynchronous processing configuration. **Must be set to `enable`**.

 Available options:enable ### Body
application/json [​ ](#model) modelenum<string> required The model name. Valid values: wan2.7-image-pro, wan2.7-image.

 Available options:wan2.7-image-pro,wan2.7-image Example:wan2.7-image-pro [​ ](#input) inputobject required Input data containing the messages array.

 Show child attributes

 [​ ](#inputmessages) input.messagesobject[] required An array of request content. Currently, only single-turn conversations are supported. This means you can pass only one set of role and content parameters. Multi-turn conversations are not supported.

 Required range:items: 1–1 Show child attributes

 [​ ](#inputmessagesrole) input.messages.roleenum<string> required Message role. Must be `user`.

 Available options:user [​ ](#inputmessagescontent) input.messages.contentobject[] required Message content array. Must contain exactly one `text` object and 0 to 9 `image` objects.


When using multiple images, include multiple `image` objects in the array. Image order is determined by array position.

 Show child attributes

 [​ ](#inputmessagescontenttext) input.messages.content.textstring The user-entered prompt. Supports Chinese and English. The length cannot exceed 5,000 characters. Each Chinese character, letter, number, or symbol counts as one character. Any excess is automatically truncated. The `content` array must contain exactly one `text` object.

 Example:Spray the graffiti from image 2 onto the car in image 1 Required range:length <= 5000 [​ ](#inputmessagescontentimage) input.messages.content.imagestring Input image as a public URL (HTTP/HTTPS) or Base64-encoded string (`data:{mime_type};base64,{data}`).


**Image constraints:**


- Formats: JPEG, JPG, PNG (alpha channel not supported), BMP, WEBP.

- Resolution: Width and height each between 240 and 8,000 pixels. Aspect ratio [1:8, 8:1].

- File size: Max 20 MB.

- Quantity: 0 to 9 images per request.


 Example:https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20251229/pjeqdf/car.webp [​ ](#parameters) parametersobject Image processing parameters.

 Show child attributes

 [​ ](#parameterssize) parameters.sizestring Output image resolution. Two specification methods are available; they cannot be used together.


**wan2.7-image-pro:**


- Method 1 (recommended): `1K`, `2K` (default), or `4K`.

Scope: Text-to-image (no image input, not generating an image set) supports 1K, 2K, and 4K. Other scenarios support 1K and 2K only.

- Total pixels: 1K = 1024×1024, 2K = 2048×2048, 4K = 4096×4096.

- Aspect ratio: With image input, output matches the aspect ratio of the last input image scaled to the selected resolution. Without image input, output is square.


- Method 2: Specify `width*height` in pixels, aspect ratio [1:8, 8:1].

Text-to-image: total pixels in [768×768, 4096×4096].

- Other scenarios: total pixels in [768×768, 2048×2048].


**wan2.7-image:**


- Method 1 (recommended): `1K` or `2K` (default). 4K is not supported.

- Method 2: Specify `width*height` in pixels. All scenarios: total pixels in [768×768, 2048×2048], aspect ratio [1:8, 8:1].


The pixel value of the output image may differ slightly from the specified value.

 Example:2K [​ ](#parametersn) parameters.ninteger default1 Number of images to generate.


**Note:** The value of `n` directly affects the cost. Cost = Unit Price × Number of successfully generated images.


- When image set mode is disabled (`enable_sequential=false`): This value represents the number of images to generate. Range: 1–4. Default: 1.

- When image set mode is enabled (`enable_sequential=true`): This value represents the maximum number of images to generate. Range: 1–12. Default: 12. The actual number is determined by the model and will not exceed `n`.


 Required range:1 <= x <= 12 [​ ](#parametersenable-sequential) parameters.enable_sequentialboolean defaultfalse Controls the image generation mode.


- `false`: Default value.

- `true`: Enables image set output mode.


 [​ ](#parametersthinking-mode) parameters.thinking_modeboolean defaulttrue Specifies whether to enable thinking mode. The default is `true` (enabled). This parameter is effective only when image set mode is disabled and there is no image input. When enabled, the model enhances its inference capabilities to improve image quality, but this increases generation time.

 [​ ](#parametersbbox-list) parameters.bbox_listinteger[][][] The selected area for interactive editing.


- Correspondence: The length of the list must match the number of input images. If an image does not require editing, pass an empty list `[]` at the corresponding position.

- Coordinate format: `[x1, y1, x2, y2]` (top-left x, top-left y, bottom-right x, bottom-right y). Use absolute pixel coordinates of the original image. The top-left coordinate is (0, 0).

- Condition: A single image supports a maximum of 2 bounding boxes.


 [​ ](#parameterscolor-palette) parameters.color_paletteobject[] A custom color theme. An array of objects containing color (`hex`) and proportion (`ratio`). It must include 3 to 10 colors. We recommend setting it to 8. Available only when image set mode is disabled (`enable_sequential=false`).

 Example: Copy [
 {
 "hex": "#C2D1E6",
 "ratio": "60.00%"
 },
 {
 "hex": "#636574",
 "ratio": "25.00%"
 },
 {
 "hex": "#CBD4E4",
 "ratio": "15.00%"
 }
]
 Required range:items: 3–10 Show child attributes

 [​ ](#parameterscolor-palettehex) parameters.color_palette.hexstring required The color value in hexadecimal (HEX) format. Example: `#C2D1E6`.

 [​ ](#parameterscolor-paletteratio) parameters.color_palette.ratiostring required The percentage of the color. It must be accurate to two decimal places (for example, `"25.00%"`). The sum of all `ratio` values must be 100.00%.

 [​ ](#parameterswatermark) parameters.watermarkboolean defaultfalse Adds a watermark label in the bottom-right corner of the image with fixed text "AI Generated".

 [​ ](#parametersseed) parameters.seedinteger Random number seed. Valid range: [0, 2147483647]. Using the same seed yields similar outputs. If omitted, the algorithm uses a random seed. Note: Image generation is probabilistic. Even with the same seed, results may vary.

 Required range:0 <= x <= 2147483647 ### Response
200-application/json [​ ](#request-id) request_idstring default"ccf4b2f4-bf30-9e13-9461-3a28c6a7bxxx" Unique request identifier.

 Example:ccf4b2f4-bf30-9e13-9461-3a28c6a7bxxx [​ ](#output) outputobject Show child attributes

 [​ ](#outputtask-id) output.task_idstring Task identifier. Use this to poll the result endpoint. Valid for 24 hours.

 Example:8811b4a4-00ac-4aa2-a2fd-017d3b90cxxx [​ ](#outputtask-status) output.task_statusstring Initial task status. Always `PENDING` immediately after creation.

 Example:PENDING
