package poker

import (
	"fmt"
	"io"
	"reflect"
	"testing"
	"time"
)

type ScheduledAlert struct {
	At     time.Duration
	Amount int
}

func (s ScheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.Amount, s.At)
}

type SpyBlindAlerter struct {
	Alerts []ScheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
	s.Alerts = append(s.Alerts, ScheduledAlert{duration, amount})
}

type StubPlayerStore struct {
	Scores   map[string]int
	WinCalls []string
	League   League
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.Scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.WinCalls = append(s.WinCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.League
}

type SpyGame struct {
	StartedWith  int
	FinishedWith string
	BlindAlert   []byte

	StartCalled  bool
	FinishCalled bool
}

func (g *SpyGame) Start(numberOfPlayers int, out io.Writer) {
	g.StartCalled = true
	g.StartedWith = numberOfPlayers
	out.Write(g.BlindAlert)
}

func (g *SpyGame) Finish(winner string) {
	g.FinishCalled = true
	g.FinishedWith = winner
}

func AssertScheduledAt(t *testing.T, want, got ScheduledAlert) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("invalid scheduledAlert: want %#v, got %#v", want, got)
	}
}

func AssertPlayerWin(t *testing.T, wantWinner string, store *StubPlayerStore) {
	t.Helper()
	if 1 != len(store.WinCalls) {
		t.Fatal("a win call expected")
	}
	if got := store.WinCalls[0]; wantWinner != got {
		t.Errorf("recorded winner is invalid: want %q, got %q", wantWinner, got)
	}
}
