---
title: go-slice
date: 2024-04-10 21:28:26
tags:
- Go
categories:
- Go
- Slice
keywords:
copyright: Guader
copyright_author_href:
copyright_info:
---

看了 Rust 的智能指针我就在想 Go 中的 Slice 实现，因为一直记的是 `append` 超量会重新分配一个新底层数组，所以之前的 slice 地址会变，所以要 `s = append(s, v)` 。

结果和我记的有点误差，这个分配在栈上的 slice struct 地址是没变的，类似于实现了 `Deref trait` 直连到底层数组。

```go

func main() {
	arr := []int{1}
	arrp := &arr
	fmt.Printf("arrp = %p\n", arrp)
	fmt.Println(*arrp)
	fmt.Printf("arr pointer: %p\n\n", arr)

	arr = append(arr, 1, 2, 3)
	fmt.Printf("arrp = %p\n", arrp)
	fmt.Println(*arrp)
	fmt.Printf("arr pointer: %p\n\n", arr)

	arr = append(arr, 11, 22, 33)
	fmt.Printf("arrp = %p\n", arrp)
	fmt.Println(*arrp)
	fmt.Printf("arr pointer: %p\n\n", arr)
    
    _ = append(arr, 9, 8, 7, 6, 5, 4, 3, 2, 1)
	fmt.Printf("arrp = %p\n", arrp)
	fmt.Println(*arrp)

	*arrp = append(arr, 9, 8, 7, 6, 5, 4, 3, 2, 1)
	fmt.Printf("arrp = %p\n", arrp)
	fmt.Println(*arrp)
}

/*
arrp = 0xc00012c000
[1]
arr pointer: 0xc00011a010

arrp = 0xc00012c000
[1 1 2 3]
arr pointer: 0xc000136000

arrp = 0xc00012c000
[1 1 2 3 11 22 33]
arr pointer: 0xc000138000

arrp = 0xc00012c000
[1 1 2 3 11 22 33]
arrp = 0xc00012c000
[1 1 2 3 11 22 33 9 8 7 6 5 4 3 2 1]
*/
```
