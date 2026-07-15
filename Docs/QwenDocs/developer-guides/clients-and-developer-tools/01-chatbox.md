# Chatbox

> **Source:** https://docs.qwencloud.com/developer-guides/clients-and-developer-tools/chatbox

Cross-platform AI chat client

 Copy page Chatbox is a cross-platform AI client. Connect it to Qwen Cloud via Token Plan (Team Edition), Coding Plan, or pay-as-you-go billing.
## [​ ](#install-chatbox) Install Chatbox

Go to the [Chatbox website](https://chatboxai.app/en) and download the installer for your operating system, or use the web version directly.
## [​ ](#configure-access-credentials) Configure access credentials

In Chatbox, click **Settings** at the bottom left, click **Model Provider**, and click **Add** at the bottom. In the pop-up, enter a **Name**, set **API Mode** to **OpenAI API Compatible**, and click **Add**. Then fill in the API Key and API Host according to your billing plan.
### [​ ](#token-plan-team-edition) Token Plan (Team Edition)

ParameterDescription**API Key**Enter the Token Plan (Team Edition) dedicated [API Key](https://home.qwencloud.com/api-keys).**API Host**`https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1`. Leave the API Path empty.**Model**Click **New** and enter a model name in **Model ID**. For available models, see [Token Plan supported models](/token-plan/overview#supported-models). 
### [​ ](#coding-plan) Coding Plan

ParameterDescription**API Key**Enter the Coding Plan dedicated [API Key](https://home.qwencloud.com/api-keys).**API Host**`https://coding-intl.dashscope.aliyuncs.com/compatible-mode/v1`. Leave the API Path empty.**Model**Click **New** and enter a model name in **Model ID**. For available models, see [Coding Plan supported models](/coding-plan/overview#plan-details). 
### [​ ](#pay-as-you-go) Pay-as-you-go

ParameterDescription**API Key**Enter the [Qwen Cloud API Key](https://home.qwencloud.com/api-keys).**API Host**`https://dashscope-intl.aliyuncs.com/compatible-mode/v1`. Leave the API Path empty.**Model**Click **New** and enter a [supported model](/developer-guides/getting-started/text-generation-models) name in **Model ID**. 
## [​ ](#verify-configuration) Verify configuration

After completing the configuration, type "Hello" in the chat box and send it. If the model responds normally, the configuration is successful.
## [​ ](#faq) FAQ

### [​ ](#error-codes) Error codes

If you encounter errors during configuration, refer to the FAQ documentation for your billing plan:

- Pay-as-you-go: [Error codes and troubleshooting](/api-reference/preparation/error-messages)

- Coding Plan: [Coding Plan FAQ](/coding-plan/faq)

- Token Plan (Team Edition): [Token Plan (Team Edition) FAQ](/token-plan/faq)


 [Previous ](/developer-guides/clients-and-developer-tools/cherry-studio)[Cline VSCode AI coding extension Next ](/developer-guides/clients-and-developer-tools/cline)
