package main

import (
	"testing"
)

func TestDNSChooseBest(t *testing.T) {
	for i := 0; i < 10; i++ {
		DNSChooseBest("223.5.5.5", "223.6.6.6", "8.8.8.8")
	}
}
