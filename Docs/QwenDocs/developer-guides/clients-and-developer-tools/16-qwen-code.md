# Qwen Code

> **Source:** https://docs.qwencloud.com/developer-guides/clients-and-developer-tools/qwen-code

Official terminal-based AI coding tool

 Copy page Qwen Code is a terminal-based AI coding tool. Connect it to Qwen Cloud via Token Plan (Team Edition), Coding Plan, or pay-as-you-go billing.
## [​ ](#install-qwen-code) Install Qwen Code

- macOS/Linux 
- Windows 

 Open a terminal and run the following command to install Qwen Code.Copy ```\nbash -c "$(curl -fsSL https://qwen-code-assets.oss-cn-hangzhou.aliyuncs.com/installation/install-qwen.sh)" -s --source bailian

``` Type `cmd` in the taskbar search box, select **Run as administrator**, and run the following command in the `cmd` window to install Qwen Code.Copy ```\ncurl -fsSL -o %TEMP%\install-qwen.bat https://qwen-code-assets.oss-cn-hangzhou.aliyuncs.com/installation/install-qwen.bat && %TEMP%\install-qwen.bat --source bailian

``` On Windows, after installation is complete, close the current `cmd` window to apply the environment variables. Reopen `cmd` and run the following command to check the installed version.Copy ```\nqwen --version

``` 
## [​ ](#configure-access-credentials) Configure access credentials

Launch Qwen Code and type `/auth` for visual configuration. Qwen Cloud offers three billing plans. Choose based on your needs:

- **Token Plan (Team Edition)**: Subscription per seat, with token consumption deducted from Credits.

- **Coding Plan**: Fixed monthly subscription billed by number of model calls.

- **Pay-as-you-go**: Post-paid based on actual usage.


### [​ ](#token-plan-team-edition) Token Plan (Team Edition)

