package parser

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/cnych/stardust/encodingx/base64x"
	"io/ioutil"
	"net/http"
	"pogo/common/logs"
	"regexp"
)

var log = logs.Log

const (
	aid = 1768
	ua = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36"
)

type Xigua struct {
	Url string
	vid string
	businessToken string
	authToken string
}


func (xg *Xigua) GetVideoInfo() (*VideoInfo, error) {
	req, err := http.NewRequest(http.MethodGet, xg.Url, nil)
	if err != nil {
		return nil, err
	}

	header := http.Header{}
	header.Add("user-agent", ua)
	header.Add("referer", "https://www.ixigua.com/")
	req.Header = header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resultBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//log.Debug(string(resultBytes))
	html := string(resultBytes)

	//"title":"乔恩为了让加菲猫干活，竟然让博士对他催眠，这就有点狠了！","tag"
	title := MathRegexpOf1(`"user_bury":0,"title":"(.*)","tag"`, html)

	videoInfo := VideoInfo{}
	videoInfo.Title = title
	videoInfo.Url = xg.Url
	videoInfo.Site = "西瓜视频"

	//"vid":"v02004910000bj377nc1n3e63t2qmc00",
	xg.vid = MathRegexpOf1(`"vid":"(.*)","user_digg"`, html)

	//"businessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjU0MjI5MTMsInZlciI6InYxIiwiYWsiOiJjZmMwNjdiYjM5ZmVmZjU5MmFmODIwODViNDJlNmRjMyIsInN1YiI6InBnY19ub3JtYWwifQ.aaUsIAdV5yqVRZtv4A9G9ijV_GGP261ww-2gK2Asyt0","authToken"
	xg.businessToken = MathRegexpOf1(`"businessToken":"(.*)","authToken"`, html)

	//"authToken":"HMAC-SHA1:2.0:1565422913273604497:cfc067bb39feff592af82085b42e6dc3:RuTnxlMBKOtk8z4p0J\u002F3aWuc27o=","is_original"
	xg.authToken = MathRegexpOf1(`"authToken":"(.*)","is_original"`, html)

	//log.Debug("vid=%s, pToken=%s, authToken=%s", vid, pToken, authToken)
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
		downloadInfo["hd2"] = videoUrl
	}
	if vjson, exists := videoListJson.CheckGet("video_3"); exists {
		videoUrl, err := getVideoUrl(vjson)
		if err != nil {
			return nil, err
		}
		downloadInfo["hd1"] = videoUrl
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

	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	header := http.Header{}
	header.Add("user-agent", ua)
	header.Add("referer", xg.Url)
	header.Add("Authorization", xg.authToken)
	req.Header = header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resultBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return simplejson.NewJson(resultBytes)
}

func MatchRegexp(pattern, content string, index int) string {
	compile := regexp.MustCompile(pattern)
	matches := compile.FindStringSubmatch(content)
	for i, v := range matches {
		if i == index {
			return v
		}
	}
	return ""
}

func MathRegexpOf1(pattern, content string) string {
	return MatchRegexp(pattern, content, 1)
}
