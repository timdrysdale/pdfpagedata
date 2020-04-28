package pdfpagedata

import (
	"encoding/json"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/timdrysdale/unipdf/v3/creator"
	"github.com/timdrysdale/unipdf/v3/extractor"
	pdf "github.com/timdrysdale/unipdf/v3/model"
)

func GetLen(input map[int][]PageData) int {
	items := 0
	for _, v := range input {
		for _ = range v {
			items++
		}
	}
	return items
}

func GetPageDataFromFile(inputPath string) (map[int][]PageData, error) {

	docData := make(map[int][]PageData)

	texts, err := OutputPdfText(inputPath)
	if err != nil {
		return docData, err
	}

	var pds []PageData

	//one text per page
	for i, text := range texts {

		strs := ExtractPageData(text)

		for _, str := range strs {
			var pd PageData

			if err := json.Unmarshal([]byte(str), &pd); err != nil {
				continue
			}

			pds = append(pds, pd)
		}

		docData[i] = pds
	}

	return docData, nil

}

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

func MarshalPageData(c *creator.Creator, pd *PageData) error {

	token, err := json.Marshal(pd)
	if err != nil {
		return err
	}

	WritePageData(c, string(token))

	return nil

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

// this function is for use in a co-operative
// environment - you can slip one past the gaolie
// in the custom fields in Questions/Processing/Custom
func StripAuthorIdentity(pd PageData) PageData {

	safe := PageData{}

	safe.Exam = pd.Exam
	safe.Author = AuthorDetails{Anonymous: pd.Author.Anonymous}
	safe.Page = pd.Page
	safe.Contact = pd.Contact
	safe.Submission = SubmissionDetails{} //nothing!
	safe.Questions = pd.Questions
	safe.Processing = pd.Processing
	safe.Custom = pd.Custom

	return safe
}

// outputPdfText produces array of strings, one string per page
// mod from https://github.com/unidoc/unipdf-examples/blob/master/text/pdf_extract_text.go
func OutputPdfText(inputPath string) ([]string, error) {

	texts := []string{}

	f, err := os.Open(inputPath)
	if err != nil {
		return texts, err
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return texts, err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return texts, err
	}

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return texts, err
		}

		ex, err := extractor.New(page)
		if err != nil {
			return texts, err
		}

		text, err := ex.ExtractText()
		if err != nil {
			return texts, err
		}

		texts = append(texts, text)
	}

	return texts, nil
}
