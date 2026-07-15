# Data security &amp; privacy

> **Source:** https://docs.qwencloud.com/developer-guides/security-compliance/data-security

How Qwen Cloud secures your data

 Copy page Qwen Cloud is designed to protect your data at every layer — from API requests to model inference. This page describes the key security measures in place.
## [​ ](#encryption) Encryption

LayerStandardData in transitTLS 1.2 or later for all API connectionsData at restAES-256 encryption for stored data (API keys, account information) 
## [​ ](#api-key-security) API key security

API keys authenticate every request to Qwen Cloud. Follow these practices to keep your keys safe:

- **Never hardcode keys** in source code or commit them to version control. Use environment variables or a secrets manager.

- **Use separate keys** for development and production environments.

- **Rotate keys** regularly. You can create new keys and delete old ones from the [API Keys](https://home.qwencloud.com/api-keys) page.

- **Restrict access** by creating keys in specific [workspaces](/developer-guides/administration/workspace) with appropriate permissions.


For detailed API key management instructions, see [API keys](/developer-guides/administration/api-keys).
## [​ ](#content-moderation) Content moderation

All API requests pass through automatic content moderation that screens both inputs and outputs for harmful, illegal, or inappropriate content.
## [​ ](#responsible-ai-practices) Responsible AI practices

As a developer building on Qwen Cloud, you share responsibility for the safety and security of your application:

- Implement input validation before passing user content to the API.

- Set appropriate `max_tokens` limits for your use case.

- Add rate limiting at the application layer to prevent abuse.

- Monitor both inputs and outputs for content safety using Qwen Cloud&#x27;s built-in content moderation.


## [​ ](#learn-more) Learn more


- [Audit & access Logs](/developer-guides/security-compliance/audit-logs): Track API usage for compliance.


 [Previous ](/developer-guides/deployment/manage-deployments)[Audit & access logs Track API usage for compliance Next ](/developer-guides/security-compliance/audit-logs)
