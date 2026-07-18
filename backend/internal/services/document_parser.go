package services

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/ledongthuc/pdf"
)

// parsedDocument keeps the import text and the optional publication title
// together. Chapter splitting must retain the source headings, while a title
// should be taken from document metadata or explicit title markup rather than
// from the uploaded filename.
type parsedDocument struct {
	text  string
	title string
}

// fileTypeOf returns the lowercased file extension without the leading dot,
// e.g. "manuscript.PDF" -> "pdf". Empty for filenames with no extension.
func fileTypeOf(filename string) string {
	ext := filepath.Ext(filename)
	return strings.ToLower(strings.TrimPrefix(ext, "."))
}

// parseDocument dispatches on file extension and returns plain text suitable
// for chapter splitting. Raw binary content must never reach the caller —
// unsupported formats and parse failures return an error instead.
func parseDocument(filename string, content []byte) (string, error) {
	parsed, err := parseDocumentDetails(filename, content)
	return parsed.text, err
}

func parseDocumentDetails(filename string, content []byte) (parsedDocument, error) {
	switch fileTypeOf(filename) {
	case "md":
		text := stripBOM(string(content))
		return parsedDocument{text: text, title: markdownDocumentTitle(text)}, nil
	case "txt":
		text := stripBOM(string(content))
		return parsedDocument{text: text, title: inferredDocumentTitle(text)}, nil
	case "docx":
		return parseDOCXDetails(content)
	case "pdf":
		text, err := parsePDF(content)
		if err != nil {
			return parsedDocument{}, err
		}
		return parsedDocument{text: text, title: pdfDocumentTitle(content, text)}, nil
	default:
		return parsedDocument{}, fmt.Errorf("unsupported file type %q (only .md, .txt, .docx, and .pdf are supported)", fileTypeOf(filename))
	}
}

func markdownDocumentTitle(text string) string {
	return markdownDocumentTitleWithH1Priority(text, true)
}

// markdownDocumentTitleWithH1Priority distinguishes author-written Markdown
// from DOCX headings that were normalized into Markdown markers on extraction.
func markdownDocumentTitleWithH1Priority(text string, explicitH1Priority bool) string {
	trimmed := strings.TrimSpace(text)
	if strings.HasPrefix(trimmed, "---\n") || strings.HasPrefix(trimmed, "---\r\n") {
		lines := strings.Split(trimmed, "\n")
		for index := 1; index < len(lines); index++ {
			line := strings.TrimSpace(lines[index])
			if line == "---" {
				break
			}
			if key, value, found := strings.Cut(line, ":"); found && strings.EqualFold(strings.TrimSpace(key), "title") {
				return cleanDocumentTitle(value)
			}
		}
	}
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			// A leading H1 is useful title evidence only when it is not the
			// first chapter. Imports commonly start directly at "# Chapter
			// One"; using that as the work title would replace the filename
			// with a chapter label.
			title := cleanDocumentTitle(strings.TrimPrefix(line, "# "))
			if !looksLikeChapterHeading(title) && (explicitH1Priority || !hasRepeatedTitleCaseHeadings(strings.Split(text, "\n"))) {
				return title
			}
			return ""
		}
	}
	return ""
}

func inferredDocumentTitle(text string) string {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) < 2 {
		return ""
	}
	first := cleanDocumentTitle(lines[0])
	if first == "" || looksLikeChapterHeading(first) || hasRepeatedTitleCaseHeadings(lines) {
		return ""
	}
	// A standalone first line is a cautious fallback for plain-text/PDF books.
	// Prose on the same line is not treated as a title.
	if strings.TrimSpace(lines[1]) == "" {
		return first
	}
	return ""
}

func cleanDocumentTitle(value string) string {
	value = strings.Trim(strings.TrimSpace(value), "\"'")
	if len([]rune(value)) == 0 || len([]rune(value)) > 160 {
		return ""
	}
	return value
}

func looksLikeChapterHeading(value string) bool {
	normalized := strings.ToLower(strings.TrimSpace(value))
	return strings.HasPrefix(normalized, "chapter ") || strings.HasPrefix(normalized, "capítulo ") || strings.HasPrefix(normalized, "capitulo ")
}

