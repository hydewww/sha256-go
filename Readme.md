# Go语言爆破Sha256前四字符

>  sha256(XXXX+DEHNCHyUEO8kVZBT) == 3354de5346a962dd0f344de80cd3c8e5c2d3ce1a18437141b6a645df9b357c91

做 CTF 时需要求 sha256 原字符串的前四位，最近正好在学 golang ，就打算用 golang 实现一下。期间总觉得 golang 应该更快，就开始优化代码，从一开始的 9.5 秒优化到最后的 1.5 秒。

**下面的代码为了测试性能，都没有在算出结果后就结束程序。**

## easy

第一版代码如下。*sha256参考： https://godoc.org/crypto/sha256* 

```go
// easy
var (
  chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789" // A-Z a-z 0-9
	tail   = "DEHNCHyUEO8kVZBT" // 原字符串的尾部
	result = "3354de5346a962dd0f344de80cd3c8e5c2d3ce1a18437141b6a645df9b357c91" // hash值
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
```

**运行时间 9.5 秒** 

## goroutine 并发

因为第一版对比 python 没有明显性能优势，因此尝试用 goroutine 并发执行。*WaitGroup参考：https://golang.org/pkg/sync/#example_WaitGroup*

```go
// ez_goroutine
var wg sync.WaitGroup

func sha(...) {
  ...
  wg.Done()			// 完成任务
}
func main () {
  ...
  for ... {
    for ... {
      for ... {
        for ... {
          wg.Add(1)		// 新增1个任务
          go sha(...)	// 新建goroutine执行sha
        }
      }
    }
  }
  wg.Wait()			// 阻塞，等待所有任务执行完成
  ...
}
```

**运行时间 7.1 秒** 

## string -> byte

感觉还是不够快，于是开始做性能测试。

在单线程代码的基础上测量 sha256 以外的运算所使用的时间

```go
// test/nohash
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
...
```

运行时间 4.3 秒。由于 sha256 运算使用的都是 byte 数组，而代码输入输出中使用的都是字符串，存在很多的强制类型转换，因此占用了较长时间。下面将代码运算改为 byte 类型

```go
// byte
var (
	chars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")
	tail = []byte("DEHNCHyUEO8kVZBT")
	result, _ = hex.DecodeString("3354de5346a962dd0f344de80cd3c8e5c2d3ce1a18437141b6a645df9b357c91")
)

func sha(head []byte) {
	...
	h.Write(head)
  h.Write(tail) // hash特性 H(s1+s2) = H(H(s1)+s2)
  if bytes.Equal(h.Sum(nil), result) { // byte数组比较使用 bytes.Equal()
		fmt.Println(string(head))
	}
}
func main() {
  ...
	for ... {
		for ... {
			for ... {
				for ... {
					sha([]byte{ch1, ch2, ch3, ch4})
				}
			}
		}
	}
  ...
}
```

**运行时间减少为 4.8 秒**。

## 并发优化

将 `goroutine` 代码也同上优化（ `byte_goroutine` ），运行时间 5.9 秒 ，虽然快了 1 秒多但却比单线程慢。

猜测是 goroutine 分配使用的时间，将 sha 运算忽略进行测试。

```go
// test/goroutine
func sha(head []byte) {
	wg.Done()
}

func main() {
	...
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
	...
}
```

运行时间 5.2s， 确实是 goroutine 的分配和回收占用了大量时间。

根据 cpu 核心数分配相应数量的 goroutine，才能使效率更高，如下

```go
// final_goroutine
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
}
```

**运行时间 1.5 秒** ，舒服了！

想到一开始的 `goroutine` 运行时间长应该不完全是因为类型转换，使用这个分配逻辑（ `re_goroutine` ）以后运行时间为 3 秒。