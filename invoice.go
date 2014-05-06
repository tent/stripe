package stripe

import "net/url"

// Invoice represents statements of what a customer owes for a particular
// billing period, including subscriptions, invoice items, and any automatic
// proration adjustments if necessary.
//
// see https://stripe.com/docs/api#invoice_object
type Invoice struct {
	ID                 string        `json:"id"`
	AmountDue          int           `json:"amount_due"`
	AttemptCount       int           `json:"attempt_count"`
	Attempted          bool          `json:"attempted"`
	Closed             bool          `json:"closed"`
	Paid               bool          `json:"paid"`
	PeriodEnd          UnixTime      `json:"period_end"`
	PeriodStart        UnixTime      `json:"period_start"`
	Subtotal           int           `json:"subtotal"`
	Total              int           `json:"total"`
	Currency           string        `json:"currency"`
	Charge             string        `json:"charge,omitempty"`
	Customer           string        `json:"customer"`
	Date               UnixTime      `json:"date"`
	Discount           *Discount     `json:"discount,omitempty"`
	Lines              *InvoiceLines `json:"lines"`
	StartingBalance    int           `json:"starting_balance"`
	EndingBalance      int           `json:"ending_balance"`
	NextPaymentAttempt *UnixTime     `json:"next_payment_attempt,omitempty"`
	ApplicationFee     int           `json:"application_fee,omitempty"`
	Livemode           bool          `json:"livemode"`
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
	Quantity    int               `json:"quantity"`
}

type Period struct {
	Start UnixTime `json:"start"`
	End   UnixTime `json:"end"`
}

// InvoiceClient encapsulates operations for querying invoices using the Stripe
// REST API.
type InvoiceClient struct{}

// Retrieves the invoice with the given ID.
//
// see https://stripe.com/docs/api#retrieve_invoice
func (InvoiceClient) Retrieve(id string) (*Invoice, error) {
	invoice := Invoice{}
	path := "/invoices/" + url.QueryEscape(id)
	err := query("GET", path, nil, &invoice)
	return &invoice, err
}

// Retrieves the upcoming invoice the given customer ID.
//
// see https://stripe.com/docs/api#retrieve_customer_invoice
func (InvoiceClient) RetrieveCustomer(cid string) (*Invoice, error) {
	invoice := Invoice{}
	values := url.Values{"customer": {cid}}
	err := query("GET", "/invoices/upcoming", values, &invoice)
	return &invoice, err
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
