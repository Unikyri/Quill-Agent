# Wan 2.7 -- Generate a video from image

> **Source:** https://docs.qwencloud.com/api-reference/video-generation/wan27-image-to-video/create-task

POST/services/aigc/video-generation/video-synthesis cURL

 cURL - First frame + audio

 Copy ```\ncurl --location 'https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis' \
 -H 'X-DashScope-Async: enable' \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY" \
 -H 'Content-Type: application/json' \
 -d '{
 "model": "wan2.7-i2v",
 "input": {
 "prompt": "A scene of urban fantasy art. A dynamic graffiti art character. A boy made of spray paint comes to life on a concrete wall. He sings an English rap song at high speed while striking a classic, energetic rapper pose. The scene is set under an urban railway bridge at night. The light comes from a single street lamp, creating a cinematic atmosphere full of high energy and amazing detail. The audio of the video consists entirely of the rap, with no other dialogue or noise.",
 "media": [
 {
 "type": "first_frame",
 "url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250925/wpimhv/rap.png"
 },
 {
 "type": "driving_audio",
 "url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250925/ozwpvi/rap.mp3"
 }
 ]
 },
 "parameters": {
 "resolution": "720P",
 "duration": 10,
 "prompt_extend": true,
 "watermark": true
 }
}'
``` 200 400 Copy ```\n{
 "request_id": "4909100c-7b5a-9f92-bfe5-xxxxxx",
 "output": {
 "task_id": "0385dc79-5ff8-4d82-bcb6-xxxxxx",
 "task_status": "PENDING"
 }
}
``` Generate videos up to 15 seconds at 1080P from images, audio, and video clips, with optional audio sync and first-last frame control.
## [​ ](#changes-from-wan2-6) Changes from wan2.6


- **Unified API**: First-frame, first-last-frame, and video continuation share one endpoint via the `media` array -- no separate APIs.

- **Audio-video sync**: Provide a `driving_audio` file for lip-syncing. If omitted, the model auto-generates matching sound effects.

- **Resolution control**: Set `resolution` (720P/1080P) instead of exact pixel `size`.

- **Longer prompts**: Up to 5,000 characters (was 800).

- **Negative prompt moved**: Now under `input.negative_prompt` instead of `parameters.negative_prompt`.

- **Watermark off by default**: `watermark` defaults to `false` (was `true`).


 ### Authorizations
 [​ ](#authorization) Authorizationstring header required DashScope API Key. Get one from the [Qwen Cloud console](https://home.qwencloud.com/api-keys).

 ### Header Parameters
 [​ ](#x-dashscope-async) X-DashScope-Asyncenum<string> required Must be set to `enable` for asynchronous task submission.

 Available options:enable ### Body
application/json [​ ](#model) modelenum<string> required Model identifier. Supported values: `wan2.7-i2v`, `wan2.7-i2v-2026-04-25`.

 Available options:wan2.7-i2v,wan2.7-i2v-2026-04-25 Example:wan2.7-i2v [​ ](#input) inputobject required The basic input information, such as the prompt.

 Show child attributes

 [​ ](#inputmedia) input.mediaobject[] required A list of media assets. Specifies reference materials (images, audio, and video) for video generation. Each element in the array is a media object that contains the `type` and `url`. Each `type` can appear at most once.


Only the following asset combinations are supported. Invalid combinations result in an error.


- `first_frame`

- `first_frame` + `driving_audio`

- `first_frame` + `last_frame`

- `first_frame` + `last_frame` + `driving_audio`

- `first_clip`

- `first_clip` + `last_frame`


 Show child attributes

 [​ ](#inputmediatype) input.media.typeenum<string> required The type of media asset. Valid values:


- `first_frame`: The URL of the first or last frame, or Base64-encoded data. Video generation from the first frame and video generation from the first and last frames are supported. JPEG/JPG/PNG/BMP/WEBP, 240-8000 px per side, 1:8 to 8:1 ratio, max 20 MB.

- `last_frame`: Last frame image. Same limits as `first_frame`.

- `driving_audio`: The URL of the audio file. Pass audio: The model uses the audio as a driving source to generate the video, such as for lip-syncing and action timing. Do not pass audio: The model automatically generates matching background music or sound effects based on the video content. WAV/MP3, 2-30 s, max 15 MB.

- `first_clip`: The URL of the video file. The model generates a continuation based on the video content. The `duration` parameter controls the maximum duration of the continuation. MP4/MOV, 2-10 s, 240-4096 px per side, 1:8 to 8:1 ratio, max 100 MB.


 Available options:first_frame,last_frame,driving_audio,first_clip [​ ](#inputmediaurl) input.media.urlstring required The URL of the media asset. Assets include images, audio, and video.

 [​ ](#inputprompt) input.promptstring The text prompt. Describes the elements and visual features of the video you want. Supports Chinese and English, up to 5,000 characters (auto-truncated if exceeded).

 Example:A kitten runs on the grass. [​ ](#inputnegative-prompt) input.negative_promptstring The negative prompt. Describes content you do not want in the video (e.g., `low quality, blurry, extra fingers`). Supports Chinese and English. The prompt can be up to 500 characters long (auto-truncated if exceeded).

 Example:low resolution, error, worst quality, low quality, deformed, extra fingers, bad proportions [​ ](#parameters) parametersobject Video processing parameters, such as resolution, duration, prompt rewriting, and watermarks.

 Show child attributes

 [​ ](#parametersresolution) parameters.resolutionenum<string> default"1080P" The resolution tier for the generated video. Controls the video definition (total pixels). The resolution directly affects the cost. Before you make a call, confirm the resolution you need. The model automatically scales the video to a total pixel count close to the selected resolution tier. The video&#x27;s aspect ratio should be as consistent as possible with the input material (first frame or first video segment).

 Available options:720P,1080P [​ ](#parametersduration) parameters.durationinteger default5 Video duration in seconds (integer, 2-15). Longer videos cost more -- billing is per second. For video continuation (`first_clip`), this is the total output duration including the input clip. For example, if `duration=15` and the input is 3 s, the model generates 12 s of new content.

 Required range:2 <= x <= 15 [​ ](#parametersprompt-extend) parameters.prompt_extendboolean defaulttrue Rewrite the prompt using a Large Language Model before generation. Improves results for short or vague prompts but adds latency. Set to `false` to use your prompt verbatim.

 [​ ](#parameterswatermark) parameters.watermarkboolean defaultfalse Add an "AI Generated" watermark in the lower-right corner.

 [​ ](#parametersseed) parameters.seedinteger Seed for reproducible results. Same seed + same parameters produces similar (not identical) output.

 Required range:0 <= x <= 2147483647 ### Response
200-application/json [​ ](#request-id) request_idstring Unique request identifier for tracing and troubleshooting.

 Example:4909100c-7b5a-9f92-bfe5-xxxxxx [​ ](#output) outputobject Show child attributes

 [​ ](#outputtask-id) output.task_idstring Task ID for polling status. Use with `GET /tasks/{task_id}`.

 Example:0385dc79-5ff8-4d82-bcb6-xxxxxx [​ ](#outputtask-status) output.task_statusenum<string> Initial task status, typically `PENDING`.

 Available options:PENDING,RUNNING,SUCCEEDED,FAILED
