package env_test

import (
	"fmt"
	"testing"

	"github.com/teemuteemu/batman/pkg/env"
	"github.com/teemuteemu/batman/pkg/files"
)

var testEnv = make(env.Env)

func TestPrepareURL(t *testing.T) {
	expectedHost := "foo.bar.com"
	expectedPath1 := "foo"
	expectedPath2 := "bar"
	expectedURL := fmt.Sprintf("https://%s/%s/%s", expectedHost, expectedPath1, expectedPath2)
	testEnv["base_url"] = expectedHost
	testEnv["path_1"] = expectedPath1
	testEnv["path_2"] = expectedPath2

	result, err := testEnv.PrepareURL("https://{{.base_url}}/{{.path_1}}/{{.path_2}}")
	if err != nil {
		t.Error(err)
	}

	if result != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, result)
	}
}

func TestPrepareHeader(t *testing.T) {
	expectedHeaderKey1 := "test_header_key_1"
	expectedHeaderKey2 := "test_header_key_2"
	expectedHeaderValue1 := "test_header_value_1"
	expectedHeaderValue2 := "test_header_value_2"
	expectedHeader := files.Header{
		expectedHeaderKey1: expectedHeaderValue1,
		expectedHeaderKey2: expectedHeaderValue2,
	}
	testEnv["header_value_1"] = expectedHeaderValue1
	testEnv["header_value_2"] = expectedHeaderValue2
	inputHeader := files.Header{
		expectedHeaderKey1: "{{.header_value_1}}",
		expectedHeaderKey2: "{{.header_value_2}}",
	}

	result, err := testEnv.PrepareHeader(inputHeader)
	if err != nil {
		t.Error(err)
	}

	for key, value := range result {
		if expectedHeader[key] != value {
			t.Errorf("Expected header key %s to have value %s, got %s", key, expectedHeader[key], value)
		}
	}
}

func TestPrepareBody(t *testing.T) {
	expectedKey := "test_key"
	expectedValue := "test_value"
	expectedBody := fmt.Sprintf(`{"%s": "%s"}`, expectedKey, expectedValue)
	testEnv["key"] = expectedKey
	testEnv["value"] = expectedValue

	result, err := testEnv.PrepareBody(`{"{{.key}}": "{{.value}}"}`)
	if err != nil {
		t.Error(err)
	}

	if result.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, result.String())
	}
}
