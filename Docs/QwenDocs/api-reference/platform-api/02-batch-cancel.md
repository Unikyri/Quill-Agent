# Cancel batch

> **Source:** https://docs.qwencloud.com/api-reference/platform-api/batch/cancel-batch

POST/batches/{batch_id}/cancel Python Python

 Copy ```\nfrom openai import OpenAI
import os

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
)

batch = client.batches.cancel("batch_abc123")

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

 ### Path Parameters
 [​ ](#batch-id) batch_idstring required The ID of the batch job to cancel.

 Example:batch_abc123 ### Response
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
