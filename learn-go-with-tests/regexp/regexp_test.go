package main

import (
	"fmt"
	"regexp"
	"testing"
)

func TestExtractPathParams(t *testing.T) {
	re, _ := regexp.Compile("/(:[a-zA-Z0-9_]+)")

	fmt.Println(re.FindAllStringSubmatch("/players/:id/:uid/:pid", -1))
	fmt.Println(re.ReplaceAllString("/players/:id/:uid/:pid", re.String()))
}
