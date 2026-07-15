# Real-time audio and video translation

> **Source:** https://docs.qwencloud.com/developer-guides/speech/realtime-translation

Real-time speech translation with 3-second latency

 Copy page ## [​ ](#model-details) Model details

qwen3.5-livetranslate-flash-realtime is a vision-enhanced real-time translation model supporting 60 languages (29 with audio + text, 31 text-only). It processes audio and image input from video streams or local files, uses visual context to improve accuracy, and outputs translated text and audio in real time.
Key features:

- **Multi-language support**: Translates between 60 languages — 29 with audio and text output, 31 with text-only output — including Chinese, English, French, German, Russian, Japanese, Korean, Spanish, Portuguese, and Arabic.

- **Visual enhancement**: Analyzes visual cues, such as lip movements, gestures, and on-screen text, to improve translation accuracy, especially in noisy environments or for ambiguous words.

- **3-second latency**: Delivers simultaneous interpretation with latency as low as 3 seconds.

- **Lossless simultaneous interpretation**: Predicts semantic units to resolve cross-language word order differences, achieving quality comparable to offline translation.

- **Natural voice**: Matches the intonation and emotion of the source audio automatically.

- **Hotword configuration**: Configurable hotwords improve translation accuracy for specific terms.

- **Voice cloning**: Clones the speaker&#x27;s voice for translated output. Supports server-side real-time cloning and pre-cloned voice profiles.


ModelVersionContext windowMax inputMax output**qwen3.5-livetranslate-flash-realtime** (Alias for qwen3.5-livetranslate-flash-realtime-2026-05-19)Stable53,24849,1524,096qwen3.5-livetranslate-flash-realtime-2026-05-19Snapshot53,24849,1524,096**qwen3-livetranslate-flash-realtime** (Alias for qwen3-livetranslate-flash-realtime-2025-09-22)Stable53,24849,1524,096qwen3-livetranslate-flash-realtime-2025-09-22Snapshot53,24849,1524,096 
## [​ ](#getting-started) Getting started

### [​ ](#prepare-the-environment) Prepare the environment

Requires Python 3.10 or later.
First, install pyaudio.
- macOS 
- Debian/Ubuntu 
- CentOS 
- Windows 

 Copy ```\nbrew install portaudio && pip install pyaudio

``` Copy ```\nsudo apt-get install python3-pyaudio
# or
pip install pyaudio

``` Copy ```\nsudo yum install -y portaudio portaudio-devel && pip install pyaudio

``` Copy ```\npip install pyaudio

``` 
Then install the WebSocket dependencies:
Copy ```\npip install websocket-client==1.8.0 websockets

``` 
### [​ ](#create-the-client) Create the client

