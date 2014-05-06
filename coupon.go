package stripe

import (
	"net/url"
	"strconv"
)

// Coupon Durations
const (
	DurationForever   = "forever"
	DurationOnce      = "once"
	DurationRepeating = "repeating"
)

// Coupon represents percent-off discount you might want to apply to a customer.
//
// see https://stripe.com/docs/api#coupon_object
type Coupon struct {
	ID               string            `json:"id"`
	Duration         string            `json:"duration"`
	AmountOff        int               `json:"amount_off,omitempty"`
	PercentOff       int               `json:"percent_off,omitempty"`
	DurationInMonths int               `json:"duration_in_months,omitempty"`
	MaxRedemptions   int               `json:"max_redemptions,omitempty"`
	RedeemBy         *UnixTime         `json:"redeem_by,omitempty"`
	TimesRedeemed    int               `json:"times_redeemed"`
	Livemode         bool              `json:"livemode"`
	Created          UnixTime          `json:"created"`
	Metadata         map[string]string `json:"metadata"`
	Valid            bool              `json:"valid"`
}

// CouponClient encapsulates operations for creating, updating, deleting and
// querying coupons using the Stripe REST API.
type CouponClient struct{}

// CouponParams encapsulates options for creating a new Coupon.
type CouponParams struct {
	// (Optional) Unique string of your choice that will be used to identify
	// this coupon when applying it a customer.
	ID string

	// A positive integer between 1 and 100 that represents the discount the
	// coupon will apply.
	PercentOff int

	// Specifies how long the discount will be in effect. Can be forever, once,
	// or repeating.
	Duration string

	// A positive integer representing the amount to subtract from an invoice
	// total (required if percent_off is not passed)
	AmountOff int

	// Currency of the amount_off parameter (required if amount_off is passed)
	Currency string

	// (Optional) If duration is repeating, a positive integer that specifies
	// the number of months the discount will be in effect.
	DurationInMonths int

	// (Optional) A positive integer specifying the number of times the coupon
	// can be redeemed before it's no longer valid. For example, you might have
	// a 50% off coupon that the first 20 readers of your blog can use.
	MaxRedemptions int

	// (Optional) UTC timestamp specifying the last time at which the coupon can
	// be redeemed. After the redeem_by date, the coupon can no longer be
	// applied to new customers.
	RedeemBy *UnixTime

	Metadata map[string]string
}

// Creates a new Coupon.
//
// see https://stripe.com/docs/api#create_coupon
func (CouponClient) Create(params *CouponParams) (*Coupon, error) {
	coupon := Coupon{}
	values := url.Values{
		"duration":    {params.Duration},
		"percent_off": {strconv.Itoa(params.PercentOff)},
	}

	if len(params.ID) != 0 {
		values.Add("id", params.ID)
	}
	if params.DurationInMonths != 0 {
		values.Add("duration_in_months", strconv.Itoa(params.DurationInMonths))
	}
	if params.MaxRedemptions != 0 {
		values.Add("max_redemptions", strconv.Itoa(params.MaxRedemptions))
	}

	if params.AmountOff != 0 {
		values.Add("amount_off", strconv.Itoa(params.AmountOff))
		values.Add("currency", params.Currency)
	}
	if params.RedeemBy != nil {
		values.Add("redeem_by", strconv.FormatInt(params.RedeemBy.Unix(), 10))
	}
	appendMetadata(values, params.Metadata)

	err := query("POST", "/coupons", values, &coupon)
	return &coupon, err
}

// Retrieves the coupon with the given ID.
//
// see https://stripe.com/docs/api#retrieve_coupon
func (CouponClient) Retrieve(id string) (*Coupon, error) {
	coupon := Coupon{}
	path := "/coupons/" + url.QueryEscape(id)
	err := query("GET", path, nil, &coupon)
	return &coupon, err
}

// Deletes the coupon with the given ID.
//
// see https://stripe.com/docs/api#delete_coupon
func (CouponClient) Delete(id string) (bool, error) {
	resp := DeleteResp{}
	path := "/coupons/" + url.QueryEscape(id)
	if err := query("DELETE", path, nil, &resp); err != nil {
		return false, err
	}
	return resp.Deleted, nil
}

// Returns a list of your coupons at the specified range.
//
// see https://stripe.com/docs/api#list_coupons
func (CouponClient) List(limit int, before, after string) ([]*Coupon, bool, error) {
	res := struct {
		ListObject
		Data []*Coupon
	}{}
	err := query("GET", "/coupons", listParams(limit, before, after), &res)
	return res.Data, res.More, err
}
