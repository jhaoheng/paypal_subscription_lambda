# 如何手動製作事件
> 因為 delivery 是使用 sns 進行通知, 且必須整合 paypal webhook header (驗證用)

1. 前往 `https://webhook.site/`, copy Your unique URL=> `https://webhook.site/......`
2. 前往 `https://developer.paypal.com/developer/webhooksSimulator/`
  1. 將 webhook url 填入
  2. 設定事件類型

# sns 已整合的物件
- 參考 xxxx.json

# 整合
1. sns format

```
{
  "Records": [
    {
      "Sns": {
        "Message": "{\"header\":\"paypal_header_物件\", \"body\":\"paypal_webhook_事件\"}"
      }
    }
  ]
}
```

2. 將 paypal 的 header 整合到 header

```
{
  "PAYPAL-AUTH-ALGO": "SHA256withRSA",
  "PAYPAL-AUTH-VERSION": "v2",
  "PAYPAL-CERT-URL": "https://api.paypal.com/v1/notifications/certs/CERT-360caa42-fca2a594-5edc0ebc",
  "PAYPAL-TRANSMISSION-ID": "b6558690-0e46-11eb-a35d-895307448f11",
  "PAYPAL-TRANSMISSION-SIG": "JwLDflXZbe1Ux4OjGOD7AztYMi3ohPyu9yg9ieCQT4BGqJyk0UzZbOm1nTgOELaTsP7z+pS1xc1rcU37p3af7gkk7oOH8xe6SrCEYF3ucKXPLiv6geDBURYmWt3weusUjAhyhCfgiaq8UR1po7bgFS+3cj5MDSK/P72uIAZttTyxaJk1O+/Oqwy2NZQHHTlhW3kNBZyU7ecWpt7O6zbt5Mwx16zejwF7ArVFKlpRUoF6I3i5WjHYeY9GIY66gJowOPim9wUctlQkBnBMZPahspeIE6IfKFlc9kwG8yYza3+O5uJCNDnUXWbaew1qv5H2Nyec2HduUczoJiGTwYYkcg==",
  "PAYPAL-TRANSMISSION-TIME": "2020-10-14T17:57:23Z"
}
```

3. 將 paypal webhook 事件, 整合到 body

```
{
  "id": "WH-TY567577T725889R5-1E6T55435R66166TR",
  "create_time": "2018-19-12T22:20:32.000Z",
  "event_type": "BILLING.SUBSCRIPTION.CREATED",
  "event_version": "1.0",
  "resource_type": "subscription",
  "resource_version": "2.0",
  "summary": "A billing subscription was created.",
  "resource": {
    ...
  },
  "links": [
    ...
  ]
}
```