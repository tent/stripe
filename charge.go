package stripe

import (
	"net/url"
	"strconv"
)

// ISO 3-digit Currency Codes for major currencies (not the full list).
const (
	USD = "usd" // US Dollar ($)
	EUR = "eur" // Euro (€)
	GBP = "gbp" // British Pound Sterling (UK£)
	JPY = "jpy" // Japanese Yen (¥)
	CAD = "cad" // Canadian Dollar (CA$)
	HKD = "hkd" // Hong Kong Dollar (HK$)
	CNY = "cny" // Chinese Yuan (CN¥)
	AUD = "aud" // Australian Dollar (A$)
)

// Charge represents details about a credit card charge in Stripe.
//
// see https://stripe.com/docs/api#charge_object
type Charge struct {
	ID                 string            `json:"id"`
	Description        string            `json:"description,omitempty"`
	Amount             int               `json:"amount"`
	Card               *Card             `json:"card"`
	Currency           string            `json:"currency"`
	Created            UnixTime          `json:"created"`
	Customer           string            `json:"customer,omitempty"`
	Invoice            string            `json:"invoice,omitempty"`
	Paid               bool              `json:"paid"`
	Refunded           bool              `json:"refunded,omitempty"`
	AmountRefunded     int               `json:"amount_refunded,omitempty"`
	BalanceTransaction string            `json:"balance_transaction"`
	Dispute            *Dispute          `json:"dispute,omitempty"`
	FailureMessage     string            `json:"failure_message,omitempty"`
	FailureCode        string            `json:"failure_code,omitempty"`
	Metadata           map[string]string `json:"metadata,omitempty"`
	Livemode           bool              `json:"livemode"`
}

type Dispute struct {
	Charge             string    `json:"charge"`
	Livemode           bool      `json:"livemode"`
	Amount             int       `json:"amount"`
	Created            UnixTime  `json:"created"`
	Currency           string    `json:"currency"`
	Reason             string    `json:"reason"`
	Status             string    `json:"status"`
	BalanceTransaction string    `json:"balance_transaction"`
	Evidence           string    `json:"evidence,omitempty"`
	EvidenceDueBy      *UnixTime `json:"evidence_due_by,omitempty"`
	Protected          bool      `json:"is_protected,omitempty"`
}

// ChargeParams encapsulates options for creating a new Charge.
type ChargeParams struct {
	// A positive integer in cents representing how much to charge the card.
	// The minimum amount is 50 cents.
	Amount int

	// 3-letter ISO code for currency. Currently, only 'usd' is supported.
	Currency string

	// (Optional) Either customer or card is required, but not both The ID of an
	// existing customer that will be charged in this request.
	Customer string

	// (Optional) Credit Card that should be charged.
	Card *CardParams

	// (Optional) Credit Card token that should be charged.
	Token string

	// An arbitrary string which you can attach to a charge object. It is
	// displayed when in the web interface alongside the charge. It's often a
	// good idea to use an email address as a description for tracking later.
	Description string

	// Whether or not to immediately capture the charge. Default is true.
	Capture *bool

	// An arbitrary string to be displayed alongside your company name on your
	// customer's credit card statement. This may be up to 15 characters.
	StatementDescription string

	Metadata map[string]string
}

// ChargeClient encapsulates operations for creating, updating, deleting and
// querying charges using the Stripe REST API.
type ChargeClient struct{}

// Creates a new credit card Charge.
//
// see https://stripe.com/docs/api#create_charge
func (ChargeClient) Create(params *ChargeParams) (*Charge, error) {
	charge := Charge{}
	values := url.Values{
		"amount":   {strconv.Itoa(params.Amount)},
		"currency": {params.Currency},
	}

	if params.Description != "" {
		values.Add("description", params.Description)
	}
	if params.Capture != nil && !*params.Capture {
		values.Add("capture", "false")
	}
	if params.StatementDescription != "" {
		values.Add("statement_description", params.StatementDescription)
	}
	appendMetadata(values, params.Metadata)

	// add optional credit card details, if specified
	if params.Card != nil {
		appendCardParams(values, true, params.Card)
	} else if len(params.Token) > 0 {
		values.Add("card", params.Token)
	} else {
		// if no credit card is provide we need to specify the customer
		values.Add("customer", params.Customer)
	}

	err := query("POST", "/charges", values, &charge)
	return &charge, err
}

// Retrieves the details of a charge with the given ID.
//
// see https://stripe.com/docs/api#retrieve_charge
func (ChargeClient) Retrieve(id string) (*Charge, error) {
	charge := Charge{}
	path := "/charges/" + url.QueryEscape(id)
	err := query("GET", path, nil, &charge)
	return &charge, err
}

// Refunds a charge for the full amount.
//
// see https://stripe.com/docs/api#refund_charge
func (ChargeClient) Refund(id string) (*Charge, error) {
	values := url.Values{}
	charge := Charge{}
	path := "/charges/" + url.QueryEscape(id) + "/refund"
	err := query("POST", path, values, &charge)
	return &charge, err
}

// Refunds a charge for the specified amount.
//
// see https://stripe.com/docs/api#refund_charge
func (ChargeClient) RefundAmount(id string, amt int) (*Charge, error) {
	values := url.Values{
		"amount": {strconv.Itoa(amt)},
	}
	charge := Charge{}
	path := "/charges/" + url.QueryEscape(id) + "/refund"
	err := query("POST", path, values, &charge)
	return &charge, err
}

// Returns a list of your Charges with the specified range.
//
// see https://stripe.com/docs/api#list_charges
func (c ChargeClient) List(limit int, before, after string) ([]*Charge, bool, error) {
	return c.list("", limit, before, after)
}

// Returns a list of your Charges with the given Customer ID.
//
// see https://stripe.com/docs/api#list_charges
func (c ChargeClient) CustomerList(id string, limit int, before, after string) ([]*Charge, bool, error) {
	return c.list(id, limit, before, after)
}

func (ChargeClient) list(id string, limit int, before, after string) ([]*Charge, bool, error) {
	res := struct {
		ListObject
		Data []*Charge
	}{}
	params := listParams(limit, before, after)
	if id != "" {
		params.Add("customer", id)
	}
	err := query("GET", "/charges", params, &res)
	return res.Data, res.More, err
}
