package stripe

import (
	"net/url"
	"strconv"
)

// Customer encapsulates details about a Customer registered in Stripe.
//
// see https://stripe.com/docs/api#customer_object
type Customer struct {
	ID            string
	Description   string
	Email         string
	Created       UnixTime
	Balance       int `json:"account_balance"`
	Currency      string
	Delinquent    bool
	Cards         *CardList
	Discount      *Discount
	Subscriptions *SubscriptionList
	Livemode      bool
	DefaultCard   string `json:"default_card"`
	Metadata      map[string]string
}

type ListObject struct {
	Count int  `json:"total_count"`
	More  bool `json:"has_more"`
}

type SubscriptionList struct {
	ListObject
	Data []*Subscription
}

type CardList struct {
	ListObject
	Data []*Card
}

// Discount represents the actual application of a coupon to a particular
// customer.
//
// see https://stripe.com/docs/api#discount_object
type Discount struct {
	Customer     string
	Start        UnixTime
	End          UnixTime
	Coupon       *Coupon
	Subscription string
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
func (c *CustomerClient) Create(cust *CustomerParams) (*Customer, error) {
	customer := Customer{}
	params := make(url.Values)
	appendCustomerParams(params, cust)

	err := query("POST", "/customers", params, &customer)
	return &customer, err
}

// Retrieves a Customer with the given ID.
//
// see https://stripe.com/docs/api#retrieve_customer
func (c *CustomerClient) Retrieve(id string) (*Customer, error) {
	customer := Customer{}
	path := "/customers/" + url.QueryEscape(id)
	err := query("GET", path, nil, &customer)
	return &customer, err
}

// Updates a Customer with the given ID.
//
// see https://stripe.com/docs/api#update_customer
func (c *CustomerClient) Update(id string, cust *CustomerParams) (*Customer, error) {
	customer := Customer{}
	params := make(url.Values)
	appendCustomerParams(params, cust)

	err := query("POST", "/customers/"+url.QueryEscape(id), params, &customer)
	return &customer, err
}

// Deletes a Customer (permanently) with the given ID.
//
// see https://stripe.com/docs/api#delete_customer
func (c *CustomerClient) Delete(id string) (bool, error) {
	resp := DeleteResp{}
	path := "/customers/" + url.QueryEscape(id)
	if err := query("DELETE", path, nil, &resp); err != nil {
		return false, err
	}
	return resp.Deleted, nil
}

// Returns a list of your Customers.
//
// see https://stripe.com/docs/api#list_customers
func (c *CustomerClient) List() ([]*Customer, error) {
	return c.ListN(10, 0)
}

// Returns a list of your Customers at the specified range.
//
// see https://stripe.com/docs/api#list_customers
func (c *CustomerClient) ListN(count int, offset int) ([]*Customer, error) {
	// define a wrapper function for the Customer List, so that we can
	// cleanly parse the JSON
	type listCustomerResp struct{ Data []*Customer }
	resp := listCustomerResp{}

	// add the count and offset to the list of url values
	values := url.Values{
		"count":  {strconv.Itoa(count)},
		"offset": {strconv.Itoa(offset)},
	}

	err := query("GET", "/customers", values, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
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
		appendCardParams(values, c.Card)
	} else if c.Token != "" {
		values.Add("card", c.Token)
	}
}

func appendCardParams(values url.Values, c *CardParams) {
	if c.Number != "" {
		values.Add("card[number]", c.Number)
	}
	if c.ExpMonth != 0 {
		values.Add("card[exp_month]", strconv.Itoa(c.ExpMonth))
	}
	if c.ExpMonth != 0 {
		values.Add("card[exp_year]", strconv.Itoa(c.ExpYear))
	}
	if c.Name != "" {
		values.Add("card[name]", c.Name)
	}
	if c.CVC != "" {
		values.Add("card[cvc]", c.CVC)
	}
	if c.Address1 != "" {
		values.Add("card[address_line1]", c.Address1)
	}
	if c.Address2 != "" {
		values.Add("card[address_line2]", c.Address2)
	}
	if c.AddressZip != "" {
		values.Add("card[address_zip]", c.AddressZip)
	}
	if c.AddressState != "" {
		values.Add("card[address_state]", c.AddressState)
	}
	if c.AddressCountry != "" {
		values.Add("card[address_country]", c.AddressCountry)
	}
}
