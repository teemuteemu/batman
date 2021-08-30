package output

import (
	"fmt"

	"github.com/teemuteemu/batman/pkg/client"
)

type RequestFormatter struct{}

func (f RequestFormatter) RenderRequest(request *client.Request) string {
	return fmt.Sprintf("%s\t%s", request.Method, request.URL)
}

func (f RequestFormatter) RenderResponse(response *client.Response) string {
	return "\n"
}

func (f RequestFormatter) Render(call *client.Call) string {
	req := f.RenderRequest(call.Request)
	res := f.RenderResponse(call.Response)

	return fmt.Sprintf("%s%s", req, res)
}
