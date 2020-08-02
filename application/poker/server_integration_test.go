package poker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	t.Parallel()

	db, cleanDatabase := createTempFile(t, "[]")
	defer cleanDatabase()
	store, err := NewFileSystemPlayerStore(db)

	assertNoError(t, err)

	server, _ := NewPlayerServer(store)
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

		want := League{{"Pepper", 3}}
		got := getLeagueFromRequest(t, res.Body)

		assertLeague(t, want, got)
	})
}
