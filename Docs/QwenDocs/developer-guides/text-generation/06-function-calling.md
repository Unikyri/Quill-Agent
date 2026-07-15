# Function Calling

> **Source:** https://docs.qwencloud.com/developer-guides/text-generation/function-calling

Connect models to external tools

 Copy page Large Language Models (LLMs) cannot access real-time data or external systems. Function Calling enables models to call external tools -- such as APIs, databases, and user-defined functions -- so a model can retrieve information or perform actions beyond its built-in capabilities. You define tools, the model decides when to call them, and your application executes the calls.
## [​ ](#how-it-works) How it works

Function Calling works through a multi-step interaction between your application and the LLM.

- **Make the first model call**


The application sends the user&#x27;s question and a list of available tools to the LLM.

- **Receive tool calling instructions from the model (tool name and input parameters)**


If the model decides to call an external tool, it returns a JSON instruction that specifies the function name and input parameters to pass.

If the model decides not to call a tool, it returns a natural-language response.


- **Run the tool in your application**


Your application runs the specified tool and obtains its output.

- **Make the second model call**


Add the tool&#x27;s output to the model&#x27;s context (the messages array), then call the model again.

- **Receive the final model response**


The model combines the tool&#x27;s output with the user&#x27;s question to generate a natural-language response.
The following diagram illustrates the workflow:
 
## [​ ](#supported-models) Supported models

