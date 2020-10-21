package paypal_services

import "testing"

func Test_Webhook_verify_signature(t *testing.T) {
	svc := NewPaypalSVC()
	svc.Auth()

	//
	w := WebhookVerifySignature{
		Auth_algo:         "SHA256withRSA",
		Cert_url:          "https://api.paypal.com/v1/notifications/certs/CERT-360caa42-fca2a594-5edc0ebc",
		Transmission_id:   "58243060-0e3e-11eb-8282-c109c51d2419",
		Transmission_sig:  "H9EmNtJLnxxT++qHjF6sxKYDZi34BEeXM2p2uoMsH0ipi512e8PlqQHvoI2Bo+dR70J+QBpEyoCZKPoitwtauWzQXkzFoH/Cw47lYRcIf8mHno+VUjtRM192atnZDt0en032Ejglgc5KlVO3WGl+SjxGE+ISslYOq1NA1LUYE8Dw5vIKp68MDqTgdsSP1klL3t7VCtEyaC9ocavc8TjDUE66Hi3COA/pb5d8CbS0Uzog8gIxeAI0yMnqrc86Df2JoIekZu5utnxBqh2kKQX2fpMMn3/BocCzDnpTyGbzCYHJq5E2HZXI3aq+x63Asvfmb9Uink6+1x4+hbuWYQxAjQ==",
		Transmission_time: "2020-10-14T16:57:29Z",
		Webhook_event:     "...",
	}
	svc.Webhook_verify_signature(w)
}
