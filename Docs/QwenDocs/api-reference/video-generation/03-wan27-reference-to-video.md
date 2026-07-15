# Wan 2.7 -- Generate from reference

> **Source:** https://docs.qwencloud.com/api-reference/video-generation/wan27-reference-to-video/create-task

POST/services/aigc/video-generation/video-synthesis cURL

 cURL - Multi-subject reference

 Copy ```\ncurl --location 'https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis' \
 -H 'X-DashScope-Async: enable' \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY" \
 -H 'Content-Type: application/json' \
 -d '{
 "model": "wan2.7-r2v-2026-06-12",
 "input": {
 "prompt": "Video 2 holds Image 3 and plays a soothing American country ballad in a coffee shop, while Video 1 smiles, watches Video 2, and slowly walks towards him",
 "media": [
 {"type": "reference_video", "url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20260129/hfugmr/wan-r2v-role1.mp4"},
 {"type": "reference_video", "url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20260129/qigswt/wan-r2v-role2.mp4"},
 {"type": "reference_image", "url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20260129/qpzxps/wan-r2v-object4.png"}
 ]
 },
 "parameters": {
 "resolution": "720P",
 "duration": 10,
 "prompt_extend": false,
 "watermark": true
 }
}'
``` 200 400 Copy ```\n{
 "request_id": "&#x3C;string>",
 "output": {
 "task_id": "&#x3C;string>",
 "task_status": "PENDING"
 }
}
``` Generate natural, lifelike performance videos from multimodal input (text, image, video) using the Wan 2.7 model (`wan2.7-r2v`).

- **Character portrayal**: Replicate a character&#x27;s appearance from a reference image or video. Reference videos also replicate voice timbre. Supports single or multi-character performances with up to 5 reference assets.

- **Media array input**: Provide reference images, videos, or a first frame via the `media` array. Use `Video 1`/`Image 1` in prompts to reference characters by their order. Images and videos are counted separately.

- **Multi-panel storyboard**: Describe multi-shot narratives with time segments (e.g., `Scene 1 [0-3s]: ...`). Provide key shots and the model automatically recognizes the panel logic.

- **Voice cloning**: Provide a `reference_voice` audio file to set the voice timbre. If not specified, audio from the reference video is used by default.

- **Resolution and ratio**: Set output quality with `resolution` (720P/1080P) and aspect ratio with `ratio` (16:9, 9:16, 1:1, 4:3, 3:4). When a `first_frame` image is provided, `ratio` is inferred from the image.

