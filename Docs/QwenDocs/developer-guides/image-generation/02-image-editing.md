# Image editing

> **Source:** https://docs.qwencloud.com/developer-guides/image-generation/image-editing

Modify images via text

 Copy page ## [​ ](#getting-started) Getting started

This example shows how to use `qwen-image-2.0-pro` to generate two edited images from three input images and a prompt.
 Input prompt: The girl in Image 1 wears the black dress from Image 2 and sits in the pose from Image 3. 
**Input image 1****Input image 2****Input image 3****Output images (multiple images)** 
Before making a call, [get an API key](/api-reference/preparation/api-key) and export the API key as an environment variable.
To call the API using the SDK, [install the DashScope SDK](/api-reference/preparation/install-sdk). The SDK is available for Python and Java.
The Qwen image editing models support one to three input images. The `qwen-image-2.0`, `qwen-image-edit-max`, and `qwen-image-edit-plus` series can generate one to six images. `qwen-image-edit` can generate only one image. The URLs for the generated images are **valid for 24 hours**. Download the images to your local device promptly.
- Python 
- Java 
- curl 

 Copy ```\nimport json
import os
import dashscope
from dashscope import MultiModalConversation

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

# The model supports one to three input images.
messages = [
 {
 "role": "user",
 "content": [
 {"image": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250925/thtclx/input1.png"},
 {"image": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250925/iclsnx/input2.png"},
 {"image": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250925/gborgw/input3.png"},
 {"text": "Make the girl from Image 1 wear the black dress from Image 2 and sit in the pose from Image 3."}
 ]
 }
]

# Replace with your API key if the environment variable is not set: api_key="sk-xxx"
api_key = os.getenv("DASHSCOPE_API_KEY")

# qwen-image-2.0, qwen-image-edit-max, and qwen-image-edit-plus series support outputting 1 to 6 images. This example shows how to output 2 images.
response = MultiModalConversation.call(
 api_key=api_key,
 model="qwen-image-2.0-pro",
 messages=messages,
 stream=False,
 n=2,
 watermark=False,
 negative_prompt=" ",
 prompt_extend=True,
 size="1024*1536",
)

if response.status_code == 200:
 # To view the full response, uncomment the following line.
 # print(json.dumps(response, ensure_ascii=False))
 for i, content in enumerate(response.output.choices[0].message.content):
 print(f"URL of output image {i+1}: {content[&#x27;image&#x27;]}")
else:
 print(f"HTTP status code: {response.status_code}")
 print(f"Error code: {response.code}")
 print(f"Error message: {response.message}")
 print("For more information, see the documentation: https://docs.qwencloud.com/api-reference/preparation/error-messages")

``` Response example

 Copy ```\n{
 "status_code": 200,
 "request_id": "fa41f9f9-3cb6-434d-a95d-4ae6b9xxxxxx",
 "code": "",
 "message": "",
 "output": {
 "text": null,
 "finish_reason": null,
 "choices": [
 {
 "finish_reason": "stop",
 "message": {
 "role": "assistant",
 "content": [
 {
 "image": "https://dashscope-result-hz.oss-cn-hangzhou.aliyuncs.com/xxx.png?Expires=xxx"
 },
 {
 "image": "https://dashscope-result-hz.oss-cn-hangzhou.aliyuncs.com/xxx.png?Expires=xxx"
 }
 ]
 }
 }
 ],
 "audio": null
 },
 "usage": {
 "input_tokens": 0,
 "output_tokens": 0,
 "characters": 0,
 "height": 1536,
 "image_count": 2,
 "width": 1024
 }
}

``` Copy ```\nimport com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversation;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationParam;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationResult;
import com.alibaba.dashscope.common.MultiModalMessage;
import com.alibaba.dashscope.common.Role;
import com.alibaba.dashscope.exception.ApiException;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.alibaba.dashscope.exception.UploadFileException;
import com.alibaba.dashscope.utils.JsonUtils;
import com.alibaba.dashscope.utils.Constants;

import java.io.IOException;
import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import java.util.List;

public class QwenImageEdit {

 static {
 Constants.baseHttpApiUrl = "https://dashscope-intl.aliyuncs.com/api/v1";
 }

 // If you have not configured the environment variable, replace the following line with your API key: apiKey="sk-xxx"
 static String apiKey = System.getenv("DASHSCOPE_API_KEY");

 public static void call() throws ApiException, NoApiKeyException, UploadFileException, IOException {

 MultiModalConversation conv = new MultiModalConversation();

 // The model supports one to three input images.
 MultiModalMessage userMessage = MultiModalMessage.builder().role(Role.USER.getValue())
 .content(Arrays.asList(
 Collections.singletonMap("image", "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250925/thtclx/input1.png"),
 Collections.singletonMap("image", "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250925/iclsnx/input2.png"),
 Collections.singletonMap("image", "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250925/gborgw/input3.png"),
 Collections.singletonMap("text", "Make the girl from Image 1 wear the black dress from Image 2 and sit in the pose from Image 3.")
 )).build();
 // qwen-image-2.0, qwen-image-edit-max, and qwen-image-edit-plus series support outputting 1 to 6 images. This example shows how to output 2 images.
 Map<String, Object> parameters = new HashMap<>();
 parameters.put("watermark", false);
 parameters.put("negative_prompt", " ");
 parameters.put("n", 2);
 parameters.put("prompt_extend", true);
 parameters.put("size", "1024*1536");

 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(apiKey)
 .model("qwen-image-2.0-pro")
 .messages(Collections.singletonList(userMessage))
 .parameters(parameters)
 .build();

 MultiModalConversationResult result = conv.call(param);
 // To view the complete response, uncomment the following line.
 // System.out.println(JsonUtils.toJson(result));
 List<Map<String, Object>> contentList = result.getOutput().getChoices().get(0).getMessage().getContent();
 int imageIndex = 1;
 for (Map<String, Object> content : contentList) {
 if (content.containsKey("image")) {
 System.out.println("URL of output image " + imageIndex + ": " + content.get("image"));
 imageIndex++;
 }
 }
 }

 public static void main(String[] args) {
 try {
 call();
 } catch (ApiException | NoApiKeyException | UploadFileException | IOException e) {
 System.out.println(e.getMessage());
 }
 }
}

``` Sample response

 Copy ```\n{
 "requestId": "46281da9-9e02-941c-ac78-be88b8xxxxxx",
 "usage": {
 "image_count": 2,
 "width": 1024,
 "height": 1536
 },
 "output": {
 "choices": [
 {
 "finish_reason": "stop",
 "message": {
 "role": "assistant",
 "content": [
 {
 "image": "https://dashscope-result-sz.oss-cn-shenzhen.aliyuncs.com/xxx.png?Expires=xxx"
 },
 {
 "image": "https://dashscope-result-sz.oss-cn-shenzhen.aliyuncs.com/xxx.png?Expires=xxx"
 }
 ]
 }
 }
 ]
 }
}

``` Copy ```\ncurl --location &#x27;https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation&#x27; \
--header &#x27;Content-Type: application/json&#x27; \
--header "Authorization: Bearer $DASHSCOPE_API_KEY" \
--data &#x27;{
 "model": "qwen-image-2.0-pro",
 "input": {
 "messages": [
 {
 "role": "user",
 "content": [
 {
 "image": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250925/thtclx/input1.png"
 },
 {
 "image": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250925/iclsnx/input2.png"
 },
 {
 "image": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250925/gborgw/input3.png"
 },
 {
 "text": "Make the girl from Image 1 wear the black dress from Image 2 and sit in the pose from Image 3."
 }
 ]
 }
 ]
 },
 "parameters": {
 "n": 2,
 "negative_prompt": " ",
 "prompt_extend": true,
 "watermark": false,
 "size": "1024*1536"
 }
}&#x27;

``` Response example

 Copy ```\n{
 "output": {
 "choices": [
 {
 "finish_reason": "stop",
 "message": {
 "role": "assistant",
 "content": [
 {
 "image": "https://dashscope-result-sz.oss-cn-shenzhen.aliyuncs.com/xxx.png?Expires=xxx"
 },
 {
 "image": "https://dashscope-result-sz.oss-cn-shenzhen.aliyuncs.com/xxx.png?Expires=xxx"
 }
 ]
 }
 }
 ]
 },
 "usage": {
 "width": 1536,
 "image_count": 2,
 "height": 1024
 },
 "request_id": "bf37ca26-0abe-98e4-8065-xxxxxx"
}

``` 
 Download images to your local device

 - Python 
