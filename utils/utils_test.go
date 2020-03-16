package utils

import (
	"net/http"
	"os"
	"testing"
)

var folder = "images/test/folder"

func TestMakeFolder(t *testing.T) {
	MakeFolder(folder)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		t.Errorf("Folder \"%s\" does not exists\n", folder)
	}
}

func TestSaveImage(t *testing.T) {
	fileLink := "https://i.morioh.com/36be53336e.png"
	fileName := "36be53336e.png"

	response, err := http.Get(fileLink)
	if err != nil {
		t.Errorf("Error when download remote image - %v\n", err)
	}

	size, err := SaveImage(folder + "/" + fileName, response)
	if err != nil {
		t.Errorf("Cant save image to folder - %v\n", err)
	}

	if size <= 0 {
		t.Errorf("Error to copy file-%s, in folder-%s", fileLink, folder)
	}

	os.RemoveAll("images")
}
