package output

import (
	"fmt"

	"github.com/teemuteemu/batman/pkg/client"
	"gopkg.in/yaml.v2"
)

type YAMLFormatter struct{}

func (f YAMLFormatter) RenderRequest(request *client.Request) string {
	res, err := yaml.Marshal(request)
	if err != nil {
		fmt.Println(err)
	}

	return string(res)
}

func (f YAMLFormatter) RenderResponse(response *client.Response) string {
	res, err := yaml.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}

	return string(res)
}

func (f YAMLFormatter) Render(call *client.Call) string {
	res, err := yaml.Marshal(call)
	if err != nil {
		fmt.Println(err)
	}

	return string(res)
}
