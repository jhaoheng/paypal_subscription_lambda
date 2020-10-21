package paypal_services

type WEBHOOK_EVENT struct {
	Id               string      `json:"id"`
	Create_time      string      `json:"create_time"`
	Resource_type    string      `json:"resource_type"` // 判斷此事件類型即可
	Event_type       string      `json:"event_type"`
	Summary          string      `json:"summary"`
	Resource         interface{} `json:"resource"` // 不太需要這個的全部類型(因不同類型內容不同), 整個事件直接儲存到資料庫
	Links            []LINKS     `json:"links"`
	Resource_version string      `json:"resource_version"`
	Event_version    string      `json:"event_version"`
}

type AMOUNT struct {
	Currency_code string `json:"currency_code"`
	Value         string `json:"value"`
}

type LINKS struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}
