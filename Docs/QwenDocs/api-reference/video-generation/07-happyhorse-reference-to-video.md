# HappyHorse -- Generate a video from reference images

> **Source:** https://docs.qwencloud.com/api-reference/video-generation/happyhorse-reference-to-video/create-task

POST/services/aigc/video-generation/video-synthesis cURL

 cURL - Reference-to-video (multi-image)

 Copy ```\ncurl --location 'https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis' \
 -H 'X-DashScope-Async: enable' \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY" \
 -H 'Content-Type: application/json' \
 -d '{
 "model": "happyhorse-1.1-r2v",
 "input": {
 "prompt": "The woman in the red qipao from [Image 1] is first framed in a profile medium shot, highlighting the tailored cut and S-curve silhouette of the dress. The camera then cuts to a low-angle upward shot, capturing the moment she gracefully raises her hand to unfold the folding fan from [Image 2], while the tassel earrings from [Image 3] sway delicately with the turn of her head.",
 "media": [
 {
 "type": "reference_image",
 "url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20260424/mvzfud/hh-v2v-girl.jpg"
 },
 {
 "type": "reference_image",
 "url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20260424/fvuihk/hh-v2v2-folding-fan.jpg"
 },
 {
 "type": "reference_image",
 "url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20260424/imerii/hh-v2v-earrings.jpg"
 }
 ]
 },
 "parameters": {
 "resolution": "720P",
 "ratio": "16:9",
 "duration": 5
 }
}'
``` 200 400 Copy ```\n{
 "request_id": "4909100c-7b5a-9f92-bfe5-xxxxxx",
 "output": {
 "task_id": "0385dc79-5ff8-4d82-bcb6-xxxxxx",
 "task_status": "PENDING"
 }
}
``` Generate videos up to 15 seconds at 1080P from 1-9 reference images using the HappyHorse model. Provide multiple reference images and a text prompt to generate a video that combines subjects from the images into a scene. Use `[Image 1]`, `[Image 2]`, etc. in the prompt to reference images by their order in the `media` array. ### Authorizations
 [​ ](#authorization) Authorizationstring header required DashScope API Key. Get one from the [Qwen Cloud console](https://home.qwencloud.com/api-keys).

 ### Header Parameters
 [​ ](#x-dashscope-async) X-DashScope-Asyncenum<string> required Must be set to `enable` for asynchronous task submission.

 Available options:enable ### Body
application/json [​ ](#model) modelenum<string> required Model identifier. Supported values: `happyhorse-1.1-r2v`, `happyhorse-1.0-r2v`.

 Available options:happyhorse-1.1-r2v,happyhorse-1.0-r2v Example:happyhorse-1.1-r2v [​ ](#input) inputobject required Input content including reference images and text prompt.

 Show child attributes

 [​ ](#inputprompt) input.promptstring required Text prompt describing the desired video. Use `[Image 1]`, `[Image 2]`, etc. to reference images in the `media` array (order matches array order). You must specify the object from the reference image, for example, "the woman in a red qipao from [Image 1]". Supports any language. Maximum 5,000 non-Chinese characters or 2,500 Chinese characters (auto-truncated if exceeded).

 Example:The woman in the red qipao from [Image 1] stands elegantly. [​ ](#inputmedia) input.mediaobject[] required List of reference images (1-9). The array order defines the image reference order in the `prompt`: the first `reference_image` maps to `[Image 1]`, the second to `[Image 2]`, and so on.

 Required range:items: 1–9 Show child attributes

 [​ ](#inputmediatype) input.media.typeenum<string> required Media asset type. Fixed value: `reference_image`.

 Available options:reference_image [​ ](#inputmediaurl) input.media.urlstring required URL or Base64-encoded data of the reference image.


**Public URL**: HTTP/HTTPS (e.g., `https://xxx/xxx.jpg`).


**Base64**: `data:{MIME_type};base64,{base64_data}` (supported MIME types: `image/jpeg`, `image/png`, `image/webp`).


Image constraints: Formats: JPEG, JPG, PNG, WEBP. Shortest side at least 400 px (720P+ recommended). File size: up to 20 MB.

 [​ ](#parameters) parametersobject Video generation parameters.

 Show child attributes

 [​ ](#parametersresolution) parameters.resolutionenum<string> default"1080P" Video clarity tier. Higher resolution costs more.

 Available options:720P,1080P [​ ](#parametersratio) parameters.ratioenum<string> default"16:9" Aspect ratio of the generated video.

 Available options:16:9,9:16,1:1,4:3,3:4,4:5,5:4,9:21,21:9 [​ ](#parametersduration) parameters.durationinteger default5 Video duration in seconds (integer, 3-15). Longer videos cost more -- billing is per second.

 Required range:3 <= x <= 15 [​ ](#parameterswatermark) parameters.watermarkboolean defaulttrue Whether to add a watermark to the generated video. `true` (default): adds a "Happy Horse" watermark in the bottom-right corner. `false`: no watermark.

 [​ ](#parametersseed) parameters.seedinteger Seed for reproducible results. Same seed + same parameters produces similar (not identical) output.

 Required range:0 <= x <= 2147483647 ### Response
200-application/json [​ ](#request-id) request_idstring Unique request identifier for tracing and troubleshooting.

 Example:4909100c-7b5a-9f92-bfe5-xxxxxx [​ ](#output) outputobject Show child attributes

 [​ ](#outputtask-id) output.task_idstring Task ID for polling status. Use with `GET /tasks/{task_id}`.

 Example:0385dc79-5ff8-4d82-bcb6-xxxxxx [​ ](#outputtask-status) output.task_statusenum<string> Initial task status, typically `PENDING`.

 Available options:PENDING,RUNNING,SUCCEEDED,FAILED
