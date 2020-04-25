package pdfpagedata

import (
	"bytes"
	"testing"

	"github.com/mattetti/filebuffer"
	"github.com/timdrysdale/unipdf/v3/creator"
	"github.com/timdrysdale/unipdf/v3/model"
	"github.com/timdrysdale/unipdf/v3/model/optimize"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func TestWriteReadDouble(t *testing.T) {

}

func TestWriteReadOtherText(t *testing.T) {}

func TestWriteRead(t *testing.T) {

	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	c.SetPageSize(creator.PageSizeA4)
	c.NewPage()
	text1 := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":1}"
	WritePageString(c, text1)
	text2 := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":2}"
	c.NewPage()
	WritePageString(c, text2)

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

	textp1, err := ReadPageString(page)

	if err != nil {
		t.Error(err)
	}

	page2, err := pdfReader.GetPage(2)
	if err != nil {
		t.Error(err)
	}

	textp2, err := ReadPageString(page2)
	if err != nil {
		t.Error(err)
	}

	assertEqual(t, text1, textp1)
	assertEqual(t, text2, textp2)

}

func TestWriteReadOptimiser(t *testing.T) {

	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	c.SetPageSize(creator.PageSizeA4)
	c.NewPage()
	text1 := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":1}"
	WritePageString(c, text1)
	text2 := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":2}"
	c.NewPage()
	WritePageString(c, text2)

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

	textp1, err := ReadPageString(page)

	if err != nil {
		t.Error(err)
	}

	page2, err := pdfReader.GetPage(2)
	if err != nil {
		t.Error(err)
	}

	textp2, err := ReadPageString(page2)
	if err != nil {
		t.Error(err)
	}

	assertEqual(t, text1, textp1)
	assertEqual(t, text2, textp2)

}
