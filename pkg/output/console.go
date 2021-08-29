package output

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/teemuteemu/batman/pkg/client"
	"github.com/teemuteemu/batman/pkg/renderer"
)

type ConsoleFormatter struct{}

func (f ConsoleFormatter) RenderRequest(request *client.Request) string {
	return fmt.Sprintf("%s\t%s", request.Method, request.URL)
}

func (f ConsoleFormatter) RenderResponse(response *client.Response) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\t%d - %s\n\n", response.StatusCode, http.StatusText(response.StatusCode)))

	for key, val := range response.Header {
		sb.WriteString(fmt.Sprintf("%s: %s\n", key, val))
	}
	sb.WriteString("\n")

	jsonOutput, err := renderer.FormatJSON(response.Body)
	if err != nil {
		panic(err)
	}

	sb.WriteString(jsonOutput)
	sb.WriteString("\n")

	return sb.String()
}

func (f ConsoleFormatter) Render(call *client.Call) string {
	req := f.RenderRequest(call.Request)
	res := f.RenderResponse(call.Response)

	return fmt.Sprintf("%s%s", req, res)
}
