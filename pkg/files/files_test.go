package files_test

import (
	"testing"

	"github.com/teemuteemu/batman/pkg/files"
)

type UnmarshalTestCase struct {
	name        string
	input       string
	output      files.Collection
	expectError bool
}

type FindRequestTestCase struct {
	name            string
	inputCollection files.Collection
	requestName     string
	outputRequest   *files.Request
	expectError     bool
}

const successfullInput = `
version: 42
name: test collection 1
description: test description
requests:
  - name: test get
    method: GET
    url: http://foo.com/get
    header:
      foo: bar
  - name: test post
    method: POST
    url: http://foo.com/post
    header:
      Content-type: "application/json; charset=UTF-8"
      Other-header: "value-123"
    body: >-
      {"test": "post body"}
`

const badYamlInput = `
foobar123
`

var successfullOutput = files.Collection{
	Version:     42,
	Name:        "test collection 1",
	Description: "test description",
	Requests: []files.Request{
		files.Request{
			Name:   "test get",
			Method: "GET",
			URL:    "http://foo.com/get",
			Header: files.Header{
				"foo": "bar",
			},
			Body: "",
		},
		files.Request{
			Name:   "test post",
			Method: "POST",
			URL:    "http://foo.com/post",
			Header: files.Header{
				"Content-type": "application/json; charset=UTF-8",
				"Other-header": "value-123",
			},
			Body: `{"test": "post body"}`,
		},
	},
}

func TestUnmarshalCollection(t *testing.T) {
	testCases := []UnmarshalTestCase{
		UnmarshalTestCase{
			name:        "Successfull input",
			input:       successfullInput,
			output:      successfullOutput,
			expectError: false,
		},
		UnmarshalTestCase{
			name:        "Bad YAML input",
			input:       badYamlInput,
			output:      files.Collection{},
			expectError: true,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.name)

		func(testCase UnmarshalTestCase) {
			testCollection, err := files.UnmarshalCollection([]byte(testCase.input))

			if err != nil && !testCase.expectError {
				t.Error(err)
			}

			if testCollection.Version != testCase.output.Version {
				t.Errorf("Expected version %d, got %d", testCase.output.Version, testCollection.Version)
			}

			if testCollection.Name != testCase.output.Name {
				t.Errorf("Expected name %s, got %s", testCase.output.Name, testCollection.Name)
			}

			if testCollection.Description != testCase.output.Description {
				t.Errorf("Expected description %s, got %s", testCase.output.Description, testCollection.Description)
			}

			if len(testCollection.Requests) != len(testCase.output.Requests) {
				t.Errorf("Expected %d requests, got %d", len(testCase.output.Requests), len(testCollection.Requests))
			}

			for i, request := range testCase.output.Requests {
				outputRequest := testCollection.Requests[i]

				if request.Name != outputRequest.Name {
					t.Errorf("Expected request name %s, got %s", request.Name, outputRequest.Name)
				}

				if request.Method != outputRequest.Method {
					t.Errorf("Expected request method %s, got %s", request.Method, outputRequest.Method)
				}

				if request.URL != outputRequest.URL {
					t.Errorf("Expected request URL %s, got %s", request.URL, outputRequest.URL)
				}

				if request.Body != outputRequest.Body {
					t.Errorf("Expected request body %s, got %s", request.Body, outputRequest.Body)
				}

				if len(request.Header) != len(outputRequest.Header) {
					t.Errorf("Expected %d headers, got %d", len(request.Header), len(outputRequest.Header))
				}

				for key, value := range request.Header {
					if outputRequest.Header[key] != value {
						t.Errorf("Expected %s header, got %s", value, outputRequest.Header[key])
					}
				}
			}
		}(testCase)
	}
}

func TestFindRequest(t *testing.T) {
	testCases := []FindRequestTestCase{
		FindRequestTestCase{
			name:            "Successfully find GET request",
			inputCollection: successfullOutput,
			outputRequest:   &successfullOutput.Requests[0],
			requestName:     "test get",
			expectError:     false,
		},
		FindRequestTestCase{
			name:            "Successfully find POST request",
			inputCollection: successfullOutput,
			outputRequest:   &successfullOutput.Requests[1],
			requestName:     "test post",
			expectError:     false,
		},
		FindRequestTestCase{
			name:            "Request not found",
			inputCollection: successfullOutput,
			outputRequest:   nil,
			requestName:     "request not found",
			expectError:     true,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.name)

		func(testCase FindRequestTestCase) {
			request, err := testCase.inputCollection.FindRequest(testCase.requestName)

			if err != nil && !testCase.expectError {
				t.Error(err)
			}

			if request == nil && testCase.expectError {
				return
			}

			if request.Name != testCase.outputRequest.Name {
				t.Errorf("Expected to find request %s, got %s", testCase.outputRequest.Name, request.Name)
			}

		}(testCase)
	}
}
