package paypal_services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type CreateDraftIncovice struct {
	Detail             CreateDraftInvoice_Detail               `json:"detail"`
	Invoicer           CreateDraftIncovice_Invoicer            `json:"invoicer"`
	Primary_recipients []CreateDraftIncovice_Primary_recipient `json:"primary_recipients"`
	Items              []CreateDraftIncovice_Item              `json:"items"`
}

type CreateDraftInvoice_Detail struct {
	Reference            string `json:"reference,omitempty"`
	Currency_code        string `json:"currency_code"` // 必要, https://developer.paypal.com/docs/api/reference/currency-codes/
	Note                 string `json:"note,omitempty"`
	Terms_and_conditions string `json:"terms_and_conditions,omitempty"`
	Memo                 string `json:"memo,omitempty"`
	Invoice_number       string `json:"invoice_number"` // 必要
	Invoice_date         string `json:"invoice_date"`   // 必要, 必須是 UTC, 指定要寄送的日期, ex : 2018-11-22
}

// 商家資訊
type CreateDraftIncovice_Invoicer struct {
	Business_name string `json:"business_name"` // 必要
	Name          struct {
		Given_name string `json:"given_name,omitempty"`
		Surname    string `json:"surname,omitempty"`
	} `json:"name,omitempty"`
	Address struct {
		Address_line_1 string `json:"address_line_1,omitempty"`
		Address_line_2 string `json:"address_line_2,omitempty"`
		Admin_area_2   string `json:"admin_area_2,omitempty"`
		Admin_area_1   string `json:"admin_area_1,omitempty"`
		Postal_code    string `json:"postal_code,omitempty"`
		Country_code   string `json:"country_code,omitempty"` // TW
	} `json:"address,omitempty"`
	Email_address string `json:"email_address"` // 必要, 必須是該 paypal 下有效的商家 email
	Phones        []struct {
		Country_code    string `json:"country_code"`    // 必要,
		National_number string `json:"national_number"` // 必要
		Phone_type      string `json:"phone_type,omitempty"`
	} `json:"phones,omitempty"`
	Website          string `json:"website,omitempty"`
	Tax_id           string `json:"tax_id,omitempty"`
	Additional_notes string `json:"additional_notes,omitempty"` // 可寫入營業時間
	Logo_url         string `json:"logo_url,omitempty"`
}

// 買家資訊
type CreateDraftIncovice_Primary_recipient struct {
	Billing_info struct {
		Business_name string `json:"business_name"` // 必要
		Name          struct {
			Given_name string `json:"given_name,omitempty"`
			Surname    string `json:"surname,omitempty"`
		} `json:"name,omitempty"`
		Email_address   string `json:"email_address,omitempty"` // 會寄送 invoice 到此 email, 如果忽略, 則不會寄送
		Additional_info string `json:"additional_info,omitempty"`
		Language        string `json:"language"` // 寄送 email 內的訊息語言
	} `json:"billing_info"`
}

// 商品資訊
type CreateDraftIncovice_Item struct {
	Name        string `json:"name"` // 必要
	Description string `json:"description,omitempty"`
	Quantity    string `json:"quantity"` // 必要
	Unit_amount struct {
		Currency_code string `json:"currency_code"` // 必要
		Value         string `json:"value"`         // 必要
	} `json:"unit_amount"`
	Tax struct {
		Name    string `json:"name,omitempty"`
		Percent string `json:"percent,omitempty"`
	} `json:"tax,omitempty"`
	Item_date string `json:"item_date"` // 服務何時可使用, UTC
}

// 執行建立 invoice draft 的 api
func (svc *PaypalSVC) Create_draft_incovice(cdi CreateDraftIncovice) (invoice_href, invoice_id string, err error) {
	b, err := json.Marshal(cdi)
	if err != nil {
		return "", "", err
	}
	request_url := fmt.Sprintf("%s%s", Base_url, url_create_draft_invoice)
	contentType := "application/json"
	var body io.Reader = bytes.NewBuffer(b)
	//
	req, err := http.NewRequest("POST", request_url, body)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", svc.BearerToken)
	//
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		respbody, _ := ioutil.ReadAll(resp.Body)
		text := fmt.Sprintf("Error: %s", string(respbody))
		return "", "", errors.New(text)
	}

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	r_obj := RespCreateDraftIncovice{}
	err = json.Unmarshal(respbody, &r_obj)
	if err != nil {
		return "", "", err
	}
	invoice_href = r_obj.Href
	//
	myUrl, err := url.Parse(invoice_href)
	if err != nil {
		return "", "", err
	}
	invoice_id = path.Base(myUrl.Path)
	return invoice_href, invoice_id, nil
}

type RespCreateDraftIncovice struct {
	Rel    string `json:"rel"`
	Href   string `json:"href"`
	Method string `json:"method"`
}
