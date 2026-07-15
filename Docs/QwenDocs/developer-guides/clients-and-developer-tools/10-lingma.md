# Qoder CN (formerly Lingma)

> **Source:** https://docs.qwencloud.com/developer-guides/clients-and-developer-tools/lingma

Alibaba Cloud intelligent coding assistant IDE

 Copy page Qoder CN (formerly Lingma) is Alibaba Cloud&#x27;s intelligent coding assistant that provides a standalone IDE, and can connect to Qwen Cloud via Token Plan, Coding Plan, or pay-as-you-go.
 Qoder CN Personal Community Edition and Personal Professional Edition both support connecting to Qwen Cloud. Enterprise Edition is not supported. 
## [​ ](#install) Install


- Go to the [Qoder CN official website](https://lingma.aliyun.com/download) to download and install Qoder CN.

- Complete the initial configuration after the first launch.

- On the Alibaba Cloud login page, log in with your Alibaba Cloud account.


## [​ ](#configure-access-credentials) Configure access credentials


- Open Qoder CN settings from the top-right corner of the interface, select **Model**, and click **Add**.

- Model configuration details are as follows:


Configuration ItemDescriptionProviderSelect **Alibaba Cloud Model Studio - International** from the dropdown menuTypeSelect **Token Plan**, **Coding Plan**, or **Pay-as-you-go** based on your billing planModelSelect a model from the dropdown menu.API keyEnter the dedicated API key for your chosen plan: Token Plan (Team Edition) - [Get API Key](https://home.qwencloud.com/api-keys); Coding Plan - [Get API Key](https://home.qwencloud.com/api-keys); Pay-as-you-go - [Get API Key](/api-reference/preparation/api-key) 

- Click **Add**. The model configuration is complete once validation passes.

- Select the corresponding model in the Qoder CN chat panel to start using it.


## [​ ](#learn-more) Learn more

To learn more about Qoder CN&#x27;s agents, MCP, Skills, and other extension capabilities, refer to the [Qoder CN official documentation](https://www.alibabacloud.com/help/en/lingma/product-overview/introduction-of-lingma).
## [​ ](#faq) FAQ

### [​ ](#error-codes) Error codes

If you encounter errors during configuration, refer to the FAQ documentation for the corresponding billing plan:

- Coding Plan: [Coding Plan FAQ](/coding-plan/faq)

- Token Plan (Team Edition): [Token Plan (Team Edition) FAQ](/token-plan/faq)

- Pay-as-you-go: [Error Codes](/api-reference/preparation/error-messages)


### [​ ](#why-is-the-model-option-missing-in-qoder-cn-settings) Why is the Model option missing in Qoder CN settings?

Possible reasons:

- **Not signed in**: You must sign in first before you can chat and configure models.

- **Unsupported edition**: Connecting to Qwen Cloud requires Qoder CN Personal Community Edition or Personal Professional Edition. Enterprise Edition is not supported.


### [​ ](#chat-error-custom-model-service-exception-please-try-again-later-or-switch-to-another-model-unknown-custom-model-exception) Chat error: "Custom model service exception, please try again later or switch to another model. Unknown Custom model Exception"

This error is Qoder CN&#x27;s generic message when it receives an unrecognized backend response. Common causes include:

- **Provider or Type does not match the actual plan**: In the Qoder CN model configuration, **Provider** and **Type** must match the purchased plan. For example, using a Token Plan (Team Edition) API key while **Type** is set to Coding Plan.

- **Selected model is not supported by the plan**: Only text generation models covered by the current plan are supported. For example, the supported model list for Token Plan (Team Edition) is available on the [Token Plan (Team Edition) page](/token-plan/overview#supported-models).

- **Temporary network or service fluctuation**: Try again later.


### [​ ](#api-key-authentication-failed-http-401) API Key Authentication Failed (HTTP 401)

Check the following:

- Verify that the API key matches your billing plan. API keys for Token Plan (Team Edition) and Coding Plan are not interchangeable.

- Verify that the plan has not expired.

- Ensure the API key is copied completely without extra spaces. If the error persists, reset the API key from the corresponding management page.


 [Previous ](/developer-guides/clients-and-developer-tools/qoder)[Kilo CLI AI coding in terminal Next ](/developer-guides/clients-and-developer-tools/kilo-cli)
