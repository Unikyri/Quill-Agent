# DashScope multimodal embedding

> **Source:** https://docs.qwencloud.com/api-reference/multimodal-embedding/dashscope-multimodal-embedding

POST/services/embeddings/multimodal-embedding/multimodal-embedding cURL

 cURL (Independent Vectors)

 Copy ```\ncurl --location --request POST \
 'https://dashscope-intl.aliyuncs.com/api/v1/services/embeddings/multimodal-embedding/multimodal-embedding' \
--header "Authorization: Bearer $DASHSCOPE_API_KEY" \
--header 'Content-Type: application/json' \
--data '{
 "model": "tongyi-embedding-vision-plus",
 "input": {
 "contents": [
 {"text": "Multimodal embedding model"},
 {"image": "https://example.com/image.jpg"},
 {"video": "https://example.com/video.mp4"}
 ]
 }
}'
``` 200 400 401 429 Copy ```\n{
 "output": {
 "embeddings": [
 {
 "index": 0,
 "embedding": [
 0
 ],
 "type": "text"
 }
 ]
 },
 "usage": {
 "input_tokens": 0,
 "input_tokens_details": {
 "image_tokens": 0,
 "text_tokens": 0
 },
 "output_tokens": 0,
 "total_tokens": 0,
 "image_tokens": 0
 },
 "request_id": "1fff9502-a6c5-9472-9ee1-73930fdd04c5"
}
``` Convert text, images, and video into numerical vectors in a unified semantic space for cross-modal retrieval, similarity search, and content classification.
 Before you begin: [get an API key](/api-reference/preparation/api-key), [set it as an environment variable](/api-reference/preparation/export-api-key-env), and [install the DashScope SDK](/api-reference/preparation/install-sdk) if you use the SDK. 
## [​ ](#endpoint) Endpoint


- HTTP: `POST https://dashscope-intl.aliyuncs.com/api/v1/services/embeddings/multimodal-embedding/multimodal-embedding`

- SDK `base_http_api_url`: `https://dashscope-intl.aliyuncs.com/api/v1`


## [​ ](#model-overview) Model overview

ModelModalitiesDimensionsImage size per imagetongyi-embedding-vision-plusText, Image, Video, Multi-images64, 128, 256, 512, 1024, 1152 (default)10 MBtongyi-embedding-vision-flashText, Image, Video, Multi-images64, 128, 256, 512, 768 (default)5 MB 
## [​ ](#notes) Notes


- **Image input**: Public URL or Base64 data URI (`data:image/{format};base64,{data}`).

- **Multi-images**: Key `multi_images`. Value is a list of image URLs, max 8 images.

- **Video input**: Must be a public URL. Use the `fps` parameter in `parameters` to control frame sampling rate (range [0, 1], default 1.0).


 ### Authorizations
 [​ ](#authorization) Authorizationstring header required DashScope API Key. Create one in the [Qwen Cloud console](https://home.qwencloud.com/api-keys). Alternatively, you can pass the API Key via the `X-DashScope-ApiKey` request header.

 ### Body
application/json [​ ](#model) modelenum<string> required Model name for multimodal embedding.

 Available options:tongyi-embedding-vision-plus,tongyi-embedding-vision-flash Example:tongyi-embedding-vision-plus [​ ](#input) inputobject required Input data containing the content items.

 Show child attributes

 [​ ](#inputcontents) input.contentsobject[] required Content items. Each item is an object with one or more modality keys (`text`, `image`, `video`, `multi_images`). For independent vectors, use one modality per object. For fused vectors, combine modalities in a single object.

 Show child attributes

 [​ ](#inputcontentstext) input.contents.textstring Text content to embed.

 [​ ](#inputcontentsimage) input.contents.imagestring Image URL (public HTTP/HTTPS) or Base64 data URI (`data:image/{format};base64,{data}`).

 [​ ](#inputcontentsvideo) input.contents.videostring Video URL (must be a public URL).

 [​ ](#inputcontentsmulti-images) input.contents.multi_imagesstring[] List of image URLs for multi-image embedding. Max 8 images. Only supported by `tongyi-embedding-vision-plus` and `tongyi-embedding-vision-flash`.

 Required range:items <= 8 [​ ](#parameters) parametersobject Parameters for multimodal embedding.

 Show child attributes

 [​ ](#parametersoutput-type) parameters.output_typeenum<string> default"dense" Output format. Only `dense` is supported.

 Available options:dense [​ ](#parametersdimension) parameters.dimensioninteger Output vector dimension. Supported values vary by model. See the model overview table for defaults and options.

 [​ ](#parametersfps) parameters.fpsnumber default1 Video frame sampling rate. Range [0, 1]. Default: 1.0.

 Required range:0 <= x <= 1 [​ ](#parametersinstruct) parameters.instructstring Custom task instruction. English recommended. Typically yields 1-5% improvement in retrieval tasks.

 ### Response
200-application/json [​ ](#output) outputobject Show child attributes

 [​ ](#outputembeddings) output.embeddingsobject[] List of embedding results.

 Show child attributes

 [​ ](#outputembeddingsindex) output.embeddings.indexinteger Position index in the input contents list.

 [​ ](#outputembeddingsembedding) output.embeddings.embeddingnumber[] Vector of floating-point numbers.

 [​ ](#outputembeddingstype) output.embeddings.typeenum<string> Content type of this embedding.

 Available options:text,image,video [​ ](#usage) usageobject Token usage statistics. Fields vary by model: `tongyi-embedding-vision-*` models return `input_tokens` (combined text and image token count), `input_tokens_details`, `output_tokens`, and `total_tokens`; other models may return different fields — see individual field descriptions.

 Show child attributes

 [​ ](#usageinput-tokens) usage.input_tokensinteger Number of input tokens consumed. For `tongyi-embedding-vision-*` models, this value includes both text and image/video tokens.

 [​ ](#usageinput-tokens-details) usage.input_tokens_detailsobject Detailed breakdown of input tokens. Only returned by `tongyi-embedding-vision-*` models.

 Show child attributes

 [​ ](#usageinput-tokens-detailsimage-tokens) usage.input_tokens_details.image_tokensinteger Tokens consumed by image/video content in the input.

 [​ ](#usageinput-tokens-detailstext-tokens) usage.input_tokens_details.text_tokensinteger Tokens consumed by text content in the input.

 [​ ](#usageoutput-tokens) usage.output_tokensinteger Number of output tokens. Only returned by `tongyi-embedding-vision-*` models.

 [​ ](#usagetotal-tokens) usage.total_tokensinteger Total token count (input_tokens + output_tokens).

 [​ ](#usageimage-tokens) usage.image_tokensinteger Number of image or video tokens in the input. For video input, the system extracts frames (up to a system-configured limit) and calculates tokens based on the extracted frames.

 [​ ](#request-id) request_idstring Unique request identifier.

 Example:1fff9502-a6c5-9472-9ee1-73930fdd04c5
