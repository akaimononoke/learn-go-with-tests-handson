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

func AssertScheduledAt(t *testing.T, want, got ScheduledAlert) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("invalid scheduledAlert: want %#v, got %#v", want, got)
	}
}

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   League
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.league
}

func AssertPlayerWin(t *testing.T, wantWinner string, store *StubPlayerStore) {
	t.Helper()
	if 1 != len(store.winCalls) {
		t.Fatal("a win call expected")
	}
	if got := store.winCalls[0]; wantWinner != got {
		t.Errorf("recorded winner is invalid: want %q, got %q", wantWinner, got)
	}
}

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
