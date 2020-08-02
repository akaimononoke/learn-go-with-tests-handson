package application

import (
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	in := strings.NewReader("Chris wins\n")

	playerStore := &StubPlayerStore{}
	cli := &CLI{playerStore, in}
	cli.PlayPoker()

	if len(playerStore.winCalls) != 1 {
		t.Fatal("a win call expected")
	}

	want := "Chris"
	got := playerStore.winCalls[0]

	if want != got {
		t.Errorf("recorded winner is invalid: want %q, got %q", want, got)
	}
}
