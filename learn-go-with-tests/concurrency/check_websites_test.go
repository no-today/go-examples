package main

import (
	"reflect"
	"testing"
)

func TestCheckWebsites(t *testing.T) {
	urls := []string{
		"https://google.com",
		"https://cathub.me",
		"waat://null.com",
	}

	want := map[string]bool{
		"https://google.com": true,
		"https://cathub.me":  true,
		"waat://null.com":    false,
	}

	got := CheckWebsites(CheckWebsite, urls)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v, got %v", want, got)
	}
}
