package paypal_services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func (svc *PaypalSVC) Send_invoice(invoice_id string) (InvoiceViewUrl string, err error) {

	b := []byte(`{
		"send_to_invoicer": true,
		"send_to_recipient": true
	  }`)
	//
	url := fmt.Sprintf("%s%s", Base_url, url_send_invoice)
	url = strings.Replace(url, "{{invoice_id}}", invoice_id, 1)
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

	if resp.StatusCode != 200 {
		respbody, _ := ioutil.ReadAll(resp.Body)
		text := fmt.Sprintf("Error: %s", string(respbody))
		return "", errors.New(text)
	}
	respbody, _ := ioutil.ReadAll(resp.Body)
	r_obj := RespSendInvoice{}
	json.Unmarshal(respbody, &r_obj)
	InvoiceViewUrl = r_obj.Href
	return InvoiceViewUrl, nil
}

type RespSendInvoice struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}
