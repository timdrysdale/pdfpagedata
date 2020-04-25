package pdfpagedata

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/mattetti/filebuffer"
	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/unipdf/v3/creator"
	"github.com/timdrysdale/unipdf/v3/model"
	"github.com/timdrysdale/unipdf/v3/model/optimize"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

// Mod from array to slice,
// from https://www.golangprograms.com/golang-check-if-array-element-exists.html
func itemExists(sliceType interface{}, item interface{}) bool {
	slice := reflect.ValueOf(sliceType)

	if slice.Kind() != reflect.Slice {
		panic("Invalid data-type")
	}

	for i := 0; i < slice.Len(); i++ {
		if slice.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func TestWriteReadDouble(t *testing.T) {

	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	c.SetPageSize(creator.PageSizeA4)

	c.NewPage()
	text1a := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":1,\"Batch\":\"a\"}"
	text1b := "{\"exam\":\"ENGI99886\",\"number\":\"B12345\",\"page\":1,\"Batch\":\"xx\"}"
	WritePageData(c, text1a)
	WritePageData(c, text1b)

	c.NewPage()
	text2a := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":2,\"Batch\":\"a\"}"
	text2b := "{\"exam\":\"ENGI99897\",\"number\":\"B12345\",\"page\":2,\"Batch\":\"b\"}"
	WritePageData(c, text2a)
	WritePageData(c, text2b)

	// write to memory instead of a file
	var buf bytes.Buffer

	c.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    90,
		ImageUpperPPI:                   150,
	}))

	err := c.Write(&buf)
	if err != nil {
		t.Error(err)
	}

	// convert buffer to readseeker
	var bufslice []byte
	fbuf := filebuffer.New(bufslice)
	fbuf.Write(buf.Bytes())

	// read in from memory
	pdfReader, err := model.NewPdfReader(fbuf)
	if err != nil {
		t.Error(err)
	}
	page, err := pdfReader.GetPage(1)
	if err != nil {
		t.Error(err)
	}

	textp1, err := ReadPageData(page)

	if err != nil {
		t.Error(err)
	}

	page2, err := pdfReader.GetPage(2)
	if err != nil {
		t.Error(err)
	}

	textp2, err := ReadPageData(page2)
	if err != nil {
		t.Error(err)
	}

	if len(textp1) == 2 {

		assert.True(t, itemExists(textp1, text1a))
		assert.True(t, itemExists(textp1, text1b))
	} else {
		t.Error("Wrong number of page data tokens")
	}

	if len(textp2) == 2 {

		assert.True(t, itemExists(textp2, text2a))
		assert.True(t, itemExists(textp2, text2a))

	} else {
		t.Error("Wrong number of page data tokens")
	}

}

func TestWriteReadDoubleLong(t *testing.T) {

	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	c.SetPageSize(creator.PageSizeA4)

	c.NewPage()
	text1a := strings.Repeat("X", 9999)
	text1b := strings.Repeat("Y", 9999)
	WritePageData(c, text1a)
	WritePageData(c, text1b)

	c.NewPage()
	text2a := strings.Repeat("A", 9999)
	text2b := strings.Repeat("B", 9999)
	WritePageData(c, text2a)
	WritePageData(c, text2b)

	// write to memory instead of a file
	var buf bytes.Buffer

	c.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    90,
		ImageUpperPPI:                   150,
	}))

	err := c.Write(&buf)
	if err != nil {
		t.Error(err)
	}

	// convert buffer to readseeker
	var bufslice []byte
	fbuf := filebuffer.New(bufslice)
	fbuf.Write(buf.Bytes())

	// read in from memory
	pdfReader, err := model.NewPdfReader(fbuf)
	if err != nil {
		t.Error(err)
	}
	page, err := pdfReader.GetPage(1)
	if err != nil {
		t.Error(err)
	}

	textp1, err := ReadPageData(page)

	if err != nil {
		t.Error(err)
	}

	page2, err := pdfReader.GetPage(2)
	if err != nil {
		t.Error(err)
	}

	textp2, err := ReadPageData(page2)
	if err != nil {
		t.Error(err)
	}

	if len(textp1) == 2 {

		assert.True(t, itemExists(textp1, text1a))
		assert.True(t, itemExists(textp1, text1b))
	} else {
		t.Error("Wrong number of page data tokens")
	}

	if len(textp2) == 2 {

		assert.True(t, itemExists(textp2, text2a))
		assert.True(t, itemExists(textp2, text2a))

	} else {
		t.Error("Wrong number of page data tokens")
	}

}

