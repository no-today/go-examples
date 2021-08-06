package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func BenchmarkCheckWebsites(b *testing.B) {
	urls := make([]string, 100)
	for i := 0; i < len(urls); i++ {
		if i%2 == 0 {
			urls[i] = buildDelayedServer(20 * time.Millisecond).URL
		} else {
			urls[i] = buildDelayedServer(0 * time.Millisecond).URL
		}
	}

	for i := 0; i < b.N; i++ {
		CheckWebsites(CheckWebsite, urls)
	}
}

// 构建延迟响应服务
func buildDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(delay)
		writer.WriteHeader(http.StatusOK)
	}))
}
