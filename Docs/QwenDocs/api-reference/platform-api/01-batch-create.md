# Create batch

> **Source:** https://docs.qwencloud.com/api-reference/platform-api/batch/create-batch

POST/batches Python Python

 Copy ```\nfrom openai import OpenAI
import os

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)

batch = client.batches.create(
 input_file_id="file-abc123",
 endpoint="/v1/chat/completions",
 completion_window="24h"
)

print(f"Batch ID: {batch.id}")
print(f"Status: {batch.status}")
``` 200 400 Copy ```\n{
 "id": "batch_abc123",
 "object": "batch",
 "endpoint": "/v1/chat/completions",
 "errors": {
 "object": "list",
 "data": [
 {
 "code": "&#x3C;string>",
 "message": "&#x3C;string>",
 "param": "&#x3C;string>",
 "line": 0
 }
 ]
 },
 "input_file_id": "file-abc123",
 "completion_window": "24h",
 "status": "validating",
 "output_file_id": "file-xyz789",
 "error_file_id": "file-err456",
 "created_at": 1735113344,
 "in_progress_at": 0,
 "expires_at": 0,
 "finalizing_at": 0,
 "completed_at": 0,
 "failed_at": 0,
 "expired_at": 0,
 "cancelled_at": 0,
 "cancelling_at": 0,
 "request_counts": {
 "total": 0,
 "completed": 0,
 "failed": 0
 },
 "metadata": {
 "ds_name": "&#x3C;string>",
 "ds_description": "&#x3C;string>"
 }
}
``` ### Authorizations
 [​ ](#authorization) Authorizationstring header required Qwen Cloud API Key. Create one in the [console](https://home.qwencloud.com/api-keys).

 ### Body
application/json [​ ](#input-file-id) input_file_idstring required The ID of the uploaded input file. Obtain this from the file upload response.

 Example:file-abc123 [​ ](#endpoint) endpointenum<string> required The API endpoint for batch requests. Must match the `url` field in the input file. Use `/v1/chat/completions` for text generation, `/v1/embeddings` for embeddings, or `/v1/chat/ds-test` for the test model.

 Available options:/v1/chat/completions,/v1/embeddings,/v1/chat/ds-test Example:/v1/chat/completions [​ ](#completion-window) completion_windowstring required Maximum time allowed for the batch to complete. Range: `24h` to `336h` (14 days). Supports hour (`h`) and day (`d`) units, e.g., `24h`, `7d`. Integers only.

 Example:24h [​ ](#metadata) metadataobject | null Optional key-value metadata for the batch job.

 Example: Copy ```\n{
 "ds_name": "Task name",
 "ds_description": "Task description"
}
``` Show child attributes

 [​ ](#metadatads-name) metadata.ds_namestring Task name. Max 100 characters. If defined multiple times, the last value is used.

 Example:nightly evaluation [​ ](#metadatads-description) metadata.ds_descriptionstring Task description. Max 200 characters. If defined multiple times, the last value is used.

 Example:Daily model evaluation batch [​ ](#metadatads-batch-finish-callback) metadata.ds_batch_finish_callbackstring Publicly accessible callback URL. The system sends a POST request with task status on completion.

 Example:https://example.com/callback ### Response
200-application/json [​ ](#id) idstring Unique batch job identifier.

 Example:batch_abc123 [​ ](#object) objectenum<string> Always `"batch"`.

 Available options:batch [​ ](#endpoint) endpointstring The API endpoint used for this batch.

 Example:/v1/chat/completions [​ ](#errors) errorsobject | null Errors encountered during batch processing.

 Show child attributes

 [​ ](#errorsobject) errors.objectenum<string> Always `"list"`.

 Available options:list [​ ](#errorsdata) errors.dataobject[] List of error details.

 Show child attributes

 [​ ](#errorsdatacode) errors.data.codestring Error code.

 [​ ](#errorsdatamessage) errors.data.messagestring Human-readable error message.

 [​ ](#errorsdataparam) errors.data.paramstring | null Parameter that caused the error.

 [​ ](#errorsdataline) errors.data.lineinteger | null Line number in the input file that caused the error.

 [​ ](#input-file-id) input_file_idstring ID of the input file.

 Example:file-abc123 [​ ](#completion-window) completion_windowstring The completion window for the batch job.

 Example:24h [​ ](#status) statusenum<string> Current status of the batch job. `validating`: input file is being validated. `in_progress`: batch is being processed. `finalizing`: results are being compiled. `completed`: all requests finished. `failed`: job failed. `expired`: job exceeded the completion window. `cancelling`: cancellation in progress. `cancelled`: job was cancelled.

 Available options:validating,in_progress,finalizing,completed,failed,expired,cancelling,cancelled [​ ](#output-file-id) output_file_idstring | null ID of the file containing successful results. Available when status is `completed`. Use with the download file content endpoint.

 Example:file-xyz789 [​ ](#error-file-id) error_file_idstring | null ID of the file containing error details. Available when some requests failed. Use with the download file content endpoint.

 Example:file-err456 [​ ](#created-at) created_atinteger Unix timestamp (seconds) when the batch was created.

 Example:1735113344 [​ ](#in-progress-at) in_progress_atinteger | null Unix timestamp (seconds) when the batch started processing.

 [​ ](#expires-at) expires_atinteger | null Unix timestamp (seconds) when the batch will expire.

 [​ ](#finalizing-at) finalizing_atinteger | null Unix timestamp (seconds) when the batch started finalizing.

 [​ ](#completed-at) completed_atinteger | null Unix timestamp (seconds) when the batch completed.

 [​ ](#failed-at) failed_atinteger | null Unix timestamp (seconds) when the batch failed.

 [​ ](#expired-at) expired_atinteger | null Unix timestamp (seconds) when the batch expired.

 [​ ](#cancelled-at) cancelled_atinteger | null Unix timestamp (seconds) when the batch was cancelled.

 [​ ](#cancelling-at) cancelling_atinteger | null Unix timestamp (seconds) when the batch entered cancelling status.

 [​ ](#request-counts) request_countsobject Request processing counts.

 Show child attributes

 [​ ](#request-countstotal) request_counts.totalinteger Total number of requests in the batch.

 [​ ](#request-countscompleted) request_counts.completedinteger Number of requests that completed successfully.

 [​ ](#request-countsfailed) request_counts.failedinteger Number of requests that failed.

 [​ ](#metadata) metadataobject | null Key-value metadata attached to the batch.

 Show child attributes

 [​ ](#metadatads-name) metadata.ds_namestring Task name.

 [​ ](#metadatads-description) metadata.ds_descriptionstring Task description.
