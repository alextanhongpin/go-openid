package main

import "fmt"

type Error struct {
	Message     string `json:"error_message"`
	Description string `json:"error_description"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Description)
}