// pdfDocumentTitle prefers the standard Info dictionary title. PDF metadata
// is not guaranteed, so it falls back only to a credible standalone title
// line. In particular, a character-name chapter heading must leave the work
// at its filename-derived title rather than becoming the book title.
func pdfDocumentTitle(content []byte, text string) string {
	if title := pdfInfoTitle(content); title != "" {
		return title
	}
	return inferredPDFDocumentTitle(text)
}

// pdfInfoTitle resolves the /Info indirect reference from the final trailer.
// A PDF may contain /Title entries for outlines or bookmarks before its Info
// dictionary; they must never be used as the document title.
func pdfInfoTitle(content []byte) string {
	objectNumber, generation, found := pdfTrailerInfoReference(content)
	if !found {
		return ""
	}
	return pdfLiteralTitle(pdfIndirectObject(content, objectNumber, generation))
}

func pdfTrailerInfoReference(content []byte) (objectNumber, generation int, found bool) {
	trailer := bytes.LastIndex(content, []byte("trailer"))
	if trailer < 0 {
		return 0, 0, false
	}
	trailerEnd := bytes.Index(content[trailer:], []byte("startxref"))
	if trailerEnd < 0 {
		return 0, 0, false
	}
	info := bytes.Index(content[trailer:trailer+trailerEnd], []byte("/Info"))
	if info < 0 {
		return 0, 0, false
	}
	index := trailer + info + len("/Info")
	index = skipPDFWhitespace(content, index)
	objectNumber, index, found = parsePDFDecimal(content, index)
	if !found {
		return 0, 0, false
	}
	index = skipPDFWhitespace(content, index)
	generation, index, found = parsePDFDecimal(content, index)
	if !found {
		return 0, 0, false
	}
	index = skipPDFWhitespace(content, index)
	return objectNumber, generation, index < len(content) && content[index] == 'R'
}

func skipPDFWhitespace(content []byte, index int) int {
	for index < len(content) && (content[index] == ' ' || content[index] == '\t' || content[index] == '\r' || content[index] == '\n' || content[index] == '\f' || content[index] == 0) {
		index++
	}
	return index
}

func parsePDFDecimal(content []byte, index int) (value, next int, found bool) {
	start := index
	for index < len(content) && content[index] >= '0' && content[index] <= '9' {
		value = value*10 + int(content[index]-'0')
		index++
	}
	return value, index, index > start
}

func pdfIndirectObject(content []byte, objectNumber, generation int) []byte {
	marker := []byte(fmt.Sprintf("%d %d obj", objectNumber, generation))
	start := bytes.Index(content, marker)
	if start < 0 {
		return nil
	}
	start += len(marker)
	end := bytes.Index(content[start:], []byte("endobj"))
	if end < 0 {
		return nil
	}
	return content[start : start+end]
}

// pdfLiteralTitle decodes the top-level /Title key in a single Info
// dictionary. Literal strings and comments are skipped while scanning so text
// such as a Subject containing "/Title (...)" cannot masquerade as a key.
func pdfLiteralTitle(content []byte) string {
	dictionaryDepth := 0
	for index := 0; index < len(content); {
		switch content[index] {
		case '%':
			for index < len(content) && content[index] != '\r' && content[index] != '\n' {
				index++
			}
		case '(':
			_, next, _ := decodePDFLiteral(content, index)
			index = next
		case '<':
			if index+1 < len(content) && content[index+1] == '<' {
				dictionaryDepth++
				index += 2
			} else {
				index++
			}
		case '>':
			if index+1 < len(content) && content[index+1] == '>' {
				dictionaryDepth--
				index += 2
			} else {
				index++
			}
		case '/':
			name, next := parsePDFName(content, index+1)
			if dictionaryDepth == 1 && name == "Title" {
				valueStart := skipPDFWhitespace(content, next)
				if valueStart < len(content) && content[valueStart] == '(' {
					title, _, ok := decodePDFLiteral(content, valueStart)
					if ok {
						return cleanDocumentTitle(title)
					}
				}
			}
			index = next
		default:
			index++
		}
	}
	return ""
}

func parsePDFName(content []byte, index int) (string, int) {
	start := index
	for index < len(content) && !isPDFNameDelimiter(content[index]) {
		index++
	}
	return string(content[start:index]), index
}

