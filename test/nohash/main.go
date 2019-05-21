package main

import (
	"fmt"
	"time"
)

var (
	tail   = "DEHNCHyUEO8kVZBT"
	result = "3354de5346a962dd0f344de80cd3c8e5c2d3ce1a18437141b6a645df9b357c91"
)

const (
	chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "1234567890"
)

func sha(head string) {
	var hash = []byte(head + tail)
	// h := sha256.New()
	// h.Write([]byte(head + tail))
	// hash := h.Sum(nil)
	str := fmt.Sprintf("%x", hash)
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
