# Postman

> **Source:** https://docs.qwencloud.com/developer-guides/clients-and-developer-tools/postman

API testing tool

 Copy page Postman is a graphical HTTP testing tool that makes it easy to test Qwen Cloud APIs. Use it to quickly validate API endpoints, test async operations for image/video generation, and prototype integrations before writing code.
## [​ ](#quick-start) Quick start

Get running in a few minutes:
Copy ```\n# 1. Install
Download Postman from postman.com/downloads

# 2. Create request (New → HTTP Request)
Method: POST
URL: https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/text-generation/generation

# 3. Configure (Headers tab)
Authorization: Bearer sk-xxx
Content-Type: application/json

# 4. Test (Body tab → raw → JSON)
{
 "model": "qwen-plus",
 "input": {
 "messages": [{"role": "user", "content": "Hello, who are you?"}]
 }
}

``` 
You should see: JSON response with the model&#x27;s reply
## [​ ](#configuration) Configuration

### [​ ](#basic-setup) Basic setup

Configure Postman for Qwen Cloud APIs:

- API endpoints: `https://dashscope-intl.aliyuncs.com`

- Authentication: Bearer token with API key

- Content-Type: `application/json`


 **Free quota and billing:**
- First-time users get a free quota. See [Free quota](/resources/free-quota) for details.

- Enable [Free quota only](https://home.qwencloud.com/benefits) to prevent unexpected charges.


 
### [​ ](#api-types) API types

Qwen Cloud offers two API patterns:
TypeUse forResponse pattern**Synchronous**Text generation, embeddingsImmediate response**Asynchronous**Image/video generationTask ID → Poll for result 
## [​ ](#synchronous-apis) Synchronous APIs

### [​ ](#text-generation-example) Text generation example

 1 Create request

New → HTTP Request → POST 2 Set URL

Copy ```\nhttps://dashscope-intl.aliyuncs.com/api/v1/services/aigc/text-generation/generation

``` 3 Add headers

KeyValueAuthorizationBearer YOUR_API_KEYContent-Typeapplication/json 4 Add body

Copy ```\n{
 "model": "qwen-plus",
 "input": {
 "messages": [
 {"role": "user", "content": "Write a haiku about coding"}
 ]
 },
 "parameters": {
 "temperature": 0.7
 }
}

``` 5 Send request

Click **Send** → View response 
## [​ ](#asynchronous-apis) Asynchronous APIs

For time-consuming tasks (images, videos), use the async pattern:
### [​ ](#step-1-create-task) Step 1: Create task

 1 Configure request

Method: **POST**
URL: `https://dashscope-intl.aliyuncs.com/api/v1/services/aigc/text2image/image-synthesis` 2 Add headers

KeyValueX-DashScope-AsyncenableAuthorizationBearer YOUR_API_KEYContent-Typeapplication/json 3 Add body

Copy ```\n{
 "model": "wan2.5-t2i-preview",
 "input": {
 "prompt": "A serene mountain landscape at sunset"
 },
 "parameters": {
 "size": "1024*1024",
 "n": 1
 }
}

``` 4 Send and save task_id

Click **Send**. The response contains `output.task_id` - save this value. The `task_id` is valid for 24 hours. Retrieve the result before it expires.Copy ```\n{
 "request_id": "896b2ccd-a0cd-40a8-a557-bb73cee5cf95",
 "output": {
 "task_id": "42442de9-917d-4c41-80a7-37fb7ad25ed2",
 "task_status": "PENDING"
 }
}

``` 
### [​ ](#step-2-query-result) Step 2: Query result

 1 Configure query

Method: **GET**
URL: `https://dashscope-intl.aliyuncs.com/api/v1/tasks/{task_id}`
Replace `{task_id}` with actual ID 2 Add headers

KeyValueAuthorizationBearer YOUR_API_KEY 3 Poll for completion

Send request every 3 to 5 seconds until `output.task_status` is `SUCCEEDED`. The response then contains the image URL:Copy ```\n{
 "output": {
 "task_status": "SUCCEEDED",
 "results": [
 {
 "orig_prompt": "A serene mountain landscape at sunset",
 "url": "https://dashscope-result-wlcb.oss-cn-wulanchabu.aliyuncs.com/..."
 }
 ]
 }
}

``` The `task_id` is valid for 24 hours. Download generated files before they expire. 
## [​ ](#curl-to-postman-mapping) cURL to Postman mapping

Converting cURL examples to Postman:
cURLPostmanLocation`curl -X POST`POSTMethod dropdownURLURLURL field`-H &#x27;Key: Value&#x27;`HeadersHeaders tab`-d &#x27;{...}&#x27;`BodyBody tab (raw JSON)`$VARIABLE`{{variable}}Environment variables 
## [​ ](#environment-variables) Environment variables

Set up reusable variables:
 1 Create environment

Environments → Create New → Name it "Qwen Cloud" 2 Add variables

VariableValueapi_keyYOUR_API_KEYbase_url[https://dashscope-intl.aliyuncs.com](https://dashscope-intl.aliyuncs.com)modelqwen-plus 3 Use in requests


- Headers: `Bearer {{api_key}}`

- URL: `{{base_url}}/api/v1/...`

- Body: `"model": "{{model}}"`


 
## [​ ](#collections) Collections

Organize related requests:

- **Create collection**: Collections → New Collection

- **Add requests**: Drag requests into collection

- **Share**: Export as JSON or share link

- **Run all**: Runner → Select collection → Run


## [​ ](#testing-tips) Testing tips

### [​ ](#response-validation) Response validation

Add tests in **Tests** tab:
Copy ```\npm.test("Status is 200", () => {
 pm.response.to.have.status(200);
});

pm.test("Has output", () => {
 const json = pm.response.json();
 pm.expect(json).to.have.property("output");
});

``` 
### [​ ](#async-polling-automation) Async polling automation

Automate task polling with scripts:
Copy ```\n// In Tests tab of create task request
const taskId = pm.response.json().output.task_id;
pm.environment.set("task_id", taskId);

// Set next request to query task
postman.setNextRequest("Query Task Status");

``` 
## [​ ](#troubleshooting) Troubleshooting

**401 Unauthorized**

Solution:

- Check API key is correct

- Verify "Bearer " prefix in Authorization header

- Ensure API key has quota


**400 Bad Request**

Solution:

- Validate JSON syntax in body

- Check required fields are present

- Verify model name is correct


**Task stuck in PENDING**

Solution:

- Image/video generation can take minutes

- Continue polling every 3-5 seconds

- Check task_metrics for progress


**Connection timeout**

Solution:

- Increase timeout in Settings → General

- Check network connectivity

- Try a simpler request first


## [​ ](#production-notes) Production notes

 Postman is for testing only. In production:
- Use official SDKs for your language

- Implement proper error handling

- Add retry logic for async tasks

- Store API keys securely


 
## [​ ](#related-resources) Related resources


- **API Reference**: [Complete API documentation →](/api-reference/chat/dashscope)

- **Models**: [Available models →](/developer-guides/getting-started/text-generation-models) | [Pricing →](/developer-guides/getting-started/pricing)

- **SDKs**: [Official client libraries →](/api-reference/preparation/install-sdk)

- **Postman docs**: [Official Postman guide →](https://learning.postman.com/docs)


 [Previous ](/developer-guides/clients-and-developer-tools/kilo-cli)[Dify Low-code LLM app platform Next ](/developer-guides/clients-and-developer-tools/dify)
