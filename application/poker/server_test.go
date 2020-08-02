package poker_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/akaimononoke/learn-go-with-tests-handson/application/poker"
	"github.com/gorilla/websocket"
)

const jsonContentType = "application/json"

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func getLeagueFromRequest(t *testing.T, body io.Reader) poker.League {
	t.Helper()
	league, err := poker.NewLeague(body)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player: %v", body, err)
	}
	return league
}

func newGameRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/game", nil)
	return req
}

func writeWebSocketMessage(t *testing.T, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("failed to send message over websocket connection: %v", err)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("error occurred unexpectedly: %v", err)
	}
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

func assertLeague(t *testing.T, want, got poker.League) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want league %#v, got %#v", want, got)
	}
}

func assertContentType(t *testing.T, want string, res *httptest.ResponseRecorder) {
	t.Helper()
	if got := res.Result().Header.Get("content-type"); want != got {
		t.Errorf("content-type is invalid: wanted %v, got %v", want, got)
	}
}

func mustMakePlayerServer(t *testing.T, store poker.PlayerStore) *poker.PlayerServer {
	server, err := poker.NewPlayerServer(store)
	if err != nil {
		t.Fatalf("failed to create player server: %v", err)
	}
	return server
}

func mustDialWebSocket(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("failed to open websocket connection on %s: %v", url, err)
	}
	return ws
}

func TestGetPlayers(t *testing.T) {
	t.Parallel()

	store := &poker.StubPlayerStore{
		Scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
	}
	server, _ := poker.NewPlayerServer(store)

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

func TestStoreWins(t *testing.T) {
	t.Parallel()

	store := &poker.StubPlayerStore{
		Scores: map[string]int{},
	}
	server, _ := poker.NewPlayerServer(store)

	t.Run("it records wins when POST", func(t *testing.T) {
		player := "Pepper"

		req := newPostWinRequest(player)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		assertStatus(t, http.StatusAccepted, res.Code)
		poker.AssertPlayerWin(t, "Pepper", store)
	})
}

func TestLeague(t *testing.T) {
	t.Parallel()

	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := poker.League{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}
		store := &poker.StubPlayerStore{nil, nil, wantedLeague}
		server, _ := poker.NewPlayerServer(store)

		req := newLeagueRequest()
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		got := getLeagueFromRequest(t, res.Body)
		assertStatus(t, http.StatusOK, res.Code)
		assertLeague(t, wantedLeague, got)
		assertContentType(t, jsonContentType, res)
	})
}

func TestGame(t *testing.T) {
	t.Parallel()

	t.Run("returns 200 for GET /game", func(t *testing.T) {
		server, _ := poker.NewPlayerServer(&poker.StubPlayerStore{})
		req := newGameRequest()
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		assertStatus(t, http.StatusOK, res.Code)
	})

	t.Run("receive winner of a game over websocket", func(t *testing.T) {
		winner := "Go"

		store := &poker.StubPlayerStore{}
		playerServer := mustMakePlayerServer(t, store)
		server := httptest.NewServer(playerServer)
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
		ws := mustDialWebSocket(t, wsURL)
		writeWebSocketMessage(t, ws, winner)

		time.Sleep(10 * time.Millisecond)
		poker.AssertPlayerWin(t, winner, store)
	})
}
