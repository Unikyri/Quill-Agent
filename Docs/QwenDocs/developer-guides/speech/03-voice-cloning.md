# Voice cloning

> **Source:** https://docs.qwencloud.com/developer-guides/speech/voice-cloning

Clone a voice from audio samples for use with CosyVoice, Qwen-TTS, or Qwen-Omni models.

 Copy page Clone a voice from 10-20 seconds of audio. The API returns a voice identifier instantly -- no training required.
## [​ ](#how-it-works) How it works


- **Clone a voice** -- Call the voice cloning API with an audio sample and a `target_model`. The API returns a `voice` identifier instantly.

- **Use the cloned voice** -- Pass the `voice` identifier to the target model&#x27;s API. The model must match the `target_model` from step 1.


 The `target_model` set during voice creation must match the model used in subsequent API calls. Mismatched models cause synthesis to fail. 
## [​ ](#supported-models) Supported models

**Voice cloning model**:

- **Qwen-TTS / Qwen-Omni**: `qwen-voice-enrollment`

- **CosyVoice**: `voice-enrollment`


**Target models** (`target_model`):

- **CosyVoice** -- `cosyvoice-v3-plus`, `cosyvoice-v3-flash`. Use cloned voices with [CosyVoice TTS](/developer-guides/speech/realtime-streaming).

- **Qwen-Omni-Realtime** -- `qwen3.5-omni-plus-realtime`, `qwen3.5-omni-flash-realtime`. Use cloned voices with the [Realtime Multimodal API](/developer-guides/speech/realtime-multimodal-speech).

- **Qwen-Omni** -- `qwen3.5-omni-plus`, `qwen3.5-omni-flash`. Use cloned voices with the [Omni API (non-realtime)](/developer-guides/speech/multimodal-speech).

- **Qwen-TTS** -- **Qwen3-TTS-VC-Realtime**, **Qwen3-TTS-VC**. Use cloned voices with [Qwen TTS](/api-reference/speech-synthesis/qwen-tts) or [Realtime streaming TTS](/developer-guides/speech/realtime-streaming). For model IDs and snapshot versions, see [Text-to-speech models](/developer-guides/speech/tts-models).


 Qwen-TTS-VC models only support custom cloned voices, not system voices like Chelsie, Serena, Ethan, or Cherry. 
## [​ ](#audio-requirements) Audio requirements

The quality of the input audio directly affects the cloning result. Each model family has different audio requirements.
- CosyVoice 
- Qwen-TTS / Qwen-Omni 

 ItemRequirement**Format**WAV (16-bit), MP3, M4A**Duration**10 -- 20 seconds recommended. 60 seconds maximum.**File size**10 MB or less**Sample rate**16 kHz or higher**Channels**Mono or stereo. For stereo audio, only the first channel is processed. Make sure the first channel contains valid speech.**Content**At least 5 seconds of continuous, clear speech. Brief pauses must not exceed 2 seconds. No background music, ambient noise, or other voices. Use normal-speed spoken audio; do not use singing.**Language**Varies by `target_model`. **cosyvoice-v3-flash**: Chinese (Mandarin, Cantonese, and regional dialects), English, French, German, Japanese, Korean, Russian, Portuguese, Thai, Indonesian, Vietnamese. **cosyvoice-v3-plus**: Chinese (Mandarin), English, French, German, Japanese, Korean, Russian. ItemRequirement**Format**WAV (16-bit), MP3, M4A**Duration**10 -- 20 seconds recommended. 60 seconds maximum.**File size**Less than 10 MB**Sample rate**24 kHz or higher**Channels**Mono**Content**At least 3 seconds of continuous, clear speech. Short pauses (up to 2 seconds) are acceptable. No background music, ambient noise, or overlapping voices. Do not use singing or song audio.**Language**Chinese (zh), English (en), German (de), Italian (it), Portuguese (pt), Spanish (es), Japanese (ja), Korean (ko), French (fr), Russian (ru) 
### [​ ](#recording-tips) Recording tips

#### [​ ](#quick-start-checklist) Quick-start checklist

Use this checklist in a standard bedroom or similar small room:

- Close all windows and doors to block external noise.

- Turn off air conditioners, fans, and other electrical devices.

- Draw curtains to reduce glass reflections.

- Cover your desk with clothing or a blanket to reduce surface reflections.

- Read through your script. Define your character&#x27;s tone and practice delivering naturally.

- Position the recording device approximately 10 cm from your mouth. Too close causes plosive distortion; too far produces a weak signal.

- Start recording.


#### [​ ](#recording-devices) Recording devices

Use a smartphone, digital voice recorder, or professional audio recorder.
#### [​ ](#set-up-your-recording-environment) Set up your recording environment

**Choose the right room:**
RequirementDetailsRoom sizeRecord in a small enclosed space (max 10 m²).Acoustic treatmentChoose a room with sound-absorbing materials: acoustic foam, carpets, or curtains.Spaces to avoidAvoid auditoriums, conference rooms, and classrooms — these large spaces cause strong reverberation that degrades clone quality. 
**Control noise:**
Noise sourceMitigationOutdoor noiseClose all windows and doors. Avoid recording near traffic or construction.Indoor noiseTurn off air conditioners, fans, and fluorescent lamp ballasts before recording. 
 Record a few seconds of ambient sound on your smartphone, then play it back at high volume to identify hidden noise sources. 
