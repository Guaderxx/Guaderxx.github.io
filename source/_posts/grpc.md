---
title: grpc
date: 2023-08-12 16:17:05
tags:
categories:
- Intro
- gRPC
keywords:
copyright: Guader
copyright_author_href:
copyright_info:
---

# gRPC简介


## 生成代码

例如代码结构如下：

```bash
project
├── cmd
│   ├── client
│   └── server
├── go.mod
├── go.sum
├── internal
│   └── protos
│       └── person.proto
└── pkg
        └── pb
```

将`internal/protos/*.proto`编译到`pkg/pb`下，可写个脚本：

{% include_code Makefile lang:Makefile grpc/Makefile %}

生成后的文件结构为：

```bash
project
├── cmd
│   ├── client
│   └── server
├── go.mod
├── go.sum
├── internal
│   └── protos
│       └── person.proto
├── Makefile
└── pkg
    └── protos
        ├── person_grpc.pb.go
        └── person.pb.go
```


## 服务定义

和许多RPC系统一样，gRPC基于定义服务的思想，指定可以远程调用的方法及其参数和返回类型。默认情况下，gRPC使用protocol buffer作为接口定义语言IDL，来描述服务接口和有效负载消息的结构。如果需要，可以使用其他替代方案。

gRPC允许定义四种服务方法：

- 一元RPC，客户端向服务器发送单个请求并获取单个响应，就像正常的函数调用

`rpc SayHello(HelloRequest) returns (HelloResponse);`

- 服务端流式RPC。客户端向服务器发送请求，并获取一个流来读取一系列信息。gRPC保证了单个RPC调用中的消息排序

`rpc LotsOfReplies(HelloRequest) returns (stream HelloResponse)`

- 客户端流式RPC。客户端写入一系列数据并发送至服务器，同样使用提供的流。客户端写完消息后，等待服务器读取信息并返回响应。gRPC再次保证了单个RPC调用中的消息排序。

`rpc LotsOfGreetings(stream HelloRequest) returns (HelloResponse)`

- 双向流RPC。双方使用读写流发送一系列消息。这两个流独立运行，因此客户端和服务器可以按照自己喜欢的任何顺序读写：例如，服务器可以等到接受到所有客户端信息后再进行响应，也可以交替读其消息进行响应，或其他读写组合。每个消息流中的信息顺序都会被保留。

`rpc BidiHello(stream HelloRequest) returns (stream HelloResponse)`



## 同步与异步

同步 RPC 调用会阻塞，直到服务器发出响应，这与 RPC 所追求的存储过程调用抽象最接近。另一方面，网络本质上是异步的，在许多情况下，启动 RPC 时不阻塞当前线程是非常有用的。

大多数语言中的 gRPC 编程 API 都有同步和异步两种风格。您可以在每种语言的教程和参考文档中找到更多信息。


## 生成代码参考

**线程安全：** 请注意，客户端 RPC 调用和服务器端 RPC 处理程序是线程安全的，可以在并发程序上运行。但也要注意，对于单个数据流而言，传入和传出数据是双向的，但却是串行的；因此，单个数据流不支持并发读取或并发写入（但读取与写入是安全并发的）。


## 生成的服务器接口上的方法

在服务端侧，每个proto文件中的`service Bar`会生成：

`func RegisterBarServer(s *grpc.Server, srv BarServer)`

应用可以使用此函数定义BarServer接口的具体实现，并将其注册到`grpc.Server`实例（在启动服务器实例之前）。


**一元方法**

这些方法在生成的服务接口上有如下签名：

`Foo(context.Context, *MsgA) (*MsgB, error)`

这里，MsgA就是客户端发送的protobuf消息，MsgB是从服务器响应的protobuf消息。

**服务端流方法**

有如下签名：

`Foo(*MsgA, <ServiceName>_FooServer) error`

这里，MsgA是来自客户端的单个请求，`<ServiceName>_FooServer`参数表示服务器到客户端的MsgB消息流。

`<ServiceName>_FooServer`有一个嵌入的`grpc.ServerStream`和以下接口：

```go
type <ServiceName>_FooServer interface {
    Send(*MsgB) error
    grpc.ServerStream
}
```

服务端处理程序通过此参数的`Send`方法向客户端发送`protobuf`消息流。流结束是由处理程序方法的结束（`return`）。


**客户端流方法**

有如下签名：

`Foo(<ServiceName)_FooServer) error`

这里，`<ServiceName>_FooServer`既可以用于读取客户端到服务器的消息流，用用作发送单个服务器响应消息。

```go
type <ServiceName>_FooServer interface {
    SendAndClose(*MsgA) error
    Recv() (*MsgB, error)
    grpc.ServerStream
}
```

服务端处理程序可以重复调用`Recv`来完整的获取客户端发送的消息流。`Recv`返回值为`(nil, io.EOF)`，一旦获取完毕，就可以调用`SendAndClose`来响应消息。注意：`SendAndClose`必须且只能调用一次。


**双向流方法**

签名如下：

`Foo(<ServiceName_FooServer) error`

```go
type <ServiceName>_FooServer interface {
    Send(*MsgA) error
    Recv() (*MsgB, error)
    grpc.ServerStream
}
```

服务器端处理程序可以重复调用此参数的 Recv 以读取客户端到服务器的消息流，一旦到达客户端到服务器流的末尾，`Recv`将返回 `(nil, io.EOF)`。通过重复调用此`<ServiceName>_FooServer`参数上的`Send`方法来发送响应服务器到客户端的消息流。服务器到客户端的流的结束由程序结束指定(`return`)。



## 生成的客户端接口上的方法

对于客户端的使用，proto 文件中的每个服务`Bar`也会产生一个函数：`func BarClient(cc *grpc.ClientConn) BarClient`，它会返回`BarClient`接口的具体实现（这个具体实现也存在于生成的`.pb.go`文件中）。


**一元方法**

如下签名：

`(ctx context.Context, in *MsgA, opts ...grpc.CallOption) (*MsgB, error)`


**服务端流方法**

```go
// 签名
Foo(ctx context.Context, in *MsgA, opts ...grpc.CallOption) (<ServiceName>_FooClient, error)


type <ServiceName>_FooClient interface {
	Recv() (*MsgB, error)
	grpc.ClientStream
}
```

**客户端流**

```go
Foo(ctx context.Context, opts ...grpc.CallOption) (<ServiceName>_FooClient, error)


type <ServiceName>_FooClient interface {
	Send(*MsgA) error
	CloseAndRecv() (*MsgB, error)
	grpc.ClientStream
}
```


**双向流方法**

```go
Foo(ctx context.Context, opts ...grpc.CallOption) (<ServiceName>_FooClient, error)


type <ServiceName>_FooClient interface {
	Send(*MsgA) error
	Recv() (*MsgB, error)
	grpc.ClientStream
}
```

*客户端到服务端的流结束是调用`CloseSend`*

