package main

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"telex-chat/internal/data"
	"telex-chat/internal/env"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// get file content from file
func GetContentFromFile(fileCon []byte, ext string) (string, error) {
	switch ext {
	case ".pdf":
		return data.ExtractTextFromPDF(fileCon)
	case ".docx":
		return data.ExtractTextFromDocx(fileCon)
	case ".txt":
		return data.ExtractTextFromTxt(fileCon)
	default:
		return "", fmt.Errorf("file ext not supported")
	}

}

// save the file to either cloudinary or disk
func (app *application) SaveFiletoCloudOrDisk(envi, ext string, header *multipart.FileHeader, file multipart.File) (string, error) {
	filePath := fmt.Sprintf("uploads/%s_%s.%s", strings.Split(header.Filename, ".")[0], env.GetID(), ext)
	// save file to local disk for local development sorry m coding this way
	if envi == "DEV" {
		outFile, err := os.Create(filePath)
		if err != nil {
			return "", err
		}
		defer outFile.Close()
		_, err = io.Copy(outFile, file)

		if err != nil {
			return "", err
		}
	} else {

		// save file to cloudinary
		ctx := context.Background()

		upload, err := app.cld.Upload.Upload(ctx, file, uploader.UploadParams{
			ResourceType: "auto",
		})

		if err != nil {
			return "", err
		}

		fmt.Println(upload.PublicID)
		filePath = upload.SecureURL

	}
	return filePath, nil
}
