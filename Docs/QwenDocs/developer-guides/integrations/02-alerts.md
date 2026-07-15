# Monitoring alerts

> **Source:** https://docs.qwencloud.com/developer-guides/integrations/alerts

Configure alert rules for model API metrics

 Copy page Configure custom alert rules to monitor model API call metrics across your workspace. When anomalies occur, receive real-time notifications to respond quickly and maintain business continuity.
Go to [Settings > Alerts](https://home.qwencloud.com/settings/alerts) to get started.
## [​ ](#alert-rules) Alert rules

### [​ ](#create-an-alert-rule) Create an alert rule

On the [Rules](https://home.qwencloud.com/settings/alerts/rules) tab, click **Create Alert Rule**. The creation wizard has three steps:
**Step 1: Basic Info**
FieldDescription**Alert Name**A descriptive name for the rule (max 50 characters)**Alert Template**Select a pre-configured template or create a custom rule**Model Selection**Choose up to 10 models to monitor 
**Step 2: Alert Settings**
FieldDescription**Alert Level**Info or Critical**Alert Content**Notification message template. Supports variables: `{{$tags.model}}`, `{{$tags.workspace_id}}`, `{{ printf "%.2f" $value }}`**Duration**Generate alert immediately, or alert based on conditions**Check Cycle**How often to check the metric (in seconds, default 60) 
**Step 3: Notification Strategy**
FieldDescription**Alert Notification**Email address(es) to receive alerts**Time Window**Active hours for notifications (default 00:00–23:59)**Escalation**No escalation (notify once) or send repeated notifications until resolved 
### [​ ](#manage-alert-rules) Manage alert rules

On the Rules tab, you can:

- **Search** rules by name or ID

- **Filter** by status (All, Active, Inactive)

- **Actions** per rule: Start, Stop, Edit, Delete


## [​ ](#alert-templates) Alert templates

On the [Templates](https://home.qwencloud.com/settings/alerts/templates) tab, create reusable templates with pre-configured rule strategies and metrics.
Templates let you quickly create multiple alert rules with consistent configurations. Each template defines a rule strategy and the metric to monitor.
To create a template, click **Create Alert Template** and configure the template name, rule strategy, and metric.
## [​ ](#alert-history) Alert history

On the [History](https://home.qwencloud.com/settings/alerts/history) tab, review past alerts and notification delivery records.
### [​ ](#alert-history-2) Alert History

View all triggered alerts with the following details:
ColumnDescription**Alert Instance**The specific alert occurrence**Alert Level**Info or Critical**Alert Time**When the alert was triggered**Alert Count**Number of times the alert fired**Alert Rule**The rule that triggered the alert**Status**Current alert status**Notification Target**Who was notified 
Use filters to narrow results by time range, alert level, alert rule, or status.
### [​ ](#notification-history) Notification History

Click **Notification History** to view detailed delivery records for each alert, including:

- Alert target (recipient)

- Alert content (the rendered notification message)

- Timestamp of each notification sent


 [Previous ](/developer-guides/integrations/mlops-observability)[Datasets overview Manage training datasets for fine-tuning models on Qwen Cloud. Next ](/developer-guides/datasets/overview)
