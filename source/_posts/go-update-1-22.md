---
title: go update 1.22
date: 2024-02-11 12:26:51
tags:
- doc
categories:
- Golang
- doc
keywords:
copyright: Guader
copyright_author_href:
copyright_info:
---


最近 Golang 更新到 1.22 了，这里简单介绍下。
我目前是 1.21.5，正好作为对比。


## `for-loop` 变量问题

之前的 `range` 只会生成一个变量，现在每次迭代都会生成新变量。

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    names := []string{"Jane", "Max", "Joe"}
    
    for _, name := range names {
        go func() {
            fmt.Println(name)
        }()
    }
    
    time.Sleep(time.Second)
}

// go version: 1.21.5
// output:
// Joe
// Joe
// Joe

// go version: 1.22
// output:
// Joe
// Max
// Jane
// 这里顺序不一定，不过保证会输出所有值
```


以及可以 `range` 整数了。


## 更好的标准库 HTTP Routing

之前的 `http.ServeMux` 只接受常规路径，不接受参数和方法。


### 比如指定方法的路由

```go
// go version: 1.21.5
// Only accept HTTP GET method.
package main

import (
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    
    mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            w.WriteHeader(http.StatusMethodNotAllowed)
            w.Write([]byte(""))
            return
        }
        
        w.Write([]byte("Hello"))
    })
    
    if err := http.ListenAndServe(":8080", mux); err != nil {
        panic(err)
    }
}
```


```bash
curl -v -X POST http://localhost:8080/hello

# Output:
# ...
# < HTTP/1.1 405 Method Not Allowed
# < Date: Sun, 11 Feb 2024 07:25:34 GMT
# < Content-Length: 0
# <
# * Connection #0 to host localhost left intact


curl -v -X GET http://localhost:8080/hello

# Output:
# ...
# < HTTP/1.1 200 OK
# < Date: Sun, 11 Feb 2024 07:25:47 GMT
# < Content-Length: 5
# < Content-Type: text/plain; charset=utf-8
# <
# * Connection #0 to host localhost left intact
# Hello%
```


```go
// go version: 1.22
package main

import (
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    
    mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello"))
    })
    
    if err != http.ListenAndServe(":8080", mux); err != nil {
        panic(err)
    }
}
```

同样的行为，现在是更好，更灵活了。
只接受指定的 `GET` 方法，其他的会返回 405。



### 路径中的参数 - 通配符

路径参数也算是个常见的需求，这也是很多人选择 `Gorilla` 或是其他库的原因。

路径参数是 URL 的一部分，需要请求中的值，比如说我们希望允许用户添加一个名字来打招呼，
那么请求是 `/hello/$NAME`。

旧的处理方式是：

```go
package main

import (
    "net/http"
    "fmt"
    "strings"
)

func main() {
    mux := http.NewServeMux()
    
    mux.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
        path := r.URL.Path
        
        parts := strings.Split(path, "/")
        
        if len(parts) < 3 {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }
        
        name := parts[2]
        
        w.Write([]byte(fmt.Sprintf("Hello %s!", name)))
    })
    
    if err := http.ListenAndServe(":8080", mux); err != nil {
        panic(err)
    }
}
```

非常乱，而且这只是一个变量。

新的 `ServeMux` 允许将参数包裹在 `{}` 中来指定名称。
因此，可以通过将路由设置为 `/hello/{name}` 来添加名称参数。

为了能获取参数，HTTP 请求包含一个名为 `PathValue` 的函数。

现在：

```go
package main

import (
    "net/http"
    "fmt"
)

func main() {
    mux := http.NewServeMux()
    
    mux.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) {
        name := r.PathValue("name")
        
        w.Write([]byte(fmt.Sprintf("Hello %s!", name)))
    })
    
    if err := http.ListenAndServe(":8080", mux); err != nil {
        panic(err)
    }
}
```

- 如果发送不带参数的请求，返回 HTTP 404
- 带多个参数，例如 `:8080/hello/ader/max`
  - 这种就需要更改通配符，`"GET/ hello/{name...}"`，通过增加 `...` 匹配所有后续参数



### 使用尾部斜杠匹配精确模式

如果路由是 `hello/` ， 那么会对任何以 `hello/` 开头的路由进行匹配。
这里的关键就是末尾的 `/` 。

可以在末尾增加 `{$}` 来进行精确匹配。
`GET /hello/{$}`。 这样后续再在路径末尾增加参数也不会匹配了。


**冲突解决和优先顺序**

随着这些新规则的出现，我们的新问题是一个请求可以匹配多个路由。

比如：

- `/hello`
- `/hello/{name}`

这里是通过选择最具体的路线来解决的。



## Others

- Go 中第一个 v2 包 - `math/rand/v2` 
  移除了 `Read` 方法
  更快的算法
  新的 `rand.N` 函数
- slog 增加了 `SetLogLoggerLevel` 
- slices 增加了新的 `Concat` 函数
- ...
