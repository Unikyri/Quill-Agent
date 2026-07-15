# Conversations

> **Source:** https://docs.qwencloud.com/api-reference/platform-api/conversations

Auto-managed chat history

 Copy page The Conversations API automatically manages multi-turn context across devices and sessions. Use it with the Responses API to inject historical context without manual synchronization.
For usage examples and multi-turn conversation patterns, see [Multi-turn conversations](/developer-guides/text-generation/multi-turn#using-conversations).
## [​ ](#service-endpoints) Service endpoints

`base_url` for the SDK: `https://dashscope-intl.aliyuncs.com/api/v2/apps/protocols/compatible-mode/v1`
HTTP base endpoint: `https://dashscope-intl.aliyuncs.com/api/v2/apps/protocols/compatible-mode/v1/conversations`
## [​ ](#endpoints) Endpoints

EndpointDescription[Create conversation](/api-reference/platform-api/conversations/create-conversation)Create a conversation with optional initial message items[Retrieve conversation](/api-reference/platform-api/conversations/retrieve-conversation)Retrieve a conversation by ID[Update conversation](/api-reference/platform-api/conversations/update-conversation)Update a conversation&#x27;s metadata[Delete conversation](/api-reference/platform-api/conversations/delete-conversation)Delete a conversation[Create items](/api-reference/platform-api/conversations/create-items)Add message items to a conversation[List items](/api-reference/platform-api/conversations/list-items)List message items in a conversation[Retrieve item](/api-reference/platform-api/conversations/retrieve-item)Retrieve a message item by ID[Delete item](/api-reference/platform-api/conversations/delete-item)Delete a message item 
## [​ ](#key-limits) Key limits


- Maximum 20 items per `items` array in create and add operations.

- Maximum 16 metadata key-value pairs (keys: max 64 chars, values: max 512 chars).


 [Previous ](/api-reference/rerank/dashscope-rerank)[Create conversation Next ](/api-reference/platform-api/conversations/create-conversation)
