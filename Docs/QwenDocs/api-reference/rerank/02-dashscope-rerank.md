# DashScope reranking

> **Source:** https://docs.qwencloud.com/api-reference/rerank/dashscope-rerank

POST/services/rerank/text-rerank/text-rerank cURL

 qwen3-rerank

 Copy ```\ncurl --request POST \
 --url https://dashscope-intl.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank \
 --header "Authorization: Bearer $DASHSCOPE_API_KEY" \
 --header "Content-Type: application/json" \
 --data '{
 "model": "qwen3-rerank",
 "input": {
 "query": "What is a rerank model?",
 "documents": [
 "Rerank models are widely used in search engines and recommendation systems. They sort candidate documents based on text relevance.",
 "Quantum computing is a cutting-edge field of computer science.",
 "The development of pre-trained language models has brought new advancements to rerank models."
 ]
 },
 "parameters": {
 "return_documents": true,
 "top_n": 2
 }
}'
``` 200 400 401 429 Copy ```\n{
 "output": {
 "results": [
 {
 "document": {
 "text": "&#x3C;string>"
 },
 "index": 0,
 "relevance_score": 0.9334521178273196
 }
 ]
 },
 "usage": {
 "total_tokens": 0
 },
 "request_id": "85ba5752-1900-47d2-8896-23f99b13f6e1"
}
``` Rerank documents by semantic relevance to a query using qwen3-rerank. Uses a nested request structure with `input` and `parameters` wrappers.
 The gte-rerank model will be discontinued on May 30, 2026. Switch to qwen3-rerank for continued service. 
 Before you begin: [get an API key](/api-reference/preparation/api-key), [set it as an environment variable](/api-reference/preparation/export-api-key-env), and [install the DashScope SDK](/api-reference/preparation/install-sdk) if you use the SDK. 
## [​ ](#endpoint) Endpoint


- HTTP: `POST https://dashscope-intl.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank`

- SDK `base_http_api_url`: `https://dashscope-intl.aliyuncs.com/api/v1`


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
application/json [​ ](#model) modelenum<string> required Model name. Must be `qwen3-rerank`.

 Available options:qwen3-rerank Example:qwen3-rerank [​ ](#input) inputobject required Input data containing the query and documents to rank.

 Show child attributes

 [​ ](#inputquery) input.querystring required Query text. Max 4,000 tokens.

 Example:What is a reranking model [​ ](#inputdocuments) input.documentsstring[] required Documents to rank. An array of strings. Max 500 documents.

 [​ ](#parameters) parametersobject Parameters for the reranking request. Must be wrapped inside this `parameters` object.

 Show child attributes

 [​ ](#parameterstop-n) parameters.top_ninteger Return only the top N results. Defaults to returning all documents.

 Example:2 Required range:x >= 1 [​ ](#parametersreturn-documents) parameters.return_documentsboolean defaultfalse Include original document text in results. Default: `false`.

 Example:true [​ ](#parametersinstruct) parameters.instructstring Custom ranking task instruction. English recommended. Default behavior is QA retrieval: `"Given a web search query, retrieve relevant passages that answer the query."`

 Example:Given a web search query, retrieve relevant passages that answer the query. ### Response
200-application/json [​ ](#output) outputobject Output wrapper containing ranked results.

 Show child attributes

 [​ ](#outputresults) output.resultsobject[] Ranked results, sorted by `relevance_score` descending.

 Show child attributes

 [​ ](#outputresultsdocument) output.results.documentobject Original document. Only returned when `return_documents` is `true`.

 Show child attributes

 [​ ](#outputresultsdocumenttext) output.results.document.textstring Text content of the document.

 [​ ](#outputresultsindex) output.results.indexinteger Original position in the input `documents` list.

 Example:0 [​ ](#outputresultsrelevance-score) output.results.relevance_scorenumber Relevance score between 0.0 and 1.0. Higher means more relevant. This is a relative score for the current request and should not be compared across requests.

 Example:0.9334521178273196 [​ ](#usage) usageobject Token usage statistics.

 Show child attributes

 [​ ](#usagetotal-tokens) usage.total_tokensinteger Total tokens consumed by this request.

 [​ ](#request-id) request_idstring Unique request identifier.

 Example:85ba5752-1900-47d2-8896-23f99b13f6e1
