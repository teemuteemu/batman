package client

import (
	"bytes"
	"io"
	"net/http"

	"github.com/teemuteemu/batman/pkg/files"
)

type ResponseHeader map[string]string

type Response struct {
	StatusCode int
	Header     ResponseHeader
	Body       string
}

func renderBody(reader io.ReadCloser) (string, error) {
	defer reader.Close()

	body, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func ExecuteRequest(method, url string, headers files.Header, body *bytes.Buffer) (*Response, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for key, val := range headers {
		request.Header.Add(key, val)
	}

	client := &http.Client{}

	httpResp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	respBody, err := renderBody(httpResp.Body)
	if err != nil {
		return nil, err
	}

	resp := &Response{
		StatusCode: httpResp.StatusCode,
		Header:     make(ResponseHeader),
		Body:       respBody,
	}

	for key, val := range httpResp.Header {
		resp.Header[key] = val[0]
	}

	return resp, nil
}
