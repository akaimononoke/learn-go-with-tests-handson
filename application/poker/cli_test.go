package poker_test

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/akaimononoke/learn-go-with-tests-handson/application/poker"
)

type scheduledAlert struct {
	at     time.Duration
	amount int
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduledAlert{duration, amount})
}

var (
	dummySpyBlindAlerter = &SpyBlindAlerter{}
	dummyPlayerStore     = &poker.StubPlayerStore{}
	dummyStdIn           = &bytes.Buffer{}
	dummyStdOut          = &bytes.Buffer{}
)

func assertScheduledAt(t *testing.T, want, got scheduledAlert) {
	if !reflect.DeepEqual(want, got) {
		t.Errorf("invalid scheduledAlert: want %#v, got %#v", want, got)
	}
}

func TestCLI(t *testing.T) {
	t.Parallel()

	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in, dummyStdOut, dummySpyBlindAlerter)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, "Chris", playerStore)
	})

	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in, dummyStdOut, dummySpyBlindAlerter)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, "Cleo", playerStore)
	})

	t.Run("schedule printing of blind values", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}

		cli := poker.NewCLI(playerStore, in, dummyStdOut, blindAlerter)
		cli.PlayPoker()

		for i, want := range []scheduledAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		} {
			t.Run(fmt.Sprint(want), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				assertScheduledAt(t, want, blindAlerter.alerts[i])
			})
		}
	})

	t.Run("prompts user to enter the number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		cli := poker.NewCLI(dummyPlayerStore, dummyStdIn, stdout, dummySpyBlindAlerter)
		cli.PlayPoker()

		want := "Please enter the number of players: "
		got := stdout.String()
		if want != got {
			t.Errorf("want %q, got %q", want, got)
		}
	})
}
