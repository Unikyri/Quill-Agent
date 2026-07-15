# Generate text

> **Source:** https://docs.qwencloud.com/developer-guides/text-generation/quickstart

Make your first text generation call

 Copy page Text generation models take natural language input and generate text for tasks such as question answering, writing, summarization, translation, and structured output generation.
## [​ ](#request-structure) Request structure

Text generation requests are typically sent as a `messages` array. Each message includes a `role` and `content`.

- **System message**: Provides high-level instructions or sets the model&#x27;s behavior.

- **User message**: Contains the user&#x27;s input or task.

- **Assistant message**: Contains the model&#x27;s response.


A typical request includes a `user` message and can optionally include a `system` message for more stable or more controllable output.
 The `system` message is optional, but recommended when you want more consistent behavior. 
Copy ```\n[
 {"role": "system", "content": "You are a helpful assistant. Answer clearly and concisely."},
 {"role": "user", "content": "Summarize the benefits of solar energy in three bullet points."}
]

``` 
The model returns its reply as an `assistant` message.
Copy ```\n{
 "role": "assistant",
 "content": "- Reduces reliance on fossil fuels.\n- Lowers long-term electricity costs.\n- Produces electricity with minimal operating emissions."
}

``` 
## [​ ](#make-your-first-call) Make your first call

Before you begin, [get an API key](/api-reference/preparation/api-key), [set it as an environment variable](/api-reference/preparation/export-api-key-env), and, if needed, [install the OpenAI or DashScope SDK](/api-reference/preparation/install-sdk).
Choose the API style that matches your stack:

- Start with **OpenAI Compatible -- Responses API** for new integrations.

- Use **OpenAI Compatible -- Chat Completions API** if you are migrating existing OpenAI-compatible code.

- Use **[Anthropic Messages API](/api-reference/chat/anthropic)** if you are migrating from Anthropic. Supports features such as thinking and tool calling.

- Use **DashScope** if you prefer the native SDK.


- OpenAI Compatible -- Responses API 
- OpenAI Compatible -- Chat Completions API 
- DashScope 
- DashScope -- Multimodal API 

 For usage notes, code examples, and migration guidance, see [OpenAI compatible - Responses](/api-reference/chat/openai-responses).Python Node.js curl Copy ```\nimport os
from openai import OpenAI

try:
 client = OpenAI(
 # If you have not set an environment variable, replace the next line with your API key: api_key="sk-xxx",
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
 )

 response = client.responses.create(
 model="qwen3.7-plus",
 input="Summarize the benefits of solar energy in three bullet points."
 )

 print(response)
except Exception as e:
 print(f"Error message: {e}")

``` **Response**The response includes these main fields:
- 
`id`: The response ID.


- 
`output`: The output list. Includes `reasoning` and `message`.
 The `reasoning` field appears only when [deep thinking](/developer-guides/text-generation/thinking) is enabled (for example, it is enabled by default in the Qwen3.5 and Qwen3.6 series). 


- 
`usage`: Token usage statistics.


Example text output:Copy ```\n- Reduces reliance on fossil fuels.
- Lowers long-term electricity costs.
- Produces electricity with minimal operating emissions.

``` Full JSON response

 Copy ```\n{
 "created_at": 1772249518,
 "id": "7ad48c6b-3cc4-904f-9284-5f419c6c5xxx",
 "model": "qwen3.7-plus",
 "object": "response",
 "output": [
 {
 "id": "msg_94805179-2801-45da-ac1c-a87e8ea20xxx",
 "summary": [
 {
 "text": "The user wants a concise answer in exactly three bullet points. Focus on the most broadly useful benefits of solar energy: reduced reliance on fossil fuels, long-term cost savings, and lower operating emissions. Keep the wording simple and direct.\n",
 "type": "summary_text"
 }
 ],
 "type": "reasoning"
 },
 {
 "content": [
 {
 "annotations": [],
 "text": "- Reduces reliance on fossil fuels.\n- Lowers long-term electricity costs.\n- Produces electricity with minimal operating emissions.",
 "type": "output_text"
 }
 ],
 "id": "msg_35be06c6-ca4d-4f2b-9677-7897e488dxxx",
 "role": "assistant",
 "status": "completed",
 "type": "message"
 }
 ],
 "parallel_tool_calls": false,
 "status": "completed",
 "tool_choice": "auto",
 "tools": [],
 "usage": {
 "input_tokens": 54,
 "input_tokens_details": {
 "cached_tokens": 0
 },
 "output_tokens": 662,
 "output_tokens_details": {
 "reasoning_tokens": 447
 },
 "total_tokens": 716,
 "x_details": [
 {
 "input_tokens": 54,
 "output_tokens": 662,
 "output_tokens_details": {
 "reasoning_tokens": 447
 },
 "total_tokens": 716,
 "x_billing_type": "response_api"
 }
 ]
 }
}

``` Python Java Node.js Go C# PHP curl Copy ```\nimport os
from openai import OpenAI

