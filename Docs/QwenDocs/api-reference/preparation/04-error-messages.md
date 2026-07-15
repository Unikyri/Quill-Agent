# Error messages

> **Source:** https://docs.qwencloud.com/api-reference/preparation/error-messages

API error code reference

 Copy page Error messages from Qwen Cloud and their solutions.
## [​ ](#400-invalidparameter) 400-InvalidParameter

### [​ ](#parameter-enable-thinking-must-be-set-to-false-for-non-streaming-calls-parameter-enable-thinking-only-support-stream-call) parameter.enable_thinking must be set to false for non-streaming calls/ parameter.enable_thinking only support stream call

**Cause:** The model does not support `enable_thinking` with non-streaming calls.
**Solution:** Set `enable_thinking` to `false`, use [streaming](/developer-guides/text-generation/streaming), or switch to a model that supports non-streaming thinking (such as qwen3.7-max, qwen3-max, and qwen3.5-plus).
### [​ ](#the-thinking-budget-parameter-must-be-a-positive-integer-and-not-greater-than-xxx) The thinking_budget parameter must be a positive integer and not greater than xxx

**Cause:** `thinking_budget` is out of range.
**Solution:** Set to a positive integer within the model&#x27;s maximum reasoning length. See [Models](https://www.qwencloud.com/models).
### [​ ](#this-model-only-support-stream-mode-please-enable-the-stream-parameter-to-access-the-model) This model only support stream mode, please enable the stream parameter to access the model.

**Cause:** Model requires [streaming output](/developer-guides/text-generation/streaming).
**Solution:** Enable [streaming output](/developer-guides/text-generation/streaming).
### [​ ](#this-model-does-not-support-enable-search) This model does not support enable_search.

**Cause:** Model does not support [web search](/developer-guides/text-generation/web-search).
**Solution:** Use a model that supports web search, or set `enable_search` to `false`.
### [​ ](#current-language-settings-are-not-supported) Current language settings are not supported!

**Cause:** Invalid `source_lang` or `target_lang` format, or unsupported language.
**Solution:** Use the correct English name or language code.
### [​ ](#the-incremental-output-parameter-must-be-true-when-enable-thinking-is-true) The incremental_output parameter must be "true" when enable_thinking is true

**Cause:** Thinking mode requires incremental streaming output.
**Solution:** Set `incremental_output` to `true`.
### [​ ](#the-incremental-output-parameter-of-this-model-cannot-be-set-to-false) The incremental_output parameter of this model cannot be set to False.

**Cause:** Model requires incremental streaming output.
**Solution:** Set `incremental_output` to `true`.
### [​ ](#range-of-input-length-should-be-1-xxx) Range of input length should be [1, xxx]

**Cause:** Input exceeds the model&#x27;s token limit.
**Solution:**

- 
API calls: Keep messages array within the model&#x27;s token limit.


- 
Chat conversations: Start a new conversation when history exceeds the limit.


### [​ ](#range-of-max-tokens-should-be-1-xxx) Range of max_tokens should be [1, xxx]

**Cause:** `max_tokens` is out of range.
**Solution:** Set to a value between 1 and the model&#x27;s maximum output tokens. See [Models](https://www.qwencloud.com/models).
### [​ ](#temperature-should-be-in-0-0-2-0-temperature-must-be-float) Temperature should be in [0.0, 2.0)/&#x27;temperature&#x27; must be Float

**Cause:** `temperature` is out of range.
**Solution:** Set to a value in [0.0, 2.0).
### [​ ](#range-of-top-p-should-be-0-0-1-0-top-p-must-be-float) Range of top_p should be (0.0, 1.0]/&#x27;top_p&#x27; must be Float

**Cause:** `top_p` is out of range.
**Solution:** Set to a value in (0.0, 1.0].
### [​ ](#parameter-top-k-be-greater-than-or-equal-to-0) Parameter top_k be greater than or equal to 0

**Cause:** `top_k` is less than 0.
**Solution:** Set to a value ≥ 0.
### [​ ](#repetition-penalty-should-be-greater-than-0-0) Repetition_penalty should be greater than 0.0

**Cause:** `repetition_penalty` is ≤ 0.
**Solution:** Set to a value > 0.
### [​ ](#presence-penalty-should-be-in-2-0-2-0) Presence_penalty should be in [-2.0, 2.0]

**Cause:** `presence_penalty` is out of range.
**Solution:** Set to a value in [-2.0, 2.0].
### [​ ](#range-of-n-should-be-1-4) Range of n should be [1, 4]

**Cause:** `n` is out of range.
**Solution:** Set to a value in [1, 4].
### [​ ](#range-of-seed-should-be-0-9223372036854775807) Range of seed should be [0, 9223372036854775807]

**Cause:** `seed` is out of range (DashScope protocol).
**Solution:** Set to a value in [0, 9223372036854775807].
### [​ ](#request-method-get-is-not-supported) Request method &#x27;GET&#x27; is not supported.

**Cause:** API does not support `GET`.
**Solution:** Use `POST`. See the API reference.
### [​ ](#messages-with-role-tool-must-be-a-response-to-a-preceding-message-with-tool-calls) Messages with role "tool" must be a response to a preceding message with "tool_calls"

**Cause:** The tool calling sequence is missing an assistant message.
**Solution:** Include the assistant message from the model&#x27;s first response before the tool message.
### [​ ](#required-body-invalid-check-the-request-body-format) Required body invalid, check the request body format.

**Cause:** Invalid request body format.
**Solution:** Ensure the request body is valid JSON. Check for extra commas or unclosed brackets.
### [​ ](#input-content-must-be-a-string) input content must be a string.

**Cause:** Text-only models do not support array-type content.
**Solution:** Use string content, not arrays like `[{"type": "text","text": "..."}]`.
### [​ ](#the-content-field-is-a-required-field) The content field is a required field.

**Cause:** Missing `content` parameter (such as `{"role": "user"}`).
**Solution:** Add `content`, such as `{"role": "user","content": "Who are you"}`.
### [​ ](#current-user-api-does-not-support-http-call) current user api does not support http call.

**Cause:** Model does not support non-streaming output.
**Solution:** Use [streaming output](/developer-guides/text-generation/streaming).
### [​ ](#either-prompt-or-messages-must-exist-and-cannot-both-be-none) Either "prompt" or "messages" must exist and cannot both be none

**Cause:** Missing `messages` or `prompt` parameter, or wrong format.
**Solution:** Specify `messages`. For DashScope-HTTP, place `messages` inside the `input` object. See [Qwen](/api-reference/chat/dashscope) API reference.
### [​ ](#messages-must-contain-the-word-json-in-some-form-to-use-response-format-of-type-json-object) &#x27;messages&#x27; must contain the word &#x27;json&#x27; in some form, to use &#x27;response_format&#x27; of type &#x27;json_object&#x27;.

**Cause:** [Structured output](/developer-guides/text-generation/structured-output) requires "json" in the prompt.
**Solution:** Add "json" (case-insensitive) to the prompt, such as "Output in json format."
### [​ ](#json-mode-response-is-not-supported-when-enable-thinking-is-true) Json mode response is not supported when enable_thinking is true

**Cause:** [Structured output](/developer-guides/text-generation/structured-output) does not work with thinking mode.
**Solution:** Set `enable_thinking` to `false`. See [How do models in thinking mode produce structured output?](/developer-guides/text-generation/structured-output)
### [​ ](#tool-names-are-not-allowed-to-be-search) Tool names are not allowed to be [search]

**Cause:** Tool name cannot be "search".
**Solution:** Use a different tool name.
### [​ ](#unknown-format-of-response-format-response-format-should-be-a-dict-includes-type-and-an-optional-key-json-schema-the-response-format-type-from-user-is-xxx) Unknown format of response_format, response_format should be a dict, includes &#x27;type&#x27; and an optional key &#x27;json_schema&#x27;. The response_format type from user is xxx.

**Cause:** Invalid `response_format`.
**Solution:** Set to `{"type": "json_object"}` for [structured output](/developer-guides/text-generation/structured-output).
### [​ ](#the-value-of-the-enable-thinking-parameter-is-restricted-to-true) The value of the enable_thinking parameter is restricted to True.

**Cause:** Model requires thinking mode (such as `qwen3-235b-a22b-thinking-2507`).
**Solution:**

- 
In code: Set `enable_thinking` to `true`.


- 
In third-party tools (such as Cherry Studio): Enable thinking in the input box.


### [​ ](#audio-output-only-support-with-stream-true) &#x27;audio&#x27; output only support with stream=true

**Cause:** Qwen-Omni requires streaming output.
**Solution:** Set `stream` to `true`.
### [​ ](#tool-choice-is-one-of-the-strings-that-should-be-none-auto) tool_choice is one of the strings that should be ["none", "auto"]

**Cause:** Invalid `tool_choice` value.
**Solution:** Use "auto" (model chooses) or "none" (no tool).
### [​ ](#model-not-exist) Model not exist.

**Cause:** Invalid or nonexistent `model` parameter.
**Solution:**

- 
**Model name format:** Check casing and remove extra spaces.


