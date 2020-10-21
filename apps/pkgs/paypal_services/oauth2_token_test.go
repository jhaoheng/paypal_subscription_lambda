package paypal_services

import (
	"fmt"
	"testing"
)

func Test_auth(t *testing.T) {
	svc := NewPaypalSVC()
	fmt.Println(svc.Auth())
}
