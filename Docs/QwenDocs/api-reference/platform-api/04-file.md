# File

> **Source:** https://docs.qwencloud.com/api-reference/platform-api/file

File management API

 Copy page Upload files for document analysis or batch processing using the OpenAI-compatible Files API.
## [​ ](#usage) Usage

Use the OpenAI SDK (Python/Node.js) or HTTP to upload, query, list, and delete files.
**Prerequisites:**

- 
An API key: [Create an API key](/api-reference/preparation/api-key) and [Export the API key as environment variable](/api-reference/preparation/export-api-key-env).


- 
To use the OpenAI SDK, install the [OpenAI SDK](/api-reference/preparation/install-sdk).


## [​ ](#model-availability) Model availability

File IDs work with:

- [Batch processing](/developer-guides/text-generation/batch): Batch file upload.

- [Fine-tuning](/developer-guides/fine-tuning/overview): Training data file upload.


## [​ ](#endpoints) Endpoints

EndpointDescription[Upload file](/api-reference/platform-api/file/upload-file)Upload a file for document analysis or batch processing[Retrieve file](/api-reference/platform-api/file/retrieve-file)Retrieve file details by ID[List files](/api-reference/platform-api/file/list-files)List all files in your account[Delete file](/api-reference/platform-api/file/delete-file)Delete a file by ID 
## [​ ](#getting-started) Getting started

### [​ ](#upload-a-file) Upload a file

You can store up to 10,000 files and 100 GB total. Files never expire.
### [​ ](#for-document-analysis) For document analysis

Set `purpose` to `file-extract`. Supported formats: text files (TXT, DOCX, PDF, XLSX, EPUB, MOBI, MD, CSV, JSON) and images (BMP, PNG, JPG/JPEG, GIF, scanned PDFs). Maximum file size: **150 MB**.
#### [​ ](#request-examples) Request examples

- Python 
- Node.js 
- cURL 

 Copy ```\nimport os
from pathlib import Path
from openai import OpenAI

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)

# test.txt is a local sample file.
file_object = client.files.create(file=Path("test.txt"), purpose="file-extract")

print(file_object.model_dump_json())

``` Copy ```\nimport OpenAI from "openai";
import fs from "fs";

const client = new OpenAI({
 apiKey: process.env.DASHSCOPE_API_KEY,
 baseURL: "https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
});

const fileObject = await client.files.create({
 file: fs.createReadStream("test.txt"),
 purpose: "file-extract",
});
console.log(fileObject);

``` Copy ```\ncurl -X POST https://dashscope-intl.aliyuncs.com/compatible-mode/v1/files \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY" \
 -F &#x27;file=@"test.txt"&#x27; \
 -F &#x27;purpose="file-extract"&#x27;

``` 
#### [​ ](#sample-response) Sample response

Copy ```\n{
 "id": "file-fe-xxx",
 "bytes": 2055,
 "created_at": 1729065448,
 "filename": "test.txt",
 "object": "file",
 "purpose": "file-extract",
 "status": "processed",
 "status_details": null
}

``` 
### [​ ](#for-batch-processing) For batch processing

Set `purpose` to `batch`. Upload a JSONL file conforming to [Batch file requirements](/developer-guides/text-generation/batch#input-file-format). Maximum file size: **500 MB**.
 For batch calls, see [Batch API](/developer-guides/text-generation/batch). 
#### [​ ](#request-examples-2) Request examples

- Python 
- Node.js 
- cURL 

 Copy ```\nimport os
from pathlib import Path
from openai import OpenAI

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)

# test.jsonl is a local sample file.
file_object = client.files.create(file=Path("test.jsonl"), purpose="batch")

print(file_object.model_dump_json())

``` Copy ```\nimport OpenAI from "openai";
import fs from "fs";

const client = new OpenAI({
 apiKey: process.env.DASHSCOPE_API_KEY,
 baseURL: "https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
});

const fileObject = await client.files.create({
 file: fs.createReadStream("test.jsonl"),
 purpose: "batch",
});
console.log(fileObject);

``` Copy ```\ncurl -X POST https://dashscope-intl.aliyuncs.com/compatible-mode/v1/files \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY" \
 -F &#x27;file=@"test.jsonl"&#x27; \
 -F &#x27;purpose="batch"&#x27;

``` 
#### [​ ](#sample-response-2) Sample response

Copy ```\n{
 "id": "file-batch-xxx",
 "bytes": 231,
 "created_at": 1729065815,
 "filename": "test.jsonl",
 "object": "file",
 "purpose": "batch",
 "status": "processed",
 "status_details": null
}

``` 
## [​ ](#billing) Billing

File operations (upload, storage, query, delete) are free. You are charged only for model inference tokens (input + output).
## [​ ](#rate-limits) Rate limits

Rate limits: upload (3 QPS), query/list/delete (10 QPS combined).
## [​ ](#going-live) Going live


- 
**Periodic cleanup**: Delete unused files to stay under the 10,000-file limit.


- 
**Status verification**: Confirm `status="processed"` before using uploaded files.


- 
**Rate limit awareness**: Upload (3 QPS), query/list/delete (10 QPS). Implement retry with exponential backoff.


- 
**Error handling**: Handle network timeouts, API errors (invalid format, size exceeded), and file limit errors. See [Error codes](/api-reference/preparation/error-messages).


## [​ ](#faq) FAQ

### [​ ](#1-what-if-file-status-stays-processing-after-upload) 1. What if file status stays "processing" after upload?

Processing typically completes in seconds. If it persists:

- Verify the file format is supported.

- Verify the file size is within limits (file-extract: 150 MB, batch: 500 MB).

- Poll status with the `retrieve` API.


### [​ ](#2-can-file-ids-be-shared-across-accounts) 2. Can file IDs be shared across accounts?

No. File IDs are scoped to the Qwen Cloud account that created them.
### [​ ](#3-are-uploaded-files-stored-permanently) 3. Are uploaded files stored permanently?

Yes. Files persist until you delete them.
### [​ ](#4-why-did-file-upload-fail) 4. Why did file upload fail?


- Invalid or missing API key (`$DASHSCOPE_API_KEY` not set).

- Unsupported file format.

- File size exceeds limit (file-extract: 150 MB, batch: 500 MB).

- Storage limit reached (10,000 files or 100 GB total).

- Rate limit exceeded (3 QPS for upload).


### [​ ](#5-when-should-i-use-file-extract-vs-batch) 5. When should I use `file-extract` vs `batch`?


- `file-extract`: Document analysis and data extraction.

- `batch`: Batch inference tasks (requires JSONL format).


 [Previous ](/api-reference/platform-api/conversations/delete-item)[Upload file Next ](/api-reference/platform-api/file/upload-file)
