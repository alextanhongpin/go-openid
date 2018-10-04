package openid

import (
	"crypto/subtle"
	"sort"
	"strings"
)

func sortstr(s string) string {
	ss := strings.Split(s, " ")
	sort.Strings(ss)
	return strings.Join(ss, " ")
}

func cmpstr(s string, cmp string, required bool) bool {
	if s == "" {
		return !required
	}
	return subtle.ConstantTimeCompare([]byte(s), []byte(cmp)) == 1
}
