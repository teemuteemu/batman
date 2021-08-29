package env

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/joho/godotenv"
	"github.com/teemuteemu/batman/pkg/client"
)

type Env map[string]string

func GetEnv(customLocation string) (Env, error) {
	env := make(Env)

	if len(customLocation) > 0 {
		err := godotenv.Load(customLocation)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Didn't find .env file\n")
		}
	} else {
		err := godotenv.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Didn't find .env file\n")
		}
	}

	for _, e := range os.Environ() {
		kv := strings.SplitN(e, "=", 2)
		env[kv[0]] = kv[1]
	}

	return env, nil
}

func (e *Env) PrepareURL(requestURL string) (string, error) {
	t, err := template.New("request").Parse(requestURL)
	if err != nil {
		return "", err
	}

	var output bytes.Buffer
	err = t.Execute(&output, *e)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func (e *Env) PrepareHeader(requestHeader client.Header) (client.Header, error) {
	headers := client.Header{}

	for key, val := range requestHeader {
		t, err := template.New(key).Parse(val)
		if err != nil {
			return headers, err
		}

		var output bytes.Buffer
		err = t.Execute(&output, *e)
		if err != nil {
			return headers, err
		}

		headers[key] = output.String()
	}

	return headers, nil
}

func (e *Env) PrepareBody(requestBody string) (*bytes.Buffer, error) {
	t, err := template.New("body").Parse(requestBody)
	if err != nil {
		return nil, err
	}

	var output bytes.Buffer
	err = t.Execute(&output, *e)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
