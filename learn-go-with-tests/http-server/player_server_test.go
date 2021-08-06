package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayers() []Player {
	panic("implement me")
}

func (s *StubPlayerStore) GetPlayerScore(uid string) int {
	return s.scores[uid]
}

func (s *StubPlayerStore) RecordWin(uid string) {
	s.winCalls = append(s.winCalls, uid)
}

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"1001": 20,
			"1002": 10,
		},
		nil,
	}

	server := NewPlayerServer(&store)

	cases := []struct {
		tn                 string
		player             string
		expectedHTTPStatus int
		expectedScore      string
	}{
		{
			"获取 1001 玩家的得分",
			"1001",
			http.StatusOK,
			"20",
		},
		{
			"获取 1002 玩家的得分",
			"1002",
			http.StatusOK,
			"10",
		},
		{
			"获取 1003 玩家的得分",
			"1003",
			http.StatusOK,
			"0",
		},
	}

	for _, tt := range cases {
		t.Run(tt.tn, func(t *testing.T) {
			request := newGetScoreRequest(tt.player)
			response := httptest.NewRecorder()

			server.ShowScore(response, request)

			assertStatus(t, tt.expectedHTTPStatus, response.Code)
			assertResponseBody(t, tt.expectedScore, response.Body.String())
		})
	}
}

func TestStoreWin(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,
	}

	server := NewPlayerServer(&store)

	t.Run("记录玩家获胜", func(t *testing.T) {
		player := "Player"

		request := newRecordWinRequest(player)
		response := httptest.NewRecorder()

		server.ProcessWin(response, request)

		if len(store.winCalls) != 1 {
			t.Fatalf("expected: 1, got: %v", len(store.winCalls))
		}

		if store.winCalls[0] != player {
			t.Errorf("expected: %v, got: %v", player, store.winCalls[0])
		}
	})
}

func assertResponseBody(t *testing.T, expected, got string) {
	if got != expected {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func assertStatus(t *testing.T, expected, got int) {
	if got != expected {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func newGetScoreRequest(uid string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/playsers/%s", uid), nil)
	return request
}

func newRecordWinRequest(uid string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/playsers/%s", uid), nil)
	return request
}
