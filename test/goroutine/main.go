package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	chars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")
	wg    sync.WaitGroup
)

func sha(head []byte) {
	wg.Done()
}

func main() {
	start := time.Now()
	for range chars {
		for range chars {
			for range chars {
				for range chars {
					wg.Add(1)
					go sha([]byte{})
				}
			}
		}
	}
	wg.Wait()
	end := time.Since(start)
	fmt.Println(end)
}
