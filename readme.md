# paypal services
```bash
.
├── apps                        <-- Source code for a lambda function
│   ├── apiPlansList            <-- api, 產生 paypal 訂閱計畫列表
│   └── apiSubscriptionCreate   <-- api, 產生 paypal 訂閱制訂單
│   └── webhookReceiver         <-- api, paypal webhook 入口
│   └── webhookVerify           <-- 驗證 paypal webhook
│   └── webhookProcess          <-- 處理 webhook 事件
├── docs                        <-- paypal 環境初始化文件說明
├── local                       <-- 本地化環境參數設定
├── postman                     <-- postman api 參考
└── template.yaml
```

# flow
```bash
- api
    1. 列出所有訂閱方案
    2. 用戶選擇方案, 轉跳到 paypal 付款頁面
    3. 用戶
        - 完成付款, 跳轉回付款完成頁面, 等待系統完成付款確認
        - 付款取消, 跳轉回取消頁面
- webhook
    1. paypal 產生 event 送到 webhookReceiver
    2. webhookReceiver 將事件送到 sns
    3. sns -> webhookVerify, 驗證交易授權
        - 若成功, 則 -> sqs = verified queue
        - 若失敗, 則 -> sqs = error queue
    4. verified queue -> webhookProcess
        - 若事件為 『付款完成』, 則產生 invoice, 且 paypal 會自動寄送 email (invoice 上)
    5. webhookProcess -> sns, 記錄所有事件 (後續可使用 email / slack / db 來記錄相關所有事件)
```

- 此服務實踐 paypal 訂閱機制 + AWS Services(lambda/sqs/sns/cw_logs/api_gateway/route53/acm)
- 前端列表 plans 與後端整合 db, 需要額外整合

# Q

## 開發、測試、部署
- 本地測試 api
    1. `sam build && sam local start-api --env-vars ./local/local.envs.json --parameter-overrides $(cat ./local/local.sam-params)`
    2. 執行 api
        - 取得清單 : `curl -X GET http://127.0.0.1:3000/plans/list`
        - 產生訂單 : 
            1. `payload='{"plan_id":"P-3C9505567S506292EL6GTKTY", "custom_id":"max-custom-idhello", "give_name":"max", "surname":"hu", "email_address":"max@astra.cloud"}'`
            2. `curl -X POST http://127.0.0.1:3000/subscription/create -d "$payload"`
        - 執行 webhook : `curl -X POST http://127.0.0.1:3000/webhook/receiver`

- 本地測試 非 api
    1. `sam build`
    2. lambda : 
        - verify : `sam local invoke "paypal-webhook-verify" -e ./apps/webhookVerify/testEvent/BILLING.SUBSCRIPTION.CREATED.json`
        - process : `sam local invoke "paypal-webhook-process" -e ./apps/webhookProcess/testEvent/PAYMENT.SALE.COMPLETED.json`

- 部署 : `sam deploy -g --config-env prod`

## 如何整合既有 backend?
- 在建立訂閱的時候, 帶入參數 custom_id, 在產生 webhook event 時, 會回傳 custom 參數
    - Maximum length: 127

## 如何檢測 `RATE_LIMIT_REACHED`?
- 目前使用 log 來整合通知

## subscription 會遇到的 webhook 事件
> https://developer.paypal.com/docs/api-basics/notifications/webhooks/event-names/#subscriptions

- 第三方 webhook 檢測工具 : `https://webhook.site/`
- 主要客戶訂閱交易事件
    - BILLING.SUBSCRIPTION.CREATED
    - BILLING.SUBSCRIPTION.ACTIVATED
    - PAYMENT.SALE.COMPLETED
- 其他請參閱連結
- 若收款狀態 = pending/on-hold, 請參考其他文章


## payment status
> https://docs.easydigitaldownloads.com/article/1180-what-do-the-different-payment-statuses-mean

## payment on-hold
> https://www.paypal.com/us/brc/article/funds-availability
> https://www.paypal.com/us/smarthelp/article/my-payment-is-on-hold.-why-is-this-happening-faq3297

- 情境 : 
    - 您是首次賣家
    - 您已經有一段時間沒有賣了
    - 多個客戶提出退款，爭議或拒付
    - 您的銷售方式異常或已更改
    - 您正在銷售風險較高的物品
- PayPal 將保留您的資金多長時間？
    - 如果您完成了訂單但不提供任何運送信息，則只要您的買家未報告交易中的任何問題，您的付款應在21天后可用。

## payment pending
> https://docs.easydigitaldownloads.com/article/1234-what-does-pending-mean-is-there-a-problem#:~:text=According%20to%20our%20documentation%20on,t%20completed%20their%20payment%20yet.

- 情境 : 客戶付款, 但 paypal 顯示 pending

## 關於費用與台灣相關事項
> https://www.paypal.com/tw/webapps/mpp/paypal-fees?locale.x=zh_TW

- 台灣賣家每筆跨國交易手續費 4.4% + $0.30 USD
- 台灣境內交易限制, ex : 賣方台灣帳號、買方台灣帳號
    - https://www.paypal.com/tw/webapps/mpp/system-enhancement-faq?locale.x=zh_TW

```
因應台灣相關法規規範，我們將更新 PayPal 交易平台系統的流程以停止處理台灣境內對境內的付款和收款交易。

如果買賣雙方皆使用台灣註冊的 PayPal 帳戶，將無法進行付款和收款。但你仍然可以使用台灣 PayPal 帳戶接受跨國交易款項，或支付商品或服務的款項給國外賣家。

如造成任何不便之處，敬請見諒。我們會持續改善交易平台系統並提升服務品質以滿足各種跨國交易需求。

PayPal 謹啟
```

## 關於 invoices

- 可客製化 invoices, 建議先查看 api 可用參數, 再更新 code