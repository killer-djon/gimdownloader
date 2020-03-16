package main

import (
	"github.com/killer-djon/gimdownloader/utils"
	"testing"
)

func TestGetConfig(t *testing.T) {
	config := utils.GetConfig("config.json")
	if config == nil {
		t.Error("Config file must be exists in folder")
	}
}

func TestSetRequestParams(t *testing.T) {
	var request = utils.NewRequest("https://www.google.com", "/customsearch/v1", "test")

	request.AddQuery("key", "key")
	request.AddQuery("cx", "cx")
	request.AddQuery("q", "query")
	request.AddQuery("num", "10")
	request.AddQuery("searchType", "image")
	request.AddQuery("imgSize", "large")
	request.AddQuery("imgColorType", "color")
	request.AddQuery("fileType", "jpeg")
	request.AddQuery("start", "1")

	queries := request.GetQueries()
	// 9 is count of all needed params
	if len(queries) != 9 {
		t.Error("Count of required params is not equal")
	}
}
