package openid

import "strings"

// ResponseType represents the enum for response type.
type ResponseType int

const (
	ResponseTypeNone ResponseType = 1 << iota
	ResponseTypeCode
	ResponseTypeIDToken
	ResponseTypeToken
)

// Has returns true if the response type belong to one of the enum.
func (r ResponseType) Has(rr ResponseType) bool {
	return r&rr != 0
}

// Is returns true if the enum matches exactly one of the enum.
func (r ResponseType) Is(rr ResponseType) bool {
	return r&rr == r|rr
}

var responsetypemap = map[string]ResponseType{
	"code":     ResponseTypeCode,
	"id_token": ResponseTypeIDToken,
	"token":    ResponseTypeToken,
}

func NewResponseType(responseType string) (i ResponseType) {
	rs := strings.Split(responseType, " ")
	for _, r := range rs {
		if v, exist := responsetypemap[r]; exist {
			i |= v
		}
	}
	if i == 0 {
		i |= ResponseTypeNone
	}
	return
}
