package client_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/teemuteemu/batman/pkg/client"
	"github.com/teemuteemu/batman/pkg/files"
)

type TestCase struct {
	name            string
	reqMethod       string
	reqHeader       files.Header
	reqBody         *bytes.Buffer
	respStatusCode  int
	respBody        string
	expectError     bool
	expectedHeaders int
}

func TestExecute(t *testing.T) {
	testCases := []TestCase{
		TestCase{
			name:            "Successfull GET",
			reqMethod:       "GET",
			reqHeader:       files.Header{},
			reqBody:         bytes.NewBufferString(""),
			respStatusCode:  200,
			respBody:        "test body",
			expectError:     false,
			expectedHeaders: 0,
		},
		TestCase{
			name:            "Successfull POST",
			reqMethod:       "POST",
			reqHeader:       files.Header{},
			reqBody:         bytes.NewBufferString("POST body"),
			respStatusCode:  201,
			respBody:        "test body",
			expectError:     false,
			expectedHeaders: 0,
		},
		TestCase{
			name:            "Successfull GET with error status",
			reqMethod:       "GET",
			reqHeader:       files.Header{},
			reqBody:         bytes.NewBufferString(""),
			respStatusCode:  503,
			respBody:        "fail",
			expectError:     false,
			expectedHeaders: 0,
		},
		TestCase{
			name:      "Successfull GET with added headers",
			reqMethod: "GET",
			reqHeader: files.Header{
				"test_key1": "test_value1",
				"test_key2": "test_value2",
			},
			reqBody:         bytes.NewBufferString(""),
			respStatusCode:  503,
			respBody:        "fail",
			expectError:     false,
			expectedHeaders: 2,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.name)

		func(testCase TestCase) {
			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(testCase.respStatusCode)
				res.Write([]byte(testCase.respBody))
			}))
			defer func() { testServer.Close() }()

			resp, err := client.ExecuteRequest(testCase.reqMethod, testServer.URL, testCase.reqHeader, testCase.reqBody)
			if err != nil && !testCase.expectError {
				t.Error(err)
			}

			if resp.StatusCode != testCase.respStatusCode {
				t.Errorf("Expected status code %d, got %d", testCase.respStatusCode, resp.StatusCode)
			}

			/*
				if len(resp.Header) != testCase.expectedHeaders {
					t.Errorf("Expected %d headers, got %d", testCase.expectedHeaders, len(resp.Header))
				}
			*/

			if string(resp.Body) != testCase.respBody {
				t.Errorf("Expected response body %s, got %s", testCase.respBody, resp.Body)
			}
		}(testCase)
	}

}
