package parser

import (
	browser "github.com/EDDYCJY/fake-useragent"
	"pogo/common/call"
	"pogo/common/convert"
	"pogo/common/nets/fetch"
	"pogo/common/strs"
)

type Douyin struct {
	BaseVideo
	vid string
}

func DouyinRegister()  {
	dy := new(Douyin)
	dy.Name = "douyin"
//http://v.douyin.com/u2f91P/
	dy.VideoPatterns = []string{`v\.douyin\.com/(\w+)`}
	Spiders[dy.Name] = dy
}

func (dy *Douyin) GetVideoInfo() (info *VideoInfo, err error)  {
	_ = call.Retry(3, func() error {
		info, err = getDouyinOnce(dy)
		return err
	})
	return
}

func getDouyinOnce(dy *Douyin) (info *VideoInfo, err error) {
	header := map[string]string {
		"user-agent": browser.Mobile(),
	}
	req := fetch.DefaultRequest(dy.Url, header)

	resp, err := fetch.Fetch(req)
	if err != nil {
		return nil, err
	}

	html, err := resp.AsText("UTF-8")
	if err != nil {
		return nil, err
	}

	//<div class="user-title">满院果蔬熟透，豇豆长势尤其猛，来做些干豇豆吧 </div><div class="user-avator"
	title := strs.MatchRegexpOf1(`<div class="user-title">(.*)</div><div class="user-avator"`, html)
	//<video id="theVideo" class="video-player" src="xxxxx" preload="auto"
	videoUrl := strs.MatchRegexpOf1(`class="video-player" src="(.*)" preload="auto"`, html)

	downloadInfo := make(map[string]interface{})
	downloadInfo["normal"] = videoUrl

	// 获取 Duration
	duration, err := convert.NewFFMpeg().Duration(videoUrl)
	if err != nil {
		duration = 0
	}

	videoInfo := VideoInfo{}
	videoInfo.Title = title
	videoInfo.Url = dy.Url
	videoInfo.Site = "抖音短视频"
	videoInfo.DownloadInfo = downloadInfo
	videoInfo.Duration = int64(duration)

	return &videoInfo, nil
}
