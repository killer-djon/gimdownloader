package utils

import (
	"log"
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
	Client *http.Request
	Query url.Values
}

func NewRequest(uri, path string) *Request {
	req, err := http.NewRequest("GET", uri+path, nil)
	if err != nil {
		return nil
	}
	return &Request{
		Client: req,
	}
}

func (req *Request) AddQuery(key, value string) {
	req.Query = req.Client.URL.Query()
	req.Query.Set(key, value)
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

	for i := 0; i < total/MAX_PER_PAGE; i++ {
		wg.Add(1)
		go func(i int) {
			req.AddQuery("start", strconv.Itoa(i))
			req.AddQuery("num", strconv.Itoa(MAX_PER_PAGE))

			log.Println("First request", i, req.Client.URL.Query())

			wg.Done()
		}(i)
	}

	wg.Wait()

}