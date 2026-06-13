package service

import "testing"

func TestIsValidLuhn(t *testing.T) {
	tests := []struct {
		number string
		want   bool
	}{
		{"12345678903", true},
		{"9278923470", true},
		{"346436439", true},
		{"2377225624", true},

		{"0", true},
		{"", false},

		{"12345678900", false},
		{"1234567890", false},

		{"123abc", false},
		{"123 456", false},
		{"-1", false},
	}

	for _, tc := range tests {
		got := isValidLuhn(tc.number)
		if got != tc.want {
			t.Errorf("isValidLuhn(%q) = %v, want %v", tc.number, got, tc.want)
		}
	}
}
