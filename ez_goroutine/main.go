package main

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

var (
	chars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789"
	tail   = "DEHNCHyUEO8kVZBT"
	result = "3354de5346a962dd0f344de80cd3c8e5c2d3ce1a18437141b6a645df9b357c91"
	wg     sync.WaitGroup
)

func sha(head string) {
	h := sha256.New()
	h.Write([]byte(head + tail))
	str := fmt.Sprintf("%x", h.Sum(nil))
	if str == result {
		fmt.Println(head)
	}
	wg.Done() // 完成任务
}

func main() {
	start := time.Now()
	for _, ch1 := range chars {
		for _, ch2 := range chars {
			for _, ch3 := range chars {
				for _, ch4 := range chars {
					wg.Add(1)                                                     // 新增1个任务
					go sha(string(ch1) + string(ch2) + string(ch3) + string(ch4)) // 新建goroutine执行sha
				}
			}
		}
	}
	wg.Wait() // 阻塞，等待所有任务执行完毕
	end := time.Since(start)
	fmt.Println(end)
}
