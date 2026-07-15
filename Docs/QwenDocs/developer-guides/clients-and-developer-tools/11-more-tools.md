# More tools

> **Source:** https://docs.qwencloud.com/developer-guides/clients-and-developer-tools/more-tools

Connect any OpenAI or Anthropic compatible programming tool

 Copy page Qwen Cloud supports any third-party programming tool compatible with OpenAI or Anthropic API protocols that allows custom endpoints. Connect via pay-as-you-go, Coding Plan, or Token Plan (Team Edition).
## [​ ](#configure-credentials) Configure credentials

Qwen Cloud offers three billing plans:

- **Token Plan (Team Edition)**: Seat-based subscription. Token consumption is deducted from Credits.

- **Coding Plan**: Fixed monthly subscription billed by number of model calls.

- **Pay-as-you-go**: Post-paid based on actual usage.


### [​ ](#token-plan-team-edition) Token Plan (Team Edition)

API protocolBase URLAPI KeySupported modelsOpenAI`https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1`Token Plan (Team Edition) dedicated [API Key](https://home.qwencloud.com/api-keys)[Supported models](/token-plan/overview#supported-models) (text generation only)Anthropic`https://token-plan.ap-southeast-1.maas.aliyuncs.com/apps/anthropic` 
### [​ ](#coding-plan) Coding Plan

API protocolBase URLAPI KeySupported modelsOpenAI`https://coding-intl.dashscope.aliyuncs.com/v1`Coding Plan dedicated [API Key](https://home.qwencloud.com/api-keys)[Supported models](/coding-plan/overview#plan-details)Anthropic`https://coding-intl.dashscope.aliyuncs.com/apps/anthropic` 
### [​ ](#pay-as-you-go) Pay-as-you-go

API protocolBase URLAPI KeySupported modelsOpenAI`https://dashscope-intl.aliyuncs.com/compatible-mode/v1`Your [Qwen Cloud API Key](https://home.qwencloud.com/api-keys)[Supported models](/api-reference/chat/openai-chat)Anthropic`https://dashscope-intl.aliyuncs.com/apps/anthropic`[Supported models](/api-reference/chat/anthropic) 
## [​ ](#unsupported-tool-types) Unsupported tool types

Token Plan (Team Edition) and Coding Plan support only AI programming tools and OpenClaw-type agents. The following tools are **not supported**:

- **Workflow and automation platforms**: Dify, n8n, Coze, etc.

- **API testing tools**: Postman, Insomnia, etc.

- **Custom applications**: automated scripts, application backends calling the API directly, etc.


 Using a plan API key for calls outside the permitted scope will be considered a violation or abuse, and may result in subscription suspension or API key ban. 
## [​ ](#use-in-ides) Use in IDEs

To use Qwen Cloud in VS Code variants or JetBrains IDEs, install one of the following extensions:
ExtensionSupported IDEsDocumentation**Claude Code**VS Code, JetBrains IDEs[Claude Code](/developer-guides/clients-and-developer-tools/claude-code)**Kilo Code**VS Code, JetBrains IDEs[Kilo Code](/developer-guides/clients-and-developer-tools/kilo-cli)**Qwen Code**VS Code[Qwen Code](/developer-guides/clients-and-developer-tools/qwen-code) [Previous ](/developer-guides/clients-and-developer-tools/dify)[MLOps & observability Production AI monitoring Next ](/developer-guides/integrations/mlops-observability)
