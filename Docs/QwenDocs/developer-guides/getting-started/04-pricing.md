# Pricing

> **Source:** https://docs.qwencloud.com/developer-guides/getting-started/pricing

Pay-as-you-go pricing for API usage

 Copy page **Pricing coverage**This page lists pricing for selected representative models only. Qwen Cloud supports many more models than shown here. For models not listed in this documentation, see the Model Marketplace for complete and up-to-date pricing.
- Model Marketplace: [https://www.qwencloud.com/models](https://www.qwencloud.com/models)

- Specific model pricing page: `https://www.qwencloud.com/models/{model-id}`

Example: qwen3.7-plus → [https://www.qwencloud.com/models/qwen3.7-plus](https://www.qwencloud.com/models/qwen3.7-plus)

- Encode `/` as `%2F` in model IDs, e.g. siliconflow/deepseek-v3.1-terminus → [https://www.qwencloud.com/models/siliconflow%2Fdeepseek-v3.1-terminus](https://www.qwencloud.com/models/siliconflow%2Fdeepseek-v3.1-terminus)


 
 New users get a free quota to try models at no cost. See [Free quota](/resources/free-quota) for details. 
Billing varies by model type: text models charge per token, image generation per image, video generation per second, and speech models per character or per second of audio. Failed API calls are not charged and do not consume free quota.
 The prices listed below are list prices. For current promotions and discounted pricing, visit the [Model Marketplace](https://www.qwencloud.com/models). 
## [​ ](#text-generation) Text generation

Billed per million tokens. Models with long-context support use tiered pricing — the more input tokens in a single request, the higher the per-token rate.
ModelInput per requestInputOutputqwen3.7-max0 – 991K$2.50$7.50qwen3.6-max-preview≤ 128K$1.30$7.80128K – 256K$2.00$12.00qwen3.7-plus≤ 256K$0.40$1.60256K – 1M$1.20$4.80qwen3.6-flash≤ 256K$0.25$1.50256K – 1M$1.00$4.00 
For complete text model pricing, see [Model Marketplace](https://www.qwencloud.com/models).
## [​ ](#images-videos) Images & videos

### [​ ](#understanding) Understanding

Vision understanding is billed per token. Qwen text generation models (qwen3.7-plus, etc.) support vision input at the same token price listed above. Dedicated vision models have separate pricing:
ModelInput per requestInputOutputqwen3-vl-plus≤ 32K$0.20$1.6032K – 128K$0.30$2.40128K – 256K$0.60$4.80qwen3-vl-flash≤ 32K$0.05$0.4032K – 128K$0.075$0.60128K – 256K$0.12$0.96 
Image and video inputs are automatically converted to tokens. The conversion varies by model:
Model familyImage conversionExample (1024×1024)Qwen (qwen3.7-plus, etc.)1 token per 32×32 pixels≈ 256 tokensQwen-VL (qwen3-vl, etc.)1 token per 32×32 pixels≈ 256 tokensQwen3.5-Omni1 token per 32×32 pixels≈ 256 tokensQwen3-Omni-Flash1 token per 32×32 pixels≈ 256 tokens 
Video tokens = sampled frames × tokens per frame. See [Token counting →](/developer-guides/run-and-scale/token-counting) for details.
### [​ ](#generation) Generation

Image generation is billed per image (resolution-independent). Video generation is billed per second of output video.
**Image generation**
ModelPrice per imageqwen-image-2.0-pro$0.075qwen-image-2.0$0.035qwen-image-edit$0.045wan2.6-t2i$0.03wan2.6-image$0.03z-image-turbo$0.015 (prompt rewrite off) / $0.03 (on) 
**Video generation**
ModelPrice per secondwan2.6-t2v$0.10wan2.6-i2v$0.10wan2.6-i2v-flash$0.05 
For all image and video model pricing, see [Model Marketplace](https://www.qwencloud.com/models).
## [​ ](#audio-speech) Audio & speech

### [​ ](#text-to-speech) Text to speech

Billed per 10,000 characters of input text.
ModelPrice per 10K charscosyvoice-v3-plus$0.26cosyvoice-v3-flash$0.13qwen3-tts-flash$0.10 
### [​ ](#speech-to-text) Speech to text

Billed per second of audio input.
ModelPrice per secondfun-asr$0.000035fun-asr-realtime$0.00009qwen3-asr-flash$0.000035 
### [​ ](#speech-to-speech) Speech to speech

 Qwen-Omni is a multimodal model that handles text, audio, and image/video in a single call. All modality prices are listed in the table below. 
Billed per million tokens, with different rates per modality.
**Token conversion**
Input typeConversion rateTextStandard tokenizerAudio input≈ 7 tokens/sec (Qwen3.5-Omni) or 12.5 tokens/sec (Qwen3-Omni-Flash) or 25 tokens/sec (Qwen-Omni-Turbo)Audio output≈ 12.5 tokens/sec (Qwen3.5-Omni) or 12.5 tokens/sec (Qwen3-Omni-Flash)Image/VideoSee [Understanding](#understanding) section above 
**qwen3.5-omni pricing**
Price per 1M tokens:
ModelText/Image/Video inputAudio inputText outputText + Audio outputqwen3.5-omni-plus$1.4$11$8.3$44qwen3.5-omni-flash$0.4$3$2.2$11.9 
For all speech model pricing, see [Model Marketplace](https://www.qwencloud.com/models).
## [​ ](#embedding-reranking) Embedding & reranking

Billed per million input tokens (output is not charged). Multimodal embedding models may charge different rates for image vs text input. Image/video token conversion for embedding models is handled internally — check the `usage` field in the API response for actual token counts.
ModelModalityPrice per 1M tokenstext-embedding-v4Text$0.07tongyi-embedding-vision-plusAll$0.09tongyi-embedding-vision-flashImage/Video$0.03Text$0.09qwen3-rerankText$0.10 
For all embedding and reranking model pricing, see [Model Marketplace](https://www.qwencloud.com/models).
## [​ ](#built-in-tools) Built-in tools

Some built-in tools incur per-call fees in addition to model token costs.
ToolFeeNotes[Web Search](/developer-guides/text-generation/web-search)$10 / 1K calls[Web Extractor](/developer-guides/text-generation/web-scraping)FREELimited time[Code Interpreter](/developer-guides/text-generation/code-interpreter)FREELimited time[Image Search](/developer-guides/text-generation/image-search)$8 / 1K callsText-to-image and image-to-image 
[Function calling](/developer-guides/text-generation/function-calling) and [MCP](/developer-guides/text-generation/mcp) have no tool fees — tool descriptions count as input tokens.
## [​ ](#free-quota) Free quota

New users get free quota upon sign-up, typically valid for 90 days. Applies to real-time API calls only. [Learn more →](/resources/free-quota)
## [​ ](#save-on-costs) Save on costs


- **Batch API** — 50% off for async workloads. [Learn more →](/developer-guides/text-generation/batch)

- **Context caching** — Reuse long prompts at reduced cost. [Learn more →](/developer-guides/text-generation/context-cache)

- **Model selection** — Match model tier to task complexity. [Compare models →](https://www.qwencloud.com/models)


 Batch and cache discounts cannot be combined on the same request. 
For worked examples and advanced strategies, see [Cost optimization →](/developer-guides/run-and-scale/cost-optimization).
## [​ ](#learn-more) Learn more


- [Model Marketplace](https://www.qwencloud.com/models) — Complete pricing for all models

- [Free quota](/resources/free-quota) — Eligibility and activation

- [Cost optimization](/developer-guides/run-and-scale/cost-optimization) — Advanced strategies

- [Token Plan](/token-plan/overview) — Credits-based pricing for AI coding tools

- [Coding Plan](/coding-plan/overview) — Fixed monthly pricing for AI coding tools

- [Billing FAQ](/resources/faq-billing) — Common questions

- [Bill management](/resources/bill-query) — View usage and invoices


 [Previous ](/developer-guides/getting-started/model-selection)[Text generation models Choose a model for AI agents, chatbots, document processing, and more. Next ](/developer-guides/getting-started/text-generation-models)
