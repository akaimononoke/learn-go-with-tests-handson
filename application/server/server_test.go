package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() []Player {
	return s.league
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
		nil,
	}
	server := NewPlayerServer(store)

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
		nil,
	}
	server := NewPlayerServer(store)

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

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func getLeagueFromRequest(t *testing.T, body io.Reader) []Player {
	t.Helper()
	league, err := NewLeague(body)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player: %v", body, err)
	}
	return league
}

func assertLeague(t *testing.T, want, got []Player) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want league %#v, got %#v", want, got)
	}
}

const jsonContentType = "application/json"

func assertContentType(t *testing.T, want string, res *httptest.ResponseRecorder) {
	t.Helper()
	if got := res.Result().Header.Get("content-type"); want != got {
		t.Errorf("content-type is invalid: wanted %v, got %v", want, got)
	}
}

func TestLeague(t *testing.T) {
	t.Parallel()

	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}
		store := &StubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(store)

		req := newLeagueRequest()
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		got := getLeagueFromRequest(t, res.Body)
		assertStatus(t, http.StatusOK, res.Code)
		assertLeague(t, wantedLeague, got)
		assertContentType(t, jsonContentType, res)
	})
}

func TestFileSystemPlayerStore(t *testing.T) {
	t.Parallel()

	t.Run("/league from a reader", func(t *testing.T) {
		db := strings.NewReader(`[{"Name": "Cleo", "Wins": 10}, {"Name": "Chris", "Wins": 20}]`)
		store := FileSystemPlayerStore{db}

		want := []Player{
			{"Cleo", 10},
			{"Chris", 20},
		}
		got := store.GetLeague()
		assertLeague(t, want, got)

		// read again
		got = store.GetLeague()
		assertLeague(t, want, got)
	})

	t.Run("get player score", func(t *testing.T) {
		db := strings.NewReader(`[{"Name": "Cleo", "Wins": 10}, {"Name": "Chris", "Wins": 33}]`)
		store := FileSystemPlayerStore{db}

		want := 33
		got := store.GetPlayerScore("Chris")

		if want != got {
			t.Errorf("score is invalid: want %d, got %d", want, got)
		}
	})
}

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	t.Parallel()

	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		res := httptest.NewRecorder()
		server.ServeHTTP(res, newGetScoreRequest(player))

		assertStatus(t, http.StatusOK, res.Code)
		assertResponseBody(t, "3", res.Body.String())
	})

	t.Run("get league", func(t *testing.T) {
		res := httptest.NewRecorder()
		server.ServeHTTP(res, newLeagueRequest())

		assertStatus(t, http.StatusOK, res.Code)

		want := []Player{{"Pepper", 3}}
		got := getLeagueFromRequest(t, res.Body)

		assertLeague(t, want, got)
	})
}

func TestPlayerStoreConcurrentlySafety(t *testing.T) {
	numProc := 1000

	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
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
}
