package crudlib

import "testing"

func TestSaysHello(t *testing.T) {
	c := NewClient()

	result := c.SayHello()
	expected := "Hello"

	if result != expected {
		t.Errorf("Got %q but expected %q", result, expected)
	}
}
