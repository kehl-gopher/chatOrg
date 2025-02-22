package data

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"

	"github.com/nguyenthenguyen/docx"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

func ExtractTextFromPDF(fileBlob []byte) (string, error) {

	// f, err := os.Open(file)

	// if err != nil {
	// 	return "", err
	// }

	// defer f.Close()

	f := bytes.NewReader(fileBlob)
	pdfReader, err := model.NewPdfReader(f)

	if err != nil {
		return "", err
	}

	numPages, err := pdfReader.GetNumPages()

	if err != nil {
		return "", err
	}

	var textBuilder strings.Builder

	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)

		if err != nil {
			log.Printf("Warning Failed to get page %d\n", i)
			continue
		}

		textExtractor, err := extractor.New(page)

		if err != nil {
			log.Printf("Warning failed to extract page %d\n", i)
			continue
		}

		text, err := textExtractor.ExtractText()

		if err != nil {
			log.Printf("Failed to extract text from page %d\n", i)
			continue
		}

		if strings.TrimSpace(text) != "" {
			textBuilder.WriteString(text + "\n")
		}
	}

	fmt.Println(textBuilder.String(), "textBuilder")
	return textBuilder.String(), nil
}

func ExtractTextFromDocx(fileBlob []byte) (string, error) {
	reader := bytes.NewReader(fileBlob)

	doc, err := docx.ReadDocxFromMemory(reader, int64(len(fileBlob)))
	if err != nil {
		return "", err
	}
	defer doc.Close()

	// Get editable content
	content := doc.Editable()

	// Get all text content
	text := content.GetContent()

	// Clean up the text (remove extra whitespace)
	re := regexp.MustCompile("<[^>]*>")
	text = re.ReplaceAllString(text, "")

	// Clean up the text
	text = strings.TrimSpace(text)                  // Remove leading/trailing whitespace
	text = strings.ReplaceAll(text, "\n\n\n", "\n") // Remove extra newlines
	text = strings.ReplaceAll(text, "  ", " ")      // Remove double spaces

	return text, nil

}

// func ExtractTextFromDocx(fileBlob []byte) (string, error) {
// 	// doc, err := document.Open(file)

// 	// if err != nil {
// 	// 	return "", err
// 	// }

// 	// defer doc.Close()

// 	f := bytes.NewReader(fileBlob)

// 	doc, err := document.Read(f, int64(len(fileBlob)))

// 	if err != nil {
// 		return "", err
// 	}
// 	var textBuilder strings.Builder

// 	for _, para := range doc.Paragraphs() {

// 		for _, run := range para.Runs() {
// 			textBuilder.WriteString(run.Text() + "\n")
// 		}
// 	}
// 	return textBuilder.String(), nil
// }

func ExtractTextFromTxt(fileBlob []byte) (string, error) {
	// f, err := os.Open(file)

	// if err != nil {
	// 	return "", err
	// }

	// defer f.Close()

	nr := bytes.NewReader(fileBlob)
	buf := make([]byte, 1024)

	var textBuilder strings.Builder
	for {
		n, err := nr.Read(buf)

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return "", err
		}

		textBuilder.WriteString(string(buf[:n]))
	}
	return textBuilder.String(), nil

}
