package main

import (
	"apps/pkgs/paypal_services"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var svc = paypal_services.NewPaypalSVC()

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//
	err := svc.Auth()
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, err
	}
	//
	plans, err := svc.List_plans()
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       plans,
	}, nil
}

func main() {
	lambda.Start(handler)
}
