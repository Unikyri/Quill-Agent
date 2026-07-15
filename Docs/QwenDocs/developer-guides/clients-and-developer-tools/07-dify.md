# Dify

> **Source:** https://docs.qwencloud.com/developer-guides/clients-and-developer-tools/dify

Low-code LLM app platform

 Copy page Dify is a low-code platform for building LLM applications with visual workflows. Connect it to Qwen Cloud&#x27;s pay-as-you-go API to create chat assistants, agents, and knowledge bases powered by Qwen models.
## [​ ](#quick-start) Quick start

Get running in a few minutes:
Copy ```\n# 1. Install plugin
Go to Dify marketplace → Models → Find "TONGYI" → Install

# 2. Configure (Settings → Model Providers → TONGYI → Settings)
API Key: sk-xxx
Use international endpoint: Yes

# 3. Test (Create blank app → Chat assistant)
Select model: qwen3.5-plus
Type message: "Write a Python hello world program"

``` 
You should see: The model responds with Python code
## [​ ](#configuration) Configuration

### [​ ](#basic-setup) Basic setup

Configure Dify to use Qwen Cloud:

- Plugin: TONGYI (from Dify marketplace)

- API endpoint: International (set to "Yes")

- Authentication: API key required

- Model selection: Choose from available models in plugin


 **Free quota and billing:**
- First-time users get a free quota. See [Free quota](/resources/free-quota) for details.

- Enable [Free quota only](https://home.qwencloud.com/benefits) to prevent unexpected charges.


 
### [​ ](#step-by-step-configuration) Step-by-step configuration

 1 Install TONGYI plugin

Go to [Dify marketplace](https://cloud.dify.ai/plugins?category=discover) → **Models** → **TONGYI** → Install 2 Configure API key

Click profile → **Settings** → **Model Providers** → **TONGYI** → **Settings**
- API Key: Your [API key](https://home.qwencloud.com/api-keys)

- Use international endpoint: **Yes**


 3 Enable models

Click models on TONGYI card → Toggle on desired models 
 For newest models not in TONGYI plugin: Use **OpenAI-API-compatible** plugin with endpoint `https://dashscope-intl.aliyuncs.com/compatible-mode/v1` 
### [​ ](#limitations) Limitations


- Plugin maintenance: TONGYI plugin is maintained by Dify, not Qwen Cloud

- Model availability: Some newest models may require OpenAI-compatible plugin


## [​ ](#examples) Examples

- Chat assistant 
- Workflow with LLM node 
- Knowledge base 
- Vision models 

 1 Create app

Workspace → **Create blank app** → **Chat assistant** 2 Configure model

Select **qwen3.5-plus** → Enable thinking mode if available 3 Test conversation

Type: "Explain how neural networks work" 1 Create workflow

Workspace → Create **Chatflow** or **Workflow** 2 Add LLM node

Drag LLM node to canvas → Select **qwen3.5-plus** 3 Run node

Add user message → Click run button 1 Create knowledge base

Go to [Knowledge bases](https://cloud.dify.ai/datasets) → Create new 2 Upload documents

Select data source → Upload files 3 Configure embedding

Select **text-embedding-v4** for text segmentation 1 Select vision model

Choose **qwen3.5-plus** or **qwen3-vl-plus** 2 Enable vision

Enable **Vision** toggle in chat or LLM node 3 Upload images

Upload images in conversation for analysis 
## [​ ](#troubleshooting) Troubleshooting

**"Invalid API-key provided" error**

Solution:

- Try earlier TONGYI plugin version

- Use API key from default workspace (not sub-workspace)

- Verify "Use international endpoint" is set to Yes


**Models with `-latest` suffix not available**

Solution: Use OpenAI-API-compatible plugin with:

- Endpoint: `https://dashscope-intl.aliyuncs.com/compatible-mode/v1`

- API key: Your DashScope key

- Model: Enter model ID manually


**High token consumption**

Solution:

- Use appropriate models for tasks

- Configure reasonable context windows

- Clear conversation history regularly


**Vision toggle not appearing**

Solution: Ensure you&#x27;ve selected a vision-capable model (`qwen3.5-plus` or `qwen3-vl-plus`)

## [​ ](#advanced-features) Advanced features

### [​ ](#thinking-mode) Thinking mode

For models that support reasoning:

- Select model with thinking support

- Enable thinking mode toggle

- Set to "True" for step-by-step reasoning


### [​ ](#code-execution-nodes) Code execution nodes

Extract reasoning from responses:

- Use regex in code execution nodes

- Separate thinking process from final answer

- Format output as needed


## [​ ](#faq) FAQ

### [​ ](#using-qwen-omni-and-qwen-ocr-models) Using Qwen-Omni and Qwen-OCR models

These models cannot be configured directly in Dify. Integrate them using an HTTP node in a Chatflow or workflow. For integration details, refer to the cURL command examples in each model&#x27;s documentation.
 Use streaming output for API calls in HTTP nodes to reduce the risk of timeouts. 
### [​ ](#using-wan-models-for-image-and-video-generation) Using Wan models for image and video generation

Dify does not offer a dedicated plugin for Wan models, but you can achieve text-to-image and text-to-video generation using nodes in a Dify Chatflow or workflow:

- **Create a workflow**: In the [workspace](https://cloud.dify.ai/apps), create a new **Workflow**. Add HTTP Request nodes for the Wan API&#x27;s POST (create task) and GET (query result) endpoints.

- **Configure environment variables**: On the workflow page, find the environment variables icon and set `DASHSCOPE_API_KEY` to your API key.

- **Test the output**: Click **Run** to generate the output. Text-to-video generation typically takes five minutes or more.

- **Publish as a tool (optional)**: To use these capabilities in other applications, click **Publish** and select **Publish as tool**.


 Use models `wan2.2-t2i-flash` (text-to-image) and `wan2.1-t2v-turbo` (text-to-video). For API details, see the [image generation](/developer-guides/image-generation/text-to-image) and [video generation](/developer-guides/video-generation/text-to-video) guides. 
## [​ ](#related-resources) Related resources


- **Models**: [Available models](/developer-guides/getting-started/text-generation-models)

- **Vision models**: [Image understanding guide](/developer-guides/multimodal/vision)

- **Embeddings**: [Text embedding models](/developer-guides/embeddings/text-embedding)

- **API docs**: [OpenAI-compatible reference](/api-reference/chat/openai-chat)


 [Previous ](/developer-guides/clients-and-developer-tools/postman)[More tools Connect any OpenAI or Anthropic compatible programming tool Next ](/developer-guides/clients-and-developer-tools/more-tools)
