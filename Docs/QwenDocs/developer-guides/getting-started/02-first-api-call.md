# First API call

> **Source:** https://docs.qwencloud.com/developer-guides/getting-started/first-api-call

Get started in a few minutes

 Copy page Set up your account and make your first API call to Qwen models.
 New users get a free quota to try models at no cost. See [Free quota](/resources/free-quota) for details. 
## [​ ](#prerequisites) Prerequisites

 1 Create an account

Go to [Qwen Cloud](https://home.qwencloud.com/) and sign in with GitHub or email. 2 Get your API key

Navigate to [**API Keys**](https://home.qwencloud.com/api-keys), click **Create API key**, and copy your key (starts with `sk-`). [Detailed guide →](/api-reference/preparation/api-key) Keep your API key secret! Never commit it to version control or share it publicly. 3 Set your environment variable

Store your API key so your code can access it:macOS/Linux Windows PowerShell Copy ```\nexport DASHSCOPE_API_KEY="sk-your-api-key-here"

``` For permanent setup across sessions, see [Configure your API key →](/api-reference/preparation/export-api-key-env). 
## [​ ](#make-your-first-call) Make your first call

- Python 
- Node.js 
- curl 

 Install the OpenAI SDK:Copy ```\npip install openai

``` Create a file `hello_qwen.py`:Copy ```\nimport os
from openai import OpenAI

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)

completion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=[
 {"role": "user", "content": "Hello! Tell me a fun fact about AI."}
 ]
)

print(completion.choices[0].message.content)

``` Run it:Copy ```\npython hello_qwen.py

``` Install the OpenAI SDK:Copy ```\nnpm install openai

``` Create a file `hello_qwen.mjs`:Copy ```\nimport OpenAI from "openai";

const openai = new OpenAI({
 apiKey: process.env.DASHSCOPE_API_KEY,
 baseURL: "https://dashscope-intl.aliyuncs.com/compatible-mode/v1"
});

const completion = await openai.chat.completions.create({
 model: "qwen3.7-plus",
 messages: [
 { role: "user", content: "Hello! Tell me a fun fact about AI." }
 ]
});

console.log(completion.choices[0].message.content);

``` Run it:Copy ```\nnode hello_qwen.mjs

``` macOS/Linux Windows PowerShell Copy ```\ncurl -X POST https://dashscope-intl.aliyuncs.com/compatible-mode/v1/chat/completions \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY" \
 -H "Content-Type: application/json" \
 -d &#x27;{
 "model": "qwen3.7-plus",
 "messages": [
 {
 "role": "user",
 "content": "Hello! Tell me a fun fact about AI."
 }
 ]
 }&#x27;

``` 
 For Java, Go, PHP, C#, and other languages, use the OpenAI-compatible endpoint shown in the curl tab with your language&#x27;s HTTP client. For Java, you can also use the [DashScope Java SDK](/developer-guides/text-generation/quickstart). 
## [​ ](#whats-next) What&#x27;s next?


- [Text generation guide](/developer-guides/text-generation/quickstart) — Streaming, function calling, and more

- [Vision models](/developer-guides/multimodal/vision) — Analyze images and videos

- [Model selection](/developer-guides/getting-started/model-selection) — Choose the right model for your use case

- [Try AI](https://home.qwencloud.com/try-ai) — Try models interactively


 [Previous ](/developer-guides/getting-started/introduction)[Choose models Match your specific use case Next ](/developer-guides/getting-started/model-selection)