Create a file named `livetranslate_client.py` with the following code:
 Client code - livetranslate_client.py

 Copy ```\nimport os
import time
import base64
import asyncio
import json
import websockets
import pyaudio
import queue
import threading
import traceback

class LiveTranslateClient:
 def __init__(self, api_key: str, target_language: str = "en", *, audio_enabled: bool = True):
 if not api_key:
 raise ValueError("API key cannot be empty.")

 self.api_key = api_key
 self.target_language = target_language
 self.audio_enabled = audio_enabled
 self.ws = None
 self.api_url = "wss://dashscope-intl.aliyuncs.com/api-ws/v1/realtime?model=qwen3.5-livetranslate-flash-realtime"

 # Audio input configuration (from microphone)
 self.input_rate = 16000
 self.input_chunk = 1600
 self.input_format = pyaudio.paInt16
 self.input_channels = 1

 # Audio output configuration (for playback)
 self.output_rate = 24000
 self.output_chunk = 2400
 self.output_format = pyaudio.paInt16
 self.output_channels = 1

 # State management
 self.is_connected = False
 self.audio_player_thread = None
 self.audio_playback_queue = queue.Queue()
 self.pyaudio_instance = pyaudio.PyAudio()

 async def connect(self):
 """Establish a WebSocket connection to the translation service."""
 headers = {"Authorization": f"Bearer {self.api_key}"}
 try:
 self.ws = await websockets.connect(self.api_url, additional_headers=headers)
 self.is_connected = True
 print(f"Successfully connected to the server: {self.api_url}")
 await self.configure_session()
 except Exception as e:
 print(f"Connection failed: {e}")
 self.is_connected = False
 raise

 async def configure_session(self):
 """Configure the translation session, setting the target language, voice, etc."""
 config = {
 "event_id": f"event_{int(time.time() * 1000)}",
 "type": "session.update",
 "session": {
 # &#x27;modalities&#x27; controls the output type.
 # ["text", "audio"]: Returns both translated text and synthesized audio (recommended).
 # ["text"]: Returns only the translated text.
 "modalities": ["text", "audio"] if self.audio_enabled else ["text"],
 "input_audio_format": "pcm",
 "output_audio_format": "pcm",
 # &#x27;input_audio_transcription&#x27; configures source language recognition.
 # Set &#x27;model&#x27; to &#x27;qwen3-asr-flash-realtime&#x27; to also output the source language recognition result.
 # "input_audio_transcription": {
 # "model": "qwen3-asr-flash-realtime",
 # "language": "zh" # source language, default &#x27;en&#x27;
 # },
 "translation": {
 "language": self.target_language,
 # &#x27;corpus&#x27; configures hotwords to improve the translation accuracy of specific terms.
 # "corpus": {
 # "phrases": {
 # "Artificial Intelligence": "Artificial Intelligence",
 # "Machine Learning": "Machine Learning"
 # }
 # }
 }
 }
 }
 print(f"Sending session configuration: {json.dumps(config, indent=2, ensure_ascii=False)}")
 await self.ws.send(json.dumps(config))

 async def send_audio_chunk(self, audio_data: bytes):
 """Encode and send an audio chunk to the server."""
 if not self.is_connected:
 return

 event = {
 "event_id": f"event_{int(time.time() * 1000)}",
 "type": "input_audio_buffer.append",
 "audio": base64.b64encode(audio_data).decode()
 }
 await self.ws.send(json.dumps(event))

 async def send_image_frame(self, image_bytes: bytes, *, event_id: str | None = None):
 # Send an image frame to the server.
 if not self.is_connected:
 return

 if not image_bytes:
 raise ValueError("image_bytes cannot be empty.")

 # Encode to Base64
 image_b64 = base64.b64encode(image_bytes).decode()

 event = {
 "event_id": event_id or f"event_{int(time.time() * 1000)}",
 "type": "input_image_buffer.append",
 "image": image_b64,
 }

 await self.ws.send(json.dumps(event))

 def _audio_player_task(self):
 stream = self.pyaudio_instance.open(
 format=self.output_format,
 channels=self.output_channels,
 rate=self.output_rate,
 output=True,
 frames_per_buffer=self.output_chunk,
 )
 try:
 while self.is_connected or not self.audio_playback_queue.empty():
 try:
 audio_chunk = self.audio_playback_queue.get(timeout=0.1)
 if audio_chunk is None: # Termination signal
 break
 stream.write(audio_chunk)
 self.audio_playback_queue.task_done()
 except queue.Empty:
 continue
 finally:
 stream.stop_stream()
 stream.close()

 def start_audio_player(self):
 """Start the audio player thread (only when audio output is enabled)."""
 if not self.audio_enabled:
 return
 if self.audio_player_thread is None or not self.audio_player_thread.is_alive():
 self.audio_player_thread = threading.Thread(target=self._audio_player_task, daemon=True)
 self.audio_player_thread.start()

 async def handle_server_messages(self, on_text_received):
 """Handle incoming messages from the server in a loop."""
 try:
 async for message in self.ws:
 event = json.loads(message)
 event_type = event.get("type")
 if event_type == "response.audio.delta" and self.audio_enabled:
 audio_b64 = event.get("delta", "")
 if audio_b64:
 audio_data = base64.b64decode(audio_b64)
 self.audio_playback_queue.put(audio_data)

 elif event_type == "response.done":
 print("\n[INFO] Response round complete.")
 usage = event.get("response", {}).get("usage", {})
 if usage:
 print(f"[INFO] token usage: {json.dumps(usage, indent=2, ensure_ascii=False)}")
 # Process source language recognition results (requires enabling input_audio_transcription.model)
 # elif event_type == "conversation.item.input_audio_transcription.text":
 # stash = event.get("stash", "") # Pending recognition text
 # print(f"[Recognizing] {stash}")
 # elif event_type == "conversation.item.input_audio_transcription.completed":
 # transcript = event.get("transcript", "") # Complete recognition result
 # print(f"[Source language] {transcript}")
 elif event_type == "response.text.text":
 # Streaming translated text in text-only modality
 text = event.get("text", "")
 stash = event.get("stash", "")
 print(f"\r[Translating] {text}{stash}", end="", flush=True)
 elif event_type == "response.audio_transcript.done":
 print("\n[INFO] Translation complete.")
 text = event.get("transcript", "")
 if text:
 print(f"[INFO] Translated text: {text}")
 elif event_type == "response.text.done":
 print("\n[INFO] Translation complete.")
 text = event.get("text", "")
 if text:
 print(f"[INFO] Translated text: {text}")

 except websockets.exceptions.ConnectionClosed as e:
 print(f"[WARNING] Connection closed: {e}")
 self.is_connected = False
 except Exception as e:
 print(f"[ERROR] An unexpected error occurred while processing messages: {e}")
 traceback.print_exc()
 self.is_connected = False

 async def start_microphone_streaming(self):
 """Capture audio from the microphone and stream it to the server."""
 stream = self.pyaudio_instance.open(
 format=self.input_format,
 channels=self.input_channels,
 rate=self.input_rate,
 input=True,
 frames_per_buffer=self.input_chunk
 )
 print("Microphone is on. Start speaking...")
 try:
 while self.is_connected:
 audio_chunk = await asyncio.get_event_loop().run_in_executor(
 None, stream.read, self.input_chunk
 )
 await self.send_audio_chunk(audio_chunk)
 finally:
 stream.stop_stream()
 stream.close()

 async def close(self):
 """Gracefully close the connection and release resources."""
 self.is_connected = False
 if self.ws:
 await self.ws.close()
 print("WebSocket connection closed.")

 if self.audio_player_thread:
 self.audio_playback_queue.put(None) # Send termination signal
 self.audio_player_thread.join(timeout=1)
 print("Audio player thread stopped.")

 self.pyaudio_instance.terminate()
 print("PyAudio instance released.")

``` 
### [​ ](#interact-with-the-model) Interact with the model

