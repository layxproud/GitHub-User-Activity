package main

import (
	"slices"
	"testing"
)

func TestCollapse(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{"empty", nil, nil},
		{"single", []string{"a"}, []string{"a"}},
		{"no dups", []string{"a", "b"}, []string{"a", "b"}},
		{"pair at end", []string{"a", "b", "b"}, []string{"a", "b x2"}},
		{"two pairs", []string{"a", "a", "b", "b"}, []string{"a x2", "b x2"}},
		{"triple", []string{"a", "a", "a"}, []string{"a x3"}},
		{"non-adjacent dups", []string{"a", "b", "a"}, []string{"a", "b", "a"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collapse(tt.input)
			if !slices.Equal(got, tt.want) {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
