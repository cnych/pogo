package fetch

import (
	"fmt"
	"github.com/cnych/stardust/timex"
	"io"
	"io/ioutil"
	"net/http"
	"pogo/common/call"
	"strings"
	"time"
)

func Fetch(req *Request) (*Response, error) {
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}
	if req.URL == "" {
		return nil, fmt.Errorf("request url is empty")
	}

	if req.Retries > 0 {
		var resp *Response
		var err error

		_ = call.Retry(req.Retries, func() error {
			resp, err = fetchOne(req)
			return err
		})

		return resp, err
	}

	return fetchOne(req)

}

func fetchOne(req *Request) (*Response, error) {
	if req.Method == "" {
		req.Method = http.MethodGet
	}

	var body io.Reader
	if req.Body != "" {
		body = strings.NewReader(req.Body)
	}

	newReq, err := http.NewRequest(req.Method, req.URL, body)
	if err != nil {
		return nil, err
	}

	// header
	header := http.Header{}
	for k, v := range req.Header {
		header.Add(k, v)
	}
	newReq.Header = header

	var timeout time.Duration
	if req.Timeout > 0 {
		timeout = timex.DurationMS(req.Timeout)
	}
	client := http.Client{Timeout: timeout}

	newResp, err := client.Do(newReq)
	if err != nil {
		return nil, err
	}
	defer newResp.Body.Close()

	resp := Response{}
	resp.StatusCode = newResp.StatusCode
	resp.URL = req.URL
	if len(newResp.Header) >0 {
		resp.Header = map[string]string{}
	}
	for k, vl := range newResp.Header {
		resp.Header[k] = vl[0]
	}

	respBody, err := ioutil.ReadAll(newResp.Body)
	if err != nil {
		return nil, err
	}
	resp.SetBody(respBody)

	return &resp, nil
}
