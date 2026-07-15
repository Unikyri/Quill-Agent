# Async task management

> **Source:** https://docs.qwencloud.com/developer-guides/run-and-scale/async-task-management

Two async patterns on Qwen Cloud: task-based for media generation, batch for high-volume text processing

 Copy page Qwen Cloud offers two async processing systems. Choose based on what you&#x27;re building:
Task APIBatch API**Use case**Image generation, video generation, file transcriptionHigh-volume text generation, embeddings**How it works**Submit one request → get a task ID → poll for resultUpload a JSONL file of requests → poll for completion → download results**Endpoint**`dashscope-intl.aliyuncs.com/api/v1/tasks/``dashscope-intl.aliyuncs.com/compatible-mode/v1/batches`**Cost**Standard pricing**50% discount****Completion time**Seconds to minutesUp to 24h (configurable to 336h)**Result retention**24 hours30 days**Max per request**1 task (with optional multi-output like `n=4`)50,000 requests per file, 500 MB 
## [​ ](#task-api) Task API

Image generation, video generation, and file transcription use the Task API: you submit a single request, receive a task ID, and poll for the result.
### [​ ](#how-it-works) How it works

Every task follows the same two-step pattern:

- **Create a task** — Send a POST request and receive a `task_id`

- **Query the result** — Poll `GET /api/v1/tasks/{task_id}` until the task reaches a terminal state


Copy ```\nPOST (create) ──→ task_id
 │
 ┌─────────────┘
 ▼
 ┌─────────┐ ┌─────────┐ ┌───────────┐
 │ PENDING │───→│ RUNNING │───→│ SUCCEEDED │
 └─────────┘ └─────────┘ └───────────┘
 │ │
 ▼ ▼
 ┌──────────┐ ┌────────┐
 │ CANCELED │ │ FAILED │
 └──────────┘ └────────┘

``` 
#### [​ ](#task-states) Task states

StateMeaning`PENDING`Queued, waiting to start`RUNNING`Actively processing`SUCCEEDED`Completed successfully`FAILED`Encountered an error`CANCELED`Manually canceled (only from `PENDING` state)`UNKNOWN`Task ID expired or cannot be queried 
### [​ ](#create-a-task) Create a task

Some APIs (like image generation) support both synchronous and asynchronous modes. Add the `X-DashScope-Async` header to force async execution:
Copy ```\ncurl -X POST https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/image-generation/generation \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY" \
 -H "Content-Type: application/json" \
 -H "X-DashScope-Async: enable" \
 -d &#x27;{
 "model": "wan2.6-t2i",
 "input": {
 "messages": [
 {
 "role": "user",
 "content": [
 {"text": "A cat sitting on a windowsill at sunset"}
 ]
 }
 ]
 }
 }&#x27;

``` 
The response returns a `task_id`:
Copy ```\n{
 "request_id": "e5d5af82-9a08-xxx",
 "output": {
 "task_id": "86ecf553-d340-xxx",
 "task_status": "PENDING"
 }
}

``` 
 Video generation APIs are always asynchronous — you do not need the `X-DashScope-Async` header for these. 
### [​ ](#query-the-result) Query the result

Poll the task endpoint with the returned `task_id`:
- curl 
- Python 

 Copy ```\ncurl -X GET "https://dashscope-intl.aliyuncs.com/api/v1/tasks/86ecf553-d340-xxx" \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY"

``` Copy ```\nimport requests
import os

task_id = "86ecf553-d340-xxx"
url = f"https://dashscope-intl.aliyuncs.com/api/v1/tasks/{task_id}"
headers = {"Authorization": f"Bearer {os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;)}"}

response = requests.get(url, headers=headers)
result = response.json()
print(result["output"]["task_status"])

``` 
A completed task returns the result in the `output` field:
Copy ```\n{
 "request_id": "e5d5af82-9a08-xxx",
 "output": {
 "task_id": "86ecf553-d340-xxx",
 "task_status": "SUCCEEDED",
 "submit_time": "2025-01-15 10:30:00.000",
 "scheduled_time": "2025-01-15 10:30:01.000",
 "end_time": "2025-01-15 10:30:12.000",
 "finished": true,
 "choices": [
 {
 "finish_reason": "stop",
 "message": {
 "role": "assistant",
 "content": [
 {
 "image": "https://dashscope-result-xxx.oss-xxx.aliyuncs.com/...",
 "type": "image"
 }
 ]
 }
 }
 ]
 },
 "usage": {
 "image_count": 1,
 "input_tokens": 0,
 "output_tokens": 0
 }
}

``` 
 Older models (Wan 2.5 and earlier) return a different response structure with `results` and `task_metrics` fields instead of `choices`. See the [text-to-image API reference](/api-reference/image-generation/wan-text-to-image-v2/query-result) for details on both formats. 
