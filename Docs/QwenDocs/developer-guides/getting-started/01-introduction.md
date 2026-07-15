# Build with Qwen Cloud

> **Source:** https://docs.qwencloud.com/developer-guides/getting-started/introduction

AI models for text, vision, speech, and image & video generation.

 Copy page 
## Developer quickstart
Make your first API request in minutes. Seamlessly integrate with any OpenAI SDK or client.

[Get started →](/developer-guides/getting-started/first-api-call) Copy ```\nimport os
from openai import OpenAI

client = OpenAI(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1"
)
completion = client.chat.completions.create(
 model="qwen3.7-plus",
 messages=[{"role": "user", "content": "Summarize the benefits of solar energy in three bullet points."}]
)
print(completion.choices[0].message.content)

``` 

## [​ ](#models) Models

Start with **qwen3.7-plus** for a balance of quality and speed. Choose **qwen3.7-max** for the hardest reasoning and coding tasks, or **qwen3.6-flash** for cost efficiency. All models share the same API — just change the `model` parameter. [Browse all models →](/developer-guides/getting-started/model-selection)
[ ## qwen3.7-max
Complex reasoning and coding ](https://www.qwencloud.com/models/qwen3.7-max)[ ## qwen3.7-plus
Balanced performance, speed, and cost ](https://www.qwencloud.com/models/qwen3.7-plus)[ ## qwen3.6-flash
Fast and cost-effective ](https://www.qwencloud.com/models/qwen3.6-flash) 
## [​ ](#start-building) Start building

[ ## Read and generate text
Prompt models to generate text, summarize, translate, or write code ](/developer-guides/text-generation/quickstart)[ ## Understand images and video
Analyze images, extract text from screenshots, or reproduce designs from mockups ](/developer-guides/multimodal/vision)[ ## Generate images
Create and edit images from text prompts with Wan and Flux models ](/developer-guides/image-generation/text-to-image)[ ## Generate videos
Animate images into video clips or generate videos from text descriptions ](/developer-guides/video-generation/text-to-video)[ ## Synthesize speech
Convert text to natural speech with built-in voices, voice cloning, or voice design ](/developer-guides/speech/tts-models)[ ## Build agentic applications
Connect models to external tools and APIs with function calling ](/developer-guides/text-generation/function-calling)[ ## Tackle complex tasks with thinking
Use reasoning models to solve multi-step math, logic, and coding problems ](/developer-guides/text-generation/thinking)[ ## Get structured data from models
Extract JSON that conforms to a schema from any model response ](/developer-guides/text-generation/structured-output) 

---

[Pricing](/developer-guides/getting-started/pricing) | [API Reference](/api-reference/preparation/api-key) | [Free Quota](/resources/free-quota) Previous [First API call Get started in a few minutes Next ](/developer-guides/getting-started/first-api-call)
