# OpenClaw

> **Source:** https://docs.qwencloud.com/developer-guides/clients-and-developer-tools/openclaw

Open-source AI assistant platform

 Copy page OpenClaw is an open-source personal AI assistant platform that lets you interact with AI through various messaging channels. You can configure it to access AI models from Qwen Cloud. It supports three access methods: pay-as-you-go, Coding Plan, and Token Plan Team Edition.
## [​ ](#install-openclaw) Install OpenClaw

OpenClaw requires Node.js 22.19.0 or later. Run the following command to check your Node.js version:
Copy ```\nnode --version

``` 
If Node.js is not installed or your version is older than 22.19.0, visit the [Node.js website](https://nodejs.org/) to download and install it.
- macOS / Linux 
- Windows 

 We recommend using the official installation script:Copy ```\ncurl -fsSL https://openclaw.ai/install.sh | bash

``` Or install globally via npm:Copy ```\nnpm install -g openclaw@latest

``` Run the following in PowerShell:Copy ```\niwr -useb https://openclaw.ai/install.ps1 | iex

``` Or install globally via npm:Copy ```\nnpm install -g openclaw@latest

``` 
After first-time installation, OpenClaw automatically starts the setup wizard. You can also manually run the `openclaw onboard` command to start the wizard.
SettingRecommendedI understand this is powerful and inherently risky. Continue?Select **Yes**Onboarding modeSelect **QuickStart**Model/auth providerSelect **Skip for now** (configure Qwen Cloud models later)Filter models by providerSelect **All providers**Default modelSelect **Keep current**Select channel (QuickStart)Select **Skip for now** (configure channels later)Configure skills now? (recommended)Select **No**Enable hooks?Press Space to select options, then press Enter to proceedHow do you want to hatch your bot?Select **Do this later** 
## [​ ](#configure-access-credentials) Configure access credentials

OpenClaw provides two configuration methods: terminal (recommended) or Web UI.
 The example disables gateway authentication (`auth.mode: none`), suitable only for single-machine local use. For shared or remote access, run `openclaw doctor --fix` to enable token authentication. 
- Method 1: Terminal 
- Method 2: Web UI 

 Edit `~/.openclaw/openclaw.json` and add the configuration for your chosen plan. 
- In your terminal, run the following command to start the Web UI:


Copy ```\nopenclaw dashboard

``` The browser will automatically open the console page (usually `http://127.0.0.1:18789`).
- 
In the left menu, choose **Configuration** > **Settings** > **Advanced**, and click **Open** to open the configuration editor.


- 
Replace the `"agents": {...}` section with the configuration for your chosen plan, and replace `YOUR_API_KEY` with your API Key.


 To keep existing settings, do not overwrite the entire file. Merge the `models` and `agents` sections carefully. 
### [​ ](#token-plan-team-edition) Token Plan Team Edition

Replace `YOUR_API_KEY` with the Token Plan Team Edition dedicated [API Key](https://home.qwencloud.com/api-keys). For available models, see Token Plan Team Edition [supported models](/token-plan/overview#supported-models).
ParameterValue**API Key**Token Plan Team Edition dedicated [API Key](https://home.qwencloud.com/api-keys)**Base URL**`https://token-plan.ap-southeast-1.maas.aliyuncs.com/apps/anthropic/v1`**Available models**Token Plan Team Edition [supported models](/token-plan/overview#supported-models) 
Copy ```\n"models": {
 "mode": "merge",
 "providers": {
 "bailian-token-plan": {
 "baseUrl": "https://token-plan.ap-southeast-1.maas.aliyuncs.com/apps/anthropic/v1",
 "apiKey": "YOUR_API_KEY",
 "api": "anthropic-messages",
 "models": [
 {
 "id": "qwen3.7-max",
 "name": "qwen3.7-max",
 "reasoning": false,
 "input": ["text"],
 "contextWindow": 1000000,
 "maxTokens": 65536,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 },
 "compat": { "thinkingFormat": "openai" }
 },
 {
 "id": "qwen3.7-plus",
 "name": "qwen3.7-plus",
 "reasoning": false,
 "input": ["text", "image"],
 "contextWindow": 1000000,
 "maxTokens": 65536,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 },
 "compat": { "thinkingFormat": "openai" }
 },
 {
 "id": "qwen3.6-plus",
 "name": "qwen3.6-plus",
 "reasoning": false,
 "input": ["text", "image"],
 "contextWindow": 1000000,
 "maxTokens": 65536,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 },
 "compat": { "thinkingFormat": "openai" }
 },
 {
 "id": "qwen3.6-flash",
 "name": "qwen3.6-flash",
 "reasoning": false,
 "input": ["text", "image"],
 "contextWindow": 1000000,
 "maxTokens": 32768,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 },
 "compat": { "thinkingFormat": "openai" }
 },
 {
 "id": "deepseek-v4-pro",
 "name": "deepseek-v4-pro",
 "reasoning": false,
 "input": ["text"],
 "contextWindow": 163840,
 "maxTokens": 32768,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 }
 },
 {
 "id": "deepseek-v4-flash",
 "name": "deepseek-v4-flash",
 "reasoning": false,
 "input": ["text"],
 "contextWindow": 163840,
 "maxTokens": 16384,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 }
 },
 {
 "id": "deepseek-v3.2",
 "name": "deepseek-v3.2",
 "reasoning": false,
 "input": ["text"],
 "contextWindow": 163840,
 "maxTokens": 16384,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 },
 "compat": { "thinkingFormat": "openai" }
 },
 {
 "id": "kimi-k2.7-code",
 "name": "kimi-k2.7-code",
 "reasoning": false,
 "input": ["text", "image"],
 "contextWindow": 262144,
 "maxTokens": 32768,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 },
 "compat": { "thinkingFormat": "openai" }
 },
 {
 "id": "kimi-k2.6",
 "name": "kimi-k2.6",
 "reasoning": false,
 "input": ["text", "image"],
 "contextWindow": 262144,
 "maxTokens": 16384,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 },
 "compat": { "thinkingFormat": "openai" }
 },
 {
 "id": "kimi-k2.5",
 "name": "kimi-k2.5",
 "reasoning": false,
 "input": ["text", "image"],
 "contextWindow": 262144,
 "maxTokens": 16384,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 },
 "compat": { "thinkingFormat": "openai" }
 },
 {
 "id": "glm-5.2",
 "name": "glm-5.2",
 "reasoning": false,
 "input": ["text"],
 "contextWindow": 1000000,
 "maxTokens": 16384,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 },
 "compat": { "thinkingFormat": "openai" }
 },
 {
 "id": "glm-5.1",
 "name": "glm-5.1",
 "reasoning": false,
 "input": ["text"],
 "contextWindow": 202752,
 "maxTokens": 16384,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 },
 "compat": { "thinkingFormat": "openai" }
 },
 {
 "id": "glm-5",
 "name": "glm-5",
 "reasoning": false,
 "input": ["text"],
 "contextWindow": 202752,
 "maxTokens": 16384,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 },
 "compat": { "thinkingFormat": "openai" }
 },
 {
 "id": "MiniMax-M2.5",
 "name": "MiniMax-M2.5",
 "reasoning": false,
 "input": ["text"],
 "contextWindow": 204800,
 "maxTokens": 131072,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 }
 }
 ]
 }
 }
},
"agents": {
 "defaults": {
 "model": {
 "primary": "bailian-token-plan/qwen3.7-plus"
 },
 "models": {
 "bailian-token-plan/qwen3.7-max": {},
 "bailian-token-plan/qwen3.7-plus": {},
 "bailian-token-plan/qwen3.6-plus": {},
 "bailian-token-plan/qwen3.6-flash": {},
 "bailian-token-plan/deepseek-v4-pro": {},
 "bailian-token-plan/deepseek-v4-flash": {},
 "bailian-token-plan/deepseek-v3.2": {},
 "bailian-token-plan/kimi-k2.7-code": {},
 "bailian-token-plan/kimi-k2.6": {},
 "bailian-token-plan/kimi-k2.5": {},
 "bailian-token-plan/glm-5.2": {},
 "bailian-token-plan/glm-5.1": {},
 "bailian-token-plan/glm-5": {},
 "bailian-token-plan/MiniMax-M2.5": {}
 }
 }
},

``` 
 To add more models, add model definitions in `providers.bailian-token-plan.models`, and add `agents.defaults.models` entries in `"bailian-token-plan/model-ID": {}`. 
### [​ ](#coding-plan) Coding Plan

Replace `YOUR_API_KEY` with the Coding Plan dedicated [API Key](https://home.qwencloud.com/api-keys). For available models, see Coding Plan [supported models](/coding-plan/overview#plan-details).
ParameterValue**API Key**Coding Plan dedicated [API Key](https://home.qwencloud.com/api-keys)**Base URL**`https://coding-intl.dashscope.aliyuncs.com/apps/anthropic/v1`**Available models**Coding Plan [supported models](/coding-plan/overview#plan-details) 
Copy ```\n"models": {
 "mode": "merge",
 "providers": {
 "bailian-coding-plan": {
 "baseUrl": "https://coding-intl.dashscope.aliyuncs.com/apps/anthropic/v1",
 "apiKey": "YOUR_API_KEY",
 "api": "anthropic-messages",
 "models": [
 {
 "id": "qwen3.7-plus",
 "name": "qwen3.7-plus",
 "reasoning": false,
 "input": ["text", "image"],
 "contextWindow": 1000000,
 "maxTokens": 65536,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 },
 "compat": { "thinkingFormat": "openai" }
 },
 {
 "id": "qwen3.6-plus",
 "name": "qwen3.6-plus",
 "reasoning": false,
 "input": ["text", "image"],
 "contextWindow": 1000000,
 "maxTokens": 65536,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 },
 "compat": { "thinkingFormat": "openai" }
 },
 {
 "id": "qwen3-coder-plus",
 "name": "qwen3-coder-plus",
 "reasoning": false,
 "input": ["text"],
 "contextWindow": 131072,
 "maxTokens": 65536,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 }
 }
 ]
 }
 }
},
"agents": {
 "defaults": {
 "model": {
 "primary": "bailian-coding-plan/qwen3.7-plus"
 },
 "models": {
 "bailian-coding-plan/qwen3.7-plus": {},
 "bailian-coding-plan/qwen3.6-plus": {},
 "bailian-coding-plan/qwen3-coder-plus": {}
 }
 }
},

``` 
### [​ ](#pay-as-you-go) Pay-as-you-go

Replace `YOUR_API_KEY` with your [Qwen Cloud API Key](https://home.qwencloud.com/api-keys). For available models, see [supported models](/developer-guides/getting-started/text-generation-models).
ParameterValue**API Key**Your [Qwen Cloud API Key](https://home.qwencloud.com/api-keys)**Base URL**`https://dashscope-intl.aliyuncs.com/apps/anthropic/v1`**Available models**[Supported models](/developer-guides/getting-started/text-generation-models) 
Copy ```\n"models": {
 "mode": "merge",
 "providers": {
 "qwencloud": {
 "baseUrl": "https://dashscope-intl.aliyuncs.com/apps/anthropic/v1",
 "apiKey": "YOUR_API_KEY",
 "api": "anthropic-messages",
 "models": [
 {
 "id": "qwen3.7-plus",
 "name": "qwen3.7-plus",
 "reasoning": false,
 "input": ["text", "image"],
 "contextWindow": 1000000,
 "maxTokens": 65536,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 },
 "compat": { "thinkingFormat": "openai" }
 },
 {
 "id": "qwen3.6-plus",
 "name": "qwen3.6-plus",
 "reasoning": false,
 "input": ["text", "image"],
 "contextWindow": 1000000,
 "maxTokens": 65536,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 },
 "compat": { "thinkingFormat": "openai" }
 },
 {
 "id": "MiniMax-M2.5",
 "name": "MiniMax-M2.5",
 "reasoning": false,
 "input": ["text"],
 "contextWindow": 204800,
 "maxTokens": 131072,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 }
 },
 {
 "id": "glm-5",
 "name": "glm-5",
 "reasoning": false,
 "input": ["text"],
 "contextWindow": 202752,
 "maxTokens": 16384,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 },
 "compat": { "thinkingFormat": "openai" }
 },
 {
 "id": "deepseek-v3.2",
 "name": "deepseek-v3.2",
 "reasoning": false,
 "input": ["text"],
 "contextWindow": 163840,
 "maxTokens": 16384,
 "cost": { "input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0 },
 "compat": { "thinkingFormat": "openai" }
 }
 ]
 }
 }
},
"agents": {
 "defaults": {
 "model": {
 "primary": "qwencloud/qwen3.7-plus"
 },
 "models": {
 "qwencloud/qwen3.7-plus": {},
 "qwencloud/qwen3.6-plus": {},
 "qwencloud/MiniMax-M2.5": {},
 "qwencloud/glm-5": {},
 "qwencloud/deepseek-v3.2": {}
 }
 }
},

``` 
## [​ ](#verify-configuration) Verify configuration

### [​ ](#enable-the-local-gateway) Enable the local gateway

Run the following command to enable the local gateway:
Copy ```\nopenclaw config set gateway.mode local

``` 
### [​ ](#start-and-test) Start and test

Copy ```\nopenclaw gateway # Start the gateway
openclaw tui # Open terminal UI

``` 
Type a message to confirm the model responds normally. If it does, the configuration is complete.
## [​ ](#common-commands) Common commands

CommandDescriptionExample`/help`Display available commands`/help``/status`View current model, session, and gateway status`/status``/model <model_name>`Switch the model for the current session`/model qwen3.7-max``/new`Start a new session`/new``/compact`Compress conversation history to free up context window`/compact``/think <level>`Set thinking (reasoning) depth: off, low, medium, high`/think high``/skills`Display all available Skills`/skills` 
## [​ ](#faq) FAQ

### [​ ](#error-codes) Error codes

If you encounter errors during configuration, refer to the FAQ documentation for the corresponding billing plan:

- Pay-as-you-go: [Error codes and troubleshooting](/api-reference/preparation/error-messages)

- Coding Plan: [Coding Plan FAQ](/coding-plan/faq)

- Token Plan Team Edition: [Token Plan Team Edition FAQ](/token-plan/faq)


### [​ ](#gateway-start-blocked-error) "Gateway start blocked" error

Run `openclaw config set gateway.mode local` to enable the local gateway, or pass the `--allow-unconfigured` flag: `openclaw gateway --allow-unconfigured`.
### [​ ](#how-to-view-configured-models) How to view configured models

Run `openclaw tui` to open the terminal UI, then type `/model` to view the model list. Press Enter to select a model, or press Esc to exit.
### [​ ](#safely-adding-new-plan-models-without-losing-existing-configurations) Safely adding new plan models without losing existing configurations

Do not overwrite the entire `~/.openclaw/openclaw.json` file. Make partial modifications instead — only edit the fields that need to change and keep existing configuration intact. After saving, restart the gateway with `openclaw gateway restart`.
### [​ ](#device-identity-required-error) "Device identity required" error

Error message:
Copy ```\ncode=1008 reason=device identity required

``` 
The client did not provide device identity information when connecting to the gateway. This usually happens on first browser visit, after clearing browser cache, or when key files in `~/.openclaw/identity/` are missing.
**Solution**: Run the following commands to approve the device and regenerate the browser access URL:
Copy ```\nopenclaw devices approve --latest
openclaw dashboard --no-open

``` 
If the issue persists, clear pending device records and try again:
Copy ```\nopenclaw devices clear --pending --yes
openclaw dashboard --no-open

``` 
### [​ ](#token-consumption-occurs-even-when-openclaw-is-not-actively-in-use) Token consumption occurs even when OpenClaw is not actively in use

**Cause:** OpenClaw has a built-in heartbeat mechanism. While the gateway is running, it automatically calls the configured model at a fixed interval (every 30 minutes by default) to check for pending tasks. Each heartbeat consumes a small number of tokens.
**How to verify:** Check the session log files (.jsonl) in the `~/.openclaw/agents/main/sessions/` directory. Heartbeat calls are marked with `[OpenClaw heartbeat poll]`.
**Solutions:**

- **Stop the gateway**: Run `openclaw gateway stop` when not in use. This stops the heartbeat immediately.

- **Increase the heartbeat interval**: Set `agents.defaults.heartbeat.every` in `~/.openclaw/openclaw.json`. For example, `"2h"` sets the interval to every 2 hours.


 [Previous ](/developer-guides/administration/rate-limits)[Hermes Agent Terminal AI coding assistant by Nous Research Next ](/developer-guides/clients-and-developer-tools/hermes-agent)
