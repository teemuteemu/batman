package renderer_test

import (
	"testing"

	"github.com/teemuteemu/batman/pkg/renderer"
)

func TestFormatJSON(t *testing.T) {
	inputJSON := `{"some_object":{"key":"value","another_key":123}}`
	expectedJSON := `{
  "some_object": {
    "key": "value",
    "another_key": 123
  }
}`
	result, err := renderer.FormatJSON(inputJSON)
	if err != nil {
		t.Error(err)
	}

	if result != expectedJSON {
		t.Errorf("Expected JSON %s, got %s", expectedJSON, result)
	}
}
