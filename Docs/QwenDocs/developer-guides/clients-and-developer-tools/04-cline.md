# Cline

> **Source:** https://docs.qwencloud.com/developer-guides/clients-and-developer-tools/cline

VSCode AI coding extension

 Copy page Cline is a VS Code extension for AI-assisted coding. Connect it to Qwen Cloud via Token Plan (Team Edition), Coding Plan, or pay-as-you-go billing.
## [​ ](#install-cline) Install Cline


- Download and install [VSCode](https://code.visualstudio.com/).

- Open VSCode, search for `Cline` in the extension marketplace, and install it.


## [​ ](#configure-credentials) Configure credentials

After installation, click the Cline icon in the left sidebar to open the configuration interface. Click **Bring my own API key**, select **OpenAI Compatible** as the API Provider, and fill in the parameters for your billing plan. If you have used Cline before, click the settings icon in the upper-right corner to open the configuration interface.
### [​ ](#token-plan-team-edition) Token Plan (Team Edition)

ParameterDescription**API Provider**Select **OpenAI Compatible**.**Base URL**`https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1`**API Key**Enter the Token Plan (Team Edition) dedicated [API key](https://home.qwencloud.com/api-keys).**Model ID**Enter a [supported model](/token-plan/overview#supported-models), such as `qwen3.7-max`. 
### [​ ](#coding-plan) Coding Plan

ParameterDescription**API Provider**Select **OpenAI Compatible**.**Base URL**`https://coding-intl.dashscope.aliyuncs.com/v1`**API Key**Enter the Coding Plan dedicated [API key](https://home.qwencloud.com/api-keys).**Model ID**Enter a [supported model](/coding-plan/overview#plan-details), such as `qwen3.7-plus`. 
### [​ ](#pay-as-you-go) Pay-as-you-go

ParameterDescription**API Provider**Select **OpenAI Compatible**.**Base URL**`https://dashscope-intl.aliyuncs.com/compatible-mode/v1`**API Key**Enter your [Qwen Cloud API Key](https://home.qwencloud.com/api-keys).**Model ID**Enter a [supported model](/developer-guides/getting-started/text-generation-models). 
If you use Qwen3 (thinking mode) or the QwQ model, click `MODEL CONFIGURATION` in the settings interface and check **Enable R1 messages format**.
## [​ ](#faq) FAQ

### [​ ](#error-codes) Error codes

If you encounter errors during configuration, refer to the FAQ documentation for the corresponding billing plan:

- Pay-as-you-go: [Error codes and troubleshooting](/api-reference/preparation/error-messages)

- Coding Plan: [Coding Plan FAQ](/coding-plan/faq)

- Token Plan (Team Edition): [Token Plan (Team Edition) FAQ](/token-plan/faq)


### [​ ](#error-401-incorrect-api-key-provided) Error: 401 Incorrect API key provided

Possible causes:

- The API key does not match the Base URL. API keys are not interchangeable across billing plans. Verify that the API key and Base URL belong to the same plan.


### [​ ](#error-400-internalerror-algo-invalidparameter) Error: 400 InternalError.Algo.InvalidParameter

Click `MODEL CONFIGURATION` in the settings interface and check **Enable R1 messages format**. [Previous ](/developer-guides/clients-and-developer-tools/chatbox)[Qoder Agentic coding platform with IDE, CLI, and JetBrains plugin Next ](/developer-guides/clients-and-developer-tools/qoder)
