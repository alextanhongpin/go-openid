package controller

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
)

// M represents simple map interface.
type M map[string]interface{}

// -- helpers

func decodeBase64(in string) (string, error) {
	b, err := base64.URLEncoding.DecodeString(in)
	return string(b), err
}

func encodeBase64(in string) string {
	return base64.URLEncoding.EncodeToString([]byte(in))
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(M{
		"error": err.Error(),
	})
}

func buildURL(uri string, q url.Values) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// getHost tries its best to return the request host.
func getHost(r *http.Request) *url.URL {
	u := r.URL

	// The scheme is http because that's the only protocol your server handles.
	u.Scheme = "http"

	// If client specified a host header, then use it for the full URL.
	u.Host = r.Host

	// Otherwise, use your server's host name.
	if u.Host == "" {
		u.Host = "your-host-name.com"
	}
	// if r.URL.IsAbs() {
	//         host := r.Host
	//         // Slice off any port information.
	//         if i := strings.Index(host, ":"); i != -1 {
	//                 host = host[:i]
	//         }
	//         u.Host = host
	// }

	return u
}
