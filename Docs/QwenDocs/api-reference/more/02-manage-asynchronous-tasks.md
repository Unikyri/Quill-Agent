# Manage asynchronous tasks

> **Source:** https://docs.qwencloud.com/api-reference/more/manage-asynchronous-tasks

Track and control async jobs

 Copy page Some models (like image and video generation) run asynchronously: you create a task, get an ID, then query the result with that ID. Use these APIs to query results, check status in batches, and cancel queued tasks.
## [​ ](#prerequisites) Prerequisites

Call these APIs over HTTP.
[Get an API key](/api-reference/preparation/api-key), then [set it as an environment variable](/api-reference/preparation/export-api-key-env).
## [​ ](#query-the-result-of-an-asynchronous-task) Query the result of an asynchronous task

**API description**: Query a task&#x27;s status and result by `task_id`.
**Rate limit**: 20 QPS per Qwen Cloud account.
 
- 
You can query all tasks under the Qwen Cloud account that owns the API key, including tasks from any API key under that account. You cannot query tasks from other accounts.


- 
Completed task data is retained for 24 hours by default (check the specific API reference for the exact duration) before automatic deletion.


 
### [​ ](#request-endpoint) Request endpoint

Copy ```\nGET https://dashscope-intl.aliyuncs.com/api/v1/tasks/{task_id}

``` 
### [​ ](#request-parameters) Request parameters

Parameter passingFieldTypeRequiredDescriptionExampleHeaderAuthorizationStringYesThe API key, in the format Bearer sk-xxx.Bearer sk-xxxPathtask_idStringYesThe task ID to query.a8532587-xxxx-xxxx-xxxx-0c46b17950d1 
### [​ ](#response-parameters) Response parameters

FieldTypeDescriptionExamplerequest_idStringThe unique ID for this request.7574ee8f-xxxx-xxxx-xxxx-11c33ab46e51outputObjectContains the result on success or the error code and message on failure. For multi-subtask requests, may include both results and errors.-output.task_idStringThe queried task ID.a8532587-xxxx-xxxx-xxxx-0c46b17950d1output.task_statusStringThe task status. A multi-subtask job is marked SUCCEEDED if at least one subtask succeeds. Failed subtasks show their errors in the output.PENDING, RUNNING, SUCCEEDED, FAILED, UNKNOWNoutput.submit_timeStringThe time the task was submitted.2023-12-20 21:36:31.896output.scheduled_timeStringThe time the task started running.2023-12-20 21:36:39.009output.end_timeStringThe time the task ended.2023-12-20 21:36:45.913output.codeStringThe error code. Returned only on failure.-output.messageStringThe error message. Returned only on failure.-output.task_metricsObjectSubtask status statistics.`{ "TOTAL": 4, "SUCCEEDED": 3, "FAILED": 1 }`usageObjectBilling information for this request. Varies by task type.`"usage": {"image_count": 1}` 
### [​ ](#request-example) Request example

Copy ```\ncurl -X GET &#x27;https://dashscope-intl.aliyuncs.com/api/v1/tasks/73205176-xxxx-xxxx-xxxx-16bd5d902219&#x27; \
--header "Authorization: Bearer $DASHSCOPE_API_KEY"

``` 
 If you haven&#x27;t set the API key as an environment variable, replace `$DASHSCOPE_API_KEY` with your actual key. Example: `--header "Authorization: Bearer sk-xxx"`. 
### [​ ](#response-example) Response example

