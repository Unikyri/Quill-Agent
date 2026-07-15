# DashScope embedding

> **Source:** https://docs.qwencloud.com/api-reference/text-embedding/dashscope-embedding

POST/api/v1/services/embeddings/text-embedding/text-embedding Python

 Input string

 Copy import dashscope
from http import HTTPStatus

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

resp = dashscope.TextEmbedding.call(
 model=dashscope.TextEmbedding.Models.text_embedding_v4,
 input=&#x27;Semantic search finds documents by meaning rather than exact keyword matching. Text embeddings map words and sentences into high-dimensional vector spaces. Retrieval-augmented generation combines search results with language models. Document clustering groups similar texts based on their vector representations.&#x27;,
 dimension=1024,
 output_type="dense&sparse"
)

print(resp) if resp.status_code == HTTPStatus.OK else print(resp)
 200 400 Copy {
 "status_code": 200,
 "request_id": "1ba94ac8-e058-99bc-9cc1-7fdb37940a46",
 "code": "",
 "message": "",
 "output": {
 "embeddings": [
 {
 "sparse_embedding": [
 {
 "index": 7149,
 "value": 0.829,
 "token": "wind"
 },
 {
 "index": 111290,
 "value": 0.9004,
 "token": "sorrow"
 }
 ],
 "embedding": [
 -0.006929283495992422,
 -0.005336422007530928
 ],
 "text_index": 0
 }
 ]
 },
 "usage": {
 "total_tokens": 27
 }
} Convert text to vectors for semantic search, recommendations, clustering, and classification.
 Before you begin: [get an API key](/api-reference/preparation/api-key), [set it as an environment variable](/api-reference/preparation/export-api-key-env), and [install the DashScope SDK](/api-reference/preparation/install-sdk) if you use the SDK. 
## [​ ](#set-the-sdk-base-url) Set the SDK base URL

**Python SDK:**
Copy ```\nimport dashscope
dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

``` 
**Java SDK:**
Copy ```\nimport com.alibaba.dashscope.utils.Constants;
Constants.baseHttpApiUrl = "https://dashscope-intl.aliyuncs.com/api/v1";

``` 
## [​ ](#supported-models) Supported models

ModelDimensionsMax tokensBatch sizeLanguagestext-embedding-v42048, 1536, 1024 (default), 768, 512, 256, 128, 648,19210100+text-embedding-v31024 (default), 768, 5128,1921050+ 
For pricing, see [Models](https://www.qwencloud.com/models).
## [​ ](#input-formats) Input formats


- **Single string**: up to 8,192 tokens

- **Array**: up to 10 strings, each up to 8,192 tokens

- **Text file**: up to 10 lines, each line up to 8,192 tokens


## [​ ](#dashscope-specific-features) DashScope-specific features


- `text_type`: set to `query` or `document` for asymmetric tasks like retrieval.

- `output_type`: return sparse vectors (`dense&sparse`) for hybrid search (v3/v4 only).

- `instruct`: add a task description to improve accuracy by ~1-5% (v4 only; English recommended).


See [Rate limits](/developer-guides/administration/rate-limits). ### Authorizations
 [​ ](#authorization) Authorizationstring header required DashScope API key. Obtain from the Qwen Cloud console.

 ### Body
application/json [​ ](#model) modelstring required The model to call. Supported values: `text-embedding-v4`, `text-embedding-v3`.

 Example:text-embedding-v4 [​ ](#input) inputobject required The input object containing texts to embed.

 Show child attributes

 [​ ](#inputtexts) input.textsstring required The input text to process. The value can be a string or an array of strings. A string input can contain up to 8,192 tokens. A list of strings can contain up to 10 items, each with up to 8,192 tokens.

 [​ ](#parameters) parametersobject Optional parameters for the embedding request.

 Show child attributes

 [​ ](#parameterstext-type) parameters.text_typeenum<string> default"document" Specifies the text type for asymmetric tasks. For retrieval, distinguish between `query` and `document` for better results. For symmetric tasks (clustering, categorization), use the default `document`.

 Available options:query,document [​ ](#parametersdimension) parameters.dimensionenum<integer> default1024 The embedding dimensions. Valid values: 2048 (text-embedding-v4 only), 1536 (text-embedding-v4 only), 1024, 768, 512, 256 (text-embedding-v4 only), 128 (text-embedding-v4 only), or 64 (text-embedding-v4 only). Default: 1024.

 Available options:2048,1536,1024,768,512,256,128,64 [​ ](#parametersoutput-type) parameters.output_typeenum<string> default"dense" Specifies whether to output a sparse vector representation. Applies only to text-embedding-v3 and text-embedding-v4. Default: `dense` (returns only dense vector).

 Available options:dense,sparse,dense&sparse [​ ](#parametersinstruct) parameters.instructstring A custom task description. Takes effect only with text-embedding-v4 when `text_type` is `query`. English descriptions are recommended (~1-5% performance improvement).

 ### Response
200-application/json [​ ](#status-code) status_codeinteger default"200" The status code. A value of 200 indicates a successful request.

 [​ ](#request-id) request_idstring default"1ba94ac8-e058-99bc-9cc1-7fdb37940a46" The unique request ID. Use this ID to trace and troubleshoot requests.

 [​ ](#code) codestring default"" If the request fails, this indicates the error code. Empty on success.

 [​ ](#message) messagestring default"" If the request fails, this indicates the detailed error message. Empty on success.

 [​ ](#output) outputobject The output data for the task.

 Show child attributes

 [​ ](#outputembeddings) output.embeddingsobject[] An array of structures. Each structure contains the output for the corresponding input text.

 Show child attributes

 [​ ](#outputembeddingssparse-embedding) output.embeddings.sparse_embeddingobject[] The sparse vector representation. Applies only to text-embedding-v3 and text-embedding-v4 when `output_type` includes `sparse`.

 Show child attributes

 [​ ](#outputembeddingssparse-embeddingindex) output.embeddings.sparse_embedding.indexinteger The index of the word or character in the vocabulary.

 [​ ](#outputembeddingssparse-embeddingvalue) output.embeddings.sparse_embedding.valuenumber The weight or importance score of the token. Higher values indicate greater importance.

 [​ ](#outputembeddingssparse-embeddingtoken) output.embeddings.sparse_embedding.tokenstring The actual text unit or word from the vocabulary.

 [​ ](#outputembeddingsembedding) output.embeddings.embeddingnumber[] The dense vector representation (dense embedding).

 [​ ](#outputembeddingstext-index) output.embeddings.text_indexinteger The index of the input text in the input array that corresponds to the result in this structure.

 [​ ](#usage) usageobject default"{\"total_tokens\":27}" Token usage statistics.

 Show child attributes

 [​ ](#usagetotal-tokens) usage.total_tokensinteger The total number of tokens in the request. Calculated based on the model&#x27;s tokenizer.
