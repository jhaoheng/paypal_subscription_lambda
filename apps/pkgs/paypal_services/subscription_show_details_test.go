package paypal_services

import (
	"fmt"
	"testing"
)

func Test_show_subscription_details(t *testing.T) {
	svc := NewPaypalSVC()
	svc.Auth()
	//
	subscription_id := "I-4EY29CDTFF9P"
	resp_details, err := svc.Show_subscription_details(subscription_id)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp_details)
}
