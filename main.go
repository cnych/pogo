package main

import (
	"flag"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"net/http"
	"pogo/common/logs"
	"pogo/common/nets/download"
	"pogo/parser"
)

var (
	log = logs.Log
)

func init() {
	flag.Parse()
}

func main() {
	if flag.NArg() == 0 {
		log.Error("no video url parameter!")
		return
	}
	url := flag.Arg(0)

	_, spider := parser.GetParser(url)
	videoInfo, err := spider.GetVideoInfo()
	if err != nil {
		log.Error(err.Error())
		return
	}

	fmt.Printf("Title:			%s\n", videoInfo.Title)
	fmt.Printf("Site:			%s\n", videoInfo.Site)
	fmt.Printf("Url:			%s\n", videoInfo.Url)
	fmt.Printf("Duration:		%ds\n", videoInfo.Duration)

	for _, hd := range []string {"hd2", "hd1", "normal"} {
		if dl, ok := videoInfo.DownloadInfo[hd]; ok {
			filename := fmt.Sprintf("%s_%s_%s.mp4", videoInfo.Site, hd, videoInfo.Title)
			header := http.Header{}
			header.Add("user-agent", browser.Computer())
			err := download.Download(filename, dl.(string), header)
			if err != nil {
				log.Error("download %s video error: %s", hd, err.Error())
			} else {
				break
			}
		}
	}

}