**Reduce reverberation:**
Reverberation blurs speech and reduces definition, directly impacting clone fidelity.

- Draw curtains, open closet doors, or cover desks/cabinets with clothing or bed sheets to reduce reflections from smooth surfaces.

- Place irregular objects (bookshelves, upholstered furniture) to scatter sound waves.


#### [​ ](#prepare-your-script) Prepare your script

GuidelineDetailsContentNo strict restrictions apply. Align content with your target use case.Sentence structureUse complete sentences. Avoid short phrases ("Hello", "Yes") that lack vocal information for cloning.ContinuityMaintain semantic continuity — pause infrequently and aim for 3+ seconds of uninterrupted speech per segment.Emotional expressionAdd appropriate emotional expression (warmth, friendliness, seriousness). Monotone delivery reduces clone naturalness.Content restrictionsDo not include sensitive words (politics, pornography, violence). Recordings with this content will fail cloning. 
## [​ ](#end-to-end-examples) End-to-end examples

Create a cloned voice from a local audio file, then use it with a matching model. Both steps must use the same `target_model`.
Replace `voice.mp3` with the path to your own audio file.
### [​ ](#cosyvoice) CosyVoice

Clone a voice and use it with CosyVoice TTS. Applies to `cosyvoice-v3-plus` and `cosyvoice-v3-flash`. For details, see [CosyVoice TTS](/developer-guides/speech/realtime-streaming).
 CosyVoice voice cloning uses a different API than Qwen-TTS/Qwen-Omni: the cloning model is `voice-enrollment` (not `qwen-voice-enrollment`), the action is `create_voice`, and the audio is passed as a public URL (not base64). 
**Step 1: Create a voice**
The `url` parameter must be a publicly accessible URL of your audio file. The `prefix` parameter sets a prefix for the voice name.
Copy ```\ncurl -X POST https://dashscope-intl.aliyuncs.com/api/v1/services/audio/tts/customization \
-H "Authorization: Bearer $DASHSCOPE_API_KEY" \
-H "Content-Type: application/json" \
-d &#x27;{
 "model": "voice-enrollment",
 "input": {
 "action": "create_voice",
 "target_model": "cosyvoice-v3-plus",
 "prefix": "myvoice",
 "url": "https://your-audio-url.wav",
 "language_hints": ["en"]
 }
}&#x27;

``` 
**Step 2: Synthesize speech with the cloned voice**
Replace `YOUR_VOICE_ID` with the voice value returned in the previous step.
Copy ```\ncurl -X POST https://dashscope-intl.aliyuncs.com/api/v1/services/audio/tts/SpeechSynthesizer \
-H "Authorization: Bearer $DASHSCOPE_API_KEY" \
-H "Content-Type: application/json" \
-d &#x27;{
 "model": "cosyvoice-v3-plus",
 "input": {
 "text": "How is the weather today?",
 "voice": "YOUR_VOICE_ID",
 "format": "wav",
 "sample_rate": 24000
 }
}&#x27;

``` 
For SDK examples (Python, Java), see [CosyVoice voice cloning SDK](/api-reference/speech-synthesis/voice-cloning/cosyvoice-sdk/python-sdk).
### [​ ](#qwen-omni-realtime-conversation) Qwen-Omni: Realtime conversation

