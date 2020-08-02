package poker_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/akaimononoke/learn-go-with-tests-handson/application/poker"
)

var (
	dummySpyBlindAlerter = &poker.SpyBlindAlerter{}
	dummyPlayerStore     = &poker.StubPlayerStore{}
	dummyStdIn           = &bytes.Buffer{}
	dummyStdOut          = &bytes.Buffer{}
)

type SpyGame struct {
	StartedWith  int
	FinishedWith string

	StartCalled bool
}

func (g *SpyGame) Start(numberOfPlayers int) {
	g.StartedWith = numberOfPlayers
	g.StartCalled = true
}

func (g *SpyGame) Finish(winner string) {
	g.FinishedWith = winner
}

func TestCLI(t *testing.T) {
	t.Parallel()

	t.Run("prints error when a non numeric value is entered and doesn't not start game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("Dummy\n")
		game := &SpyGame{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		if game.StartCalled {
			t.Error("game should not have started")
		}

		wantPrompt := poker.PlayerPrompt + "you're so silly"
		gotPrompt := stdout.String()
		if wantPrompt != gotPrompt {
			t.Errorf("want %q, got %q", wantPrompt, gotPrompt)
		}
	})

	t.Run("it prompts the user to enter the number of players and starts the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		game := &SpyGame{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		gotPrompt := stdout.String()
		wantPrompt := poker.PlayerPrompt

		if gotPrompt != wantPrompt {
			t.Errorf("got %q, want %q", gotPrompt, wantPrompt)
		}

		if game.StartedWith != 7 {
			t.Errorf("wanted Start called with 7 but got %d", game.StartedWith)
		}
	})
}
