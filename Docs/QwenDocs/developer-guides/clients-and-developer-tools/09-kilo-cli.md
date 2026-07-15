# Kilo CLI

> **Source:** https://docs.qwencloud.com/developer-guides/clients-and-developer-tools/kilo-cli

AI coding in terminal

 Copy page Kilo CLI brings AI-powered coding assistance directly to your terminal. Connect it to Qwen Cloud via Token Plan (Team Edition), Coding Plan, or pay-as-you-go billing.
## [​ ](#quick-start) Quick start

Get running in a few minutes:
Copy ```\n# 1. Install (requires Node.js v18+)
npm install -g @kilocode/cli
kilo --version

# 2. Configure: Edit ~/.config/kilo/config.json
# Add provider configuration with your API key

# 3. Test (run kilo and type)
kilo
Ask: "Write a function to calculate fibonacci numbers"

``` 
You should see: Kilo generates the fibonacci function code
## [​ ](#configuration) Configuration

Open `~/.config/kilo/config.json` with a text editor and add the configuration for your chosen plan.
 **Free quota and billing:**
- First-time users get a free quota. See [Free quota](/resources/free-quota) for details.

- Enable [Free quota only](https://home.qwencloud.com/benefits) to prevent unexpected charges.


 
### [​ ](#token-plan-team-edition) Token Plan (Team Edition)

You must first purchase a Token Plan (Team Edition) with an active subscription on the [Token Plan page](https://home.qwencloud.com/billing/subscription/token-plan).
Replace `YOUR_API_KEY` with the Token Plan (Team Edition) dedicated [API Key](https://home.qwencloud.com/api-keys). For available models, see Token Plan (Team Edition) [supported models](/token-plan/overview#supported-models).
Copy ```\n{
 "$schema": "https://kilo.ai/config.json",
 "provider": {
 "bailian-token-plan": {
 "npm": "@ai-sdk/openai-compatible",
 "name": "Qwen Cloud",
 "options": {
 "baseURL": "https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1",
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
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 8192
 }
 }
 },
 "kimi-k2.6": {
 "name": "Kimi K2.6",
 "options": {
 "thinking": {
 "type": "enabled",
 "budgetTokens": 8192
 }
 }
 },
 "kimi-k2.5": {
 "name": "Kimi K2.5",
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

Replace `YOUR_API_KEY` with the Coding Plan dedicated [API Key](https://home.qwencloud.com/api-keys). For available models, see Coding Plan [supported models](/coding-plan/overview#plan-details).
Copy ```\n{
 "$schema": "https://kilo.ai/config.json",
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
 }
 }
 }
 }
}

``` 
### [​ ](#pay-as-you-go) Pay-as-you-go

Replace `YOUR_API_KEY` with your [Qwen Cloud API Key](https://home.qwencloud.com/api-keys). For available models, see [supported models](/developer-guides/getting-started/text-generation-models).
Copy ```\n{
 "$schema": "https://kilo.ai/config.json",
 "provider": {
 "qwencloud": {
 "npm": "@ai-sdk/openai-compatible",
 "name": "Qwen Cloud",
 "options": {
 "baseURL": "https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
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
 }
 }
 }
 }
}

``` 
To add more models, append them in the same format within `models`.
### [​ ](#model-recommendations) Model recommendations

TaskModelWhy**Simple tasks**`qwen3-coder-plus`Fast responses, low cost**Standard coding**`qwen3-coder-plus`Balanced performance**Complex algorithms**`qwen3.6-plus`Strong reasoning**Architecture design**`qwen3.6-plus`Deep code understanding 
## [​ ](#verify-configuration) Verify configuration

After saving the configuration, restart Kilo CLI, type `/models`, search for "Qwen Cloud", and select the model you want to use.
For more tips and common commands, see the [Kilo Code official documentation](https://kilo.ai/docs).
## [​ ](#limitations) Limitations


- Terminal only: No GUI interface

- Token usage: Multi-file edits consume more tokens

- Model compatibility: Not all features work with all models


## [​ ](#troubleshooting) Troubleshooting

**"Invalid API key" error**

Solution:

- Verify API key is correct and matches your plan

- Ensure API key has quota or active subscription


**"Model not found" error**

Solution:

- Check model ID spelling

- See [Model list](/developer-guides/getting-started/model-selection)


**High token consumption**

Solution:

- Work in specific directories

- Use precise prompts

- Clear context with new sessions

- Choose appropriate models for tasks


**Slow responses**

Solution:

- Use faster models like `qwen3-coder-next`

- Check network connection

- Reduce context size


## [​ ](#token-optimization) Token optimization

Reduce usage by:

- **Focused directories**: Navigate to specific project folders

- **Clear prompts**: Be specific about what you need

- **Model selection**: Use lightweight models for simple tasks

- **Context management**: Start new sessions for unrelated tasks

- **Incremental changes**: Make small, focused edits


## [​ ](#faq) FAQ

If you encounter errors, refer to the FAQ documentation for the corresponding billing plan:

- Pay-as-you-go: [Error codes and troubleshooting](/api-reference/preparation/error-messages)

- Coding Plan: [Coding Plan FAQ](/coding-plan/faq)

- Token Plan (Team Edition): [Token Plan FAQ](/token-plan/faq)


## [​ ](#related-resources) Related resources


- **Token Plan**: [Setup with subscription](/token-plan/quickstart)

- **Coding Plan**: [Setup with subscription](/coding-plan/overview)

- **Models**: [Available models](/developer-guides/getting-started/text-generation-models)

- **API docs**: [OpenAI-compatible reference](/api-reference/chat/openai-chat)

- **Official docs**: [Kilo CLI documentation](https://kilo.ai/docs)


 [Previous ](/developer-guides/clients-and-developer-tools/lingma)[Postman API testing tool Next ](/developer-guides/clients-and-developer-tools/postman)
