package oidc_test

import (
	"errors"
	"strings"
	"testing"
)

var displaysMap = map[string]struct{}{
	"page":  struct{}{},
	"popup": struct{}{},
	"touch": struct{}{},
	"wap":   struct{}{},
}

var displaysSlice = [...]string{"page", "popup", "touch", "wap"}

func maptest(display string) error {
	if _, ok := displaysMap[display]; !ok {
		return errors.New("invalid display")
	}
	return nil
}

func containstest(display string) error {
	if !strings.Contains("page popup touch wap", display) {
		return errors.New("invalid display")
	}
	return nil
}

func rangetestequal(display string) error {
	for _, d := range displaysSlice {
		if strings.EqualFold(d, display) {
			return nil
		}
	}
	return errors.New("invalid display")
}

func rangetest(display string) error {
	for _, d := range displaysSlice {
		if d == display {
			return nil
		}
	}
	return errors.New("invalid display")
}

func benchmarkhelp(fn func(display string) error, display string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		if err := fn(display); err != nil {
			panic(err)
		}
	}
}

func BenchmarkMap(b *testing.B)        { benchmarkhelp(maptest, "touch", b) }
func BenchmarkContains(b *testing.B)   { benchmarkhelp(containstest, "touch", b) }
func BenchmarkRangeEqual(b *testing.B) { benchmarkhelp(rangetestequal, "touch", b) }
func BenchmarkRange(b *testing.B)      { benchmarkhelp(rangetest, "touch", b) }
