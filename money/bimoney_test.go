package money

import (
	"strconv"
	"testing"
)

func TestBimoney_PercentageChange(t *testing.T) {
	tests := map[string]struct {
		inputMoneyStrOriginal string
		inputMoneyStrNew      string
		expected              Bimoney
		expectedFmt           string
	}{
		"112.5 => 52": {
			inputMoneyStrOriginal: "112.5",
			inputMoneyStrNew:      "171.00",
			expected:              Bimoney(5200000000),
			expectedFmt:           "52.00000000",
		},
		"2016 => 48.8": {
			inputMoneyStrOriginal: "2016",
			inputMoneyStrNew:      "3000",
			expected:              Bimoney(4880000000),
			expectedFmt:           "48.80000000",
		},
		"2016 => -35.51": {
			inputMoneyStrOriginal: "2016",
			inputMoneyStrNew:      "1300",
			expected:              Bimoney(-3551000000),
			expectedFmt:           "-35.51000000",
		},
		"1 => 49": {
			inputMoneyStrOriginal: "1",
			inputMoneyStrNew:      "50",
			expected:              Bimoney(490000000000),
			expectedFmt:           "4900.00000000",
		},
		"0 => 50": {
			inputMoneyStrOriginal: "0",
			inputMoneyStrNew:      "50",
			expected:              Bimoney(10000000000),
			expectedFmt:           "100.00000000",
		},
		"1.56660000 => -85.25": {
			inputMoneyStrOriginal: "1.56660000",
			inputMoneyStrNew:      "0.231000000",
			expected:              Bimoney(-8525000000),
			expectedFmt:           "-85.25000000",
		},
		"0.26660000 => -13.35": {
			inputMoneyStrOriginal: "0.26660000",
			inputMoneyStrNew:      "0.231000000",
			expected:              Bimoney(-1335000000),
			expectedFmt:           "-13.35000000",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			inputMoneyOriginal, err := ParseBimoney(tt.inputMoneyStrOriginal)
			if err != nil {
				t.Fatal(err)
			}

			inputMoneyNew, err := ParseBimoney(tt.inputMoneyStrNew)
			if err != nil {
				t.Fatal(err)
			}

			got := inputMoneyOriginal.PercentageChange(inputMoneyNew)

			if tt.expected != got {
				t.Fatalf("expected: %d, got: %d", tt.expected, got)
			}

			if tt.expectedFmt != got.FormatBimoney(false) {
				t.Fatalf("expected: %s, got: %s", tt.expectedFmt, got.FormatBimoney(false))
			}
		})
	}
}

func TestBimoney_AmountFromPercentage(t *testing.T) {
	tests := map[string]struct {
		inputMoneyStr      string
		inputPercentageInt int
		expected           Bimoney
		expectedFmt        string
	}{
		"0.21874356 (25%)": {
			inputMoneyStr:      "0.21874356",
			inputPercentageInt: 25,
			expected:           Bimoney(5468589),
			expectedFmt:        "0.05468589",
		},
		"112.5 (52%)": {
			inputMoneyStr:      "112.50",
			inputPercentageInt: 52,
			expected:           Bimoney(5850000000),
			expectedFmt:        "58.50000000",
		},
		"215.3589 (150%)": {
			inputMoneyStr:      "215.3589",
			inputPercentageInt: 150,
			expected:           Bimoney(32303835000),
			expectedFmt:        "323.03835000",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			inputMoney, err := ParseBimoney(tt.inputMoneyStr)
			if err != nil {
				t.Fatal(err)
			}

			inputPercentage, err := ParseBimoney(strconv.Itoa(tt.inputPercentageInt))
			if err != nil {
				t.Fatal(err)
			}

			got := inputMoney.AmountFromPercentage(inputPercentage)

			if tt.expected != got {
				t.Fatalf("expected: %d, got: %d", tt.expected, got)
			}

			if tt.expectedFmt != got.FormatBimoney(false) {
				t.Fatalf("expected: %s, got: %s", tt.expectedFmt, got.FormatBimoney(false))
			}
		})
	}
}
