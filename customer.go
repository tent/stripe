package stripe

import (
	"net/url"
	"strconv"
)

// Customer encapsulates details about a Customer registered in Stripe.
//
// see https://stripe.com/docs/api#customer_object
type Customer struct {
	ID            string            `json:"id"`
	Description   string            `json:"description,omitempty"`
	Email         string            `json:"email,omitempty"`
	Created       UnixTime          `json:"created"`
	Balance       int               `json:"account_balance,omitempty"`
	Currency      string            `json:"currency"`
	Delinquent    bool              `json:"delinquent,omitempty"`
	Cards         *CardList         `json:"cards,omitempty"`
	Discount      *Discount         `json:"discount,omitempty"`
	Subscriptions *SubscriptionList `json:"subscriptions,omitempty"`
	Livemode      bool              `json:"livemode"`
	DefaultCard   string            `json:"default_card"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

type ListObject struct {
	Count int  `json:"total_count"`
	More  bool `json:"has_more"`
}

type SubscriptionList struct {
	ListObject
	Data []*Subscription `json:"data"`
}

type CardList struct {
	ListObject
	Data []*Card `json:"data"`
}

// Discount represents the actual application of a coupon to a particular
// customer.
//
// see https://stripe.com/docs/api#discount_object
type Discount struct {
	Customer     string    `json:"customer"`
	Start        UnixTime  `json:"start"`
	End          *UnixTime `json:"end,omitempty"`
	Coupon       *Coupon   `json:"coupon"`
	Subscription string    `json:"subscription,omitempty"`
}

// CustomerParams encapsulates options for creating and updating Customers.
type CustomerParams struct {
	// (Optional) The customer's email address.
	Email string

	// (Optional) An arbitrary string which you can attach to a customer object.
	Description string

	// (Optional) Customer's Active Credit Card
	Card *CardParams

	// (Optional) Customer's Active Credid Card, using a Card Token
	Token string

	// (Optional) If you provide a coupon code, the customer will have a
	// discount applied on all recurring charges.
	Coupon string

	// (Optional) The identifier of the plan to subscribe the customer to. If
	// provided, the returned customer object has a 'subscription' attribute
	// describing the state of the customer's subscription.
	Plan string

	// (Optional) The quantity you’d like to apply to the subscription you’re creating.
	Quantity int

	// (Optional) timestamp representing the end of the trial period
	// the customer will get before being charged for the first time.
	TrialEnd *UnixTime

	// (Optional) Customer's account balance. Negative is credit, positive is added to the next invoice.
	Balance *int

	// (Optional) Customer's default card id.
	DefaultCard string

	// (Optional) Metadata.
	Metadata map[string]string
}

// CustomerClient encapsulates operations for creating, updating, deleting and
// querying customers using the Stripe REST API.
type CustomerClient struct{}

// Creates a new Customer.
//
// see https://stripe.com/docs/api#create_customer
func (CustomerClient) Create(cust *CustomerParams) (*Customer, error) {
	customer := Customer{}
	params := make(url.Values)
	appendCustomerParams(params, cust)

	err := query("POST", "/customers", params, &customer)
	return &customer, err
}

// Retrieves a Customer with the given ID.
//
// see https://stripe.com/docs/api#retrieve_customer
func (CustomerClient) Get(id string) (*Customer, error) {
	customer := Customer{}
	path := "/customers/" + url.QueryEscape(id)
	err := query("GET", path, nil, &customer)
	return &customer, err
}

// Updates a Customer with the given ID.
//
// see https://stripe.com/docs/api#update_customer
func (CustomerClient) Update(id string, cust *CustomerParams) (*Customer, error) {
	customer := Customer{}
	params := make(url.Values)
	appendCustomerParams(params, cust)

	err := query("POST", "/customers/"+url.QueryEscape(id), params, &customer)
	return &customer, err
}

// Deletes a Customer (permanently) with the given ID.
//
// see https://stripe.com/docs/api#delete_customer
func (CustomerClient) Delete(id string) (bool, error) {
	resp := DeleteResp{}
	err := query("DELETE", "/customers/"+url.QueryEscape(id), nil, &resp)
	return resp.Deleted, err
}

// Returns a list of your Customers at the specified range.
//
// see https://stripe.com/docs/api#list_customers
func (CustomerClient) List(limit int, before, after string) ([]*Customer, bool, error) {
	res := struct {
		ListObject
		Data []*Customer
	}{}
	err := query("GET", "/customers", listParams(limit, before, after), &res)
	return res.Data, res.More, err
}

////////////////////////////////////////////////////////////////////////////////
// Helper Function(s)

func appendCustomerParams(values url.Values, c *CustomerParams) {
	// add optional parameters, if specified
	if c.Email != "" {
		values.Add("email", c.Email)
	}
	if c.Description != "" {
		values.Add("description", c.Description)
	}
	if c.Coupon != "" {
		values.Add("coupon", c.Coupon)
	}
	if c.Plan != "" {
		values.Add("plan", c.Plan)
	}
	if c.TrialEnd != nil {
		values.Add("trial_end", strconv.FormatInt(c.TrialEnd.Unix(), 10))
	}
	if c.Balance != nil {
		values.Add("account_balance", strconv.Itoa(*c.Balance))
	}
	if c.DefaultCard != "" {
		values.Add("default_card", c.DefaultCard)
	}
	appendMetadata(values, c.Metadata)

	// add optional credit card details, if specified
	if c.Card != nil {
		appendCardParams(values, true, c.Card)
	} else if c.Token != "" {
		values.Add("card", c.Token)
	}
}

func appendCardParams(values url.Values, nested bool, c *CardParams) {
	p := func(s string) string {
		if nested {
			return "card[" + s + "]"
		}
		return s
	}
	if c.Number != "" {
		values.Add(p("number"), c.Number)
	}
	if c.ExpMonth != 0 {
		values.Add(p("exp_month"), strconv.Itoa(c.ExpMonth))
	}
	if c.ExpMonth != 0 {
		values.Add(p("exp_year"), strconv.Itoa(c.ExpYear))
	}
	if c.Name != "" {
		values.Add(p("name"), c.Name)
	}
	if c.CVC != "" {
		values.Add(p("cvc"), c.CVC)
	}
	if c.Address1 != "" {
		values.Add(p("address_line1"), c.Address1)
	}
	if c.Address2 != "" {
		values.Add(p("address_line2"), c.Address2)
	}
	if c.AddressZip != "" {
		values.Add(p("address_zip"), c.AddressZip)
	}
	if c.AddressState != "" {
		values.Add(p("address_state"), c.AddressState)
	}
	if c.AddressCountry != "" {
		values.Add(p("address_country"), c.AddressCountry)
	}
}
