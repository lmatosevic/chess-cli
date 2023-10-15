package utils

import "testing"

func AssertTestCondition(t *testing.T, expected any, received any, errorMessage string) {
	if expected != received {
		t.Fatalf("%s\nExpected: %v\nReceived: %v\n", errorMessage, expected, received)
	}
}
