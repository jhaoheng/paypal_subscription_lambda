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

type AppAuth struct {
	Access_token string `json:"access_token"`
	Token_type   string `json:"token_type"`
	App_id       string `json:"app_id"`
	Expires_in   int    `json:"expires_in"`
	Nonce        string `json:"nonce"`
}

func (svc *PaypalSVC) Auth() error {

	// 檢查 token 是否過期, 若過期才去要求 new bearer token
	if !svc.Verify_bearerToken_is_expried() {
		return nil
	} else {
		url := fmt.Sprintf("%s%s", Base_url, url_auth)
		contentType := "application/x-www-form-urlencoded"
		jsonStr := []byte(`grant_type=client_credentials`)
		var body io.Reader = bytes.NewBuffer(jsonStr)

		//
		req, err := http.NewRequest("POST", url, body)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", svc.Authorization)

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
		//
		respbody, _ := ioutil.ReadAll(resp.Body)
		var appAuth = AppAuth{}
		json.Unmarshal(respbody, &appAuth)
		//
		svc.AppAuth = appAuth
		svc.BearerToken = fmt.Sprintf("%s %s", appAuth.Token_type, appAuth.Access_token)
		//
		svc.set_expired_at()
	}
	fmt.Printf("[Paypal Config] %#v\n\n", svc)
	return nil
}

/*
To detect when an access token expires, write code to either:

- Keep track of the expires_in value in the token response.
- Handle the HTTP 401 Unauthorized status code. The API endpoint issues this status code when it detects an expired token.
*/
func (svc *PaypalSVC) set_expired_at() {
	svc.BearerTokenExpiredAt = time.Now().UTC().Unix() + int64(svc.AppAuth.Expires_in)
}
func (svc *PaypalSVC) Verify_bearerToken_is_expried() bool {
	// 現在時間 > 過期時間
	if time.Now().UTC().Unix() >= svc.BearerTokenExpiredAt {
		return true
	}
	return false
}
