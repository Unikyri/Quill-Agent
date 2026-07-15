# Visual understanding models

> **Source:** https://docs.qwencloud.com/developer-guides/getting-started/vision-models

Choose a model for image analysis, video understanding, OCR, and more.

 Copy page ## [​ ](#image-and-video-understanding) Image and video understanding

Start with `qwen3.7-plus` — strongest accuracy, 1M context, 2-hour video support, and the full feature set including function calling and built-in tools. Once your use case works well, try `qwen3.6-flash` to reduce cost — near-flagship quality with the same context and features.
### [​ ](#image-resolution) Image resolution

Most models support up to 16M pixels per image. Higher resolution costs more tokens: each image uses `h × w / (32 × 32) + 2` tokens.
### [​ ](#video-support) Video support


- Up to 2 hours / 2GB → `qwen3.7-plus`, `qwen3.6-plus`, `qwen3.6-flash`, `qwen3.5-plus`, `qwen3.5-flash`

- Up to 1 hour / 2GB → `qwen3.5-omni-plus`, `qwen3.5-omni-flash` (also accepts audio — see [Speech models](/developer-guides/speech/s2s-models))


### [​ ](#function-calling-built-in-tools) Function calling + built-in tools

Let the model take actions based on what it sees in images or video.

- Function calling: Qwen3.6, Qwen3.5, Qwen3.5-Omni (including realtime models), and Qwen3-VL models

- Built-in tools (web search, code execution — no setup): `qwen3.7-plus`, `qwen3.6-plus`, `qwen3.6-flash`, `qwen3.5-plus`, `qwen3.5-flash` only


### [​ ](#structured-output) Structured output

Get valid JSON from visual input — e.g., extract product info from a photo.
Available on Qwen3.6 and Qwen3.5 in non-thinking mode.
## [​ ](#recommended-models) Recommended models

