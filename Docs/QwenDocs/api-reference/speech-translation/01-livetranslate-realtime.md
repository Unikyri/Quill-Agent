# LiveTranslate client events

> **Source:** https://docs.qwencloud.com/api-reference/speech-translation/livetranslate-realtime/client-events

WebSocket client reference

 Copy page This topic describes the client events for the qwen3.5-livetranslate-flash-realtime API.
 Reference: [Speech translation](/developer-guides/speech/realtime-translation). 
## [​ ](#connect) Connect

Establish a WebSocket connection to start a session. The server sends a `session.created` event when the connection is ready.
ConfigurationValueEndpoint`wss://dashscope-intl.aliyuncs.com/api-ws/v1/realtime`Query parameter`model=qwen3.5-livetranslate-flash-realtime`Auth header`Authorization: Bearer $DASHSCOPE_API_KEY`ProtocolJSON text frames 
Full URL:
Copy ```\nwss://dashscope-intl.aliyuncs.com/api-ws/v1/realtime?model=qwen3.5-livetranslate-flash-realtime

``` 
## [​ ](#session-update) session.update

Updates the session configuration after you connect. The server validates parameters and returns the full configuration, or an error if any value is invalid.
Example Copy ```\n{
 "event_id": "event_ToPZqeobitzUJnt3QqtWg",
 "type": "session.update",
 "session": {
 "modalities": [
 "text",
 "audio"
 ],
 "voice": "Tina",
 "input_audio_format": "pcm16",
 "output_audio_format": "pcm24",
 "input_audio_transcription": {
 "model": "qwen3-asr-flash-realtime",
 "language": "zh"
 },
 "translation": {
 "language": "en"
 }
 }
}

``` 
[​ ](#param-type) typestring body required Always `"session.update"`. 
[​ ](#param-session) sessionobject body Session configuration.Show properties

 [​ ](#param-modalities) modalitiesarray body Output types. Valid values:
- `["text"]` — Text only.

- `["text", "audio"]` (default) — Text and audio.


 [​ ](#param-voice) voicestring body Voice for audio output. See [Supported voices](/developer-guides/speech/realtime-translation#supported-voices). Default: `Tina` for Qwen3.5-LiveTranslate-Flash-Realtime, or `Cherry` for Qwen3-LiveTranslate-Flash-Realtime. [​ ](#param-input_audio_transcription) input_audio_transcriptionobject body Input audio settings.Show properties

 [​ ](#param-model) modelstring body Speech recognition model. When set, the server returns source-language text through `conversation.item.input_audio_transcription.text` and `conversation.item.input_audio_transcription.completed` events.Valid value: `qwen3-asr-flash-realtime`. [​ ](#param-language) languagestring body Source language. See [Supported languages](/developer-guides/speech/realtime-translation#supported-languages). Default: `en`. [​ ](#param-input_audio_format) input_audio_formatstring body Input audio format. Must be `pcm16`. [​ ](#param-output_audio_format) output_audio_formatstring body Output audio format. Must be `pcm24`. [​ ](#param-translation) translationobject body Translation settings.Show properties

 [​ ](#param-language_2) languagestring body Target language. See [Supported languages](/developer-guides/speech/realtime-translation#supported-languages). Default: `en`. [​ ](#param-enable_voice_clone) enable_voice_cloneboolean body Whether to enable voice cloning. Default: `false`. When enabled, the model clones the speaker&#x27;s voice from the input audio for translated output. In this case, `voice` no longer accepts system preset voices and must be set to `default` or a voice ID previously created through the Voice Clone API. [​ ](#param-voice_clone_options) voice_clone_optionsobject body Voice cloning options. Only applies when `enable_voice_clone` is `true`.Show properties

 [​ ](#param-frequency) frequencystring body Controls when voice cloning occurs. Valid values: `never` (use pre-cloned voice profile), `once` (clone at session start), `always` (clone before each response). Default: `once`. 
## [​ ](#input-audio-buffer-append) input_audio_buffer.append

Appends audio bytes to the input buffer. The server uses this buffer for speech detection and submission timing.
Example Copy ```\n{
 "event_id": "event_xxx",
 "type": "input_audio_buffer.append",
 "audio": "xxx"
}

``` 
[​ ](#param-type_2) typestring body required Always `"input_audio_buffer.append"`. 
[​ ](#param-audio) audiostring body required Base64-encoded audio data. 
## [​ ](#input-image-buffer-append) input_image_buffer.append

Adds image data to the buffer from a local file or a real-time video stream.
Image limits:

- Format: JPG or JPEG. Recommended resolution: 480p or 720p. Maximum: 1080p.

- Maximum size: 500 KB (before Base64 encoding).

- Must be Base64-encoded.

- Maximum rate: 2 images per second.

- You must send at least one `input_audio_buffer.append` event first.


Example Copy ```\n{
 "event_id": "event_xxx",
 "type": "input_image_buffer.append",
 "image": "xxx"
}

``` 
[​ ](#param-type_3) typestring body required Always `"input_image_buffer.append"`. 
[​ ](#param-image) imagestring body required Base64-encoded image data. 
## [​ ](#session-finish) session.finish

Ends the session. The server responds based on whether it detected speech:

- **Speech detected:** The server finishes recognition and sends [conversation.item.input_audio_transcription.completed](/api-reference/speech-translation/livetranslate-realtime/server-events) with the result, then sends [session.finished](/api-reference/speech-translation/livetranslate-realtime/server-events).

- **No speech detected:** The server sends [session.finished](/api-reference/speech-translation/livetranslate-realtime/server-events) directly.


Disconnect after you receive `session.finished`.
Example Copy ```\n{
 "event_id": "event_xxx",
 "type": "session.finish"
}

``` 
[​ ](#param-type_4) typestring body required Always `"session.finish"`. [Previous ](/api-reference/real-time-multimodal/realtime-java-sdk)[LiveTranslate server events WebSocket server reference Next ](/api-reference/speech-translation/livetranslate-realtime/server-events)
