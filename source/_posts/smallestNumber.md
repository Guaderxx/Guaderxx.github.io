---
title: smallestNumber
date: 2023-09-05 22:04:43
tags:
categories:
- Algorithm
keywords:
copyright: Guader
copyright_author_href:
copyright_info:
---

## 力扣每日一题

> 2605: Form Smallest Number From Two Digit Arrays      easy

虽然是简单题，错了很多遍，感觉脑子退化了

```go
func MinNumber(nums1 []int, nums2 []int) int {
    exists := make(map[int]bool)
    min1 := nums1[0]
    res := 11
    for i := 0; i < len(nums1); i++ {
        exists[nums1[i]] = true
        if min1 > nums1[i] {
            min1 = nums1[i]
        }
    }
    min2 := nums2[0]
    for i := 0; i < len(nums2); i++ {
        if exists[nums2[i]] && res > nums2[i] {
            res = nums2[i]
        }
        if min2 > nums2[i] {
            min2 = nums2[i]
        }
    }
    if res != 11 {
        return res
    }
    if min1 > min2 {
        return min2 * 10 + min1
    }
    return min1 * 10 + min2
}
```


---

{% asset_img photo.png %}
