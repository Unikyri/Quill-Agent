# Wan 2.7 -- Generate a video from text

> **Source:** https://docs.qwencloud.com/api-reference/video-generation/wan27-text-to-video/create-task

POST/services/aigc/video-generation/video-synthesis cURL

 cURL - Multi-shot narrative

 Copy ```\ncurl --location 'https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis' \
 -H 'X-DashScope-Async: enable' \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY" \
 -H 'Content-Type: application/json' \
 -d '{
 "model": "wan2.7-t2v-2026-06-12",
 "input": {
 "prompt": "A tense detective story with cinematic storytelling. Shot 1 [0\u20133 seconds] wide shot: Rainy New York street at night, neon lights flicker, a detective in a black trench coat walks briskly. Shot 2 [3\u20136 seconds] medium shot: The detective enters an old building, rain wets his coat, the door closes slowly behind him. Shot 3 [6\u20139 seconds] close-up: The detective\u2019s focused eyes, distant sirens sound, he frowns slightly. Shot 4 [9\u201312 seconds] medium shot: The detective moves carefully down a dim hallway, his flashlight illuminating the way. Shot 5 [12\u201315 seconds] close-up: The detective discovers a key clue, his face shows sudden realization."
 },
 "parameters": {
 "resolution": "720P",
 "ratio": "16:9",
 "prompt_extend": true,
 "watermark": true,
 "duration": 15
 }
}'
``` 200 400 Copy ```\n{
 "request_id": "4909100c-7b5a-9f92-bfe5-xxxxxx",
 "output": {
 "task_id": "0385dc79-5ff8-4d82-bcb6-xxxxxx",
 "task_status": "PENDING"
 }
}
``` Generate videos up to 15 seconds at 1080P from text prompts, with optional audio sync and multi-shot narrative.
## [​ ](#changes-from-wan2-6) Changes from wan2.6


- **Resolution control**: Set `resolution` (720P/1080P) + `ratio` (16:9, 9:16, etc.) instead of exact pixel `size`.

- **Longer prompts**: Up to 5,000 characters (was 1,500).

- **Negative prompt moved**: Now under `input.negative_prompt` instead of `parameters.negative_prompt`.

- **No `shot_type` parameter**: Describe shots directly in the prompt text.

- **Watermark off by default**: `watermark` defaults to `false` (was `true`).


 ### Authorizations
 [​ ](#authorization) Authorizationstring header required DashScope API Key. Get one from the [Qwen Cloud console](https://home.qwencloud.com/api-keys).

 ### Header Parameters
 [​ ](#x-dashscope-async) X-DashScope-Asyncenum<string> required Must be set to `enable` for asynchronous task submission.

 Available options:enable ### Body
application/json [​ ](#model) modelenum<string> required Model identifier. Supported values: `wan2.7-t2v`, `wan2.7-t2v-2026-04-25`, `wan2.7-t2v-2026-06-12`.

 Available options:wan2.7-t2v,wan2.7-t2v-2026-04-25,wan2.7-t2v-2026-06-12 Example:wan2.7-t2v [​ ](#input) inputobject required Input content for video generation.

 Show child attributes

 [​ ](#inputprompt) input.promptstring required Describe the video you want. Supports Chinese and English, up to 5,000 characters (auto-truncated if exceeded). For multi-shot videos, describe each shot with timestamps: `Shot 1 [0-3 seconds] wide shot: ...`.

 Example:A kitten running in the moonlight. [​ ](#inputnegative-prompt) input.negative_promptstring Describe what to exclude from the video (e.g., `low quality, blurry, extra fingers`). Supports Chinese and English, up to 500 characters (auto-truncated if exceeded).

 Example:low resolution, error, worst quality, low quality, deformed, extra fingers, bad proportions [​ ](#inputaudio-url) input.audio_urlstring Audio file URL for lip-sync and action-to-audio alignment. The model matches mouth movements and actions to the audio track. Supports WAV and MP3 via HTTP/HTTPS, 2-30 seconds, up to 15 MB. Audio longer than `duration` is truncated; audio shorter than the video leaves the remainder silent. If omitted, the model generates matching background music or sound effects automatically.

 [​ ](#parameters) parametersobject Video generation parameters.

 Show child attributes

 [​ ](#parametersresolution) parameters.resolutionenum<string> default"1080P" Video clarity tier. Higher resolution costs more.


Actual output dimensions depend on `ratio`:


- **720P**: 16:9=1280x720, 9:16=720x1280, 1:1=960x960, 4:3=1104x832, 3:4=832x1104

- **1080P**: 16:9=1920x1080, 9:16=1080x1920, 1:1=1440x1440, 4:3=1648x1248, 3:4=1248x1648


 Available options:720P,1080P [​ ](#parametersratio) parameters.ratioenum<string> default"16:9" Aspect ratio of the generated video. Default: `16:9`.

 Available options:16:9,9:16,1:1,4:3,3:4 [​ ](#parametersduration) parameters.durationinteger default5 Video duration in seconds (integer, 2-15). Longer videos cost more -- billing is per second.

 Required range:2 <= x <= 15 [​ ](#parametersprompt-extend) parameters.prompt_extendboolean defaulttrue Rewrite the prompt using a Large Language Model before generation. Improves results for short or vague prompts but adds latency. Set to `false` to use your prompt verbatim.

 [​ ](#parameterswatermark) parameters.watermarkboolean defaultfalse Add an "AI Generated" watermark in the lower-right corner.

 [​ ](#parametersseed) parameters.seedinteger Seed for reproducible results. Same seed + same parameters produces similar (not identical) output.

 Required range:0 <= x <= 2147483647 ### Response
200-application/json [​ ](#request-id) request_idstring Unique request identifier for tracing and troubleshooting.

 Example:4909100c-7b5a-9f92-bfe5-xxxxxx [​ ](#output) outputobject Show child attributes

 [​ ](#outputtask-id) output.task_idstring Task ID for polling status. Use with `GET /tasks/{task_id}`.

 Example:0385dc79-5ff8-4d82-bcb6-xxxxxx [​ ](#outputtask-status) output.task_statusenum<string> Initial task status, typically `PENDING`.

 Available options:PENDING,RUNNING,SUCCEEDED,FAILED
