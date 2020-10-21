package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var (
	// 是否本地測試, 關閉處理 aws 的外部服務
	isLocalTest = os.Getenv("IsLocalTest")
	//
	internalNotificationArn = os.Getenv("InternalNotificationArn")
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//
	fmt.Printf("isLocalTest : %s\n", isLocalTest)
	fmt.Printf("internalNotificationArn : %s\n", internalNotificationArn)
	fmt.Printf("[request.Body] %s\n\n", request.Body)
	//
	paypal_header := map[string]string{}
	for key, value := range request.Headers {
		if strings.HasPrefix(key, "PAYPAL") {
			paypal_header[strings.ToLower(key)] = value
		}
	}

	// sns
	header, _ := json.Marshal(paypal_header)
	encoded := base64.StdEncoding.EncodeToString([]byte(request.Body))
	obj := SNSMessage{
		Paypal_header: string(header),
		Webhook_event: encoded,
	}
	m, _ := json.Marshal(obj)
	message := string(m)
	fmt.Println("SNS Message =>", message)

	if strings.ToLower(isLocalTest) != "true" {
		err := SNS_Process(message)
		if err != nil {
			fmt.Println(err)
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
			}, err
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       request.Body,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

type SNSMessage struct {
	Paypal_header string `json:"paypal_header"`
	Webhook_event string `json:"webhook_event"`
}

func SNS_Process(message string) error {
	mySession := session.Must(session.NewSession())
	svc := sns.New(mySession)
	publishInput := sns.PublishInput{
		TargetArn: aws.String(internalNotificationArn),
		Message:   aws.String(message),
		// MessageStructure: aws.String("json"),
	}
	_, err := svc.Publish(&publishInput)
	return err
}
