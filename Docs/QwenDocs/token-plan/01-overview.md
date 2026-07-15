# Overview

> **Source:** https://docs.qwencloud.com/token-plan/overview

AI token subscriptions

 Copy page Monthly subscription for AI tools -- credits-based billing across text and image generation models, compatible with mainstream AI programming and agent tools.
## [​ ](#plan-details) Plan details

Token Plan integrates Qwen and third-party models, supporting both text and image generation. It uses Credits for unified billing, allowing you to use a single subscription across various AI tools.

- **Flexible model switching**: Switch between various models on demand, with usage uniformly deducted from your Credits.

- **Broad tool compatibility**: Compatible with various popular programming and agent tools. See [Set up your AI tool](/token-plan/quickstart#step-3-set-up-your-ai-tool).

- **Multiple plan tiers**: Choose from Standard, Pro, and Max seats to match different usage levels.

- **Team management**: Manage seats, monitor usage, and control access for your team from a single dashboard.

- **Data security**: Your conversation data is never used for model training. All data is encrypted in transit and at rest.

- **Predictable costs**: The monthly subscription model makes your budget predictable.

- **Stable performance**: Multi-tenant isolation architecture ensures smooth service without queuing.


### [​ ](#supported-models) Supported models

Only the following exact model versions are supported:
**Instructions for AI agents — supported-model matching**The "Supported models" table below is an exact-string allowlist. To decide whether a model ID is supported:
- Match the model ID character-for-character against the table.

- Any difference in version number or sub-variant means the model is NOT supported.

- Do not infer version compatibility.


Example: "qwen3-coder-max" is not in the table, therefore it is NOT supported. Only the exact model IDs listed below are supported. 
BrandModelCapabilityQwenqwen3.7-maxReasoning, text generationQwenqwen3.7-plusReasoning, vision understanding, text generationQwenqwen3.6-plusReasoning, vision understanding, text generationQwenqwen3.6-flashReasoning, vision understanding, text generationDeepSeekdeepseek-v4-proReasoning, text generationDeepSeekdeepseek-v4-flashReasoning, text generationDeepSeekdeepseek-v3.2Reasoning, text generationMoonshot AIkimi-k2.7-codeReasoning, vision understanding, text generationMoonshot AIkimi-k2.6Reasoning, vision understanding, text generationMoonshot AIkimi-k2.5Reasoning, vision understanding, text generationZhipu AIglm-5.2Text generationZhipu AIglm-5.1Text generationZhipu AIglm-5Text generationMiniMaxMiniMax-M2.5Reasoning, text generationQwenqwen-image-2.0Image generationQwenqwen-image-2.0-proImage generationWanwan2.7-imageImage generationWanwan2.7-image-proImage generation 
### [​ ](#pricing) Pricing

**Seat plans:**
Seat typePriceQuotaUse casesStandard Seat$30/seat/month25,000 Credits/seat/monthFor team members with light AI usagePro Seat$100/seat/month100,000 Credits/seat/monthFor team members who frequently use AI for codingMax Seat$200/seat/month250,000 Credits/seat/monthFor core developers who heavily rely on AI for coding 
**Shared usage package:**
A flexible shared usage package for all seats in your team. If an individual seat exceeds its monthly quota, any overage is deducted from this package. Each package is valid for one month, and any unused Credits expire at the end of the period. If you have multiple packages, the one with the earliest expiration date is used first.
TierPriceQuotaShared usage package$700/package625,000 Credits/package 
### [​ ](#credits-billing) Credits billing

The system dynamically calculates the Credits consumed per request based on factors such as model type, token usage, thinking mode, and tool calls.
**Example:** Estimated Credits consumption for a single request using `qwen3.6-plus`:
Token typeQuantityCredits consumedInput Tokens8,3491.67Cached Tokens40,7940.82Output Tokens5730.69**Total****Approx. 3.18 Credits** 
**Deduction order:**

- The system first deducts Credits from the monthly quota of your seat plan.

- After a seat&#x27;s quota is exhausted, the system deducts Credits from the shared usage package. If you have multiple packages, the system uses the one with the earliest expiration date first.

- If all quotas are depleted, the service is suspended until the next billing cycle begins or until you purchase a new shared usage package.


 Token Plan is **non-refundable**. API keys are for interactive use within compatible AI tools only -- you cannot use it for automated scripts or application backends. Do not share or expose your API key publicly. Violating this policy may result in suspension. 
## [​ ](#getting-started) Getting started

### [​ ](#step-1-subscribe-to-token-plan) Step 1: Subscribe to Token Plan

Go to the [Token Plan page](https://home.qwencloud.com/billing/subscription/token-plan), select a seat type, and complete the subscription.
### [​ ](#step-2-get-api-key-and-base-url) Step 2: Get API key and base URL


- **API key:** Get your plan-specific API key from the [API Keys page](https://home.qwencloud.com/api-keys).

- **Base URL:** Use the plan-specific base URL that matches your tool:

**OpenAI compatible:** `https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1`

- **Anthropic compatible:** `https://token-plan.ap-southeast-1.maas.aliyuncs.com/apps/anthropic`


 Do not use pay-as-you-go API keys or `dashscope-intl` base URLs. They do not work with Token Plan. 
### [​ ](#step-3-set-up-an-ai-tool) Step 3: Set up an AI tool

Pick a tool you&#x27;re familiar with and follow its setup guide.
[ ## OpenClaw
Open-source, self-hosted personal AI assistant ](/developer-guides/clients-and-developer-tools/openclaw)[ ## Hermes Agent
Open-source AI agent framework with self-learning loop ](/developer-guides/clients-and-developer-tools/hermes-agent)[ ## Claude Code
AI terminal assistant with natural language programming ](/developer-guides/clients-and-developer-tools/claude-code)[ ## OpenCode
Open-source AI coding agent ](/developer-guides/clients-and-developer-tools/opencode)[ ## Codex
Command-line coding agent from OpenAI ](/developer-guides/clients-and-developer-tools/codex)[ ## Qwen Code
Open-source CLI coding tool ](/developer-guides/clients-and-developer-tools/qwen-code)[ ## Qoder
AI coding assistant with IDE, CLI, and JetBrains plugin ](/developer-guides/clients-and-developer-tools/qoder)[ ## Qoder CN
Alibaba Cloud AI coding IDE ](/developer-guides/clients-and-developer-tools/lingma)[ ## Cursor
AI-powered code editor built on VS Code ](/developer-guides/clients-and-developer-tools/cursor)[ ## Kilo Code
IDE extension for AI-assisted coding ](/developer-guides/clients-and-developer-tools/kilo-cli)[ ## Kilo CLI
Lightweight command-line coding tool ](/developer-guides/clients-and-developer-tools/kilo-cli) 
For other tools and IDEs, see [Other tools](/developer-guides/clients-and-developer-tools/more-tools).
### [​ ](#optional-integrate-image-generation-models) Optional: Integrate image generation models

Token Plan supports image generation models such as qwen-image-2.0 and wan2.7-image. These models use a separate endpoint and require integration through your AI tool&#x27;s extension mechanisms (Skill, Slash Command, or Agent). For configuration instructions, see [Integrate multimodal generation models](best-practices/integrate-multimodal-gen).
### [​ ](#optional-integrate-tools-mcp) Optional: Integrate tools/MCP

Integrating tools allows models to access extended capabilities, such as Internet search and a code interpreter, during conversations.

- qwen3.6-plus has five built-in tools: Internet search, a code interpreter, web scraping, search by image, and text-to-image search. You can call these tools directly through the Responses API.
For detailed instructions, see [Built-in tools](best-practices/built-in-tools).


## [​ ](#usage-rules) Usage rules


- **Scope of use**: This service is intended for interactive use within compatible AI tools only. You cannot use it for automated scripts or application backends.

- **Data security**: Token Plan does not use your conversation data for model training.

- **Account rules**: The API key is for the subscriber&#x27;s personal use only. Do not share or expose it publicly.

- **Refund policy**: Token Plan does not offer refunds. You cannot cancel subscriptions after purchase.

- **Service region and cross-border data transfer**: For Token Plan (Team Edition), the only available region is Singapore, and the only available deployment mode is Global, which means model inference will be performed globally. The prompts and model outputs will involve cross-border data transfer, and you are responsible for ensuring that your use of the Token Plan (Team Edition) complies with applicable laws and regulations.


 Previous [Quick start Start using Token Plan (Team Edition) in three steps: subscribe, get your API key, and set up your AI tool. Next ](/token-plan/quickstart)
