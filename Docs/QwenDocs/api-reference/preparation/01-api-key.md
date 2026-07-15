# Get an API key

> **Source:** https://docs.qwencloud.com/api-reference/preparation/api-key

First step to using models

 Copy page ## [​ ](#create-an-api-key) Create an API key

 1 Create an API key

Go to [API Keys](https://home.qwencloud.com/api-keys) and click **Create API key**. 2 Add a description

Enter a description to help identify your key, then click **Generate Key**. 3 Copy your API key

Copy and save your API key immediately. For security reasons, the full key is only shown once. The key list displays a masked version after creation. 
## [​ ](#use-an-api-key) Use an API key

 This topic describes the pay-as-you-go API key. If you are using a [Token Plan](/token-plan/overview) or [Coding Plan](/coding-plan/overview), use the corresponding dedicated API key (`sk-sp-xxxxx` format) from the [API Keys page](https://home.qwencloud.com/api-keys). 

- 
**Method 1: In third-party tools such as [Chatbox](/developer-guides/clients-and-developer-tools/chatbox)**
When calling models in third-party tools, provide:

API key (from this page)

- Base URL: `https://dashscope-intl.aliyuncs.com/compatible-mode/v1`

- Model name (such as `qwen3.7-plus`)


- 
**Method 2: Through code**
[Set the API key as an environment variable](/api-reference/preparation/export-api-key-env) to avoid hardcoding it in source code.


 When you call a model from code or a third-party tool, in addition to the API key you must specify a **service endpoint** (the **API Host** shown in the creation success dialog, which corresponds to the `base_url` in your SDK or HTTP request). Qwen Cloud provides both **OpenAI-compatible** and **Anthropic-compatible** protocol interfaces. The `base_url` differs between the two protocols and varies by region. Refer to the documentation for the protocol you use:
- OpenAI-compatible protocol: [OpenAI-compatible - Chat](/api-reference/chat/openai-chat)

- Anthropic-compatible protocol: [Anthropic-compatible Messages](/api-reference/chat/anthropic)


 
Never share your API key—unauthorized use causes security risks and financial loss.
## [​ ](#manage-api-keys) Manage API keys

From the [API Keys](https://home.qwencloud.com/api-keys) page, you can:

- **Search**: Find keys by description.

- **Edit**: Update the description of an existing key.

- **Delete**: Permanently remove a key. Deleted keys cannot be recovered and any applications using them will stop working.


## [​ ](#validity) Validity

API keys remain valid permanently unless you manually delete them.
## [​ ](#error-codes) Error codes

If a model call fails and returns an error message, see [Error codes](/api-reference/preparation/error-messages) for troubleshooting.
## [​ ](#faq) FAQ

### [​ ](#pay-as-you-go-api-key-vs-coding-plan-api-key) Pay-as-you-go API key vs Coding Plan API key

General API keys (`sk-xxxxx`) use pay-as-you-go billing for API calls. Coding Plan API keys (`sk-sp-xxxxx`) are tied to your subscription and provide higher rate limits and premium features.
### [​ ](#echo-works-but-code-reports-no-api-key-found) echo works but code reports "no API key found"

**Q: I used the `echo` command and confirmed the environment variable was set correctly. Why does my code still report that the API key cannot be found?**
A: This can occur for the following reasons:

- Scenario 1: **A temporary environment variable was set**. A temporary variable is valid only within the current terminal session and does not affect running IDEs or other applications. Refer to [Configure your API key](/api-reference/preparation/export-api-key-env) to set a permanent environment variable.

- Scenario 2: **The IDE, command-line tool, or application was not restarted**. You must restart your IDE (such as VS Code) or terminal to load the new environment variables. If you set the environment variable after deploying an application, you must restart the application service.

- Scenario 3: **The variable is missing from a service configuration file**. If your application is started by a service manager (such as systemd or supervisord), you may need to add the environment variable to the service manager&#x27;s configuration file.

- Scenario 4: **Using the `sudo` command**. `sudo` does not inherit all environment variables by default. Use `sudo -E python xx.py` where the `-E` parameter ensures that environment variables are passed.


 Previous [Configure your API key Avoid hardcoding secrets Next ](/api-reference/preparation/export-api-key-env)
