package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"runtime"
	"sync"
	"time"
)

var (
	chars     = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")
	tail      = []byte("DEHNCHyUEO8kVZBT")
	result, _ = hex.DecodeString("3354de5346a962dd0f344de80cd3c8e5c2d3ce1a18437141b6a645df9b357c91")
	wg        sync.WaitGroup
)

func sha(s []byte) {
	for _, ch1 := range s {
		for _, ch2 := range chars {
			for _, ch3 := range chars {
				for _, ch4 := range chars {
					head := []byte{ch1, ch2, ch3, ch4}
					h := sha256.New()
					h.Write(head)
					h.Write(tail)
					if bytes.Equal(h.Sum(nil), result) {
						fmt.Println(string(head))
					}
				}
			}
		}
	}
	wg.Done()
}

func main() {
	threads := runtime.NumCPU() // 获取cpu逻辑核心数（包括超线程）
	start := time.Now()

	/* len(chars) = sum * sthreads + (sum+1) * (threads-sthreads) */
	snum := len(chars) / threads
	sthreads := threads*(1+snum) - len(chars)

	wg.Add(threads)
	for i := 0; i < threads; i++ {
		if i < sthreads {
			go sha(chars[snum*i : snum*(i+1)])
		} else {
			base := snum * sthreads
			go sha(chars[base+(snum+1)*(i-sthreads) : base+(snum+1)*(i-sthreads+1)])
		}
	}
	wg.Wait()
	end := time.Since(start)
	fmt.Println(end)
}
