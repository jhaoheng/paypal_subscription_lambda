package paypal_services

func check_reponse_error_handle(status_code int, b []byte) {
	switch status_code {
	case 400:
		// Bad Request
	case 401:
		// Unauthorized
	case 403:
		// Forbidden
	case 404:
		// Not Found
	case 405:
		// Method Not Allowed
	case 406:
		// Not Acceptable
	case 415:
		// Unsupported Media Type
	case 422:
		// Unprocessable Entity
	case 429:
		// RATE_LIMIT_REACHED
	}
}
