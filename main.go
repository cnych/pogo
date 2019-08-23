package main

import (
	"flag"
	"fmt"
	"net/http"
	"pogo/common/logs"
	"pogo/common/nets/download"
	"pogo/parser"
)

var (
	log = logs.Log
	ua = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36"
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
	log.Debug(url)

	xg := parser.Xigua{Url: url}
	videoInfo, err := xg.GetVideoInfo()
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Debug("Title: %s", videoInfo.Title)
	log.Debug("Url: %s", videoInfo.Url)
	log.Debug("Site: %s", videoInfo.Site)
	log.Debug("Duration: %d", videoInfo.Duration)
	log.Debug("DownloadInfo: %v", videoInfo.DownloadInfo)

	for _, hd := range []string{"hd3", "hd2", "hd1", "normal"} {
		if dl, ok := videoInfo.DownloadInfo[hd]; ok {
			filename := fmt.Sprintf("%s_%s_%s.mp4", videoInfo.Site, hd, videoInfo.Title)
			header := http.Header{}
			header.Add("user-agent", ua)
			result, err := download.Download(filename, dl.(string), header)
			if err != nil {
				log.Error("download %s video error: %s", hd, err.Error())
			} else {
				log.Debug("download %s video completed, location %s", hd, result)
				break
			}
		}
	}

}