ModelContextMax pixels/imageMax video durationMax video sizeMax images (URL)Max images (Base64)Max videosFunction callingBuilt-in toolsStructured output`qwen3.7-plus`1M16M2h2GB2,04825064✓✓✓`qwen3.6-plus`1M16M2h2GB25625064✓✓✓`qwen3.6-flash`1M16M2h2GB25625064✓✓✓`qwen3.5-flash`1M16M2h2GB25625064✓✓✓`qwen3.5-omni-plus`256k—1h2GB2,0482501✓——`qwen3.5-omni-flash`256k—1h2GB2,0482501✓—— 
## [​ ](#all-models) All models

 Qwen3.7

 Model IDInputOutputContextMax OutputMax images (URL)Max images (Base64)Max videosFunction callingBuilt-in toolsStructured output`qwen3.7-max-2026-06-08`Text, image, videoText1M64k2,04825064✓✓—`qwen3.7-plus`Text, image, videoText1M64k2,04825064✓✓✓`qwen3.7-plus-2026-05-26`Text, image, videoText1M64k2,04825064✓✓✓ Qwen3.6

 Model IDInputOutputContextMax OutputMax images (URL)Max images (Base64)Max videosFunction callingBuilt-in toolsStructured output`qwen3.6-plus`Text, image, videoText1M64k25625064✓✓✓`qwen3.6-plus-2026-04-02`Text, image, videoText1M64k25625064✓✓✓`qwen3.6-flash`Text, image, videoText1M64k25625064✓✓✓`qwen3.6-35b-a3b`Text, image, videoText32k8k25625064✓✓✓`qwen3.6-27b`Text, image, videoText32k8k25625064✓—✓ Qwen3.5

 Model IDInputOutputContextMax OutputMax images (URL)Max images (Base64)Max videosFunction callingBuilt-in toolsStructured output`qwen3.5-plus`Text, image, videoText1M64k25625064✓✓✓`qwen3.5-plus-2026-04-20`Text, image, videoText1M64k25625064✓✓✓`qwen3.5-plus-2026-02-15`Text, image, videoText1M64k25625064✓✓✓`qwen3.5-flash`Text, image, videoText1M64k25625064✓✓✓`qwen3.5-flash-2026-02-23`Text, image, videoText1M64k25625064✓✓✓`qwen3.5-397b-a17b`Text, image, videoText32k8k25625064✓✓✓`qwen3.5-122b-a10b`Text, image, videoText32k8k25625064✓✓✓`qwen3.5-27b`Text, image, videoText32k8k25625064✓✓✓`qwen3.5-35b-a3b`Text, image, videoText32k8k25625064✓✓✓ Qwen3.5-Omni

 Unlike other models on this page, Qwen3.5-Omni accepts audio input and can output both text and speech.**Standard**Model IDInputOutputContextMax OutputMax images (URL)Max images (Base64)Max videosFunction callingBuilt-in toolsStructured output`qwen3.5-omni-plus`Text, image, audio, videoText, audio256k64k2,0482501✓——`qwen3.5-omni-plus-2026-03-15`Text, image, audio, videoText, audio256k64k2,0482501✓——`qwen3.5-omni-flash`Text, image, audio, videoText, audio256k64k2,0482501✓——`qwen3.5-omni-flash-2026-03-15`Text, image, audio, videoText, audio256k64k2,0482501✓—— **Realtime** — streaming audio input with built-in Voice Activity Detection (VAD).Model IDInputOutputContextMax OutputFunction calling`qwen3.5-omni-plus-realtime`Text, image, audio (streaming)Text, audio256k64k✓`qwen3.5-omni-plus-realtime-2026-03-15`Text, image, audio (streaming)Text, audio256k64k✓`qwen3.5-omni-flash-realtime`Text, image, audio (streaming)Text, audio256k64k✓`qwen3.5-omni-flash-realtime-2026-03-15`Text, image, audio (streaming)Text, audio256k64k✓ **Captioner** (open source) — audio captioning model.Model IDInputOutputContextMax OutputMax images (URL)Max images (Base64)Max videosFunction callingBuilt-in toolsStructured output`qwen3-omni-30b-a3b-captioner`AudioText64k32k—————— Legacy

 Older model versions retained for backward compatibility. We recommend Qwen3.5 or Qwen3.5-Omni for new projects.Model IDInputOutputContextMax OutputMax images (URL)Max images (Base64)Max videosFunction callingBuilt-in toolsStructured output`qwen-vl-ocr`Text, imageText38k8k256250————`qwen-vl-ocr-2025-11-20`Text, imageText38k8k256250————`qwen3-vl-plus`Text, image, videoText256k32k25625064✓—✓`qwen3-vl-plus-2025-12-19`Text, image, videoText256k32k25625064✓—✓`qwen3-vl-plus-2025-09-23`Text, image, videoText256k32k25625064✓—✓`qwen3-vl-flash`Text, image, videoText256k32k25625064✓—✓`qwen3-vl-flash-2026-01-22`Text, image, videoText256k32k25625064✓—✓`qwen3-vl-flash-2025-10-15`Text, image, videoText256k32k25625064✓—✓`qwen3-omni-flash`Text, image, audio, videoText, audio64k16k2,0482501✓——`qwen3-omni-flash-2025-12-01`Text, image, audio, videoText, audio64k16k2,0482501✓——`qwen3-omni-flash-2025-09-15`Text, image, audio, videoText, audio64k16k2,0482501✓——`qwen3-omni-flash-realtime`Text, image, audio (streaming)Text, audio64k16k——————`qwen3-omni-flash-realtime-2025-12-01`Text, image, audio (streaming)Text, audio64k16k——————`qwen3-omni-flash-realtime-2025-09-15`Text, image, audio (streaming)Text, audio64k16k——————`qwen-omni-turbo`Text, image, audio, videoText, audio32k2k2,0482501———`qwen-omni-turbo-latest`Text, image, audio, videoText, audio32k2k2,0482501———`qwen-omni-turbo-2025-03-26`Text, image, audio, videoText, audio32k2k2,0482501———`qwen-omni-turbo-realtime`Text, audio (streaming)Text, audio32k2k——————`qwen-omni-turbo-realtime-latest`Text, audio (streaming)Text, audio32k2k——————`qwen-omni-turbo-realtime-2025-05-08`Text, audio (streaming)Text, audio32k2k——————`qwen3-vl-235b-a22b-thinking`Text, image, videoText128k8k25625064✓——`qwen3-vl-235b-a22b-instruct`Text, image, videoText128k8k25625064✓—✓`qwen3-vl-32b-thinking`Text, image, videoText128k8k25625064✓——`qwen3-vl-32b-instruct`Text, image, videoText128k8k25625064✓—✓`qwen3-vl-30b-a3b-thinking`Text, image, videoText128k8k25625064✓——`qwen3-vl-30b-a3b-instruct`Text, image, videoText128k8k25625064✓—✓`qwen3-vl-8b-thinking`Text, image, videoText128k8k25625064✓——`qwen3-vl-8b-instruct`Text, image, videoText128k8k25625064✓—✓`qwen2.5-vl-72b-instruct`Text, image, videoText128k8k25625064✓—✓`qwen2.5-vl-32b-instruct`Text, image, videoText128k8k25625064✓—✓`qwen2.5-vl-7b-instruct`Text, image, videoText128k8k25625064✓—✓`qwen2.5-vl-3b-instruct`Text, image, videoText128k8k25625064✓—✓`qwen2.5-omni-7b`Text, image, audio, videoText, audio32k2k2,0482501———`qwen-vl-max`Text, imageText32k8k256250————`qwen-vl-max-latest`Text, imageText128k8k256250————`qwen-vl-max-2025-08-13`Text, imageText128k8k256250————`qwen-vl-max-2025-04-08`Text, imageText128k8k256250————`qwen-vl-plus`Text, imageText128k8k256250————`qwen-vl-plus-latest`Text, imageText128k8k256250————`qwen-vl-plus-2025-08-15`Text, imageText128k8k256250————`qwen-vl-plus-2025-05-07`Text, imageText128k8k256250————`qwen-vl-plus-2025-01-25`Text, imageText128k8k256250————`qvq-max`Text, imageText128k8k256250————`qvq-max-latest`Text, imageText128k8k256250————`qvq-max-2025-03-25`Text, imageText128k8k256250———— 

---

## [​ ](#learn-more) Learn more

[ ## Vision understanding guide
Learn how to send images and videos to multimodal models. ](/developer-guides/multimodal/vision)[ ## Try free
Try models in the browser — just sign in, no API key setup required. ](https://home.qwencloud.com/try-ai) [Previous ](/developer-guides/third-party-models/kimi)[Analyze images and videos Generate content from visual inputs Next ](/developer-guides/multimodal/vision)
