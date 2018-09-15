package encrypt

import (
	"log"
	"testing"
)

func TestHashEmptyString(t *testing.T) {
	_, err := HashPassword("")

	log.Println(err)
	expected := "Password cannot be empty"
	got := err.Error()
	if got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestCompareDifferentPassword(t *testing.T) {

	hashed, _ := HashPassword("123456")
	isSamePassword := ComparePassword("1", hashed)

	expected := false
	got := isSamePassword
	if got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}
func TestHashAndComparePassword(t *testing.T) {
	hashed, _ := HashPassword("123456")
	isSamePassword := ComparePassword("123456", hashed)

	expected := true
	got := isSamePassword
	if got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}
