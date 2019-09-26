package parser

import "fmt"

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
	fmt.Printf("downloading douyin video\n")
	return nil, nil
}