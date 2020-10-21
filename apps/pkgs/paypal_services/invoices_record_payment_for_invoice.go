package paypal_services

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func (svc *PaypalSVC) Record_payment_for_invoice(invoice_id string) (err error) {
	b := []byte(`{
		"method": "PAYPAL"
	  }`)
	//
	url := fmt.Sprintf("%s%s", Base_url, url_record_payment_for_invoice)
	url = strings.Replace(url, "{{invoice_id}}", invoice_id, 1)
	contentType := "application/json"
	var body io.Reader = bytes.NewBuffer(b)
	//
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", svc.BearerToken)
	//
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respbody, _ := ioutil.ReadAll(resp.Body)
		text := fmt.Sprintf("Error: %s", string(respbody))
		return errors.New(text)
	}
	return nil
}
