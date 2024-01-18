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

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	region := os.Getenv("AWS_REGION")
	currSession, _ := session.NewSession(&aws.Config{ // Use aws sdk to connect to dynamoDB
		Region: &region,
	})

	svc := lambda2.New(currSession)
	body, err := json.Marshal(map[string]interface{}{
		"key1": 2,
		"key2": "value2",
		"key3": "value3",
	})

	type Payload struct {
		// You can also include more objects in the structure like below,
		// but for my purposes body was all that was required
		// Method string `json:"httpMethod"`
		Body string `json:"body"`
	}
	p := Payload{
		// Method: "POST",
		Body: string(body),
	}
	payload, err := json.Marshal(p)
	// Result should be: {"body":"{\"name\":\"Jimmy\"}"}
	// This is the required format for the lambda request body.

	if err != nil {
		fmt.Println("Json Marshalling error")
	}
	//fmt.Println(string(payload))

	input := &lambda2.InvokeInput{
		FunctionName:   aws.String("fundamental_data_interpreter"),
		InvocationType: aws.String("RequestResponse"),
		LogType:        aws.String("Tail"),
		Payload:        payload,
	}
	result, err := svc.Invoke(input)
	if err != nil {
		fmt.Println("error")
		fmt.Println(err.Error())
	}
	var m map[string]interface{}
	json.Unmarshal(result.Payload, &m)

	//type Body struct {
	//	Foo string `json:"foo"`
	//}
	//
	//type Message struct {
	//	Body Body `json:"body"`
	//}
	//
	//message := Message{}
	//json.Unmarshal(result.Payload, &message)
	//fmt.Println(message.Body)
	//fmt.Println(message.Body.Foo)

	invokeReponse, err := json.Marshal(m["body"])
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

func main() {
	lambda.Start(Handler)
}
