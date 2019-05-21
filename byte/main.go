package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

var (
	chars     = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")
	tail      = []byte("DEHNCHyUEO8kVZBT")
	result, _ = hex.DecodeString("3354de5346a962dd0f344de80cd3c8e5c2d3ce1a18437141b6a645df9b357c91")
)

func sha(head []byte) {
	h := sha256.New()
	h.Write(head)
	h.Write(tail)
	if bytes.Equal(h.Sum(nil), result) {
		fmt.Println(string(head))
	}
}

func main() {
	start := time.Now()
	for _, ch1 := range chars {
		for _, ch2 := range chars {
			for _, ch3 := range chars {
				for _, ch4 := range chars {
					sha([]byte{ch1, ch2, ch3, ch4})
				}
			}
		}
	}
	end := time.Since(start)
	fmt.Println(end)
}
