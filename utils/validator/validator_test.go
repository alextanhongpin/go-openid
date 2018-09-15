package validator

import "testing"

func TestEmptyID(t *testing.T) {
	_, err := ValidateID("")

	expected := errIDRequired.Error()
	got := err.Error()

	if got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}
func TestInvalidID(t *testing.T) {
	_, err := ValidateID("123")

	expected := errIDInvalid.Error()
	got := err.Error()

	if got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}
