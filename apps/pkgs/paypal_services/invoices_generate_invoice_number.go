package paypal_services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (svc *PaypalSVC) Generate_invoices_number() (invoice_number string, err error) {
	url := fmt.Sprintf("%s%s", Base_url, url_generate_invoice_number)
	contentType := "application/json"
	//
	req, err := http.NewRequest("POST", url, nil)
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

	if resp.StatusCode != 200 {
		respbody, _ := ioutil.ReadAll(resp.Body)
		text := fmt.Sprintf("Error: %s", string(respbody))
		return "", errors.New(text)
	}
	respbody, _ := ioutil.ReadAll(resp.Body)
	r_obj := RespGenerateInvoicesNumber{}
	json.Unmarshal(respbody, &r_obj)
	invoice_number = r_obj.Invoice_number
	return invoice_number, nil
}

type RespGenerateInvoicesNumber struct {
	Invoice_number string `json:"invoice_number"`
}
