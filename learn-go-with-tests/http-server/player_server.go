package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// PlayerStore 存储玩家信息
type PlayerStore interface {

	// GetPlayerScore 获取玩家得分
	GetPlayerScore(uid string) int

	// RecordWin 记录玩家获胜次数
	RecordWin(uid string)

	// GetPlayers 获取玩家信息
	GetPlayers() []Player
}

// Player 玩家信息
type Player struct {
	Uid  string
	Wins int
}

// PlayerServer 玩家服务
type PlayerServer struct {
	store   PlayerStore
	handler http.Handler
}

const jsonContentType = "application/json"

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)
	p.store = store

	router := NewHttpRoute()
	router.GET("/players/:uid", p.ShowScore)
	router.POST("/players/:uid", p.ProcessWin)
	router.GET("/players", p.GetPlayers)

	p.handler = router

	return p
}

func (p *PlayerServer) ShowScore(w http.ResponseWriter, r *http.Request) {
	score := p.store.GetPlayerScore(getUid(r))
	fmt.Fprint(w, score)
}

func (p *PlayerServer) ProcessWin(w http.ResponseWriter, r *http.Request) {
	p.store.RecordWin(getUid(r))
	w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) GetPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(p.store.GetPlayers())
}

func getUid(r *http.Request) string {
	split := strings.Split(r.URL.Path, "/")
	uid := split[len(split)-1]
	return uid
}
