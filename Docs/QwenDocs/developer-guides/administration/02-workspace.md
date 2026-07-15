# Workspaces

> **Source:** https://docs.qwencloud.com/developer-guides/administration/workspace

Organize users and access

 Copy page ## [​ ](#overview) Overview

A workspace is the basic unit for organizing resources, permissions, and billing in Qwen Cloud. Each workspace has its own set of API keys, model access, and usage quotas.
Workspaces let you:

- Isolate resources between departments, teams, or projects

- Control which models can be called and set per-model rate limits

- Manage API keys independently per workspace

- Track usage and billing separately


## [​ ](#create-a-workspace) Create a workspace

 1 Go to the Workspaces page

Go to **Settings** > [Workspaces](https://home.qwencloud.com/settings/workspaces). 2 Create a workspace

Click **Create Workspace**. Enter a workspace name and click **Save Changes**. 
Your account includes a default workspace that is created automatically. You can create additional workspaces as needed. The workspace limit is displayed on the **Create Workspace** button (such as 5/20).
## [​ ](#manage-workspaces) Manage workspaces

From the [Workspaces](https://home.qwencloud.com/settings/workspaces) page, you can:

- **Search**: Find workspaces by name or workspace ID.

- **Edit**: Click **Edit** next to a workspace to open the editing panel.

- **Switch**: Use the **Current Workspace** dropdown to filter the view, or use the workspace switcher at the bottom of the sidebar to change your active workspace.


Each workspace shows its name, workspace ID, creation date, and the number of API keys associated with it.
## [​ ](#edit-a-workspace) Edit a workspace

Click **Edit** on any workspace to open the editing panel with two sections:
### [​ ](#basic-info) Basic info


- **Workspace Name**: Rename the workspace (max 30 characters).

- **Total API Keys**: Shows the current number of API keys and the maximum allowed (such as 0/20).


### [​ ](#model-permission) Model permission

Control which models can be called using API keys in this workspace.

- 
**Add models**: Click **All Models** to open the model selection dialog. You can search by model name, filter by **Model Vendor** or **Model Type**, and select models to add. Use **Select All Results** to check every model on the current page.


- 
**Rate limits**: For each model you add, set custom rate limits:

**Times / min**: Maximum number of API requests per minute (RPM).

- **Token / min**: Maximum number of tokens processed per minute (TPM). Some model types (such as image generation) do not use token-based metering and may not show this field.


The RPM and TPM values you set cannot exceed the account-level limits for that model. You can check each model&#x27;s current limits on the [Analytics](https://home.qwencloud.com/analytics) page.

- 
Click **Select (N)** to confirm your selection in the modal, then click **Save Changes** on the panel to persist the changes.


 The model selection dialog is paginated. Selections on one page may not carry over when you navigate to another page. Use search and filters to narrow results, or save after each page. 
The default workspace has access to all models with account-level rate limits and is not editable. To customize model access or rate limits, create a sub-workspace.
## [​ ](#switch-workspaces) Switch workspaces

You can switch your active workspace in two ways:

- **Workspace switcher**: Click the workspace name at the bottom of the sidebar to open the switcher. Select the workspace you want to use.

- **Workspaces page**: Use the **Current Workspace** dropdown to filter and select a workspace.


When you switch workspaces, the API Keys page, Analytics, and Monitoring data all reflect the selected workspace.
## [​ ](#api-keys-and-workspaces) API keys and workspaces

Each API key belongs to one workspace. When calling models through the API, the key determines which workspace&#x27;s model access and rate limits apply.

- Default workspace keys can call all available models.

- Sub-workspace keys can only call models that have been granted to that workspace.


To manage API keys for a specific workspace, switch to it first, then go to the [API Keys](https://home.qwencloud.com/api-keys) page.
## [​ ](#limits) Limits

ResourceLimitWorkspaces per account20 (includes the default workspace)API keys per workspace20Workspace name length30 characters 
The default workspace is created automatically and cannot be deleted. It can be renamed.
## [​ ](#faq) FAQ

### [​ ](#how-do-i-get-the-workspace-id) How do I get the workspace ID?

Go to the [Workspaces](https://home.qwencloud.com/settings/workspaces) page. The workspace ID is shown in the **Workspace ID** column of the table.
### [​ ](#how-do-i-call-models-from-a-specific-workspace) How do I call models from a specific workspace?

Use an API key that belongs to that workspace. The key determines the workspace context for the API call.
### [​ ](#how-do-i-restrict-which-models-a-team-can-use) How do I restrict which models a team can use?

Create a workspace for the team, click **Edit**, then use **Model Permission** to add only the models they need. Any API key in that workspace can only call the allowed models.
### [​ ](#how-do-i-set-rate-limits-for-a-specific-model) How do I set rate limits for a specific model?

When adding models to a workspace, set the **Times / min** (RPM) and **Token / min** (TPM) fields for each model before clicking **Select (N)**, then **Save Changes**.
### [​ ](#why-is-model-permission-missing-from-the-editing-panel) Why is Model Permission missing from the editing panel?

You are editing the default workspace. The default workspace has access to all models and does not expose a model permission section.
### [​ ](#can-i-delete-the-default-workspace) Can I delete the default workspace?

No. The default workspace cannot be deleted, but you can rename it. [Previous ](/developer-guides/administration/api-keys)[Rate limits Understand and manage API rate limits Next ](/developer-guides/administration/rate-limits)
