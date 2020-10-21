# 參數清單

- Route53HostedZoneId : Route53, HostedZoneId
- DomainName : 決定 Domain Name, 不用額外設定, template 會建立
- DomainSSLCertArn : 取得 ACM Certificate Arn

- paypal app 參數 (參考 Paypal App 設定)
    - PaypalApi : 
        - sandbox : https://api.sandbox.paypal.com
        - prod : https://api.paypal.com
    - PaypalAppClientId     : 
    - PaypalAppSecret       : 
    - PaypalAppWebhookId    : 填寫 Domain Name 後取得
    - PaypalAppBrandName    : 
    - PaypalAppReturnUrl    : 訂單完成返回 url
    - PaypalAppCancelUrl    : 訂單取消返回 url

- SQS
    - VerifiedQueueName : 
    - ErrorQueueName : 
- SNS
    - InternalNotificationName : 
    - EventFinishNotificationName : 


## Paypal App 設定

1. 建立 app : `https://developer.paypal.com/developer/applications`
2. 取得 app 參數

## AWS 前期設定

- 確定 route53 有 domain, ex : `super.com`, 必須取得 HostedZoneId
- 前期需要到 ACM 中, 設定 `paypal.super.com` 的 certification, 需要取得 arn