# Improve recognition accuracy

> **Source:** https://docs.qwencloud.com/developer-guides/speech/improve-recognition-accuracy

Improve speech recognition accuracy using custom hotwords and context enhancement.

 Copy page Qwen Cloud provides two methods to improve ASR accuracy: **custom hotwords** for term-level biasing and **context enhancement** for conversation-aware recognition.
FeatureHow it worksBest forCustom hotwordsBoost specific terms with priority weightsFixed terminology: product names, proper nouns, medical termsContext enhancementPass conversation history to the ASR modelDynamic context: names, locations, domain terms from ongoing conversations 
## [​ ](#prerequisites) Prerequisites


- [Get your API key](/api-reference/preparation/api-key) and set it as an environment variable.

- [Install the DashScope SDK](/api-reference/preparation/install-sdk).


## [​ ](#custom-hotwords) Custom hotwords

### [​ ](#supported-scope) Supported scope

Hotwords are supported by Fun-ASR models. The following models are available:

- **Real-time speech recognition**: fun-asr-realtime, fun-asr-realtime-2025-11-07

- **Non-real-time speech recognition**: fun-asr, fun-asr-2025-11-07, fun-asr-2025-08-25, fun-asr-mtl, fun-asr-mtl-2025-08-25, fun-asr-flash-2026-06-15


For the full model list, see [Speech-to-text models](/developer-guides/speech/speech-to-text-models).
### [​ ](#quick-start) Quick start

**Workflow:**

- **Create a hotword list**: Call the [Create API](/api-reference/speech-recognition/custom-hotwords/http-api) to define a list of hotwords and set `target_model` to the speech recognition model you plan to use.

- **Use the hotword list**: Pass the hotword list ID (`vocabulary_id`) in the speech recognition request parameters. Ensure that `target_model` matches the model being called.


Audio file used in the examples: [asr_example.wav](https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/en-US/20250805/hsiqyf/asr_example.wav).
- Python 
- Java 

 Copy ```\nimport dashscope
from dashscope.audio.asr import *
import os

dashscope.api_key = os.environ.get(&#x27;DASHSCOPE_API_KEY&#x27;)

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;
dashscope.base_websocket_api_url = &#x27;wss://dashscope-intl.aliyuncs.com/api-ws/v1/inference&#x27;
prefix = &#x27;testpfx&#x27;
target_model = "fun-asr-realtime"

my_vocabulary = [
 {"text": "Speech Laboratory", "weight": 4}
]

service = VocabularyService()
vocabulary_id = service.create_vocabulary(
 prefix=prefix,
 target_model=target_model,
 vocabulary=my_vocabulary)

try:
 if service.query_vocabulary(vocabulary_id)[&#x27;status&#x27;] == &#x27;OK&#x27;:
 recognition = Recognition(model=target_model,
 format=&#x27;wav&#x27;,
 sample_rate=16000,
 callback=None,
 vocabulary_id=vocabulary_id)
 result = recognition.call(&#x27;asr_example.wav&#x27;)
 print(result.output)
finally:
 service.delete_vocabulary(vocabulary_id)

``` Copy ```\nimport com.alibaba.dashscope.audio.asr.recognition.Recognition;
import com.alibaba.dashscope.audio.asr.recognition.RecognitionParam;
import com.alibaba.dashscope.audio.asr.vocabulary.Vocabulary;
import com.alibaba.dashscope.audio.asr.vocabulary.VocabularyService;
import com.alibaba.dashscope.exception.InputRequiredException;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.alibaba.dashscope.utils.Constants;
import com.google.gson.JsonArray;
import com.google.gson.JsonObject;

import java.io.File;
import java.util.ArrayList;
import java.util.List;

public class Main {
 public static String apiKey = System.getenv("DASHSCOPE_API_KEY");

 public static void main(String[] args) throws NoApiKeyException, InputRequiredException {
 Constants.baseHttpApiUrl = "https://dashscope-intl.aliyuncs.com/api/v1";
 Constants.baseWebsocketApiUrl = "wss://dashscope-intl.aliyuncs.com/api-ws/v1/inference";

 String targetModel = "fun-asr-realtime";

 JsonArray vocabularyJson = new JsonArray();
 List<Hotword> wordList = new ArrayList<>();
 wordList.add(new Hotword("Speech Laboratory", 4));

 for (Hotword word : wordList) {
 JsonObject jsonObject = new JsonObject();
 jsonObject.addProperty("text", word.text);
 jsonObject.addProperty("weight", word.weight);
 vocabularyJson.add(jsonObject);
 }

 VocabularyService service = new VocabularyService(apiKey);
 Vocabulary vocabulary = service.createVocabulary(targetModel, "testpfx", vocabularyJson);

 try {
 if ("OK".equals(service.queryVocabulary(vocabulary.getVocabularyId()).getStatus())) {
 Recognition recognizer = new Recognition();
 RecognitionParam param =
 RecognitionParam.builder()
 .model(targetModel)
 .apiKey(apiKey)
 .format("wav")
 .sampleRate(16000)
 .vocabularyId(vocabulary.getVocabularyId())
 .build();

 try {
 System.out.println("Recognition result: " + recognizer.call(param, new File("asr_example.wav")));
 } catch (Exception e) {
 e.printStackTrace();
 } finally {
 recognizer.getDuplexApi().close(1000, "bye");
 }
 }
 } finally {
 service.deleteVocabulary(vocabulary.getVocabularyId());
 }
 System.exit(0);
 }
}

class Hotword {
 String text;
 int weight;

 public Hotword(String text, int weight) {
 this.text = text;
 this.weight = weight;
 }
}

``` 
### [​ ](#hotword-format) Hotword format

