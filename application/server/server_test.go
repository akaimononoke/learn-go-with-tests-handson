package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertStatus(t *testing.T, want, got int) {
	t.Helper()
	if want != got {
		t.Errorf("want status %d, got %d", want, got)
	}
}

func assertResponseBody(t *testing.T, want, got string) {
	t.Helper()
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestGetPlayers(t *testing.T) {
	t.Parallel()

	store := &StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
	}
	server := &PlayerServer{store}

	t.Run("returns Pepper's score", func(t *testing.T) {
		req := newGetScoreRequest("Pepper")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		assertStatus(t, http.StatusOK, res.Code)
		assertResponseBody(t, "20", res.Body.String())
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		req := newGetScoreRequest("Floyd")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		assertStatus(t, http.StatusOK, res.Code)
		assertResponseBody(t, "10", res.Body.String())
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		req := newGetScoreRequest("Apollo")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		assertStatus(t, http.StatusNotFound, res.Code)
	})
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func TestStoreWins(t *testing.T) {
	t.Parallel()

	store := &StubPlayerStore{
		map[string]int{},
		nil,
	}
	server := &PlayerServer{store}

	t.Run("it records wins when POST", func(t *testing.T) {
		player := "Pepper"

		req := newPostWinRequest(player)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		assertStatus(t, http.StatusAccepted, res.Code)

		if 1 != len(store.winCalls) {
			t.Errorf("want calls %d, got %d", 1, len(store.winCalls))
		}
		if player != store.winCalls[0] {
			t.Errorf("collected store winner is invalid: want %q, got %q", player, store.winCalls[0])
		}
	})
}

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	t.Parallel()

	t.Run("records 3 times wins", func(t *testing.T) {
		store := NewInMemoryPlayerStore()
		server := &PlayerServer{store}
		player := "Pepper"

		server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
		server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
		server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

		res := httptest.NewRecorder()
		server.ServeHTTP(res, newGetScoreRequest(player))

		assertStatus(t, http.StatusOK, res.Code)
		assertResponseBody(t, "3", res.Body.String())
	})

	t.Run("records wins safely concurrently", func(t *testing.T) {
		numProc := 1000

		store := NewInMemoryPlayerStore()
		server := &PlayerServer{store}
		player := "Pepper"

		var wg sync.WaitGroup
		wg.Add(numProc)

		for i := 0; i < numProc; i++ {
			go func(w *sync.WaitGroup) {
				server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
				w.Done()
			}(&wg)
		}
		wg.Wait()

		res := httptest.NewRecorder()
		server.ServeHTTP(res, newGetScoreRequest(player))

		assertResponseBody(t, "1000", res.Body.String())
	})
}
