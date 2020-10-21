package paypal_services

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (svc *PaypalSVC) List_plans() (plans string, err error) {
	//
	url := fmt.Sprintf("%s%s", Base_url, url_list_plans)
	contentType := "application/json"
	//
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", svc.BearerToken)
	req.Header.Set("Prefer", "return=representation")
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
	plans = string(respbody)
	return plans, nil
}
