# Qoder

> **Source:** https://docs.qwencloud.com/developer-guides/clients-and-developer-tools/qoder

Agentic coding platform with IDE, CLI, and JetBrains plugin

 Copy page Qoder is an agentic coding platform for software development that supports a desktop IDE, CLI, and JetBrains plugin, and can connect to Qwen Cloud via Coding Plan or Token Plan (Team Edition).
 Connecting to Qwen Cloud requires Qoder Pro Trial, Pro, Pro+, or Ultra version. Free and Teams versions are not supported. 
## [​ ](#qoder-ide) Qoder IDE

### [​ ](#install) Install


- Go to the [Qoder official website](https://qoder.com/) to download and install Qoder.

- Complete the initial configuration and log in to your Qoder account after the first launch.


### [​ ](#configure-access-credentials) Configure access credentials


- 
Open Qoder settings from the top-right corner of the interface, select **Models**, and click **Add**.


- 
Model configuration details are as follows:


Configuration ItemDescriptionProviderSelect **Alibaba Cloud Model Studio - International** from the dropdown menuTypeSelect **Token Plan** or **Coding Plan** based on your billing planModelSelect a model from the dropdown menu. Only text generation models are supported.API KeyEnter the dedicated API key for your chosen plan: Token Plan (Team Edition) - [Get API Key](https://home.qwencloud.com/api-keys); Coding Plan - [Get API Key](https://home.qwencloud.com/api-keys) 

- 
Click **Add**. The model configuration is complete once validation passes.


- 
Select the corresponding model from the model list to start using it.


## [​ ](#qoder-cli) Qoder CLI

### [​ ](#install-2) Install


- Run the following command in the terminal to install.


Copy ```\ncurl -fsSL https://qoder.com/install | bash

``` 

- Verify that the installation was successful.


Copy ```\nqodercli --version

``` 
If a version number is displayed, the installation is successful.
### [​ ](#log-in-to-qoder) Log in to Qoder

Authentication is required before use. There are two methods:

- 
**Log in via TUI (Recommended)**

Run `qodercli` to enter the interactive interface, and type `/login` in the dialog.

- Select `login with browser` or `login with qoder personal access token` to complete the login.


- 
**Log in via environment variable**
Suitable for non-interactive environments (such as CI/CD pipelines). Replace `your_personal_access_token_here` with your actual token, which can be obtained from the [Integrations page](https://qoder.com/account/integrations).


- macOS/Linux 
- Windows 

 Copy ```\nexport QODER_PERSONAL_ACCESS_TOKEN="your_personal_access_token_here"

``` Copy ```\nset QODER_PERSONAL_ACCESS_TOKEN=your_personal_access_token_here

``` 
### [​ ](#configure-access-credentials-2) Configure access credentials


- Type `/model` in the dialog, and use the Tab key to switch to `Custom`.

- Press Enter to select Add custom model, choose **Alibaba Cloud Model Studio - International** as the provider, and select **Token Plan** or **Coding Plan** as the type based on your billing plan.

- Select a model and enter the dedicated API key for your chosen plan. Confirm and wait for the configuration to take effect.


### [​ ](#use-qoder-cli) Use Qoder CLI


- Restart Qoder CLI.


Copy ```\nqodercli

``` 

- Type `/model` in the dialog, use the Tab key to switch to `Custom`, and select the configured model to start using it.


## [​ ](#jetbrains-plugin) JetBrains plugin


- Open a JetBrains IDE (such as IntelliJ IDEA, PyCharm, etc.), search for `Qoder` in the extensions marketplace and install it.

- Click Qoder in the right navigation bar and complete the login in the Qoder chat panel.

- Click the settings icon in the top-right corner, select **Plugin Settings**, and click **Add Model** in the popup.

- Configuration details are as follows:


Configuration ItemDescriptionProviderSelect **Alibaba Cloud Model Studio - International** from the dropdown menuTypeSelect **Token Plan** or **Coding Plan** based on your billing planModelSelect a model from the dropdown menu. Only text generation models are supported.API KeyEnter the dedicated API key for your chosen plan: Token Plan (Team Edition) - [Get API Key](https://home.qwencloud.com/api-keys); Coding Plan - [Get API Key](https://home.qwencloud.com/api-keys) 
After completing the configuration, click **OK** and wait for the configuration to take effect.

- Select the configured model from Custom models to start a conversation.


## [​ ](#learn-more) Learn more

To learn more about Qoder&#x27;s agents, MCP, Skills, and other extension capabilities, refer to the [Qoder official documentation](https://docs.qoder.com/zh).
## [​ ](#faq) FAQ

### [​ ](#error-codes) Error codes

If you encounter errors during configuration, refer to the FAQ documentation for the corresponding billing plan:

- Coding Plan: [Coding Plan FAQ](/coding-plan/faq)

- Token Plan (Team Edition): [Token Plan (Team Edition) FAQ](/token-plan/faq)


### [​ ](#why-cant-i-find-the-model-option-in-qoder-settings) Why can&#x27;t I find the model option in Qoder settings?

Possible reasons:

- **Not logged in**: You must log in first before you can chat and configure models.

- **Unsupported version**: Connecting to Qwen Cloud requires Qoder Pro Trial, Pro, Pro+, or Ultra version. Free and Teams versions are not supported. Please upgrade to a supported version.


### [​ ](#previously-able-to-use-qwen-cloud-models-but-after-some-time-unable-to-switch-and-configuration-is-not-editable) Previously able to use Qwen Cloud models, but after some time, unable to switch and configuration is not editable

**Cause**: New Qoder accounts can try Qoder Pro Trial for free for two weeks. After the trial expires, the account automatically reverts to the free version, which no longer supports Qwen Cloud integration.
**Solution**: Upgrade to Qoder Pro Trial, Pro, Pro+, or Ultra version.
For more questions, refer to [Coding Plan FAQ](/coding-plan/faq) or [Token Plan (Team Edition) FAQ](/token-plan/faq). [Previous ](/developer-guides/clients-and-developer-tools/cline)[Qoder CN (formerly Lingma) Alibaba Cloud intelligent coding assistant IDE Next ](/developer-guides/clients-and-developer-tools/lingma)
