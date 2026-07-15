# HappyHorse -- Generate a video from an image

> **Source:** https://docs.qwencloud.com/api-reference/video-generation/happyhorse-image-to-video/create-task

POST/services/aigc/video-generation/video-synthesis cURL

 cURL - Image-to-video (first frame)

 Copy ```\ncurl --location 'https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis' \
 -H 'X-DashScope-Async: enable' \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY" \
 -H 'Content-Type: application/json' \
 -d '{
 "model": "happyhorse-1.1-i2v",
 "input": {
 "prompt": "A cat running on the grass",
 "media": [
 {
 "type": "first_frame",
 "url": "https://cdn.translate.alibaba.com/r/wanx-demo-1.png"
 }
 ]
 },
 "parameters": {
 "resolution": "720P",
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
``` Generate videos up to 15 seconds at 1080P from a first-frame image using the HappyHorse model. The output video aspect ratio automatically follows the input image. ### Authorizations
 [​ ](#authorization) Authorizationstring header required DashScope API Key. Get one from the [Qwen Cloud console](https://home.qwencloud.com/api-keys).

 ### Header Parameters
 [​ ](#x-dashscope-async) X-DashScope-Asyncenum<string> required Must be set to `enable` for asynchronous task submission.

 Available options:enable ### Body
application/json [​ ](#model) modelenum<string> required Model identifier. Supported values: `happyhorse-1.1-i2v`, `happyhorse-1.0-i2v`.

 Available options:happyhorse-1.1-i2v,happyhorse-1.0-i2v Example:happyhorse-1.1-i2v [​ ](#input) inputobject required Input content for video generation.

 Show child attributes

 [​ ](#inputmedia) input.mediaobject[] required Input media list. Must contain exactly one `first_frame` image. The aspect ratio of the generated video automatically follows the input image.

 Show child attributes

 [​ ](#inputmediatype) input.media.typeenum<string> required Media asset type. Fixed value: `first_frame`.

 Available options:first_frame [​ ](#inputmediaurl) input.media.urlstring required URL or Base64-encoded data of the first frame image.


**Public URL**: HTTP/HTTPS (e.g., `https://xxx/xxx.png`).


**Base64**: `data:{MIME_type};base64,{base64_data}` (supported MIME types: `image/jpeg`, `image/png`, `image/webp`).


Image constraints: Formats: JPEG, JPG, PNG, WEBP. Resolution: width and height at least 300 px. Aspect ratio: 1:2.5 to 2.5:1. File size: up to 20 MB.

 [​ ](#inputprompt) input.promptstring Text description of the video to generate. Optional -- if omitted, the model infers motion from the image alone. Supports any language. Maximum 5,000 non-Chinese characters or 2,500 Chinese characters (auto-truncated if exceeded).

 Example:A cat running on the grass [​ ](#parameters) parametersobject Video generation parameters.

 Show child attributes

 [​ ](#parametersresolution) parameters.resolutionenum<string> default"1080P" Video clarity tier. The model automatically scales to approximate the total pixel count for the selected resolution. The aspect ratio of the output video closely matches the input first frame.

 Available options:720P,1080P [​ ](#parametersduration) parameters.durationinteger default5 Video duration in seconds (integer, 3-15). Longer videos cost more -- billing is per second.

 Required range:3 <= x <= 15 [​ ](#parameterswatermark) parameters.watermarkboolean defaulttrue Whether to add a watermark to the generated video. `true` (default): adds a "HappyHorse" watermark in the lower-right corner. `false`: no watermark.

 [​ ](#parametersseed) parameters.seedinteger Seed for reproducible results. Same seed + same parameters produces similar (not identical) output.

 Required range:0 <= x <= 2147483647 ### Response
200-application/json [​ ](#request-id) request_idstring Unique request identifier for tracing and troubleshooting.

 Example:4909100c-7b5a-9f92-bfe5-xxxxxx [​ ](#output) outputobject Show child attributes

 [​ ](#outputtask-id) output.task_idstring Task ID for polling status. Use with `GET /tasks/{task_id}`.

 Example:0385dc79-5ff8-4d82-bcb6-xxxxxx [​ ](#outputtask-status) output.task_statusenum<string> Initial task status, typically `PENDING`.

 Available options:PENDING,RUNNING,SUCCEEDED,FAILED