All general-purpose text generation models support function calling, including third-party models (DeepSeek, Kimi, GLM, MiniMax). For vision, Qwen3-VL, qwen3.5-omni-plus, qwen3.5-omni-flash, and qwen3-omni-flash also support it. See [Models](/developer-guides/getting-started/text-generation-models) for the full list.
For realtime speech models, qwen3.5-omni-plus-realtime and qwen3.5-omni-flash-realtime support function calling via WebSocket. See [Realtime](#realtime) below for implementation details.
## [​ ](#getting-started) Getting started

This section shows how to use function calling with a weather lookup scenario.
- OpenAI compatible 
- DashScope 

 Python Node.js Copy ```\nfrom openai import OpenAI
from datetime import datetime
import json
import os
import random

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)
# Simulate user question
USER_QUESTION = "What&#x27;s the weather in Singapore?"
# Define tool list
tools = [
 {
 "type": "function",
 "function": {
 "name": "get_current_weather",
 "description": "Useful when you want to check the weather for a specific city.",
 "parameters": {
 "type": "object",
 "properties": {
 "location": {
 "type": "string",
 "description": "City or county, such as Singapore or New York.",
 }
 },
 "required": ["location"],
 },
 },
 },
]


# Simulate weather lookup tool
def get_current_weather(arguments):
 weather_conditions = ["sunny", "cloudy", "rainy"]
 random_weather = random.choice(weather_conditions)
 location = arguments["location"]
 return f"{location} is {random_weather} today."


# Wrap model response function
def get_response(messages):
 completion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=messages,
 tools=tools,
 )
 return completion


messages = [{"role": "user", "content": USER_QUESTION}]
response = get_response(messages)
assistant_output = response.choices[0].message
if assistant_output.content is None:
 assistant_output.content = ""
messages.append(assistant_output)
# If no tool is needed, output content directly
if assistant_output.tool_calls is None:
 print(f"No weather tool call needed. Direct reply: {assistant_output.content}")
else:
 # Enter tool calling loop
 while assistant_output.tool_calls is not None:
 tool_call = assistant_output.tool_calls[0]
 tool_call_id = tool_call.id
 func_name = tool_call.function.name
 arguments = json.loads(tool_call.function.arguments)
 print(f"Calling tool [{func_name}], parameters: {arguments}")
 # Execute tool
 tool_result = get_current_weather(arguments)
 # Build tool response message
 tool_message = {
 "role": "tool",
 "tool_call_id": tool_call_id,
 "content": tool_result,
 }
 print(f"Tool returned: {tool_message[&#x27;content&#x27;]}")
 messages.append(tool_message)
 # Call model again to get summarized natural-language reply
 response = get_response(messages)
 assistant_output = response.choices[0].message
 if assistant_output.content is None:
 assistant_output.content = ""
 messages.append(assistant_output)
 print(f"Final assistant reply: {assistant_output.content}")

``` `qwen3.7-plus` uses the `MultiModalConversation` interface. For models like `qwen-plus` and `qwen3-max`, use the `Generation` interface -- see [Text generation quickstart](/developer-guides/text-generation/quickstart) for examples. Python Java Copy ```\nimport os
from dashscope import MultiModalConversation
import dashscope
import json
import random
dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

# 1. Define tool list
tools = [
 {
 "type": "function",
 "function": {
 "name": "get_current_weather",
 "description": "Useful when you want to check the weather for a specific city.",
 "parameters": {
 "type": "object",
 "properties": {
 "location": {
 "type": "string",
 "description": "City or county, such as Singapore or New York.",
 }
 },
 "required": ["location"],
 },
 },
 }
]

# 2. Simulate weather lookup tool
def get_current_weather(arguments):
 weather_conditions = ["sunny", "cloudy", "rainy"]
 random_weather = random.choice(weather_conditions)
 location = arguments["location"]
 return f"{location} is {random_weather} today."

# 3. Wrap model response function
def get_response(messages):
 response = MultiModalConversation.call(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 model="qwen3.7-plus",
 messages=messages,
 tools=tools,
 )
 return response

# 4. Initialize conversation history
messages = [
 {
 "role": "user",
 "content": "What&#x27;s the weather in Singapore?"
 }
]

# 5. First model call
response = get_response(messages)
assistant_output = response.output.choices[0].message
messages.append(assistant_output)

# 6. Determine whether to call a tool
if "tool_calls" not in assistant_output or not assistant_output["tool_calls"]:
 print(f"No tool call needed. Direct reply: {assistant_output[&#x27;content&#x27;][0][&#x27;text&#x27;]}")
else:
 # 7. Enter tool calling loop
 while "tool_calls" in assistant_output and assistant_output["tool_calls"]:
 tool_call = assistant_output["tool_calls"][0]
 func_name = tool_call["function"]["name"]
 arguments = json.loads(tool_call["function"]["arguments"])
 tool_call_id = tool_call.get("id")
 print(f"Calling tool [{func_name}], parameters: {arguments}")
 tool_result = get_current_weather(arguments)
 tool_message = {
 "role": "tool",
 "content": tool_result,
 "tool_call_id": tool_call_id
 }
 print(f"Tool returned: {tool_message[&#x27;content&#x27;]}")
 messages.append(tool_message)
 response = get_response(messages)
 assistant_output = response.output.choices[0].message
 messages.append(assistant_output)
 # 8. Output final natural-language reply
 print(f"Final assistant reply: {assistant_output[&#x27;content&#x27;][0][&#x27;text&#x27;]}")

``` Running the code produces the following output:Output Copy ```\nCalling tool [get_current_weather], parameters: {&#x27;location&#x27;: &#x27;Singapore&#x27;}
Tool returned: Singapore is cloudy today.
Final assistant reply: Today&#x27;s weather in Singapore is cloudy.

``` 
## [​ ](#tool-schema-reference) Tool schema reference

Each tool is a JSON object with the following structure:
Copy ```\n{
 "type": "function",
 "function": {
 "name": "get_current_weather",
 "description": "Useful when you want to check the weather for a specific city.",
 "parameters": {
 "type": "object",
 "properties": {
 "location": {
 "type": "string",
 "description": "City or county, such as Singapore or New York."
 }
 },
 "required": ["location"]
 }
 }
}

``` 
FieldDescription`type`Always `"function"`.`function.name`The function name. Must match the function your application executes.`function.description`Describes the tool&#x27;s purpose. The model uses this to decide when to call the tool.`function.parameters`A JSON Schema object describing the function&#x27;s input parameters. Omit if the function takes no input.`function.parameters.properties`Each key is a parameter name. Values describe the parameter type and purpose.`function.parameters.required`Array of required parameter names. 
 Write clear, specific descriptions. The model relies on `description` fields to select the right tool and extract the right parameters. 
## [​ ](#specify-tool-calling-behavior) Specify tool calling behavior

### [​ ](#parallel-tool-calling) Parallel tool calling

By default, the model returns one tool call per response. If the user&#x27;s request requires multiple independent tool calls -- such as "What&#x27;s the weather in Beijing and Shanghai?" -- set `parallel_tool_calls` to `true`:
Python Node.js Copy ```\ncompletion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=messages,
 tools=tools,
 parallel_tool_calls=True
)

``` 
The `tool_calls` array then contains multiple entries:
Copy ```\n{
 "role": "assistant",
 "tool_calls": [
 {
 "function": { "name": "get_current_weather", "arguments": "{\"location\": \"Beijing\"}" },
 "index": 0, "id": "call_c2d8a3a2...", "type": "function"
 },
 {
 "function": { "name": "get_current_weather", "arguments": "{\"location\": \"Shanghai\"}" },
 "index": 1, "id": "call_dc7f2f67...", "type": "function"
 }
 ]
}

``` 
 Use parallel tool calling only when tasks are independent. If tasks depend on each other (such as when tool A&#x27;s input relies on tool B&#x27;s output), use the while loop from Getting started to call tools serially. 
### [​ ](#forced-tool-calling-tool-choice) Forced tool calling (tool_choice)

The `tool_choice` parameter controls whether the model calls a tool. The default is `"auto"` (model decides).

- 
**Force a specific tool**: Set `tool_choice` to `{"type": "function", "function": {"name": "get_current_weather"}}`. The model skips tool selection and always calls the specified function.


- 
**Force at least one tool**: Set `tool_choice` to `"required"`. The model always returns a tool call, so `tool_calls` is never empty. Use this when every question in the scenario requires a tool. Confirm the questions are relevant to the available tools first, or the results may be unexpected.


- 
**Block all tools**: Set `tool_choice` to `"none"` (or omit the `tools` parameter). The model replies directly without calling any tool.


Python Node.js Copy ```\n# Force a specific tool
completion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=messages,
 tools=tools,
 tool_choice={"type": "function", "function": {"name": "get_current_weather"}},
 extra_body={"enable_thinking": False}
)

# Force at least one tool
completion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=messages,
 tools=tools,
 tool_choice="required",
 extra_body={"enable_thinking": False}
)

# Block all tools
completion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=messages,
 tools=tools,
 tool_choice="none"
)

``` 
 Models like `qwen3.7-plus` enable thinking mode by default. Thinking mode only supports `tool_choice` set to `"auto"` or `"none"`. To force a specific tool, disable thinking mode first by setting `enable_thinking` to `false`. 

Remove the `tool_choice` parameter when summarizing tool outputs. Otherwise, the API still returns tool call information instead of a natural-language response.

## [​ ](#multi-turn-conversations) Multi-turn conversations

A user might ask "What&#x27;s the weather in Beijing?" in round one and "What about Shanghai?" in round two. Without round-one context, the model cannot determine which tool to call. In multi-turn conversations, keep the messages array after each round. Then add the new user message and invoke function calling again. The messages structure looks like this:
Copy ```\n[
 "System message -- Strategy guiding the model to call tools",
 "User message -- User&#x27;s question",
 "Assistant message -- Tool call information returned by the model",
 "Tool message -- Tool output",
 "Assistant message -- Model&#x27;s summary of tool call information",
 "User message -- User&#x27;s second question"
]

``` 
## [​ ](#streaming) Streaming

 For general streaming concepts (SSE protocol, how to enable streaming, billing, and token usage), see [Streaming output](/developer-guides/text-generation/streaming). 
When streaming with function calling, tool call information arrives in chunks:

- **Tool function name**: returned in the first chunk.

- **Tool call arguments**: returned incrementally across subsequent chunks.


You must aggregate the argument deltas before parsing JSON. Add `stream=True` to your request, then join the chunks:
Python Node.js Copy ```\ntool_calls = {}
for response_chunk in stream:
 delta_tool_calls = response_chunk.choices[0].delta.tool_calls
 if delta_tool_calls:
 for tool_call_chunk in delta_tool_calls:
 call_index = tool_call_chunk.index
 tool_call_chunk.function.arguments = tool_call_chunk.function.arguments or ""
 if call_index not in tool_calls:
 tool_calls[call_index] = tool_call_chunk
 else:
 tool_calls[call_index].function.arguments += tool_call_chunk.function.arguments
print(tool_calls[0].model_dump_json())

``` 
When building the assistant message for the second model call, replace the `tool_calls` field with the aggregated content.
## [​ ](#responses-api) Responses API

The [Responses API](/api-reference/chat/openai-responses) uses a different response format for tool calls. Instead of `tool_calls` on the assistant message, tool calls appear as `function_call` items in the `output` array.
**Step 1 -- Send tools with your request:**
Python Copy ```\nfrom openai import OpenAI
import os

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)

tools = [
 {
 "type": "function",
 "name": "get_current_weather",
 "description": "Useful when you want to check the weather for a specific city.",
 "parameters": {
 "type": "object",
 "properties": {
 "location": {
 "type": "string",
 "description": "City or county, such as Singapore or New York.",
 }
 },
 "required": ["location"],
 },
 },
]

response = client.responses.create(
 model="qwen3.7-plus",
 tools=tools,
 input="What&#x27;s the weather in Singapore?",
)

``` 
**Step 2 -- Parse the function_call output item:**
The response `output` array contains a `function_call` item:
Copy ```\n{
 "type": "function_call",
 "id": "fc_12345",
 "call_id": "call_xxx",
 "name": "get_current_weather",
 "arguments": "{\"location\": \"Singapore\"}"
}

``` 
**Step 3 -- Return the tool result:**
After executing the function, pass the result back using a `function_call_output` item:
Python Copy ```\ntool_result = get_current_weather({"location": "Singapore"})
response = client.responses.create(
 model="qwen3.7-plus",
 tools=tools,
 input=[
 {"type": "message", "role": "user", "content": "What&#x27;s the weather in Singapore?"},
 response.output[0], # The function_call item
 {
 "type": "function_call_output",
 "call_id": response.output[0].call_id,
 "output": tool_result,
 },
 ],
)
print(response.output_text)

``` 
 In the Responses API, tool definitions use a flat structure (`name` and `parameters` at the top level) rather than the nested `function` wrapper used in Chat Completions. See [Responses API](/api-reference/chat/openai-responses) for details. 
## [​ ](#qwen-omni-models) Qwen-Omni models

### [​ ](#non-realtime) Non-realtime

During the tool information retrieval phase, `Qwen3.5-Omni` and `Qwen3-Omni-Flash` differ from other models in two ways:

- 
**Streaming output is required:** Set `stream=True` when retrieving tool information.


- 
**Output text only (recommended):** Set `modalities=["text"]` to avoid unnecessary audio output during tool selection.


See [Audio and video file understanding](/developer-guides/speech/multimodal-speech) for details on Qwen-Omni models.

Python Node.js Copy ```\nfrom openai import OpenAI
import os

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)
tools = [
 {
 "type": "function",
 "function": {
 "name": "get_current_weather",
 "description": "Useful when you want to check the weather for a specific city.",
 "parameters": {
 "type": "object",
 "properties": {
 "location": {
 "type": "string",
 "description": "City or county, such as Singapore or New York.",
 }
 },
 "required": ["location"],
 },
 },
 },
]

completion = client.chat.completions.create(
 model="qwen3-omni-flash",
 messages=[{"role": "user", "content": "What&#x27;s the weather in Singapore?"}],
 modalities=["text"],
 stream=True,
 tools=tools
)

for chunk in completion:
 if chunk.choices:
 delta = chunk.choices[0].delta
 print(delta.tool_calls)

``` 
Running the code gives this output:
Output Copy ```\n[ChoiceDeltaToolCall(index=0, id=&#x27;call_391c8e5787bc4972a388aa&#x27;, function=ChoiceDeltaToolCallFunction(arguments=None, name=&#x27;get_current_weather&#x27;), type=&#x27;function&#x27;)]
[ChoiceDeltaToolCall(index=0, id=&#x27;call_391c8e5787bc4972a388aa&#x27;, function=ChoiceDeltaToolCallFunction(arguments=&#x27; {"location": "Singapore"}&#x27;, name=None), type=&#x27;function&#x27;)]
None

``` 
See the Streaming section above for code to aggregate argument chunks.
### [​ ](#realtime) Realtime

The Qwen3.5-Omni-Plus-Realtime and Qwen3.5-Omni-Flash-Realtime series support tool calling for voice conversations. You can call them through the DashScope SDK or the native WebSocket protocol.
**Workflow:**
After you establish a WebSocket connection, you can pass the tool definitions through `session.update` to start the following interaction flow:
**Phase 1: Speech input and tool calling**

- The user asks a question with their voice. The client captures the audio and sends it to the server (corresponding to the `append_audio()` method).

- After the server&#x27;s Voice Activity Detection (VAD) detects the end of speech, it performs model inference and determines that a tool needs to be called.

- The server returns the tool call information to the client (corresponding to the `response.function_call_arguments.done` event), including the function name (`name`), function parameters (`arguments`), and call ID (`call_id`). Example:


Copy ```\n{
 "type": "response.function_call_arguments.done",
 "response_id": "resp_JnTOsWXlFhKcFohZbtfz6",
 "item_id": "item_Rhcms7CauTNsQprV5S4Hr",
 "output_index": 0,
 "name": "get_current_weather",
 "call_id": "call_2be200f4cafe419b9530dd",
 "arguments": "{\"location\": \"Hangzhou\"}"
}

``` 

- The client executes the corresponding tool function locally based on the function name and parameters, and gets the result.


**Phase 2: Client sends back tool results and triggers the final response**

- The client sends the tool execution result back to the server (corresponding to the `conversation.item.create` event), including the call ID (`call_id`) and execution result (`output`). Example:


Copy ```\n{
 "type": "conversation.item.create",
 "item": {
 "type": "function_call_output",
 "call_id": "call_2be200f4cafe419b9530dd",
 "output": "The weather in Hangzhou today is sunny, with a temperature of 25°C and a light breeze."
 }
}

``` 

- The client continues to send a `response.create` event to trigger the server to generate the final voice response based on the tool execution result.

- The client receives the voice and text returned by the server (corresponding to the `response.audio.delta` and `response.audio_transcript.delta` events) and plays the voice response to the user.


 The Qwen-Omni-Realtime series does not support the `tool_choice` and `parallel_tool_calls` parameters. 
For more details about Qwen-Omni Realtime models and the WebSocket protocol, see [Real-time multimodal speech](/developer-guides/speech/realtime-multimodal-speech), [Client events](/api-reference/real-time-multimodal/client-events), and [Server events](/api-reference/real-time-multimodal/server-events).
DashScope Python SDK DashScope Java SDK WebSocket (Python) Copy ```\nimport os
import uuid
import threading
import traceback
import json
import base64
import signal
import sys
import time
from typing import Dict, Any, Optional, List
import pyaudio
import queue
import contextlib
import dashscope
from dashscope.audio.qwen_omni import *

# ==================== Constant definitions ====================
VOICE = &#x27;Tina&#x27;
MODEL = "qwen3.5-omni-plus-realtime"
WS_URL = "wss://dashscope-intl.aliyuncs.com/api-ws/v1/realtime"
# Configure the API key. If the environment variable is not set, replace the line below with dashscope.api_key = "sk-xxx"
dashscope.api_key = os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;)
AUDIO_SAMPLE_RATE = 16000
AUDIO_CHUNK_SIZE = 3200
OUTPUT_AUDIO_SAMPLE_RATE = 24000


# ==================== Tool definitions ====================
def get_train_price(src: str, dst: str) -> str:
 """Query train ticket prices"""
 return f"The train ticket price from {src} to {dst} is 100-200 CNY."


def get_flight_price(src: str, dst: str) -> str:
 """Query flight ticket prices"""
 return f"The flight ticket price from {src} to {dst} is 200-300 USD."


def get_current_weather(location: str) -> str:
 """Query weather for a specified city"""
 return f"The weather in {location} today is hazy turning to sunny, with a temperature of 4/-4°C and a light breeze."


# Unified OpenAI-format tool definitions
TOOLS = [
 {
 "type": "function",
 "function": {
 "name": "get_current_weather",
 "description": "Useful when you want to check the weather in a specific city.",
 "parameters": {
 "type": "object",
 "properties": {
 "location": {
 "type": "string",
 "description": "City or county, such as Beijing, Hangzhou, or Yuhang District.",
 }
 },
 "required": ["location"],
 },
 },
 },
 {
 "type": "function",
 "function": {
 "name": "get_flight_price",
 "description": "Useful when you want to query flight ticket prices.",
 "parameters": {
 "type": "object",
 "properties": {
 "src": {
 "type": "string",
 "description": "The departure city of the flight, such as Beijing or Hangzhou.",
 },
 "dst": {
 "type": "string",
 "description": "The arrival city of the flight, such as Beijing or Hangzhou.",
 },
 },
 "required": ["src", "dst"],
 },
 },
 },
 {
 "type": "function",
 "function": {
 "name": "get_train_price",
 "description": "Useful when you want to query train ticket prices.",
 "parameters": {
 "type": "object",
 "properties": {
 "src": {
 "type": "string",
 "description": "The departure city of the train, such as Beijing or Hangzhou.",
 },
 "dst": {
 "type": "string",
 "description": "The arrival city of the train, such as Beijing or Hangzhou.",
 },
 },
 "required": ["src", "dst"],
 },
 },
 },
]

# Mapping from tool names to functions
TOOL_FUNCTIONS = {
 "get_current_weather": get_current_weather,
 "get_flight_price": get_flight_price,
 "get_train_price": get_train_price,
}


# ==================== Tool call handling ====================
def handle_tool_call(tool_call_response: Dict[str, Any]) -> Dict[str, Any]:
 """
 Handle tool call requests

 Args:
 tool_call_response: Tool call information including name, arguments, and call_id

 Returns:
 Updated tool call response including the output field
 """
 try:
 function_name = tool_call_response[&#x27;name&#x27;]
 tool_call_arguments = json.loads(tool_call_response[&#x27;arguments&#x27;])

 print(f&#x27;[Tool Call] Start processing: name={function_name}, args={tool_call_arguments}&#x27;)

 # Find the corresponding function
 if function_name not in TOOL_FUNCTIONS:
 tool_call_response[&#x27;output&#x27;] = f"Client did not find the tool: {function_name}"
 print(f&#x27;[Tool Call] Error: Tool not found {function_name}&#x27;)
 return tool_call_response

 # Call the function
 func = TOOL_FUNCTIONS[function_name]
 result = func(**tool_call_arguments)
 tool_call_response[&#x27;output&#x27;] = result

 print(f&#x27;[Tool Call] Complete: {result}&#x27;)
 return tool_call_response

 except Exception as e:
 error_msg = f"Tool call failed: {str(e)}"
 tool_call_response[&#x27;output&#x27;] = error_msg
 print(f&#x27;[Tool Call] Exception: {error_msg}&#x27;)
 traceback.print_exc()
 return tool_call_response


def send_tool_call_response(conversation: OmniRealtimeConversation, response: Dict[str, Any]) -> None:
 """Send tool call response to the server"""
 conversation.create_item({
 "id": &#x27;item_&#x27; + uuid.uuid4().hex,
 "type": "function_call_output",
 "call_id": response[&#x27;call_id&#x27;],
 "output": response["output"],
 })

# ==================== PCM audio player ====================
class PCMPlayer:
 """
 PCM audio player

 Uses a dual-thread architecture for real-time audio playback:
 - Decoder thread: Decodes base64-encoded audio data into raw PCM data
 - Player thread: Writes PCM data to the audio output device

 Supports dynamic addition of audio data, cancellation of playback, saving audio files, etc.
 """

 def __init__(self, pya: pyaudio.PyAudio, sample_rate=24000, chunk_size_ms=100, save_file=False):
 """
 Initialize the PCM player

 Args:
 pya: pyaudio.PyAudio instance
 sample_rate: Audio sampling rate (Hz), default 24000
 chunk_size_ms: Audio chunk size (milliseconds), affects playback cancellation latency, default 100ms
 save_file: Whether to save the played audio to a file (result.pcm), default False
 """

 self.pya = pya
 self.sample_rate = sample_rate
 self.chunk_size_bytes = chunk_size_ms * sample_rate * 2 // 1000
 self.player_stream = pya.open(format=pyaudio.paInt16,
 channels=1,
 rate=sample_rate,
 output=True)

 self.raw_audio_buffer: queue.Queue = queue.Queue()
 self.b64_audio_buffer: queue.Queue = queue.Queue()
 self.status_lock = threading.Lock()
 self.status = &#x27;playing&#x27;
 self.decoder_thread = threading.Thread(target=self.decoder_loop)
 self.player_thread = threading.Thread(target=self.player_loop)
 self.decoder_thread.start()
 self.player_thread.start()
 self.complete_event: threading.Event = None
 self.save_file = save_file
 if self.save_file:
 self.out_file = open(&#x27;result.pcm&#x27;, &#x27;wb&#x27;)

 def decoder_loop(self):
 """Decoder thread: Decodes base64 audio data into raw PCM data"""
 while self.status != &#x27;stop&#x27;:
 recv_audio_b64 = None
 with contextlib.suppress(queue.Empty):
 recv_audio_b64 = self.b64_audio_buffer.get(timeout=0.1)
 if recv_audio_b64 is None:
 continue
 recv_audio_raw = base64.b64decode(recv_audio_b64)
 # push raw audio data into queue by chunk
 for i in range(0, len(recv_audio_raw), self.chunk_size_bytes):
 chunk = recv_audio_raw[i:i + self.chunk_size_bytes]
 self.raw_audio_buffer.put(chunk)
 if self.save_file:
 self.out_file.write(chunk)

 def player_loop(self):
 """Player thread: Writes PCM data to the audio output device"""
 while self.status != &#x27;stop&#x27;:
 recv_audio_raw = None
 with contextlib.suppress(queue.Empty):
 recv_audio_raw = self.raw_audio_buffer.get(timeout=0.1)
 if recv_audio_raw is None:
 if self.complete_event:
 self.complete_event.set()
 continue
 # write chunk to pyaudio audio player, wait until finish playing this chunk.
 self.player_stream.write(recv_audio_raw)

 def cancel_playing(self):
 """Cancel playback: Clear all buffer queues"""
 self.b64_audio_buffer.queue.clear()
 self.raw_audio_buffer.queue.clear()

 def add_data(self, data):
 """Add base64-encoded audio data to the playback queue"""
 self.b64_audio_buffer.put(data)

 def wait_for_complete(self):
 """Wait for playback to complete"""
 self.complete_event = threading.Event()
 self.complete_event.wait()
 self.complete_event = None

 def shutdown(self):
 """Shut down the player and release resources"""
 self.status = &#x27;stop&#x27;
 self.decoder_thread.join()
 self.player_thread.join()
 self.player_stream.close()
 if self.save_file:
 self.out_file.close()

# ==================== Audio manager ====================
class AudioManager:
 """Manages audio input and output resources"""

 def __init__(self):
 self.pya: Optional[pyaudio.PyAudio] = None
 self.mic_stream: Optional[pyaudio.Stream] = None
 self.player: Optional[PCMPlayer] = None

 def initialize(self) -> None:
 """Initialize audio devices"""
 print(&#x27;Initializing audio devices...&#x27;)
 self.pya = pyaudio.PyAudio()
 self.mic_stream = self.pya.open(
 format=pyaudio.paInt16,
 channels=1,
 rate=AUDIO_SAMPLE_RATE,
 input=True
 )
 self.player = PCMPlayer(self.pya, sample_rate=OUTPUT_AUDIO_SAMPLE_RATE)
 print(&#x27;Audio devices initialized&#x27;)

 def read_audio_chunk(self) -> Optional[bytes]:
 """Read an audio data chunk"""
 if not self.mic_stream:
 return None
 try:
 return self.mic_stream.read(AUDIO_CHUNK_SIZE, exception_on_overflow=False)
 except Exception as e:
 print(f&#x27;[Error] Failed to read audio data: {e}&#x27;)
 return None

 def cleanup(self) -> None:
 """Clean up audio resources"""
 print(&#x27;Cleaning up audio resources...&#x27;)
 if self.player:
 self.player.shutdown()
 if self.mic_stream:
 self.mic_stream.close()
 if self.pya:
 self.pya.terminate()
 print(&#x27;Audio resources cleaned up&#x27;)


# ==================== Callback handler ====================
class OmniCallback(OmniRealtimeCallback):
 """Omni real-time conversation callback handler"""

 def __init__(self, audio_manager: AudioManager):
 self.audio_manager = audio_manager
 self.tool_calls: Dict[str, Dict[str, Any]] = {}
 self.all_response_text: str = &#x27;&#x27;
 self.last_package_time: float = 0
 self.is_first_text: bool = True
 self.is_first_audio: bool = True
 self.conversation: Optional[OmniRealtimeConversation] = None

 def set_conversation(self, conversation: OmniRealtimeConversation) -> None:
 """Set the conversation instance reference"""
 self.conversation = conversation

 def on_open(self) -> None:
 """Callback on connection establishment"""
 print(&#x27;Connection established&#x27;)
 self.audio_manager.initialize()
 self.last_package_time = time.time() * 1000
 self.is_first_text = True
 self.is_first_audio = True
 self.tool_calls = {}
 self.all_response_text = &#x27;&#x27;

 def on_close(self, close_status_code: int, close_msg: str) -> None:
 """Callback on connection closure"""
 print(f&#x27;Connection closed: code={close_status_code}, msg={close_msg}&#x27;)
 self.audio_manager.cleanup()
 sys.exit(0)

 def on_event(self, response: Dict[str, Any]) -> None:
 """Handle event callbacks"""
 try:
 event_type = response.get(&#x27;type&#x27;, &#x27;&#x27;)

 # Session created
 if event_type == &#x27;session.created&#x27;:
 print(f&#x27;Session started: {response["session"]["id"]}&#x27;)

 # Speech-to-text completed
 elif event_type == &#x27;conversation.item.input_audio_transcription.completed&#x27;:
 print(f&#x27;User question: {response.get("transcript", "")}&#x27;)

 # Incremental text response
 elif event_type in (&#x27;response.audio_transcript.delta&#x27;, &#x27;response.text.delta&#x27;):
 if self.is_first_text:
 self.is_first_text = False
 latency = time.time() * 1000 - self.last_package_time
 print(f&#x27;Time to first token (VAD end): {latency:.0f} ms&#x27;)

 text = response.get(&#x27;delta&#x27;, &#x27;&#x27;)
 self.all_response_text += text

 # Incremental audio response
 elif event_type == &#x27;response.audio.delta&#x27;:
 if self.is_first_audio:
 self.is_first_audio = False
 latency = time.time() * 1000 - self.last_package_time
 print(f&#x27;Time to first audio chunk (VAD end): {latency:.0f} ms&#x27;)

 audio_interval = time.time() * 1000 - self.last_package_time
 print(f&#x27;Audio interval: {audio_interval:.0f} ms&#x27;)
 self.last_package_time = time.time() * 1000

 recv_audio_b64 = response.get(&#x27;delta&#x27;, &#x27;&#x27;)
 if self.audio_manager.player:
 self.audio_manager.player.add_data(recv_audio_b64)

 # VAD detected speech start
 elif event_type == &#x27;input_audio_buffer.speech_started&#x27;:
 print(&#x27;====== VAD detected speech start ======&#x27;)
 if self.audio_manager.player:
 self.audio_manager.player.cancel_playing()

 # VAD detected speech end
 elif event_type == &#x27;input_audio_buffer.speech_stopped&#x27;:
 print(&#x27;====== VAD detected speech end ======&#x27;)
 self.last_package_time = time.time() * 1000
 self.is_first_text = True
 self.is_first_audio = True
 self.tool_calls = {}

 # Function call arguments done
 elif event_type == &#x27;response.function_call_arguments.done&#x27;:
 print(&#x27;====== Received tool call request ======&#x27;)
 call_id = response.get(&#x27;call_id&#x27;, &#x27;&#x27;)
 self.tool_calls[call_id] = response.copy()
 self.tool_calls[call_id][&#x27;processed&#x27;] = False

 # Response done
 elif event_type == &#x27;response.done&#x27;:
 print(&#x27;====== Response done ======&#x27;)
 print(f&#x27;Full response: {self.all_response_text}&#x27;)

 if self.conversation:
 response_id = self.conversation.get_last_response_id()
 text_delay = self.conversation.get_last_first_text_delay()
 audio_delay = self.conversation.get_last_first_audio_delay()

 # Print detailed metrics only when all are available
 if response_id is not None and text_delay is not None and audio_delay is not None:
 print(f&#x27;[Metric] Response ID: {response_id}, &#x27;
 f&#x27;Time to first token: {text_delay:.0f}ms, &#x27;
 f&#x27;Time to first audio chunk: {audio_delay:.0f}ms&#x27;)
 else:
 print(&#x27;[Metric] Metric info not available (might be a response after a tool call)&#x27;)

 self.all_response_text = &#x27;&#x27;

 except Exception as e:
 print(f&#x27;[Error] Exception in event handling: {e}&#x27;)
 traceback.print_exc()

 def process_pending_tool_calls(self) -> bool:
 """
 Process pending tool calls

 Returns:
 Whether there are new tool calls to respond to
 """
 has_pending = False

 for call_id, tool_call in self.tool_calls.items():
 if not tool_call.get(&#x27;processed&#x27;, False):
 has_pending = True
 tool_call[&#x27;processed&#x27;] = True

 # Handle the tool call
 result = handle_tool_call(tool_call)

 # Send the result to the server
 if self.conversation:
 send_tool_call_response(self.conversation, result)

 return has_pending


# ==================== Main program ====================
def main():
 """Main function"""
 print(&#x27;Initializing Omni real-time conversation...&#x27;)

 # Create audio manager
 audio_manager = AudioManager()

 # Create callback handler
 callback = OmniCallback(audio_manager)

 # Create conversation instance
 conversation = OmniRealtimeConversation(
 api_key=dashscope.api_key,
 url=WS_URL,
 model=MODEL,
 callback=callback,
 )

 # Set conversation reference in callback
 callback.set_conversation(conversation)

 # Establish connection
 conversation.connect()

 # Configure session parameters
 omni_output_modalities = [MultiModality.AUDIO, MultiModality.TEXT]

 conversation.update_session(
 output_modalities=omni_output_modalities,
 voice=VOICE,
 input_audio_format=AudioFormat.PCM_16000HZ_MONO_16BIT,
 output_audio_format=AudioFormat.PCM_24000HZ_MONO_16BIT,
 enable_input_audio_transcription=True,
 enable_turn_detection=True,
 turn_detection_type=&#x27;server_vad&#x27;,
 tools=TOOLS,
 )

 # Set signal handler
 def signal_handler(sig, frame):
 print(&#x27;\nCtrl+C received, stopping...&#x27;)
 conversation.close()
 audio_manager.cleanup()
 print(&#x27;Omni real-time conversation stopped&#x27;)
 sys.exit(0)

 signal.signal(signal.SIGINT, signal_handler)
 print("Press Ctrl+C to stop the conversation...\n")

 # Main loop: continuously send audio and check for tool calls
 try:
 while True:
 # Process pending tool calls
 has_tool_calls = callback.process_pending_tool_calls()

 if has_tool_calls:
 print("*** Tool call complete, creating new response ***")
 conversation.create_response(
 instructions=None,
 output_modalities=omni_output_modalities
 )
 print(&#x27;====== Tool call processing complete ======\n&#x27;)

 # Read and send audio data
 audio_data = audio_manager.read_audio_chunk()
 if audio_data:
 audio_b64 = base64.b64encode(audio_data).decode(&#x27;ascii&#x27;)
 conversation.append_audio(audio_b64)
 else:
 break

 except KeyboardInterrupt:
 signal_handler(signal.SIGINT, None)
 except Exception as e:
 print(f&#x27;[Error] Main loop exception: {e}&#x27;)
 traceback.print_exc()
 finally:
 conversation.close()
 audio_manager.cleanup()


if __name__ == &#x27;__main__&#x27;:
 main()

``` 
## [​ ](#thinking-mode) Thinking mode

Deep thinking models reason before generating tool calls. Set `enable_thinking=True` to see the model&#x27;s reasoning process. The response includes `reasoning_content` before tool calls:

- The model outputs `reasoning_content` showing its analysis of user intent, tool selection, and parameter planning.

- The model then outputs the `tool_calls` as usual.


When passing the assistant message back in subsequent requests, you must include the `reasoning_content` field.
 With thinking mode enabled, the `tool_choice` parameter only supports `"auto"` (default) or `"none"`. 
Python Copy ```\ncompletion = client.chat.completions.create(
 model="qwen3.7-plus", # Use a thinking-capable model
 messages=messages,
 tools=tools,
 extra_body={"enable_thinking": True},
 stream=True,
)

``` 
For complete streaming parse code with `reasoning_content` handling, see [Thinking mode](/developer-guides/text-generation/thinking).
## [​ ](#best-practices) Best practices


- **Test tool selection accuracy**: Build an evaluation dataset mirroring real scenarios. Track tool selection accuracy, parameter extraction accuracy, and end-to-end success rate.

- **Optimize tool descriptions**: When the model selects the wrong tool or extracts wrong parameters, refine descriptions and system prompts before upgrading models.

- **Keep candidate tool sets small**: Limit tools passed to the model to no more than 20. Use a routing layer (semantic search, keyword filtering, or a lightweight LLM router) to pre-filter large tool libraries.

- **Apply least-privilege principle**: Default to read-only tools. Never give the model direct access to dangerous operations (code execution, file deletion, financial transfers).

- **Add human confirmation for write operations**: The model can generate action requests, but irreversible actions (sending email, modifying data) should require user confirmation.

- **Set timeouts and provide fallback responses**: Assign independent timeouts to each step. On failure, return a clear message such as "Sorry, I cannot retrieve that information right now. Please try again later."

- **Show progress in the UI**: Display messages like "Looking up weather for you..." when starting tool execution.


## [​ ](#pass-tool-information-via-system-message) Pass tool information via system message

 Alternative: embed tool info in system prompt

 We recommend using the `tools` parameter (shown throughout this guide). If you need direct control over the tool prompt, embed tool definitions in the system message using this template:Python Copy ```\nimport os
from openai import OpenAI
import json

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)

tools = [
 {
 "type": "function",
 "function": {
 "name": "get_current_time",
 "description": "Useful when you want to know the current time.",
 "parameters": {}
 }
 },
 {
 "type": "function",
 "function": {
 "name": "get_current_weather",
 "description": "Useful when you want to check the weather for a specific city.",
 "parameters": {
 "type": "object",
 "properties": {
 "location": {
 "type": "string",
 "description": "City or county, such as Singapore or New York."
 }
 },
 "required": ["location"]
 }
 }
 }
]

custom_prompt = "You are a helpful assistant."
tools_content = "\n".join(json.dumps(tool, ensure_ascii=False) for tool in tools)

system_prompt = f"""{custom_prompt}

# Tools

You may call one or more functions to assist with the user query.

You are provided with function signatures within <tools></tools> XML tags:
<tools>
{tools_content}
</tools>

For each function call, return a json object with function name and arguments within <tool_call></tool_call> XML tags:
<tool_call>
{{"name": <function-name>, "arguments": <args-json-object>}}
</tool_call>"""

messages = [
 {"role": "system", "content": system_prompt},
 {"role": "user", "content": "What time is it?"}
]

completion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=messages,
)
print(completion.choices[0].message.content)

``` After running the code, use an XML parser to extract tool call information from between the `<tool_call>` and `</tool_call>` tags. 
## [​ ](#billing) Billing

In addition to tokens in the messages array, tool descriptions also count as input tokens and are billed as part of the prompt.
## [​ ](#error-codes) Error codes

If a call fails, see [Error messages](/api-reference/preparation/error-messages). [Previous ](/developer-guides/embeddings/reranking)[Web search Ground model responses in real-time web data Next ](/developer-guides/text-generation/web-search)
