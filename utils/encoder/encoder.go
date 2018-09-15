package encoder

import (
	"encoding/json"
	"net/http"
)

// JSON returns an interface as json response
func JSON(w http.ResponseWriter, response interface{}, code int) {
	if code == 0 {
		code = http.StatusOK
	}
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Error returns the error
type Error struct {
	Error            int    `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// ErrorJSON returns the error as json response
func ErrorJSON(w http.ResponseWriter, err error, code int) {
	if code == 0 {
		code = http.StatusBadRequest
	}
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Error{
		Error:            code,
		ErrorDescription: err.Error(),
	})
}
