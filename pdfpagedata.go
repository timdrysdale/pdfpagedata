package pdfpagedata

import (
	"strings"

	"github.com/timdrysdale/unipdf/v3/creator"
	"github.com/timdrysdale/unipdf/v3/extractor"
	pdf "github.com/timdrysdale/unipdf/v3/model"
)

const (
	StartTag       = "<gradex-pagedata>"
	EndTag         = "</gradex-pagedata>"
	StartTagOffset = len(StartTag)
	EndTagOffset   = len(EndTag)
)

func ExtractPageData(pageText string) []string {

	var tokens []string

LOOP:
	for {

		startIndex := strings.Index(pageText, StartTag)
		if startIndex < 0 {
			break LOOP
		}

		endIndex := strings.Index(pageText, EndTag)
		if endIndex < 0 {
			break LOOP
		}

		token := pageText[startIndex+StartTagOffset : endIndex]

		tokens = append(tokens, token)

		pageText = pageText[endIndex+EndTagOffset : len(pageText)]

	}

	return tokens
}

func ReadPageString(page *pdf.PdfPage) (string, error) {

	ex, err := extractor.New(page)
	if err != nil {
		return "", err
	}

	text, err := ex.ExtractText()
	return text, err
}

func WritePageString(c *creator.Creator, text string) {
	p := c.NewParagraph(text)
	p.SetFontSize(0.1)
	p.SetPos(0.1, 0.1)
	c.Draw(p)
}