Clone a voice and use it in a realtime conversation. Applies to `qwen3.5-omni-plus-realtime` and `qwen3.5-omni-flash-realtime`. For details, see [Real-time multimodal speech](/developer-guides/speech/realtime-multimodal-speech).
- Python 
- Java 

 Copy ```\n# Requirements: dashscope >= 1.23.9, pyaudio
import os
import requests
import base64
import pathlib
import time
import pyaudio
from dashscope.audio.qwen_omni import MultiModality, OmniRealtimeCallback, OmniRealtimeConversation
import dashscope

TARGET_MODEL = "qwen3.5-omni-plus-realtime"
PREFERRED_NAME = "guanyu"
VOICE_FILE_PATH = "voice.mp3"

def create_voice(file_path: str) -> str:
 api_key = os.getenv("DASHSCOPE_API_KEY")
 file_path_obj = pathlib.Path(file_path)
 base64_str = base64.b64encode(file_path_obj.read_bytes()).decode()
 data_uri = f"data:audio/mpeg;base64,{base64_str}"

 resp = requests.post(
 "https://dashscope-intl.aliyuncs.com/api/v1/services/audio/tts/customization",
 headers={"Authorization": f"Bearer {api_key}", "Content-Type": "application/json"},
 json={
 "model": "qwen-voice-enrollment",
 "input": {
 "action": "create",
 "target_model": TARGET_MODEL,
 "preferred_name": PREFERRED_NAME,
 "audio": {"data": data_uri}
 }
 }
 )
 if resp.status_code != 200:
 raise RuntimeError(f"Failed to create voice: {resp.status_code}, {resp.text}")
 return resp.json()["output"]["voice"]

class SimpleCallback(OmniRealtimeCallback):
 def __init__(self, pya):
 self.pya = pya
 self.out = None
 def on_open(self):
 self.out = self.pya.open(format=pyaudio.paInt16, channels=1, rate=24000, output=True)
 def on_event(self, response):
 if response[&#x27;type&#x27;] == &#x27;response.audio.delta&#x27;:
 self.out.write(base64.b64decode(response[&#x27;delta&#x27;]))
 elif response[&#x27;type&#x27;] == &#x27;conversation.item.input_audio_transcription.completed&#x27;:
 print(f"[User] {response[&#x27;transcript&#x27;]}")
 elif response[&#x27;type&#x27;] == &#x27;response.audio_transcript.done&#x27;:
 print(f"[LLM] {response[&#x27;transcript&#x27;]}")

if __name__ == &#x27;__main__&#x27;:
 dashscope.api_key = os.getenv("DASHSCOPE_API_KEY")

 # Step 1: Clone a voice
 voice = create_voice(VOICE_FILE_PATH)
 print(f"Voice cloning complete. Voice: {voice}")

 # Step 2: Start a conversation with the cloned voice
 pya = pyaudio.PyAudio()
 callback = SimpleCallback(pya)
 conv = OmniRealtimeConversation(
 model=TARGET_MODEL, callback=callback,
 url="wss://dashscope-intl.aliyuncs.com/api-ws/v1/realtime"
 )
 conv.connect()
 conv.update_session(
 output_modalities=[MultiModality.AUDIO, MultiModality.TEXT],
 voice=voice
 )
 mic = pya.open(format=pyaudio.paInt16, channels=1, rate=16000, input=True)
 print("Conversation started. Speak into your microphone (Ctrl+C to exit)...")
 try:
 while True:
 audio_data = mic.read(3200, exception_on_overflow=False)
 conv.append_audio(base64.b64encode(audio_data).decode())
 time.sleep(0.01)
 except KeyboardInterrupt:
 conv.close()
 mic.close()
 callback.out.close()
 pya.terminate()
 print("\nConversation ended")

``` Copy ```\nimport com.alibaba.dashscope.audio.omni.*;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.google.gson.Gson;
import com.google.gson.JsonObject;

import javax.sound.sampled.*;
import java.io.*;
import java.net.HttpURLConnection;
import java.net.URL;
import java.nio.ByteBuffer;
import java.nio.file.*;
import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.Base64;
import java.util.Queue;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.atomic.AtomicBoolean;

public class Main {
 private static final String TARGET_MODEL = "qwen3.5-omni-plus-realtime";
 private static final String PREFERRED_NAME = "guanyu";
 private static final String AUDIO_FILE = "voice.mp3";
 private static final String AUDIO_MIME_TYPE = "audio/mpeg";

 public static String toDataUrl(String filePath) throws IOException {
 byte[] bytes = Files.readAllBytes(Paths.get(filePath));
 String encoded = Base64.getEncoder().encodeToString(bytes);
 return "data:" + AUDIO_MIME_TYPE + ";base64," + encoded;
 }

 public static String createVoice() throws Exception {
 String apiKey = System.getenv("DASHSCOPE_API_KEY");
 String jsonPayload =
 "{"
 + "\"model\": \"qwen-voice-enrollment\","
 + "\"input\": {"
 + "\"action\": \"create\","
 + "\"target_model\": \"" + TARGET_MODEL + "\","
 + "\"preferred_name\": \"" + PREFERRED_NAME + "\","
 + "\"audio\": {"
 + "\"data\": \"" + toDataUrl(AUDIO_FILE) + "\""
 + "}"
 + "}"
 + "}";

 String url = "https://dashscope-intl.aliyuncs.com/api/v1/services/audio/tts/customization";
 HttpURLConnection con = (HttpURLConnection) new URL(url).openConnection();
 con.setRequestMethod("POST");
 con.setRequestProperty("Authorization", "Bearer " + apiKey);
 con.setRequestProperty("Content-Type", "application/json");
 con.setDoOutput(true);

 try (OutputStream os = con.getOutputStream()) {
 os.write(jsonPayload.getBytes(StandardCharsets.UTF_8));
 }

 int status = con.getResponseCode();
 try (BufferedReader br = new BufferedReader(
 new InputStreamReader(status >= 200 && status < 300 ? con.getInputStream() : con.getErrorStream(),
 StandardCharsets.UTF_8))) {
 StringBuilder response = new StringBuilder();
 String line;
 while ((line = br.readLine()) != null) {
 response.append(line);
 }
 if (status == 200) {
 JsonObject jsonObj = new Gson().fromJson(response.toString(), JsonObject.class);
 return jsonObj.getAsJsonObject("output").get("voice").getAsString();
 }
 throw new IOException("Failed to create voice: " + status + " - " + response);
 }
 }

 static class SimpleAudioPlayer {
 private final SourceDataLine line;
 private final Queue<byte[]> audioQueue = new ConcurrentLinkedQueue<>();
 private final Thread playerThread;
 private final AtomicBoolean shouldStop = new AtomicBoolean(false);

 public SimpleAudioPlayer() throws LineUnavailableException {
 AudioFormat format = new AudioFormat(24000, 16, 1, true, false);
 line = AudioSystem.getSourceDataLine(format);
 line.open(format);
 line.start();
 playerThread = new Thread(() -> {
 while (!shouldStop.get()) {
 byte[] audio = audioQueue.poll();
 if (audio != null) {
 line.write(audio, 0, audio.length);
 } else {
 try { Thread.sleep(10); } catch (InterruptedException ignored) {}
 }
 }
 }, "AudioPlayer");
 playerThread.start();
 }

 public void play(String base64Audio) {
 audioQueue.add(Base64.getDecoder().decode(base64Audio));
 }

 public void close() {
 shouldStop.set(true);
 try { playerThread.join(1000); } catch (InterruptedException ignored) {}
 line.drain();
 line.close();
 }
 }

 public static void main(String[] args) {
 try {
 String voice = createVoice();
 System.out.println("Voice cloning complete. Voice: " + voice);

 SimpleAudioPlayer player = new SimpleAudioPlayer();
 AtomicBoolean shouldStop = new AtomicBoolean(false);

 OmniRealtimeParam param = OmniRealtimeParam.builder()
 .model(TARGET_MODEL)
 .apikey(System.getenv("DASHSCOPE_API_KEY"))
 .url("wss://dashscope-intl.aliyuncs.com/api-ws/v1/realtime")
 .build();

 OmniRealtimeConversation conversation = new OmniRealtimeConversation(param, new OmniRealtimeCallback() {
 @Override public void onOpen() { System.out.println("Connection established"); }
 @Override public void onClose(int code, String reason) {
 System.out.println("Connection closed (" + code + "): " + reason);
 shouldStop.set(true);
 }
 @Override public void onEvent(JsonObject event) {
 String type = event.get("type").getAsString();
 if ("response.audio.delta".equals(type)) {
 player.play(event.get("delta").getAsString());
 } else if ("conversation.item.input_audio_transcription.completed".equals(type)) {
 System.out.println("[User] " + event.get("transcript").getAsString());
 } else if ("response.audio_transcript.done".equals(type)) {
 System.out.println("[LLM] " + event.get("transcript").getAsString());
 }
 }
 });

 conversation.connect();
 conversation.updateSession(OmniRealtimeConfig.builder()
 .modalities(Arrays.asList(OmniRealtimeModality.AUDIO, OmniRealtimeModality.TEXT))
 .voice(voice)
 .enableTurnDetection(true)
 .enableInputAudioTranscription(true)
 .build()
 );

 System.out.println("Conversation started. Speak into the microphone (Ctrl+C to exit)...");
 AudioFormat format = new AudioFormat(16000, 16, 1, true, false);
 TargetDataLine mic = AudioSystem.getTargetDataLine(format);
 mic.open(format);
 mic.start();

 ByteBuffer buffer = ByteBuffer.allocate(3200);
 while (!shouldStop.get()) {
 int bytesRead = mic.read(buffer.array(), 0, buffer.capacity());
 if (bytesRead > 0) {
 conversation.appendAudio(Base64.getEncoder().encodeToString(buffer.array()));
 }
 Thread.sleep(20);
 }

 conversation.close(1000, "Normal exit");
 player.close();
 mic.close();
 System.out.println("\nConversation ended");
 } catch (NoApiKeyException e) {
 System.err.println("API KEY not found. Set the DASHSCOPE_API_KEY environment variable.");
 } catch (Exception e) {
 e.printStackTrace();
 }
 System.exit(0);
 }
}

``` 
### [​ ](#qwen-omni-non-realtime-conversation) Qwen-Omni: Non-realtime conversation

