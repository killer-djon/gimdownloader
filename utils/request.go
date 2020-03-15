package utils

import (
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
	Client *http.Request
	Query  url.Values
}

func NewRequest(uri, path string) *Request {
	req, err := http.NewRequest("GET", uri+path, nil)
	if err != nil {
		return nil
	}
	return &Request{
		Client: req,
		Query:  nil,
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

	totalPage := math.Ceil(float64(total) / float64(MAX_PER_PAGE))
	log.Println("Total pages", total, totalPage)

	remain := total - (total/MAX_PER_PAGE) * MAX_PER_PAGE
	log.Println("Remain num", remain)
	for i := 0; i < int(totalPage); i++ {
		wg.Add(1)
		go func(i int) {
			req.AddQuery("start", strconv.Itoa(i*MAX_PER_PAGE+1))
			log.Println("First request", i, req.Client.URL.Query())

			wg.Done()
		}(i)
	}

	wg.Wait()

}
