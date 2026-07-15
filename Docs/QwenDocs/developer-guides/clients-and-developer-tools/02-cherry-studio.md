# Cherry Studio

> **Source:** https://docs.qwencloud.com/developer-guides/clients-and-developer-tools/cherry-studio

Open-source AI desktop client

 Copy page Cherry Studio is an open-source AI desktop client. Connect it to Qwen Cloud via Token Plan (Team Edition), Coding Plan, or pay-as-you-go billing.
## [​ ](#install-cherry-studio) Install Cherry Studio

Go to the [Cherry Studio download page](https://www.cherry-ai.com/download), download the installation package for your operating system, and complete the installation.
## [​ ](#configure-credentials) Configure credentials

Open Cherry Studio, click the settings button in the upper-right corner, and in the **Model Service** section click **Add**. Enter a provider name (for example, Token Plan (Team Edition)) and select OpenAI as the provider type.
### [​ ](#token-plan-team-edition) Token Plan (Team Edition)

ParameterDescription**API key**Enter the Token Plan (Team Edition) dedicated [API key](https://home.qwencloud.com/api-keys).**API address**`https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1`**Model**For available models, see Token Plan (Team Edition) [supported models](/token-plan/overview#supported-models). 
### [​ ](#coding-plan) Coding Plan

ParameterDescription**API key**Enter the Coding Plan dedicated [API key](https://home.qwencloud.com/api-keys).**API address**`https://coding-intl.dashscope.aliyuncs.com/compatible-mode/v1`**Model**For available models, see Coding Plan [supported models](/coding-plan/overview#plan-details). 
### [​ ](#pay-as-you-go) Pay-as-you-go

ParameterDescription**API key**Enter your [Qwen Cloud API Key](https://home.qwencloud.com/api-keys).**API address**`https://dashscope-intl.aliyuncs.com/compatible-mode/v1`**Model**Enter a [supported model](/developer-guides/getting-started/text-generation-models). 
## [​ ](#verify-configuration) Verify configuration

In the **Model ID** field, enter the model you want to use (for example, `qwen3.7-max`), and click **Add**. Return to the chat interface, enter any question, and confirm that the model responds normally. If it does, the configuration is complete.
## [​ ](#faq) FAQ

### [​ ](#error-codes) Error codes

If you encounter errors during configuration, refer to the FAQ documentation for the corresponding billing plan:

- Pay-as-you-go: [Error codes and troubleshooting](/api-reference/preparation/error-messages)

- Coding Plan: [Coding Plan FAQ](/coding-plan/faq)

- Token Plan (Team Edition): [Token Plan (Team Edition) FAQ](/token-plan/faq)


### [​ ](#error-the-value-of-the-enable-thinking-parameter-is-restricted-to-true) Error: The value of the enable_thinking parameter is restricted to True

**Cause**: This model only supports thinking mode, but thinking mode was not enabled when the model was called.
**Solution**: Enable thinking mode in the client.
### [​ ](#charges-incurred-despite-having-a-free-quota-when-using-pay-as-you-go) Charges incurred despite having a free quota when using Pay-as-you-go

Possible causes:

- **Incorrect API address**: The free quota requires using the correct Qwen Cloud API address (`https://dashscope-intl.aliyuncs.com/compatible-mode/v1`). Check that the **API address** in your configuration is correct. For details, see [Free quota](/resources/free-quota).

- **Per-model free quota**: The free quota for each model is calculated independently and cannot be shared across models.

- **Data update delay**: The console updates free quota data hourly. Your free quota may already be exhausted even if the console still shows a remaining balance.


 [Previous ](/developer-guides/clients-and-developer-tools/qwenpaw)[Chatbox Cross-platform AI chat client Next ](/developer-guides/clients-and-developer-tools/chatbox)
