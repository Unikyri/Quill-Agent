# Multi-turn conversations

> **Source:** https://docs.qwencloud.com/developer-guides/text-generation/multi-turn

Manage chat context

 Copy page The Qwen API is stateless and does not save conversation history. To implement multi-turn conversations, you must pass the conversation history in each request. You can also use strategies, such as truncation, summarization, and retrieval, to efficiently manage the context and reduce token consumption.
 This topic describes how to implement multi-turn conversation using OpenAI compatible Chat Completion or DashScope API. The Responses API provides a more convenient alternative, see [OpenAI compatible - Responses](/api-reference/chat/openai-responses). 
## [​ ](#how-it-works) How it works

To implement a multi-turn conversation, you must maintain a `messages` array. In each round, append the user&#x27;s latest question and the model&#x27;s response to this array. Then, use the updated array as the input for the next request.
Example of how the `messages` array changes during a conversation:
 1 First round

Add the user&#x27;s question to the `messages` array.Copy ```\n// Use a text model
[
 {"role": "user", "content": "Recommend a sci-fi movie about space exploration."}
]

// Use a multimodal model, for example, Qwen-VL
// {"role": "user",
// "content": [{"type": "image_url","image_url": {"url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20251031/ownrof/f26d201b1e3f4e62ab4a1fc82dd5c9bb.png"}},
// {"type": "text", "text": "What products are shown in the image?"}]
// }

``` 2 Second round

Add the model&#x27;s response and the user&#x27;s latest question to the `messages` array.Copy ```\n// Use a text model
[
 {"role": "user", "content": "Recommend a sci-fi movie about space exploration."},
 {"role": "assistant", "content": "I recommend &#x27;XXX&#x27;. It is a classic sci-fi work."},
 {"role": "user", "content": "Who is the director of this movie?"}
]

// Use a multimodal model, for example, Qwen-VL
//[
// {"role": "user", "content": [
// {"type": "image_url","image_url": {"url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20251031/ownrof/f26d201b1e3f4e62ab4a1fc82dd5c9bb.png"}},
// {"type": "text", "text": "What products are shown in the image?"}]},
// {"role": "assistant", "content": "The image shows three items: a pair of light blue overalls, a blue and white striped short-sleeve shirt, and a pair of white sneakers."},
// {"role": "user", "content": "What style are they?"}
//]

``` 
## [​ ](#getting-started) Getting started

- Responses API 
- OpenAI compatible 
- DashScope 

 The Responses API simplifies multi-turn conversations. Pass `previous_response_id` to link context automatically—no manual message history needed. For advanced session management, see [Using Conversations](#using-conversations). Use the response `id` (UUID format, such as `f0dbb153-117f-9bbf-8176-5284b47f3xxx`) as `previous_response_id`. Do not use a message `id` from the `output` array (such as `msg_56c860c4-3ad8-4a96-8553-d2f94c259xxx`). The response `id` expires in 7 days. Python Node.js curl Copy ```\nimport os
from openai import OpenAI

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)

# First round
response1 = client.responses.create(
 model="qwen3.7-plus",
 input="My name is John, please remember it."
)
print(f"First response: {response1.output_text}")

# Second round - use previous_response_id to link context
response2 = client.responses.create(
 model="qwen3.7-plus",
 input="Do you remember my name?",
 previous_response_id=response1.id
)
print(f"Second response: {response2.output_text}")

``` **Response example (second round):**Copy ```\n{
 "id": "f0dbb153-117f-9bbf-8176-5284b47f3xxx",
 "model": "qwen3.7-plus",
 "status": "completed",
 "output": [
 {
 "type": "message",
 "role": "assistant",
 "content": [
 {
 "type": "output_text",
 "text": "Yes, John! I remember your name. How can I assist you today?"
 }
 ]
 }
 ],
 "usage": {
 "input_tokens": 78,
 "output_tokens": 16,
 "total_tokens": 94
 }
}

``` Python Node.js curl Copy ```\nimport os
from openai import OpenAI


