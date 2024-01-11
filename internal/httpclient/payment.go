package httpclient

import (
	"context"
	"fmt"

	"github.com/ucok-man/pixelrental/internal/config"
	"github.com/xendit/xendit-go/v4"
	"github.com/xendit/xendit-go/v4/common"
	"github.com/xendit/xendit-go/v4/invoice"
)

type PaymentService struct {
	cfg *config.Config
}

func (s *PaymentService) CreateInvoice(id int, amount float64, method *string) (*invoice.Invoice, error) {
	createInvoiceRequest := *invoice.NewCreateInvoiceRequest(fmt.Sprint(id), amount)
	createInvoiceRequest.SetCurrency("IDR")
	if method != nil {
		createInvoiceRequest.SetPaymentMethods([]string{*method})
	}
	xenditClient := xendit.NewClient(s.cfg.External.Xendit.API_KEY)

	resp, _, err := xenditClient.InvoiceApi.CreateInvoice(context.Background()).
		CreateInvoiceRequest(createInvoiceRequest).
		Execute()

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *PaymentService) GetInvoice(invoiceID string) (*invoice.Invoice, *common.XenditSdkError) {
	xenditClient := xendit.NewClient(s.cfg.External.Xendit.API_KEY)

	resp, _, err := xenditClient.
		InvoiceApi.
		GetInvoiceById(context.Background(), invoiceID).
		Execute()

	if err != nil {
		return nil, err
	}
	return resp, nil
}