try:
 client = OpenAI(
 # If you have not set an environment variable, replace the next line with your API key: api_key="sk-xxx",
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
 )

 completion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=[
 {"role": "system", "content": "You are a helpful assistant."},
 {"role": "user", "content": "Summarize the benefits of solar energy in three bullet points."},
 ],
 )
 print(completion.choices[0].message.content)
 # To view the full response, uncomment the following line
 # print(completion.model_dump_json())
except Exception as e:
 print(f"Error message: {e}")

``` **Response**Copy ```\n- Reduces reliance on fossil fuels.
- Lowers long-term electricity costs.
- Produces electricity with minimal operating emissions.

``` Full JSON response

 Copy ```\n{
 "choices": [
 {
 "message": {
 "role": "assistant",
 "content": "- Reduces reliance on fossil fuels.\n- Lowers long-term electricity costs.\n- Produces electricity with minimal operating emissions."
 },
 "finish_reason": "stop",
 "index": 0,
 "logprobs": null
 }
 ],
 "object": "chat.completion",
 "usage": {
 "prompt_tokens": 26,
 "completion_tokens": 66,
 "total_tokens": 92
 },
 "created": 1726127645,
 "system_fingerprint": null,
 "model": "qwen3.7-plus",
 "id": "chatcmpl-81951b98-28b8-9659-ab07-xxxxxx"
}

``` qwen3.7-max, qwen3.7-max-2026-05-20, and qwen3.6-max-preview only support the text interface (`Generation`). qwen3.7-max-2026-06-08, the Qwen3.6 and Qwen3.5 series require the multimodal interface (`MultiModalConversation`). The examples in this tab use `qwen-plus` via the text interface. For models that require the multimodal interface, see the **Multimodal API** tab. Python Java Node.js Go C# PHP curl Copy ```\nimport json
import os
from dashscope import Generation
import dashscope

dashscope.base_http_api_url = "https://dashscope-intl.aliyuncs.com/api/v1"

messages = [
 {"role": "system", "content": "You are a helpful assistant."},
 {"role": "user", "content": "Summarize the benefits of solar energy in three bullet points."},
]
response = Generation.call(
 # If you have not set an environment variable, replace the next line with: api_key = "sk-xxx",
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 model="qwen-plus",
 messages=messages,
 result_format="message",
)

if response.status_code == 200:
 print(response.output.choices[0].message.content)
 # To view the full response, uncomment the following line
 # print(json.dumps(response, default=lambda o: o.__dict__, indent=4))
else:
 print(f"HTTP status code: {response.status_code}")
 print(f"Error code: {response.code}")
 print(f"Error message: {response.message}")

``` **Response**Copy ```\n- Reduces reliance on fossil fuels.
- Lowers long-term electricity costs.
- Produces electricity with minimal operating emissions.

``` Full JSON response

 Copy ```\n{
 "output": {
 "choices": [
 {
 "finish_reason": "stop",
 "message": {
 "role": "assistant",
 "content": "- Reduces reliance on fossil fuels.\n- Lowers long-term electricity costs.\n- Produces electricity with minimal operating emissions."
 }
 }
 ]
 },
 "usage": {
 "total_tokens": 92,
 "output_tokens": 66,
 "input_tokens": 26
 },
 "request_id": "09dceb20-ae2e-999b-85f9-xxxxxx",
 "model": "qwen-plus"
}

