package main

type URIs []string

func (uris URIs) Contains(s string) bool {
	for _, u := range uris {
		if u == s {
			return true
		}
	}
	return false
}
