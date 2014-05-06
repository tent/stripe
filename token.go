package stripe

import (
	"net/url"
)

// Token represents a unique identifier for a credit card that can be safely
// stored without having to hold sensitive card information on your own servers.
//
// see https://stripe.com/docs/api#token_object
type Token struct {
	ID       string   `json:"id"`
	Card     *Card    `json:"card"`
	Created  UnixTime `json:"created"`
	Used     bool     `json:"used"`
	Livemode bool     `json:"livemode"`
}

// TokenClient encapsulates operations for creating and querying tokens using
// the Stripe REST API.
type TokenClient struct{}

// TokenParams encapsulates options for creating a new Card Token.
type TokenParams struct {
	Card *CardParams
}

// Creates a single use token that wraps the details of a credit card.
// This token can be used in place of a credit card hash with any API method.
// These tokens can only be used once: by creating a new charge object, or
// attaching them to a customer.
//
// see https://stripe.com/docs/api#create_token
func (c *TokenClient) Create(params *TokenParams) (*Token, error) {
	token := &Token{}
	values := make(url.Values)
	appendCardParams(values, params.Card)

	err := query("POST", "/tokens", values, token)
	return token, err
}

// Retrieves the card token with the given ID.
//
// see https://stripe.com/docs/api#retrieve_token
func (c *TokenClient) Retrieve(id string) (*Token, error) {
	token := Token{}
	path := "/tokens/" + url.QueryEscape(id)
	err := query("GET", path, nil, &token)
	return &token, err
}
