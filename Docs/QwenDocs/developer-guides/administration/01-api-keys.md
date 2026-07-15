# API keys

> **Source:** https://docs.qwencloud.com/developer-guides/administration/api-keys

Create and manage API keys

 Copy page ## [​ ](#overview) Overview

Every API request to Qwen Cloud requires an API key for authentication. The [API Keys](https://home.qwencloud.com/api-keys) page provides centralized management of all your access credentials, including pay-as-you-go API keys and Token Plan API keys.

- **Pay-as-you-go API keys**: Each key belongs to one [workspace](/developer-guides/administration/workspace) and inherits the workspace&#x27;s model access and [rate limits](/developer-guides/administration/rate-limits).

- **Token Plan API keys**: Displays the subscriber&#x27;s access credentials for model API calls. Seat assignment and member management are handled in the [Token Plan management portal](https://tokenplan-enterprise.qwencloud.com).


For a quick guide on getting your first API key, see [Get an API key](/api-reference/preparation/api-key).
## [​ ](#create-an-api-key) Create an API key

 1 Select a workspace

Go to the [API Keys](https://home.qwencloud.com/api-keys) page. Use the workspace switcher at the bottom of the sidebar to select the workspace where you want to create the key. 2 Create the key

Click **Create API key**. Enter a description to help you identify the key later (such as "Production API key for main application"), then click **Generate Key**. 3 Copy the key

Copy the API key immediately — **it is shown only once**. After you close the dialog, you can only see a masked version of the key. 
Each workspace supports up to 20 API keys.
## [​ ](#manage-pay-as-you-go-api-keys) Manage pay-as-you-go API keys

In the **Pay-As-You-Go** section of the [API Keys](https://home.qwencloud.com/api-keys) page, you can:

- **Search**: Find keys by description.

- **Edit**: Click **Edit** to update the key&#x27;s description.

- **Delete**: Click **Delete** to permanently revoke the key. This action cannot be undone — all requests using this key will immediately fail.

- **Enable/Disable**: Toggle the status switch to enable or disable an API key. Disabled keys reject all requests; re-enabling restores access immediately.

- **View cost**: Click **View cost** to see the cost breakdown for a specific API key.


The API keys table shows each key&#x27;s **ID**, **API Key** (masked), **Description**, **Created**, **Cost**, **Status**, and available **Actions**.
## [​ ](#manage-token-plan-api-keys) Manage Token Plan API keys

In the **Token Plan** section of the [API Keys](https://home.qwencloud.com/api-keys) page, the subscriber&#x27;s Token Plan credentials are displayed in a table with **API Key** (masked), **Status**, **Seat type**, **Remaining days**, and **Actions** (Reset, Manage).

- **Reset**: Regenerate the API key. The old key is invalidated immediately.

- **Manage**: Navigate to the [Token Plan management portal](https://tokenplan-enterprise.qwencloud.com) for seat assignment, member management, and other operations.


 Full Token Plan management (seat assignment, SSO login, member management) is available in the [Token Plan management portal](https://tokenplan-enterprise.qwencloud.com). See [Team management](/token-plan/team-management) for details. 
## [​ ](#api-keys-and-workspaces) API keys and workspaces

API keys are scoped to the workspace where they are created:

- Keys in the **default workspace** can call all available models with account-level rate limits.

- Keys in a **sub-workspace** can only call models that have been granted to that workspace, with the rate limits configured for that workspace.


To manage API keys for a different workspace, switch to it using the workspace switcher at the bottom of the sidebar.
## [​ ](#api-keys-and-billing) API keys and billing

Pay-as-you-go API keys and Token Plan API keys must be used with their respective endpoints. Do not mix them.
 API calls via DashScope endpoints use pay-as-you-go billing and will not consume your Token Plan subscription quota. 
## [​ ](#security-best-practices) Security best practices


- **Store keys securely**: Use environment variables or a secrets manager. Never hardcode API keys in source code.

- **Don&#x27;t expose keys client-side**: API keys should only be used in server-side code. Never include them in frontend applications, mobile apps, or browser code.

- **Rotate keys regularly**: Delete compromised or unused keys and create new ones.

- **Use separate keys per environment**: Create different keys (or workspaces) for development, testing, and production.


 [Previous ](/developer-guides/run-and-scale/safety)[Workspaces Organize users and access Next ](/developer-guides/administration/workspace)
