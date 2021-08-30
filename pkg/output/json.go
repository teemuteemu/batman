package output

import (
	"encoding/json"
	"fmt"

	"github.com/teemuteemu/batman/pkg/client"
)

type JSONFormatter struct{}

func (f JSONFormatter) RenderRequest(request *client.Request) string {
	res, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err)
	}

	return string(res)
}

func (f JSONFormatter) RenderResponse(response *client.Response) string {
	res, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}

	return fmt.Sprintf("%s\n", string(res))
}

func (f JSONFormatter) Render(call *client.Call) string {
	res, err := json.Marshal(call)
	if err != nil {
		fmt.Println(err)
	}

	return fmt.Sprintf("%s\n", string(res))
}
