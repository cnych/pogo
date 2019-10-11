package parser

import (
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/bitly/go-simplejson"
	"github.com/cnych/stardust/encodingx/base64x"
	"pogo/common/call"
	"pogo/common/logs"
	"pogo/common/nets/fetch"
	"pogo/common/strs"
)

var log = logs.Log

const (
	aid = 1768
)

type Xigua struct {
	BaseVideo
	vid string
	businessToken string
	authToken string
}

func XiguaRegister()  {
	xg := new(Xigua)
	xg.Name = "xigua"
	xg.VideoPatterns = []string{`ixigua\.com\/i(\d+)`}
	Spiders[xg.Name] = xg
}

func (xg *Xigua) GetVideoInfo() (info *VideoInfo, err error) {
	_ = call.Retry(3, func() error {
		info, err = getXiguaOnce(xg)
		return err
	})
	return
}

func getXiguaOnce(xg *Xigua) (*VideoInfo, error) {
	header := map[string]string{
		"user-agent": browser.Computer(),
		"referer": "https://www.ixigua.com/",
	}
	req := fetch.DefaultRequest(xg.Url, header)

	resp, err := fetch.Fetch(req)
	if err != nil {
		return nil, err
	}

	html, err := resp.AsText("UTF-8")
	if err != nil {
		return nil, err
	}

	//"title":"乔恩为了让加菲猫干活，竟然让博士对他催眠，这就有点狠了！","tag"
	title := strs.MatchRegexpOf1(`"user_bury":0,"title":"(.*)","tag"`, html)

	videoInfo := VideoInfo{}
	videoInfo.Title = title
	videoInfo.Url = xg.Url
	videoInfo.Site = "西瓜视频"

	//"vid":"v02004910000bj377nc1n3e63t2qmc00",
	xg.vid = strs.MatchRegexpOf1(`"vid":"(.*)","user_digg"`, html)

	//"businessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjU0MjI5MTMsInZlciI6InYxIiwiYWsiOiJjZmMwNjdiYjM5ZmVmZjU5MmFmODIwODViNDJlNmRjMyIsInN1YiI6InBnY19ub3JtYWwifQ.aaUsIAdV5yqVRZtv4A9G9ijV_GGP261ww-2gK2Asyt0","authToken"
	xg.businessToken = strs.MatchRegexpOf1(`"businessToken":"(.*)","authToken"`, html)

	//"authToken":"HMAC-SHA1:2.0:1565422913273604497:cfc067bb39feff592af82085b42e6dc3:RuTnxlMBKOtk8z4p0J\u002F3aWuc27o=","is_original"
	xg.authToken = strs.MatchRegexpOf1(`"authToken":"(.*)","is_original"`, html)

	videoJson, err := parseVideoUrl(xg)
	if err != nil {
		return nil, err
	}

	duration := videoJson.Get("data").Get("video_duration").MustFloat64()
	videoInfo.Duration = int64(duration)

	downloadInfo := make(map[string]interface{})
	videoListJson := videoJson.Get("data").Get("video_list")
	if vjson, exists := videoListJson.CheckGet("video_1"); exists {
		videoUrl, err := getVideoUrl(vjson)
		if err != nil {
			return nil, err
		}
		downloadInfo["normal"] = videoUrl
	}
	if vjson, exists := videoListJson.CheckGet("video_2"); exists {
		videoUrl, err := getVideoUrl(vjson)
		if err != nil {
			return nil, err
		}
		downloadInfo["hd1"] = videoUrl
	}
	if vjson, exists := videoListJson.CheckGet("video_3"); exists {
		videoUrl, err := getVideoUrl(vjson)
		if err != nil {
			return nil, err
		}
		downloadInfo["hd2"] = videoUrl
	}

	videoInfo.DownloadInfo = downloadInfo

	return &videoInfo, nil
}

func getVideoUrl(vjson *simplejson.Json) (string, error) {
	mainUrl := vjson.Get("main_url").MustString()
	resultBytes, err := base64x.DecodeString(mainUrl)
	if err != nil {
		mainUrl = vjson.Get("backup_url_1").MustString()
		resultBytes, err = base64x.DecodeString(mainUrl)
		if err != nil {
			return "", err
		}
	}
	return string(resultBytes), nil
}

func parseVideoUrl(xg *Xigua) (*simplejson.Json, error) {
	//https://vas.snssdk.com/video/openapi/v1/?aid=1768&action=GetPlayInfo&video_id=v03004fa0000bl6vmnljcra2v8631obg&nobase64=false&ptoken=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjU0MjQ1NzgsInZlciI6InYxIiwiYWsiOiJjZmMwNjdiYjM5ZmVmZjU5MmFmODIwODViNDJlNmRjMyIsInN1YiI6InBnY19ub3JtYWwifQ.EBAM3YOFY0SlGu4bieHAM0iI8nzJXAvCsHp5qvvn6Lg&vfrom=xgplayer
	apiUrl := fmt.Sprintf("https://vas.snssdk.com/video/openapi/v1/?aid=%d&" +
		"action=GetPlayInfo&video_id=%s" +
		"&nobase64=false&ptoken=%s&" +
		"vfrom=xgplayer", aid, xg.vid, xg.businessToken)

	header := map[string]string{
		"user-agent": browser.Computer(),
		"referer": xg.Url,
		"Authorization": xg.authToken,
	}
	req := fetch.DefaultRequest(apiUrl, header)

	resp, err := fetch.Fetch(req)
	if err != nil {
		return nil, err
	}

	resultJSON, err := resp.ParseSimpleJSON()

	//log.Debug("%v", resultJSON)

	code := resultJSON.Get("code").MustInt()
	if code != 0 {
		return nil, fmt.Errorf(resultJSON.Get("message").MustString())
	}

	return resultJSON, err
}
