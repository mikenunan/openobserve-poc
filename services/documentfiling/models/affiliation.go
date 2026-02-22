package models

// Affiliation represents an association between a document and a party (customer).
type Affiliation struct {
	PartyID      int    `json:"partyId"`
	DocumentName string `json:"documentName"`
}
