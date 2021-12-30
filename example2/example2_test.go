package example2

import "testing"

func TestIsEven(t *testing.T) {
	testTable := []struct {
		num      int
		expected bool
	}{
		{
			num:      5,
			expected: true,
		},
		{
			num:      0,
			expected: false,
		},
		{
			num:      -5,
			expected: true,
		},
		{
			num:      -2,
			expected: false,
		},
		{
			num:      4,
			expected: false,
		},
	}

	for _, testCase := range testTable {
		result := isEven(testCase.num)

		if result != testCase.expected {
			t.Errorf("Incorrect result. Expect %d, got %t", testCase.num, result)
		}
	}
}