Clone a voice and use it in a non-realtime conversation. Applies to `qwen3.5-omni-plus` and `qwen3.5-omni-flash`. For details, see [Multimodal speech](/developer-guides/speech/multimodal-speech).
- Python 
- Java 

 Copy ```\nimport os
import requests
import base64
import pathlib
import numpy as np
import soundfile as sf
import dashscope

TARGET_MODEL = "qwen3.5-omni-plus"
PREFERRED_NAME = "guanyu"
VOICE_FILE_PATH = "voice.mp3"

def create_voice(file_path: str) -> str:
 api_key = os.getenv("DASHSCOPE_API_KEY")
 file_path_obj = pathlib.Path(file_path)
 base64_str = base64.b64encode(file_path_obj.read_bytes()).decode()
 data_uri = f"data:audio/mpeg;base64,{base64_str}"

 resp = requests.post(
 "https://dashscope-intl.aliyuncs.com/api/v1/services/audio/tts/customization",
 headers={"Authorization": f"Bearer {api_key}", "Content-Type": "application/json"},
 json={
 "model": "qwen-voice-enrollment",
 "input": {
 "action": "create",
 "target_model": TARGET_MODEL,
 "preferred_name": PREFERRED_NAME,
 "audio": {"data": data_uri}
 }
 }
 )
 if resp.status_code != 200:
 raise RuntimeError(f"Failed to create voice: {resp.status_code}, {resp.text}")
 return resp.json()["output"]["voice"]

if __name__ == &#x27;__main__&#x27;:
 dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

 voice = create_voice(VOICE_FILE_PATH)
 print(f"Voice cloning complete. Voice: {voice}")

 messages = [{"role": "user", "content": [{"text": "Hello, please introduce yourself"}]}]

 response = dashscope.MultiModalConversation.call(
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 model=TARGET_MODEL,
 messages=messages,
 modalities=["text", "audio"],
 audio={"voice": voice, "format": "wav"},
 stream=True
 )

 print("Model response:")
 audio_base64_string = ""
 for r in response:
 try:
 content = r.output.choices[0].message.content
 for item in content:
 if isinstance(item, dict):
 if "audio" in item:
 audio_base64_string += item["audio"].get("data", "")
 elif "text" in item:
 print(item["text"], end="")
 except Exception:
 pass

 if audio_base64_string:
 wav_bytes = base64.b64decode(audio_base64_string)
 audio_np = np.frombuffer(wav_bytes, dtype=np.int16)
 sf.write("audio_cloned.wav", audio_np, samplerate=24000)
 print("\nAudio saved to: audio_cloned.wav")

``` Copy ```\nimport com.google.gson.Gson;
import com.google.gson.JsonObject;

import java.io.*;
import java.net.HttpURLConnection;
import java.net.URL;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.file.*;
import java.nio.charset.StandardCharsets;
import java.util.Base64;

public class Main {
 private static final String TARGET_MODEL = "qwen3.5-omni-plus";
 private static final String PREFERRED_NAME = "guanyu";
 private static final String AUDIO_FILE = "voice.mp3";
 private static final String AUDIO_MIME_TYPE = "audio/mpeg";

 public static void writeWav(String path, byte[] pcmData, int sampleRate) throws IOException {
 int channels = 1, bitsPerSample = 16;
 int byteRate = sampleRate * channels * bitsPerSample / 8;
 int blockAlign = channels * bitsPerSample / 8;
 ByteBuffer header = ByteBuffer.allocate(44).order(ByteOrder.LITTLE_ENDIAN);
 header.put("RIFF".getBytes()); header.putInt(36 + pcmData.length);
 header.put("WAVE".getBytes()); header.put("fmt ".getBytes());
 header.putInt(16); header.putShort((short) 1); header.putShort((short) channels);
 header.putInt(sampleRate); header.putInt(byteRate);
 header.putShort((short) blockAlign); header.putShort((short) bitsPerSample);
 header.put("data".getBytes()); header.putInt(pcmData.length);
 try (FileOutputStream fos = new FileOutputStream(path)) {
 fos.write(header.array());
 fos.write(pcmData);
 }
 }

 public static String toDataUrl(String filePath) throws IOException {
 byte[] bytes = Files.readAllBytes(Paths.get(filePath));
 String encoded = Base64.getEncoder().encodeToString(bytes);
 return "data:" + AUDIO_MIME_TYPE + ";base64," + encoded;
 }

 public static String createVoice() throws Exception {
 String apiKey = System.getenv("DASHSCOPE_API_KEY");
 String jsonPayload =
 "{"
 + "\"model\": \"qwen-voice-enrollment\","
 + "\"input\": {"
 + "\"action\": \"create\","
 + "\"target_model\": \"" + TARGET_MODEL + "\","
 + "\"preferred_name\": \"" + PREFERRED_NAME + "\","
 + "\"audio\": {"
 + "\"data\": \"" + toDataUrl(AUDIO_FILE) + "\""
 + "}"
 + "}"
 + "}";

 String url = "https://dashscope-intl.aliyuncs.com/api/v1/services/audio/tts/customization";
 HttpURLConnection con = (HttpURLConnection) new URL(url).openConnection();
 con.setRequestMethod("POST");
 con.setRequestProperty("Authorization", "Bearer " + apiKey);
 con.setRequestProperty("Content-Type", "application/json");
 con.setDoOutput(true);

 try (OutputStream os = con.getOutputStream()) {
 os.write(jsonPayload.getBytes(StandardCharsets.UTF_8));
 }

 int status = con.getResponseCode();
 try (BufferedReader br = new BufferedReader(
 new InputStreamReader(status >= 200 && status < 300 ? con.getInputStream() : con.getErrorStream(),
 StandardCharsets.UTF_8))) {
 StringBuilder response = new StringBuilder();
 String line;
 while ((line = br.readLine()) != null) {
 response.append(line);
 }
 if (status == 200) {
 JsonObject jsonObj = new Gson().fromJson(response.toString(), JsonObject.class);
 return jsonObj.getAsJsonObject("output").get("voice").getAsString();
 }
 throw new IOException("Failed to create voice: " + status + " - " + response);
 }
 }

 public static void main(String[] args) {
 try {
 String apiKey = System.getenv("DASHSCOPE_API_KEY");

 String voice = createVoice();
 System.out.println("Voice cloning complete. Voice: " + voice);

 String requestBody = "{"
 + "\"model\": \"" + TARGET_MODEL + "\","
 + "\"messages\": [{\"role\": \"user\", \"content\": \"Hello, please introduce yourself\"}],"
 + "\"modalities\": [\"text\", \"audio\"],"
 + "\"audio\": {\"voice\": \"" + voice + "\", \"format\": \"wav\"},"
 + "\"stream\": true,"
 + "\"stream_options\": {\"include_usage\": true}"
 + "}";

 String url = "https://dashscope-intl.aliyuncs.com/compatible-mode/v1/chat/completions";
 HttpURLConnection con = (HttpURLConnection) new URL(url).openConnection();
 con.setRequestMethod("POST");
 con.setRequestProperty("Authorization", "Bearer " + apiKey);
 con.setRequestProperty("Content-Type", "application/json");
 con.setDoOutput(true);

 try (OutputStream os = con.getOutputStream()) {
 os.write(requestBody.getBytes(StandardCharsets.UTF_8));
 }

 StringBuilder audioBase64 = new StringBuilder();
 System.out.println("Model response:");

 try (BufferedReader br = new BufferedReader(
 new InputStreamReader(con.getInputStream(), StandardCharsets.UTF_8))) {
 String line;
 while ((line = br.readLine()) != null) {
 if (!line.startsWith("data: ") || line.equals("data: [DONE]")) continue;
 String json = line.substring(6);
 JsonObject chunk = new Gson().fromJson(json, JsonObject.class);
 if (!chunk.has("choices") || chunk.getAsJsonArray("choices").size() == 0) continue;

 JsonObject delta = chunk.getAsJsonArray("choices").get(0)
 .getAsJsonObject().getAsJsonObject("delta");
 if (delta.has("content") && !delta.get("content").isJsonNull()) {
 System.out.print(delta.get("content").getAsString());
 }
 if (delta.has("audio") && !delta.get("audio").isJsonNull()) {
 JsonObject audio = delta.getAsJsonObject("audio");
 if (audio.has("data")) {
 audioBase64.append(audio.get("data").getAsString());
 }
 }
 }
 }

 if (audioBase64.length() > 0) {
 byte[] pcmBytes = Base64.getDecoder().decode(audioBase64.toString());
 writeWav("audio_cloned.wav", pcmBytes, 24000);
 System.out.println("\nAudio saved to: audio_cloned.wav");
 }
 } catch (Exception e) {
 e.printStackTrace();
 }
 }
}

``` 
### [​ ](#qwen-tts-bidirectional-streaming-real-time) Qwen-TTS: Bidirectional streaming (real-time)

