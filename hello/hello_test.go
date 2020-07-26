package main

import (
	"testing"
)

func TestHello(t *testing.T) {
	for desc, c := range map[string]struct {
		Value string
	}{
		"success": {"Hello, world!"},
	} {
		t.Log(desc)
		got := Hello()
		want := c.Value
		if got != want {
			t.Errorf("Hello() is invalid: got %q, want %q", got, want)
		}
	}
}