- Java 

 Copy ```\nimport requests


def download_image(image_url, save_path=&#x27;output.png&#x27;):
 try:
 response = requests.get(image_url, stream=True, timeout=300) # Set timeout
 response.raise_for_status() # Raise an exception if the HTTP status code is not 200
 with open(save_path, &#x27;wb&#x27;) as f:
 for chunk in response.iter_content(chunk_size=8192):
 f.write(chunk)
 print(f"Image downloaded successfully to: {save_path}")

 except requests.exceptions.RequestException as e:
 print(f"Image download failed: {e}")


image_url = "https://dashscope-result-sz.oss-cn-shenzhen.aliyuncs.com/xxx.png?Expires=xxx"
download_image(image_url, save_path=&#x27;output.png&#x27;)

``` Copy ```\nimport java.io.FileOutputStream;
import java.io.InputStream;
import java.net.HttpURLConnection;
import java.net.URL;

public class ImageDownloader {
 public static void downloadImage(String imageUrl, String savePath) {
 try {
 URL url = new URL(imageUrl);
 HttpURLConnection connection = (HttpURLConnection) url.openConnection();
 connection.setConnectTimeout(5000);
 connection.setReadTimeout(300000);
 connection.setRequestMethod("GET");
 InputStream inputStream = connection.getInputStream();
 FileOutputStream outputStream = new FileOutputStream(savePath);
 byte[] buffer = new byte[8192];
 int bytesRead;
 while ((bytesRead = inputStream.read(buffer)) != -1) {
 outputStream.write(buffer, 0, bytesRead);
 }
 inputStream.close();
 outputStream.close();

 System.out.println("Image downloaded successfully to: " + savePath);
 } catch (Exception e) {
 System.err.println("Image download failed: " + e.getMessage());
 }
 }

 public static void main(String[] args) {
 String imageUrl = "https://dashscope-result-sz.oss-cn-shenzhen.aliyuncs.com/xxx.png?Expires=xxx";
 String savePath = "output.png";
 downloadImage(imageUrl, savePath);
 }
}

``` 
## [​ ](#input-instructions) Input instructions

