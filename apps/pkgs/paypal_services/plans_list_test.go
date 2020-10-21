package paypal_services

import (
	"fmt"
	"testing"
)

func Test_list_plans(t *testing.T) {
	svc := NewPaypalSVC()
	svc.Auth()
	//
	plans, _ := svc.List_plans()
	fmt.Println(plans)
}
