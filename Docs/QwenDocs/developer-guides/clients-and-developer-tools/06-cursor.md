# Cursor

> **Source:** https://docs.qwencloud.com/developer-guides/clients-and-developer-tools/cursor

AI-powered code editor

 Copy page Cursor is an AI coding IDE. Connect it to Qwen Cloud via Token Plan (Team Edition), Coding Plan, or pay-as-you-go billing.
## [​ ](#install-cursor) Install Cursor

Download and install Cursor from the [Cursor official website](https://cursor.com/features).
## [​ ](#configure-access-credentials) Configure access credentials

In Cursor, click the settings icon and go to **Cursor Settings** > **Models**. Enable **OpenAI API Key** and **Override OpenAI Base URL**, then enter the API Key, Base URL, and model name corresponding to your selected plan.
### [​ ](#token-plan-team-edition) Token Plan (Team Edition)

ParameterValue**API Key**Token Plan (Team Edition) dedicated [API Key](https://home.qwencloud.com/api-keys)**Base URL**`https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1`**Available models**Token Plan (Team Edition) [supported models](/token-plan/overview#supported-models). Some model names need adjustment: kimi-k2.6 should be written as **kimi-k2-6**, kimi-k2.5 as **kimi-k2-5**, glm-5.2 as **glm-5-2**, glm-5.1 as **glm-5-1**, glm-5 as **glm-5-0**. 
### [​ ](#coding-plan) Coding Plan

ParameterValue**API Key**Coding Plan dedicated [API Key](https://home.qwencloud.com/api-keys)**Base URL**`https://coding-intl.dashscope.aliyuncs.com/v1`**Available models**Coding Plan [supported models](/coding-plan/overview#plan-details). Some model names need adjustment: glm-5 → **glm-5-0**, glm-5.1 → **glm-5-1**, glm-5.2 → **glm-5-2**, glm-4.7 → **glm-4-7**, kimi-k2.6 → **kimi-k2-6**, kimi-k2.5 → **kimi-k2-5**. 
### [​ ](#pay-as-you-go) Pay-as-you-go

ParameterValue**API Key**Your [Qwen Cloud API Key](https://home.qwencloud.com/api-keys)**Base URL**`https://dashscope-intl.aliyuncs.com/compatible-mode/v1`**Available models**Enter a [supported model](/developer-guides/getting-started/text-generation-models). Some model names need to be adjusted: write kimi-k2.6 as **kimi-k2-6**, kimi-k2.5 as **kimi-k2-5**, glm-5.2 as **glm-5-2**, glm-5.1 as **glm-5-1**, glm-5 as **glm-5-0**. 
## [​ ](#faq) FAQ

### [​ ](#error-codes) Error codes

If you encounter errors during configuration, refer to the FAQ documentation for the corresponding billing plan:

- Pay-as-you-go: [Error codes and troubleshooting](/api-reference/preparation/error-messages)

- Coding Plan: [Coding Plan FAQ](/coding-plan/faq)

- Token Plan (Team Edition): [Token Plan (Team Edition) FAQ](/token-plan/faq)


### [​ ](#cannot-invoke-an-added-model-in-cursor) Cannot invoke an added model in Cursor

Error messages:

- The model xxx does not work with your current plan or api key.

- Named models unavailable Free plans can only use Auto. Switch to Auto or upgrade plans to continue.


**Cause**: The free version of Cursor only supports Auto mode and does not support invoking custom models.
**Solution**: Upgrade to **Cursor Pro or a higher plan**.
### [​ ](#cannot-find-the-added-model-after-configuration) Cannot find the added model after configuration

In the chat panel, click to disable **Auto** mode, then select the desired model from the model dropdown.
### [​ ](#model-invocation-error-we-re-having-trouble-connecting-to-the-model-provider-or-unauthorized-user-api-key) Model invocation error: "We&#x27;re having trouble connecting to the model provider." or "Unauthorized User API key"

Troubleshoot the following items:

- Verify that the API Key, Base URL, and model name match the selected billing plan. Credentials are not interchangeable between plans.

- Some model names conflict with Cursor&#x27;s built-in model names or contain dots that need replacement. Refer to the available models section above for details.


 [Previous ](/developer-guides/clients-and-developer-tools/opencode)[Codex OpenAI&#x27;s terminal AI coding assistant Next ](/developer-guides/clients-and-developer-tools/codex)