Applies to Qwen3-TTS-VC-Realtime models. For parameter details, see [Realtime streaming TTS](/developer-guides/speech/realtime-streaming).
- Python 
- Java 

 Copy ```\n# pyaudio installation:
# macOS: brew install portaudio && pip install pyaudio
# Ubuntu: sudo apt-get install python3-pyaudio (or pip install pyaudio)
# CentOS: sudo yum install -y portaudio portaudio-devel && pip install pyaudio
# Windows: python -m pip install pyaudio

import pyaudio
import os
import requests
import base64
import pathlib
import threading
import time
import dashscope
from dashscope.audio.qwen_tts_realtime import QwenTtsRealtime, QwenTtsRealtimeCallback, AudioFormat

TARGET_MODEL = "qwen3-tts-vc-realtime-2026-01-15"
VOICE_FILE = "voice.mp3" # Replace with your audio file

TEXT_TO_SYNTHESIZE = [
 &#x27;Today we explore the wonders of speech synthesis.&#x27;,
 &#x27;Each voice carries a unique character.&#x27;,
 &#x27;With voice cloning, you can bring any text to life.&#x27;,
 "Let&#x27;s create something amazing together."
]

def create_voice(file_path: str) -> str:
 """Create a cloned voice and return the voice identifier."""
 api_key = os.getenv("DASHSCOPE_API_KEY")
 file_path_obj = pathlib.Path(file_path)
 base64_str = base64.b64encode(file_path_obj.read_bytes()).decode()
 data_uri = f"data:audio/mpeg;base64,{base64_str}"

 response = requests.post(
 "https://dashscope-intl.aliyuncs.com/api/v1/services/audio/tts/customization",
 headers={"Authorization": f"Bearer {api_key}", "Content-Type": "application/json"},
 json={
 "model": "qwen-voice-enrollment",
 "input": {
 "action": "create",
 "target_model": TARGET_MODEL,
 "preferred_name": "myvoice",
 "audio": {"data": data_uri}
 }
 }
 )
 return response.json()["output"]["voice"]

class MyCallback(QwenTtsRealtimeCallback):
 def __init__(self):
 self.complete_event = threading.Event()
 self._player = pyaudio.PyAudio()
 self._stream = self._player.open(format=pyaudio.paInt16, channels=1, rate=24000, output=True)

 def on_event(self, response: dict) -> None:
 if response.get("type") == "response.audio.delta":
 audio_data = base64.b64decode(response["delta"])
 self._stream.write(audio_data)
 elif response.get("type") == "session.finished":
 self.complete_event.set()

if __name__ == "__main__":
 dashscope.api_key = os.getenv("DASHSCOPE_API_KEY")
 callback = MyCallback()
 tts = QwenTtsRealtime(model=TARGET_MODEL, callback=callback,
 url="wss://dashscope-intl.aliyuncs.com/api-ws/v1/realtime")
 tts.connect()
 tts.update_session(voice=create_voice(VOICE_FILE),
 response_format=AudioFormat.PCM_24000HZ_MONO_16BIT, mode="server_commit")

 for text in TEXT_TO_SYNTHESIZE:
 tts.append_text(text)
 time.sleep(0.1)

 tts.finish()
 callback.complete_event.wait()

``` Copy ```\nimport com.alibaba.dashscope.audio.qwen_tts_realtime.*;
import com.google.gson.Gson;
import com.google.gson.JsonObject;
import java.io.*;
import java.net.HttpURLConnection;
import java.net.URL;
import java.nio.file.*;
import java.util.Base64;
import java.util.concurrent.CountDownLatch;

public class Main {
 private static final String TARGET_MODEL = "qwen3-tts-vc-realtime-2026-01-15";
 private static final String AUDIO_FILE = "voice.mp3"; // Replace with your audio file

 public static String createVoice() throws Exception {
 String apiKey = System.getenv("DASHSCOPE_API_KEY");
 byte[] bytes = Files.readAllBytes(Paths.get(AUDIO_FILE));
 String encoded = Base64.getEncoder().encodeToString(bytes);
 String dataUri = "data:audio/mpeg;base64," + encoded;

 String jsonPayload = "{\"model\":\"qwen-voice-enrollment\",\"input\":{"
 + "\"action\":\"create\",\"target_model\":\"" + TARGET_MODEL + "\","
 + "\"preferred_name\":\"myvoice\",\"audio\":{\"data\":\"" + dataUri + "\"}}}";

 HttpURLConnection con = (HttpURLConnection) new URL(
 "https://dashscope-intl.aliyuncs.com/api/v1/services/audio/tts/customization").openConnection();
 con.setRequestMethod("POST");
 con.setRequestProperty("Authorization", "Bearer " + apiKey);
 con.setRequestProperty("Content-Type", "application/json");
 con.setDoOutput(true);
 try (OutputStream os = con.getOutputStream()) {
 os.write(jsonPayload.getBytes("UTF-8"));
 }

 BufferedReader br = new BufferedReader(new InputStreamReader(con.getInputStream(), "UTF-8"));
 StringBuilder response = new StringBuilder();
 String line;
 while ((line = br.readLine()) != null) response.append(line);
 return new Gson().fromJson(response.toString(), JsonObject.class)
 .getAsJsonObject("output").get("voice").getAsString();
 }

 public static void main(String[] args) throws Exception {
 CountDownLatch latch = new CountDownLatch(1);
 QwenTtsRealtimeParam param = QwenTtsRealtimeParam.builder()
 .model(TARGET_MODEL)
 .url("wss://dashscope-intl.aliyuncs.com/api-ws/v1/realtime")
 .apikey(System.getenv("DASHSCOPE_API_KEY"))
 .build();

 QwenTtsRealtime tts = new QwenTtsRealtime(param, new QwenTtsRealtimeCallback() {
 public void onEvent(JsonObject msg) {
 if (msg.get("type").getAsString().equals("session.finished")) latch.countDown();
 }
 });
 tts.connect();

 QwenTtsRealtimeConfig config = QwenTtsRealtimeConfig.builder()
 .voice(createVoice())
 .responseFormat(QwenTtsRealtimeAudioFormat.PCM_24000HZ_MONO_16BIT)
 .mode("server_commit").build();
 tts.updateSession(config);

 for (String text : new String[]{
 "Today we explore the wonders of speech synthesis.",
 "Each voice carries a unique character.",
 "With voice cloning, you can bring any text to life.",
 "Let&#x27;s create something amazing together."}) {
 tts.appendText(text);
 Thread.sleep(100);
 }
 tts.finish();
 latch.await();
 }
}

``` 
### [​ ](#qwen-tts-non-streaming-synthesis) Qwen-TTS: Non-streaming synthesis

