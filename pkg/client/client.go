package client

import (
	"io"
	"net/http"
	"strings"
)

type Formatter interface {
	RenderRequest(request *Request) string
	RenderResponse(response *Response) string
	Render(call *Call) string
}

type Formaters map[string]Formatter

type Header map[string]string

type Request struct {
	Name   string `yaml:"name" json:"name"`
	Method string `yaml:"method" json:"method"`
	URL    string `yaml:"url" json:"url"`
	Header Header `yaml:"header" json:"header"`
	Body   string `yaml:"body" json:"body"`
}

type Response struct {
	StatusCode int    `yaml:"status" json:"status"`
	Header     Header `yaml:"header" json:"header"`
	Body       string `yaml:"body" json:"body"`
}

type Call struct {
	Request  *Request  `yaml:"request" json:"request"`
	Response *Response `yaml:"response" json:"response"`
}

func renderBody(reader io.ReadCloser) (string, error) {
	defer reader.Close()

	body, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func ExecuteRequest(request *Request) (*Response, error) {
	httpReq, err := http.NewRequest(request.Method, request.URL, strings.NewReader(request.Body))
	if err != nil {
		return nil, err
	}

	for key, val := range request.Header {
		httpReq.Header.Add(key, val)
	}

	client := &http.Client{}

	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	respBody, err := renderBody(httpResp.Body)
	if err != nil {
		return nil, err
	}

	resp := &Response{
		StatusCode: httpResp.StatusCode,
		Header:     make(Header),
		Body:       respBody,
	}

	for key, val := range httpResp.Header {
		resp.Header[key] = val[0]
	}

	return resp, nil
}