### [​ ](#input-images-messages) Input images (messages)

The `messages` parameter is an array that must contain a single object. This object must include the `role` and `content` properties. The `role` property must be set to `user`. The `content` property must include both `image` (one to three images) and `text` (one editing instruction).
The input images must meet the following requirements:

- 
The supported image formats are JPG, JPEG, PNG, BMP, TIFF, WEBP, and GIF.
 The output image is in PNG format. For animated GIFs, only the first frame is processed. 


- 
For best results, the image resolution should be between 384 and 3072 pixels for both width and height. A low resolution may result in a blurry output, while a high resolution increases processing time.


- 
The size of a single image file cannot exceed 10 MB.


Copy ```\n"messages": [
 {
 "role": "user",
 "content": [
 { "image": "Public URL or Base64 data of Image 1" },
 { "image": "Public URL or Base64 data of Image 2" },
 { "image": "Public URL or Base64 data of Image 3" },
 { "text": "Your editing instruction, for example: &#x27;The girl in Image 1 wears the black dress from Image 2 and sits in the pose from Image 3&#x27;" }
 ]
 }
]

``` 
### [​ ](#image-input-order) Image input order

**Input image 1****Input image 2****Output image** Image 1 Image 2 Prompt: Move Image 1 onto Image 2 Prompt: Move Image 2 onto Image 1
### [​ ](#image-input-methods) Image input methods

**Public URL**

- You can provide a publicly accessible image URL that supports the HTTP or HTTPS protocol.

- Example value: `https://xxxx/img.png`.


**Base64 encoding**
Convert the image file to a Base64-encoded string and concatenate it in the following format: `data:<mime_type>;base64,<base64_data>`.

- `<mime_type>`: The media type of the image, which must correspond to the file format.

- `<base64_data>`: The Base64-encoded string of the file.

- Example value: `data:image/jpeg;base64,GDU7MtCZz...` (The example is truncated for demonstration purposes.)


