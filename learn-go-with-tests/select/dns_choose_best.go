package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var tenSecondTimeout = 10 * time.Second

// DNSChooseBest DNS 选优
func DNSChooseBest(ip1, ip2, ip3 string) (ip string, err error) {
	var cost time.Duration

	select {
	case c1 := <-ping(ip1):
		ip = ip1
		cost = c1
	case c2 := <-ping(ip2):
		ip = ip2
		cost = c2
	case c3 := <-ping(ip3):
		ip = ip3
		cost = c3
	case <-time.After(tenSecondTimeout):
		return "", errors.New("all dns ip unavailable")
	}

	fmt.Println("The fastest DNS ip is:", ip, ", time consuming:", cost)

	return ip, nil
}

// ping 测试连通性, 在连接通畅时往通道写入耗时
func ping(ip string) chan time.Duration {
	ch := make(chan time.Duration)
	start := time.Now()

	go func() {
		millisecond := time.Duration(rand.Int63n(500))
		time.Sleep(millisecond * time.Millisecond)

		// 十分之一概率失败(模拟ping不通)
		if rand.Int31n(10)%10 == 0 {
			return
		}

		// 写入耗时
		ch <- time.Since(start)
	}()

	return ch
}
