package output

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/teemuteemu/batman/pkg/files"
)

type CurlFormatter struct{}

func (cf CurlFormatter) Format(method string, url string, headers files.Header, body *bytes.Buffer) (string, error) {
	var sb strings.Builder

	sb.WriteString("curl")
	sb.WriteString(" -v")
	sb.WriteString(fmt.Sprintf(" -X %s", method))

	for headerKey, headerValue := range headers {
		sb.WriteString(fmt.Sprintf(` -H "%s: %s"`, headerKey, headerValue))
	}

	if body != nil && len(body.Bytes()) > 0 {
		bodyStr := strings.ReplaceAll(body.String(), "\n", "")
		bodyStr = strings.ReplaceAll(bodyStr, " ", "")

		sb.WriteString(fmt.Sprintf(" -d '%s'", bodyStr))
	}

	sb.WriteString(" ")
	sb.WriteString(url)

	return sb.String(), nil
}
