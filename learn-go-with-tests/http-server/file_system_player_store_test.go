package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("获取玩家列表", func(t *testing.T) {
		file, cleanDatabase := createTempFile(t, `[
			{"Uid": "p1", "Wins": 10},
			{"Uid": "p2", "Wins": 33}
		]`)
		defer cleanDatabase()

		want := []Player{
			{"p1", 10},
			{"p2", 33},
		}

		store, err := NewFileSystemPlayerStore(file)
		assertNotError(t, err)

		got := store.GetPlayers()
		assertPlayers(t, got, want)
	})
}

func assertPlayers(t *testing.T, got []Player, want []Player) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func assertNotError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("expected no error, but appeared: %v", err)
	}
}

// 创建临时文件
func createTempFile(t *testing.T, initialDate string) (*os.File, func()) {
	t.Helper()

	file, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("create temp file fail: %v", err)
	}

	file.Write([]byte(initialDate))

	removeFile := func() {
		os.Remove(file.Name())
	}

	return file, removeFile
}
