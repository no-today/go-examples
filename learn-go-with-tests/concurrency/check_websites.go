package main

import (
	"fmt"
	"time"
)

// WebsiteChecker 检测该 URL 是否可用
type WebsiteChecker func(url string) bool
type result struct {
	url string
	ok  bool
}

// CheckWebsites 检测切面内所以的 URL 是否可用, 然后返回结果
// 使用 WebsiteChecker 函数检测 URL 是否可用
func CheckWebsites(checker WebsiteChecker, urls []string) map[string]bool {
	start := time.Now()

	results := make(map[string]bool)
	channel := make(chan result)

	for _, url := range urls {
		go func(u string) {
			channel <- result{u, checker(u)}
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		result := <-channel
		results[result.url] = result.ok
	}

	fmt.Println("take size:", len(urls), "cost:", time.Since(start))

	return results
}
