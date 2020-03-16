package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Response Root struct for response
// by download
type Response struct {
	StatusCode int
	Status     string
	Body       *Images
}

// Images struct with slice
// of downloaded images list
type Images struct {
	Error *Error
	Items []*Image `json:"items"`
}

// Error struct
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Image Simple image struct
type Image struct {
	Title      string     `json:"title"`
	Link       string     `json:"link"`
	Name       string     `json:"displayLink"`
	FileFormat string     `json:"fileFormat"`
	Meta       *ImageMeta `json:"image"`
}

// ImageMeta Image metadata for all
// images in list
type ImageMeta struct {
	Height int `json:"height"`
	Width  int `json:"width"`
	Size   int `json:"byteSize"`
}

// NewResponse Create new struct Response
func NewResponse(resp *http.Response) (*Response, error) {
	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("Error to get response bytes: %v", err)
	}

	var images = &Images{}
	err = json.Unmarshal(bodyByte, images)
	if err != nil {
		log.Printf("Error when decode json response images array: %v", err)
		return nil, err
	}

	if images.Error != nil {
		log.Printf("Error occurred when download files: %d - %s\n", images.Error.Code, images.Error.Message)
		return nil, err
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       images,
	}, nil
}

// GetImages Helper methods for get items for response
// Get images list
func (response Response) GetImages() []*Image {
	return response.Body.Items
}

// GetStatusCode Get status code
func (response Response) GetStatusCode() int {
	return response.StatusCode
}

// GetStatus Get status text
func (response Response) GetStatus() string {
	return response.Status
}
