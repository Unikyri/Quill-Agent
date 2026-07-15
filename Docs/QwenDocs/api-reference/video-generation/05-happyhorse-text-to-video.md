# HappyHorse -- Generate a video from text

> **Source:** https://docs.qwencloud.com/api-reference/video-generation/happyhorse-text-to-video/create-task

POST/services/aigc/video-generation/video-synthesis cURL

 cURL - Text-to-video

 Copy ```\ncurl --location 'https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis' \
 -H 'X-DashScope-Async: enable' \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY" \
 -H 'Content-Type: application/json' \
 -d '{
 "model": "happyhorse-1.1-t2v",
 "input": {
 "prompt": "A miniature city built from cardboard and bottle caps comes to life at night. A cardboard train slowly passes through, with small lights dotting the scene and illuminating the way ahead."
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
``` Generate videos up to 15 seconds at 1080P from text prompts using the HappyHorse model. ### Authorizations
 [​ ](#authorization) Authorizationstring header required DashScope API Key. Get one from the [Qwen Cloud console](https://home.qwencloud.com/api-keys).

 ### Header Parameters
 [​ ](#x-dashscope-async) X-DashScope-Asyncenum<string> required Must be set to `enable` for asynchronous task submission.

 Available options:enable ### Body
application/json [​ ](#model) modelenum<string> required Model identifier. Supported values: `happyhorse-1.1-t2v`, `happyhorse-1.0-t2v`.

 Available options:happyhorse-1.1-t2v,happyhorse-1.0-t2v Example:happyhorse-1.1-t2v [​ ](#input) inputobject required Input content for video generation.

 Show child attributes

 [​ ](#inputprompt) input.promptstring required Describe the video you want. Supports any language. Maximum 5,000 non-Chinese characters or 2,500 Chinese characters (auto-truncated if exceeded).

 Example:A miniature city built from cardboard and bottle caps comes to life at night. [​ ](#parameters) parametersobject Video generation parameters.

 Show child attributes

 [​ ](#parametersresolution) parameters.resolutionenum<string> default"1080P" Video clarity tier. Higher resolution costs more.

 Available options:720P,1080P [​ ](#parametersratio) parameters.ratioenum<string> default"16:9" Aspect ratio of the generated video.

 Available options:16:9,9:16,1:1,4:3,3:4,4:5,5:4,9:21,21:9 [​ ](#parametersduration) parameters.durationinteger default5 Video duration in seconds (integer, 3-15). Longer videos cost more -- billing is per second.

 Required range:3 <= x <= 15 [​ ](#parameterswatermark) parameters.watermarkboolean defaulttrue Whether to add a watermark to the generated video. `true` (default): adds a "HappyHorse" watermark in the lower-right corner. `false`: no watermark.

 [​ ](#parametersseed) parameters.seedinteger Seed for reproducible results. Same seed + same parameters produces similar (not identical) output.

 Required range:0 <= x <= 2147483647 ### Response
200-application/json [​ ](#request-id) request_idstring Unique request identifier for tracing and troubleshooting.

 Example:4909100c-7b5a-9f92-bfe5-xxxxxx [​ ](#output) outputobject Show child attributes

 [​ ](#outputtask-id) output.task_idstring Task ID for polling status. Use with `GET /tasks/{task_id}`.

 Example:0385dc79-5ff8-4d82-bcb6-xxxxxx [​ ](#outputtask-status) output.task_statusenum<string> Initial task status, typically `PENDING`.

 Available options:PENDING,RUNNING,SUCCEEDED,FAILED
