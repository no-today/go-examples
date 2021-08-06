package main

import (
	"testing"
)

const key = "Hello"
const val = "你好"

func TestGet(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		map_ := Map{}
		map_.Put(key, val)
		assertExists(t, map_, key, val)
	})

	t.Run("NotFound", func(t *testing.T) {
		map_ := Map{}
		assertNotExists(t, map_, key)
	})
}

func assertNotExists(t *testing.T, map_ Map, key string) {
	_, err := map_.Get(key)

	assertError(t, err, ErrorNotFound.Error())
}

func assertExists(t *testing.T, map_ Map, key, val string) {
	t.Helper()

	got, err := map_.Get(key)

	assertStrings(t, got, val)
	assertNotError(t, err)
}

func TestPut(t *testing.T) {
	map_ := Map{}
	assertPut(t, map_, key, val)
}

func assertPut(t *testing.T, map_ Map, key, val string) {
	map_.Put(key, val)

	assertExists(t, map_, key, val)
}

func TestRemove(t *testing.T) {
	map_ := Map{}
	assertPut(t, map_, key, val)
	map_.Remove(key)
	assertNotExists(t, map_, key)
}

func assertStrings(t *testing.T, got string, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func assertError(t *testing.T, err error, want string) {
	if err == nil {
		t.Fatal("wanted an error but didn't get one")
	}

	if err.Error() != want {
		t.Errorf("got %s, want %s", err.Error(), want)
	}
}

func assertNotError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("no errors expected")
	}
}
