# Qwen Cloud Documentation Index

> Documentación completa de Qwen Cloud guardada localmente en formato Markdown.
> Fuente original: [https://docs.qwencloud.com](https://docs.qwencloud.com)

---

## 🧭 Navegación rápida

| Sección | Descripción |
|---------|-------------|
| [🚀 Getting Started](#-getting-started) | Primeros pasos, modelos y precios |
| [📝 Developer Guide](#-developer-guide) | Guías completas de uso de modelos |
| [📖 API & SDK Reference](#-api--sdk-reference) | Referencia técnica de APIs |
| [💰 Resources](#-resources) | Facturación y cuotas gratuitas |
| [🎫 Token Plan](#-token-plan) | Planes y FAQ de tokens |
| [📋 Changelog](#-changelog) | Historial de cambios |

---

## 🚀 Getting Started

### Overview & First Steps
- [Overview / Introduction](developer-guides/getting-started/01-introduction.md)
- [First API Call](developer-guides/getting-started/02-first-api-call.md)
- [Choose Models](developer-guides/getting-started/03-model-selection.md)
- [Pricing](developer-guides/getting-started/04-pricing.md)

### Model Reference Cards
- [Text Generation Models](developer-guides/getting-started/05-text-generation-models.md)
- [Vision Models](developer-guides/getting-started/06-vision-models.md)
- [Image Models](developer-guides/getting-started/07-image-models.md)
- [Video Models](developer-guides/getting-started/08-video-models.md)

---

## 📝 Developer Guide

### 💬 Text Generation
- [Quickstart](developer-guides/text-generation/01-quickstart.md)
- [Streaming](developer-guides/text-generation/02-streaming.md)
- [Thinking](developer-guides/text-generation/03-thinking.md)
- [Structured Output](developer-guides/text-generation/04-structured-output.md)
- [Batch API](developer-guides/text-generation/05-batch.md)
- [Function Calling](developer-guides/text-generation/06-function-calling.md)
- [Multi-turn Conversations](developer-guides/text-generation/07-multi-turn.md)
- [Partial Mode](developer-guides/text-generation/08-partial-mode.md)

### 👁️ Multimodal
- [Vision (Image Understanding)](developer-guides/multimodal/01-vision.md)
- [OCR](developer-guides/multimodal/02-ocr.md)

### 🖼️ Image Generation
- [Text to Image](developer-guides/image-generation/01-text-to-image.md)
- [Image Editing](developer-guides/image-generation/02-image-editing.md)

### 🎬 Video Generation
- [Text to Video](developer-guides/video-generation/01-text-to-video.md)
- [Image to Video](developer-guides/video-generation/02-image-to-video.md)
- [Reference Video](developer-guides/video-generation/03-reference-video.md)
- [Video Editing](developer-guides/video-generation/04-video-editing.md)

### 🔊 Audio & Speech
- [TTS Models](developer-guides/speech/01-tts-models.md)
- [Text-to-Speech (TTS)](developer-guides/speech/02-tts.md)
- [Voice Cloning](developer-guides/speech/03-voice-cloning.md)
- [Voice Design](developer-guides/speech/04-voice-design.md)
- [Realtime Streaming ASR](developer-guides/speech/05-realtime-streaming.md)
- [LaTeX to Speech](developer-guides/speech/06-latex-to-speech.md)
- [Improve Recognition Accuracy](developer-guides/speech/07-improve-recognition-accuracy.md)
- [Realtime Translation](developer-guides/speech/08-realtime-translation.md)

### 🎯 Accuracy Tuning
- [Image Generation Tuning](developer-guides/accuracy-tuning/01-image-generation.md)
- [Video Generation Tuning](developer-guides/accuracy-tuning/02-video-generation.md)

### ⚙️ Run and Scale
- [Async Task Management](developer-guides/run-and-scale/01-async-task-management.md)
- [Connection Pooling](developer-guides/run-and-scale/02-connection-pooling.md)
- [Latency Optimization](developer-guides/run-and-scale/03-latency-optimization.md)
- [Cost Optimization](developer-guides/run-and-scale/04-cost-optimization.md)
- [Safety](developer-guides/run-and-scale/05-safety.md)

### 🔌 Integrations
- [MLOps & Observability](developer-guides/integrations/01-mlops-observability.md)
- [Alerts](developer-guides/integrations/02-alerts.md)

### 🛠️ Clients & Developer Tools
- [Chatbox](developer-guides/clients-and-developer-tools/01-chatbox.md)
- [Cherry Studio](developer-guides/clients-and-developer-tools/02-cherry-studio.md)
- [Claude Code](developer-guides/clients-and-developer-tools/03-claude-code.md)
- [Cline](developer-guides/clients-and-developer-tools/04-cline.md)
- [Codex](developer-guides/clients-and-developer-tools/05-codex.md)
- [Cursor](developer-guides/clients-and-developer-tools/06-cursor.md)
- [Dify](developer-guides/clients-and-developer-tools/07-dify.md)
- [Hermes Agent](developer-guides/clients-and-developer-tools/08-hermes-agent.md)
- [Kilo CLI](developer-guides/clients-and-developer-tools/09-kilo-cli.md)
- [Lingma](developer-guides/clients-and-developer-tools/10-lingma.md)
- [More Tools](developer-guides/clients-and-developer-tools/11-more-tools.md)
- [OpenClaw](developer-guides/clients-and-developer-tools/12-openclaw.md)
- [OpenCode](developer-guides/clients-and-developer-tools/13-opencode.md)
- [Postman](developer-guides/clients-and-developer-tools/14-postman.md)
- [Qoder](developer-guides/clients-and-developer-tools/15-qoder.md)
- [Qwen Code](developer-guides/clients-and-developer-tools/16-qwen-code.md)
- [QwenPaw](developer-guides/clients-and-developer-tools/17-qwenpaw.md)

### 🔒 Administration
- [API Keys](developer-guides/administration/01-api-keys.md)
- [Workspaces](developer-guides/administration/02-workspace.md)
- [Rate Limits](developer-guides/administration/03-rate-limits.md)

### 🛡️ Security & Compliance
- [Data Security](developer-guides/security-compliance/01-data-security.md)
- [Audit Logs](developer-guides/security-compliance/02-audit-logs.md)

---

## 📖 API & SDK Reference

### 🔧 Preparation
- [API Key Setup](api-reference/preparation/01-api-key.md)
- [Install SDK](api-reference/preparation/02-install-sdk.md)
- [Export API Key as Env Variable](api-reference/preparation/03-export-api-key-env.md)
- [Error Messages](api-reference/preparation/04-error-messages.md)

### 💬 Chat Completions
- [OpenAI-compatible Chat](api-reference/chat/01-openai-chat.md)
- [OpenAI Responses API](api-reference/chat/02-openai-responses.md)
- [DashScope Chat](api-reference/chat/03-dashscope.md)
- [Anthropic-compatible Chat](api-reference/chat/04-anthropic.md)

### 🔢 Text Embedding
- [OpenAI-compatible Embedding](api-reference/text-embedding/01-openai-embedding.md)
- [DashScope Embedding](api-reference/text-embedding/02-dashscope-embedding.md)

### 🔄 Reranking
- [OpenAI-compatible Rerank](api-reference/rerank/01-openai-rerank.md)
- [DashScope Rerank](api-reference/rerank/02-dashscope-rerank.md)

### 🖼️ Multimodal Embedding
- [DashScope Multimodal Embedding](api-reference/multimodal-embedding/01-dashscope-multimodal-embedding.md)

### 🎨 Image Generation
- [Qwen Text-to-Image](api-reference/image-generation/01-qwen-text-to-image.md)
- [WAN2.7 Image Gen & Edit](api-reference/image-generation/02-wan27-image-gen-edit.md)
- [WAN2.6 Image Gen & Edit](api-reference/image-generation/03-wan26-image-gen-edit.md)
- [WAN Text-to-Image v2](api-reference/image-generation/04-wan-text-to-image-v2.md)
- [Z-Image](api-reference/image-generation/05-z-image.md)

### 🎬 Video Generation
- [WAN2.7 Text-to-Video](api-reference/video-generation/01-wan27-text-to-video.md)
- [WAN2.7 Image-to-Video](api-reference/video-generation/02-wan27-image-to-video.md)
- [WAN2.7 Reference-to-Video](api-reference/video-generation/03-wan27-reference-to-video.md)
- [WAN Text-to-Video](api-reference/video-generation/04-wan-text-to-video.md)
- [HappyHorse Text-to-Video](api-reference/video-generation/05-happyhorse-text-to-video.md)
- [HappyHorse Image-to-Video](api-reference/video-generation/06-happyhorse-image-to-video.md)
- [HappyHorse Reference-to-Video](api-reference/video-generation/07-happyhorse-reference-to-video.md)
- [HappyHorse Video Editing](api-reference/video-generation/08-happyhorse-video-editing.md)

### 🔊 Speech Synthesis
- [CosyVoice Voice List](api-reference/speech-synthesis/01-cosyvoice-voice-list.md)
- [Voice Cloning - Create Voice](api-reference/speech-synthesis/02-voice-cloning-create.md)

### 🌐 Speech Translation
- [LiveTranslate Realtime](api-reference/speech-translation/01-livetranslate-realtime.md)

### 🗂️ Platform API
- [Create Batch](api-reference/platform-api/01-batch-create.md)
- [Cancel Batch](api-reference/platform-api/02-batch-cancel.md)
- [Conversations](api-reference/platform-api/03-conversations.md)
- [File Management](api-reference/platform-api/04-file.md)

### 🧰 Toolkit & Frameworks
- [OpenAI-compatible Overview](api-reference/toolkitframework/01-openai-compatible-overview.md)

### ➕ More APIs
- [Generate Temporary API Key](api-reference/more/01-generate-temporary-api-key.md)
- [Manage Asynchronous Tasks](api-reference/more/02-manage-asynchronous-tasks.md)

---

## 💰 Resources

- [Billing Overview](resources/01-billing-overview.md)
- [Free Quota](resources/02-free-quota.md)

---

## 🎫 Token Plan

- [Token Plan Overview](token-plan/01-overview.md)
- [Token Plan FAQ](token-plan/02-faq.md)

---

## 📋 Changelog

- [Models Changelog](changelog/01-models.md)
- [Platform Changelog](changelog/02-platform.md)

---

*Documentación descargada el 2026-07-14. Para la versión más actualizada, visita [docs.qwencloud.com](https://docs.qwencloud.com).*