- 
**Correct model name:** Use Qwen Cloud model IDs (such as `qwen3-235b-a22b-instruct-2507`), not community names (such as `Qwen/Qwen3-235B-A22B-Instruct-2507`). See [Models](https://www.qwencloud.com/models).


### [​ ](#the-result-format-parameter-must-be-message-when-enable-thinking-is-true) The result_format parameter must be "message" when enable_thinking is true

**Cause:** Thinking mode requires `result_format` = `"message"`.
**Solution:** Set `result_format` to `"message"`.
### [​ ](#the-audio-is-empty) The audio is empty

**Cause:** Audio is too short.
**Solution:** Use a longer audio file.
### [​ ](#file-parsing-in-progress-try-again-later) File parsing in progress, try again later.

**Cause:** File is still being parsed for document processing.
**Solution:** Wait for parsing to finish, then retry.
### [​ ](#the-stop-parameter-must-be-of-type-str-list-str-list-int-or-list-list-int-and-all-elements-within-the-list-must-be-of-the-same-type) The "stop" parameter must be of type "str", "list[str]", "list[int]", or "list[list[int]]", and all elements within the list must be of the same type.

**Cause:** Invalid `stop` parameter format.
**Solution:** Use `str`, `list[str]`, `list[int]`, or `list[list[int]]`. See [Qwen](/api-reference/chat/openai-chat) API reference.
### [​ ](#value-error-batch-size-is-invalid-it-should-not-be-larger-than-xxx) Value error, batch size is invalid, it should not be larger than xxx.

**Cause:** Too many texts for the embedding model.
**Solution:** Keep the input count within the model&#x27;s batch size. See [Embedding](/developer-guides/embeddings/text-embedding).
### [​ ](#is-too-short) [] is too short

**Cause:** Empty messages array.
**Solution:** Add at least one message.
### [​ ](#the-tool-call-is-not-supported) The tool call is not supported.

**Cause:** Model does not support tool calling.
**Solution:** Switch to a Qwen or DeepSeek model that supports Function Calling.
### [​ ](#the-provided-messages-input-is-invalid-the-error-info-is-unexpected-item-type-in-content-input-should-be-a-valid-string-input-should-be-a-valid-dictionary-or-instance-of-input-messages-x-content) The provided messages input is invalid. The error info is [Unexpected item type in content] / Input should be a valid string / Input should be a valid dictionary or instance of ... (input.messages.x.content...)

**Cause:** The `content` field of a message in the `messages` array contains an item of an unsupported type. When `content` is an array, each element must be a string or a valid object (such as `{"type": "text", "text": "..."}` or `{"type": "image_url", ...}`). Passing a number, boolean, nested array, or an object whose `type` is not supported triggers this error. The corresponding HTTP status code is 400 and the error code is `InternalError.Algo.InvalidParameter`.
**Common scenario:** When you use a text-only model (such as the Qwen-Max text series, including `qwen3-max`) and the `messages` (especially the multi-turn conversation history) contain multimodal `content` items such as images (`image_url`), this error is also triggered, because text-only models do not accept image or other modality inputs. In some integrated tools or agents, this error may be surfaced as a message like "the model provider returned empty content".
**Solution:**

- 
**Text-only models:** Set `content` to a string instead of an array or any other type. If your use case requires image or other multimodal input, switch to a multimodal model (such as the Qwen-VL or Qwen3-VL series). If you must keep using a text-only model, remove multimodal items such as images (`image_url`) from `messages` and the conversation history before calling.


- 
**Multimodal models:** Each element in the `content` array must be a valid object whose `type` is one of the modality types supported by the model, such as `text`, `image_url`, `video_url`, or `video`. Do not include numbers, booleans, nested arrays, or elements that are missing `type` or whose `type` value is invalid.


### [​ ](#repetitive-tool-calls-detected-in-the-conversation-history-the-same-tool-call-with-identical-name-and-arguments-has-been-repeated-across-multiple-consecutive-rounds-please-modify-your-request-or-adjust-the-tool-call-arguments-to-avoid-infinite-loops) Repetitive tool calls detected in the conversation history. The same tool call with identical name and arguments has been repeated across multiple consecutive rounds. Please modify your request or adjust the tool call arguments to avoid infinite loops.

**Cause:** A repetitive tool call was detected in the conversation history. A tool call with the exact same name and arguments was repeated in multiple consecutive rounds, which indicates that the model may be in a loop. The corresponding HTTP status code is 400, and the error code is `InternalError.Algo.InvalidParameter`.
**Solution:**

- 
After each tool call, append the tool&#x27;s execution result as a tool message to the `messages` array before you send the next request. This helps the model proceed instead of repeating the same call.


- 
Check whether the tool is returning results correctly. If the tool consistently returns the same or invalid results, add a tool call limit or termination logic on the application side to prevent infinite loops.


- 
If necessary, optimize the prompt to clarify the task completion conditions to prevent the model from repeatedly making the same tool call.


### [​ ](#required-parameter-xxx-missing-or-invalid-check-the-request-parameters) Required parameter(xxx) missing or invalid, check the request parameters.

**Cause:** Invalid API parameters.
**Solution:** Provide all required parameters in the correct format.
### [​ ](#input-must-contain-file-urls) input must contain file_urls

**Cause:** Missing `file_urls` parameter (Paraformer).
**Solution:** Provide `file_urls` in the request.
### [​ ](#the-provided-url-does-not-appear-to-be-valid-ensure-it-is-correctly-formatted) The provided URL does not appear to be valid. Ensure it is correctly formatted.

**Cause:** Invalid URL or file path for multimodal input.
**Solution:**

- 
URLs: Must start with `http://`, `https://`, or `data:`. For `data:` URLs, include `"base64"` before encoded data.


- 
Local paths: Must start with `file://`.


- 
Temporary URLs: Add `X-DashScope-OssResourceResolve: enable` to request headers (HTTP), or use DashScope SDK (not OpenAI SDK).


### [​ ](#input-should-be-a-valid-dictionary-or-instance-of-gpt3message) Input should be a valid dictionary or instance of GPT3Message

**Cause:** Invalid `messages` format (mismatched brackets, missing keys).
**Solution:** Verify the JSON structure is correct.
### [​ ](#value-error-contents-is-neither-str-nor-list-of-str-input-contents) Value error, contents is neither str nor list of str. : input.contents

**Cause:** Embedding input must be a string or an array of strings.
**Solution:** Use correct input format.
### [​ ](#the-video-modality-input-does-not-meet-the-requirements-because-the-range-of-sequence-images-shoule-be-4-512-4-80) The video modality input does not meet the requirements because: the range of sequence images shoule be (4, 512)./(4,80).

**Cause:** Invalid number of images for video input.
**Solution:** Qwen3-VL/Qwen2.5-VL: 4-512 images. Other models: 4-80 images. See [Image and video understanding](/developer-guides/multimodal/vision).
### [​ ](#exceeded-limit-on-max-bytes-per-data-uri-item-10485760-multimodal-file-size-is-too-large) Exceeded limit on max bytes per data-uri item : 10485760&#x27;. / Multimodal file size is too large

**Cause:** File exceeds size limit for multimodal models.
**Solution:**

- 
Local files (Base64): Max 10 MB after encoding.


- 
Images (URL): Max 10 MB.


- 
Videos (URL): Max 2 GB (Qwen3-VL, qwen-vl-max), 1 GB (qwen-vl-plus series), or 150 MB (other models).


See [How do I compress an image or video to the required size?](/developer-guides/multimodal/vision#faq)
### [​ ](#input-should-be-cherry-serena-ethan-or-chelsie-parameters-audio-voice) Input should be &#x27;Cherry&#x27;, &#x27;Serena&#x27;, &#x27;Ethan&#x27; or &#x27;Chelsie&#x27;: parameters.audio.voice

**Cause:** Invalid `voice` parameter.
**Solution:** Use &#x27;Cherry&#x27;, &#x27;Serena&#x27;, &#x27;Ethan&#x27;, or &#x27;Chelsie&#x27;.
### [​ ](#the-image-length-and-width-do-not-meet-the-model-restrictions) The image length and width do not meet the model restrictions.

**Cause:** Image dimensions out of range.
**Solution:** Width and height must be ≥ 10 pixels. Aspect ratio must be between 1:200 and 200:1.
### [​ ](#failed-to-decode-the-image-during-the-data-inspection) Failed to decode the image during the data inspection.

**Cause:** Image decoding failed.
**Solution:** Verify the image is not corrupted and meets format requirements.
### [​ ](#the-file-format-is-illegal-and-cannot-be-opened-the-audio-format-is-illegal-and-cannot-be-opened-the-media-format-is-not-supported-or-incorrect-for-the-data-inspection) The file format is illegal and cannot be opened. / The audio format is illegal and cannot be opened. / The media format is not supported or incorrect for the data inspection.

**Cause:** Unsupported or unreadable file format.
**Solution:** Verify the file is not corrupted, extension matches actual format, and format is supported.
### [​ ](#the-input-messages-do-not-contain-elements-with-the-role-of-user) The input messages do not contain elements with the role of user.

**Cause:** No user message was passed to the model, or when calling a Qwen Cloud workflow application through an API, the parameters in the start node were passed through `user_prompt_params` instead of `biz_params`.
**Solution:** Pass a user message to the model, or pass custom parameters through `biz_params`.
### [​ ](#failed-to-download-multimodal-content-download-the-media-resource-timed-out-during-the-data-inspection-process-unable-to-download-the-media-resource-during-the-data-inspection-process) Failed to download multimodal content. / Download the media resource timed out during the data inspection process./ Unable to download the media resource during the data inspection process.

**Cause:** Failed to download the image, video, or audio, or the download timed out.

- 
**Connectivity issues:** You are using an internal network address from Object Storage Service.


- 
**Network latency:** Cross-region access causes timeouts.


- 
**Service instability:** The source storage service is slow or unreachable.


**Solution:**

- **Change the storage service**


Use a storage service in the same region as the model service. Use Object Storage Service to generate public URLs. Do not use private network addresses.

- **Adjust the transfer method**


If you cannot use a public URL, switch to a method from the table below. See [Pass local files (Base64 encoding or file path)](/developer-guides/multimodal/vision).
TypeSpecificationsDashScope SDK (Python, Java)OpenAI-compatible / DashScope HTTPImageGreater than 7 MB and less than 10 MBPass the local pathOnly public URLs are supported. Use Object Storage Service.ImageLess than 7 MBPass the local pathBase64 encodingVideoGreater than 100 MBOnly public URLs are supported. Use Object Storage Service.Only public URLs are supported. Use Object Storage Service.VideoGreater than 7 MB and less than 100 MBPass the local pathOnly public URLs are supported. Use Object Storage Service.VideoLess than 7 MBPass the local pathBase64 encodingAudioGreater than 10 MBOnly public URLs are supported. Use Object Storage Service.Only public URLs are supported. Use Object Storage Service.AudioGreater than 7 MB and less than 10 MBPass the local pathOnly public URLs are supported. Use Object Storage Service.AudioLess than 7 MBPass the local pathBase64 encoding 
 
- Base64 encoding increases the data size. Therefore, the original file must be smaller than 7 MB.

- Using Base64 encoding or a local file path can prevent server-side download timeouts and improve stability.


 
### [​ ](#url-error-check-url) url error, check url!


- **Reason 1: The model name does not match the API endpoint.** For example, you use a multimodal endpoint to call a text model, or a text endpoint to call a multimodal model.


**Solution:**

- When you use multimodal models such as qwen3.7-plus and qwen3-vl-plus through DashScope, use the MultiModalConversation.call() method or the multimodal-generation endpoint. See [Image and video understanding](/developer-guides/multimodal/vision).


 If you use the spring-ai-alibaba framework, confirm that you set the multimodal parameter [withMultiModel](https://github.com/spring-ai-alibaba/examples/blob/c66ffdec789defe4adf86b34bac0084df3b71e92/spring-ai-alibaba-multi-model-example/dashscope-multi-model/src/main/java/com/alibaba/cloud/ai/example/multi/controller/MultiModelController.java#L82). 

- When you use text models such as qwen3-max and qwen3.7-plus through DashScope, use Generation.call() or the text-generation endpoint. See [Text generation overview](/developer-guides/text-generation/quickstart).


- **Reason 2: The DashScope SDK version is outdated.** Older versions cannot resolve the correct server address for image or video generation models.


**Solution:** [Install the SDK](/api-reference/preparation/install-sdk)
### [​ ](#dont-have-authorization-to-access-the-media-resource-during-the-data-inspection-process) Don&#x27;t have authorization to access the media resource during the data inspection process.

**Cause:** The signed file URL in OSS has expired.
**Solution:** Access the file URL within its validity period.
### [​ ](#the-item-of-content-should-be-a-message-of-a-certain-modal) The item of content should be a message of a certain modal.

**Cause:** When using the DashScope SDK to call a multimodal model, the key for each element in the `content` array must be `image`, `video`, `audio`, or `text`.
**Solution:** Use the correct `content` parameter.
### [​ ](#invalid-video-file) Invalid video file.

**Cause:** The input video file is invalid.
**Solution:** Check whether the video file is corrupted or in an unsupported format.
### [​ ](#the-video-modality-input-does-not-meet-the-requirements-because-the-video-file-is-too-long) The video modality input does not meet the requirements because: The video file is too long.

**Cause:** The video duration exceeds the limit for the Qwen-VL or Qwen-Omni model.
**Solution:**

- 
Qwen2.5-VL: 2 seconds to 10 minutes.


- 
Other Qwen-VL or Qwen-Omni models: 2 to 40 seconds.


### [​ ](#field-required-xxx) Field required: xxx

**Cause:** An input parameter is missing.
**Solution:** Add the corresponding parameter based on the error message `xxx`.
### [​ ](#the-request-is-missing-required-parameters-or-in-a-wrong-format-check-the-parameters-that-you-send) The request is missing required parameters or in a wrong format, check the parameters that you send.

**Cause:** A required input parameter is missing or has a wrong format.
**Solution:** Verify that request parameters are complete and in the correct format.
### [​ ](#missing-training-files) Missing training files.

**Cause:** A parameter is incorrect, missing, or has a format issue.
### [​ ](#the-style-is-invalid) The style is invalid.

**Cause:** The `style` value is not in the allowed range.
**Solution:** Check that the `style` parameter value is valid.
### [​ ](#parameters-video-ratio-must-be-9-16-or-3-4) parameters.video_ratio must be 9:16 or 3:4.

**Cause:** The `video_ratio` parameter accepts only "9:16" or "3:4".
**Solution:** Change the `video_ratio` parameter to "9:16" or "3:4".
### [​ ](#input-json-error) input json error.

**Cause:** The input JSON is incorrect.
**Solution:** Check that the JSON format of the request is correct.
### [​ ](#read-image-error) read image error.

**Cause:** Failed to read the image.
**Solution:** Check whether the image file is corrupted or in an unsupported format.
### [​ ](#the-parameters-must-conform-to-the-specification-xxx) the parameters must conform to the specification: xxx.

**Cause:** An input parameter is out of range.
**Solution:** Correct the parameter value based on the error message `xxx`.
### [​ ](#the-size-of-person-image-and-coarse-image-are-not-the-same) The size of person image and coarse_image are not the same.

**Cause:** The resolution of `coarse_image` differs from `person_image`.
**Solution:** Set the resolution of `coarse_image` to match `person_image`.
### [​ ](#the-request-is-missing-required-parameters-or-the-parameters-are-out-of-the-specified-range-check-the-parameters-that-you-send) The request is missing required parameters or the parameters are out of the specified range, check the parameters that you send.

**Cause:** A required parameter is missing or out of bounds.
**Solution:** Correct the request parameters.
### [​ ](#image-format-error) image format error

**Cause:** The image format is incorrect.
**Solution:** The value must be an image URL or a Base64 string.
### [​ ](#no-messages-found-in-input) No messages found in input

**Cause:** The request parameters must contain the messages field.
**Solution:** See [Qwen - image editing](/developer-guides/image-generation/image-editing).
### [​ ](#invalid-image-format-or-corrupted-file) Invalid image format or corrupted file

**Cause:** The input image format is incorrect or the file is corrupted.
**Solution:** Verify the image is not corrupted and meets format requirements.
### [​ ](#download-image-failed) download image failed

**Cause:** The image cannot be downloaded.
**Solution:** Verify the file can be downloaded.
### [​ ](#messages-length-only-support-1) messages length only support 1

**Cause:** The length of the messages array can only be 1.
**Solution:** You can pass only one message. See [Qwen - image editing](/developer-guides/image-generation/image-editing).
### [​ ](#content-length-only-support-2) content length only support 2

**Cause:** The length of the content array can only be 2.
**Solution:** You can pass only one set of text and image. See [Qwen - image editing](/developer-guides/image-generation/image-editing).
### [​ ](#lack-of-image-or-text) lack of image or text

**Cause:** The request parameters are missing the image or text field.
**Solution:** See [Qwen - image editing](/developer-guides/image-generation/image-editing).
### [​ ](#num-images-per-prompt-must-be-1) num_images_per_prompt must be 1.

**Cause:** A request parameter is invalid. The parameter `n` can only be 1.
**Solution:** Set `n` to 1.
### [​ ](#input-files-format-not-supported) Input files format not supported.

**Cause:** The audio or image format is unsupported.
**Solution:** Supported audio formats: mp3, wav, aac. Supported image formats: jpg, jpeg, png, bmp, webp.
### [​ ](#failed-to-download-input-files) Failed to download input files.

**Cause:** Failed to download the input files.
**Solution:** Verify the file URL is accessible and the network is working.
### [​ ](#oss-download-error) oss download error.

**Cause:** Failed to download the input image.
**Solution:** Verify the OSS link is correct and accessible.
### [​ ](#the-image-content-does-not-comply-with-green-network-verification) The image content does not comply with green network verification.

**Cause:** The image content is not compliant.
**Solution:** Replace the image with one that passes content moderation.
### [​ ](#read-video-error) read video error.

**Cause:** Failed to read the video.
**Solution:** Check whether the video file is corrupted or in an unsupported format.
### [​ ](#the-size-of-input-image-is-too-small-or-too-large) the size of input image is too small or too large.

**Cause:** The input image size is out of range.
**Solution:** Adjust the image size to meet the API requirements.
### [​ ](#the-type-or-value-of-parameter-is-out-of-definition) The type or value of {parameter} is out of definition.

**Cause:** The parameter type or value does not meet the requirements.
**Solution:** Check the API documentation for valid parameter values.
### [​ ](#the-request-parameter-is-invalid-check-the-request-parameter) The request parameter is invalid, check the request parameter.

**Cause:** The aspect ratio parameter is invalid.
**Solution:** Use "1:1" or "3:4".
### [​ ](#request-timeout-after-23-seconds) request timeout after 23 seconds.

**Cause:** No data was sent to the service for more than 23 seconds. This occurs with [speech recognition (Paraformer)](/developer-guides/speech/asr) or [speech synthesis (CosyVoice)](/developer-guides/speech/tts).
**Solution:** Check why no data was sent for an extended period. End the task if you do not plan to send data within 23 seconds.
### [​ ](#ensure-input-text-is-valid) Ensure input text is valid.

**Cause:** When using [speech synthesis (CosyVoice)](/developer-guides/speech/tts), the text to be synthesized was not sent. Possible causes: the `text` parameter has no value, or a code exception prevented the assignment.
**Solution:** Verify that the `text` parameter is assigned and sent.
### [​ ](#missing-required-parameter-payload-model-please-follow-the-protocol) Missing required parameter &#x27;payload.model&#x27;! Please follow the protocol!

**Cause:** When using [speech synthesis (CosyVoice)](/developer-guides/speech/tts), the `model` parameter was not specified.
**Solution:** Specify the `model` parameter.
### [​ ](#tts-engine-return-error-code-418) [tts:]Engine return error code: 418

**Cause:** When using [speech synthesis (CosyVoice)](/developer-guides/speech/tts), the `voice` parameter is wrong, or the `model` version does not match the `voice` version.
**Solution:**

- 
**Check the `voice` parameter:**

For a default voice, verify it against the "voice parameter" in the [CosyVoice voice list](/api-reference/speech-synthesis/cosyvoice/voice-list).


- 
**Check version matching:** v2 models require v2 voices. v1 models require v1 voices. Do not mix them.


### [​ ](#request-voice-is-invalid) Request voice is invalid!

**Cause:** When using [speech synthesis (CosyVoice)](/developer-guides/speech/tts), the `voice` parameter was not set.
**Solution:** Assign a value to the `voice` parameter. For the [WebSocket API](/api-reference/speech-synthesis/cosyvoice/websocket-api), configure parameters in the correct JSON format.
### [​ ](#ref-images-url-and-obj-or-bg-must-be-the-same-length) ref_images_url and obj_or_bg must be the same length.

**Cause:** When using the multi-image reference feature of [Wan - general video editing](/developer-guides/video-generation/video-editing), the `ref_images_url` and `obj_or_bg` arrays have different lengths.
**Solution:** Set `ref_images_url` and `obj_or_bg` to the same length.
### [​ ](#check-input-data-style) check input data style.

**Cause:** An input parameter does not meet requirements.
**Solution:** Correct the input parameters.
### [​ ](#an-error-during-model-pre-process) An error during model pre-process.

**Cause:** The `content` field format is wrong.
**Solution:**

- Do not set content to an array type, such as `[{"type": "text", "text": "Who are you?"}]`.


### [​ ](#the-image-size-is-not-supported-for-the-data-inspection) The image size is not supported for the data inspection.

**Cause:**

- 
The image dimensions do not meet the Qwen-VL model&#x27;s requirements.


- 
The output image size exceeds 10 MB.


**Solution:**

- 
Image dimensions must meet these requirements:


Width and height must both be at least 10 pixels.


- 
Aspect ratio must not exceed 200:1 or 1:200.


- 
Adjust the image generation parameters.


### [​ ](#wrong-content-type-of-multimodal-url) Wrong Content-Type of multimodal url

**Cause** : The `Content-Type` in the URL response header is wrong.
 The Qwen-VL model supports these Content-Types: image/bmp, image/icns, image/x-icon, image/jpeg, image/jp2, image/png, image/sgi, image/tiff, and image/webp. See [Images supported by the Qwen-VL model](/developer-guides/multimodal/vision). 
**Solution** :
To view the `Content-Type` field:

- 
Open a browser such as Chrome or Firefox.


- 
Press F12 or right-click and select "Inspect" to open developer tools.


- 
Switch to the Network tab.


- 
Enter the image URL in the address bar.


- 
Find the request, then check the `Content-Type` in the "Response Headers" section.


### [​ ](#text-request-limit-violated-expected-1) Text request limit violated, expected 1.

**Cause:** In a [WebSocket API](/api-reference/speech-synthesis/cosyvoice/websocket-api) call to CosyVoice, you sent the `continue-task` instruction multiple times after setting `enable_ssml` to `true`.
**Solution:** When `enable_ssml` is `true`, send the continue-task instruction only once.
### [​ ](#ssml-text-is-not-supported-at-the-moment) SSML text is not supported at the moment!

**Cause:** The current CosyVoice model or voice does not support SSML, or SSML is used incorrectly.
**Solution:** Check the [Limitations](/api-reference/speech-synthesis/cosyvoice/websocket-api).
### [​ ](#at-least-one-of-lyrics-or-prompt-must-be-provided) At least one of &#x27;lyrics&#x27; or &#x27;prompt&#x27; must be provided.

**Cause:** Neither the `lyrics` nor the `prompt` parameter was provided when using the Fun-Music model.
**Solution:** Provide at least one of the `lyrics` or `prompt` parameter in your request.
### [​ ](#lyrics-content-is-illegal-and-cannot-be-used-for-music-generation) Lyrics content is illegal and cannot be used for music generation.

**Cause:** The lyrics content did not pass the content moderation check when using the Fun-Music model. The lyrics may contain infringing content.
**Solution:** Modify the lyrics to remove any infringing or non-compliant content, then try again.
## [​ ](#400-invalid-request-error-invalid-value) 400-invalid_request_error-invalid_value

### [​ ](#1-is-lesser-than-the-minimum-of-0-seed-seed-must-be-integer) -1 is lesser than the minimum of 0 - &#x27;seed&#x27;/&#x27;seed&#x27; must be Integer

**Cause:** When using the OpenAI-compatible protocol, `seed` is not in the range [0, 231-1].
**Solution:** Set `seed` to a value in [0, 231-1].
## [​ ](#400-invalid-request-error) 400-invalid_request_error

### [​ ](#you-must-provide-a-model-parameter) you must provide a model parameter.

**Cause:** The `model` parameter was not provided in the request.
**Solution:** Add the `model` parameter to the request.
## [​ ](#400-invalidparameter-notsupportenablethinking) 400-InvalidParameter.NotSupportEnableThinking

### [​ ](#the-model-xxx-does-not-support-enable-thinking) The model xxx does not support enable_thinking.

**Cause:** The model does not support `enable_thinking`.
**Solution:** Remove `enable_thinking` from the request, or use a model that supports thinking mode.
## [​ ](#400-invalid-value) 400-invalid_value

### [​ ](#the-requested-voice-xxx-is-not-supported) The requested voice &#x27;xxx&#x27; is not supported.

**Cause:** The selected voice belongs to a different model.
**Solution:** Set `target_model` for voice cloning to match the `model` for speech synthesis.
## [​ ](#400-arrearage) 400-Arrearage

### [​ ](#access-denied-make-sure-your-account-is-in-good-standing) Access denied, make sure your account is in good standing.

**Cause:** Account has an overdue payment.
**Solution:** Check [Billing Overview](https://home.qwencloud.com/billing/overview). If no payment is overdue, verify the API key belongs to your account. If overdue, settle the outstanding balance and wait for it to update before retrying.
### [​ ](#api-provider-returned-a-billing-error-your-api-key-has-run-out-of-credits-or-has-an-insufficient-balance-check-your-providers-billing-dashboard-and-top-up-or-switch-to-a-different-api-key) API provider returned a billing error -- your API key has run out of credits or has an insufficient balance. Check your provider&#x27;s billing dashboard and top up or switch to a different API key.

**Cause:** When you call Qwen Cloud through a third-party client such as OpenClaw, if the underlying account has an overdue payment or insufficient balance, the client aggregates the server-side billing error into this message. The underlying server-side error code is Arrearage (the **Access denied, make sure your account is in good standing.** message above). It is not a billing issue of the client itself.
**Solution:** Check [Billing Overview](https://home.qwencloud.com/billing/overview) for overdue payments. If overdue, settle the balance and wait for it to update before retrying. If no payment is overdue, verify the API key belongs to your account.
## [​ ](#400-datainspectionfailed-data-inspection-failed) 400-DataInspectionFailed/data_inspection_failed

### [​ ](#input-or-output-data-may-contain-inappropriate-content-input-data-may-contain-inappropriate-content-output-data-may-contain-inappropriate-content) Input or output data may contain inappropriate content. / Input data may contain inappropriate content. / Output data may contain inappropriate content.

**Cause:** Content blocked by moderation.
**Solution:** Modify the input and retry.
### [​ ](#input-xxx-data-may-contain-inappropriate-content) Input xxx data may contain inappropriate content.

**Cause:** The input data, such as the prompt or image, may contain sensitive content.
**Solution:** Check and modify the input, then try again.
### [​ ](#qwen-rejected-the-input-image-before-model-inference-no-actual-ad-porn-ocr-result-was-produced) Qwen rejected the input image before model inference; no actual ad/porn/OCR result was produced.

**Cause:** Before the input image entered model inference, it was flagged by the content-safety pre-check as suspected to contain sensitive or non-compliant content and was therefore blocked. The model did not run inference on the image, so no recognition, OCR, or analysis result is returned. This is a pre-inference content-compliance check on multimodal input images.
**Solution:** Replace or modify the input image and retry. If the image is confirmed to be compliant but is still consistently blocked, submit a ticket for further verification.
## [​ ](#400-apiconnectionerror) 400-APIConnectionError

### [​ ](#connection-error) Connection error.

**Cause:** Network connection error (often proxy-related).
**Solution:** Disable or restart proxy.
## [​ ](#400-invalidfile-downloadfailed) 400-InvalidFile.DownloadFailed

### [​ ](#the-audio-file-cannot-be-downloaded) The audio file cannot be downloaded.

**Cause:** Failed to download audio file (Paraformer).
**Solution:** Verify the URL is publicly accessible.
## [​ ](#400-invalidfile-audiolengtherror) 400-InvalidFile.AudioLengthError

### [​ ](#audio-length-must-be-between-1s-and-300s) Audio length must be between 1s and 300s.

**Cause:** Audio duration out of range.
**Solution:** Ensure duration is 1-300 seconds.
### [​ ](#audio-length-must-be-between-1s-and-180s) Audio length must be between 1s and 180s.

**Cause:** Audio duration out of range.
**Solution:** Ensure duration is 1-180 seconds.
## [​ ](#400-invalidfile-nohuman) 400-InvalidFile.NoHuman

### [​ ](#the-input-image-has-no-human-body-upload-another-image-with-single-person) The input image has no human body. Upload another image with single person.

**Cause:** There is no person in the input image, or no face was detected.
**Solution:** Upload a photo of a single person.
## [​ ](#400-invalidfile-bodyproportion) 400-InvalidFile.BodyProportion

### [​ ](#the-proportion-of-the-detected-person-in-the-picture-is-too-large-or-too-small-upload-another-image) The proportion of the detected person in the picture is too large or too small, upload another image.

**Cause:** The person&#x27;s proportion in the image does not meet the requirement.
**Solution:** Upload an image with a suitable person proportion.
## [​ ](#400-invalidfile-facepose) 400-InvalidFile.FacePose

### [​ ](#the-pose-of-the-detected-face-is-invalid-upload-another-image-with-whole-face-and-expected-orientation) The pose of the detected face is invalid, upload another image with whole face and expected orientation.

**Cause:** The facial pose does not meet the requirement. The face must be visible and the head must not be severely skewed.
**Solution:** Upload an image with a visible, properly oriented face.
### [​ ](#the-pose-of-the-detected-face-is-invalid-upload-another-image-with-the-expected-oriention) The pose of the detected face is invalid, upload another image with the expected oriention.

**Cause:** The head orientation is too skewed.
**Solution:** Upload an image where the face is not skewed.
### [​ ](#the-pose-of-the-detected-face-is-invalid-upload-another-image-with-the-expected-orientation) The pose of the detected face is invalid, upload another image with the expected orientation.

**Cause:** The head orientation is too skewed.
**Solution:** Upload an image where the face is not skewed.
## [​ ](#400-invalidfile-resolution) 400-InvalidFile.Resolution

### [​ ](#the-image-resolution-is-invalid-make-sure-that-the-largest-length-of-image-is-smaller-than-7000-and-the-smallest-length-of-image-is-larger-than-400) The image resolution is invalid, make sure that the largest length of image is smaller than 7000, and the smallest length of image is larger than 400.

**Cause:** The image resolution does not meet the requirement.
**Solution:** Resolution must be between 400*400 and 7,000*7,000.
### [​ ](#the-image-resolution-is-invalid-make-sure-that-the-largest-length-of-image-is-smaller-than-4096-and-the-smallest-length-of-image-is-larger-than-224) The image resolution is invalid, make sure that the largest length of image is smaller than 4096, and the smallest length of image is larger than 224.

**Cause:** The image resolution does not meet the requirement.
**Solution:** The longest side must be less than 4,096 pixels. The shortest side must be greater than 224 pixels.
### [​ ](#the-image-resolution-is-invalid-make-sure-that-the-largest-length-of-image-is-smaller-than-xxx-and-the-smallest-length-of-image-is-larger-than-yyy) The image resolution is invalid, make sure that the largest length of image is smaller than xxx, and the smallest length of image is larger than yyy.

**Cause:** The image resolution does not meet the requirement.
**Solution:** Resolution must be between yyy*yyy and xxx*xxx.
### [​ ](#the-image-resolution-is-invalid-make-sure-that-the-aspect-ratio-is-smaller-than-xxx-and-largest-length-of-image-is-smaller-than-yyy) The image resolution is invalid, make sure that the aspect ratio is smaller than xxx, and largest length of image is smaller than yyy.

**Cause:** The image resolution does not meet the requirement.
**Solution:** The aspect ratio must be less than xxx, and resolution must not exceed yyy*yyy.
### [​ ](#invalid-video-resolution-the-height-or-width-of-video-must-be-xxx-yyy) Invalid video resolution. The height or width of video must be xxx ~ yyy.

**Cause:** The video resolution does not meet the requirement.
**Solution:** Each side must be between xxx and yyy pixels.
## [​ ](#400-invalidfile-fps) 400-InvalidFile.FPS

### [​ ](#invalid-video-fps-the-video-fps-must-be-15-60) Invalid video FPS. The video FPS must be 15 ~ 60.

**Cause:** The video frame rate does not meet the requirement.
**Solution:** Frame rate must be between 15 and 60 fps.
## [​ ](#400-invalidfile-value) 400-InvalidFile.Value

### [​ ](#the-value-of-the-image-is-invalid-upload-other-clearer-image) The value of the image is invalid, upload other clearer image.

**Cause:** The image is too dark.
**Solution:** Upload a clearer image with adequate lighting.
## [​ ](#400-invalidfile-frontbody) 400-InvalidFile.FrontBody

### [​ ](#the-pose-of-the-detected-person-is-invalid-upload-another-image-with-the-front-view) The pose of the detected person is invalid, upload another image with the front view.

**Cause:** The person has their back to the camera.
**Solution:** Upload an image where the person faces the camera.
## [​ ](#400-invalidfile-fullface) 400-InvalidFile.FullFace

### [​ ](#the-pose-of-the-detected-face-is-invalid-upload-another-image-with-whole-face) The pose of the detected face is invalid, upload another image with whole face.

**Cause:** The face is not fully visible.
**Solution:** Upload an image where the entire face is visible and not obstructed.
## [​ ](#400-invalidfile-facenotmatch) 400-InvalidFile.FaceNotMatch

### [​ ](#there-are-no-matched-face-in-the-video-with-the-provided-reference-image) There are no matched face in the video with the provided reference image.

**Cause:** The face in the reference image does not match the face in the video.
**Solution:** Ensure the reference image and video contain the same person.
## [​ ](#400-invalidfile-content) 400-InvalidFile.Content

### [​ ](#the-first-frame-of-input-video-has-no-human-body-please-choose-another-clip) The first frame of input video has no human body. Please choose another clip.

**Cause:** The first frame must contain a person.
**Solution:** Select a video clip with a person in the first frame.
### [​ ](#the-human-is-too-small-in-the-first-frame-of-input-video-please-choose-another-clip) The human is too small in the first frame of input video. Please choose another clip.

**Cause:** The person in the first frame is too small.
**Solution:** Select a video with a larger person in the first frame.
### [​ ](#the-human-is-not-clear-in-the-first-frame-of-input-video-please-choose-another-clip) The human is not clear in the first frame of input video. Please choose another clip.

**Cause:** The person in the first frame is not clear.
**Solution:** Select a video with a clear person in the first frame.
### [​ ](#the-input-image-has-no-human-body-or-multi-human-bodies-upload-another-image-with-single-person) The input image has no human body or multi human bodies. Upload another image with single person.

**Cause:** The input image has no person or multiple people.
**Solution:** Upload a photo of a single person.
### [​ ](#the-input-image-has-no-human-body-or-has-unclear-human-body-upload-another-image) The input image has no human body or has unclear human body. Upload another image.

**Cause:** The human body in the input image is incomplete or absent.
**Solution:** Upload an image with a complete, clear human body.
### [​ ](#the-input-image-has-multi-human-bodies-upload-another-image-with-single-person) The input image has multi human bodies. Upload another image with single person.

**Cause:** The input image has multiple people.
**Solution:** Upload a photo of a single person.
## [​ ](#400-invalidfile-fullbody) 400-InvalidFile.FullBody

### [​ ](#the-human-is-not-fullbody-in-the-first-frame-of-input-video-please-choose-another-clip) The human is not fullbody in the first frame of input video. Please choose another clip.

**Cause:** The person in the first frame is not a full body shot.
**Solution:** The entire body must be visible in the first frame.
### [​ ](#the-pose-of-the-detected-person-is-invalid-upload-another-image-with-whole-body-or-change-the-ratio-parameter-to-1-1) The pose of the detected person is invalid, upload another image with whole body, or change the ratio parameter to 1:1.

**Cause:** The person&#x27;s pose does not meet the requirement.
**Solution:** Upload an image that meets the requirements. For a head shot, the entire head must be visible. For a half-body shot, the body from the hips up must be visible. Or change the aspect ratio to 1:1.
## [​ ](#400-invalidfile-bodypose) 400-InvalidFile.BodyPose

### [​ ](#the-pose-of-the-detected-person-is-invalid-upload-another-image-with-whole-body-and-expected-orientation) The pose of the detected person is invalid, upload another image with whole body and expected orientation.

**Cause:** The person&#x27;s pose does not meet the requirement.
**Solution:** Upload an image where the shoulders and ankles are visible, the person is not sitting or turned away, and the orientation is not severely skewed.
## [​ ](#400-invalidfile-size) 400-InvalidFile.Size

### [​ ](#invalid-file-size-the-video-file-size-must-be-less-than-200mb-and-the-audio-file-size-must-be-less-than-15mb) Invalid file size. The video file size must be less than 200MB, and the audio file size must be less than 15MB.

**Cause:** The file size exceeds the limit.
**Solution:** Video must be under 200 MB. Audio must be under 15 MB.
### [​ ](#invalid-file-size-the-image-file-size-must-be-smaller-than-5mb) Invalid file size, The image file size must be smaller than 5MB.

**Cause:** The file size exceeds the limit.
**Solution:** Image must be under 5 MB.
### [​ ](#invalid-file-size-the-video-audio-image-file-size-must-be-less-than-xxxmb) Invalid file size. The video/audio/image file size must be less than xxxMB.

**Cause:** The file size exceeds the limit.
**Solution:** Keep the file under the specified size limit.
## [​ ](#400-invalidfile-duration) 400-InvalidFile.Duration

### [​ ](#invalid-file-duration-the-file-duration-must-be-xxx-s-yyy-s) Invalid file duration. The file duration must be xxx s ~ yyy s.

**Cause:** The file duration is out of range.
**Solution:** Duration must be between xxx and yyy seconds.
## [​ ](#400-invalidfile-imagesize) 400-InvalidFile.ImageSize

### [​ ](#the-size-of-image-is-beyond-limit) The size of image is beyond limit.

**Cause:** The image size exceeds the limit.
**Solution:** Aspect ratio must not exceed 2, and the longest side must not exceed 4096.
## [​ ](#400-invalidfile-aspectratio) 400-InvalidFile.AspectRatio

### [​ ](#invalid-file-ratio-the-file-aspect-ratio-height-width-must-be-between-3-1-and-1-3) Invalid file ratio. The file aspect ratio (height/width) must be between 3:1 and 1:3.

**Cause:** The file aspect ratio is out of range.
**Solution:** Aspect ratio must be between 3:1 and 1:3.
### [​ ](#invalid-file-ratio-the-file-aspect-ratio-height-width-must-be-between-2-0-and-0-5) Invalid file ratio. The file aspect ratio (height/width) must be between 2.0 and 0.5.

**Cause:** The file aspect ratio is out of range.
**Solution:** Aspect ratio must be between 2.0 and 0.5.
## [​ ](#400-invalidfile-openerror) 400-InvalidFile.Openerror

### [​ ](#invalid-file-cannot-open-file-as-video-audio-image) Invalid file, cannot open file as video/audio/image.

**Cause:** The file cannot be opened.
**Solution:** Check whether the file is corrupted or in an unsupported format.
## [​ ](#400-invalidfile-template-content) 400-InvalidFile.Template.Content

### [​ ](#invalid-template-content) Invalid template content.

**Cause:** You lack permission for the action template, or the template content is invalid.
**Solution:** Check the template permissions and content.
## [​ ](#400-invalidfile-format) 400-InvalidFile.Format

### [​ ](#invalid-file-format-the-request-file-format-is-one-of-the-following-types-mp4-avi-mov-mp3-wav-aac-jpeg-jpg-png-bmp-and-webp) Invalid file format，the request file format is one of the following types: MP4, AVI, MOV, MP3, WAV, AAC, JPEG, JPG, PNG, BMP, and WEBP.

**Cause:** The file format is unsupported.
**Solution:** Supported formats: video (mp4, avi, mov), audio (mp3, wav, aac), image (jpg, jpeg, png, bmp, webp).
## [​ ](#400-invalidfile-multihuman) 400-InvalidFile.MultiHuman

### [​ ](#the-input-image-has-multi-human-bodies-upload-another-image-with-single-person-2) The input image has multi human bodies. Upload another image with single person.

**Cause:** The input image has multiple people.
**Solution:** Upload a photo of a single person.
## [​ ](#400-invalidperson) 400-InvalidPerson

### [​ ](#the-input-image-has-no-human-body-or-multi-human-bodies-upload-another-image-with-single-person-2) The input image has no human body or multi human bodies. Upload another image with single person.

**Cause:** The input image has no person or multiple people.
**Solution:** Upload a photo of a single person.
## [​ ](#400-invalidparameter-datainspection) 400-InvalidParameter.DataInspection

### [​ ](#unable-to-download-the-media-resource-during-the-data-inspection-process) Unable to download the media resource during the data inspection process.

**Cause:** The download of an image or audio file timed out.
**Solution:** Store files in an OSS bucket in the same region as the model service.
## [​ ](#400-flownotpublished) 400-FlowNotPublished

### [​ ](#flow-has-not-published-yet-please-publish-flow-and-try-again) Flow has not published yet, please publish flow and try again.

**Cause:** The flow is not published.
**Solution:** Publish the flow and try again.
## [​ ](#400-invalidimage-imagesize) 400-InvalidImage.ImageSize

### [​ ](#the-size-of-image-is-beyond-limit-2) The size of image is beyond limit.

**Cause:** The image size exceeds the limit.
**Solution:** The image aspect ratio must not be greater than 2, and the longest side must not be greater than 4096.
## [​ ](#400-invalidimage-nohumanface) 400-InvalidImage.NoHumanFace

### [​ ](#no-human-face-detected) No human face detected.

**Cause:** No human face was detected. This occurs only with asynchronous query APIs for generation tasks.
**Solution:** Upload an image that contains a clear human face.
## [​ ](#400-invalidimageresolution) 400-InvalidImageResolution

### [​ ](#the-input-image-resolution-is-too-large-or-small) The input image resolution is too large or small.

**Cause:** The input image resolution is too large or too small.
**Solution:** Resolution must be between 256 x 256 and 5760 x 3240 pixels.
## [​ ](#400-invalidimageformat) 400-InvalidImageFormat

### [​ ](#the-input-image-is-in-invalid-format) The input image is in invalid format.

**Cause:** The image format is unsupported.
**Solution:** Use JPEG, PNG, JPG, BMP, or WEBP format.
## [​ ](#400-invalidurl) 400-InvalidURL

### [​ ](#invalid-url-provided-in-your-request) Invalid URL provided in your request.

**Cause:** The URL is invalid.
**Solution:** Use a valid URL.
### [​ ](#required-url-is-missing-or-invalid-check-the-request-url) Required URL is missing or invalid, check the request URL.

**Cause:** The input URL is invalid or missing.
**Solution:** Provide the correct URL.
### [​ ](#the-request-url-is-invalid-make-sure-the-url-is-correct-and-is-an-image) The request URL is invalid, make sure the url is correct and is an image.

**Cause:** The input URL is invalid.
**Solution:** Verify the URL is correct and points to an image file.
### [​ ](#the-input-audio-is-longer-than-xxs) The input audio is longer than xxs.

**Cause:** The input audio exceeds the maximum duration of xx seconds.
**Solution:** Trim the audio to under xx seconds.
### [​ ](#file-size-is-larger-than-15mb) File size is larger than 15MB.

**Cause:** The input audio exceeds 15 MB.
**Solution:** Compress the audio to under 15 MB.
### [​ ](#file-type-is-not-supported-allowed-types-are-wav-mp3) File type is not supported. Allowed types are: .wav, .mp3.

**Cause:** The audio format is unsupported.
**Solution:** Use wav or mp3 format.
### [​ ](#the-request-url-is-invalid-check-the-request-url-is-available-and-the-request-image-format-is-one-of-the-following-types-jpeg-jpg-png-bmp-and-webp) The request URL is invalid, check the request URL is available and the request image format is one of the following types: JPEG, JPG, PNG, BMP, and WEBP.

**Cause:** The image is inaccessible or the format is unsupported.
**Solution:** Verify the URL is accessible and the image format is JPEG, JPG, PNG, BMP, or WEBP.
## [​ ](#400-invalidimage-fileformat) 400-InvalidImage.FileFormat

### [​ ](#invalid-image-type-ensure-the-uploaded-file-is-a-valid-image) Invalid image type. Ensure the uploaded file is a valid image.

**Cause:** The image format is unsupported.
**Solution:** Use JPG, JPEG, PNG, BMP, or WEBP.
## [​ ](#400-invalidurl-connectionrefused) 400-InvalidURL.ConnectionRefused

### [​ ](#connection-to-xxx-refused-please-provide-available-url) Connection to xxx refused, please provide available URL.

**Cause:** The download was refused.
**Solution:** Provide an available URL.
## [​ ](#400-invalidurl-timeout) 400-InvalidURL.Timeout

### [​ ](#download-xxx-timeout-check-network-connection) Download xxx timeout, check network connection.

**Cause:** The download timed out.
**Solution:** Check network connectivity.
## [​ ](#400-badrequestexception) 400-BadRequestException

### [​ ](#invalid-part-type) Invalid part type.

**Cause:** Document processing does not support this file type.
**Solution:** Upload a supported file type.
## [​ ](#400-badrequest-emptyinput) 400-BadRequest.EmptyInput

### [​ ](#required-input-parameter-missing-from-request) Required input parameter missing from request.

**Cause:** The `input` parameter was not added to the request.
**Solution:** Add the `input` parameter to the request.
## [​ ](#400-badrequest-emptyparameters) 400-BadRequest.EmptyParameters

### [​ ](#required-parameter-parameters-missing-from-request) Required parameter "parameters" missing from request.

**Cause:** The `parameters` parameter was not added to the request.
**Solution:** Add the `parameters` parameter to the request.
## [​ ](#400-badrequest-emptymodel) 400-BadRequest.EmptyModel

### [​ ](#required-parameter-model-missing-from-request) Required parameter "model" missing from request.

**Cause:** The `model` parameter was not provided in the request.
**Solution:** Add the `model` parameter to the request.
## [​ ](#400-badrequest-illegalinput) 400-BadRequest.IllegalInput

### [​ ](#the-input-parameter-requires-json-format) The input parameter requires json format.

**Cause:** The input format is not valid JSON.
**Solution:** Verify the input is standard JSON.
## [​ ](#400-badrequest-inputdownloadfailed) 400-BadRequest.InputDownloadFailed

### [​ ](#failed-to-download-the-input-file-xxx) Failed to download the input file: xxx.

**Cause:** Failed to download the input file. Possible causes: download timeout, download failure, or file exceeding the size quota.
**Solution:** Troubleshoot based on the error message `xxx`.
### [​ ](#failed-to-download-the-input-file) Failed to download the input file.

**Cause:** When using Qwen-TTS for voice cloning, the server failed to download the audio.
**Solution:** Verify the audio file can be downloaded and is under 10 MB.
## [​ ](#400-badrequest-unsupportedfileformat) 400-BadRequest.UnsupportedFileFormat

### [​ ](#input-file-format-is-not-supported) Input file format is not supported.

**Cause:** The file format is unsupported.
**Solution:** Use a supported format.
## [​ ](#400-badrequest-toolarge) 400-BadRequest.TooLarge

### [​ ](#payload-too-large) Payload Too Large.

**Cause:** The file size exceeds the limit.
**Solution:**

- 
When the "purpose" parameter is "file-extract", documents cannot exceed 150 MB and images cannot exceed 20 MB.


- 
When the "purpose" parameter is "batch", files cannot exceed 500 MB. Split the file and [upload the files](/api-reference/platform-api/file) in batches.


## [​ ](#400-badrequest-resourcenotexist) 400-BadRequest.ResourceNotExist

### [​ ](#the-required-resource-not-exist) The Required resource not exist.

**Cause:** When calling update, query, or delete interfaces for [CosyVoice voice cloning](/developer-guides/speech/tts), the corresponding voice does not exist.
## [​ ](#400-throttling-allocationquota) 400-Throttling.AllocationQuota

### [​ ](#free-allocated-quota-exceeded) Free allocated quota exceeded.

**Cause:** The number of custom vocabularies exceeds the limit (default: 10 per account).
**Solution:** Delete some custom vocabularies.
### [​ ](#maximum-voice-storage-limit-exceeded-please-delete-existing-voices) Maximum voice storage limit exceeded, please delete existing voices.

**Cause:** When using Qwen-TTS for voice cloning, your account exceeds the voice limit.
**Solution:** Delete some voices or request a quota increase.
## [​ ](#400-invalidgarment) 400-InvalidGarment

### [​ ](#missing-clothing-image-please-input-at-least-one-top-garment-or-bottom-garment-image) Missing clothing image.Please input at least one top garment or bottom garment image.

**Cause:** A clothing image is missing.
**Solution:** Provide at least one image of a top garment (top_garment_url) or a bottom garment (bottom_garment_url).
## [​ ](#400-invalidschema) 400-InvalidSchema

### [​ ](#database-schema-is-invalid-for-text2sql) Database schema is invalid for text2sql.

**Cause:** Database schema information was not provided.
**Solution:** Enter the database schema information.
## [​ ](#400-invalidschemaformat) 400-InvalidSchemaFormat

### [​ ](#database-schema-format-is-invalid-for-text2sql) Database schema format is invalid for text2sql.

**Cause:** The data table information format is invalid.
**Solution:** Check and correct the format of the data table information.
## [​ ](#400-invalidinputlength) 400-InvalidInputLength

### [​ ](#the-image-resolution-is-invalid-make-sure-that-the-largest-length-of-image-is-smaller-than-4096-and-the-smallest-length-of-image-is-larger-than-150-and-the-size-of-image-ranges-from-5kb-to-5mb) The image resolution is invalid, make sure that the largest length of image is smaller than 4096, and the smallest length of image is larger than 150. and the size of image ranges from 5KB to 5MB.

**Cause:** The image dimensions or file size do not meet the requirements.
**Solution:** See [Input image requirements](/developer-guides/multimodal/vision).
## [​ ](#400-faqruleblocked) 400-FaqRuleBlocked

### [​ ](#input-or-output-data-is-blocked-by-faq-rule) Input or output data is blocked by faq rule.

**Cause:** The FAQ rule intervention module was hit.
## [​ ](#400-clientdisconnect) 400-ClientDisconnect

### [​ ](#client-disconnected-before-task-finished) Client disconnected before task finished!

**Cause:** The client disconnected before the task finished. This occurs with speech synthesis or recognition services.
**Solution:** Do not disconnect before the task completes.
## [​ ](#400-serviceunavailableerror) 400-ServiceUnavailableError

### [​ ](#role-must-be-user-or-assistant-and-content-length-must-be-greater-than-0) Role must be user or assistant and Content length must be greater than 0.

**Cause:** The input content length is 0 or the `role` is incorrect.
**Solution:** Verify the input content length is greater than 0 and the `role` parameter meets API requirements.
## [​ ](#400-ipinfringementsuspect) 400-IPInfringementSuspect

### [​ ](#input-data-is-suspected-of-being-involved-in-ip-infringement) Input data is suspected of being involved in IP infringement.

**Cause:** The input data is suspected of intellectual property infringement.
**Solution:** Check the input to ensure it does not contain content that poses an infringement risk.
## [​ ](#400-unsupportedoperation) 400-UnsupportedOperation

### [​ ](#the-operation-is-unsupported-on-the-referee-object) The operation is unsupported on the referee object.

**Cause:** The associated object does not support the operation.
**Solution:** Verify the operation object and type match.
## [​ ](#400-customroleblocked) 400-CustomRoleBlocked

### [​ ](#input-or-output-data-may-contain-inappropriate-content-with-custom-rule) Input or output data may contain inappropriate content with custom rule.

**Cause:** The request or response did not pass the custom policy.
**Solution:** Check the content or adjust the custom policy.
## [​ ](#400-audio-preprocesserror) 400-Audio.PreprocessError

### [​ ](#audio-preprocess-error) Audio preprocess error.

**Cause:** When using Qwen-TTS for voice cloning, audio preprocessing failed. Possible causes: the `text` parameter differs significantly from the audio, the effective speech is too short, or the audio contains no sound.
**Solution:** Adjust the `text` parameter. If the problem persists, refer to our recording guide and record the audio again.
### [​ ](#no-segments-meet-minimum-duration-requirement) No segments meet minimum duration requirement

**Cause:** When using Qwen-TTS for voice cloning, the effective speech in the audio is too short.
**Solution:** Follow the recording guide and record again.
## [​ ](#400-badrequest-voicenotfound) 400-BadRequest.VoiceNotFound

### [​ ](#voice-s-not-found) Voice &#x27;%s&#x27; not found.

**Cause:** When using Qwen-TTS for voice cloning, you called the delete voice API, but the voice is already deleted or does not exist.
**Solution:** Verify the `voice` parameter is correct.
## [​ ](#400-audio-decodererror) 400-Audio.DecoderError

### [​ ](#decoder-audio-file-failed) Decoder audio file failed.

**Cause:** Audio decoding failed when using Qwen-TTS or CosyVoice for voice cloning.
**Solution:** Check the audio file for corruption and verify it meets the required format. Qwen-TTS has specific requirements. CosyVoice supports WAV (16-bit), MP3, or M4A.
## [​ ](#400-audio-audiorateerror) 400-Audio.AudioRateError

### [​ ](#file-sample-rate-unsupported) File sample rate unsupported.

**Cause:** The audio sample rate is unsupported when using Qwen-TTS or CosyVoice for voice cloning.
**Solution:** Use a sample rate of at least 24000 Hz.
## [​ ](#400-audio-durationlimiterror) 400-Audio.DurationLimitError

### [​ ](#audio-duration-exceeds-maximum-allowed-limit) Audio duration exceeds maximum allowed limit.

**Cause:** The audio for Qwen-TTS voice cloning is too long.
**Solution:** Keep the audio under 60 seconds.
## [​ ](#401-invalidapikey-invalid-api-key) 401-InvalidApiKey/invalid_api_key

### [​ ](#invalid-api-key-provided-incorrect-api-key-provided) Invalid API-key provided. / Incorrect API key provided.

**Cause:** The API key is incorrect.
**Solution:** Common causes and solutions:

- 
**Incorrect environment variable**


**Incorrect:**`api_key=os.getenv("sk-xxx") `. `sk-xxx` is treated as an environment variable name, not a key.


- 
**Correct:**

**If the environment variable is set:** Use `api_key=os.getenv("DASHSCOPE_API_KEY")`.


 Make sure you have set `DASHSCOPE_API_KEY` 

- **If the environment variable is not set:** Use `api_key = "sk-xxx"`.


 Do not use this method in the production environment. 

- 
**Incorrect entry** : API keys start with `sk-`. You may have entered a key from another provider.


- 
**Coding Plan dedicated API key** : The Coding Plan provides a dedicated API key (starting with `sk-sp-`). This key **must be used with the dedicated API endpoint** (such as [https://coding-intl.dashscope.aliyuncs.com/v1](https://coding-intl.dashscope.aliyuncs.com/v1)). Update both the API key and the base URL.


- 
**Incorrect endpoint or API key** : Ensure your base URL is set to `https://dashscope-intl.aliyuncs.com/compatible-mode/v1` (OpenAI compatible) or `https://dashscope-intl.aliyuncs.com/api/v1` (DashScope), and that your API key was created in Qwen Cloud. Verify your API key on the [API Keys](https://home.qwencloud.com/api-keys) page.


- 
**Tool compatibility issues** : A third-party tool is misconfigured. For example, [Dify](/developer-guides/clients-and-developer-tools/dify)&#x27;s latest plugin version is unstable -- install an older TONGYI plugin version. Or in older versions of [Cline](/developer-guides/clients-and-developer-tools/cline), **API Provider** must be **OpenAI Compatible**, not **Alibaba Qwen**.


If none of the above apply, the API key may have been deleted. Get a new key and retry.
## [​ ](#401-not-authorized) 401-NOT AUTHORIZED

### [​ ](#access-denied-either-you-are-not-authorized-to-access-this-workspace-or-the-workspace-does-not-exist-please-nverify-the-workspace-configuration-ncheck-your-api-endpoint-settings-ensure-you-are-targeting-the-correct-environment) Access denied: Either you are not authorized to access this workspace, or the workspace does not exist. Please:\nVerify the workspace configuration.\nCheck your API endpoint settings. Ensure you are targeting the correct environment.

**Cause:**

- 
The WorkspaceId is invalid, or the account is not a member of this workspace.


- 
The endpoint is incorrect.


**Solution:**

- Confirm the WorkspaceId is correct and the account is a workspace member, then retry.


## [​ ](#401-invalid-access-token-or-token-expired) 401-invalid access token or token expired

### [​ ](#invalid-access-token-or-token-expired) invalid access token or token expired.

**Possible cause:** The Base URL for the Coding Plan or another plan was used with a Token Plan API key.
**Solution:** Use the dedicated Base URL for Token Plan:

- Anthropic compatible endpoint: `https://token-plan.ap-southeast-1.maas.aliyuncs.com/apps/anthropic`

- OpenAI compatible endpoint: `https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1`


## [​ ](#403-accessdenied-access-denied) 403-AccessDenied/access_denied

### [​ ](#current-user-api-does-not-support-asynchronous-calls) Current user api does not support asynchronous calls.

**Cause:** The API does not support asynchronous invocation.
**Solution:** Remove `X-DashScope-Async` from the request header, or set its value to `disable`.
### [​ ](#current-user-api-does-not-support-synchronous-calls) current user api does not support synchronous calls.

**Cause:** The API does not support synchronous calls.
**Solution:** Set `X-DashScope-Async: enable` in the request header.
### [​ ](#invalid-according-to-policy-policy-expired) Invalid according to Policy: Policy expired.

**Cause:** The file upload credential has expired.
**Solution:** Call the [file upload API](/api-reference/platform-api/file/upload-file) to generate a new credential.
### [​ ](#access-denied) Access denied.

**Reason:** Access denied. You lack the required permissions, or the model&#x27;s free quota is exhausted and it does not support paid usage.
## [​ ](#403-accessdenied-unpurchased) 403-AccessDenied.Unpurchased

### [​ ](#access-to-model-denied-make-sure-you-are-eligible-for-using-the-model) Access to model denied. Make sure you are eligible for using the model.

**Cause:** You have not logged in to Qwen Cloud, or you do not have a Qwen Cloud account.
**Solution:** Sign in to Qwen Cloud:

- 
**Register an account** : If you do not have a Qwen Cloud account, [register](https://account.alibabacloud.com/sso/login.htm?response_type=code&client_id=qwencloud&scope=openid&state=fb705aee7c3e462d9f7d66b0d173b7b7&redirect_uri=https%3A%2F%2Faccount.qwencloud.com%2Fsso%2FssoLogin&accounttraceid=2c18fe007ff24898a5dbf4eae036a1d8nvpg&cspNonce=bzpy020jk1) one.


- 
**Sign in:** Go to [Qwen Cloud](https://home.qwencloud.com/try-ai) with your Qwen Cloud account.


## [​ ](#403-model-accessdenied) 403-Model.AccessDenied

### [​ ](#model-access-denied) Model access denied.

**Cause:** You lack permission to call this model.
**Solution:**

- 
**Calling standard models:** When using a sub-workspace API key to call standard models (such as `qwen-plus`), the sub-workspace must have calling permissions.


- 
**Calling custom models:** After successful deployment, custom models can only be called with the API key of their workspace and do not require model calling authorization.


## [​ ](#403-app-accessdenied) 403-App.AccessDenied

### [​ ](#app-access-denied) App access denied.

**Cause:** You lack permission to access the model.
**Solution:**

- 
Confirm that the workspace has access authorization.


- 
Verify the API KEY is correct.


- 
If Claude Code reports an error, use the API key of the default workspace.


## [​ ](#403-workspace-accessdenied) 403-Workspace.AccessDenied

### [​ ](#workspace-access-denied) Workspace access denied.

**Cause:** You lack permission to access the workspace&#x27;s model.
**Solution:**

- 
Verify you are a member of the workspace.


- 
Use the Qwen Cloud account&#x27;s API key, which has permissions for all workspaces.


## [​ ](#403-allocationquota-freetieronly) 403-AllocationQuota.FreeTierOnly

### [​ ](#the-free-tier-of-the-model-has-been-exhausted-if-you-want-to-continue-access-the-model-on-a-paid-basis-please-disable-the-use-free-tier-only-mode) The free tier of the model has been exhausted. If you want to continue access the model on a paid basis, please disable the "use free tier only" mode.

**Cause 1**: New user free quota is exhausted.
**Solution**: [Complete your profile](https://home.qwencloud.com/settings/account), then retry.
**Cause 2**: [Free quota only](/resources/free-quota) is enabled and the free quota is exhausted.
 The free quota displayed in the console is updated every minute. Refresh the page to see the latest data. 
**Solution** :

- 
To start paid calls, you can disable the [Free quota only](https://home.qwencloud.com/benefits) switch at any time without waiting until the free quota is used up.


- 
If you encounter this with the Coding Plan, it is likely a configuration issue. The Coding Plan requires its own Base URL and [API key](https://home.qwencloud.com/api-keys). See [Getting started](/coding-plan/overview).


## [​ ](#404-modelnotfound-model-not-found) 404-ModelNotFound/model_not_found

### [​ ](#the-provided-model-xxx-is-not-supported-by-the-batch-api) The provided model xxx is not supported by the Batch API.

**Cause:** The model does not support Batch calls, or the model name is misspelled.
**Solution:** See [Batch API](/developer-guides/text-generation/batch) for supported models and correct names.
### [​ ](#model-can-not-be-found-the-model-xxx-does-not-exist-the-model-xxx-does-not-exist-or-you-do-not-have-access-to-it) Model can not be found. / The model xxx does not exist. / The model xxx does not exist or you do not have access to it.

**Cause:** The model does not exist, or Qwen Cloud is not activated.
**Solution:**

- 
Verify the model name in [Models](https://www.qwencloud.com/models).


- 
Log in to [Qwen Cloud](https://www.qwencloud.com/).


## [​ ](#404-model-not-supported) 404-model_not_supported

### [​ ](#unsupported-model-xxx-for-openai-compatibility-mode) Unsupported model xxx for OpenAI compatibility mode.

**Cause:** The model does not support OpenAI-compatible mode.
**Solution:** Use the DashScope method instead.
## [​ ](#404-workspacenotfound) 404-WorkSpaceNotFound

### [​ ](#workspace-can-not-be-found) WorkSpace can not be found.

**Cause:** The workspace does not exist.
## [​ ](#404-notfound) 404-NotFound

### [​ ](#not-found) Not found!

**Cause:** The resource does not exist.
**Solution:** Verify the resource ID.
### [​ ](#request-path-not-found) Request path not found.

**Cause:** The service address does not exist when using the Fun-Music model.
**Solution:** Verify that the API endpoint path is correct and does not contain invalid characters.
## [​ ](#429-throttling) 429-Throttling

### [​ ](#requests-throttling-triggered) Requests throttling triggered.

**Cause:** The API call triggered rate limiting.
**Solution:** Reduce the call frequency or retry later.
### [​ ](#too-many-requests-in-route-try-again-later) Too many requests in route. Try again later.

**Cause** : Too many requests triggered rate limiting.
**Solution** : Retry later.
### [​ ](#all-models-are-temporarily-rate-limited-please-try-again-in-a-few-minutes) All models are temporarily rate-limited. Please try again in a few minutes.

**Cause:** When using services such as Coding Plan that call multiple models through a client (such as Claude Code), this message is returned by the client when all available models have triggered 429 rate limiting. The underlying server-side error codes are Throttling.RateQuota or Throttling.AllocationQuota.
**Solution:**

- 
Wait a few minutes and retry. The rate limit will be lifted automatically.


- 
Reduce request concurrency to avoid sending too many requests in a short period.


- 
To increase your rate limit, see [Rate limits](/developer-guides/administration/rate-limits) to request a quota increase.


## [​ ](#429-throttling-ratequota-limitrequests-limit-requests) 429-Throttling.RateQuota/LimitRequests/limit_requests

### [​ ](#you-have-exceeded-your-request-limit-requests-rate-limit-exceeded-try-again-later-you-exceeded-your-current-requests-list) You have exceeded your request limit./Requests rate limit exceeded, try again later. /You exceeded your current requests list.

**Cause:** The call frequency (RPS or RPM) triggered rate limiting.
**Solution:** See [Rate limits](/developer-guides/administration/rate-limits) and reduce the call frequency.
## [​ ](#429-throttling-burstrate-limit-burst-rate) 429-Throttling.BurstRate/limit_burst_rate

### [​ ](#request-rate-increased-too-quickly-to-ensure-system-stability-please-adjust-your-client-logic-to-scale-requests-more-smoothly-over-time) Request rate increased too quickly. To ensure system stability, please adjust your client logic to scale requests more smoothly over time.

**Cause** : The call frequency increased sharply, triggering the system stability protection.
**Solution** : Adopt a smooth request strategy such as uniform scheduling, exponential backoff, or request queue buffering. Distribute requests evenly within the time window to avoid instantaneous peaks.
## [​ ](#429-throttling-allocationquota-insufficient-quota) 429-Throttling.AllocationQuota/insufficient_quota

### [​ ](#allocated-quota-exceeded-please-increase-your-quota-limit-you-exceeded-your-current-quota-check-your-plan-and-billing-details) Allocated quota exceeded, please increase your quota limit./ You exceeded your current quota, check your plan and billing details.

**Cause:** The token consumption rate (TPS/TPM) triggered [Rate limits](/developer-guides/administration/rate-limits).
**Solution:** See the [Rate limits](/developer-guides/administration/rate-limits) documentation for model limits and adjust your call strategy. If default quotas are insufficient, request a temporary TPM increase in the console.
### [​ ](#too-many-requests-batch-requests-are-being-throttled-due-to-system-capacity-limits-try-again-later) Too many requests. Batch requests are being throttled due to system capacity limits. Try again later.

**Cause:** Too many batch requests triggered rate limiting.
**Solution:** Retry later.
### [​ ](#free-allocated-quota-exceeded-2) Free allocated quota exceeded.

**Cause:** The free quota has expired or been used up, and the model does not support pay-as-you-go.
**Solution:** Use another model.
## [​ ](#429-commoditynotpurchased) 429-CommodityNotPurchased

### [​ ](#commodity-has-not-purchased-yet) Commodity has not purchased yet.

**Cause:** The workspace is not subscribed.
**Solution:** Subscribe to the workspace service first.
## [​ ](#429-prepaidbilloverdue) 429-PrepaidBillOverdue

### [​ ](#the-prepaid-bill-is-overdue) The prepaid bill is overdue.

**Cause:** The workspace subscription bill has expired.
## [​ ](#429-postpaidbilloverdue) 429-PostpaidBillOverdue

### [​ ](#the-postpaid-bill-is-overdue) The postpaid bill is overdue.

**Cause:** The model inference product has expired.
## [​ ](#500-internalerror-internal-error) 500-InternalError/internal_error

### [​ ](#an-internal-error-has-occured-try-again-later-or-contact-service-support) An internal error has occured, try again later or contact service support.

**Cause:** An internal error occurred.
**Solution:**

- If using the [Qwen-Omni model](/developer-guides/speech/multimodal-speech), use streaming output.


### [​ ](#internal-server-error) Internal server error!

**Cause:** An internal algorithm error occurred.
**Solution:** Retry later.
### [​ ](#request-asr-failed) request asr failed

**Cause:** When using [CosyVoice voice cloning](/developer-guides/speech/tts), the audio file is invalid because it does not contain a valid human voice, the audio is unclear, or it contains excessive background noise.
**Solution:** Follow the recording guide and record the audio again, then retry.
## [​ ](#500-internalerror-fileupload) 500-InternalError.FileUpload

### [​ ](#oss-upload-error) oss upload error.

**Cause:** File upload failed.
**Solution:** Check your OSS configuration and network.
## [​ ](#500-internalerror-upload) 500-InternalError.Upload

### [​ ](#failed-to-upload-result) Failed to upload result.

**Cause:** Failed to upload the generation result.
**Solution:** Check the storage configuration or retry later.
## [​ ](#500-internalerror-algo) 500-InternalError.Algo

### [​ ](#inference-internal-error) inference internal error.

**Cause:** A service exception occurred.
**Solution:** Try again first to rule out a sporadic issue.
### [​ ](#expecting-delimiter-line-x-column-xxx-char-xxx) Expecting &#x27;,&#x27; delimiter: line x column xxx (char xxx)

**Cause:** The model generated invalid JSON, preventing a normal tool call.
**Solution:** Switch to the latest model or optimize the prompt and retry.
### [​ ](#missing-content-length-of-multimodal-url) Missing Content-Length of multimodal url.

**Cause:** The `Content-Length` field is missing from the response header of the URL request.
**Solution:** Try a different image link if the problem persists.
To view the `Content-Length` field:

- 
Open a browser such as Chrome or Firefox.


- 
Press F12 or right-click and select "Inspect" to open developer tools.


- 
Switch to the Network tab.


- 
Enter the image URL in the address bar.


- 
Find the request, then check the `Content-Length` in the "Response Headers" section.


### [​ ](#an-error-occurred-in-model-serving-error-message-is-request-rejected-by-inference-engine) An error occurred in model serving, error message is: [Request rejected by inference engine!]

**Cause:** An error occurred on the model service server.
**Solution:** Retry later.
### [​ ](#an-internal-error-has-occured-during-algorithm-execution) An internal error has occured during algorithm execution.

**Cause:** An error occurred during algorithm runtime.
**Solution:** Retry later.
### [​ ](#inference-error-inference-error) Inference error: Inference error.

**Cause:** An inference error occurred.
**Solution:** Check the input image for corruption. Verify the person image contains a complete, clear face.
### [​ ](#role-must-be-in-user-assistant) Role must be in [user, assistant]

**Cause:** When using Qwen-MT, the messages array contains a role other than `user`.
**Solution:** The messages array must contain only one element with role `user`.
### [​ ](#embedding-pipeline-error-xxx) Embedding_pipeline_Error: xxx

**Cause:** An error occurred in image or video pre-processing.
**Solution:** Verify the uploaded image or video and request code meet requirements, then retry.
### [​ ](#receive-batching-backend-response-failed) Receive batching backend response failed!

**Cause:** An internal service error occurred.
**Solution:** Retry later.
### [​ ](#music-receive-batching-backend-response-failed) [music]Receive batching backend response failed!

**Cause:** The Fun-Music model service exceeded the concurrency limit.
**Solution:** Reduce the number of concurrent requests and try again.
### [​ ](#other-kinds-of-server-error) Other kinds of server error.

**Cause:** An unknown internal error occurred when using the Fun-Music model.
**Solution:** Provide the request ID to technical support for troubleshooting.
### [​ ](#an-internal-error-has-occured-during-execution-try-again-later-or-contact-service-support-algorithm-process-error-inference-error-an-internal-error-occurs-during-computation-please-try-this-model-later) An internal error has occured during execution, try again later or contact service support. / algorithm process error. / inference error. / An internal error occurs during computation, please try this model later.

**Cause:** An internal algorithm error occurred.
**Solution:** Retry later.
### [​ ](#list-index-out-of-range) list index out of range

**Cause:** The last element in the messages array must be a User Message.
**Solution:** Adjust the `messages` array so the last element is `{"role": "user", ...}`.
## [​ ](#500-internalerror-timeout) 500-InternalError.Timeout

### [​ ](#an-internal-timeout-error-has-occured-during-execution-try-again-later-or-contact-service-support) An internal timeout error has occured during execution, try again later or contact service support.

**Cause:** An asynchronous task returned no result within 3 hours.
**Solution:** Check the task execution status or contact technical support.
## [​ ](#500-systemerror) 500-SystemError

### [​ ](#an-system-error-has-occured-try-again-later) An system error has occured, try again later.

**Cause:** A system error occurred.
**Solution:** Retry later.
## [​ ](#500-modelservicefailed) 500-ModelServiceFailed

### [​ ](#failed-to-request-model-service) Failed to request model service.

**Cause:** The model service call failed.
**Solution:** Retry later.
## [​ ](#500-requesttimeout) 500-RequestTimeOut

### [​ ](#request-timed-out-try-again-later-response-timeout-i-o-error-on-post-request-for-https-dashscope-intl-aliyuncs-com-compatible-mode-v1-chat-completions-timeout) Request timed out, try again later. / Response timeout! / I/O error on POST request for "[https://dashscope-intl.aliyuncs.com/compatible-mode/v1/chat/completions](https://dashscope-intl.aliyuncs.com/compatible-mode/v1/chat/completions)": timeout

**Cause:**

- 
Language model: The request timed out (timeout is 300 seconds).


- 
Speech recognition (Paraformer): No audio or only silent audio was sent for a long time.


- 
Image generation/editing model: Processing time exceeded the limit due to large image dimensions or high complexity.


**Solution:**

- 
Language model: Use streaming output. See [Streaming output](/developer-guides/text-generation/streaming).


- 
Paraformer: Set `heartbeat` to `true`, or end the recognition task promptly.


- 
Image model: Reduce image resolution, simplify editing requirements, or retry later.


## [​ ](#500-responsetimeout) 500-ResponseTimeout

### [​ ](#response-stream-timeout) Response stream timeout

**Cause:** An internal execution timeout occurred when using the Fun-Music model.
**Solution:** Try again later.
## [​ ](#500-invokepluginfailed) 500-InvokePluginFailed

### [​ ](#failed-to-invoke-plugin) Failed to invoke plugin.

**Cause:** The plug-in call failed.
**Solution:** Check the plug-in configuration and availability.
## [​ ](#500-appprocessfailed) 500-AppProcessFailed

### [​ ](#failed-to-proceed-application-request) Failed to proceed application request.

**Cause:** The application flow failed.
**Solution:** Check the application configuration and flow nodes.
## [​ ](#500-rewritefailed) 500-RewriteFailed

### [​ ](#failed-to-rewrite-content-for-prompt) Failed to rewrite content for prompt.

**Cause:** The prompt rewriting model call failed.
**Solution:** Retry later.
## [​ ](#500-retrievalfailed) 500-RetrievalFailed

### [​ ](#failed-to-retrieve-data-from-documents) Failed to retrieve data from documents.

**Cause:** Document retrieval failed.
**Solution:** Check the document index and retrieval configuration.
## [​ ](#500-503-modelservingerror) 500/503-ModelServingError

### [​ ](#too-many-requests-your-requests-are-being-throttled-due-to-system-capacity-limits-try-again-later) Too many requests. Your requests are being throttled due to system capacity limits. Try again later.

**Cause:** Network resources are saturated.
**Solution:** Retry later.
## [​ ](#503-modelunavailable) 503-ModelUnavailable

### [​ ](#model-is-unavailable-try-again-later) Model is unavailable, try again later.

**Cause:** The model is temporarily unavailable.
**Solution:** Retry later.
## [​ ](#sdk-errors) SDK errors

### [​ ](#error-authenticationerror-no-api-key-provided-you-can-set-by-dashscope-api-key-your-api-key-in-code-or-you-can-set-it-via-environment-variable-dashscope-api-key-your-api-key) error.AuthenticationError: No api key provided. You can set by dashscope.api_key = your_api_key in code, or you can set it via environment variable DASHSCOPE_API_KEY= your_api_key.

**Cause:** No API key was provided when using the DashScope SDK.
**Solution:** See [Configure your API key](/api-reference/preparation/export-api-key-env).
### [​ ](#openai-openaierror-the-api-key-client-option-must-be-set-either-by-passing-api-key-to-the-client-or-by-setting-the-openai-api-key-environment-variable) openai.OpenAIError: The api_key client option must be set either by passing api_key to the client or by setting the OPENAI_API_KEY environment variable

**Cause:** No API key was provided.
**Solution:**

- **Pass the API key using an environment variable (recommended)**


Set `DASHSCOPE_API_KEY` (see [Configure your API key](/api-reference/preparation/export-api-key-env)). When initializing the `client`, read the key with `os.getenv`:
`client = OpenAI(api_key=os.getenv("DASHSCOPE_API_KEY"),...)`

- **Pass the API key in plaintext (testing only)**


Pass the API key directly to `api_key`:
`client = OpenAI(api_key="sk-xxx", ...)`
**Note:** This poses a security risk. Do not use in production.
### [​ ](#bad-request-for-url-xxx) Bad Request for url: xxx

**Cause:** Adding `response.raise_for_status()` with the Python requests library prevents the server&#x27;s error details from being returned.
**Solution:** Use `print(response.json())` to view the server response.
### [​ ](#cannot-resolve-symbol-ttsv2) Cannot resolve symbol &#x27;ttsv2&#x27;

**Cause:** The DashScope SDK is outdated for [speech synthesis (CosyVoice)](/developer-guides/speech/tts).
**Solution:** [Install the latest DashScope SDK](/api-reference/preparation/install-sdk).
## [​ ](#networkerror) NetworkError

### [​ ](#noapikeyexception-can-not-find-api-key) NoApiKeyException: Can not find api-key.

**Cause:** The environment variable has not taken effect.
**Solution:** Restart the client or IDE and retry. See [FAQ](/resources/faq-inference).
### [​ ](#connectexception-failed-to-connect-to-dashscope-intl-aliyuncs-com) ConnectException: Failed to connect to dashscope-intl.aliyuncs.com

**Cause:** The local network has an issue.
**Solution:** Check local network settings. You may have a certificate issue preventing HTTPS access, or a misconfigured firewall. Try a different network environment or server.
### [​ ](#inputrequiredexception-parameter-invalid-text-is-null) InputRequiredException: Parameter invalid: text is null

**Cause** : When using [speech synthesis (CosyVoice)](/developer-guides/speech/tts), the text to be synthesized was not sent.
**Solution:** Assign a value to the `text` parameter.
### [​ ](#multimodalconversation-call-missing-1-required-positional-argument-messages) MultiModalConversation.call() missing 1 required positional argument: &#x27;messages&#x27;

**Cause** : The DashScope SDK version is outdated.
**Solution:** [Install the latest DashScope SDK](/api-reference/preparation/install-sdk).
## [​ ](#mismatched-model) mismatched_model

### [​ ](#the-model-xxx-for-this-request-does-not-match-the-rest-of-the-batch-each-batch-must-contain-requests-for-a-single-model) The model &#x27;xxx&#x27; for this request does not match the rest of the batch. Each batch must contain requests for a single model.

**Cause:** In a single batch task, all requests must use the same model.
**Solution:** Check your input file according to the [input file format](/developer-guides/text-generation/batch#input-file-format).
## [​ ](#duplicate-custom-id) duplicate_custom_id

### [​ ](#the-custom-id-xxx-for-this-request-is-a-duplicate-of-another-request-the-custom-id-parameter-must-be-unique-for-each-request-in-a-batch) The custom_id &#x27;xxx&#x27; for this request is a duplicate of another request. The custom_id parameter must be unique for each request in a batch.

**Cause:** In a single batch task, the ID of each request must be unique.
**Solution:** Verify the input file per the [input file format](/developer-guides/text-generation/batch#input-file-format) and ensure all request IDs are unique.
### [​ ](#upload-file-capacity-exceed-limit-upload-file-number-exceed-limit) Upload file capacity exceed limit. / Upload file number exceed limit.

**Cause:** The file upload failed. The storage space under the current account is full or nearly full.
**Solution:** Use the [OpenAI compatible - File](/api-reference/platform-api/file) API to delete unnecessary files. Storage supports up to 10,000 files and 100 GB total.
## [​ ](#websocket-errors) WebSocket errors

### [​ ](#the-decoded-text-message-was-too-big-for-the-output-buffer-and-the-endpoint-does-not-support-partial-messages) The decoded text message was too big for the output buffer and the endpoint does not support partial messages

**Cause:** When using streaming speech recognition with Paraformer, the recognition result data is too large.
**Solution:** Send audio in segments. Each segment should be about 100 milliseconds, with data between 1 KB and 16 KB.
### [​ ](#timeouterror-websocket-connection-could-not-established-within-5s-check-your-network-connection-firewall-settings-or-server-status) TimeoutError: websocket connection could not established within 5s. Check your network connection, firewall settings, or server status.

**Cause:** When using CosyVoice speech synthesis, a WebSocket connection could not be established within 5 seconds.
**Solution:** Check local network and firewall settings, or try a different network environment.
### [​ ](#unsupported-audio-format-xxx) unsupported audio format:xxx

**Cause:** The CosyVoice voice cloning audio format does not meet requirements.
**Solution:** Use WAV (16-bit), MP3, or M4A. The file extension alone does not guarantee the format. Use a tool like ffprobe or mediainfo to confirm the encoding format.
### [​ ](#internal-unknown-error) internal unknown error

**Cause:** The CosyVoice voice cloning audio format may not meet requirements.
**Solution:** Use WAV (16-bit), MP3, or M4A. Use a tool like ffprobe or mediainfo to confirm the encoding format.
### [​ ](#invalid-backend-response-received-missing-status-name) Invalid backend response received (missing status name)

**Cause:** When using the RESTful API for audio file recognition with Paraformer, a request parameter is misspelled.
**Solution:** Check your code against the API reference.
### [​ ](#no-input-audio-error) NO_INPUT_AUDIO_ERROR

**Cause:** No valid speech was detected.
**Solution:** For real-time speech recognition with Paraformer:

- 
Check whether there is audio input.


- 
Check whether the audio format is correct. Supported formats: pcm, wav, mp3, opus, speex, aac, amr.


### [​ ](#success-with-no-valid-fragment) SUCCESS_WITH_NO_VALID_FRAGMENT

**Cause:** The Paraformer recognition query succeeded, but VAD detected no valid speech.
**Solution:** Verify the audio contains valid speech. Pure silence produces no recognition result.
### [​ ](#asr-response-have-no-words) ASR_RESPONSE_HAVE_NO_WORDS

**Cause:** The Paraformer recognition query succeeded, but the final result is empty.
**Solution:** Check whether the audio contains valid speech, or whether all speech is filler words filtered out by `disfluency_removal_enabled`.
### [​ ](#file-download-failed) FILE_DOWNLOAD_FAILED

**Cause:** The Paraformer audio file failed to download.
**Solution:** Verify the audio file path is correct and externally accessible.
### [​ ](#file-check-failed) FILE_CHECK_FAILED

**Cause:** The Paraformer audio file format is incorrect.
**Solution:** Use single-track/dual-track WAV or MP3 format.
### [​ ](#file-too-large) FILE_TOO_LARGE

**Cause:** The Paraformer audio file is too large.
**Solution:** If the file exceeds 2 GB, segment it into smaller files.
### [​ ](#file-normalize-failed) FILE_NORMALIZE_FAILED

**Cause:** The Paraformer audio file failed to normalize.
**Solution:** Check whether the audio file is corrupted and can be played.
### [​ ](#file-parse-failed) FILE_PARSE_FAILED

**Cause:** The Paraformer audio file failed to parse.
**Solution:** Check whether the audio file is corrupted and can be played.
### [​ ](#mkv-parse-failed) MKV_PARSE_FAILED

**Cause:** The Paraformer MKV file failed to parse.
**Solution:** Check whether the audio file is corrupted and can be played.
### [​ ](#file-trans-task-expired) FILE_TRANS_TASK_EXPIRED

**Cause:** The Paraformer audio file recognition task has expired.
**Solution:** The TaskId does not exist or has expired. Resubmit the task.
### [​ ](#request-invalid-file-url-value) REQUEST_INVALID_FILE_URL_VALUE

**Cause:** The Paraformer `file_url` parameter is invalid.
**Solution:** Verify the `file_url` parameter format.
### [​ ](#content-length-check-failed) CONTENT_LENGTH_CHECK_FAILED

**Cause:** The Paraformer `content-length` check failed.
**Solution:** Verify `content-length` in the HTTP response matches the actual file size.
### [​ ](#file-404-not-found) FILE_404_NOT_FOUND

**Cause:** The Paraformer file does not exist.
**Solution:** Verify the file URL.
### [​ ](#file-403-forbidden) FILE_403_FORBIDDEN

**Cause:** You lack permission to download the Paraformer audio file.
**Solution:** Check the file access permissions.
### [​ ](#file-server-error) FILE_SERVER_ERROR

**Cause:** The file server is unavailable.
**Solution:** Retry later or check the file server status.
### [​ ](#audio-duration-too-long) AUDIO_DURATION_TOO_LONG

**Cause:** The Paraformer audio file exceeds 12 hours.
**Solution:** Segment the audio and submit multiple recognition tasks. Use tools like FFmpeg to split the audio.
### [​ ](#decode-error) DECODE_ERROR

**Cause:** Paraformer failed to detect the audio file information.
**Solution:** Verify the file is in a supported audio format.
### [​ ](#client-error-qwen-tts-engine-return-error-code-411) CLIENT_ERROR -[qwen-tts:]Engine return error code: 411

**Cause:** When using Qwen-TTS real-time speech synthesis with `qwen-tts-vc-realtime-2025-08-20`, the default voice was used. This model supports only cloned voices.
**Solution:** Use a voice generated by voice cloning, not the default voice.
### [​ ](#no-valid-audio-error) NO_VALID_AUDIO_ERROR

**Cause:** When using Paraformer, the audio is invalid.
**Solution:** Verify the audio format, sample rate, and other parameters meet requirements.
### [​ ](#invalidparameter-task-can-not-be-null) InvalidParameter: task can not be null

**Cause:** The `input` field is missing from the payload of the run-task or finish-task instruction, or the `input.text` field is missing from the continue-task instruction in the CosyVoice WebSocket API.
**Solution:**

- 
run-task instruction: Verify the payload contains `"input": {}` (an empty object). The `input` field is required.


- 
continue-task instruction: Verify `payload.input` contains the `text` field with a non-empty value.


- 
finish-task instruction: Verify the payload contains `"input": {}`.


## [​ ](#200-bailiangateway-workspace-notauthorised) 200-BailianGateway.Workspace.NotAuthorised

**Cause:** This error may occur due to: (1) The accessed URL contains special characters or uses a non-standard format, which causes workspace authorization validation to fail. (2) A sub-account user attempted to operate on a workspace without the required permissions.
**Solution:** (1) Return to the Qwen Cloud homepage and navigate to the target page. (2) The account owner or a user with administrator permissions must grant the required workspace permissions to the sub-account user. [Previous ](/api-reference/preparation/install-sdk)[OpenAI chat Compatible Chat API Next ](/api-reference/chat/openai-chat)
