package paypal_services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type WebhookVerifySignature struct {
	// 用戶填寫
	Auth_algo         string `json:"auth_algo"`         // PAYPAL-AUTH-ALGO
	Cert_url          string `json:"cert_url"`          // PAYPAL-CERT-URL
	Transmission_id   string `json:"transmission_id"`   // PAYPAL-TRANSMISSION-ID
	Transmission_sig  string `json:"transmission_sig"`  // PAYPAL-TRANSMISSION-SIG
	Transmission_time string `json:"transmission_time"` // PAYPAL-TRANSMISSION-TIME
	Webhook_event     string `json:"webhook_event"`     // A webhook event notification.

	// 系統自帶
	Webhook_id string `json:"webhook_id"` // The ID of the webhook as configured in your Developer Portal account.
}

func (svc *PaypalSVC) Webhook_verify_signature(w WebhookVerifySignature) (ok bool, err error) {
	w.Webhook_id = Webhook_id

	/* 注意 :
	- 除了 webhook_event, 其他都是字串
	- 保留 webhook_event 型態, 驗證會錯
	*/
	s := fmt.Sprintf(`{"webhook_id":"%s","auth_algo":"%s","cert_url":"%s","transmission_id":"%s","transmission_sig":"%s","transmission_time":"%s", "webhook_event":%s}`,
		w.Webhook_id,
		w.Auth_algo,
		w.Cert_url,
		w.Transmission_id,
		w.Transmission_sig,
		w.Transmission_time,
		w.Webhook_event,
	)
	fmt.Println("request body =>", s)

	//
	url := fmt.Sprintf("%s%s", Base_url, url_webhook_verify_signature)
	contentType := "application/json"
	var body io.Reader = strings.NewReader(s)

	//
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", svc.BearerToken)

	//
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respbody, _ := ioutil.ReadAll(resp.Body)
		text := fmt.Sprintf("Error: %s", string(respbody))
		return false, errors.New(text)
	}
	respbody, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("[Webhook_verify_signature]\nbody : %s\n", string(respbody))
	respWebhookVerifySignature := RespWebhookVerifySignature{}
	json.Unmarshal(respbody, &respWebhookVerifySignature)
	if respWebhookVerifySignature.Verification_status == "FAILURE" {
		return false, errors.New("verification_status = FAILURE")
	}

	return true, nil
}

type RespWebhookVerifySignature struct {
	Verification_status string `json:"verification_status"`
}
