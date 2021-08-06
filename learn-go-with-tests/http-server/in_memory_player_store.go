package main

import (
	"sync"
)

// InMemoryPlayerStore 基于内存的存储实现
type InMemoryPlayerStore struct {
	mu     sync.RWMutex
	scores map[string]int
}

// NewInMemoryPlayerStore 初始化一个空的玩家信息存储
func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		sync.RWMutex{},
		make(map[string]int)}
}

// GetPlayerScore 获取玩家得分
func (m *InMemoryPlayerStore) GetPlayerScore(uid string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.scores[uid]
}

// RecordWin 记录玩家获胜得分
func (m *InMemoryPlayerStore) RecordWin(uid string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.scores[uid]++
}

func (m *InMemoryPlayerStore) GetPlayers() (players []Player) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for uid, score := range m.scores {
		players = append(players, Player{uid, score})
	}

	if players == nil {
		players = []Player{}
	}

	return players
}