In the same directory, create a file named `main.py` with the following code:
 main.py

 Copy ```\nimport os
import asyncio
from livetranslate_client import LiveTranslateClient

def print_banner():
 print("=" * 60)
 print(" Powered by Qwen qwen3.5-livetranslate-flash-realtime")
 print("=" * 60 + "\n")

def get_user_config():
 """Get user configuration."""
 print("Select a mode:")
 print("1. Voice + Text [Default] | 2. Text Only")
 mode_choice = input("Enter your choice (press Enter for Voice + Text): ").strip()
 audio_enabled = (mode_choice != "2")

 if audio_enabled:
 lang_map = {
 "1": "en", "2": "zh", "3": "ru", "4": "fr", "5": "de", "6": "pt",
 "7": "es", "8": "it", "9": "ko", "10": "ja", "11": "yue"
 }
 print("Select the target language (Voice + Text mode):")
 print("1. English | 2. Chinese | 3. Russian | 4. French | 5. German | 6. Portuguese | 7. Spanish | 8. Italian | 9. Korean | 10. Japanese | 11. Cantonese")
 else:
 lang_map = {
 "1": "en", "2": "zh", "3": "ru", "4": "fr", "5": "de", "6": "pt", "7": "es", "8": "it",
 "9": "id", "10": "ko", "11": "ja", "12": "vi", "13": "th", "14": "ar",
 "15": "yue", "16": "hi", "17": "el", "18": "tr"
 }
 print("Select the target language (Text Only mode):")
 print("1. English | 2. Chinese | 3. Russian | 4. French | 5. German | 6. Portuguese | 7. Spanish | 8. Italian | 9. Indonesian | 10. Korean | 11. Japanese | 12. Vietnamese | 13. Thai | 14. Arabic | 15. Cantonese | 16. Hindi | 17. Greek | 18. Turkish")

 choice = input("Enter your choice (defaults to the first option): ").strip()
 target_language = lang_map.get(choice, next(iter(lang_map.values())))

 return target_language, audio_enabled

async def main():
 """Main program entry point."""
 print_banner()

 api_key = os.environ.get("DASHSCOPE_API_KEY")
 if not api_key:
 print("[ERROR] Please set the DASHSCOPE_API_KEY environment variable.")
 print(" For example: export DASHSCOPE_API_KEY=&#x27;your_api_key_here&#x27;")
 return

 target_language, audio_enabled = get_user_config()
 print("\nConfiguration complete:")
 print(f" - Target language: {target_language}")
 if not audio_enabled:
 print(" - Output mode: Text Only")

 client = LiveTranslateClient(api_key=api_key, target_language=target_language, audio_enabled=audio_enabled)

 # Define the callback function.
 def on_translation_text(text):
 print(text, end="", flush=True)

 try:
 print("Connecting to the translation service...")
 await client.connect()

 # Start audio playback based on the mode.
 client.start_audio_player()

 print("\n" + "-" * 60)
 print("Connection successful! Speak into the microphone.")
 print("The program will translate your speech in real time and play the translated audio. Press Ctrl+C to exit.")
 print("-" * 60 + "\n")

 # Run message handling and microphone recording concurrently.
 message_handler = asyncio.create_task(client.handle_server_messages(on_translation_text))
 tasks = [message_handler]
 # Capture audio from the microphone for translation, regardless of whether audio output is enabled.
 microphone_streamer = asyncio.create_task(client.start_microphone_streaming())
 tasks.append(microphone_streamer)

 await asyncio.gather(*tasks)

 except KeyboardInterrupt:
 print("\n\nUser interrupted. Exiting...")
 except Exception as e:
 print(f"\nA critical error occurred: {e}")
 finally:
 print("\nCleaning up resources...")
 await client.close()
 print("Program exited.")

if __name__ == "__main__":
 asyncio.run(main())

``` 
Run `main.py` and speak into your microphone. The model outputs translated audio and text in real time. The system automatically detects speech and sends it to the server.
## [​ ](#how-to-use) How to use

