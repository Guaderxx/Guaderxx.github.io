---
title: forloop-performance
date: 2023-09-20 19:12:42
tags:
- Golang
categories:
- Benchmark
keywords:
copyright: Guader
copyright_author_href:
copyright_info:
---


# 嵌套循环时不同顺序对性能的影响

> 看CSAPP第一课讲了这个，就做了个Go的测试。只能说确实每留意过这个会带来性能提升。
> （但这个具体是和一级缓存有关还是什么的等到后面再补充。）


{% gist 6c14f78a0949b2355aaf51f00a8ce5d0 result.txt %}


{% gist 6c14f78a0949b2355aaf51f00a8ce5d0 main.go %}


{% gist 6c14f78a0949b2355aaf51f00a8ce5d0 main_test.go %}
