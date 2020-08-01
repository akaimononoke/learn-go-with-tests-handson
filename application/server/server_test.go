package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPlayers(t *testing.T) {
	t.Parallel()

	t.Run("returns Pepper's score", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/players/Pepper", nil)
		res := httptest.NewRecorder()

		PlayerServer(res, req)

		want := "20"
		got := res.Body.String()

		if got != want {
			t.Errorf("want %q, want %q", want, got)
		}
	})
}
