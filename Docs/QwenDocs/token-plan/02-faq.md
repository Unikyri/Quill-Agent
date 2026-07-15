# FAQ

> **Source:** https://docs.qwencloud.com/token-plan/faq

Fix Token Plan issues fast

 Copy page ## [​ ](#what-is-the-difference-between-token-plan-and-coding-plan) What is the difference between Token Plan and Coding Plan?

Token PlanCoding PlanUse casesDay-to-day work for solo founders, teams, and enterprisesIndividual development scenariosSupported modelsText generation and image generation modelsText generation modelsBilling methodCredits deducted based on token consumptionBilled per model callUsage frequencyUnlimitedQuota per 5 hours / per weekAPI Key and Base URLGet from the [Token Plan page](https://home.qwencloud.com/billing/subscription/token-plan). See [Quick start](/token-plan/quickstart)Get from the Coding Plan pagePeak-time performanceMulti-tenant isolationRequests may queue during peak hoursData securityYour data is not used to train modelsUser data authorization required 
## [​ ](#important-resources) Important resources


- **Token Plan page**: Manage your subscription and view usage on the [Token Plan page](https://home.qwencloud.com/billing/subscription/token-plan)

- **API Keys page**: View and reset your API keys on the [API Keys page](https://home.qwencloud.com/api-keys)

- **Supported models**: See the [full list of models](/token-plan/overview#supported-models)

- **Setup guides**: Configuration instructions for [Claude Code](/developer-guides/clients-and-developer-tools/claude-code), [OpenCode](/developer-guides/clients-and-developer-tools/opencode), [OpenClaw](/developer-guides/clients-and-developer-tools/openclaw), and more


## [​ ](#start-here-60-second-diagnosis) Start here: 60-second diagnosis

If your tool is not working and you see errors like `401`, `403`, `404`, `invalid access token`, `invalid api-key`, `quota exceeded`, or connection failures, check these first:

- **API key**: Must be the Token Plan dedicated key, not a pay-as-you-go `sk-` key or a Coding Plan `sk-sp-` key

- **Base URL**: Must contain `token-plan`, not `dashscope-intl` or `coding`

- **Authentication header**: Must use `Authorization: Bearer`, not `x-api-key`

- **Subscription status**: Check whether your subscription is active and credits are not depleted

- **Updated config**: Restart the tool after updating the API key or base URL


If you are still blocked, continue with the questions below.
## [​ ](#connection-and-usage) Connection and usage

### [​ ](#which-text-generation-models-are-supported) Which text generation models are supported?

Token Plan supports the following text generation models:

- qwen3.7-plus

- qwen3.7-max

- qwen3.6-plus

- qwen3.6-flash

- deepseek-v4-pro

- deepseek-v4-flash

- deepseek-v3.2

- kimi-k2.6

- kimi-k2.5

- glm-5.1

- glm-5

- MiniMax-M2.5


For the complete and up-to-date list, see [Supported models](/token-plan/overview#supported-models).
### [​ ](#common-errors-and-solutions) Common errors and solutions

Error messagePossible causesSolution**401 InvalidApiKey: No API-key provided.**The request did not include the API key in either the `Authorization: Bearer` or `x-api-key` header.Go to the [Token Plan page](https://home.qwencloud.com/billing/subscription/token-plan), copy your dedicated API key, and configure it in your tool.**401 InvalidApiKey: Invalid API-key provided.**1. You used a Qwen Cloud pay-as-you-go key (`sk-xxx`) or a Coding Plan key (`sk-sp-xxx`). 2. Your Token Plan subscription has expired. 3. The API key is incomplete or contains spaces.1. Confirm you are using the Token Plan dedicated API key. Ensure the key is copied completely without any spaces. 2. Verify whether the subscription has expired. 3. If the error persists, reset the API key, then configure the new key.**400 Model not exist.**The model name does not match any available model.Check the model name against [Supported models](/token-plan/overview#supported-models) and ensure it is spelled exactly as documented.**model &#x27;xxx&#x27; not found or not supported**The model name is misspelled or uses the wrong case.Ensure the model name is case-sensitive and matches the model ID listed in [Supported models](/token-plan/overview#supported-models).**400 InvalidParameter: url error, please check url!**The Base URL path does not match the protocol. For example, the OpenAI-compatible path is configured on the Anthropic endpoint, or vice versa.Choose the endpoint matching the protocol your tool uses: Anthropic-compatible (Claude Code and similar) ends with `/apps/anthropic`. OpenAI-compatible (Cursor, Qwen Code, and similar) ends with `/compatible-mode/v1`.**400 InvalidParameter: Range of max_tokens should be [1, xxxx]**The `max_tokens` value in the request (or the maximum output length in the tool configuration) exceeds the maximum output tokens supported by the current model.Set `max_tokens` to a value not exceeding the limit shown in the error message.**400 invalid_parameter_error: The thinking_budget parameter must be a positive integer and not greater than xxxxx**The thinking length configured in the tool (such as `thinking_budget` or `budgetTokens`) exceeds the limit supported by the current model.Set the thinking length to a value not exceeding the limit shown in the error, or remove the parameter for models that do not support thinking mode.**invalid access token or token expired**You used the base URL for the Coding Plan or another plan.Use the correct Token Plan base URL. OpenAI compatible: `https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1`. Anthropic compatible: `https://token-plan.ap-southeast-1.maas.aliyuncs.com/apps/anthropic`.**Incorrect API key provided**You used the general Qwen Cloud base URL (`dashscope-intl.aliyuncs.com`).Use the correct Token Plan base URL. OpenAI compatible: `https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1`. Anthropic compatible: `https://token-plan.ap-southeast-1.maas.aliyuncs.com/apps/anthropic`.**Range of input length should be [1, xxx]**Your input (including conversation history and code context) exceeds the model&#x27;s maximum context length.Start a new session or shorten the context. You can also switch to a model with a larger context window.**429 API-Key Requests rate limit exceeded**Requests were too dense within a short period and triggered the model rate limit.Wait one minute and retry. If this occurs frequently, reduce the request frequency and verify the API key is not being shared with others.**429 Throttling.AllocationQuota / insufficient_quota**This error can be triggered by either of two causes: **Plan quota depleted**: both the seat quota and the shared quota pack are exhausted. **Model rate limit triggered**: even when the plan quota is sufficient, this error occurs when the tokens consumed per second or per minute (TPS/TPM) exceed the model&#x27;s rate limit threshold. Rate limits are calculated at the Alibaba Cloud account level — the combined usage of all RAM users, workspaces, and API keys under the account counts together. Even if the per-minute total stays within the limit, a short burst of requests can still trigger it.**Quota depleted**: wait for the quota to reset at the next billing cycle, purchase a shared usage package, or purchase additional seats to add more credits. **Rate limit triggered**: wait about one minute and retry, and smooth your request rate (for example, use even scheduling, exponential backoff, or a request queue) to avoid instantaneous spikes.**Connection error**The base URL domain is misspelled, or you have a network connection issue.Check the spelling of the base URL domain and your network connection. 
## [​ ](#image-generation) Image generation

### [​ ](#which-image-generation-models-are-supported) Which image generation models are supported?

Token Plan supports the following image generation models:

- qwen-image-2.0

- qwen-image-2.0-pro

- wan2.7-image

- wan2.7-image-pro


For details, see [Supported models](/token-plan/overview#supported-models).
### [​ ](#can-i-use-image-models-in-coding-tools) Can I use image models in coding tools?

Image generation models use separate APIs and cannot be called directly through the text model&#x27;s base URL. You must integrate them through your tool&#x27;s skill or extension mechanism. For configuration steps, see the image generation model integration section in each tool&#x27;s setup guide.
### [​ ](#how-do-i-use-image-generation-models-in-coding-tools) How do I use image generation models in coding tools?

Most coding tools (such as Claude Code, Cursor, and Cline) support extension mechanisms like Skills, Slash Commands, or MCP servers. You can use these mechanisms to call image generation models within your coding workflow. For step-by-step instructions, see [Integrate multimodal generation models](/token-plan/best-practices/integrate-multimodal-gen).
## [​ ](#common-errors-fixes) Common errors & fixes

### [​ ](#which-api-key-and-base-url-should-i-use) Which API key and base URL should I use?

Use a **Token Plan dedicated key** and a **Token Plan base URL**:

- API key: Get it from the [Token Plan page](https://home.qwencloud.com/billing/subscription/token-plan)

- OpenAI compatible: `https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1`

- Anthropic compatible: `https://token-plan.ap-southeast-1.maas.aliyuncs.com/apps/anthropic`


Pay-as-you-go `dashscope-intl` URLs and Coding Plan `coding-intl` URLs do not work with Token Plan keys.
### [​ ](#i-got-401-or-403-is-my-api-key-invalid) I got 401 or 403. Is my API key invalid?

If you see `401 invalid access token`, `InvalidApiKey`, or `403 invalid api-key`, the problem is usually a wrong key type, an expired subscription, a malformed key, or a key and URL mismatch.
Check these points:

- Confirm the key is your Token Plan dedicated key

- Re-copy the key to remove hidden spaces or line breaks

- Confirm your subscription is active on the [Token Plan page](https://home.qwencloud.com/billing/subscription/token-plan)

- Make sure the base URL contains `token-plan`


In most cases:

- `401 invalid access token` means wrong key type, expired subscription, or malformed key

- `403 invalid api-key` means you are using a Token Plan key with a non-Token Plan URL

- `Incorrect API key provided` means you are using a Token Plan key with the general `dashscope-intl` URL


### [​ ](#why-do-i-get-404-or-connection-errors) Why do I get 404 or connection errors?

If you get `404 status code (no body)`, endpoint not found, or connection errors, the base URL path is usually wrong for your tool.
Use the correct path for your protocol:

- Anthropic-compatible tools use `/apps/anthropic`

- OpenAI-compatible tools use `/compatible-mode/v1`


Also check that the domain is exactly `token-plan.ap-southeast-1.maas.aliyuncs.com`.
### [​ ](#why-is-my-model-not-supported) Why is my model "not supported"?

If you get `model &#x27;xxx&#x27; is not supported` or `model &#x27;xxx&#x27; not found`, the model ID is usually incorrect, case-sensitive, or copied with extra spaces.
Check these points:

- Use the exact model name from [Supported models](/token-plan/overview#supported-models)

- Remove any leading or trailing spaces

- Verify capitalization exactly as documented


### [​ ](#how-to-handle-input-length-limits) How to handle input length limits?

If you get `400 InvalidParameter: Range of input length should be [1, xxx]`, your context or request exceeds model limits.
Try these fixes:

- Create a new session

- Switch to a model with a larger context window

- In your tool, configure context length limits if available


### [​ ](#what-should-i-do-if-i-see-data-inspection-failed) What should I do if I see data_inspection_failed?

If you see `data_inspection_failed`, refer to [Troubleshooting data inspection errors](/api-reference/preparation/error-messages#400-datainspectionfailed-data-inspection-failed).
## [​ ](#product-features) Product features

### [​ ](#can-i-mix-token-plan-coding-plan-and-pay-as-you-go-keys) Can I mix Token Plan, Coding Plan, and pay-as-you-go keys?

No. The API keys and base URLs for Token Plan, Coding Plan, and Qwen Cloud pay-as-you-go are not interchangeable. Do not mix them. Using an incorrect API key will not deduct credits from your Token Plan quota.
### [​ ](#can-i-use-one-subscription-across-multiple-tools) Can I use one subscription across multiple tools?

Yes. You can use the same API key in all compatible AI tools, and all usage draws from the same quota.
### [​ ](#are-there-usage-restrictions) Are there usage restrictions?

The service is for interactive use in compatible AI programming and agent tools only. It cannot be used for automated scripts or application backends. Violating these terms may result in subscription suspension or revocation of your API key.
### [​ ](#are-the-models-full-featured-or-quantized) Are the models full-featured or quantized?

All supported models in Token Plan are full-featured, unquantized versions.
### [​ ](#how-do-i-reset-or-manage-my-api-key) How do I reset or manage my API key?

Use the Reset button on the [API Keys page](https://home.qwencloud.com/api-keys) to reset your key.
Keep these rules in mind:

- After reset, update the key in all tools

- If you resubscribe after expiration, the key changes

- Token Plan supports one key per subscription


### [​ ](#what-is-the-difference-between-revoking-a-seat-changing-a-role-and-removing-a-member-from-the-organization) What is the difference between revoking a seat, changing a role, and removing a member from the organization?


- **Revoke seat**: Cancels the member&#x27;s usage rights and releases the seat back to the available pool. The member loses access to the seat but remains in the organization.

- **Change role**: Updates the member&#x27;s permissions (for example, between administrator and regular member) without affecting seat assignment or organization membership.

- **Remove from organization**: Completely removes the member from the team. The seat is automatically reclaimed and the member is removed from the member list.


### [​ ](#why-can-the-api-key-only-be-viewed-once-and-what-should-i-do-if-i-lose-it) Why can the API key only be viewed once, and what should I do if I lose it?

To prevent API key misuse across teams and avoid billing confusion, the API key is displayed only when it is first generated or reset — it cannot be viewed or copied again afterward. If you lose the API key, go to the **Member Management** page and click **Reset** to generate a new key. The original key becomes invalid immediately.
### [​ ](#what-are-the-restrictions-on-plan-conversion-and-resubscription) What are the restrictions on plan conversion and resubscription?

**Plan conversion**: Token Plan (Team Edition) and Coding Plan are two independent subscription plans and cannot be converted to each other. Switching from a purchased Token Plan (Team Edition) to a Coding Plan, or vice versa, is not supported, even by paying a price difference. However, you can subscribe to both plans simultaneously under the same account.
**Resubscription after cancellation**: If you cancel and repurchase a subscription, the API key and base URL will change. You must reconfigure the new dedicated API key in your tools to resume normal usage.
## [​ ](#billing-and-quota) Billing and quota

### [​ ](#why-am-i-still-billed-after-subscribing) Why am I still billed after subscribing?

If you still see pay-as-you-go charges after subscribing to Token Plan, it is usually one of the following:

- **Wrong credentials (most common)**: Your tool is using a pay-as-you-go API key (`sk-` format) and a base URL without `token-plan`. Use your Token Plan dedicated key and ensure the base URL contains `token-plan`.

- **Tool falling back to general credentials**: If your tool has both general and Token Plan credentials configured, it may route requests through the general credentials. Remove the general API configuration and keep only Token Plan credentials.

- **Client cache not cleared**: After updating credentials, clear the tool&#x27;s cache and restart it to ensure it uses the new configuration.


Use the [Token Plan page](https://home.qwencloud.com/billing/subscription/token-plan) to verify usage. If needed, reset your key on the [API Keys page](https://home.qwencloud.com/api-keys).
### [​ ](#how-are-credits-deducted) How are credits deducted?

Credit consumption is calculated based on the input, cached, and output tokens in each request. The system first deducts credits from your seat quota. After the seat quota is depleted, credits are deducted from the shared usage package. If all credits are used, the service is suspended until the next billing cycle or until you purchase a shared usage package to add more credits.
### [​ ](#what-happens-when-my-quota-runs-out) What happens when my quota runs out?

Once your seat quota is depleted, the system automatically begins deducting from the shared usage package. If you run out of credits entirely, the service is suspended. You can resume the service in one of the following ways:

- Purchase a shared usage package to add more credits

- Wait for your quota to automatically reset at the start of the next billing cycle


There is no automatic fallback to pay-as-you-go when credits are exhausted.
### [​ ](#when-do-my-credits-reset) When do my credits reset?

The seat quota resets at the start of each new subscription month. Renewing or resubscribing mid-cycle does not immediately add or reset credits — the renewed quota becomes available only when the next billing cycle begins. Unused credits do not roll over to the next month. Similarly, unused credits in the shared usage package expire at the end of each billing cycle.
### [​ ](#i-just-renewed-my-plan-but-my-credits-are-still-0-why) I just renewed my plan but my credits are still 0. Why?

Credits reset only at the start of each subscription month, not at the moment of renewal. If you renew mid-cycle, your renewal extends the subscription, but the new quota does not take effect until the next billing cycle begins. If your current cycle quota is exhausted and you need to restore service immediately, you can only do so through the following options:

- Purchase a shared usage package

- Upgrade to a higher-tier seat

- Purchase additional seats


Subscribing to quota for the next billing cycle does not immediately restore service.
### [​ ](#how-can-i-view-my-usage) How can I view my usage?

Check the [Token Plan page](https://home.qwencloud.com/billing/subscription/token-plan) for overall usage, including seat quota and shared usage package details.
### [​ ](#can-i-upgrade-or-downgrade-my-plan) Can I upgrade or downgrade my plan?

No. Switching between seat tiers is not supported. If you need more credits, purchase a shared usage package.
### [​ ](#can-i-buy-multiple-subscriptions) Can I buy multiple subscriptions?

No. Each Qwen Cloud account is limited to one Token Plan subscription.
### [​ ](#can-i-get-a-refund) Can I get a refund?

No. Token Plan does not support refunds. Please review the [Usage rules](/token-plan/overview#usage-rules) for details before you subscribe.
### [​ ](#does-an-overdue-payment-affect-my-plan) Does an overdue payment affect my plan?

Token Plan is a prepaid subscription product. As long as your plan&#x27;s credits are not depleted and the subscription is still valid, an overdue payment on your Qwen Cloud account will not affect your use of the service.
### [​ ](#can-multiple-people-share-one-key) Can multiple people share one key?

No. Token Plan API keys are for the subscriber&#x27;s personal use only. Public exposure or account sharing may result in automatic key suspension.
## [​ ](#data-security) Data security

### [​ ](#how-is-my-data-protected) How is my data protected?

Token Plan does not use your conversation data for model training. [Previous ](/token-plan/team-management)[Migrate from Coding Plan Transition from Coding Plan to Token Plan Next ](/token-plan/migrate-from-coding-plan)
