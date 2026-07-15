# Configure your API key

> **Source:** https://docs.qwencloud.com/api-reference/preparation/export-api-key-env

Avoid hardcoding secrets

 Copy page ## [​ ](#prerequisites) Prerequisites

[Create an API key](/api-reference/preparation/api-key) first.
## [​ ](#steps) Steps

- Linux 
- macOS 
- Windows 

 ### Permanent environment variable
Set a permanent environment variable for the current user: 1 Add the environment variable

Add the variable to `~/.bashrc`:Copy ```\n# Replace YOUR_DASHSCOPE_API_KEY with your API key.
echo "export DASHSCOPE_API_KEY=&#x27;YOUR_DASHSCOPE_API_KEY&#x27;" >> ~/.bashrc

``` Edit manually

 Open `~/.bashrc`:Copy ```\nnano ~/.bashrc

``` Add the following content to the file:Copy ```\n# Replace YOUR_DASHSCOPE_API_KEY with your API key.
export DASHSCOPE_API_KEY="YOUR_DASHSCOPE_API_KEY"

``` In the nano editor, press Ctrl+X and then Y. Press Enter to save and close the file. 2 Apply the changes

Apply changes:Copy ```\nsource ~/.bashrc

``` 3 Verify

Verify in a new session:Copy ```\necho $DASHSCOPE_API_KEY

``` ### Temporary environment variable
Set a temporary variable (current session only): 1 Set the variable

Copy ```\n# Replace YOUR_DASHSCOPE_API_KEY with your API key.
export DASHSCOPE_API_KEY="YOUR_DASHSCOPE_API_KEY"

``` 2 Verify

Copy ```\necho $DASHSCOPE_API_KEY

``` ### Permanent environment variable
Set a permanent environment variable for the current user: 1 Check your default shell type

Copy ```\necho $SHELL

``` 2 Add the environment variable

- Zsh 
- Bash 

 Add the variable to `~/.zshrc`:Copy ```\n# Replace YOUR_DASHSCOPE_API_KEY with your API key.
echo "export DASHSCOPE_API_KEY=&#x27;YOUR_DASHSCOPE_API_KEY&#x27;" >> ~/.zshrc

``` Edit manually

 Open `~/.zshrc`:Copy ```\nnano ~/.zshrc

``` Add the following content to the file:Copy ```\n# Replace YOUR_DASHSCOPE_API_KEY with your API key.
export DASHSCOPE_API_KEY="YOUR_DASHSCOPE_API_KEY"

``` In the nano editor, press Ctrl+X and then Y. Press Enter to save and close the file. Add the variable to `~/.bash_profile`:Copy ```\n# Replace YOUR_DASHSCOPE_API_KEY with your API key.
echo "export DASHSCOPE_API_KEY=&#x27;YOUR_DASHSCOPE_API_KEY&#x27;" >> ~/.bash_profile

``` Edit manually

 Open `~/.bash_profile`:Copy ```\nnano ~/.bash_profile

``` Add the following content to the file:Copy ```\n# Replace YOUR_DASHSCOPE_API_KEY with your API key.
export DASHSCOPE_API_KEY="YOUR_DASHSCOPE_API_KEY"

``` In the nano editor, press Ctrl+X and then Y. Press Enter to save and close the file. 3 Apply the changes

- Zsh 
- Bash 

 Copy ```\nsource ~/.zshrc

``` Copy ```\nsource ~/.bash_profile

``` 4 Verify

Verify in a new session:Copy ```\necho $DASHSCOPE_API_KEY

``` ### Temporary environment variable
Set a temporary variable (current session only): 1 Set the variable

Copy ```\n# Replace YOUR_DASHSCOPE_API_KEY with your API key.
export DASHSCOPE_API_KEY="YOUR_DASHSCOPE_API_KEY"

``` 2 Verify

Copy ```\necho $DASHSCOPE_API_KEY

``` Set environment variables through System Properties, CMD, or PowerShell.- System Properties 
- CMD 
- PowerShell 

 
- Permanent environment variable (requires admin permissions)

- Takes effect in new sessions only — restart terminals, IDEs, and apps


 1 Open System Properties

Press `Win+Q`, search for "**Edit the system environment variables**", and open **System Properties**. 2 Add the environment variable

Click **Environment Variables** > **System Variables** > **New**. Set name to `DASHSCOPE_API_KEY` and value to your API key. 3 Confirm

Click **OK** on all three dialog boxes. 4 Verify

Verify in CMD or PowerShell:
- CMD:


Copy ```\necho %DASHSCOPE_API_KEY%

``` 
- Windows PowerShell:


Copy ```\necho $env:DASHSCOPE_API_KEY

``` #### Permanent environment variable
 1 Set the variable

Copy ```\n# Replace YOUR_DASHSCOPE_API_KEY with your API key.
setx DASHSCOPE_API_KEY "YOUR_DASHSCOPE_API_KEY"

``` 2 Open a new session

Open a new CMD session. 3 Verify

Copy ```\necho %DASHSCOPE_API_KEY%

``` #### Temporary environment variable
 1 Run the following command

Copy ```\n# Replace YOUR_DASHSCOPE_API_KEY with your API key.
set DASHSCOPE_API_KEY="YOUR_DASHSCOPE_API_KEY"

``` 2 Verify

Copy ```\necho %DASHSCOPE_API_KEY%

``` #### Permanent environment variable
 1 Set the variable

Copy ```\n# Replace YOUR_DASHSCOPE_API_KEY with your API key.
[Environment]::SetEnvironmentVariable("DASHSCOPE_API_KEY", "YOUR_DASHSCOPE_API_KEY", [EnvironmentVariableTarget]::User)

``` 2 Open a new session

Open a new PowerShell session. 3 Verify

Copy ```\necho $env:DASHSCOPE_API_KEY

``` #### Temporary environment variable
 1 Run the following command

Copy ```\n# Replace YOUR_DASHSCOPE_API_KEY with your API key.
$env:DASHSCOPE_API_KEY = "YOUR_DASHSCOPE_API_KEY"

``` 2 Verify

Copy ```\necho $env:DASHSCOPE_API_KEY

``` 
## [​ ](#faq) FAQ

### [​ ](#echo-works-but-code-reports-no-api-key-found) echo works but code reports "no API key found"

Common causes:

- 
**Non-permanent variable**: Temporary variables work in the current session only. Set a permanent variable instead.


- 
**Need restart**: Restart IDE, terminal, or app. Service-managed apps may need service restart.


- 
**Service manager config**: For service-managed apps (systemd, supervisord), add the variable to the service config file.


- 
**Using sudo**: `sudo` doesn&#x27;t inherit environment variables. Use `sudo -E python xx.py` (the `-E` flag passes variables) or run without `sudo` if permissions allow.


- 
**Base URL required**: Set Qwen Cloud&#x27;s base URL:

In code:


Copy ```\ndashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

``` 

- As environment variable:


Copy ```\nexport DASHSCOPE_HTTP_BASE_URL=&#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

``` [Previous ](/api-reference/preparation/api-key)[Install the SDK Python and Java setup Next ](/api-reference/preparation/install-sdk)
