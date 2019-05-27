package openid

// Email represents the fields for the email scope.
type Email struct {
	Email         string `json:"email,omitempty"` // Preferred e-mail address.
	EmailVerified bool   `json:"email_verified"`  // True if the e-mail address has been verified; otherwise false.
}

// VerifyEmailAddress checks with the given string if the email is valid.
func (e *Email) VerifyEmailAddress(email string, required bool) bool {
	return cmpstr(e.Email, email, required)
}