Launch Qwen Code and type `/auth`, then select **Subscription Plan** > **Alibaba Cloud Model Studio Token Plan**, and enter the Token Plan (Team Edition) dedicated [API Key](https://home.qwencloud.com/api-keys) to complete the configuration. For available models, see Token Plan (Team Edition) [supported models](/token-plan/overview#supported-models).
 Advanced configuration: via settings.json

 Edit or create a `settings.json` file and replace `YOUR_API_KEY` with your Token Plan (Team Edition) dedicated API Key. The file path is as follows:
- macOS/Linux: `~/.qwen/settings.json`

- Windows: `C:\Users\<Windows username>\.qwen\settings.json`


Copy ```\n{
 "env": {
 "BAILIAN_TOKEN_PLAN_API_KEY": "YOUR_API_KEY"
 },
 "modelProviders": {
 "openai": [
 {
 "id": "qwen3.7-max",
 "name": "[Token Plan Team Edition] qwen3.7-max",
 "baseUrl": "https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1",
 "envKey": "BAILIAN_TOKEN_PLAN_API_KEY",
 "generationConfig": {
 "extra_body": {
 "enable_thinking": true
 }
 }
 },
 {
 "id": "qwen3.7-plus",
 "name": "[Token Plan Team Edition] qwen3.7-plus",
 "baseUrl": "https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1",
 "envKey": "BAILIAN_TOKEN_PLAN_API_KEY",
 "generationConfig": {
 "extra_body": {
 "enable_thinking": true
 }
 }
 },
 {
 "id": "qwen3.6-plus",
 "name": "[Token Plan Team Edition] qwen3.6-plus",
 "baseUrl": "https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1",
 "envKey": "BAILIAN_TOKEN_PLAN_API_KEY",
 "generationConfig": {
 "extra_body": {
 "enable_thinking": true
 }
 }
 },
 {
 "id": "qwen3.6-flash",
 "name": "[Token Plan Team Edition] qwen3.6-flash",
 "baseUrl": "https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1",
 "envKey": "BAILIAN_TOKEN_PLAN_API_KEY",
 "generationConfig": {
 "extra_body": {
 "enable_thinking": true
 }
 }
 },
 {
 "id": "deepseek-v4-pro",
 "name": "[Token Plan Team Edition] deepseek-v4-pro",
 "baseUrl": "https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1",
 "envKey": "BAILIAN_TOKEN_PLAN_API_KEY"
 },
 {
 "id": "deepseek-v4-flash",
 "name": "[Token Plan Team Edition] deepseek-v4-flash",
 "baseUrl": "https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1",
 "envKey": "BAILIAN_TOKEN_PLAN_API_KEY"
 },
 {
 "id": "deepseek-v3.2",
 "name": "[Token Plan Team Edition] deepseek-v3.2",
 "baseUrl": "https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1",
 "envKey": "BAILIAN_TOKEN_PLAN_API_KEY"
 },
 {
 "id": "kimi-k2.7-code",
 "name": "[Token Plan Team Edition] kimi-k2.7-code",
 "baseUrl": "https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1",
 "envKey": "BAILIAN_TOKEN_PLAN_API_KEY",
 "generationConfig": {
 "extra_body": {
 "enable_thinking": true
 }
 }
 },
 {
 "id": "kimi-k2.6",
 "name": "[Token Plan Team Edition] kimi-k2.6",
 "baseUrl": "https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1",
 "envKey": "BAILIAN_TOKEN_PLAN_API_KEY",
 "generationConfig": {
 "extra_body": {
 "enable_thinking": true
 }
 }
 },
 {
 "id": "kimi-k2.5",
 "name": "[Token Plan Team Edition] kimi-k2.5",
 "baseUrl": "https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1",
 "envKey": "BAILIAN_TOKEN_PLAN_API_KEY",
 "generationConfig": {
 "extra_body": {
 "enable_thinking": true
 }
 }
 },
 {
 "id": "glm-5.2",
 "name": "[Token Plan Team Edition] glm-5.2",
 "baseUrl": "https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1",
 "envKey": "BAILIAN_TOKEN_PLAN_API_KEY",
 "generationConfig": {
 "extra_body": {
 "enable_thinking": true
 }
 }
 },
 {
 "id": "glm-5.1",
 "name": "[Token Plan Team Edition] glm-5.1",
 "baseUrl": "https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1",
 "envKey": "BAILIAN_TOKEN_PLAN_API_KEY",
 "generationConfig": {
 "extra_body": {
 "enable_thinking": true
 }
 }
 },
 {
 "id": "glm-5",
 "name": "[Token Plan Team Edition] glm-5",
 "baseUrl": "https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1",
 "envKey": "BAILIAN_TOKEN_PLAN_API_KEY",
 "generationConfig": {
 "extra_body": {
 "enable_thinking": true
 }
 }
 },
 {
 "id": "MiniMax-M2.5",
 "name": "[Token Plan Team Edition] MiniMax-M2.5",
 "baseUrl": "https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1",
 "envKey": "BAILIAN_TOKEN_PLAN_API_KEY"
 }
 ]
 },
 "security": {
 "auth": {
 "selectedType": "openai"
 }
 },
 "tokenPlan": {
 "region": "global"
 },
 "model": {
 "name": "qwen3.7-plus"
 },
 "$version": 3
}

``` 
### [​ ](#coding-plan) Coding Plan

Launch Qwen Code and type `/auth`, then select **Subscription Plan** > **Alibaba Cloud Model Studio Coding Plan**, choose the Coding Plan region (global), and enter the Coding Plan dedicated [API Key](https://home.qwencloud.com/api-keys) to complete the configuration. For available models, see Coding Plan [supported models](/coding-plan/overview#plan-details).
 Advanced configuration: via settings.json

 Edit or create a `settings.json` file and replace `YOUR_API_KEY` with your Coding Plan dedicated API Key. The file path is as follows:
- macOS/Linux: `~/.qwen/settings.json`

- Windows: `C:\Users\<Windows username>\.qwen\settings.json`


Copy ```\n{
 "env": {
 "BAILIAN_CODING_PLAN_API_KEY": "YOUR_API_KEY"
 },
 "modelProviders": {
 "openai": [
 {
 "id": "qwen3.7-plus",
 "name": "[Coding Plan] qwen3.7-plus",
 "baseUrl": "https://coding-intl.dashscope.aliyuncs.com/v1",
 "envKey": "BAILIAN_CODING_PLAN_API_KEY",
 "generationConfig": {
 "extra_body": {
 "enable_thinking": true
 }
 }
 },
 {
 "id": "qwen3.6-plus",
 "name": "[Coding Plan] qwen3.6-plus",
 "baseUrl": "https://coding-intl.dashscope.aliyuncs.com/v1",
 "envKey": "BAILIAN_CODING_PLAN_API_KEY",
 "generationConfig": {
 "extra_body": {
 "enable_thinking": true
 }
 }
 },
 {
 "id": "qwen3-coder-plus",
 "name": "[Coding Plan] qwen3-coder-plus",
 "baseUrl": "https://coding-intl.dashscope.aliyuncs.com/v1",
 "envKey": "BAILIAN_CODING_PLAN_API_KEY"
 }
 ]
 },
 "security": {
 "auth": {
 "selectedType": "openai"
 }
 },
 "codingPlan": {
 "region": "global"
 },
 "model": {
 "name": "qwen3.7-plus"
 },
 "$version": 3
}

``` 
### [​ ](#pay-as-you-go) Pay-as-you-go

Launch Qwen Code and type `/auth`, select **Use your own API Key** > **Standard API Key**, then choose the region and enter the [Qwen Cloud API Key](https://home.qwencloud.com/api-keys) to complete the configuration. For available models, see [OpenAI compatible - Supported models](/api-reference/chat/openai-chat).
 Advanced configuration: via settings.json

 Edit or create a `settings.json` file and replace `YOUR_API_KEY` with your Qwen Cloud API Key. The file path is as follows:
- macOS/Linux: `~/.qwen/settings.json`

- Windows: `C:\Users\<Windows username>\.qwen\settings.json`


Copy ```\n{
 "env": {
 "BAILIAN_API_KEY": "YOUR_API_KEY"
 },
 "modelProviders": {
 "openai": [
 {
 "id": "qwen3.6-plus",
 "name": "[Qwen Cloud] qwen3.6-plus",
 "baseUrl": "https://dashscope-intl.aliyuncs.com/compatible-mode/v1",
 "envKey": "BAILIAN_API_KEY",
 "generationConfig": {
 "extra_body": {
 "enable_thinking": true
 }
 }
 }
 ]
 },
 "security": {
 "auth": {
 "selectedType": "openai"
 }
 },
 "model": {
 "name": "qwen3.7-plus"
 },
 "$version": 3
}

``` To add [other models](/api-reference/chat/openai-chat), append them in the same format within `modelProviders.openai`. 
## [​ ](#verify-configuration) Verify configuration

After configuration is complete, run the following command in your project directory to launch Qwen Code and start a conversation.
Copy ```\nqwen

``` 
## [​ ](#qwen-code-ide-plugins) Qwen Code IDE plugins

Qwen Code supports usage as a plugin in VS Code and JetBrains IDE, providing AI coding capabilities within the IDE.
### [​ ](#vs-code) VS Code

Ensure your VS Code version is 1.85.0 or higher before use.

- Open VS Code, search for `Qwen Code Companion` in the Extensions Marketplace and install it.

- The CLI and IDE plugin share the same `settings.json`. If you have already completed the configuration above, skip this step. Otherwise, follow the [Configure access credentials](#configure-access-credentials) section above.

- Click the icon in the upper-right corner to launch Qwen Code. Type or click `/`, then select `Switch model` to switch models.


### [​ ](#jetbrains-ide) JetBrains IDE

Ensure your JetBrains IDE supports Agent Client Protocol (ACP) before use.

- 
Qwen Code IDE depends on Qwen Code CLI. First, follow the [Install Qwen Code](#install-qwen-code) section in this document to install and complete [Configure access credentials](#configure-access-credentials).


- 
Open JetBrains IDE, go to the **AI Chat** page, and click **Install Plugin**. IDEA will install the JetBrains AI Assistant plugin.


- 
Click the three-dot menu in the upper-right corner of the **AI Chat** window, select **Add Custom Agent**, and fill in the following configuration. Replace `/path/to/qwen` with the Qwen Code installation path. Run the following command to find the path:

macOS/Linux: `which qwen`

- Windows: `where qwen` (CMD) or `Get-Command qwen` (PowerShell)


Copy ```\n{
 "agent_servers": {
 "qwen": {
 "command": "/path/to/qwen",
 "args": ["--acp"],
 "env": {}
 }
 }
}

``` 

- After configuration, Qwen Code will appear in the **AI Chat** panel. You can switch models in the lower-right corner.


## [​ ](#common-commands) Common commands

 The following commands apply to Qwen Code CLI. The IDE plugin only supports some commands. Please refer to actual usage. 
CommandDescriptionExample`/model`Switch the model used in the current session.`/model``/auth`Change the authentication method.`/auth``/init`Analyze the current directory and create an initial context file (QWEN.md).`/init``/clear`Clear the terminal screen and start a new conversation.`/clear``/compress`Replace chat history with a summary to save tokens.`/compress``/settings`Open the settings editor to configure language, theme, and more.`/settings``/summary`Generate a project summary based on the conversation history.`/summary``/resume`Resume a previous conversation session.`/resume``/stats`Display detailed statistics for the current session.`/stats``/help`Display help information for available commands.`/help` or `/?``/quit`Exit Qwen Code.`/quit` 
For more advanced features of Qwen Code, see the [Qwen Code official documentation](https://qwenlm.github.io/qwen-code-docs/en/users/features/commands/).
## [​ ](#learn-more) Learn more


- For advanced features such as sub-agents, MCP, and Skills in Qwen Code, see the [Qwen Code official documentation](https://qwenlm.github.io/qwen-code-docs/en/users/overview/).

- For Qwen Code use cases, see [Use Cases](https://qwenlm.github.io/qwen-code-docs/en/showcase/).


## [​ ](#faq) FAQ

### [​ ](#error-codes) Error codes

If you encounter errors during configuration, refer to the FAQ documentation for your billing plan:

- Pay-as-you-go: [Error codes and troubleshooting](/api-reference/preparation/error-messages)

- Coding Plan: [Coding Plan FAQ](/coding-plan/faq)

- Token Plan (Team Edition): [Token Plan (Team Edition) FAQ](/token-plan/faq)


### [​ ](#how-to-switch-models) How to switch models?

Type the `/model` command directly in Qwen Code to select and switch from the list of models configured in `settings.json`. To add a new model, add the corresponding configuration under `settings.json` in the `modelProviders` section, then restart Qwen Code. [Previous ](/developer-guides/clients-and-developer-tools/codex)[QwenPaw Open-source personal AI assistant from the AgentScope team Next ](/developer-guides/clients-and-developer-tools/qwenpaw)
