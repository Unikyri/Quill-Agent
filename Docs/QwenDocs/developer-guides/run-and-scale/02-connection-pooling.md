# Connection reuse and pooling

> **Source:** https://docs.qwencloud.com/developer-guides/run-and-scale/connection-pooling

HTTP connection reuse and WebSocket connection pooling for high-concurrency workloads.

 Copy page Reusing connections reduces resource consumption and improves throughput. The strategy depends on the protocol:

- **HTTP APIs** (text generation, multimodal, embeddings): reuse TCP connections via connection pool configuration (Java) or Session objects (Python).

- **WebSocket APIs** (TTS, real-time speech): pool synthesizer objects that hold long-lived WebSocket connections.


## [​ ](#prerequisites) Prerequisites


- Obtain and configure your API Key as the `DASHSCOPE_API_KEY` environment variable.

- Install the latest DashScope SDK:

**Python SDK**: >= 1.25.2

- **Java SDK**: >= 2.16.6


## [​ ](#http-connection-reuse) HTTP connection reuse

 The DashScope endpoint differs by model type:
- **Text models** (`qwen-plus`, `qwen3-max`, etc.): use the `Generation` class, which routes to `/services/aigc/text-generation/generation`.

- **Multimodal models** (`qwen3.7-plus`, `qwen3-vl-plus`, etc.): use the `MultiModalConversation` class, which routes to `/services/aigc/multimodal-generation/generation`.


 
### [​ ](#java-sdk) Java SDK

Connection pooling is enabled by default. Adjust the following parameters as needed.
ParameterDescriptionDefaultUnitNotes`connectTimeout`Timeout for establishing a connection.120secondsShorter timeouts reduce wait time in low-latency scenarios.`readTimeout`Timeout for reading data.300seconds`writeTimeout`Timeout for writing data.60seconds`connectionIdleTimeout`Timeout for idle connections.300secondsLonger idle timeouts avoid frequent reconnections under high concurrency.`connectionPoolSize`Maximum connections in the pool.32itemsToo few connections cause blocking; too many increase server load.`maximumAsyncRequests`Maximum concurrent requests across all hosts. Must be ≤ `connectionPoolSize`.32requests`maximumAsyncRequestsPerHost`Maximum concurrent requests per host. Must be ≤ `maximumAsyncRequests`.32items 
Configure connection pool parameters and call a model service:
Copy ```\n// Recommended DashScope SDK version >= 2.12.0
import java.time.Duration;
import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;

import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversation;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationParam;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationResult;
import com.alibaba.dashscope.common.MultiModalMessage;
import com.alibaba.dashscope.common.Role;
import com.alibaba.dashscope.exception.ApiException;
import com.alibaba.dashscope.exception.InputRequiredException;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.alibaba.dashscope.protocol.ConnectionConfigurations;
import com.alibaba.dashscope.protocol.Protocol;
import com.alibaba.dashscope.utils.Constants;

public class Main {
 public static MultiModalConversationResult callWithMessage() throws ApiException, NoApiKeyException, InputRequiredException {
 MultiModalConversation conv = new MultiModalConversation(Protocol.HTTP.getValue(), "https://dashscope-intl.aliyuncs.com/api/v1");
 Map<String, Object> textContent = new HashMap<>();
 textContent.put("text", "Who are you?");
 MultiModalMessage userMsg = MultiModalMessage.builder()
 .role(Role.USER.getValue())
 .content(Collections.singletonList(textContent))
 .build();
 MultiModalConversationParam param = MultiModalConversationParam.builder()
 // If you have not configured the environment variable, replace with your API key: .apiKey("sk-xxx")
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model("qwen3.7-plus")
 .messages(Collections.singletonList(userMsg))
 .build();

 return conv.call(param);
 }
 public static void main(String[] args) {
 // Connection pool configuration
 Constants.connectionConfigurations = ConnectionConfigurations.builder()
 .connectTimeout(Duration.ofSeconds(10)) // Timeout for establishing a connection, default 120s
 .readTimeout(Duration.ofSeconds(300)) // Timeout for reading data, default 300s
 .writeTimeout(Duration.ofSeconds(60)) // Timeout for writing data, default 60s
 .connectionIdleTimeout(Duration.ofSeconds(300)) // Timeout for idle connections, default 300s
 .connectionPoolSize(256) // Maximum connections in the connection pool, default 32
 .maximumAsyncRequests(256) // Maximum concurrent requests, default 32
 .maximumAsyncRequestsPerHost(256) // Maximum concurrent requests per host, default 32
 .build();

 try {
 MultiModalConversationResult result = callWithMessage();
 System.out.println(result.getOutput().getChoices().get(0).getMessage().getContent().get(0).get("text"));
 } catch (ApiException | NoApiKeyException | InputRequiredException e) {
 System.err.println("An error occurred while calling the service: " + e.getMessage());
 }
 System.exit(0);
 }
}

``` 
### [​ ](#python-sdk) Python SDK

