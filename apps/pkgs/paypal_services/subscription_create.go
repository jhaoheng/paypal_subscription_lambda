package paypal_services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type CreateSubscription_Subscriber struct {
	Name struct {
		Given_name string `json:"given_name"`
		Surname    string `json:"surname"`
	} `json:"name"`
	Email_address string `json:"email_address"`
}

type CreateSubscription_ApplicationContext struct {
	Brand_name          string `json:"brand_name"`
	Locale              string `json:"locale"`
	Shipping_preference string `json:"shipping_preference"`
	User_action         string `json:"user_action"`
	Payment_method      struct {
		Payer_selected  string `json:"payer_selected"`
		Payee_preferred string `json:"payee_preferred"`
	} `json:"payment_method"`
	Return_url string `json:"return_url"`
	Cancel_url string `json:"cancel_url"`
}

type CreateSubscription struct {
	Plan_id             string                                `json:"plan_id"`
	Start_time          string                                `json:"start_time"` // ex : 2020-10-15T00:00:00Z, default is `Current time`
	Quantity            string                                `json:"quantity"`
	Subscriber          CreateSubscription_Subscriber         `json:"subscriber"`
	Application_context CreateSubscription_ApplicationContext `json:"application_context"`
	Custom_id           string                                `json:"custom_id"` // 自定義 id
}

/*
start_time : UTC, 2006-01-02T15:04:05Z
*/
func (svc *PaypalSVC) Create_subscription(plan_id string, subscriber CreateSubscription_Subscriber, custom_id string) (approve_url string, err error) {
	var start_time string = time.Now().Add(60 * time.Second).UTC().Format("2006-01-02T15:04:05Z")
	//
	application_context := CreateSubscription_ApplicationContext{
		Brand_name:          Brand_name,
		Locale:              "en-US",
		Shipping_preference: "NO_SHIPPING",
		User_action:         "SUBSCRIBE_NOW",
		Payment_method: struct {
			Payer_selected  string `json:"payer_selected"`
			Payee_preferred string `json:"payee_preferred"`
		}{
			Payer_selected:  "PAYPAL",
			Payee_preferred: "IMMEDIATE_PAYMENT_REQUIRED",
		},
		Return_url: Return_url,
		Cancel_url: Cancel_url,
	}
	//
	obj := CreateSubscription{
		Plan_id:             plan_id,
		Start_time:          start_time,
		Quantity:            "1",
		Subscriber:          subscriber,
		Application_context: application_context,
		Custom_id:           custom_id,
	}

	//
	b, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	//
	url := fmt.Sprintf("%s%s", Base_url, url_create_subscription)
	contentType := "application/json"
	var body io.Reader = bytes.NewBuffer(b)
	//
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", svc.BearerToken)
	//
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		respbody, _ := ioutil.ReadAll(resp.Body)
		text := fmt.Sprintf("Error: %s", string(respbody))
		return "", errors.New(text)
	}

	respbody, _ := ioutil.ReadAll(resp.Body)
	r_obj := RespCreateSubscription{}
	json.Unmarshal(respbody, &r_obj)

	links := r_obj.Links
	for _, link := range links {
		if link.Rel == "approve" {
			approve_url = link.Href
			break
		}
	}
	return
}

type RespCreateSubscription struct {
	Id         string `json:"id"`
	Status     string `json:"status"`
	Plan_id    string `json:"plan_id"`
	Start_time string `json:"start_time"`
	Links      []struct {
		Href   string `json:"href"`
		Rel    string `json:"rel"`
		Method string `json:"method"`
	} `json:"links"`
}
