package oidc

import (
	"sort"
	"strings"
)

func sortstr(s string) string {
	ss := strings.Split(s, " ")
	sort.Strings(ss)
	return strings.Join(ss, " ")
}
