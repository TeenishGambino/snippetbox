package assert

import (
	"testing"
)

// Generic Equal function for testing
func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got %q; want %q", actual, "17 Mar 2022 at 10:15")	
	}
}