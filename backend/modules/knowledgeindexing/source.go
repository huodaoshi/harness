package knowledgeindexing

import (
	"bytes"
	"fmt"
	"path"
	"strings"
)

// SourceFormat identifies document content format for parser selection.
type SourceFormat string

const (
	SourceFormatText     SourceFormat = "text"
	SourceFormatMarkdown SourceFormat = "markdown"
	SourceFormatJSON     SourceFormat = "json"
	SourceFormatHTML     SourceFormat = "html"
	SourceFormatExcel    SourceFormat = "excel"
	SourceFormatPDF      SourceFormat = "pdf"
	SourceFormatCSV      SourceFormat = "csv"
)

// ResolvedSource is a fetched document with detected format.
type ResolvedSource struct {
	Body        []byte
	ContentType string
	FinalURL    string
	FileName    string
	Format      SourceFormat
}

func detectSourceFormat(sourceURL, contentType string, body []byte) SourceFormat {
	ext := strings.ToLower(path.Ext(sourceURL))
	switch ext {
	case ".md", ".markdown":
		return SourceFormatMarkdown
	case ".json":
		return SourceFormatJSON
	case ".html", ".htm":
		return SourceFormatHTML
	case ".txt":
		return SourceFormatText
	}

	ct := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	switch {
	case strings.Contains(ct, "markdown"):
		return SourceFormatMarkdown
	case strings.Contains(ct, "json"):
		return SourceFormatJSON
	case strings.Contains(ct, "html"):
		return SourceFormatHTML
	case strings.Contains(ct, "text/"):
		return SourceFormatText
	}

	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return SourceFormatText
	}
	if trimmed[0] == '{' || trimmed[0] == '[' {
		return SourceFormatJSON
	}
	return SourceFormatText
}

func resolveSource(raw []byte, sourceURL, contentType, finalURL string) ResolvedSource {
	format := detectSourceFormat(sourceURL, contentType, raw)
	fileName := path.Base(sourceURL)
	if fileName == "." || fileName == "/" {
		fileName = ""
	}
	return ResolvedSource{
		Body:        raw,
		ContentType: contentType,
		FinalURL:    finalURL,
		FileName:    fileName,
		Format:      format,
	}
}

func unsupportedFormatErr(format SourceFormat) error {
	return fmt.Errorf("unsupported format %q", format)
}
