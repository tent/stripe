package stripe

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Credit Card Types accepted by the Stripe API.
const (
	AmericanExpress = "American Express"
	DinersClub      = "Diners Club"
	Discover        = "Discover"
	JCB             = "JCB"
	MasterCard      = "MasterCard"
	Visa            = "Visa"
	UnknownCard     = "Unknown"
)

// Card represents details about a Credit Card entered into Stripe.
type Card struct {
	ID                string `json:"id"`
	Name              string `json:"name,omitempty"`
	Type              string `json:"type"`
	ExpMonth          int    `json:"exp_month"`
	ExpYear           int    `json:"exp_year"`
	Last4             string `json:"last4"`
	Fingerprint       string `json:"fingerprint"`
	Country           string `json:"country,omitempty"`
	Address1          string `json:"address_line1,omitempty"`
	Address2          string `json:"address_line2,omitempty"`
	AddressCountry    string `json:"address_country,omitempty"`
	AddressState      string `json:"address_state,omitempty"`
	AddressZip        string `json:"address_zip,omitempty"`
	AddressLine1Check string `json:"address_line1_check,omitempty"`
	AddressZipCheck   string `json:"address_zip_check,omitempty"`
	CVCCheck          string `json:"cvc_check,omitempty"`
	Customer          string `json:"customer,omitempty"`
}

// CardParams encapsulates options for Creating or Updating Credit Cards.
type CardParams struct {
	// (Optional) Cardholder's full name.
	Name string

	// The card number, as a string without any separators.
	Number string

	// The card's expiration month.
	ExpMonth int

	// The card's expiration year.
	ExpYear int

	// Card security code
	CVC string

	// (Optional) Billing address line 1
	Address1 string

	// (Optional) Billing address line 2
	Address2 string

	// (Optional) Billing address country
	AddressCountry string

	// (Optional) Billing address state
	AddressState string

	// (Optional) Billing address zip code
	AddressZip string
}

type CardClient struct{}

func (c CardClient) path(customerID, cardID string) string {
	p := fmt.Sprintf("/customers/%s/cards", url.QueryEscape(customerID))
	if cardID != "" {
		p += "/" + url.QueryEscape(cardID)
	}
	return p
}

func (c CardClient) Create(customerID, token string, card *CardParams) (*Card, error) {
	params := make(url.Values)
	if token != "" {
		params.Add("card", token)
	} else {
		appendCardParams(params, card)
	}
	res := &Card{}
	return res, query("POST", c.path(customerID, ""), params, res)
}

func (c CardClient) Update(customerID, cardID string, card *CardParams) (*Card, error) {
	params := make(url.Values)
	appendCardParams(params, card)
	res := &Card{}
	return res, query("POST", c.path(customerID, cardID), params, res)
}

func (c CardClient) Delete(customerID, cardID string) (bool, error) {
	res := &DeleteResp{}
	err := query("DELETE", c.path(customerID, cardID), nil, res)
	return res.Deleted, err
}

func (c CardClient) Retrieve(customerID, cardID string) (*Card, error) {
	res := &Card{}
	return res, query("GET", c.path(customerID, cardID), nil, res)
}

func (c CardClient) List(customerID string, limit int, before, after string) ([]*Card, bool, error) {
	res := struct {
		ListObject
		Data []*Card
	}{}
	err := query("GET", c.path(customerID, ""), listParams(limit, before, after), &res)
	return res.Data, res.More, err
}

// IsLuhnValid uses the Luhn Algorithm (also known as the Mod 10 algorithm) to
// verify a credit cards checksum, which helps flag accidental data entry
// errors.
//
// see http://en.wikipedia.org/wiki/Luhn_algorithm
func IsLuhnValid(card string) (bool, error) {
	var sum = 0
	var digits = strings.Split(card, "")

	// iterate through the digits in reverse order
	for i, even := len(digits)-1, false; i >= 0; i, even = i-1, !even {
		// convert the digit to an integer
		digit, err := strconv.Atoi(digits[i])
		if err != nil {
			return false, err
		}

		// we multiply every other digit by 2, adding the product to the sum.
		// note: if the product is double digits (i.e. 14) we add the two digits
		//       to the sum (14 -> 1+4 = 5). A simple shortcut is to subtract 9
		//       from a double digit product (14 -> 14 - 9 = 5).
		switch {
		case even && digit > 4:
			sum += (digit * 2) - 9
		case even:
			sum += digit * 2
		case !even:
			sum += digit
		}
	}

	// if the sum is divisible by 10, it passes the check
	return sum%10 == 0, nil
}

// GetCardType is a simple algorithm to determine the Card Type (ie Visa,
// Discover) based on the Credit Card Number. If the Number is not recognized, a
// value of "Unknown" will be returned.
func GetCardType(card string) string {
	switch card[0:1] {
	case "4":
		return Visa
	case "2", "1":
		switch card[0:4] {
		case "2131", "1800":
			return JCB
		}
	case "6":
		switch card[0:4] {
		case "6011":
			return Discover
		}
	case "5":
		switch card[0:2] {
		case "51", "52", "53", "54", "55":
			return MasterCard
		}
	case "3":
		switch card[0:2] {
		case "34", "37":
			return AmericanExpress
		case "36":
			return DinersClub
		case "30":
			switch card[0:3] {
			case "300", "301", "302", "303", "304", "305":
				return DinersClub
			}
		default:
			return JCB
		}
	}

	return UnknownCard
}
