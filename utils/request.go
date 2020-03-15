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

const (
	MAX_PER_PAGE = 10
	MAX_PAGES    = 10
)

type Request struct {
	Client  *http.Request
	Query   url.Values
	Uri     string
	Path    string
	TagName string
}

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

func (req *Request) AddQuery(key, value string) {
	req.Query = req.Client.URL.Query()
	req.Query.Set(key, value)
	req.Client.URL.RawQuery = req.Query.Encode()
}

func (req *Request) DelQuery(key string) {
	req.Query = req.Client.URL.Query()
	req.Query.Del(key)
	req.Client.URL.RawQuery = req.Query.Encode()
}

func (req Request) DownloadImages(folder string) {
	var wg sync.WaitGroup

	total, _ := strconv.Atoi(req.Client.URL.Query().Get("num"))

	if total > MAX_PAGES*MAX_PER_PAGE {
		log.Fatalf("Number items for request result is greater then %d\n", MAX_PAGES*MAX_PER_PAGE)
		return
	}

	MakeFolder(folder)

	totalPage := math.Ceil(float64(total) / float64(MAX_PER_PAGE))
	log.Println("Total pages", total, totalPage)

	remain := total - (total/MAX_PER_PAGE)*MAX_PER_PAGE

	for i := 0; i < int(totalPage); i++ {
		wg.Add(1)
		go func(i int) {
			if (i + 1) == int(totalPage) {
				if remain > 0 {
					req.AddQuery("num", strconv.Itoa(remain))
				}else{
					req.AddQuery("num", strconv.Itoa(MAX_PER_PAGE))
				}

			} else {
				req.AddQuery("num", strconv.Itoa(MAX_PER_PAGE))
			}
			req.AddQuery("start", strconv.Itoa(i*MAX_PER_PAGE+1))

			req.getPage(i+1, folder, req.Client.URL.Query())
			wg.Done()
		}(i)
	}

	wg.Wait()
}

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
