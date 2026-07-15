# Hermes Agent

> **Source:** https://docs.qwencloud.com/developer-guides/clients-and-developer-tools/hermes-agent

Terminal AI coding assistant by Nous Research

 Copy page Hermes Agent is a terminal AI coding tool that connects to Qwen Cloud through pay-as-you-go, Coding Plan, or Token Plan (Team Edition).
## [​ ](#install-hermes-agent) Install Hermes Agent


- Run this command to install Hermes Agent. The script automatically installs dependencies like Python and Git.


Copy ```\ncurl -fsSL https://raw.githubusercontent.com/NousResearch/hermes-agent/main/scripts/install.sh | bash

``` 
 Windows does not support native installation. Install [WSL2](https://learn.microsoft.com/en-us/windows/wsl/install) first, then run the above command in WSL2. 

- After installation, reload the terminal environment.


Copy ```\nsource ~/.bashrc # If using zsh, change to source ~/.zshrc

``` 

- Verify the installation. A version number confirms success.


Copy ```\nhermes --version

``` 
## [​ ](#configure-access-credentials) Configure access credentials

Run `hermes config set` to configure the Base URL and API Key for your chosen plan:

- **Token Plan (Team Edition)**: Seat-based subscription; token consumption deducts Credits.

- **Coding Plan**: Fixed monthly subscription, metered by model invocation count.

- **Pay-as-you-go**: Post-paid based on actual usage.


 The examples in this topic use the **Anthropic-compatible protocol**: the Base URL ends with `/apps/anthropic` and `api_mode` is set to `anthropic_messages`. Hermes Agent also supports the **OpenAI-compatible protocol**: replace the trailing `/apps/anthropic` in the Base URL with `/compatible-mode/v1` and remove the `api_mode` setting. For example, the OpenAI-compatible Base URL for Pay-as-you-go (Singapore) is `https://dashscope-intl.aliyuncs.com/compatible-mode/v1`.In addition to the command-line version, Hermes Agent also offers a desktop version (Hermes Desktop). You can download the installer from the [Hermes website](https://hermes-agent.nousresearch.com/), or run `hermes desktop` after installing the command-line version. Hermes Desktop and the command-line version share the same `~/.hermes/config.yaml` configuration file, so the access parameters are identical to those described here. When connecting through a Custom Endpoint in the desktop app, use the OpenAI-compatible Base URL described above. 
### [​ ](#token-plan-team-edition) Token Plan (Team Edition)

Replace `YOUR_API_KEY` with your Token Plan (Team Edition) [API Key](https://home.qwencloud.com/api-keys). Available models are listed in [supported models](/token-plan/overview#supported-models).
Copy ```\nhermes config set model.provider custom
hermes config set model.base_url https://token-plan.ap-southeast-1.maas.aliyuncs.com/apps/anthropic
hermes config set model.api_mode anthropic_messages
hermes config set model.api_key YOUR_API_KEY
hermes config set model.default qwen3.7-max

``` 
Show config.yaml example

 These commands write to `~/.hermes/config.yaml`. You can also edit the file directly:Copy ```\nmodel:
 default: qwen3.7-max
 provider: custom
 base_url: https://token-plan.ap-southeast-1.maas.aliyuncs.com/apps/anthropic
 api_mode: anthropic_messages
 api_key: YOUR_API_KEY

``` 
### [​ ](#coding-plan) Coding Plan

Replace `YOUR_API_KEY` with your Coding Plan [API Key](https://home.qwencloud.com/api-keys). Available models are listed in [supported models](/coding-plan/overview#plan-details).
Copy ```\nhermes config set model.provider custom
hermes config set model.base_url https://coding-intl.dashscope.aliyuncs.com/apps/anthropic
hermes config set model.api_mode anthropic_messages
hermes config set model.api_key YOUR_API_KEY
hermes config set model.default qwen3.7-plus

``` 
Show config.yaml example

 These commands write to `~/.hermes/config.yaml`. You can also edit the file directly:Copy ```\nmodel:
 default: qwen3.7-plus
 provider: custom
 base_url: https://coding-intl.dashscope.aliyuncs.com/apps/anthropic
 api_mode: anthropic_messages
 api_key: YOUR_API_KEY

``` 
### [​ ](#pay-as-you-go) Pay-as-you-go

Replace `YOUR_API_KEY` with your [Qwen Cloud API Key](https://home.qwencloud.com/api-keys). Available models are listed in [Anthropic compatible API](/api-reference/chat/openai-chat).
Set `base_url` for your region. The API Key must match the selected region:

- Singapore: `https://dashscope-intl.aliyuncs.com/apps/anthropic`


Copy ```\nhermes config set model.provider alibaba
hermes config set model.base_url https://dashscope-intl.aliyuncs.com/apps/anthropic
hermes config set model.api_mode anthropic_messages
hermes config set model.api_key YOUR_API_KEY
hermes config set model.default qwen3.7-max

``` 
Show config.yaml example

 These commands write to `~/.hermes/config.yaml`. You can also edit the file directly:Copy ```\nmodel:
 default: qwen3.7-max
 provider: alibaba
 base_url: https://dashscope-intl.aliyuncs.com/apps/anthropic
 api_mode: anthropic_messages
 api_key: YOUR_API_KEY

``` 
## [​ ](#verify-configuration) Verify configuration

Send a test message to verify your configuration:
Copy ```\nhermes chat -q "Hello"

``` 
A successful response confirms the configuration works. To switch models, use the `-m` parameter:
Copy ```\nhermes chat -m qwen3.7-max

``` 
## [​ ](#faq) FAQ

### [​ ](#error-codes) Error codes

For configuration errors, check the FAQ for your billing plan:

- Pay-as-you-go: [Error codes and troubleshooting](/api-reference/preparation/error-messages)

- Coding Plan: [Coding Plan FAQ](/coding-plan/faq)

- Token Plan (Team Edition): [Token Plan (Team Edition) FAQ](/token-plan/faq)


 [Previous ](/developer-guides/clients-and-developer-tools/openclaw)[Claude Code Anthropic&#x27;s terminal AI coding assistant Next ](/developer-guides/clients-and-developer-tools/claude-code)
