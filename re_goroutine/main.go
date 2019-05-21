package main

import (
	"crypto/sha256"
	"fmt"
	"runtime"
	"sync"
	"time"
)

var (
	chars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789"
	tail   = "DEHNCHyUEO8kVZBT"
	result = "3354de5346a962dd0f344de80cd3c8e5c2d3ce1a18437141b6a645df9b357c91"
	wg     sync.WaitGroup
)

func sha(s string) {
	for _, ch1 := range s {
		for _, ch2 := range chars {
			for _, ch3 := range chars {
				for _, ch4 := range chars {
					head := string(ch1) + string(ch2) + string(ch3) + string(ch4)
					h := sha256.New()
					h.Write([]byte(head + tail))
					str := fmt.Sprintf("%x", h.Sum(nil))
					if str == result {
						fmt.Println(head)
					}
				}
			}
		}
	}
	wg.Done() // 完成任务
}

func main() {
	threads := runtime.NumCPU()
	start := time.Now()

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
