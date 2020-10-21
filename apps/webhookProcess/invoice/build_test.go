package invoice

import (
	"apps/pkgs/paypal_services"
	"encoding/json"
	"fmt"
	"testing"
)

func Test_generate_invoice_template(t *testing.T) {
	invoice_unmber := "001"
	cdi := paypal_services.CreateDraftIncovice{}
	generate_invoice_template(&cdi, invoice_unmber)
	b, _ := json.Marshal(cdi)
	fmt.Println(string(b))
}
