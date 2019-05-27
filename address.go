package openid

// Address represents the fields for the address scope.
type Address struct {
	Country       string `json:"country,omitempty"`
	Formatted     string `json:"formatted,omitempty"`
	Locality      string `json:"locality,omitempty"`
	PostalCode    string `json:"postal_code,omitempty"`
	Region        string `json:"region,omitempty"`
	StreetAddress string `json:"street_address,omitempty"`
}
