# Create a cloned voice

> **Source:** https://docs.qwencloud.com/api-reference/speech-synthesis/voice-cloning/qwen/create-voice

POST/services/audio/tts/customization cURL cURL

 Copy ```\ncurl -X POST https://dashscope-intl.aliyuncs.com/api/v1/services/audio/tts/customization \
-H "Authorization: Bearer $DASHSCOPE_API_KEY" \
-H "Content-Type: application/json" \
-d '{
 "model": "qwen-voice-enrollment",
 "input": {
 "action": "create",
 "target_model": "qwen3-tts-vc-realtime-2026-01-15",
 "preferred_name": "guanyu",
 "audio": {
 "data": "https://your-audio-url.wav"
 }
 }
}'
``` 200 400 Copy ```\n{
 "output": {
 "voice": "qwen-omni-vc-guanyu-voice-20250812105009984-838b",
 "target_model": "qwen3.5-omni-plus-realtime",
 "fallback_mode": false,
 "fallback_reason": "&#x3C;string>"
 },
 "usage": {
 "count": 1
 },
 "request_id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
``` ### Authorizations
 [​ ](#authorization) Authorizationstring header required DashScope API key. Get one at [API key](/api-reference/preparation/api-key).

 ### Body
application/json [​ ](#model) modelenum<string> required Fixed to `qwen-voice-enrollment`.

 Available options:qwen-voice-enrollment Example:qwen-voice-enrollment [​ ](#input) inputobject required Show child attributes

 [​ ](#inputaction) input.actionenum<string> required Fixed to `create`.

 Available options:create Example:create [​ ](#inputtarget-model) input.target_modelenum<string> required Target model for the cloned voice. Must match the model in subsequent API calls.

 Available options:qwen3.5-omni-plus-realtime,qwen3.5-omni-flash-realtime,qwen3-tts-vc-realtime-2026-01-15,qwen3-tts-vc-realtime-2025-11-27,qwen3-tts-vc-2026-01-22 Example:qwen3.5-omni-plus-realtime [​ ](#inputaudio) input.audioobject required Audio data for cloning.

 Show child attributes

 [​ ](#inputaudiodata) input.audio.datastring required **Data URL** (Base64): `data:<mediatype>;base64,<data>` where `<mediatype>` is `audio/wav`, `audio/mpeg`, or `audio/mp4`. Keep encoded data under 10 MB. **Audio URL**: publicly accessible URL (no auth required).

 Example:https://your-audio-url.wav [​ ](#inputpreferred-name) input.preferred_namestring Voice name keyword (digits, letters, underscores, max 16 characters). Appears in the generated voice name. Example: `guanyu` produces `qwen-tts-vc-guanyu-voice-20250812105009984-838b`.

 Example:guanyu Required range:length <= 16pattern: ^[a-zA-Z0-9_]+$ [​ ](#inputtext) input.textstring Text matching the audio content. The server validates the match and returns `Audio.PreprocessError` if significantly different.

 Example:Hello, this is a sample text for voice cloning. [​ ](#inputlanguage) input.languageenum<string> Audio language code. Must match the audio language if specified.

 Available options:zh,en,de,it,pt,es,ja,ko,fr,ru Example:zh ### Response
200-application/json [​ ](#output) outputobject Show child attributes

 [​ ](#outputvoice) output.voicestring Generated voice name. Pass as the `voice` parameter in subsequent Qwen TTS or Realtime Multimodal API calls.

 Example:qwen-omni-vc-guanyu-voice-20250812105009984-838b [​ ](#outputtarget-model) output.target_modelstring Target model bound to this voice.

 Example:qwen3.5-omni-plus-realtime [​ ](#outputfallback-mode) output.fallback_modeboolean `true` if the voice was created in degraded mode due to poor audio quality or text mismatch.

 Example:false [​ ](#outputfallback-reason) output.fallback_reasonstring Reason for degradation. Possible values: `no_merged_segments`, `no_valid_asr_segments`, etc. Only returned when `fallback_mode` is `true`.

 [​ ](#usage) usageobject Show child attributes

 [​ ](#usagecount) usage.countinteger Billed voice creation operations. Always `1`.

 Example:1 [​ ](#request-id) request_idstring Request ID for troubleshooting.

 Example:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
