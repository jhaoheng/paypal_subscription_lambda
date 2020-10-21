package paypal_services

import (
	"encoding/base64"
	"fmt"
	"os"
)

const (
	url_auth                     = "/v1/oauth2/token"
	url_webhook_verify_signature = "/v1/notifications/verify-webhook-signature"
	//
	url_create_subscription       = "/v1/billing/subscriptions"
	url_show_subscription_details = "/v1/billing/subscriptions/{{subscription_id}}"
	//
	url_list_plans = "/v1/billing/plans?total_required=true"
	//invoices
	url_generate_invoice_number    = "/v2/invoicing/generate-next-invoice-number"
	url_create_draft_invoice       = "/v2/invoicing/invoices"
	url_send_invoice               = "/v2/invoicing/invoices/{{invoice_id}}/send"
	url_record_payment_for_invoice = "/v2/invoicing/invoices/{{invoice_id}}/payments"
)

var (
	Base_url  = os.Getenv("Paypal_Api")
	Client_id = os.Getenv("Paypal_App_Client_Id")
	Secret    = os.Getenv("Paypal_App_Secret")
	//
	Webhook_id = os.Getenv("Paypal_App_Webhook_Id")
	Brand_name = os.Getenv("Paypal_App_Brand_Name")
	Return_url = os.Getenv("Paypal_App_Return_Url")
	Cancel_url = os.Getenv("Paypal_App_Cancel_Url")
)

type PaypalSVC struct {
	Client_ID     string
	Secret        string
	Authorization string // Basic Auth
	//
	BearerToken          string // Bearer Token
	BearerTokenExpiredAt int64  // Bearer Token Expired At
	//
	AppAuth
}

func NewPaypalSVC() *PaypalSVC {

	fmt.Println("[paypal env]")
	fmt.Printf("Base_url 	: %s\n", Base_url)
	fmt.Printf("Client_id 	: %s\n", Client_id)
	fmt.Printf("Secret 		: %s\n", Secret)
	fmt.Printf("Webhook_id 	: %s\n", Webhook_id)
	fmt.Printf("Brand_name 	: %s\n", Brand_name)
	fmt.Printf("Cancel_url 	: %s\n", Cancel_url)

	appkey := fmt.Sprintf("%s:%s", Client_id, Secret)
	keyencoded := base64.StdEncoding.EncodeToString([]byte(appkey))
	authorization := fmt.Sprintf("Basic %s", keyencoded)

	return &PaypalSVC{
		Client_ID:     Client_id,
		Secret:        Secret,
		Authorization: authorization,
	}
}
