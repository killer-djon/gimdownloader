package utils

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

// Defaults constant variables
// for pagination when download images
const (
	MAX_PER_PAGE = 10
	MAX_PAGES    = 10
)

// Request Main request struct
// which have exported methods
type Request struct {
	Client  *http.Request
	Query   url.Values
	Uri     string
	Path    string
	TagName string
}

// NewRequest Create blank struct Request
// for generate initial query params
func NewRequest(uri, path, tag string) *Request {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		return nil
	}
	return &Request{
		Client:  req,
		Query:   nil,
		Uri:     uri,
		Path:    path,
		TagName: tag,
	}
}

// AddQuery Add/Set query item to slice
func (req *Request) AddQuery(key, value string) {
	req.Query = req.Client.URL.Query()
	req.Query.Set(key, value)
	req.Client.URL.RawQuery = req.Query.Encode()
}

// DelQuery Remove query item from slice
func (req *Request) DelQuery(key string) {
	req.Query = req.Client.URL.Query()
	req.Query.Del(key)
	req.Client.URL.RawQuery = req.Query.Encode()
}

// GetQueries Get all query values
func (req Request) GetQueries() url.Values {
	return req.Query
}

// GetQuery Get query item by key from slice queries
func (req Request) GetQuery(key string) string {
	return req.Query.Get(key)
}

// DownloadImages Start download all images by config params
func (req Request) DownloadImages(folder string) {
	total, _ := strconv.Atoi(req.Client.URL.Query().Get("num"))

	if total > MAX_PAGES*MAX_PER_PAGE {
		log.Fatalf("Number items for request result is greater then %d\n", MAX_PAGES*MAX_PER_PAGE)
		return
	}

	MakeFolder(folder)

	totalPage := math.Ceil(float64(total) / float64(MAX_PER_PAGE))
	log.Println("Total pages", total, totalPage)

	remain := total - (total/MAX_PER_PAGE)*MAX_PER_PAGE

	var wg sync.WaitGroup
	for i := 0; i < int(totalPage); i++ {
		wg.Add(1)
		go func(i int) {
			count := MAX_PER_PAGE
			if (i + 1) == int(totalPage) {
				if remain > 0 {
					count = remain
				}
			}

			req.AddQuery("num", strconv.Itoa(count))
			req.AddQuery("start", strconv.Itoa(i*MAX_PER_PAGE+1))

			req.getPage(i+1, folder, req.Client.URL.Query())
			wg.Done()
		}(i)
	}

	wg.Wait()
}

// getPage Private method for get page with images
func (req *Request) getPage(grNum int, folder string, values url.Values) {

	uri, _ := url.ParseRequestURI(req.Uri + req.Path)
	uri.RawQuery = values.Encode()

	resp, err := http.Get(uri.String())
	if err != nil {
		log.Fatalf("Can't get page uri: %v", err)
	}
	log.Println("Request uri", uri.String())
	response, _ := NewResponse(resp)

	if len(response.GetImages()) > 0 {
		i := 0
		for index, itemImg := range response.GetImages() {
			responseFile, err := http.Get(itemImg.Link)
			if err != nil {
				log.Printf("Error occurred when download file image: %v\n", err)
				continue
			}

			if responseFile.StatusCode != 200 {
				fmt.Printf("File not found by this link: %s\n", itemImg.Link)
				continue
			}

			fileName := fmt.Sprintf("%s_000%d_00%d.%s", req.TagName, grNum, index, values.Get("fileType"))
			filePath := folder + "/" + fileName
			size, err := SaveImage(filePath, responseFile)
			if err != nil {
				log.Printf("Error for download this file: %s, %v\n", itemImg.Link, err)
				continue
			}
			i++
			log.Printf("File saved as name: %s, with size: %dKb\n", filePath, size)
		}
		log.Printf("Count saved images: %d\n", i)
		log.Printf("Count unsaved images by some error: %d\n", len(response.GetImages())-i)

	}
}
