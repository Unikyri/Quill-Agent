# Structured output

> **Source:** https://docs.qwencloud.com/developer-guides/text-generation/structured-output

Guaranteed valid JSON

 Copy page Structured output guarantees the model returns valid JSON. Set `response_format` to `{"type": "json_object"}` and include the word "JSON" in your prompt.
- OpenAI compatible 
- DashScope 

 Python Node.js curl Copy ```\nfrom openai import OpenAI
import os

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)

completion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=[
 {"role": "system", "content": "Extract the name and age. Return JSON."},
 {"role": "user", "content": "My name is Alex Brown, I am 34 years old."},
 ],
 response_format={"type": "json_object"}, # <-- force valid JSON output
)
print(completion.choices[0].message.content)

``` Python Java curl Copy ```\nimport os
import dashscope
from dashscope import MultiModalConversation

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

response = MultiModalConversation.call(
 api_key=os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;),
 model="qwen3.7-plus",
 messages=[
 {"role": "system", "content": [{"text": "Extract the name and age. Return JSON."}]},
 {"role": "user", "content": [{"text": "My name is Alex Brown, I am 34 years old."}]},
 ],
 response_format={&#x27;type&#x27;: &#x27;json_object&#x27;}, # <-- force valid JSON output
)
print(response.output.choices[0].message.content[0]["text"])

``` 
**Output:**
Copy ```\n{"name": "Alex Brown", "age": 34}

``` 
 For more reliable results, describe the expected fields, types, required/optional status, and provide examples in your prompt. 
## [​ ](#thinking-mode-workaround) Thinking mode workaround

Structured output is not supported in [thinking mode](/developer-guides/text-generation/thinking). As a workaround, parse the thinking-mode output and fall back to a fast model for JSON repair.
Python Copy ```\nimport json
from openai import OpenAI
import os

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)

# Step 1: Get output from thinking mode (streaming required)
completion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=[
 {"role": "system", "content": system_prompt},
 {"role": "user", "content": user_input},
 ],
 extra_body={"enable_thinking": True}, # <-- no response_format here
 stream=True,
)
json_string = ""
for chunk in completion:
 if chunk.choices[0].delta.content:
 json_string += chunk.choices[0].delta.content

# Step 2: Parse, or repair with a fast model
try:
 result = json.loads(json_string)
except json.JSONDecodeError:
 repair = client.chat.completions.create(
 model="qwen3.5-flash", # <-- fast model fixes the JSON
 messages=[
 {"role": "system", "content": "Fix this to valid JSON."},
 {"role": "user", "content": json_string},
 ],
 response_format={"type": "json_object"},
 )
 result = json.loads(repair.choices[0].message.content)

``` 
## [​ ](#notes) Notes


- **Supported models**: Qwen3.7, Qwen3.6, Qwen3.5, Qwen3, Qwen3-Coder, Qwen2.5, and legacy (Plus/Max/Flash/Turbo) — non-thinking mode. Also works with vision models (qwen3-vl-plus, etc.) when passing images or video as described in [Image and video understanding](/developer-guides/multimodal/vision).

- **Don&#x27;t set `max_tokens`**: Truncated output produces invalid JSON that breaks downstream parsing. Use the default (model&#x27;s maximum output limit).

- **Validate output**: JSON Object mode guarantees valid JSON but not schema conformance. Validate with jsonschema (Python), Ajv (JavaScript), or Everit (Java) before passing to downstream systems.


 [Previous ](/developer-guides/text-generation/mcp)[Thinking Solve complex tasks with step-by-step thinking Next ](/developer-guides/text-generation/thinking)
