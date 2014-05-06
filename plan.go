package stripe

import (
	"net/url"
	"strconv"
)

// Plan Intervals
const (
	IntervalMonth = "month"
	IntervalYear  = "year"
)

// Plan holds details about pricing information for different products and
// feature levels on your site. For example, you might have a $10/month plan
// for basic features and a different $20/month plan for premium features.
//
// see https://stripe.com/docs/api#plan_object
type Plan struct {
	ID                   string            `json:"id"`
	Name                 string            `json:"name"`
	Amount               int               `json:"amount"`
	Interval             string            `json:"interval"`
	IntervalCount        int               `json:"interval_count"`
	Currency             string            `json:"currency"`
	TrialPeriodDays      int               `json:"trial_period_days"`
	StatementDescription string            `json:"statement_description,omitempty"`
	Livemode             bool              `json:"livemode"`
	Created              UnixTime          `json:"created"`
	Metadata             map[string]string `json:"metadata"`
}

// PlanClient encapsulates operations for creating, updating, deleting and
// querying plans using the Stripe REST API.
type PlanClient struct{}

// PlanParams encapsulates options for creating a new Plan.
type PlanParams struct {
	// Unique string of your choice that will be used to identify this plan
	// when subscribing a customer.
	ID string

	// A positive integer in cents (or 0 for a free plan) representing how much
	// to charge (on a recurring basis)
	Amount int

	// 3-letter ISO code for currency. Currently, only 'usd' is supported.
	Currency string

	// Specifies billing frequency. Either month or year.
	Interval string

	// The number of intervals between each subscription billing.
	IntervalCount int

	// Name of the plan, to be displayed on invoices and in the web interface.
	Name string

	// (Optional) Specifies a trial period in (an integer number of) days. If
	// you include a trial period, the customer won't be billed for the first
	// time until the trial period ends. If the customer cancels before the
	// trial period is over, she'll never be billed at all.
	TrialPeriodDays int

	// An arbitrary string to be displayed on your customers' credit card
	// statements (alongside your company name) for charges created by this
	// plan.
	StatementDescription *string

	Metadata map[string]string
}

// Creates a new Plan.
//
// see https://stripe.com/docs/api#create_plan
func (PlanClient) Create(params *PlanParams) (*Plan, error) {
	plan := Plan{}
	values := url.Values{
		"id":       {params.ID},
		"name":     {params.Name},
		"amount":   {strconv.Itoa(params.Amount)},
		"interval": {params.Interval},
		"currency": {params.Currency},
	}

	// trial_period_days is optional, add if specified
	if params.TrialPeriodDays != 0 {
		values.Add("trial_period_days", strconv.Itoa(params.TrialPeriodDays))
	}
	if params.IntervalCount > 1 {
		values.Add("interval_count", strconv.Itoa(params.IntervalCount))
	}
	if params.StatementDescription != nil {
		values.Add("statement_description", *params.StatementDescription)
	}
	appendMetadata(values, params.Metadata)

	err := query("POST", "/plans", values, &plan)
	return &plan, err
}

// Retrieves the plan with the given ID.
//
// see https://stripe.com/docs/api#retrieve_plan
func (PlanClient) Retrieve(id string) (*Plan, error) {
	plan := Plan{}
	path := "/plans/" + url.QueryEscape(id)
	err := query("GET", path, nil, &plan)
	return &plan, err
}

// Updates the name of a plan. Other plan details (price, interval, etc.) are,
// by design, not editable.
//
// see https://stripe.com/docs/api#update_plan
func (PlanClient) Update(id string, params *PlanParams) (*Plan, error) {
	values := make(url.Values)
	if params.Name != "" {
		values.Add("name", params.Name)
	}
	if params.StatementDescription != nil {
		values.Add("statement_description", *params.StatementDescription)
	}
	appendMetadata(values, params.Metadata)

	plan := Plan{}
	path := "/plans/" + url.QueryEscape(id)
	err := query("POST", path, values, &plan)
	return &plan, err
}

// Deletes a plan with the given ID.
//
// see https://stripe.com/docs/api#delete_plan
func (PlanClient) Delete(id string) (bool, error) {
	resp := DeleteResp{}
	path := "/plans/" + url.QueryEscape(id)
	if err := query("DELETE", path, nil, &resp); err != nil {
		return false, err
	}
	return resp.Deleted, nil
}

// Returns a list of your Plans.
//
// see https://stripe.com/docs/api#list_Plans
func (PlanClient) List(limit int, before, after string) ([]*Plan, bool, error) {
	res := struct {
		ListObject
		Data []*Plan
	}{}
	err := query("GET", "/plans", listParams(limit, before, after), &res)
	return res.Data, res.More, err
}
