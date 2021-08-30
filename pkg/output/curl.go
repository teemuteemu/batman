package output

import (
	"fmt"
	"strings"

	"github.com/teemuteemu/batman/pkg/client"
)

type CurlFormatter struct{}

func (f CurlFormatter) RenderRequest(request *client.Request) string {
	var sb strings.Builder

	sb.WriteString("curl")
	sb.WriteString(" -v")
	sb.WriteString(fmt.Sprintf(" -X %s", request.Method))

	for headerKey, headerValue := range request.Header {
		sb.WriteString(fmt.Sprintf(` -H "%s: %s"`, headerKey, headerValue))
	}

	if len(request.Body) > 0 {
		bodyStr := strings.ReplaceAll(request.Body, "\n", "")
		bodyStr = strings.ReplaceAll(bodyStr, " ", "")

		sb.WriteString(fmt.Sprintf(" -d '%s'", bodyStr))
	}

	sb.WriteString(" ")
	sb.WriteString(request.URL)
	sb.WriteString("\n")

	return sb.String()
}

func (f CurlFormatter) RenderResponse(response *client.Response) string {
	return ConsoleFormatter{}.RenderResponse(response)
}

func (f CurlFormatter) Render(call *client.Call) string {
	req := f.RenderRequest(call.Request)
	res := f.RenderResponse(call.Response)

	return fmt.Sprintf("%s%s", req, res)
}
