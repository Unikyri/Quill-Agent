# Video generation models

> **Source:** https://docs.qwencloud.com/developer-guides/getting-started/video-models

Choose a model for text-to-video, image-to-video, reference-to-video, video editing, and more.

 Copy page ## [​ ](#text-to-video) Text-to-video

Generate videos from text prompts with audio → `happyhorse-1.1-t2v`. 1080P, 3--15s. To upload custom audio input (narration, background music), use `wan2.7-t2v-2026-04-25`.
## [​ ](#image-to-video) Image-to-video

Animate a still image into motion. First-frame to video → `happyhorse-1.1-i2v` (1080P, 3--15s, audio). For first+last frame control, building long videos by chaining segments, or uploading custom audio input, use `wan2.7-i2v-2026-04-25`.
## [​ ](#reference-to-video) Reference-to-video

Consistent characters across scenes from reference images → `happyhorse-1.1-r2v` (1--9 reference images, 1080P, 3--15s). For custom audio input or video-based references, use `wan2.7-r2v`.
## [​ ](#video-editing) Video editing

Edit existing videos with text instructions and reference images → `happyhorse-1.0-video-edit` (style transfer, element replacement, 1080P). For effect replication or camera motion replication, use `wan2.7-videoedit`.
## [​ ](#character-animation) Character animation

Transfer motion from a reference video onto a person in a still image → `wan2.2-animate-move`. Replace a person in a video with someone from an image → `wan2.2-animate-mix`.
## [​ ](#recommended-models) Recommended models

ModelUse this when...Output audioMax resolutionMax duration`happyhorse-1.1-t2v`Text-to-videoYes720P, 1080P3--15s`wan2.7-t2v-2026-04-25`Text-to-video, custom audio inputYes720P, 1080P2--15s`happyhorse-1.1-i2v`Image-to-video (first frame)Yes720P, 1080P3--15s`wan2.7-i2v-2026-04-25`Image-to-video (first frame, first+last frame, video continuation)Yes720P, 1080P2--15s`happyhorse-1.1-r2v`Consistent characters from reference images (1--9)Yes720P, 1080P3--15s`wan2.7-r2v`Reference images and video references, voice cloningYes720P, 1080P2--10s`happyhorse-1.0-video-edit`Edit videos with prompts and reference imagesYes720P, 1080P3--15s`wan2.7-videoedit`Effect replication, camera replicationYes720P, 1080Pup to 10s`wan2.2-animate-move`Motion transfer onto a still personYes720P2--30s`wan2.2-animate-mix`Replace a person in a videoYes720P2--30s 
## [​ ](#all-models) All models

 HappyHorse

 ModelCapabilityFeaturesOutput`happyhorse-1.1-t2v`Text-to-videoAudio, aspect ratio control720P, 1080P. 3--15s. 24 fps, MP4`happyhorse-1.1-i2v`Image-to-videoAudio, first-frame input720P, 1080P. 3--15s. 24 fps, MP4`happyhorse-1.1-r2v`Reference-to-videoAudio, multi-character (1--9 refs), aspect ratio control720P, 1080P. 3--15s. 24 fps, MP4`happyhorse-1.0-t2v`Text-to-videoAudio, aspect ratio control720P, 1080P. 3--15s. 24 fps, MP4`happyhorse-1.0-i2v`Image-to-videoAudio, first-frame input720P, 1080P. 3--15s. 24 fps, MP4`happyhorse-1.0-r2v`Reference-to-videoAudio, multi-character (1--9 refs), aspect ratio control720P, 1080P. 3--15s. 24 fps, MP4`happyhorse-1.0-video-edit`Video editingAudio (auto/original), reference images (0--5)720P, 1080P. 3--15s. 24 fps, MP4 Wan 2.7

 ModelCapabilityFeaturesOutput`wan2.7-t2v`Text-to-videoAudio sync, multi-shot narrative, aspect ratio control720P, 1080P. 2–15s. 30 fps, MP4`wan2.7-t2v-2026-04-25`Text-to-video (snapshot)Same as `wan2.7-t2v`720P, 1080P. 2–15s. 30 fps, MP4`wan2.7-i2v`Image-to-videoAudio sync, first-frame, first-last-frame, video continuation720P, 1080P. 2–15s. 30 fps, MP4`wan2.7-i2v-2026-04-25`Image-to-video (snapshot)Same as `wan2.7-i2v`720P, 1080P. 2–15s. 30 fps, MP4`wan2.7-r2v`Reference-to-videoMulti-character, image/video references720P, 1080P. 2–10s. 30 fps, MP4`wan2.7-videoedit`Video editingInstruction editing, video migration720P, 1080P. Up to 10s. 30 fps, MP4 Character Animation

 ModelCapabilityFeaturesOutput`wan2.2-animate-move`Motion transfer`wan-std` / `wan-pro` modes720P. 2s--30s. 15/25 fps. MP4`wan2.2-animate-mix`Face swap`wan-std` / `wan-pro` modes720P. 2s--30s. 15/25 fps. MP4 Legacy

 Previous generation models. We recommend HappyHorse or Wan 2.7 for new projects.### Wan 2.6
