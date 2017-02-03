package main

import (
	"testing"
)

func TestFormatExecutionTime(t *testing.T) {
	testcases := []struct {
		executionTime uint64
		expected      string
	}{
		{
			executionTime: 0,
			expected:      "n/a",
		},
		{
			executionTime: 100,
			expected:      "Less than a second",
		},
		{
			executionTime: 1 * 1000 * 1000,
			expected:      "1 second",
		},
		{
			executionTime: 20 * 1000 * 1000,
			expected:      "20 seconds",
		},
		{
			executionTime: 60 * 1000 * 1000,
			expected:      "1 minute",
		},
		{
			executionTime: 61 * 1000 * 1000,
			expected:      "1 minute 1 second",
		},
		{
			executionTime: 62 * 1000 * 1000,
			expected:      "1 minute 2 seconds",
		},
		{
			executionTime: 120 * 1000 * 1000,
			expected:      "2 minutes",
		},
		{
			executionTime: 121 * 1000 * 1000,
			expected:      "2 minutes 1 second",
		},
		{
			executionTime: 142 * 1000 * 1000,
			expected:      "2 minutes 22 seconds",
		},
		{
			executionTime: 192 * 60 * 1000 * 1000,
			expected:      "192 minutes",
		},
	}

	for _, testcase := range testcases {
		got := formatExecutionTime(testcase.executionTime)
		if got != testcase.expected {
			t.Errorf("got: %q, expected: %q", got, testcase.expected)
		}
	}
}
