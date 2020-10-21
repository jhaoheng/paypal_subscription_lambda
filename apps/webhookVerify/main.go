package main

import (
	"apps/pkgs/paypal_services"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var (
	// 是否本地測試, 關閉處理 aws 的外部服務
	isLocalTest = os.Getenv("IsLocalTest")
	//
	verifiedQueueUrl = os.Getenv("VerifiedQueueUrl")
	errorQueueUrl    = os.Getenv("ErrorQueueUrl")
)

func handler(snsEvent events.SNSEvent) error {
	//
	fmt.Printf("isLocalTest : %s\n", isLocalTest)
	fmt.Printf("verifiedQueueUrl : %s\n", verifiedQueueUrl)
	fmt.Printf("errorQueueUrl : %s\n", errorQueueUrl)
	fmt.Printf("[snsEvent.Records] %#v\n\n", snsEvent.Records)
	//

	var (
		queueUrl    string
		messageBody string
	)

	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		fmt.Printf("snsRecord.Message = %s \n", snsRecord.Message)
		snsMessage := SNSMessage{}
		json.Unmarshal([]byte(snsRecord.Message), &snsMessage)
		fmt.Printf("message obj = %s \n", snsMessage)
		//
		paypal_header := PaypalHeader{}
		json.Unmarshal([]byte(snsMessage.Paypal_header), &paypal_header)
		fmt.Println("paypal_header=>", snsMessage.Paypal_header)
		//
		decoded, _ := base64.StdEncoding.DecodeString(snsMessage.Webhook_event)
		webhook_event := string(decoded)
		fmt.Println("webhook_event =>", webhook_event)
		//
		if ok, err := paypal_verify(paypal_header, webhook_event); !ok {
			fmt.Println("error => ", err)
			queueUrl = errorQueueUrl
			messageBody = snsRecord.Message
		} else {
			queueUrl = verifiedQueueUrl
			messageBody = webhook_event
		}

		// sqs
		if strings.ToLower(isLocalTest) == "true" {
			aws_services_process(queueUrl, messageBody)
		}
	}
	return nil
}

func main() {
	lambda.Start(handler)
}

// sqs
var sqs_agent = sqs.New(session.Must(session.NewSession()))

func aws_services_process(queueUrl, messageBody string) {
	input := sqs.SendMessageInput{
		QueueUrl:    aws.String(queueUrl),
		MessageBody: aws.String(messageBody),
	}
	fmt.Println("sqs input =>", input)
	_, err := sqs_agent.SendMessage(&input)
	if err != nil {
		fmt.Println(err)
	}
}

// paypal
type SNSMessage struct {
	Paypal_header string `json:"paypal_header"`
	Webhook_event string `json:"webhook_event"`
}

type PaypalHeader struct {
	PAYPAL_AUTH_ALGO         string `json:"paypal-auth-algo"`
	PAYPAL_CERT_URL          string `json:"paypal-cert-url"`
	PAYPAL_AUTH_VERSION      string `json:"paypal-auth-version"`
	PAYPAL_TRANSMISSION_ID   string `json:"paypal-transmission-id"`
	PAYPAL_TRANSMISSION_SIG  string `json:"paypal-transmission-sig"`
	PAYPAL_TRANSMISSION_TIME string `json:"paypal-transmission-time"`
}

var svc = paypal_services.NewPaypalSVC()

func paypal_verify(header PaypalHeader, webhook_event string) (ok bool, err error) {
	err = svc.Auth()
	if err != nil {
		return false, err
	}
	//
	w := paypal_services.WebhookVerifySignature{
		Auth_algo:         header.PAYPAL_AUTH_ALGO,
		Cert_url:          header.PAYPAL_CERT_URL,
		Transmission_id:   header.PAYPAL_TRANSMISSION_ID,
		Transmission_sig:  header.PAYPAL_TRANSMISSION_SIG,
		Transmission_time: header.PAYPAL_TRANSMISSION_TIME,
		Webhook_event:     webhook_event,
	}
	return svc.Webhook_verify_signature(w)
}
