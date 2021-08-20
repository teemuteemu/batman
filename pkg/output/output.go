package output

import (
	"bytes"
	"fmt"

	"github.com/teemuteemu/batman/pkg/env"
	"github.com/teemuteemu/batman/pkg/files"
)

type Formatter interface {
	Format(method string, url string, headers files.Header, body *bytes.Buffer) (string, error)
}

type Formaters map[string]Formatter

type Outputter struct {
	Collection *files.Collection
	Env        *env.Env
	Formatters Formaters
}

func New(c *files.Collection, e *env.Env) *Outputter {
	formatters := Formaters{
		"curl": CurlFormatter{},
	}

	return &Outputter{
		Collection: c,
		Env:        e,
		Formatters: formatters,
	}
}

func (o *Outputter) OutputRequests(requestNames []string, format string) error {
	for _, requestName := range requestNames {
		output, err := o.getOutput(requestName, format)
		if err != nil {
			return err
		}

		fmt.Println(output)
	}

	return nil
}

func (o *Outputter) getOutput(requestName, format string) (string, error) {
	request, err := o.Collection.FindRequest(requestName)
	if err != nil {
		return "", err
	}
	if request == nil {
		return "", fmt.Errorf(`Request "%s" not found`, requestName)
	}

	url, err := o.Env.PrepareURL(request.URL)
	if err != nil {
		return "", err
	}

	headers, err := o.Env.PrepareHeader(request.Header)
	if err != nil {
		return "", err
	}

	body, err := o.Env.PrepareBody(request.Body)
	if err != nil {
		return "", err
	}

	if formatter, ok := o.Formatters[format]; ok {
		output, err := formatter.Format(request.Method, url, headers, body)
		if err != nil {
			return "", err
		}

		return output, nil
	}

	return "", fmt.Errorf("Couldn't find output format for %s", format)
}
