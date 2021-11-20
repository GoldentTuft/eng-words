package main

import (
	"reflect"
	"testing"
)

func TestReadWords(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"a", []string{"a"}},
		{"a b", []string{"a", "b"}},
		{"aa b", []string{"aa", "b"}},
		{"a bb", []string{"a", "bb"}},
		{"aa bb", []string{"aa", "bb"}},
		{"\"a\"b", []string{"a", "b"}},
		{"s “abc”", []string{"s", "abc"}},
	}
	for _, tc := range cases {
		got := readWords(tc.in)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("in %v, want %v, got %v", tc.in, tc.want, got)
		}
	}
}
