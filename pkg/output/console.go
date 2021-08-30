package output

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/teemuteemu/batman/pkg/client"
)

type ConsoleFormatter struct{}

func (f ConsoleFormatter) RenderRequest(request *client.Request) string {
	return fmt.Sprintf("%s\t%s\t", request.Method, request.URL)
}

func (f ConsoleFormatter) RenderResponse(response *client.Response) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%d - %s\n\n", response.StatusCode, http.StatusText(response.StatusCode)))

	for key, val := range response.Header {
		sb.WriteString(fmt.Sprintf("%s: %s\n", key, val))
	}
	sb.WriteString("\n")

	sb.WriteString(response.Body)
	sb.WriteString("\n")

	return sb.String()
}

func (f ConsoleFormatter) Render(call *client.Call) string {
	req := f.RenderRequest(call.Request)
	res := f.RenderResponse(call.Response)

	return fmt.Sprintf("%s%s", req, res)
}