def get_response(messages):
 client = OpenAI(
 # If you have not configured the environment variable, replace the following line with your API key: api_key="sk-xxx",
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
 )
 completion = client.chat.completions.create(model="qwen3.7-plus", messages=messages)
 return completion

# Initialize a messages array
messages = [
 {
 "role": "system",
 "content": """You are a salesperson at the Bailian phone store. You are responsible for recommending phones to users. The phones have two parameters: screen size (including 6.1-inch, 6.5-inch, and 6.7-inch) and resolution (including 2K and 4K).
 You can only ask the user for one parameter at a time. If the user does not provide complete information, you need to ask a follow-up question to get the missing parameter. When all parameters are collected, you must say: I have understood your purchase intention. Please wait.""",
 }
]
assistant_output = "Welcome to the Bailian phone store. What screen size are you looking for?"
print(f"Model output: {assistant_output}\n")
while "I have understood your purchase intention" not in assistant_output:
 user_input = input("Please enter: ")
 # Add the user&#x27;s question to the messages list
 messages.append({"role": "user", "content": user_input})
 assistant_output = get_response(messages).choices[0].message.content
 # Add the model&#x27;s response to the messages list
 messages.append({"role": "assistant", "content": assistant_output})
 print(f"Model output: {assistant_output}")
 print("\n")

``` Python Java curl Copy ```\nimport os
from dashscope import Generation
import dashscope
dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

def get_response(messages):
 response = Generation.call(
 # If you have not configured the environment variable, replace the following line with your API key: api_key="sk-xxx",
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 model="qwen-plus",
 messages=messages,
 result_format="message",
 )
 return response


messages = [
 {
 "role": "system",
 "content": """You are a salesperson at the Bailian phone store. You are responsible for recommending phones to users. The phones have two parameters: screen size (including 6.1-inch, 6.5-inch, and 6.7-inch) and resolution (including 2K and 4K).
 You can only ask the user for one parameter at a time. If the user does not provide complete information, you need to ask a follow-up question to get the missing parameter. When all parameters are collected, you must say: I have understood your purchase intention. Please wait.""",
 }
]

assistant_output = "Welcome to the Bailian phone store. What screen size are you looking for?"
print(f"Model output: {assistant_output}\n")
while "I have understood your purchase intention" not in assistant_output:
 user_input = input("Please enter: ")
 # Add the user&#x27;s question to the messages list
 messages.append({"role": "user", "content": user_input})
 assistant_output = get_response(messages).output.choices[0].message.content
 # Add the model&#x27;s response to the messages list
 messages.append({"role": "assistant", "content": assistant_output})
 print(f"Model output: {assistant_output}")
 print("\n")

``` 
## [​ ](#for-multimodal-models) For multimodal models

 
- This section applies to multimodal models such as Qwen3-VL and Qwen3.5. For `Qwen-Omni`, see [Non-Realtime](/developer-guides/speech/multimodal-speech).

- Qwen3-Omni-Captioner is designed for single-turn tasks and does not support multi-turn conversations.


 
Multi-turn conversations for multimodal models differ from text models:

- **Construction of user messages**: User messages for multimodal models can contain multimodal information, such as images and audio, in addition to text.

- **DashScope SDK interface:** When you use the DashScope Python SDK, call the `MultiModalConversation` interface. When you use the DashScope Java SDK, call the `MultiModalConversation` class.


- OpenAI compatible 
- DashScope 

 Python Node.js curl Copy ```\nfrom openai import OpenAI
import os

client = OpenAI(
 # If you have not configured the environment variable, replace the following line with your API key: api_key="sk-xxx"
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1"
)
messages = [
 {"role": "user",
 "content": [
 {
 "type": "image_url",
 "image_url": {
 "url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20251031/ownrof/f26d201b1e3f4e62ab4a1fc82dd5c9bb.png"
 },
 },
 {"type": "text", "text": "What products are shown in the image?"},
 ],
 }
]
completion = client.chat.completions.create(
 model="qwen3-vl-plus", # You can replace this with other multimodal models and modify the messages as needed
 messages=messages,
 )
print(f"First round output: {completion.choices[0].message.content}")

