package data

// get file content from file
func GetContentFromFile(filepath, ext string) (string, error) {

	switch ext {
	case ".pdf":
		return ExtractTextFromPDF(filepath)
	case ".docx":
		return ExtractTextFromDocx(filepath)
	case ".txt":

	}

	return "", nil
}
