# Convert LaTeX formulas to speech (Chinese language only)

> **Source:** https://docs.qwencloud.com/developer-guides/speech/latex-to-speech

Speak math formulas aloud

 Copy page CosyVoice converts mathematical formulas embedded in text to natural speech for audiobooks, online education, and other audio content in mathematics and physics.
 This feature only supports **Chinese**. Formulas may not be pronounced correctly in other languages. Pronunciation examples in this guide are English translations and do not represent actual synthesis performance. 
## [​ ](#steps) Steps

Wrap the formulas in your text with specific separators, and then call the speech synthesis API.
 1 Mark formulas with separators

Wrap formulas with any of these separators (all produce identical results):
- `$...$`

- `$$...$$`

- `\(...\)`

- `\[...\]`


Example:Copy ```\nHere is the quadratic formula: $x = \frac{-b \pm \sqrt{b^2-4ac}}{2a}$. Calculate carefully.

``` 2 Call the API to request speech synthesis

Call the speech synthesis API with your marked formulas. In JSON or string literals, the backslash (`\`) is an escape character — write it as `\\`.Example call (in Python):Copy ```\n# coding=utf-8

import os
import dashscope
from dashscope.audio.tts_v2 import *

# Set your API key (format: sk-xxx). If not set as env var, uncomment the next line:
# dashscope.api_key = "sk-xxx"
dashscope.api_key = os.environ.get(&#x27;DASHSCOPE_API_KEY&#x27;)

dashscope.base_websocket_api_url=&#x27;wss://dashscope-intl.aliyuncs.com/api-ws/v1/inference&#x27;

# Model
model = "cosyvoice-v3-flash"
# Voice (must be compatible with your model version)
# cosyvoice-v3-flash/plus: longanyang, etc.
voice = "longanyang"

# Instantiate synthesizer with model and voice
synthesizer = SpeechSynthesizer(model=model, voice=voice)
# Send text and get binary audio
audio = synthesizer.call("This is the quadratic formula: $x = \\frac{-b \\pm \\sqrt{b^2-4ac}}{2a}$. Calculate it carefully.")
# First request includes WebSocket connection setup time
print(&#x27;[Metric] requestId: {}, first-package delay: {} ms&#x27;.format(
 synthesizer.get_last_request_id(),
 synthesizer.get_first_package_delay()))

# Save audio to a local file
with open(&#x27;output.mp3&#x27;, &#x27;wb&#x27;) as f:
 f.write(audio)

``` 
## [​ ](#supported-tags-and-symbols) Supported tags and symbols

The currently supported tags and symbols are listed below.
### [​ ](#basic-mathematics) Basic mathematics

Tag or symbolFunctionFormula content exampleFormula input examplePronunciation+Add2 + 3 = 5`$2 + 3 = 5$`Two plus three equals five-Subtract3 - 2 = 1`$3 - 2 = 1$`Three minus two equals one\pmPlus or minus / Positive or negative\pm 1 \pm 2`$\pm 1\pm 2$`Plus or minus one, plus or minus two\timesMultiply2 \times 3 = 6`$2 \times 3 = 6$`Two times three equals six×Multiply2 × 3 = 6`$$2 × 3 = 6$$`Two times three equals six*Multiply2 * 3 = 6`\(2 * 3 = 6\)`Two times three equals six\divDivide6\div2=3`\[6\div2=3\]`Six divided by two equals three÷Divide6÷2=3`$6÷2=3$`Six divided by two equals three/Divide6/2=3`$6/2=3$`Six divided by two equals three=Equals3+5=8`$3+5=8$`Three plus five equals eight<Less than1< 2`$1< 2$`One is less than two≤Less than or equal to3≤5`$3≤5$`Three is less than or equal to five<=Less than or equal to3<=5`$3<=5$`Three is less than or equal to five\leqLess than or equal to3\leq5`$3\leq 5$`Three is less than or equal to five\leLess than or equal to3\le5`$3\le 5$`Three is less than or equal to five\leqqLess than or equal to3\leqq5`$3\leqq 5$`Three is less than or equal to five\leqslantLess than or equal to3\leqslant5`$3\leqslant 5$`Three is less than or equal to five>Greater than2>1`$2>1$`Two is greater than one≥Greater than or equal to5≥3`$5≥3$`Five is greater than or equal to three>=Greater than or equal to5>=3`$5>=3$`Five is greater than or equal to three\geqGreater than or equal to5\geq3`$5\geq 3$`Five is greater than or equal to three\geGreater than or equal to5\ge3`$5\ge 3$`Five is greater than or equal to three\geqqGreater than or equal to5\geqq3`$5\geqq 3$`Five is greater than or equal to three\geqslantGreater than or equal to5\geqslant3`$5\geqslant 3$`Five is greater than or equal to three\fracFraction`2\frac3``$\frac {2}{3}$`Two-thirds^Power`2^1``$2^{1}$`Two to the first power\sqrtRoot`\sqrt{9} = 3``$\sqrt {9} = 3$`The square root of nine equals three\sqrtRoot`\sqrt[3]{8} = 2``$\sqrt[3]{8} = 2$`The cube root of eight equals two%Percentage`5\%``$5\%$`Five percent|Absolute value`∣3∣=3``$|3| =3$`The absolute value of three equals three\vertAbsolute value`3\vert=3``$\vert 3\vert =3$`The absolute value of three equals three\lgLogarithm`lg {10}``$\lg {10}$`log ten\logLogarithm`\log{5}``$\log{5}$`log five\lnNatural logarithm`\lnX``$ln {10}$`LN 10!Factorial5!`$5!$`Five factorial()Parentheses(2+1)`$(2+1)$`Two plus one in parentheses`\{ \}`Braces`\{2+1\}``$\{2+1\}$`Open brace two plus one close brace 
### [​ ](#special-mathematical-symbols) Special mathematical symbols

Tag or symbolTransformFormula content exampleFormula input examplePronunciation\alphaalpha\alpha`$\alpha$`alpha\Alphaalpha\Alpha`$\Alpha$`alpha\betabeta\beta`$\beta$`beta\Betabeta\Beta`$\Beta$`beta\gammagamma\gamma`$\gamma$`gamma\Gammagamma\Gamma`$\Gamma$`gamma\deltadelta\delta`$\delta$`delta\Deltadelta\Delta`$\Delta$`delta\inftyinfinity (Chinese)\infty`$\infty$`Infinity∞infty (English)∞`$∞$`Infinity 
### [​ ](#geometry) Geometry

Tag or symbolFunctionFormula content exampleFormula input examplePronunciation\piPi\pi=3.14159`$\pi =3.14159$`Pi equals 3.14159\sinTrigonometric function`\sin 30^\circ=\frac{1}{2}``$\sin 30^\circ =\frac {1}{2}$`The sine of 30 degrees equals 1/2\cosTrigonometric function`\cos 30^\circ=\frac{\sqrt{2}}{2}``$\cos 30^\circ =\frac {\sqrt {2}}{2}$`The cosine of 30 degrees equals the square root of two over two\tanTrigonometric function`\tan 30^\circ=\frac{\sin 30^\circ}{\cos 30^\circ}``$\tan 30^\circ =\frac {\sin 30^\circ}{\cos 30^\circ}$`The tangent of 30 degrees equals the sine of 30 degrees over the cosine of 30 degrees\cscTrigonometric function\csc A`$\csc A$`cosecant A\secTrigonometric function\sec A`$\sec A$`secant A\cotTrigonometric function\cot A`$\cot A$`cotangent A\angleAngle\angle AB`$\angle AB$`Angle AB∠Angle∠AB`$∠AB$`Angle AB^\circDegree∠AB = 30^\circ`$∠AB = 30^\circ$`Angle AB equals 30 degrees\odotCircle\odot`$\odot$`Circle`\overset\frown`Arc`\overset\frown {BC}``$\overset\frown {BC}$`Arc BC`\rm{Rt}`Right angle`\because \rm{Rt}\triangle ABC``$\because \rm{Rt}\triangle ABC$`Because triangle ABC is a right triangle`\mathrm{Rt}`Right angle`\therefore AB \perp BC``$\therefore AB \perp BC$`Therefore, AB is perpendicular to BC\triangleTriangle\triangle ABC`$\triangle ABC$`Triangle ABC△Triangle△ABC`$△ABC$`Triangle ABC\parallelogramParallelogram\parallelogram ABCD`$\parallelogram ABCD$`Parallelogram ABCD\perpPerpendicularAB \perp BC`$AB \perp BC$`AB is perpendicular to BC\botPerpendicularAB \bot BC`$AB \bot BC$`AB is perpendicular to BC⊥PerpendicularAB ⊥ BC`$AB ⊥ BC$`AB is perpendicular to BC\parallelParallelA\parallel B`$A\parallel B$`A is parallel to B\equalparallelParallel and equal toA\equalparallel B`$A\equalparallel B$`A is parallel and equal to B\congCongruent△ABC\cong△DEF`$△ABC\cong△DEF$`Triangle ABC is congruent to triangle DEF 
### [​ ](#conditions) Conditions

Tag or symbolFunctionFormula content exampleFormula input examplePronunciation\impliesImplies\implies 1+1=2`$\implies 1+1=2$`This implies that one plus one equals two\iffEquivalent top\iffq`$p\iffq$`p is equivalent to q\becauseBecause\because a = b \therefore b=a`$\because a = b \therefore b=a$`Because a equals b, b equals a\thereforeTherefore\because a = b \therefore b=a`$\because a = b \therefore b=a$`Because a equals b, b equals a 
### [​ ](#units) Units

Units must be wrapped with `\unit`, `\quantity`, `\mathit`, `\mathrm`, or `\rm` tags (such as `\unit{cm}`).
Tag or symbolPronunciationFormula content exampleFormula input exampleExample pronunciationmmmillimeter`5\quantity{mm}``$5\quantity{mm}$`5 millimeterscmcentimeter`5\quantity{cm}``$5\quantity{cm}$`5 centimetersdmdecimeter`5\quantity{dm}``$5\quantity{dm}$`5 decimetersmmeter`5\quantity{m}``$5\quantity{m}$`5 meterskmkilometer`5\quantity{km}``$5\quantity{km}$`5 kilometersggram`5\quantity{g}``$5\quantity{g}$`5 gramskgkilogram`5\quantity{kg}``$5\quantity{kg}$`5 kilogramstton`5\quantity{t}``$5\quantity{t}$`5 tonsmm^2square millimeter`5\quantity{mm^2}``$5\quantity{mm^2}$`5 square millimeterscm^2square centimeter`5\quantity{cm^2}``$5\quantity{cm^2}$`5 square centimetersdm^2square decimeter`5\quantity{dm^2}``$5\quantity{dm^2}$`5 square decimetersm^2square meter`5\quantity{m^2}``$5\quantity{m^2}$`5 square meterskm^2square kilometer`5\quantity{km^2}``$5\quantity{km^2}$`5 square kilometersmm^3cubic millimeter`5\quantity{mm^3}``$5\quantity{mm^3}$`5 cubic millimeterscm^3cubic centimeter`5\quantity{cm^3}``$5\quantity{cm^3}$`5 cubic centimetersdm^3cubic decimeter`5\quantity{dm^3}``$5\quantity{dm^3}$`5 cubic decimetersm^3cubic meter`5\quantity{m^3}``$5\quantity{m^3}$`5 cubic meterskm^3cubic kilometer`5\quantity{km^3}``$5\quantity{km^3}$`5 cubic kilometersmlmilliliter`5\quantity{ml}``$5\quantity{ml}$`5 millilitersssecond`5\quantity{s}``$5\quantity{s}$`5 secondsminminute`5\quantity{min}``$5\quantity{min}$`5 minuteshhour`5\quantity{h}``$5\quantity{h}$`5 hourskm/hkilometers per hour`5\quantity{km/h}``$5\quantity{km/h}$`5 kilometers per hourg/lgrams per liter`5\quantity{g/l}``$5\quantity{g/l}$`5 grams per liter 
## [​ ](#limitations) Limitations


- **Chinese only:** Other languages are not supported

- **Content limits:**

Use only the tags and symbols in [Supported tags and symbols](#supported-tags-and-symbols)

- Markdown math blocks (````math ... ````) are not supported

- Include only formulas within separators - other content may cause inaccurate synthesis


- **Compatible models:** `cosyvoice-v3-flash`, `cosyvoice-v3-plus`


## [​ ](#faq) FAQ

### [​ ](#why-is-the-formula-i-entered-not-read) Why is the formula I entered not read?


- **Separators:** Confirm the formula is wrapped in `$...$`, `$$...$$`, `\(...\)`, or `\[...\]`

- **Formula complexity:** Ensure the formula uses only [supported tags and symbols](#supported-tags-and-symbols)

- **Escape characters:** Confirm backslashes (`\`) are escaped as `\\` in API requests


### [​ ](#how-do-i-handle-backslash-in-my-code) How do I handle backslash (`\`) in my code?

The backslash (`\`) is an escape character in string literals and JSON. Escape it as `\\`. Example: write `\frac` as `\\frac` in Python, Java, JavaScript, and similar languages. [Previous ](/developer-guides/speech/voice-design)[Speech-to-speech models Choose a model for voice conversation, speech translation, or simultaneous interpretation. Next ](/developer-guides/speech/s2s-models)
