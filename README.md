# Google image downloader util

This util is helper for download many images from custom google search service and put them to specified folder

[![License MIT](https://img.shields.io/apm/l/vim-mode.svg)](https://en.wikipedia.org/wiki/MIT_License)


### Main usage
You can use this library in your project
```go
package main

import "github.com/killer-djon/gimdownloader/utils"

var (
    url = "https://www.googleapis.com"
    path = "/customsearch/v1"
    tag = "some_tag_name" // required param
    folder = "required_folder_to_save (default: ./images)"
)
// When tag is prefix for file name
// tag - required param
var request = utils.NewRequest(url, path, tag)

request.AddQuery("key", key)
request.AddQuery("cx", cx)
request.AddQuery("q", query)
request.AddQuery("num", strconv.Itoa(num))
request.AddQuery("searchType", "image")
request.AddQuery("imgSize", imgSize)
request.AddQuery("imgColorType", imgColorType)
request.AddQuery("fileType", imgType)
request.AddQuery("start", "1")

request.DownloadImages(folder)
```