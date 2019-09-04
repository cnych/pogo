package fetch

import (
	"github.com/bitly/go-simplejson"
	"github.com/cnych/stardust/unicodex"
	"net/http"
)

type Request struct {
	Method string `json:"method"`
	URL string `json:"url"`
	Header map[string]string `json:"header"`
	Body string `json:"body,omitempty"`
	Retries int `json:"retries,omitempty"`
	Timeout int64 `json:"timeout,omitempty"`
}

func NewGetRequest(url string, retries int, timeout int64, header map[string]string) *Request {
	return &Request{
		URL: url,
		Method: http.MethodGet,
		Header: header,
		Retries: retries,
		Timeout: timeout,
	}
}

// DefaultRequest ... 定义的一个默认的Request对象
func DefaultRequest(url string, header map[string]string) *Request {
	return NewGetRequest(url, 3, 3000, header)
}

type Response struct {
	StatusCode int `json:"status_code"`
	URL string `json:"url"`
	Header map[string]string `json:"header"`
	Body string `json:"body,omitempty"`
}

func (resp *Response) AsBytes() ([]byte, error) {
	if resp.Body == "" {
		return nil, nil
	}
	return []byte(resp.Body), nil
}

func (resp *Response) AsText(charset string) (string, error) {
	buff, err := resp.AsBytes()
	if err != nil {
		return "", err
	}
	return unicodex.Decode(buff, charset)
}

func (resp *Response) ParseSimpleJSON() (*simplejson.Json, error) {
	bodyBytes, err := resp.AsBytes()
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(bodyBytes)
}

func (resp *Response) SetBody(data []byte) {
	if len(data) == 0 {
		return
	}
	resp.Body = string(data)
}