``` qwen3.7-max-2026-06-08, the Qwen3.6 and Qwen3.5 series require the multimodal DashScope API (`MultiModalConversation`), not the text interface (`Generation`). Running the examples from the previous tab directly will result in a url error. User message content must be an array of objects.Python Java curl Copy ```\nimport os
import dashscope
from dashscope import MultiModalConversation

dashscope.base_http_api_url = "https://dashscope-intl.aliyuncs.com/api/v1"

messages = [
 {"role": "system", "content": "You are a helpful assistant."},
 {
 "role": "user",
 "content": [{"text": "Summarize the benefits of solar energy in three bullet points."}],
 },
]
response = MultiModalConversation.call(
 # If you have not set an environment variable, replace the next line with: api_key="sk-xxx",
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 model="qwen3.7-plus",
 messages=messages,
)

if response.status_code == 200:
 print(response.output.choices[0].message.content[0]["text"])
 # To view the full response, uncomment the following line
 # import json; print(json.dumps(response, default=lambda o: o.__dict__, indent=4))
else:
 print(f"HTTP status code: {response.status_code}")
 print(f"Error code: {response.code}")
 print(f"Error message: {response.message}")

``` **Response**Copy ```\n- Reduces reliance on fossil fuels.
- Lowers long-term electricity costs.
- Produces electricity with minimal operating emissions.

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
 "text": "- Reduces reliance on fossil fuels.\n- Lowers long-term electricity costs.\n- Produces electricity with minimal operating emissions."
 }
 ]
 }
 }
 ]
 },
 "usage": {
 "input_tokens": 25,
 "output_tokens": 613,
 "total_tokens": 638
 },
 "request_id": "1486945b-ebc7-93a1-af4d-651f8e18e76f"
}

``` 
## [​ ](#handle-requests-asynchronously) Handle requests asynchronously

Once a basic synchronous request is working, asynchronous calls can improve throughput for high-concurrency workloads.
- OpenAI Compatible -- Chat Completions API 
- DashScope 

 Python Java Copy ```\nimport os
import asyncio
from openai import AsyncOpenAI
import platform

# Create an asynchronous client instance
client = AsyncOpenAI(
 # If you have not set an environment variable, replace the line below with: api_key="sk-xxx",
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1"
)

# Define an asynchronous task list
async def task(question):
 print(f"Send question: {question}")
 response = await client.chat.completions.create(
 messages=[
 {"role": "user", "content": question}
 ],
 model="qwen3.7-plus",
 )
 print(f"Model response: {response.choices[0].message.content}")

# Main asynchronous function
async def main():
 questions = [
 "Summarize the benefits of solar energy in three bullet points.",
 "Write a subject line for a product launch email.",
 "Translate \"Welcome to our platform\" into Spanish."
 ]
 tasks = [task(q) for q in questions]
 await asyncio.gather(*tasks)

if __name__ == &#x27;__main__&#x27;:
 # Set event loop policy
 if platform.system() == &#x27;Windows&#x27;:
 asyncio.set_event_loop_policy(asyncio.WindowsSelectorEventLoopPolicy())
 # Run the main coroutine
 asyncio.run(main(), debug=False)


``` The DashScope SDK supports asynchronous text generation calls only in Python.Copy ```\n# DashScope Python SDK version must be at least 1.19.0
import asyncio
import platform
from dashscope.aigc.generation import AioGeneration
import os
import dashscope
dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

# Define an asynchronous task list
async def task(question):
 print(f"Send question: {question}")
 response = await AioGeneration.call(
 # If you have not set an environment variable, replace the line below with: api_key="sk-xxx",
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 model="qwen-plus",
 messages=[{"role": "system", "content": "You are a helpful assistant."},
 {"role": "user", "content": question}],
 result_format="message",
 )
 print(f"Model response: {response.output.choices[0].message.content}")

# Main asynchronous function
async def main():
 questions = [
 "Summarize the benefits of solar energy in three bullet points.",
 "Write a subject line for a product launch email.",
 "Translate \"Welcome to our platform\" into Spanish."
 ]
 tasks = [task(q) for q in questions]
 await asyncio.gather(*tasks)

