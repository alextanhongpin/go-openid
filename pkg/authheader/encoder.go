package authheader

import "encoding/base64"

// EncodeBase64 encodes the given arguments into a base64 authorization header
// string.
func EncodeBase64(a, b string) string {
	data := append([]byte(a), append([]byte(sep), []byte(b)...)...)
	return base64.URLEncoding.EncodeToString(data)
}
