package main

import (
	"apps/pkgs/paypal_services"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Plan_id *string `json:"plan_id"`
	//
	Custom_id *string `json:"custom_id"`
	// buyer
	Give_name     *string `json:"give_name"`
	Surname       *string `json:"surname"`
	Email_address *string `json:"email_address"`
}

func checkEvent(b []byte, e *Event) (bool, error) {
	err := json.Unmarshal(b, e)
	if err != nil {
		return false, errors.New("json format Unmarshal fail")
	}

	if e.Custom_id == nil {
		return false, errors.New("Custom_id empty")
	} else if e.Plan_id == nil {
		return false, errors.New("Plan_id empty")
	} else if e.Give_name == nil {
		return false, errors.New("Give_name empty")
	} else if e.Surname == nil {
		return false, errors.New("Surname empty")
	} else if e.Email_address == nil {
		return false, errors.New("Email_address empty")
	}
	return true, nil
}

var svc = paypal_services.NewPaypalSVC()

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("\n[request.Body]\n%s\n\n", request.Body)

	e := &Event{}
	if ok, err := checkEvent([]byte(request.Body), e); !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, err
	}

	//
	err := svc.Auth()
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, err
	}
	//
	plan_id := *e.Plan_id
	subscriber := paypal_services.CreateSubscription_Subscriber{
		Name: struct {
			Given_name string `json:"given_name"`
			Surname    string `json:"surname"`
		}{
			Given_name: *e.Give_name,
			Surname:    *e.Surname,
		},
		Email_address: *e.Email_address,
	}
	custom_id := *e.Custom_id
	approval_url, err := svc.Create_subscription(plan_id, subscriber, custom_id)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
		}, err
	}

	//
	type ReturnBody struct {
		Approval_url string `json:"approval_url"`
	}
	returnBody := ReturnBody{
		Approval_url: approval_url,
	}
	b, _ := json.Marshal(returnBody)

	return events.APIGatewayProxyResponse{
		Body:       string(b),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