The Python SDK supports connection reuse via a custom Session. Two methods are available: [async (aiohttp)](#async-aiohttp) and [sync (requests.Session)](#sync-requests-session).
#### [​ ](#async-aiohttp) Async (aiohttp)

Use `aiohttp.ClientSession` with `aiohttp.TCPConnector` for async connection reuse.
ParameterDescriptionDefaultNotes`limit`Total connection limit100Higher values improve concurrency.`limit_per_host`Connection limit per host0 (unlimited)Prevents excessive load on a single host.`ssl`SSL context configurationNoneSSL certificate validation for HTTPS connections. 
Copy ```\nimport asyncio
import aiohttp
import ssl
import certifi
from dashscope import AioMultiModalConversation
import dashscope
import os

async def main():
 dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

 # If you have not configured the environment variable, replace with your API key: dashscope.api_key = "sk-xxx"
 dashscope.api_key = os.getenv("DASHSCOPE_API_KEY")

 # Configure connection parameters
 connector = aiohttp.TCPConnector(
 limit=100, # Total connection limit
 limit_per_host=30, # Connection limit per host
 ssl=ssl.create_default_context(cafile=certifi.where()),
 )

 # Create a custom Session and pass it to the call method
 async with aiohttp.ClientSession(connector=connector) as session:
 response = await AioMultiModalConversation.call(
 model=&#x27;qwen3.7-plus&#x27;,
 messages=[{&#x27;role&#x27;: &#x27;user&#x27;, &#x27;content&#x27;: [{&#x27;text&#x27;: &#x27;Hello, please introduce yourself&#x27;}]}],
 session=session, # Pass the custom Session
 )
 print(response)

asyncio.run(main())

``` 
#### [​ ](#sync-requests-session) Sync (requests.Session)

Use `requests.Session` for sync connection reuse. Requests within the same Session reuse the TCP connection.
Copy ```\nimport requests
from dashscope import MultiModalConversation
import dashscope
import os

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

# If you have not configured the environment variable, replace with your API key: dashscope.api_key = "sk-xxx"
dashscope.api_key = os.getenv("DASHSCOPE_API_KEY")

# Use a with statement to ensure the Session closes correctly
with requests.Session() as session:
 response = MultiModalConversation.call(
 model=&#x27;qwen3.7-plus&#x27;,
 messages=[{&#x27;role&#x27;: &#x27;user&#x27;, &#x27;content&#x27;: [{&#x27;text&#x27;: &#x27;Hello&#x27;}]}],
 session=session # Pass the custom Session
 )
 print(response)

``` 
Reuse the same Session across multiple calls:
Copy ```\nimport requests
from dashscope import MultiModalConversation
import dashscope
import os

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

# If you have not configured the environment variable, replace with your API key: dashscope.api_key = "sk-xxx"
dashscope.api_key = os.getenv("DASHSCOPE_API_KEY")

# Create a Session object
session = requests.Session()

try:
 # Reuse the same Session for multiple calls
 response1 = MultiModalConversation.call(
 model=&#x27;qwen3.7-plus&#x27;,
 messages=[{&#x27;role&#x27;: &#x27;user&#x27;, &#x27;content&#x27;: [{&#x27;text&#x27;: &#x27;Hello&#x27;}]}],
 session=session
 )
 print(response1)

 response2 = MultiModalConversation.call(
 model=&#x27;qwen3.7-plus&#x27;,
 messages=[{&#x27;role&#x27;: &#x27;user&#x27;, &#x27;content&#x27;: [{&#x27;text&#x27;: &#x27;Introduce yourself&#x27;}]}],
 session=session
 )
 print(response2)
finally:
 # Ensure the Session closes correctly
 session.close()

``` 
## [​ ](#websocket-connection-pooling) WebSocket connection pooling

TTS services use WebSocket connections for real-time streaming. In production, creating a new connection per request wastes resources and adds latency. This section covers connection pooling, object pooling, and concurrent request management for high-throughput TTS workloads.
### [​ ](#python-object-pool) Python: Object pool

The Python SDK provides `SpeechSynthesizerObjectPool` to manage and reuse `SpeechSynthesizer` instances. The pool pre-creates objects and establishes WebSocket connections at initialization, eliminating per-request connection overhead.
**Pool sizing**: Set `max_size` to 1.5x-2x your peak concurrency. Do not exceed your account&#x27;s QPS limit.
Copy ```\nimport os
import threading
import dashscope
from dashscope.audio.tts_v2 import *

dashscope.api_key = os.getenv("DASHSCOPE_API_KEY")
dashscope.base_websocket_api_url = &#x27;wss://dashscope-intl.aliyuncs.com/api-ws/v1/inference&#x27;

# Create a global object pool (one-time cost at startup)
pool = SpeechSynthesizerObjectPool(max_size=20)

def synthesize(text, task_id):
 complete_event = threading.Event()

 class Callback(ResultCallback):
 def on_open(self):
 self.file = open(f&#x27;result_{task_id}.mp3&#x27;, &#x27;wb&#x27;)

 def on_complete(self):
 complete_event.set()

 def on_error(self, message):
 print(f&#x27;[task_{task_id}] Error: {message}&#x27;)

 def on_data(self, data):
 self.file.write(data)

 def on_close(self):
 if hasattr(self, &#x27;file&#x27;):
 self.file.close()

 callback = Callback()

 # Borrow a pre-connected synthesizer from the pool
 synth = pool.borrow_synthesizer(
 model=&#x27;cosyvoice-v3-flash&#x27;,
 voice=&#x27;longanyang&#x27;,
 callback=callback
 )

 try:
 synth.call(text)
 complete_event.wait()
 print(f&#x27;[task_{task_id}] First packet delay: &#x27;
 f&#x27;{synth.get_first_package_delay()} ms&#x27;)
 # Return the synthesizer to the pool for reuse
 pool.return_synthesizer(synth)
 except Exception as e:
 print(f&#x27;[task_{task_id}] Failed: {e}&#x27;)
 synth.close() # Do not return failed objects

# Run concurrent tasks
texts = ["First sentence.", "Second sentence.", "Third sentence."]
threads = [threading.Thread(target=synthesize, args=(t, i))
 for i, t in enumerate(texts)]
for t in threads:
 t.start()
for t in threads:
 t.join()

pool.shutdown()

``` 
 Never return a synthesizer to the pool if the task failed or is still running. Close it manually instead. 
### [​ ](#java-connection-pool-object-pool) Java: Connection pool + object pool

The Java SDK uses OkHttp3 connection pooling (enabled by default) plus an optional Apache Commons Pool2 object pool for `SpeechSynthesizer` instances.
**Step 1: Configure connection pool via environment variables**
VariableDefaultRecommendation`DASHSCOPE_CONNECTION_POOL_SIZE`322x peak concurrency`DASHSCOPE_MAXIMUM_ASYNC_REQUESTS`32Match connection pool size`DASHSCOPE_MAXIMUM_ASYNC_REQUESTS_PER_HOST`32Match connection pool size 
Copy ```\nexport DASHSCOPE_CONNECTION_POOL_SIZE=2000
export DASHSCOPE_MAXIMUM_ASYNC_REQUESTS=2000
export DASHSCOPE_MAXIMUM_ASYNC_REQUESTS_PER_HOST=2000

``` 
**Step 2: Add commons-pool2 dependency**
- Maven 
- Gradle 

 Copy ```\n<dependency>
 <groupId>org.apache.commons</groupId>
 <artifactId>commons-pool2</artifactId>
 <version>the-latest-version</version>
</dependency>

``` Copy ```\nimplementation group: &#x27;org.apache.commons&#x27;,
 name: &#x27;commons-pool2&#x27;, version: &#x27;the-latest-version&#x27;

``` 
**Step 3: Create and use the object pool**
VariableDefaultRecommendation`SAMBERT_OBJECTPOOL_SIZE` (Sambert)5001.5x-2x peak concurrency, must not exceed connection pool size 
Copy ```\nimport com.alibaba.dashscope.audio.tts.SpeechSynthesisParam;
import com.alibaba.dashscope.audio.tts.SpeechSynthesizer;
import org.apache.commons.pool2.BasePooledObjectFactory;
import org.apache.commons.pool2.PooledObject;
import org.apache.commons.pool2.impl.DefaultPooledObject;
import org.apache.commons.pool2.impl.GenericObjectPool;
import org.apache.commons.pool2.impl.GenericObjectPoolConfig;

// Factory
class SynthesizerFactory extends BasePooledObjectFactory<SpeechSynthesizer> {
 public SpeechSynthesizer create() { return new SpeechSynthesizer(); }
 public PooledObject<SpeechSynthesizer> wrap(SpeechSynthesizer obj) {
 return new DefaultPooledObject<>(obj);
 }
}

// Pool (global singleton)
GenericObjectPoolConfig<SpeechSynthesizer> config = new GenericObjectPoolConfig<>();
config.setMaxTotal(1200);
config.setMaxIdle(1200);
config.setMinIdle(1200);
GenericObjectPool<SpeechSynthesizer> pool =
 new GenericObjectPool<>(new SynthesizerFactory(), config);

// Usage in each task
SpeechSynthesizer synth = pool.borrowObject();
try {
 // ... configure params and call synth
 pool.returnObject(synth);
} catch (Exception e) {
 synth = null; // Do not return on failure
}

``` 
 Reference server sizing: a 4-core 8 GiB machine can handle ~600 concurrent Sambert TTS tasks with an object pool of 1200 and a connection pool of 2000. 
## [​ ](#best-practices) Best practices


- **Java SDK**: Set `connectionPoolSize` and `maximumAsyncRequests` based on your concurrent workload. Too few connections cause blocking; too many increase server load.

- **Python SDK**: Use `with` statements to manage the Session lifecycle and ensure proper resource cleanup.

- **Choose the right method**: Use async calls for async applications (like asyncio or FastAPI). Use sync calls for traditional applications.

- **WebSocket object pools**: Never return a synthesizer to the pool if the task failed or is still running. Close it manually instead.


## [​ ](#performance-monitoring) Performance monitoring

Track these metrics to maintain healthy production TTS services:
MetricDescriptionTargetFirst packet delayTime from request to first audio chunk< 500 msEnd-to-end latencyTotal time for complete synthesisDepends on text lengthError ratePercentage of failed requests< 0.1%Pool utilizationBorrowed objects / pool size60%-80% at peakConnection reuse ratioReused connections / total requests> 95% 
Access these metrics from the SDK:
Copy ```\n# TTS
print(f"Request ID: {synthesizer.get_last_request_id()}")
print(f"First packet delay: {synthesizer.get_first_package_delay()} ms")

``` 
Copy ```\n// TTS
System.out.println("Request ID: " + synthesizer.getLastRequestId());
System.out.println("First packet delay: " + synthesizer.getFirstPackageDelay() + " ms");

``` 
## [​ ](#production-checklist) Production checklist

Before going live, verify the following:

- API Key stored in environment variable, not hardcoded.

- Connection pool and object pool sizes configured for expected peak load.

- Pool sizes do not exceed your account&#x27;s QPS limit.

- Error handling returns failed objects to disposal (not back to pool).

- Graceful shutdown calls `pool.shutdown()` (Python) or pool close (Java).

- WebSocket connections use the correct endpoint (`wss://dashscope-intl.aliyuncs.com/...`).

- Monitoring dashboards track first-packet delay, error rate, and pool utilization.

- Load tested with 2x expected peak concurrency.

- Retry logic with exponential backoff for transient failures.


## [​ ](#related) Related


- [Text to Speech](/developer-guides/speech/tts) -- TTS models, parameters, and streaming modes.

- [Realtime streaming](/developer-guides/speech/realtime-streaming) -- realtime TTS streaming guide.

- [Improve recognition accuracy](/developer-guides/accuracy-tuning/speech-recognition) -- ASR optimization including high-concurrency ASR patterns.


 [Previous ](/developer-guides/run-and-scale/async-task-management)[Accuracy tuning Maximize correctness and consistent behavior across text, image, video, speech, and vision models on Qwen Cloud Next ](/developer-guides/accuracy-tuning/overview)
