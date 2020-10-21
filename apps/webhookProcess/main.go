package main

import (
	"apps/pkgs/paypal_services"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var (
	// 是否本地測試, 關閉處理 aws 的外部服務
	isLocalTest = os.Getenv("IsLocalTest")
	//
	verifiedQueueUrl           = os.Getenv("VerifiedQueueUrl")
	eventFinishNotificationArn = os.Getenv("EventFinishNotificationArn")
)

var queueUrl = verifiedQueueUrl
var svc_sqs = sqs.New(session.Must(session.NewSession()))
var svc_sns = sns.New(session.Must(session.NewSession()))

func handler(sqsEvent events.SQSEvent) (err error) {
	//
	fmt.Printf("isLocalTest : %s\n", isLocalTest)
	fmt.Printf("verifiedQueueUrl : %s\n", verifiedQueueUrl)
	fmt.Printf("eventFinishNotificationArn : %s\n", eventFinishNotificationArn)
	fmt.Printf("[sqsEvent.Records] %#v\n\n", sqsEvent.Records)
	//

	for _, message := range sqsEvent.Records {
		fmt.Printf("message => %s \n", message.Body)
		webhook_event := paypal_services.WEBHOOK_EVENT{}
		err := json.Unmarshal([]byte(message.Body), &webhook_event)
		if err != nil {
			return err
		}

		// 事件處理
		/* webhook events
		- https://developer.paypal.com/docs/api-basics/notifications/webhooks/event-names/#subscriptions
		*/
		if webhook_event.Event_type == "PAYMENT.SALE.COMPLETED" {
			fmt.Println("==> 付款成功, 但未執行發票寄送, 因需要額外設定 invoice 上的 items 與 收貨人訊息, 請更新 code")

			/*
				// 付款完成, 執行發票寄送
				// 1. 取得 subscription 資訊
				// 2. 建立 invoice
				subscription_id := "請在 custom_id 中夾帶 subscription_id"
				resp_detail, err := invoice.Get_subsctiption_detail(subscription_id)
				if err != nil {
					return err
				}
				var billing_info = invoice.Billing_Info{}
				var item = invoice.Item{}
				invoice_view_url, err := invoice.Build(billing_info, item)
				if err != nil {
					return err
				}
				fmt.Println("發票 url =>", invoice_view_url)
			*/
		}

		if strings.ToLower(isLocalTest) != "true" {
			aws_services_process(message)
		}
	}

	return err
}

func main() {
	lambda.Start(handler)
}

func aws_services_process(message events.SQSMessage) {
	// 所有事件, 送到 sns 中, 通知管理員(slack / email)
	sns_input := sns.PublishInput{
		Message:   aws.String(message.Body),
		TargetArn: aws.String(eventFinishNotificationArn),
	}
	svc_sns.Publish(&sns_input)
	// 從 sqs 中刪除事件
	sqs_input := sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueUrl),
		ReceiptHandle: aws.String(message.ReceiptHandle),
	}
	svc_sqs.DeleteMessage(&sqs_input)
}
