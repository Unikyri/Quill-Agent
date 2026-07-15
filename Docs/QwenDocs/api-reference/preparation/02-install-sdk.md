# Install the SDK

> **Source:** https://docs.qwencloud.com/api-reference/preparation/install-sdk

Python and Java setup

 Copy page Qwen Cloud provides DashScope SDKs (Python, Java) and supports OpenAI-compatible calls. OpenAI provides SDKs for Python, Node.js, Java, and Go.
## [​ ](#prepare-the-environment) Prepare the environment

Skip this section if you already have Python, Java, Node.js, or Go configured locally.
- Python 
- Java 
- Node.js 
- Go 

 ### Check your Python version
Check if Python and pip are installed:Copy ```\npython -V
pip --version

``` Python 3.8 or higher required. Install from [python.org](https://www.python.org/downloads/) if needed. #### `python -V` or `pip --version` returns "command not found"?
- Windows 
- Linux and macOS 

 
- 
Install Python and add it to PATH. See [Install Python](https://www.python.org/downloads/). 


- 
If the error persists after installing Python and setting PATH, restart your terminal and try again.


 
- 
Install Python. See [Install Python](https://www.python.org/downloads/).


- 
If the error persists, check if `python` and `pip` exist:


If output shows `/usr/bin/python` and `/usr/bin/pip`, restart your terminal.


- 
If "no python" appears, try `which python3 pip3`:


Copy ```\n/usr/bin/which: no python in (/root/.local/bin:/root/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin)
/usr/bin/which: no pip in (/root/.local/bin:/root/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin)

``` If found at `/usr/bin/python3` and `/usr/bin/pip3`, use `python3 -V` and `pip3 --version`.Copy ```\n/usr/bin/python3
/usr/bin/pip3

``` ### Configure a virtual environment (optional)
Create a virtual environment to isolate SDK dependencies (recommended).
- **Create a virtual environment**


Create a virtual environment named **.venv**:Copy ```\n# If the command fails to run, you can replace python with python3 and run it again.
python -m venv .venv

``` 
- **Activate the virtual environment**


Activate the virtual environment:
- Windows:


Copy ```\n.venv\Scripts\activate

``` 
- macOS/Linux:


Copy ```\nsource .venv/bin/activate

``` ### Check your Java version
Run the following command in your terminal:Copy ```\njava -version
# (Optional) If using Maven, verify it&#x27;s installed
mvn --version

``` Java 8 or higher required. Install from [Java Downloads](https://www.oracle.com/java/technologies/downloads/) if needed. ### Check Node.js installation
Check if Node.js and npm are installed:Copy ```\nnode -v
npm -v

``` If Node.js is not installed, download it from the [Node.js official website](https://nodejs.org/en/download/package-manager). ### Check your Go version
Run the following command in your terminal:Copy ```\ngo version

``` If Go is not installed, download it from the [official Go website](https://go.dev/doc/install). The OpenAI Go SDK requires Go 1.22 or higher. ### Create a project and initialize the module
Create project and initialize module:Copy ```\n# Create a project folder (adjust path/commands for your OS)
mkdir D:\your_project_folder && cd /d D:\your_project_folder

# Initialize the module. example.com is an example. Any name in this format is acceptable. A real domain name is not required.
go mod init example.com/your_project_folder

``` 
## [​ ](#install-the-sdk) Install the SDK

- Python 
- Java 
- Node.js 
- Go 

 You can call Qwen Cloud APIs using the OpenAI Python SDK or the DashScope Python SDK.### Install the OpenAI Python SDK
Install the OpenAI Python SDK:Copy ```\n# If the command fails to run, you can replace pip with pip3 and run it again.
pip install -U openai

``` Success message: `Successfully installed ... openai-x.x.x`### Install the DashScope Python SDK
Install the DashScope Python SDK:Copy ```\n# If the command fails to run, you can replace pip with pip3 and run it again.
pip install -U dashscope

``` Success message: `Successfully installed ... dashscope-x.x.x` Pip version warnings are safe to ignore. ### Install the DashScope Java SDK
Add the [DashScope Java SDK](https://mvnrepository.com/artifact/com.alibaba/dashscope-sdk-java) dependency. Replace `the-latest-version` with the latest version number.#### XML

- 
Open the `pom.xml` file of your Maven project.


- 
Add the following dependency information within the `<dependencies>` tag.


Copy ```\n<dependency>
 <groupId>com.alibaba</groupId>
 <artifactId>dashscope-sdk-java</artifactId>
 <!-- Replace &#x27;the-latest-version&#x27; with the latest version number from: https://mvnrepository.com/artifact/com.alibaba/dashscope-sdk-java -->
 <version>the-latest-version</version>
</dependency>

``` 
- 
Save the `pom.xml` file.


- 
Run a Maven command such as `mvn compile` or `mvn clean install` to update dependencies. Maven automatically downloads and adds the DashScope Java SDK to your project.


Example in IntelliJ IDEA: #### Gradle

- 
Open the `build.gradle` file of your Gradle project.


- 
Add the following dependency information within the `dependencies` block.


Copy ```\ndependencies {
 // Replace &#x27;the-latest-version&#x27; with the latest version number from: https://mvnrepository.com/artifact/com.alibaba/dashscope-sdk-java
 implementation group: &#x27;com.alibaba&#x27;, name: &#x27;dashscope-sdk-java&#x27;, version: &#x27;the-latest-version&#x27;
}

``` 
- 
Save the `build.gradle` file.


- 
Run this command from your project root to update dependencies:


Copy ```\n./gradlew build --refresh-dependencies

``` Example in IntelliJ IDEA: ### Install the OpenAI Java SDK
Add the OpenAI Java SDK dependency.#### XML

- 
Open the `pom.xml` file of your Maven project.


- 
Add this dependency (version 4.30.0+) within the `<dependencies>` tag:


Copy ```\n<dependency>
 <groupId>com.openai</groupId>
 <artifactId>openai-java</artifactId>
 <version>4.30.0</version>
</dependency>

``` 
- 
Save the `pom.xml` file.


- 
Use a Maven command, such as `mvn compile` or `mvn clean install`, to update the project dependencies. Maven automatically downloads and adds the OpenAI Java SDK to your project.


 Install via npm or yarn:Copy ```\nnpm install --save openai
# or
yarn add openai

``` If installation fails, configure a mirror source:Copy ```\nnpm config set registry https://registry.npmmirror.com/

``` Then retry installation. Success message: `added xx package in xxs`. Check version: `npm list openai` Install the OpenAI Go SDK:Copy ```\ngo get github.com/openai/openai-go/v3@v3.30.0

``` Success message: `go: added github.com/openai/openai-go/v3 v3.30.0` 
- 
Version `v3.30.0` has been tested and is considered stable.


- 
The SDK is in testing phase.


- 
If server times out, use our mirror:


Copy ```\n# Set our mirror
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

``` 
## [​ ](#next-steps) Next steps

Run the code examples for the OpenAI SDK or the DashScope SDK. [Previous ](/api-reference/preparation/export-api-key-env)[Error messages API error code reference Next ](/api-reference/preparation/error-messages)