Copy ```\n{
 "request_id": "45ac7f13-xxxx-xxxx-xxxx-e03c35068d83",
 "output": {
 "task_id": "73205176-xxxx-xxxx-xxxx-16bd5d902219",
 "task_status": "SUCCEEDED",
 "submit_time": "2023-12-20 21:36:31.896",
 "scheduled_time": "2023-12-20 21:36:39.009",
 "end_time": "2023-12-20 21:36:45.913",
 "results": [
 {
 "url": "https://dashscope-result-bj.oss-cn-beijing.aliyuncs.com/xxx1.png"
 },
 {
 "url": "https://dashscope-result-bj.oss-cn-beijing.aliyuncs.com/xxx2.png"
 },
 {
 "url": "https://dashscope-result-bj.oss-cn-beijing.aliyuncs.com/xxx3.png"
 },
 {
 "code": "DataInspectionFailed",
 "message": "Output data may contain inappropriate content."
 }
 ],
 "task_metrics": {
 "TOTAL": 4,
 "SUCCEEDED": 3,
 "FAILED": 1
 }
 },
 "usage": {
 "image_count": 3
 }
}

``` 
## [​ ](#query-the-status-of-multiple-asynchronous-tasks) Query the status of multiple asynchronous tasks

**API description**: Query the status of multiple tasks at once by time range, model, or status.
**Rate limit**: 20 QPS per Qwen Cloud account.
 
- 
You can query all tasks under the Qwen Cloud account that owns the API key, including tasks from any API key under that account. You cannot query tasks from other accounts.


- 
Completed task data is retained for 24 hours by default (check the specific API reference for the exact duration) before automatic deletion.


 
### [​ ](#request-endpoint-2) Request endpoint

Copy ```\nGET https://dashscope-intl.aliyuncs.com/api/v1/tasks/

``` 
### [​ ](#request-parameters-2) Request parameters

Parameter passingFieldTypeRequiredDescriptionExampleHeaderAuthorizationStringYesThe API key, in the format Bearer sk-xxx.Bearer sk-xxxParamstask_idStringNoA specific task ID. If set, only that task is returned. Omit to query multiple tasks.a8532587-xxxx-xxxx-xxxx-0c46b17950d1Paramsstart_timeStringNoStart time (format: YYYYMMDDhhmmss). Defaults to 24 hours before end_time, or the last 24 hours if end_time is also omitted. Range cannot exceed 24 hours.20230420193058 represents 19:30:58 on April 20, 2023.Paramsend_timeStringNoEnd time (format: YYYYMMDDhhmmss). Defaults to 24 hours after start_time. Range cannot exceed 24 hours.20230420193058 represents 19:30:58 on April 20, 2023.Paramsmodel_nameStringNoThe model name.wan2.6-t2vParamsstatusStringNoThe task status: PENDING, RUNNING, SUCCEEDED, FAILED, CANCELED, UNKNOWNParamspage_noIntegerNoThe page number. Default: 1.-Paramspage_sizeIntegerNoThe number of entries per page. Default: 10.- 
### [​ ](#response-parameters-2) Response parameters

FieldTypeDescriptionExamplerequest_idStringThe unique ID for this request.7574ee8f-xxxx-xxxx-xxxx-11c33ab46e51dataArrayA list of query results.See example belowdata[].api_key_idStringThe API key ID.See example belowdata[].caller_parent_idStringThe Qwen Cloud account ID.See example belowdata[].caller_uidStringThe Qwen Cloud account UID.See example belowdata[].gmt_createLongThe task creation time in milliseconds since the epoch.See example belowdata[].start_timeLongThe task start time in milliseconds since the epoch.See example belowdata[].end_timeLongThe task end time in milliseconds since the epoch.See example belowdata[].regionStringThe region.See example belowdata[].request_idStringThe request ID for the task submission.See example belowdata[].statusStringThe task status: PENDING, RUNNING, SUCCEEDED, FAILED, CANCELED, UNKNOWNSee example belowdata[].task_idStringThe task ID.See example belowdata[].user_api_unique_keyStringA unique index derived from the model&#x27;s API parameters at submission.See example belowdata[].model_nameStringThe model name.See example belowpage_noIntegerThe current page number.`"page_no": 1`page_sizeIntegerThe number of entries per page.`"page_size": 10`total_pageIntegerThe total number of pages.`"total_page": 4`totalIntegerThe total number of entries.`"total": 39`codeStringThe error code. Returned only on failure.`"code": "Throttling.RateQuota"`messageStringThe error message. Returned only on failure.`"message": "Requests rate limit exceeded, please try again later."` 
### [​ ](#request-example-2) Request example

