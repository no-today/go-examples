package main

import (
	"net/http"
	"time"
)

var (
	CheckTimeout = 5 * time.Second
)

// 设定超时时间
var client = http.Client{Timeout: CheckTimeout}

// CheckWebsite URL 响应码为 200 时返回 ture, 否则返回 false
func CheckWebsite(url string) bool {
	// 我们只需要验证连接通畅, 不需要返回内容
	response, err := client.Head(url)
	if err != nil {
		return false
	}

	return response.StatusCode == http.StatusOK
}
