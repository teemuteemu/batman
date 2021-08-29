package client_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/teemuteemu/batman/pkg/client"
)

type TestCase struct {
	name           string
	request        client.Request
	respStatusCode int
	respBody       string
	expectError    bool
}

func TestExecute(t *testing.T) {
	testCases := []TestCase{
		{
			name: "Successfull GET",
			request: client.Request{
				Method: "GET",
				Header: client.Header{},
				Body:   "",
			},
			respStatusCode: 200,
			respBody:       "test body",
			expectError:    false,
		},
		{
			name: "Successfull POST",
			request: client.Request{
				Method: "POST",
				Header: client.Header{},
				Body:   "POST body",
			},
			respStatusCode: 201,
			respBody:       "test body",
			expectError:    false,
		},
		{
			name: "Successfull GET with error status",
			request: client.Request{
				Method: "GET",
				Header: client.Header{},
				Body:   "",
			},
			respStatusCode: 503,
			respBody:       "fail",
			expectError:    false,
		},
		{
			name: "Successfull GET with added headers",
			request: client.Request{
				Method: "GET",
				Header: client.Header{
					"test_key1": "test_value1",
					"test_key2": "test_value2",
				},
				Body: "",
			},
			respStatusCode: 503,
			respBody:       "fail",
			expectError:    false,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.name)

		func(testCase TestCase) {
			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(testCase.respStatusCode)
				_, err := res.Write([]byte(testCase.respBody))
				if err != nil {
					t.Fatal(err)
				}
			}))
			defer func() { testServer.Close() }()

			testCase.request.URL = testServer.URL

			resp, err := client.ExecuteRequest(&testCase.request)
			if err != nil && !testCase.expectError {
				t.Error(err)
			}

			if resp.StatusCode != testCase.respStatusCode {
				t.Errorf("Expected status code %d, got %d", testCase.respStatusCode, resp.StatusCode)
			}

			if string(resp.Body) != testCase.respBody {
				t.Errorf("Expected response body %s, got %s", testCase.respBody, resp.Body)
			}
		}(testCase)
	}

}
