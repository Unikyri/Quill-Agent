# OpenAI compatible reranking

> **Source:** https://docs.qwencloud.com/api-reference/rerank/openai-rerank

POST/reranks qwen3-rerank cURL

 Copy ```\ncurl --request POST \
 --url https://dashscope-intl.aliyuncs.com/compatible-api/v1/reranks \
 --header "Authorization: Bearer $DASHSCOPE_API_KEY" \
 --header "Content-Type: application/json" \
 --data '{
 "model": "qwen3-rerank",
 "documents": [
 "Rerank models are widely used in search engines and recommendation systems. They sort candidate documents based on text relevance.",
 "Quantum computing is a cutting-edge field of computer science.",
 "The development of pre-trained language models has brought new advancements to rerank models."
 ],
 "query": "What is a rerank model?",
 "top_n": 2,
 "instruct": "Given a web search query, retrieve relevant passages that answer the query."
}'
``` 200 400 401 429 Copy ```\n{
 "id": "&#x3C;string>",
 "object": "list",
 "model": "qwen3-rerank",
 "results": [
 {
 "document": {
 "text": "&#x3C;string>"
 },
 "index": 0,
 "relevance_score": 0.9334521178273196
 }
 ],
 "usage": {
 "total_tokens": 0
 }
}
``` Rerank documents by semantic relevance to a query using qwen3-rerank.
 The gte-rerank model will be discontinued on May 30, 2026. Switch to qwen3-rerank for continued service. 
 Before you call the API, [get an API key](/api-reference/preparation/api-key) and [set it as an environment variable](/api-reference/preparation/export-api-key-env). If you use the OpenAI SDK, [install it](/api-reference/preparation/install-sdk) first. 
**Supported model:** qwen3-rerank only.
## [​ ](#endpoint) Endpoint


- HTTP: `POST https://dashscope-intl.aliyuncs.com/compatible-api/v1/reranks`

- SDK `base_url`: `https://dashscope-intl.aliyuncs.com/compatible-api/v1`


## [​ ](#model-overview) Model overview

ModelMax DocumentsMax Tokens/DocMax Request TokensLanguagesUse Casesqwen3-rerank5004,000120,000100+ languagesText semantic search, RAG 
For pricing, see [Models](https://www.qwencloud.com/models).
**Parameter definitions**:

- **Max Tokens/Doc**: Maximum token count per query or document. Content exceeding this limit is truncated, which may affect ranking accuracy.

- **Max Documents**: Maximum number of documents per request.

- **Max Request Tokens**: Calculated as `Query Tokens x Document Count + Total Document Tokens`. Must not exceed the limit.


 ### Authorizations
 [​ ](#authorization) Authorizationstring header required Qwen Cloud API Key. Create one in the [console](https://home.qwencloud.com/api-keys).

 ### Body
application/json [​ ](#model) modelenum<string> required Model name. Must be `qwen3-rerank` for the text reranking endpoint.

 Available options:qwen3-rerank Example:qwen3-rerank [​ ](#query) querystring required Query text. Max 4,000 tokens.

 Example:What is a reranking model [​ ](#documents) documentsstring[] required Documents to rank. An array of strings. Max 500 documents.

 Example: Copy ```\n[
 "Reranking models are widely used in search engines and recommendation systems to sort candidates by relevance",
 "Quantum computing is a frontier field of computer science",
 "The development of pre-trained language models has brought new advances to reranking"
]
``` [​ ](#top-n) top_ninteger Return only the top N results. Defaults to returning all documents.

 Example:2 Required range:x >= 1 [​ ](#instruct) instructstring Custom ranking task instruction. English recommended. Default behavior is QA retrieval: `"Given a web search query, retrieve relevant passages that answer the query."`

 Example:Given a web search query, retrieve relevant passages that answer the query. ### Response
200-application/json [​ ](#id) idstring Unique request identifier.

 [​ ](#object) objectstring Object type. Always `list`.

 Example:list [​ ](#model) modelstring Model used for reranking.

 Example:qwen3-rerank [​ ](#results) resultsobject[] Ranked results, sorted by `relevance_score` descending.

 Show child attributes

 [​ ](#resultsdocument) results.documentobject Original document. Only returned when `return_documents` is `true`.

 Show child attributes

 [​ ](#resultsdocumenttext) results.document.textstring Text content of the document.

 [​ ](#resultsindex) results.indexinteger Original position in the input `documents` list.

 Example:0 [​ ](#resultsrelevance-score) results.relevance_scorenumber Relevance score between 0.0 and 1.0. Higher means more relevant. This is a relative score for the current request and should not be compared across requests.

 Example:0.9334521178273196 [​ ](#usage) usageobject Token usage statistics.

 Show child attributes

 [​ ](#usagetotal-tokens) usage.total_tokensinteger Total tokens consumed by this request.
