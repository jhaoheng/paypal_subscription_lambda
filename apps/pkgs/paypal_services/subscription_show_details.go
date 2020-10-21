package paypal_services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (svc *PaypalSVC) Show_subscription_details(subscription_id string) (r *RespShowSubscriptionDetails, err error) {
	//
	url := fmt.Sprintf("%s%s", Base_url, strings.Replace(url_show_subscription_details, "{{subscription_id}}", subscription_id, 1))
	contentType := "application/json"
	//
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", svc.BearerToken)
	//
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//
	if resp.StatusCode != 200 {
		respbody, _ := ioutil.ReadAll(resp.Body)
		text := fmt.Sprintf("Error: %s", string(respbody))
		return nil, errors.New(text)
	}
	respbody, _ := ioutil.ReadAll(resp.Body)
	resp_details := RespShowSubscriptionDetails{}
	json.Unmarshal(respbody, &resp_details)
	return &resp_details, nil
}

type RespShowSubscriptionDetails struct {
	Status             string                                 `json:"status,omitempty"` // APPROVAL_PENDING, APPROVED, ACTIVE, SUSPENDED, CANCELLED, EXPIRED
	Status_change_note string                                 `json:"status_change_note,omitempty"`
	Status_update_time string                                 `json:"status_update_time,omitempty"`
	Id                 string                                 `json:"id,omitempty"`
	Plan_id            string                                 `json:"plan_id,omitempty"`
	Start_time         string                                 `json:"start_time,omitempty"`
	Quantity           string                                 `json:"quantity,omitempty"`
	Shipping_amount    interface{}                            `json:"shipping_amount,omitempty"`
	Subscriber         RespShowSubscriptionDetails_Subscriber `json:"subscriber,omitempty"`
	Billing_info       interface{}                            `json:"billing_info,omitempty"`
	Create_time        string                                 `json:"create_time,omitempty"`
	Update_time        string                                 `json:"update_time,omitempty"`
	Custom_id          string                                 `json:"custom_id,omitempty"` // The custom id for the subscription.
	Plan               interface{}                            `json:"plan,omitempty"`
	Links              interface{}                            `json:"links,omitempty"`
}

type RespShowSubscriptionDetails_Subscriber struct {
	Name struct {
		Given_name string `json:"given_name,omitempty"`
		Surname    string `json:"surname,omitempty"`
	} `json:"name,omitempty"`
	Email_address    string      `json:"email_address,omitempty"`
	Payer_id         string      `json:"payer_id,omitempty"`
	Phone            interface{} `json:"phone,omitempty"`
	Shipping_address interface{} `json:"shipping_address,omitempty"`
	Payment_source   interface{} `json:"payment_source,omitempty"`
}
