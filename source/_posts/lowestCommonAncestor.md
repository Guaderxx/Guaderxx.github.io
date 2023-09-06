---
title: lowestCommonAncestor
date: 2023-09-06 23:10:55
tags:
categories:
- Algorithm
keywords:
copyright: Guader
copyright_author_href:
copyright_info:
---


## 力扣每日一题

> 1123. Lowest Common Ancestor of Deepest Leaves


很典型的一个树相关的题，倒序解，但其实解完我也不是很理解题意。

```go
func lcaDeepestLeaves(root *TreeNode) *TreeNode {
    if root.Left == nil && root.Right == nil {
        return root
    }

    var postorder func(*TreeNode) (*TreeNode, int)
    postorder = func(node *TreeNode) (*TreeNode, int) {
        if node == nil {
            return nil, 0
        }
        l, lDepth := postorder(node.Left)
        r, rDepth := postorder(node.Right)

        if lDepth == rDepth {
            return node, lDepth + 1
        } else if lDepth > rDepth {
            return l, lDepth + 1
        } else {
            return r, rDepth + 1
        }
    }

    res, _ := postorder(root)
    return res
}
```


---

{% asset_img photo.png %}