func isPDFNameDelimiter(value byte) bool {
	return value == ' ' || value == '\t' || value == '\r' || value == '\n' || value == '\f' || value == 0 || value == '(' || value == ')' || value == '<' || value == '>' || value == '[' || value == ']' || value == '{' || value == '}' || value == '/' || value == '%'
}

func decodePDFLiteral(content []byte, index int) (string, int, bool) {
	var value strings.Builder
	depth := 1
	for index = index + 1; index < len(content); index++ {
		current := content[index]
		if current == '\\' {
			if index+1 >= len(content) {
				break
			}
			index++
			escaped := content[index]
			if escaped >= '0' && escaped <= '7' {
				octal := escaped - '0'
				for digits := 1; digits < 3 && index+1 < len(content); digits++ {
					next := content[index+1]
					if next < '0' || next > '7' {
						break
					}
					index++
					octal = octal*8 + content[index] - '0'
				}
				value.WriteByte(octal)
				continue
			}
			switch escaped {
			case 'n':
				value.WriteByte('\n')
			case 'r':
				value.WriteByte('\r')
			case 't':
				value.WriteByte('\t')
			case 'b':
				value.WriteByte('\b')
			case 'f':
				value.WriteByte('\f')
			default:
				value.WriteByte(escaped)
			}
			continue
		}
		switch current {
		case '(':
			depth++
			value.WriteByte(current)
		case ')':
			depth--
			if depth == 0 {
				return value.String(), index + 1, true
			}
			value.WriteByte(current)
		default:
			value.WriteByte(current)
		}
	}
	return "", len(content), false
}

func inferredPDFDocumentTitle(text string) string {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) < 2 {
		return ""
	}
	first := cleanDocumentTitle(lines[0])
	if first == "" {
		return ""
	}
	// Metadata-free book covers commonly start with a multi-word ALL-CAPS
	// title followed immediately by author/translator credits. This shape is
	// stronger evidence than a single ALL-CAPS chapter label.
	if likelyPDFCoverTitle(first, lines[1]) {
		return first
	}
	if strings.TrimSpace(lines[1]) != "" || likelyPDFChapterHeading(lines, first) {
		return ""
	}
	return first
}

func likelyPDFCoverTitle(value, followingLine string) bool {
	return len(strings.Fields(value)) >= 3 &&
		isAllCapsHeadingLine(value) &&
		!looksLikeChapterHeading(value) &&
		looksLikePDFCreditLine(followingLine)
}

func looksLikePDFCreditLine(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" || strings.HasSuffix(value, ".") || strings.HasSuffix(value, "!") || strings.HasSuffix(value, "?") {
		return false
	}
	words := strings.Fields(value)
	if len(words) < 2 || len(words) > 6 {
		return false
	}
	for _, word := range words {
		word = strings.Trim(word, ".,;:")
		runes := []rune(word)
		if len(runes) == 0 || !unicode.IsUpper(runes[0]) {
			return false
		}
		for _, r := range runes[1:] {
			if unicode.IsLetter(r) && !unicode.IsLower(r) {
				return false
			}
		}
	}
	return true
}

func likelyPDFChapterHeading(lines []string, title string) bool {
	if looksLikeChapterHeading(title) || isAllCapsHeadingLine(title) {
		return true
	}
	// A title-cased character heading is ambiguous on its own, but not when
	// the PDF repeats the same standalone-heading shape for later chapters.
	return looksLikeTitleCaseHeading(title) && hasRepeatedTitleCaseHeadings(lines)
}

// hasRepeatedTitleCaseHeadings identifies title-cased chapter labels in both
// plain text/PDF extraction and markdown-normalized DOCX content. A candidate
// must be standalone unless it is explicitly marked with one to three #s.
func hasRepeatedTitleCaseHeadings(lines []string) bool {
	count := 0
	for index, line := range lines {
		candidate, marked := titleCaseHeadingCandidate(line)
		if !looksLikeTitleCaseHeading(candidate) || looksLikeChapterHeading(candidate) {
			continue
		}
		if marked || ((index == 0 || strings.TrimSpace(lines[index-1]) == "") && (index == len(lines)-1 || strings.TrimSpace(lines[index+1]) == "")) {
			count++
			if count >= 2 {
				return true
			}
		}
	}
	return false
}

