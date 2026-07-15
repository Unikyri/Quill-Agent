# Text generation models

> **Source:** https://docs.qwencloud.com/developer-guides/getting-started/text-generation-models

Choose a model for AI agents, chatbots, document processing, and more.

 Copy page ## [​ ](#using-openclaw-or-claude-code) Using OpenClaw or Claude Code?

`qwen3.7-plus` — balanced performance and cost, full tool support, 1M context for large codebases. For strongest reasoning, choose `qwen3.7-max`.
## [​ ](#migrate-from-closed-source-models) Migrate from closed-source models

Map your current GPT, Claude, or Gemini model to an equivalent Qwen Cloud model.
TierClosed-source examplesQwen Cloud recommendationHighest capabilityGPT-5.5, Claude Opus 4.7, Gemini 3.1 Pro`qwen3.7-max`BalancedGPT-5.4, Claude Sonnet 4.6, Gemini 3 Pro`qwen3.7-plus`, `deepseek-v4-pro`Lightweight & low-costGPT-5.4-mini, Claude Haiku 4.5, Gemini 3.1 Flash`qwen3.6-flash`, `deepseek-v4-flash` 
## [​ ](#for-other-applications) For other applications

Chatbots, content generation, summarization, document processing — start with `qwen3.7-plus` for strongest accuracy, 1M context, and the full feature set. Once your use case works well, try `qwen3.6-flash` to reduce cost — near-flagship quality with the same context and features.
### [​ ](#context-window) Context window

1M tokens is roughly 750,000 words or 10 novels.

- Long documents or large codebases → `qwen3.7-max` / `qwen3.7-plus` / `qwen3.6-flash` (1M)

- Standard tasks → 128k–256k is plenty


### [​ ](#thinking-mode) Thinking mode

Step-by-step reasoning for multi-step math, debugging, architecture planning, or legal cross-referencing.
Toggle with `enable_thinking`. All Qwen3+ models support it — most are hybrid, so you can switch per request.
### [​ ](#function-calling-built-in-tools) Function calling + built-in tools

Let the model take actions: check weather, query a database, book a meeting.

- Function calling (you define tools, model calls them): all general-purpose models

- Built-in tools (web search, code execution — no setup): `qwen3.7-max`, `qwen3.7-plus`, `qwen3.6-flash`, `qwen3.5-plus`, `qwen3.5-flash`, `qwen3-max` series only


## [​ ](#recommended-models) Recommended models

ModelContextThinkingFunction callingBuilt-in toolsStructured output`qwen3.7-max`1M✓✓✓—`qwen3.7-plus`1M✓✓✓✓`qwen3.6-flash`1M✓✓✓✓`deepseek-v4-pro`1M✓✓——`deepseek-v4-flash`1M✓✓—— 
## [​ ](#all-models) All models

 Qwen3.7

 Model IDContextMax OutputThinking BudgetFunction callingBuilt-in toolsStructured output`qwen3.7-max`1M64k256k✓✓—`qwen3.7-max-2026-06-08`1M64k256k✓✓—`qwen3.7-max-2026-05-20`1M64k256k✓✓—`qwen3.7-max-preview`1M64k256k✓✓—`qwen3.7-max-2026-05-17`1M64k256k✓✓—`qwen3.7-plus`1M64k256k✓✓✓`qwen3.7-plus-2026-05-26`1M64k256k✓✓✓ Qwen3.6

 Model IDContextMax OutputThinking BudgetFunction callingBuilt-in toolsStructured output`qwen3.6-max-preview`256k64k128k✓—✓`qwen3.6-flash`1M64k80k✓✓✓`qwen3.6-flash-2026-04-16`1M64k80k✓✓✓`qwen3.6-35b-a3b`256k64k80k✓✓✓`qwen3.6-27b`256k64k80k✓—✓`qwen3.6-plus`1M64k80k✓✓✓`qwen3.6-plus-2026-04-02`1M64k80k✓✓✓ Qwen3.5

 Model IDContextMax OutputThinking BudgetFunction callingBuilt-in toolsStructured output`qwen3.5-plus`1M64k80k✓✓✓`qwen3.5-plus-2026-04-20`1M64k80k✓✓✓`qwen3.5-plus-2026-02-15`1M64k80k✓✓✓`qwen3.5-flash`1M64k80k✓✓✓`qwen3.5-flash-2026-02-23`1M64k80k✓✓✓`qwen3.5-397b-a17b`256k64k80k✓✓✓`qwen3.5-122b-a10b`256k64k80k✓✓✓`qwen3.5-27b`256k64k80k✓✓✓`qwen3.5-35b-a3b`256k64k80k✓✓✓ Specialized

 ### Translation
