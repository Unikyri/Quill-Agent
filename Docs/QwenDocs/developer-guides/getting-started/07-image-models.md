# Image generation models

> **Source:** https://docs.qwencloud.com/developer-guides/getting-started/image-models

Choose a model for text-to-image generation, image editing, and more.

 Copy page ## [​ ](#text-to-image) Text-to-image

Start with `wan2.7-image-pro` — it covers text rendering, brand color control, multi-image sets with consistent characters, and image editing in one model. Up to 4096×4096 for text-to-image, 2048×2048 for editing.
### [​ ](#when-to-use-z-image-turbo-instead) When to use `z-image-turbo` instead


- You only need generation (no editing)

- Speed or cost is the priority — 10x faster, ~1/5 the price

- Realistic portraits and product shots


### [​ ](#when-to-use-qwen-image-2-0-pro-instead) When to use `qwen-image-2.0-pro` instead


- You need negative prompts to exclude specific elements from output

- You need up to 6 image variants per call (wan supports up to 4 in standard mode)


### [​ ](#not-sure-which-fits) Not sure which fits?

[Try models in Try AI](https://home.qwencloud.com/try-ai) — same prompt, compare results. Style preference is personal.
## [​ ](#image-editing) Image editing

Start with `wan2.7-image-pro` — it supports multi-image reference (up to 9 input images), bounding-box interactive editing, and character-consistent multi-image sets.
### [​ ](#when-to-use-qwen-image-2-0-pro-instead-2) When to use `qwen-image-2.0-pro` instead

Need negative prompts during editing → `qwen-image-2.0-pro` (same model ID for generation + editing).
## [​ ](#recommended-models) Recommended models

ModelUse this whent2iEditMax outputMax resolution`wan2.7-image-pro`Text rendering, brand colors, multi-image sets, multi-image editing✓✓4 (12 sequential)4096×4096 (t2i) / 2048×2048 (edit)`wan2.7-image`Same capabilities, faster generation, up to 2K✓✓4 (12 sequential)2048×2048`z-image-turbo`Fast generation, low cost, realistic portraits✓12048×2048`qwen-image-2.0-pro`Negative prompts, up to 6 image variants✓✓62048×2048`qwen-image-2.0`Faster variant of qwen-image-2.0-pro✓✓62048×2048 
## [​ ](#all-models) All models

 Wan

 Model IDt2iEditMax outputMax resolution`wan2.7-image-pro`✓✓4 (12 sequential)4096×4096 (t2i) / 2048×2048 (edit)`wan2.7-image`✓✓4 (12 sequential)2048×2048 Qwen Image

 Model IDt2iEditMax outputMax resolution`qwen-image-2.0-pro`✓✓62048×2048`qwen-image-2.0-pro-2026-04-22`✓✓62048×2048`qwen-image-2.0-pro-2026-03-03`✓✓62048×2048`qwen-image-2.0`✓✓62048×2048`qwen-image-2.0-2026-03-03`✓✓62048×2048 Z-Image

 Model IDt2iEditMax outputMax resolution`z-image-turbo`✓12048×2048 Legacy

 Previous generation models. We recommend Wan 2.7 or Qwen Image 2.0 for new projects.### Wan
Model IDt2iEditMax outputMax resolution`wan2.6-t2i`✓4~1440×1440`wan2.6-image`✓*✓4~1440×1440`wan2.5-t2i-preview`✓4~1440×1440`wan2.5-i2i-preview`✓41280×1280`wan2.2-t2i-plus`✓4~1440×1440`wan2.2-t2i-flash`✓4~1440×1440`wan2.1-t2i-plus`✓4~1440×1440`wan2.1-t2i-turbo`✓4~1440×1440 * Requires `enable_interleave=true` and `stream=true`. See [text-to-image guide](/developer-guides/image-generation/text-to-image).### Qwen Image
Model IDt2iEditMax outputMax resolution`qwen-image-max`✓11664×928`qwen-image-max-2025-12-30`✓11664×928`qwen-image-plus`✓11664×928`qwen-image-plus-2026-01-09`✓11664×928`qwen-image`✓11664×928`qwen-image-edit-max`✓62048×2048`qwen-image-edit-max-2026-01-16`✓62048×2048`qwen-image-edit-plus`✓62048×2048`qwen-image-edit-plus-2025-12-15`✓62048×2048`qwen-image-edit-plus-2025-10-30`✓62048×2048`qwen-image-edit`✓11024×1024 

---

## [​ ](#learn-more) Learn more

[ ## Image generation guide
Learn how to generate images via API. ](/developer-guides/image-generation/text-to-image)[ ## Try free
Try models in the browser — no API key needed. ](https://home.qwencloud.com/try-ai) [Previous ](/developer-guides/multimodal/ocr)[Text-to-image Generate images from text prompts. Next ](/developer-guides/image-generation/text-to-image)
