package paypal_services

import (
	"fmt"
	"testing"
)

func Test_create_subscription(t *testing.T) {
	svc := NewPaypalSVC()
	svc.Auth()
	//
	plan_id := "P-5MM560960P6012142L6DGUWY"
	subscriber := CreateSubscription_Subscriber{
		Name: struct {
			Given_name string `json:"given_name"`
			Surname    string `json:"surname"`
		}{
			Given_name: "John",
			Surname:    "Doe",
		},
		Email_address: "buyer@astra.sandbox",
	}
	custom_id := "40d773f5-a4fa-4a32-8cc6-65a1d1fba17c"
	approval_url, _ := svc.Create_subscription(plan_id, subscriber, custom_id)
	fmt.Println(approval_url)
}
