package httpclient

import "github.com/ucok-man/pixelrental/internal/config"

type HTTPClient struct {
	Payment PaymentService
}

func New(cfg *config.Config) *HTTPClient {
	return &HTTPClient{
		Payment: PaymentService{cfg: cfg},
	}
}
