# Batch API

> **Source:** https://docs.qwencloud.com/developer-guides/text-generation/batch

Process bulk requests asynchronously at 50% off

 Copy page Process bulk `qwen-max`, `qwen-plus`, `qwen-flash`, or `qwen-turbo` requests asynchronously at **50% of the real-time price**. Results are delivered within 24 hours. You can create and manage batch jobs using the Qwen Cloud console or the API.
## [​ ](#input-file-format) Input file format

Each line in the JSONL input file is one request:
Copy ```\n{"custom_id":"req-1","method":"POST","url":"/v1/chat/completions","body":{"model":"qwen-plus","messages":[{"role":"user","content":"Summarize quantum computing in two sentences."}]}}
{"custom_id":"req-2","method":"POST","url":"/v1/chat/completions","body":{"model":"qwen-plus","messages":[{"role":"user","content":"What is 2+2?"}]}}

``` 
Set `url` to `/v1/chat/completions` for all requests. Up to **50,000** requests per file, **500 MB** total, **6 MB** per line. All requests must use the same model. Each `custom_id` must be unique.
## [​ ](#use-the-qwen-cloud-console) Use the Qwen Cloud console

Open [Batch API](https://home.qwencloud.com/model-production/batch-api) in the Qwen Cloud console.
### [​ ](#create-a-batch-job) Create a batch job


- Click **Create batch job**.

- Fill in the **Task name** and **Description**, set the **Max wait time** (1–14 days), and upload your JSONL input file.

- Click **Create batch job** to submit.


 Click **Sample File** on the right to download a template JSONL file. 
### [​ ](#monitor-and-manage-tasks) Monitor and manage tasks

On the task list, view each task&#x27;s **progress** (processed / total requests) and **status**. Filter by status to locate a task.
Click **Cancel** to stop a task that is `validating` or `in_progress`. Click **Detail** to view job configuration, statistics, and files.
### [​ ](#download-results) Download results

After the task reaches `completed` status, open the job detail page to download from **Input & Output Files**:

- **Output file**: Successful requests with their responses.

- **Error file** (if any): Failed requests with error details.


Both files include `custom_id` for matching against the original input.
### [​ ](#view-usage) View usage

On the [Pay-As-You-Go](https://home.qwencloud.com/billing/pay-as-you-go) page, view spending by model. Batch usage appears as a line item in the **Spending Trends** table. Data may lag by up to 1–2 hours.
## [​ ](#use-the-api) Use the API

### [​ ](#upload-file) Upload file

Python Node.js curl Copy ```\nimport os
from pathlib import Path
from openai import OpenAI
client = OpenAI(api_key=os.getenv("DASHSCOPE_API_KEY"), base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1")

file_object = client.files.create(
 file=Path("input.jsonl"),
 purpose="batch"
)
print(file_object.id) # <-- use this in the next step

``` 
Response (key fields):
Copy ```\n{"id": "file-batch-xxx", "status": "uploaded", "purpose": "batch"}

``` 
 **Reuse an existing file ID:** The ID returned after uploading a file (e.g., `file-batch-xxx`) can be reused. If the input content remains the same, skip re-uploading and directly create a task with the existing ID:Copy ```\nbatch = client.batches.create(
 input_file_id="file-batch-xxx", # Reuse existing file ID, no need to re-upload
 endpoint="/v1/chat/completions",
 completion_window="24h"
)

``` You can retrieve historical file IDs through the `client.files.list(purpose="batch")` API. 
### [​ ](#create-batch) Create batch

Python Node.js curl Copy ```\nbatch = client.batches.create(
 input_file_id="file-batch-xxx", # <-- from upload step
 endpoint="/v1/chat/completions", # <-- must match url in input file
 completion_window="24h", # <-- 24h to 336h (14 days)
 metadata={
 "ds_name": "My batch job", # <-- optional: task name (max 100 chars)
 "ds_description": "Weekly report", # <-- optional: task description (max 200 chars)
 }
)
print(batch.id)

``` 
 **Dry-run with the test model**: Use model `batch-test-model` with endpoint `/v1/chat/ds-test` to validate your file format without inference costs. Limits: 1 MB file, 100 lines, 2 concurrent tasks. 
### [​ ](#check-status) Check status

Python Node.js curl Copy ```\nbatch = client.batches.retrieve("batch_xxx")
print(batch.status) # <-- see status lifecycle below

``` 
Status lifecycle: `validating` → `in_progress` → `finalizing` → `completed`. Terminal states: `completed`, `failed`, `expired`, `cancelled`. Poll every 1–2 minutes.
Response (key fields):
Copy ```\n{
 "id": "batch_xxx",
 "status": "completed",
 "output_file_id": "file-batch_output-xxx", // <-- download this
 "error_file_id": "file-batch_error-xxx", // <-- failed requests (if any)
 "request_counts": {"total": 100, "completed": 98, "failed": 2}
}

``` 
### [​ ](#download-results-2) Download results

Python Node.js curl Copy ```\ncontent = client.files.content("file-batch_output-xxx") # <-- output_file_id from above
content.write_to_file("result.jsonl")

``` 
Each line in the output JSONL maps to a request by `custom_id`:
Copy ```\n{"id": "batch_req_xxx", "custom_id": "req-1", "response": {"status_code": 200, "body": {"choices": [{"message": {"content": "..."}}], "usage": {...}}}}

``` 
Download the error file (`error_file_id`) the same way to inspect failed requests. See [error codes](/api-reference/preparation/error-messages) for details.
## [​ ](#manage-batches) Manage batches

### [​ ](#list-batches) List batches

Python curl Copy ```\nbatches = client.batches.list(limit=10)

``` 
Filter parameters (query string): `ds_name` (fuzzy match), `input_file_ids` (comma-separated, max 20), `status` (comma-separated), `create_after` / `create_before` (format: `yyyyMMddHHmmss`), `after` (cursor), `limit` (page size).
### [​ ](#cancel-a-batch) Cancel a batch

Python curl Copy ```\nclient.batches.cancel("batch_xxx")

``` 
Status moves to `cancelling`, then `cancelled` after in-flight requests finish. Completed requests before cancellation are still billed.
## [​ ](#utility-scripts) Utility scripts

 CSV to JSONL converter

 Copy ```\nimport csv, json

def build_messages(content):
 return [{"role": "user", "content": content}]

with open("input.csv") as fin, open("input.jsonl", "w") as fout:
 for row in csv.reader(fin):
 request = {
 "custom_id": row[0],
 "method": "POST",
 "url": "/v1/chat/completions",
 "body": {"model": "qwen-plus", "messages": build_messages(row[1])}
 }
 fout.write(json.dumps(request, ensure_ascii=False) + "\n")

``` 
 JSONL results to CSV converter

 Copy ```\nimport json, csv

columns = ["custom_id", "status_code", "content", "usage"]

def get(obj, path):
 for key in path:
 obj = obj[key] if obj and key in obj else None
 return obj

with open("result.jsonl") as fin, open("result.csv", "w") as fout:
 writer = csv.writer(fout)
 writer.writerow(columns)
 for line in fin:
 r = json.loads(line)
 writer.writerow([
 r.get("custom_id"),
 get(r, ["response", "status_code"]),
 get(r, ["response", "body", "choices", 0, "message", "content"]),
 get(r, ["response", "body", "usage"]),
 ])

``` 
## [​ ](#notes) Notes


- **50% discount**: Input and output tokens are billed at half the real-time price. Only successful requests are billed. See [pricing](/developer-guides/getting-started/pricing).

- **Thinking tokens**: Models like `qwen3.6-plus`, `qwen3.5-plus`, and `qwen3.5-flash` enable thinking by default, generating extra tokens at output price. Set `enable_thinking` based on task needs. See [thinking](/developer-guides/text-generation/thinking).

- **`enable_thinking` placement**: In the JSONL request body, `enable_thinking` is a top-level parameter of `body` and must be placed at the same level as `model`. Do not place it inside `extra_body`.

- **Not stackable**: Batch discount does not stack with context cache or other discounts.

- **File storage**: 10,000 files / 100 GB per account. Delete old files to free space.

- **Rate limits**: Create 1,000/min (1,000 concurrent), query 1,000/min, list 100/min, cancel 1,000/min.

- **Task retention**: Only tasks from the last 30 days are queryable via list.


## [​ ](#api-reference) API reference


- [Create batch](/api-reference/platform-api/batch/create-batch)

- [Retrieve batch](/api-reference/platform-api/batch/retrieve-batch)

- [List batches](/api-reference/platform-api/batch/list-batches)

- [Cancel batch](/api-reference/platform-api/batch/cancel-batch)


 [Previous ](/developer-guides/text-generation/thinking)[Multi-turn conversations Manage chat context Next ](/developer-guides/text-generation/multi-turn)