Submit a JSON array of hotword objects.
**Example**: Improve movie title recognition (Fun-ASR and Paraformer series models)
Copy ```\n[
 {"text": "赛德克巴莱", "weight": 4, "lang": "zh"},
 {"text": "Seediq Bale", "weight": 4, "lang": "en"},
 {"text": "夏洛特烦恼", "weight": 4, "lang": "zh"},
 {"text": "Goodbye Mr. Loser", "weight": 4, "lang": "en"},
 {"text": "阙里人家", "weight": 4, "lang": "zh"},
 {"text": "Confucius&#x27; Family", "weight": 4, "lang": "en"}
]

``` 
**Field descriptions**:
FieldTypeRequiredDescriptiontextstringYesThe hotword text. Must be supported by the selected model. Use actual words, not random characters. See length rules below.weightintYesPriority weight, an integer from 1 to 5. Start with 4. Increase if results are weak, but too high a weight can hurt recognition of other words.langstringNoLanguage code. Boosts hotwords for a specific language. Leave empty for auto-detection. See the model&#x27;s API reference for supported codes. If you set `language_hints`, only matching hotwords take effect. 
**Hotword text length rules**:

- 
**Contains non-ASCII characters**: Maximum 15 characters total, including non-ASCII characters (Chinese, Japanese kana, Korean Hangul, Russian Cyrillic) and ASCII characters.
Examples:

`"厄洛替尼盐酸盐"` (7 Chinese characters)

- `"EGFR抑制剂"` (3 Chinese characters and 4 ASCII characters, for a total of 7 characters)

- `"こんにちは"` (5 characters)

- `"Фенибут Белфарм"` (15 characters, including the space)

- `"Клофелин Белмедпрепараты"` (24 characters) -- exceeds limit


- 
**Contains only ASCII characters**: Maximum 7 segments. A segment is a sequence of characters separated by spaces.
Examples:

`"Exothermic reaction"` -- 2 segments

- `"Human immunodeficiency virus type 1"` -- 5 segments

- `"The effect of temperature variations on enzyme activity in biochemical reactions"` -- 11 segments, exceeds limit


### [​ ](#tune-hotword-performance) Tune hotword performance

#### [​ ](#adjust-hotword-weights) Adjust hotword weights

