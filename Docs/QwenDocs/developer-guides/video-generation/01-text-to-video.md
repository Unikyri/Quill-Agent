# Text-to-video

> **Source:** https://docs.qwencloud.com/developer-guides/video-generation/text-to-video

Generate video from text

 Copy page The Wan text-to-video model supports **multimodal input** — including text and audio — and generates videos up to 15 seconds long at 1080P resolution.

- **Core capabilities**: Supports integer video durations (2-15 seconds), custom video resolutions (720P or 1080P), aspect ratio control, prompt rewriting, and watermarking.

- **Audio capabilities**: Supports automatic dubbing or custom audio files for synchronized audio and video. **(Supported by wan2.5 and later)**

- **Multi-shot narrative**: Generates videos with multiple shots while keeping the main subject consistent across shot transitions. **(Supported by wan2.6 and wan2.7)**


**Quick access:** [Try it online](https://home.qwencloud.com/try-ai) **|** API reference: [wan2.7](/api-reference/video-generation/wan27-text-to-video/create-task), [wan2.6](/api-reference/video-generation/wan-text-to-video/create-task) **|** [Prompt guide](/developer-guides/accuracy-tuning/video-generation)
## [​ ](#getting-started) Getting started

**Input prompt****Output video (multi-shot, audio-enabled)**A thrilling detective chase story with cinematic storytelling. Shot 1 [0-3 s]: Wide shot of a rainy New York street at night, neon lights flickering, a detective in a black trench coat walking briskly. Shot 2 [3-6 s]: Medium shot of the detective entering an old building, rain soaking his coat, the door closing slowly behind him. Shot 3 [6-9 s]: Close-up of the detective&#x27;s focused, determined eyes as distant sirens wail and he frowns slightly in thought. Shot 4 [9-12 s]: Medium shot of the detective moving carefully down a dim hallway, his flashlight illuminating the path ahead. Shot 5 [12-15 s]: Close-up of the detective discovering a key clue, his face lighting up with sudden realization.
Before calling the API, [get an API key](/api-reference/preparation/api-key). Then [set your API key as an environment variable](/api-reference/preparation/export-api-key-env). To use the SDK, [install the DashScope SDK](/api-reference/preparation/install-sdk).
 Wan 2.7 uses `resolution` + `ratio` instead of `size`, and describes multi-shot directly in the prompt (no `shot_type` parameter). 
**Step 1: Create a task to get the task ID**
Copy ```\ncurl --location &#x27;https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis&#x27; \
 -H &#x27;X-DashScope-Async: enable&#x27; \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY" \
 -H &#x27;Content-Type: application/json&#x27; \
 -d &#x27;{
 "model": "wan2.7-t2v",
 "input": {
 "prompt": "A thrilling detective chase story with cinematic storytelling. Shot 1 [0-3 s]: Wide shot of a rainy New York street at night, neon lights flickering, a detective in a black trench coat walking briskly. Shot 2 [3-6 s]: Medium shot of the detective entering an old building, rain soaking his coat, the door closing slowly behind him. Shot 3 [6-9 s]: Close-up of the detective focused, determined eyes as distant sirens wail and he frowns slightly in thought. Shot 4 [9-12 s]: Medium shot of the detective moving carefully down a dim hallway, his flashlight illuminating the path ahead. Shot 5 [12-15 s]: Close-up of the detective discovering a key clue, his face lighting up with sudden realization."
 },
 "parameters": {
 "resolution": "1080P",
 "ratio": "16:9",
 "prompt_extend": true,
 "duration": 15
 }
}&#x27;

``` 
**Step 2: Get the result using the task ID**
Replace `task_id` with the `task_id` value returned by the previous API call.
Copy ```\ncurl -X GET https://dashscope-intl.aliyuncs.com/api/v1/tasks/{task_id} \
--header "Authorization: Bearer $DASHSCOPE_API_KEY"

``` 
 wan2.6 examples (Python SDK, Java SDK, curl)

 To use the SDK, [install the DashScope SDK](/api-reference/preparation/install-sdk).- Python SDK 
- Java SDK 
- curl 

 Make sure your DashScope Python SDK version is **at least** `1.25.8` before running the code below.If your version is too low, you may see errors such as "url error, please check url!". [Install the SDK](/api-reference/preparation/install-sdk). Copy ```\nimport os
from http import HTTPStatus
from dashscope import VideoSynthesis
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;
api_key = os.getenv("DASHSCOPE_API_KEY", "YOUR_API_KEY")

print(&#x27;please wait...&#x27;)
rsp = VideoSynthesis.call(api_key=api_key,
 model=&#x27;wan2.6-t2v&#x27;,
 prompt=&#x27;A thrilling detective chase story with cinematic storytelling. Shot 1 [0–3 s]: Wide shot of a rainy New York street at night, neon lights flickering, a detective in a black trench coat walking briskly. Shot 2 [3–6 s]: Medium shot of the detective entering an old building, rain soaking his coat, the door closing slowly behind him. Shot 3 [6–9 s]: Close-up of the detective\&#x27;s focused, determined eyes as distant sirens wail and he frowns slightly in thought. Shot 4 [9–12 s]: Medium shot of the detective moving carefully down a dim hallway, his flashlight illuminating the path ahead. Shot 5 [12–15 s]: Close-up of the detective discovering a key clue, his face lighting up with sudden realization.&#x27;,
 size="1280*720",
 duration=15,
 shot_type="multi",
 prompt_extend=True,
 watermark=True)
print(rsp)
if rsp.status_code == HTTPStatus.OK:
 print("video_url:", rsp.output.video_url)
else:
 print(&#x27;Failed, status_code: %s, code: %s, message: %s&#x27; % (rsp.status_code, rsp.code, rsp.message))

``` Make sure your DashScope Java SDK version is **at least** `2.22.6` before running the code below.If your version is too low, you may see errors such as "url error, please check url!". [Install the SDK](/api-reference/preparation/install-sdk). Copy ```\nimport com.alibaba.dashscope.aigc.videosynthesis.VideoSynthesis;
import com.alibaba.dashscope.aigc.videosynthesis.VideoSynthesisParam;
import com.alibaba.dashscope.aigc.videosynthesis.VideoSynthesisResult;
import com.alibaba.dashscope.exception.ApiException;
import com.alibaba.dashscope.exception.InputRequiredException;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.alibaba.dashscope.utils.JsonUtils;
import com.alibaba.dashscope.utils.Constants;

public class Text2Video {

 static {
 Constants.baseHttpApiUrl = "https://dashscope-intl.aliyuncs.com/api/v1";
 }

 // If you have not set an environment variable, replace the line below with: apiKey="sk-xxx"
 static String apiKey = System.getenv("DASHSCOPE_API_KEY");

 public static void text2video() throws ApiException, NoApiKeyException, InputRequiredException {
 VideoSynthesis vs = new VideoSynthesis();
 VideoSynthesisParam param =
 VideoSynthesisParam.builder()
 .apiKey(apiKey)
 .model("wan2.6-t2v")
 .prompt("A thrilling detective chase story with cinematic storytelling. Shot 1 [0–3 s]: Wide shot of a rainy New York street at night, neon lights flickering, a detective in a black trench coat walking briskly. Shot 2 [3–6 s]: Medium shot of the detective entering an old building, rain soaking his coat, the door closing slowly behind him. Shot 3 [6–9 s]: Close-up of the detective&#x27;s focused, determined eyes as distant sirens wail and he frowns slightly in thought. Shot 4 [9–12 s]: Medium shot of the detective moving carefully down a dim hallway, his flashlight illuminating the path ahead. Shot 5 [12–15 s]: Close-up of the detective discovering a key clue, his face lighting up with sudden realization.")
 .duration(15)
 .size("1280*720")
 .shotType("multi")
 .promptExtend(true)
 .watermark(true)
 .build();
 System.out.println("please wait...");
 VideoSynthesisResult result = vs.call(param);
 System.out.println(JsonUtils.toJson(result));
 }

 public static void main(String[] args) {
 try {
 text2video();
 } catch (ApiException | NoApiKeyException | InputRequiredException e) {
 System.out.println(e.getMessage());
 }
 System.exit(0);
 }
}

``` **Step 1: Create a task to get the task ID**Copy ```\ncurl --location &#x27;https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis&#x27; \
 -H &#x27;X-DashScope-Async: enable&#x27; \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY" \
 -H &#x27;Content-Type: application/json&#x27; \
 -d &#x27;{
 "model": "wan2.6-t2v",
 "input": {
 "prompt": "A thrilling detective chase story with cinematic storytelling. Shot 1 [0–3 s]: Wide shot of a rainy New York street at night, neon lights flickering, a detective in a black trench coat walking briskly. Shot 2 [3–6 s]: Medium shot of the detective entering an old building, rain soaking his coat, the door closing slowly behind him. Shot 3 [6–9 s]: Close-up of the detective&#x27;s focused, determined eyes as distant sirens wail and he frowns slightly in thought. Shot 4 [9–12 s]: Medium shot of the detective moving carefully down a dim hallway, his flashlight illuminating the path ahead. Shot 5 [12–15 s]: Close-up of the detective discovering a key clue, his face lighting up with sudden realization."
 },
 "parameters": {
 "size": "1280*720",
 "prompt_extend": true,
 "watermark": true,
 "duration": 15,
 "shot_type":"multi"
 }
}&#x27;

``` **Step 2: Get the result using the task ID**Replace `task_id` with the `task_id` value returned by the previous API call.Copy ```\ncurl -X GET https://dashscope-intl.aliyuncs.com/api/v1/tasks/{task_id} \
--header "Authorization: Bearer $DASHSCOPE_API_KEY"

``` 
**Sample output**
 `video_url` expires after 24 hours. Download the video promptly. 
Copy ```\n{
 "request_id": "c1209113-8437-424f-a386-xxxxxx",
 "output": {
 "task_id": "966cebcd-dedc-4962-af88-xxxxxx",
 "task_status": "SUCCEEDED",
 "video_url": "https://dashscope-result-sh.oss-accelerate.aliyuncs.com/xxx.mp4?Expires=xxx",
 ...
 },
 ...
}

``` 
## [​ ](#core-capabilities) Core capabilities

### [​ ](#create-multi-shot-videos) Create multi-shot videos

**Supported models**: `wan2.7`, `wan2.6 series`.
**Description**: The model automatically switches between shots — for example, from a wide shot to a close-up — ideal for music videos and similar use cases.
**Parameters**:

- **wan2.7**: Describe shots directly in the prompt text (e.g., `Shot 1 [0-3 s]: ...`). No `shot_type` parameter needed.

- **wan2.6**: Set `shot_type` to `"multi"`.

- `prompt_extend`: Set to `true` (enables prompt rewriting to optimize shot descriptions).


**Input prompt****Output video (multi-shot video)**A vision of harmony between future technology and nature. Shot 1 [0-2 s]: Wide shot of an aerial garden in a futuristic city, floating plants swaying gently in the breeze. Shot 2 [2-4 s]: A robot gardener carefully trims plants with precise, graceful movements. Shot 3 [4-7 s]: Sunlight streams through a transparent dome, illuminating the entire garden and showcasing perfect fusion of technology and nature. Shot 4 [7-10 s]: The camera pulls back to reveal the grand scale of the entire futuristic city, with the aerial garden just one part of it.
- Python SDK 
- Java SDK 
- curl 

 Make sure your DashScope Python SDK version is at least `1.25.8`. [Install the SDK](/api-reference/preparation/install-sdk). Copy ```\nimport os
from http import HTTPStatus
from dashscope import VideoSynthesis
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

# If you have not set an environment variable, replace the line below with: api_key="sk-xxx"
api_key = os.getenv("DASHSCOPE_API_KEY")

def sample_async_call_t2v():
 # Asynchronous call returns a task_id
 rsp = VideoSynthesis.async_call(api_key=api_key,
 model=&#x27;wan2.6-t2v&#x27;,
 prompt=&#x27;A vision of harmony between future technology and nature. Shot 1 [0–2 s]: Wide shot of an aerial garden in a futuristic city, floating plants swaying gently in the breeze. Shot 2 [2–4 s]: A robot gardener carefully trims plants with precise, graceful movements. Shot 3 [4–7 s]: Sunlight streams through a transparent dome, illuminating the entire garden and showcasing perfect fusion of technology and nature. Shot 4 [7–10 s]: The camera pulls back to reveal the grand scale of the entire futuristic city, with the aerial garden just one part of it.&#x27;,
 size=&#x27;1280*720&#x27;,
 shot_type="multi", # Multi-shot
 duration=10,
 prompt_extend=True,
 watermark=True,
 negative_prompt="",
 seed=12345)
 print(rsp)
 if rsp.status_code == HTTPStatus.OK:
 print("task_id: %s" % rsp.output.task_id)
 else:
 print(&#x27;Failed, status_code: %s, code: %s, message: %s&#x27; % (rsp.status_code, rsp.code, rsp.message))

 # Wait for asynchronous task to complete
 rsp = VideoSynthesis.wait(task=rsp, api_key=api_key)
 print(rsp)
 if rsp.status_code == HTTPStatus.OK:
 print(rsp.output.video_url)
 else:
 print(&#x27;Failed, status_code: %s, code: %s, message: %s&#x27; % (rsp.status_code, rsp.code, rsp.message))


if __name__ == &#x27;__main__&#x27;:
 sample_async_call_t2v()

``` Make sure your DashScope Java SDK version is at least `2.22.6`. [Install the SDK](/api-reference/preparation/install-sdk). Copy ```\nimport com.alibaba.dashscope.aigc.videosynthesis.VideoSynthesis;
import com.alibaba.dashscope.aigc.videosynthesis.VideoSynthesisParam;
import com.alibaba.dashscope.aigc.videosynthesis.VideoSynthesisResult;
import com.alibaba.dashscope.exception.ApiException;
import com.alibaba.dashscope.exception.InputRequiredException;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.alibaba.dashscope.utils.JsonUtils;
import com.alibaba.dashscope.utils.Constants;

public class Text2Video {
 static {
 Constants.baseHttpApiUrl = "https://dashscope-intl.aliyuncs.com/api/v1";
 }

 // If you have not set an environment variable, replace the line below with: apiKey="sk-xxx"
 static String apiKey = System.getenv("DASHSCOPE_API_KEY");

 public static void text2Video() throws ApiException, NoApiKeyException, InputRequiredException {
 VideoSynthesis vs = new VideoSynthesis();
 VideoSynthesisParam param =
 VideoSynthesisParam.builder()
 .apiKey(apiKey)
 .model("wan2.6-t2v")
 .prompt("A vision of harmony between future technology and nature. Shot 1 [0–2 s]: Wide shot of an aerial garden in a futuristic city, floating plants swaying gently in the breeze. Shot 2 [2–4 s]: A robot gardener carefully trims plants with precise, graceful movements. Shot 3 [4–7 s]: Sunlight streams through a transparent dome, illuminating the entire garden and showcasing perfect fusion of technology and nature. Shot 4 [7–10 s]: The camera pulls back to reveal the grand scale of the entire futuristic city, with the aerial garden just one part of it.")
 .negativePrompt("")
 .size("1280*720")
 .shotType("multi")
 .duration(10)
 .promptExtend(true)
 .watermark(true)
 .seed(12345)
 .build();
 // Asynchronous call
 VideoSynthesisResult task = vs.asyncCall(param);
 System.out.println(JsonUtils.toJson(task));
 System.out.println("please wait...");

 // Get result
 VideoSynthesisResult result = vs.wait(task, apiKey);
 System.out.println(JsonUtils.toJson(result));
 }

 public static void main(String[] args) {
 try {
 text2video();
 } catch (ApiException | NoApiKeyException | InputRequiredException e) {
 System.out.println(e.getMessage());
 }
 System.exit(0);
 }
}

``` **Step 1: Create a task to get the task ID**Copy ```\ncurl --location &#x27;https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis&#x27; \
 -H &#x27;X-DashScope-Async: enable&#x27; \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY" \
 -H &#x27;Content-Type: application/json&#x27; \
 -d &#x27;{
 "model": "wan2.6-t2v",
 "input": {
 "prompt": "A vision of harmony between future technology and nature. Shot 1 [0-2 s]: Wide shot of an aerial garden in a futuristic city, floating plants swaying gently in the breeze. Shot 2 [2-4 s]: A robot gardener carefully trims plants with precise, graceful movements. Shot 3 [4-7 s]: Sunlight streams through a transparent dome, illuminating the entire garden and showcasing perfect fusion of technology and nature. Shot 4 [7-10 s]: The camera pulls back to reveal the grand scale of the entire futuristic city, with the aerial garden just one part of it."
 },
 "parameters": {
 "size": "1280*720",
 "prompt_extend": true,
 "duration": 10,
 "shot_type":"multi"
 }
}&#x27;

``` **Step 2: Get the result using the task ID**Replace `task_id` with the `task_id` value returned by the previous API call.Copy ```\ncurl -X GET https://dashscope-intl.aliyuncs.com/api/v1/tasks/{task_id} \
--header "Authorization: Bearer $DASHSCOPE_API_KEY"

``` 
### [​ ](#synchronize-audio-and-video) Synchronize audio and video

**Supported models**: `wan2.7`, `wan2.6 series`, `wan2.5 series`.
**Description**: Make characters in photos speak or sing, with mouth movements matching the audio. For more examples, see [Video audio generation](/developer-guides/accuracy-tuning/video-generation).
**Parameters**:

- **Provide an audio file**: Pass an `audio_url`. The model aligns mouth movement to the audio.

- **Automatic dubbing**: Audio-enabled video is generated by default. Do not pass `audio_url`. The model auto-generates background sound effects, music, or voice based on the scene.


**Input example****Output video (audio-enabled video)****Input prompt**: Shot from a low angle, in a medium close-up, with warm tones, mixed lighting (the practical light from the desk lamp blends with the overcast light from the window), side lighting, and a central composition. In a classic detective office, wooden bookshelves are filled with old case files and ashtrays. A green desk lamp illuminates a case file spread out in the center of the desk. A fox, wearing a dark brown trench coat and a light gray fedora, sits in a leather chair, its fur crimson, its tail resting lightly on the edge, its fingers slowly turning yellowed pages. Outside, a steady drizzle falls beneath a blue sky, streaking the glass with meandering streaks. It slowly raises its head, its ears twitching slightly, its amber eyes gazing directly at the camera, its mouth clearly moving as it speaks in a smooth, cynical voice: **&#x27;The case was cold, colder than a fish in winter. But every chicken has its secrets, and I, for one, intended to find them &#x27;**. **Input audio**: 
- Python SDK 
- Java SDK 
- curl 

 Make sure your DashScope Python SDK version is at least `1.25.8`. [Install the SDK](/api-reference/preparation/install-sdk). Copy ```\nimport os
from http import HTTPStatus
from dashscope import VideoSynthesis
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

# If you have not set an environment variable, replace the line below with: api_key="sk-xxx"
api_key = os.getenv("DASHSCOPE_API_KEY")

def sample_async_call_t2v():
 # Asynchronous call returns a task_id
 rsp = VideoSynthesis.async_call(api_key=api_key,
 model=&#x27;wan2.6-t2v&#x27;,
 prompt="Shot from a low angle, in a medium close-up, with warm tones, mixed lighting (the practical light from the desk lamp blends with the overcast light from the window), side lighting, and a central composition. In a classic detective office, wooden bookshelves are filled with old case files and ashtrays. A green desk lamp illuminates a case file spread out in the center of the desk. A fox, wearing a dark brown trench coat and a light gray fedora, sits in a leather chair, its fur crimson, its tail resting lightly on the edge, its fingers slowly turning yellowed pages. Outside, a steady drizzle falls beneath a blue sky, streaking the glass with meandering streaks. It slowly raises its head, its ears twitching slightly, its amber eyes gazing directly at the camera, its mouth clearly moving as it speaks in a smooth, cynical voice: &#x27;The case was cold, colder than a fish in winter. But every chicken has its secrets, and I, for one, intended to find them &#x27;.",
 audio_url=&#x27;https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250929/stjqnq/%E7%8B%90%E7%8B%B8.mp3&#x27;,
 size=&#x27;1280*720&#x27;,
 duration=10,
 shot_type="multi", # Multi-shot
 prompt_extend=True,
 watermark=True,
 negative_prompt="",
 seed=12345)
 print(rsp)
 if rsp.status_code == HTTPStatus.OK:
 print("task_id: %s" % rsp.output.task_id)
 else:
 print(&#x27;Failed, status_code: %s, code: %s, message: %s&#x27; % (rsp.status_code, rsp.code, rsp.message))

 # Wait for asynchronous task to complete
 rsp = VideoSynthesis.wait(task=rsp, api_key=api_key)
 print(rsp)
 if rsp.status_code == HTTPStatus.OK:
 print(rsp.output.video_url)
 else:
 print(&#x27;Failed, status_code: %s, code: %s, message: %s&#x27; % (rsp.status_code, rsp.code, rsp.message))


if __name__ == &#x27;__main__&#x27;:
 sample_async_call_t2v()

``` Make sure your DashScope Java SDK version is at least `2.22.6`. [Install the SDK](/api-reference/preparation/install-sdk). Copy ```\nimport com.alibaba.dashscope.aigc.videosynthesis.VideoSynthesis;
import com.alibaba.dashscope.aigc.videosynthesis.VideoSynthesisParam;
import com.alibaba.dashscope.aigc.videosynthesis.VideoSynthesisResult;
import com.alibaba.dashscope.exception.ApiException;
import com.alibaba.dashscope.exception.InputRequiredException;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.alibaba.dashscope.utils.JsonUtils;
import com.alibaba.dashscope.utils.Constants;

public class Text2Video {
 static {
 Constants.baseHttpApiUrl = "https://dashscope-intl.aliyuncs.com/api/v1";
 }

 // If you have not set an environment variable, replace the line below with: apiKey="sk-xxx"
 static String apiKey = System.getenv("DASHSCOPE_API_KEY");

 public static void text2Video() throws ApiException, NoApiKeyException, InputRequiredException {
 VideoSynthesis vs = new VideoSynthesis();
 VideoSynthesisParam param =
 VideoSynthesisParam.builder()
 .apiKey(apiKey)
 .model("wan2.6-t2v")
 .prompt("Shot from a low angle, in a medium close-up, with warm tones, mixed lighting (the practical light from the desk lamp blends with the overcast light from the window), side lighting, and a central composition. In a classic detective office, wooden bookshelves are filled with old case files and ashtrays. A green desk lamp illuminates a case file spread out in the center of the desk. A fox, wearing a dark brown trench coat and a light gray fedora, sits in a leather chair, its fur crimson, its tail resting lightly on the edge, its fingers slowly turning yellowed pages. Outside, a steady drizzle falls beneath a blue sky, streaking the glass with meandering streaks. It slowly raises its head, its ears twitching slightly, its amber eyes gazing directly at the camera, its mouth clearly moving as it speaks in a smooth, cynical voice: &#x27;The case was cold, colder than a fish in winter. But every chicken has its secrets, and I, for one, intended to find them &#x27;.")
 .audioUrl("https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250929/stjqnq/%E7%8B%90%E7%8B%B8.mp3")
 .negativePrompt("")
 .size("1280*720")
 .shotType("multi")
 .duration(10)
 .promptExtend(true)
 .watermark(true)
 .seed(12345)
 .build();
 // Asynchronous call
 VideoSynthesisResult task = vs.asyncCall(param);
 System.out.println(JsonUtils.toJson(task));
 System.out.println("please wait...");

 // Get result
 VideoSynthesisResult result = vs.wait(task, apiKey);
 System.out.println(JsonUtils.toJson(result));
 }

 public static void main(String[] args) {
 try {
 text2video();
 } catch (ApiException | NoApiKeyException | InputRequiredException e) {
 System.out.println(e.getMessage());
 }
 System.exit(0);
 }
}

``` **Step 1: Create a task to get the task ID**Copy ```\ncurl --location &#x27;https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis&#x27; \
 -H &#x27;X-DashScope-Async: enable&#x27; \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY" \
 -H &#x27;Content-Type: application/json&#x27; \
 -d &#x27;{
 "model": "wan2.6-t2v",
 "input": {
 "prompt": "Shot from a low angle, in a medium close-up, with warm tones, mixed lighting (the practical light from the desk lamp blends with the overcast light from the window), side lighting, and a central composition. In a classic detective office, wooden bookshelves are filled with old case files and ashtrays. A green desk lamp illuminates a case file spread out in the center of the desk. A fox, wearing a dark brown trench coat and a light gray fedora, sits in a leather chair, its fur crimson, its tail resting lightly on the edge, its fingers slowly turning yellowed pages. Outside, a steady drizzle falls beneath a blue sky, streaking the glass with meandering streaks. It slowly raises its head, its ears twitching slightly, its amber eyes gazing directly at the camera, its mouth clearly moving as it speaks in a smooth, cynical voice: \"The case was cold, colder than a fish in winter. But every chicken has its secrets, and I, for one, intended to find them \". ",
 "audio_url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250929/stjqnq/%E7%8B%90%E7%8B%B8.mp3"
 },
 "parameters": {
 "size": "1280*720",
 "prompt_extend": true,
 "duration": 10,
 "shot_type":"multi"
 }
}&#x27;

``` **Step 2: Get the result using the task ID**Replace `task_id` with the `task_id` value returned by the previous API call.Copy ```\ncurl -X GET https://dashscope-intl.aliyuncs.com/api/v1/tasks/{task_id} \
--header "Authorization: Bearer $DASHSCOPE_API_KEY"

``` 
### [​ ](#generate-silent-videos) Generate silent videos

**Supported models**: `wan2.2 series`, `wan2.1 series`.
**Description**: Ideal for visual-only use cases like animated posters or silent short videos.
**Parameters**: Silent video is the default output for wan2.2 and earlier versions. No extra configuration is needed.
**Input prompt****Output video (silent video)**Low contrast. A street musician performs in a retro 1970s-style subway station, bathed in dim colors and rough textures. He wears a vintage jacket and plays guitar with intense focus. Commuters rush past. A small crowd gradually gathers to listen. The camera pans slowly right, capturing the interplay of instrument sounds and city noise, with vintage subway signs and peeling walls in the background.
- Python SDK 
- Java SDK 
- curl 

 Ensure that the DashScope SDK for Python version is at least `1.25.8`. For instructions on how to update, see [Installing the SDK](/api-reference/preparation/install-sdk). Copy ```\nimport os
from http import HTTPStatus
from dashscope import VideoSynthesis
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

# If you have not set an environment variable, replace the line below with: api_key="sk-xxx"
api_key = os.getenv("DASHSCOPE_API_KEY")

def sample_async_call_t2v():
 # Asynchronous call returns a task_id
 rsp = VideoSynthesis.async_call(api_key=api_key,
 model=&#x27;wan2.2-t2v-plus&#x27;,
 prompt=&#x27;Low contrast. A street musician performs in a retro 1970s-style subway station, bathed in dim colors and rough textures. He wears a vintage jacket and plays guitar with intense focus. Commuters rush past. A small crowd gradually gathers to listen. The camera pans slowly right, capturing the interplay of instrument sounds and city noise, with vintage subway signs and peeling walls in the background.&#x27;,
 prompt_extend=True,
 size=&#x27;832*480&#x27;,
 negative_prompt="",
 watermark=True,
 seed=12345)
 print(rsp)
 if rsp.status_code == HTTPStatus.OK:
 print("task_id: %s" % rsp.output.task_id)
 else:
 print(&#x27;Failed, status_code: %s, code: %s, message: %s&#x27; % (rsp.status_code, rsp.code, rsp.message))

 # Wait for asynchronous task to complete
 rsp = VideoSynthesis.wait(task=rsp, api_key=api_key)
 print(rsp)
 if rsp.status_code == HTTPStatus.OK:
 print(rsp.output.video_url)
 else:
 print(&#x27;Failed, status_code: %s, code: %s, message: %s&#x27; % (rsp.status_code, rsp.code, rsp.message))


if __name__ == &#x27;__main__&#x27;:
 sample_async_call_t2v()

``` Ensure that the DashScope Java SDK version is at least `2.22.6`. To update the SDK, see [Install the SDK](/api-reference/preparation/install-sdk). Copy ```\nimport com.alibaba.dashscope.aigc.videosynthesis.VideoSynthesis;
import com.alibaba.dashscope.aigc.videosynthesis.VideoSynthesisParam;
import com.alibaba.dashscope.aigc.videosynthesis.VideoSynthesisResult;
import com.alibaba.dashscope.exception.ApiException;
import com.alibaba.dashscope.exception.InputRequiredException;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.alibaba.dashscope.utils.JsonUtils;
import com.alibaba.dashscope.utils.Constants;

public class Text2Video {
 static {
 Constants.baseHttpApiUrl = "https://dashscope-intl.aliyuncs.com/api/v1";
 }

 // If you have not set an environment variable, replace the line below with: apiKey="sk-xxx"
 static String apiKey = System.getenv("DASHSCOPE_API_KEY");

 public static void text2video() throws ApiException, NoApiKeyException, InputRequiredException {
 VideoSynthesis vs = new VideoSynthesis();
 VideoSynthesisParam param =
 VideoSynthesisParam.builder()
 .apiKey(apiKey)
 .model("wan2.2-t2v-plus")
 .prompt("Low contrast. A street musician performs in a retro 1970s-style subway station, bathed in dim colors and rough textures. He wears a vintage jacket and plays guitar with intense focus. Commuters rush past. A small crowd gradually gathers to listen. The camera pans slowly right, capturing the interplay of instrument sounds and city noise, with vintage subway signs and peeling walls in the background.")
 .size("832*480")
 .promptExtend(true)
 .watermark(true)
 .seed(12345)
 .build();
 // Asynchronous call
 VideoSynthesisResult task = vs.asyncCall(param);
 System.out.println(JsonUtils.toJson(task));
 System.out.println("please wait...");

 // Get result
 VideoSynthesisResult result = vs.wait(task, apiKey);
 System.out.println(JsonUtils.toJson(result));
 }

 public static void main(String[] args) {
 try {
 text2video();
 } catch (ApiException | NoApiKeyException | InputRequiredException e) {
 System.out.println(e.getMessage());
 }
 System.exit(0);
 }
}

``` **Step 1: Create a task to get the task ID**Copy ```\ncurl --location &#x27;https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis&#x27; \
 -H &#x27;X-DashScope-Async: enable&#x27; \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY" \
 -H &#x27;Content-Type: application/json&#x27; \
 -d &#x27;{
 "model": "wan2.2-t2v-plus",
 "input": {
 "prompt": "Low contrast. A street musician performs in a retro 1970s-style subway station, bathed in dim colors and rough textures. He wears a vintage jacket and plays guitar with intense focus. Commuters rush past. A small crowd gradually gathers to listen. The camera pans slowly right, capturing the interplay of instrument sounds and city noise, with vintage subway signs and peeling walls in the background."
 },
 "parameters": {
 "size": "832*480",
 "prompt_extend": true,
 "watermark": true
 }
}&#x27;

``` **Step 2: Get the result using the task ID**Replace `task_id` with the `task_id` value returned by the previous API call.Copy ```\ncurl -X GET https://dashscope-intl.aliyuncs.com/api/v1/tasks/{task_id} \
--header "Authorization: Bearer $DASHSCOPE_API_KEY"

``` 
## [​ ](#input-audio) Input audio


- **Number of files**: One.

- **Input methods**:

**Public URL**: Supports HTTP or HTTPS protocols.


## [​ ](#output-video) Output video


- **Number of files**: One.

- **Format**: MP4. See [Video generation models](/developer-guides/getting-started/video-models#all-models) for output specifications per model.

- **URL expiration**: **24 hours**.

- **Dimensions**:

**wan2.7**: Set by `resolution` and `ratio`. For example, `resolution=1080P` + `ratio=16:9` outputs a **1920x1080** video.

- **wan2.6 and earlier**: Set by the `size` parameter. For example, `size=1280*720` outputs a **16:9** video.


## [​ ](#billing-and-rate-limits) Billing and rate limits


- For free quota and pricing details, see [Model invocation pricing](/developer-guides/getting-started/pricing).

- For model rate limits, see [Rate limits](/developer-guides/administration/rate-limits).

- Billing details:

Input is free. Output is billed per successfully generated **second of video**.

- Failed model calls or processing errors incur no charge and do not consume your [free quota](/resources/free-quota).


## [​ ](#api-reference) API reference


- [Wan 2.7 text-to-video API reference](/api-reference/video-generation/wan27-text-to-video/create-task)

- [Wan 2.6 text-to-video API reference](/api-reference/video-generation/wan-text-to-video/create-task)


## [​ ](#faq) FAQ

### [​ ](#how-do-i-set-the-video-aspect-ratio-for-example-16-9) How do I set the video aspect ratio (for example, 16:9)?

**wan2.7**: Use the `ratio` parameter directly (e.g., `"16:9"`, `"9:16"`, `"1:1"`, `"4:3"`, `"3:4"`), combined with `resolution` (`"720P"` or `"1080P"`).
**wan2.6 and earlier**: Use the `size` parameter to specify the video resolution in pixels. The system calculates the aspect ratio automatically. For example, `size=1280*720` outputs a **16:9** video.
### [​ ](#sdk-error-url-error-please-check-url) SDK error: "url error, please check url!"

Make sure:

- Your DashScope Python SDK version is at least `1.25.8`.

- Your DashScope Java SDK version is at least `2.22.6`.


If your version is too low, you may see the "url error, please check url!" error. [Upgrade the SDK](/api-reference/preparation/install-sdk).
### [​ ](#why-does-the-call-fail-with-model-not-exist) Why does the call fail with "Model not exist"?

Check these items:

- Is the model name spelled correctly?


For a list of supported models, see [Video generation models](/developer-guides/getting-started/video-models). [Previous ](/developer-guides/getting-started/video-models)[Image-to-video: first frame Animate from a single image Next ](/developer-guides/video-generation/image-to-video)
