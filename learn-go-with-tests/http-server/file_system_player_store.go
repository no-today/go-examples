package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// FileSystemPlayerStore 基于文件的存储实现
type FileSystemPlayerStore struct {
	mu       sync.RWMutex
	database *json.Encoder
	players  []Player
}

func (f *FileSystemPlayerStore) GetPlayers() []Player {
	return f.players
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	err := initDatabase(file)
	if err != nil {
		return nil, fmt.Errorf("problem initialising db file, %v", err)
	}

	players, err := readPlayersFromJsonFile(file, err)
	if err != nil {
		return nil, fmt.Errorf("problem decode player from json, %v", err)
	}

	return &FileSystemPlayerStore{
		sync.RWMutex{},
		json.NewEncoder(&tape{file}),
		players,
	}, nil
}

func readPlayersFromJsonFile(file *os.File, err error) ([]Player, error) {
	var players []Player
	err = json.NewDecoder(file).Decode(&players)

	if err != nil {
		return nil, fmt.Errorf("preblem parsing players from json, %v", err)
	}
	return players, nil
}

func initDatabase(file *os.File) error {
	file.Seek(0, 0)

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}
	return nil
}
