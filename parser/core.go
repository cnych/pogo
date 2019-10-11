package parser

import "pogo/common/strs"

type VideoInfo struct {
	Site string `json:"site"`
	Title string `json:"title"`
	Duration int64 `json:"duration"`
	Url string `json:"url"`
	DownloadInfo map[string]interface{} `json:"download_info"`
}

type BaseVideo struct {
	Url string
	Name string
	VideoPatterns []string
}

func (base *BaseVideo) MatchUrl(url string) bool {
	//	https://www.ixigua.com/i6688067268607214088/
	for _, pattern := range base.VideoPatterns {
		if rcontent := strs.MatchRegexpOf1(pattern, url); len(rcontent) > 0 {
			base.Url = url
			return true
		}
	}
	return false
}

type Spider interface {
	MatchUrl(url string) bool
	GetVideoInfo() (info *VideoInfo, err error)
}

var Spiders = make(map[string]Spider)

func init()  {
	XiguaRegister()
	DouyinRegister()
}

func GetParser(url string) (key string, spider Spider) {
	//url "ixigua.com/abc"
	for a, b := range Spiders {
		if b.MatchUrl(url) {
			key = a
			spider = b
			return
		}
	}
	return
}