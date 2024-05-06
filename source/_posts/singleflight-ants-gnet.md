---
title: singleflight|ants|gnet
date: 2024-05-06 15:12:55
tags:
- Golang
categories:
- Golang
- pkg
keywords:
copyright: Guader
copyright_author_href:
copyright_info:
---


## Singleflight

[singliflight](https://cs.opensource.google/go/x/sync/+/master:singleflight/): 提供重复函数调用抑制机制

前两天有人单飞，然后说什么 singlefly （我听叉了），~~我寻思双飞我知道，单飞是啥~~

看了一下，还挺不错的

```go
// copy from singleflight.go
// ...
// call is an in-flight or completed singleflight.Do call
type call struct {
	wg sync.WaitGroup

	// These fields are written once before the WaitGroup is done
	// and are only read after the WaitGroup is done.
	val interface{}
	err error

	// These fields are read and written with the singleflight
	// mutex held before the WaitGroup is done, and are read but
	// not written after the WaitGroup is done.
	dups  int
	chans []chan<- Result
}

// Group represents a class of work and forms a namespace in
// which units of work can be executed with duplicate suppression.
type Group struct {
	mu sync.Mutex       // protects m
	m  map[string]*call // lazily initialized
}

// ...

func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		c.dups++
		g.mu.Unlock()
		c.wg.Wait()

		if e, ok := c.err.(*panicError); ok {
			panic(e)
		} else if c.err == errGoexit {
			runtime.Goexit()
		}
		return c.val, c.err, true
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	g.doCall(c, key, fn)
	return c.val, c.err, c.dups > 0
}

// ...

// doCall handles the single call for a key.
func (g *Group) doCall(c *call, key string, fn func() (interface{}, error)) {
	normalReturn := false
	recovered := false

	// use double-defer to distinguish panic from runtime.Goexit,
	// more details see https://golang.org/cl/134395
	defer func() {
		// the given function invoked runtime.Goexit
		if !normalReturn && !recovered {
			c.err = errGoexit
		}

		g.mu.Lock()
		defer g.mu.Unlock()
		c.wg.Done()
		if g.m[key] == c {
			delete(g.m, key)
		}

		if e, ok := c.err.(*panicError); ok {
			// In order to prevent the waiting channels from being blocked forever,
			// needs to ensure that this panic cannot be recovered.
			if len(c.chans) > 0 {
				go panic(e)
				select {} // Keep this goroutine around so that it will appear in the crash dump.
			} else {
				panic(e)
			}
		} else if c.err == errGoexit {
			// Already in the process of goexit, no need to call again
		} else {
			// Normal return
			for _, ch := range c.chans {
				ch <- Result{c.val, c.err, c.dups > 0}
			}
		}
	}()

	func() {
		defer func() {
			if !normalReturn {
				// Ideally, we would wait to take a stack trace until we've determined
				// whether this is a panic or a runtime.Goexit.
				//
				// Unfortunately, the only way we can distinguish the two is to see
				// whether the recover stopped the goroutine from terminating, and by
				// the time we know that, the part of the stack trace relevant to the
				// panic has been discarded.
				if r := recover(); r != nil {
					c.err = newPanicError(r)
				}
			}
		}()

		c.val, c.err = fn()
		normalReturn = true
	}()

	if !normalReturn {
		recovered = true
	}
}
```

上面几乎是所有逻辑了，根据 `key` 设置一个内存缓存，类似于闭包的 singleton，并在函数完成（包括 panic ）后删除 `key`。

很巧妙且有效；不过根据 `fn` 的调用时间，`key` 的设置需要注意，如果调用时间稍长，应该会有 *数据一致性* 的问题。


## Ants

[ants](https://github.com/panjf2000/ants): Goroutine 池

**Goroutine 虽然轻量但不应该无限制的使用**

不过这个库的 benchmark test 结果不尽如人意，而在我将

```go
// ants_benchmark_test.go
func demoFunc() {
	time.Sleep(time.Duration(BenchParam) * time.Millisecond)
}
```

改为 json decode 后，基准测试又一直失败。

~~可能这个库更适合 IO 密集型的系统~~
~~目前我更多依靠手动处理这些问题，不算是个好习惯.~~ 可以少量试用吧


## Gnet

[gnet](https://github.com/panjf2000/gnet): 高性能、轻量级、非阻塞的事件驱动 Go 网络框架

和 ants 是一个作者， ~~好像也在公众号看过~~ ，顺便看了下。

看完它在 [TechEmpower 上的基准测试](https://www.techempower.com/benchmarks/#hw=ph&test=plaintext&section=data-r22) 就没往后面看。

七个测试中 plaintext 毫无争议的排在了第一位，希望早日补全剩下的测试。

---

这个测试挺有意思的：

综合评价中前三都是 Rust 的，第一是 [ntex](https://github.com/ntex-rs/ntex)

**Rust 天下第一**
