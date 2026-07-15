# MLOps &amp; observability

> **Source:** https://docs.qwencloud.com/developer-guides/integrations/mlops-observability

Production AI monitoring

 Copy page ## [​ ](#overview) Overview

The [Analytics](https://home.qwencloud.com/analytics) page provides observability for your model deployments, including token consumption, request counts, latency, success rates, and per-model performance metrics.
## [​ ](#analytics) Analytics

Go to the [Analytics](https://home.qwencloud.com/analytics) page to view usage and analytics for your workspace.
### [​ ](#filters) Filters


- **Time range**: Select the time window (such as 7 days).

- **Models**: Filter by specific model or view all models.

- **API Keys**: Filter usage by a specific API key.

- **Granularity**: Choose the aggregation interval (By Day or By Hour).


### [​ ](#tabs) Tabs


- **Usage**: Displays usage trend charts and per-model metrics.

- **Logs**: View detailed request logs with filtering by source, status, and more.


### [​ ](#metrics) Metrics

The page shows the following usage metrics:
MetricDescription**Tokens**Total token consumption**Count**Total number of API requests**Seconds**Usage for models billed by duration (audio/video)**Characters**Usage for models billed by character count**API Calls**Total API call count 
Below the usage metrics, trend charts show requests, average TTFT (time to first token), average latency, and success rate over time.
 Cost includes all consumption across the entire platform. Refer to billing data for details. 
### [​ ](#usage-units-by-model-type) Usage units by model type

TypeSubcategoryUnitBilling basis**Large language model**Text generation, Deep thinking, Vision understandingTokenBilled by input and output token count**Vision model**Image generationImage (count)Billed by successfully generated images**Vision model**Video generationSecondsBilled by successfully generated video duration**Speech model**TTS, Realtime TTS, File ASR, Realtime ASR, Audio/video translationSeconds, characters, or tokensVaries by model -- may bill by audio duration, text characters, or token count**Omni-modal model**Omni-modal, Realtime multimodalTokenText billed by tokens; other modalities (audio, image, video) billed by corresponding token count 
### [​ ](#per-model-metrics) Per-model metrics

Below the usage charts, a per-model table breaks down throughput (TPM/RPM), call volume, success rate, time to first token, and latency for each model. Use this to identify underperforming models or unexpected error spikes.
### [​ ](#request-logs) Request logs

Go to the [Logs](https://home.qwencloud.com/analytics/request-logs) tab to inspect individual API requests. Logs are retained for the last **14 days**.
#### [​ ](#filters-2) Filters


- **Time range**: Narrow the log window to a specific period.

- **Model**: Filter by a specific model.

- **API Keys**: Filter by a specific API key.

- **Sources**: Distinguish between API calls and web-based calls (such as the Try AI page) to pinpoint traffic sources.

- **Status**: Filter by HTTP status code (such as 200, 400, 429).

- **Request ID**: Search for a specific request by its ID.


#### [​ ](#log-table) Log table

Each row shows:
ColumnDescription**Request ID**Unique identifier for the request**Timestamp**When the request was made**Model**Model used for the request**Source**API call or web-based call**API Key**The API key used for the request**Usage**Token breakdown (total, input, output)**TTFT**Time to first token**Latency**Total response time**Status**HTTP status code 
#### [​ ](#request-details) Request details

Click **Details** to open a side panel with the full request breakdown. You can view the data in **List** or **JSON** format.
Click **Export** to download logs as a file for offline analysis.

---

## [​ ](#alerts) Alerts

Set up custom alert rules to monitor API call metrics and receive real-time notifications when anomalies occur. For details, see [Monitoring alerts](/developer-guides/integrations/alerts). [Previous ](/developer-guides/clients-and-developer-tools/more-tools)[Monitoring alerts Configure alert rules for model API metrics Next ](/developer-guides/integrations/alerts)
