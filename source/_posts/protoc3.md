---
title: protoc3
date: 2023-08-12 11:21:56
tags:
categories:
- Intro
- pb3
keywords:
copyright: Guader
copyright_author_href:
copyright_info:
---

# Proto 3

> 一些简介  
> Protocol Buffer是一种与语言无关，与平台无关的可扩展机制，用于序列化结构化数据


## Style Guide

[原文](https://protobuf.dev/programming-guides/style/)



## Enum Behavior

在GO中，enum会被编码为int32。

枚举有两种不同的风格：`open`和`closed`。
除了处理未知值意外，其他行为相同。

```proto
enum Enum {
    A = 0;
    B = 1;
}

message Msg {
    optional Enum enum = 1;
}
```

*`open`和`closed`的区别可以概括为一个问题*  
当程序解析包含enum的Msg时，其值为2，会发生什么？

- `open`  enums 将解析值 2 并将其直接存储在字段中。访问器将报告该字段已设置，并将返回代表 2 的内容
- `closed`  enums 将解析值 2 并将其存储在消息的未知字段集中。访问器将报告该字段未设置，并将返回枚举的默认值。


而所有已知的GO版本都不符合要求。GO将所有枚举视为`open`。



## Well-Known Types and Common Types

Well-Known Types

- [duration](https://github.com/protocolbuffers/protobuf/blob/main/src/google/protobuf/duration.proto) 带符号的固定长度时间跨度（例如：42S）
- [timestamp](https://github.com/protocolbuffers/protobuf/blob/main/src/google/protobuf/timestamp.proto) 独立于任何时区或日历的时间点 （例如：2017-01-15T01:30:15.01Z）
- [field_mask](https://github.com/protocolbuffers/protobuf/blob/main/src/google/protobuf/field_mask.proto) 一组符号字段路径 （例如：`f.b.d`）


Common Types

- [interval](https://github.com/googleapis/googleapis/blob/master/google/type/interval.proto) 独立于时区或日历的时间间隔 （例如：2017-01-15T01:30:15.01Z - 2017-01-16T02:30:15.01Z）
- [date](https://github.com/googleapis/googleapis/blob/master/google/type/date.proto) 整个日期时间 （例如：2005-09-19）
- [dayofweek](https://github.com/googleapis/googleapis/blob/master/google/type/dayofweek.proto) 一周的日子 （例如： Monday）
- [timeofday](https://github.com/googleapis/googleapis/blob/master/google/type/timeofday.proto) 一天的时间 （例如：10:22:23）
- [latlng](https://github.com/googleapis/googleapis/blob/master/google/type/latlng.proto) 是纬度/经度对（例如，37.386051 纬度和 -122.083855 经度
- [money](https://github.com/googleapis/googleapis/blob/master/google/type/money.proto) 具有货币类型的金额 （例如， 42USD）
- [postal_address](https://github.com/googleapis/googleapis/blob/master/google/type/postal_address.proto) 是邮政地址（例如，1600 Amphitheatre Parkway Mountain View, CA 94043 USA）
- [color](https://github.com/googleapis/googleapis/blob/master/google/type/color.proto) RGBA颜色空间中的颜色
- [month](https://github.com/googleapis/googleapis/blob/master/google/type/month.proto) 一年的月份 （例如，April）



## 定义一个消息类型

```proto
syntax = "proto3";    // 指定使用的是proto3协议

message SearchRequest {    // 定义一个有三个字段的消息
    string query = 1;
    int32 page_number = 2;
    int32 results_per_page = 3;
}
```

你必须给消息中的每个字段分配一个在`1`到`536,870,911`之间的数字，同时还有下面这些限制：

- 分配的数字在这个消息体中**必须是唯一的**
- `19,000`到`19,999`之间的数字分配给了Protocol Buffer实现，不可用
- 你不能用任何之前保留的数字或是分配给扩展的数字

**一旦信息类型在使用，该数字就无法更改**，因为它标识[消息传输格式](https://protobuf.dev/programming-guides/encoding)中的字段。“更改”字段编号相当于删除该字段并创建一个具有相同类型但新编号的新字段。查看[删除字段](https://protobuf.dev/programming-guides/proto3/#deleting)了解如何正确执行此操作的信息。

字段数字**永远不应该重复使用**。不要将[保留列表](https://protobuf.dev/programming-guides/proto3/#fieldreserved)中的数字取出供新字段重复使用。查阅[重复使用字段编号的后果](https://protobuf.dev/programming-guides/proto3/#consequences)。

你应该将字段编号1到15用在最常用的字段。较低的字段数字会占用较少的空间。例如，1到15范围内的字段编号需要一个字节进行编码，而16到2047范围的编号需要两个字节。你可以在[Protocol Buffer编码](https://protobuf.dev/programming-guides/encoding#structure)中了解更多。


## 重复使用字段编号的后果

重复使用字段编号会使用解码消息变得不明确。

Protobuf消息格式很精简，并没有提供一种方法来检测使用一种定义编码并使用另一种定义解码的字段。

使用一个定义对字段进行编码，再使用另一个不同的定义对同一字段进行解码会导致：

- 开发人员因调试而损失的时间
- 解析/合并错误（最好的情况）
- PII/SPII泄漏
- 数据损坏


字段号重复使用的常见原因：

- 对字段重新编号（有时是为了使字段编号顺序更美观）。重新编号实际上会删除并重新添加新编号中涉及的所有字段，从而导致不兼容的格式更改。
- 删除字段且不保留编号以防止将来重复使用


最大字段为29位(bit)，而不是更常见的32位，因为有3个较低位用于有线格式。了解更多可查看[编码主题](https://protobuf.dev/programming-guides/encoding#structure).



## 指定字段标签

消息字段可以增加以下的标签：

- `optional`  一个`optional`字段可能是以下两种状态
  - 字段存在，并包含一个显式设置或从线路解析的值。它将被序列化到线路上
  - 字段不存在，会返回默认值。不会被序列化到线路上。
  可以检查该值是否明确设置
- `repeated`  该字段类型可以在格式正确的消息中重复零次或多次。会保留重复值的顺序。
- `map`  这是成对的key:value字段类型。查看[Maps](https://protobuf.dev/programming-guides/encoding#maps)了解更多。
- 如果没有使用显式字段标签，则假定使用默认字段标签，称为“隐式字段存在”。（无法显式将字段设置为该状态）。格式良好的消息可以用零个或一个次字段（但不能超过一个）。你也无法确定是否从线路中解析了该类型的字段。隐式存在字段将被序列化到线路，除非它是默认值。查看[字段存在](https://protobuf.dev/programming-guides/field_presence)了解更多。


在proto3，标量数值类型的`repeated`字段默认使用`packed`编码。可以在[Protocol Buffer编码](https://protobuf.dev/programming-guides/encoding#packed)中了解更多`packed`。



## 添加更多消息类型

在单个`.proto`文件中可以定义多个消息类型。这对于定义多个相关的消息很有用。例如，如果你想定义与`SearchRequest`对应的消息类型`SearchResponse`，可以直接在同文件中添加。

```proto
message SearchRequest {
    ...
}

message SearchResponse {
    ...
}
```


## 添加注释

要在你的`.proto`文件中添加注释，使用`C/C++`格式的`//`或是`/* ... */`语法。

```proto
/* Comment
 * Multiple line comments
 *
 */
 
message SearchRequest {  // single line comment
    ...
}
```



## 删除字段

如果操作不当，删除字段可能会导致严重问题。

当你不再需要一个字段且客户端相关的代码已经删掉了，可以从消息中删除该字段定义。但是，你**必须**[保留已删除的字段编号](https://protobuf.dev/programming-guides/proto3/#fieldreserved)。如果你不保留字段编号，开发人员将来可以重复使用该编号。

你还应该保留字段名称，以允许消息的JSON和TextFormat编码继续解析。



## 保留字段

如果你通过删除一个字段或是注释掉来[更新](https://protobuf.dev/programming-guides/proto3/#updating)一个消息类型，后来的开发可能会在更新类型时重用字段编号。这会导致很多问题，就像在[重复使用字段编号](https://protobuf.dev/programming-guides/proto3/#consequences)中描述的。

为了确保这不会发生，将你删除的字段添加到`reserved`列表。为了确保消息的JSON和TextFormat实例依然可以被解析，将字段名也添加到`reserved`列表。

Protocol buffer编译器会在新开发尝试使用这些保留的字段或是字段名时报错。

```proto
message Foo {
    reserved 2, 15, 9 to 11;
    reserved "foo", "bar";
}
```

保留字段编号可以使用范围（`9 to 11`就是`9, 10, 11`）。
注意：不能在同一个`reserved`语句中混用字段名和字段编号。



## 从proto文件可以生成什么？

当使用[Protocol Buffer编译器](https://protobuf.dev/programming-guides/proto3/#generating)编译proto文件时，编译器会根据你选择的语言生成在文件中描述的消息类型，包括`getting`和`setting`字段值，将消息序列化到输出流，从输入流解析数据。

- 对于GO，编译器生成`.pb.go`文件，包含了proto文件中定义的每个消息类型


有关更多API详细信息，参阅相关[API参考](https://protobuf.dev/reference/)。



## 标量值类型

[原文](https://protobuf.dev/programming-guides/proto3/#scalar)


可以在[Protobuf Buffer编码](https://protobuf.dev/programming-guides/encoding)中查看更多关于在编码中序列化消息。



## 默认值

解析消息时，如果编码的消息不包含特定的单一元素，则解析对象中的相应字段将设置为该字段的默认值。这些默认值是特定于类型的：

- 对于`strings`，默认值就是空字符串
- 对于`bytes`，默认值就是空的字节
- 对于`bools`，默认值就是false
- 对于数值类型，默认值就是0
- 对于`enums`，默认值是**第一个定义的enum值**，也就是0
- 对于消息字段，未设置该字段。明确值取决于语言。查看[生成代码向导](https://protobuf.dev/reference/)了解更多细节。


对于`repeated`字段的空默认值（通常是相关语言的空列表）。


注意：对于标量消息字段，一旦解析消息，就无法判断字段是否已显式设置为默认值（例如布尔值是否设置为false）或根本未设置： 在定义消息类型时应该注意这一点。例如，如果你不希望默认情况下也发生某些行为，则没有一个布尔值可以在设置为false时开启某些行为。另请注意：如果标量消息字段设置为其默认值，则该值将不会在线上序列化。


有关默认值如何在生成的代码中工作的更多详细信息，参阅所用语言的[生成代码指南](https://protobuf.dev/reference/)



## 枚举


```proto
enum Corpus {
    CORPUS_UNSPECIFIED = 0;
    CORPUS_UNIVERSAL = 1;
    CORPUS_WEB = 2;
    CORPUS_IMAGES = 3;
    CORPUS_LOCAL = 4;
    CORPUS_NEWS = 5;
    CORPUS_PRODUCTS = 6;
    CORPUS_VIDEO = 7;
}

message SearchRequest {
    string query = 1;
    int32 page_number = 2;
    int32 results_per_page = 3;
    Corpus corpus = 4;
}
```

Corpus枚举的第一个常量映射到0：每个枚举值定义必须包含一个映射到0的值作为第一个常量，这是因为：

- 必须有0值，这样我们可以作默认值
- 0值需要时第一个值，为了和proto2兼容


可以通过将相同的值分配给不同的枚举常量来定义别名。为此，需要将`allow_alias`选项设置为true，否则编译器会发出警告。尽管所有别名值在反序列化期间都有效，但在序列化时始终使用第一个值。

```proto
enum EnumAlloingAlias {
    option allow_alias = true;
    EAA_UNSPECIFIED = 0;
    EAA_STARTED = 1;
    EAA_RUNNING = 1;
    EAA_FINISHED = 2;
}

enum EnumNotAllowingAlias {
    ENAA_UNSPECIFIED = 0;
    ENAA_STARTED = 1;
    // ENAA_RUNNING = 1; // 不注释这行会导致警告信息
    ENAA_FINISHED = 2;
}
```

枚举常量必须在32位整数范围内。因为枚举值使用的[varint编码](https://protobuf.dev/programming-guides/encoding)，负值效率低下，因此不推荐。



## 保留值

如果通过完全删除枚举条目或将其注释掉来[更新](https://protobuf.dev/programming-guides/proto3/#updating)枚举类型，则将来的用户在对该类型进行自己的更新时可以重用该数值。如果他们稍后加载同一 .proto 的旧版本，这可能会导致严重问题，包括数据损坏、隐私错误等。确保不会发生这种情况的一种方法是指定保留已删除条目的数值（和/或名称，这也可能导致 JSON 序列化问题）。如果任何未来的用户尝试使用这些标识符，Protocol Buffer编译器将会警告。您可以使用 max 关键字指定保留的数值范围达到最大可能值。

```proto
enum Foo {
    reserved 2, 15, 9 to 11, 40 to max;
    reserved "FOO", "BAR";
}
```


## 使用其他消息类型

可以使用其他消息类型作为字段类型。例如，假设你想在每个`SearchResponse`消息中包含`Result`消息，你可以在同一文件中定义`Result`消息，然后在`SearchResponse`中指定`Result`类型的字段。

```proto
message Result {
    string url = 1;
    string title = 2;
    repeated string snippets = 3;
}

message SearchResponse {
    repeated Result results = 1;
}
```



## 导入定义

在上面的例子中，`Result`消息和`SearchResponse`定义在一个文件中 - 如果你想用已经在其他文件中定义的消息作为字段类型呢？

你可以导入它们再使用。语法是这样：

`import "myproject/other_protos.proto";`

默认情况下，只能使用直接导入的`proto`文件中的定义。但是，有时候可能需要将proto文件移动到新位置。可以将占位符proto文件放在旧位置，使用导入公共概念将所有导入转发到新位置，而不是直接移动proto文件并在更改中更新所有调用。

导入包含`import public`语句的原型的任何代码都可以传递依赖`import public`依赖项。例如：

```proto
// new.proto
// All definitions are moved here
```

```proto
// old.proto
// This is the proto that all clients are importing
import public "new.proto";
import "other.proto";
```

```proto
// client.proto
import "old.proto";
// You use definitions from old.proto and new.proto, but not other.proto
```

Protocol编译器使用`-I/--proto_path`标志在命令行指定的一组目录中搜索导入的文件。如果没有指定位置，它会在调用的目录下查找。一般来说，应该将`--proto_path`标志设置为项目的根目录，并为所有导入使用完全限定名称。



## 嵌套类型

你可以在其他消息中定义并使用消息，就像下面的例子 - 在`SearchResponse`中定义`Result`消息：

```proto
message SearchResponse {
    message Result {
        string url = 1;
        string title = 2;
        repeated string snippets = 3;
    }
    repeated Result results = 1;
}
```

如果你想在其他消息中复用这个内部定义的消息，可以这样`_Parent_._Type_`:

```proto
message SomeOtherMessage {
    SearchResponse.Result result = 1;
}
```

你可以按照自己的喜好随意嵌套消息：

```proto
message Outer {                // level 0
    message MiddleAA {         // level 1
        message Inner {        // level 2
            int64 ival = 1;
            bool booly = 2;
        }
    }
    
    message MiddleBB {         // level 1
        message Inner {        // level 2
            int32 ival = 1;
            bool booly = 2;
        }
    }
}
```



## 更新消息类型

如果一个已存在的消息不满足需要了，例如，你希望消息有额外的字段，但你仍然希望使用旧格式创建的代码，不要担心。使用二进制有线格式时，更新消息类型很简单，不会破坏任何现有代码。

> 注意：  
> 如果你使用JSON或proto文本格式来存储protocol buffer消息，则你可以在proto定义中进行的更改会由所不同。


检查[Proto最佳实践](https://protobuf.dev/programming-guides/dos-donts)和以下规则：

- 不要更改任何现有字段的字段编号
- 如果添加新字段，则使用旧消息格式的代码序列化的任何消息仍可以由新生成的代码进行解析。你应该记住这些元素的默认值，以便新代码可以和旧代码生成的消息正确交互。同样，新代码创建的消息旧代码也可解析：旧的二进制文件解析时只是忽略新字段。
- 只要在更新的消息类型中不再使用字段编号，就可以删除字段。你可能想重命名该字段，也许添加前缀`OBSOLETE_`，或者保留字段编号，以便未来的开发不会意外的重复使用该编号。
- `int32, uint32, int64, uint64, bool`都是兼容的。这意味着可以将这些类型的字段更改为另一个类型而不会破坏兼容性。如果从线路中解析的数字不适应该类型，将会获得在C++中将该数字转换为该类型的效果。
- `sint32, sint64`彼此兼容，但与其他整数类型不兼容
- `string, bytes`只要字节是有效的UTF8，就可以兼容
- 如果字节包含消息的编码版本，则嵌入消息和字节兼容
- `fixed32`和`sfixed32, fixed64, sfixed64`兼容
- 对于`string, bytes`，以及消息字段，`optional`和`repeated`兼容。
  给定重复字段的序列化数据作为输入，希望该字段为可选的客户端会采用最后一个输入值（如果它是原始类型字段）或合并所有输入元素（如果是消息类型字段）。
  注意： 这对于数字类型（包括枚举和布尔）通常并不安全。数字类型的重复字段可以`packed`格式序列化，当需要可选字段时，将无法正确解析。
- `enum`和`int32, uint32, int64, uint64`在有线格式方面兼容(注意：如果不合适，值会被截断)。
- 将单个`optional`字段或`extension`更改为新的`oneof`的成员是二进制兼容的，但对于某些语言（尤其是GO），生成的代码的API将以不兼容的方式发生变化。



## 未知字段

未知字段是格式正确的协议缓冲区序列化数据，表示解析器无法识别的字段。例如，当旧二进制文件使用新字段解析新二进制文件发送的数据时，这些新字段将成为旧二进制文件中的未知字段。

最初，proto3 消息在解析过程中总是丢弃未知字段，但在 3.5 版本中，我们重新引入了保留未知字段以匹配 proto2 行为。在 3.5 及更高版本中，未知字段在解析期间被保留并包含在序列化输出中。



## Any

Any消息类型允许你将消息用作嵌入类型，而无需其proto定义。Any包含作为字节的任意序列化消息，以及充当该消息类型的全局唯一标识符并解析为该消息类型的URL。要使用Any类型，需要导入`google/protobuf/any.proto`

```proto
import "google/protobuf/any.proto";

message ErrorStatus {
    string message = 1;
    repeated google.protobuf.Any details = 2;
}
```

不同的语言实现将支持运行时库助手以类型安全的方式打包和解包任何值。



## Oneof

如果您的消息包含多个字段，并且最多同时设置一个字段，则可以使用 oneof 功能强制执行此行为并节省内存。

oneof 字段与常规字段类似，只是所有字段都位于 oneof 共享内存中，并且最多可以同时设置一个字段。设置 oneof 的任何成员都会自动清除所有其他成员。您可以使用特殊的`case()`或`WhichOneof()`方法检查 oneof 中设置了哪个值（如果有），具体取决于您选择的语言。

**Go应该不支持 TODO**

## Maps

如果您想创建关联映射作为数据定义的一部分，protocol buffer提供了一种方便的快捷语法:

`map<key_type, value_type> map_field = N;`

其中`key_type`可以是任何整数类型或字符串类型（因此，除浮点类型和字节之外的任何标量类型）。注意，枚举不行。`value_type`可以是除另一个映射外的任何类型。

例如，如果你想创建一个项目映射，其中每个项目消息都和一个字符串键关联，可以这样定义：

```proto
message Project {
    ...
}

map<string, Project> projects = 3;
```

- Map字段不能`repeated`
- map值的顺序未定义，因此不能依赖特定顺序的映射项
- 为proto生成文本格式时，map按键排序，数字键按数字排序
- 当从线路解析或合并时，如果存在重复的map键，使用最后收到的键。从文本解析map时，如果存在重复的键，解析可能失败
- 如果为map字段提供键但未提供值，则序列化字段时的行为取决于语言。



### 向后兼容性

Map语法相当于以下内容，因此不支持map的protocol buffer实现依然可以处理你的数据：

```proto
message MapFieldEntry {
    key_type key = 1;
    value_type value = 2;
}

repeated MapFieldEntry map_field = N;
```

任何支持map的protocol buffer实现都必须生成并接受上述定义可以接受的数据




## 包

你可以想proto文件添加可选的包说明符，以防止protocol消息类型间发生名称冲突

```proto
package foo.bar;
message Open { ... }
```

你可以在定义消息字段时使用包说明符：

```proto
message Foo {
    ...
    foo.bar.Open open = 1;
    ...
}
```

包说明符影响生成代码的方式取决于使用的语言：

- GO  除非在proto文件中显式指定了`go_package`，否则会被用作包名



## 定义服务

如果要在RPC（Remote Procedure Call）系统中使用消息类型，可以在proto文件中定义一个RPC服务接口，protocol buffer编译器会根据选择的语言生成服务接口代码。例如，如果要定义一个接收`SearchRequest`并返回`SearchResponse`的RPC服务，可以像下面这样定义：

```proto
service SearchService {
    rpc Search(SearchRequest) returns (SearchResponse);
}
```

使用protocol buffer的最直接的RPC系统就是GRPC。

还有很多正在进行的三方项目来开发protocol buffer的RPC实现，查看[三方附加组件wiki页面](https://github.com/protocolbuffers/protobuf/blob/master/docs/third_party.md)



## JSON映射

[原文](https://protobuf.dev/programming-guides/proto3/#json)




## 生成代码

对于 Go，您还需要为编译器安装一个特殊的代码生成器插件：您可以在 GitHub 上的 [golang/protobuf](https://github.com/golang/protobuf/) 存储库中找到此插件和安装说明。



