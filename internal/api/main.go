package main

import (
	"fmt"
	"net/url"
	"strings"
)

func main() {
	fmt.Println("Hello, playground")
}

func stringIsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func ValidateURI(uri string) error {
	_, err := url.Parse(uri)
	return err
}

// type Controller struct {
//         service service.OpenIDConnect
// }
//
// func (c *Controller) Authenticate(w http.ResponseWriter, r *http.Request) {
//         // ParseRequest
//         req := ParseRequest()
//         ctx := r.Context()
//
//         // Call Service
//         res, err := c.service.Authenticate(ctx, req)
//
//         // In Service:
//         // - Validate Required Fields (Pre-Work)
//         // - Validate Request
//         // - DoWorkWithRequest (Work)
//         // - BuildResponse (Post-Work)
//
//         // Log Error
//         // ParseResponse
// }
