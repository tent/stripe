package stripe

import (
	"fmt"
	"net/url"
)

// Invoice represents statements of what a customer owes for a particular
// billing period, including subscriptions, invoice items, and any automatic
// proration adjustments if necessary.
//
// see https://stripe.com/docs/api#invoice_object
type Invoice struct {
	ID                 string            `json:"id"`
	AmountDue          int               `json:"amount_due"`
	AttemptCount       int               `json:"attempt_count"`
	Attempted          bool              `json:"attempted"`
	Closed             bool              `json:"closed"`
	Paid               bool              `json:"paid"`
	PeriodEnd          UnixTime          `json:"period_end"`
	PeriodStart        UnixTime          `json:"period_start"`
	Subtotal           int               `json:"subtotal"`
	Total              int               `json:"total"`
	Currency           string            `json:"currency"`
	Charge             string            `json:"charge,omitempty"`
	Customer           string            `json:"customer"`
	Date               UnixTime          `json:"date"`
	Discount           *Discount         `json:"discount,omitempty"`
	Lines              *InvoiceLines     `json:"lines"`
	StartingBalance    int               `json:"starting_balance"`
	EndingBalance      int               `json:"ending_balance"`
	NextPaymentAttempt *UnixTime         `json:"next_payment_attempt,omitempty"`
	Livemode           bool              `json:"livemode"`
	Metadata           map[string]string `json:"metadata"`
	Description        string            `json:"omitempty"`
}

// InvoiceLines represents an individual line items that is part of an invoice.
type InvoiceLines struct {
	ListObject
	Data []*InvoiceLineItem `json:"data"`
}

type InvoiceLineItem struct {
	ID          string            `json:"id"`
	Livemode    bool              `json:"livemode"`
	Amount      int               `json:"amount"`
	Currency    string            `json:"currency"`
	Period      Period            `json:"period"`
	Proration   bool              `json:"proration"`
	Type        string            `json:"type"`
	Description string            `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata"`
	Plan        *Plan             `json:"plan,omitempty"`
	Quantity    int               `json:"quantity,omitempty"`
}

type Period struct {
	Start UnixTime `json:"start"`
	End   UnixTime `json:"end"`
}

type InvoiceParams struct {
	// The customer ID to invoice
	Customer string

	// (Optional) Invoice description
	Description string

	// (Optional) Invoice metadata
	Metadata map[string]string

	// (Optional) The ID of the subscription to invoice. If not set, the created
	// invoice will include all pending invoice items for the customer.
	Subscription string

	// (Optional) Boolean representing whether an invoice is closed or not.
	Closed *bool
}

// InvoiceClient encapsulates operations for querying invoices using the Stripe
// REST API.
type InvoiceClient struct{}

// Retrieves the invoice with the given ID.
//
// see https://stripe.com/docs/api#retrieve_invoice
func (InvoiceClient) Get(id string) (*Invoice, error) {
	res := &Invoice{}
	return res, query("GET", "/invoices/"+url.QueryEscape(id), nil, res)
}

func (InvoiceClient) Create(params *InvoiceParams) (*Invoice, error) {
	res := &Invoice{}
	return res, query("POST", "/invoices", invoiceValues(params), res)
}

func (InvoiceClient) Update(id string, params *InvoiceParams) (*Invoice, error) {
	res := &Invoice{}
	return res, query("POST", "/invoices/"+url.QueryEscape(id), invoiceValues(params), res)
}

func (InvoiceClient) Pay(id string) (*Invoice, error) {
	res := &Invoice{}
	return res, query("POST", fmt.Sprintf("/invoices/%s/pay", url.QueryEscape(id)), nil, res)
}

// Retrieves the upcoming invoice the given customer ID.
//
// see https://stripe.com/docs/api#retrieve_customer_invoice
func (InvoiceClient) Upcoming(customerID string) (*Invoice, error) {
	res := &Invoice{}
	return res, query("GET", "/invoices/upcoming", url.Values{"customer": {customerID}}, res)
}

// Returns a list of Invoices at the specified range.
//
// see https://stripe.com/docs/api#list_customer_invoices
func (c InvoiceClient) List(limit int, before, after string) ([]*Invoice, bool, error) {
	return c.list("", limit, before, after)
}

// Returns a list of Invoices with the given Customer ID.
//
// see https://stripe.com/docs/api#list_customer_invoices
func (c InvoiceClient) CustomerList(id string, limit int, before, after string) ([]*Invoice, bool, error) {
	return c.list(id, limit, before, after)
}

func (InvoiceClient) list(id string, limit int, before, after string) ([]*Invoice, bool, error) {
	res := struct {
		ListObject
		Data []*Invoice
	}{}
	params := listParams(limit, before, after)
	// query for customer id, if provided
	if id != "" {
		params.Add("customer", id)
	}
	err := query("GET", "/invoices", params, &res)
	return res.Data, res.More, err
}

func invoiceValues(inv *InvoiceParams) url.Values {
	values := make(url.Values)
	if inv.Customer != "" {
		values.Add("customer", inv.Customer)
	}
	if inv.Description != "" {
		values.Add("description", inv.Description)
	}
	if inv.Subscription != "" {
		values.Add("subscription", inv.Subscription)
	}
	if inv.Closed != nil {
		values.Add("closed", fmt.Sprintf("%t", *inv.Closed))
	}
	appendMetadata(values, inv.Metadata)
	return values
}