func TestWriteReadOtherText(t *testing.T) {

	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	c.SetPageSize(creator.PageSizeA4)
	c.NewPage()
	text1 := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":1}"
	WritePageData(c, text1)
	text2 := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":2}"
	c.NewPage()
	WritePageData(c, text2)

	p := c.NewParagraph("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	p.SetFontSize(12)
	p.SetPos(0, 0)
	c.Draw(p)

	// write to memory instead of a file
	var buf bytes.Buffer

	c.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    90,
		ImageUpperPPI:                   150,
	}))

	err := c.Write(&buf)
	if err != nil {
		t.Error(err)
	}

	// convert buffer to readseeker
	var bufslice []byte
	fbuf := filebuffer.New(bufslice)
	fbuf.Write(buf.Bytes())

	// read in from memory
	pdfReader, err := model.NewPdfReader(fbuf)
	if err != nil {
		t.Error(err)
	}
	page, err := pdfReader.GetPage(1)
	if err != nil {
		t.Error(err)
	}

	textp1, err := ReadPageData(page)

	if err != nil {
		t.Error(err)
	}

	page2, err := pdfReader.GetPage(2)
	if err != nil {
		t.Error(err)
	}

	textp2, err := ReadPageData(page2)
	if err != nil {
		t.Error(err)
	}

	assertEqual(t, text1, textp1[0])
	assertEqual(t, text2, textp2[0])

}

func TestWriteRead(t *testing.T) {

	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	c.SetPageSize(creator.PageSizeA4)
	c.NewPage()
	text1 := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":1}"
	WritePageData(c, text1)
	text2 := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":2}"
	c.NewPage()
	WritePageData(c, text2)
	c.NewPage()
	// write to memory instead of a file
	var buf bytes.Buffer

	err := c.Write(&buf)
	if err != nil {
		t.Error(err)
	}

	// convert buffer to readseeker
	var bufslice []byte
	fbuf := filebuffer.New(bufslice)
	fbuf.Write(buf.Bytes())

	// read in from memory
	pdfReader, err := model.NewPdfReader(fbuf)
	if err != nil {
		t.Error(err)
	}

	page, err := pdfReader.GetPage(1)
	if err != nil {
		t.Error(err)
	}

	textp1, err := ReadPageData(page)

	if err != nil {
		t.Error(err)
	}

	page2, err := pdfReader.GetPage(2)
	if err != nil {
		t.Error(err)
	}

	textp2, err := ReadPageData(page2)
	if err != nil {
		t.Error(err)
	}

	assertEqual(t, text1, textp1[0])
	assertEqual(t, text2, textp2[0])

}

func TestWriteReadOptimiser(t *testing.T) {

	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	c.SetPageSize(creator.PageSizeA4)
	c.NewPage()
	text1 := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":1}"
	WritePageData(c, text1)
	text2 := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":2}"
	c.NewPage()
	WritePageData(c, text2)

	// write to memory instead of a file
	var buf bytes.Buffer

	c.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    90,
		ImageUpperPPI:                   150,
	}))

	err := c.Write(&buf)
	if err != nil {
		t.Error(err)
	}

	// convert buffer to readseeker
	var bufslice []byte
	fbuf := filebuffer.New(bufslice)
	fbuf.Write(buf.Bytes())

	// read in from memory
	pdfReader, err := model.NewPdfReader(fbuf)
	if err != nil {
		t.Error(err)
	}
	page, err := pdfReader.GetPage(1)
	if err != nil {
		t.Error(err)
	}

	textp1, err := ReadPageData(page)

	if err != nil {
		t.Error(err)
	}

	page2, err := pdfReader.GetPage(2)
	if err != nil {
		t.Error(err)
	}

	textp2, err := ReadPageData(page2)
	if err != nil {
		t.Error(err)
	}

	assertEqual(t, text1, textp1[0])
	assertEqual(t, text2, textp2[0])

}
