package models

// Customer represents a bank customer record.
type Customer struct {
	PartyID      int    `json:"partyId"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	City         string `json:"city"`
	Postcode     string `json:"postcode"`
	Country      string `json:"country"`
}
