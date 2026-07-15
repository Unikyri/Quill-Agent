# HappyHorse -- Edit a video

> **Source:** https://docs.qwencloud.com/api-reference/video-generation/happyhorse-video-editing/create-task

POST/services/aigc/video-generation/video-synthesis cURL

 cURL - Video editing (instruction + reference image)

 Copy ```\ncurl --location 'https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis' \
 -H 'X-DashScope-Async: enable' \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY" \
 -H 'Content-Type: application/json' \
 -d '{
 "model": "happyhorse-1.0-video-edit",
 "input": {
 "prompt": "Make the horse-headed humanoid character in the video wear the striped sweater from the image",
 "media": [
 {
 "type": "video",
 "url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20260409/dozxak/Wan_Video_Edit_33_1.mp4"
 },
 {
 "type": "reference_image",
 "url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20260415/hynnff/wan-video-edit-clothes.webp"
 }
 ]
 },
 "parameters": {
 "resolution": "720P"
 }
}'
``` 200 400 Copy ```\n{
 "request_id": "4909100c-7b5a-9f92-bfe5-xxxxxx",
 "output": {
 "task_id": "0385dc79-5ff8-4d82-bcb6-xxxxxx",
 "task_status": "PENDING"
 }
}
``` Edit videos using the HappyHorse model with text instructions and optional reference images. Supports style transfer, local replacement, and other editing tasks. Input videos can be 3-60 seconds; output is capped at 15 seconds. ### Authorizations
 [​ ](#authorization) Authorizationstring header required DashScope API Key. Get one from the [Qwen Cloud console](https://home.qwencloud.com/api-keys).

 ### Header Parameters
 [​ ](#x-dashscope-async) X-DashScope-Asyncenum<string> required Must be set to `enable` for asynchronous task submission.

 Available options:enable ### Body
application/json [​ ](#model) modelenum<string> required Model identifier. Fixed value: `happyhorse-1.0-video-edit`.

 Available options:happyhorse-1.0-video-edit Example:happyhorse-1.0-video-edit [​ ](#input) inputobject required Input content including the video to edit, optional reference images, and the editing prompt.

 Show child attributes

 [​ ](#inputprompt) input.promptstring required Text prompt describing the intended edit (style transfer, local replacement, etc.). Supports any language. Maximum 5,000 non-Chinese characters or 2,500 Chinese characters (auto-truncated if exceeded).

 Example:Make the horse-headed humanoid character in the video wear the striped sweater from the image [​ ](#inputmedia) input.mediaobject[] required List of media assets. Must contain exactly 1 `video` element and optionally 0-5 `reference_image` elements.

 Show child attributes

 [​ ](#inputmediatype) input.media.typeenum<string> required Media asset type: `video` (the video to edit, exactly 1 required) or `reference_image` (optional, 0-5).

 Available options:video,reference_image [​ ](#inputmediaurl) input.media.urlstring required URL or Base64-encoded data of the media asset.


**Video** (`type=video`): Publicly accessible URL (HTTP/HTTPS). MP4/MOV (H.264 recommended). Duration: 3-60 seconds. Resolution: longer side <= 4,096 px, shorter side >= 360 px. Aspect ratio: 1:2.5 to 2.5:1. File size: up to 100 MB. Frame rate: > 8 fps. Output duration is capped at 15 seconds (longer inputs are auto-trimmed).


**Image** (`type=reference_image`): URL or Base64-encoded data. JPEG/JPG/PNG/WEBP. Width and height >= 300 px. Aspect ratio: 1:2.5 to 2.5:1. File size: up to 20 MB.


**Base64 format**: `data:{MIME_type};base64,{base64_data}` (supported MIME types: `image/jpeg`, `image/png`, `image/webp`).

 [​ ](#parameters) parametersobject Video editing parameters.

 Show child attributes

 [​ ](#parametersresolution) parameters.resolutionenum<string> default"1080P" Resolution of the generated video.

 Available options:720P,1080P [​ ](#parameterswatermark) parameters.watermarkboolean defaulttrue Whether to add a watermark to the generated video. `true` (default): adds a "Happy Horse" watermark in the bottom-right corner. `false`: no watermark.

 [​ ](#parametersaudio-setting) parameters.audio_settingenum<string> default"auto" Audio control for the output video.

 Available options:auto,origin [​ ](#parametersseed) parameters.seedinteger Seed for reproducible results. Same seed + same parameters produces similar (not identical) output.

 Required range:0 <= x <= 2147483647 ### Response
200-application/json [​ ](#request-id) request_idstring Unique request identifier for tracing and troubleshooting.

 Example:4909100c-7b5a-9f92-bfe5-xxxxxx [​ ](#output) outputobject Show child attributes

 [​ ](#outputtask-id) output.task_idstring Task ID for polling status. Use with `GET /tasks/{task_id}`.

 Example:0385dc79-5ff8-4d82-bcb6-xxxxxx [​ ](#outputtask-status) output.task_statusenum<string> Initial task status, typically `PENDING`.

 Available options:PENDING,RUNNING,SUCCEEDED,FAILED