if __name__ == &#x27;__main__&#x27;:
 # Set event loop policy
 if platform.system() == &#x27;Windows&#x27;:
 asyncio.set_event_loop_policy(asyncio.WindowsSelectorEventLoopPolicy())
 # Run the main coroutine
 asyncio.run(main(), debug=False)

``` 
**Response**
 Because the call is asynchronous, the order of the responses may differ from the example. 
Copy ```\nSend question: Summarize the benefits of solar energy in three bullet points.
Send question: Write a subject line for a product launch email.
Send question: Translate "Welcome to our platform" into Spanish.
Model response: - Reduces reliance on fossil fuels.
- Lowers long-term electricity costs.
- Produces electricity with minimal operating emissions.
Model response: Meet our newest product launch
Model response: Bienvenido a nuestra plataforma.

``` 
## [​ ](#going-live) Going live

### [​ ](#build-better-context) Build better context

Feeding raw data directly to a large language model increases costs and reduces quality because of context-length limits. Context engineering improves output quality and efficiency by dynamically loading precise knowledge. Core techniques:

- **Prompt engineering**: Design and refine text instructions (prompts) to guide the model toward the desired outputs. For more information, see [Text-to-text prompt guide](/developer-guides/accuracy-tuning/text-generation).

- **Retrieval-augmented generation (RAG)**: Use this technique when the model must answer questions using external knowledge bases, such as product documentation or technical manuals.

- **Tool calling**: Allows the model to fetch real-time data, such as weather or traffic, or perform actions, such as calling an API or sending an email.

- **Memory mechanisms**: Provide the model with short-term and long-term memory to understand conversation history.


### [​ ](#explore-more-text-generation-features) Explore more text generation features

For complex scenarios:

- [Multi-turn conversations](/developer-guides/text-generation/multi-turn): Use this feature for follow-up questions or information gathering that requires continuous dialogue.

- [Streaming output](/developer-guides/text-generation/streaming): Use this feature for chatbots or real-time code generation to improve the user experience and avoid timeouts caused by long responses.

- [Deep thinking](/developer-guides/text-generation/thinking): Use this feature for complex reasoning or policy analysis that requires high-quality, structured answers.

- [Structured output](/developer-guides/text-generation/structured-output): Use this feature when you need the model to reply in a stable JSON format for programmatic use or data parsing.

- [Partial mode](/developer-guides/text-generation/partial-mode): Use this feature for code completion or long-form writing where the model continues from existing text.


## [​ ](#reference) Reference

For a complete list of model invocation parameters, see [OpenAI Compatible API Reference](/api-reference/chat/openai-chat), [Anthropic Compatible API Reference](/api-reference/chat/anthropic), and [DashScope API Reference](/api-reference/chat/dashscope).
## [​ ](#faq) FAQ

**Why is the input token count higher than the token count of the text I sent?**
When processing a conversation, the system uses a Chat Template to wrap the raw input text, adding control markers such as role identifiers and message boundaries. These system-generated markers are also counted as tokens.
For example, when you send the message `{"role": "user", "content": "Hi"}` to `qwen3.7-max`, the text "Hi" corresponds to only 1 token after tokenization. However, during system processing, the actual full input text is formatted as follows: `<|im_start|>user\nHi<|im_end|>\n<|im_start|>assistant\n<think>`. After tokenization, this full text increases the total input token count to 11.
**Why can&#x27;t the Qwen API analyze web links?**
The Qwen API cannot directly access or parse web links. You can use [Function calling](/developer-guides/text-generation/function-calling), or combine them with web scraping tools such as Python&#x27;s Beautiful Soup to read webpage content.
**Why do Qwen web app and API responses differ?**
The Qwen web app includes additional engineering optimizations beyond the Qwen API, enabling features such as webpage parsing, web search, image drawing, and PPT creation. These capabilities are not part of the core large language model API. You can replicate them using [Function calling](/developer-guides/text-generation/function-calling) to enhance model performance.
**Can the model directly generate Word, Excel, PDF, or PPT files?**
No, they cannot. Qwen Cloud text generation models output only plain text. You can convert the text to your desired format using code or third-party libraries. [Previous ](/developer-guides/getting-started/text-generation-models)[Partial mode Continue from a prefix Next ](/developer-guides/text-generation/partial-mode)