### [​ ](#more-parameters) More parameters

Adjust the generation results using the following **optional** parameters:

- 
**n**: The number of images to generate. The default value is 1. The qwen-image-2.0, qwen-image-edit-max, and qwen-image-edit-plus series of models support generating one to six images. The `qwen-image-edit` model supports generating only one image.


- 
**negative_prompt**: Describes content to exclude from the image, such as "blur" or "extra fingers". This parameter helps optimize the quality of the generated image.


- 
**watermark**: Specifies whether to add a "Qwen-Image" watermark to the bottom-right corner of the image. The default value is `false`.


- 
**seed**: The random number seed. The value must be an integer from `[0, 2147483647]`. If this parameter is not specified, the algorithm generates a random number to use as the seed. Using the same seed value helps ensure consistent generation results.


The following **optional** parameters are available only for the qwen-image-2.0, qwen-image-edit-max, and qwen-image-edit-plus series of models:

- **size**: The resolution of the output image. The format is `width*height`, such as `"1024*2048"`. For the qwen-image-2.0 series models, you can set the width and height freely. The total pixels of the output image must be between 512 x 512 and 2048 x 2048. By default, the resolution is the same as the input image (the last image if multiple are provided). For the qwen-image-edit-max and qwen-image-edit-plus series models, the width and height can range from 512 to 2048 pixels. By default, the output image has a resolution close to `1024*1024` and an aspect ratio similar to the original image.

- **prompt_extend:** Enables or disables the prompt rewriting feature. The default value is `true`. If enabled, the model optimizes the prompt. This feature can significantly improve the results for simple or less descriptive prompts.


For a complete list of parameters, see [Qwen-Image-Edit API reference](/api-reference/image-generation/qwen-image-editing).
## [​ ](#overview) Overview

### [​ ](#multi-image-fusion) Multi-image fusion

**Input image 1****Input image 2****Input image 3****Output image** The girl in Image 1 wears the necklace from Image 2 and carries the bag from Image 3 on her left shoulder.
### [​ ](#subject-consistency) Subject consistency

**Input image****Output image 1****Output image 2** Change the image to an ID photo with a blue background. The person is wearing a white shirt, a black suit, and a striped tie. The person is wearing a white shirt, a gray suit, and a striped tie. One hand rests on the tie. The background is light-colored. The air conditioner is placed in a living room next to a sofa. Mist is added from the air conditioner&#x27;s vent, extending over the sofa. Green leaves are also added.
### [​ ](#sketch-creation) Sketch creation

**Input image****Output image** Generate an image that matches the detailed shape outlined in Image 1 and follows this description: A young woman smiles on a sunny day. She wears round brown sunglasses with a leopard print frame. Her hair is neatly tied up, she wears pearl earrings, a dark blue scarf with purple star patterns, and a black leather jacket. Generate an image that matches the detailed shape outlined in Image 1 and follows this description: An elderly man smiles at the camera. His face is wrinkled, his hair is messy in the wind, and he wears round-framed reading glasses. He has a worn-out red scarf with star patterns around his neck and is wearing a cotton-padded jacket.
### [​ ](#creative-product-generation) Creative product generation

**Input image****Output image** Make this bear sit under the moon (represented by a light gray crescent outline on a white background), holding a guitar, with small stars and speech bubbles with phrases such as "Be Kind" floating around. Print this design on a T-shirt and a paper tote bag. A female model is displaying these items. The woman is also wearing a baseball cap with "Be kind" written on it. A hyper-realistic 1/7 scale character model, designed as a commercial finished product, is placed on a desk with an iMac that has a white keyboard. The model stands on a clean, round, transparent acrylic base with no labels or text. Professional studio lighting highlights the sculpted details. On the iMac screen in the background, the ZBrush modeling process for the same model is displayed. Next to the model, place a packaging box with a transparent window on the front, showing only the clear plastic shell inside. The box is slightly taller than the model and reasonably sized to hold it. This bear is wearing an astronaut suit and pointing into the distance. This bear is wearing a gorgeous ball gown, with its arms spread in an elegant dance pose. This bear is wearing sportswear, holding a basketball, with one leg bent.
### [​ ](#generate-image-from-depth-map) Generate image from depth map

