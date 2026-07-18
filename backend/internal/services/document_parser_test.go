package services

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// buildTestDOCX generates a minimal .docx (a zip containing word/document.xml)
// in-test for determinism, instead of committing a binary fixture — matches
// the design's stated approach for this format.
func buildTestDOCX(t *testing.T, paragraphs []struct{ style, text string }) []byte {
	t.Helper()
	var body strings.Builder
	for _, p := range paragraphs {
		body.WriteString(`<w:p>`)
		if p.style != "" {
			body.WriteString(fmt.Sprintf(`<w:pPr><w:pStyle w:val="%s"/></w:pPr>`, p.style))
		}
		body.WriteString(fmt.Sprintf(`<w:r><w:t>%s</w:t></w:r>`, p.text))
		body.WriteString(`</w:p>`)
	}
	documentXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
<w:body>` + body.String() + `</w:body>
</w:document>`

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	f, err := zw.Create("word/document.xml")
	if err != nil {
		t.Fatalf("create zip entry: %v", err)
	}
	if _, err := f.Write([]byte(documentXML)); err != nil {
		t.Fatalf("write zip entry: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
	return buf.Bytes()
}

func buildTestDOCXWithCoreTitle(t *testing.T, title string, paragraphs []struct{ style, text string }) []byte {
	t.Helper()
	content := buildTestDOCX(t, paragraphs)
	zr, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		t.Fatalf("open generated docx: %v", err)
	}
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, file := range zr.File {
		reader, openErr := file.Open()
		if openErr != nil {
			t.Fatalf("open generated entry: %v", openErr)
		}
		data, readErr := io.ReadAll(reader)
		reader.Close()
		if readErr != nil {
			t.Fatalf("read generated entry: %v", readErr)
		}
		writer, createErr := zw.Create(file.Name)
		if createErr != nil {
			t.Fatalf("copy generated entry: %v", createErr)
		}
		if _, writeErr := writer.Write(data); writeErr != nil {
			t.Fatalf("write generated entry: %v", writeErr)
		}
	}
	core, err := zw.Create("docProps/core.xml")
	if err != nil {
		t.Fatalf("create core properties: %v", err)
	}
	if _, err := core.Write([]byte(`<?xml version="1.0"?><cp:coreProperties xmlns:cp="urn:core" xmlns:dc="http://purl.org/dc/elements/1.1/"><dc:title>` + title + `</dc:title></cp:coreProperties>`)); err != nil {
		t.Fatalf("write core properties: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close generated docx: %v", err)
	}
	return buf.Bytes()
}

// buildTestPDF generates a minimal single-page valid PDF containing the given
// text, computing correct xref byte offsets — hand-crafted structure per the
// design, generated in-code (rather than a committed binary) so the exact
// byte offsets are always correct and the fixture is trivially editable.
func buildTestPDF(text string) []byte {
	return buildTestPDFWithTitle(text, "")
}

func buildTestPDFWithTitle(text, title string) []byte {
	return buildTestPDFWithTitleLiteral(text, strings.NewReplacer("\\", "\\\\", "(", "\\(", ")", "\\)").Replace(title))
}

func buildTestPDFWithTitleLiteral(text, titleLiteral string) []byte {
	return buildTestPDFWithOutlineAndTitleLiteral(text, "", titleLiteral)
}

func buildTestPDFWithOutlineAndTitleLiteral(text, outlineTitleLiteral, titleLiteral string) []byte {
	return buildTestPDFWithOutlineSubjectAndTitleLiteral(text, outlineTitleLiteral, "", titleLiteral)
}

func buildTestPDFWithSubjectAndTitleLiteral(text, subjectLiteral, titleLiteral string) []byte {
	return buildTestPDFWithOutlineSubjectAndTitleLiteral(text, "", subjectLiteral, titleLiteral)
}

func buildTestPDFWithOutlineSubjectAndTitleLiteral(text, outlineTitleLiteral, subjectLiteral, titleLiteral string) []byte {
	var buf bytes.Buffer
	objectCount := 6
	if outlineTitleLiteral != "" {
		objectCount = 7
	}
	offsets := make([]int, objectCount+1) // index 1..objectCount used

	write := func(s string) { buf.WriteString(s) }

	write("%PDF-1.4\n")

	offsets[1] = buf.Len()
	catalog := "1 0 obj\n<< /Type /Catalog /Pages 2 0 R"
	if outlineTitleLiteral != "" {
		catalog += " /Outlines 7 0 R"
	}
	write(catalog + " >>\nendobj\n")

	offsets[2] = buf.Len()
	write("2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")

	offsets[3] = buf.Len()
	write("3 0 obj\n<< /Type /Page /Parent 2 0 R /Resources << /Font << /F1 4 0 R >> >> /MediaBox [0 0 612 792] /Contents 5 0 R >>\nendobj\n")

	offsets[4] = buf.Len()
	write("4 0 obj\n<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>\nendobj\n")

	var stream strings.Builder
	stream.WriteString("BT /F1 24 Tf 72 712 Td ")
	for index, line := range strings.Split(text, "\n") {
		if index > 0 {
			stream.WriteString("T* ")
		}
		stream.WriteString("(")
		stream.WriteString(strings.NewReplacer("\\", "\\\\", "(", "\\(", ")", "\\)").Replace(line))
		stream.WriteString(") Tj ")
	}
	stream.WriteString("ET")
	offsets[5] = buf.Len()
	write(fmt.Sprintf("5 0 obj\n<< /Length %d >>\nstream\n%s\nendstream\nendobj\n", stream.Len(), stream.String()))

	if outlineTitleLiteral != "" {
		offsets[7] = buf.Len()
		write(fmt.Sprintf("7 0 obj\n<< /Type /Outline /Title (%s) >>\nendobj\n", outlineTitleLiteral))
	}

	offsets[6] = buf.Len()
	info := "6 0 obj\n<<"
	if subjectLiteral != "" {
		info += fmt.Sprintf(" /Subject (%s)", subjectLiteral)
	}
	info += fmt.Sprintf(" /Title (%s) >>\nendobj\n", titleLiteral)
	write(info)

	xrefStart := buf.Len()
	write(fmt.Sprintf("xref\n0 %d\n", objectCount+1))
	write("0000000000 65535 f \n")
	for i := 1; i <= objectCount; i++ {
		write(fmt.Sprintf("%010d 00000 n \n", offsets[i]))
	}
	write(fmt.Sprintf("trailer\n<< /Size %d /Root 1 0 R /Info 6 0 R >>\n", objectCount+1))
	write(fmt.Sprintf("startxref\n%d\n%%%%EOF", xrefStart))

	return buf.Bytes()
}

func TestFileTypeOf(t *testing.T) {
	cases := map[string]string{
		"manuscript.md":   "md",
		"manuscript.TXT":  "txt",
		"manuscript.docx": "docx",
		"manuscript.PDF":  "pdf",
		"manuscript.doc":  "doc",
		"noext":           "",
	}
	for name, want := range cases {
		if got := fileTypeOf(name); got != want {
			t.Errorf("fileTypeOf(%q) = %q, want %q", name, got, want)
		}
	}
}

func TestParseDocumentPassthrough(t *testing.T) {
	for _, ext := range []string{"md", "txt"} {
		filename := "manuscript." + ext
		text, err := parseDocument(filename, []byte("Hello world."))
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", filename, err)
		}
		if text != "Hello world." {
			t.Errorf("%s: text = %q, want %q", filename, text, "Hello world.")
		}
	}
}

func TestParseDocumentDetailsUsesExplicitDocumentTitles(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		content  []byte
		want     string
	}{
		{
			name:     "markdown frontmatter",
			filename: "upload.md",
			content:  []byte("---\ntitle: The Long Way Home\n---\n\n# Chapter One\n\nText."),
			want:     "The Long Way Home",
		},
		{
			name:     "markdown h1",
			filename: "upload.md",
			content:  []byte("# The Long Way Home\n\n## Chapter One\n\nText."),
			want:     "The Long Way Home",
		},
		{
			name:     "markdown chapter heading is not a work title",
			filename: "upload.md",
			content:  []byte("# Chapter One\n\nText."),
			want:     "",
		},
		{
			name:     "explicit markdown h1 takes precedence over later title case headings",
			filename: "upload.md",
			content:  []byte("# The Long Way Home\n\n## Holden\n\nHe arrived at dusk.\n\n## Miller\n\nNobody followed."),
			want:     "The Long Way Home",
		},
		{
			name:     "text repeated title case headings fall back to filename",
			filename: "upload.txt",
			content:  []byte("Holden\n\nHe arrived at dusk.\n\nMiller\n\nNobody followed."),
			want:     "",
		},
		{
			name:     "docx core property",
			filename: "upload.docx",
			content:  buildTestDOCXWithCoreTitle(t, "The Long Way Home", []struct{ style, text string }{{style: "Heading1", text: "Chapter One"}, {text: "Text."}}),
			want:     "The Long Way Home",
		},
		{
			name:     "docx heading one chapter is not a work title",
			filename: "upload.docx",
			content:  buildTestDOCX(t, []struct{ style, text string }{{style: "Heading1", text: "Chapter One"}, {text: "Text."}}),
			want:     "",
		},
		{
			name:     "docx repeated title case headings fall back to filename",
			filename: "upload.docx",
			content:  buildTestDOCX(t, []struct{ style, text string }{{style: "Heading1", text: "Holden"}, {text: "He arrived at dusk."}, {style: "Heading1", text: "Miller"}, {text: "Nobody followed."}}),
			want:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := parseDocumentDetails(tt.filename, tt.content)
			if err != nil {
				t.Fatalf("parse document details: %v", err)
			}
			if parsed.title != tt.want {
				t.Fatalf("title = %q, want %q", parsed.title, tt.want)
			}
			if parsed.text == "" {
				t.Fatal("expected extracted text")
			}
		})
	}
}

func TestParseDocumentStripsBOM(t *testing.T) {
	content := append([]byte(utf8BOM), []byte("Hello world.")...)
	text, err := parseDocument("manuscript.txt", content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "Hello world." {
		t.Errorf("text = %q, want BOM stripped %q", text, "Hello world.")
	}
}

func TestParseDocumentRejectsLegacyDoc(t *testing.T) {
	_, err := parseDocument("manuscript.doc", []byte("binary junk"))
	if err == nil {
		t.Fatal("expected error for .doc, got nil")
	}
}

func TestParseDocumentRejectsUnknownExtension(t *testing.T) {
	_, err := parseDocument("manuscript.rtf", []byte("binary junk"))
	if err == nil {
		t.Fatal("expected error for unknown extension, got nil")
	}
}

func TestParseDOCX(t *testing.T) {
	content := buildTestDOCX(t, []struct{ style, text string }{
		{style: "Heading1", text: "Chapter One"},
		{text: "Bilbo was going to have a birthday party."},
		{style: "Heading2", text: "A Shorter Heading"},
		{text: "Frodo learns about the Ring."},
	})

	text, err := parseDocument("manuscript.docx", content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(text, "# Chapter One") {
		t.Errorf("expected Heading1 prefixed with '# ', got: %q", text)
	}
	if !strings.Contains(text, "## A Shorter Heading") {
		t.Errorf("expected Heading2 prefixed with '## ', got: %q", text)
	}
	if !strings.Contains(text, "Bilbo was going to have a birthday party.") {
		t.Errorf("missing plain paragraph text, got: %q", text)
	}
}

// TestParseDOCXDecompressionBomb verifies a .docx whose word/document.xml
// decompresses past the ceiling is rejected with a clear error rather than
// streamed unbounded into memory (OOM). The ceiling is lowered to a few KB for
// the test so no multi-GB fixture is needed — the highly compressible body is
// tiny on disk but its uncompressed size trips the cap.
func TestParseDOCXDecompressionBomb(t *testing.T) {
	orig := maxDecompressedDOCXBytes
	maxDecompressedDOCXBytes = 4 << 10 // 4KB test-only ceiling
	defer func() { maxDecompressedDOCXBytes = orig }()

	// A valid document.xml wrapper around a large, highly compressible run of
	// bytes — its uncompressed size far exceeds the lowered ceiling.
	documentXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
<w:body><w:p><w:r><w:t>` + strings.Repeat("A", 64<<10) + `</w:t></w:r></w:p></w:body>
</w:document>`

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	f, err := zw.Create("word/document.xml")
	if err != nil {
		t.Fatalf("create zip entry: %v", err)
	}
	if _, err := f.Write([]byte(documentXML)); err != nil {
		t.Fatalf("write zip entry: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}

	// The compressed .docx is only a few KB — a real bomb would be too, which
	// is exactly why an on-disk size check is not enough.
	if buf.Len() > 8<<10 {
		t.Fatalf("test fixture unexpectedly large (%d bytes); repeated bytes should compress tiny", buf.Len())
	}

	_, err = parseDocument("bomb.docx", buf.Bytes())
	if err == nil {
		t.Fatal("expected an error for an over-cap .docx, got nil (would OOM in prod)")
	}
	if !strings.Contains(err.Error(), "too large") {
		t.Errorf("expected a 'too large when decompressed' error, got: %v", err)
	}
}

func TestParseDOCXNotAZip(t *testing.T) {
	_, err := parseDocument("manuscript.docx", []byte("not a zip file"))
	if err == nil {
		t.Fatal("expected error for non-zip content, got nil")
	}
}

func TestParseDOCXMissingDocumentXML(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	f, _ := zw.Create("word/other.xml")
	f.Write([]byte("<xml/>"))
	zw.Close()

	_, err := parseDocument("manuscript.docx", buf.Bytes())
	if err == nil {
		t.Fatal("expected error for missing word/document.xml, got nil")
	}
}

func TestParsePDFValid(t *testing.T) {
	pdfBytes := buildTestPDF("Hello World")
	text, err := parseDocument("manuscript.pdf", pdfBytes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(text, "Hello World") {
		t.Errorf("expected extracted text to contain %q, got: %q", "Hello World", text)
	}
}

func TestParsePDFDocumentTitleSkipsChapterHeadings(t *testing.T) {
	tests := []struct {
		name string
		pdf  []byte
		want string
	}{
		{
			name: "metadata title wins over character headings",
			pdf:  buildTestPDFWithTitle("AMOS\n\nHe arrived at dusk.\n\nBOB\n\nNobody followed.", "The Glass City"),
			want: "The Glass City",
		},
		{
			name: "trailer info title wins over preceding outline title",
			pdf:  buildTestPDFWithOutlineAndTitleLiteral("AMOS\n\nHe arrived at dusk.", "Chapter 1", "Actual Book Title"),
			want: "Actual Book Title",
		},
		{
			name: "subject literal title does not mask info title",
			pdf:  buildTestPDFWithSubjectAndTitleLiteral("AMOS\n\nHe arrived at dusk.", "Subject mentions /Title \\(Fake Chapter\\)", "Actual Book Title"),
			want: "Actual Book Title",
		},
		{
			name: "metadata title decodes escaped parentheses",
			pdf:  buildTestPDFWithTitle("AMOS\n\nHe arrived at dusk.", "The (Glass) City"),
			want: "The (Glass) City",
		},
		{
			name: "metadata title decodes octal parentheses",
			pdf:  buildTestPDFWithTitleLiteral("AMOS\n\nHe arrived at dusk.", "The \\050Glass\\051 City"),
			want: "The (Glass) City",
		},
		{
			name: "all caps character heading falls back to filename",
			pdf:  buildTestPDF("AMOS\n\nHe arrived at dusk.\n\nBOB\n\nNobody followed."),
			want: "",
		},
		{
			name: "repeated title case character headings fall back to filename",
			pdf:  buildTestPDF("Holden\n\nHe arrived at dusk.\n\nMiller\n\nNobody followed."),
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := parseDocumentDetails("manuscript.pdf", tt.pdf)
			if err != nil {
				t.Fatalf("parse PDF: %v", err)
			}
			if parsed.title != tt.want {
				t.Fatalf("title = %q, want %q (text: %q)", parsed.title, tt.want, parsed.text)
			}
		})
	}
}

func TestParsePDFPreservesLinesForCharacterHeadingDetection(t *testing.T) {
	parsed, err := parseDocumentDetails("manuscript.pdf", buildTestPDF("AMOS\n\nHe met Amos at dusk.\n\nBOB\n\nShe met Bob at dawn."))
	if err != nil {
		t.Fatalf("parse PDF: %v", err)
	}
	chunks := (&IngestionService{}).splitChunks(parsed.text)
	if len(chunks) != 2 {
		t.Fatalf("chunks = %d, want AMOS and BOB from parser output %q", len(chunks), parsed.text)
	}
	if chunks[0].title != "AMOS" || chunks[1].title != "BOB" {
		t.Errorf("chunk titles = %q, %q; want AMOS, BOB", chunks[0].title, chunks[1].title)
	}
}

func TestParseRealLeviatanPDFPreservesNarrativeHeadings(t *testing.T) {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve acceptance test source path")
	}
	pdfPath := filepath.Join(filepath.Dir(sourceFile), "..", "..", "..", "Docs", "El-Despertar-Del-Leviatan-Parte-01.pdf")
	content, err := os.ReadFile(pdfPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skipf("skipping real Leviatan PDF acceptance test: optional local fixture is absent at %s", pdfPath)
		}
		t.Fatalf("read real PDF: %v", err)
	}
	parsed, err := parseDocumentDetails(filepath.Base(pdfPath), content)
	if err != nil {
		t.Fatalf("parse real PDF: %v", err)
	}
	if parsed.title != "EL DESPERTAR DEL LEVIATÁN" {
		t.Errorf("metadata-free PDF title = %q, want cover title", parsed.title)
	}
	chunks := (&IngestionService{}).splitChunks(parsed.text)
	titles := make([]string, len(chunks))
	for i, chunk := range chunks {
		titles[i] = chunk.title
	}
	t.Logf("real PDF title=%q sections=%d titles=%q", parsed.title, len(chunks), titles)
	wantTitles := []string{"Front Matter", "Prólogo: Julie", "Holden", "Miller", "Holden", "Miller"}
	if len(titles) != len(wantTitles) {
		t.Fatalf("section count = %d, want %d: %q", len(titles), len(wantTitles), titles)
	}
	for i, want := range wantTitles {
		if titles[i] != want {
			t.Errorf("section %d title = %q, want %q", i, titles[i], want)
		}
	}

	required := map[string]bool{"Front Matter": false, "Prólogo: Julie": false, "Holden": false, "Miller": false}
	for _, chunk := range chunks {
		if _, ok := required[chunk.title]; ok && strings.TrimSpace(chunk.content) != "" {
			required[chunk.title] = true
		}
	}
	for title, hasContent := range required {
		if !hasContent {
			t.Errorf("%q is missing or has no body", title)
		}
	}
}

