package poker

import (
	"sync"
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

	t.Run("get player score", func(t *testing.T) {
		store := &InMemoryPlayerStore{
			store: map[string]int{
				"Cleo":  10,
				"Chris": 20,
			},
		}

		assertScoreEquals(t, 20, store.GetPlayerScore("Chris"))
	})

	t.Run("save wins for existing player", func(t *testing.T) {
		store := &InMemoryPlayerStore{
			store: map[string]int{
				"Cleo":  10,
				"Chris": 20,
			},
		}

		store.RecordWin("Chris")

		assertScoreEquals(t, 21, store.GetPlayerScore("Chris"))
	})

	t.Run("save wins for new player", func(t *testing.T) {
		store := &InMemoryPlayerStore{
			store: map[string]int{
				"Cleo":  10,
				"Chris": 20,
			},
		}

		store.RecordWin("Bob")

		assertScoreEquals(t, 1, store.GetPlayerScore("Bob"))
	})

	t.Run("concurrently safety", func(t *testing.T) {
		numProc := 1000

		store := NewInMemoryPlayerStore()

		var wg sync.WaitGroup
		wg.Add(numProc)

		for i := 0; i < numProc; i++ {
			go func(wg *sync.WaitGroup) {
				store.RecordWin("Pepper")
				wg.Done()
			}(&wg)
		}
		wg.Wait()

		assertScoreEquals(t, 1000, store.GetPlayerScore("Pepper"))
	})
}