assistant_message = completion.choices[0].message
messages.append(assistant_message.model_dump())
messages.append({
 "role": "user",
 "content": [
 {
 "type": "text",
 "text": "What style are they?"
 }
 ]
 })
completion = client.chat.completions.create(
 model="qwen3-vl-plus",
 messages=messages,
 )

print(f"Second round output: {completion.choices[0].message.content}")

``` Python Java curl Copy ```\nimport os
from dashscope import MultiModalConversation
import dashscope
dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

messages = [
 {
 "role": "user",
 "content": [
 {
 "image": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20251031/ownrof/f26d201b1e3f4e62ab4a1fc82dd5c9bb.png"
 },
 {"text": "What products are shown in the image?"},
 ],
 }
]
response = MultiModalConversation.call(
 # If you have not configured the environment variable, replace the following line with: api_key="sk-xxx",
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 model=&#x27;qwen3-vl-plus&#x27;, # You can replace this with other multimodal models and modify the messages as needed
 messages=messages
 )

print(f"Model first round output {response.output.choices[0].message.content[0][&#x27;text&#x27;]}")
messages.append(response[&#x27;output&#x27;][&#x27;choices&#x27;][0][&#x27;message&#x27;])
user_msg = {"role": "user", "content": [{"text": "What style are they?"}]}
messages.append(user_msg)
response = MultiModalConversation.call(
 # If the environment variable is not configured, please replace the following line with: api_key="sk-xxx",
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 model=&#x27;qwen3-vl-plus&#x27;,
 messages=messages
 )

print(f"Model second round output {response.output.choices[0].message.content[0][&#x27;text&#x27;]}")

``` 
## [​ ](#for-thinking-models) For thinking models

Thinking models return two fields: `reasoning_content` (the thinking process) and `content` (the response). When you update the messages array, retain only the `content` field and ignore the `reasoning_content` field.
Copy ```\n[
 {"role": "user", "content": "Recommend a sci-fi movie about space exploration."},
 {"role": "assistant", "content": "I recommend &#x27;XXX&#x27;. It is a classic sci-fi work."}, # Do not add the reasoning_content field when you add to the context
 {"role": "user", "content": "Who is the director of this movie?"}
]

``` 
 For more information about thinking models, see [Thinking](/developer-guides/text-generation/thinking) and [Vision](/developer-guides/multimodal/vision). For multi-turn conversations with Qwen3-Omni-Flash (thinking mode), see [Non-Realtime](/developer-guides/speech/multimodal-speech). 
- OpenAI compatible 
- DashScope 

 Python Node.js HTTP Copy ```\nfrom openai import OpenAI
import os

# Initialize the OpenAI client
client = OpenAI(
 # If you have not configured the environment variable, replace the following line with your API key: api_key="sk-xxx"
 api_key = os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1"
)


messages = []
conversation_idx = 1
while True:
 reasoning_content = "" # Define the complete thinking process
 answer_content = "" # Define the complete response
 is_answering = False # Determine whether to end the thinking process and start responding
 print("="*20+f"Conversation Round {conversation_idx}"+"="*20)
 conversation_idx += 1
 user_msg = {"role": "user", "content": input("Enter your message: ")}
 messages.append(user_msg)
 # Create a chat completion request
 completion = client.chat.completions.create(
 # You can replace this with other deep thinking models as needed
 model="qwen3.7-plus",
 messages=messages,
 extra_body={"enable_thinking": True},
 stream=True,
 # stream_options={
 # "include_usage": True
 # }
 )
 print("\n" + "=" * 20 + "Thinking Process" + "=" * 20 + "\n")
 for chunk in completion:
 # If chunk.choices is empty, print usage
 if not chunk.choices:
 print("\nUsage:")
 print(chunk.usage)
 else:
 delta = chunk.choices[0].delta
 # Print the thinking process
 if hasattr(delta, &#x27;reasoning_content&#x27;) and delta.reasoning_content != None:
 print(delta.reasoning_content, end=&#x27;&#x27;, flush=True)
 reasoning_content += delta.reasoning_content
 else:
 # Start responding
 if delta.content != "" and is_answering is False:
 print("\n" + "=" * 20 + "Complete Response" + "=" * 20 + "\n")
 is_answering = True
 # Print the response process
 print(delta.content, end=&#x27;&#x27;, flush=True)
 answer_content += delta.content
 # Add the content of the model&#x27;s response to the context
 messages.append({"role": "assistant", "content": answer_content})
 print("\n")

``` Python Java HTTP Copy ```\nimport os
from dashscope import MultiModalConversation
import dashscope
dashscope.base_http_api_url = "https://dashscope-intl.aliyuncs.com/api/v1/"

messages = []
conversation_idx = 1
while True:
 print("=" * 20 + f"Conversation Round {conversation_idx}" + "=" * 20)
 conversation_idx += 1
 user_msg = {"role": "user", "content": [{"text": input("Enter your message: ")}]}
 messages.append(user_msg)
 response = MultiModalConversation.call(
 # If you have not configured the environment variable, replace the following line with your API key: api_key="sk-xxx",
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 # qwen3.7-plus requires the multimodal endpoint. For qwen3-max, qwen-plus, etc., use Generation.call instead.
 model="qwen3.7-plus",
 messages=messages,
 enable_thinking=True,
 result_format="message",
 stream=True,
 incremental_output=True
 )
 # Define the complete thinking process
 reasoning_content = ""
 # Define the complete response
 answer_content = ""
 # Determine whether to end the thinking process and start responding
 is_answering = False
 print("=" * 20 + "Thinking Process" + "=" * 20)
 for chunk in response:
 # If both the thinking process and the response are empty, ignore
 if (chunk.output.choices[0].message.content == "" and
 chunk.output.choices[0].message.reasoning_content == ""):
 pass
 else:
 # If it is currently the thinking process
 if (chunk.output.choices[0].message.reasoning_content != "" and
 chunk.output.choices[0].message.content == ""):
 print(chunk.output.choices[0].message.reasoning_content, end="",flush=True)
 reasoning_content += chunk.output.choices[0].message.reasoning_content
 # If it is currently the response
 elif chunk.output.choices[0].message.content != "":
 if not is_answering:
 print("\n" + "=" * 20 + "Complete Response" + "=" * 20)
 is_answering = True
 print(chunk.output.choices[0].message.content, end="",flush=True)
 answer_content += chunk.output.choices[0].message.content
 # Add the content of the model&#x27;s response to the context
 messages.append({"role": "assistant", "content": answer_content})
 print("\n")
 # To print the complete thinking process and complete response, uncomment and run the following code
 # print("=" * 20 + "Complete Thinking Process" + "=" * 20 + "\n")
 # print(f"{reasoning_content}")
 # print("=" * 20 + "Complete Response" + "=" * 20 + "\n")
 # print(f"{answer_content}")

``` 
## [​ ](#using-conversations) Using conversations

`previous_response_id` works well for simple chained conversations. For server-side session management, cross-device continuity, or manual message control, use the Conversations API with the `conversation` parameter.
### [​ ](#create-a-conversation-and-chat) Create a conversation and chat

First, create a conversation using the Conversations API, then pass the `conversation` parameter and `instructions` (system prompt) to `responses.create`. The server manages context automatically.
Python Node.js curl Copy ```\nimport os
from openai import OpenAI

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)

# Create a conversation
conversation = client.conversations.create()

# First round
response1 = client.responses.create(
 conversation=conversation.id,
 model="qwen3.7-plus",
 instructions="You are a travel advisor who specializes in recommending travel destinations.",
 input="Recommend a city suitable for summer travel.",
)
print(f"First response: {response1.output_text}")

# Second round - context managed automatically by server
response2 = client.responses.create(
 conversation=conversation.id,
 model="qwen3.7-plus",
 instructions="You are a travel advisor who specializes in recommending travel destinations.",
 input="What are the must-try foods there?",
)
print(f"Second response: {response2.output_text}")

``` 
### [​ ](#add-messages-to-a-conversation) Add messages to a conversation

You can manually add message items to a conversation (such as supplementary user messages or external knowledge).
Python Node.js curl Copy ```\nitems = client.conversations.items.create(
 "conv_xxx", # Replace with your conversation id
 items=[
 {
 "type": "message",
 "role": "user",
 "content": [{"type": "input_text", "text": "Additional info: I prefer coastal cities."}],
 }
 ],
)
print(items.data)

``` 
### [​ ](#view-conversation-history) View conversation history

List all message items in a conversation to view the complete dialogue history.
Python Node.js curl Copy ```\nitems = client.conversations.items.list("conv_xxx") # Replace with your conversation id
print(items.data)

``` 
### [​ ](#important-notes) Important notes


- **ID validity**: Response `id` and message items in a conversation are valid for 7 days. The `conversation` itself has no expiration, but expired items no longer participate in context.

- **Correct ID source**: Use the response top-level `id`, not the `id` of messages inside the `output` array.

- **Cross-turn context**: Each time you pass `previous_response_id`, the system automatically links the full context from the initial conversation to the current turn.

- **Mutual exclusivity**: `previous_response_id` and `conversation` cannot be used together. Otherwise, you&#x27;ll get error `[400] INVALID_REQUEST: Mutually exclusive parameters: Ensure you are only providing one of: previous_response_id or conversation.`

- **Conversation message expiry**: The `conversation` itself has no expiry and can be used continuously. However, message items within it expire after 7 days and will no longer appear in the conversation context. We recommend passing system instructions via the `instructions` parameter rather than through items when creating a conversation, to avoid losing system instructions due to expiry.


### [​ ](#which-approach-to-choose) Which approach to choose?

ApproachBest for`previous_response_id`Simple chained multi-turn conversations without creating a separate session`conversation`Server-side session management, cross-device continuity, or manual message add/delete 
For more Conversations API operations (update conversation, delete conversation, delete messages, etc.), see [Conversations](/api-reference/platform-api/conversations).
## [​ ](#going-live) Going live

Multi-turn conversations can consume many tokens and may exceed the model&#x27;s context limit. Use these strategies to manage context and control costs.
### [​ ](#1-context-management) 1. Context management

The `messages` array grows with each round and may exceed the token limit.
#### [​ ](#1-1-context-truncation) 1.1. Context truncation

When the conversation history becomes too long, keep only the most recent N rounds of conversation. This method is simple to implement but results in the loss of earlier conversation information.
#### [​ ](#1-2-rolling-summary) 1.2. Rolling summary

To dynamically compress the conversation history and control the context length without losing core information, summarize the context as the conversation progresses:
a. When the conversation history reaches a certain length, such as 70% of the maximum context length, extract an earlier part of the history, such as the first half. Then, make a separate API call to the model to generate a "memory summary" of this part.
b. When you construct the next request, replace the lengthy conversation history with the "memory summary" and append the most recent conversation rounds.
#### [​ ](#1-3-vectorized-retrieval) 1.3. Vectorized retrieval

A rolling summary can cause some information loss. To allow the model to recall relevant information from a large volume of conversation history, switch from linear context passing to on-demand retrieval:
a. After each conversation round, store the conversation in a vector database.
b. When a user asks a question, retrieve relevant conversation records based on similarity.
c. Combine the retrieved conversation records with the most recent user input and send the combined content to the model.
### [​ ](#2-cost-control) 2. Cost control

Input tokens increase with each round, raising costs.
#### [​ ](#2-1-reduce-input-tokens) 2.1. Reduce input tokens

Use the context management strategies described previously to reduce input tokens and lower costs.
#### [​ ](#2-2-use-models-that-support-context-cache) 2.2. Use models that support context cache

The `messages` array is repeatedly processed and billed. [Context cache](/developer-guides/text-generation/context-cache) (available for select Qwen models including Qwen-Max, Qwen-Plus, Qwen-Flash, and Qwen-Coder) reduces costs and improves response speed.
 The context cache feature is enabled automatically. No code changes are required. 
## [​ ](#error-codes) Error codes

If a call fails, see [Error messages](/api-reference/preparation/error-messages). [Previous ](/developer-guides/text-generation/batch)[Context cache Cut cost with prefix reuse Next ](/developer-guides/text-generation/context-cache)
