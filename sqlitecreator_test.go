package sqlitecreator

import (
	"testing"
)

func TestIsDateFormatValid(t *testing.T) {
	testCases := []struct {
		dateStr   string
		valid     bool
	}{
		{"2023-11-03", true},
		{"2023/11/03", false}, // Invalid date format
		{"03-11-2023", false}, // Invalid date format
		{"2023-20-01", false}, // Invalid date format
		{"2023-02-29", false}, // Invalid date format (not a leap year)
		{"2023-02-2911", false},
		{"20230229", false},
		{"", false},           // Invalid date format (empty string)
	}

	for _, tc := range testCases {
		t.Run(tc.dateStr, func(t *testing.T) {
			valid := isDateFormatValid(tc.dateStr)

			if valid != tc.valid {
				t.Errorf("For date %s, want valid=%t, got valid=%t", tc.dateStr, tc.valid, valid)
			}
		})
	}
}