Weight controls how strongly the model favors a hotword. Set it appropriately to improve target word accuracy without introducing false recognitions.
WeightEffectBest for1-2Slight preferenceHotwords that sound similar to common words, where overcorrection must be avoided3-4Clear preference (recommended)The best starting point for most scenarios5Forced preferenceUse only when the term appears frequently in the audio and is unlikely to be confused with other words. An excessively high weight can cause phonetically similar words to be misrecognized as the hotword. 
Start with `weight=4` and adjust incrementally based on recognition results.
#### [​ ](#design-hotword-lists) Design hotword lists


- **Group by scenario**: Create separate vocabulary lists for different business scenarios (for example, one for medical terms and another for product names) to simplify maintenance and reuse.

- **Mix multiple languages**: A single vocabulary list can contain terms in different languages. Use the `lang` field to distinguish them. When `language_hints` is specified during speech recognition, only hotwords that match the specified language take effect.

- **Clean up regularly**: Delete unused vocabulary lists to free up quota. Each account supports up to 10 lists.


### [​ ](#limits-and-billing) Limits and billing

LimitDescriptionNumber of vocabulary lists10 per account, shared across all models.Hotwords per listUp to 500 hotwords per vocabulary list.BillingFree of charge. 
## [​ ](#context-enhancement) Context enhancement

### [​ ](#supported-scope-2) Supported scope

Context enhancement is supported by:

- **Non-real-time speech recognition**: fun-asr-flash-2026-06-15


### [​ ](#quick-start-2) Quick start

Context enhancement requires no pre-created resources. Simply pass context parameters in your speech recognition request:

- **Non-real-time speech recognition:** Pass context messages in `input.messages` of the HTTP request, placed before the audio message.


**Use cases:** By passing conversation history or domain terminology as context, you can significantly improve transcription accuracy for specialized terms such as names, locations, and product terminology. This feature supports the following scenarios:

- **Word list enhancement:** Pass domain-specific word lists or terminology through `user` (`input_text`) to help the model accurately recognize specialized vocabulary.

- **Multi-turn conversation context:** In voice interaction scenarios that combine ASR with large language models, pass prior recognition results (`user` / `input_text`) and model responses (`assistant` / `text`) to improve recognition accuracy for the current turn.


 
- **Message count limit:** The engine retains up to the 5 most recent turns of context. Word list enhancement typically requires only 1 message and is not affected by this limit. When exceeded, earlier messages are automatically ignored without returning an error.

- **Text length limit:** The total text length per turn of context (the combined length of all `user` and `assistant` `text` fields in the same turn) must not exceed 400 characters (counted per character, where each character counts as 1, including letters, Chinese characters, digits, spaces, and punctuation). Excess content is silently truncated from the end without returning an error. In multi-turn contexts, each turn is counted independently.

- **Context mechanism:** Context works primarily through word matching. The `text` field should contain the exact words to be recognized in the audio (such as "Kubernetes" or "Bulge Bracket"). Passing only semantically related descriptions without the target words will have limited effect.


 
## [​ ](#non-real-time-speech-recognition) Non-real-time speech recognition

