package main

import (
	"flag"
	"pogo/common/logs"
	"pogo/parser"
)

var log = logs.Log

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
	log.Debug("go to download video")

}
