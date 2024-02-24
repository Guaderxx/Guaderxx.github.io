---
title: daily note
date: 2024-02-24 12:13:13
tags:
categories:
- Daily
keywords:
- Note
copyright: Guader
copyright_author_href:
copyright_info:
---

近来一直在看乐子，没想到自己成乐子了，MotherFuck.

昨天做个小面试，做个题：

用两个 goroutine 分别输出 A 1 B 2 ... Z 26。

然后我满脑子都是：

```go
str := "A 1 B 2 ... Z 26"
```

过了十几分钟脑子里还是这个，属实被自己整笑了。

不过结束后我想了想，其实也不是不行，就是奇怪了点（反正这题也挺奇怪的）

> 如果要控制输出次序，那就要手动调度，最后本质上是个串行代码，不需要默认的 goroutine 调度。  
> 也不需要锁。这个概念就很适合单线程的 Js 的事件驱动，EventLoop。  
> 所以这个题的解法就可以乱七八糟的

比如：

```go
func main() {
    res := make(chan rune)
    wg := sync.WaitGroup{}
    wg.Add(2)
    
    go func() {
        defer close(res)
        defer wg.Done()
        
        var i rune
        for i = 1; i < 27; i++ {
            res <- (i + 64)
            res <- i
        }
    }()
    
    go func() {
        defer wg.Done()
        for val := range res {
            if val < 27 {
                fmt.Printf("%d ", val)
            } else if val > 64 {
                fmt.Printf("%c ", val)
            }
        }
    }()
    
    wg.Wait()
    fmt.Printf("\n")
}
```

这又何尝不是以两个 goroutine 交替输出呢.
