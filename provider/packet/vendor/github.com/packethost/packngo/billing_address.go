// Copyright IBM Corp. 2017, 2026

package packngo

type BillingAddress struct {
	StreetAddress string `json:"street_address,omitempty"`
	PostalCode    string `json:"postal_code,omitempty"`
	CountryCode   string `json:"country_code_alpha2,omitempty"`
}
