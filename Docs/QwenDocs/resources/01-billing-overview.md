# Billing overview

> **Source:** https://docs.qwencloud.com/resources/billing-overview

View bills, manage payments, and track spending

 Copy page ## [​ ](#overview) Overview

The [Billing Overview](https://home.qwencloud.com/billing/overview) page shows your overall consumption and outstanding charges for Pay-as-you-go, Token Plan, and Coding Plan billing modes.
The page contains these sections:

- **Total Spend** — View your total spending for the selected month, including tax-exclusive amount and tax

- **Total Amount Due** — View outstanding balance and make payments

- **Spending Trend** — Visualize spending history with filterable charts


## [​ ](#bill-settlement-and-repayment) Bill settlement and repayment

For pay-as-you-go resources, pay your bills promptly to avoid them becoming overdue and prevent business interruptions.
### [​ ](#automatic-payment) Automatic payment

By default, the system automatically deducts payment from your default payment method. There is no fixed repayment time.

- If you link a bank card, PayPal, Google Pay, or Apple Pay as your default payment method, the system automatically initiates a deduction when your accumulated pay-as-you-go fees reach a threshold:

**Bank cards**: Threshold is USD 1,000

- **PayPal**: Threshold varies from USD 8 to USD 500, depending on the account


- At the end of each month, the system makes a single deduction for any bills that have not yet reached the threshold


 When you use PayPal as your default payment method, a pre-authorization check is triggered the first time you activate a pay-as-you-go product. 
### [​ ](#manual-payment) Manual payment

If automatic payment fails (for example, due to insufficient available credit) or if you prefer to pay manually, you can manually repay your outstanding bills.

- 
**Go to Qwen Cloud**: On the **Billing Overview** page, click **Pay Now** in the **Total Amount Due** section.


- 
**Set the amount**: In the **Repay Outstanding Amount** dialog, you can see the total outstanding amount and modify it if needed.


- 
**Select a payment method**: Choose from the available payment methods:

**Credit & Debit Cards** — Supports Visa, Mastercard, American Express, UnionPay, Diners Club, Discover, and JCB

- **PayPal** — Associate your account with your PayPal account

- **Google Pay** — One-tap payment with Google Pay

- **Apple Pay** — One-tap payment with Apple Pay


- 
**Complete the payment**: Click **Confirm Repayment** to complete the payment.


- 
**View the result**: The bill payment status takes about 30 minutes to update. You can refresh the page to check the status. Do not make a duplicate payment during this time.


### [​ ](#add-a-payment-method) Add a payment method

To enable automatic payments or make manual payments, you first need to add a payment method to your Qwen Cloud account.

- 
On the **Billing Overview** page, click **Set Payment Method** in the **Total Amount Due** section.


- 
In the **Payment Methods** dialog, select a payment method:

**Credit & Debit Cards** — Supports Visa, Mastercard, American Express, UnionPay, Diners Club, Discover, and JCB

- **PayPal** — Associate your account with your PayPal account

- **Google Pay** — Associate your account with Google Pay

- **Apple Pay** — Associate your account with Apple Pay


#### [​ ](#link-a-bank-card) Link a bank card

Before linking a bank card, ensure it meets these conditions:

- The bank card is enabled for online payments and international transactions.

- The card has a sufficient balance or credit limit to cover a $1.00 USD pre-authorization.

- The bank card is enabled for 3D Secure authentication. In some countries or regions, such as Malaysia and EU member states, 3D Secure authentication is required when linking a bank card due to regulatory requirements.


- 
Fill in your card information:

**Card number** — Enter your bank card number

- **Exp date (MM/YY)** — Expiration date. You cannot add a card that expires in the current month.

- **CVV** — Security code (3 or 4 digits on the back of the card)

- **Name on card** — Must match the name on the bank card


- 
Click **Confirm**. The system initiates a $1.00 USD pre-authorization to verify the card. Your bank statement shows ALIBABACLOUD.COM. The issuing bank usually cancels it automatically within one to five business days.


- 
After successful pre-authorization, the card is linked. The first linked bank card is set as the **Preferred External Payment Method** by default. You can use it once its **Payment Method Status** changes to "Active".


 **Service fee**Qwen Cloud does not charge a service fee for bank card payments. If the bank card&#x27;s settlement currency differs from the transaction currency, the issuing bank may charge a foreign currency conversion fee. Contact your issuing bank for details. 
 **3D Secure (3DS) authentication**3D Secure is an online payment security protocol. It triggers authentication when you **link a bank card** or **place an order** (purchase, upgrade, or renew instances) if the system detects potential security risks.When triggered, you are automatically redirected to your issuing bank&#x27;s verification page. Complete identity verification as prompted—for example, enter a dynamic verification code or password.Keep the Qwen Cloud page open during verification. Do not close the original window. If your browser blocks pop-ups, configure it to allow them.After successful verification, you are automatically redirected back to Qwen Cloud. Continue with your operation. 
#### [​ ](#link-a-paypal-account) Link a PayPal account


- 
Click **Next**. You are redirected to the PayPal website.


- 
Log on to your PayPal account and confirm the payment agreement. The system automatically adds PayPal as a payment method and initiates a $1.00 USD pre-authorization.


 
- If your PayPal account has no bank card or other funding source, PayPal may require you to add one. For related issues, contact PayPal customer service.

- If PayPal authorization succeeds but Qwen Cloud does not show successful linking, the payment agreement may not have been approved. Contact PayPal customer service to confirm.


 
 **One payment method per account**Each bank card, PayPal, Google Pay, or Apple Pay account can be linked to only one Qwen Cloud account. If your payment method is already linked to another account (including a deactivated account still in its cool-down period), you must first unlink it from that account. 
## [​ ](#related-resources) Related resources


- **[Model pricing](/developer-guides/getting-started/pricing)** — Compare Pay-as-you-go, Token Plan, and Coding Plan options, view rate tables for all models

- **[Free quota](/resources/free-quota)** — Check eligibility and remaining free quota

- **[Billing and cost management](/resources/bill-query)** — View detailed bills, set spending limits, and manage costs

- **[Overdue payment protection](/resources/overdue-payment-protection)** — Learn about service suspension protection during overdue payments


 Previous [Free quota New user free quota Next ](/resources/free-quota)