ModelCapabilityFeaturesOutput`wan2.6-t2v`Text-to-videoAudio sync, multi-shot narrative720P, 1080P. 2–15s. 30 fps, MP4`wan2.6-i2v`Image-to-videoAudio sync, multi-shot narrative720P, 1080P. 2–15s. 30 fps, MP4`wan2.6-i2v-flash`Image-to-videoAudio, multi-shot, fast720P, 1080P. 2–15s. 30 fps, MP4`wan2.6-r2v`Reference-to-videoAudio sync, multi-character, narrative720P, 1080P. 2–10s. 30 fps, MP4`wan2.6-r2v-flash`Reference-to-videoMulti-character, fast720P, 1080P. 2–10s. 30 fps, MP4 ### Wan 2.5
ModelCapabilityFeaturesOutput`wan2.5-t2v-preview`Text-to-videoAudio sync480P, 720P, 1080P. 5s, 10s. 30 fps, MP4`wan2.5-i2v-preview`Image-to-videoAudio sync480P, 720P, 1080P. 5s, 10s. 30 fps, MP4 ### Wan 2.2
ModelCapabilityFeaturesOutput`wan2.2-t2v-plus`Text-to-videoNo audio480P, 1080P. 5s. 30 fps, MP4`wan2.2-i2v-plus`Image-to-videoNo audio480P, 1080P. 5s. 30 fps, MP4`wan2.2-i2v-flash`Image-to-videoNo audio, 50% faster than 2.1480P, 720P, 1080P. 5s. 30 fps, MP4`wan2.2-kf2v-flash`First & last framesNo audio480P, 720P, 1080P. 5s. 30 fps, MP4 ### Wan 2.1
ModelCapabilityFeaturesOutput`wan2.1-t2v-plus`Text-to-videoNo audio720P. 5s. 30 fps, MP4`wan2.1-t2v-turbo`Text-to-videoNo audio480P, 720P. 5s. 30 fps, MP4`wan2.1-i2v-plus`Image-to-videoNo audio720P. 5s. 30 fps, MP4`wan2.1-i2v-turbo`Image-to-videoNo audio480P, 720P. 3s–5s. 30 fps, MP4`wan2.1-kf2v-plus`First & last framesNo audio720P. 5s. 30 fps, MP4`wan2.1-vace-plus`Video editingNo audio720P. Up to 5s. 30 fps, MP4 

---

## [​ ](#learn-more) Learn more

[ ## Text-to-video
Generate videos from text prompts. ](/developer-guides/video-generation/text-to-video)[ ## Image-to-video: first frame
Animate from a single image. ](/developer-guides/video-generation/image-to-video)[ ## First & last frames
Animate between two frames. ](/developer-guides/video-generation/image-to-video-first-last)[ ## Reference-to-video
Generate videos with character consistency. ](/developer-guides/video-generation/reference-video)[ ## Video editing
Edit, extend, and redraw videos. ](/developer-guides/video-generation/video-editing) [Previous ](/developer-guides/image-generation/wan-image-editing)[Text-to-video Generate video from text Next ](/developer-guides/video-generation/text-to-video)
