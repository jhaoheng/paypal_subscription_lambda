package invoice

import "apps/pkgs/paypal_services"

func Get_subsctiption_detail(subscription_id string) (resp *paypal_services.RespShowSubscriptionDetails, err error) {
	resp_detail, err := svc.Show_subscription_details(subscription_id)
	if err != nil {
		return nil, err
	}
	return resp_detail, nil
}
