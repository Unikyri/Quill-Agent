# OpenCode

> **Source:** https://docs.qwencloud.com/developer-guides/clients-and-developer-tools/opencode

Terminal AI coding assistant

 Copy page OpenCode is a terminal AI coding tool. Connect it to Qwen Cloud via Token Plan (Team Edition), Coding Plan, or pay-as-you-go billing.
## [​ ](#install-opencode) Install OpenCode


- 
Install or update [Node.js](https://nodejs.org/en/download/) (v18.0 or later).


- 
Run the following command in the terminal to install OpenCode.


Copy ```\nnpm install -g opencode-ai

``` 
Run the following command to verify the installation. The installation is successful if a version number is displayed.
Copy ```\nopencode -v

``` 
## [​ ](#configure-access-credentials) Configure access credentials

Create and open the configuration file `opencode.json` at one of the following paths:

- macOS / Linux: `~/.config/opencode/opencode.json`

- Windows: `C:\Users\<Your-Username>\.config\opencode\opencode.json`


Add the configuration for your billing plan.
### [​ ](#token-plan-team-edition) Token Plan (Team Edition)

Replace `YOUR_API_KEY` with the Token Plan (Team Edition) dedicated [API Key](https://home.qwencloud.com/api-keys). For available models, see [Token Plan supported models](/token-plan/overview#supported-models).
Copy ```\n{
 "$schema": "https://opencode.ai/config.json",
 "provider": {
 "bailian-token-plan": {
 "npm": "@ai-sdk/anthropic",
 "name": "Qwen Cloud",
 "options": {
 "baseURL": "https://token-plan.ap-southeast-1.maas.aliyuncs.com/apps/anthropic/v1",
 "apiKey": "YOUR_API_KEY"
 },
 "models": {
 "qwen3.7-max": {
 "name": "Qwen3.7 Max",
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 8192
 }
 }
 },
 "qwen3.7-plus": {
 "name": "Qwen3.7 Plus",
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 8192
 }
 }
 },
 "qwen3.6-plus": {
 "name": "Qwen3.6 Plus",
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 8192
 }
 }
 },
 "qwen3.6-flash": {
 "name": "Qwen3.6 Flash",
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 8192
 }
 }
 },
 "deepseek-v4-pro": {
 "name": "DeepSeek V4 Pro"
 },
 "deepseek-v4-flash": {
 "name": "DeepSeek V4 Flash"
 },
 "deepseek-v3.2": {
 "name": "DeepSeek V3.2"
 },
 "kimi-k2.7-code": {
 "name": "Kimi K2.7 Code",
 "modalities": {
 "input": ["text", "image"],
 "output": ["text"]
 },
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 8192
 }
 }
 },
 "kimi-k2.6": {
 "name": "Kimi K2.6",
 "modalities": {
 "input": ["text", "image"],
 "output": ["text"]
 },
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 8192
 }
 }
 },
 "kimi-k2.5": {
 "name": "Kimi K2.5",
 "modalities": {
 "input": ["text", "image"],
 "output": ["text"]
 },
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 8192
 }
 }
 },
 "glm-5.2": {
 "name": "GLM-5.2",
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 8192
 }
 }
 },
 "glm-5.1": {
 "name": "GLM-5.1",
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 8192
 }
 }
 },
 "glm-5": {
 "name": "GLM-5",
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 8192
 }
 }
 },
 "MiniMax-M2.5": {
 "name": "MiniMax M2.5"
 }
 }
 }
 }
}

``` 
### [​ ](#coding-plan) Coding Plan

Replace `YOUR_API_KEY` with the Coding Plan dedicated [API Key](https://home.qwencloud.com/api-keys). For available models, see [Coding Plan supported models](/coding-plan/overview#plan-details).
Copy ```\n{
 "$schema": "https://opencode.ai/config.json",
 "provider": {
 "bailian-coding-plan": {
 "npm": "@ai-sdk/anthropic",
 "name": "Qwen Cloud",
 "options": {
 "baseURL": "https://coding-intl.dashscope.aliyuncs.com/apps/anthropic/v1",
 "apiKey": "YOUR_API_KEY"
 },
 "models": {
 "qwen3.7-plus": {
 "name": "Qwen3.7 Plus",
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 1024
 }
 }
 },
 "qwen3.6-plus": {
 "name": "Qwen3.6 Plus",
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 1024
 }
 }
 },
 "qwen3.5-plus": {
 "name": "Qwen3.5 Plus",
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 1024
 }
 }
 },
 "qwen3-max-2026-01-23": {
 "name": "Qwen3 Max 0123"
 },
 "qwen3-coder-next": {
 "name": "Qwen3 Coder Next"
 },
 "qwen3-coder-plus": {
 "name": "Qwen3 Coder Plus"
 },
 "MiniMax-M2.5": {
 "name": "MiniMax M2.5"
 },
 "glm-5": {
 "name": "GLM-5",
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 1024
 }
 }
 },
 "glm-4.7": {
 "name": "GLM-4.7",
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 1024
 }
 }
 },
 "kimi-k2.5": {
 "name": "Kimi K2.5",
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 1024
 }
 }
 }
 }
 }
 }
}

``` 
### [​ ](#pay-as-you-go) Pay-as-you-go

Replace `YOUR_API_KEY` with the [Qwen Cloud API Key](https://home.qwencloud.com/api-keys). For available models, see [Anthropic compatible API](/api-reference/chat/anthropic).
Copy ```\n{
 "$schema": "https://opencode.ai/config.json",
 "provider": {
 "bailian-payg": {
 "npm": "@ai-sdk/anthropic",
 "name": "Qwen Cloud",
 "options": {
 "baseURL": "https://dashscope-intl.aliyuncs.com/apps/anthropic/v1",
 "apiKey": "YOUR_API_KEY"
 },
 "models": {
 "qwen3.7-max": {
 "name": "Qwen3.7 Max",
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 8192
 }
 }
 },
 "qwen3.7-plus": {
 "name": "Qwen3.7 Plus",
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 1024
 }
 }
 },
 "qwen3.6-plus": {
 "name": "Qwen3.6 Plus",
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 8192
 }
 }
 },
 "deepseek-v3.2": {
 "name": "DeepSeek V3.2"
 }
 }
 }
 }
}

``` 
## [​ ](#verify-configuration) Verify configuration

After saving the configuration, restart OpenCode, type `/models`, search for `Qwen Cloud`, and select the model you want to use.
## [​ ](#faq) FAQ

### [​ ](#error-codes) Error codes

If you encounter errors during configuration, refer to the FAQ documentation for your billing plan:

- Pay-as-you-go: [Error codes and troubleshooting](/api-reference/preparation/error-messages)

- Coding Plan: [Coding Plan FAQ](/coding-plan/faq)

- Token Plan (Team Edition): [Token Plan (Team Edition) FAQ](/token-plan/faq)


 [Previous ](/developer-guides/clients-and-developer-tools/claude-code)[Cursor AI-powered code editor Next ](/developer-guides/clients-and-developer-tools/cursor)
