package authheader_test

import (
	"testing"
	"testing/quick"

	"github.com/alextanhongpin/go-openid/pkg/authheader"
)

func TestDecoder(t *testing.T) {
	f := func(a, b string) bool {
		data := authheader.EncodeBase64(a, b)
		c, d, err := authheader.DecodeBase64(data)
		if err != nil {
			return c == d
		}
		return a == c && b == d
	}

	if err := quick.Check(f, nil); err != nil {
		t.Fatal(err)
	}
}
