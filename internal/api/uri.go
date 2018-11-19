package main

import "net/url"

type URI string

func (u URI) Validate() error {
	_, err := url.Parse(string(u))
	return err
}

type URIs []string

func (uris URIs) Contains(s string) bool {
	for _, u := range uris {
		if u == s {
			return true
		}
	}
	return false
}
