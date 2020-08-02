package poker

import (
	"testing"
)

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
