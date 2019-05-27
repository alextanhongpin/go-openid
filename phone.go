package openid

// Phone represents the fields for the scope phone.
type Phone struct {
	PhoneNumber         string `json:"phone_number,omitempty"` // Preferred telephone number.
	PhoneNumberVerified bool   `json:"phone_number_verified"`  // True if the phone number has been verified; otherwise false.
}
