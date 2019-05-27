package openid

import "time"

// Profile represents the fields for the scope profile.
type Profile struct {
	Birthdate         string `json:"birth_date,omitempty"`         // Birthday.
	FamilyName        string `json:"family_name,omitempty"`        // Surname(s) or first name(s).
	Gender            string `json:"gender,omitempty"`             // Gender.
	GivenName         string `json:"given_name,omitempty"`         // Given name(s) or first name(s).
	Locale            string `json:"locale,omitempty"`             // Locale.
	MiddleName        string `json:"middle_name,omitempty"`        // Middle name(s).
	Name              string `json:"name,omitempty"`               // Full name.
	Nickname          string `json:"nickname,omitempty"`           // Casual name.
	Picture           string `json:"picture,omitempty"`            // Profile picture URL.
	PreferredUsername string `json:"preferred_username,omitempty"` // Shorthand name by which the End-User wishes to be referred to.
	Profile           string `json:"profile,omitempty"`            // Profile page URL.
	UpdatedAt         int64  `json:"updated_at,omitempty"`         // Time the information was last updated.
	ZoneInfo          string `json:"zone_info,omitempty"`          // Time zone.
	Website           string `json:"website,omitempty"`            // Web page or blog URL.
}

// Update updates the profile last updated time.
func (p *Profile) Update() {
	p.UpdatedAt = time.Now().UTC().Unix()
}
