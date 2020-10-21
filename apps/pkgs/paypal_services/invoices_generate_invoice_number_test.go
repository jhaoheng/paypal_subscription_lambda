package paypal_services

import (
	"fmt"
	"testing"
)

func Test_generate_invoices_number(t *testing.T) {
	svc := NewPaypalSVC()
	svc.Auth()
	//
	invoices_number, _ := svc.Generate_invoices_number()
	fmt.Println(invoices_number)
}
