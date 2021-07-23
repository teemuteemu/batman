package renderer

import (
	"bytes"
	"encoding/json"
)

func FormatJSON(jsonBody string) (string, error) {
	var out bytes.Buffer

	err := json.Indent(&out, []byte(jsonBody), "", "  ")
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
