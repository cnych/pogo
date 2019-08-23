package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"pogo/common/logs"
	"strconv"
)

var log = logs.Log

func Download(filename, url string, header http.Header) (result string, err error) {
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header = header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	contentLen, err := strconv.ParseInt(resp.Header["Content-Length"][0], 10, 64)
	if err != nil {
		return "", err
	}
	io.Copy(file, resp.Body)
	defer resp.Body.Close()

	if fileInfo, err := os.Stat(filename); err == nil {
		fileLength := fileInfo.Size()
		if contentLen == fileLength {
			result = filename
		} else {
			log.Debug("file size not equal %s,%d,%d ", filename, fileLength, contentLen)
			return "", fmt.Errorf("download video error")
		}

	} else {
		log.Debug("file %s not exists ", filename)
		return "", err
	}
	return
}
