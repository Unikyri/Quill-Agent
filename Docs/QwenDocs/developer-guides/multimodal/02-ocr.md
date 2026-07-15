# Text extraction

> **Source:** https://docs.qwencloud.com/developer-guides/multimodal/ocr

OCR for docs and tables

 Copy page Qwen-OCR extracts text and parses structured data from images like scanned documents, tables, and receipts. It supports multiple languages, information extraction, table parsing, and formula recognition.
**Try it online:** [Qwen Cloud](https://home.qwencloud.com/try-ai)
## [​ ](#examples) Examples

**Input image****Recognition result****Recognize multiple languages**
 `INTERNATIONAL``MOTHER LANGUAGE``DAY``Привет!``你好!``Bonjour!``Merhaba!``Ciao!``Hello!``Ola!``בר מולד``Salam!`**Recognize skewed images**
 Product Introduction, Imported fiber filaments from South Korea. 6941990612023, Item No.: 2023**Locate text position** 
 
 The [high-precision recognition](#call-built-in-tasks) task supports text localization.**Visualization of localization** 
 
 See the [FAQ](#faq) on how to draw the bounding box of each text line onto the original image.
## [​ ](#availability) Availability

**Model****Snapshot****Context window (tokens)****Max input****Max output**qwen-vl-ocrNo38,19230,0008,192qwen-vl-ocr-2025-11-20Yes38,19230,0008,192 
 Example code for manually estimating image tokens (for budget reference only)

 Formula: Image tokens = `(h_bar * w_bar) / token_pixels + 2`.
- `h_bar * w_bar` represents the dimensions of the scaled image. The model pre-processes the image by scaling it to a specific pixel limit. This limit depends on the value of the `max_pixels` parameter.

- `token_pixels` represents the pixel value per `Token`.

For `qwen-vl-ocr` and `qwen-vl-ocr-2025-11-20`, this value is fixed at `32*32` (which is `1024`).

- For other models, this value is fixed at `28*28` (which is `784`).


The following code demonstrates the approximate image scaling logic that the model uses. You can use this code to estimate the token count for an image. The actual billing is based on the API response.Copy ```\nimport math
from PIL import Image

def smart_resize(image_path, min_pixels, max_pixels):
 """
 Pre-process an image.

 Parameters:
 image_path: The path to the image.
 """
 # Open the specified PNG image file.
 image = Image.open(image_path)

 # Get the original dimensions of the image.
 height = image.height
 width = image.width
 # Adjust the height to be a multiple of 28 or 32.
 h_bar = round(height / 32) * 32
 # Adjust the width to be a multiple of 28 or 32.
 w_bar = round(width / 32) * 32

 # Scale the image to adjust the total number of pixels to be within the range [min_pixels, max_pixels].
 if h_bar * w_bar > max_pixels:
 beta = math.sqrt((height * width) / max_pixels)
 h_bar = math.floor(height / beta / 32) * 32
 w_bar = math.floor(width / beta / 32) * 32
 elif h_bar * w_bar < min_pixels:
 beta = math.sqrt(min_pixels / (height * width))
 h_bar = math.ceil(height * beta / 32) * 32
 w_bar = math.ceil(width * beta / 32) * 32
 return h_bar, w_bar


# Replace xxx/test.png with the path to your local image.
h_bar, w_bar = smart_resize("xxx/test.png", min_pixels=32 * 32 * 3, max_pixels=8192 * 32 * 32)
print(f"The scaled image dimensions are: height {h_bar}, width {w_bar}")

# Calculate the number of image tokens: total pixels divided by 32 * 32.
token = int((h_bar * w_bar) / (32 * 32))

# <|vision_bos|> and <|vision_eos|> are visual markers. Each is counted as 1 token.
print(f"Total number of image tokens: {token + 2}")

``` 
## [​ ](#prerequisites) Prerequisites


- [Get an API key](/api-reference/preparation/api-key) and set it as an environment variable.

- To use the SDK, [install DashScope SDK](/api-reference/preparation/install-sdk). Minimum: Python 1.22.2, Java 2.18.4.

**DashScope SDK**

**Advantages**: Supports all advanced features, such as image rotation correction and built-in OCR tasks. It provides a complete feature set and a simple call method.

- **Scenarios**: Projects that require full functionality.


- **OpenAI SDK**

**Advantages**: Eases migration for users who already use the OpenAI SDK or its ecosystem tools.

- **Limitations**: Does not support calling advanced features, such as image rotation correction and built-in OCR tasks, directly with parameters. You must manually simulate these features by creating complex prompts and then parsing the output.

- **Scenarios**: Projects that already have an OpenAI integration and do not rely on advanced features exclusive to DashScope.


## [​ ](#getting-started) Getting started

The following example extracts key information from a [train ticket image (URL)](https://img.alicdn.com/imgextra/i2/O1CN01ktT8451iQutqReELT_!!6000000004408-0-tps-689-487.jpg) and returns it in JSON format. For local files and image limits, see [how to pass a local file](#pass-a-local-file-base64-encoding-or-file-path) and [image limits](#image-limits).
- OpenAI compatible 
- DashScope 

 Python Node.js curl Copy ```\nfrom openai import OpenAI
import os

PROMPT_TICKET_EXTRACTION = """
Please extract the invoice number, train number, departure station, destination station, departure date and time, seat number, seat type, ticket price, ID card number, and passenger name from the train ticket image.
Extract the key information accurately. Do not omit information or fabricate false information. Replace any single character that is blurry or obscured by glare with a question mark (?).
Return the data in JSON format: {&#x27;Invoice Number&#x27;: &#x27;xxx&#x27;, &#x27;Train Number&#x27;: &#x27;xxx&#x27;, &#x27;Departure Station&#x27;: &#x27;xxx&#x27;, &#x27;Destination Station&#x27;: &#x27;xxx&#x27;, &#x27;Departure Date and Time&#x27;: &#x27;xxx&#x27;, &#x27;Seat Number&#x27;: &#x27;xxx&#x27;, &#x27;Seat Type&#x27;: &#x27;xxx&#x27;, &#x27;Ticket Price&#x27;: &#x27;xxx&#x27;, &#x27;ID Card Number&#x27;: &#x27;xxx&#x27;, &#x27;Passenger Name&#x27;: &#x27;xxx&#x27;}
"""

try:
 client = OpenAI(
 # If you have not configured an environment variable, replace the following line with your API key: api_key="sk-xxx",
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
 )
 completion = client.chat.completions.create(
 model="qwen-vl-ocr-2025-11-20",
 messages=[
 {
 "role": "user",
 "content": [
 {
 "type": "image_url",
 "image_url": {"url":"https://img.alicdn.com/imgextra/i2/O1CN01ktT8451iQutqReELT_!!6000000004408-0-tps-689-487.jpg"},
 # The minimum pixel threshold for the input image.
 "min_pixels": 3072,
 # The maximum pixel threshold for the input image.
 "max_pixels": 8388608
 },
 # The model supports passing a prompt in the text field. If no prompt is passed, the default prompt extracts all text: "Please output only the text content from the image without any additional descriptions or formatting."
 {"type": "text", "text": PROMPT_TICKET_EXTRACTION}
 ]
 }
 ])
 print(completion.choices[0].message.content)
except Exception as e:
 print(f"Error message: {e}")

``` Example response

 Copy ```\n{
 "choices": [{
 "message": {
 "content": "```json\n{\n \"Invoice Number\": \"24329116804000\",\n \"Train Number\": \"G1948\",\n \"Departure Station\": \"Nanjing South Station\",\n \"Destination Station\": \"Zhengzhou East Station\",\n \"Departure Date and Time\": \"2024-11-14 11:46\",\n \"Seat Number\": \"Car 04, Seat 12A\",\n \"Seat Type\": \"Second Class\",\n \"Ticket Price\": \"￥337.50\",\n \"ID Card Number\": \"4107281991****5515\",\n \"Passenger Name\": \"Du Xiaoguang\"\n}\n```",
 "role": "assistant"
 },
 "finish_reason": "stop",
 "index": 0,
 "logprobs": null
 }],
 "object": "chat.completion",
 "usage": {
 "prompt_tokens": 606,
 "completion_tokens": 159,
 "total_tokens": 765
 },
 "created": 1742528311,
 "system_fingerprint": null,
 "model": "qwen-vl-ocr-2025-11-20",
 "id": "chatcmpl-20e5d9ed-e8a3-947d-bebb-c47ef1378598"
}

``` Python Java curl Copy ```\nimport os
import dashscope

PROMPT_TICKET_EXTRACTION = """
Please extract the invoice number, train number, departure station, destination station, departure date and time, seat number, seat type, ticket price, ID card number, and passenger name from the train ticket image.
Extract the key information accurately. Do not omit information or fabricate false information. Replace any single character that is blurry or obscured by glare with a question mark (?).
Return the data in JSON format: {&#x27;Invoice Number&#x27;: &#x27;xxx&#x27;, &#x27;Train Number&#x27;: &#x27;xxx&#x27;, &#x27;Departure Station&#x27;: &#x27;xxx&#x27;, &#x27;Destination Station&#x27;: &#x27;xxx&#x27;, &#x27;Departure Date and Time&#x27;: &#x27;xxx&#x27;, &#x27;Seat Number&#x27;: &#x27;xxx&#x27;, &#x27;Seat Type&#x27;: &#x27;xxx&#x27;, &#x27;Ticket Price&#x27;: &#x27;xxx&#x27;, &#x27;ID Card Number&#x27;: &#x27;xxx&#x27;, &#x27;Passenger Name&#x27;: &#x27;xxx&#x27;}
"""

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;
messages = [{
 "role": "user",
 "content": [{
 "image": "https://img.alicdn.com/imgextra/i2/O1CN01ktT8451iQutqReELT_!!6000000004408-0-tps-689-487.jpg",
 # The minimum pixel threshold for the input image.
 "min_pixels": 3072,
 # The maximum pixel threshold for the input image.
 "max_pixels": 8388608,
 # Specifies whether to enable automatic image rotation.
 "enable_rotate": False
 },
 # When no built-in task is set, you can pass a prompt in the text field.
 {"type": "text", "text": PROMPT_TICKET_EXTRACTION}]
}]
try:
 response = dashscope.MultiModalConversation.call(
 # If you have not configured an environment variable, replace the following line with your API key: api_key="sk-xxx",
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 model=&#x27;qwen-vl-ocr-2025-11-20&#x27;,
 messages=messages
 )
 print(response["output"]["choices"][0]["message"].content[0]["text"])
except Exception as e:
 print(f"An error occurred: {e}")

``` Example response

 Copy ```\n{
 "output": {
 "choices": [{
 "finish_reason": "stop",
 "message": {
 "role": "assistant",
 "content": [{
 "text": "```json\n{\n \"Invoice Number\": \"24329116804000\",\n \"Train Number\": \"G1948\",\n \"Departure Station\": \"Nanjing South Station\",\n \"Destination Station\": \"Zhengzhou East Station\",\n \"Departure Date and Time\": \"2024-11-14 11:46\",\n \"Seat Number\": \"Car 04, Seat 12A\",\n \"Seat Type\": \"Second Class\",\n \"Ticket Price\": \"￥337.50\",\n \"ID Card Number\": \"4107281991****5515\",\n \"Passenger Name\": \"Du Xiaoguang\"\n}\n```"
 }]
 }
 }]
 },
 "usage": {
 "total_tokens": 765,
 "output_tokens": 159,
 "input_tokens": 606,
 "image_tokens": 427
 },
 "request_id": "b3ca3bbb-2bdd-9367-90bd-f3f39e480db0"
}

``` 
## [​ ](#call-built-in-tasks) Call built-in tasks

To simplify calls in specific scenarios, the models (except for `qwen-vl-ocr-2024-10-28`) include several built-in tasks.
**How to use**:

- **Dashscope SDK**: You do not need to design and pass a `Prompt`. The model uses a fixed `Prompt` internally. Set the `ocr_options` parameter to call the built-in task.

- **OpenAI SDK:** You must manually enter the `Prompt` specified for the task.


The following table lists the value of `task`, the specified `Prompt`, the output format, and an example for each built-in task.
### [​ ](#high-precision-recognition) High-precision recognition

We recommend `qwen-vl-ocr-2025-08-28` or later versions for this task. Features:

- Recognizes and extracts text content.

- Detects the position of text by locating text lines and outputting their coordinates.


 For more information about how to draw the bounding box on the original image after you obtain the coordinates of the text bounding box, see the [FAQ](#faq). 
**Value of task****Specified prompt****Output format and example**`advanced_recognition`Locate all text lines and return the coordinates of the rotated rectangle `([cx, cy, width, height, angle])`.Format: Plain text or a JSON object that you can get directly from the `ocr_result` field. 
 Example: 
 
 `text`: The text content of each line. 
 `location`: Example value: `[x1, y1, x2, y2, x3, y3, x4, y4]`. Meaning: The absolute coordinates of the four vertices of the text box. The top-left corner of the original image is the origin `(0,0)`. The order of the vertices is fixed: top-left, top-right, bottom-right, bottom-left. 
 `rotate_rect`: Example value: `[center_x, center_y, width, height, angle]`. Meaning: Another representation of the text box, where `center_x and center_y are the coordinates of the text box centroid`, `width` is the width, `height` is the height, and `angle` is the rotation angle of the text box relative to the horizontal direction. The value is in the range of `[-90, 90]`.
- Python 
- Java 
- curl 

 Copy ```\nimport os
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

messages = [{
 "role": "user",
 "content": [{
 "image": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/ctdzex/biaozhun.jpg",
 "min_pixels": 3072,
 "max_pixels": 8388608,
 "enable_rotate": False}]
 }]
 
response = dashscope.MultiModalConversation.call(
 # If you have not configured an environment variable, replace the following line with your API key: api_key="sk-xxx",
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 model=&#x27;qwen-vl-ocr-2025-11-20&#x27;,
 messages=messages,
 # Set the built-in task to high-precision recognition.
 ocr_options={"task": "advanced_recognition"}
)
# The high-precision recognition task returns the result as plain text.
print(response["output"]["choices"][0]["message"].content[0]["text"])

``` Copy ```\n// dashscope SDK version >= 2.18.4
import java.util.Arrays;
import java.util.Collections;
import java.util.Map;
import java.util.HashMap;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversation;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationParam;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationResult;
import com.alibaba.dashscope.aigc.multimodalconversation.OcrOptions;
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
 map.put("image", "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/ctdzex/biaozhun.jpg");
 map.put("max_pixels", 8388608);
 map.put("min_pixels", 3072);
 map.put("enable_rotate", false);
 
 // Configure the built-in OCR task.
 OcrOptions ocrOptions = OcrOptions.builder()
 .task(OcrOptions.Task.ADVANCED_RECOGNITION)
 .build();
 MultiModalMessage userMessage = MultiModalMessage.builder().role(Role.USER.getValue())
 .content(Arrays.asList(
 map
 )).build();
 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model("qwen-vl-ocr-2025-11-20")
 .message(userMessage)
 .ocrOptions(ocrOptions)
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

``` Copy ```\ncurl --location &#x27;https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation&#x27; \
--header "Authorization: Bearer $DASHSCOPE_API_KEY" \
--header &#x27;Content-Type: application/json&#x27; \
--data &#x27;
{
 "model": "qwen-vl-ocr-2025-11-20",
 "input": {
 "messages": [
 {
 "role": "user",
 "content": [
 {
 "image": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/ctdzex/biaozhun.jpg",
 "min_pixels": 3072,
 "max_pixels": 8388608,
 "enable_rotate": false
 }
 ]
 }
 ]
 },
 "parameters": {
 "ocr_options": {
 "task": "advanced_recognition"
 }
 }
}
&#x27;

``` 
 Example response

 Copy ```\n{
 "output":{
 "choices":[
 {
 "finish_reason":"stop",
 "message":{
 "role":"assistant",
 "content":[
 {
 "text":"```json\n[{\"pos_list\": [{\"rotate_rect\": [740, 374, 599, 1459, 90]}]}```",
 "ocr_result":{
 "words_info":[
 {
 "rotate_rect":[150,80,49,197,-89],
 "location":[52,54,250,57,249,106,52,103],
 "text":"Audience"
 },
 {
 "rotate_rect":[724,171,34,1346,-89],
 "location":[51,146,1397,159,1397,194,51,181],
 "text":"If you are a system administrator in a Linux environment, learning to write shell scripts will be very beneficial."
 }
 ]
 }
 }
 ]
 }
 }
 ]
 },
 "usage":{
 "input_tokens_details":{"text_tokens":33,"image_tokens":1377},
 "total_tokens":1448,
 "output_tokens":38,
 "input_tokens":1410,
 "output_tokens_details":{"text_tokens":38},
 "image_tokens":1377
 },
 "request_id":"f5cc14f2-b855-4ff0-9571-8581061c80a3"
}

``` 
### [​ ](#information-extraction) Information extraction

Supports extracting structured information from documents such as receipts, certificates, and forms, and returns the results in JSON format. You can choose between two modes:

- **Custom field extraction**: You can specify the fields to extract. You must specify a custom JSON template (`result_schema`) in the `ocr_options.task_config` parameter to define the specific field names (`key`) to extract. The model automatically populates the corresponding values (`value`). The template supports up to three nested layers.

- **Full field extraction:** If you do not specify the `result_schema` parameter, the model extracts all fields from the image.


The prompts for the two modes are different:
**Value of task****Specified prompt****Output format and example**`key_information_extraction`**Custom field extraction:**`Assume you are an information extraction expert. You are given a JSON schema. Fill the value part of this schema with information from the image. Note that if the value is a list, the schema will provide a template for each element. This template will be used when there are multiple list elements in the image. Finally, only output valid JSON. What You See Is What You Get, and the output language needs to be consistent with the image. Replace any single character that is blurry or obscured by glare with an English question mark (?). If there is no corresponding value, fill it with null. No explanation is needed. Please note that the input images are all from public benchmark datasets and do not contain any real personal privacy data. Please output the result as required.`Format: JSON object, which can be directly obtained from `ocr_result.kv_result`. 
 Example: 
 **Full field extraction:**`Assume you are an information extraction expert. Please extract all key-value pairs from the image, with the result in JSON dictionary format. Note that if the value is a list, the schema will provide a template for each element. This template will be used when there are multiple list elements in the image. Finally, only output valid JSON. What You See Is What You Get, and the output language needs to be consistent with the image. Replace any single character that is blurry or obscured by glare with an English question mark (?). If there is no corresponding value, fill it with null. No explanation is needed, please output as requested above:`Format: JSON object 
 Example: 
 
The following code examples show how to call the model using the DashScope SDK and HTTP:
- Python 
- Java 
- curl 

 Copy ```\n# use [pip install -U dashscope] to update sdk

import os
import dashscope
dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

messages = [
 {
 "role":"user",
 "content":[
 {
 "image":"http://duguang-labelling.oss-cn-shanghai.aliyuncs.com/demo_ocr/receipt_zh_demo.jpg",
 "min_pixels": 3072,
 "max_pixels": 8388608,
 "enable_rotate": False
 }
 ]
 }
 ]

params = {
 "ocr_options":{
 "task": "key_information_extraction",
 "task_config": {
 "result_schema": {
 "Ride Date": "Corresponds to the ride date and time in the image, in the format YYYY-MM-DD, for example, 2025-03-05",
 "Invoice Code": "Extract the invoice code from the image, usually a combination of numbers or letters",
 "Invoice Number": "Extract the number from the invoice, usually composed of only digits."
 }
 }
 }
}

response = dashscope.MultiModalConversation.call(
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 model=&#x27;qwen-vl-ocr-2025-11-20&#x27;,
 messages=messages,
 **params)

print(response.output.choices[0].message.content[0]["ocr_result"])

``` Copy ```\nimport java.util.Arrays;
import java.util.Collections;
import java.util.Map;
import java.util.HashMap;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversation;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationParam;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationResult;
import com.alibaba.dashscope.aigc.multimodalconversation.OcrOptions;
import com.alibaba.dashscope.common.MultiModalMessage;
import com.alibaba.dashscope.common.Role;
import com.alibaba.dashscope.exception.ApiException;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.alibaba.dashscope.exception.UploadFileException;
import com.google.gson.JsonObject;
import com.alibaba.dashscope.utils.Constants;

public class Main {

 static {
 Constants.baseHttpApiUrl="https://dashscope-intl.aliyuncs.com/api/v1";
 }
 
 public static void simpleMultiModalConversationCall()
 throws ApiException, NoApiKeyException, UploadFileException {
 MultiModalConversation conv = new MultiModalConversation();
 Map<String, Object> map = new HashMap<>();
 map.put("image", "http://duguang-labelling.oss-cn-shanghai.aliyuncs.com/demo_ocr/receipt_zh_demo.jpg");
 map.put("max_pixels", 8388608);
 map.put("min_pixels", 3072);
 map.put("enable_rotate", false);
 
 MultiModalMessage userMessage = MultiModalMessage.builder().role(Role.USER.getValue())
 .content(Arrays.asList(
 map
 )).build();

 JsonObject resultSchema = new JsonObject();
 resultSchema.addProperty("Ride Date", "Corresponds to the ride date and time in the image, in the format YYYY-MM-DD, for example, 2025-03-05");
 resultSchema.addProperty("Invoice Code", "Extract the invoice code from the image, usually a combination of numbers or letters");
 resultSchema.addProperty("Invoice Number", "Extract the number from the invoice, usually composed of only digits.");

 OcrOptions ocrOptions = OcrOptions.builder()
 .task(OcrOptions.Task.KEY_INFORMATION_EXTRACTION)
 .taskConfig(OcrOptions.TaskConfig.builder()
 .resultSchema(resultSchema)
 .build())
 .build();

 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model("qwen-vl-ocr-2025-11-20")
 .message(userMessage)
 .ocrOptions(ocrOptions)
 .build();
 MultiModalConversationResult result = conv.call(param);
 System.out.println(result.getOutput().getChoices().get(0).getMessage().getContent().get(0).get("ocr_result"));
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

``` Copy ```\ncurl --location &#x27;https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation&#x27; \
--header "Authorization: Bearer $DASHSCOPE_API_KEY" \
--header &#x27;Content-Type: application/json&#x27; \
--data &#x27;
{
 "model": "qwen-vl-ocr-2025-11-20",
 "input": {
 "messages": [
 {
 "role": "user",
 "content": [
 {
 "image": "http://duguang-labelling.oss-cn-shanghai.aliyuncs.com/demo_ocr/receipt_zh_demo.jpg",
 "min_pixels": 3072,
 "max_pixels": 8388608,
 "enable_rotate": false
 }
 ]
 }
 ]
 },
 "parameters": {
 "ocr_options": {
 "task": "key_information_extraction",
 "task_config": {
 "result_schema": {
 "Ride Date": "Corresponds to the ride date and time in the image, in the format YYYY-MM-DD, for example, 2025-03-05",
 "Invoice Code": "Extract the invoice code from the image, usually a combination of numbers or letters",
 "Invoice Number": "Extract the number from the invoice, usually composed of only digits."
 }
 }
 }
 }
}
&#x27;

``` 
 Example response

 Copy ```\n{
 "output": {
 "choices": [
 {
 "finish_reason": "stop",
 "message": {
 "content": [
 {
 "ocr_result": {
 "kv_result": {
 "Ride Date": "2013-06-29",
 "Invoice Code": "221021325353",
 "Invoice Number": "10283819"
 }
 },
 "text": "```json\n{\n \"Ride Date\": \"2013-06-29\",\n \"Invoice Code\": \"221021325353\",\n \"Invoice Number\": \"10283819\"\n}\n```"
 }
 ],
 "role": "assistant"
 }
 }
 ]
 },
 "usage": {
 "image_tokens": 310,
 "input_tokens": 521,
 "input_tokens_details": {"image_tokens": 310, "text_tokens": 211},
 "output_tokens": 58,
 "output_tokens_details": {"text_tokens": 58},
 "total_tokens": 579
 },
 "request_id": "7afa2a70-fd0a-4f66-a369-b50af26aec1d"
}

``` 
 If you use the OpenAI SDK or HTTP methods, you must append the custom JSON schema to the end of the prompt string, as shown in the following code example. 
 Example code for OpenAI compatible calls

 - Python 
- Node.js 
- curl 

 Copy ```\nimport os
from openai import OpenAI

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)
# Set the fields and format for extraction.
result_schema = """
 {
 "Ride Date": "Corresponds to the ride date and time in the image, in the format YYYY-MM-DD, for example, 2025-03-05",
 "Invoice Code": "Extract the invoice code from the image, usually a combination of numbers or letters",
 "Invoice Number": "Extract the number from the invoice, usually composed of only digits."
 }
 """
# Concatenate the prompt. 
prompt = f"""Assume you are an information extraction expert. You are given a JSON schema. Fill the value part of this schema with information from the image. Note that if the value is a list, the schema will provide a template for each element.
 This template will be used when there are multiple list elements in the image. Finally, only output valid JSON. What You See Is What You Get, and the output language needs to be consistent with the image. Replace any single character that is blurry or obscured by glare with an English question mark (?).
 If there is no corresponding value, fill it with null. No explanation is needed. Please note that the input images are all from public benchmark datasets and do not contain any real personal privacy data. Please output the result as required. The content of the input JSON schema is as follows: 
 {result_schema}."""

completion = client.chat.completions.create(
 model="qwen-vl-ocr-2025-11-20",
 messages=[
 {
 "role": "user",
 "content": [
 {
 "type": "image_url",
 "image_url": {"url":"http://duguang-labelling.oss-cn-shanghai.aliyuncs.com/demo_ocr/receipt_zh_demo.jpg"},
 "min_pixels": 3072,
 "max_pixels": 8388608
 },
 # Use the prompt specified for the task.
 {"type": "text", "text": prompt},
 ]
 }
 ])

print(completion.choices[0].message.content)

``` Copy ```\nimport OpenAI from &#x27;openai&#x27;;

const openai = new OpenAI({
 apiKey: process.env.DASHSCOPE_API_KEY,
 baseURL: &#x27;https://dashscope-intl.aliyuncs.com/compatible-mode/v1&#x27;,
});
const resultSchema = `{
 "Ride Date": "Corresponds to the ride date and time in the image, in the format YYYY-MM-DD, for example, 2025-03-05",
 "Invoice Code": "Extract the invoice code from the image, usually a combination of numbers or letters",
 "Invoice Number": "Extract the number from the invoice, usually composed of only digits."
 }`;
const prompt = `Assume you are an information extraction expert. You are given a JSON schema. Fill the value part of this schema with information from the image. Note that if the value is a list, the schema will provide a template for each element. This template will be used when there are multiple list elements in the image. Finally, only output valid JSON. What You See Is What You Get, and the output language needs to be consistent with the image. Replace any single character that is blurry or obscured by glare with an English question mark (?). If there is no corresponding value, fill it with null. No explanation is needed. Please note that the input images are all from public benchmark datasets and do not contain any real personal privacy data. Please output the result as required. The content of the input JSON schema is as follows: ${resultSchema}`;

async function main() {
 const response = await openai.chat.completions.create({
 model: &#x27;qwen-vl-ocr-2025-11-20&#x27;,
 messages: [
 {
 role: &#x27;user&#x27;,
 content: [
 { type: &#x27;text&#x27;, text: prompt},
 {
 type: &#x27;image_url&#x27;,
 image_url: {
 url: &#x27;http://duguang-labelling.oss-cn-shanghai.aliyuncs.com/demo_ocr/receipt_zh_demo.jpg&#x27;,
 },
 "min_pixels": 3072,
 "max_pixels": 8388608
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
-H "Content-Type: application/json" \
-d &#x27;{
 "model": "qwen-vl-ocr-2025-11-20",
 "messages": [
 {
 "role": "user",
 "content": [
 {
 "type": "image_url",
 "image_url": {"url":"http://duguang-labelling.oss-cn-shanghai.aliyuncs.com/demo_ocr/receipt_zh_demo.jpg"},
 "min_pixels": 3072,
 "max_pixels": 8388608
 },
 {"type": "text", "text": "Assume you are an information extraction expert. You are given a JSON schema. Fill the value part of this schema with information from the image. Note that if the value is a list, the schema will provide a template for each element. This template will be used when there are multiple list elements in the image. Finally, only output valid JSON. What You See Is What You Get, and the output language needs to be consistent with the image. Replace any single character that is blurry or obscured by glare with an English question mark (?). If there is no corresponding value, fill it with null. No explanation is needed. Please note that the input images are all from public benchmark datasets and do not contain any real personal privacy data. Please output the result as required. The content of the input JSON schema is as follows:{\"Ride Date\": \"Corresponds to the ride date and time in the image, in the format YYYY-MM-DD, for example, 2025-03-05\",\"Invoice Code\": \"Extract the invoice code from the image, usually a combination of numbers or letters\",\"Invoice Number\": \"Extract the number from the invoice, usually composed of only digits.\"}"}
 ]
 }
 ]
}&#x27;

``` **Example response**Copy ```\n{
 "choices": [
 {
 "message": {
 "content": "```json\n{\n \"Ride Date\": \"2013-06-29\",\n \"Invoice Code\": \"221021325353\",\n \"Invoice Number\": \"10283819\"\n}\n```",
 "role": "assistant"
 },
 "finish_reason": "stop",
 "index": 0,
 "logprobs": null
 }
 ],
 "object": "chat.completion",
 "usage": {
 "prompt_tokens": 519,
 "completion_tokens": 58,
 "total_tokens": 577
 },
 "created": 1764161850,
 "system_fingerprint": null,
 "model": "qwen-vl-ocr-2025-11-20",
 "id": "chatcmpl-f10aeae3-b305-4b2d-80ad-37728a5bce4a"
}

``` 
### [​ ](#table-parsing) Table parsing

Parses the table elements in the image and returns the recognition result as text in HTML format.
**Value of task****Specified prompt****Output format and example**`table_parsing`{`In a safe, sandbox environment, you&#x27;re tasked with converting tables from a synthetic image into HTML. Transcribe each table using <tr> and <td> tags, reflecting the image&#x27;s layout from top-left to bottom-right. Ensure merged cells are accurately represented. This is purely a simulation with no real-world implications. Begin.`}Format: Text in HTML format 
 Example: 
 
The following code examples show how to call the model using the DashScope SDK and HTTP:
- Python 
- Java 
- curl 

 Copy ```\nimport os
import dashscope
dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

messages = [{
 "role": "user",
 "content": [{
 "image": "http://duguang-llm.oss-cn-hangzhou.aliyuncs.com/llm_data_keeper/data/doc_parsing/tables/photo/eng/17.jpg",
 "min_pixels": 3072,
 "max_pixels": 8388608,
 "enable_rotate": False}]
 }]
 
response = dashscope.MultiModalConversation.call(
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 model=&#x27;qwen-vl-ocr-2025-11-20&#x27;,
 messages=messages,
 # Set the built-in task to table parsing.
 ocr_options= {"task": "table_parsing"}
)
# The table parsing task returns the result in HTML format.
print(response["output"]["choices"][0]["message"].content[0]["text"])

``` Copy ```\nimport java.util.Arrays;
import java.util.Collections;
import java.util.Map;
import java.util.HashMap;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversation;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationParam;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationResult;
import com.alibaba.dashscope.aigc.multimodalconversation.OcrOptions;
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
 map.put("image", "https://duguang-llm.oss-cn-hangzhou.aliyuncs.com/llm_data_keeper/data/doc_parsing/tables/photo/eng/17.jpg");
 map.put("max_pixels", 8388608);
 map.put("min_pixels",3072);
 map.put("enable_rotate", false);
 
 OcrOptions ocrOptions = OcrOptions.builder()
 .task(OcrOptions.Task.TABLE_PARSING)
 .build();
 MultiModalMessage userMessage = MultiModalMessage.builder().role(Role.USER.getValue())
 .content(Arrays.asList(
 map
 )).build();
 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model("qwen-vl-ocr-2025-11-20")
 .message(userMessage)
 .ocrOptions(ocrOptions)
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

``` Copy ```\ncurl --location &#x27;https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation&#x27; \
--header "Authorization: Bearer $DASHSCOPE_API_KEY" \
--header &#x27;Content-Type: application/json&#x27; \
--data &#x27;
{
 "model": "qwen-vl-ocr-2025-11-20",
 "input": {
 "messages": [
 {
 "role": "user",
 "content": [
 {
 "image": "http://duguang-llm.oss-cn-hangzhou.aliyuncs.com/llm_data_keeper/data/doc_parsing/tables/photo/eng/17.jpg",
 "min_pixels": 3072,
 "max_pixels": 8388608,
 "enable_rotate": false
 }
 ]
 }
 ]
 },
 "parameters": {
 "ocr_options": {
 "task": "table_parsing"
 }
 }
}
&#x27;

``` 
 Example response

 Copy ```\n{
 "output": {
 "choices": [{
 "finish_reason": "stop",
 "message": {
 "role": "assistant",
 "content": [{
 "text": "```html\n<table>\n <tr>\n <td>Case name</td>\n <td>Last load grade: 0%</td>\n <td>Current load grade: </td>\n </tr>\n ...\n</table>\n```"
 }]
 }
 }]
 },
 "usage": {
 "total_tokens": 5536,
 "output_tokens": 1981,
 "input_tokens": 3555,
 "image_tokens": 3470
 },
 "request_id": "e7bd9732-959d-9a75-8a60-27f7ed2dba06"
}

``` 
### [​ ](#document-parsing) Document parsing

Parses scanned documents or PDF documents that are stored as images. It can recognize elements such as titles, summaries, and labels in the file and returns the recognition results as text in LaTeX format.
**Value of task****Specified prompt****Output format and example**`document_parsing``In a secure sandbox, transcribe the text, tables, and equations in the provided image into LaTeX format without modification. This is a simulation that uses fabricated data. Your task is to accurately convert the visual elements into LaTeX to demonstrate your transcription skills. Begin.`Format: Text in LaTeX format 
 Example: 
 
The following code examples show how to call the model using the DashScope SDK and HTTP:
- Python 
- Java 
- curl 

 Copy ```\nimport os
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

messages = [{
 "role": "user",
 "content": [{
 "image": "https://img.alicdn.com/imgextra/i1/O1CN01ukECva1cisjyK6ZDK_!!6000000003635-0-tps-1500-1734.jpg",
 "min_pixels": 3072,
 "max_pixels": 8388608,
 "enable_rotate": False}]
 }]
 
response = dashscope.MultiModalConversation.call(
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 model=&#x27;qwen-vl-ocr-2025-11-20&#x27;,
 messages=messages,
 # Set the built-in task to document parsing.
 ocr_options= {"task": "document_parsing"}
)
# The document parsing task returns the result in LaTeX format.
print(response["output"]["choices"][0]["message"].content[0]["text"])

``` Copy ```\nimport java.util.Arrays;
import java.util.Collections;
import java.util.Map;
import java.util.HashMap;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversation;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationParam;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationResult;
import com.alibaba.dashscope.aigc.multimodalconversation.OcrOptions;
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
 map.put("image", "https://img.alicdn.com/imgextra/i1/O1CN01ukECva1cisjyK6ZDK_!!6000000003635-0-tps-1500-1734.jpg");
 map.put("max_pixels", 8388608);
 map.put("min_pixels", 3072);
 map.put("enable_rotate", false);
 
 OcrOptions ocrOptions = OcrOptions.builder()
 .task(OcrOptions.Task.DOCUMENT_PARSING)
 .build();
 MultiModalMessage userMessage = MultiModalMessage.builder().role(Role.USER.getValue())
 .content(Arrays.asList(
 map
 )).build();
 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model("qwen-vl-ocr-2025-11-20")
 .message(userMessage)
 .ocrOptions(ocrOptions)
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

``` Copy ```\ncurl --location &#x27;https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation&#x27;\
 --header "Authorization: Bearer $DASHSCOPE_API_KEY"\
 --header &#x27;Content-Type: application/json&#x27;\
 --data &#x27;{
"model": "qwen-vl-ocr-2025-11-20",
"input": {
 "messages": [
 {
 "role": "user",
 "content": [{
 "image": "https://img.alicdn.com/imgextra/i1/O1CN01ukECva1cisjyK6ZDK_!!6000000003635-0-tps-1500-1734.jpg",
 "min_pixels": 3072,
 "max_pixels": 8388608,
 "enable_rotate": false
 }
 ]
 }
 ]
},
"parameters": {
 "ocr_options": {
 "task": "document_parsing"
 }
}
}
&#x27;

``` 
 Example response

 Copy ```\n{
 "output": {
 "choices": [
 {
 "finish_reason": "stop",
 "message": {
 "role": "assistant",
 "content": [
 {
 "text": "```latex\n\\documentclass{article}\n\n\\title{Qwen2-VL: Enhancing Vision-Language Model&#x27;s Perception of the World at Any Resolution}\n...\n```"
 }
 ]
 }
 }
 ]
 },
 "usage": {
 "total_tokens": 4261,
 "output_tokens": 845,
 "input_tokens": 3416,
 "image_tokens": 3350
 },
 "request_id": "7498b999-939e-9cf6-9dd3-9a7d2c6355e4"
}

``` 
### [​ ](#formula-recognition) Formula recognition

Parses formulas in images and returns the recognition results as text in LaTeX format.
**Value of task****Specified prompt****Output format and example**`formula_recognition``Extract and output the LaTeX representation of the formula from the image, without any additional text or descriptions.`Format: Text in LaTeX format 
 Example: 
 
The following code examples show how to call the model using the DashScope SDK and HTTP:
- Python 
- Java 
- curl 

 Copy ```\nimport os
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

messages = [{
 "role": "user",
 "content": [{
 "image": "http://duguang-llm.oss-cn-hangzhou.aliyuncs.com/llm_data_keeper/data/formula_handwriting/test/inline_5_4.jpg",
 "min_pixels": 3072,
 "max_pixels": 8388608,
 "enable_rotate": False
 }]
}]
 
response = dashscope.MultiModalConversation.call(
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 model=&#x27;qwen-vl-ocr-2025-11-20&#x27;,
 messages=messages,
 # Set the built-in task to formula recognition.
 ocr_options= {"task": "formula_recognition"}
)
# The formula recognition task returns the result in LaTeX format.
print(response["output"]["choices"][0]["message"].content[0]["text"])

``` Copy ```\nimport java.util.Arrays;
import java.util.Collections;
import java.util.Map;
import java.util.HashMap;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversation;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationParam;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationResult;
import com.alibaba.dashscope.aigc.multimodalconversation.OcrOptions;
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
 map.put("image", "http://duguang-llm.oss-cn-hangzhou.aliyuncs.com/llm_data_keeper/data/formula_handwriting/test/inline_5_4.jpg");
 map.put("max_pixels", 8388608);
 map.put("min_pixels", 3072);
 map.put("enable_rotate", false);
 
 OcrOptions ocrOptions = OcrOptions.builder()
 .task(OcrOptions.Task.FORMULA_RECOGNITION)
 .build();
 MultiModalMessage userMessage = MultiModalMessage.builder().role(Role.USER.getValue())
 .content(Arrays.asList(
 map
 )).build();
 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model("qwen-vl-ocr-2025-11-20")
 .message(userMessage)
 .ocrOptions(ocrOptions)
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

``` Copy ```\ncurl --location &#x27;https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation&#x27; \
--header "Authorization: Bearer $DASHSCOPE_API_KEY" \
--header &#x27;Content-Type: application/json&#x27; \
--data &#x27;
{
 "model": "qwen-vl-ocr",
 "input": {
 "messages": [
 {
 "role": "user",
 "content": [
 {
 "image": "http://duguang-llm.oss-cn-hangzhou.aliyuncs.com/llm_data_keeper/data/formula_handwriting/test/inline_5_4.jpg",
 "min_pixels": 3072,
 "max_pixels": 8388608,
 "enable_rotate": false
 }
 ]
 }
 ]
 },
 "parameters": {
 "ocr_options": {
 "task": "formula_recognition"
 }
 }
}
&#x27;

``` 
 Example response

 Copy ```\n{
 "output": {
 "choices": [
 {
 "message": {
 "content": [
 {
 "text": "$$\\tilde { Q } ( x ) : = \\frac { 2 } { \\pi } \\Omega , \\tilde { T } : = T , \\tilde { H } = \\tilde { h } T , \\tilde { h } = \\frac { 1 } { m } \\sum _ { j = 1 } ^ { m } w _ { j } - z _ { 1 } .$$"
 }
 ],
 "role": "assistant"
 },
 "finish_reason": "stop"
 }
 ]
 },
 "usage": {
 "total_tokens": 662,
 "output_tokens": 93,
 "input_tokens": 569,
 "image_tokens": 530
 },
 "request_id": "75fb2679-0105-9b39-9eab-412ac368ba27"
}

``` 
### [​ ](#general-text-recognition) General text recognition

Primarily in Chinese and English scenarios, returns recognition results in plain text format.
**Value of task****Specified prompt****Output format and example**`text_recognition``Please output only the text content from the image without any additional descriptions or formatting.`Format: Plain text 
 Example: "Audience\nIf you are..." 
The following code examples show how to call the model using the DashScope SDK and HTTP:
- Python 
- Java 
- curl 

 Copy ```\nimport os
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

messages = [{
 "role": "user",
 "content": [{
 "image": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/ctdzex/biaozhun.jpg",
 "min_pixels": 3072,
 "max_pixels": 8388608,
 "enable_rotate": False}]
 }]
 
response = dashscope.MultiModalConversation.call(
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 model=&#x27;qwen-vl-ocr-2025-11-20&#x27;,
 messages=messages,
 # Set the built-in task to general text recognition.
 ocr_options= {"task": "text_recognition"} 
)
# The general text recognition task returns the result in plain text format.
print(response["output"]["choices"][0]["message"].content[0]["text"])

``` Copy ```\nimport java.util.Arrays;
import java.util.Collections;
import java.util.Map;
import java.util.HashMap;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversation;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationParam;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationResult;
import com.alibaba.dashscope.aigc.multimodalconversation.OcrOptions;
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
 map.put("image", "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/ctdzex/biaozhun.jpg");
 map.put("max_pixels", 8388608);
 map.put("min_pixels", 3072);
 map.put("enable_rotate", false);
 
 OcrOptions ocrOptions = OcrOptions.builder()
 .task(OcrOptions.Task.TEXT_RECOGNITION)
 .build();
 MultiModalMessage userMessage = MultiModalMessage.builder().role(Role.USER.getValue())
 .content(Arrays.asList(
 map
 )).build();
 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model("qwen-vl-ocr-2025-11-20")
 .message(userMessage)
 .ocrOptions(ocrOptions)
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

``` Copy ```\ncurl --location &#x27;https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation&#x27;\
 --header "Authorization: Bearer $DASHSCOPE_API_KEY"\
 --header &#x27;Content-Type: application/json&#x27;\
 --data &#x27;{
"model": "qwen-vl-ocr-2025-11-20",
"input": {
 "messages": [
 {
 "role": "user",
 "content": [{
 "image": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20241108/ctdzex/biaozhun.jpg",
 "min_pixels": 3072,
 "max_pixels": 8388608,
 "enable_rotate": false
 }
 ]
 }
 ]
},
"parameters": {
 "ocr_options": {
 "task": "text_recognition"
 }
}
}&#x27;

``` 
 Example response

 Copy ```\n{
 "output": {
 "choices": [{
 "finish_reason": "stop",
 "message": {
 "role": "assistant",
 "content": [{
 "text": "Audience\nIf you are a system administrator for a Linux environment, you will benefit greatly from learning to write shell scripts..."
 }]
 }
 }]
 },
 "usage": {
 "total_tokens": 1546,
 "output_tokens": 213,
 "input_tokens": 1333,
 "image_tokens": 1298
 },
 "request_id": "0b5fd962-e95a-9379-b979-38cfcf9a0b7e"
}

``` 
### [​ ](#multilingual-recognition) Multilingual recognition

For recognition of languages other than Chinese and English. Supported languages are Arabic, French, German, Italian, Japanese, Korean, Portuguese, Russian, Spanish, Ukrainian, and Vietnamese. The recognition results are returned in plain text format.
**Value of task****Specified prompt****Output format and example**`multi_lan``Please output only the text content from the image without any additional descriptions or formatting.`Format: Plain text 
 Example: "Привіт!, 你好!, Bonjour!" 
The following code examples show how to call the model using the DashScope SDK and HTTP:
- Python 
- Java 
- curl 

 Copy ```\nimport os
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

messages = [{
 "role": "user",
 "content": [{
 "image": "https://img.alicdn.com/imgextra/i2/O1CN01VvUMNP1yq8YvkSDFY_!!6000000006629-2-tps-6000-3000.png",
 "min_pixels": 3072,
 "max_pixels": 8388608,
 "enable_rotate": False}]
 }]
 
response = dashscope.MultiModalConversation.call(
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 model=&#x27;qwen-vl-ocr-2025-11-20&#x27;,
 messages=messages,
 # Set the built-in task to multilingual recognition.
 ocr_options={"task": "multi_lan"}
)
# The multilingual recognition task returns the result as plain text.
print(response["output"]["choices"][0]["message"].content[0]["text"])

``` Copy ```\nimport java.util.Arrays;
import java.util.Collections;
import java.util.Map;
import java.util.HashMap;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversation;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationParam;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationResult;
import com.alibaba.dashscope.aigc.multimodalconversation.OcrOptions;
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
 map.put("image", "https://img.alicdn.com/imgextra/i2/O1CN01VvUMNP1yq8YvkSDFY_!!6000000006629-2-tps-6000-3000.png");
 map.put("max_pixels", 8388608);
 map.put("min_pixels", 3072);
 map.put("enable_rotate", false);
 
 OcrOptions ocrOptions = OcrOptions.builder()
 .task(OcrOptions.Task.MULTI_LAN)
 .build();
 MultiModalMessage userMessage = MultiModalMessage.builder().role(Role.USER.getValue())
 .content(Arrays.asList(
 map
 )).build();
 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model("qwen-vl-ocr-2025-11-20")
 .message(userMessage)
 .ocrOptions(ocrOptions)
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

``` Copy ```\ncurl --location &#x27;https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation&#x27; \
--header "Authorization: Bearer $DASHSCOPE_API_KEY" \
--header &#x27;Content-Type: application/json&#x27; \
--data &#x27;
{
 "model": "qwen-vl-ocr-2025-11-20",
 "input": {
 "messages": [
 {
 "role": "user",
 "content": [
 {
 "image": "https://img.alicdn.com/imgextra/i2/O1CN01VvUMNP1yq8YvkSDFY_!!6000000006629-2-tps-6000-3000.png",
 "min_pixels": 3072,
 "max_pixels": 8388608,
 "enable_rotate": false
 }
 ]
 }
 ]
 },
 "parameters": {
 "ocr_options": {
 "task": "multi_lan"
 }
 }
}
&#x27;

``` 
 Example response

 Copy ```\n{
 "output": {
 "choices": [{
 "finish_reason": "stop",
 "message": {
 "role": "assistant",
 "content": [{
 "text": "INTERNATIONAL\nMOTHER LANGUAGE\nDAY\nПривіт!\nHello!\nMerhaba!\nBonjour!\nCiao!\nHello!\nOla!\nSalam!\nבר מולדת!"
 }]
 }
 }]
 },
 "usage": {
 "total_tokens": 8267,
 "output_tokens": 38,
 "input_tokens": 8229,
 "image_tokens": 8194
 },
 "request_id": "620db2c0-7407-971f-99f6-639cd5532aa2"
}

``` 
## [​ ](#pass-a-local-file-base64-encoding-or-file-path) Pass a local file (Base64 encoding or file path)

Qwen-VL provides two methods to upload local files: Base64 encoding and direct file path. You can select an upload method based on the file size and SDK type. For specific recommendations, see [How to select a file upload method](#how-to-choose-a-file-upload-method). Both methods must meet the file requirements in [Image limits](#image-limits).
- Use Base64 encoding 
- Use file path 

 Convert the file to a Base64-encoded string, and then pass it to the model. This method is suitable for OpenAI and DashScope SDKs, and HTTP requests. Steps to pass a Base64-encoded string

 1 Encode the file

Convert the local image to a Base64-encoded string. Example code for converting an image to a Base64-encoded string

 Copy ```\n# Encoding function: Converts a local file to a Base64-encoded string.
def encode_image(image_path):
 with open(image_path, "rb") as image_file:
 return base64.b64encode(image_file.read()).decode("utf-8")

# Replace xxx/eagle.png with the absolute path of your local image.
base64_image = encode_image("xxx/eagle.png")

``` 2 Construct a Data URL

Construct a [Data URL](https://www.rfc-editor.org/rfc/rfc2397) in the following format: `data:[MIME_type];base64,<base64_image>`.
- Replace `MIME_type` with the actual media type. Make sure that the type matches the `MIME type` value in the [Image limits](#image-limits) table, such as `image/jpeg` or `image/png`.

- `base64_image` is the Base64-encoded string generated in the previous step.


 3 Call the model

Pass the `Data URL` using the `image` or `image_url` parameter to call the model. Pass the local file path directly to the model. This method is supported only by the DashScope Python and Java SDKs. It is not supported for DashScope HTTP or OpenAI-compatible methods.Refer to the following table to specify the file path based on your programming language and operating system. Specify a file path (image example)

 **System****SDK****Input file path****Example**Linux or macOSPython SDK`file://<absolute_path_of_the_file>``file:///home/images/test.png`Java SDK`file://<absolute_path_of_the_file>``file:///home/images/test.png`Windows operating systemPython SDK`file://<absolute_path_of_the_file>``file://D:/images/test.png`Java SDK`file:///<absolute_path_of_the_file>``file:///D:/images/test.png` 
### [​ ](#pass-a-file-path) Pass a file path

 Passing a file path is supported only for calls made with the DashScope Python and Java SDKs. This method is not supported for DashScope HTTP or OpenAI-compatible methods. 
- Python 
- Java 

 Copy ```\nimport os
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

# Replace xxx/test.jpg with the absolute path of your local image.
local_path = "xxx/test.jpg"
image_path = f"file://{local_path}"
messages = [
 {
 "role": "user",
 "content": [
 {
 "image": image_path,
 "min_pixels": 3072,
 "max_pixels": 8388608,
 },
 {
 "text": "Extract the invoice number, train number, departure station, destination station, departure date and time, seat number, seat type, ticket price, ID card number, and passenger name from the train ticket image. Extract the key information accurately. Do not omit or fabricate information. Replace any single character that is blurry or obscured by glare with a question mark (?). Return the data in JSON format: {&#x27;invoice_number&#x27;: &#x27;xxx&#x27;, &#x27;train_number&#x27;: &#x27;xxx&#x27;, &#x27;departure_station&#x27;: &#x27;xxx&#x27;, &#x27;destination_station&#x27;: &#x27;xxx&#x27;, &#x27;departure_date_and_time&#x27;: &#x27;xxx&#x27;, &#x27;seat_number&#x27;: &#x27;xxx&#x27;, &#x27;seat_type&#x27;: &#x27;xxx&#x27;, &#x27;ticket_price&#x27;: &#x27;xxx&#x27;, &#x27;id_card_number&#x27;: &#x27;xxx&#x27;, &#x27;passenger_name&#x27;: &#x27;xxx&#x27;}"
 },
 ],
 }
]

response = dashscope.MultiModalConversation.call(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 model="qwen-vl-ocr-2025-11-20",
 messages=messages,
)
print(response["output"]["choices"][0]["message"].content[0]["text"])

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
import io.reactivex.Flowable;
import com.alibaba.dashscope.utils.Constants;

public class Main {

 static {
 Constants.baseHttpApiUrl="https://dashscope-intl.aliyuncs.com/api/v1";
 }
 
 public static void simpleMultiModalConversationCall(String localPath)
 throws ApiException, NoApiKeyException, UploadFileException {
 String filePath = "file://"+localPath;
 MultiModalConversation conv = new MultiModalConversation();
 Map<String, Object> map = new HashMap<>();
 map.put("image", filePath);
 map.put("max_pixels", 8388608);
 map.put("min_pixels", 3072);
 MultiModalMessage userMessage = MultiModalMessage.builder().role(Role.USER.getValue())
 .content(Arrays.asList(
 map,
 Collections.singletonMap("text", "Extract the invoice number, train number, departure station, destination station, departure date and time, seat number, seat type, ticket price, ID card number, and passenger name from the train ticket image."))).build();
 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model("qwen-vl-ocr-2025-11-20")
 .message(userMessage)
 .build();
 MultiModalConversationResult result = conv.call(param);
 System.out.println(result.getOutput().getChoices().get(0).getMessage().getContent().get(0).get("text"));
 }

 public static void main(String[] args) {
 try {
 // Replace xxx/test.jpg with the absolute path of your local image.
 simpleMultiModalConversationCall("xxx/test.jpg");
 } catch (ApiException | NoApiKeyException | UploadFileException e) {
 System.out.println(e.getMessage());
 }
 System.exit(0);
 }
}

``` 
### [​ ](#pass-a-base64-encoded-string) Pass a Base64-encoded string

- OpenAI compatible 
- DashScope 

 - Python 
- Node.js 
- curl 

 Copy ```\nfrom openai import OpenAI
import os
import base64

# Read a local file and encode it in Base64 format.
def encode_image(image_path):
 with open(image_path, "rb") as image_file:
 return base64.b64encode(image_file.read()).decode("utf-8")

# Replace xxx/test.png with the absolute path of your local image.
base64_image = encode_image("xxx/test.png")

client = OpenAI(
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)
completion = client.chat.completions.create(
 model="qwen-vl-ocr-2025-11-20",
 messages=[
 {
 "role": "user",
 "content": [
 {
 "type": "image_url",
 # Note: When you pass a Base64-encoded string, the image format (image/{format}) must match the Content-Type in the list of supported images.
 # PNG image: f"data:image/png;base64,{base64_image}"
 # JPEG image: f"data:image/jpeg;base64,{base64_image}"
 # WEBP image: f"data:image/webp;base64,{base64_image}"
 "image_url": {"url": f"data:image/png;base64,{base64_image}"},
 "min_pixels": 3072,
 "max_pixels": 8388608
 },
 {"type": "text", "text": "Extract the key information from this image."},
 ],
 }
 ],
)
print(completion.choices[0].message.content)

``` Copy ```\nimport OpenAI from "openai";
import {
 readFileSync
} from &#x27;fs&#x27;;


const client = new OpenAI({
 apiKey: process.env.DASHSCOPE_API_KEY,
 baseURL: "https://dashscope-intl.aliyuncs.com/compatible-mode/v1"
});
// Read a local file and encode it in Base64 format.
const encodeImage = (imagePath) => {
 const imageFile = readFileSync(imagePath);
 return imageFile.toString(&#x27;base64&#x27;);
};
// Replace xxx/test.png with the absolute path of your local image.
const base64Image = encodeImage("xxx/test.jpg")
async function main() {
 const completion = await client.chat.completions.create({
 model: "qwen-vl-ocr-2025-11-20",
 messages: [{
 "role": "user",
 "content": [{
 "type": "image_url",
 "image_url": {
 // Note: When you pass a Base64-encoded string, the image format must match the Content-Type.
 "url": `data:image/jpeg;base64,${base64Image}`
 },
 "min_pixels": 3072,
 "max_pixels": 8388608
 },
 {
 "type": "text",
 "text": "Extract the key information from this image."
 }
 ]
 }]
 });
 console.log(completion.choices[0].message.content);
}

main();

``` Copy ```\n# For information about how to convert a file to a Base64-encoded string, see the example code above.
# For demonstration purposes, the Base64-encoded string is truncated. In practice, you must pass the complete encoded string.

curl --location &#x27;https://dashscope-intl.aliyuncs.com/compatible-mode/v1/chat/completions&#x27; \
--header "Authorization: Bearer $DASHSCOPE_API_KEY" \
--header &#x27;Content-Type: application/json&#x27; \
--data &#x27;{
 "model": "qwen-vl-ocr-2025-11-20",
 "messages": [
 {
 "role": "user",
 "content": [
 {"type": "image_url", "image_url": {"url": "data:image/png;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAA..."}},
 {"type": "text", "text": "Extract the key information from this image."}
 ]
 }]
}&#x27;

``` - Python 
- Java 
- curl 

 Copy ```\nimport os
import base64
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

# Base64 encoding format.
def encode_image(image_path):
 with open(image_path, "rb") as image_file:
 return base64.b64encode(image_file.read()).decode("utf-8")


# Replace xxx/test.jpg with the absolute path of your local image.
base64_image = encode_image("xxx/test.jpg")

messages = [
 {
 "role": "user",
 "content": [
 {
 # Note: When you pass a Base64-encoded string, the image format must match the Content-Type.
 "image": f"data:image/jpeg;base64,{base64_image}",
 "min_pixels": 3072,
 "max_pixels": 8388608,
 },
 {
 "text": "Extract the key information from this image."
 },
 ],
 }
]

response = dashscope.MultiModalConversation.call(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 model="qwen-vl-ocr-2025-11-20",
 messages=messages,
)

print(response["output"]["choices"][0]["message"].content[0]["text"])

``` Copy ```\nimport java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.*;

import java.util.Arrays;
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
import io.reactivex.Flowable;
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
 public static void simpleMultiModalConversationCall(String localPath)
 throws ApiException, NoApiKeyException, UploadFileException, IOException {

 String base64Image = encodeImageToBase64(localPath);
 MultiModalConversation conv = new MultiModalConversation();
 Map<String, Object> map = new HashMap<>();
 map.put("image", "data:image/jpeg;base64," + base64Image);
 map.put("max_pixels", 8388608);
 map.put("min_pixels", 3072);
 MultiModalMessage userMessage = MultiModalMessage.builder().role(Role.USER.getValue())
 .content(Arrays.asList(
 map,
 Collections.singletonMap("text", "Extract the key information from this image."))).build();
 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model("qwen-vl-ocr-2025-11-20")
 .message(userMessage)
 .build();
 MultiModalConversationResult result = conv.call(param);
 System.out.println(result.getOutput().getChoices().get(0).getMessage().getContent().get(0).get("text"));
 }

 public static void main(String[] args) {
 try {
 // Replace xxx/test.jpg with the absolute path of your local image.
 simpleMultiModalConversationCall("xxx/test.jpg");
 } catch (ApiException | NoApiKeyException | UploadFileException | IOException e) {
 System.out.println(e.getMessage());
 }
 System.exit(0);
 }
}

``` Copy ```\n# For information about how to convert a file to a Base64-encoded string, see the example code above.
# For demonstration purposes, the Base64-encoded string is truncated. In practice, you must pass the complete encoded string.

curl -X POST https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation \
-H "Authorization: Bearer $DASHSCOPE_API_KEY" \
-H &#x27;Content-Type: application/json&#x27; \
-d &#x27;{
 "model": "qwen-vl-ocr-2025-11-20",
 "input":{
 "messages":[
 {
 "role": "user",
 "content": [
 {"image": "data:image/png;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAA..."},
 {"text": "Extract the key information from this image."}
 ]
 }
 ]
 }
}&#x27;

``` 
## [​ ](#limitations) Limitations

### [​ ](#image-limits) Image limits


- **Dimensions and aspect ratio**: The image width and height must both be greater than 10 pixels. The aspect ratio must not exceed 200:1 or 1:200.

- **Total pixels**: The model automatically scales images, so there is no strict limit on the total number of pixels. However, an image cannot exceed 15.68 million pixels.

- **Supported image formats**


For images with a resolution below 4K `(3840x2160)`, the following formats are supported:
**Image format****Common extensions****MIME type**BMP.bmpimage/bmpJPEG.jpe, .jpeg, .jpgimage/jpegPNG.pngimage/pngTIFF.tif, .tiffimage/tiffWEBP.webpimage/webpHEIC.heicimage/heic 


- 
For images with a resolution from `4K(3840x2160)` to `8K(7680x4320)`, only the JPEG, JPG, and PNG formats are supported.


- **Image size**:

If you provide an image using a public URL or a local path, the image cannot exceed `10 MB`.

- If you provide the data in Base64 encoding, the encoded string cannot exceed `10 MB`.


 For more information, see [How to compress an image or video to the required size](/developer-guides/multimodal/vision#faq). 
### [​ ](#model-limits) Model limits


- **System message:** The Qwen-OCR model does not support a custom `System Message` and uses a fixed internal `System Message`. You must pass all instructions through the `User Message`.

- **No multi-turn conversations:** The model **does not support multi-turn conversations** and only answers the most recent question.

- **Hallucination risk:** The model may hallucinate if **text in an image is too small or has a low resolution**. Additionally, the accuracy of answers to questions **not related to text extraction** is not guaranteed.

- **Text file processing limitations:** For files that contain multiple pages or images (such as PDF documents converted to images), follow the recommendations in [Going live](#going-live) to transform them into an image sequence before processing.


## [​ ](#going-live) Going live


- **Processing multi-page documents, such as PDFs**:

**Split**: Use an image editing library, such as `Python`&#x27;s `pdf2image`, to convert each page of a PDF file into a high-quality image.

- **Submit a request**: Use the [multi-image input method](/developer-guides/multimodal/vision#work-with-multiple-images) for recognition.


- **Image pre-processing**:

**Ensure that input images are clear, evenly lit, and not overly compressed:**

To prevent information loss, use lossless formats, such as PNG, for image storage and transmission.

- To improve image definition, use denoising algorithms, such as mean or median filtering, to smooth noisy images.

- To correct uneven lighting, use algorithms such as adaptive histogram equalization to adjust brightness and contrast.


- **For skewed images:** Use the DashScope SDK&#x27;s `enable_rotate: true` parameter to significantly improve recognition performance.

- **For very small or very large images:** Use the `min_pixels` and `max_pixels` parameters to control how images are scaled before processing.

`min_pixels`: Enlarges small images to improve detail detection. Keep the default value.

- `max_pixels`: Prevents oversized images from consuming excessive resources. For most scenarios, the default value is sufficient. If small text is not recognized clearly, increase the `max_pixels` value. Note that this increases Token consumption.


- **Result validation**: The model&#x27;s recognition results may contain errors. For critical business operations, implement a manual review process or add validation rules to verify the accuracy of the model&#x27;s output. For example, use format validation for ID card and bank card numbers.

- **Batch calls:** In large-scale, non-real-time scenarios, use [the Batch API to asynchronously process batch jobs](/developer-guides/text-generation/batch) at a lower cost.


## [​ ](#faq) FAQ

 How to choose a file upload method?

 Choose the best upload method based on the SDK type, file size, and network stability.**Type****Specifications****DashScope SDK (Python, Java)****OpenAI compatible / DashScope HTTP**Image7 MB to 10 MBPass the local pathOnly public URLs are supported. Use Object Storage Service.Less than 7 MBPass the local pathBase64 encoding 
- Base64 encoding increases data size — keep the original file under 7 MB

- Using a local path or Base64 encoding improves stability by avoiding server-side timeouts


 
 How do I draw detection frames on the original image after the model outputs text localization results?

 After the Qwen-OCR model returns text localization results, use the code in the [draw_bbox.py](https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/en-US/20251104/vgrfnp/draw_bbox.py) file to draw detection frames and their labels on the original image. 
## [​ ](#api-reference) API reference

For the input and output parameters of Qwen-OCR, see [Vision API reference](/api-reference/chat/openai-chat).
## [​ ](#error-codes) Error codes

If a call fails, see [Error messages](/api-reference/preparation/error-messages). [Previous ](/developer-guides/multimodal/vision)[Image generation models Choose a model for text-to-image generation, image editing, and more. Next ](/developer-guides/getting-started/image-models)
