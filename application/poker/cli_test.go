package poker_test

import (
	"bytes"
	"io"
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

	StartCalled  bool
	FinishCalled bool
}

func (g *SpyGame) Start(numberOfPlayers int, to io.Writer) {
	g.StartCalled = true
	g.StartedWith = numberOfPlayers
}

func (g *SpyGame) Finish(winner string) {
	g.FinishCalled = true
	g.FinishedWith = winner
}

func userSends(messages ...string) io.Reader {
	return strings.NewReader(strings.Join(messages, "\n"))
}

func assertMessagesSentToUser(t *testing.T, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if want != got {
		t.Errorf("messages are invalid: want %+v, got %q", messages, got)
	}
}

func assertGameNotStarted(t *testing.T, game *SpyGame) {
	t.Helper()
	if game.StartCalled {
		t.Errorf("game should not have started")
	}
}

func assertGameStartedWith(t *testing.T, wantNumberOfPlayers int, game *SpyGame) {
	t.Helper()
	if wantNumberOfPlayers != game.StartedWith {
		t.Errorf("number of players is invalid: want %d, got %d", wantNumberOfPlayers, game.StartedWith)
	}
}

func assertGameNotFinished(t *testing.T, game *SpyGame) {
	t.Helper()
	if game.FinishCalled {
		t.Errorf("game should not have finished")
	}
}

func assertGameFinishedWith(t *testing.T, wantWinner string, game *SpyGame) {
	t.Helper()
	if wantWinner != game.FinishedWith {
		t.Errorf("winner is invalid: want %q, got %q", wantWinner, game.FinishedWith)
	}
}

func assertScheduledAt(t *testing.T, want, got poker.ScheduledAlert) {
	t.Helper()
	if want != got {
		t.Errorf("want %+v, got %+v", want, got)
	}
}

func TestCLI(t *testing.T) {
	t.Parallel()

	t.Run("start with 3 players, finish with 'Chris' as winner", func(t *testing.T) {
		stdin := userSends("3", "Chris wins")
		stdout := &bytes.Buffer{}
		game := &SpyGame{}
		cli := poker.NewCLI(stdin, stdout, game)

		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
		assertGameStartedWith(t, 3, game)
		assertGameFinishedWith(t, "Chris", game)
	})

	t.Run("start with 8 players, finish with 'Cleo' as winner", func(t *testing.T) {
		stdin := userSends("8", "Cleo wins")
		game := &SpyGame{}
		cli := poker.NewCLI(stdin, dummyStdOut, game)

		cli.PlayPoker()

		assertGameStartedWith(t, 8, game)
		assertGameFinishedWith(t, "Cleo", game)
	})

	t.Run("non numeric value is sent as number of players, print error", func(t *testing.T) {
		stdin := userSends("pies")
		stdout := &bytes.Buffer{}
		game := &SpyGame{}
		cli := poker.NewCLI(stdin, stdout, game)

		cli.PlayPoker()

		assertGameNotStarted(t, game)
		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.BadPlayerInputErrMsg)
	})

	t.Run("invalid syntax of winner declaration, print error", func(t *testing.T) {
		stdin := userSends("8", "Floyd is a winner")
		stdout := &bytes.Buffer{}
		game := &SpyGame{}
		cli := poker.NewCLI(stdin, stdout, game)

		cli.PlayPoker()

		assertGameNotFinished(t, game)
		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.BadWinnerInputMsg)
	})
}
