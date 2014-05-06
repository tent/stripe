package stripe

import (
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

// Subscribes a customer to a new plan.
//
// see https://stripe.com/docs/api#update_subscription
func (SubscriptionClient) Update(customerID string, params *SubscriptionParams) (*Subscription, error) {
	values := url.Values{"plan": {params.Plan}}

	// set optional parameters
	if len(params.Coupon) != 0 {
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
	// attach a new card, if requested
	if len(params.Token) != 0 {
		values.Add("card", params.Token)
	} else if params.Card != nil {
		appendCardParams(values, params.Card)
	}

	s := Subscription{}
	path := "/customers/" + url.QueryEscape(customerID) + "/subscription"
	err := query("POST", path, values, &s)
	return &s, err
}

// Cancels the customer's subscription if it exists.  It cancels the
// subscription immediately.
//
// see https://stripe.com/docs/api#cancel_subscription
func (SubscriptionClient) Cancel(customerID string) (*Subscription, error) {
	s := Subscription{}
	path := "/customers/" + url.QueryEscape(customerID) + "/subscription"
	err := query("DELETE", path, nil, &s)
	return &s, err
}

// Cancels the customer's subscription at the end of the billing period.
//
// see https://stripe.com/docs/api#cancel_subscription
func (SubscriptionClient) CancelAtPeriodEnd(customerID string) (*Subscription, error) {
	values := url.Values{}
	values.Add("at_period_end", "true")

	s := Subscription{}
	path := "/customers/" + url.QueryEscape(customerID) + "/subscription"
	err := query("DELETE", path, values, &s)
	return &s, err
}
