package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
	print("test run")
	return
}

func TestHandler(t *testing.T) {
	tests := map[string]struct {
		request        events.APIGatewayProxyRequest
		expectedOutput events.APIGatewayProxyResponse
		expectedError  assert.ErrorAssertionFunc
	}{
		"happy path": {
			request: events.APIGatewayProxyRequest{
				Body: `
					{
						"raw_data": [1, 2.0, 3.4],
						"key2": "value2"
					}
				`,
			},
			expectedOutput: events.APIGatewayProxyResponse{
				StatusCode:      200,
				Headers:         nil,
				Body:            "",
				IsBase64Encoded: false,
			},
			expectedError: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err)
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			os.Setenv("AWS_REGION", "us-east-1")

			got, err := Handler(tt.request)

			if !tt.expectedError(t, err) {
				return
			}

			assert.Equal(t, tt.expectedOutput, got)
		})
	}
}
