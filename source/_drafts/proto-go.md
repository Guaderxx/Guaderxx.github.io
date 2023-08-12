---
title: proto-go
date: 2023-08-12 15:09:56
tags:
categories:
keywords:
copyright: Guader
copyright_author_href:
copyright_info:
---

# 使用Protocol Buffer的基本介绍

首先，定义proto文件

```proto
// addressbook.proto
syntax = "proto3";
package tutorial;

import "google/protobuf/timestamp.proto";

option go_package = "github.com/protocolbuffers/protobuf/examples/go/tutorialpb";

enum PhoneType {
    PHONE_TYPE_UNSPECIFIED = 0;
    PHONE_TYPE_MOBILE = 1;
    PHONE_TYPE_HOME = 2;
    PHONE_TYPE_WORK = 3;
}

message PhoneNumber {
    string number = 1;
    PhoneType type = 2;
}

message Person {
    string name = 1;
    int32 id = 2;    // unique ID number for this person.
    string email = 3;
    
    repeated PhoneNumber phones = 4;
    
    google.protobuf.Timestamp last_updated = 5;
}

// Our address book file is just one of these.
message AddressBook {
    repeated Person perple = 1;
}
```

在上面的示例中，Person消息包含PhoneNumber消息，而AddressBook包含Person消息。


## 编译你的proto文件

1. 下载protoc编译器
2. 下载GO插件，`go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
3. 运行protoc编译器，`protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/addressbook.proto`



## 编译器调用

调用时如果使用了`go_out`标志，则会产生GO输出。`go_out`的参数就是希望编译器写入GO输出的目录，每个proto文件生成一个GO文件，名称名用`.pb.go`替换掉`.proto`扩展名。

生成的`.pb.go`文件在输出目录的位置取决于编译器参数。有多种输出模式：

- 如果指定了`paths=import`参数，输出文件放置在以GO包的导入路径命名的文件夹中。例如，输入文件`protos/buzz.proto`，GO导入路径为`example.com/project/protos/fizz`，那么输出文件在`example.com/project/protos/fizz/buzz.pb.go`。如果没有指定`paths`参数，这也是默认的输出模式。
- 如果`module=$PREFIX`参数指定了，输出文件放置在以GO包的导入路径命名的目录中，但从输出文件名中删除指定的目录前缀。例如，输入文件`protos/buzz.proto`，GO导入路径为`example.com/project/protos/fizz`，同时将`example.com/project`指定为`module`前缀，则输出文件在`protos/fizz/buzz.pb.go`。在模块路径外生成任何GO包都会导致错误。这个模式适用于将生成的文件直接输出到GO模块中。
- 如果`paths=source_relative`参数指定了，输出文件放置在和输入文件相同的相对文件夹中。例如，输入文件`protos/buzz.proto`，输出文件在`protos/buzz.pb.go`。

特定于`protoc-gen-go`的标志是通过调用`protoc`是传递`go_opt`标志提供的。可以传递多个`go_opt`标志，例如：

`protoc -I=src --go_out=out --go_opt=paths=source_relative foo.proto bar/baz.proto`

编译器会从`src`文件夹读取`foo.proto, bar/baz.proto`，然后输出文件`foo.pb.go, bar/baz.pb.go`到`out`文件夹。编译器会自动创建输出的嵌套子文件夹，但不会自动创建输出文件夹。



## 包

为了生成GO代码，必须为每个proto文件（包括正在生成的proto文件传递依赖的文件）提供GO包的导入路径。有两种方式指定GO导入路径。

- 在`proto`文件内声明，或
- 调用`protoc`时声明

建议在proto文件中声明，方便proto文件的GO包和proto文件对应上，且简化调用protoc时传递的参数。命令行传递的优先级会更高。

通过使用GO包的完整导入路径声明`go_package`选项，在proto文件中指定：

```proto
option go_package = "example.com/project/protos/fizz";
```

也可以在调用编译器时指定一个或多个`M${PROTO_FILE}=${GO_IMPORT_PATH}`：

```bash
protoc -I=src \
    --go_opt=Mprotos/buzz.proto=example.com/project/protos/fizz \
    --go_opt=Mprotos/bar.proto=example.com/project/protos/foo \
    protos/buzz.proto protos/bar.proto
```


GO导入路径和proto文件的包标识符没有关系。后者只与protobuf命名空间相关，而前者只与GO命名空间相关。此外，GO导入路径和proto导入路径之间也没有关联。



## 类型转换

[原文](https://protobuf.dev/reference/go/go-generated/#message)