func titleCaseHeadingCandidate(line string) (string, bool) {
	line = strings.TrimSpace(line)
	markerLen := 0
	for markerLen < len(line) && markerLen < 3 && line[markerLen] == '#' {
		markerLen++
	}
	if markerLen > 0 && markerLen < len(line) && line[markerLen] == ' ' {
		return cleanDocumentTitle(line[markerLen+1:]), true
	}
	return cleanDocumentTitle(line), false
}

func looksLikeTitleCaseHeading(value string) bool {
	words := strings.Fields(value)
	if len(words) == 0 || len(words) > 8 {
		return false
	}
	for _, word := range words {
		runes := []rune(word)
		if len(runes) == 0 || !unicode.IsUpper(runes[0]) {
			return false
		}
		for _, r := range runes[1:] {
			if !unicode.IsLower(r) {
				return false
			}
		}
	}
	return true
}

// utf8BOM is the 3-byte UTF-8 byte-order-mark, written as an escaped byte
// sequence (not a literal rune) so gofmt/the compiler don't mistake it for a
// file-leading BOM, which is only meaningful at byte offset 0.
const utf8BOM = "\xEF\xBB\xBF"

// stripBOM removes a leading UTF-8 byte-order-mark, if present.
func stripBOM(s string) string {
	return strings.TrimPrefix(s, utf8BOM)
}

// ponytail: 50MB of decompressed document.xml is far beyond any real
// manuscript. This ceiling guards against a decompression bomb (a tiny .docx
// that expands to gigabytes and OOMs the shared ingestion process). It caps
// both the declared uncompressed size and the bytes actually streamed, so a
// lying central directory can't bypass it. Package-level var so tests can
// lower it without writing a multi-GB fixture.
var maxDecompressedDOCXBytes int64 = 50 << 20

// maxDOCXZipEntries rejects an absurd central directory (zip-bomb defense in
// depth) before we even scan for word/document.xml. A real .docx has a
// handful of entries, not thousands.
const maxDOCXZipEntries = 1000

// parseDOCX extracts plain text from a .docx file's word/document.xml using
// only the stdlib (archive/zip + encoding/xml) — no external dependency.
// Paragraphs styled Heading1/Heading2/Heading3 are prefixed with markdown
// heading markers so the existing splitChunks markdown pattern picks them up.
func parseDOCX(content []byte) (string, error) {
	parsed, err := parseDOCXDetails(content)
	return parsed.text, err
}

func parseDOCXDetails(content []byte) (parsedDocument, error) {
	zr, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return parsedDocument{}, fmt.Errorf("not a valid .docx (legacy .doc? Save as .docx): %w", err)
	}

	if len(zr.File) > maxDOCXZipEntries {
		return parsedDocument{}, fmt.Errorf("not a valid .docx: too many zip entries (%d)", len(zr.File))
	}

	var docFile *zip.File
	for _, f := range zr.File {
		if f.Name == "word/document.xml" {
			docFile = f
			break
		}
	}
	if docFile == nil {
		return parsedDocument{}, fmt.Errorf("not a valid .docx: missing word/document.xml")
	}
	if docFile.UncompressedSize64 > uint64(maxDecompressedDOCXBytes) {
		return parsedDocument{}, fmt.Errorf("document too large when decompressed (%d bytes exceeds %d byte cap)", docFile.UncompressedSize64, maxDecompressedDOCXBytes)
	}

	rc, err := docFile.Open()
	if err != nil {
		return parsedDocument{}, fmt.Errorf("open word/document.xml: %w", err)
	}
	defer rc.Close()

	var buf strings.Builder
	// Defense in depth: cap the bytes actually streamed in case the central
	// directory under-reported UncompressedSize64.
	dec := xml.NewDecoder(io.LimitReader(rc, maxDecompressedDOCXBytes+1))

	var curStyle string
	var curText strings.Builder
	inText := false

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return parsedDocument{}, fmt.Errorf("parse word/document.xml: %w", err)
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "p":
				curStyle = ""
				curText.Reset()
			case "pStyle":
				for _, a := range t.Attr {
					if a.Name.Local == "val" {
						curStyle = a.Value
					}
				}
			case "t":
				inText = true
			}
		case xml.CharData:
			if inText {
				curText.Write(t)
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "t":
				inText = false
			case "p":
				text := curText.String()
				switch curStyle {
				case "Heading1":
					buf.WriteString("# ")
				case "Heading2":
					buf.WriteString("## ")
				case "Heading3":
					buf.WriteString("## ")
				}
				buf.WriteString(text)
				buf.WriteString("\n\n")
				if int64(buf.Len()) > maxDecompressedDOCXBytes {
					return parsedDocument{}, fmt.Errorf("document too large when decompressed (exceeds %d byte cap)", maxDecompressedDOCXBytes)
				}
			}
		}
	}

	text := strings.TrimSpace(buf.String())
	title := docxCoreTitle(zr)
	if title == "" {
		title = markdownDocumentTitleWithH1Priority(text, false)
	}
	return parsedDocument{text: text, title: title}, nil
}

