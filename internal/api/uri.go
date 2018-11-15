package main

import "strings"

type URIs []string

func (uris URIs) Contains(s string) bool {
	for _, u := range uris {
		if strings.EqualFold(u, s) {
			return true
		}
	}
	return false
}
