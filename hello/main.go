package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	lambda2 "github.com/aws/aws-sdk-go/service/lambda"
)

type RequestBody struct {
	RawData    []float64 `json:"raw_data"`
	SomeString string    `json:"some_string"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// unmarshal API request
	var requestBody RequestBody
	json.Unmarshal([]byte(request.Body), &requestBody)

	// process raw data
	rawDataMean := calculateMean(requestBody.RawData)

	// prepare request to python lambda
	region := os.Getenv("AWS_REGION")
	currSession, _ := session.NewSession(&aws.Config{ // Use aws sdk to connect to dynamoDB
		Region: &region,
	})

	svc := lambda2.New(currSession)
	body, err := json.Marshal(map[string]interface{}{
		"key1": rawDataMean,
		"key2": "value2",
		"key3": "value3",
	})

	type Payload struct {
		Body string `json:"body"`
	}
	p := Payload{
		// Method: "POST",
		Body: string(body),
	}
	payload, err := json.Marshal(p)
	// Result should be: {"body":"{\"key1\":\"2\"}"}
	// This is the required format for the lambda request body.

	if err != nil {
		fmt.Println("Json Marshalling error")
	}

	input := &lambda2.InvokeInput{
		FunctionName:   aws.String("fundamental_data_interpreter"),
		InvocationType: aws.String("RequestResponse"),
		LogType:        aws.String("Tail"),
		Payload:        payload,
	}

	// call python lambda
	result, err := svc.Invoke(input)
	if err != nil {
		fmt.Println("error")
		fmt.Println(err.Error())
	}

	// unmarshal response from python lambda
	var m map[string]interface{}
	json.Unmarshal(result.Payload, &m)
	invokeReponse, err := json.Marshal(m["body"])

	// send API response
	resp := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(invokeReponse),
	}
	fmt.Println(resp)

	return resp, nil
}

func calculateMean(data []float64) float64 {
	var sum float64
	for _, number := range data {
		sum += number
	}
	return sum / float64(len(data))
}

func main() {
	lambda.Start(Handler)
}
