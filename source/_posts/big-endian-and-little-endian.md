---
title: big-endian and little-endian
date: 2024-10-21 13:34:23
tags:
categories:
keywords:
- endianness
copyright: Guader
copyright_author_href:
copyright_info:
---

- 大端序（Big Endian）：高字节存储在低地址，低字节存储在高地址。  
- 小端序（Little Endian）：低字节存储在低地址，高字节存储在高地址。

- 高字节： 一个多字节数据中，数值较大的字节。比如一个16位的整数，它的高字节代表了数值的较高位部分。
- 低字节： 一个多字节数据中，数值较小的字节。比如一个16位的整数，它的低字节代表了数值的较低位部分。

- 高地址： 内存中数值较大的地址。
- 低地址： 内存中数值较小的地址。

以整数 300 为例。

```go
// 00000001 00101100
var num uint16 = 0x012C
// 可以用以下两个字节表示     
// 00101100 00000001
var num1, num2 uint8 = 0x01, 0x2C

// big-endian
// 低地址-高字节：0x01 高地址-低字节：0x2C

// little-endian
// 低地址-低字节：0x2C 高地址-高字节：0x01
```

## Example

```go
// 以解析 [3]byte 到一个 uint32 为例
// 有 UI24，但是 Go 中没这类型，因此用 uint32 表示
func BE_U24(n []byte) uint32 {
    return uint32(n[0]) << 16 |
        uint32(n[1]) << 8 |
        uint32(n[2])
}

// 小端就是反过来
func LE_U24(n []byte) uint32 {
    return uint32(n[2]) << 16 |
        uint32(n[1]) << 8 |
        uint32(n[0])
}
```

> 写解析器会用到，相同的字节切片按不同的字节序会解析出不同的数字。  
> 常用的在 `encoding/binary` 里有。
