package utils

import (
	"log"
	"net/http"
	"strconv"
	"sync"
)

const (
	MAX_PER_PAGE = 10
	MAX_PAGES    = 10
)

type Request struct {
	Client *http.Request
}

func NewRequest(link, path string) *Request {
	req, err := http.NewRequest("GET", link+path, nil)

	if err != nil {
		log.Fatalf("Error for creating request: %v\n", err)
		return nil
	}

	return &Request{
		Client: req,
	}
}

func (req Request) AddQuery(key, val string) {
	query := req.Client.URL.Query()
	query.Set(key, val)

	req.Client.URL.RawQuery = query.Encode()
}

func (req Request) HasQuery(key string) bool {
	query := req.Client.URL.Query()
	if query.Get(key) != "" {
		return true
	}

	return false
}

func (req Request) GetQuery(key string) string {
	query := req.Client.URL.Query()
	return query.Get(key)
}

func (req Request) AddHeader(key, val string) {
	head := req.Client.Header
	head.Set(key, val)
}

func (req Request) HasHeader(key string) bool {
	head := req.Client.Header
	if head.Get(key) != "" {
		return true
	}
	return false
}

func (req Request) GetHeader(key string) string {
	head := req.Client.Header
	return head.Get(key)
}

func (req Request) GetRequest() string {
	return req.Client.URL.String()
}

func (req *Request) DownloadImages(parentFolder string) {
	num, _ := strconv.Atoi(req.GetQuery("num"))
	if num > MAX_PAGES*MAX_PER_PAGE {
		log.Fatalf("Number items for request result is greater then %d\n", MAX_PAGES*MAX_PER_PAGE)
		return
	}

	MakeFolder(parentFolder)

	var wg sync.WaitGroup
	wg.Add(num/MAX_PER_PAGE)

	for i := 0; i < num/MAX_PER_PAGE; i++ {
		go func(req *Request, i int) {
			defer wg.Done()

			nextNum := strconv.Itoa(i)
			req.AddQuery("start", nextNum)
			log.Println("Page", i, req.GetRequest())
		}(req, i)
	}

	wg.Wait()
}
