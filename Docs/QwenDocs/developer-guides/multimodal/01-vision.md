# Analyze images and videos

> **Source:** https://docs.qwencloud.com/developer-guides/multimodal/vision

Generate content from visual inputs

 Copy page Vision models understand images and videos to answer questions, extract text, solve problems, and generate descriptions. These multimodal models combine visual understanding with language capabilities for tasks ranging from OCR to creative writing.
## [​ ](#visual-input-structure) Visual input structure

Vision models accept images and videos alongside text prompts. Each message can contain multiple content types:

- **Text prompt**: Your question or instruction about the visual content

- **Image URL**: Direct link to an online image

- **Base64 image**: Encoded image data for local files

- **Video URL**: Direct link to video content (select models)


## [​ ](#make-your-first-vision-call) Make your first vision call

**Prerequisites**

- [Get an API key](/api-reference/preparation/api-key) and [export it as an environment variable](/api-reference/preparation/install-sdk)

- Install the [SDK](/api-reference/preparation/install-sdk) if using one (Python SDK 1.24.6+, Java SDK 2.21.10+)


**Which API to use?**

- **OpenAI Compatible**: Best for new integrations and migrating from OpenAI

- **DashScope**: Use if you prefer the native SDK or need specific DashScope features


- OpenAI compatible 
- DashScope 

 Python Node.js curl Copy ```\nfrom openai import OpenAI
import os

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1"
)

completion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=[
 {
 "role": "user",
 "content": [
 {
 "type": "image_url",
 "image_url": {
 "url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241022/emyrja/dog_and_girl.jpeg"
 },
 },
 {"type": "text", "text": "Describe what you see in this image"},
 ],
 },
 ],
)
print(completion.choices[0].message.content)

``` **Response**Copy ```\nThis is a photo taken on a beach. In the photo, a person and a dog are sitting on the sand, with the sea and sky in the background. The person and dog appear to be interacting, with the dog&#x27;s front paw resting on the person&#x27;s hand. Sunlight is coming from the right side of the frame, adding a warm atmosphere to the scene.

``` Full JSON response

 Copy ```\n{
 "choices": [
 {
 "message": {
 "content": "This image depicts a heartwarming scene on a sandy beach...",
 "reasoning_content": "The user wants a description of the image.\n\n1. **Identify the main subjects:** A woman and a dog...",
 "role": "assistant"
 },
 "finish_reason": "stop",
 "index": 0,
 "logprobs": null
 }
 ],
 "object": "chat.completion",
 "usage": {
 "prompt_tokens": 2520,
 "completion_tokens": 777,
 "total_tokens": 3297,
 "completion_tokens_details": {
 "reasoning_tokens": 539,
 "text_tokens": 238
 },
 "prompt_tokens_details": {
 "image_tokens": 2503,
 "text_tokens": 17
 }
 },
 "created": 1774322504,
 "system_fingerprint": null,
 "model": "qwen3.7-plus",
 "id": "chatcmpl-be9bf2d1-2e70-91c4-b8bc-c7f5bbd30320"
}

``` Python Java curl Copy ```\nimport os
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

messages = [
{
 "role": "user",
 "content": [
 {"image": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241022/emyrja/dog_and_girl.jpeg"},
 {"text": "Describe what you see in this image"}]
}]

response = dashscope.MultiModalConversation.call(
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 model=&#x27;qwen3.7-plus&#x27;,
 messages=messages
)

print(response.output.choices[0].message.content[0]["text"])

``` **Response**Copy ```\nThis is a photo taken on a beach. In the photo, there is a woman and a dog. The woman is sitting on the sand, smiling and interacting with the dog. The dog is wearing a collar and appears to be shaking hands with the woman. The background is the sea and the sky, and the sunlight shining on them creates a warm atmosphere.

``` Full JSON response

 Copy ```\n{
 "output": {
 "choices": [
 {
 "finish_reason": "stop",
 "message": {
 "role": "assistant",
 "content": [
 {
 "text": "This is a photo taken on a beach. In the photo, there is a person in a plaid shirt and a dog with a collar. They are sitting on the sand, with the sea and sky in the background. Sunlight is coming from the right side of the frame, adding a warm atmosphere to the scene."
 }
 ]
 }
 }
 ]
 },
 "usage": {
 "output_tokens": 55,
 "input_tokens": 1271,
 "image_tokens": 1247
 },
 "request_id": "ccf845a3-dc33-9cda-b581-20fe7dc23f70"
}

``` 
## [​ ](#compare-model-performance) Compare model performance

### [​ ](#answer-questions-about-images) Answer questions about images

Describe the content of an image or classify and label it, such as identifying people, places, animals, and plants.
InputOutput If the sun is glaring, what item from this picture should I use?When the sun is glaring, you should use the pink sunglasses from the picture. Sunglasses can effectively block strong light, reduce UV damage to your eyes, and help protect your vision and improve visual comfort in bright sunlight.
### [​ ](#generate-creative-content-from-images) Generate creative content from images

Generate vivid text descriptions based on image or video content. This is suitable for creative scenarios such as story writing, copywriting, and short video scripts.
InputOutput Please help me write an interesting social media post based on the content of the picture.Merry Christmas from our little winter wonderland! We&#x27;re getting ready for the holidays with warm lights, pinecones, and plenty of rustic charm. Hope your season is filled with this much warmth and joy!
### [​ ](#extract-text-and-information) Extract text and information