func TestAlignPDFRowsUsesCanonicalTextAndDropsLayoutGlyphs(t *testing.T) {
	got := alignPDFRows(
		[]string{"EL DESPERTAR DEL LEVIATÁNV", "Prólogo: Julieu", "Holdenu"},
		"EL DESPERTAR DEL LEVIATÁN Prólogo: Julie Holden",
	)
	want := []string{"EL DESPERTAR DEL LEVIATÁN", "Prólogo: Julie", "Holden"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("aligned rows = %q, want %q", got, want)
	}
}

func TestInferredPDFDocumentTitleRecognizesCoverButNotChapterLabel(t *testing.T) {
	cases := []struct {
		name string
		text string
		want string
	}{
		{
			name: "multi word all caps cover with credits",
			text: "EL DESPERTAR DEL LEVIATÁN\nJames S. A. Corey\n\nPrólogo: Julie",
			want: "EL DESPERTAR DEL LEVIATÁN",
		},
		{
			name: "single all caps chapter label",
			text: "HOLDEN\n\nHe woke before dawn.",
			want: "",
		},
		{
			name: "all caps chapter heading followed by prose",
			text: "THE FINAL BATTLE\nThe armies advanced at dawn.",
			want: "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := inferredPDFDocumentTitle(tc.text); got != tc.want {
				t.Errorf("inferredPDFDocumentTitle() = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestParsePDFCorrupt exercises the recover() path in parsePDF: the
// ledongthuc/pdf parser is known to panic on malformed input, and a truncated
// PDF (valid header, then garbage) is a reliable way to trigger that.
func TestParsePDFCorrupt(t *testing.T) {
	valid := buildTestPDF("Hello World")
	corrupt := valid[:len(valid)/2]

	_, err := parseDocument("manuscript.pdf", corrupt)
	if err == nil {
		t.Fatal("expected error for corrupt PDF, got nil")
	}
}

func TestParsePDFNotAPDF(t *testing.T) {
	_, err := parseDocument("manuscript.pdf", []byte("this is not a pdf at all"))
	if err == nil {
		t.Fatal("expected error for non-PDF content, got nil")
	}
}

func TestParseDocumentEmptyExtraction(t *testing.T) {
	text, err := parseDocument("manuscript.md", []byte("   \n\n  "))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(text) != "" {
		t.Errorf("expected whitespace-only passthrough, got: %q", text)
	}
	// Note: parseDocument itself does not reject empty/whitespace-only text —
	// that check lives in runWorker (D1), which is format-agnostic and also
	// catches parse results like an empty DOCX/PDF.
}
