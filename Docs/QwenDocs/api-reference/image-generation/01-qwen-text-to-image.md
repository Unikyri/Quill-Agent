# Qwen synchronous

> **Source:** https://docs.qwencloud.com/api-reference/image-generation/qwen-text-to-image

POST/services/aigc/multimodal-generation/generation cURL cURL

 Copy ```\ncurl --location 'https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/multimodal-generation/generation' \
--header 'Content-Type: application/json' \
--header "Authorization: Bearer $DASHSCOPE_API_KEY" \
--data '{
 "model": "qwen-image-2.0-pro",
 "input": {
 "messages": [
 {
 "role": "user",
 "content": [
 {
 "text": "Healing-style hand-drawn poster featuring three puppies playing with a ball on lush green grass. The main title Come Play Ball! is prominently displayed at the top in bold, blue cartoon font."
 }
 ]
 }
 ]
 },
 "parameters": {
 "negative_prompt": "",
 "prompt_extend": true,
 "watermark": false,
 "size": "2048*2048"
 }
}'
``` 200 400 Copy ```\n{
 "output": {
 "choices": [
 {
 "finish_reason": "stop",
 "message": {
 "content": [
 {
 "image": "https://dashscope-result-sh.oss-cn-shanghai.aliyuncs.com/xxx.png?Expires=xxx"
 }
 ],
 "role": "assistant"
 }
 }
 ]
 },
 "usage": {
 "height": 928,
 "image_count": 1,
 "width": 1664
 },
 "request_id": "d0250a3d-b07f-49e1-bdc8-6793f4929xxx"
}
``` [Get an API key](/api-reference/preparation/api-key) and [set it as an environment variable](/api-reference/preparation/export-api-key-env). If using the SDK, [install it](/api-reference/preparation/install-sdk). ### Authorizations
 [​ ](#authorization) Authorizationstring header required DashScope API Key. Create one in the [Qwen Cloud console](https://home.qwencloud.com/api-keys).

 ### Body
application/json [​ ](#model) modelenum<string> required Model name.

 Available options:qwen-image-2.0-pro,qwen-image-2.0-pro-2026-06-22,qwen-image-2.0-pro-2026-04-22,qwen-image-2.0-pro-2026-03-03,qwen-image-2.0,qwen-image-max,qwen-image-plus,qwen-image Example:qwen-image-2.0-pro [​ ](#input) inputobject required Input data containing the messages array.

 Show child attributes

 [​ ](#inputmessages) input.messagesobject[] required Single-turn only. Exactly one message with role `user`.

 Required range:items: 1–1 Show child attributes

 [​ ](#inputmessagesrole) input.messages.roleenum<string> required Must be `user`.

 Available options:user [​ ](#inputmessagescontent) input.messages.contentobject[] required Message content array. Must contain exactly one text object.

 Required range:items: 1–1 Show child attributes

 [​ ](#inputmessagescontenttext) input.messages.content.textstring required Positive prompt describing the desired content, style, and composition. The qwen-image-2.0 series accept up to 1,300 tokens. Other models accept up to 800 tokens. The system truncates excess tokens.

 Example:Healing-style hand-drawn poster featuring three puppies playing with a ball on lush green grass. [​ ](#parameters) parametersobject Image generation parameters.

 Show child attributes

 [​ ](#parametersnegative-prompt) parameters.negative_promptstring Describes content you do NOT want in the image. Max 500 characters. Excess is auto-truncated.

 Required range:length <= 500 [​ ](#parameterssize) parameters.sizestring default"2048*2048" Output resolution as `width*height`. **qwen-image-2.0 series**: total pixels between 512*512 and 2048*2048, default `2048*2048`. **qwen-image-max**: supports custom resolution (total pixels between 512*512 and 2048*2048) and fixed sizes. **qwen-image-plus/image**: fixed sizes only. Fixed sizes for qwen-image-max/plus/image: `1664*928` (16:9, default), `1472*1104` (4:3), `1328*1328` (1:1), `1104*1472` (3:4), `928*1664` (9:16).

 Example:2048*2048 [​ ](#parametersn) parameters.ninteger default1 Number of images to generate. Default: 1. qwen-image-2.0 series: 1-6. qwen-image-max/plus series: fixed at 1.

 Required range:1 <= x <= 6 [​ ](#parametersprompt-extend) parameters.prompt_extendboolean defaulttrue Enable prompt rewriting. `true` (default): model optimizes the prompt. `false`: use your prompt as-is.

 [​ ](#parameterswatermark) parameters.watermarkboolean defaultfalse Add a "Qwen-Image" watermark to the bottom-right corner. Default: `false`.

 [​ ](#parametersseed) parameters.seedinteger Random number seed. Range: [0, 2147483647]. Same seed produces more consistent (but not identical) results. If omitted, a random seed is used.

 Required range:0 <= x <= 2147483647 ### Response
200-application/json [​ ](#output) outputobject Show child attributes

 [​ ](#outputchoices) output.choicesobject[] List of generated results. Contains one element per generated image.

 Show child attributes

 [​ ](#outputchoicesfinish-reason) output.choices.finish_reasonstring `stop` means normal completion.

 Example:stop [​ ](#outputchoicesmessage) output.choices.messageobject Show child attributes

 [​ ](#outputchoicesmessagerole) output.choices.message.roleenum<string> Always `assistant`.

 Available options:assistant [​ ](#outputchoicesmessagecontent) output.choices.message.contentobject[] Response content array.

 Show child attributes

 [​ ](#outputchoicesmessagecontentimage) output.choices.message.content.imagestring URL of the generated image (PNG). **Valid for 24 hours.** Download promptly.

 [​ ](#outputtask-metric) output.task_metricobject Task result statistics. Not returned for qwen-image-2.0 series.

 Show child attributes

 [​ ](#outputtask-metrictotal) output.task_metric.TOTALinteger Total number of tasks.

 [​ ](#outputtask-metricsucceeded) output.task_metric.SUCCEEDEDinteger Number of successful tasks.

 [​ ](#outputtask-metricfailed) output.task_metric.FAILEDinteger Number of failed tasks.

 [​ ](#usage) usageobject Usage statistics (successful results only).

 Show child attributes

 [​ ](#usageimage-count) usage.image_countinteger Number of images generated.

 [​ ](#usagewidth) usage.widthinteger Width of the generated image in pixels.

 [​ ](#usageheight) usage.heightinteger Height of the generated image in pixels.

 [​ ](#request-id) request_idstring default"d0250a3d-b07f-49e1-bdc8-6793f4929xxx" Unique request identifier for tracing and troubleshooting.

 Example:abf1645b-b630-433a-92f6-xxxxxx