Recognize text and formulas in images or extract information from receipts, certificates, and forms, with support for formatted text output. The Qwen3-VL model has expanded its language support to 33 languages. For a list of supported languages, see [Vision models](/developer-guides/getting-started/vision-models).
InputOutput Extract the following from the image: [&#x27;Invoice Code&#x27;, &#x27;Invoice Number&#x27;, &#x27;Destination&#x27;, &#x27;Fuel Surcharge&#x27;, &#x27;Fare&#x27;, &#x27;Travel Date&#x27;, &#x27;Departure Time&#x27;, &#x27;Train Number&#x27;, &#x27;Seat Number&#x27;]. Please output in JSON format.{`{"Invoice Code": "221021325353", "Invoice Number": "10283819", "Destination": "Development Zone", "Fuel Surcharge": "2.0", "Fare": "8.00<Full>", "Travel Date": "2013-06-29", "Departure Time": "Serial", "Train Number": "040", "Seat Number": "371"}`}
### [​ ](#solve-complex-visual-problems) Solve complex visual problems

Solve problems in images, such as math, physics, and chemistry problems. This feature is suitable for primary, secondary, university, and adult education.
InputOutput Please solve the math problem in the image step by step. 
### [​ ](#generate-code-from-visual-designs) Generate code from visual designs

Generate code from images or videos. This can be used to create HTML, CSS, and JS code from design drafts, website screenshots, and more.
InputOutput Design a webpage using HTML and CSS based on my sketch, with black as the main color. **Webpage preview**
### [​ ](#locate-objects-in-images) Locate objects in images

The model supports 2D and 3D localization to determine object orientation, perspective changes, and occlusion relationships. Qwen3-VL adds 3D localization.
 For Qwen2.5-VL, object detection is robust within 480x480 to 2560x2560 resolution. Outside this range, accuracy may decrease with occasional bounding box drift. To draw localization results on the original image, see [FAQ](#faq). 
InputOutput**2D localization** 
- Return Box (bounding box) coordinates: Detect all food items in the image and output their bbox coordinates in JSON format.

- Return Point (centroid) coordinates: Locate all food items in the image as points and output their point coordinates in XML format.


**Visualization of 2D localization results** 
 3D localization

 Detect the car in the image and predict its 3D position. Output JSON: `[{"bbox_3d": [x_center, y_center, z_center, x_size, y_size, z_size, roll, pitch, yaw], "label": "category"}]`. 
 Visualization of 3D localization results

 
### [​ ](#parse-documents-and-pdfs) Parse documents and PDFs

Parse image-based documents, such as scans or image PDFs, into QwenVL HTML or QwenVL Markdown format. This format not only accurately recognizes text but also obtains the position information of elements such as images and tables. The Qwen3-VL model adds the ability to parse documents into Markdown format.
 Recommended prompts are: `qwenvl html` (to parse into HTML format) or `qwenvl markdown` (to parse into Markdown format). 
InputOutput qwenvl markdown. **Visualization of results**
### [​ ](#analyze-video-content) Analyze video content

Analyze video content, such as locating specific events and obtaining timestamps, or generating summaries of key time periods.
InputOutputPlease describe the series of actions of the person in the video. Output in JSON format with start_time, end_time, and event. Use HH:mm:ss for timestamps.{`{"events": [{"start_time": "00:00:00", "end_time": "00:00:05", "event": "The person walks towards the table holding a cardboard box and places it on the table."}, {"start_time": "00:00:05", "end_time": "00:00:15", "event": "The person picks up a scanner and scans the label on the cardboard box."}, {"start_time": "00:00:15", "end_time": "00:00:21", "event": "The person puts the scanner back in its place and then picks up a pen to write information in a notebook."}]}`}
## [​ ](#work-with-visual-content) Work with visual content

### [​ ](#thinking-mode) Thinking mode

 For enable/disable, streaming output, and `thinking_budget`, see [Thinking](/developer-guides/text-generation/thinking). 
Vision defaults: thinking is **off** for `qwen3-vl-plus` and `qwen3-vl-flash`, **on** for `qwen3.6` and `qwen3.5`. Models with a `-thinking` suffix always think.
### [​ ](#work-with-multiple-images) Work with multiple images

Pass multiple images in a single request for tasks like **product comparison and multi-page document processing**. Include multiple image objects in the `user message`&#x27;s `content` array.
 Per request: up to **256** images when passed as a **public URL** or **local file path**, and up to **250** images when passed as **Base64-encoded** images. Independently, total tokens for all images and text must stay **below** the model&#x27;s maximum input (combined image-and-text token limit). 
- OpenAI compatible 
- DashScope 

 - Python 
- Node.js 
- curl 

 Copy ```\nimport os
from openai import OpenAI

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)

completion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=[
 {"role": "user","content": [
 {"type": "image_url","image_url": {"url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241022/emyrja/dog_and_girl.jpeg"},},
 {"type": "image_url","image_url": {"url": "https://dashscope.oss-cn-beijing.aliyuncs.com/images/tiger.png"},},
 {"type": "text", "text": "What do these images depict?"},
 ],
 }
 ],
)

print(completion.choices[0].message.content)

``` **Response**Copy ```\nImage 1 shows a scene of a woman and a Labrador retriever interacting on a beach. The woman is wearing a plaid shirt and sitting on the sand, shaking hands with the dog. The background is ocean waves and the sky, and the whole picture is filled with a warm and pleasant atmosphere.

Image 2 shows a scene of a tiger walking in a forest. The tiger&#x27;s coat is orange with black stripes, and it is stepping forward. The surroundings are dense trees and vegetation, and the ground is covered with fallen leaves. The whole picture gives a feeling of wild nature.

``` Copy ```\nimport OpenAI from "openai";

const openai = new OpenAI(
 {
 apiKey: process.env.DASHSCOPE_API_KEY,
 baseURL: "https://dashscope-intl.aliyuncs.com/compatible-mode/v1"
 }
);

async function main() {
 const response = await openai.chat.completions.create({
 model: "qwen3.7-plus",
 messages: [
 {role: "user",content: [
 {type: "image_url",image_url: {"url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241022/emyrja/dog_and_girl.jpeg"}},
 {type: "image_url",image_url: {"url": "https://dashscope.oss-cn-beijing.aliyuncs.com/images/tiger.png"}},
 {type: "text", text: "What do these images depict?" },
 ]}]
 });
 console.log(response.choices[0].message.content);
}

main()

``` **Response**Copy ```\nIn the first image, a person and a dog are interacting on a beach. The person is wearing a plaid shirt, and the dog is wearing a collar. They seem to be shaking hands or giving a high-five.

In the second image, a tiger is walking in a forest. The tiger&#x27;s coat is orange with black stripes, and the background is green trees and vegetation.

``` Copy ```\ncurl -X POST https://dashscope-intl.aliyuncs.com/compatible-mode/v1/chat/completions \
-H "Authorization: Bearer $DASHSCOPE_API_KEY" \
-H &#x27;Content-Type: application/json&#x27; \
-d &#x27;{
 "model": "qwen3.7-plus",
 "messages": [
 {
 "role": "user",
 "content": [
 {
 "type": "image_url",
 "image_url": {
 "url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241022/emyrja/dog_and_girl.jpeg"
 }
 },
 {
 "type": "image_url",
 "image_url": {
 "url": "https://dashscope.oss-cn-beijing.aliyuncs.com/images/tiger.png"
 }
 },
 {
 "type": "text",
 "text": "What do these images depict?"
 }
 ]
 }
 ]
}&#x27;

``` **Response**Copy ```\nImage 1 shows a scene of a woman and a Labrador retriever interacting on a beach. The woman is wearing a plaid shirt and sitting on the sand, shaking hands with the dog. The background is ocean views and a sunset sky, and the whole picture looks very warm and harmonious.

Image 2 shows a scene of a tiger walking in a forest. The tiger&#x27;s coat is orange with black stripes, and it is stepping forward. The surroundings are dense trees and vegetation, and the ground is covered with fallen leaves. The whole picture is full of natural wildness and vitality.

``` Full JSON response

 Copy ```\n{
 "choices": [
 {
 "message": {
 "content": "Image 1 shows a scene of a woman and a Labrador retriever interacting on a beach. The woman is wearing a plaid shirt and sitting on the sand, shaking hands with the dog. The background is ocean views and a sunset sky, and the whole picture looks very warm and harmonious.\n\nImage 2 shows a scene of a tiger walking in a forest. The tiger&#x27;s coat is orange with black stripes, and it is stepping forward. The surroundings are dense trees and vegetation, and the ground is covered with fallen leaves. The whole picture is full of natural wildness and vitality.",
 "role": "assistant"
 },
 "finish_reason": "stop",
 "index": 0,
 "logprobs": null
 }
 ],
 "object": "chat.completion",
 "usage": {
 "prompt_tokens": 2497,
 "completion_tokens": 109,
 "total_tokens": 2606
 },
 "created": 1725948561,
 "system_fingerprint": null,
 "model": "qwen3.7-plus",
 "id": "chatcmpl-0fd66f46-b09e-9164-a84f-3ebbbedbac15"
}

``` - Python 
- Java 
- curl 

 Copy ```\nimport os
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

messages = [
 {
 "role": "user",
 "content": [
 {"image": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241022/emyrja/dog_and_girl.jpeg"},
 {"image": "https://dashscope.oss-cn-beijing.aliyuncs.com/images/tiger.png"},
 {"text": "What do these images depict?"}
 ]
 }
]

response = dashscope.MultiModalConversation.call(
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 model=&#x27;qwen3.7-plus&#x27;,
 messages=messages
)

print(response.output.choices[0].message.content[0]["text"])

``` **Response**Copy ```\nThese images show some animals and natural scenes. In the first image, a person and a dog are interacting on a beach. The second image is of a tiger walking in a forest.

``` Copy ```\nimport java.util.Arrays;
import java.util.Collections;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversation;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationParam;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationResult;
import com.alibaba.dashscope.common.MultiModalMessage;
import com.alibaba.dashscope.common.Role;
import com.alibaba.dashscope.exception.ApiException;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.alibaba.dashscope.exception.UploadFileException;
import com.alibaba.dashscope.utils.Constants;

public class Main {
 static {
 Constants.baseHttpApiUrl="https://dashscope-intl.aliyuncs.com/api/v1";
 }
 public static void simpleMultiModalConversationCall()
 throws ApiException, NoApiKeyException, UploadFileException {
 MultiModalConversation conv = new MultiModalConversation();
 MultiModalMessage userMessage = MultiModalMessage.builder().role(Role.USER.getValue())
 .content(Arrays.asList(
 Collections.singletonMap("image", "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241022/emyrja/dog_and_girl.jpeg"),
 Collections.singletonMap("image", "https://dashscope.oss-cn-beijing.aliyuncs.com/images/tiger.png"),
 Collections.singletonMap("text", "What do these images depict?"))).build();
 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model("qwen3.7-plus")
 .messages(Arrays.asList(userMessage))
 .build();
 MultiModalConversationResult result = conv.call(param);
 System.out.println(result.getOutput().getChoices().get(0).getMessage().getContent().get(0).get("text")); }
 public static void main(String[] args) {
 try {
 simpleMultiModalConversationCall();
 } catch (ApiException | NoApiKeyException | UploadFileException e) {
 System.out.println(e.getMessage());
 }
 System.exit(0);
 }
}

``` **Response**Copy ```\nThese images show some animals and natural scenes.

1. First image: A woman and a dog are interacting on a beach. The woman is wearing a plaid shirt and sitting on the sand, and the dog is wearing a collar and extending its paw to shake hands with the woman.
2. Second image: A tiger is walking in a forest. The tiger&#x27;s coat is orange with black stripes, and the background is trees and leaves.

``` Copy ```\ncurl --location &#x27;https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation&#x27; \
--header "Authorization: Bearer $DASHSCOPE_API_KEY" \
--header &#x27;Content-Type: application/json&#x27; \
--data &#x27;{
 "model": "qwen3.7-plus",
 "input":{
 "messages":[
 {
 "role": "user",
 "content": [
 {"image": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241022/emyrja/dog_and_girl.jpeg"},
 {"image": "https://dashscope.oss-cn-beijing.aliyuncs.com/images/tiger.png"},
 {"text": "What do these images show?"}
 ]
 }
 ]
 }
}&#x27;

``` **Response**Copy ```\nThese images show some animals and natural scenes. In the first image, a person and a dog are interacting on a beach. The second image is of a tiger walking in a forest.

``` Full JSON response

 Copy ```\n{
 "output": {
 "choices": [
 {
 "finish_reason": "stop",
 "message": {
 "role": "assistant",
 "content": [
 {
 "text": "These images show some animals and natural scenes. In the first image, a person and a dog are interacting on a beach. The second image is of a tiger walking in a forest."
 }
 ]
 }
 }
 ]
 },
 "usage": {
 "output_tokens": 81,
 "input_tokens": 1277,
 "image_tokens": 2497
 },
 "request_id": "ccf845a3-dc33-9cda-b581-20fe7dc23f70"
}

``` 
### [​ ](#analyze-video-content-2) Analyze video content

Visual understanding models support understanding video content. You can provide files in the form of an image list (video frames) or a video file. The following is an example of code for understanding an online video or image list specified by a URL. For more information about video limits or the number of images that can be passed in an image list, see the [Video limits](#input-file-limits) section.
 We recommend using the latest or a recent snapshot version of the model for better performance in understanding video files. 
- Video file 
- Image list 

 Visual understanding models analyze content by extracting a sequence of frames from a video. You can control the frame extraction policy using the following two parameters:
- **fps**: Controls the frequency. One frame every `1/fps` seconds. The value range is [0.1, 10] and the default value is 2.0.

High-speed motion scenes: Set a higher fps value to capture more detail.

- Static or long videos: Set a lower fps value for efficiency.


- **max_frames:** The upper limit of frames extracted. When the number calculated based on fps exceeds max_frames, the system automatically and evenly samples frames to stay within the limit. **This parameter is active only for the DashScope SDK.**


- OpenAI compatible 
- DashScope 

 When you directly input a video file to a visual understanding model using the OpenAI SDK or HTTP method, you must set the `"type"` parameter in the user message to `"video_url"`. - Python 
- Node.js 
- curl 

 Copy ```\nimport os
from openai import OpenAI

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)
completion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=[
 {
 "role": "user",
 "content": [
 {
 "type": "video_url",
 "video_url": {
 "url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241115/cqqkru/1.mp4"
 },
 "fps": 2
 },
 {
 "type": "text",
 "text": "Summarize what happens in this video"
 }
 ]
 }
 ]
)

print(completion.choices[0].message.content)

``` Copy ```\nimport OpenAI from "openai";

const openai = new OpenAI(
 {
 apiKey: process.env.DASHSCOPE_API_KEY,
 baseURL: "https://dashscope-intl.aliyuncs.com/compatible-mode/v1"
 }
);

async function main() {
 const response = await openai.chat.completions.create({
 model: "qwen3.7-plus",
 messages: [
 {
 role: "user",
 content: [
 {
 type: "video_url",
 video_url: {
 "url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241115/cqqkru/1.mp4"
 },
 "fps": 2
 },
 {
 type: "text",
 text: "Summarize what happens in this video"
 }
 ]
 }
 ]
 });

 console.log(response.choices[0].message.content);
}

main();

``` Copy ```\ncurl -X POST https://dashscope-intl.aliyuncs.com/compatible-mode/v1/chat/completions \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY" \
 -H &#x27;Content-Type: application/json&#x27; \
 -d &#x27;{
 "model": "qwen3.7-plus",
 "messages": [
 {
 "role": "user",
 "content": [
 {
 "type": "video_url",
 "video_url": {
 "url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241115/cqqkru/1.mp4"
 },
 "fps":2
 },
 {
 "type": "text",
 "text": "Summarize what happens in this video"
 }
 ]
 }
 ]
 }&#x27;

``` - Python 
- Java 
- curl 

 Copy ```\nimport dashscope
import os

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;
messages = [
 {"role": "user",
 "content": [
 {"video": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241115/cqqkru/1.mp4","fps":2},
 {"text": "Summarize what happens in this video"}
 ]
 }
]

response = dashscope.MultiModalConversation.call(
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 model=&#x27;qwen3.7-plus&#x27;,
 messages=messages
)

print(response.output.choices[0].message.content[0]["text"])

``` Copy ```\nimport java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;

import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversation;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationParam;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationResult;
import com.alibaba.dashscope.common.MultiModalMessage;
import com.alibaba.dashscope.common.Role;
import com.alibaba.dashscope.exception.ApiException;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.alibaba.dashscope.exception.UploadFileException;
import com.alibaba.dashscope.utils.Constants;

public class Main {
 static {
 Constants.baseHttpApiUrl="https://dashscope-intl.aliyuncs.com/api/v1";
 }
 public static void simpleMultiModalConversationCall()
 throws ApiException, NoApiKeyException, UploadFileException {
 MultiModalConversation conv = new MultiModalConversation();
 Map<String, Object> params = new HashMap<>();
 params.put("video", "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241115/cqqkru/1.mp4");
 params.put("fps", 2);
 MultiModalMessage userMessage = MultiModalMessage.builder().role(Role.USER.getValue())
 .content(Arrays.asList(
 params,
 Collections.singletonMap("text", "Summarize what happens in this video"))).build();
 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model("qwen3.7-plus")
 .messages(Arrays.asList(userMessage))
 .build();
 MultiModalConversationResult result = conv.call(param);
 System.out.println(result.getOutput().getChoices().get(0).getMessage().getContent().get(0).get("text"));
 }
 public static void main(String[] args) {
 try {
 simpleMultiModalConversationCall();
 } catch (ApiException | NoApiKeyException | UploadFileException e) {
 System.out.println(e.getMessage());
 }
 System.exit(0);
 }
}

``` Copy ```\ncurl -X POST https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation \
-H "Authorization: Bearer $DASHSCOPE_API_KEY" \
-H &#x27;Content-Type: application/json&#x27; \
-d &#x27;{
 "model": "qwen3.7-plus",
 "input":{
 "messages":[
 {"role": "user","content": [{"video": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241115/cqqkru/1.mp4","fps":2},
 {"text": "Summarize what happens in this video"}]}]}
}&#x27;

``` When a video is passed as a list of images (pre-extracted video frames), you can use the `fps` parameter to inform the model of the time interval between video frames. This helps the model better understand the sequence, duration, and dynamic changes of events. The model supports specifying the original video&#x27;s frame rate using the `fps` parameter, which indicates that the video frames were extracted from the original video every `1/fps` seconds. This parameter is supported by **Qwen3.6**, **Qwen3-VL**, and **Qwen2.5-VL** models.- OpenAI compatible 
- DashScope 

 When you input a video as a list of images to a visual understanding model using the OpenAI SDK or an HTTP request, you must set the `"type"` parameter in the user message to `"video"`. - Python 
- Node.js 
- curl 

 Copy ```\nimport os
from openai import OpenAI

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)

completion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=[{"role": "user","content": [
 {"type": "video","video": [
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/xzsgiz/football1.jpg",
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/tdescd/football2.jpg",
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/zefdja/football3.jpg",
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/aedbqh/football4.jpg"],
 "fps":2},
 {"type": "text","text": "Describe the specific process of this video"},
 ]}]
)

print(completion.choices[0].message.content)

``` Copy ```\nimport OpenAI from "openai";

const openai = new OpenAI({
 apiKey: process.env.DASHSCOPE_API_KEY,
 baseURL: "https://dashscope-intl.aliyuncs.com/compatible-mode/v1"
});

async function main() {
 const response = await openai.chat.completions.create({
 model: "qwen3.7-plus",
 messages: [{
 role: "user",
 content: [
 {
 type: "video",
 video: [
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/xzsgiz/football1.jpg",
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/tdescd/football2.jpg",
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/zefdja/football3.jpg",
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/aedbqh/football4.jpg"],
 "fps": 2
 },
 {
 type: "text",
 text: "Describe the specific process of this video"
 }
 ]
 }]
 });
 console.log(response.choices[0].message.content);
}

main();

``` Copy ```\ncurl -X POST https://dashscope-intl.aliyuncs.com/compatible-mode/v1/chat/completions \
-H "Authorization: Bearer $DASHSCOPE_API_KEY" \
-H &#x27;Content-Type: application/json&#x27; \
-d &#x27;{
 "model": "qwen3.7-plus",
 "messages": [{"role": "user","content": [{"type": "video","video": [
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/xzsgiz/football1.jpg",
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/tdescd/football2.jpg",
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/zefdja/football3.jpg",
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/aedbqh/football4.jpg"],
 "fps":2},
 {"type": "text","text": "Describe the specific process of this video"}]}]
}&#x27;

``` - Python 
- Java 
- curl 

 Copy ```\nimport os
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;
messages = [{"role": "user",
 "content": [
 {"video":["https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/xzsgiz/football1.jpg",
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/tdescd/football2.jpg",
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/zefdja/football3.jpg",
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/aedbqh/football4.jpg"],
 "fps":2},
 {"text": "Describe the specific process of this video"}]}]
response = dashscope.MultiModalConversation.call(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 model=&#x27;qwen3.7-plus&#x27;,
 messages=messages
)
print(response.output.choices[0].message.content[0]["text"])

``` Copy ```\n// DashScope SDK version must be 2.21.10 or later
import java.util.Arrays;
import java.util.Collections;
import java.util.Map;

import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversation;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationParam;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationResult;
import com.alibaba.dashscope.common.MultiModalMessage;
import com.alibaba.dashscope.common.Role;
import com.alibaba.dashscope.exception.ApiException;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.alibaba.dashscope.exception.UploadFileException;
import com.alibaba.dashscope.utils.Constants;

public class Main {
 static {
 Constants.baseHttpApiUrl="https://dashscope-intl.aliyuncs.com/api/v1";
 }
 private static final String MODEL_NAME = "qwen3.7-plus";
 public static void videoImageListSample() throws ApiException, NoApiKeyException, UploadFileException {
 MultiModalConversation conv = new MultiModalConversation();
 Map<String, Object> params = new HashMap<>();
 params.put("video", Arrays.asList("https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/xzsgiz/football1.jpg",
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/tdescd/football2.jpg",
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/zefdja/football3.jpg",
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/aedbqh/football4.jpg"));
 params.put("fps", 2);
 MultiModalMessage userMessage = MultiModalMessage.builder()
 .role(Role.USER.getValue())
 .content(Arrays.asList(
 params,
 Collections.singletonMap("text", "Describe the specific process of this video")))
 .build();
 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model(MODEL_NAME)
 .messages(Arrays.asList(userMessage)).build();
 MultiModalConversationResult result = conv.call(param);
 System.out.print(result.getOutput().getChoices().get(0).getMessage().getContent().get(0).get("text"));
 }
 public static void main(String[] args) {
 try {
 videoImageListSample();
 } catch (ApiException | NoApiKeyException | UploadFileException e) {
 System.out.println(e.getMessage());
 }
 System.exit(0);
 }
}

``` Copy ```\ncurl -X POST https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation \
-H "Authorization: Bearer $DASHSCOPE_API_KEY" \
-H &#x27;Content-Type: application/json&#x27; \
-d &#x27;{
 "model": "qwen3.7-plus",
 "input": {
 "messages": [
 {
 "role": "user",
 "content": [
 {
 "video": [
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/xzsgiz/football1.jpg",
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/tdescd/football2.jpg",
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/zefdja/football3.jpg",
 "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/aedbqh/football4.jpg"
 ],
 "fps":2
 },
 {
 "text": "Describe the specific process of this video"
 }
 ]
 }
 ]
 }
}&#x27;

``` 
### [​ ](#use-local-files) Use local files

Visual understanding models provide two ways to upload local files: Base64 encoding and direct file path upload. Choose the upload method based on the file size and SDK type. For specific recommendations, see [How to choose a file upload method](#faq). Both methods must meet the file requirements described in [Image limits](#input-file-limits).
- Base64 encoding upload 
- File path upload 

 Convert the file to a Base64 encoded string and then pass it to the model. This method is supported by the OpenAI and DashScope SDKs, and HTTP requests. Steps to pass a Base64-encoded string (image example)

 1 File encoding

Convert the local image to a Base64 encoding. Example code for converting an image to Base64 encoding

 Copy ```\n# Encoding function: Converts a local file to a Base64 encoded string
import base64
def encode_image(image_path):
 with open(image_path, "rb") as image_file:
 return base64.b64encode(image_file.read()).decode("utf-8")

# Replace xxxx/eagle.png with the absolute path of your local image
base64_image = encode_image("xxx/eagle.png")

``` 2 Construct a Data URL

The format is as follows: `data:[MIME_type];base64,<base64_image>`.
- Replace `MIME_type` with the actual media type. Ensure it matches the `MIME type` value in the [Supported image formats](#input-file-limits) table, such as `image/jpeg` or `image/png`.

- `base64_image` is the Base64 string generated in the previous step.


 3 Call the model

Pass the `Data URL` through the `image` or `image_url` parameter and call the model. Pass the local file path directly to the model. This method is supported only by the DashScope Python and Java SDKs, not by DashScope HTTP or OpenAI-compatible methods.Refer to the table below to specify the file path based on your programming language and operating system. Specify the file path (using an image as an example)

 **System****SDK****Input file path****Example**Linux or macOSPython SDK`file://<absolute path of the file>``file:///home/images/test.png`Linux or macOSJava SDK`file://<absolute path of the file>``file:///home/images/test.png`WindowsPython SDK`file://<absolute path of the file>``file://D:/images/test.png`WindowsJava SDK`file:///<absolute path of the file>``file:///D:/images/test.png` 
The code examples below show how to pass local images, videos, and image lists using both Base64 encoding and file path methods. Due to the large number of examples, they are organized by file type.
 Image - Pass by file path (DashScope only)

 - Python 
- Java 

 Copy ```\nimport os
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

# Replace xxx/eagle.png with the absolute path of your local image
local_path = "xxx/eagle.png"
image_path = f"file://{local_path}"
messages = [
 {&#x27;role&#x27;:&#x27;user&#x27;,
 &#x27;content&#x27;: [{&#x27;image&#x27;: image_path},
 {&#x27;text&#x27;: &#x27;Describe what you see in this image&#x27;}]}]
response = dashscope.MultiModalConversation.call(
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 model=&#x27;qwen3.7-plus&#x27;,
 messages=messages)
print(response.output.choices[0].message.content[0]["text"])

``` Copy ```\nimport java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversation;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationParam;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationResult;
import com.alibaba.dashscope.common.MultiModalMessage;
import com.alibaba.dashscope.common.Role;
import com.alibaba.dashscope.exception.ApiException;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.alibaba.dashscope.exception.UploadFileException;
import com.alibaba.dashscope.utils.Constants;

public class Main {

 static {
 Constants.baseHttpApiUrl="https://dashscope-intl.aliyuncs.com/api/v1";
 }
 
 public static void callWithLocalFile(String localPath)
 throws ApiException, NoApiKeyException, UploadFileException {
 String filePath = "file://"+localPath;
 MultiModalConversation conv = new MultiModalConversation();
 MultiModalMessage userMessage = MultiModalMessage.builder().role(Role.USER.getValue())
 .content(Arrays.asList(new HashMap<String, Object>(){{put("image", filePath);}},
 new HashMap<String, Object>(){{put("text", "Describe what you see in this image");}})).build();
 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model("qwen3.7-plus")
 .messages(Arrays.asList(userMessage))
 .build();
 MultiModalConversationResult result = conv.call(param);
 System.out.println(result.getOutput().getChoices().get(0).getMessage().getContent().get(0).get("text"));}

 public static void main(String[] args) {
 try {
 // Replace xxx/eagle.png with the absolute path of your local image
 callWithLocalFile("xxx/eagle.png");
 } catch (ApiException | NoApiKeyException | UploadFileException e) {
 System.out.println(e.getMessage());
 }
 System.exit(0);
 }
}

``` 
 Image - Pass in Base64 encoding

 - OpenAI compatible 
- DashScope 

 - Python 
- Node.js 
- curl 

 Copy ```\nfrom openai import OpenAI
import os
import base64

def encode_image(image_path):
 with open(image_path, "rb") as image_file:
 return base64.b64encode(image_file.read()).decode("utf-8")

base64_image = encode_image("xxx/eagle.png")
client = OpenAI(
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)
completion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=[
 {
 "role": "user",
 "content": [
 {
 "type": "image_url",
 "image_url": {"url": f"data:image/png;base64,{base64_image}"},
 },
 {"type": "text", "text": "Describe what you see in this image"},
 ],
 }
 ],
)
print(completion.choices[0].message.content)

``` Copy ```\nimport OpenAI from "openai";
import { readFileSync } from &#x27;fs&#x27;;

const openai = new OpenAI(
 {
 apiKey: process.env.DASHSCOPE_API_KEY,
 baseURL: "https://dashscope-intl.aliyuncs.com/compatible-mode/v1"
 }
);

const encodeImage = (imagePath) => {
 const imageFile = readFileSync(imagePath);
 return imageFile.toString(&#x27;base64&#x27;);
 };
const base64Image = encodeImage("xxx/eagle.png")
async function main() {
 const completion = await openai.chat.completions.create({
 model: "qwen3.7-plus",
 messages: [
 {"role": "user",
 "content": [{"type": "image_url",
 "image_url": {"url": `data:image/png;base64,${base64Image}`},},
 {"type": "text", "text": "Describe what you see in this image"}]}]
 });
 console.log(completion.choices[0].message.content);
} 

main();

``` 
- For display purposes, the Base64 encoded string in the code is truncated. You must pass the complete encoded string.


Copy ```\ncurl --location &#x27;https://dashscope-intl.aliyuncs.com/compatible-mode/v1/chat/completions&#x27; \
--header "Authorization: Bearer $DASHSCOPE_API_KEY" \
--header &#x27;Content-Type: application/json&#x27; \
--data &#x27;{
 "model": "qwen3.7-plus",
 "messages": [
 {
 "role": "user",
 "content": [
 {"type": "image_url", "image_url": {"url": "data:image/jpg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAA"}},
 {"type": "text", "text": "Describe what you see in this image"}
 ]
 }]
}&#x27;

``` - Python 
- Java 
- curl 

 Copy ```\nimport base64
import os
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

def encode_image(image_path):
 with open(image_path, "rb") as image_file:
 return base64.b64encode(image_file.read()).decode("utf-8")

base64_image = encode_image("xxxx/eagle.png")

messages = [
 {
 "role": "user",
 "content": [
 {"image": f"data:image/png;base64,{base64_image}"},
 {"text": "Describe what you see in this image"},
 ],
 },
]

response = dashscope.MultiModalConversation.call(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 model="qwen3.7-plus",
 messages=messages,
)
print(response.output.choices[0].message.content[0]["text"])

``` Copy ```\nimport java.io.IOException;
import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.Base64;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

import com.alibaba.dashscope.aigc.multimodalconversation.*;
import com.alibaba.dashscope.common.MultiModalMessage;
import com.alibaba.dashscope.common.Role;
import com.alibaba.dashscope.exception.ApiException;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.alibaba.dashscope.exception.UploadFileException;
import com.alibaba.dashscope.utils.Constants;

public class Main {

 static {
 Constants.baseHttpApiUrl="https://dashscope-intl.aliyuncs.com/api/v1";
 }

 private static String encodeImageToBase64(String imagePath) throws IOException {
 Path path = Paths.get(imagePath);
 byte[] imageBytes = Files.readAllBytes(path);
 return Base64.getEncoder().encodeToString(imageBytes);
 }

 public static void callWithLocalFile(String localPath) throws ApiException, NoApiKeyException, UploadFileException, IOException {

 String base64Image = encodeImageToBase64(localPath);

 MultiModalConversation conv = new MultiModalConversation();
 MultiModalMessage userMessage = MultiModalMessage.builder().role(Role.USER.getValue())
 .content(Arrays.asList(
 new HashMap<String, Object>() {{ put("image", "data:image/png;base64," + base64Image); }},
 new HashMap<String, Object>() {{ put("text", "Describe what you see in this image"); }}
 )).build();

 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model("qwen3.7-plus")
 .messages(Arrays.asList(userMessage))
 .build();

 MultiModalConversationResult result = conv.call(param);
 System.out.println(result.getOutput().getChoices().get(0).getMessage().getContent().get(0).get("text"));
 }

 public static void main(String[] args) {
 try {
 callWithLocalFile("xxx/eagle.png");
 } catch (ApiException | NoApiKeyException | UploadFileException | IOException e) {
 System.out.println(e.getMessage());
 }
 System.exit(0);
 }
}

``` 
- For display purposes, the Base64 encoded string in the code is truncated. You must pass the complete encoded string.


Copy ```\ncurl -X POST https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation \
-H "Authorization: Bearer $DASHSCOPE_API_KEY" \
-H &#x27;Content-Type: application/json&#x27; \
-d &#x27;{
 "model": "qwen3.7-plus",
 "input":{
 "messages":[
 {
 "role": "user",
 "content": [
 {"image": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAA..."},
 {"text": "Describe what you see in this image"}
 ]
 }
 ]
 }
}&#x27;

``` 
 Video file - Pass by file path (DashScope only)

 This example uses a locally saved [test.mp4](https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/en-US/20250415/nvwkcj/test.mp4) file.- Python 
- Java 

 Copy ```\nimport os
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

local_path = "xxx/test.mp4"
video_path = f"file://{local_path}"
messages = [
 {&#x27;role&#x27;:&#x27;user&#x27;,
 &#x27;content&#x27;: [{&#x27;video&#x27;: video_path,"fps":2},
 {&#x27;text&#x27;: &#x27;What scene does this video depict?&#x27;}]}]
response = dashscope.MultiModalConversation.call(
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 model=&#x27;qwen3.7-plus&#x27;, 
 messages=messages)
print(response.output.choices[0].message.content[0]["text"])

``` Copy ```\nimport java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversation;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationParam;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationResult;
import com.alibaba.dashscope.common.MultiModalMessage;
import com.alibaba.dashscope.common.Role;
import com.alibaba.dashscope.exception.ApiException;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.alibaba.dashscope.exception.UploadFileException;
import com.alibaba.dashscope.utils.Constants;

public class Main {

 static {
 Constants.baseHttpApiUrl="https://dashscope-intl.aliyuncs.com/api/v1";
 }
 
 public static void callWithLocalFile(String localPath)
 throws ApiException, NoApiKeyException, UploadFileException {
 String filePath = "file://"+localPath;
 MultiModalConversation conv = new MultiModalConversation();
 MultiModalMessage userMessage = MultiModalMessage.builder().role(Role.USER.getValue())
 .content(Arrays.asList(new HashMap<String, Object>()
 {{
 put("video", filePath);
 put("fps", 2);
 }}, 
 new HashMap<String, Object>(){{put("text", "What scene does this video depict?");}})).build();
 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model("qwen3.7-plus") 
 .messages(Arrays.asList(userMessage))
 .build();
 MultiModalConversationResult result = conv.call(param);
 System.out.println(result.getOutput().getChoices().get(0).getMessage().getContent().get(0).get("text"));}

 public static void main(String[] args) {
 try {
 callWithLocalFile("xxx/test.mp4");
 } catch (ApiException | NoApiKeyException | UploadFileException e) {
 System.out.println(e.getMessage());
 }
 System.exit(0);
 }
}

``` 
 Video file - Base64-encoded input

 - OpenAI compatible 
- DashScope 

 - Python 
- Node.js 
- curl 

 Copy ```\nfrom openai import OpenAI
import os
import base64

def encode_video(video_path):
 with open(video_path, "rb") as video_file:
 return base64.b64encode(video_file.read()).decode("utf-8")

base64_video = encode_video("xxx/test.mp4")
client = OpenAI(
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)
completion = client.chat.completions.create(
 model="qwen3.7-plus", 
 messages=[
 {
 "role": "user",
 "content": [
 {
 "type": "video_url",
 "video_url": {"url": f"data:video/mp4;base64,{base64_video}"},
 "fps":2
 },
 {"type": "text", "text": "What scene does this video depict?"},
 ],
 }
 ],
)
print(completion.choices[0].message.content)

``` Copy ```\nimport OpenAI from "openai";
import { readFileSync } from &#x27;fs&#x27;;

const openai = new OpenAI(
 {
 apiKey: process.env.DASHSCOPE_API_KEY,
 baseURL: "https://dashscope-intl.aliyuncs.com/compatible-mode/v1"
 }
);

const encodeVideo = (videoPath) => {
 const videoFile = readFileSync(videoPath);
 return videoFile.toString(&#x27;base64&#x27;);
 };
const base64Video = encodeVideo("xxx/test.mp4")
async function main() {
 const completion = await openai.chat.completions.create({
 model: "qwen3.7-plus", 
 messages: [
 {"role": "user",
 "content": [{
 "type": "video_url", 
 "video_url": {"url": `data:video/mp4;base64,${base64Video}`},
 "fps":2},
 {"type": "text", "text": "What scene does this video depict?"}]}]
 });
 console.log(completion.choices[0].message.content);
}

main();

``` 
- For display purposes, the Base64 encoded string in the code is truncated. You must pass the complete encoded string.


Copy ```\ncurl --location &#x27;https://dashscope-intl.aliyuncs.com/compatible-mode/v1/chat/completions&#x27; \
--header "Authorization: Bearer $DASHSCOPE_API_KEY" \
--header &#x27;Content-Type: application/json&#x27; \
--data &#x27;{
 "model": "qwen3.7-plus",
 "messages": [
 {
 "role": "user",
 "content": [
 {"type": "video_url", "video_url": {"url": "data:video/mp4;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAA..."},"fps":2},
 {"type": "text", "text": "What scene is depicted in the image?"}
 ]
 }]
}&#x27;

``` - Python 
- Java 
- curl 

 Copy ```\nimport base64
import os
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

def encode_video(video_path):
 with open(video_path, "rb") as video_file:
 return base64.b64encode(video_file.read()).decode("utf-8")

base64_video = encode_video("xxxx/test.mp4")

messages = [{&#x27;role&#x27;:&#x27;user&#x27;,
 &#x27;content&#x27;: [{&#x27;video&#x27;: f"data:video/mp4;base64,{base64_video}","fps":2},
 {&#x27;text&#x27;: &#x27;What scene does this video depict?&#x27;}]}]
response = dashscope.MultiModalConversation.call(
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 model=&#x27;qwen3.7-plus&#x27;,
 messages=messages)

print(response.output.choices[0].message.content[0]["text"])

``` Copy ```\nimport java.io.IOException;
import java.util.*;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

import com.alibaba.dashscope.aigc.multimodalconversation.*;
import com.alibaba.dashscope.common.MultiModalMessage;
import com.alibaba.dashscope.common.Role;
import com.alibaba.dashscope.exception.ApiException;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.alibaba.dashscope.exception.UploadFileException;
import com.alibaba.dashscope.utils.Constants;

public class Main {

 static {
 Constants.baseHttpApiUrl="https://dashscope-intl.aliyuncs.com/api/v1";
 }
 
 private static String encodeVideoToBase64(String videoPath) throws IOException {
 Path path = Paths.get(videoPath);
 byte[] videoBytes = Files.readAllBytes(path);
 return Base64.getEncoder().encodeToString(videoBytes);
 }

 public static void callWithLocalFile(String localPath)
 throws ApiException, NoApiKeyException, UploadFileException, IOException {

 String base64Video = encodeVideoToBase64(localPath);

 MultiModalConversation conv = new MultiModalConversation();
 MultiModalMessage userMessage = MultiModalMessage.builder().role(Role.USER.getValue())
 .content(Arrays.asList(new HashMap<String, Object>()
 {{
 put("video", "data:video/mp4;base64," + base64Video);
 put("fps", 2);
 }},
 new HashMap<String, Object>(){{put("text", "What scene does this video depict?");}})).build();

 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model("qwen3.7-plus")
 .messages(Arrays.asList(userMessage))
 .build();

 MultiModalConversationResult result = conv.call(param);
 System.out.println(result.getOutput().getChoices().get(0).getMessage().getContent().get(0).get("text"));
 }

 public static void main(String[] args) {
 try {
 callWithLocalFile("xxx/test.mp4");
 } catch (ApiException | NoApiKeyException | UploadFileException | IOException e) {
 System.out.println(e.getMessage());
 }
 System.exit(0);
 }
}

``` 
- For display purposes, the Base64 encoded string in the code is truncated. You must pass the complete encoded string.


Copy ```\ncurl -X POST https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation \
-H "Authorization: Bearer $DASHSCOPE_API_KEY" \
-H &#x27;Content-Type: application/json&#x27; \
-d &#x27;{
 "model": "qwen3.7-plus",
 "input":{
 "messages":[
 {
 "role": "user",
 "content": [
 {"video": "data:video/mp4;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAA..."},
 {"text": "What scene does this video depict? "}
 ]
 }
 ]
 }
}&#x27;

``` 
 Image list - Pass by file path (DashScope only)

 This example uses locally saved files: [football1.jpg](https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/en-US/20250415/spqrrx/football1.jpg), [football2.jpg](https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/en-US/20250415/vtnhyr/football2.jpg), [football3.jpg](https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/en-US/20250415/ykaoih/football3.jpg), and [football4.jpg](https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/en-US/20250415/vkuupi/football4.jpg).- Python 
- Java 

 Copy ```\nimport os
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

local_path1 = "football1.jpg"
local_path2 = "football2.jpg"
local_path3 = "football3.jpg"
local_path4 = "football4.jpg"

image_path1 = f"file://{local_path1}"
image_path2 = f"file://{local_path2}"
image_path3 = f"file://{local_path3}"
image_path4 = f"file://{local_path4}"

messages = [{&#x27;role&#x27;:&#x27;user&#x27;,
 &#x27;content&#x27;: [{&#x27;video&#x27;: [image_path1,image_path2,image_path3,image_path4],"fps":2},
 {&#x27;text&#x27;: &#x27;What scene does this video depict?&#x27;}]}]
response = dashscope.MultiModalConversation.call(
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 model=&#x27;qwen3.7-plus&#x27;,
 messages=messages)

print(response.output.choices[0].message.content[0]["text"])

``` Copy ```\n// DashScope SDK version must be 2.21.10 or later
import java.util.Arrays;
import java.util.Map;
import java.util.Collections;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversation;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationParam;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationResult;
import com.alibaba.dashscope.common.MultiModalMessage;
import com.alibaba.dashscope.common.Role;
import com.alibaba.dashscope.exception.ApiException;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.alibaba.dashscope.exception.UploadFileException;
import com.alibaba.dashscope.utils.Constants;

public class Main {

 static {
 Constants.baseHttpApiUrl="https://dashscope-intl.aliyuncs.com/api/v1";
 }
 
 private static final String MODEL_NAME = "qwen3.7-plus";
 public static void videoImageListSample(String localPath1, String localPath2, String localPath3, String localPath4)
 throws ApiException, NoApiKeyException, UploadFileException {
 MultiModalConversation conv = new MultiModalConversation();
 String filePath1 = "file://" + localPath1;
 String filePath2 = "file://" + localPath2;
 String filePath3 = "file://" + localPath3;
 String filePath4 = "file://" + localPath4;
 Map<String, Object> params = new HashMap<>();
 params.put("video", Arrays.asList(filePath1,filePath2,filePath3,filePath4));
 params.put("fps", 2);
 MultiModalMessage userMessage = MultiModalMessage.builder()
 .role(Role.USER.getValue())
 .content(Arrays.asList(params,
 Collections.singletonMap("text", "Describe the specific process of this video")))
 .build();
 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model(MODEL_NAME)
 .messages(Arrays.asList(userMessage)).build();
 MultiModalConversationResult result = conv.call(param);
 System.out.print(result.getOutput().getChoices().get(0).getMessage().getContent().get(0).get("text"));
 }
 public static void main(String[] args) {
 try {
 videoImageListSample(
 "xxx/football1.jpg",
 "xxx/football2.jpg",
 "xxx/football3.jpg",
 "xxx/football4.jpg");
 } catch (ApiException | NoApiKeyException | UploadFileException e) {
 System.out.println(e.getMessage());
 }
 System.exit(0);
 }
}

``` 
 Image list - Base64-encoded input

 - OpenAI compatible 
- DashScope 

 - Python 
- Node.js 
- curl 

 Copy ```\nimport os
from openai import OpenAI
import base64

def encode_image(image_path):
 with open(image_path, "rb") as image_file:
 return base64.b64encode(image_file.read()).decode("utf-8")

base64_image1 = encode_image("football1.jpg")
base64_image2 = encode_image("football2.jpg")
base64_image3 = encode_image("football3.jpg")
base64_image4 = encode_image("football4.jpg")
client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)
completion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=[ 
 {"role": "user","content": [
 {"type": "video","video": [
 f"data:image/jpeg;base64,{base64_image1}",
 f"data:image/jpeg;base64,{base64_image2}",
 f"data:image/jpeg;base64,{base64_image3}",
 f"data:image/jpeg;base64,{base64_image4}",]},
 {"type": "text","text": "Describe the specific process of this video"},
 ]}]
)
print(completion.choices[0].message.content)

``` Copy ```\nimport OpenAI from "openai";
import { readFileSync } from &#x27;fs&#x27;;

const openai = new OpenAI(
 {
 apiKey: process.env.DASHSCOPE_API_KEY,
 baseURL: "https://dashscope-intl.aliyuncs.com/compatible-mode/v1"
 }
);

const encodeImage = (imagePath) => {
 const imageFile = readFileSync(imagePath);
 return imageFile.toString(&#x27;base64&#x27;);
 };
 
const base64Image1 = encodeImage("football1.jpg")
const base64Image2 = encodeImage("football2.jpg")
const base64Image3 = encodeImage("football3.jpg")
const base64Image4 = encodeImage("football4.jpg")
async function main() {
 const completion = await openai.chat.completions.create({
 model: "qwen3.7-plus",
 messages: [
 {"role": "user",
 "content": [{"type": "video",
 "video": [
 `data:image/jpeg;base64,${base64Image1}`,
 `data:image/jpeg;base64,${base64Image2}`,
 `data:image/jpeg;base64,${base64Image3}`,
 `data:image/jpeg;base64,${base64Image4}`]},
 {"type": "text", "text": "What scene does this video depict?"}]}]
 });
 console.log(completion.choices[0].message.content);
}

main();

``` 
- For display purposes, the Base64-encoded strings in the code are truncated. You must pass the complete encoded strings.


Copy ```\ncurl -X POST https://dashscope-intl.aliyuncs.com/compatible-mode/v1/chat/completions \
-H "Authorization: Bearer $DASHSCOPE_API_KEY" \
-H &#x27;Content-Type: application/json&#x27; \
-d &#x27;{
 "model": "qwen3.7-plus",
 "messages": [{"role": "user",
 "content": [{"type": "video",
 "video": [
 "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAA...",
 "data:image/jpeg;base64,nEpp6jpnP57MoWSyOWwrkXMJhHRCWYeFYb...",
 "data:image/jpeg;base64,JHWQnJPc40GwQ7zERAtRMK6iIhnWw4080s...",
 "data:image/jpeg;base64,adB6QOU5HP7dAYBBOg/Fb7KIptlbyEOu58..."
 ]},
 {"type": "text",
 "text": "Describe the specific process of this video"}]}]
}&#x27;

``` - Python 
- Java 
- curl 

 Copy ```\nimport base64
import os
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

def encode_image(image_path):
 with open(image_path, "rb") as image_file:
 return base64.b64encode(image_file.read()).decode("utf-8")

base64_image1 = encode_image("football1.jpg")
base64_image2 = encode_image("football2.jpg")
base64_image3 = encode_image("football3.jpg")
base64_image4 = encode_image("football4.jpg")

messages = [{&#x27;role&#x27;:&#x27;user&#x27;,
 &#x27;content&#x27;: [
 {&#x27;video&#x27;:
 [f"data:image/jpeg;base64,{base64_image1}",
 f"data:image/jpeg;base64,{base64_image2}",
 f"data:image/jpeg;base64,{base64_image3}",
 f"data:image/jpeg;base64,{base64_image4}"],
 "fps":2},
 {&#x27;text&#x27;: &#x27;Describe the specific process of this video&#x27;}]}]

response = dashscope.MultiModalConversation.call(
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 model=&#x27;qwen3.7-plus&#x27;,
 messages=messages)

print(response.output.choices[0].message.content[0]["text"])

``` Copy ```\n// DashScope SDK version must be 2.21.10 or later
import java.util.*;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

import com.alibaba.dashscope.aigc.multimodalconversation.*;
import com.alibaba.dashscope.common.MultiModalMessage;
import com.alibaba.dashscope.common.Role;
import com.alibaba.dashscope.exception.ApiException;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.alibaba.dashscope.exception.UploadFileException;
import com.alibaba.dashscope.utils.Constants;

public class Main {
 static {
 Constants.baseHttpApiUrl="https://dashscope-intl.aliyuncs.com/api/v1";
 }

 private static String encodeImageToBase64(String imagePath) throws IOException {
 Path path = Paths.get(imagePath);
 byte[] imageBytes = Files.readAllBytes(path);
 return Base64.getEncoder().encodeToString(imageBytes);
 }
 
 public static void simpleMultiModalConversationCall()
 throws ApiException, NoApiKeyException, UploadFileException, IOException {
 String base64Image1 = encodeImageToBase64("football1.jpg");
 String base64Image2 = encodeImageToBase64("football2.jpg");
 String base64Image3 = encodeImageToBase64("football3.jpg");
 String base64Image4 = encodeImageToBase64("football4.jpg");

 MultiModalConversation conv = new MultiModalConversation();
 Map<String, Object> params = new HashMap<>();
 params.put("video", Arrays.asList(
 "data:image/jpeg;base64," + base64Image1,
 "data:image/jpeg;base64," + base64Image2,
 "data:image/jpeg;base64," + base64Image3,
 "data:image/jpeg;base64," + base64Image4));
 params.put("fps", 2);
 MultiModalMessage userMessage = MultiModalMessage.builder()
 .role(Role.USER.getValue())
 .content(Arrays.asList(params,
 Collections.singletonMap("text", "Describe the specific process of this video")))
 .build();
 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model("qwen3.7-plus")
 .messages(Arrays.asList(userMessage)).build();
 MultiModalConversationResult result = conv.call(param);
 System.out.print(result.getOutput().getChoices().get(0).getMessage().getContent().get(0).get("text"));
 }
 public static void main(String[] args) {
 try {
 simpleMultiModalConversationCall();
 } catch (ApiException | NoApiKeyException | UploadFileException | IOException e) {
 System.out.println(e.getMessage());
 }
 System.exit(0);
 }
}

``` 
- For display purposes, the Base64 encoded string in the code is truncated. You must pass the complete encoded string.


Copy ```\ncurl -X POST https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation \
-H "Authorization: Bearer $DASHSCOPE_API_KEY" \
-H &#x27;Content-Type: application/json&#x27; \
-d &#x27;{
 "model": "qwen3.7-plus",
 "input":{
 "messages":[
 {
 "role": "user",
 "content": [
 {"video": ["data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAA...",
 "data:image/jpeg;base64,nEpp6jpnP57MoWSyOWwrkXMJhHRCWYeFYb...",
 "data:image/jpeg;base64,JHWQnJPc40GwQ7zERAtRMK6iIhnWw4080s...",
 "data:image/jpeg;base64,adB6QOU5HP7dAYBBOg/Fb7KIptlbyEOu58..."],
 "fps":2},
 {"text": "What scene does this video depict?"}
 ]
 }
 ]
 }
}&#x27;

``` 
### [​ ](#handle-high-resolution-images) Handle high-resolution images

The visual understanding model API has a limit on the number of visual tokens for a single image after encoding. With default configurations, high-resolution images are compressed, which may result in a loss of detail and affect understanding accuracy. Enable `vl_high_resolution_images` or adjust `max_pixels` to increase the number of visual tokens, which preserves more image details and improves understanding.
 View the pixels per visual token, token limit, and pixel limit for each model

 If an input image has more pixels than the model&#x27;s pixel limit, the image is scaled down to fit within the limit. **Model****Pixels per token****vl_high_resolution_images****max_pixels****Token limit****Pixel limit**`Qwen3.5` and `Qwen3-VL` series models`32*32``true``max_pixels` is invalid`16384 tokens``16777216` (which is `16384*32*32`)`Qwen3.5` and `Qwen3-VL` series models`32*32``false` (default)Customizable. The default is `2621440`, and the maximum is `16777216`.Determined by `max_pixels`, which is `max_pixels/32/32`.`max_pixels``qwen-vl-max`, `qwen-vl-max-latest`, `qwen-vl-max-2025-08-13`, `qwen-vl-plus`, `qwen-vl-plus-latest`, `qwen-vl-plus-2025-08-15``32*32``true``max_pixels` is invalid`16384 tokens``16777216` (which is `16384*32*32`)Same Qwen2.5-VL models above`32*32``false` (default)Customizable. The default is `2621440`, and the maximum is `16777216`.Determined by `max_pixels`, which is `max_pixels/32/32`.`max_pixels``QVQ` and other `Qwen2.5-VL` models`28*28`Not supportedCustomizable. The default is `1003520`, and the maximum is `12845056`.Determined by `max_pixels`, which is `max_pixels/28/28`.`max_pixels` 
- OpenAI compatible 
- DashScope 

 - Python 
- Node.js 
- curl 

 Copy ```\nimport os
import time
from openai import OpenAI

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)

completion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=[
 {"role": "user","content": [
 {"type": "image_url","image_url": {"url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250212/earbrt/vcg_VCG211286867973_RF.jpg"},
 # max_pixels represents the maximum pixel threshold for the input image. It is invalid when vl_high_resolution_images=True, but customizable when vl_high_resolution_images=False. The maximum value varies by model.
 # "max_pixels": 16384 * 32 * 32
 },
 {"type": "text", "text": "What festival atmosphere does this picture show?"},
 ],
 }
 ],
 extra_body={"vl_high_resolution_images":True}

)
print(f"Model output: {completion.choices[0].message.content}")
print(f"Total input tokens: {completion.usage.prompt_tokens}")

``` Copy ```\nimport OpenAI from "openai";

const openai = new OpenAI(
 {
 apiKey: process.env.DASHSCOPE_API_KEY,
 baseURL: "https://dashscope-intl.aliyuncs.com/compatible-mode/v1"
 }
);

const response = await openai.chat.completions.create({
 model: "qwen3.7-plus",
 messages: [
 {role: "user",content: [
 {type: "image_url",
 image_url: {"url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250212/earbrt/vcg_VCG211286867973_RF.jpg"},
 // max_pixels represents the maximum pixel threshold for the input image. It is not effective when vl_high_resolution_images=True, but customizable when vl_high_resolution_images=False. The maximum value varies by model.
 // "max_pixels": 2560 * 32 * 32
 },
 {type: "text", text: "What festival atmosphere does this picture show?" },
 ]}],
 vl_high_resolution_images:true
 })


console.log("Model output:",response.choices[0].message.content);
console.log("Total input tokens",response.usage.prompt_tokens);

``` Copy ```\ncurl -X POST https://dashscope-intl.aliyuncs.com/compatible-mode/v1/chat/completions \
-H "Authorization: Bearer $DASHSCOPE_API_KEY" \
-H &#x27;Content-Type: application/json&#x27; \
-d &#x27;{
 "model": "qwen3.7-plus",
 "messages": [
 {
 "role": "user",
 "content": [
 {
 "type": "image_url",
 "image_url": {
 "url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250212/earbrt/vcg_VCG211286867973_RF.jpg"
 }
 },
 {
 "type": "text",
 "text": "What festival atmosphere does this picture show?"
 }
 ]
 }
 ],
 "vl_high_resolution_images":true
}&#x27;

``` - Python 
- Java 
- curl 

 Copy ```\nimport os
import time

import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

messages = [
 {
 "role": "user",
 "content": [
 {"image": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250212/earbrt/vcg_VCG211286867973_RF.jpg",
 # max_pixels represents the maximum pixel threshold for the input image. It is invalid when vl_high_resolution_images=True, but customizable when vl_high_resolution_images=False. The maximum value varies by model.
 # "max_pixels": 16384 * 32 * 32
 },
 {"text": "What festival atmosphere does this picture show?"}
 ]
 }
]

response = dashscope.MultiModalConversation.call(
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 model=&#x27;qwen3.7-plus&#x27;,
 messages=messages,
 vl_high_resolution_images=True
 )

print("Model output",response.output.choices[0].message.content[0]["text"])
print("Total input tokens:",response.usage.input_tokens)

``` Copy ```\nimport java.util.Arrays;
import java.util.Collections;
import java.util.Map;
import java.util.HashMap;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversation;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationParam;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationResult;
import com.alibaba.dashscope.common.MultiModalMessage;
import com.alibaba.dashscope.common.Role;
import com.alibaba.dashscope.exception.ApiException;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.alibaba.dashscope.exception.UploadFileException;
import com.alibaba.dashscope.utils.Constants;

public class Main {

 static {
 Constants.baseHttpApiUrl="https://dashscope-intl.aliyuncs.com/api/v1";
 }

 public static void simpleMultiModalConversationCall()
 throws ApiException, NoApiKeyException, UploadFileException {
 MultiModalConversation conv = new MultiModalConversation();
 Map<String, Object> map = new HashMap<>();
 map.put("image", "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250212/earbrt/vcg_VCG211286867973_RF.jpg");
 // max_pixels represents the maximum pixel threshold for the input image. It is invalid when vl_high_resolution_images=True, but customizable when vl_high_resolution_images=False. The maximum value varies by model.
 // map.put("max_pixels", 2621440);
 MultiModalMessage userMessage = MultiModalMessage.builder().role(Role.USER.getValue())
 .content(Arrays.asList(
 map,
 Collections.singletonMap("text", "What festival atmosphere does this picture show?"))).build();
 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model("qwen3.7-plus")
 .message(userMessage)
 .vlHighResolutionImages(true)
 .build();
 MultiModalConversationResult result = conv.call(param);
 System.out.println(result.getOutput().getChoices().get(0).getMessage().getContent().get(0).get("text"));
 System.out.println(result.getUsage().getInputTokens());
 }

 public static void main(String[] args) {
 try {
 simpleMultiModalConversationCall();
 } catch (ApiException | NoApiKeyException | UploadFileException e) {
 System.out.println(e.getMessage());
 }
 System.exit(0);
 }
}

``` Copy ```\ncurl -X POST https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation \
-H "Authorization: Bearer $DASHSCOPE_API_KEY" \
-H &#x27;Content-Type: application/json&#x27; \
-d &#x27;{
 "model": "qwen3.7-plus",
 "input":{
 "messages":[
 {
 "role": "user",
 "content": [
 {"image": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250212/earbrt/vcg_VCG211286867973_RF.jpg"},
 {"text": "What festival atmosphere does this picture show?"}
 ]
 }
 ]
 },
 "parameters": {
 "vl_high_resolution_images": true
 }
}&#x27;

``` 
### [​ ](#advanced-features) Advanced features


- [Multi-turn conversation](/developer-guides/text-generation/multi-turn)

- [Streaming output](/developer-guides/text-generation/streaming)


## [​ ](#limits) Limits

### [​ ](#input-file-limits) Input file limits

- Image limits 
- Video limits 

 
- 
**Image resolution:**

Minimum size: The width and height of the image must both be greater than `10` pixels.

- Aspect ratio: The ratio of the long side to the short side of the image cannot exceed `200:1`.

- Pixel limit:

We recommend keeping the image resolution within `8K (7680x4320)`. Images that exceed this resolution may cause API call timeouts because of large file sizes and long network transmission times.

- Automatic scaling: The model can adjust the image size using `max_pixels` and `min_pixels`. Therefore, providing ultra-high-resolution images does not improve recognition accuracy but increases the risk of call failures. We recommend scaling the image to a reasonable size on the client in advance.


- 
**Supported image formats**


For resolutions below `4K (3840x2160)`, the supported image formats are as follows:
**Image format****Common extensions****MIME type**BMP.bmpimage/bmpJPEG.jpe, .jpeg, .jpgimage/jpegPNG.pngimage/pngTIFF.tif, .tiffimage/tiffWEBP.webpimage/webpHEIC.heicimage/heic 


- 
For resolutions between `4K (3840x2160)` and `8K (7680x4320)`, only the JPEG, JPG, and PNG formats are supported.


- 
**Image size:**

When passed as a public URL: A single image cannot exceed `20 MB` for Qwen3.5, and `10 MB` for other models.

- When passed as a local path: A single image cannot exceed `10 MB`.

- When passed as a Base64-encoded string: The encoded string cannot exceed `10 MB`.


 For more information about how to compress the file size, see [How to compress an image or video to the required size](#faq). 

- 
**Number of supported images:**

`qwen3.7-plus` series: up to **2048** images per request when passed as a **public URL** or **local file path**.

- Other models: up to **256** images per request when passed as a **public URL** or **local file path**.

- When passed as **Base64-encoded** strings: up to **250** images per request.


These per-request caps are not the only constraint: the combined **image and text** token usage must stay within the model&#x27;s maximum input. The total tokens from **all** images and **all** text must be **less than** the model&#x27;s maximum input length.
 For example, if you use the `qwen3-vl-plus` model in thinking mode, the maximum input is `258048` `tokens`. If the input text consumes `100` tokens and each image consumes `2560` tokens, you can pass a maximum of `(258048 - 100) / 2560 = 100` images. 


 
- 
**When passed as an image list, the number of images in the list is limited as follows:**

`qwen3.5` series: A minimum of 4 images and a maximum of 8,000 images.

- `qwen3-vl-plus` series, `qwen3-vl-flash` series, `qwen3-vl-235b-a22b-thinking`, and `qwen3-vl-235b-a22b-instruct`: A minimum of 4 images and a maximum of 2,000 images.

- Other `Qwen3-VL` open source, `Qwen2.5-VL` (including commercial and open source versions), and `QVQ` series models: A minimum of 4 images and a maximum of 512 images.

- Other models: A minimum of 4 images and a maximum of 80 images.


- 
**When passed as a video file:**


**Video size:**

When passed as a public URL:

`qwen3.5` series, `Qwen3-VL` series, and `qwen-vl-max` (including `qwen-vl-max-latest`, `qwen-vl-max-2025-04-08`, and all subsequent versions): Cannot exceed 2 GB.

- `qwen-vl-plus` series, other `qwen-vl-max` models, `Qwen2.5-VL` open source series, and `QVQ` series models: Cannot exceed 1 GB.

- Other models cannot exceed 150 MB.


- When passed as a Base64-encoded string: The encoded string must be less than 10 MB.

- When passed as a local file path: The video file cannot exceed 100 MB.


 For more information about how to compress the file size, see [How to compress an image or video to the required size](#faq). 

- 
**Video duration:**

`qwen3.5` series: 2 seconds to 2 hours.

- `qwen3-vl-plus` series, `qwen3-vl-flash` series, `qwen3-vl-235b-a22b-thinking`, and `qwen3-vl-235b-a22b-instruct`: 2 seconds to 1 hour.

- Other `Qwen3-VL` open source series and `qwen-vl-max` (including `qwen-vl-max-latest`, `qwen-vl-max-2025-04-08`, and later updated versions): 2 seconds to 20 minutes.

- `qwen-vl-plus` series, other `qwen-vl-max` models, `Qwen2.5-VL` open source series, and `QVQ` series models: 2 seconds to 10 minutes.

- Other models: 2 seconds to 40 seconds.


- 
**Video format:** MP4, AVI, MKV, MOV, FLV, WMV, and more.


- 
**Video dimensions:** No specific limit. The model can automatically adjust video dimensions using `max_pixels` and `min_pixels`. Larger video files do not result in better understanding.


- 
**Number of videos:** Regardless of whether you use a **public URL**, **Base64-encoded** string, or **local file path**, you can pass at most **64** videos per request. In addition to this count limit, the total tokens from **all** videos and **all** text must be **less than** the model&#x27;s maximum input length.


- 
**Audio understanding:** The model does not support understanding the audio from video files.


 
### [​ ](#file-input-methods) File input methods


- **Public URL**: Provide a publicly accessible file address that supports the HTTP or HTTPS protocol. For optimal stability and performance, upload the file to OSS to get a public URL.


 To ensure that the model can successfully download the file, the request header of the public URL **must** include Content-Length (file size) and Content-Type (media type, such as image/jpeg). If either field is missing or incorrect, the file download fails. 

- **Pass as a Base64-encoded string:** Convert the file to a Base64-encoded string and then pass it.

- **Pass as a local file path (DashScope SDK only):** Pass the path of the local file.


 For recommendations on file input methods, see [How to choose a file upload method?](#faq) 
## [​ ](#deploy-to-production) Deploy to production


- 
**Image/video pre-processing:** Visual understanding models have size limits for input files. For more information about how to compress files, see [Image or video compression methods](#faq).


- 
**Process text files:** Visual understanding models support processing files only in image format and cannot directly process text files. Convert the text file to an image format. We recommend using an image processing library, such as `Python`&#x27;s `pdf2image`, to convert the file page by page into multiple high-quality images. Then pass them to the model using the [multiple image input](#work-with-multiple-images) method.


- 
**Fault tolerance and stability**

Timeout handling: In non-streaming calls, if the model does not finish outputting within 180 seconds, a timeout error is usually triggered. To improve the user experience, the response body returns any content already generated after a timeout. If the response header contains `x-dashscope-partialresponse:true`, the response triggered a timeout. You can use the [partial mode](/developer-guides/text-generation/partial-mode) feature (supported by some models) to add the generated content to the messages array and send the request again. This lets the large model continue generating content. For more information, see [Continue writing based on incomplete output](/developer-guides/text-generation/partial-mode).

- Retry mechanism: Design a reasonable API call retry logic, such as exponential backoff, to handle network fluctuations or temporary service unavailability.


## [​ ](#billing-rules) Billing rules


- **Billing:** The total cost is based on the total number of input and output tokens. For input and output prices, see [Models](https://www.qwencloud.com/models).

**Token composition:** Input tokens consist of text tokens and tokens converted from images or videos. Output tokens are the text that the model generates. In thinking mode, the model&#x27;s thought process also counts toward the output tokens. If the thought process is not an output in thinking mode, the price for non-thinking mode applies.

- **Calculate image and video tokens:** Use the following code to estimate the token consumption for an image or video. The estimate is for reference only. The actual usage is based on the API response.


 Calculate image and video tokens

 - Image 
- Video 

 Formula: `Image Token = h_bar * w_bar / token_pixels + 2`
- `h_bar, w_bar`: The height and width of the scaled image. Before processing an image, the model pre-processes it by scaling it down to a specific pixel limit. The pixel limit depends on the values of the `max_pixels` and `vl_high_resolution_images` parameters. For more information, see [Process high-resolution images](#handle-high-resolution-images).

- `token_pixels`: The pixel value that corresponds to each visual `token`. This value varies by model:

`Qwen3.5`, `Qwen3-VL`, `qwen-vl-max`, `qwen-vl-max-latest`, `qwen-vl-max-2025-08-13`, `qwen-vl-plus`, `qwen-vl-plus-latest`, `qwen-vl-plus-2025-08-15`: Each `token` corresponds to `32x32` pixels.

- `QVQ` and other `Qwen2.5-VL` models: Each token corresponds to `28x28` pixels.


The following code shows the approximate image scaling logic within the model. Use it to estimate the token count for an image. The actual billing is based on the API response.Copy ```\nimport math
# Use the following command to install the Pillow library: pip install Pillow
from PIL import Image

def token_calculate(image_path, max_pixels, vl_high_resolution_images):
 # Open the specified image file.
 image = Image.open(image_path)

 # Get the original dimensions of the image.
 height = image.height
 width = image.width

 # Adjust the width and height to be multiples of 32 or 28, depending on the model.
 h_bar = round(height / 32) * 32
 w_bar = round(width / 32) * 32

 # Lower limit for image tokens: 4 tokens.
 min_pixels = 4 * 32 * 32
 # If vl_high_resolution_images is set to True, the upper limit for input image tokens is 16,386, and the corresponding maximum pixel value is 16384 * 32 * 32 or 16384 * 28 * 28. Otherwise, it is the value set for max_pixels.
 if vl_high_resolution_images:
 max_pixels = 16384 * 32 * 32
 else:
 max_pixels = max_pixels

 # Scale the image so that the total number of pixels is within the range of [min_pixels, max_pixels].
 if h_bar * w_bar > max_pixels:
 beta = math.sqrt((height * width) / max_pixels)
 h_bar = math.floor(height / beta / 32) * 32
 w_bar = math.floor(width / beta / 32) * 32
 elif h_bar * w_bar < min_pixels:
 beta = math.sqrt(min_pixels / (height * width))
 h_bar = math.ceil(height * beta / 32) * 32
 w_bar = math.ceil(width * beta / 32) * 32
 return h_bar, w_bar

if __name__ == "__main__":
 # Replace xxx/test.jpg with the path to your local image.
 h_bar, w_bar = token_calculate("xxx/test.jpg", max_pixels=16384*32*32, vl_high_resolution_images=False)
 print(f"Scaled image dimensions: height {h_bar}, width {w_bar}")
 # The system automatically adds the <vision_bos> and <vision_eos> visual markers (1 token each).
 token = int((h_bar * w_bar) / (32 * 32))+2
 print(f"Number of tokens for the image: {token}")

``` 
- **Video file:**


When the model processes a video file, it first extracts the video frames and then calculates the total number of tokens for all the frames. You can use the following code to estimate the total number of tokens that a video consumes by providing the video path:Copy ```\n# Before use, install: pip install opencv-python
import math
import os
import logging
import cv2

logger = logging.getLogger(__name__)

FRAME_FACTOR = 2
IMAGE_FACTOR = 32
MAX_RATIO = 200
VIDEO_MIN_PIXELS = 4 * 32 * 32
VIDEO_MAX_PIXELS = 640 * 32 * 32
FPS = 2.0
FPS_MIN_FRAMES = 4
FPS_MAX_FRAMES = 2000
VIDEO_TOTAL_PIXELS = int(float(os.environ.get(&#x27;VIDEO_MAX_PIXELS&#x27;, 131072 * 32 * 32)))

def round_by_factor(number: int, factor: int) -> int:
 return round(number / factor) * factor

def ceil_by_factor(number: int, factor: int) -> int:
 return math.ceil(number / factor) * factor

def floor_by_factor(number: int, factor: int) -> int:
 return math.floor(number / factor) * factor

def extract_vision_info(conversations):
 vision_infos = []
 if isinstance(conversations[0], dict):
 conversations = [conversations]
 for conversation in conversations:
 for message in conversation:
 if isinstance(message["content"], list):
 for ele in message["content"]:
 if (
 "image" in ele
 or "image_url" in ele
 or "video" in ele
 or ele.get("type","") in ("image", "image_url", "video")
 ):
 vision_infos.append(ele)
 return vision_infos

def smart_nframes(ele,total_frames,video_fps):
 assert not ("fps" in ele and "nframes" in ele), "Only accept either `fps` or `nframes`"
 fps = ele.get("fps", FPS)
 min_frames = ceil_by_factor(ele.get("min_frames", FPS_MIN_FRAMES), FRAME_FACTOR)
 max_frames = floor_by_factor(ele.get("max_frames", min(FPS_MAX_FRAMES, total_frames)), FRAME_FACTOR)
 duration = total_frames / video_fps if video_fps != 0 else 0
 if duration-int(duration)>(1/fps):
 total_frames = math.ceil(duration * video_fps)
 else:
 total_frames = math.ceil(int(duration)*video_fps)
 nframes = total_frames / video_fps * fps
 if nframes > total_frames:
 logger.warning(f"smart_nframes: nframes[{nframes}] > total_frames[{total_frames}]")
 nframes = int(min(min(max(nframes, min_frames), max_frames), total_frames))
 if not (FRAME_FACTOR <= nframes and nframes <= total_frames):
 raise ValueError(f"nframes should in interval [{FRAME_FACTOR}, {total_frames}], but got {nframes}.")
 return nframes

def get_video(video_path):
 cap = cv2.VideoCapture(video_path)
 frame_width = int(cap.get(cv2.CAP_PROP_FRAME_WIDTH))
 frame_height = int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT))
 total_frames = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
 video_fps = cap.get(cv2.CAP_PROP_FPS)
 return frame_height, frame_width, total_frames, video_fps

def smart_resize(ele, path, factor=IMAGE_FACTOR):
 height, width, total_frames, video_fps = get_video(path)
 min_pixels = VIDEO_MIN_PIXELS
 total_pixels = VIDEO_TOTAL_PIXELS
 nframes = smart_nframes(ele, total_frames, video_fps)
 max_pixels = max(min(VIDEO_MAX_PIXELS, total_pixels / nframes * FRAME_FACTOR),int(min_pixels * 1.05))

 if max(height, width) / min(height, width) > MAX_RATIO:
 raise ValueError(
 f"absolute aspect ratio must be smaller than {MAX_RATIO}, got {max(height, width) / min(height, width)}"
 )

 h_bar = max(factor, round_by_factor(height, factor))
 w_bar = max(factor, round_by_factor(width, factor))
 if h_bar * w_bar > max_pixels:
 beta = math.sqrt((height * width) / max_pixels)
 h_bar = floor_by_factor(height / beta, factor)
 w_bar = floor_by_factor(width / beta, factor)
 elif h_bar * w_bar < min_pixels:
 beta = math.sqrt(min_pixels / (height * width))
 h_bar = ceil_by_factor(height * beta, factor)
 w_bar = ceil_by_factor(width * beta, factor)
 return h_bar, w_bar

def token_calculate(video_path, fps):
 messages = [{"content": [{"video": video_path, "fps": fps}]}]
 vision_infos = extract_vision_info(messages)[0]
 resized_height, resized_width = smart_resize(vision_infos, video_path)
 height, width, total_frames, video_fps = get_video(video_path)
 num_frames = smart_nframes(vision_infos, total_frames, video_fps)
 print(f"Original video dimensions: {height}*{width}, Dimensions for model input: {resized_height}*{resized_width}, Total video frames: {total_frames}, Total frames extracted when fps is {fps}: {num_frames}", end=", ")
 video_token = int(math.ceil(num_frames / 2) * resized_height / 32 * resized_width / 32)
 video_token += 2
 return video_token

video_token = token_calculate("xxx/test.mp4", 1)
print("Video tokens:", video_token)

``` 
- **Image list:**


If you provide the video as a list of images, use the following code to calculate the number of tokens consumed:Copy ```\n# Install before use: pip install Pillow
import math
import os
import logging
from typing import Tuple
from PIL import Image

logger = logging.getLogger(__name__)

FRAME_FACTOR = 2
IMAGE_FACTOR = 32
TOKEN_DIVISOR = 32
VISION_SPECIAL_TOKENS = 2
MAX_RATIO = 200
VIDEO_MIN_PIXELS = 4 * 32 * 32
VIDEO_MAX_PIXELS = 640 * 32 * 32
VIDEO_TOTAL_PIXELS = int(float(os.environ.get(&#x27;VIDEO_MAX_PIXELS&#x27;, 131072 * 32 * 32)))

def round_by_factor(number: int, factor: int) -> int:
 return round(number / factor) * factor

def ceil_by_factor(number: int, factor: int) -> int:
 return math.ceil(number / factor) * factor

def floor_by_factor(number: int, factor: int) -> int:
 return math.floor(number / factor) * factor

def get_image_size(image_path: str) -> Tuple[int, int]:
 if not os.path.exists(image_path):
 raise FileNotFoundError(f"Image file not found: {image_path}")
 try:
 image = Image.open(image_path)
 height = image.height
 width = image.width
 image.close()
 return height, width
 except Exception as e:
 raise ValueError(f"Cannot read the image file {image_path}: {str(e)}")

def smart_resize(height: int, width: int, nframes: int, factor: int = IMAGE_FACTOR) -> Tuple[int, int]:
 min_pixels = VIDEO_MIN_PIXELS
 total_pixels = VIDEO_TOTAL_PIXELS
 max_pixels = max(min(VIDEO_MAX_PIXELS, total_pixels / nframes * FRAME_FACTOR), int(min_pixels * 1.05))

 aspect_ratio = max(height, width) / min(height, width)
 if aspect_ratio > MAX_RATIO:
 raise ValueError(
 f"The aspect ratio of the image must be less than {MAX_RATIO}:1. The current ratio is {aspect_ratio:.2f}:1."
 )

 h_bar = max(factor, round_by_factor(height, factor))
 w_bar = max(factor, round_by_factor(width, factor))
 if h_bar * w_bar > max_pixels:
 beta = math.sqrt((height * width) / max_pixels)
 h_bar = floor_by_factor(height / beta, factor)
 w_bar = floor_by_factor(width / beta, factor)
 elif h_bar * w_bar < min_pixels:
 beta = math.sqrt(min_pixels / (height * width))
 h_bar = ceil_by_factor(height * beta, factor)
 w_bar = ceil_by_factor(width * beta, factor)
 return h_bar, w_bar

def calculate_video_tokens(image_path: str, nframes: int = 1, factor: int = IMAGE_FACTOR, verbose: bool = True) -> int:
 height, width = get_image_size(image_path)
 resized_height, resized_width = smart_resize(height, width, nframes, factor)
 video_token = int(
 math.ceil(nframes / 2) *
 (resized_height / TOKEN_DIVISOR) *
 (resized_width / TOKEN_DIVISOR)
 )
 video_token += VISION_SPECIAL_TOKENS
 if verbose:
 print(f"Original video frame dimensions: {height}x{width}, dimensions input to the model: {resized_height}x{resized_width}, ", end="")
 return video_token

if __name__ == "__main__":
 try:
 video_token = calculate_video_tokens("xxx/test.jpg", nframes=30)
 print(f"Video tokens: {video_token}\n")
 except Exception as e:
 print(f"Error: {str(e)}\n")

``` 

- **View bills:** View your bills in the [Billing section](https://home.qwencloud.com/billing/overview).


## [​ ](#reference) Reference

For the input and output parameters of visual understanding models, see [Chat API](/api-reference/chat/openai-chat).
## [​ ](#faq) FAQ

 How do I choose a file upload method?

 Choose the most suitable upload method based on the SDK type, file size, and network stability.**Type****Specifications****DashScope SDK (Python, Java)****OpenAI compatible / DashScope HTTP**ImageGreater than 7 MB and less than 10 MBPass the local pathOnly public URLs are supported. We recommend using Object Storage Service.ImageLess than 7 MBPass the local pathBase64 encodingVideoGreater than 100 MBOnly public URLs are supported. We recommend using Object Storage Service.Only public URLs are supported. We recommend using Object Storage Service.VideoGreater than 7 MB and less than 100 MBPass the local pathOnly public URLs are supported. We recommend using Object Storage Service.VideoLess than 7 MBPass the local pathBase64 encoding Base64 encoding increases data size, so the original file must be under 7 MB. Using Base64 or a local path avoids server-side download timeouts and improves stability. 
 How do I compress an image or video to the required size?

 Visual understanding models have size limits for input files. Compress files using the following methods.**Image compression methods**
- Online tools: Use online tools such as [CompressJPEG](https://compressjpeg.com/) or [TinyPng](https://free.tinypng.site/).

- Local software: Use software such as Photoshop to adjust the quality during export.

- Code implementation:


Copy ```\n# pip install pillow

from PIL import Image
def compress_image(input_path, output_path, quality=85):
 with Image.open(input_path) as img:
 img.save(output_path, "JPEG", optimize=True, quality=quality)

# Pass a local image
compress_image("/xxx/before-large.jpeg","/xxx/after-min.jpeg")

``` **Video compression methods**
- Online tools: Use online tools such as [FreeConvert](https://www.freeconvert.com).

- Local software: Use software such as [HandBrake](https://handbrake.fr/).

- Code implementation: Use the FFmpeg tool. For more information, see the [FFmpeg official website](https://ffmpeg.org/download.html).


Copy ```\n# Basic conversion command
# -i: input file path
# -vcodec: video encoder (libx264 recommended)
# -crf: controls video quality. Value range: [18-28]. Smaller values = higher quality.
# --preset: controls encoding speed vs compression. Common values: slow, fast, faster.
# -y: overwrite existing file.

ffmpeg -i input.mp4 -vcodec libx264 -crf 28 -preset slow output.mp4

``` 
 After the model outputs object localization results, how do I draw the bounding boxes on the original image?

 After the visual understanding model outputs object localization results, use the following code to draw the bounding boxes and their labels on the original image.
- Qwen2.5-VL: Returns coordinates as absolute values in pixels. These coordinates are relative to the top-left corner of the scaled image. To draw the bounding boxes, see the code in [qwen2_5_vl_2d.py](https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/en-US/20251103/bgyust/qwen2_5-vl-2d.py).

- Qwen3-VL: Returns relative coordinates that are normalized to the range `[0, 999]`. To draw the bounding boxes, see the code in [qwen3_vl_2d.py](https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/en-US/20251103/wpucdo/qwen3-vl-2d.py) (for 2D localization) or [qwen3_vl_3d.zip](https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/en-US/20251103/mjjrka/qwen3_vl_3d.zip) (for 3D localization).


 
## [​ ](#error-codes) Error codes

If a call fails, see [Error messages](/api-reference/preparation/error-messages). [Previous ](/developer-guides/getting-started/vision-models)[Text extraction OCR for docs and tables Next ](/developer-guides/multimodal/ocr)