### [​ ](#1-configure-the-connection) 1. Configure the connection

The qwen3.5-livetranslate-flash-realtime model uses the WebSocket protocol. The connection requires the following parameters:
ParameterDescriptionendpoint`wss://dashscope-intl.aliyuncs.com/api-ws/v1/realtime`query parameterThe `model` query parameter must be set to the model name. Example: `?model=qwen3.5-livetranslate-flash-realtime`message headerUse a Bearer Token for authentication: `Authorization: Bearer DASHSCOPE_API_KEY` 
 DASHSCOPE_API_KEY is your API key from Qwen Cloud. 
Sample Python code for establishing a connection:
 Python sample code for WebSocket connection

 Copy ```\n# pip install websocket-client
import json
import websocket
import os

API_KEY=os.getenv("DASHSCOPE_API_KEY")
API_URL = "wss://dashscope-intl.aliyuncs.com/api-ws/v1/realtime?model=qwen3.5-livetranslate-flash-realtime"

headers = [
 "Authorization: Bearer " + API_KEY
]

def on_open(ws):
 print(f"Connected to server: {API_URL}")
def on_message(ws, message):
 data = json.loads(message)
 print("Received event:", json.dumps(data, indent=2))
def on_error(ws, error):
 print("Error:", error)

ws = websocket.WebSocketApp(
 API_URL,
 header=headers,
 on_open=on_open,
 on_message=on_message,
 on_error=on_error
)

ws.run_forever()

``` 
### [​ ](#2-configure-language-modality-and-voice) 2. Configure language, modality, and voice

