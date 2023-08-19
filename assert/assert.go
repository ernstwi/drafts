// https://www.alexedwards.net/blog/easy-test-assertions-with-go-generics

package assert

import (
	"reflect"
	"testing"
)

func Equal[T comparable](t *testing.T, expected, actual T) {
	t.Helper()

	if expected != actual {
		t.Errorf("want: %v; got: %v", expected, actual)
	}
}

func EqualSlice[T comparable](t *testing.T, expected, actual []T) {
	t.Helper()

	if len(expected) != len(actual) {
		t.Fatalf("len want: %v; len got: %v", len(expected), len(actual))
	}

	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("want: %v; got: %v", expected, actual)
		}
	}
}

func DeepEqual[T comparable](t *testing.T, expected, actual T) {
	t.Helper()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("want: %v; got: %v", expected, actual)
	}
}
