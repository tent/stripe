package stripe

import (
	"net/url"
	"strconv"
)

// InvoiceItem represents a charge (or credit) that should be applied to the
// customer at the end of a billing cycle.
//
// see https://stripe.com/docs/api#invoiceitem_object
type InvoiceItem struct {
	ID           string            `json:"id"`
	Amount       int               `json:"amount"`
	Currency     string            `json:"currency"`
	Customer     string            `json:"customer"`
	Date         UnixTime          `json:"date"`
	Description  string            `json:"description,omitempty"`
	Invoice      string            `json:"invoice,omitempty"`
	Subscription string            `json:"subscription,omitempty"`
	Proration    bool              `json:"proration"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	Livemode     bool              `json:"livemode"`
}

// InvoiceItemParams encapsulates options for creating a new Invoice Items.
type InvoiceItemParams struct {
	// The ID of the customer who will be billed when this invoice item is
	// billed.
	Customer string

	// The integer amount in cents of the charge to be applied to the upcoming
	// invoice. If you want to apply a credit to the customer's account, pass a
	// negative amount.
	Amount int

	// 3-letter ISO code for currency.
	Currency string

	// (Optional) An arbitrary string which you can attach to the invoice item.
	// The description is displayed in the invoice for easy tracking.
	Description string

	// (Optional) The ID of an existing invoice to add this invoice item to.
	// When left blank, the invoice item will be added to the next upcoming
	// scheduled invoice.
	Invoice string

	// (Optional) The ID of a subscription to add this invoice item to.
	Subscription string

	Metadata map[string]string
}

// InvoiceItemClient encapsulates operations for creating, updating, deleting
// and querying invoices using the Stripe REST API.
type InvoiceItemClient struct{}

// Create adds an arbitrary charge or credit to the customer's upcoming invoice.
//
// see https://stripe.com/docs/api#invoiceitem_object
func (InvoiceItemClient) Create(params *InvoiceItemParams) (*InvoiceItem, error) {
	item := InvoiceItem{}
	values := url.Values{
		"amount":   {strconv.Itoa(params.Amount)},
		"currency": {params.Currency},
		"customer": {params.Customer},
	}

	// add optional parameters
	if params.Description != "" {
		values.Add("description", params.Description)
	}
	if params.Invoice != "" {
		values.Add("invoice", params.Invoice)
	}
	if params.Subscription != "" {
		values.Add("subscription", params.Subscription)
	}
	appendMetadata(values, params.Metadata)

	err := query("POST", "/invoiceitems", values, &item)
	return &item, err
}

// Retrieves the Invoice Item with the given ID.
//
// see https://stripe.com/docs/api#retrieve_invoiceitem
func (InvoiceItemClient) Retrieve(id string) (*InvoiceItem, error) {
	item := InvoiceItem{}
	path := "/invoiceitems/" + url.QueryEscape(id)
	err := query("GET", path, nil, &item)
	return &item, err
}

// Update changes the amount or description of an Invoice Item on an upcoming
// invoice, using the given Invoice Item ID.
//
// see https://stripe.com/docs/api#update_invoiceitem
func (InvoiceItemClient) Update(id string, params *InvoiceItemParams) (*InvoiceItem, error) {
	item := InvoiceItem{}
	values := make(url.Values)

	if params.Description != "" {
		values.Add("description", params.Description)
	}
	if params.Amount != 0 {
		values.Add("invoice", strconv.Itoa(params.Amount))
	}
	appendMetadata(values, params.Metadata)

	err := query("POST", "/invoiceitems/"+url.QueryEscape(id), values, &item)
	return &item, err
}

// Removes an Invoice Item with the given ID.
//
// see https://stripe.com/docs/api#delete_invoiceitem
func (InvoiceItemClient) Delete(id string) (bool, error) {
	resp := DeleteResp{}
	path := "/invoiceitems/" + url.QueryEscape(id)
	if err := query("DELETE", path, nil, &resp); err != nil {
		return false, err
	}
	return resp.Deleted, nil
}

// Returns a list of Invoice Items.
//
// see https://stripe.com/docs/api#list_invoiceitems
func (c InvoiceItemClient) List(limit int, before, after string) ([]*InvoiceItem, error) {
	return c.list("", limit, before, after)
}

// Returns a list of Invoice Items for the specified Customer ID.
//
// see https://stripe.com/docs/api#list_invoiceitems
func (c InvoiceItemClient) CustomerList(id string, limit int, before, after string) ([]*InvoiceItem, error) {
	return c.list(id, limit, before, after)
}

func (InvoiceItemClient) list(id string, limit int, before, after string) ([]*InvoiceItem, error) {
	res := struct{ Data []*InvoiceItem }{}
	params := listParams(limit, before, after)
	if id != "" {
		params.Add("customer", id)
	}
	err := query("GET", "/invoiceitems", params, &res)
	return res.Data, err
}