Model IDContextMax OutputFunction callingBuilt-in toolsStructured output`qwen-mt-plus`16k8k———`qwen-mt-turbo`16k8k———`qwen-mt-flash`16k8k———`qwen-mt-lite`16k8k——— ### Character roleplay
Model IDContextMax OutputFunction callingBuilt-in toolsStructured output`qwen-plus-character`32k4k———`qwen-plus-character-ja`8k4k———`qwen-flash-character`8k4k——— Third-party

 Non-Qwen models available through the same API.Model IDContextMax OutputThinking BudgetFunction callingBuilt-in toolsStructured output`deepseek-v4-pro`1M384k **✓——`deepseek-v4-flash`1M384k **✓——`deepseek-v3.2`128k64k32k✓—— * DeepSeek V4 models share a 384k total budget across output and thinking. Legacy

 Previous generation models. We recommend Qwen3.6 for new projects.### Qwen3
Model IDContextMax OutputThinking BudgetFunction callingBuilt-in toolsStructured output`qwen3-max`256k64k80k✓✓✓`qwen3-max-2026-01-23`256k64k80k✓✓✓`qwen3-max-preview`256k64k80k✓✓✓`qwen3-max-2025-09-23`256k64k—✓✓✓`qwen3-235b-a22b`128k16k38k✓—✓`qwen3-235b-a22b-thinking-2507`128k32k80k✓——`qwen3-235b-a22b-instruct-2507`128k32k—✓—✓`qwen3-next-80b-a3b-thinking`128k32k80k✓——`qwen3-next-80b-a3b-instruct`128k32k—✓—✓`qwen3-32b`128k16k38k✓—✓`qwen3-30b-a3b`128k16k38k✓—✓`qwen3-30b-a3b-thinking-2507`128k32k80k✓——`qwen3-30b-a3b-instruct-2507`128k32k—✓—✓`qwen3-14b`128k8k38k✓—✓`qwen3-8b`128k8k38k✓—✓`qwen3-4b`128k8k38k✓—✓`qwen3-1.7b`32k8k30k✓—✓`qwen3-0.6b`32k8k30k✓—✓ ### Qwen3-Coder
Model IDContextMax OutputFunction callingBuilt-in toolsStructured output`qwen3-coder-plus`1M64k✓—✓`qwen3-coder-plus-2025-09-23`1M64k✓—✓`qwen3-coder-plus-2025-07-22`1M64k✓—✓`qwen3-coder-flash`1M64k✓—✓`qwen3-coder-flash-2025-07-28`1M64k✓—✓`qwen3-coder-next`256k64k✓—✓`qwen3-coder-480b-a35b-instruct`256k64k✓—✓`qwen3-coder-30b-a3b-instruct`256k64k✓—✓ ### Qwen2.5 (open source)
Model IDContextMax OutputFunction callingBuilt-in toolsStructured output`qwen2.5-omni-7b`32k8k✓—✓`qwen2.5-vl-72b-instruct`128k8k✓—✓`qwen2.5-vl-32b-instruct`128k8k✓—✓`qwen2.5-vl-7b-instruct`128k8k✓—✓`qwen2.5-vl-3b-instruct`128k8k✓—✓`qwen2.5-72b-instruct`32k8k✓—✓`qwen2.5-32b-instruct`32k8k✓—✓`qwen2.5-14b-instruct`32k8k✓—✓`qwen2.5-14b-instruct-1m`1M8k✓—✓`qwen2.5-7b-instruct`32k8k✓—✓`qwen2.5-7b-instruct-1m`1M8k✓—✓ ### Legacy (qwen-plus/max/flash/turbo)
Model IDContextMax OutputThinking BudgetFunction callingBuilt-in toolsStructured output`qwen-plus`1M32k80k✓—✓`qwen-plus-latest`1M32k80k✓—✓`qwen-plus-2025-12-01`1M32k80k✓—✓`qwen-plus-2025-09-11`1M32k80k✓—✓`qwen-plus-2025-07-28`1M32k80k✓—✓`qwen-plus-2025-07-14`128k16k80k✓—✓`qwen-plus-2025-04-28`128k16k80k✓—✓`qwen-plus-2025-01-25`128k8k—✓—✓`qwen-max`32k8k—✓—✓`qwen-max-latest`32k8k—✓—✓`qwen-max-2025-01-25`32k8k—✓—✓`qwen-flash`1M32k80k✓—✓`qwen-flash-2025-07-28`1M32k80k✓—✓`qwen-turbo`128k16k38k✓—✓`qwen-turbo-latest`128k16k38k✓—✓`qwen-turbo-2025-04-28`128k16k38k✓—✓`qwen-turbo-2024-11-01`1M8k—✓—✓`qwq-plus`128k8k32k———`qvq-max`128k8k80k———`qvq-max-latest`128k8k80k———`qvq-max-2025-03-25`128k8k80k———`qwen-omni-turbo`32k2k80k———`qwen-omni-turbo-latest`32k2k80k———`qwen-omni-turbo-2025-03-26`32k2k80k——— 

---

## [​ ](#learn-more) Learn more

[ ## Model selection guide
Not looking for text generation models? Start here. ](/developer-guides/getting-started/model-selection)[ ## Try free
Try models in the browser — no API key needed. ](https://home.qwencloud.com/try-ai) [Previous ](/developer-guides/getting-started/pricing)[Generate text Make your first text generation call Next ](/developer-guides/text-generation/quickstart)
