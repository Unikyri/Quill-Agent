# QwenPaw

> **Source:** https://docs.qwencloud.com/developer-guides/clients-and-developer-tools/qwenpaw

Open-source personal AI assistant from the AgentScope team

 Copy page QwenPaw (formerly CoPaw) is an open-source personal AI assistant from the AgentScope team. It supports local and cloud deployment, and integrates with Qwen Cloud via Token Plan Team Edition, Coding Plan, or pay-as-you-go.
## [​ ](#install-qwenpaw) Install QwenPaw

Use the pip package or one-click installation script. For Docker, desktop app, or ModelScope online runtime, see the [QwenPaw official documentation](https://qwenpaw.agentscope.io/).
- pip install 
- One-click script 

 Requires Python 3.10 ~ 3.13:Copy ```\npip install qwenpaw
qwenpaw init --defaults
qwenpaw app

``` The script installs uv, creates a virtual environment, and downloads dependencies automatically. No manual Python setup required. Choose the command for your operating system:
- macOS / Linux:


Copy ```\ncurl -fsSL https://qwenpaw.agentscope.io/install.sh | bash

``` 
- Windows (CMD):


Copy ```\ncurl -fsSL https://qwenpaw.agentscope.io/install.bat -o install.bat && install.bat

``` 
- Windows (PowerShell):


Copy ```\nirm https://qwenpaw.agentscope.io/install.ps1 | iex

``` After installation, run in a new terminal:Copy ```\nqwenpaw init --defaults
qwenpaw app

``` 
After launch, visit `http://127.0.0.1:8088/` to open QwenPaw Console.
## [​ ](#configure-credentials) Configure credentials

In Console, click **Settings** > **Models** and configure the provider for your billing plan.
### [​ ](#token-plan-team-edition) Token Plan Team Edition

QwenPaw does not include a built-in provider for Token Plan Team Edition. On the **Providers** page, click **Add Provider**. Set **Protocol** to **OpenAI-compatible (Chat Completions)** (Provider ID and Name are customizable, e.g. `bailian-token-plan`). After saving, open the **Settings** page and fill in the table below.
**Configuration****Description****Base URL**`https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1`**API Key**Enter the Token Plan Team Edition dedicated [API Key](https://home.qwencloud.com/api-keys).**Models**On the provider **Models** page, click **Add Model** and set **Model ID** to a [supported Token Plan Team Edition model](/token-plan/overview#supported-models). 
### [​ ](#coding-plan) Coding Plan

Open the built-in **Aliyun Coding Plan (International)** provider **Settings** page and fill in the API Key.
**Configuration****Description****API Key**Enter the Coding Plan dedicated [API Key](https://home.qwencloud.com/api-keys).**Models**Common models are pre-configured. Verify them via **Test Connection** on the **Models** page. To add a model, click **Add Model** and set **Model ID** to a [supported Coding Plan model](/coding-plan/overview#plan-details). 
### [​ ](#pay-as-you-go) Pay-as-you-go

Open the **DashScope** provider **Settings** page and fill in the API Key.
**Configuration****Description****API Key**Enter the [Qwen Cloud API Key](/api-reference/preparation/api-key).**Base URL**`https://dashscope-intl.aliyuncs.com/compatible-mode/v1`**Models**Common models are pre-configured. To add a model, click **Add Model** and set **Model ID** to a [supported model](/api-reference/chat/openai-chat). 
## [​ ](#set-the-default-model) Set the default model

Open **Settings** > **Models** > **Default LLM**, select a model, and click **Save**. The dropdown in the top-right of the chat page switches the provider and model for the current session.
## [​ ](#faq) FAQ

### [​ ](#error-codes) Error codes

Troubleshoot by billing plan:

- Pay-as-you-go: [Error code troubleshooting](/api-reference/preparation/error-messages)

- Coding Plan: [Coding Plan FAQ](/coding-plan/faq)

- Token Plan Team Edition: [Token Plan Team Edition FAQ](/token-plan/faq)


### [​ ](#error-401-incorrect-api-key-provided) Error 401 Incorrect API key provided

Possible causes:

- API Keys are not interchangeable across the three billing plans. Make sure the API Key and Base URL come from the same plan.

- The pay-as-you-go API Key and Base URL are in different regions.


### [​ ](#context-length-exceeded-during-long-conversations-or-tool-calls) Context length exceeded during long conversations or tool calls

On the provider **Settings** page for that model, expand **Advanced Configuration** and adjust generation parameters such as `max_tokens` in JSON format, then save:
Copy ```\n{
 "temperature": 0.7,
 "top_p": 0.9,
 "max_tokens": 4096
}

``` [Previous ](/developer-guides/clients-and-developer-tools/qwen-code)[Cherry Studio Open-source AI desktop client Next ](/developer-guides/clients-and-developer-tools/cherry-studio)
