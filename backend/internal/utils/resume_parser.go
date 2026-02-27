package utils

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"

	"ithozyeva/internal/models"
)

type ParsedResumeData struct {
	WorkExperience  string
	DesiredPosition string
	WorkFormat      models.WorkFormat
	Confidence      float64
}

func ParseResume(filename string, content []byte) (*ParsedResumeData, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	var text string
	var err error

	switch ext {
	case ".pdf":
		text, err = extractTextFromPDF(content)
	case ".docx":
		text, err = extractTextFromDocx(content)
	default:
		return &ParsedResumeData{}, fmt.Errorf("format %s is not supported for parsing", ext)
	}

	if err != nil {
		return nil, err
	}

	data := analyzeText(text)
	return data, nil
}

func extractTextFromPDF(content []byte) (string, error) {
	reader, err := pdf.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	numPages := reader.NumPage()
	for i := 1; i <= numPages; i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}
		textExtractor, err := page.GetPlainText(nil)
		if err != nil {
			return "", err
		}
		builder.WriteString(textExtractor)
		builder.WriteString("\n")
	}
	return builder.String(), nil
}

func extractTextFromDocx(content []byte) (string, error) {
	readerAt := bytes.NewReader(content)
	zr, err := zip.NewReader(readerAt, int64(len(content)))
	if err != nil {
		return "", err
	}

	for _, file := range zr.File {
		if file.Name != "word/document.xml" {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()

		var builder strings.Builder
		decoder := xml.NewDecoder(rc)
		for {
			tok, err := decoder.Token()
			if err != nil {
				if err == io.EOF {
					break
				}
				return "", err
			}

			switch t := tok.(type) {
			case xml.CharData:
				builder.WriteString(string(t))
			case xml.EndElement:
				if strings.HasSuffix(t.Name.Local, "p") {
					builder.WriteString("\n")
				}
			}
		}
		return builder.String(), nil
	}

	return "", fmt.Errorf("document.xml not found in docx")
}

func analyzeText(text string) *ParsedResumeData {
	lines := strings.Split(text, "\n")
	var (
		workExp string
		role    string
		format  models.WorkFormat
		matches int
	)

	for _, line := range lines {
		clean := strings.TrimSpace(line)
		if clean == "" {
			continue
		}
		lower := strings.ToLower(clean)

		if workExp == "" && (strings.Contains(lower, "опыт") || strings.Contains(lower, "experience")) {
			workExp = clean
			matches++
		}

		if role == "" && (strings.Contains(lower, "должн") || strings.Contains(lower, "position") || strings.Contains(lower, "role")) {
			role = clean
			matches++
		}

		if format == "" {
			switch {
			case strings.Contains(lower, "удален"):
				format = models.WorkFormatRemote
				matches++
			case strings.Contains(lower, "гибрид"):
				format = models.WorkFormatHybrid
				matches++
			case strings.Contains(lower, "офис"):
				format = models.WorkFormatOffice
				matches++
			}
		}

		if matches == 3 {
			break
		}
	}

	confidence := float64(matches) / 3
	return &ParsedResumeData{
		WorkExperience:  workExp,
		DesiredPosition: role,
		WorkFormat:      format,
		Confidence:      confidence,
	}
}