**Input image****Output image** Generate an image that matches the depth map outlined in Image 1 and follows this description: A blue bicycle is parked in a side alley, with a few weeds growing from cracks in the stone in the background. Generate an image that matches the depth map outlined in Image 1 and follows this description: A worn-out red bicycle is parked on a muddy path, with a dense primeval forest in the background.
### [​ ](#generate-image-from-keypoints) Generate image from keypoints

**Input image****Output image** Generate an image that matches the human pose outlined in Image 1 and follows this description: A Chinese woman in a Hanfu is holding an oil-paper umbrella in the rain, with a Suzhou garden in the background. Generate an image that matches the human pose outlined in Image 1 and follows this description: A young man stands on a subway platform. He wears a baseball cap, a T-shirt, and jeans. A train is speeding by behind him.
### [​ ](#text-editing) Text editing

**Input image****Output image****Input image****Output image** Replace &#x27;HEALTH INSURANCE&#x27; on the Scrabble tiles with &#x27;**Tomorrow will be better**&#x27;. Change the phrase "Take a Breather" on the note to "**Relax and Recharge**".
**Input image****Output image** Change "Qwen-Image" to a black ink-drip font. Change "Qwen-Image" to a black handwriting font. Change "Qwen-Image" to a black pixel font. Change "Qwen-Image" to red. Change "Qwen-Image" to a blue-purple gradient. Change "Qwen-Image" to candy colors. Change the material of "Qwen-Image" to metal. Change the material of "Qwen-Image" to clouds. Change the material of "Qwen-Image" to glass.
### [​ ](#add-delete-replace-and-modify) Add, delete, replace, and modify

**Input image****Output image****Input image****Output image** Add a small wooden sign in front of the penguin that says "Welcome to Penguin Beach". Remove the hair from the plate.
### [​ ](#viewpoint-transformation) Viewpoint transformation

**Input image****Output image****Input image****Output image** Get a front view. Face left. Get a rear view. Face right.
### [​ ](#old-photo-processing) Old photo processing

**Input image****Output image** Restore the old photo, remove scratches, reduce noise, enhance details, high resolution, realistic image, natural skin tone, clear facial features, no distortion. Intelligently colorize the image based on its content to make it more vivid.
## [​ ](#billing-and-rate-limits) Billing and rate limits

For free quota and pricing, see [pricing](/developer-guides/getting-started/pricing).
For rate limits, see [Rate limits](/developer-guides/administration/rate-limits).
**Billing details:**

- Billing is based on the **number of successfully generated images**. Failed model calls or processing errors do not incur fees or consume the free quota.

- You can enable the &#x27;Free quota only&#x27; feature to avoid extra charges after your free quota is used up. For more information, see [Free quota for new users](/resources/free-quota).


## [​ ](#api-reference) API reference

For the input and output parameters of the API, see [Qwen - image editing](/api-reference/image-generation/qwen-image-editing).
## [​ ](#error-codes) Error codes

If a call fails, see [Error messages](/api-reference/preparation/error-messages).
## [​ ](#faq) FAQ

**Q: What languages do the Qwen image editing models support?**
They officially support Simplified Chinese and **English**. Other languages may work, but results are not guaranteed.
**Q: How do I view model invocation metrics?**
One hour after a model invocation completes, go to the [**Analytics**](https://home.qwencloud.com/analytics) page to view metrics such as invocation count and success rate. For more information, see [Billing and cost management](/resources/bill-query).
**Q: How do I get the domain name whitelist for image storage?**
Images generated by models are stored in OSS. The API returns a temporary public URL. **To configure a firewall whitelist for this download URL**, note the following: The underlying storage may change dynamically. This topic does not provide a fixed OSS domain name whitelist to prevent access issues caused by outdated information. If you have security control requirements, contact your account manager to obtain the latest OSS domain name list.
For more information, see [Images & videos FAQ](/resources/faq-images-videos). [Previous ](/developer-guides/image-generation/text-to-image)[Image editing - Wan2.7/2.6/2.5 Edit images using text instructions with Wan2.7, 2.6, and 2.5 models. Next ](/developer-guides/image-generation/wan-image-editing)
