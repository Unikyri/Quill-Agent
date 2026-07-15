# OpenAI compatible embedding

> **Source:** https://docs.qwencloud.com/api-reference/text-embedding/openai-embedding

POST/compatible-mode/v1/embeddings Python

 Input string

 Copy ```\nimport os
from openai import OpenAI

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"), # If you have not configured an environment variable, replace the placeholder with your API key.
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1" 
)

completion = client.embeddings.create(
 model="text-embedding-v4",
 input='The clothes are of good quality and look good, definitely worth the wait. I love them.',
 dimensions=1024,
 encoding_format="float"
)

print(completion.model_dump_json())
``` 200 400 Copy ```\n{
 "data": [
 {
 "embedding": [
 -0.0695386752486229,
 0.030681096017360687
 ],
 "index": 0,
 "object": "embedding"
 },
 {
 "embedding": [
 -0.06348952651023865,
 0.060446035116910934
 ],
 "index": 5,
 "object": "embedding"
 }
 ],
 "model": "text-embedding-v4",
 "object": "list",
 "usage": {
 "prompt_tokens": 184,
 "total_tokens": 184
 },
 "id": "73591b79-d194-9bca-8bb5-xxxxxxxxxxxx"
}
``` Before you call the API, [get an API key](/api-reference/preparation/api-key) and [set it as an environment variable](/api-reference/preparation/export-api-key-env). If you use the OpenAI SDK, [install it](/api-reference/preparation/install-sdk) first. 
## [​ ](#supported-models) Supported models

ModelDimensionsMax tokensBatch sizeLanguagestext-embedding-v42048, 1536, 1024 (default), 768, 512, 256, 128, 648,19210100+ major languagestext-embedding-v31024 (default), 768, 5128,1921050+ languages 
For pricing, see [Models](https://www.qwencloud.com/models).
## [​ ](#endpoints) Endpoints

 International API keys can only use international endpoints. China API keys can only use China endpoints. 
API key fromEndpointInternational`https://dashscope-intl.aliyuncs.com/compatible-mode/v1/embeddings`China`https://dashscope.aliyuncs.com/compatible-mode/v1/embeddings` 
SDKs use the base URL without `/embeddings`. ### Authorizations
 [​ ](#authorization) Authorizationstring header required DashScope API key. Obtain from the Qwen Cloud console.

 ### Body
application/json [​ ](#model) modelstring required The model to call. Supported values: `text-embedding-v4`, `text-embedding-v3`.

 Example:text-embedding-v4 [​ ](#input) inputstring required The input text to process. The value can be a string, an array of strings, or a file. When the input is a string, it can contain up to 8,192 tokens. When the input is a list of strings or a file, it can contain up to 10 items (lines), each with a maximum of 8,192 tokens.

 [​ ](#dimensions) dimensionsenum<integer> default1024 The embedding dimensions. Valid values: 2048 (text-embedding-v4 only), 1536 (text-embedding-v4 only), 1024, 768, 512, 256 (text-embedding-v4 only), 128 (text-embedding-v4 only), or 64 (text-embedding-v4 only). Default: 1024.

 Available options:2048,1536,1024,768,512,256,128,64 [​ ](#encoding-format) encoding_formatenum<string> default"float" The format of the returned embedding. Only `float` is supported.

 Available options:float ### Response
200-application/json [​ ](#data) dataobject[] The output data for the task.

 Show child attributes

 [​ ](#dataembedding) data.embeddingnumber[] The embedding vector, which is an array of floating-point numbers.

 [​ ](#dataindex) data.indexinteger The index of the input text in the input array that corresponds to the result in this structure.

 [​ ](#dataobject) data.objectstring default"embedding" The type of object returned by the call. Default: `embedding`.

 [​ ](#model) modelstring default"text-embedding-v4" The model that was called.

 [​ ](#object) objectstring default"list" The type of data returned by the call. Default: `list`.

 [​ ](#usage) usageobject default"{\"prompt_tokens\":184,\"total_tokens\":184}" Token usage statistics.

 Show child attributes

 [​ ](#usageprompt-tokens) usage.prompt_tokensinteger The number of tokens in the input text.

 [​ ](#usagetotal-tokens) usage.total_tokensinteger The total number of tokens in the request. Calculated based on the model&#x27;s tokenizer.

 [​ ](#id) idstring default"73591b79-d194-9bca-8bb5-xxxxxxxxxxxx" The unique request ID. Use this ID to trace and troubleshoot requests.
