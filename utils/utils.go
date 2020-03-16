package utils

import (
	"io"
	"log"
	"net/http"
	"os"
)

// MakeFolder create typed folder
func MakeFolder(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err := os.MkdirAll(folder, 0775)
		if err != nil {
			log.Panicf("Error when try to create folder: %v", err)
		}
	}
}

// SaveImage try to save in folder downloaded image
func SaveImage(filepath string, resp *http.Response) (int64, error) {
	file, err := os.Create(filepath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Write the body to file
	size, err := io.Copy(file, resp.Body)
	return size, nil
}
