package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

var (
	chars     = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")
	tail      = []byte("DEHNCHyUEO8kVZBT")
	result, _ = hex.DecodeString("3354de5346a962dd0f344de80cd3c8e5c2d3ce1a18437141b6a645df9b357c91")
	wg        sync.WaitGroup
)

func sha(head []byte) {
	h := sha256.New()
	h.Write(head)
	h.Write(tail)
	if bytes.Equal(h.Sum(nil), result) {
		fmt.Println(string(head))
	}
	wg.Done()
}

func main() {
	start := time.Now()
	for _, ch1 := range chars {
		for _, ch2 := range chars {
			for _, ch3 := range chars {
				for _, ch4 := range chars {
					wg.Add(1)
					go sha([]byte{ch1, ch2, ch3, ch4})
				}
			}
		}
	}
	wg.Wait()
	end := time.Since(start)
	fmt.Println(end)
}
