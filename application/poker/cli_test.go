package poker_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/akaimononoke/learn-go-with-tests-handson/application/poker"
)

type SpyBlindAlerter struct {
	alerts []struct {
		scheduledAt time.Duration
		amount      int
	}
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	s.alerts = append(s.alerts, struct {
		scheduledAt time.Duration
		amount      int
	}{duration, amount})
}

var dummySpyBlindAlerter = &SpyBlindAlerter{}

func TestCLI(t *testing.T) {
	t.Parallel()

	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in, dummySpyBlindAlerter)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, "Chris", playerStore)
	})

	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in, dummySpyBlindAlerter)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, "Cleo", playerStore)
	})

	t.Run("schedule printing of blind values", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}

		cli := poker.NewCLI(playerStore, in, blindAlerter)
		cli.PlayPoker()

		for i, c := range []struct {
			wantScheduleTime time.Duration
			wantAmount       int
		}{
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
			t.Run(fmt.Sprintf("%d scheduled for %v", c.wantAmount, c.wantScheduleTime), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				alert := blindAlerter.alerts[i]

				if gotAmount := alert.amount; c.wantAmount != gotAmount {
					t.Errorf("invalid amount: want %d, got %d", c.wantAmount, gotAmount)
				}
				if gotScheduleTime := alert.scheduledAt; c.wantScheduleTime != gotScheduleTime {
					t.Errorf("invalid schedule time: want %v, got %v", c.wantScheduleTime, gotScheduleTime)
				}
			})
		}
	})
}
