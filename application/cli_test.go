package application

import (
	"testing"
)

func TestCLI(t *testing.T) {
	playerStore := &StubPlayerStore{}
	cli := &CLI{playerStore}
	cli.PlayPoker()

	if len(playerStore.winCalls) != 1 {
		t.Fatal("a win call expected")
	}
}
