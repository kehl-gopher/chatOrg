package data

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

func ExtractTextFromPDF(file string) (string, error) {

	f, err := os.Open(file)

	if err != nil {
		return "", err
	}

	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)

	if err != nil {
		return "", err
	}

	numPages, err := pdfReader.GetNumPages()

	if err != nil {
		return "", err
	}

	var textBuilder strings.Builder

	for i := 0; i < numPages; i++ {
		pageNum := i + 1
		page, err := pdfReader.GetPage(pageNum)

		if err != nil {
			return "", err
		}

		textExtractor, err := extractor.New(page)

		if err != nil {
			return "", err
		}

		text, err := textExtractor.ExtractText()

		if err != nil {
			return "", err
		}
		textBuilder.WriteString(text + "\n")

	}

	fmt.Println(textBuilder.String(), "textBuilder")
	return textBuilder.String(), nil
}

func ExtractTextFromDocx(file string) (string, error) {
	doc, err := document.Open(file)

	if err != nil {
		return "", err
	}

	defer doc.Close()

	var textBuilder strings.Builder

	for _, para := range doc.Paragraphs() {

		for _, run := range para.Runs() {
			textBuilder.WriteString(run.Text() + "\n")
		}
	}

	fmt.Println(textBuilder.String(), "doc builder")
	return textBuilder.String(), nil
}

func ExtractTextFromTxt(file string) (string, error) {
	f, err := os.Open(file)

	if err != nil {
		return "", err
	}

	defer f.Close()
	buf := make([]byte, 1024)

	var textBuilder strings.Builder
	for {
		n, err := f.Read(buf)

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return "", err
		}

		textBuilder.WriteString(string(buf[:n]))
		textBuilder.WriteString("\n")
	}
	return textBuilder.String(), nil

}
