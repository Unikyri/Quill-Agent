# Real-time speech synthesis

> **Source:** https://docs.qwencloud.com/developer-guides/speech/realtime-streaming

Stream TTS in real time

 Copy page Real-time speech synthesis converts text into natural speech over a WebSocket connection. It supports streaming input and output, voice cloning, voice design, and fine-grained audio control for use cases such as voice assistants, audiobooks, and intelligent customer service.
## [​ ](#overview) Overview

Low-latency real-time speech synthesis over WebSocket, built for voice assistants, intelligent customer service, live captioning, and other scenarios that require instant responses.

- Streaming input and output (full-duplex WebSocket) with low time to first audio, ideal for real-time conversations such as voice assistants and intelligent customer service

- Adjustable speech rate, pitch, volume, and bitrate for fine-grained voice control

- Compatible with mainstream audio formats (PCM, WAV, MP3, Opus) and supports up to 48 kHz sample rate output

- Supports [instruction-based control](#instruction-based-control), allowing natural-language instructions to control voice expressiveness

- Supports [voice cloning](/api-reference/speech-synthesis/voice-cloning/qwen/create-voice) and [Voice Design](/api-reference/speech-synthesis/voice-design/qwen/create-voice) voice customization


If you don&#x27;t need real-time output, use non-real-time speech synthesis (HTTP API), which is suited for batch scenarios such as audiobooks and courseware dubbing. For model selection guidance, see [Speech synthesis models](/developer-guides/speech/tts-models).
## [​ ](#prerequisites) Prerequisites


- [Configure an API key](/api-reference/preparation/api-key) and [set it as an environment variable](/api-reference/preparation/export-api-key-env).

- If you call the service through the DashScope SDK, [install the latest SDK](/api-reference/preparation/install-sdk).


## [​ ](#quick-start) Quick start

The following examples demonstrate speech synthesis for each model. For more examples and parameter descriptions, see the API reference of each model.
- CosyVoice 
- Qwen-TTS 

 The following example shows how to synthesize speech with a system voice (see [CosyVoice voice list](/api-reference/speech-synthesis/cosyvoice/voice-list)).To use instruction control, configure instructions through the `instructions` parameter.- Python 
- Java 

 Copy ```\n# coding=utf-8
import os
import dashscope
from dashscope.audio.tts_v2 import *
# Obtain an API key from the Qwen Cloud console.
# If the environment variable is not configured, replace the following line with your Qwen Cloud API key: dashscope.api_key = "sk-xxx"
dashscope.api_key = os.environ.get(&#x27;DASHSCOPE_API_KEY&#x27;)
dashscope.base_websocket_api_url=&#x27;wss://dashscope-intl.aliyuncs.com/api-ws/v1/inference&#x27;
# Model
# Different model versions require corresponding voice types:
# cosyvoice-v3-flash/cosyvoice-v3-plus: Use voices such as longanyang.
# Each voice supports different languages. When synthesizing non-Chinese languages such as Japanese or Korean, select a voice that supports the target language. For details, see the CosyVoice voice list.
model = "cosyvoice-v3-flash"
# Voice
voice = "longanyang"
# Instantiate SpeechSynthesizer and pass request parameters such as model and voice in the constructor
synthesizer = SpeechSynthesizer(model=model, voice=voice)
# Send text for synthesis and obtain binary audio
audio = synthesizer.call("How is the weather today?")
# The first text submission requires establishing a WebSocket connection, so the first-packet latency includes connection setup time
print(&#x27;[Metric] requestId: {}, first package delay: {} ms&#x27;.format(
 synthesizer.get_last_request_id(),
 synthesizer.get_first_package_delay()))
# Save audio to a local file
with open(&#x27;output.mp3&#x27;, &#x27;wb&#x27;) as f:
 f.write(audio)

``` Copy ```\nimport com.alibaba.dashscope.audio.ttsv2.SpeechSynthesisParam;
import com.alibaba.dashscope.audio.ttsv2.SpeechSynthesizer;
import com.alibaba.dashscope.utils.Constants;
import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.nio.ByteBuffer;
public class Main {
 // Model
 // Different model versions require corresponding voice types:
 // cosyvoice-v3-flash/cosyvoice-v3-plus: Use voices such as longanyang.
 // Each voice supports different languages. When synthesizing non-Chinese languages such as Japanese or Korean, select a voice that supports the target language. For details, see the CosyVoice voice list.
 private static String model = "cosyvoice-v3-flash";
 // Voice
 private static String voice = "longanyang";
 public static void streamAudioDataToSpeaker() {
 // Request parameters
 SpeechSynthesisParam param =
 SpeechSynthesisParam.builder()
 // Obtain an API key from the Qwen Cloud console.
 // If the environment variable is not configured, replace the following line with your Qwen Cloud API key: .apiKey("sk-xxx")
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .model(model) // Model
 .voice(voice) // Voice
 .build();
 // Synchronous mode: disable callback (second parameter is null)
 SpeechSynthesizer synthesizer = new SpeechSynthesizer(param, null);
 ByteBuffer audio = null;
 try {
 // Block until audio is returned
 audio = synthesizer.call("How is the weather today?");
 } catch (Exception e) {
 throw new RuntimeException(e);
 } finally {
 // Close the WebSocket connection after the task completes
 synthesizer.getDuplexApi().close(1000, "bye");
 }
 if (audio != null) {
 // Save the audio data to a local file "output.mp3"
 File file = new File("output.mp3");
 // The first text submission requires establishing a WebSocket connection, so the first-packet latency includes connection setup time
 // Note: getFirstPackageDelay() requires dashscope-sdk-java 2.18.0 or later
 System.out.println(
 "[Metric] requestId: "
 + synthesizer.getLastRequestId()
 + ", first package delay (ms): "
 + synthesizer.getFirstPackageDelay());
 try (FileOutputStream fos = new FileOutputStream(file)) {
 fos.write(audio.array());
 } catch (IOException e) {
 throw new RuntimeException(e);
 }
 }
 }
 public static void main(String[] args) {
 Constants.baseWebsocketApiUrl = "wss://dashscope-intl.aliyuncs.com/api-ws/v1/inference";
 streamAudioDataToSpeaker();
 System.exit(0);
 }
}

``` The following example shows how to synthesize speech with a system voice (see [Supported voices](#supported-voices)).To use [instruction-based control](#instruction-based-control), set `model` to `qwen3-tts-instruct-flash-realtime` and configure the instruction through the `instructions` parameter.- Python 
- Java 

 - Server commit mode 
- Commit mode 

 Copy ```\nimport os
import base64
import threading
import time
import dashscope
from dashscope.audio.qwen_tts_realtime import *
qwen_tts_realtime: QwenTtsRealtime = None
text_to_synthesize = [
 &#x27;Right? I love supermarkets like this.&#x27;,
 &#x27;Especially during Chinese New Year,&#x27;,
 &#x27;I go shopping at supermarkets.&#x27;,
 &#x27;And I feel&#x27;,
 &#x27;absolutely thrilled!&#x27;,
 &#x27;I want to buy so many things!&#x27;
]
DO_VIDEO_TEST = False
def init_dashscope_api_key():
 """
 Set your DashScope API key. More information:
 https://github.com/aliyun/alibabacloud-bailian-speech-demo/blob/master/PREREQUISITES.md
 """
 # Obtain an API key from the Qwen Cloud console.
 if &#x27;DASHSCOPE_API_KEY&#x27; in os.environ:
 dashscope.api_key = os.environ[
 &#x27;DASHSCOPE_API_KEY&#x27;] # Load API key from environment variable DASHSCOPE_API_KEY
 else:
 dashscope.api_key = &#x27;your-dashscope-api-key&#x27; # Set API key manually
class MyCallback(QwenTtsRealtimeCallback):
 def __init__(self):
 self.complete_event = threading.Event()
 self.file = open(&#x27;result_24k.pcm&#x27;, &#x27;wb&#x27;)
 def on_open(self) -> None:
 print(&#x27;connection opened, init player&#x27;)
 def on_close(self, close_status_code, close_msg) -> None:
 self.file.close()
 print(&#x27;connection closed with code: {}, msg: {}, destroy player&#x27;.format(close_status_code, close_msg))
 def on_event(self, response: str) -> None:
 try:
 global qwen_tts_realtime
 type = response[&#x27;type&#x27;]
 if &#x27;session.created&#x27; == type:
 print(&#x27;start session: {}&#x27;.format(response[&#x27;session&#x27;][&#x27;id&#x27;]))
 if &#x27;response.audio.delta&#x27; == type:
 recv_audio_b64 = response[&#x27;delta&#x27;]
 self.file.write(base64.b64decode(recv_audio_b64))
 if &#x27;response.done&#x27; == type:
 print(f&#x27;response {qwen_tts_realtime.get_last_response_id()} done&#x27;)
 if &#x27;session.finished&#x27; == type:
 print(&#x27;session finished&#x27;)
 self.complete_event.set()
 except Exception as e:
 print(&#x27;[Error] {}&#x27;.format(e))
 return
 def wait_for_finished(self):
 self.complete_event.wait()
if __name__ == &#x27;__main__&#x27;:
 init_dashscope_api_key()
 print(&#x27;Initializing ...&#x27;)
 callback = MyCallback()
 qwen_tts_realtime = QwenTtsRealtime(
 # To use instruction control, replace the model with qwen3-tts-instruct-flash-realtime
 model=&#x27;qwen3-tts-flash-realtime&#x27;,
 callback=callback,
 url=&#x27;wss://dashscope-intl.aliyuncs.com/api-ws/v1/realtime&#x27;
 )
 qwen_tts_realtime.connect()
 qwen_tts_realtime.update_session(
 voice = &#x27;Cherry&#x27;,
 response_format = AudioFormat.PCM_24000HZ_MONO_16BIT,
 # To use instruction control, uncomment the following lines and replace the model with qwen3-tts-instruct-flash-realtime
 # instructions=&#x27;Speak quickly with a rising intonation, suitable for introducing fashion products.&#x27;,
 # optimize_instructions=True,
 mode = &#x27;server_commit&#x27; 
 )
 for text_chunk in text_to_synthesize:
 print(f&#x27;send text: {text_chunk}&#x27;)
 qwen_tts_realtime.append_text(text_chunk)
 time.sleep(0.1)
 qwen_tts_realtime.finish()
 callback.wait_for_finished()
 print(&#x27;[Metric] session: {}, first audio delay: {}&#x27;.format(
 qwen_tts_realtime.get_session_id(), 
 qwen_tts_realtime.get_first_audio_delay(),
 ))

``` Copy ```\nimport base64
import os
import threading
import dashscope
from dashscope.audio.qwen_tts_realtime import *
qwen_tts_realtime: QwenTtsRealtime = None
text_to_synthesize = [
 &#x27;This is the first sentence.&#x27;,
 &#x27;This is the second sentence.&#x27;,
 &#x27;This is the third sentence.&#x27;,
]
DO_VIDEO_TEST = False
def init_dashscope_api_key():
 """
 Set your DashScope API key. More information:
 https://github.com/aliyun/alibabacloud-bailian-speech-demo/blob/master/PREREQUISITES.md
 """
 # Obtain an API key from the Qwen Cloud console.
 if &#x27;DASHSCOPE_API_KEY&#x27; in os.environ:
 dashscope.api_key = os.environ[
 &#x27;DASHSCOPE_API_KEY&#x27;] # Load API key from environment variable DASHSCOPE_API_KEY
 else:
 dashscope.api_key = &#x27;your-dashscope-api-key&#x27; # Set API key manually
class MyCallback(QwenTtsRealtimeCallback):
 def __init__(self):
 super().__init__()
 self.response_counter = 0
 self.complete_event = threading.Event()
 self.file = open(f&#x27;result_{self.response_counter}_24k.pcm&#x27;, &#x27;wb&#x27;)
 def reset_event(self):
 self.response_counter += 1
 self.file = open(f&#x27;result_{self.response_counter}_24k.pcm&#x27;, &#x27;wb&#x27;)
 self.complete_event = threading.Event()
 def on_open(self) -> None:
 print(&#x27;connection opened, init player&#x27;)
 def on_close(self, close_status_code, close_msg) -> None:
 print(&#x27;connection closed with code: {}, msg: {}, destroy player&#x27;.format(close_status_code, close_msg))
 def on_event(self, response: str) -> None:
 try:
 global qwen_tts_realtime
 type = response[&#x27;type&#x27;]
 if &#x27;session.created&#x27; == type:
 print(&#x27;start session: {}&#x27;.format(response[&#x27;session&#x27;][&#x27;id&#x27;]))
 if &#x27;response.audio.delta&#x27; == type:
 recv_audio_b64 = response[&#x27;delta&#x27;]
 self.file.write(base64.b64decode(recv_audio_b64))
 if &#x27;response.done&#x27; == type:
 print(f&#x27;response {qwen_tts_realtime.get_last_response_id()} done&#x27;)
 self.complete_event.set()
 self.file.close()
 if &#x27;session.finished&#x27; == type:
 print(&#x27;session finished&#x27;)
 self.complete_event.set()
 except Exception as e:
 print(&#x27;[Error] {}&#x27;.format(e))
 return
 def wait_for_response_done(self):
 self.complete_event.wait()
if __name__ == &#x27;__main__&#x27;:
 init_dashscope_api_key()
 print(&#x27;Initializing ...&#x27;)
 callback = MyCallback()
 qwen_tts_realtime = QwenTtsRealtime(
 # To use instruction control, replace the model with qwen3-tts-instruct-flash-realtime
 model=&#x27;qwen3-tts-flash-realtime&#x27;,
 callback=callback, 
 url=&#x27;wss://dashscope-intl.aliyuncs.com/api-ws/v1/realtime&#x27;
 )
 qwen_tts_realtime.connect()
 qwen_tts_realtime.update_session(
 voice = &#x27;Cherry&#x27;,
 response_format = AudioFormat.PCM_24000HZ_MONO_16BIT,
 # To use instruction control, uncomment the following lines and replace the model with qwen3-tts-instruct-flash-realtime
 # instructions=&#x27;Speak quickly with a rising intonation, suitable for introducing fashion products.&#x27;,
 # optimize_instructions=True,
 mode = &#x27;commit&#x27; 
 )
 print(f&#x27;send text: {text_to_synthesize[0]}&#x27;)
 qwen_tts_realtime.append_text(text_to_synthesize[0])
 qwen_tts_realtime.commit()
 callback.wait_for_response_done()
 callback.reset_event()
 print(f&#x27;send text: {text_to_synthesize[1]}&#x27;)
 qwen_tts_realtime.append_text(text_to_synthesize[1])
 qwen_tts_realtime.commit()
 callback.wait_for_response_done()
 callback.reset_event()
 print(f&#x27;send text: {text_to_synthesize[2]}&#x27;)
 qwen_tts_realtime.append_text(text_to_synthesize[2])
 qwen_tts_realtime.commit()
 callback.wait_for_response_done()
 qwen_tts_realtime.finish()
 print(&#x27;[Metric] session: {}, first audio delay: {}&#x27;.format(
 qwen_tts_realtime.get_session_id(), 
 qwen_tts_realtime.get_first_audio_delay(),
 ))

``` - Server commit mode 
- Commit mode 

 `appendText()`Copy ```\nimport com.alibaba.dashscope.audio.qwen_tts_realtime.*;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.google.gson.JsonObject;
import javax.sound.sampled.LineUnavailableException;
import javax.sound.sampled.SourceDataLine;
import javax.sound.sampled.AudioFormat;
import javax.sound.sampled.DataLine;
import javax.sound.sampled.AudioSystem;
import java.io.*;
import java.util.Base64;
import java.util.Queue;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.atomic.AtomicReference;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.atomic.AtomicBoolean;
public class Main {
 static String[] textToSynthesize = {
 "Right? I really love this kind of supermarket.",
 "Especially during the Chinese New Year.",
 "Going to the supermarket.",
 "It just makes me feel.",
 "Super, super happy!",
 "I want to buy so many things!"
 };
 public static QwenTtsRealtimeAudioFormat ttsFormat = QwenTtsRealtimeAudioFormat.PCM_24000HZ_MONO_16BIT;
 // Real-time PCM audio player
 public static class RealtimePcmPlayer {
 private int sampleRate;
 private SourceDataLine line;
 private AudioFormat audioFormat;
 private Thread decoderThread;
 private Thread playerThread;
 private AtomicBoolean stopped = new AtomicBoolean(false);
 private Queue<String> b64AudioBuffer = new ConcurrentLinkedQueue<>();
 private Queue<byte[]> RawAudioBuffer = new ConcurrentLinkedQueue<>();
 private ByteArrayOutputStream totalAudioStream = new ByteArrayOutputStream();
 // Initialize the audio format and audio line.
 public RealtimePcmPlayer(int sampleRate) throws LineUnavailableException {
 this.sampleRate = sampleRate;
 this.audioFormat = new AudioFormat(this.sampleRate, 16, 1, true, false);
 DataLine.Info info = new DataLine.Info(SourceDataLine.class, audioFormat);
 line = (SourceDataLine) AudioSystem.getLine(info);
 line.open(audioFormat);
 line.start();
 decoderThread = new Thread(new Runnable() {
 @Override
 public void run() {
 while (!stopped.get()) {
 String b64Audio = b64AudioBuffer.poll();
 if (b64Audio != null) {
 byte[] rawAudio = Base64.getDecoder().decode(b64Audio);
 RawAudioBuffer.add(rawAudio);
 // Write audio data to totalAudioStream.
 try {
 totalAudioStream.write(rawAudio);
 } catch (IOException e) {
 throw new RuntimeException(e);
 }
 } else {
 try {
 Thread.sleep(100);
 } catch (InterruptedException e) {
 throw new RuntimeException(e);
 }
 }
 }
 }
 });
 playerThread = new Thread(new Runnable() {
 @Override
 public void run() {
 while (!stopped.get()) {
 byte[] rawAudio = RawAudioBuffer.poll();
 if (rawAudio != null) {
 try {
 playChunk(rawAudio);
 } catch (IOException e) {
 throw new RuntimeException(e);
 } catch (InterruptedException e) {
 throw new RuntimeException(e);
 }
 } else {
 try {
 Thread.sleep(100);
 } catch (InterruptedException e) {
 throw new RuntimeException(e);
 }
 }
 }
 }
 });
 decoderThread.start();
 playerThread.start();
 }
 // Play an audio chunk and block until playback completes.
 private void playChunk(byte[] chunk) throws IOException, InterruptedException {
 if (chunk == null || chunk.length == 0) return;
 int bytesWritten = 0;
 while (bytesWritten < chunk.length) {
 bytesWritten += line.write(chunk, bytesWritten, chunk.length - bytesWritten);
 }
 int audioLength = chunk.length / (this.sampleRate*2/1000);
 // Wait for the buffered audio to finish playing.
 Thread.sleep(audioLength - 10);
 }
 public void write(String b64Audio) {
 b64AudioBuffer.add(b64Audio);
 }
 public void cancel() {
 b64AudioBuffer.clear();
 RawAudioBuffer.clear();
 }
 public void waitForComplete() throws InterruptedException {
 while (!b64AudioBuffer.isEmpty() || !RawAudioBuffer.isEmpty()) {
 Thread.sleep(100);
 }
 line.drain();
 }
 public void shutdown() throws InterruptedException, IOException {
 stopped.set(true);
 decoderThread.join();
 playerThread.join();
 // Save the complete audio file.
 File file = new File("TotalAudio_"+ttsFormat.getSampleRate()+"."+ttsFormat.getFormat());
 try (FileOutputStream fos = new FileOutputStream(file)) {
 fos.write(totalAudioStream.toByteArray());
 }
 if (line != null && line.isRunning()) {
 line.drain();
 line.close();
 }
 }
 }
 public static void main(String[] args) throws InterruptedException, LineUnavailableException, IOException {
 QwenTtsRealtimeParam param = QwenTtsRealtimeParam.builder()
 // To use instruction control, replace the model with qwen3-tts-instruct-flash-realtime.
 .model("qwen3-tts-flash-realtime")
 .url("wss://dashscope-intl.aliyuncs.com/api-ws/v1/realtime")
 // Obtain an API key from the Qwen Cloud console.
 .apikey(System.getenv("DASHSCOPE_API_KEY"))
 .build();
 AtomicReference<CountDownLatch> completeLatch = new AtomicReference<>(new CountDownLatch(1));
 final AtomicReference<QwenTtsRealtime> qwenTtsRef = new AtomicReference<>(null);
 // Create a real-time audio player instance.
 RealtimePcmPlayer audioPlayer = new RealtimePcmPlayer(24000);
 QwenTtsRealtime qwenTtsRealtime = new QwenTtsRealtime(param, new QwenTtsRealtimeCallback() {
 @Override
 public void onOpen() {
 // Handle connection establishment.
 }
 @Override
 public void onEvent(JsonObject message) {
 String type = message.get("type").getAsString();
 switch(type) {
 case "session.created":
 // Handle session creation.
 if (message.has("session")) {
 String eventId = message.get("event_id").getAsString();
 String sessionId = message.get("session").getAsJsonObject().get("id").getAsString();
 System.out.println("[onEvent] session.created, session_id: "
 + sessionId + ", event_id: " + eventId);
 }
 break;
 case "response.audio.delta":
 String recvAudioB64 = message.get("delta").getAsString();
 // Play audio in real time.
 audioPlayer.write(recvAudioB64);
 break;
 case "response.done":
 // Handle response completion.
 break;
 case "session.finished":
 // Handle session termination.
 completeLatch.get().countDown();
 default:
 break;
 }
 }
 @Override
 public void onClose(int code, String reason) {
 // Handle connection closure.
 }
 });
 qwenTtsRef.set(qwenTtsRealtime);
 try {
 qwenTtsRealtime.connect();
 } catch (NoApiKeyException e) {
 throw new RuntimeException(e);
 }
 QwenTtsRealtimeConfig config = QwenTtsRealtimeConfig.builder()
 .voice("Cherry")
 .responseFormat(ttsFormat)
 .mode("server_commit")
 // To use instruction control, uncomment the following lines and replace the model with qwen3-tts-instruct-flash-realtime.
 // .instructions("")
 // .optimizeInstructions(true)
 .build();
 qwenTtsRealtime.updateSession(config);
 for (String text:textToSynthesize) {
 qwenTtsRealtime.appendText(text);
 Thread.sleep(100);
 }
 qwenTtsRealtime.finish();
 completeLatch.get().await();
 qwenTtsRealtime.close();
 // Wait for audio playback to complete, then shut down the player.
 audioPlayer.waitForComplete();
 audioPlayer.shutdown();
 System.exit(0);
 }
}

``` `commit()`Copy ```\nimport com.alibaba.dashscope.audio.qwen_tts_realtime.*;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.google.gson.JsonObject;
import javax.sound.sampled.LineUnavailableException;
import javax.sound.sampled.SourceDataLine;
import javax.sound.sampled.AudioFormat;
import javax.sound.sampled.DataLine;
import javax.sound.sampled.AudioSystem;
import java.io.*;
import java.util.Base64;
import java.util.Queue;
import java.util.Scanner;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.atomic.AtomicReference;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.atomic.AtomicBoolean;
public class Main {
 public static QwenTtsRealtimeAudioFormat ttsFormat = QwenTtsRealtimeAudioFormat.PCM_24000HZ_MONO_16BIT;
 // Real-time PCM audio player
 public static class RealtimePcmPlayer {
 private int sampleRate;
 private SourceDataLine line;
 private AudioFormat audioFormat;
 private Thread decoderThread;
 private Thread playerThread;
 private AtomicBoolean stopped = new AtomicBoolean(false);
 private Queue<String> b64AudioBuffer = new ConcurrentLinkedQueue<>();
 private Queue<byte[]> RawAudioBuffer = new ConcurrentLinkedQueue<>();
 private ByteArrayOutputStream totalAudioStream = new ByteArrayOutputStream();
 // Initialize the audio format and audio line.
 public RealtimePcmPlayer(int sampleRate) throws LineUnavailableException {
 this.sampleRate = sampleRate;
 this.audioFormat = new AudioFormat(this.sampleRate, 16, 1, true, false);
 DataLine.Info info = new DataLine.Info(SourceDataLine.class, audioFormat);
 line = (SourceDataLine) AudioSystem.getLine(info);
 line.open(audioFormat);
 line.start();
 decoderThread = new Thread(new Runnable() {
 @Override
 public void run() {
 while (!stopped.get()) {
 String b64Audio = b64AudioBuffer.poll();
 if (b64Audio != null) {
 byte[] rawAudio = Base64.getDecoder().decode(b64Audio);
 RawAudioBuffer.add(rawAudio);
 // Write audio data to totalAudioStream.
 try {
 totalAudioStream.write(rawAudio);
 } catch (IOException e) {
 throw new RuntimeException(e);
 }
 } else {
 try {
 Thread.sleep(100);
 } catch (InterruptedException e) {
 throw new RuntimeException(e);
 }
 }
 }
 }
 });
 playerThread = new Thread(new Runnable() {
 @Override
 public void run() {
 while (!stopped.get()) {
 byte[] rawAudio = RawAudioBuffer.poll();
 if (rawAudio != null) {
 try {
 playChunk(rawAudio);
 } catch (IOException e) {
 throw new RuntimeException(e);
 } catch (InterruptedException e) {
 throw new RuntimeException(e);
 }
 } else {
 try {
 Thread.sleep(100);
 } catch (InterruptedException e) {
 throw new RuntimeException(e);
 }
 }
 }
 }
 });
 decoderThread.start();
 playerThread.start();
 }
 // Play an audio chunk and block until playback completes.
 private void playChunk(byte[] chunk) throws IOException, InterruptedException {
 if (chunk == null || chunk.length == 0) return;
 int bytesWritten = 0;
 while (bytesWritten < chunk.length) {
 bytesWritten += line.write(chunk, bytesWritten, chunk.length - bytesWritten);
 }
 int audioLength = chunk.length / (this.sampleRate*2/1000);
 // Wait for the buffered audio to finish playing.
 Thread.sleep(audioLength - 10);
 }
 public void write(String b64Audio) {
 b64AudioBuffer.add(b64Audio);
 }
 public void cancel() {
 b64AudioBuffer.clear();
 RawAudioBuffer.clear();
 }
 public void waitForComplete() throws InterruptedException {
 // Wait for all buffered audio data to finish playing.
 while (!b64AudioBuffer.isEmpty() || !RawAudioBuffer.isEmpty()) {
 Thread.sleep(100);
 }
 // Wait for the audio line to drain.
 line.drain();
 }
 public void shutdown() throws InterruptedException {
 stopped.set(true);
 decoderThread.join();
 playerThread.join();
 // Save the complete audio file.
 File file = new File("TotalAudio_"+ttsFormat.getSampleRate()+"."+ttsFormat.getFormat());
 try (FileOutputStream fos = new FileOutputStream(file)) {
 fos.write(totalAudioStream.toByteArray());
 } catch (FileNotFoundException e) {
 throw new RuntimeException(e);
 } catch (IOException e) {
 throw new RuntimeException(e);
 }
 if (line != null && line.isRunning()) {
 line.drain();
 line.close();
 }
 }
 }
 public static void main(String[] args) throws InterruptedException, LineUnavailableException, FileNotFoundException {
 Scanner scanner = new Scanner(System.in);
 QwenTtsRealtimeParam param = QwenTtsRealtimeParam.builder()
 // To use instruction control, replace the model with qwen3-tts-instruct-flash-realtime.
 .model("qwen3-tts-flash-realtime")
 .url("wss://dashscope-intl.aliyuncs.com/api-ws/v1/realtime")
 // Obtain an API key from the Qwen Cloud console.
 .apikey(System.getenv("DASHSCOPE_API_KEY"))
 .build();
 AtomicReference<CountDownLatch> completeLatch = new AtomicReference<>(new CountDownLatch(1));
 // Create a real-time audio player instance.
 RealtimePcmPlayer audioPlayer = new RealtimePcmPlayer(24000);
 final AtomicReference<QwenTtsRealtime> qwenTtsRef = new AtomicReference<>(null);
 QwenTtsRealtime qwenTtsRealtime = new QwenTtsRealtime(param, new QwenTtsRealtimeCallback() {
 @Override
 public void onOpen() {
 System.out.println("connection opened");
 System.out.println("Enter text and press Enter to send. Enter &#x27;quit&#x27; to exit the program.");
 }
 @Override
 public void onEvent(JsonObject message) {
 String type = message.get("type").getAsString();
 switch(type) {
 case "session.created":
 System.out.println("start session: " + message.get("session").getAsJsonObject().get("id").getAsString());
 break;
 case "response.audio.delta":
 String recvAudioB64 = message.get("delta").getAsString();
 byte[] rawAudio = Base64.getDecoder().decode(recvAudioB64);
 // Play audio in real time.
 audioPlayer.write(recvAudioB64);
 break;
 case "response.done":
 System.out.println("response done");
 // Wait for audio playback to complete.
 try {
 audioPlayer.waitForComplete();
 } catch (InterruptedException e) {
 throw new RuntimeException(e);
 }
 // Prepare for the next input.
 completeLatch.get().countDown();
 break;
 case "session.finished":
 System.out.println("session finished");
 if (qwenTtsRef.get() != null) {
 System.out.println("[Metric] response: " + qwenTtsRef.get().getResponseId() +
 ", first audio delay: " + qwenTtsRef.get().getFirstAudioDelay() + " ms");
 }
 completeLatch.get().countDown();
 default:
 break;
 }
 }
 @Override
 public void onClose(int code, String reason) {
 System.out.println("connection closed code: " + code + ", reason: " + reason);
 try {
 // Wait for playback to complete, then shut down the player.
 audioPlayer.waitForComplete();
 audioPlayer.shutdown();
 } catch (InterruptedException e) {
 throw new RuntimeException(e);
 }
 }
 });
 qwenTtsRef.set(qwenTtsRealtime);
 try {
 qwenTtsRealtime.connect();
 } catch (NoApiKeyException e) {
 throw new RuntimeException(e);
 }
 QwenTtsRealtimeConfig config = QwenTtsRealtimeConfig.builder()
 .voice("Cherry")
 .responseFormat(ttsFormat)
 .mode("commit")
 // To use instruction control, uncomment the following lines and replace the model with qwen3-tts-instruct-flash-realtime.
 // .instructions("")
 // .optimizeInstructions(true)
 .build();
 qwenTtsRealtime.updateSession(config);
 // Read user input in a loop.
 while (true) {
 System.out.print("Enter the text to synthesize: ");
 String text = scanner.nextLine();
 // Exit when the user enters &#x27;quit&#x27;.
 if ("quit".equalsIgnoreCase(text.trim())) {
 System.out.println("Closing the connection...");
 qwenTtsRealtime.finish();
 completeLatch.get().await();
 break;
 }
 // Skip empty input.
 if (text.trim().isEmpty()) {
 continue;
 }
 // Re-initialize the countdown latch.
 completeLatch.set(new CountDownLatch(1));
 // Send the text.
 qwenTtsRealtime.appendText(text);
 qwenTtsRealtime.commit();
 // Wait for the current synthesis to complete.
 completeLatch.get().await();
 }
 // Clean up resources.
 audioPlayer.waitForComplete();
 audioPlayer.shutdown();
 scanner.close();
 System.exit(0);
 }
}

``` 
## [​ ](#advanced-features) Advanced features

### [​ ](#qwen-tts-interaction-modes) Qwen-TTS interaction modes

The Qwen-TTS Realtime API provides two interaction modes:

- **server_commit mode**: The server intelligently handles text segmentation and synthesis timing. This mode suits continuous synthesis of large text blocks. The client only needs to append text without managing segmentation or submission.

- **commit mode**: The client manually submits the text buffer to trigger synthesis. This mode suits scenarios that require precise control over synthesis timing, such as turn-by-turn synthesis in conversational AI.


**Switch the interaction mode**:

- **WebSocket**: Set the `mode` field in the `session.update` event.


Copy ```\n{
 "type": "session.update",
 "session": {
 "mode": "server_commit"
 }
}

``` 

- **Python SDK**: Set the `mode` parameter in the `update_session` method.


Copy ```\nqwen_tts_realtime.update_session(
 voice=&#x27;Cherry&#x27;,
 response_format=AudioFormat.PCM_24000HZ_MONO_16BIT,
 mode=&#x27;server_commit&#x27;
)

``` 

- **Java SDK**: Set the `mode` parameter through `QwenTtsRealtimeConfig.builder()`.


Copy ```\nQwenTtsRealtimeConfig config = QwenTtsRealtimeConfig.builder()
 .voice("Cherry")
 .responseFormat(ttsFormat)
 .mode("server_commit")
 .build();
qwenTtsRealtime.updateSession(config);

``` 
For complete SDK code examples, see [Python SDK](/api-reference/speech-synthesis/cosyvoice/python-sdk) and [Java SDK](/api-reference/speech-synthesis/cosyvoice/java-sdk). For the WebSocket event lifecycle and connection reuse, see the WebSocket API reference.
### [​ ](#instruction-based-control) Instruction-based control

Instruction-based control lets you shape tone, speed, emotion, and timbre through natural language descriptions, without adjusting complex audio parameters.
**Instruction specifications by model**:
- CosyVoice 
- Qwen-TTS 

 **Supported models**: cosyvoice-v3-plus, cosyvoice-v3-flashDifferent models have different instruction format requirements:
- cosyvoice-v3-plus:

Voice Clone/Design timbres: Instruction control isn&#x27;t supported.

- System voices: Instructions must follow a fixed format. For details, see CosyVoice voice list.


- cosyvoice-v3-flash:

Voice Clone/Design timbres: Accept arbitrary instructions.

- System voices: Instructions must follow a fixed format. For details, see CosyVoice voice list.


**How to use**: Specify instruction content through the `instructions` parameter.**Supported languages for instruction text**:
- cosyvoice-v3-plus:

Voice Clone/Design timbres: Chinese, English, French, German, Japanese, Korean, and Russian.

- System voices: Instructions must follow a fixed format. For details, see CosyVoice voice list.


- cosyvoice-v3-flash:

Voice Clone/Design timbres: Chinese, English, French, German, Japanese, Korean, and Russian.

- System voices: Chinese.


**Instruction text length limit**: Up to 100 characters. Chinese characters (including Simplified/Traditional Chinese, Japanese Kanji, and Korean Hanja) count as 2 characters each. Other characters (punctuation, letters, numbers, Japanese Kana, Korean Hangul, etc.) count as 1 character each. **Supported models**: Only the Qwen3-TTS-Instruct-Flash-Realtime series models are supported.**How to use**: Specify instruction content through the `instruction` parameter.**Supported languages for instruction text**: Chinese and English only.**Instruction text length limit**: Up to 1,600 tokens. 
**Use cases**:

- Audiobook and radio drama voiceover

- Advertising and promotional voiceover

- Game character and animation voiceover

- Emotionally expressive voice assistants

- Documentary narration and news broadcasting


**Tips for writing high-quality voice descriptions**:

- 
**Core principles**:

**Be specific, not vague**: Use words that describe concrete vocal qualities, such as "deep," "crisp," or "slightly fast." Avoid subjective or vague terms like "nice" or "normal."

- **Be multidimensional, not single-faceted**: A good description covers multiple dimensions (gender, age, emotion, etc.). Writing only "female voice" is too broad to produce a distinctive timbre.

- **Be objective, not subjective**: Focus on the physical and perceptual qualities of the voice. For example, use "slightly high pitch with energy" rather than "my favorite voice."

- **Be original, not imitative**: Describe the vocal qualities you want, rather than requesting imitation of specific public figures (such as celebrities or actors). The model doesn&#x27;t support imitation, and it may involve copyright risks.

- **Be concise, not redundant**: Make every word count. Avoid repeating synonyms or stacking meaningless modifiers.


- 
**Description dimensions**:
Combining the following dimensions produces more accurate results. The more dimensions described, the more precise the output.
DimensionExample descriptionsGenderMale, female, neutralAgeChild (5-12), teenager (13-18), young adult (19-35), middle-aged (36-55), elderly (55+)PitchHigh, mid, low, slightly high, slightly lowSpeedFast, moderate, slow, slightly fast, slightly slowEmotionCheerful, calm, gentle, serious, lively, composed, soothingTimbreMagnetic, crisp, husky, mellow, sweet, rich, powerfulUse caseNews broadcasting, advertising, audiobook, animation character, voice assistant, documentary narration 


- 
**Examples**:

Standard broadcasting style: Clear and precise articulation with standard pronunciation

- Young, lively female voice with a slightly fast pace and a noticeable rising intonation, suitable for introducing fashion products

- Calm middle-aged male voice with a slow pace, deep and magnetic timbre, suitable for reading news or narrating documentaries

- Gentle, intellectual female voice, around 30 years old, with a calm tone, suitable for audiobook reading

- Cute child voice, about 8-year-old girl, slightly childish speech, suitable for animation character voiceover


### [​ ](#dialects) Dialects

Use the model to output speech in **Chinese dialects** such as Henan, Sichuan, or Cantonese. Configuration varies by model and voice type.
**Dialect setup by model**:
- CosyVoice 
- Qwen-TTS 

 
- **System voices**: Pick one of the following voices from CosyVoice voice list:

Voices with built-in dialect support (for example, `longshange_v3`) output that dialect without any extra configuration.

- Voices that support [instruction-based control](#instruction-based-control) and allow dialect selection (for example, `longanhuan_v3`): specify the target dialect in the instruction text.


- **Voice Clone timbres**: Use [instruction-based control](#instruction-based-control) to set the dialect — for example, set the instruction text to `请用河南话表达` ("Say it in Henan dialect").

- **Voice Design timbres**: Dialects aren&#x27;t supported yet.


**Supported dialects per model**: See the "Supported languages" entry for each model in CosyVoice.**Example**: To produce Henan dialect speech, use the `cosyvoice-v3-flash` model with the `longanhuan_v3` voice and set the instruction text to `"请用河南话表达。"` ("Say it in Henan dialect").Copy ```\n# coding=utf-8
import os
import dashscope
from dashscope.audio.tts_v2 import *
# Obtain an API key from the Qwen Cloud console.
# If the environment variable isn&#x27;t configured, replace the following line with your Qwen Cloud API key: dashscope.api_key = "sk-xxx"
dashscope.api_key = os.environ.get(&#x27;DASHSCOPE_API_KEY&#x27;)
dashscope.base_websocket_api_url=&#x27;wss://dashscope-intl.aliyuncs.com/api-ws/v1/inference&#x27;
# Model
# Different model versions require corresponding voice types:
# cosyvoice-v3-flash/cosyvoice-v3-plus: Use voices such as longanyang.
# Pick a dialect-capable voice
model = "cosyvoice-v3-flash"
# Voice
voice = "longanhuan_v3"
# Instantiate SpeechSynthesizer and pass request parameters such as model, voice, and instruction in the constructor
synthesizer = SpeechSynthesizer(model=model, voice=voice, instruction="请用河南话表达。") # "Say it in Henan dialect"
# Send text for synthesis and obtain binary audio
audio = synthesizer.call("叫你去买盐，你买回来一袋面，这不是弄啥嘞吗！") # Henan dialect sample text
# The first text submission requires establishing a WebSocket connection, so the first-packet latency includes connection setup time
print(&#x27;[Metric] requestId: {}, first package delay: {} ms&#x27;.format(
 synthesizer.get_last_request_id(),
 synthesizer.get_first_package_delay()))
# Save audio to a local file
with open(&#x27;output.mp3&#x27;, &#x27;wb&#x27;) as f:
 f.write(audio)

``` 
- **System voices**: Use a system voice that supports dialects. For the Qwen-TTS voice list, see [Supported voices](#supported-voices).

- **Voice Clone/Design timbres**: Dialects aren&#x27;t supported.


**Supported dialects per model**: See the "Supported languages" entry for each model in Qwen3-TTS. 
### [​ ](#raw-websocket-protocol) Raw WebSocket protocol

The following examples show how to connect directly to the server through the raw WebSocket protocol, for scenarios where the DashScope SDK isn&#x27;t used. These are minimal working implementations. For the WebSocket protocol specification of each model, see the corresponding API reference.
 View raw WebSocket protocol examples

 - CosyVoice 
- Qwen-TTS 

 - Go 
- C# 
- PHP 
- Node.js 
- Java 
- Python 

 Copy ```\npackage main
import (
 "encoding/json"
 "fmt"
 "net/http"
 "os"
 "strings"
 "time"
 "github.com/google/uuid"
 "github.com/gorilla/websocket"
)
const (
 wsURL = "wss://dashscope-intl.aliyuncs.com/api-ws/v1/inference/"
 outputFile = "output.mp3"
)
func main() {
 // Obtain an API key from the Qwen Cloud console.
 // If the environment variable is not configured, replace the following line with your Qwen Cloud API key: apiKey := "sk-xxx"
 apiKey := os.Getenv("DASHSCOPE_API_KEY")
 // Clear the output file
 os.Remove(outputFile)
 os.Create(outputFile)
 // Connect to WebSocket
 header := make(http.Header)
 header.Add("X-DashScope-DataInspection", "enable")
 header.Add("Authorization", fmt.Sprintf("bearer %s", apiKey))
 conn, resp, err := websocket.DefaultDialer.Dial(wsURL, header)
 if err != nil {
 if resp != nil {
 fmt.Printf("Connection failed, HTTP status code: %d\n", resp.StatusCode)
 }
 fmt.Println("Connection failed:", err)
 return
 }
 defer conn.Close()
 // Generate task ID
 taskID := uuid.New().String()
 fmt.Printf("Generated task ID: %s\n", taskID)
 // Send run-task event
 runTaskCmd := map[string]interface{}{
 "header": map[string]interface{}{
 "action": "run-task",
 "task_id": taskID,
 "streaming": "duplex",
 },
 "payload": map[string]interface{}{
 "task_group": "audio",
 "task": "tts",
 "function": "SpeechSynthesizer",
 "model": "cosyvoice-v3-flash",
 "parameters": map[string]interface{}{
 "text_type": "PlainText",
 "voice": "longanyang",
 "format": "mp3",
 "sample_rate": 22050,
 "volume": 50,
 "rate": 1,
 "pitch": 1,
 // If enable_ssml is set to true, only one continue-task event is allowed; otherwise the error "Text request limit violated, expected 1." will occur.
 "enable_ssml": false,
 },
 "input": map[string]interface{}{},
 },
 }
 runTaskJSON, _ := json.Marshal(runTaskCmd)
 fmt.Printf("Sending run-task event: %s\n", string(runTaskJSON))
 err = conn.WriteMessage(websocket.TextMessage, runTaskJSON)
 if err != nil {
 fmt.Println("Failed to send run-task:", err)
 return
 }
 textSent := false
 // Process messages
 for {
 messageType, message, err := conn.ReadMessage()
 if err != nil {
 fmt.Println("Failed to read message:", err)
 break
 }
 // Process binary messages
 if messageType == websocket.BinaryMessage {
 fmt.Printf("Received binary message, length: %d\n", len(message))
 file, _ := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
 file.Write(message)
 file.Close()
 continue
 }
 // Process text messages
 messageStr := string(message)
 fmt.Printf("Received text message: %s\n", strings.ReplaceAll(messageStr, "\n", ""))
 // Simple JSON parsing to get event type
 var msgMap map[string]interface{}
 if json.Unmarshal(message, &msgMap) == nil {
 if header, ok := msgMap["header"].(map[string]interface{}); ok {
 if event, ok := header["event"].(string); ok {
 fmt.Printf("Event type: %s\n", event)
 switch event {
 case "task-started":
 fmt.Println("=== Received task-started event ===")
 if !textSent {
 // Send continue-task events
 texts := []string{"Before my bed, moonlight shines bright, I suspect it&#x27;s frost upon the ground.", "I raise my eyes to gaze at the bright moon, then bow my head, thinking of home."}
 for _, text := range texts {
 continueTaskCmd := map[string]interface{}{
 "header": map[string]interface{}{
 "action": "continue-task",
 "task_id": taskID,
 "streaming": "duplex",
 },
 "payload": map[string]interface{}{
 "input": map[string]interface{}{
 "text": text,
 },
 },
 }
 continueTaskJSON, _ := json.Marshal(continueTaskCmd)
 fmt.Printf("Sending continue-task event: %s\n", string(continueTaskJSON))
 err = conn.WriteMessage(websocket.TextMessage, continueTaskJSON)
 if err != nil {
 fmt.Println("Failed to send continue-task:", err)
 return
 }
 }
 textSent = true
 // Delay before sending finish-task
 time.Sleep(500 * time.Millisecond)
 // Send finish-task event
 finishTaskCmd := map[string]interface{}{
 "header": map[string]interface{}{
 "action": "finish-task",
 "task_id": taskID,
 "streaming": "duplex",
 },
 "payload": map[string]interface{}{
 "input": map[string]interface{}{},
 },
 }
 finishTaskJSON, _ := json.Marshal(finishTaskCmd)
 fmt.Printf("Sending finish-task event: %s\n", string(finishTaskJSON))
 err = conn.WriteMessage(websocket.TextMessage, finishTaskJSON)
 if err != nil {
 fmt.Println("Failed to send finish-task:", err)
 return
 }
 }
 case "task-finished":
 fmt.Println("=== Task completed ===")
 return
 case "task-failed":
 fmt.Println("=== Task failed ===")
 if header["error_message"] != nil {
 fmt.Printf("Error message: %s\n", header["error_message"])
 }
 return
 case "result-generated":
 fmt.Println("Received result-generated event")
 }
 }
 }
 }
 }
}

``` Copy ```\nusing System.Net.WebSockets;
using System.Text;
using System.Text.Json;
class Program {
 // Obtain an API key from the Qwen Cloud console.
 // If the environment variable is not configured, replace the following line with your Qwen Cloud API key: private static readonly string ApiKey = "sk-xxx"
 private static readonly string ApiKey = Environment.GetEnvironmentVariable("DASHSCOPE_API_KEY") ?? throw new InvalidOperationException("DASHSCOPE_API_KEY environment variable is not set.");
 private const string WebSocketUrl = "wss://dashscope-intl.aliyuncs.com/api-ws/v1/inference/";
 // Output file path
 private const string OutputFilePath = "output.mp3";
 // WebSocket client
 private static ClientWebSocket _webSocket = new ClientWebSocket();
 // Cancellation token source
 private static CancellationTokenSource _cancellationTokenSource = new CancellationTokenSource();
 // Task ID
 private static string? _taskId;
 // Whether the task has started
 private static TaskCompletionSource<bool> _taskStartedTcs = new TaskCompletionSource<bool>();
 static async Task Main(string[] args) {
 try {
 // Clear the output file
 ClearOutputFile(OutputFilePath);
 // Connect to the WebSocket service
 await ConnectToWebSocketAsync(WebSocketUrl);
 // Start the message receiving task
 Task receiveTask = ReceiveMessagesAsync();
 // Send run-task event
 _taskId = GenerateTaskId();
 await SendRunTaskCommandAsync(_taskId);
 // Wait for the task-started event
 await _taskStartedTcs.Task;
 // Continuously send continue-task events
 string[] texts = {
 "Before my bed, moonlight shines bright",
 "I suspect it&#x27;s frost upon the ground",
 "I raise my eyes to gaze at the bright moon",
 "then bow my head, thinking of home"
 };
 foreach (string text in texts) {
 await SendContinueTaskCommandAsync(text);
 }
 // Send finish-task event
 await SendFinishTaskCommandAsync(_taskId);
 // Wait for the receiving task to complete
 await receiveTask;
 Console.WriteLine("Task completed, connection closed.");
 } catch (OperationCanceledException) {
 Console.WriteLine("Task cancelled.");
 } catch (Exception ex) {
 Console.WriteLine($"Error occurred: {ex.Message}");
 } finally {
 _cancellationTokenSource.Cancel();
 _webSocket.Dispose();
 }
 }
 private static void ClearOutputFile(string filePath) {
 if (File.Exists(filePath)) {
 File.WriteAllText(filePath, string.Empty);
 Console.WriteLine("Output file cleared.");
 } else {
 Console.WriteLine("Output file does not exist, no need to clear.");
 }
 }
 private static async Task ConnectToWebSocketAsync(string url) {
 var uri = new Uri(url);
 if (_webSocket.State == WebSocketState.Connecting || _webSocket.State == WebSocketState.Open) {
 return;
 }
 // Set WebSocket connection headers
 _webSocket.Options.SetRequestHeader("Authorization", $"bearer {ApiKey}");
 _webSocket.Options.SetRequestHeader("X-DashScope-DataInspection", "enable");
 try {
 await _webSocket.ConnectAsync(uri, _cancellationTokenSource.Token);
 Console.WriteLine("Successfully connected to the WebSocket service.");
 } catch (OperationCanceledException) {
 Console.WriteLine("WebSocket connection cancelled.");
 } catch (Exception ex) {
 Console.WriteLine($"WebSocket connection failed: {ex.Message}");
 throw;
 }
 }
 private static async Task SendRunTaskCommandAsync(string taskId) {
 var command = CreateCommand("run-task", taskId, "duplex", new {
 task_group = "audio",
 task = "tts",
 function = "SpeechSynthesizer",
 model = "cosyvoice-v3-flash",
 parameters = new
 {
 text_type = "PlainText",
 voice = "longanyang",
 format = "mp3",
 sample_rate = 22050,
 volume = 50,
 rate = 1,
 pitch = 1,
 // If enable_ssml is set to true, only one continue-task event is allowed; otherwise the error "Text request limit violated, expected 1." will occur.
 enable_ssml = false
 },
 input = new { }
 });
 await SendJsonMessageAsync(command);
 Console.WriteLine("run-task event sent.");
 }
 private static async Task SendContinueTaskCommandAsync(string text) {
 if (_taskId == null) {
 throw new InvalidOperationException("Task ID not initialized.");
 }
 var command = CreateCommand("continue-task", _taskId, "duplex", new {
 input = new {
 text
 }
 });
 await SendJsonMessageAsync(command);
 Console.WriteLine("continue-task event sent.");
 }
 private static async Task SendFinishTaskCommandAsync(string taskId) {
 var command = CreateCommand("finish-task", taskId, "duplex", new {
 input = new { }
 });
 await SendJsonMessageAsync(command);
 Console.WriteLine("finish-task event sent.");
 }
 private static async Task SendJsonMessageAsync(string message) {
 var buffer = Encoding.UTF8.GetBytes(message);
 try {
 await _webSocket.SendAsync(new ArraySegment<byte>(buffer), WebSocketMessageType.Text, true, _cancellationTokenSource.Token);
 } catch (OperationCanceledException) {
 Console.WriteLine("Message sending cancelled.");
 }
 }
 private static async Task ReceiveMessagesAsync() {
 while (_webSocket.State == WebSocketState.Open) {
 var response = await ReceiveMessageAsync();
 if (response != null) {
 var eventStr = response.RootElement.GetProperty("header").GetProperty("event").GetString();
 switch (eventStr) {
 case "task-started":
 Console.WriteLine("Task started.");
 _taskStartedTcs.TrySetResult(true);
 break;
 case "task-finished":
 Console.WriteLine("Task completed.");
 _cancellationTokenSource.Cancel();
 break;
 case "task-failed":
 Console.WriteLine("Task failed: " + response.RootElement.GetProperty("header").GetProperty("error_message").GetString());
 _cancellationTokenSource.Cancel();
 break;
 default:
 // result-generated can be handled here
 break;
 }
 }
 }
 }
 private static async Task<JsonDocument?> ReceiveMessageAsync() {
 var buffer = new byte[1024 * 4];
 var segment = new ArraySegment<byte>(buffer);
 try {
 WebSocketReceiveResult result = await _webSocket.ReceiveAsync(segment, _cancellationTokenSource.Token);
 if (result.MessageType == WebSocketMessageType.Close) {
 await _webSocket.CloseAsync(WebSocketCloseStatus.NormalClosure, "Closing", _cancellationTokenSource.Token);
 return null;
 }
 if (result.MessageType == WebSocketMessageType.Binary) {
 // Process binary data
 Console.WriteLine("Received binary data...");
 // Save binary data to file
 using (var fileStream = new FileStream(OutputFilePath, FileMode.Append)) {
 fileStream.Write(buffer, 0, result.Count);
 }
 return null;
 }
 string message = Encoding.UTF8.GetString(buffer, 0, result.Count);
 return JsonDocument.Parse(message);
 } catch (OperationCanceledException) {
 Console.WriteLine("Message receiving cancelled.");
 return null;
 }
 }
 private static string GenerateTaskId() {
 return Guid.NewGuid().ToString("N").Substring(0, 32);
 }
 private static string CreateCommand(string action, string taskId, string streaming, object payload) {
 var command = new {
 header = new {
 action,
 task_id = taskId,
 streaming
 },
 payload
 };
 return JsonSerializer.Serialize(command);
 }
}

``` The example code directory structure:Copy ```\nmy-php-project/
├── composer.json
├── vendor/
└── index.php

``` The composer.json contents are as follows. Determine the appropriate dependency versions based on your requirements:Copy ```\n{
 "require": {
 "react/event-loop": "^1.3",
 "react/socket": "^1.11",
 "react/stream": "^1.2",
 "react/http": "^1.1",
 "ratchet/pawl": "^0.4"
 },
 "autoload": {
 "psr-4": {
 "App\\": "src/"
 }
 }
}

``` The index.php contents:Copy ```\n<?php
require __DIR__ . &#x27;/vendor/autoload.php&#x27;;
use Ratchet\Client\Connector;
use React\EventLoop\Loop;
use React\Socket\Connector as SocketConnector;
// Obtain an API key from the Qwen Cloud console.
// If the environment variable is not configured, replace the following line with your Qwen Cloud API key: $api_key = "sk-xxx"
$api_key = getenv("DASHSCOPE_API_KEY");
$websocket_url = &#x27;wss://dashscope-intl.aliyuncs.com/api-ws/v1/inference/&#x27;; // WebSocket server address
$output_file = &#x27;output.mp3&#x27;; // Output file path
$loop = Loop::get();
if (file_exists($output_file)) {
 // Clear file content
 file_put_contents($output_file, &#x27;&#x27;);
}
// Create a custom connector
$socketConnector = new SocketConnector($loop, [
 &#x27;tcp&#x27; => [
 &#x27;bindto&#x27; => &#x27;0.0.0.0:0&#x27;,
 ],
 &#x27;tls&#x27; => [
 &#x27;verify_peer&#x27; => false,
 &#x27;verify_peer_name&#x27; => false,
 ],
]);
$connector = new Connector($loop, $socketConnector);
$headers = [
 &#x27;Authorization&#x27; => &#x27;bearer &#x27; . $api_key,
 &#x27;X-DashScope-DataInspection&#x27; => &#x27;enable&#x27;
];
$connector($websocket_url, [], $headers)->then(function ($conn) use ($loop, $output_file) {
 echo "Connected to WebSocket server\n";
 // Generate task ID
 $taskId = generateTaskId();
 // Send run-task event
 sendRunTaskMessage($conn, $taskId);
 // Define the function for sending continue-task events
 $sendContinueTask = function() use ($conn, $loop, $taskId) {
 // Text to send
 $texts = ["Before my bed, moonlight shines bright", "I suspect it&#x27;s frost upon the ground", "I raise my eyes to gaze at the bright moon", "then bow my head, thinking of home"];
 $continueTaskCount = 0;
 foreach ($texts as $text) {
 $continueTaskMessage = json_encode([
 "header" => [
 "action" => "continue-task",
 "task_id" => $taskId,
 "streaming" => "duplex"
 ],
 "payload" => [
 "input" => [
 "text" => $text
 ]
 ]
 ]);
 echo "Sending continue-task event: " . $continueTaskMessage . "\n";
 $conn->send($continueTaskMessage);
 $continueTaskCount++;
 }
 echo "Number of continue-task events sent: " . $continueTaskCount . "\n";
 // Send finish-task event
 sendFinishTaskMessage($conn, $taskId);
 };
 // Flag for whether the task-started event has been received
 $taskStarted = false;
 // Listen for messages
 $conn->on(&#x27;message&#x27;, function($msg) use ($conn, $sendContinueTask, $loop, &$taskStarted, $taskId, $output_file) {
 if ($msg->isBinary()) {
 // Write binary data to local file
 file_put_contents($output_file, $msg->getPayload(), FILE_APPEND);
 } else {
 // Process non-binary messages
 $response = json_decode($msg, true);
 if (isset($response[&#x27;header&#x27;][&#x27;event&#x27;])) {
 handleEvent($conn, $response, $sendContinueTask, $loop, $taskId, $taskStarted);
 } else {
 echo "Unknown message format\n";
 }
 }
 });
 // Listen for connection close
 $conn->on(&#x27;close&#x27;, function($code = null, $reason = null) {
 echo "Connection closed\n";
 if ($code !== null) {
 echo "Close code: " . $code . "\n";
 }
 if ($reason !== null) {
 echo "Close reason: " . $reason . "\n";
 }
 });
}, function ($e) {
 echo "Unable to connect: {$e->getMessage()}\n";
});
$loop->run();
/**
 * Generate task ID
 * @return string
 */
function generateTaskId(): string {
 return bin2hex(random_bytes(16));
}
/**
 * Send run-task event
 * @param $conn
 * @param $taskId
 */
function sendRunTaskMessage($conn, $taskId) {
 $runTaskMessage = json_encode([
 "header" => [
 "action" => "run-task",
 "task_id" => $taskId,
 "streaming" => "duplex"
 ],
 "payload" => [
 "task_group" => "audio",
 "task" => "tts",
 "function" => "SpeechSynthesizer",
 "model" => "cosyvoice-v3-flash",
 "parameters" => [
 "text_type" => "PlainText",
 "voice" => "longanyang",
 "format" => "mp3",
 "sample_rate" => 22050,
 "volume" => 50,
 "rate" => 1,
 "pitch" => 1,
 // If enable_ssml is set to true, only one continue-task event is allowed; otherwise the error "Text request limit violated, expected 1." will occur.
 "enable_ssml" => false
 ],
 "input" => (object) []
 ]
 ]);
 echo "Sending run-task event: " . $runTaskMessage . "\n";
 $conn->send($runTaskMessage);
 echo "run-task event sent\n";
}
/**
 * Read audio file
 * @param string $filePath
 * @return bool|string
 */
function readAudioFile(string $filePath) {
 $voiceData = file_get_contents($filePath);
 if ($voiceData === false) {
 echo "Unable to read audio file\n";
 }
 return $voiceData;
}
/**
 * Split audio data
 * @param string $data
 * @param int $chunkSize
 * @return array
 */
function splitAudioData(string $data, int $chunkSize): array {
 return str_split($data, $chunkSize);
}
/**
 * Send finish-task event
 * @param $conn
 * @param $taskId
 */
function sendFinishTaskMessage($conn, $taskId) {
 $finishTaskMessage = json_encode([
 "header" => [
 "action" => "finish-task",
 "task_id" => $taskId,
 "streaming" => "duplex"
 ],
 "payload" => [
 "input" => (object) []
 ]
 ]);
 echo "Sending finish-task event: " . $finishTaskMessage . "\n";
 $conn->send($finishTaskMessage);
 echo "finish-task event sent\n";
}
/**
 * Handle events
 * @param $conn
 * @param $response
 * @param $sendContinueTask
 * @param $loop
 * @param $taskId
 * @param $taskStarted
 */
function handleEvent($conn, $response, $sendContinueTask, $loop, $taskId, &$taskStarted) {
 switch ($response[&#x27;header&#x27;][&#x27;event&#x27;]) {
 case &#x27;task-started&#x27;:
 echo "Task started, sending continue-task events...\n";
 $taskStarted = true;
 // Send continue-task events
 $sendContinueTask();
 break;
 case &#x27;result-generated&#x27;:
 // Received result-generated event
 break;
 case &#x27;task-finished&#x27;:
 echo "Task completed\n";
 $conn->close();
 break;
 case &#x27;task-failed&#x27;:
 echo "Task failed\n";
 echo "Error code: " . $response[&#x27;header&#x27;][&#x27;error_code&#x27;] . "\n";
 echo "Error message: " . $response[&#x27;header&#x27;][&#x27;error_message&#x27;] . "\n";
 $conn->close();
 break;
 case &#x27;error&#x27;:
 echo "Error: " . $response[&#x27;payload&#x27;][&#x27;message&#x27;] . "\n";
 break;
 default:
 echo "Unknown event: " . $response[&#x27;header&#x27;][&#x27;event&#x27;] . "\n";
 break;
 }
 // If the task is completed, close the connection
 if ($response[&#x27;header&#x27;][&#x27;event&#x27;] == &#x27;task-finished&#x27;) {
 // Wait 1 second to ensure all data has been transmitted
 $loop->addTimer(1, function() use ($conn) {
 $conn->close();
 echo "Client closed connection\n";
 });
 }
 // If the task-started event was not received, close the connection
 if (!$taskStarted && in_array($response[&#x27;header&#x27;][&#x27;event&#x27;], [&#x27;task-failed&#x27;, &#x27;error&#x27;])) {
 $conn->close();
 }
}

``` Install the required dependencies:Copy ```\nnpm install ws
npm install uuid

``` Example code:Copy ```\nconst WebSocket = require(&#x27;ws&#x27;);
const fs = require(&#x27;fs&#x27;);
const uuid = require(&#x27;uuid&#x27;).v4;
// Obtain an API key from the Qwen Cloud console.
// If the environment variable is not configured, replace the following line with your Qwen Cloud API key: const apiKey = "sk-xxx"
const apiKey = process.env.DASHSCOPE_API_KEY;
const url = &#x27;wss://dashscope-intl.aliyuncs.com/api-ws/v1/inference/&#x27;;
// Output file path
const outputFilePath = &#x27;output.mp3&#x27;;
// Clear the output file
fs.writeFileSync(outputFilePath, &#x27;&#x27;);
// Create WebSocket client
const ws = new WebSocket(url, {
 headers: {
 Authorization: `bearer ${apiKey}`,
 &#x27;X-DashScope-DataInspection&#x27;: &#x27;enable&#x27;
 }
});
let taskStarted = false;
let taskId = uuid();
ws.on(&#x27;open&#x27;, () => {
 console.log(&#x27;Connected to WebSocket server&#x27;);
 // Send run-task event
 const runTaskMessage = JSON.stringify({
 header: {
 action: &#x27;run-task&#x27;,
 task_id: taskId,
 streaming: &#x27;duplex&#x27;
 },
 payload: {
 task_group: &#x27;audio&#x27;,
 task: &#x27;tts&#x27;,
 function: &#x27;SpeechSynthesizer&#x27;,
 model: &#x27;cosyvoice-v3-flash&#x27;,
 parameters: {
 text_type: &#x27;PlainText&#x27;,
 voice: &#x27;longanyang&#x27;, // Voice
 format: &#x27;mp3&#x27;, // Audio format
 sample_rate: 22050, // Sample rate
 volume: 50, // Volume
 rate: 1, // Speech rate
 pitch: 1, // Pitch
 enable_ssml: false // Whether to enable SSML. If enable_ssml is set to true, only one continue-task event is allowed; otherwise the error "Text request limit violated, expected 1." will occur.
 },
 input: {}
 }
 });
 ws.send(runTaskMessage);
 console.log(&#x27;run-task message sent&#x27;);
});
const fileStream = fs.createWriteStream(outputFilePath, { flags: &#x27;a&#x27; });
ws.on(&#x27;message&#x27;, (data, isBinary) => {
 if (isBinary) {
 // Write binary data to file
 fileStream.write(data);
 } else {
 const message = JSON.parse(data);
 switch (message.header.event) {
 case &#x27;task-started&#x27;:
 taskStarted = true;
 console.log(&#x27;Task started&#x27;);
 // Send continue-task events
 sendContinueTasks(ws);
 break;
 case &#x27;task-finished&#x27;:
 console.log(&#x27;Task completed&#x27;);
 ws.close();
 fileStream.end(() => {
 console.log(&#x27;File stream closed&#x27;);
 });
 break;
 case &#x27;task-failed&#x27;:
 console.error(&#x27;Task failed:&#x27;, message.header.error_message);
 ws.close();
 fileStream.end(() => {
 console.log(&#x27;File stream closed&#x27;);
 });
 break;
 default:
 // result-generated can be handled here
 break;
 }
 }
});
function sendContinueTasks(ws) {
 const texts = [
 &#x27;Before my bed, moonlight shines bright,&#x27;,
 &#x27;I suspect it\&#x27;s frost upon the ground.&#x27;,
 &#x27;I raise my eyes to gaze at the bright moon,&#x27;,
 &#x27;then bow my head, thinking of home.&#x27;
 ];
 texts.forEach((text, index) => {
 setTimeout(() => {
 if (taskStarted) {
 const continueTaskMessage = JSON.stringify({
 header: {
 action: &#x27;continue-task&#x27;,
 task_id: taskId,
 streaming: &#x27;duplex&#x27;
 },
 payload: {
 input: {
 text: text
 }
 }
 });
 ws.send(continueTaskMessage);
 console.log(`Sent continue-task, text: ${text}`);
 }
 }, index * 1000); // Send one every second
 });
 // Send finish-task event
 setTimeout(() => {
 if (taskStarted) {
 const finishTaskMessage = JSON.stringify({
 header: {
 action: &#x27;finish-task&#x27;,
 task_id: taskId,
 streaming: &#x27;duplex&#x27;
 },
 payload: {
 input: {}
 }
 });
 ws.send(finishTaskMessage);
 console.log(&#x27;finish-task sent&#x27;);
 }
 }, texts.length * 1000 + 1000); // Send 1 second after all continue-task events are sent
}
ws.on(&#x27;close&#x27;, () => {
 console.log(&#x27;Disconnected from WebSocket server&#x27;);
});

``` For Java, we recommend using the Java DashScope SDK. For details, see [Java SDK](/api-reference/speech-synthesis/cosyvoice/java-sdk).The following is a Java WebSocket example. Before running it, make sure you&#x27;ve imported these dependencies:
- `Java-WebSocket`

- `jackson-databind`


Use Maven or Gradle to manage dependencies. The configuration is as follows:- pom.xml 
- build.gradle 

 Copy ```\n<dependencies>
 <!-- WebSocket Client -->
 <dependency>
 <groupId>org.java-websocket</groupId>
 <artifactId>Java-WebSocket</artifactId>
 <version>1.5.3</version>
 </dependency>
 <!-- JSON Processing -->
 <dependency>
 <groupId>com.fasterxml.jackson.core</groupId>
 <artifactId>jackson-databind</artifactId>
 <version>2.13.0</version>
 </dependency>
</dependencies>

``` Copy ```\n// Other code omitted
dependencies {
 // WebSocket Client
 implementation &#x27;org.java-websocket:Java-WebSocket:1.5.3&#x27;
 // JSON Processing
 implementation &#x27;com.fasterxml.jackson.core:jackson-databind:2.13.0&#x27;
}
// Other code omitted

``` Java code:Copy ```\nimport com.fasterxml.jackson.databind.ObjectMapper;
import org.java_websocket.client.WebSocketClient;
import org.java_websocket.handshake.ServerHandshake;
import java.io.FileOutputStream;
import java.io.IOException;
import java.net.URI;
import java.nio.ByteBuffer;
import java.util.*;
public class TTSWebSocketClient extends WebSocketClient {
 private final String taskId = UUID.randomUUID().toString();
 private final String outputFile = "output_" + System.currentTimeMillis() + ".mp3";
 private boolean taskFinished = false;
 public TTSWebSocketClient(URI serverUri, Map<String, String> headers) {
 super(serverUri, headers);
 }
 @Override
 public void onOpen(ServerHandshake serverHandshake) {
 System.out.println("Connection established");
 // Send run-task event
 // If enable_ssml is set to true, only one continue-task event is allowed; otherwise the error "Text request limit violated, expected 1." will occur.
 String runTaskCommand = "{ \"header\": { \"action\": \"run-task\", \"task_id\": \"" + taskId + "\", \"streaming\": \"duplex\" }, \"payload\": { \"task_group\": \"audio\", \"task\": \"tts\", \"function\": \"SpeechSynthesizer\", \"model\": \"cosyvoice-v3-flash\", \"parameters\": { \"text_type\": \"PlainText\", \"voice\": \"longanyang\", \"format\": \"mp3\", \"sample_rate\": 22050, \"volume\": 50, \"rate\": 1, \"pitch\": 1, \"enable_ssml\": false }, \"input\": {} }}";
 send(runTaskCommand);
 }
 @Override
 public void onMessage(String message) {
 System.out.println("Received message from server: " + message);
 try {
 // Parse JSON message
 Map<String, Object> messageMap = new ObjectMapper().readValue(message, Map.class);
 if (messageMap.containsKey("header")) {
 Map<String, Object> header = (Map<String, Object>) messageMap.get("header");
 if (header.containsKey("event")) {
 String event = (String) header.get("event");
 if ("task-started".equals(event)) {
 System.out.println("Received task-started event from server");
 List<String> texts = Arrays.asList(
 "Before my bed, moonlight shines bright, I suspect it&#x27;s frost upon the ground",
 "I raise my eyes to gaze at the bright moon, then bow my head, thinking of home"
 );
 for (String text : texts) {
 // Send continue-task event
 sendContinueTask(text);
 }
 // Send finish-task event
 sendFinishTask();
 } else if ("task-finished".equals(event)) {
 System.out.println("Received task-finished event from server");
 taskFinished = true;
 closeConnection();
 } else if ("task-failed".equals(event)) {
 System.out.println("Task failed: " + message);
 closeConnection();
 }
 }
 }
 } catch (Exception e) {
 System.err.println("Exception occurred: " + e.getMessage());
 }
 }
 @Override
 public void onMessage(ByteBuffer message) {
 System.out.println("Received binary audio data, size: " + message.remaining());
 try (FileOutputStream fos = new FileOutputStream(outputFile, true)) {
 byte[] buffer = new byte[message.remaining()];
 message.get(buffer);
 fos.write(buffer);
 System.out.println("Audio data written to local file " + outputFile);
 } catch (IOException e) {
 System.err.println("Failed to write audio data to local file: " + e.getMessage());
 }
 }
 @Override
 public void onClose(int code, String reason, boolean remote) {
 System.out.println("Connection closed: " + reason + " (" + code + ")");
 }
 @Override
 public void onError(Exception ex) {
 System.err.println("Error: " + ex.getMessage());
 ex.printStackTrace();
 }
 private void sendContinueTask(String text) {
 String command = "{ \"header\": { \"action\": \"continue-task\", \"task_id\": \"" + taskId + "\", \"streaming\": \"duplex\" }, \"payload\": { \"input\": { \"text\": \"" + text + "\" } }}";
 send(command);
 }
 private void sendFinishTask() {
 String command = "{ \"header\": { \"action\": \"finish-task\", \"task_id\": \"" + taskId + "\", \"streaming\": \"duplex\" }, \"payload\": { \"input\": {} }}";
 send(command);
 }
 private void closeConnection() {
 if (!isClosed()) {
 close();
 }
 }
 public static void main(String[] args) {
 try {
 // Obtain an API key from the Qwen Cloud console.
 // If the environment variable is not configured, replace the following line with your Qwen Cloud API key: String apiKey = "sk-xxx"
 String apiKey = System.getenv("DASHSCOPE_API_KEY");
 if (apiKey == null || apiKey.isEmpty()) {
 System.err.println("Please set the DASHSCOPE_API_KEY environment variable");
 return;
 }
 Map<String, String> headers = new HashMap<>();
 headers.put("Authorization", "bearer " + apiKey);
 TTSWebSocketClient client = new TTSWebSocketClient(new URI("wss://dashscope-intl.aliyuncs.com/api-ws/v1/inference/"), headers);
 client.connect();
 while (!client.isClosed() && !client.taskFinished) {
 Thread.sleep(1000);
 }
 } catch (Exception e) {
 System.err.println("Failed to connect to WebSocket service: " + e.getMessage());
 e.printStackTrace();
 }
 }
}

``` For Python, we recommend using the Python DashScope SDK. For details, see [Python SDK](/api-reference/speech-synthesis/cosyvoice/python-sdk).The following is a Python WebSocket example. Before running it, install the dependency as follows:Copy ```\npip uninstall websocket-client
pip uninstall websocket
pip install websocket-client

``` Don&#x27;t name your Python file "websocket.py". Otherwise, an error occurs (AttributeError: module &#x27;websocket&#x27; has no attribute &#x27;WebSocketApp&#x27;. Did you mean: &#x27;WebSocket&#x27;?). Copy ```\nimport websocket
import json
import uuid
import os
import time
class TTSClient:
 def __init__(self, api_key, uri):
 """
 Initialize a TTSClient instance.
 Parameters:
 api_key (str): API Key for authentication
 uri (str): WebSocket service address
 """
 self.api_key = api_key # Replace with your API Key
 self.uri = uri # Replace with your WebSocket address
 self.task_id = str(uuid.uuid4()) # Generate a unique task ID
 self.output_file = f"output_{int(time.time())}.mp3" # Output audio file path
 self.ws = None # WebSocketApp instance
 self.task_started = False # Whether task-started has been received
 self.task_finished = False # Whether task-finished / task-failed has been received
 def on_open(self, ws):
 """
 Callback when WebSocket connection is established.
 Sends the run-task event to start the speech synthesis task.
 """
 print("WebSocket connected")
 # Construct the run-task event
 run_task_cmd = {
 "header": {
 "action": "run-task",
 "task_id": self.task_id,
 "streaming": "duplex"
 },
 "payload": {
 "task_group": "audio",
 "task": "tts",
 "function": "SpeechSynthesizer",
 "model": "cosyvoice-v3-flash",
 "parameters": {
 "text_type": "PlainText",
 "voice": "longanyang",
 "format": "mp3",
 "sample_rate": 22050,
 "volume": 50,
 "rate": 1,
 "pitch": 1,
 # If enable_ssml is set to True, only one continue-task event is allowed; otherwise the error "Text request limit violated, expected 1." will occur.
 "enable_ssml": False
 },
 "input": {}
 }
 }
 # Send the run-task event
 ws.send(json.dumps(run_task_cmd))
 print("run-task event sent")
 def on_message(self, ws, message):
 """
 Callback when a message is received.
 Handles text and binary messages separately.
 """
 if isinstance(message, str):
 # Process JSON text messages
 try:
 msg_json = json.loads(message)
 print(f"Received JSON message: {msg_json}")
 if "header" in msg_json:
 header = msg_json["header"]
 if "event" in header:
 event = header["event"]
 if event == "task-started":
 print("Task started")
 self.task_started = True
 # Send continue-task events
 texts = [
 "Before my bed, moonlight shines bright, I suspect it&#x27;s frost upon the ground",
 "I raise my eyes to gaze at the bright moon, then bow my head, thinking of home"
 ]
 for text in texts:
 self.send_continue_task(text)
 # Send finish-task after all continue-task events are sent
 self.send_finish_task()
 elif event == "task-finished":
 print("Task completed")
 self.task_finished = True
 self.close(ws)
 elif event == "task-failed":
 error_msg = msg_json.get("error_message", "Unknown error")
 print(f"Task failed: {error_msg}")
 self.task_finished = True
 self.close(ws)
 except json.JSONDecodeError as e:
 print(f"JSON parsing failed: {e}")
 else:
 # Process binary messages (audio data)
 print(f"Received binary message, size: {len(message)} bytes")
 with open(self.output_file, "ab") as f:
 f.write(message)
 print(f"Audio data written to local file {self.output_file}")
 def on_error(self, ws, error):
 """Callback when an error occurs"""
 print(f"WebSocket error: {error}")
 def on_close(self, ws, close_status_code, close_msg):
 """Callback when connection is closed"""
 print(f"WebSocket closed: {close_msg} ({close_status_code})")
 def send_continue_task(self, text):
 """Send a continue-task event with the text content to synthesize"""
 cmd = {
 "header": {
 "action": "continue-task",
 "task_id": self.task_id,
 "streaming": "duplex"
 },
 "payload": {
 "input": {
 "text": text
 }
 }
 }
 self.ws.send(json.dumps(cmd))
 print(f"Sent continue-task event, text content: {text}")
 def send_finish_task(self):
 """Send a finish-task event to end the speech synthesis task"""
 cmd = {
 "header": {
 "action": "finish-task",
 "task_id": self.task_id,
 "streaming": "duplex"
 },
 "payload": {
 "input": {}
 }
 }
 self.ws.send(json.dumps(cmd))
 print("finish-task event sent")
 def close(self, ws):
 """Actively close the connection"""
 if ws and ws.sock and ws.sock.connected:
 ws.close()
 print("Connection closed actively")
 def run(self):
 """Start the WebSocket client"""
 # Set request headers (authentication)
 header = {
 "Authorization": f"bearer {self.api_key}",
 "X-DashScope-DataInspection": "enable"
 }
 # Create a WebSocketApp instance
 self.ws = websocket.WebSocketApp(
 self.uri,
 header=header,
 on_open=self.on_open,
 on_message=self.on_message,
 on_error=self.on_error,
 on_close=self.on_close
 )
 print("Listening for WebSocket messages...")
 self.ws.run_forever() # Start long-lived connection listener
# Example usage
if __name__ == "__main__":
 # Obtain an API key from the Qwen Cloud console.
 # If the environment variable is not configured, replace the following line with your Qwen Cloud API key: API_KEY = "sk-xxx"
 API_KEY = os.environ.get("DASHSCOPE_API_KEY")
 SERVER_URI = "wss://dashscope-intl.aliyuncs.com/api-ws/v1/inference/" # Replace with your WebSocket address
 client = TTSClient(API_KEY, SERVER_URI)
 client.run()

``` 
- **Create the client**


- Python 
- Java 

 Create a Python file named `tts_realtime_client.py` and copy the following code into it:Copy ```\n# -- coding: utf-8 --
import asyncio
import websockets
import json
import base64
import time
from typing import Optional, Callable, Dict, Any
from enum import Enum
class SessionMode(Enum):
 SERVER_COMMIT = "server_commit"
 COMMIT = "commit"
class TTSRealtimeClient:
 """
 Client for interacting with the TTS Realtime API.
 This class provides methods for connecting to the TTS Realtime API, sending text data,
 receiving audio output, and managing WebSocket connections.
 Attributes:
 base_url (str):
 Base URL of the Realtime API.
 api_key (str):
 API Key for authentication.
 voice (str):
 Voice used by the server for speech synthesis.
 mode (SessionMode):
 Session mode, either server_commit or commit.
 audio_callback (Callable[[bytes], None]):
 Callback function for receiving audio data.
 language_type(str)
 Language of the synthesized speech. Options: Chinese, English, German, Italian, Portuguese, Spanish, Japanese, Korean, French, Russian, Auto
 """
 def __init__(
 self,
 base_url: str,
 api_key: str,
 voice: str = "Cherry",
 mode: SessionMode = SessionMode.SERVER_COMMIT,
 audio_callback: Optional[Callable[[bytes], None]] = None,
 language_type: str = "Auto"):
 self.base_url = base_url
 self.api_key = api_key
 self.voice = voice
 self.mode = mode
 self.ws = None
 self.audio_callback = audio_callback
 self.language_type = language_type
 # Current response state
 self._current_response_id = None
 self._current_item_id = None
 self._is_responding = False
 self._response_done_future = None
 async def connect(self) -> None:
 """Establish a WebSocket connection with the TTS Realtime API."""
 headers = {
 "Authorization": f"Bearer {self.api_key}"
 }
 self.ws = await websockets.connect(self.base_url, additional_headers=headers)
 # Set default session configuration
 await self.update_session({
 "mode": self.mode.value,
 "voice": self.voice,
 # To use the instruction control feature, uncomment the following lines and replace the model with qwen3-tts-instruct-flash-realtime in server_commit.py or commit.py
 # "instructions": "Speak quickly with a noticeable rising intonation, suitable for introducing fashion products.",
 # "optimize_instructions": true
 "language_type": self.language_type,
 "response_format": "pcm",
 "sample_rate": 24000
 })
 async def send_event(self, event) -> None:
 """Send an event to the server."""
 event[&#x27;event_id&#x27;] = "event_" + str(int(time.time() * 1000))
 print(f"Sending event: type={event[&#x27;type&#x27;]}, event_id={event[&#x27;event_id&#x27;]}")
 await self.ws.send(json.dumps(event))
 async def update_session(self, config: Dict[str, Any]) -> None:
 """Update session configuration."""
 event = {
 "type": "session.update",
 "session": config
 }
 print("Updating session configuration: ", event)
 await self.send_event(event)
 async def append_text(self, text: str) -> None:
 """Send text data to the API."""
 event = {
 "type": "input_text_buffer.append",
 "text": text
 }
 await self.send_event(event)
 async def commit_text_buffer(self) -> None:
 """Commit the text buffer to trigger processing."""
 event = {
 "type": "input_text_buffer.commit"
 }
 await self.send_event(event)
 async def clear_text_buffer(self) -> None:
 """Clear the text buffer."""
 event = {
 "type": "input_text_buffer.clear"
 }
 await self.send_event(event)
 async def finish_session(self) -> None:
 """End the session."""
 event = {
 "type": "session.finish"
 }
 await self.send_event(event)
 async def wait_for_response_done(self):
 """Wait for the response.done event"""
 if self._response_done_future:
 await self._response_done_future
 async def handle_messages(self) -> None:
 """Process messages from the server."""
 try:
 async for message in self.ws:
 event = json.loads(message)
 event_type = event.get("type")
 if event_type != "response.audio.delta":
 print(f"Received event: {event_type}")
 if event_type == "error":
 print("Error: ", event.get(&#x27;error&#x27;, {}))
 continue
 elif event_type == "session.created":
 print("Session created, ID: ", event.get(&#x27;session&#x27;, {}).get(&#x27;id&#x27;))
 elif event_type == "session.updated":
 print("Session updated, ID: ", event.get(&#x27;session&#x27;, {}).get(&#x27;id&#x27;))
 elif event_type == "input_text_buffer.committed":
 print("Text buffer committed, item ID: ", event.get(&#x27;item_id&#x27;))
 elif event_type == "input_text_buffer.cleared":
 print("Text buffer cleared")
 elif event_type == "response.created":
 self._current_response_id = event.get("response", {}).get("id")
 self._is_responding = True
 # Create a new future to wait for response.done
 self._response_done_future = asyncio.Future()
 print("Response created, ID: ", self._current_response_id)
 elif event_type == "response.output_item.added":
 self._current_item_id = event.get("item", {}).get("id")
 print("Output item added, ID: ", self._current_item_id)
 # Process audio delta
 elif event_type == "response.audio.delta" and self.audio_callback:
 audio_bytes = base64.b64decode(event.get("delta", ""))
 self.audio_callback(audio_bytes)
 elif event_type == "response.audio.done":
 print("Audio generation completed")
 elif event_type == "response.done":
 self._is_responding = False
 self._current_response_id = None
 self._current_item_id = None
 # Mark the future as done
 if self._response_done_future and not self._response_done_future.done():
 self._response_done_future.set_result(True)
 print("Response completed")
 elif event_type == "session.finished":
 print("Session ended")
 except websockets.exceptions.ConnectionClosed:
 print("Connection closed")
 except Exception as e:
 print("Error processing messages: ", str(e))
 async def close(self) -> None:
 """Close the WebSocket connection."""
 if self.ws:
 await self.ws.close()

``` Create a Java file named `TTSRealtimeClient.java` and copy the following code into it:Copy ```\nimport com.google.gson.Gson;
import com.google.gson.JsonObject;
import org.java_websocket.client.WebSocketClient;
import org.java_websocket.handshake.ServerHandshake;
import java.net.URI;
import java.util.Base64;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.CountDownLatch;
import java.util.function.Consumer;
/**
 * Client for interacting with the TTS Realtime API.
 *
 * This class provides methods for connecting to the TTS Realtime API, sending text data,
 * receiving audio output, and managing WebSocket connections.
 */
public class TTSRealtimeClient {
 public enum SessionMode {
 SERVER_COMMIT("server_commit"),
 COMMIT("commit");
 private final String value;
 SessionMode(String value) { this.value = value; }
 public String getValue() { return value; }
 }
 /**
 * Audio callback interface
 */
 public interface AudioCallback {
 void onAudio(byte[] audioData);
 }
 private final String baseUrl;
 private final String apiKey;
 private final String voice;
 private final SessionMode mode;
 private final String languageType;
 private final AudioCallback audioCallback;
 private final Gson gson = new Gson();
 private WebSocketClient ws;
 private CountDownLatch responseDoneLatch;
 private CountDownLatch sessionFinishedLatch;
 public TTSRealtimeClient(String baseUrl, String apiKey, String voice,
 SessionMode mode, AudioCallback audioCallback,
 String languageType) {
 this.baseUrl = baseUrl;
 this.apiKey = apiKey;
 this.voice = voice;
 this.mode = mode;
 this.audioCallback = audioCallback;
 this.languageType = languageType;
 }
 public TTSRealtimeClient(String baseUrl, String apiKey, String voice,
 SessionMode mode, AudioCallback audioCallback) {
 this(baseUrl, apiKey, voice, mode, audioCallback, "Auto");
 }
 /**
 * Establish a WebSocket connection with the TTS Realtime API.
 */
 public void connect() throws Exception {
 Map<String, String> headers = new HashMap<>();
 headers.put("Authorization", "Bearer " + apiKey);
 responseDoneLatch = new CountDownLatch(0);
 sessionFinishedLatch = new CountDownLatch(1);
 ws = new WebSocketClient(new URI(baseUrl), headers) {
 @Override
 public void onOpen(ServerHandshake handshake) {
 System.out.println("WebSocket connection established");
 // Send default session configuration
 JsonObject session = new JsonObject();
 session.addProperty("mode", mode.getValue());
 session.addProperty("voice", TTSRealtimeClient.this.voice);
 // To use the instruction control feature, uncomment the following lines and replace the model with qwen3-tts-instruct-flash-realtime
 // session.addProperty("instructions", "Speak quickly with a noticeable rising intonation, suitable for introducing fashion products.");
 // session.addProperty("optimize_instructions", true);
 session.addProperty("language_type", languageType);
 session.addProperty("response_format", "pcm");
 session.addProperty("sample_rate", 24000);
 updateSession(session);
 }
 @Override
 public void onMessage(String message) {
 JsonObject event = gson.fromJson(message, JsonObject.class);
 String eventType = event.has("type") ? event.get("type").getAsString() : "";
 if (!"response.audio.delta".equals(eventType)) {
 System.out.println("Received event: " + eventType);
 }
 switch (eventType) {
 case "error":
 System.err.println("Error: " + event.get("error"));
 break;
 case "session.created":
 System.out.println("Session created, ID: " +
 event.getAsJsonObject("session").get("id").getAsString());
 break;
 case "session.updated":
 System.out.println("Session updated, ID: " +
 event.getAsJsonObject("session").get("id").getAsString());
 break;
 case "input_text_buffer.committed":
 System.out.println("Text buffer committed, item ID: " + event.get("item_id"));
 break;
 case "input_text_buffer.cleared":
 System.out.println("Text buffer cleared");
 break;
 case "response.created":
 System.out.println("Response created, ID: " +
 event.getAsJsonObject("response").get("id").getAsString());
 responseDoneLatch = new CountDownLatch(1);
 break;
 case "response.output_item.added":
 System.out.println("Output item added, ID: " +
 event.getAsJsonObject("item").get("id").getAsString());
 break;
 case "response.audio.delta":
 if (audioCallback != null) {
 byte[] audioBytes = Base64.getDecoder().decode(
 event.get("delta").getAsString());
 audioCallback.onAudio(audioBytes);
 }
 break;
 case "response.audio.done":
 System.out.println("Audio generation completed");
 break;
 case "response.done":
 System.out.println("Response completed");
 responseDoneLatch.countDown();
 break;
 case "session.finished":
 System.out.println("Session ended");
 sessionFinishedLatch.countDown();
 break;
 }
 }
 @Override
 public void onClose(int code, String reason, boolean remote) {
 System.out.println("Connection closed: " + reason);
 }
 @Override
 public void onError(Exception ex) {
 System.err.println("WebSocket error: " + ex.getMessage());
 }
 };
 ws.connectBlocking();
 }
 /**
 * Send an event to the server.
 */
 public void sendEvent(JsonObject event) {
 String eventId = "event_" + System.currentTimeMillis();
 event.addProperty("event_id", eventId);
 System.out.println("Sending event: type=" + event.get("type").getAsString()
 + ", event_id=" + eventId);
 ws.send(gson.toJson(event));
 }
 /**
 * Update session configuration.
 */
 public void updateSession(JsonObject config) {
 JsonObject event = new JsonObject();
 event.addProperty("type", "session.update");
 event.add("session", config);
 System.out.println("Updating session configuration: " + event);
 sendEvent(event);
 }
 /**
 * Send text data to the API.
 */
 public void appendText(String text) {
 JsonObject event = new JsonObject();
 event.addProperty("type", "input_text_buffer.append");
 event.addProperty("text", text);
 sendEvent(event);
 }
 /**
 * Commit the text buffer to trigger processing.
 */
 public void commitTextBuffer() {
 JsonObject event = new JsonObject();
 event.addProperty("type", "input_text_buffer.commit");
 sendEvent(event);
 }
 /**
 * Clear the text buffer.
 */
 public void clearTextBuffer() {
 JsonObject event = new JsonObject();
 event.addProperty("type", "input_text_buffer.clear");
 sendEvent(event);
 }
 /**
 * End the session.
 */
 public void finishSession() {
 JsonObject event = new JsonObject();
 event.addProperty("type", "session.finish");
 sendEvent(event);
 }
 /**
 * Wait for the response.done event.
 */
 public void waitForResponseDone() throws InterruptedException {
 responseDoneLatch.await();
 }
 /**
 * Wait for the session.finished event.
 */
 public void waitForSessionFinished() throws InterruptedException {
 sessionFinishedLatch.await();
 }
 /**
 * Close the WebSocket connection.
 */
 public void close() {
 if (ws != null) {
 ws.close();
 }
 }
}

``` 
- **Choose a synthesis mode**


The Realtime API supports the following two modes:
- 
**server_commit mode**
The client sends only text. The server intelligently determines segmentation and synthesis timing. This mode suits low-latency scenarios that don&#x27;t require manual control over synthesis rhythm, such as GPS navigation.


- 
**commit mode**
The client adds text to a buffer and then triggers the server to synthesize the specified text. This mode suits scenarios that require fine-grained control over sentence breaks and pauses, such as news broadcasting.


- server_commit mode 
- commit mode 

 - Python 
- Java 

 In the same directory as `tts_realtime_client.py`, create another Python file named `server_commit.py` and copy the following code into it:Copy ```\nimport os
import asyncio
import logging
import wave
from tts_realtime_client import TTSRealtimeClient, SessionMode
import pyaudio
# QwenTTS service configuration
# To use the instruction control feature, replace the model with qwen3-tts-instruct-flash-realtime and uncomment the instructions in tts_realtime_client.py
URL = "wss://dashscope-intl.aliyuncs.com/api-ws/v1/realtime?model=qwen3-tts-flash-realtime"
# Obtain an API key from the Qwen Cloud console.
# If the environment variable is not configured, replace the following line with your Qwen Cloud API key: API_KEY="sk-xxx"
API_KEY = os.getenv("DASHSCOPE_API_KEY")
if not API_KEY:
 raise ValueError("Please set DASHSCOPE_API_KEY environment variable")
# Collect audio data
_audio_chunks = []
# Real-time playback related
_AUDIO_SAMPLE_RATE = 24000
_audio_pyaudio = pyaudio.PyAudio()
_audio_stream = None # Will be opened at runtime
def _audio_callback(audio_bytes: bytes):
 """TTSRealtimeClient audio callback: real-time playback and caching"""
 global _audio_stream
 if _audio_stream is not None:
 try:
 _audio_stream.write(audio_bytes)
 except Exception as exc:
 logging.error(f"PyAudio playback error: {exc}")
 _audio_chunks.append(audio_bytes)
 logging.info(f"Received audio chunk: {len(audio_bytes)} bytes")
def _save_audio_to_file(filename: str = "output.wav", sample_rate: int = 24000) -> bool:
 """Save collected audio data as a WAV file"""
 if not _audio_chunks:
 logging.warning("No audio data to save")
 return False
 try:
 audio_data = b"".join(_audio_chunks)
 with wave.open(filename, &#x27;wb&#x27;) as wav_file:
 wav_file.setnchannels(1) # Mono
 wav_file.setsampwidth(2) # 16-bit
 wav_file.setframerate(sample_rate)
 wav_file.writeframes(audio_data)
 logging.info(f"Audio saved to: {filename}")
 return True
 except Exception as exc:
 logging.error(f"Failed to save audio: {exc}")
 return False
async def _produce_text(client: TTSRealtimeClient):
 """Send text fragments to the server"""
 text_fragments = [
 "Qwen Cloud is an all-in-one platform for developing and building large language model applications.",
 "Both developers and business users can deeply participate in the design and development of large language model applications.",
 "You can develop a large language model application in five minutes using a simple interface,",
 "or train a dedicated model in a few hours, allowing you to focus more energy on application innovation."
 ]
 logging.info("Sending text fragments…")
 for text in text_fragments:
 logging.info(f"Sending fragment: {text}")
 await client.append_text(text)
 await asyncio.sleep(0.1) # Brief delay between fragments
 # Wait for server to complete internal processing before ending the session
 await asyncio.sleep(1.0)
 await client.finish_session()
async def _run_demo():
 """Run the complete demo"""
 global _audio_stream
 # Open PyAudio output stream
 _audio_stream = _audio_pyaudio.open(
 format=pyaudio.paInt16,
 channels=1,
 rate=_AUDIO_SAMPLE_RATE,
 output=True,
 frames_per_buffer=1024
 )
 client = TTSRealtimeClient(
 base_url=URL,
 api_key=API_KEY,
 voice="Cherry",
 mode=SessionMode.SERVER_COMMIT,
 audio_callback=_audio_callback
 )
 # Establish connection
 await client.connect()
 # Run message handling and text sending in parallel
 consumer_task = asyncio.create_task(client.handle_messages())
 producer_task = asyncio.create_task(_produce_text(client))
 await producer_task # Wait for text sending to complete
 # Wait for response.done
 await client.wait_for_response_done()
 # Close connection and cancel consumer task
 await client.close()
 consumer_task.cancel()
 # Close audio stream
 if _audio_stream is not None:
 _audio_stream.stop_stream()
 _audio_stream.close()
 _audio_pyaudio.terminate()
 # Save audio data
 os.makedirs("outputs", exist_ok=True)
 _save_audio_to_file(os.path.join("outputs", "qwen_tts_output.wav"))
def main():
 """Synchronous entry point"""
 logging.basicConfig(
 level=logging.INFO,
 format=&#x27;%(asctime)s [%(levelname)s] %(message)s&#x27;,
 datefmt=&#x27;%Y-%m-%d %H:%M:%S&#x27;
 )
 logging.info("Starting QwenTTS Realtime Client demo…")
 asyncio.run(_run_demo())
if __name__ == "__main__":
 main()

``` Run `server_commit.py` to hear the audio generated by the Realtime API in real time. In the same directory as `TTSRealtimeClient.java`, create another Java file named `ServerCommit.java` and copy the following code into it:Copy ```\nimport javax.sound.sampled.*;
import java.io.*;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.atomic.AtomicBoolean;
public class ServerCommit {
 private static final String URL = "wss://dashscope-intl.aliyuncs.com/api-ws/v1/realtime?model=qwen3-tts-flash-realtime";
 // Obtain an API key from the Qwen Cloud console.
 // If the environment variable is not configured, replace the following line with your Qwen Cloud API key: private static final String API_KEY = "sk-xxx";
 private static final String API_KEY = System.getenv("DASHSCOPE_API_KEY");
 private static final int SAMPLE_RATE = 24000;
 // Audio data cache
 private static final List<byte[]> audioChunks = new ArrayList<>();
 // Real-time playback queue
 private static final ConcurrentLinkedQueue<byte[]> playbackQueue = new ConcurrentLinkedQueue<>();
 private static final AtomicBoolean playing = new AtomicBoolean(true);
 public static void main(String[] args) throws Exception {
 if (API_KEY == null || API_KEY.isEmpty()) {
 throw new IllegalStateException("Please set the DASHSCOPE_API_KEY environment variable");
 }
 // Initialize audio playback
 AudioFormat format = new AudioFormat(SAMPLE_RATE, 16, 1, true, false);
 DataLine.Info info = new DataLine.Info(SourceDataLine.class, format);
 SourceDataLine audioLine = (SourceDataLine) AudioSystem.getLine(info);
 audioLine.open(format);
 audioLine.start();
 // Start playback thread
 Thread playerThread = new Thread(() -> {
 while (playing.get() || !playbackQueue.isEmpty()) {
 byte[] chunk = playbackQueue.poll();
 if (chunk != null) {
 audioLine.write(chunk, 0, chunk.length);
 } else {
 try { Thread.sleep(10); } catch (InterruptedException ignored) {}
 }
 }
 });
 playerThread.start();
 // Create TTS client
 // To use the instruction control feature, replace the model with qwen3-tts-instruct-flash-realtime and uncomment the instructions in TTSRealtimeClient.java
 TTSRealtimeClient client = new TTSRealtimeClient(
 URL, API_KEY, "Cherry",
 TTSRealtimeClient.SessionMode.SERVER_COMMIT,
 audioData -> {
 playbackQueue.add(audioData);
 audioChunks.add(audioData);
 System.out.println("Received audio data: " + audioData.length + " bytes");
 }
 );
 client.connect();
 // Send text fragments
 String[] textFragments = {
 "Qwen Cloud is an all-in-one platform for developing and building large language model applications.",
 "Both developers and business users can deeply participate in the design and development of large language model applications.",
 "You can develop a large language model application in five minutes using a simple interface,",
 "or train a dedicated model in a few hours, allowing you to focus more energy on application innovation."
 };
 System.out.println("Sending text...");
 for (String text : textFragments) {
 System.out.println("Sending fragment: " + text);
 client.appendText(text);
 Thread.sleep(100);
 }
 Thread.sleep(1000);
 client.finishSession();
 // Wait for response to complete
 client.waitForResponseDone();
 client.waitForSessionFinished();
 client.close();
 // Wait for playback to complete
 playing.set(false);
 playerThread.join();
 audioLine.drain();
 audioLine.close();
 // Save audio file
 saveWav("output.wav");
 System.out.println("Done");
 }
 private static void saveWav(String filename) throws IOException {
 if (audioChunks.isEmpty()) {
 System.out.println("No audio data to save");
 return;
 }
 ByteArrayOutputStream bos = new ByteArrayOutputStream();
 for (byte[] chunk : audioChunks) {
 bos.write(chunk);
 }
 byte[] allAudio = bos.toByteArray();
 AudioFormat format = new AudioFormat(SAMPLE_RATE, 16, 1, true, false);
 AudioInputStream ais = new AudioInputStream(
 new ByteArrayInputStream(allAudio), format, allAudio.length / 2);
 new File("outputs").mkdirs();
 AudioSystem.write(ais, AudioFileFormat.Type.WAVE,
 new File("outputs/" + filename));
 System.out.println("Audio saved to: outputs/" + filename);
 }
}

``` Compile and run `ServerCommit.java` to hear the audio generated by the Realtime API in real time. - Python 
- Java 

 In the same directory as `tts_realtime_client.py`, create another Python file named `commit.py` and copy the following code into it:Copy ```\nimport os
import asyncio
import logging
import wave
from tts_realtime_client import TTSRealtimeClient, SessionMode
import pyaudio
# QwenTTS service configuration
# To use the instruction control feature, replace the model with qwen3-tts-instruct-flash-realtime and uncomment the instructions in tts_realtime_client.py
URL = "wss://dashscope-intl.aliyuncs.com/api-ws/v1/realtime?model=qwen3-tts-flash-realtime"
# Obtain an API key from the Qwen Cloud console.
# If the environment variable is not configured, replace the following line with your Qwen Cloud API key: API_KEY="sk-xxx"
API_KEY = os.getenv("DASHSCOPE_API_KEY")
if not API_KEY:
 raise ValueError("Please set DASHSCOPE_API_KEY environment variable")
# Collect audio data
_audio_chunks = []
_AUDIO_SAMPLE_RATE = 24000
_audio_pyaudio = pyaudio.PyAudio()
_audio_stream = None
def _audio_callback(audio_bytes: bytes):
 """TTSRealtimeClient audio callback: real-time playback and caching"""
 global _audio_stream
 if _audio_stream is not None:
 try:
 _audio_stream.write(audio_bytes)
 except Exception as exc:
 logging.error(f"PyAudio playback error: {exc}")
 _audio_chunks.append(audio_bytes)
 logging.info(f"Received audio chunk: {len(audio_bytes)} bytes")
def _save_audio_to_file(filename: str = "output.wav", sample_rate: int = 24000) -> bool:
 """Save collected audio data as a WAV file"""
 if not _audio_chunks:
 logging.warning("No audio data to save")
 return False
 try:
 audio_data = b"".join(_audio_chunks)
 with wave.open(filename, &#x27;wb&#x27;) as wav_file:
 wav_file.setnchannels(1) # Mono
 wav_file.setsampwidth(2) # 16-bit
 wav_file.setframerate(sample_rate)
 wav_file.writeframes(audio_data)
 logging.info(f"Audio saved to: {filename}")
 return True
 except Exception as exc:
 logging.error(f"Failed to save audio: {exc}")
 return False
async def _user_input_loop(client: TTSRealtimeClient):
 """Continuously get user input and send text. When user enters empty text, send a commit event and end the session."""
 print("Enter text (press Enter directly to send a commit event and end the session, press Ctrl+C or Ctrl+D to exit the program):")
 while True:
 try:
 user_text = input("> ")
 if not user_text: # User input is empty
 # Empty input marks the end of a conversation: commit buffer -> end session -> exit loop
 logging.info("Empty input, sending commit event and ending session")
 await client.commit_text_buffer()
 # Wait briefly for server to process commit, preventing premature session end that could lose audio
 await asyncio.sleep(0.3)
 await client.finish_session()
 break # Exit user input loop directly, no need to press Enter again
 else:
 logging.info(f"Sending text: {user_text}")
 await client.append_text(user_text)
 except EOFError: # User pressed Ctrl+D
 break
 except KeyboardInterrupt: # User pressed Ctrl+C
 break
 # End session
 logging.info("Ending session...")
async def _run_demo():
 """Run the complete demo"""
 global _audio_stream
 # Open PyAudio output stream
 _audio_stream = _audio_pyaudio.open(
 format=pyaudio.paInt16,
 channels=1,
 rate=_AUDIO_SAMPLE_RATE,
 output=True,
 frames_per_buffer=1024
 )
 client = TTSRealtimeClient(
 base_url=URL,
 api_key=API_KEY,
 voice="Cherry",
 mode=SessionMode.COMMIT, # Changed to COMMIT mode
 audio_callback=_audio_callback
 )
 # Establish connection
 await client.connect()
 # Run message handling and user input in parallel
 consumer_task = asyncio.create_task(client.handle_messages())
 producer_task = asyncio.create_task(_user_input_loop(client))
 await producer_task # Wait for user input to complete
 # Wait for response.done
 await client.wait_for_response_done()
 # Close connection and cancel consumer task
 await client.close()
 consumer_task.cancel()
 # Close audio stream
 if _audio_stream is not None:
 _audio_stream.stop_stream()
 _audio_stream.close()
 _audio_pyaudio.terminate()
 # Save audio data
 os.makedirs("outputs", exist_ok=True)
 _save_audio_to_file(os.path.join("outputs", "qwen_tts_output.wav"))
def main():
 logging.basicConfig(
 level=logging.INFO,
 format=&#x27;%(asctime)s [%(levelname)s] %(message)s&#x27;,
 datefmt=&#x27;%Y-%m-%d %H:%M:%S&#x27;
 )
 logging.info("Starting QwenTTS Realtime Client demo…")
 asyncio.run(_run_demo())
if __name__ == "__main__":
 main()

``` Run `commit.py`. You can enter text multiple times. Press Enter without entering text to hear the audio returned by the Realtime API through your speakers. In the same directory as `TTSRealtimeClient.java`, create another Java file named `Commit.java` and copy the following code into it:Copy ```\nimport javax.sound.sampled.*;
import java.io.*;
import java.util.ArrayList;
import java.util.List;
import java.util.Scanner;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.atomic.AtomicBoolean;
public class Commit {
 private static final String URL = "wss://dashscope-intl.aliyuncs.com/api-ws/v1/realtime?model=qwen3-tts-flash-realtime";
 // Obtain an API key from the Qwen Cloud console.
 // If the environment variable is not configured, replace the following line with your Qwen Cloud API key: private static final String API_KEY = "sk-xxx";
 private static final String API_KEY = System.getenv("DASHSCOPE_API_KEY");
 private static final int SAMPLE_RATE = 24000;
 private static final List<byte[]> audioChunks = new ArrayList<>();
 private static final ConcurrentLinkedQueue<byte[]> playbackQueue = new ConcurrentLinkedQueue<>();
 private static final AtomicBoolean playing = new AtomicBoolean(true);
 public static void main(String[] args) throws Exception {
 if (API_KEY == null || API_KEY.isEmpty()) {
 throw new IllegalStateException("Please set the DASHSCOPE_API_KEY environment variable");
 }
 // Initialize audio playback
 AudioFormat format = new AudioFormat(SAMPLE_RATE, 16, 1, true, false);
 DataLine.Info info = new DataLine.Info(SourceDataLine.class, format);
 SourceDataLine audioLine = (SourceDataLine) AudioSystem.getLine(info);
 audioLine.open(format);
 audioLine.start();
 // Start playback thread
 Thread playerThread = new Thread(() -> {
 while (playing.get() || !playbackQueue.isEmpty()) {
 byte[] chunk = playbackQueue.poll();
 if (chunk != null) {
 audioLine.write(chunk, 0, chunk.length);
 } else {
 try { Thread.sleep(10); } catch (InterruptedException ignored) {}
 }
 }
 });
 playerThread.start();
 // Create TTS client (commit mode)
 // To use the instruction control feature, replace the model with qwen3-tts-instruct-flash-realtime and uncomment the instructions in TTSRealtimeClient.java
 TTSRealtimeClient client = new TTSRealtimeClient(
 URL, API_KEY, "Cherry",
 TTSRealtimeClient.SessionMode.COMMIT,
 audioData -> {
 playbackQueue.add(audioData);
 audioChunks.add(audioData);
 System.out.println("Received audio data: " + audioData.length + " bytes");
 }
 );
 client.connect();
 // Interactive input
 System.out.println("Enter text (press Enter directly to send a commit event and end the session, press Ctrl+D to exit the program):");
 Scanner scanner = new Scanner(System.in);
 while (true) {
 System.out.print("> ");
 if (!scanner.hasNextLine()) {
 client.finishSession();
 break;
 }
 String userText = scanner.nextLine();
 if (userText.isEmpty()) {
 // Empty input: commit buffer and end session
 System.out.println("Empty input, sending commit event and ending session");
 client.commitTextBuffer();
 Thread.sleep(300);
 client.finishSession();
 break;
 } else {
 System.out.println("Sending text: " + userText);
 client.appendText(userText);
 }
 }
 scanner.close();
 // Wait for response to complete
 client.waitForResponseDone();
 client.waitForSessionFinished();
 client.close();
 // Wait for playback to complete
 playing.set(false);
 playerThread.join();
 audioLine.drain();
 audioLine.close();
 // Save audio file
 saveWav("output.wav");
 System.out.println("Done");
 }
 private static void saveWav(String filename) throws IOException {
 if (audioChunks.isEmpty()) {
 System.out.println("No audio data to save");
 return;
 }
 ByteArrayOutputStream bos = new ByteArrayOutputStream();
 for (byte[] chunk : audioChunks) {
 bos.write(chunk);
 }
 byte[] allAudio = bos.toByteArray();
 AudioFormat format = new AudioFormat(SAMPLE_RATE, 16, 1, true, false);
 AudioInputStream ais = new AudioInputStream(
 new ByteArrayInputStream(allAudio), format, allAudio.length / 2);
 new File("outputs").mkdirs();
 AudioSystem.write(ais, AudioFileFormat.Type.WAVE,
 new File("outputs/" + filename));
 System.out.println("Audio saved to: outputs/" + filename);
 }
}

``` Compile and run `Commit.java`. You can enter text multiple times. Press Enter without entering text to hear the audio returned by the Realtime API through your speakers. 
## [​ ](#voice-customization) Voice customization

- CosyVoice 
- Qwen-TTS-Realtime 

 ### Voice cloning: input audio format requirements
High-quality input audio is the foundation for excellent cloning results.ItemRequirement**Supported formats**WAV (16-bit), MP3, M4A**Audio duration**Recommended: 10-20 seconds. Maximum: 60 seconds.**File size**≤ 10 MB**Sample rate**≥ 16 kHz**Channels**Mono or stereo. For stereo audio, only the first channel is processed, so make sure it contains clear speech.**Content**The audio must contain at least 5 seconds of continuous, clear speech with no background sound. The rest may contain only brief pauses (≤ 2 seconds). The entire clip should be free of background music, noise, or other voices to ensure high-quality core speech content. Use audio of normal speech as the input; do not upload songs or singing, to ensure accurate and usable cloning results. ### Voice design: writing high-quality voice descriptions
#### Constraints
When writing a voice description (`voice_prompt`), follow these technical constraints:
- **Length limit**: The `voice_prompt` content must not exceed 500 characters.

- **Supported languages**: The description text supports Chinese and English only.


#### Core principles
The `voice_prompt` guides the model to generate a voice with specific characteristics.When writing a voice description, follow these core principles:
- **Be specific, not vague**: Use words that describe concrete vocal qualities, such as "deep," "crisp," or "slightly fast." Avoid subjective, low-information terms like "nice" or "normal."

- **Be multidimensional, not single-faceted**: A good description usually combines multiple dimensions (such as gender, age, and emotion). A single-dimension description (such as just "female voice") is too broad to produce a distinctive result.

- **Be objective, not subjective**: Focus on the physical and perceptual qualities of the voice itself, rather than personal preferences. For example, use "slightly high pitch with energy" instead of "my favorite voice."

- **Be original, not imitative**: Describe the qualities of the voice, rather than requesting imitation of specific people (such as celebrities or actors). Such requests involve copyright risks, and the model doesn&#x27;t support direct imitation.

- **Be concise, not redundant**: Make sure every word is meaningful. Avoid repeating synonyms or meaningless intensifiers (such as "a very, very nice voice").


#### Description dimension reference
DimensionExamplesGenderMale, female, neutralAgeChild (5-12), teenager (13-18), young adult (19-35), middle-aged (36-55), elderly (55+)PitchHigh, mid, low, slightly high, slightly lowSpeedFast, moderate, slow, slightly fast, slightly slowEmotionCheerful, composed, gentle, serious, lively, cool, soothingTimbreMagnetic, crisp, husky, mellow, sweet, rich, powerfulUse caseNews broadcasting, advertising voiceover, audiobook, animation character, voice assistant, documentary narration #### Example comparison
**Good examples**:
- "A young, lively female voice with a fast pace and a noticeable rising intonation, suitable for introducing fashion products."

*Analysis*: This description combines age, personality, speed, and intonation, and specifies a use case, forming a clear voice profile.


- "A composed middle-aged male voice with a slightly slow pace, deep and magnetic, suitable for news broadcasting or documentary narration."

*Analysis*: This description clearly defines gender, age range, speed, timbre, and use case.


- "A cute child voice, about an 8-year-old girl, speaking with a slightly childish tone, suitable for animation character voiceover."

*Analysis*: This description precisely targets age and a vocal quality (childishness), with a clear use case.


- "A gentle, intellectual female, around 30 years old, with a calm tone, suitable for audiobook reading."

*Analysis*: This description effectively conveys the emotion and style of the voice through words like "intellectual" and "calm."


**Bad examples and suggested improvements**:Bad exampleMain problemSuggested improvement"A nice voice"The description is too vague and subjective, lacking actionable detail.Add concrete dimensions, such as "a young female voice with a clear timbre and a soft intonation.""A voice like a certain celebrity"Involves copyright risks, and the model doesn&#x27;t support direct imitation.Extract and describe vocal characteristics, such as "a mature, magnetic male voice with a composed pace.""A very, very, very nice female voice"The description is redundant; repeated words don&#x27;t help define the voice.Remove repetition and add effective details, such as "a female voice aged 20-24 with a light timbre, lively intonation, and sweet quality."123456Invalid input that can&#x27;t be parsed into vocal characteristics.Provide a meaningful text description; see the recommended examples above. Qwen3-TTS supports voice cloning (Qwen3-TTS-VC) and voice design (Qwen3-TTS-VD). See the [Voice cloning](/api-reference/speech-synthesis/voice-cloning/qwen/create-voice) guide for details. 
## [​ ](#connection-reuse-websocket) Connection reuse (WebSocket)

WebSocket connections support reuse: after a synthesis task completes, the next task can start on the same connection without reconnecting.
**Reuse flow**:

- **CosyVoice**: The client sends `finish-task`. After the server returns a `task-finished` event, the client can send a `run-task` event to start a new task.

- **Qwen-TTS**: The client sends `session.finish`. After the server returns `session.finished`, the client can establish a new session for the next task.


 
- Wait for the server to return the completion event (`task-finished` or `session.finished`) before starting a new task.

- CosyVoice requires a different `task_id` for each task on a reused connection.

- If a task fails, the server returns an error event and closes the connection. The connection cannot be reused.

- If no new task starts within 60 seconds, the connection automatically disconnects.


 
For event details, see the corresponding [API reference](#api-reference).
## [​ ](#high-concurrency-best-practices) High-concurrency best practices

The DashScope SDK includes built-in pooling mechanisms to reuse WebSocket connections and synthesis objects, reducing the overhead of frequent creation and destruction.
**Prerequisites**:

- [Obtain an API key](/api-reference/preparation/api-key)

- A DashScope SDK version that meets the requirements is installed. We recommend [installing the latest version](/api-reference/preparation/install-sdk):

Python SDK: Version 1.25.2 or later

- Java SDK: Version 2.16.6 or later


- Python SDK 
- Java SDK 

 The Python SDK uses `SpeechSynthesizerObjectPool` to manage and reuse `SpeechSynthesizer` objects.The object pool creates the specified number of `SpeechSynthesizer` instances and establishes WebSocket connections at initialization. Objects obtained from the pool can start requests immediately, reducing time to first audio. When returned, connections remain active for the next task.#### Implementation steps

- 
Install the dependency: Install the DashScope package (`pip install -U dashscope`).


- 
Create and configure the object pool.
Set the pool size to 1.5-2 times your peak concurrency. The pool size must not exceed your account&#x27;s QPS limit.
Create a global singleton object pool. Connections are established during initialization, which takes some time:


Copy ```\nfrom dashscope.audio.tts_v2 import SpeechSynthesizerObjectPool
synthesizer_object_pool = SpeechSynthesizerObjectPool(max_size=20)
import dashscope
dashscope.base_http_api_url = "https://dashscope-intl.aliyuncs.com/api/v1"

``` 
- `SpeechSynthesizerObjectPool` establishes WebSocket connections using the current global `dashscope.api_key` value at initialization. The API key is written to the `Authorization` header only during the WebSocket handshake; subsequent task messages (such as `run-task`) do not carry it. **Changing dashscope.api_key after the pool is created does not affect existing connections.** Objects returned by `borrow_synthesizer`—including objects reused after being returned—continue to use the API key from the original handshake. The new value is silently ignored, which can cause identity, quota, or billing attribution to differ from expectations. `borrow_synthesizer` does not support specifying an API key through its parameters.

- If you need to use multiple different API keys, maintain a **separate** `SpeechSynthesizerObjectPool` **instance** for each API key.


 
- 
Borrow a `SpeechSynthesizer` object from the pool.
If the number of unreturned objects exceeds the pool capacity, the system creates an additional object. Such objects require a new connection and do not provide the reuse benefit.


Copy ```\nspeech_synthesizer = connectionPool.borrow_synthesizer(
 model=&#x27;cosyvoice-v3-flash&#x27;,
 voice=&#x27;longanyang&#x27;,
 seed=12382,
 callback=synthesizer_callback
)

``` 
- Perform speech synthesis. Call the `SpeechSynthesizer` object&#x27;s `call` or `streaming_call` method to perform speech synthesis.

- Return the `SpeechSynthesizer` object. After the task completes, return the object to the pool for reuse. Do not return objects with incomplete or failed tasks.


Copy ```\nconnectionPool.return_synthesizer(speech_synthesizer)

``` Complete code

 Before copying: `SpeechSynthesizerObjectPool` establishes WebSocket connections and authenticates using the current global `dashscope.api_key` at initialization. **Changing dashscope.api_key after the pool is created does not affect existing connections**; the new value is silently ignored. For multi-key scenarios, maintain a separate pool instance per API key. See the important note above for details. Copy ```\n# !/usr/bin/env python3
# Copyright (C) Alibaba Group. All Rights Reserved.
# MIT License (https://opensource.org/licenses/MIT)
import os
import time
import threading
import dashscope
from dashscope.audio.tts_v2 import *
USE_CONNECTION_POOL = True
text_to_synthesize = [
 &#x27;First sentence: Welcome to Alibaba Cloud speech synthesis.&#x27;,
 &#x27;Second sentence: Welcome to Alibaba Cloud speech synthesis.&#x27;,
 &#x27;Third sentence: Welcome to Alibaba Cloud speech synthesis.&#x27;,
]
connectionPool = None
def init_dashscope_api_key():
 &#x27;&#x27;&#x27;
 Set your DashScope API-key. More information:
 https://github.com/aliyun/alibabacloud-bailian-speech-demo/blob/master/PREREQUISITES.md
 &#x27;&#x27;&#x27;
 if &#x27;DASHSCOPE_API_KEY&#x27; in os.environ:
 dashscope.api_key = os.environ[
 &#x27;DASHSCOPE_API_KEY&#x27;] # load API-key from environment variable DASHSCOPE_API_KEY
 else:
 dashscope.api_key = &#x27;<your-dashscope-api-key>&#x27; # set API-key manually
def synthesis_text_to_speech_and_play_by_streaming_mode(text, task_id):
 global USE_CONNECTION_POOL, connectionPool
 &#x27;&#x27;&#x27;
 Synthesize speech with given text by streaming mode, async call and play the synthesized audio in real-time.
 for more information, please refer to https://help.aliyun.com/document_detail/2712523.html
 &#x27;&#x27;&#x27;
 complete_event = threading.Event()
 # Define a callback to handle the result
 class Callback(ResultCallback):
 def on_open(self):
 # when using object pool, on_open will be called after task start
 self.file = open(f&#x27;result_{task_id}.mp3&#x27;, &#x27;wb&#x27;)
 print(f&#x27;[task_{task_id}] start&#x27;)
 def on_complete(self):
 print(f&#x27;[task_{task_id}] speech synthesis task complete successfully.&#x27;)
 complete_event.set()
 def on_error(self, message: str):
 print(f&#x27;[task_{task_id}] speech synthesis task failed, {message}&#x27;)
 def on_close(self):
 # when using object pool, on_close will be called after task finished
 print(f&#x27;[task_{task_id}] finished&#x27;)
 def on_event(self, message):
 # print(f&#x27;recv speech synthesis message {message}&#x27;)
 pass
 def on_data(self, data: bytes) -> None:
 # send to player
 # save audio to file
 self.file.write(data)
 # Call the speech synthesizer callback
 synthesizer_callback = Callback()
 # Initialize the speech synthesizer
 # you can customize the synthesis parameters, like voice, format, sample_rate or other parameters
 if USE_CONNECTION_POOL:
 speech_synthesizer = connectionPool.borrow_synthesizer(
 model=&#x27;cosyvoice-v3-flash&#x27;,
 voice=&#x27;longanyang&#x27;,
 seed=12382,
 callback=synthesizer_callback
 )
 else:
 speech_synthesizer = SpeechSynthesizer(model=&#x27;cosyvoice-v3-flash&#x27;,
 voice=&#x27;longanyang&#x27;,
 seed=12382,
 callback=synthesizer_callback)
 try:
 speech_synthesizer.call(text)
 except Exception as e:
 print(f&#x27;[task_{task_id}] speech synthesis task failed, {e}&#x27;)
 if USE_CONNECTION_POOL:
 # close the synthesizer connection manually if task failed when using connection pool.
 speech_synthesizer.close()
 return
 print(&#x27;[task_{}] Synthesized text: {}&#x27;.format(task_id, text))
 complete_event.wait()
 print(&#x27;[task_{}][Metric] requestId: {}, first package delay ms: {}&#x27;.format(
 task_id,
 speech_synthesizer.get_last_request_id(),
 speech_synthesizer.get_first_package_delay()))
 if USE_CONNECTION_POOL:
 connectionPool.return_synthesizer(speech_synthesizer)
# main function
if __name__ == &#x27;__main__&#x27;:
 # You must set dashscope.api_key and base_websocket_api_url before creating the SpeechSynthesizerObjectPool.
 # The pool establishes WebSocket connections using the current global dashscope.api_key at initialization.
 # Changing dashscope.api_key after the pool is created does not affect existing connections.
 dashscope.base_websocket_api_url=&#x27;wss://dashscope-intl.aliyuncs.com/api-ws/v1/inference&#x27;
 init_dashscope_api_key()
 if USE_CONNECTION_POOL:
 print(&#x27;creating connection pool&#x27;)
 start_time = time.time() * 1000
 connectionPool = SpeechSynthesizerObjectPool(max_size=3)
 end_time = time.time() * 1000
 print(&#x27;connection pool created, cost: {} ms&#x27;.format(end_time - start_time))
 task_thread_list = []
 for task_id in range(3):
 thread = threading.Thread(
 target=synthesis_text_to_speech_and_play_by_streaming_mode,
 args=(text_to_synthesize[task_id], task_id))
 task_thread_list.append(thread)
 for task_thread in task_thread_list:
 task_thread.start()
 for task_thread in task_thread_list:
 task_thread.join()
 if USE_CONNECTION_POOL:
 connectionPool.shutdown()

``` #### Resource management and error handling

- 
Task success: When a speech synthesis task completes normally, call `connectionPool.return_synthesizer(speech_synthesizer)` to return the `SpeechSynthesizer` object to the pool for reuse.
 Do not return `SpeechSynthesizer` objects with incomplete or failed tasks. 


- 
Task failure: When an SDK internal error or business logic exception interrupts the task, close the underlying WebSocket connection: `speech_synthesizer.close()`.


- 
After all speech synthesis tasks complete, shut down the object pool: `connectionPool.shutdown()`.


- 
When the service returns a TaskFailed error, no additional handling is required.


 The Java SDK achieves optimal performance through a built-in connection pool and a custom object pool working together.
- Connection pool: The OkHttp3 connection pool integrated within the SDK manages and reuses underlying WebSocket connections to reduce network handshake overhead. This feature is enabled by default.

- Object pool: Implemented based on `commons-pool2`, it maintains a group of `SpeechSynthesizer` objects with pre-established connections. Obtaining objects from the pool eliminates connection establishment latency and significantly reduces time to first audio.


#### Implementation steps

- Add dependencies. Add dashscope-sdk-java and commons-pool2 to the dependency configuration file of your project build tool. The following examples show the configurations for Maven and Gradle:


Maven Gradle Copy ```\n<dependency>
 <groupId>com.alibaba</groupId>
 <artifactId>dashscope-sdk-java</artifactId>
 <!-- Replace &#x27;the-latest-version&#x27; with version 2.16.9 or later. Find the version at: https://mvnrepository.com/artifact/com.alibaba/dashscope-sdk-java -->
 <version>the-latest-version</version>
</dependency>
<dependency>
 <groupId>org.apache.commons</groupId>
 <artifactId>commons-pool2</artifactId>
 <!-- Replace &#x27;the-latest-version&#x27; with the latest version. Find the version at: https://mvnrepository.com/artifact/org.apache.commons/commons-pool2 -->
 <version>the-latest-version</version>
</dependency>

``` 
- Configure the connection pool. Configure the key connection pool parameters by using environment variables:


Environment variableDescriptionDASHSCOPE_CONNECTION_POOL_SIZEThe connection pool size. Recommended value: at least twice the peak per-server concurrency. Default value: 32.DASHSCOPE_MAXIMUM_ASYNC_REQUESTSThe maximum number of asynchronous requests. Recommended value: the same as `DASHSCOPE_CONNECTION_POOL_SIZE`. Default value: 32.DASHSCOPE_MAXIMUM_ASYNC_REQUESTS_PER_HOSTThe maximum number of asynchronous requests per host. Recommended value: the same as `DASHSCOPE_CONNECTION_POOL_SIZE`. Default value: 32. 
- Configure the object pool. Configure the object pool size by using environment variables:


Environment variableDescriptionCOSYVOICE_OBJECTPOOL_SIZEThe object pool size. Recommended value: 1.5 to 2 times the peak per-server concurrency. Default value: 500. 
- The object pool size (`COSYVOICE_OBJECTPOOL_SIZE`) must be less than or equal to the connection pool size (`DASHSCOPE_CONNECTION_POOL_SIZE`). Otherwise, when the object pool requests objects while the connection pool is full, the calling thread will be blocked waiting for an available connection.

- The object pool size must not exceed the QPS (queries per second) limit of your account.


 Create the object pool by using the following code:Copy ```\nclass CosyvoiceObjectPool {
 // ...Other code is omitted here. For the complete example, see the complete code section.
 public static GenericObjectPool<SpeechSynthesizer> getInstance() {
 lock.lock();
 if (synthesizerPool == null) {
 // You can set the object pool size here, or configure it in the COSYVOICE_OBJECTPOOL_SIZE environment variable.
 // We recommend setting this value to 1.5 to 2 times the maximum per-server concurrency.
 int objectPoolSize = getObjectivePoolSize();
 SpeechSynthesizerObjectFactory speechSynthesizerObjectFactory =
 new SpeechSynthesizerObjectFactory();
 GenericObjectPoolConfig<SpeechSynthesizer> config =
 new GenericObjectPoolConfig<>();
 config.setMaxTotal(objectPoolSize);
 config.setMaxIdle(objectPoolSize);
 config.setMinIdle(objectPoolSize);
 synthesizerPool =
 new GenericObjectPool<>(speechSynthesizerObjectFactory, config);
 }
 lock.unlock();
 return synthesizerPool;
 }
}

``` 
- Obtain a `SpeechSynthesizer` object from the object pool. If the number of unreturned objects exceeds the maximum capacity of the object pool, the system creates an additional `SpeechSynthesizer` object. Such newly created objects need to be re-initialized and establish new WebSocket connections. They cannot leverage existing connections in the object pool and therefore do not provide the reuse benefit.


Copy ```\nsynthesizer = CosyvoiceObjectPool.getInstance().borrowObject();

``` 
- 
Perform speech synthesis. After borrowing a `SpeechSynthesizer` object from the pool, call `updateParamAndCallback(param, callback)` to bind the parameters and callback for this task, then call `streamingCall` or `call` to start synthesis.
 
In object pool scenarios, `updateParamAndCallback` is called once per borrow to update the callback and task-level parameters (such as `voice` and `format`). **The apiKey passed in each call must remain the same.** `updateParamAndCallback` updates only the local fields of the current `SpeechSynthesizer` instance; it does not rebuild the underlying WebSocket connection. The SDK writes the `apiKey` to the `Authorization` header only during the WebSocket handshake; subsequent task messages (such as `run-task`) do not carry it. Because the reused connection remains open, a different `apiKey` is never transmitted to the server—the request continues to use the `apiKey` from the original handshake, which can cause identity, quota, or billing attribution to differ from expectations.

- If you need to use multiple different API keys, maintain a **separate object pool instance** for each API key.


 

- 
Return the `SpeechSynthesizer` object. After the speech synthesis task is complete, return the `SpeechSynthesizer` object so that subsequent tasks can reuse it. Do not return objects whose tasks are incomplete or have failed.


Copy ```\nCosyvoiceObjectPool.getInstance().returnObject(synthesizer);

``` Complete code

 Before copying: In object pool scenarios, the `apiKey` passed to each `updateParamAndCallback` call **must remain the same**. The SDK does not update the `apiKey` of established connections; passing a different `apiKey` has no effect. For multi-key scenarios, maintain a separate object pool instance per API key. See the important note in the **Perform speech synthesis** section for details. Copy ```\nimport com.alibaba.dashscope.audio.tts.SpeechSynthesisResult;
import com.alibaba.dashscope.audio.ttsv2.SpeechSynthesisAudioFormat;
import com.alibaba.dashscope.audio.ttsv2.SpeechSynthesisParam;
import com.alibaba.dashscope.audio.ttsv2.SpeechSynthesizer;
import com.alibaba.dashscope.common.ResultCallback;
import com.alibaba.dashscope.exception.NoApiKeyException;
import com.alibaba.dashscope.utils.Constants;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.pool2.BasePooledObjectFactory;
import org.apache.commons.pool2.PooledObject;
import org.apache.commons.pool2.impl.DefaultPooledObject;
import org.apache.commons.pool2.impl.GenericObjectPool;
import org.apache.commons.pool2.impl.GenericObjectPoolConfig;
import java.time.LocalDateTime;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.locks.Lock;
/**
 * You need to import org.apache.commons.pool2 and DashScope related packages in your project.
 *
 * DashScope SDK 2.16.6 and later versions are optimized for high-concurrency scenarios.
 * Versions earlier than DashScope SDK 2.16.6 are not recommended for high-concurrency scenarios.
 *
 *
 * Before making high-concurrency calls to the TTS service,
 * configure the connection pool parameters by using the following environment variables.
 *
 * DASHSCOPE_MAXIMUM_ASYNC_REQUESTS
 * DASHSCOPE_MAXIMUM_ASYNC_REQUESTS_PER_HOST
 * DASHSCOPE_CONNECTION_POOL_SIZE
 *
 */
class SpeechSynthesizerObjectFactory
 extends BasePooledObjectFactory<SpeechSynthesizer> {
 public SpeechSynthesizerObjectFactory() {
 super();
 }
 @Override
 public SpeechSynthesizer create() throws Exception {
 return new SpeechSynthesizer();
 }
 @Override
 public PooledObject<SpeechSynthesizer> wrap(SpeechSynthesizer obj) {
 return new DefaultPooledObject<>(obj);
 }
}
class CosyvoiceObjectPool {
 public static GenericObjectPool<SpeechSynthesizer> synthesizerPool;
 public static String COSYVOICE_OBJECTPOOL_SIZE_ENV = "COSYVOICE_OBJECTPOOL_SIZE";
 public static int DEFAULT_OBJECT_POOL_SIZE = 500;
 private static Lock lock = new java.util.concurrent.locks.ReentrantLock();
 public static int getObjectivePoolSize() {
 try {
 Integer n = Integer.parseInt(System.getenv(COSYVOICE_OBJECTPOOL_SIZE_ENV));
 System.out.println("Using Object Pool Size In Env: "+ n);
 return n;
 } catch (NumberFormatException e) {
 System.out.println("Using Default Object Pool Size: "+ DEFAULT_OBJECT_POOL_SIZE);
 return DEFAULT_OBJECT_POOL_SIZE;
 }
 }
 public static GenericObjectPool<SpeechSynthesizer> getInstance() {
 lock.lock();
 if (synthesizerPool == null) {
 // You can set the object pool size here, or configure it in the COSYVOICE_OBJECTPOOL_SIZE environment variable.
 // We recommend setting this value to 1.5 to 2 times the maximum per-server concurrency.
 int objectPoolSize = getObjectivePoolSize();
 SpeechSynthesizerObjectFactory speechSynthesizerObjectFactory =
 new SpeechSynthesizerObjectFactory();
 GenericObjectPoolConfig<SpeechSynthesizer> config =
 new GenericObjectPoolConfig<>();
 config.setMaxTotal(objectPoolSize);
 config.setMaxIdle(objectPoolSize);
 config.setMinIdle(objectPoolSize);
 synthesizerPool =
 new GenericObjectPool<>(speechSynthesizerObjectFactory, config);
 }
 lock.unlock();
 return synthesizerPool;
 }
}
class SynthesizeTaskWithCallback implements Runnable {
 String[] textArray;
 String requestId;
 long timeCost;
 public SynthesizeTaskWithCallback(String[] textArray) {
 this.textArray = textArray;
 }
 @Override
 public void run() {
 SpeechSynthesizer synthesizer = null;
 long startTime = System.currentTimeMillis();
 // if recv onError
 final boolean[] hasError = {false};
 try {
 class ReactCallback extends ResultCallback<SpeechSynthesisResult> {
 ReactCallback() {}
 @Override
 public void onEvent(SpeechSynthesisResult message) {
 if (message.getAudioFrame() != null) {
 try {
 byte[] bytesArray = message.getAudioFrame().array();
 System.out.println("Audio received. Audio stream length: " + bytesArray.length);
 } catch (Exception e) {
 throw new RuntimeException(e);
 }
 }
 }
 @Override
 public void onComplete() {}
 @Override
 public void onError(Exception e) {
 System.out.println(e.getMessage());
 e.printStackTrace();
 hasError[0] = true;
 }
 }
 SpeechSynthesisParam param =
 SpeechSynthesisParam.builder()
 .model("cosyvoice-v3-flash")
 .voice("longanyang")
 // Get an API Key: https://home.qwencloud.com/api-keys
 // If no environment variable is configured, replace the next line with: .apiKey("sk-xxx")
 .apiKey(System.getenv("DASHSCOPE_API_KEY"))
 .format(SpeechSynthesisAudioFormat
 .MP3_22050HZ_MONO_256KBPS) // Use PCM or MP3 for streaming synthesis
 .build();
 try {
 synthesizer = CosyvoiceObjectPool.getInstance().borrowObject();
 // Important: In object pool scenarios, the apiKey passed to each updateParamAndCallback call must remain the same.
 synthesizer.updateParamAndCallback(param, new ReactCallback());
 for (String text : textArray) {
 synthesizer.streamingCall(text);
 }
 Thread.sleep(20);
 synthesizer.streamingComplete(60000);
 requestId = synthesizer.getLastRequestId();
 } catch (Exception e) {
 System.out.println("Exception e: " + e.toString());
 hasError[0] = true;
 }
 } catch (Exception e) {
 hasError[0] = true;
 throw new RuntimeException(e);
 }
 if (synthesizer != null) {
 try {
 if (hasError[0] == true) {
 // If an error occurs, close the connection and invalidate the object in the object pool.
 synthesizer.getDuplexApi().close(1000, "bye");
 CosyvoiceObjectPool.getInstance().invalidateObject(synthesizer);
 } else {
 // If the task completes normally, return the object to the pool.
 CosyvoiceObjectPool.getInstance().returnObject(synthesizer);
 }
 } catch (Exception e) {
 throw new RuntimeException(e);
 }
 long endTime = System.currentTimeMillis();
 timeCost = endTime - startTime;
 System.out.println("[Thread " + Thread.currentThread() + "] Speech synthesis task completed. Time elapsed: " + timeCost + " ms, RequestId " + requestId);
 }
 }
}
@Slf4j
public class SynthesizeTextToSpeechWithCallbackConcurrently {
 public static void checkoutEnv(String envName, int defaultSize) {
 if (System.getenv(envName) != null) {
 System.out.println("[ENV CHECK]: " + envName + " "
 + System.getenv(envName));
 } else {
 System.out.println("[ENV CHECK]: " + envName
 + " Using Default which is " + defaultSize);
 }
 }
 public static void main(String[] args)
 throws InterruptedException, NoApiKeyException {
 Constants.baseWebsocketApiUrl = "wss://dashscope-intl.aliyuncs.com/api-ws/v1/inference";
 // Check for connection pool env
 checkoutEnv("DASHSCOPE_CONNECTION_POOL_SIZE", 32);
 checkoutEnv("DASHSCOPE_MAXIMUM_ASYNC_REQUESTS", 32);
 checkoutEnv("DASHSCOPE_MAXIMUM_ASYNC_REQUESTS_PER_HOST", 32);
 checkoutEnv(CosyvoiceObjectPool.COSYVOICE_OBJECTPOOL_SIZE_ENV, CosyvoiceObjectPool.DEFAULT_OBJECT_POOL_SIZE);
 int runTimes = 3;
 // Create the pool of SpeechSynthesis objects
 ExecutorService executorService = Executors.newFixedThreadPool(runTimes);
 for (int i = 0; i < runTimes; i++) {
 // Record the task submission time
 LocalDateTime submissionTime = LocalDateTime.now();
 executorService.submit(new SynthesizeTaskWithCallback(new String[] {
 "Before my bed, moonlight gleams,", "It seems like frost upon the ground.", "I lift my gaze to watch the bright moon,", "Then bow my head, thinking of home."}));
 }
 // Shut down the ExecutorService and wait for all tasks to complete
 executorService.shutdown();
 executorService.awaitTermination(1, TimeUnit.MINUTES);
 System.exit(0);
 }
}

``` #### Recommended configuration
The following configurations are based on test results from running only the CosyVoice speech synthesis service on Alibaba Cloud servers with specified specifications. Excessively high concurrency may cause task processing delays. Here, per-server concurrency refers to the number of CosyVoice speech synthesis tasks running simultaneously at a given moment, which can also be understood as the worker thread count.Server specifications (Alibaba Cloud)Maximum per-server concurrencyObject pool sizeConnection pool size4 vCPUs, 8 GiB10050020008 vCPUs, 16 GiB150500200016 vCPUs, 32 GiB2005002000 #### Resource management and error handling

- 
Task success: When a speech synthesis task completes normally, you must call the returnObject method of GenericObjectPool to return the `SpeechSynthesizer` object back to the pool for reuse. In the current code, this corresponds to `CosyvoiceObjectPool.getInstance().returnObject(synthesizer)`.
 Do not return `SpeechSynthesizer` objects whose tasks are incomplete or have failed. 


- 
Task failure: When the SDK internally or the business logic throws an exception that interrupts a task, you must perform the following two operations:

Proactively close the underlying WebSocket connection.

- Invalidate the object from the object pool to prevent it from being reused.


Copy ```\n// Close the connection
synthesizer.getDuplexApi().close(1000, "bye");
// Invalidate the synthesizer that encountered an exception from the object pool
CosyvoiceObjectPool.getInstance().invalidateObject(synthesizer);

``` 
- When the service reports a TaskFailed error, no additional handling is required.


#### Call warm-up and latency measurement guidelines
When evaluating the DashScope Java SDK performance (such as concurrent call latency), perform sufficient warm-up operations first. This ensures that measurements reflect steady-state performance and avoids data skew caused by initial connection establishment overhead.**Connection reuse mechanism**: The DashScope Java SDK efficiently manages and reuses WebSocket connections through a global singleton connection pool, designed to reduce the overhead of frequent connection establishment and disconnection and improve processing capacity in high-concurrency scenarios. This mechanism works as follows:
- **Created on demand**: The SDK does not pre-create WebSocket connections at service startup. Instead, connections are established on demand when the first call is made.

- **Time-limited reuse**: After a request completes, the connection is retained in the pool for up to 60 seconds for reuse.

If a new request arrives within 60 seconds, the existing connection is reused, avoiding repeated handshake overhead.

- If a connection remains idle for more than 60 seconds, it is automatically closed to free up resources.


**Importance of warm-up**: In the following scenarios, the connection pool may not have reusable active connections, causing requests to establish new connections:
- The application has just started and no calls have been made yet.

- The service has been idle for more than 60 seconds, and the connections in the pool have been closed due to timeout.


In these scenarios, the first request must complete a full WebSocket connection (TCP handshake, TLS negotiation, and protocol upgrade), resulting in latency significantly higher than subsequent requests on reused connections. Without warm-up, performance test results are skewed by the inclusion of connection establishment overhead.**SDK-side latency vs actual time to first audio**: The time to first audio reported on the SDK side (such as the value obtained through `get_first_package_delay()`) includes the time spent on WebSocket connection establishment and network transmission, and does not equal the actual time to first audio of the model service. The actual time to first audio refers to the time interval from when the server receives the `run-task` command to when it returns the first `result-generated` event. This value can be viewed in server-side logs. In high-concurrency scenarios, due to the establishment of a large number of connections and resource scheduling, the latency reported on the SDK side may be significantly higher than the actual time to first audio on the server side. If the SDK reports high time to first audio:
- Compare the time to first audio in the server-side logs (from `run-task` to the first `result-generated`) to verify whether the model inference performance is normal.

- Use the object pool or connection pool mechanism described above for warm-up to eliminate WebSocket connection establishment overhead, so that the latency reported on the SDK side more closely reflects the actual time to first audio.


**Recommended practices**: To obtain reliable performance data, follow these warm-up steps before performing formal performance benchmarking or latency measurement:
- Simulate the concurrency level of your formal test by sending a certain number of calls in advance (for example, continuously for 1-2 minutes) to fully populate the connection pool.

- After confirming that the connection pool has established and maintained a sufficient number of active connections, begin the formal performance data collection.


With proper warm-up, the SDK connection pool enters a stable reuse state, enabling you to measure more representative latency metrics that accurately reflect the service performance during steady-state online operation.#### Common Java SDK exceptions
 Exception 1: Server TCP connections keep increasing despite stable business traffic

 **Cause:****Type 1:** Each SDK object creates a connection when instantiated. If an object pool is not used, each object is destroyed after its task completes. The connection then enters an unreferenced state and remains open until the server triggers a connection timeout after 61 seconds, which means the connection cannot be reused during those 61 seconds.In high-concurrency scenarios, new tasks create new connections when no reusable connections are available, which leads to the following consequences:
- The number of connections keeps increasing.

- Excessive connections exhaust server resources, causing the server to become unresponsive.

- The connection pool reaches its limit, and new tasks are blocked while waiting for available connections.


**Type 2:** When MaxIdle is set to a value smaller than MaxTotal in the object pool configuration, idle objects that exceed MaxIdle are destroyed, causing connection leaks. Leaked connections remain open until the 61-second timeout triggers a disconnection. Similar to Type 1, this causes the number of connections to keep increasing.**Solution**:For Type 1, use an object pool.For Type 2, check the object pool configuration parameters. Set MaxIdle equal to MaxTotal, and disable the automatic object destruction policy of the object pool. Exception 2: Task takes 60 seconds longer than a normal call

 Same as "**Exception 1**". The connection pool has reached its maximum connection limit, and new tasks must wait 61 seconds for unreferenced connections to time out before a connection becomes available. Exception 3: Tasks are slow at service startup but gradually return to normal

 **Cause**: In high-concurrency scenarios, the same object reuses the same WebSocket connection, so WebSocket connections are only created at service startup. Note that if high-concurrency calls begin immediately during the startup phase, creating too many WebSocket connections simultaneously can cause blocking.**Solution**: Gradually increase the concurrency level after starting the service, or add warm-up tasks. Exception 4: Server error "Invalid action(&#x27;run-task&#x27;)! Please follow the protocol!"

 **Cause**: This occurs when a client-side error happens but the server is unaware of it, leaving the connection in an active task state. When the connection and object are reused for the next task, a protocol error occurs and the next task fails.**Solution**: Actively close the WebSocket connection after an exception is thrown, and then return the object to the object pool. Exception 5: Abnormal call volume spikes despite stable business traffic

 **Cause**: Creating too many WebSocket connections simultaneously causes blocking, but business traffic continues to arrive, leading to a short-term task backlog. After the blocking resolves, all backlogged tasks are executed immediately. This causes call volume spikes and may momentarily exceed the account&#x27;s concurrency limit, resulting in partial task failures, server unresponsiveness, and other issues.This situation of creating too many WebSocket connections simultaneously often occurs during:
- The service startup phase

- Network anomalies that cause a large number of WebSocket connections to disconnect and reconnect simultaneously

- A large number of server-side errors occurring at the same time, triggering mass WebSocket reconnections. A common error is exceeding the account&#x27;s concurrency limit ("Requests rate limit exceeded, please try again later.").


**Solution**:
- Check the network conditions.

- Investigate whether a large number of other server-side errors occurred before the spike.

- Increase the account&#x27;s concurrency limit.

- Reduce the object pool and connection pool sizes to limit the maximum concurrency level through the object pool upper limit.

- Upgrade server configurations or add more machines.


 Exception 6: All tasks slow down as the concurrency level increases

 **Solution**:
- Check whether the network bandwidth limit has been reached.

- Check whether the actual concurrency level is too high.


 
## [​ ](#supported-scope) Supported scope

To call the following models, use your Qwen Cloud API key:

- **CosyVoice:** cosyvoice-v3-plus, cosyvoice-v3-flash

- **Qwen-TTS:**

**Qwen3-TTS-Instruct-Flash-Realtime**: qwen3-tts-instruct-flash-realtime (stable, currently equivalent to qwen3-tts-instruct-flash-realtime-2026-01-22), qwen3-tts-instruct-flash-realtime-2026-01-22 (latest snapshot)

- **Qwen3-TTS-VD-Realtime:** qwen3-tts-vd-realtime-2026-01-15 (latest snapshot), qwen3-tts-vd-realtime-2025-12-16 (snapshot)

- **Qwen3-TTS-VC-Realtime:** qwen3-tts-vc-realtime-2026-01-15 (latest snapshot), qwen3-tts-vc-realtime-2025-11-27 (snapshot)

- **Qwen3-TTS-Flash-Realtime:** qwen3-tts-flash-realtime (stable, currently equivalent to qwen3-tts-flash-realtime-2025-11-27), qwen3-tts-flash-realtime-2025-11-27 (latest snapshot), qwen3-tts-flash-realtime-2025-09-18 (snapshot)


## [​ ](#supported-voices) Supported voices

Different models support different voices. Set the `voice` request parameter to the value in the **voice parameter** column of the voice list.

- [CosyVoice voice list](/api-reference/speech-synthesis/cosyvoice/voice-list)

- [Qwen-TTS voice list](/api-reference/speech-synthesis/qwen-tts/voice-list)


## [​ ](#api-reference) API reference


- Real-time speech synthesis - [CosyVoice API reference](/api-reference/speech-synthesis/cosyvoice/websocket-api)

- Real-time speech synthesis - [Qwen-TTS API reference](/api-reference/speech-synthesis/qwen-tts-realtime/client-events)


## [​ ](#faq) FAQ

### [​ ](#q-how-do-i-fix-incorrect-pronunciation-in-speech-synthesis-how-do-i-control-the-pronunciation-of-polyphonic-characters) Q: How do I fix incorrect pronunciation in speech synthesis? How do I control the pronunciation of polyphonic characters?


- Replace the polyphonic character with a homophone to quickly fix the pronunciation issue.

- Use SSML markup language to control pronunciation.


### [​ ](#q-how-do-i-troubleshoot-silent-audio-when-using-a-cloned-voice) Q: How do I troubleshoot silent audio when using a cloned voice?


- 
**Verify the voice status**
Call the voice cloning/design API and confirm that the voice `status` is `OK`.


- 
**Check model version consistency**
Make sure the `target_model` parameter used during voice cloning matches the `model` parameter used during speech synthesis. For example:

If you used `cosyvoice-v3-plus` for cloning

- You must also use `cosyvoice-v3-plus` for synthesis


- 
**Verify source audio quality**
Check whether the source audio used for voice cloning meets the audio requirements and best practices:

Audio duration: 10-20 seconds

- Clear audio quality

- No background noise


- 
**Check request parameters**
Confirm that the `voice` parameter in your speech synthesis request is set to the cloned voice ID.


### [​ ](#q-what-should-i-do-if-the-cloned-voice-produces-unstable-or-incomplete-speech) Q: What should I do if the cloned voice produces unstable or incomplete speech?

If the synthesized speech from a cloned voice exhibits any of the following issues:

- Incomplete playback that only reads part of the text

- Inconsistent synthesis quality

- Abnormal pauses or silent segments in the speech


**Possible cause**: The source audio quality doesn&#x27;t meet the requirements.
**Solution**: Check whether the source audio meets the [Recording guide for voice cloning](/developer-guides/speech/voice-cloning#audio-requirements) requirements. We recommend re-recording based on the recording guidelines.
### [​ ](#q-why-does-the-actual-duration-differ-from-the-duration-displayed-in-the-wav-file-header) Q: Why does the actual duration differ from the duration displayed in the WAV file header?

Speech synthesis uses a streaming mechanism that returns data progressively as it&#x27;s generated. The duration in the saved WAV file header is an estimate and may contain inaccuracies. For precise duration, set `format` to `pcm`, wait for the complete synthesis result, and then add the WAV file header yourself.
### [​ ](#q-why-wont-the-audio-play) Q: Why won&#x27;t the audio play?

Troubleshoot based on the following scenarios:

- Audio saved as a complete file (such as xx.mp3)

Audio format consistency: The audio format in the request parameters must match the file extension (for example, if the format is set to `wav`, the file must be saved as `.wav`).

- Player compatibility: Confirm that your player supports the audio format and sample rate.


- Streaming audio playback

Save the audio stream as a complete file and try playing it with a media player. If the file won&#x27;t play, refer to the troubleshooting steps in Scenario 1.

- If the file plays correctly, the issue is in the streaming playback implementation. Confirm that your player supports streaming playback (such as ffmpeg, pyaudio, AudioFormat, or MediaSource).


### [​ ](#q-why-is-audio-playback-stuttering) Q: Why is audio playback stuttering?

Troubleshoot with the following steps:

- **Check text send rate**: Make sure the interval between text segments is short enough that the next segment arrives before the previous audio finishes playing.

- **Check callback function performance**:

Confirm that the callback function has no blocking business logic.

- The callback runs on the WebSocket thread. Blocking it delays data reception. Write audio data to a separate buffer and process it in a separate thread.


- **Check network stability**: Network fluctuations can cause audio transmission interruptions or delays.


### [​ ](#q-why-is-speech-synthesis-taking-a-long-time) Q: Why is speech synthesis taking a long time?

Troubleshoot with the following steps:

- 
Check input interval
For streaming synthesis, confirm that the interval between text segments isn&#x27;t too long. Long intervals increase total synthesis time.


- 
Analyze performance metrics

First-packet latency: typically around 500 ms.

- RTF (Real-Time Factor = total synthesis time / audio duration): should be less than 1.0 under normal conditions.


### [​ ](#q-how-do-i-restrict-an-api-key-to-speech-synthesis-only-permission-isolation) Q: How do I restrict an API key to speech synthesis only (permission isolation)?

Create a new workspace and authorize only specific models to limit the scope of an API key. For details, see Manage workspaces. [Previous ](/developer-guides/speech/tts-models)[Non-real-time speech synthesis Non-realtime speech synthesis with Qwen3-TTS Next ](/developer-guides/speech/tts)
