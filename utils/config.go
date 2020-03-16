package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Root config params
type Config struct {
	Key         string      `json:"api_key"`
	Cx          string      `json:"api_cx"`
	Url         string      `json:"api_url"`
	EndPoint    string      `json:"api_end_point"`
	QueryConfig QueryConfig `json:"query_config"`
}

// Config params for query strings
type QueryConfig struct {
	FileType     string `json:"file_type"`
	ImgColorType string `json:"img_color_type"`
	ImgSize      string `json:"img_size"`
	Num          int    `json:"num"`
	SearchType   string `json:"search_type"`
}

// get config file
func GetConfig(path string) *Config {
	jsonFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Printf("Error when try to open json config file %v", err)
		return nil
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config = &Config{}

	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		log.Printf("Cant read json root data %v", err)
		return nil
	}

	return config
}
