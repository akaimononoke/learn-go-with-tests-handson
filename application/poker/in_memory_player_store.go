package poker

import (
	"sort"
	"sync"
)

type InMemoryPlayerStore struct {
	mu    sync.RWMutex
	Store map[string]int
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		Store: map[string]int{},
	}
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.Store[name]
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.mu.Lock()
	i.Store[name]++
	i.mu.Unlock()
}

func (i *InMemoryPlayerStore) GetLeague() League {
	var league League
	for name, wins := range i.Store {
		league = append(league, Player{name, wins})
	}
	sort.Slice(league, func(i, j int) bool { return league[i].Wins > league[j].Wins })
	return league
}
