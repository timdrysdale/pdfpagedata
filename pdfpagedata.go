package pdfpagedata

import (
	"encoding/json"
	"math/rand"
	"strings"
	"time"

	"github.com/timdrysdale/unipdf/v3/creator"
	"github.com/timdrysdale/unipdf/v3/extractor"
	pdf "github.com/timdrysdale/unipdf/v3/model"
)

func UnmarshalPageData(page *pdf.PdfPage) ([]PageData, error) {

	pageDatas := []PageData{}

	tokens, err := ReadPageData(page)

	if err != nil {
		return pageDatas, err
	}

	var lastError error

	for _, token := range tokens {

		var pd PageData

		if err := json.Unmarshal([]byte(token), &pd); err != nil {
			lastError = err
			continue
		}

		pageDatas = append(pageDatas, pd)

	}

	return pageDatas, lastError

}

func ReadPageData(page *pdf.PdfPage) ([]string, error) {

	text, err := ReadPageString(page)

	if err != nil {
		return []string{text}, err
	}

	return ExtractPageData(text), nil

}

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

func WritePageData(c *creator.Creator, text string) {
	WritePageString(c, StartTag+text+EndTag)
}

func WritePageString(c *creator.Creator, text string) {
	p := c.NewParagraph(text)
	p.SetFontSize(0.000001)
	rand.Seed(time.Now().UnixNano())
	x := rand.Float64()*0.1 + 99999 //0.3
	y := rand.Float64()*999 + 99999 //0.3
	p.SetPos(x, y)
	c.Draw(p)
}
