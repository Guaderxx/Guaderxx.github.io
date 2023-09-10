---
title: string-concatenation-benchmark
date: 2023-09-10 12:14:46
tags:
- Golang
categories:
- Benchmark
keywords:
copyright: Guader
copyright_author_href:
copyright_info:
---
 

# 字符串拼接的基准测试

结果很明显，只有两个字符串时直接相加`+`还不错，数量到10个开始基本都是`strings.Builder`了。


```
goos: linux
goarch: amd64
pkg: github.com/Guaderxx/goproj/pkg/teststr
cpu: AMD Ryzen 7 5800H with Radeon Graphics         
BenchmarkStringCombine2/Plain-16  	                47900655	        25.91 ns/op	       2 B/op	       1 allocs/op
BenchmarkStringCombine2/bytes.Buffer-16         	20379310	       100.1 ns/op	      66 B/op	       2 allocs/op
BenchmarkStringCombine2/fmt.Sprintf-16          	 5183205	       245.4 ns/op	      50 B/op	       4 allocs/op
BenchmarkStringCombine2/strings.Builder-16      	55521981	        48.14 ns/op	       8 B/op	       1 allocs/op

BenchmarkStringCombine10/Plain-16               	 3460802	       316.2 ns/op	      64 B/op	       9 allocs/op
BenchmarkStringCombine10/bytes.Buffer-16        	 6272886	       190.1 ns/op	      80 B/op	       2 allocs/op
BenchmarkStringCombine10/fmt.Sprintf-16         	  664045	      1825 ns/op	     368 B/op	      28 allocs/op
BenchmarkStringCombine10/strings.Builder-16     	 8755777	       139.7 ns/op	      24 B/op	       2 allocs/op

BenchmarkStringCombine100/Plain-16              	  170568	      7688 ns/op	    9744 B/op	      99 allocs/op
BenchmarkStringCombine100/bytes.Buffer-16       	 1000000	      1191 ns/op	     640 B/op	       4 allocs/op
BenchmarkStringCombine100/fmt.Sprintf-16        	   46423	     30486 ns/op	   12936 B/op	     298 allocs/op
BenchmarkStringCombine100/strings.Builder-16    	 1000000	      1004 ns/op	     504 B/op	       6 allocs/op

BenchmarkStringCombine1000/Plain-16             	    6036	    193308 ns/op	 1494036 B/op	     999 allocs/op
BenchmarkStringCombine1000/bytes.Buffer-16      	  204086	     13486 ns/op	   11200 B/op	       8 allocs/op
BenchmarkStringCombine1000/fmt.Sprintf-16       	    2751	    505771 ns/op	 1528576 B/op	    3000 allocs/op
BenchmarkStringCombine1000/strings.Builder-16   	  241279	     10227 ns/op	    8440 B/op	      11 allocs/op

BenchmarkStringCombine10000/Plain-16            	      54	  21088288 ns/op	204471898 B/op	   10015 allocs/op
BenchmarkStringCombine10000/bytes.Buffer-16     	   10000	    152557 ns/op	  171968 B/op	      12 allocs/op
BenchmarkStringCombine10000/fmt.Sprintf-16      	      45	  31473111 ns/op	205650848 B/op	   30187 allocs/op
BenchmarkStringCombine10000/strings.Builder-16  	   15818	     78155 ns/op	  154360 B/op	      20 allocs/op
PASS
ok  	github.com/Guaderxx/goproj/pkg/teststr	39.116s
```


{% gist 1050a41afa9aa7160c854cdc9318fb66 %}
