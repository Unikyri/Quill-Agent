# Z-Image

> **Source:** https://docs.qwencloud.com/api-reference/image-generation/z-image

POST/services/aigc/multimodal-generation/generation cURL cURL

 Copy ```\ncurl --location 'https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation' \
--header 'Content-Type: application/json' \
--header "Authorization: Bearer $DASHSCOPE_API_KEY" \
--data '{
 "model": "z-image-turbo",
 "input": {
 "messages": [
 {
 "role": "user",
 "content": [
 {
 "text": "A sitting orange cat with a happy expression, lively and cute, realistic and accurate"
 }
 ]
 }
 ]
 },
 "parameters": {
 "prompt_extend": false,
 "size": "1024*1024"
 }
}'
``` 200 400 401 429 Copy ```\n{
 "output": {
 "choices": [
 {
 "finish_reason": "stop",
 "message": {
 "role": "assistant",
 "content": [
 {
 "image": "https://dashscope-result-sgp.oss-ap-southeast-1.aliyuncs.com/xxx.png?Expires=xxx"
 },
 {
 "text": "Photo of a stylish young woman..."
 }
 ],
 "reasoning_content": ""
 }
 }
 ]
 },
 "usage": {
 "width": 1024,
 "height": 1024,
 "image_count": 1,
 "input_tokens": 0,
 "output_tokens": 0,
 "total_tokens": 0
 },
 "request_id": "abf1645b-b630-433a-92f6-xxxxxx"
}
``` [Get an API key](/api-reference/preparation/api-key) and [set it as an environment variable](/api-reference/preparation/export-api-key-env). If using the SDK, [install it](/api-reference/preparation/install-sdk). 
Z-Image is a lightweight text-to-image model that generates images quickly. It renders Chinese and English text and adapts to various resolutions and aspect ratios.
**Quick links**: [Try it online](https://home.qwencloud.com/try-ai/chat?models=z-image-turbo) | [Technical blog](https://tongyi-mai.github.io/Z-Image-blog/) ### Authorizations
 [​ ](#authorization) Authorizationstring header required DashScope API key. Create one in the [Qwen Cloud console](https://home.qwencloud.com/api-keys). Alternatively, pass it via the `X-DashScope-ApiKey` request header.

 ### Body
application/json [​ ](#model) modelenum<string> required Model name.

 Available options:z-image-turbo Example:z-image-turbo [​ ](#input) inputobject required Input content.

 Show child attributes

 [​ ](#inputmessages) input.messagesobject[] required Request content array. **Single-turn only** — pass one message with `role: user`. Multi-turn is not supported.

 Required range:items: 1–1 Show child attributes

 [​ ](#inputmessagesrole) input.messages.roleenum<string> required Message role. Must be `user`.

 Available options:user [​ ](#inputmessagescontent) input.messages.contentobject[] required Message content array. Must contain exactly one text object. Passing zero or multiple text objects returns an error.

 Show child attributes

 [​ ](#inputmessagescontenttext) input.messages.content.textstring Positive prompt describing desired content, style, and composition. Supports Chinese and English. Max 800 characters (each character, letter, number, or symbol counts as one). Extra characters are truncated.

 Example:A sitting orange cat with a happy expression, lively and cute, realistic and accurate Required range:length <= 800 [​ ](#parameters) parametersobject Image generation parameters.

 Show child attributes

 [​ ](#parameterssize) parameters.sizestring default"1024*1536" Output image resolution as `width*height`. Range: 512×512 to 2048×2048 pixels. Recommended range: 1024×1024 to 1536×1536. Default: `1024*1536`.


Recommended resolutions (1024×1024 total pixels): `1024*1024` (1:1), `832*1248` (2:3), `1248*832` (3:2), `864*1152` (3:4), `1152*864` (4:3), `720*1280` (9:16), `1280*720` (16:9).


Recommended resolutions (1280×1280 total pixels): `1280*1280` (1:1), `1024*1536` (2:3), `1536*1024` (3:2), `1104*1472` (3:4), `1472*1104` (4:3), `864*1536` (9:16), `1536*864` (16:9).


Recommended resolutions (1536×1536 total pixels): `1536*1536` (1:1), `1248*1872` (2:3), `1872*1248` (3:2), `1296*1728` (3:4), `1728*1296` (4:3), `1152*2048` (9:16), `2048*1152` (16:9).

 Example:1024*1024 [​ ](#parametersprompt-extend) parameters.prompt_extendboolean defaultfalse Enable intelligent prompt rewriting via LLM optimization.


- `false` (default): Returns the image and original prompt. No extra cost.

- `true`: Returns the image, an optimized prompt, and reasoning. Increases response time and cost — see model pricing for details.


 [​ ](#parametersseed) parameters.seedinteger Random seed for reproducibility. Valid range: `[0, 2147483647]`. Using the same seed yields similar outputs. If omitted, a random seed is used.


**Note:** Image generation is probabilistic — even with the same seed, results may vary.

 Required range:0 <= x <= 2147483647 ### Response
200-application/json [​ ](#output) outputobject Model output.

 Show child attributes

 [​ ](#outputchoices) output.choicesobject[] Model output content. Array contains one element.

 Show child attributes

 [​ ](#outputchoicesfinish-reason) output.choices.finish_reasonstring Reason for completion. `stop` indicates success.

 Example:stop [​ ](#outputchoicesmessage) output.choices.messageobject Model response message.

 Show child attributes

 [​ ](#outputchoicesmessagerole) output.choices.message.roleenum<string> Message role. Always `assistant`.

 Available options:assistant [​ ](#outputchoicesmessagecontent) output.choices.message.contentobject[] Response content items. Contains an `image` object with the generated image URL and a `text` object with the prompt.

 Show child attributes

 [​ ](#outputchoicesmessagecontentimage) output.choices.message.content.imagestring Generated image URL (PNG). **Valid for 24 hours** — download promptly.

 [​ ](#outputchoicesmessagecontenttext) output.choices.message.content.textstring The original prompt (when `prompt_extend=false`) or the rewritten prompt (when `prompt_extend=true`).

 [​ ](#outputchoicesmessagereasoning-content) output.choices.message.reasoning_contentstring Model reasoning output. Only returned when `prompt_extend=true`; empty string otherwise.

 [​ ](#usage) usageobject Usage statistics. Includes data for successful generations only.

 Show child attributes

 [​ ](#usagewidth) usage.widthinteger Generated image width in pixels.

 [​ ](#usageheight) usage.heightinteger Generated image height in pixels.

 [​ ](#usageimage-count) usage.image_countinteger Number of generated images. Always 1.

 [​ ](#usageinput-tokens) usage.input_tokensinteger Input tokens consumed. `0` when `prompt_extend=false`.

 [​ ](#usageoutput-tokens) usage.output_tokensinteger Output tokens consumed. `0` when `prompt_extend=false`.

 [​ ](#usageoutput-tokens-details) usage.output_tokens_detailsobject Output token breakdown. Only returned when `prompt_extend=true`.

 Show child attributes

 [​ ](#usageoutput-tokens-detailsreasoning-tokens) usage.output_tokens_details.reasoning_tokensinteger Tokens used for reasoning.

 [​ ](#usagetotal-tokens) usage.total_tokensinteger Total tokens consumed. `0` when `prompt_extend=false`.

 [​ ](#request-id) request_idstring Unique request identifier. Use for tracing and troubleshooting.

 Example:abf1645b-b630-433a-92f6-xxxxxx
