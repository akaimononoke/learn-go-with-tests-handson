package poker

import (
	"testing"
)

func assertScoreEquals(t *testing.T, want, got int) {
	t.Helper()
	if want != got {
		t.Errorf("score is invalid: want %d, got %d", want, got)
	}
}

func TestFileSystemPlayerStore(t *testing.T) {
	t.Parallel()

	t.Run("/league from a reader", func(t *testing.T) {
		db, cleanDatabase := createTempFile(t, `[{"Name": "Cleo", "Wins": 10}, {"Name": "Chris", "Wins": 20}]`)
		defer cleanDatabase()
		store, err := NewFileSystemPlayerStore(db)

		assertNoError(t, err)

		want := League{
			{"Chris", 20},
			{"Cleo", 10},
		}
		assertLeague(t, want, store.GetLeague())

		// read again
		assertLeague(t, want, store.GetLeague())
	})

	t.Run("get player score", func(t *testing.T) {
		db, cleanDatabase := createTempFile(t, `[{"Name": "Cleo", "Wins": 10}, {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()
		store, err := NewFileSystemPlayerStore(db)

		assertNoError(t, err)
		assertScoreEquals(t, 33, store.GetPlayerScore("Chris"))
	})

	t.Run("store wins for existing players", func(t *testing.T) {
		db, cleanDatabase := createTempFile(t, `[{"Name": "Cleo", "Wins": 10}, {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()
		store, err := NewFileSystemPlayerStore(db)

		assertNoError(t, err)

		store.RecordWin("Chris")

		assertScoreEquals(t, 34, store.GetPlayerScore("Chris"))
	})

	t.Run("store wins for new player", func(t *testing.T) {
		db, cleanDatabase := createTempFile(t, `[{"Name": "Cleo", "Wins": 10}, {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()
		store, err := NewFileSystemPlayerStore(db)

		assertNoError(t, err)
		store.RecordWin("Pepper")

		assertScoreEquals(t, 1, store.GetPlayerScore("Pepper"))
	})

	t.Run("works with an empty file", func(t *testing.T) {
		db, cleanDatabase := createTempFile(t, ``)
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(db)

		assertNoError(t, err)
	})
}
