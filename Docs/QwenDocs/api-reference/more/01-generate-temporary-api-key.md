# Generate a temporary API key

> **Source:** https://docs.qwencloud.com/api-reference/more/generate-a-temporary-api-key

Short-lived access tokens

 Copy page Use a secure backend to provide temporary API keys when your application calls model services from untrusted environments like browsers or mobile apps. This prevents exposing your permanent API key.
 A temporary API key inherits the permissions of the API key that created it, such as model access restrictions. 
## [​ ](#prerequisites) Prerequisites

Create a permanent API key in [Key Management](https://home.qwencloud.com/api-keys) and set the `DASHSCOPE_API_KEY` environment variable. See [Configure your API key](/api-reference/preparation/export-api-key-env).
## [​ ](#request-example) Request example

Temporary API keys expire after 60 seconds by default. You can set `expire_in_seconds` from 1 to 1,800.
Copy ```\ncurl -X POST "https://dashscope-intl.aliyuncs.com/api/v1/tokens?expire_in_seconds=1800" \
-H "Authorization: Bearer $DASHSCOPE_API_KEY"

``` 
## [​ ](#sample-response) Sample response

### [​ ](#success-response) Success response

Copy ```\n{
 "token": "st-****",
 "expires_at": 1744080369
}

``` 
### [​ ](#response-parameters) Response parameters

ParameterTypeDescriptionExampletokenStringThe temporary API key.st-****expires_atNumberExpiration time as a UNIX timestamp in seconds.1744080369 
### [​ ](#error-response) Error response

Copy ```\n{
 "code": "InvalidApiKey",
 "message": "Invalid API-key provided.",
 "request_id": "902fee3b-f7f0-9a8c-96a1-6b4ea25af114"
}

``` 
### [​ ](#response-parameters-2) Response parameters

ParameterTypeDescriptionExamplecodeStringError code. See Error messages for details.InvalidApiKeymessageStringError message.Invalid API-key provided.request_idStringThe request ID.902fee3b-f7f0-9a8c-96a1-6b4ea25af114 
## [​ ](#faq) FAQ

### [​ ](#can-i-manually-delete-a-temporary-api-key) Can I manually delete a temporary API key?

No. Temporary API keys expire automatically and cannot be deleted manually. [Previous ](/api-reference/toolkitframework/openai-compatible/overview)[Manage asynchronous tasks Track and control async jobs Next ](/api-reference/more/manage-asynchronous-tasks)
