package oidc

type AuthenticationErrorCode int

const (
	InteractionRequired AuthenticationErrorCode = iota
	LoginRequired
	AccountSelectionRequired
	ConsentRequired
	InvalidRequestURI
	InvalidRequestObject
	RequestNotSupported
	RequestURINotSupported
	RegistrationNotSupported
)

var authenticationErrorCode = map[AuthenticationErrorCode]string{
	InteractionRequired:      "interaction_required",
	LoginRequired:            "login_required",
	AccountSelectionRequired: "account_selection_required",
	ConsentRequired:          "consent_required",
	InvalidRequestURI:        "invalid_request_uri",
	InvalidRequestObject:     "invalid_request_object",
	RequestNotSupported:      "request_not_supported",
	RequestURINotSupported:   "request_uri_not_supported",
	RegistrationNotSupported: "registration_not_supported",
}

func (a AuthenticationErrorCode) String() string {
	msg, ok := authenticationErrorCode[a]
	if !ok {
		return "error_code_missing"
	}
	return msg
}

func (a AuthenticationErrorCode) JSON(msg string) error {
	return &ErrorJSON{
		Code:        a.String(),
		Description: msg,
	}
}
