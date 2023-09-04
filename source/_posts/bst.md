---
title: bst
date: 2023-09-04 21:19:52
tags:
- Tree
categories:
- Data Structure
- Algorithm
keywords:
copyright: Guader
copyright_author_href:
copyright_info:
---

# Binary Search Tree

> 今天的每日一题是BST相关的，说实话，只记得叫`Binary Search Tree`，树相关的。  
> 我寻思整理一下吧。


## Summary

定义：

- 若任意节点的左子树不空，则左子树上所有节点的值均小于根节点的值。
- 若任意节点的右子树不空，则右子树上所有节点的值均大于根节点的值。
- 任意节点的左，右子树也分别为BST


相比其他数据结构的优势在于查找，插入的时间复杂度较低。  
是基础性数据结构，用于构建更抽象的数据结构：集合，多重集，关联数组等。



## 时间复杂度

```
算法        平均        最差
空间        O(n)        O(n)
搜索        O(log n)    O(n)
插入        O(log n)    O(n)
删除        O(log n)    O(n)
```


## 查找

```go
type BSTree struct {
    Val int
    Left *BSTree
    Right *BSTree
}

func SearchBST(t *BSTree, v int) *BSTree {
    if t == nil {
        return nil
    }
    if t.Val == v {
        return t
    }
    if v < t.Val {
        return SearchBST(t.Left, v)
    }
    if v > t.Val {
        return SearchBST(t.Right, v)
    }
}
```


## 插入

```go
func InsertBST(t *BSTree, v int) (*BSTree, error) {
    if t == nil {
        return &BSTree{ Val: v }, nil
    }
    if t.Val == v {
        return t, errors.New("Invalid Insert Value")
    }
    if t.Val > v {
        t.Left = InsertBST(t.Left, v)
    }
    if t.Val < v {
        t.Right = InsertBST(t.Right, v)
    }
    return t
}
```


## 删除

```go
# 叶子节点或有单侧子树直接删了改引用，不破坏树结构
# 若左右子树均存在
#    1. 

func DeleteBST(t *BSTree, v int) *BSTree {
    if t == nil {
        return nil
    }
    if t.Val == v {
        if t.Left == nil && t.Right == nil {
            return nil
        } else if t.Left == nil {
            return t.Right
        } else if t.Right == nil {
            return t.Left
        } else {
            tmp := t
            r := t.Left
            for r.Right != nil {
                tmp = r
                r = r.Right
            }
            t.Val = r.Val
            if tmp == t {
                t.Left = r.Left
            } else {
                tmp.Right = r.Left
            }
            return t
        }
    }
    if t.Val > v {
        t.Left = DeleteBST(t.Left, v)
    }
    if t.Val < v {
        t.Right = DeleteBST(t.Right, v)
    }
    return t
}
```


## 创建

如果是单支树性能很差，所以还有平衡树


## Daily Question

> 449. Serialize and Deserialize BST  
> 我想着根据前序，中序解析就完了，没想到这么差，明天再换个解法

{% asset_img photo.png %}
