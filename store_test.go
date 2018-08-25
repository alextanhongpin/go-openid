package openid

import "testing"

func TestCodeStoreNoKey(t *testing.T) {
	s := NewCodeStore()
	expected := s.Get("hello world")
	if expected != nil {
		t.Errorf("code store error: got %v want %v", expected, nil)
	}
}

func TestPutCodeStore(t *testing.T) {
	s := NewCodeStore()

	tests := []struct {
		id, code string
	}{
		{"1", "ABC"},
		{"2", "XYZ"},
		{"1", "MNO"},
	}

	for _, tt := range tests {
		s.Put(tt.id, tt.code)
		if expected := s.Get(tt.id); expected.Code != tt.code {
			t.Errorf("put code store error: got %v want %v", expected.Code, tt.code)
		}
	}
}

func TestDeleteCodeStore(t *testing.T) {
	s := NewCodeStore()
	s.Put("1", "hello world")
	s.Delete("1")

	if expected := s.Get("1"); expected != nil {
		t.Errorf("delete code store error: got %v want %v", expected, nil)
	}
}
