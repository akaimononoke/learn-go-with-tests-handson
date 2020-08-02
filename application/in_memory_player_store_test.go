package application

import (
	"testing"
)

func TestInMemoryPlayerStore(t *testing.T) {
	t.Parallel()
	t.Run("get league", func(t *testing.T) {
		store := &InMemoryPlayerStore{
			store: map[string]int{
				"Cleo":  10,
				"Chris": 20,
			},
		}

		want := League{
			{"Chris", 20},
			{"Cleo", 10},
		}
		assertLeague(t, want, store.GetLeague())
	})
	// t.Run("get player score", func(t *testing.T) {})
	// t.Run("save wins for existing player", func(t *testing.T) {})
	// t.Run("save wins for new player", func(t *testing.T) {})
	// t.Run("concurrently safety", func(t *testing.T) {})
}