Applies to Qwen3-TTS-VC models. For details, see [Qwen TTS](/api-reference/speech-synthesis/qwen-tts).
- Python 
- Java 

 Copy ```\nimport os
import requests
import base64
import pathlib
import dashscope

TARGET_MODEL = "qwen3-tts-vc-2026-01-22"
VOICE_FILE = "voice.mp3"

def create_voice(file_path: str) -> str:
 api_key = os.getenv("DASHSCOPE_API_KEY")
 base64_str = base64.b64encode(pathlib.Path(file_path).read_bytes()).decode()
 data_uri = f"data:audio/mpeg;base64,{base64_str}"

 response = requests.post(
 "https://dashscope-intl.aliyuncs.com/api/v1/services/audio/tts/customization",
 headers={"Authorization": f"Bearer {api_key}", "Content-Type": "application/json"},
 json={
 "model": "qwen-voice-enrollment",
 "input": {"action": "create", "target_model": TARGET_MODEL,
 "preferred_name": "myvoice", "audio": {"data": data_uri}}
 }
 )
 return response.json()["output"]["voice"]

if __name__ == "__main__":
 dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;
 response = dashscope.MultiModalConversation.call(
 model=TARGET_MODEL,
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 text="Today we explore the wonders of speech synthesis.",
 voice=create_voice(VOICE_FILE),
 stream=False
 )
 print(response)

``` Copy ```\nimport com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversation;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationParam;
import com.alibaba.dashscope.aigc.multimodalconversation.MultiModalConversationResult;
import com.alibaba.dashscope.utils.Constants;
import com.google.gson.Gson;
import com.google.gson.JsonObject;

import java.io.*;
import java.net.HttpURLConnection;
import java.net.URL;
import java.nio.file.*;
import java.nio.charset.StandardCharsets;
import java.util.Base64;

public class Main {
 private static final String TARGET_MODEL = "qwen3-tts-vc-2026-01-22";
 private static final String AUDIO_FILE = "voice.mp3"; // Replace with your audio file

 public static String createVoice() throws Exception {
 String apiKey = System.getenv("DASHSCOPE_API_KEY");
 byte[] bytes = Files.readAllBytes(Paths.get(AUDIO_FILE));
 String encoded = Base64.getEncoder().encodeToString(bytes);
 String dataUri = "data:audio/mpeg;base64," + encoded;

 String jsonPayload = "{\"model\":\"qwen-voice-enrollment\",\"input\":{"
 + "\"action\":\"create\",\"target_model\":\"" + TARGET_MODEL + "\","
 + "\"preferred_name\":\"myvoice\",\"audio\":{\"data\":\"" + dataUri + "\"}}}";

 HttpURLConnection con = (HttpURLConnection) new URL(
 "https://dashscope-intl.aliyuncs.com/api/v1/services/audio/tts/customization").openConnection();
 con.setRequestMethod("POST");
 con.setRequestProperty("Authorization", "Bearer " + apiKey);
 con.setRequestProperty("Content-Type", "application/json");
 con.setDoOutput(true);
 try (OutputStream os = con.getOutputStream()) {
 os.write(jsonPayload.getBytes(StandardCharsets.UTF_8));
 }

 BufferedReader br = new BufferedReader(new InputStreamReader(con.getInputStream(), StandardCharsets.UTF_8));
 StringBuilder response = new StringBuilder();
 String line;
 while ((line = br.readLine()) != null) response.append(line);
 return new Gson().fromJson(response.toString(), JsonObject.class)
 .getAsJsonObject("output").get("voice").getAsString();
 }

 public static void main(String[] args) {
 try {
 Constants.baseHttpApiUrl = "https://dashscope-intl.aliyuncs.com/api/v1";
 MultiModalConversation conv = new MultiModalConversation();
 MultiModalConversationParam param = MultiModalConversationParam.builder()
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model(TARGET_MODEL)
 .text("Today we explore the wonders of speech synthesis.")
 .parameter("voice", createVoice())
 .build();
 MultiModalConversationResult result = conv.call(param);
 String audioUrl = result.getOutput().getAudio().getUrl();
 System.out.println("Audio URL: " + audioUrl);

 // Download audio
 try (InputStream in = new URL(audioUrl).openStream();
 FileOutputStream out = new FileOutputStream("output.wav")) {
 byte[] buffer = new byte[1024];
 int bytesRead;
 while ((bytesRead = in.read(buffer)) != -1) {
 out.write(buffer, 0, bytesRead);
 }
 System.out.println("Audio saved to output.wav");
 }
 } catch (Exception e) {
 System.out.println("Error: " + e.getMessage());
 }
 System.exit(0);
 }
}

``` 
## [​ ](#quota-and-billing) Quota and billing

