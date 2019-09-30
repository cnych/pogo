package convert

import (
	"fmt"
	"os"
	"path"
	"pogo/common/call"
	"strconv"
	"strings"
)


type FFMpeg struct {
	ffmpegPath string
}

func NewFFMpeg() (*FFMpeg) {
	ffm := FFMpeg{ffmpegPath: ""}
	ffm.ffmpegPath = ffm.which("ffmpeg")
	return &ffm
}

func XOk(filepath string) (result bool) {
	fileInfo, err := os.Stat(filepath)
	if err != nil && os.IsNotExist(err) {
		return
	}
	mode := fileInfo.Mode().String()
	if strings.Index(mode, "x") > 0 {
		result = true
	}
	return
}

func (ffm *FFMpeg) Duration(url string) (int, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ffmpge Merge error: ", err)
		}
	}()

	result, err := call.Cmd(fmt.Sprintf("%s -i '%s' 2>&1 | grep 'Duration' | cut -d ' ' -f 4 | sed s/,//", ffm.ffmpegPath, url))
	if err != nil {
		return 0, err
	}
	duration, err := ffm.parseDuration(result)
	if err != nil {
		return 0, err
	}
	return duration, nil
}

func (ffm *FFMpeg) parseDuration(duration string) (int, error) {
	//00:00:30.05
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

func (ffm *FFMpeg) which(name string) (filepath string) {
	for _, value := range os.Environ() {
		if strings.Index(value, "PATH=") != 0 {
			continue
		}
		paths := strings.Split(value[5:], ":")
		for _, p := range paths {
			fp := path.Join(p, name)
			if XOk(fp) {
				filepath = fp
				break
			}
		}
		if len(filepath) > 0 {
			break
		}
	}
	return
}