func docxCoreTitle(zr *zip.Reader) string {
	for _, file := range zr.File {
		if file.Name != "docProps/core.xml" {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			return ""
		}
		defer rc.Close()
		decoder := xml.NewDecoder(io.LimitReader(rc, 1<<20))
		inTitle := false
		for {
			token, err := decoder.Token()
			if err != nil {
				return ""
			}
			switch current := token.(type) {
			case xml.StartElement:
				inTitle = current.Name.Local == "title"
			case xml.CharData:
				if inTitle {
					return cleanDocumentTitle(string(current))
				}
			case xml.EndElement:
				if current.Name.Local == "title" {
					inTitle = false
				}
			}
		}
	}
	return ""
}

// parsePDF extracts plain text from a PDF.
//
// Primary method: github.com/ledongthuc/pdf (pure Go). Fallback: pdftotext
// (poppler-utils) when the Go parser fails — catches PDFs with unusual headers
// or encoding that the lightweight Go library rejects.
//
// recover() is mandatory: the Go parser panics on malformed PDFs. A recovered
// panic falls through to pdftotext instead of failing the job.
func parsePDF(content []byte) (text string, err error) {
	defer func() {
		if r := recover(); r != nil {
			text = ""
			err = fmt.Errorf("malformed PDF (parser panic: %v)", r)
		}
	}()

	// ponytail: Go parser first (fast, no subprocess), pdftotext fallback on
	// any failure (catches the edge cases the Go parser can't handle).
	if t, e := parsePDFWithGo(content); e == nil {
		return t, nil
	}
	return parsePDFWithPdftotext(content)
}

// parsePDFWithGo uses the pure-Go ledongthuc/pdf library.
func parsePDFWithGo(content []byte) (string, error) {
	reader, err := pdf.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return "", fmt.Errorf("PDF is encrypted or unreadable: %w", err)
	}

	// GetPlainText can flatten positioned PDF glyphs into a single paragraph.
	// Prefer the library's row grouping so manuscript headings remain on their
	// own lines for chapter detection; keep the older plain-text path as a
	// fallback for PDFs where row extraction is unavailable.
	if text, rowErr := parsePDFRows(reader); rowErr == nil && text != "" {
		return text, nil
	}

	plain, err := reader.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("extract PDF text: %w", err)
	}

	b, err := io.ReadAll(plain)
	if err != nil {
		return "", fmt.Errorf("read extracted PDF text: %w", err)
	}

	extracted := strings.TrimSpace(string(b))
	if extracted == "" {
		return "", fmt.Errorf("no extractable text — scanned or image-only PDF? OCR is not supported")
	}

	return extracted, nil
}

func parsePDFRows(reader *pdf.Reader) (string, error) {
	var text strings.Builder
	for pageNumber := 1; pageNumber <= reader.NumPage(); pageNumber++ {
		page := reader.Page(pageNumber)
		rows := pdfGeometryRows(page)
		if plain, err := pdfPagePlainText(page); err == nil && plain != "" {
			rows = alignPDFRows(rows, plain)
		}
		for _, line := range rows {
			if line == "" {
				continue
			}
			if text.Len() > 0 {
				text.WriteByte('\n')
			}
			text.WriteString(line)
		}
	}
	extracted := strings.TrimSpace(text.String())
	if !strings.Contains(extracted, "\n") {
		// Some simple PDFs encode line movement in a form the geometry API does
		// not expose. Let GetPlainText retain its established behavior there.
		return "", nil
	}
	return extracted, nil
}

