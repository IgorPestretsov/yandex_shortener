package app

import (
	"testing"
	"unicode/utf8"
)

func TestGenerateShortLink(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{{
		name: "Short link generated",
		want: seq_length,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utf8.RuneCountInString(GenerateShortLink()); got != tt.want {
				t.Errorf("GenerateShortLink() = %v, want %v", got, tt.want)
			}
		})
	}
}
