package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ucok-man/pixelrental/api"
	xendit "github.com/xendit/xendit-go/v4"
	invoice "github.com/xendit/xendit-go/v4/invoice"
)

func main() {
	if err := api.New().Serve(); err != nil {
		log.Fatal(err)
	}

	// createinvoice()
	// getinvoice()

}

func createinvoice() {
	customer := invoice.NewCustomerObject()
	customer.SetEmail("custom@email.com")
	customer.SetGivenNames("Ali")
	customer.SetSurname("Baba")

	createInvoiceRequest := *invoice.NewCreateInvoiceRequest("ExternalId_example", float64(123)) // [REQUIRED] | CreateInvoiceRequest
	createInvoiceRequest.SetCurrency("IDR")
	// createInvoiceRequest.SetCustomer(*customer)

	// Business ID of the sub-account merchant (XP feature)
	// forUserId := "62efe4c33e45694d63f585f0" // [OPTIONAL] | string

	xenditClient := xendit.NewClient("xnd_development_J9kh4VwOMvHRgxiknN5tiCb5tVyHOO3OaOKQm6gkhtuqjova7nUAPxsoCxXiFRS")

	resp, _, err := xenditClient.InvoiceApi.CreateInvoice(context.Background()).
		CreateInvoiceRequest(createInvoiceRequest).
		// ForUserId(forUserId). // [OPTIONAL]
		Execute()

	if err != nil {
		// fmt.Fprintf(os.Stderr, "Error when calling `InvoiceApi.CreateInvoice``: %v\n", err.Error())

		b, _ := json.MarshalIndent(err.FullError(), "", "  ")
		fmt.Fprintf(os.Stderr, "Full Error Struct: %v\n", string(b))

		// fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateInvoice`: Invoice
	// fmt.Fprintf(os.Stdout, "Response from `InvoiceApi.CreateInvoice`: %v\n", resp)

	b, errjs := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatal(errjs)
	}
	fmt.Println(string(b))
}

func getinvoice() {
	// Invoice ID
	invoiceId := "659f6548fc41bb2389b9c9e5" // [REQUIRED] | string

	// Business ID of the sub-account merchant (XP feature)
	// forUserId := "62efe4c33e45694d63f585f0" // [OPTIONAL] | string

	xenditClient := xendit.NewClient("xnd_development_J9kh4VwOMvHRgxiknN5tiCb5tVyHOO3OaOKQm6gkhtuqjova7nUAPxsoCxXiFRS")

	resp, _, err := xenditClient.InvoiceApi.GetInvoiceById(context.Background(), invoiceId).
		// ForUserId(forUserId). // [OPTIONAL]
		Execute()

	if err != nil {
		// fmt.Fprintf(os.Stderr, "Error when calling `InvoiceApi.GetInvoiceById``: %v\n", err.Error())

		b, _ := json.Marshal(err.FullError())
		fmt.Fprintf(os.Stderr, "Full Error Struct: %v\n", string(b))

		// fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetInvoiceById`: Invoice
	// fmt.Fprintf(os.Stdout, "Response from `InvoiceApi.GetInvoiceById`: %v\n", resp)
	b, errjs := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatal(errjs)
	}
	fmt.Println(string(b))
}
