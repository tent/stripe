package stripe

import (
	"net/url"
	"strconv"
)

// Invoice represents statements of what a customer owes for a particular
// billing period, including subscriptions, invoice items, and any automatic
// proration adjustments if necessary.
//
// see https://stripe.com/docs/api#invoice_object
type Invoice struct {
	ID              string
	AmountDue       int  `json:"amount_due"`
	AttemptCount    int  `json:"attempt_count"`
	Attempted       bool `json:"attempted"`
	Closed          bool
	Paid            bool
	PeriodEnd       UnixTime `json:"period_end"`
	PeriodStart     UnixTime `json:"period_start"`
	Subtotal        int
	Total           int
	Charge          string
	Customer        string
	Date            UnixTime
	Discount        *Discount
	Lines           *InvoiceLines
	StartingBalance int      `json:"starting_balance"`
	EndingBalance   int      `json:"ending_balance"`
	NextPayment     UnixTime `json:"next_payment_attempt"`
	Livemode        bool
}

// InvoiceLines represents an individual line items that is part of an invoice.
type InvoiceLines struct {
	InvoiceItems  []*InvoiceItem      `json:"invoiceitems"`
	Prorations    []*InvoiceItem      `json:"prorations"`
	Subscriptions []*SubscriptionItem `json:"subscriptions"`
}

type SubscriptionItem struct {
	Amount int
	Period *Period
	Plan   *Plan
}

type Period struct {
	Start UnixTime
	End   UnixTime
}

// InvoiceClient encapsulates operations for querying invoices using the Stripe
// REST API.
type InvoiceClient struct{}

// Retrieves the invoice with the given ID.
//
// see https://stripe.com/docs/api#retrieve_invoice
func (c *InvoiceClient) Retrieve(id string) (*Invoice, error) {
	invoice := Invoice{}
	path := "/invoices/" + url.QueryEscape(id)
	err := query("GET", path, nil, &invoice)
	return &invoice, err
}

// Retrieves the upcoming invoice the given customer ID.
//
// see https://stripe.com/docs/api#retrieve_customer_invoice
func (c *InvoiceClient) RetrieveCustomer(cid string) (*Invoice, error) {
	invoice := Invoice{}
	values := url.Values{"customer": {cid}}
	err := query("GET", "/invoices/upcoming", values, &invoice)
	return &invoice, err
}

// Returns a list of Invoices.
//
// see https://stripe.com/docs/api#list_customer_invoices
func (c *InvoiceClient) List() ([]*Invoice, error) {
	return c.list("", 10, 0)
}

// Returns a list of Invoices at the specified range.
//
// see https://stripe.com/docs/api#list_customer_invoices
func (c *InvoiceClient) ListN(count int, offset int) ([]*Invoice, error) {
	return c.list("", count, offset)
}

// Returns a list of Invoices with the given Customer ID.
//
// see https://stripe.com/docs/api#list_customer_invoices
func (c *InvoiceClient) CustomerList(id string) ([]*Invoice, error) {
	return c.list(id, 10, 0)
}

// Returns a list of Invoices with the given Customer ID, at the specified range.
//
// see https://stripe.com/docs/api#list_customer_invoices
func (c *InvoiceClient) CustomerListN(id string, count int, offset int) ([]*Invoice, error) {
	return c.list(id, count, offset)
}

func (c *InvoiceClient) list(id string, count int, offset int) ([]*Invoice, error) {
	// define a wrapper function for the Invoice List, so that we can
	// cleanly parse the JSON
	type listInvoicesResp struct{ Data []*Invoice }
	resp := listInvoicesResp{}

	// add the count and offset to the list of url values
	values := url.Values{
		"count":  {strconv.Itoa(count)},
		"offset": {strconv.Itoa(offset)},
	}

	// query for customer id, if provided
	if id != "" {
		values.Add("customer", id)
	}

	err := query("GET", "/invoices", values, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
