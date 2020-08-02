package poker_test

import (
	"strings"
	"testing"

	"github.com/akaimononoke/learn-go-with-tests-handson/application/poker"
)

func TestCLI(t *testing.T) {
	t.Parallel()

	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, "Chris", playerStore)
	})

	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, "Cleo", playerStore)
	})
}