- **Prompt enhancement**: Enable `prompt_extend` to rewrite the prompt with an LLM. Improves results for shorter prompts but increases processing time.


 ### Authorizations
 [​ ](#authorization) Authorizationstring header required DashScope API Key. Create one in the [Qwen Cloud console](https://home.qwencloud.com/api-keys).

 ### Header Parameters
 [​ ](#x-dashscope-async) X-DashScope-Asyncenum<string> required Must be `enable` to create an asynchronous task.

 Available options:enable ### Body
application/json [​ ](#model) modelenum<string> required Model name.

 Available options:wan2.7-r2v,wan2.7-r2v-2026-06-12 Example:wan2.7-r2v [​ ](#input) inputobject required Input data for Wan 2.7 reference-to-video generation.

 Show child attributes

 [​ ](#inputprompt) input.promptstring required Text prompt describing the desired video content. Supports Chinese and English. Each Chinese character, letter, and punctuation mark counts as one character. Text that exceeds the limit is automatically truncated.


**Reference identifiers**: Use `Image 1`, `Image 2`, etc. to reference characters from reference images, and `Video 1`, `Video 2`, etc. to reference characters from reference videos. The numbering matches the order in the `media` array. Images and videos are counted separately — `Image 1` and `Video 1` can coexist. If there is only one reference image or video, you can use `the reference image` or `the reference video` instead.


**Scene description**: Two approaches — (1) Use identifiers directly: "Image 1 is playing in Image 2"; (2) Supplement with subject/scene context: "The cat from Image 1 is playing in the room from Image 2".


**Multi-panel storyboard**: Describe multi-shot narratives with time segments (e.g., `Scene 1 [0-3s]: ...`). You do not need to describe every panel — provide key shots, and the model automatically recognizes the panel logic.

 Example:Video 2 holds Image 3 and plays a soothing American country ballad in a coffee shop Required range:length <= 5000 [​ ](#inputmedia) input.mediaobject[] required Array of reference media objects. Each object has a `type` and `url`. Supports image and video input for visual reference. Images support multiple views, commonly for referencing characters, props, and scenes.


**Ordering**: The first `reference_video` in the array is `Video 1`, the second is `Video 2`, and so on. The first `reference_image` is `Image 1`, the second is `Image 2`. Images and videos are counted separately.


**Limits**: At least 1 reference image or video is required. Total images + videos must not exceed 5. At most 1 `first_frame` is allowed. Each reference image or video should contain a single character when used for main character portrayal.

 Required range:items: 1–5 Show child attributes

 [​ ](#inputmediatype) input.media.typeenum<string> required Type of reference media.


- `reference_image`: Reference image containing a single character or object. Supported formats: JPEG, JPG, PNG (alpha channel not supported), BMP, WEBP. Resolution: 240-8000 px per side. Aspect ratio: 1:8 to 8:1. Max file size: 20 MB.

- `reference_video`: Reference video containing a single character. Supported formats: MP4, MOV. Duration: 1-30 seconds. Resolution: 240-4096 px per side. Aspect ratio: 1:8 to 8:1. Max file size: 100 MB.

- `first_frame`: First frame image for the generated video. At most 1 allowed. Supported formats and constraints same as `reference_image`. When used with subject references, two modes apply: (1) Subject already in the first frame — use subject reference to enhance character consistency or reference voice timbre. (2) Subject not in the first frame — use subject reference to define features of a new character appearing during the video.


 Available options:reference_image,reference_video,first_frame [​ ](#inputmediaurl) input.media.urlstring required URL of the reference media file.

 [​ ](#inputnegative-prompt) input.negative_promptstring Content to exclude from the generated video. Supports Chinese and English. Maximum 500 characters. Text that exceeds the limit is automatically truncated. Example values: `low resolution, error, worst quality, low quality, disfigured, extra fingers, bad proportions`.

 Required range:length <= 500 [​ ](#inputreference-voice) input.reference_voicestring URL of a reference audio file for voice timbre. Used with `reference_image` or `reference_video` to specify the voice timbre for the main character. This audio is used only for voice timbre reference and is not related to the spoken content.


Supported formats: WAV, MP3. Duration: 1-10 seconds. Max file size: 15 MB.


**Default behavior**: If `reference_video` has audio and `reference_voice` is not specified, the original video audio is used as the voice timbre.


**Priority**: If both `reference_video` audio and `reference_voice` are provided, `reference_voice` takes priority.


For best results, the language of the reference audio should match the language of the prompt.

 [​ ](#parameters) parametersobject Generation parameters for Wan 2.7 reference-to-video.

 Show child attributes

 [​ ](#parametersresolution) parameters.resolutionenum<string> default"1080P" Video clarity tier. Higher resolution costs more.


Actual output dimensions depend on `ratio`:


- **720P**: 16:9=1280x720, 9:16=720x1280, 1:1=960x960, 4:3=1104x832, 3:4=832x1104

- **1080P**: 16:9=1920x1080, 9:16=1080x1920, 1:1=1440x1440, 4:3=1648x1248, 3:4=1248x1648


 Available options:720P,1080P [​ ](#parametersratio) parameters.ratioenum<string> default"16:9" Aspect ratio of the output video. When a `first_frame` image is provided, `ratio` is ignored and the video uses the image&#x27;s aspect ratio.

 Available options:16:9,9:16,1:1,4:3,3:4 [​ ](#parametersduration) parameters.durationinteger default5 Video duration in seconds. Longer videos cost more -- billing is per second. Conditional range: if reference material includes a video, the range is 2-10; if reference material does not include a video (images only), the range is 2-15.

 Required range:2 <= x <= 15 [​ ](#parametersprompt-extend) parameters.prompt_extendboolean defaulttrue Whether to rewrite the prompt using a large language model. Improves results for shorter prompts but increases processing time.

 [​ ](#parametersseed) parameters.seedinteger Random seed for reproducible generation. If not specified, a random seed is generated. A fixed seed improves reproducibility, but because model generation is probabilistic, the same seed does not guarantee identical results.

 Required range:0 <= x <= 2147483647 [​ ](#parameterswatermark) parameters.watermarkboolean defaultfalse Whether to add an "AI Generated" watermark to the lower-right corner of the output video.

 ### Response
200-application/json [​ ](#request-id) request_idstring Unique request identifier.

 [​ ](#output) outputobject Show child attributes

 [​ ](#outputtask-id) output.task_idstring Task identifier. Use this with `GET /tasks/{task_id}` to poll for results.

 [​ ](#outputtask-status) output.task_statusenum<string> Initial task status, typically `PENDING`.

 Available options:PENDING
