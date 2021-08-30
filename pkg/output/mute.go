package output

import (
	"fmt"

	"github.com/teemuteemu/batman/pkg/client"
)

type MuteFormatter struct{}

func (f MuteFormatter) RenderRequest(request *client.Request) string {
	return ""
}

func (f MuteFormatter) RenderResponse(response *client.Response) string {
	return ""
}

func (f MuteFormatter) Render(call *client.Call) string {
	req := f.RenderRequest(call.Request)
	res := f.RenderResponse(call.Response)

	return fmt.Sprintf("%s%s", req, res)
}