func pdfGeometryRows(page pdf.Page) []string {
	fragments := append([]pdf.Text(nil), page.Content().Text...)
	sort.Slice(fragments, func(i, j int) bool {
		if absPDFCoordinate(fragments[i].Y-fragments[j].Y) > 2 {
			return fragments[i].Y > fragments[j].Y
		}
		return fragments[i].X < fragments[j].X
	})

	rows := make([]string, 0)
	for start := 0; start < len(fragments); {
		end := start + 1
		for end < len(fragments) && absPDFCoordinate(fragments[start].Y-fragments[end].Y) <= 2 {
			end++
		}
		for _, line := range strings.Split(pdfLineText(fragments[start:end]), "\n") {
			if line = strings.TrimSpace(line); line != "" {
				rows = append(rows, line)
			}
		}
		start = end
	}
	return rows
}

func pdfPagePlainText(page pdf.Page) (string, error) {
	fonts := make(map[string]*pdf.Font)
	for _, name := range page.Fonts() {
		font := page.Font(name)
		fonts[name] = &font
	}
	return page.GetPlainText(fonts)
}

// alignPDFRows retains geometry-derived line boundaries while consuming the
// authoritative per-page lexical stream. A glyph that cannot be matched at
// the current canonical position is layout noise and is omitted; whitespace
// absent from positioned glyphs is retained from the canonical stream.
func alignPDFRows(rows []string, canonical string) []string {
	canonicalRunes := []rune(canonical)
	position := 0
	aligned := make([]string, 0, len(rows))
	for _, row := range rows {
		var line strings.Builder
		for _, glyph := range []rune(row) {
			for position < len(canonicalRunes) && unicode.IsSpace(canonicalRunes[position]) && !unicode.IsSpace(glyph) {
				line.WriteRune(canonicalRunes[position])
				position++
			}
			if position < len(canonicalRunes) && glyph == canonicalRunes[position] {
				line.WriteRune(canonicalRunes[position])
				position++
				continue
			}
			if unicode.IsSpace(glyph) {
				continue
			}
			if position+1 < len(canonicalRunes) && glyph == canonicalRunes[position+1] {
				line.WriteRune(canonicalRunes[position])
				line.WriteRune(canonicalRunes[position+1])
				position += 2
			}
		}
		if line := strings.TrimSpace(line.String()); line != "" {
			aligned = append(aligned, line)
		}
	}
	return aligned
}

func absPDFCoordinate(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}

func pdfLineText(fragments []pdf.Text) string {
	var line strings.Builder
	for index, fragment := range fragments {
		if index > 0 {
			previous := fragments[index-1]
			gap := fragment.X - (previous.X + previous.W)
			if gap > maxPDFFragmentFontSize(previous, fragment)*1.5 {
				line.WriteByte('\n')
				line.WriteString(fragment.S)
				continue
			}
			if gap > 1 && !strings.HasSuffix(previous.S, " ") && !strings.HasPrefix(fragment.S, " ") {
				line.WriteByte(' ')
			}
		}
		line.WriteString(fragment.S)
	}
	return strings.TrimSpace(line.String())
}

func maxPDFFragmentFontSize(left, right pdf.Text) float64 {
	if left.FontSize > right.FontSize {
		return left.FontSize
	}
	return right.FontSize
}

// parsePDFWithPdftotext uses the system pdftotext command (poppler-utils).
// Falls back gracefully when pdftotext is not installed.
func parsePDFWithPdftotext(content []byte) (string, error) {
	cmd := exec.Command("pdftotext", "-", "-")
	cmd.Stdin = bytes.NewReader(content)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("pdftotext failed: %w", err)
	}
	extracted := strings.TrimSpace(string(out))
	if extracted == "" {
		// ponytail: same message as the Go parser — the caller doesn't need to
		// know which backend was tried. A scan-only PDF fails the same way.
		return "", fmt.Errorf("no extractable text — scanned or image-only PDF? OCR is not supported")
	}
	return extracted, nil
}
