package output

import (
	"fmt"

	"github.com/teemuteemu/batman/pkg/client"
)

type ResponseFormatter struct{}

func (f ResponseFormatter) RenderRequest(request *client.Request) string {
	return "\n"
}

func (f ResponseFormatter) RenderResponse(response *client.Response) string {
	return ConsoleFormatter{}.RenderResponse(response)
}

func (f ResponseFormatter) Render(call *client.Call) string {
	req := f.RenderRequest(call.Request)
	res := f.RenderResponse(call.Response)

	return fmt.Sprintf("%s%s", req, res)
}
