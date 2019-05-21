package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

var (
	// A-Z a-z 0-9
	chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789"
	// 原字符串的尾部
	tail = "DEHNCHyUEO8kVZBT"
	// hash值
	result = "3354de5346a962dd0f344de80cd3c8e5c2d3ce1a18437141b6a645df9b357c91"
)

func sha(head string) {
	h := sha256.New()
	h.Write([]byte(head + tail))
	str := fmt.Sprintf("%x", h.Sum(nil))
	if str == result {
		fmt.Println(head)
	}
}

func main() {
	start := time.Now()
	for _, ch1 := range chars {
		for _, ch2 := range chars {
			for _, ch3 := range chars {
				for _, ch4 := range chars {
					sha(string(ch1) + string(ch2) + string(ch3) + string(ch4))
				}
			}
		}
	}
	end := time.Since(start)
	fmt.Println(end)
}
