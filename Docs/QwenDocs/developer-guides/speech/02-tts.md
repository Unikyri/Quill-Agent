# Non-real-time speech synthesis

> **Source:** https://docs.qwencloud.com/developer-guides/speech/tts

Non-realtime speech synthesis with Qwen3-TTS

 Copy page Non-real-time speech synthesis converts text to speech (TTS) through an HTTP API. It suits latency-tolerant scenarios such as audiobook production, e-learning narration, and content production.
## [​ ](#overview) Overview

Convert complete text to speech files through an HTTP API. Two output modes are available: non-streaming and streaming.

- Non-streaming mode returns an audio file URL that expires after 24 hours. Streaming mode returns PCM audio data in chunks.

- Supports multiple languages, including Chinese dialects.

- Supports [voice cloning](/developer-guides/speech/voice-cloning) and [voice design](/developer-guides/speech/voice-design) for custom voice creation.

- Supports [instruction control](#instruction-control), which lets you control speech expressiveness through natural-language instructions.


For low-latency streaming synthesis, see [Real-time speech synthesis](/developer-guides/speech/realtime-streaming). To choose a model, see [Text-to-speech models](/developer-guides/speech/tts-models).
## [​ ](#prerequisites) Prerequisites


- [Get an API key](/api-reference/preparation/api-key) and [set it as an environment variable](/api-reference/preparation/export-api-key-env).

- To use the SDK, [install it](/api-reference/preparation/install-sdk). The Java SDK requires version 2.21.9+. The Python SDK requires version 1.24.6+.


## [​ ](#quick-start) Quick start

The following examples demonstrate how to synthesize speech with the Qwen-TTS model family. For detailed parameter descriptions, see the [API reference](#api-reference).
#### [​ ](#use-system-voice) Use system voice

Use a [built-in voice](#built-in-voices) for speech synthesis.
**Non-streaming output**
Use the returned `url` to retrieve the synthesized audio. The URL is valid for 24 hours.
You must import the Gson dependency for Java. If you use Maven or Gradle, add the dependency as follows:
- Maven 
- Gradle 

 Add the following content to `pom.xml`:Copy ```\n<!-- https://mvnrepository.com/artifact/com.google.code.gson/gson -->
<dependency>
 <groupId>com.google.code.gson</groupId>
 <artifactId>gson</artifactId>
 <version>2.13.1</version>
</dependency>

``` Add the following content to `build.gradle`:Copy ```\n// https://mvnrepository.com/artifact/com.google.code.gson/gson
implementation("com.google.code.gson:gson:2.13.1")

``` 
Python Java cURL Copy ```\nimport os
import dashscope

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

text = "Today is a wonderful day to build something people love!"
# To use the SpeechSynthesizer interface: dashscope.audio.qwen_tts.SpeechSynthesizer.call(...)
response = dashscope.MultiModalConversation.call(
 # Replace the model with qwen3-tts-instruct-flash to use instruction control.
 model="qwen3-tts-flash",
 # If you have not configured an environment variable, replace the following line with your API key: api_key = "sk-xxx"
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 text=text,
 voice="Cherry",
 language_type="English", # We recommend that this matches the text language to ensure correct pronunciation and natural intonation.
 # To use instruction control, uncomment the following lines and replace the model with qwen3-tts-instruct-flash.
 # instructions=&#x27;Speak quickly with a noticeable rising intonation, suitable for introducing fashion products.&#x27;,
 # optimize_instructions=True,
 stream=False
)
print(response)

``` 
**Streaming output**
Stream audio data in Base64 format. The last packet contains the URL for the complete audio file.
Python Java cURL Copy ```\n# coding=utf-8
#
# Installation instructions for pyaudio:
# APPLE Mac OS X
# brew install portaudio
# pip install pyaudio
# Debian/Ubuntu
# sudo apt-get install python-pyaudio python3-pyaudio
# or
# pip install pyaudio
# CentOS
# sudo yum install -y portaudio portaudio-devel && pip install pyaudio
# Microsoft Windows
# python -m pip install pyaudio

import os
import dashscope
import pyaudio
import time
import base64
import numpy as np

dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

p = pyaudio.PyAudio()
# Create an audio stream
stream = p.open(format=pyaudio.paInt16,
 channels=1,
 rate=24000,
 output=True)


text = "Today is a wonderful day to build something people love!"
response = dashscope.MultiModalConversation.call(
 # If you have not configured an environment variable, replace the following line with your API key: api_key = "sk-xxx"
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 # Replace the model with qwen3-tts-instruct-flash to use instruction control.
 model="qwen3-tts-flash",
 text=text,
 voice="Cherry",
 language_type="English", # We recommend that this matches the text language to ensure correct pronunciation and natural intonation.
 # To use instruction control, uncomment the following lines and replace the model with qwen3-tts-instruct-flash.
 # instructions=&#x27;Speak quickly with a noticeable rising intonation, suitable for introducing fashion products.&#x27;,
 # optimize_instructions=True,
 stream=True
)

for chunk in response:
 if chunk.output is not None:
 audio = chunk.output.audio
 if audio.data is not None:
 wav_bytes = base64.b64decode(audio.data)
 audio_np = np.frombuffer(wav_bytes, dtype=np.int16)
 # Play the audio data directly
 stream.write(audio_np.tobytes())
 if chunk.output.finish_reason == "stop":
 print("finish at: {} ", chunk.output.audio.expires_at)
time.sleep(0.8)
# Clean up resources
stream.stop_stream()
stream.close()
p.terminate()

``` 
#### [​ ](#use-cloned-voice) Use cloned voice

Voice cloning does not provide preview audio. Apply the cloned voice to speech synthesis to evaluate the result.
These examples adapt the non-streaming output code, replacing the `voice` parameter with a cloned voice.

- **Key principle**: The model used for voice cloning (`target_model`) must match the model used for speech synthesis (`model`). Otherwise, synthesis fails.

- This example uses the local audio file `voice.mp3` for voice cloning. Replace this path when running the code.


Add the Gson dependency for Java. If you use Maven or Gradle, add the dependency as follows:
- Maven 
- Gradle 

 Add the following content to your `pom.xml`:Copy ```\n<!-- https://mvnrepository.com/artifact/com.google.code.gson/gson -->
<dependency>
 <groupId>com.google.code.gson</groupId>
 <artifactId>gson</artifactId>
 <version>2.13.1</version>
</dependency>

``` Add the following content to your `build.gradle`:Copy ```\n// https://mvnrepository.com/artifact/com.google.code.gson/gson
implementation("com.google.code.gson:gson:2.13.1")

``` 
 When using a custom voice generated by voice cloning for speech synthesis, set the voice as follows:Copy ```\nMultiModalConversationParam param = MultiModalConversationParam.builder()
 .parameter("voice", "your_voice") // Replace the voice parameter with the custom voice generated by cloning
 .build();

``` 
Python Java Copy ```\nimport os
import requests
import base64
import pathlib
import dashscope

# ======= Constant configuration =======
DEFAULT_TARGET_MODEL = "qwen3-tts-vc-2026-01-22" # Use the same model for voice cloning and speech synthesis
DEFAULT_PREFERRED_NAME = "guanyu"
DEFAULT_AUDIO_MIME_TYPE = "audio/mpeg"
VOICE_FILE_PATH = "voice.mp3" # Relative path to the local audio file used for voice cloning


def create_voice(file_path: str,
 target_model: str = DEFAULT_TARGET_MODEL,
 preferred_name: str = DEFAULT_PREFERRED_NAME,
 audio_mime_type: str = DEFAULT_AUDIO_MIME_TYPE) -> str:
 """
 Create a voice and return the voice parameter.
 """
 # If you haven&#x27;t configured an environment variable, replace the following line with: api_key = "sk-xxx"
 api_key = os.getenv("DASHSCOPE_API_KEY")

 file_path_obj = pathlib.Path(file_path)
 if not file_path_obj.exists():
 raise FileNotFoundError(f"Audio file does not exist: {file_path}")

 base64_str = base64.b64encode(file_path_obj.read_bytes()).decode()
 data_uri = f"data:{audio_mime_type};base64,{base64_str}"

 url = "https://dashscope-intl.aliyuncs.com/api/v1/services/audio/tts/customization"
 payload = {
 "model": "qwen-voice-enrollment", # Do not change this value
 "input": {
 "action": "create",
 "target_model": target_model,
 "preferred_name": preferred_name,
 "audio": {"data": data_uri}
 }
 }
 headers = {
 "Authorization": f"Bearer {api_key}",
 "Content-Type": "application/json"
 }

 resp = requests.post(url, json=payload, headers=headers)
 if resp.status_code != 200:
 raise RuntimeError(f"Failed to create voice: {resp.status_code}, {resp.text}")

 try:
 return resp.json()["output"]["voice"]
 except (KeyError, ValueError) as e:
 raise RuntimeError(f"Failed to parse voice response: {e}")


if __name__ == &#x27;__main__&#x27;:
 dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

 text = "How&#x27;s the weather today?"
 # SpeechSynthesizer interface usage: dashscope.audio.qwen_tts.SpeechSynthesizer.call(...)
 response = dashscope.MultiModalConversation.call(
 model=DEFAULT_TARGET_MODEL,
 # If you haven&#x27;t configured an environment variable, replace the following line with: api_key = "sk-xxx"
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 text=text,
 voice=create_voice(VOICE_FILE_PATH), # Replace the voice parameter with the custom voice generated by cloning
 stream=False
 )
 print(response)

``` 
#### [​ ](#use-designed-voice) Use designed voice

Voice design returns preview audio. Listen to the preview to confirm it meets your expectations before using it for synthesis to reduce costs.
 1 Generate a custom voice and preview the result

If you are satisfied with the result, proceed to the next step. Otherwise, generate it again.You need to import the Gson dependency for Java. If you are using Maven or Gradle, add the dependency as follows:- Maven 
- Gradle 

 Add the following content to `pom.xml`:Copy ```\n<!-- https://mvnrepository.com/artifact/com.google.code.gson/gson -->
<dependency>
 <groupId>com.google.code.gson</groupId>
 <artifactId>gson</artifactId>
 <version>2.13.1</version>
</dependency>

``` Add the following content to `build.gradle`:Copy ```\n// https://mvnrepository.com/artifact/com.google.code.gson/gson
implementation("com.google.code.gson:gson:2.13.1")

``` When using a custom voice generated by voice design for speech synthesis, you must set the voice as follows:Copy ```\nMultiModalConversationParam param = MultiModalConversationParam.builder()
 .parameter("voice", "your_voice") // Replace the voice parameter with the custom voice generated by voice design
 .build();

``` Python Java Copy ```\nimport requests
import base64
import os

def create_voice_and_play():
 # If the environment variable is not set, replace the following line with your API key: api_key = "sk-xxx"
 api_key = os.getenv("DASHSCOPE_API_KEY")

 if not api_key:
 print("Error: DASHSCOPE_API_KEY environment variable not found. Please set the API key first.")
 return None, None, None

 # Prepare request data
 headers = {
 "Authorization": f"Bearer {api_key}",
 "Content-Type": "application/json"
 }

 data = {
 "model": "qwen-voice-design",
 "input": {
 "action": "create",
 "target_model": "qwen3-tts-vd-2026-01-26",
 "voice_prompt": "A composed middle-aged male announcer with a deep, rich and magnetic voice, a steady speaking speed and clear articulation, is suitable for news broadcasting or documentary commentary.",
 "preview_text": "Dear listeners, hello everyone. Welcome to the evening news.",
 "preferred_name": "announcer",
 "language": "en"
 },
 "parameters": {
 "sample_rate": 24000,
 "response_format": "wav"
 }
 }

 url = "https://dashscope-intl.aliyuncs.com/api/v1/services/audio/tts/customization"

 try:
 # Send the request
 response = requests.post(
 url,
 headers=headers,
 json=data,
 timeout=60 # Add a timeout setting
 )

 if response.status_code == 200:
 result = response.json()

 # Get the voice name
 voice_name = result["output"]["voice"]
 print(f"Voice name: {voice_name}")

 # Get the preview audio data
 base64_audio = result["output"]["preview_audio"]["data"]

 # Decode the Base64 audio data
 audio_bytes = base64.b64decode(base64_audio)

 # Save the audio file locally
 filename = f"{voice_name}_preview.wav"

 # Write the audio data to a local file
 with open(filename, &#x27;wb&#x27;) as f:
 f.write(audio_bytes)

 print(f"Audio saved to local file: {filename}")
 print(f"File path: {os.path.abspath(filename)}")

 return voice_name, audio_bytes, filename
 else:
 print(f"Request failed with status code: {response.status_code}")
 print(f"Response content: {response.text}")
 return None, None, None

 except requests.exceptions.RequestException as e:
 print(f"A network request error occurred: {e}")
 return None, None, None
 except KeyError as e:
 print(f"Response data format error, missing required field: {e}")
 print(f"Response content: {response.text if &#x27;response&#x27; in locals() else &#x27;No response&#x27;}")
 return None, None, None
 except Exception as e:
 print(f"An unknown error occurred: {e}")
 return None, None, None

if __name__ == "__main__":
 print("Starting to create voice...")
 voice_name, audio_data, saved_filename = create_voice_and_play()

 if voice_name:
 print(f"\nSuccessfully created voice &#x27;{voice_name}&#x27;")
 print(f"Audio file saved as: &#x27;{saved_filename}&#x27;")
 print(f"File size: {os.path.getsize(saved_filename)} bytes")
 else:
 print("\nVoice creation failed")

``` 2 Use the custom voice for speech synthesis

Use the custom voice generated in the previous step for non-streaming speech synthesis.This example adapts the non-streaming output code, replacing the `voice` parameter with the custom voice generated by voice design. For streaming synthesis, see [Quick start](#quick-start).**Key principle**: The model used for voice design (`target_model`) must be the same as the model used for subsequent speech synthesis (`model`). Otherwise, the synthesis will fail.Python Java Copy ```\nimport os
import dashscope


if __name__ == &#x27;__main__&#x27;:
 dashscope.base_http_api_url = &#x27;https://dashscope-intl.aliyuncs.com/api/v1&#x27;

 text = "How is the weather today?"
 # How to use the SpeechSynthesizer interface: dashscope.audio.qwen_tts.SpeechSynthesizer.call(...)
 response = dashscope.MultiModalConversation.call(
 model="qwen3-tts-vd-2026-01-26",
 # If the environment variable is not set, replace the following line with your API key: api_key = "sk-xxx"
 api_key=os.getenv("DASHSCOPE_API_KEY"),
 text=text,
 voice="myvoice", # Replace the voice parameter with the custom voice generated by voice design
 stream=False
 )
 print(response)

``` 
## [​ ](#instruction-control) Instruction control

Control pitch, speed, emotion, and timbre using natural language instructions instead of audio parameters.
**Supported models**: Qwen3-TTS-Instruct-Flash series only.
**Usage**: Specify instructions in the `instructions` parameter. Example: "Fast-paced with rising intonation, suitable for fashion products."
**Supported languages**: Chinese and English only.
**Length limit**: Maximum 1600 tokens.
**Scenarios**:

- Audiobook and radio drama voice-overs

- Advertising and promotional video voice-overs

- Game role and animation voice-overs

- Emotionally intelligent voice assistants

- Documentary and news broadcasting


**Writing high-quality sound descriptions**
**Core principles**

- Be specific: Use descriptive words such as "deep," "crisp," or "fast-paced." Avoid vague words such as "nice" or "normal."

- Be multi-dimensional: Combine multiple dimensions such as pitch, speed, and emotion. Single-dimension descriptions such as "high-pitched" are too broad.

- Be objective: Focus on physical and perceptual features, not personal preferences. Use "high-pitched and energetic" instead of "my favorite sound."

- Be original: Describe sound qualities instead of requesting imitation of specific people. The model does not support direct imitation.

- Be concise: Ensure every word serves a purpose. Avoid repetitive synonyms or meaningless intensifiers.


**Dimension description reference**: You can combine multiple dimensions to create richer audio effects.
DimensionExamplePitchHigh, medium, low, high-pitched, low-pitchedSpeedFast, medium, slow, fast-paced, slow-pacedEmotionCheerful, calm, gentle, serious, lively, composed, soothingCharacteristicsMagnetic, crisp, hoarse, mellow, sweet, deep, powerfulUsageNews broadcast, ad voice-over, audiobook, animation role, voice assistant, documentary narration 
**Examples**

- Standard broadcast style: Clear and precise articulation, well-rounded pronunciation.

- Progressive emotional effect: Volume rapidly increases from normal conversation to a shout, with a straightforward personality and easily excited, expressive emotions.

- Special emotional state: A sobbing tone causes slightly slurred and hoarse pronunciation, with noticeable tension in the crying voice.

- Ad voice-over style: High-pitched, medium speed, full of energy and appeal, suitable for ad voice-overs.

- Gentle and soothing style: Slow-paced, with a gentle and sweet pitch, and a soothing, warm tone, like a caring friend.


## [​ ](#voice-customization) Voice customization

Qwen3-TTS supports both voice cloning (Qwen3-TTS-VC) and voice design (Qwen3-TTS-VD). See [Voice cloning (Qwen)](/api-reference/speech-synthesis/voice-cloning/qwen/create-voice) and [Voice design (Qwen)](/api-reference/speech-synthesis/voice-design/qwen/create-voice) for the API reference.
## [​ ](#api-reference) API reference


- [Speech synthesis - Qwen API reference](/api-reference/speech-synthesis/qwen-tts)

- [Voice cloning API reference](/api-reference/speech-synthesis/voice-cloning/qwen/create-voice)

- [Voice design API reference](/api-reference/speech-synthesis/voice-design/qwen/create-voice)


## [​ ](#built-in-voices) Built-in voices

See [Qwen-TTS voice list](/api-reference/speech-synthesis/qwen-tts/voice-list) for the list of supported voices, model compatibility, and audio samples.
## [​ ](#faq) FAQ

**Q: How long is the audio file URL valid?**
The audio file URL expires after 24 hours.
## [​ ](#learn-more) Learn more


- [Real-time speech synthesis](/developer-guides/speech/realtime-streaming) — Real-time streaming speech synthesis with WebSocket

- [CosyVoice voice list](/api-reference/speech-synthesis/cosyvoice/voice-list)

- [Qwen-TTS voice list](/api-reference/speech-synthesis/qwen-tts/voice-list)


 [Previous ](/developer-guides/speech/realtime-streaming)[Voice cloning Clone a voice from audio samples for use with CosyVoice, Qwen-TTS, or Qwen-Omni models. Next ](/developer-guides/speech/voice-cloning)