A failed task includes an error code and message:
Copy ```\n{
 "output": {
 "task_id": "86ecf553-d340-xxx",
 "task_status": "FAILED",
 "code": "DataInspectionFailed",
 "message": "Input or output data may contain inappropriate content."
 }
}

``` 
### [​ ](#polling-strategies) Polling strategies

Naive fixed-interval polling wastes API calls and may hit rate limits. Use exponential backoff instead.
Start with a short interval and increase gradually. Different modalities have different typical completion times:
ModalityInitial intervalIncrease factorTimeoutImage generation3 seconds1.5×2 minutesVideo generation15 seconds1.5×5 minutesFile transcription (ASR)5 seconds1.5×Varies with audio length 
- Python 
- Node.js 

 Copy ```\nimport time
import requests
import os

def poll_task(task_id, initial_interval=3, max_interval=15, timeout=120):
 """Poll a task with exponential backoff."""
 url = f"https://dashscope-intl.aliyuncs.com/api/v1/tasks/{task_id}"
 headers = {"Authorization": f"Bearer {os.getenv(&#x27;DASHSCOPE_API_KEY&#x27;)}"}

 interval = initial_interval
 elapsed = 0

 while elapsed < timeout:
 response = requests.get(url, headers=headers)
 result = response.json()
 status = result["output"]["task_status"]

 if status == "SUCCEEDED":
 return result
 elif status == "FAILED":
 raise Exception(
 f"Task failed: {result[&#x27;output&#x27;].get(&#x27;code&#x27;, &#x27;Unknown&#x27;)} - "
 f"{result[&#x27;output&#x27;].get(&#x27;message&#x27;, &#x27;&#x27;)}"
 )

 time.sleep(interval)
 elapsed += interval
 interval = min(interval * 1.5, max_interval)

 raise TimeoutError(f"Task {task_id} did not complete within {timeout}s")

# Image generation: poll every 3s, escalate to 15s, timeout at 2min
result = poll_task("86ecf553-d340-xxx", initial_interval=3, max_interval=15, timeout=120)

# Video generation: poll every 15s, escalate to 60s, timeout at 5min
result = poll_task("86ecf553-d340-xxx", initial_interval=15, max_interval=60, timeout=300)

``` Copy ```\nasync function pollTask(taskId, { initialInterval = 3000, maxInterval = 15000, timeout = 120000 } = {}) {
 const url = `https://dashscope-intl.aliyuncs.com/api/v1/tasks/${taskId}`;
 const headers = { Authorization: `Bearer ${process.env.DASHSCOPE_API_KEY}` };

 let interval = initialInterval;
 const deadline = Date.now() + timeout;

 while (Date.now() < deadline) {
 const res = await fetch(url, { headers });
 const result = await res.json();
 const { task_status, code, message } = result.output;

 if (task_status === "SUCCEEDED") return result;
 if (task_status === "FAILED") throw new Error(`Task failed: ${code} - ${message}`);

 await new Promise(resolve => setTimeout(resolve, interval));
 interval = Math.min(interval * 1.5, maxInterval);
 }

 throw new Error(`Task ${taskId} did not complete within ${timeout}ms`);
}

``` 
### [​ ](#sdk-level-abstraction) SDK-level abstraction

The DashScope Python SDK encapsulates the polling loop. Use `call()` for a synchronous experience, or `async_call()` + `wait()` for explicit control.
Copy ```\nimport dashscope
from dashscope import ImageSynthesis

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

# Option 1: Synchronous — SDK polls internally
result = ImageSynthesis.call(
 model="wan2.1-t2i-plus",
 prompt="A cat sitting on a windowsill at sunset",
 n=1
)
print(result.output.results[0].url)

# Option 2: Explicit async — you control when to wait
task = ImageSynthesis.async_call(
 model="wan2.1-t2i-plus",
 prompt="A cat sitting on a windowsill at sunset",
 n=1
)
# Do other work here...
result = ImageSynthesis.wait(task)
print(result.output.results[0].url)

``` 
The same pattern applies across modalities:
Copy ```\nfrom dashscope import VideoSynthesis

task = VideoSynthesis.async_call(model="wan2.1-t2v-plus", prompt="...")
result = VideoSynthesis.wait(task)

