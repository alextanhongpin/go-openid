package crypto

import (
	"testing"
)

func TestGenerateRandomBytes(t *testing.T) {
	inputs := []struct {
		input, expected int
	}{
		{-1, 0},
		{1, 1},
		{10, 10},
		{100, 100},
		{1000, 1000},
	}
	for _, v := range inputs {
		output, _ := GenerateRandomBytes(v.input)
		got := len(output)
		if got != v.expected {
			t.Errorf("expected %v, got %v", v.expected, got)
		}
	}
}

func TestGenerateRandomString(t *testing.T) {
	inputs := []int{-1, 0, 10}
	for _, v := range inputs {
		_, got := GenerateRandomString(v)

		if got != nil {
			t.Errorf("expected %v, got %v", nil, got)
		}
	}
}
