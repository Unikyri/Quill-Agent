# Wan 2.6 and earlier — Generate a video from text

> **Source:** https://docs.qwencloud.com/api-reference/video-generation/wan-text-to-video/create-task

POST/services/aigc/video-generation/video-synthesis Python

 Python (Synchronous call)

 Copy ```\nimport os
from http import HTTPStatus
from dashscope import VideoSynthesis
import dashscope

dashscope.base_http_api_url = 'https://dashscope-intl.aliyuncs.com/api/v1'
api_key = os.getenv("DASHSCOPE_API_KEY")

print('please wait...')
rsp = VideoSynthesis.call(api_key=api_key,
 model='wan2.6-t2v',
 prompt='A thrilling detective chase story with cinematic storytelling. Shot 1 [0\u20133 s]: Wide shot of a rainy New York street at night, neon lights flickering, a detective in a black trench coat walking briskly. Shot 2 [3\u20136 s]: Medium shot of the detective entering an old building, rain soaking his coat, the door closing slowly behind him. Shot 3 [6\u20139 s]: Close-up of the detective\'s focused, determined eyes as distant sirens wail and he frowns slightly in thought. Shot 4 [9\u201312 s]: Medium shot of the detective moving carefully down a dim hallway, his flashlight illuminating the path ahead. Shot 5 [12\u201315 s]: Close-up of the detective discovering a key clue, his face lighting up with sudden realization.',
 size="1280*720",
 duration=15,
 shot_type="multi",
 prompt_extend=True,
 watermark=True)
print(rsp)
if rsp.status_code == HTTPStatus.OK:
 print("video_url:", rsp.output.video_url)
else:
 print('Failed, status_code: %s, code: %s, message: %s' % (rsp.status_code, rsp.code, rsp.message))
``` 200 400 Copy ```\n{
 "request_id": "c1209113-8437-424f-a386-xxxxxx",
 "output": {
 "task_id": "966cebcd-dedc-4962-af88-xxxxxx",
 "task_status": "PENDING"
 }
}
``` The Wan text-to-video model accepts text, images, and audio as input and generates videos up to 15 seconds long at 1080P resolution.

- **Core capabilities**: Integer video durations (2–15 seconds), custom resolutions (480P, 720P, 1080P), prompt rewriting, and watermarking.

- **Audio capabilities**: Automatic dubbing or custom audio files for audio-video sync. **(Supported by wan2.5 and wan2.6)**

- **Multi-shot narrative**: Multiple shots with consistent main subject across transitions. **(Supported only by wan2.6)**


 ### Authorizations
 [​ ](#authorization) Authorizationstring header required DashScope API Key. Get one from the [Qwen Cloud console](https://home.qwencloud.com/api-keys).

 ### Header Parameters
 [​ ](#x-dashscope-async) X-DashScope-Asyncenum<string> required Must be set to `enable` for asynchronous task submission.

 Available options:enable ### Body
application/json [​ ](#model) modelenum<string> required Model name. See the model table in the endpoint description for supported models and their capabilities.

 Available options:wan2.6-t2v,wan2.5-t2v-preview,wan2.2-t2v-plus,wan2.1-t2v-turbo,wan2.1-t2v-plus Example:wan2.6-t2v [​ ](#input) inputobject required Input data for video generation.

 Show child attributes

 [​ ](#inputprompt) input.promptstring required Text prompt describing the desired video content. Maximum length: 1,500 characters for wan2.6 and wan2.5 series; 800 characters for wan2.2 and wan2.1 series. Text beyond the limit is auto-truncated. For multi-shot videos (wan2.6), use the format: `Shot 1 [0–3 s]: description. Shot 2 [3–6 s]: description.` etc.

 Example:A thrilling detective chase story with cinematic storytelling. [​ ](#inputaudio-url) input.audio_urlstring URL of an audio file for audio-video synchronization. The model aligns mouth movements and actions to the audio. Supports HTTP/HTTPS URLs. **Supported by wan2.5 and wan2.6 series only.** If omitted on wan2.5/wan2.6, the model auto-generates background audio (automatic dubbing).

 [​ ](#parameters) parametersobject Video generation parameters.

 Show child attributes

 [​ ](#parameterssize) parameters.sizestring Output video resolution as `width*height`. Available sizes depend on the model:


- **wan2.6-t2v**: `1280*720` (720P), `1920*1080` (1080P)

- **wan2.5-t2v-preview**: `832*480` (480P), `1280*720` (720P), `1920*1080` (1080P)

- **wan2.2-t2v-plus**: `832*480` (480P), `1920*1080` (1080P)

- **wan2.1-t2v-turbo**: `832*480` (480P), `1280*720` (720P)

- **wan2.1-t2v-plus**: `1280*720` (720P)


 Example:1280*720 [​ ](#parametersduration) parameters.durationinteger Video duration in seconds. Available durations depend on the model:


- **wan2.6-t2v**: Integer from 2 to 15

- **wan2.5-t2v-preview**: 5 or 10

- **wan2.2-t2v-plus, wan2.1 series**: Fixed at 5


 Example:15 [​ ](#parametersshot-type) parameters.shot_typeenum<string> Shot composition mode. Set to `"multi"` to enable multi-shot narrative with automatic shot transitions. **Supported by wan2.6 series only.**

 Available options:multi [​ ](#parametersprompt-extend) parameters.prompt_extendboolean defaulttrue Enable prompt rewriting. `true` (default): the model optimizes your prompt for better results. `false`: use your prompt as-is. Recommended to enable for multi-shot videos.

 [​ ](#parameterswatermark) parameters.watermarkboolean defaulttrue Add a watermark to the generated video. Default: `true`.

 [​ ](#parametersnegative-prompt) parameters.negative_promptstring Describes content you do NOT want in the video.

 [​ ](#parametersseed) parameters.seedinteger Random number seed for reproducibility. Range: [0, 2147483647]. Same seed with the same parameters produces more consistent (but not identical) results.

 Required range:0 <= x <= 2147483647 ### Response
200-application/json [​ ](#request-id) request_idstring Unique request identifier for tracing and troubleshooting.

 Example:c1209113-8437-424f-a386-xxxxxx [​ ](#output) outputobject Show child attributes

 [​ ](#outputtask-id) output.task_idstring Task ID for polling status. Use with `GET /tasks/{task_id}`.

 Example:966cebcd-dedc-4962-af88-xxxxxx [​ ](#outputtask-status) output.task_statusenum<string> Initial task status, typically `PENDING`.

 Available options:PENDING,RUNNING,SUCCEEDED,FAILED