### [​ ](#voice-quota-and-automatic-cleanup) Voice quota and automatic cleanup

**Total voice limit**: Each Qwen Cloud account has a separate limit of 1,000 custom voices for CosyVoice and 1,000 for Qwen-TTS. The two quotas are counted independently.
**Automatic cleanup**: If a voice isn&#x27;t used in any speech synthesis request for one year, the system automatically deletes it.
### [​ ](#billing-rules) Billing rules

 The prices listed below are list prices. For current promotions and discounted pricing, visit the [Model Marketplace](https://www.qwencloud.com/models). 

- 
**CosyVoice**: Voice cloning is free.


- 
**Qwen-TTS**: Each voice cloning costs USD 0.01. Failed creations aren&#x27;t charged. Voice design has [separate pricing](/developer-guides/speech/voice-design#billing-rules).
**Free quota** (Singapore region only):

You get 1,000 free voice cloning creations during the first 90 days after activating Qwen Cloud.

- Failed creations don&#x27;t consume the free quota.

- Deleting a voice doesn&#x27;t restore the free quota.

- After the free quota is used up or the 90-day window expires, voice cloning is billed at USD 0.01 per voice.


## [​ ](#troubleshooting) Troubleshooting

If you encounter errors, see [Error messages](/api-reference/preparation/error-messages).
## [​ ](#next-steps) Next steps


- [Voice cloning API reference (Qwen)](/api-reference/speech-synthesis/voice-cloning/qwen/create-voice) -- Qwen-TTS/Qwen-Omni voice cloning API

- [Voice cloning API reference (CosyVoice)](/api-reference/speech-synthesis/voice-cloning/cosyvoice/create-voice) -- CosyVoice voice cloning API

- [Real-time multimodal speech](/developer-guides/speech/realtime-multimodal-speech) -- Use cloned voices in realtime conversation

- [Multimodal speech](/developer-guides/speech/multimodal-speech) -- Use cloned voices in non-realtime conversation

- [Realtime streaming TTS](/developer-guides/speech/realtime-streaming) -- Bidirectional streaming details

- [Get an API key](/api-reference/preparation/api-key) -- Set up authentication

- [Install the DashScope SDK](/api-reference/preparation/install-sdk) -- SDK installation


 [Previous ](/developer-guides/speech/tts)[Voice design Create custom voices from text descriptions for use with CosyVoice or Qwen TTS models. Next ](/developer-guides/speech/voice-design)
