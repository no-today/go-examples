package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	player := "Player"

	server.ProcessWin(httptest.NewRecorder(), newRecordWinRequest(player))
	server.ProcessWin(httptest.NewRecorder(), newRecordWinRequest(player))
	server.ProcessWin(httptest.NewRecorder(), newRecordWinRequest(player))

	response := httptest.NewRecorder()
	server.ShowScore(response, newGetScoreRequest(player))

	assertStatus(t, http.StatusOK, response.Code)
	assertResponseBody(t, "3", response.Body.String())
}
