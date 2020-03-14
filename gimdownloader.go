package main

import (
	"flag"
	"imagedownload/gimdownloader/utils"
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
	searchType   string
	url          string
	path         string
)

func init() {
	flag.StringVar(&folder, "folder", IMAGES_FOLDER, "Folder where images will be download")
	flag.StringVar(&configFile, "configFile", "", "If is set then get config params, otherwise get by args")
	flag.StringVar(&key, "key", "", "Key API from google console")
	flag.StringVar(&cx, "cx", "", "Key for Custom search API")
	flag.IntVar(&num, "num", 0, "How match images went to get (default: 10)")
	flag.StringVar(&imgSize, "imgSize", "", "Image size for download, like medium,large,small ...")
	flag.StringVar(&imgColorType, "imgColorType", IMAGE_COLOR_TYPE, "Image color type like (color, gray, mono)")
	flag.StringVar(&query, "query", "", "Query string for search images by this query")
	flag.StringVar(&imgType, "imgType", IMAGE_TYPE, "Image type for download (default: jpeg)")

}

func main() {
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

	if num > 0 {
		count = num
	} else {
		if count > 0 {
			count = count
		}
	}

	if key == "" {
		log.Println("Your key flag is empty you must specified them")
		return
	}

	if query == "" {
		log.Println("You must set query string for search images")
		return
	}
	var request = utils.NewRequest(url, path)

	request.AddQuery("key", key)
	request.AddQuery("cx", cx)
	request.AddQuery("q", query)
	request.AddQuery("num", strconv.Itoa(count))
	request.AddQuery("searchType", "image")
	request.AddQuery("imgSize", imgSize)
	request.AddQuery("imgColorType", imgColorType)
	request.AddQuery("fileType", imgType)
	request.AddQuery("start", "1")
	
	log.Println(request.Client.URL.Query().Encode())
	request.DownloadImages("./images")
}