Send the [session.update](/api-reference/speech-translation/livetranslate-realtime/client-events#session-update) client event with the following parameters:

- 
**Language**


**Source language:** Configure using the `session.input_audio_transcription.language` parameter.
 The default value is `en` (English). 


- 
**Target language:** Configure using the `session.translation.language` parameter.
 The default value is `en` (English). 


See [Supported languages](#supported-languages).

- 
**Output source language recognition results**
To enable this feature, set the `session.input_audio_transcription.model` parameter. When set to `qwen3-asr-flash-realtime`, the server returns both the translation and the speech recognition result (original text) for the input audio.
When this feature is enabled, the server returns the following events:

`conversation.item.input_audio_transcription.text`: Streams the recognition results.

- `conversation.item.input_audio_transcription.completed`: Returns the final result after the recognition is complete.


- 
**Output modality**
Set the `session.modalities` parameter to `["text"]` (text only) or `["text","audio"]` (text and audio).


- 
**Voice**
Configure using the `session.voice` parameter. See [Supported voices](#supported-voices).


- 
**Hotword**
Configure hotwords using the `session.translation.corpus.phrases` parameter. Hotwords are key-value pairs that map source terms to target translations, improving accuracy for specific terms.
Example: Map `"artificial intelligence"` to `"Artificial Intelligence"`.


- 
**Voice cloning**
Configure using the `session.enable_voice_clone`, `session.voice_clone_options.frequency`, and `session.voice` parameters. Supports three modes: pre-cloned voice profile (`frequency`: `never`), server-side clone once at session start (`once`), or real-time clone before each response (`always`). See [Voice cloning](#voice-cloning).


### [​ ](#3-input-audio-and-images) 3. Input audio and images

Send Base64-encoded audio and image data using the [input_audio_buffer.append](/api-reference/speech-translation/livetranslate-realtime/client-events#input-audio-buffer-append) and [input_image_buffer.append](/api-reference/speech-translation/livetranslate-realtime/client-events#input-image-buffer-append) events. Audio input is required; image input is optional.
 Images can be from a local file or captured in real time from a video stream. The server automatically detects speech boundaries and triggers the model response. 
### [​ ](#4-receive-the-model-response) 4. Receive the model response

The model responds when the server detects the end of speech. The response format depends on the output modality.

- 
**Text-only output**
The server streams incremental translated text (including confirmed text and tentative predicted text) through [response.text.text](/api-reference/speech-translation/livetranslate-realtime/server-events#response-text-text) events; upon completion, the full translated text is returned in a [response.text.done](/api-reference/speech-translation/livetranslate-realtime/server-events#response-text-done) event.


- 
**Text and audio output**

**Text**: The server streams incremental translated text through [response.audio_transcript.text](/api-reference/speech-translation/livetranslate-realtime/server-events#response-audio-transcript-text) events; upon completion, the full translated text is returned in a [response.audio_transcript.done](/api-reference/speech-translation/livetranslate-realtime/server-events#response-audio-transcript-done) event.

- **Audio**: Incremental, Base64-encoded audio data is returned in [response.audio.delta](/api-reference/speech-translation/livetranslate-realtime/server-events#response-audio-delta) events.


 The real-time translation model uses the `response.text.text` event for incremental text delivery, which differs from the `response.text.delta` event used by Omni (full-duplex voice conversation) models. These events have different field structures and semantics — do not use them interchangeably. 
### [​ ](#5-end-the-session) 5. End the session

After sending all audio, send a [session.finish](/api-reference/speech-translation/livetranslate-realtime/client-events#session-finish) event, then wait for the server to return a `session.finished` event before closing the WebSocket connection.
 If you close the WebSocket without sending `session.finish`, the server&#x27;s VAD cannot detect the end of the final speech segment. This causes translation results for that segment to be lost entirely, and the connection may hang indefinitely. Always send this event before disconnecting. 
## [​ ](#voice-cloning) Voice cloning

The model clones the speaker&#x27;s voice from the input audio and uses the cloned voice for translated output, so the translation sounds like the speaker delivering it in another language. Use a pre-cloned voice profile, or let the server clone the voice in real time. This is useful in scenarios where preserving the speaker&#x27;s voice matters, such as conference interpreting, live streaming, and video dubbing.
Set the following parameters in [session.update](/api-reference/speech-translation/livetranslate-realtime/client-events#session-update) to enable voice cloning:

- `session.enable_voice_clone`: Set to `true` to enable voice cloning.

- `session.voice_clone_options.frequency`: Controls when voice cloning occurs. Accepted values:

`never`: Does not clone on the server. Uses a pre-cloned voice profile instead. Set `session.voice` to your custom cloned voice ID.

- `once`: Clones the voice from the input audio once at session start, then reuses it for all subsequent output. Best for single-speaker scenarios. Set `session.voice` to `default`.

- `always`: Clones the voice before each response, dynamically adapting to speaker changes. Best for multi-speaker conversations. Set `session.voice` to `default`.


- `session.voice`: Specifies the output voice. The value depends on the `frequency` setting:

Set to `default`: Use with `frequency` set to `once` or `always`. The server clones the speaker&#x27;s voice from the input audio. A default voice is used until cloning completes.

- Set to a custom cloned voice ID (for example, `qwen-translate-vc-xxx-yyy-zzz`): Use with `frequency` set to `never`. You must prepare the voice in advance using the Voice Cloning API with `targetModel` set to `qwen3.5-livetranslate-flash-realtime`.


 When `frequency` is set to `once` or `always`, the `voice` parameter must be set to `default`. Any other value causes the server to return an error. 
### [​ ](#voice-cloning-configuration-examples) Voice cloning configuration examples

**Pre-cloned voice profile** (consistent quality; recommended when a stable voice identity is required):
Copy ```\n{
 "type": "session.update",
 "session": {
 "modalities": ["text","audio"],
 "voice": "qwen-translate-vc-xxx-yyy-zzz",
 "translation": {
 "language": "en"
 },
 "enable_voice_clone": true,
 "voice_clone_options": {
 "frequency": "never"
 }
 }
}

``` 
**Server-side cloning, once per session** (best for single-speaker scenarios):
Copy ```\n{
 "type": "session.update",
 "session": {
 "modalities": ["text","audio"],
 "voice": "default",
 "translation": {
 "language": "en"
 },
 "enable_voice_clone": true,
 "voice_clone_options": {
 "frequency": "once"
 }
 }
}

``` 
**Server-side cloning, every response** (best for multi-speaker conversations):
Copy ```\n{
 "type": "session.update",
 "session": {
 "modalities": ["text","audio"],
 "voice": "default",
 "translation": {
 "language": "en"
 },
 "enable_voice_clone": true,
 "voice_clone_options": {
 "frequency": "always"
 }
 }
}

``` 
## [​ ](#interaction-flow) Interaction flow

Real-time speech translation follows an event-driven WebSocket model. The server automatically detects speech boundaries and responds.
LifecycleClient eventServer eventSession initializationsession.update (Session configuration)session.created (Session created), session.updated (Session configuration updated)User audio inputinput_audio_buffer.append (Append audio to the buffer)NoneServer audio outputNoneresponse.created (Signals that the server starts generating a response), response.output_item.added (Signals that a new output item is available), response.content_part.added (Signals that a new content part has been added to the assistant message), response.text.text (Incremental translated text in text-only modality), response.audio_transcript.text (Incremental translated text in audio+text modality), response.audio.delta (Contains an incremental chunk of the synthesized audio), response.text.done (Translation text complete in text-only modality), response.audio_transcript.done (Translation text complete in audio+text modality), response.audio.done (Signals that the synthesized audio is complete), response.content_part.done (Signals that a text or audio content part for the assistant message is complete), response.output_item.done (Signals that the entire output item for the assistant message is complete), response.done (Signals that the entire response is complete) 
## [​ ](#improve-translation-with-images) Improve translation with images

The qwen3.5-livetranslate-flash-realtime model uses image input to improve audio translation, helping disambiguate homonyms and recognize uncommon proper nouns. Send no more than 2 images per second.
Download the following sample images: [medical mask.png](https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/en-US/20250923/tjpeys/%E5%8F%A3%E7%BD%A9.png) and [masquerade mask.png](https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/en-US/20250923/ifqttq/%E9%9D%A2%E5%85%B7.png)
Download the following code to the same directory as `livetranslate_client.py` and run it. Say `"What is mask?"` into your microphone. The model uses the provided image to disambiguate the word "mask." For example, using the `medical mask.png` file translates the phrase as "What is a medical mask?", while using the `masquerade mask.png` file translates it as "What is a masquerade mask?".
Copy ```\nimport os
import time
import json
import asyncio
import contextlib
import functools

from livetranslate_client import LiveTranslateClient

IMAGE_PATH = "medical mask.png"
# IMAGE_PATH = "masquerade mask.png"

def print_banner():
 print("=" * 60)
 print(" Powered by Qwen qwen3.5-livetranslate-flash-realtime — single-turn interaction example (mask)")
 print("=" * 60 + "\n")

async def stream_microphone_once(client: LiveTranslateClient, image_bytes: bytes):
 pa = client.pyaudio_instance
 stream = pa.open(
 format=client.input_format,
 channels=client.input_channels,
 rate=client.input_rate,
 input=True,
 frames_per_buffer=client.input_chunk,
 )
 print(f"[INFO] Recording started. Please speak...")
 loop = asyncio.get_event_loop()
 last_img_time = 0.0
 frame_interval = 0.5 # 2 fps
 try:
 while client.is_connected:
 data = await loop.run_in_executor(None, stream.read, client.input_chunk)
 await client.send_audio_chunk(data)

 # Append an image frame every 0.5 seconds
 now = time.time()
 if now - last_img_time >= frame_interval:
 await client.send_image_frame(image_bytes)
 last_img_time = now
 finally:
 stream.stop_stream()
 stream.close()

async def main():
 print_banner()
 api_key = os.environ.get("DASHSCOPE_API_KEY")
 if not api_key:
 print("[ERROR] Please set the DASHSCOPE_API_KEY environment variable.")
 return

 client = LiveTranslateClient(api_key=api_key, target_language="zh", audio_enabled=True)

 def on_text(text: str):
 print(text, end="", flush=True)

 try:
 await client.connect()
 client.start_audio_player()
 message_task = asyncio.create_task(client.handle_server_messages(on_text))
 with open(IMAGE_PATH, "rb") as f:
 img_bytes = f.read()
 await stream_microphone_once(client, img_bytes)
 await asyncio.sleep(15)
 finally:
 await client.close()
 if not message_task.done():
 message_task.cancel()
 with contextlib.suppress(asyncio.CancelledError):
 await message_task

if __name__ == "__main__":
 asyncio.run(main())

``` 
## [​ ](#billing) Billing

**Qwen3.5-LiveTranslate-Flash-Realtime**

- **Audio**: 7 tokens per second of input audio; 12.5 tokens per second of output audio.

- **Image**: Every 32x32 pixels consumes 0.5 tokens.

- **Text**: When source language speech recognition is enabled, the service returns a transcript of the input audio in addition to the translation. This transcript is billed as output text tokens.


**Qwen3-LiveTranslate-Flash-Realtime**

- **Audio**: Each second of audio input or output consumes 12.5 tokens.

- **Image**: Every 28x28 pixels consumes 0.5 tokens.

- **Text**: When source language speech recognition is enabled, the service returns a transcript of the input audio in addition to the translation. This transcript is billed as output text tokens.


For pricing, see [Model list](/developer-guides/getting-started/model-selection).
## [​ ](#supported-languages) Supported languages

Use the following language codes to specify the source and target languages.
 Some target languages only support text. The legacy model qwen3-livetranslate-flash-realtime supports only the following 18 languages: en, zh, ru, fr, de, pt, es, it, id, ko, ja, vi, th, ar, yue, hi, el, tr. 
Language codeLanguageOutputzhChineseAudio + textenEnglishAudio + textarArabicAudio + textdeGermanAudio + textfrFrenchAudio + textesSpanishAudio + textptPortugueseAudio + textidIndonesianAudio + textitItalianAudio + textkoKoreanAudio + textruRussianAudio + textthThaiAudio + textviVietnameseAudio + textjaJapaneseAudio + texttrTurkishAudio + texthiHindiAudio + textmsMalayAudio + textnlDutchAudio + texturUrduAudio + textnbNorwegian BokmålAudio + textsvSwedishAudio + textdaDanishAudio + textheHebrewAudio + textfiFinnishAudio + textplPolishAudio + textisIcelandicAudio + textcsCzechAudio + textfilFilipinoAudio + textfaPersianAudio + textyueCantoneseTextelGreekTextafAfrikaansTextastAsturianTextbeBelarusianTextbgBulgarianTextbnBengaliTextbsBosnianTextcaCatalanTextcebCebuanoTextetEstonianTextglGalicianTextguGujaratiTexthrCroatianTexthuHungarianTextjvJavaneseTextkkKazakhTextknKannadaTextkyKyrgyzTextlvLatvianTextmkMacedonianTextmlMalayalamTextmrMarathiTextpaPunjabiTextroRomanianTextskSlovakTextslSlovenianTextswSwahiliTexttgTajikTextazAzerbaijaniTextukUkrainianText 
## [​ ](#supported-voices) Supported voices

For supported voices and the corresponding `voice` parameter values, see the [API reference](#api-reference).
## [​ ](#api-reference) API reference


- [Client events](/api-reference/speech-translation/livetranslate-realtime/client-events)

- [Server events](/api-reference/speech-translation/livetranslate-realtime/server-events)

- [Python SDK](/api-reference/speech-translation/livetranslate-realtime/python-sdk)

- [Java SDK](/api-reference/speech-translation/livetranslate-realtime/java-sdk)


 [Previous ](/developer-guides/speech/multimodal-speech)[Audio and video file translation 18-language translation Next ](/developer-guides/speech/file-translation)
