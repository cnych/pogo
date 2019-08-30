package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"pogo/common/bar"
	"strconv"
)


// Download ... 下载网络文件到本地文件
func Download(filename, url string, header http.Header) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header = header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	contentLen, err := strconv.ParseInt(resp.Header["Content-Length"][0], 10, 64)
	if err != nil {
		return err
	}

	br := bar.NewBar(contentLen)
	br.Resize = func(bar *bar.Bar) error {
		fileInfo, err := os.Stat(filename)
		if err == nil {
			br.Size = fileInfo.Size()
		}
		return err
	}
	br.Start()

	go func() {
		io.Copy(file, resp.Body)
		defer resp.Body.Close()
	}()

	br.ShowProgress()

	if fileInfo, err := os.Stat(filename); err == nil {
		fileLength := fileInfo.Size()
		if contentLen == fileLength {
			return nil
		}
		return fmt.Errorf("file size not equal %s,%d,%d", filename, fileLength, contentLen)
	}

	return fmt.Errorf("file %s not exists", filename)

}
