---
title: range-over-func
date: 2024-08-19 12:01:33
tags:
- Doc
categories:
- Golang
- Doc
keywords: 
- Golang
copyright: Guader
copyright_author_href:
copyright_info:
---

## Background

[Discussion-54245: standard iterator interface][iterator]

关于 Iterator 的提案，包括为什么现在提议这个，以及希望提供的功能。

[Discussion-56413: user-defined iteration using range over func values][range-over-func-values]

关于 `for range` 和 Push/Pull 函数相关 （这个讨论该看）


[Discussion-56010: redefining for loop variable semantics][loopVar]

Loop 变量作用域，1.22 已经包含了

[Discussion-43557: function values as iterators][43557]
[Discussion-43557-comment][43557comment]

这两个更早点。也是 Iterator 相关。

[Issue-61897: Iter pakcage][IterPkg]

关于 1.23 新增的 `iter` 包及其相关功能。



## Simple Summary

大多数语言提供了一种标准化方法来使用迭代器接口迭代存储在容器中的值。  
Go 提供了与 map、slice、string、array 和 channel 一起使用的 `for range` ，但它没有为用户编写的容器提供任何通用机制，也没有提供迭代器接口。   
这导致 Go 相关的非泛型迭代器的用法五花八门：

- [`runtime.CallersFrames`][runtimeCallersFrames]  返回一个在堆栈帧上迭代的 `runtime.Frames` ； `Frames` 有一个 `Next` 方法，它返回一个 `Frame` 和一个报告是否有更多帧的 bool 值（也就是下次调用 `Next` 方法是否有值返回）。
- [`bufio.Scanner`][bufioScanner]  是一个通过 `io.Reader` 的迭代器，其中 `Scan` 方法前进到下一个值。该值由 `Bytes` 方法返回。错误由 `Err` 方法收集并返回。
- [`database/sql.Rows`][databaseSqlRows]  迭代查询的结果，其中 `Next` 方法前进到下一个值，并且该值由 `Scan` 方法返回。 `Scan` 方法可能会返回错误。
- [`archive/tar.Reader.Next`][archiveNext]
- [`bufio.Reader.ReadByte`][bufioReadByte] 
- [`bufio.Scanner.Scan`][bufioScan] 
- [`container/ring.Ring.Do`][ringDo] 
- [`expvar.Do`][expvarDo] 
- [`flag.Visit`][flagVisit]
- [`go/token.FileSet.Iterate`][tokenIterate] 
- [`path/filepath.Walk`][filepathWalk]
- [`sync.Map.Range`][mapRange]

部分原因是在引入泛型之前，无法编写描述迭代器的接口。 

不过现在有泛型了，我们可以为具有 `E` 类型元素的容器上的迭代器编写一个接口 `Iter[E]` 。  
其他语言中迭代器的存在表明这是一个强大的工具。


### Push/Pull functions

[#56413][range-over-func-values] 讨论中关于 push/pull 函数的功能讨论的很多，包括它们的互相转换。   
不过和 iter 包现有的并不完全一样，大体概念倒是没变。


## Iter Package

目前来说做的不多，两个类型和两个函数签名

```go
type Seq[V any] func(yield func(V) bool)
type Seq2[K, V any] func(yield func(K, V) bool)

func Pull[V any](seq Seq[V]) (next func() (V, bool), stop func())
func Pull2[K, V any](seq Seq2[K, V]) (next func() (K, V, bool), stop func())
```

slices 和 maps 包中也添加了相关的一些函数，这里随便挑一个看下

```go
// All returns an iterator over index-value pairs in the slice
// in the usual order.
func All[Slice ~[]E, E any](s Slice) iter.Seq2[int, E] {
    return func(yield func(int, E) bool) {
        for i, v := range s {
            if !yield(i, v) {
                return
            }
        }
    }
}
```

用法大概就是：

```go
nums := []int{2,3,4,5}
for _, v := range slices.All(nums) {
    fmt.Println(v)
}
// 2
// 3
// 4
// 5
```

尝试写一个吧，比如获取一个 slice 中的奇数：

```go
func Odd[Slice ~[]E, E int](s Slice) iter.Seq2[int, E] {
    return func(yield func(int, E) bool) {
        for i, v := range s {
            if v % 2 != 0 {
                if !yield(i, v) {
                    return
                }
            }
        }
    }
}
```

这个例子不算太好，因为 `v % 2 != 0` 这一步约束了这个类型必须是整数，这里我就简单的标注为 `E int` 了。

现在就可以用了：

```go
nums := []int{2,3,4,5,6}
for _, v := range Odd(nums) {
    fmt.Println(v)
}
// 3
// 5
```

但这么写是有点奇怪的，挺奇怪的，我们可以直接操作 `yield` 函数

```go
slices.All(nums)(func (k, v int) bool {
    if v % 2 != 0 {
        fmt.Println(v)
    }
    return true
})
// 3
// 5
```

然后是 PULL 函数

比如说写一个交换 map k,v 对的函数

```go
func ReplaceKV[K comparable](seq iter.Seq2[K, K]) iter.Seq2[K, K] {
    return func(yield func(K, K) bool) {
        next, stop := iter.Pull2(seq)
        defer stop()
        for {
            k, v, ok := next()
            if !ok {
                return
            }
            if !yield(v, k) {
                return
            }
        }
    }
}
```

```go
m := map[string]string{"cn": "CHINA", "en": "ENGLISH", "us":"AMERICA"}
for k, v := range maps.All(m) {
    fmt.Println(k, ": ", v)
}
fmt.Println("---")
for k, v := range ReplaceKV(maps.All(m)) {
    fmt.Println(k, ": ", v)
}
// us :  AMERICA
// cn :  CHINA
// en :  ENGLISH
// ---
// ENGLISH :  en
// AMERICA :  us
// CHINA :  cn
```

这里是只用了 slice 和 map 举例，所以更像个语法糖。  
应该和结构体或者说复合类型结合起来看，效果会更好，我就不继续写了。  

而且这确实是规范了 Iterator 的写法。



[range-over-func-values]: https://github.com/golang/go/discussions/56413
[iterator]: https://github.com/golang/go/discussions/54245
[loopVar]: https://github.com/golang/go/discussions/56010
[43557]: https://github.com/golang/go/issues/43557
[43557comment]: https://github.com/golang/go/issues/43557#issuecomment-895211452
[IterPkg]: https://github.com/golang/go/issues/61897
[runtimeCallersFrames]: https://pkg.go.dev/runtime/#Frames.Next
[databaseSqlRows]: https://pkg.go.dev/database/sql#Rows
[bufioScanner]: https://pkg.go.dev/bufio#Scanner
[archiveNext]: https://pkg.go.dev/archive/tar/#Reader.Next
[bufioReadByte]: https://pkg.go.dev/bufio/#Reader.ReadByte
[bufioScan]: https://pkg.go.dev/bufio/#Scanner.Scan
[ringDo]: https://pkg.go.dev/container/ring/#Ring.Do
[expvarDo]: https://pkg.go.dev/expvar/#Do
[flagVisit]: https://pkg.go.dev/flag/#Visit
[tokenIterate]: https://pkg.go.dev/go/token/#FileSet.Iterate
[filepathWalk]: https://pkg.go.dev/path/filepath/#Walk
[mapRange]: https://pkg.go.dev/sync/#Map.Range
