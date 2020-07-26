package main

import (
	"testing"
)

func TestHello(t *testing.T) {
	for desc, c := range map[string]struct {
		Name string
		Want string
	}{
		"success.1": {"Chris", "Hello, Chris!"},
		"success.2": {"Kazumasa", "Hello, Kazumasa!"},
	} {
		t.Log(desc)
		got := Hello(c.Name)
		if got != c.Want {
			t.Errorf("Hello() is invalid: got %q, want %q", got, c.Want)
		}
	}
}
