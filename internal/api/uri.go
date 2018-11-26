package main

import "net/url"

// URI represents a URI type.
type URI string

// Validate checks if the URI is valid.
func (u URI) Validate() error {
	_, err := url.Parse(string(u))
	return err
}

// URIs represents a slice or URI.
type URIs []string

// Contains checks if the slice of URI contains the given URI.
func (uris URIs) Contains(s string) bool {
	for _, u := range uris {
		if u == s {
			return true
		}
	}
	return false
}
