# Text-to-speech models

> **Source:** https://docs.qwencloud.com/developer-guides/speech/tts-models

Choose a model for speech synthesis, voice cloning, and voice design.

 Copy page Two questions narrow the field: do you need a custom voice or will a built-in voice work, and do you need real-time streaming?
## [​ ](#built-in-or-custom-voice) Built-in or custom voice?

### [​ ](#built-in-voices) Built-in voices

Pick a voice from the library and start synthesizing immediately.

- **CosyVoice** — rich voice library, high quality, no setup beyond picking a voice

- **Qwen3-TTS** — low-latency streaming; add `-instruct` for natural-language control over speed, emotion, and style


### [​ ](#custom-voice) Custom voice

Need a voice that doesn&#x27;t exist in the library?

- **Voice Cloning** — reproduce a specific person&#x27;s voice from audio samples. Use when you have a target voice to match.

- **Voice Design** — create a new voice from a text description (e.g., "a warm, low-pitched female voice"). Use when you want a brand voice without audio samples.


## [​ ](#controlling-how-the-voice-sounds) Controlling how the voice sounds

Three approaches, ranked by flexibility:

- 
**Instruction control** (`cosyvoice-v3-flash`, `qwen3-tts-instruct-flash`, `qwen3-tts-instruct-flash-realtime`) — Describe the desired delivery in natural language. Control speed, emotion, and style per request. Most flexible. For details, see [Real-time speech synthesis > Instruction control](/developer-guides/speech/realtime-streaming#instruction-based-control).


- 
**Voice design** (`qwen3-tts-vd-*`) — Generate a custom voice from a text description. Good for creating a brand voice without audio samples.


- 
**Voice cloning** (`qwen3-tts-vc-*`) — Reproduce an existing voice from audio samples. Best when you need to match a specific person&#x27;s voice.


## [​ ](#recommended-models) Recommended models

ModelFamilyStreamingCustom voiceInstruction control`cosyvoice-v3-plus`CosyVoice✓——`qwen3-tts-flash`Qwen3-TTS✓——`qwen3-tts-flash-realtime`Qwen3-TTS✓——`qwen3-tts-instruct-flash`Qwen3-TTS✓—✓`qwen3-tts-vc-realtime-2026-01-15`Voice Cloning✓✓—`qwen3-tts-vd-realtime-2026-01-15`Voice Design✓✓— 
## [​ ](#all-models) All models

 CosyVoice

 ModelStreamingCustom voiceInstruction control`cosyvoice-v3-plus`✓——`cosyvoice-v3-flash`✓—— Qwen3-TTS

 ModelStreamingCustom voiceInstruction control`qwen3-tts-flash`✓——`qwen3-tts-flash-realtime`✓——`qwen3-tts-instruct-flash`✓—✓`qwen3-tts-instruct-flash-realtime`✓—✓ Voice Cloning & Design

 ModelStreamingCustom voiceInstruction control`qwen3-tts-vc-2026-01-22`✗✓—`qwen3-tts-vc-realtime-2026-01-15`✓✓—`qwen3-tts-vd-2026-01-26`✗✓—`qwen3-tts-vd-realtime-2026-01-15`✓✓— Legacy

 Previous generation models. We recommend the latest versions above for new projects.ModelFamilyStreamingCustom voiceInstruction control`qwen3-tts-flash-2025-11-27`Qwen3-TTS✓——`qwen3-tts-flash-2025-09-18`Qwen3-TTS✓——`qwen3-tts-flash-realtime-2025-11-27`Qwen3-TTS✓——`qwen3-tts-flash-realtime-2025-09-18`Qwen3-TTS✓——`qwen3-tts-instruct-flash-2026-01-26`Qwen3-TTS✓—✓`qwen3-tts-instruct-flash-realtime-2026-01-22`Qwen3-TTS✓—✓`qwen3-tts-vc-realtime-2025-11-27`Voice Cloning✓✓—`qwen3-tts-vd-realtime-2025-12-16`Voice Design✓✓— 
## [​ ](#learn-more) Learn more

[ ## Text-to-speech guide
Learn how to use TTS models via API. ](/developer-guides/speech/tts)[ ## Real-time streaming guide
Use real-time TTS models via WebSocket. ](/developer-guides/speech/realtime-streaming)[ ## CosyVoice voices
Browse CosyVoice voices and samples. ](/api-reference/speech-synthesis/cosyvoice/voice-list)[ ## Qwen-TTS voices
Browse Qwen-TTS voices for non-streaming models. ](/api-reference/speech-synthesis/qwen-tts/voice-list#non-real-time-speech-synthesis-voices)[ ## Qwen-TTS-Realtime voices
Browse Qwen-TTS-Realtime voices for streaming models. ](/api-reference/speech-synthesis/qwen-tts/voice-list#real-time-speech-synthesis-voices)[ ## Voice cloning
Clone a voice from audio samples. ](/api-reference/speech-synthesis/voice-cloning/qwen/create-voice) [Previous ](/developer-guides/speech/improve-recognition-accuracy)[Real-time speech synthesis Stream TTS in real time Next ](/developer-guides/speech/realtime-streaming)
