package main

import (
	"flag"
	"fmt"
	"github.com/killer-djon/gimdownloader/utils"
	"log"
	"strconv"
)

const (
	IMAGES_FOLDER    = "images"
	IMAGE_TYPE       = "jpeg"
	IMAGE_COLOR_TYPE = "color"
	DEFUALT_URL      = "https://www.googleapis.com"
	DEFAULT_PATH     = "/customsearch/v1"
)

var (
	folder       string
	configFile   string
	key          string
	cx           string
	num, count   int
	imgSize      string
	imgColorType string
	query        string
	imgType      string
	url          string
	path         string
	tag          string
)

func init() {
	flag.StringVar(&folder, "folder", IMAGES_FOLDER, "Folder where images will be download")
	flag.StringVar(&configFile, "configFile", "", "If is set then get config params, otherwise get by args")
	flag.StringVar(&key, "key", "", "Key API from google console")
	flag.StringVar(&cx, "cx", "", "Key for Custom search API")
	flag.IntVar(&num, "num", 10, "How match images went to get")
	flag.StringVar(&imgSize, "imgSize", "", "Image size for download, like medium,large,small ...")
	flag.StringVar(&imgColorType, "imgColorType", IMAGE_COLOR_TYPE, "Image color type like (color, gray, mono)")
	flag.StringVar(&query, "query", "", "Query string for search images by this query")
	flag.StringVar(&imgType, "imgType", IMAGE_TYPE, "Image type for download")
	flag.StringVar(&tag, "tag", "", "Tag name for named image for download")
}

func main() {

	flag.Usage = func() {
		fmt.Printf("Google image downloader by Leshanu Evgeniy\n")
		fmt.Println("Usage:")
		fmt.Printf("	gimdownloader [options] \n")
		fmt.Println("Options:")
		flag.PrintDefaults()

	}

	config := utils.GetConfig("./config.json")
	flag.Parse()
	url = DEFUALT_URL
	path = DEFAULT_PATH
	if config != nil {
		key = config.Key
		cx = config.Cx
		count = config.QueryConfig.Num
		imgSize = config.QueryConfig.ImgSize
		imgColorType = config.QueryConfig.ImgColorType
		imgType = config.QueryConfig.FileType
		url = config.Url
		path = config.EndPoint
	}

	if tag == "" {
		log.Println("You must type tag name for images")
		flag.Usage()
		return
	}

	if num > 0 {
		count = num
	} else {
		if count > 0 {
			count = count
		}
	}

	if key == "" {
		log.Println("Your key flag is empty you must specified them")
		flag.Usage()
		return
	}

	if query == "" {
		log.Println("You must set query string for search images")
		flag.Usage()
		return
	}
	var request = utils.NewRequest(url, path, tag)

	request.AddQuery("key", key)
	request.AddQuery("cx", cx)
	request.AddQuery("q", query)
	request.AddQuery("num", strconv.Itoa(count))
	request.AddQuery("searchType", "image")
	request.AddQuery("imgSize", imgSize)
	request.AddQuery("imgColorType", imgColorType)
	request.AddQuery("fileType", imgType)
	request.AddQuery("start", "1")

	//log.Println(request.Client.URL.Query().Encode())
	request.DownloadImages(folder)
}
