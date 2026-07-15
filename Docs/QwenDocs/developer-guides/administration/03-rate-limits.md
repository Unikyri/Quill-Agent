# Rate limits

> **Source:** https://docs.qwencloud.com/developer-guides/administration/rate-limits

Understand and manage API rate limits

 Copy page ## [​ ](#how-rate-limits-work) How rate limits work

Rate limits control how many API requests and tokens your account can consume per minute for each model. There are two types of limits:

- **RPM** (Requests Per Minute): Maximum number of API calls per minute.

- **TPM** (Tokens Per Minute): Maximum number of tokens processed per minute.


Limits are applied at the **account level** — they are shared across all workspaces and API keys under one account.
 Rate limits also apply per second: RPS = RPM / 60, TPS = TPM / 60. Burst requests within a single second can trigger throttling even if total usage stays below the per-minute limit. 
## [​ ](#view-your-rate-limits) View your rate limits

Go to [Analytics](https://home.qwencloud.com/analytics) to see the rate limits and real-time usage for every model in your account.
The **Analytics** page displays:

- **Summary metrics**: Total Models, Total Calls, Failures, Avg Time to First Token, and Avg Latency for the selected time window.

- **Per-model breakdown**: A table showing each model&#x27;s Workspace, Avg TPM, Avg RPM, Total Calls, Failed Calls, Failure Rate, Avg Time to First Token, and Avg Latency.


Use the time range selector (such as 3 Hours, 24 Hours) to adjust the monitoring window.
## [​ ](#set-rate-limits-per-workspace) Set rate limits per workspace

You can set custom RPM and TPM limits for individual models in a [workspace](/developer-guides/administration/workspace).
 1 Go to the Workspaces page

Go to **Settings** > [Workspaces](https://home.qwencloud.com/settings/workspaces) and click **Edit** on a sub-workspace. 2 Add models and set limits

Under **Model Permission**, click **All Models** to add models. For each model, set the **Times / min** (RPM) and **Token / min** (TPM) values, then click **Apply**. 3 Save changes

Click **Save Changes** to apply the new rate limits. 
The RPM and TPM values you set cannot exceed the account-level limits for that model. The default workspace uses account-level limits and cannot be modified.
## [​ ](#temporarily-increase-rate-limits) Temporarily increase rate limits

If you need higher throughput for a specific model, you can request a temporary increase through your account settings.
 1 Go to the Rate Limits page

Go to **Settings** > [Rate Limits](https://home.qwencloud.com/settings/rate-limit). 2 Request an increase

Click **Increase Rate Limit Temporarily**. Select the model, then enter the desired **Token Rate Limit** (tokens per 60 seconds). The dialog shows your current quota and the upper limit. 3 Submit

Click **Submit** to apply the temporary increase. 
 Apply quotas based on actual needs. Unused capacity may be downsized to default limits after a period of inactivity. 
The **Rate Limits** page also shows a history of all temporary increase requests, including the application time, model code, and account TPM limit for each request.
## [​ ](#rate-limit-errors) Rate limit errors

When a rate limit is triggered, the API returns HTTP status code `429`. The error message indicates which limit was hit:
Error messageCause`Requests rate limit exceeded` or `You exceeded your current requests list`RPM limit reached`Allocated quota exceeded` or `You exceeded your current quota`TPM limit reached`Request rate increased too quickly`Sudden request surge triggered stability protection, even if RPM/TPM limits were not reached 
**Limits reset within one minute.** For other errors, see [Error messages](/api-reference/preparation/error-messages).
## [​ ](#best-practices) Best practices

### [​ ](#smooth-your-request-rate) Smooth your request rate

Spread requests evenly over time rather than sending them in bursts. Use constant-rate scheduling, exponential backoff, or a request queue to avoid triggering per-second limits.
### [​ ](#use-a-backup-model) Use a backup model

If a request is rate-limited, fall back to an alternative model to maintain availability:
Copy ```\nimport os
import asyncio
from openai import AsyncOpenAI, APIStatusError

API_KEY = os.getenv("DASHSCOPE_API_KEY")
MODEL = "qwen-plus-2025-07-28"
BACKUP_MODEL = "qwen-plus-2025-07-14"
QUESTION = "Who are you?"
NUM_REQUESTS = 10

client = AsyncOpenAI(
 api_key=API_KEY,
 base_url="https://dashscope-intl.aliyuncs.com/compatible-mode/v1"
)

async def send_request(model):
 try:
 await client.chat.completions.create(
 model=model,
 messages=[{"role": "user", "content": QUESTION}]
 )
 return True
 except APIStatusError as e:
 if e.status_code == 429:
 print(f"[Rate limit triggered] Model {model}")
 return False
 raise
 except Exception as e:
 print(f"[Request failed] Model {model}, Error: {e}")
 return False

async def task(i):
 if await send_request(MODEL):
 return True
 return await send_request(BACKUP_MODEL)

async def main():
 results = await asyncio.gather(*(task(i) for i in range(NUM_REQUESTS)))
 print(f"Successful: {sum(results)}, Failed: {len(results) - sum(results)}")

if __name__ == "__main__":
 asyncio.run(main())

``` 
### [​ ](#split-large-tasks) Split large tasks

Long conversations or large documents consume many tokens quickly. Split large batch tasks into smaller batches and submit them at different times to stay within TPM limits.
### [​ ](#choose-higher-limit-models) Choose higher-limit models

Stable or latest model versions typically have higher rate limits than older snapshots. Where possible, use the latest version of a model.
### [​ ](#use-batch-inference) Use batch inference

If you don&#x27;t need real-time results, use the [Batch API](/developer-guides/text-generation/batch). Batch jobs are not subject to real-time rate limits but may have queuing and processing delays. [Previous ](/developer-guides/administration/workspace)[OpenClaw Open-source AI assistant platform Next ](/developer-guides/clients-and-developer-tools/openclaw)