``` 
### [​ ](#manage-tasks-at-scale) Manage tasks at scale

#### [​ ](#list-tasks) List tasks

Query tasks by time range, status, or model:
Copy ```\ncurl -X GET "https://dashscope-intl.aliyuncs.com/api/v1/tasks/?start_time=20250115100000&end_time=20250115120000&status=FAILED&model_name=wan2.6-t2i&page_no=1&page_size=20" \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY"

``` 
Supported filters:
ParameterDescription`start_time`Start of time range, format: `YYYYMMDDhhmmss``end_time`End of time range, format: `YYYYMMDDhhmmss``status`Filter by task status (`PENDING`, `RUNNING`, `SUCCEEDED`, `FAILED`)`model_name`Filter by model name`page_no`Page number (starting from 1)`page_size`Results per page 
 The time range cannot exceed 24 hours. If both `start_time` and `end_time` are omitted, the API returns tasks from the last 24 hours. 
#### [​ ](#cancel-a-task) Cancel a task

Cancel a task that is still in `PENDING` state:
Copy ```\ncurl -X POST "https://dashscope-intl.aliyuncs.com/api/v1/tasks/86ecf553-d340-xxx/cancel" \
 -H "Authorization: Bearer $DASHSCOPE_API_KEY"

``` 
 Only tasks in `PENDING` state can be canceled. Once a task moves to `RUNNING`, it will continue to completion. 
#### [​ ](#rate-limits) Rate limits

Task management endpoints (query, list, cancel) share a rate limit of **20 QPS per account**. If you manage many concurrent tasks, batch your status checks rather than polling each task individually at high frequency.
### [​ ](#task-api-best-practices) Task API best practices

**Download results immediately** — Task data and result URLs (images, videos) are retained for **24 hours** only. After that, both the task metadata and generated files are automatically deleted. Always download and save results to persistent storage as soon as the task succeeds.
**Handle multi-output tasks** — When you request multiple outputs (such as `n=4` images), the task is marked `SUCCEEDED` if **at least one** output is generated. Iterate over the `choices` array and check `finish_reason` to handle partial success:
Copy ```\nresult = poll_task(task_id)
choices = result["output"]["choices"]
print(f"{len(choices)} outputs generated")

for choice in choices:
 if choice["finish_reason"] == "stop":
 for item in choice["message"]["content"]:
 if "image" in item:
 download(item["image"])
 else:
 print(f"Output skipped: {choice[&#x27;finish_reason&#x27;]}")

``` 
**Avoid duplicate tasks** — Each POST request creates a new task, even with identical parameters. If your application retries a failed submission, track task IDs to avoid duplicate work and unnecessary cost.
## [​ ](#batch-api) Batch API

The Batch API processes large volumes of text generation or embedding requests at **50% of the standard price**. You upload a JSONL file of requests, wait for processing to complete, and download the results.
FeatureDetails**Supported models**Qwen text generation models, text embedding models**Input format**JSONL file, up to 50,000 requests / 500 MB**Completion window**24 hours (configurable up to 336 hours)**Result retention**30 days**Cancellation**Supported at any stage (in-progress requests finish first) 
The Batch API uses the OpenAI-compatible endpoint and SDK. For the complete guide with examples, input format, task states, and management operations, see [Batch API](/developer-guides/text-generation/batch).
## [​ ](#common-errors) Common errors

These errors apply to both the Task API and the Batch API:
ErrorCauseAction`DataInspectionFailed`Input or output blocked by content moderationModify the prompt and retry. See [Safety](/developer-guides/run-and-scale/safety)`Throttling.RateQuota`Rate limit exceededReduce request frequency; increase polling intervalTask/batch stuckProcessing taking longer than expectedFor tasks: continue polling until timeout. For batches: increase `completion_window` 
## [​ ](#next-steps) Next steps


- [Manage asynchronous tasks](/api-reference/more/manage-asynchronous-tasks) — Task API reference (query, list, cancel)

- [Batch API](/developer-guides/text-generation/batch) — Full Batch API guide with all languages

- [Text-to-image](/developer-guides/image-generation/text-to-image) — Image generation with async task examples

- [Text-to-video](/developer-guides/video-generation/text-to-video) — Video generation with async task examples

- [Rate limits](/developer-guides/administration/rate-limits) — Account-level rate limit details


 [Previous ](/developer-guides/text-generation/streaming)[Connection reuse and pooling HTTP connection reuse and WebSocket connection pooling for high-concurrency workloads. Next ](/developer-guides/run-and-scale/connection-pooling)
