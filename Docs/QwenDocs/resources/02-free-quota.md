# Free quota

> **Source:** https://docs.qwencloud.com/resources/free-quota

New user free quota

 Copy page ## [​ ](#rules) **Rules**

### [​ ](#validity-period) **Validity period**

Free quota is typically valid for 90 days, but the start date depends on the quota type:

- **New user free quota**: When you activate your account, all existing models receive a free quota. The 90-day period starts from your **account activation date**.

- **New model free quota**: For models released after you sign up, the 90-day period starts from the **model&#x27;s release date**.


After expiration or depletion, continued model inference will incur [charges](/developer-guides/getting-started/pricing).
### [​ ](#scope-of-application) **Scope of application**

Free quota only offsets costs from real-time model inference (invocation). It does not offset:

- [Batch calls](/developer-guides/text-generation/batch)

- [Built-in tool call fees](/developer-guides/getting-started/pricing#built-in-tools) (web search, image search, etc.)

- Model fine-tuning

- Model deployment

- Custom models (fine-tuned models, deployed models)


## [​ ](#get-the-free-quota) Get the free quota

[Sign up](https://home.qwencloud.com/) for Qwen Cloud to activate your free quota — no payment method required. You can try models immediately after signing up. After your free quota is exhausted, you will be guided to add a payment method to continue using paid services.
 If you did not receive your free quota after signing up, go to the [Free Tier](https://home.qwencloud.com/benefits) page and follow the on-screen prompt: "Please complete your payment information to unlock your free quota and secure your account." Once completed, your free quota will be granted automatically. 
## [​ ](#view-your-remaining-quota) View your remaining quota

Go to the [Free Tier](https://home.qwencloud.com/benefits) page to view and manage your free quota across all models.
The page provides:

- **Eligible models**: Number of models with available free quota

- **Models expiring soon**: Models with quota expiring within 7 days

- **Low balance models**: Models with 80% or more of quota consumed

- **Unavailable models**: Models with free quota expired or exhausted (data from past 180 days)


Switch between model types (Embedding, LLM, Multimodal, Audio, Vision) and view details including:

- Model Code, Free Quota, Consumed, Remaining, Expiration, Status

- Toggle **Free quota only** switch in the **Actions** column to enable auto-stop when quota ends


 The free quota displayed on the console is subject to a delay of several minutes and is not real-time data. 
## [​ ](#use-the-free-quota) **Use the free quota**

When you make real-time calls to a model, free quota is automatically deducted. For more information, see [Get started with Qwen Cloud](/developer-guides/getting-started/first-api-call).
 Unverified users cannot continue calling models after free quota is exhausted. Requests are rejected with error code `AllocationQuota.FreeTierOnly`. You can click **Set up payment** on the [Free Tier](https://home.qwencloud.com/benefits) page to complete verification and switch to pay-as-you-go billing.Verified users are billed on a pay-as-you-go basis after free quota is exhausted. You can enable [Free quota only](#free-quota-only) in advance to prevent unexpected charges. 
## [​ ](#free-quota-only) **Free quota only**

When Free quota only is enabled, model calls are blocked after quota exhaustion and no charges are incurred.

- **Unverified users** (users who have not [set up payment](https://home.qwencloud.com/benefits)): This feature is enabled by default and cannot be disabled.

- **Verified users**: This feature is disabled by default. You can enable or disable it on the [Free Tier](https://home.qwencloud.com/benefits) page.


 Unverified users will not see the **Free quota only** switch on the console. The feature is automatically enabled for these users. 
#### [​ ](#how-to-enable) **How to enable**

Go to the [Free Tier](https://home.qwencloud.com/benefits) page:

- **Single model**: Toggle on the **Free quota only** switch for the target model

- **Multiple models**:

Click the **Auto-stop when free quota runs out** dropdown at the top-right

- Select **Enable selected models** to enter bulk mode (checkboxes appear in the table)

- Check the boxes next to the models you want to enable, then click **Enable selected models**; or click **Enable all models** directly without selecting any

- Click **Exit bulk mode** when done


If the switch is not displayed for a model, it means the free quota for that model has been exhausted or has expired, or the model does not offer a free quota.
#### [​ ](#how-to-disable) **How to disable**


- **Unverified users**: This feature cannot be disabled. [Set up payment](https://home.qwencloud.com/benefits) to unlock this option.

- **Verified users**: If enabled, you can disable it on the [Free Tier](https://home.qwencloud.com/benefits) page.


## [​ ](#faq) **FAQ**

### [​ ](#will-i-be-notified-when-my-free-quota-is-used-up) **Will I be notified when my free quota is used up?**

Currently, there is no notification mechanism.
### [​ ](#what-happens-when-my-free-quota-is-used-up) **What happens when my free quota is used up?**

**Unverified users**: You cannot continue calling models after free quota is used up. You can click **Set up payment** on the [Free Tier](https://home.qwencloud.com/benefits) page to switch to pay-as-you-go billing.
**Verified users**:

- If [Free quota only](#free-quota-only) is enabled, model calls are blocked after free quota is used up. Disable this feature to switch to pay-as-you-go billing.

- If [Free quota only](#free-quota-only) is not enabled, model calls in progress complete without interruption. Tokens exceeding free quota are billed based on input/output costs in [Model invocation pricing](/developer-guides/getting-started/pricing). Charges are automatically deducted on a pay-as-you-go basis, which may result in overdue payment.


When your account has an overdue payment, all model calls are blocked — even if other models still have free quota.
Before calling a model, check its free quota and use [budget management](/resources/bill-query).
### [​ ](#why-was-i-charged) **Why was I charged?**

**Possible reasons:**

- You used a model without free quota.

- The free quota cannot be used to offset costs from [Batch](/developer-guides/text-generation/batch) calls.

- Free quota data is subject to a delay of several minutes. The display might show remaining quota when it&#x27;s actually exhausted, resulting in charges. Check status again later for latest data.


You can confirm the charge details on the [Pay-as-you-go](/resources/bill-query) page.
### [​ ](#how-do-i-avoid-charges) **How do I avoid charges?**

If you have not added a payment method, no charges can occur — your usage simply stops when free quota is exhausted.
If you have added a payment method, charges are automatically deducted from your payment method after free quota is exhausted. To manage charge risk:

- **Delete API keys**: Go to the Qwen Cloud [API-Key](https://home.qwencloud.com/api-keys) page and delete created API keys. After you delete an API key, you can no longer call models using the API, which prevents further charges.

- **Set spending alert**: Configure a [spending alert](https://home.qwencloud.com/billing/overview). You&#x27;ll receive email notifications when monthly spending exceeds the threshold.


### [​ ](#i-have-a-remaining-quota-so-why-did-my-call-fail) **I have a remaining quota, so why did my call fail?**

Your account may have an overdue payment. After the grace period expires, all model calls are blocked — even if you have remaining free quota. See [Overdue payment protection](/resources/overdue-payment-protection) for details.
### [​ ](#why-cant-i-see-the-free-quota-and-its-validity-period) **Why can&#x27;t I see the free quota and its validity period?**

If the **Free Quota** column shows **No free quota** or the section is not displayed, your free quota for the model has expired. [Previous ](/resources/billing-overview)[Billing and cost management Analyze and manage costs Next ](/resources/bill-query)
