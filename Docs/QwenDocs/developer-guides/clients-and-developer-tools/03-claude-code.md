# Claude Code

> **Source:** https://docs.qwencloud.com/developer-guides/clients-and-developer-tools/claude-code

Anthropic&#x27;s terminal AI coding assistant

 Copy page Claude Code is a command-line AI coding assistant developed by Anthropic. Connect it to Qwen Cloud via Token Plan (Team Edition), Coding Plan, or pay-as-you-go billing.
## [​ ](#install-claude-code) Install Claude Code

- macOS 
- Windows 

 
- 
Install or update [Node.js](https://nodejs.org/en/download/) (v18.0 or later).


- 
Run the following command in the terminal to install Claude Code.


Copy ```\nnpm install -g @anthropic-ai/claude-code

``` 
- Run the following command to verify the installation. If a version number is displayed, the installation is successful.


Copy ```\nclaude --version

``` To use Claude Code on Windows, install WSL or [Git for Windows](https://git-scm.com/install/windows), then run the following command in WSL or Git Bash.Copy ```\nnpm install -g @anthropic-ai/claude-code

``` For more details, see the [Windows setup guide](https://docs.anthropic.com/en/docs/claude-code/setup#windows-setup) in the official Claude Code documentation. 
### [​ ](#skip-login-verification) Skip login verification

Edit or create `~/.claude.json` (Windows path: `C:\Users\<username>\.claude.json`), and set `hasCompletedOnboarding` to `true` to skip the official Anthropic login verification.
Copy ```\n{
 "hasCompletedOnboarding": true
}

``` 
## [​ ](#configure-access-credentials) Configure access credentials

Create `~/.claude/settings.json` (Windows path: `C:\Users\<username>\.claude\settings.json`), and add the corresponding configuration based on your selected plan.
### [​ ](#token-plan-team-edition) Token Plan (Team Edition)

Replace YOUR_API_KEY with the Token Plan (Team Edition) dedicated [API Key](https://home.qwencloud.com/api-keys). For available models, see [supported models](/token-plan/overview#supported-models) for Token Plan (Team Edition).
Copy ```\n{
 "env": {
 "ANTHROPIC_AUTH_TOKEN": "YOUR_API_KEY",
 "ANTHROPIC_BASE_URL": "https://token-plan.ap-southeast-1.maas.aliyuncs.com/apps/anthropic",
 "ANTHROPIC_MODEL": "qwen3.7-max",
 "ANTHROPIC_DEFAULT_HAIKU_MODEL": "qwen3.6-flash",
 "ANTHROPIC_DEFAULT_SONNET_MODEL": "qwen3.7-max",
 "ANTHROPIC_DEFAULT_OPUS_MODEL": "qwen3.7-max",
 "CLAUDE_CODE_SUBAGENT_MODEL": "qwen3.7-max"
 }
}

``` 
### [​ ](#coding-plan) Coding Plan

Replace YOUR_API_KEY with the Coding Plan dedicated [API Key](https://home.qwencloud.com/api-keys). For available models, see [supported models](/coding-plan/overview#plan-details) for Coding Plan.
Copy ```\n{
 "env": {
 "ANTHROPIC_AUTH_TOKEN": "YOUR_API_KEY",
 "ANTHROPIC_BASE_URL": "https://coding-intl.dashscope.aliyuncs.com/apps/anthropic",
 "ANTHROPIC_MODEL": "qwen3.7-plus",
 "ANTHROPIC_DEFAULT_HAIKU_MODEL": "qwen3.7-plus",
 "ANTHROPIC_DEFAULT_SONNET_MODEL": "qwen3.7-plus",
 "ANTHROPIC_DEFAULT_OPUS_MODEL": "qwen3.7-plus",
 "CLAUDE_CODE_SUBAGENT_MODEL": "qwen3.7-plus"
 }
}

``` 
### [​ ](#pay-as-you-go) Pay-as-you-go

Replace YOUR_API_KEY with your [Qwen Cloud API Key](https://home.qwencloud.com/api-keys). For available models, see [Anthropic API compatible](/api-reference/chat/anthropic).
Copy ```\n{
 "env": {
 "ANTHROPIC_AUTH_TOKEN": "YOUR_API_KEY",
 "ANTHROPIC_BASE_URL": "https://dashscope-intl.aliyuncs.com/apps/anthropic",
 "ANTHROPIC_MODEL": "qwen3.7-max",
 "ANTHROPIC_DEFAULT_HAIKU_MODEL": "qwen3.6-flash",
 "ANTHROPIC_DEFAULT_SONNET_MODEL": "qwen3.7-max",
 "ANTHROPIC_DEFAULT_OPUS_MODEL": "qwen3.7-max",
 "CLAUDE_CODE_SUBAGENT_MODEL": "qwen3.7-max"
 }
}

``` 
## [​ ](#verify-configuration) Verify configuration

After saving the configuration, open a new terminal window and run the following command to verify whether the connection is successful:
Copy ```\nclaude "Hello"

``` 
If the model returns a normal response, the configuration is successful.
Run `claude` to enter interactive mode, which supports multi-turn conversations, file editing, and command execution. For details, see the [Claude Code official documentation](https://code.claude.com/docs/en/overview).
## [​ ](#claude-code-ide-plugins) Claude Code IDE plugins

After completing the CLI configuration above, you can install the Claude Code plugin in your IDE, which directly reuses the configuration in `settings.json`.
### [​ ](#vs-code) VS Code


- Search for `Claude Code for VS Code` in the extension marketplace and install it.

- Restart VS Code and click the icon in the upper-right corner to open Claude Code.

- Type `/` in the dialog box, select General config, and set the model in Selected Model.


### [​ ](#jetbrains) JetBrains


- Search for `Claude Code` in the extension marketplace and install it.

- Restart the IDE and click the icon in the upper-right corner to start using it.


## [​ ](#faq) FAQ

### [​ ](#error-codes) Error codes

If you encounter errors during configuration, refer to the FAQ documentation for the corresponding billing plan:

- Pay-as-you-go: [Anthropic API compatible - Error codes](/api-reference/preparation/error-messages)

- Coding Plan: [Coding Plan FAQ](/coding-plan/faq)

- Token Plan (Team Edition): [Token Plan (Team Edition) FAQ](/token-plan/faq)


### [​ ](#after-starting-claude-code-the-interface-displays-unable-to-connect-to-anthropic-services) After starting Claude Code, the interface displays "Unable to connect to Anthropic services"

This error indicates that Claude Code is attempting to connect to the official Anthropic service instead of the Qwen Cloud server, usually because the environment variables are not configured correctly or have not taken effect. Follow these steps to troubleshoot:

- 
**Check the configuration**: After starting Claude Code, run the `/status` command and verify that the values of `ANTHROPIC_BASE_URL` and `ANTHROPIC_AUTH_TOKEN` correctly point to the Qwen Cloud address. If the output is empty or points to a non-Qwen Cloud address, check whether the `settings.json` configuration is correct.


- 
**Verify hasCompletedOnboarding**: Check that `hasCompletedOnboarding` is set to `true` in the `~/.claude.json` file. Otherwise, Claude Code will attempt to connect to the official Anthropic service for login verification at startup.


- 
**Reopen the terminal**: After modifying the configuration file, open a new terminal window and run the `claude` command for the configuration to take effect.


 [Previous ](/developer-guides/clients-and-developer-tools/hermes-agent)[OpenCode Terminal AI coding assistant Next ](/developer-guides/clients-and-developer-tools/opencode)
