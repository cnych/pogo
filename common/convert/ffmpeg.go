package convert

import (
	"fmt"
	"pogo/common/call"
	"pogo/common/files"
	"strconv"
	"strings"
)

type FFMpeg struct {
	ffmpegPath string
}

func NewFFMpeg() *FFMpeg {
	ffm := FFMpeg{ffmpegPath: ""}
	ffm.ffmpegPath = files.Which("ffmpeg")
	return &ffm
}

func (ffm *FFMpeg) Duration(url string) (int, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("ffmpeg get duration error: %v\n", err)
		}
	}()

	durationCmd := fmt.Sprintf("%s -i '%s' 2>&1 | grep Duration | cut -d ' ' -f 4 | sed s/,//", ffm.ffmpegPath, url)
	result, err := call.Cmd(durationCmd)
	if err != nil {
		return 0, err
	}
	//	00:00:30.05
	duration, err := ffm.parseDuration(result)
	if err != nil {
		return 0, err
	}

	return duration, nil
}

func (ffm *FFMpeg) parseDuration(duration string) (int, error) {
	//	00:00:30.05
	//[00, 00, 30.05] [30]
	arr := strings.Split(duration, ":")
	if len(arr) == 3 {
		h, err := strconv.Atoi(arr[0])
		if err != nil {
			return 0, err
		}

		m, err := strconv.Atoi(arr[1])
		if err != nil {
			return 0, err
		}

		sArr := strings.Split(arr[2], ".")
		s, err := strconv.Atoi(sArr[0])
		if err != nil {
			return 0, err
		}

		return h * 60 * 60 + m * 60 + s, nil

	}

	return 0, fmt.Errorf("duration format error")
}

