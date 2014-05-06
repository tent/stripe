package stripe

import (
	"fmt"
	"net/url"
	"strconv"
)

// Subscription Statuses
const (
	SubscriptionTrialing = "trialing"
	SubscriptionActive   = "active"
	SubscriptionPastDue  = "past_due"
	SubscriptionCanceled = "canceled"
	SubscriptionUnpaid   = "unpaid"
)

// Subscriptions represents a recurring charge a customer's card.
//
// see https://stripe.com/docs/api#subscription_object
type Subscription struct {
	ID                 string    `json:"id"`
	Customer           string    `json:"customer"`
	Status             string    `json:"status"`
	Plan               *Plan     `json:"plan"`
	Start              UnixTime  `json:"start"`
	EndedAt            *UnixTime `json:"ended_at,omitempty"`
	CurrentPeriodStart UnixTime  `json:"current_period_start"`
	CurrentPeriodEnd   UnixTime  `json:"current_period_end"`
	TrialStart         *UnixTime `json:"trial_start,omitempty"`
	TrialEnd           *UnixTime `json:"trial_end,omitempty"`
	CanceledAt         *UnixTime `json:"canceled_at,omitempty"`
	CancelAtPeriodEnd  bool      `json:"cancel_at_period_end"`
	Quantity           int       `json:"quantity"`
	Discount           *Discount `json:"discount,omitempty"`
}

// SubscriptionClient encapsulates operations for updating and canceling
// customer subscriptions using the Stripe REST API.
type SubscriptionClient struct{}

// SubscriptionParams encapsulates options for updating a Customer's
// subscription.
type SubscriptionParams struct {
	// The identifier of the plan to subscribe the customer to.
	Plan string

	// (Optional) The code of the coupon to apply to the customer if you would
	// like to apply it at the same time as creating the subscription.
	Coupon string

	// (Optional) Flag telling us whether to prorate switching plans during a
	// billing cycle
	Prorate bool

	// (Optional) UTC integer timestamp representing the end of the trial period
	// the customer will get before being charged for the first time. If set,
	// trial_end will override the default trial period of the plan the customer
	// is being subscribed to.
	TrialEnd *UnixTime

	// (Optional) A new card to attach to the customer.
	Card *CardParams

	// (Optional) A new card Token to attach to the customer.
	Token string

	// (Optional) The quantity you'd like to apply to the subscription you're creating.
	Quantity int
}

func (c SubscriptionClient) path(customerID, subscriptionID string) string {
	p := fmt.Sprintf("/customers/%s/subscriptions", url.QueryEscape(customerID))
	if subscriptionID != "" {
		p += "/" + url.QueryEscape(subscriptionID)
	}
	return p
}

func (c SubscriptionClient) Create(customerID string, params *SubscriptionParams) (*Subscription, error) {
	res := &Subscription{}
	return res, query("POST", c.path(customerID, ""), c.values(params), res)
}

func (c SubscriptionClient) values(params *SubscriptionParams) url.Values {
	values := make(url.Values)
	if params.Plan != "" {
		values.Add("plan", params.Plan)
	}
	if params.Coupon != "" {
		values.Add("coupon", params.Coupon)
	}
	if params.Prorate {
		values.Add("prorate", "true")
	}
	if params.TrialEnd != nil {
		values.Add("trial_end", strconv.FormatInt(params.TrialEnd.Unix(), 10))
	}
	if params.Quantity != 0 {
		values.Add("quantity", strconv.Itoa(params.Quantity))
	}
	if params.Token != "" {
		values.Add("card", params.Token)
	} else if params.Card != nil {
		appendCardParams(values, params.Card)
	}
	return values
}

// Subscribes a customer to a new plan.
//
// see https://stripe.com/docs/api#update_subscription
func (c SubscriptionClient) Update(customerID, subscriptionID string, params *SubscriptionParams) (*Subscription, error) {
	res := &Subscription{}
	return res, query("POST", c.path(customerID, subscriptionID), c.values(params), res)
}

func (c SubscriptionClient) Cancel(customerID, subscriptionID string, atPeriodEnd bool) (*Subscription, error) {
	values := make(url.Values)
	if atPeriodEnd {
		values.Add("at_period_end", "true")
	}
	res := &Subscription{}
	return res, query("DELETE", c.path(customerID, subscriptionID), values, res)
}

func (c SubscriptionClient) Retrieve(customerID, subscriptionID string) (*Subscription, error) {
	res := &Subscription{}
	return res, query("GET", c.path(customerID, subscriptionID), nil, res)
}

func (c SubscriptionClient) List(customerID string, limit int, before, after string) ([]*Subscription, bool, error) {
	res := struct {
		ListObject
		Data []*Subscription
	}{}
	err := query("GET", c.path(customerID, ""), listParams(limit, before, after), res)
	return res.Data, res.More, err
}
