---
title: Decorator Pattern
date: 2023-12-15 13:48:08
tags:
categories:
- Design Pattern
keywords:
- Design Pattern
- Golang
copyright: Guader
copyright_author_href:
copyright_info:
---

装饰器模式是一种设计模式，允许将行为动态的添加到单个对象，而不影响同类中其他对象的行为。

提供了子类化之外的灵活选择。


## Example

在Golang中，可以使用接口和匿名函数来实现。

{% gist d8176075eeaf7a4075a0cfff1a2ea4af decorator.go %}

上面的代码中，我们定义了`Printer`接口以及实现了接口的结构体`SimplePrinter`。

然后，我们定义了`BoldDecorator`函数，接收一个`Printer interface`并返回一个`Printer interface`。将原来的`Print()`方法封装到一个新的方法中，返回用`<b>`括起来的新字符串。