Pass context through `input.messages`. Use the `user` role with `input_text` type for prior speech recognition results or domain-specific word lists. Use the `assistant` role for prior model responses (optional). Place context messages before the audio message. For details, see [DashScope (Fun-ASR)](/api-reference/speech-recognition/fun-asr-recording/restful-api).
- Word list enhancement 
- Multi-turn conversation context 

 Pass domain-specific terminology through `user` (`input_text`). No `assistant` message is needed.Copy ```\n{
 "model": "fun-asr-flash-2026-06-15",
 "input": {
 "messages": [
 {
 "role": "user",
 "content": [
 {
 "type": "input_text",
 "text": "Kubernetes Istio Envoy service mesh sidecar proxy"
 }
 ]
 },
 {
 "role": "user",
 "content": [
 {
 "type": "input_audio",
 "input_audio": {
 "data": "Audio URL or Base64 of the current audio to recognize"
 }
 }
 ]
 }
 ]
 },
 "parameters": {}
}

``` Pass prior recognition results (`user` / `input_text`) and model responses (`assistant` / `text`).Copy ```\n{
 "model": "fun-asr-flash-2026-06-15",
 "input": {
 "messages": [
 {
 "role": "user",
 "content": [
 {
 "type": "input_text",
 "text": "Prior user speech recognition result"
 }
 ]
 },
 {
 "role": "assistant",
 "content": [
 {
 "type": "text",
 "text": "Prior model response content"
 }
 ]
 },
 {
 "role": "user",
 "content": [
 {
 "type": "input_audio",
 "input_audio": {
 "data": "Audio URL or Base64 of the current audio to recognize"
 }
 }
 ]
 }
 ]
 },
 "parameters": {}
}

``` 
### [​ ](#effect-example) Effect example

The `text` field content format is flexible -- it can be a word list, natural language paragraph, or a mix of both. It has high tolerance for unrelated text.
An audio clip should be correctly recognized as: "The jargon within investment banking circles, how much do you know? First, the nine major foreign investment banks, Bulge Bracket, BB ..."
Without context enhancementWith context enhancementWithout context enhancement, some investment bank names are recognized incorrectly. For example, "Bird Rock" should be "Bulge Bracket". Recognition result: "...the nine major foreign investment banks, Bird Rock, BB ..."With context enhancement, investment bank names are recognized correctly. Recognition result: "...the nine major foreign investment banks, Bulge Bracket, BB ..." 
In the example above, adding a word list or natural language paragraph containing terms like "Bulge Bracket" to the `text` field achieves the enhancement effect.
## [​ ](#api-reference) API reference


- [Custom Hotword API Reference](/api-reference/speech-recognition/custom-hotwords/http-api)


## [​ ](#faq) FAQ

### [​ ](#why-dont-hotwords-improve-recognition-accuracy) Why don&#x27;t hotwords improve recognition accuracy?

Check the following in order:

- **Model mismatch**: The `target_model` specified when creating the list must match the model used by the speech recognition API. A mismatch doesn&#x27;t cause an error, and recognition still returns results, but the hotwords don&#x27;t take effect. If the results don&#x27;t contain expected hotwords, check this first.

- **Unsupported model**: The model must belong to the Fun-ASR or Paraformer family. Other families don&#x27;t support hotwords. Calling the API with an unsupported model doesn&#x27;t return an error, but the results may be empty or lack hotword enhancement. If using a model such as SenseVoice, check this first.

- **Inappropriate weight**: Increase the weight from 4 to 5 and observe the results. If phonetically similar words start being misrecognized as the hotword, reduce it back to 4.

- **Hotword list status**: Use the Query API to confirm that `status` is `OK`.


### [​ ](#are-hotwords-used-differently-in-real-time-and-file-based-recognition) Are hotwords used differently in real-time and file-based recognition?

Hotword lists are created the same way. The calling method differs:

- **Real-time speech recognition**: Pass `vocabulary_id` in the Recognition or WebSocket connection parameters.

- **File-based speech recognition**: Pass `vocabulary_id` in the Transcription request parameters.


In both cases, `target_model` must match the speech recognition model used in the API call.
### [​ ](#how-to-improve-recognition-accuracy-beyond-hotwords) How to improve recognition accuracy beyond hotwords?

In addition to hotwords and context enhancement, consider the following:

- **Audio quality**: Match the sample rate to the model requirements (16 kHz or 8 kHz) and reduce background noise.

- **Choose the right model**: Different scenarios call for different models. For details, see the [Speech-to-text](/developer-guides/speech/speech-to-text-models) model selection guide.

- **Specify the language**: Declare the audio language through `language_hints` to improve accuracy in single-language scenarios.


 [Previous ](/developer-guides/speech/asr)[Text-to-speech models Choose a model for speech synthesis, voice cloning, and voice design. Next ](/developer-guides/speech/tts-models)
