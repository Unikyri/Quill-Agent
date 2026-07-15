# Codex

> **Source:** https://docs.qwencloud.com/developer-guides/clients-and-developer-tools/codex

OpenAI&#x27;s terminal AI coding assistant

 Copy page Codex is a terminal AI coding assistant developed by OpenAI. Connect it to Qwen Cloud via Token Plan (Team Edition), Coding Plan, or pay-as-you-go billing.
## [​ ](#install-codex) Install Codex


- 
Install or update [Node.js](https://nodejs.org/en/download/) (v18.0 or later).


- 
Run the following command in a terminal to install Codex.


Copy ```\nnpm install -g @openai/codex

``` 
Run the following command to verify the installation.
Copy ```\ncodex --version

``` 
## [​ ](#configure-access-credentials) Configure access credentials

To connect, edit the configuration file `~/.codex/config.toml` and configure the environment variable `OPENAI_API_KEY`. Replace the corresponding values based on your selected billing plan.
### [​ ](#token-plan-team-edition) Token Plan (Team Edition)

For `model`, select a [supported model](/token-plan/overview#supported-models). Set the `OPENAI_API_KEY` environment variable to the Token Plan (Team Edition) dedicated [API Key](https://home.qwencloud.com/api-keys).
#### [​ ](#responses-api-qwen3-7-max-qwen3-7-plus-qwen3-6-plus-qwen3-6-flash) Responses API (qwen3.7-max, qwen3.7-plus, qwen3.6-plus, qwen3.6-flash)

qwen3.7-max, qwen3.7-plus, qwen3.6-plus, and qwen3.6-flash support the Responses API, compatible with the latest Codex version.
Copy ```\nmodel_provider = "Model_Studio_Token_Plan"
model = "qwen3.7-max"
[model_providers.Model_Studio_Token_Plan]
name = "Model_Studio_Token_Plan"
base_url = "https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1"
env_key = "OPENAI_API_KEY"
wire_api = "responses"

``` 
#### [​ ](#chat-completions-api-other-models) Chat/Completions API (other models)

Other models must be connected via the Chat/Completions API, which requires installing an older version of Codex, such as 0.80.0:
Copy ```\nmodel_provider = "Model_Studio_Token_Plan"
model = "glm-5"
[model_providers.Model_Studio_Token_Plan]
name = "Model_Studio_Token_Plan"
base_url = "https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1"
env_key = "OPENAI_API_KEY"
wire_api = "chat"

``` 
#### [​ ](#configure-environment-variables) Configure environment variables

Set the `OPENAI_API_KEY` environment variable to the Token Plan (Team Edition) dedicated API Key.
- macOS 
- Windows 

 
- Run the following command in a terminal to check the default shell type.


Copy ```\necho $SHELL

``` 
- Set the environment variable based on your shell type:


- Zsh 
- Bash 

 Copy ```\n# Replace YOUR_API_KEY with the Token Plan (Team Edition) API Key
echo &#x27;export OPENAI_API_KEY="YOUR_API_KEY"&#x27; >> ~/.zshrc

``` Copy ```\n# Replace YOUR_API_KEY with the Token Plan (Team Edition) API Key
echo &#x27;export OPENAI_API_KEY="YOUR_API_KEY"&#x27; >> ~/.bash_profile

``` 
- Run the following command to apply the environment variable.


- Zsh 
- Bash 

 Copy ```\nsource ~/.zshrc

``` Copy ```\nsource ~/.bash_profile

``` - CMD 
- PowerShell 

 
- Run the following command in CMD to set the environment variable.


Copy ```\nREM Replace YOUR_API_KEY with the Token Plan (Team Edition) API Key
setx OPENAI_API_KEY "YOUR_API_KEY"

``` 
- Open a new CMD window and run the following command to verify that the environment variable is set.


Copy ```\necho %OPENAI_API_KEY%

``` 
- Run the following command in PowerShell to set the environment variable.


Copy ```\n# Replace YOUR_API_KEY with the Token Plan (Team Edition) API Key
[Environment]::SetEnvironmentVariable("OPENAI_API_KEY", "YOUR_API_KEY", [EnvironmentVariableTarget]::User)

``` 
- Open a new PowerShell window and run the following command to verify that the environment variable is set.


Copy ```\necho $env:OPENAI_API_KEY

``` 
### [​ ](#coding-plan) Coding Plan

For `model`, select a [supported model](/coding-plan/overview#plan-details). Set the `OPENAI_API_KEY` environment variable to the Coding Plan dedicated [API Key](https://home.qwencloud.com/api-keys).
#### [​ ](#chat-completions-api) Chat/Completions API

Coding Plan only supports the Chat/Completions API. Install an older version of Codex, such as 0.80.0:
Copy ```\nmodel_provider = "Model_Studio_Coding_Plan"
model = "qwen3.7-plus"
[model_providers.Model_Studio_Coding_Plan]
name = "Model_Studio_Coding_Plan"
base_url = "https://coding-intl.dashscope.aliyuncs.com/v1"
env_key = "OPENAI_API_KEY"
wire_api = "chat"

``` 
#### [​ ](#configure-environment-variables-2) Configure environment variables

Set the `OPENAI_API_KEY` environment variable to the Coding Plan dedicated API Key.
- macOS 
- Windows 

 
- Run the following command in a terminal to check the default shell type.


Copy ```\necho $SHELL

``` 
- Set the environment variable based on your shell type:


- Zsh 
- Bash 

 Copy ```\n# Replace YOUR_API_KEY with the Coding Plan API Key
echo &#x27;export OPENAI_API_KEY="YOUR_API_KEY"&#x27; >> ~/.zshrc

``` Copy ```\n# Replace YOUR_API_KEY with the Coding Plan API Key
echo &#x27;export OPENAI_API_KEY="YOUR_API_KEY"&#x27; >> ~/.bash_profile

``` 
- Run the following command to apply the environment variable.


- Zsh 
- Bash 

 Copy ```\nsource ~/.zshrc

``` Copy ```\nsource ~/.bash_profile

``` - CMD 
- PowerShell 

 
- Run the following command in CMD to set the environment variable.


Copy ```\nREM Replace YOUR_API_KEY with the Coding Plan API Key
setx OPENAI_API_KEY "YOUR_API_KEY"

``` 
- Open a new CMD window and run the following command to verify that the environment variable is set.


Copy ```\necho %OPENAI_API_KEY%

``` 
- Run the following command in PowerShell to set the environment variable.


Copy ```\n# Replace YOUR_API_KEY with the Coding Plan API Key
[Environment]::SetEnvironmentVariable("OPENAI_API_KEY", "YOUR_API_KEY", [EnvironmentVariableTarget]::User)

``` 
- Open a new PowerShell window and run the following command to verify that the environment variable is set.


Copy ```\necho $env:OPENAI_API_KEY

``` 
### [​ ](#pay-as-you-go) Pay-as-you-go

Set the `OPENAI_API_KEY` environment variable to your [Qwen Cloud API Key](https://home.qwencloud.com/api-keys). For available models, see [supported models](/api-reference/chat/openai-chat).
Pay-as-you-go supports both the Responses API and Chat/Completions API. Choose the appropriate one based on the model you are using:
#### [​ ](#responses-api) Responses API

Applicable to models that support the [OpenAI Responses API](/api-reference/chat/openai-responses) (such as qwen3.7-max, qwen3.7-plus, qwen3.6-plus, and qwen3.6-flash), compatible with the latest Codex version.
Copy ```\nmodel_provider = "Model_Studio"
model = "qwen3.7-max"
[model_providers.Model_Studio]
name = "Model_Studio"
base_url = "https://dashscope-intl.aliyuncs.com/compatible-mode/v1"
env_key = "OPENAI_API_KEY"
wire_api = "responses"

``` 
#### [​ ](#chat-completions-api-2) Chat/Completions API

Applicable to models that only support the Chat/Completions API. Requires installing Codex 0.80.0:
Copy ```\nmodel_provider = "Model_Studio"
model = "qwen3.6-plus"
[model_providers.Model_Studio]
name = "Model_Studio"
base_url = "https://dashscope-intl.aliyuncs.com/compatible-mode/v1"
env_key = "OPENAI_API_KEY"
wire_api = "chat"

``` 
#### [​ ](#configure-environment-variables-3) Configure environment variables

Set the `OPENAI_API_KEY` environment variable to your [Qwen Cloud API Key](https://home.qwencloud.com/api-keys).
- macOS 
- Windows 

 
- Run the following command in a terminal to check the default shell type.


Copy ```\necho $SHELL

``` 
- Set the environment variable based on your shell type:


- Zsh 
- Bash 

 Copy ```\n# Replace YOUR_API_KEY with the Qwen Cloud API Key
echo &#x27;export OPENAI_API_KEY="YOUR_API_KEY"&#x27; >> ~/.zshrc

``` Copy ```\n# Replace YOUR_API_KEY with the Qwen Cloud API Key
echo &#x27;export OPENAI_API_KEY="YOUR_API_KEY"&#x27; >> ~/.bash_profile

``` 
- Run the following command to apply the environment variable.


- Zsh 
- Bash 

 Copy ```\nsource ~/.zshrc

``` Copy ```\nsource ~/.bash_profile

``` - CMD 
- PowerShell 

 
- Run the following command in CMD to set the environment variable.


Copy ```\nREM Replace YOUR_API_KEY with the Qwen Cloud API Key
setx OPENAI_API_KEY "YOUR_API_KEY"

``` 
- Open a new CMD window and run the following command to verify that the environment variable is set.


Copy ```\necho %OPENAI_API_KEY%

``` 
- Run the following command in PowerShell to set the environment variable.


Copy ```\n# Replace YOUR_API_KEY with the Qwen Cloud API Key
[Environment]::SetEnvironmentVariable("OPENAI_API_KEY", "YOUR_API_KEY", [EnvironmentVariableTarget]::User)

``` 
- Open a new PowerShell window and run the following command to verify that the environment variable is set.


Copy ```\necho $env:OPENAI_API_KEY

``` 
## [​ ](#verify-configuration) Verify configuration

After the configuration is complete, open a new terminal window and run the following command to start Codex:
Copy ```\ncodex

``` 
If the chat interface launches successfully, the configuration is correct.
## [​ ](#faq) FAQ

### [​ ](#what-should-i-do-if-a-third-party-tool-reports-domestic-models-not-supported-or-check-rejected-bad-request-400) What should I do if a third-party tool reports "domestic models not supported" or "check rejected / Bad request (400)"?

**Cause**: Some third-party management tools (such as CC-Switch) send a "health check / connection test" probe request when switching providers. The format of this probe differs from the request format Codex actually uses, so the Qwen Cloud gateway may reject it with 400 Bad request, and the tool then reports "domestic models not supported". This message only indicates that the health check probe failed; **it does not mean Qwen Cloud lacks support for domestic models, nor does it affect actual Codex usage.**
**Note**: Qwen Cloud supports using domestic models such as qwen3.7-max, qwen3.7-plus, qwen3.6-plus, qwen3.6-flash, and glm-5 through Codex. For configuration details, see [Configure access credentials](#configure-access-credentials) above.
**Solution**: Configure Codex directly in `~/.codex/config.toml` as described in Configure access credentials, without relying on the third-party tool&#x27;s health check result. After configuration, start Codex as described in [Verify configuration](#verify-configuration); if the chat interface launches normally, domestic models are working.
### [​ ](#what-should-i-do-if-i-get-a-wire-api-configuration-error) What should I do if I get a wire_api configuration error?

**Cause**: Newer versions of Codex no longer support `wire_api = "chat"`. Depending on the version, you may see one of the following errors:

- `wire_api = "chat" is no longer supported`

- `unknown configuration field wire_api`


**Solution**:

- Error `wire_api = "chat" is no longer supported`: Change `wire_api` to `responses` and verify that `base_url` is correct. See [Configure access credentials](#configure-access-credentials) for configuration examples.

- Error `unknown configuration field wire_api`: Remove the `wire_api` line from the corresponding provider section in `~/.codex/config.toml`.


### [​ ](#what-should-i-do-if-i-get-the-error-unexpected-status-401-unauthorized) What should I do if I get the error "unexpected status 401 Unauthorized"?

**Cause**:

- Using an API Key from a different plan (API Keys for Token Plan (Team Edition), Coding Plan, and pay-as-you-go are not interchangeable)

- Subscription expired

- API Key was copied incompletely, contains spaces, or has a typo


**Solution**:

- Verify that you are using the dedicated API Key for your selected plan.

- Go to the management page of your selected plan and check whether the subscription has expired.

- Re-copy the API Key and make sure it is complete and has no spaces.

- If the error persists after verifying the above, reset the API Key on the management page of your selected plan. After resetting, use the new API Key for configuration.


### [​ ](#what-should-i-do-if-i-get-the-error-unexpected-status-404-not-found) What should I do if I get the error "unexpected status 404 Not Found"?

**Cause**: The `base_url` or `wire_api` in the configuration file is incorrect.
**Solution**: Verify that `base_url` and `wire_api` match the configuration for your selected plan. See the configuration examples for your plan in [Configure access credentials](#configure-access-credentials) above.
### [​ ](#what-should-i-do-if-i-get-the-error-stream-disconnected-before-completion-stream-closed-before-response-completed) What should I do if I get the error "stream disconnected before completion: stream closed before response.completed"?

**Cause**: The streaming connection between Codex and the server was interrupted before the response completed. This commonly occurs in the following scenarios:

- The conversation thread is too long, causing a context compaction request to fail

- Unstable network causing the SSE or WebSocket connection to drop mid-stream

- Server overload or rate limiting that terminates the connection early


**Solution**:

- Start a new conversation thread to avoid excessive context accumulation in a single thread.

- Check your network connection. Try disabling VPN or proxy and retry.

- Wait and retry. Codex has a built-in retry mechanism that resolves most transient failures automatically.


 [Previous ](/developer-guides/clients-and-developer-tools/cursor)[Qwen Code Official terminal-based AI coding tool Next ](/developer-guides/clients-and-developer-tools/qwen-code)
