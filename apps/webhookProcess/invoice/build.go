package invoice

import (
	"apps/pkgs/paypal_services"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

// 帳單資訊, 客製化發票, 必要資訊
type Billing_Info struct {
	Business_name string `json:"business_name"`
	Name          struct {
		Given_name string `json:"given_name"`
		Surname    string `json:"surname"`
	} `json:"name"`
	Email_address string `json:"email_address"`
	Language      string `json:"language"`
}

// 購買物品資訊, 客製化發票, 必要資訊
type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Quantity    string `json:"quantity"`
	Unit_amount struct {
		Currency_code string `json:"currency_code"`
		Value         string `json:"value"`
	} `json:"unit_amount"`
	Tax struct {
		Name    string `json:"name"`
		Percent string `json:"percent"`
	} `json:"tax"`
	Item_date string `json:"item_date"` // The date when the item or service was provided, ex : 2020-01-01
}

func Build(billing_info Billing_Info, item Item) (invoice_view_url string, err error) {

	// paypal get auth
	err = svc.Auth()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// 1. 建立 invoice number
	invoice_unmber, err := svc.Generate_invoices_number()
	if err != nil {
		return "", err
	}
	// 2. 建立 invoice 草稿, 取得 invoice_id
	cdi := paypal_services.CreateDraftIncovice{}
	generate_invoice_template(&cdi, invoice_unmber)
	_, invoice_id, err := svc.Create_draft_incovice(cdi)
	if err != nil {
		return "", err
	}
	// 3. 寄送 invoice
	invoice_view_url, err = svc.Send_invoice(invoice_id)
	if err != nil {
		return "", err
	}
	// 4. 將 invoice 標記為已付款
	err = svc.Record_payment_for_invoice(invoice_id)
	if err != nil {
		return "", err
	}
	return invoice_view_url, nil
}

func generate_invoice_template(cdi *paypal_services.CreateDraftIncovice, invoice_unmber string) {
	buffer, err := ioutil.ReadFile("./invoice.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(buffer, &cdi)
	if err != nil {
		panic(err)
	}
	//
	cdi.Detail.Invoice_number = invoice_unmber
	cdi.Detail.Invoice_date = time.Now().UTC().Format("2006-01-02") //2020-01-13
	//
	cdi.Primary_recipients[0].Billing_info.Business_name = "ooxx"
	cdi.Primary_recipients[0].Billing_info.Name.Given_name = "ooxx"
	cdi.Primary_recipients[0].Billing_info.Name.Surname = "ooxx"
	cdi.Primary_recipients[0].Billing_info.Email_address = "ooxx"
	cdi.Primary_recipients[0].Billing_info.Language = "US"
	//
	cdi.Items[0].Name = "ooxx"
	cdi.Items[0].Description = "ooxx"
	cdi.Items[0].Quantity = "1"
	cdi.Items[0].Unit_amount.Currency_code = "USD"
	cdi.Items[0].Unit_amount.Value = "100"
	// cdi.Items[0].Tax.Name = "ooxx"
	// cdi.Items[0].Tax.Percent = "0"
	// fmt.Printf("%#v\n", cdi)
}