Copy ```\ncurl -X GET &#x27;https://dashscope-intl.aliyuncs.com/api/v1/tasks/?start_time=xxx&end_time=xxx&status=xxx&#x27; \
--header "Authorization: Bearer $DASHSCOPE_API_KEY"

``` 
### [​ ](#response-example-2) Response example

Copy ```\n{
 "total": 2,
 "data": [
 {
 "api_key_id": "15xxxx",
 "caller_parent_id": "xxxxxxxxx",
 "caller_uid": "xxxxxxxxx",
 "gmt_create": 1745568428109,
 "model_name": "wanx2.1-kf2v-plus",
 "region": "ap-southeast-1",
 "request_id": "1abfc3c8-dd25-98da-ad0b-xxxxxx",
 "start_time": 1745568428138,
 "status": "RUNNING",
 "task_id": "50e2ccea-abc4-43d7-a0dc-xxxxxx",
 "user_api_unique_key": "apikey:v1:aigc:image2video:video-synthesis:wanx2.1-kf2v-plus"
 },
 {
 "api_key_id": "15xxxx",
 "caller_parent_id": "xxxxxxxxx",
 "caller_uid": "xxxxxxxxx",
 "end_time": 1745568302481,
 "gmt_create": 1745568293253,
 "model_name": "wanx2.1-t2i-turbo",
 "region": "ap-southeast-1",
 "request_id": "f6bf34d9-bf87-9e8b-9ed4-xxxxxx",
 "start_time": 1745568293273,
 "status": "SUCCEEDED",
 "task_id": "3c777dbc-8cc6-4d80-aa90-xxxxxx",
 "user_api_unique_key": "apikey:v1:aigc:text2image:image-synthesis:wanx2.1-t2i-turbo"
 }
 ],
 "total_page": 1,
 "page_no": 1,
 "request_id": "f6756b7e-d0bb-9b74-813a-xxxxxx",
 "page_size": 10
}

``` 
## [​ ](#cancel-an-asynchronous-task) Cancel an asynchronous task

**API description**: Cancel a task in `PENDING` state. Tasks in other states cannot be canceled.
**Rate limit**: 20 QPS per Qwen Cloud account.
 
- You can cancel any task under the Qwen Cloud account that owns the API key, including tasks from any API key under that account. You cannot cancel tasks from other accounts.


 
### [​ ](#request-endpoint-3) Request endpoint

Copy ```\nPOST https://dashscope-intl.aliyuncs.com/api/v1/tasks/{task_id}/cancel

``` 
### [​ ](#request-parameters-3) Request parameters

Parameter passingFieldTypeRequiredDescriptionExampleHeaderAuthorizationStringYesThe API key, in the format Bearer sk-xxx.Bearer sk-xxxPathtask_idStringYesThe task ID to cancel.a8532587-xxxx-xxxx-xxxx-0c46b17950d1 
### [​ ](#response-parameters-3) Response parameters

FieldTypeDescriptionExamplerequest_idStringThe unique ID for this request.7574ee8f-xxxx-xxxx-xxxx-11c33ab46e51codeStringThe error code. Returned only on failure.`"code": "Throttling.RateQuota"`messageStringThe error message. Returned only on failure.`"message": "Requests rate limit exceeded, please try again later."` 
### [​ ](#request-example-3) Request example

Copy ```\ncurl -X POST &#x27;https://dashscope-intl.aliyuncs.com/api/v1/tasks/73205176-xxxx-xxxx-xxxx-16bd5d902219/cancel&#x27; \
--header "Authorization: Bearer $DASHSCOPE_API_KEY"

``` 
### [​ ](#response-example-3) Response example

Copy ```\n{
 "request_id": "45ac7f13-xxxx-xxxx-xxxx-e03c35068d83"
}

``` [Previous ](/api-reference/more/generate-a-temporary-api-key)Next
