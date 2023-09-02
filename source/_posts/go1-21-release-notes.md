---
title: go1.21-release-notes
date: 2023-08-22 20:20:38
tags:
- Golang
categories:
- Intro
keywords:
copyright: Guader
copyright_author_href:
copyright_info:
---

# GO1.21 Release Notes

> 前两天更新Clash，编译时候提示我要更新到1.21，正好详细看看有哪些新东西  
> [原文](https://go.dev/doc/go1.21)


## 介绍

大部分变化体现在工具链`toolchain`，运行时(`runtime`)，以及库(`libraries`)的实现。还有版本号的一点小变化，从`N.N`到`N.N.N`，详情查看[Go versions](https://go.dev/doc/toolchain#version)。


## 语言的变化

增加了三个新的内键函数。

- `min, max`: 计算给定数的最大最小值。[详情](https://go.dev/ref/spec#Min_and_max)
- `clear`: 删除一个`map`的所有对象或是将`slice`的所有对象零值化。[详情](https://go.dev/ref/spec#Clear)


现在更精确的指定了包初始化的顺序。新的算法是：

- 根据导入路径排序所有的包
- 重复以下操作，直到包列表为空：
  - 找到列表中所有导入已经初始化完成的第一个包
  - 初始化该包并从列表中移除


这可能会改变某些程序的行为，因为它们依赖于特定的初始化顺序，而这种顺序并没有通过显式导入来表达。在过去的版本中，规范对此类程序的行为并没有明确定义。新规则提供了明确的定义。

已经进行了多项改进，提高了类型推断的能力和精度。

- 现在可以调用一个（可能部分实例化的）泛型函数，其参数本身也是（可能部分实例化的）泛型函数。编译器将尝试推断被调用者缺少的类型参数（如前所述），并为每个未完全实例化的泛型函数参数推断其缺少的类型参数（新参数）。典型的使用场景是调用在容器上运行的泛型函数（例如`slices.IndexFunc`），在这种情况下，函数参数也可能是泛型的，被调用函数的类型参数及其参数是从容器类型中推断出来的。更常见的情况下，如果可以从赋值中推断出参数的类型，则当将泛型函数分配给变量或作为结果值返回时，就可以使用该函数，而无需显式实例化。
- 类型推断现在也会在为接口赋值时考虑方法：方法签名中使用的类型参数可以从匹配方法的相应参数类型中推断出来。
- 同样，由于类型参数必须实现其相应约束条件的所有方法，因此类型参数和约束条件的方法必须匹配，这可能导致推断出更多的类型参数。
- 如果将多个不同类型的无类型常量参数（例如一个无类型的int和一个无类型的float）传递给具有相同（未指定）类型参数类型的参数，现在类型推断将使用与具有无类型常量操作数的运算符相同的方法来确定类型，而不是出错。这个变化使从无类型常量参数推断出的类型与常量表达式的类型保持一致。
- 在赋值中匹配相应类型时，类型推断现在变得精确了：组件类型（如slice的元素，或是函数签名中的参数类型）必须相同（给定合适的类型参数）才能匹配，否则推理失败。这一变化使错误信息更加准确：在过去，类型推断可能会错误的成功并导致无效赋值，而现在，如果两个类型不可能匹配，编译器就会报告推断错误。

更广泛的说，语言规范中对[类型推论](https://go.dev/ref/spec#Type_inference)的描述已得到澄清。所有这些变化使得类型推断功能更加强大，推断失败也不再那么令人惊讶。

GO1.21包含我们考虑再GO未来版本中进行的语言更改预览：将`for`循环变量设置为每次迭代，而不是每个循环，以避免意外共享bug。详细信息，参阅[LoopvarExperiment wiki页](https://go.dev/wiki/LoopvarExperiment)。

GO1.21现在规定，如果一个`goroutine panic`了，而`recover`被一个`defer`函数直接调用，则保证`recover`的返回值不为`nil`。为确保这一点，使用接口值为`nil`（或未键入的`nil`）调用`panic`会导致[`*runtime.PanicNilError`](https://go.dev/pkg/runtime/#PanicNilError)类型的运行时`panic`。

为支持为旧版GO编写的程序，可通过设置`GODEBUG=panicnil=1`重新启用`nil panic`。如果编译程序的主软件包位与声明GO1.20或更早版本的模块中，则会自动启用此设置。


## 工具

GO1.21在GO工具链中改进了对向后兼容性和向前兼容性的支持。

为了提高向后兼容性，GO1.21正式规定GO使用`GODEBUG`环境变量来控制那些根据兼容性策略是非破坏性的，但可能导致现有程序破坏的更改的默认行为。（例如，依赖于错误行为的程序在错误修复后可能会被破坏，但错误修复不被视为破坏性更改。）当GO必须改变这种行为时，它会根据工作区`go.work`文件中的`go`行或主模块的`go.mod`文件在新旧行为之间做出选择。升级到新的GO工具链，但将`go`行设置为原来的（较旧的）GO版本，可以保留较旧工具链的行为。有了这种兼容性支持，最新的GO工具链应该总是旧版本GO的最佳，最安全的实现。查看[Go，向后兼容性和GODBUG](https://go.dev/doc/godebug)了解更多。

为了提高向前兼容性，GO1.21现在将读取`go.work`或`go.mod`文件中的`go`行作为严格的最低要求：`go 1.21.0`意味着该工作区或模块不能在GO 1.20或 GO 1.21rc1中使用。这样，依赖于GO后期版本修正的项目就可以确保它们不会被早期版本使用。它还能为使用GO新功能的项目提供更好的错误报告：当问题是需要更新的GO版本时，该问题会被清楚的报告，而不是在尝试构建代码时打印出未解决的导入或语法错误。

为了使这些新的更严格的版本要求更易于管理，`go`命令现在不仅可以调用捆绑在其自身版本中的工具链，还可以调用在`PATH`中找到的或按需下载的其他GO工具链版本。如果`go.mod`或`go.work`的`go`行声明了对GO更新版本的最低要求，`go`命令将自动查找并运行该版本。新工具链指令设置了建议使用的最小工具链，它可能比严格的`go`最小工具链更新。详情查看[Go Toolchains](https://go.dev/doc/toolchain)。



## Go 命令

现在，`-pgo`构建标志默认为`-pgo=auto`，在命令行中指定单个主软件包的限制也已取消。如果主软件包文件夹中存在名为`default.pgo`的文件，`go`命令将使用该文件启用配置文件引导的优化，以构建相应的程序。

现在，使用`-C dir`标志时，必须将其作为命令行的第一个标志。

新的`go test`选项`-fullpath`会在测试日志信息中打印完整的路径名，而不仅仅是基名。

`go test -c`标志现在支持为多个软件包写入测试二进制文件，每个二进制文件都写入`pkg.test`，其中`pkg`是软件包名称。如果正在编译的测试软件包中有一个以上的软件包名称，则会出错。

`go test -o`现在可以接收文件夹参数，在这种情况下，测试二进制文件将写入该文件夹，而不是当前文件夹。


## CGO

在`import "C"`的文件中，如果尝试在C类型上声明GO方法，GO工具链现在会正确报错。


## Runtime

在打印非常深的堆栈时，运行时现在会打印前50帧（最内层），然后打印最底层50帧（最外层），而不是只打印前100帧。这使得查看递归堆栈开始的深度变得更容易，对于调试堆栈溢出尤其有价值。

在支持透明大页(`transparent huge pages`)的Linux平台上，GO运行时现在可以更明确的管理堆的哪些部分可以由大页支持。这样就能更好的利用内存：小型堆应减少内存使用量（在病理情况下最多可减少50%），而大型堆应减少堆中密集部分被破坏的大页，CPU使用率和延迟最多可改善1%。

通过运行时内部垃圾收集调整，应用程序的尾部延迟最多可减少40%，内存使用量也会略有减少。某些应用可能还会观察到吞吐量的少量损失。内存使用量的减少应与吞吐量的损失成正比，因此，只要稍微增加`GOGC`[和/或]`GOMEMLIMIT`，就可以恢复上一版本的吞吐量/内存权衡（延迟变化不大）。

在C创建的线程上从C调用GO需要一些设置来准备GO的执行。在UNIX平台，这种设置现在可以在同一线程的多次调用中得到保留。这大大减少了后续C到GO调用的开销，从每次调用约1～3微秒减少到每次调用约100-200纳秒。


## Compiler

策略引导优化（PGO），在GO 1.20中作为预览版加入，现在已可普遍使用。PGO可对生产工作负载配置文件确认为常用的代码进行额外优化。如前文GO命令部分所述，对于主软件包文件夹中包含`default.pgo`配置文件的二进制文件，默认情况下会启用PGO。性能的提升因应用程序行为而异，在一组具有代表性的围棋程序中，大多数程序都能通过启用PGO实现2%到7%的提升。查看[PGO user guide](https://go.dev/doc/pgo)了解更多。

PGO构建现在可以将某些接口方法调用虚拟化，为最常见的被调用者添加具体调用。这样就可以进一步优化，例如内联被调用者。

GO1.21将编译速度提高了6%，这主要归功于使用PGO构建编译器本身。



## Assembler

> 这部分和Linker部分，我不太了解，翻译的可能有问题。

在`amd64`架构，无框架`nosplit`装配函数不再自动标记为`NOFRAME`。作为替代，如果需要，必须明确指定`NOFRAME`属性，这在其他支持帧指针的架构上已经是一种行为。这样，运行时就能为堆栈转换维护帧指针。

改进了在`amd64`上进行动态链接时检查R15使用错误的验证器。


## Linker

在windows/amd64上，链接器（在编译器的帮助下）现在默认会输出SEH开卷数据，这改进了GO应用程序与Windows调试器和其他工具的集成。

在GO 1.21中，如果变量初始化器中的条目数量足够多，并且初始化器表达式没有副作用，链接器（在编译器的帮助下）现在能够删除死的（未引用的）全局映射变量。



## 核心库


### 新的`log/slog`包

新的[`log/slog`](https://go.dev/pkg/log/slog)包提供带级别的结构化日志。结构化日志记录会发出键值对,以实现快速、准确的处理大量日志数据。该包支持与流行的日志分析工具和服务集成。


### 新的`testing/slogtest`包

新的[`testing/slogtest`](https://go.dev/pkg/testing/slogtest)可以帮助验证[`slog.Handler`](https://go.dev/pkg/log/slog#Handler)的实现。


### 新的`slices`包

新的[`slices`](https://go.dev/pkg/slices)包使用可处理任何元素类型的泛型函数提供了很多切片的常见操作。


### 新的`maps`包

新的[`maps`](https://go.dev/pkg/maps/)包使用可处理任何元素类型的泛型函数提供了很多map的常用操作。


### 新的`cmp`包

新的[`cmp`](https://go.dev/pkg/cmp/)定义了类型约束[Ordered](https://go.dev/pkg/cmp/#Ordered)和两个新的泛型函数[Less](https://go.dev/pkg/cmp/#Less)和[Compare](https://go.dev/pkg/cmp/#Compare)，它们对于[有序类型](https://go.dev/ref/spec/#Comparison_operators)很有用。


### 库的小改动

- [`archive/tar`](https://go.dev/pkg/archive/tar/)

> 由[`Header.FileInfo`](https://go.dev/pkg/archive/tar/#Header.FileInfo)返回的[`io/fs.FileInfo`](https://go.dev/pkg/io/fs/#FileInfo)接口的实现现在实现了一个调用[`io/fs.FormatFileInfo`](https://go.dev/pkg/io/fs/#FormatFileInfo)的`String`方法。

- [`archive/zip`](https://go.dev/pkg/archive/zip/)

> 由[`FileHeader.FileInfo`](https://go.dev/pkg/archive/zip/#FileHeader.FileInfo)返回的[`io/fs.FileInfo`](https://go.dev/pkg/io/fs/#FileInfo)接口的实现现在实现了一个调用[`io/fs.FormatFileInfo`](https://go.dev/pkg/io/fs/#FormatFileInfo)的`String`方法。
> 由[`Reader.Open`](https://go.dev/pkg/archive/zip/#Reader.Open)返回的[`io/fs.File`](https://go.dev/pkg/io/fs/#File)的[`io/fs.ReadDirFile.ReadDir`](https://go.dev/pkg/io/fs/#ReadDirFile.ReadDir)方法返回的[`io/fs.DirEntry`](https://go.dev/pkg/io/fs/#DirEntry)接口的实现现在实现了一个调用[`io/fs.FormatDirEntry`](https://go.dev/pkg/io/fs/#FormatDirEntry)的`String`方法。

- [`bytes`](https://go.dev/pkg/bytes/)

> [`Buffer`](https://go.dev/pkg/bytes/#Buffer)类型有两个新方法：[`Available`](https://go.dev/pkg/bytes/#Buffer.Available)和[`AvailableBuffer`](https://go.dev/pkg/bytes/#Buffer.AvailableBuffer)。这些方法可与[`Write`](https://go.dev/pkg/bytes/#Buffer.Write)方法一起使用，直接向`Buffer`追加内容。

- [`context`](https://go.dev/pkg/context/)

> 新的[`WithoutCancel`](https://go.dev/pkg/context/#WithoutCancel)函数返回一个上下文副本，该副本在原始上下文被取消时不会被取消。
> 新的[`WithDeadlineCause`](https://go.dev/pkg/context/#WithDeadlineCause)和[`WithTimeoutCause`](https://go.dev/pkg/context/#WithTimeoutCause)函数提供了一种在最后期限或计时器到期时设置上下文取消原因的方法。可以使用[`Cause`](https://go.dev/pkg/context/#Cause)函数查看原因。
> 新的[`AfterFunc`](https://go.dev/pkg/context/#AfterFunc)函数用于注册一个函数，以便在取消上下文后运行。
> 优化意味着调用[Background](https://go.dev/pkg/context/#Background)和[TODO](https://go.dev/pkg/context/#TODO)并将其转换为共享类型的结果可视为相同。在以前的版本中，它们总是不同的。比较[上下文](https://go.dev/pkg/context/#Context)值是否相等从来没有明确的定义，因此这不被认为是不兼容的更改。

- [`crypto/ecdsa`](https://go.dev/pkg/crypto/ecdsa/)

> [`PublicKey.Equal`](https://go.dev/pkg/crypto/ecdsa/#PublicKey.Equal)和[`PrivateKey.Equal`](https://go.dev/pkg/crypto/ecdsa/#PrivateKey.Equal)现在可在恒定时间内执行。

- [`crypto/elliptic`](https://go.dev/pkg/crypto/elliptic/)

> 所有的[`Curve`](https://go.dev/pkg/crypto/elliptic/#Curve)方法以及[`GenerateKey`](https://go.dev/pkg/crypto/elliptic/#GenerateKey)，[`Marshal`](https://go.dev/pkg/crypto/elliptic/#Marshal)，[`Unmarshal`](https://go.dev/pkg/crypto/elliptic/#Unmarshal)都已废弃。对于`ECDH`操作，应使用新的[`crypto/ecdh`](https://go.dev/pkg/crypto/ecdh/)包。对于较底层的操作，可使用三方模块，如[`filippo.io/nistec`](https://pkg.go.dev/filippo.io/nistec)。

- [`crypto/rand`](https://go.dev/pkg/crypto/rand/)

> 该包现在可用在`NetBSD 10.0`及更高版本上使用`getrandom`系统调用。

- [`crypto/rsa`](https://go.dev/pkg/crypto/rsa/)

> 对于`GOARCH=amd64`和`GOARCH=arm64`，私有RSA操作（解密和签名）的性能比GO 1.19更好。在GO 1.20中，性能有所下降。
> 由于在[`PrecomputedValues`](https://go.dev/pkg/crypto/rsa/#PrecomputedValues)中添加了私有字段，因此即使反序列化（例如从JSON）先前已预计算的私钥，也必须调用[`PrivateKey.Precompute`](https://go.dev/pkg/crypto/rsa/#PrivateKey.Precompute)以获得最佳性能。
> [`PublicKey.Equal`](https://go.dev/pkg/crypto/rsa/#PublicKey.Equal)和[`PrivateKey.Equal`](https://go.dev/pkg/crypto/rsa/#PrivateKey.Equal)现在可在恒定时间内执行。
> [`GenerateMultiPrimeKey`](https://go.dev/pkg/crypto/rsa/#GenerateMultiPrimeKey)函数和[`PrecomputedValues.CRTValues`](https://go.dev/pkg/crypto/rsa/#PrecomputedValues.CRTValues)字段已废弃。在调用[`PrivateKey.Precompute`](https://go.dev/pkg/crypto/rsa/#PrivateKey.Precompute)时，[`PrecomputedValues.CRTValues`](https://go.dev/pkg/crypto/rsa/#PrecomputedValues.CRTValues)仍将被填充，但在解密操作中不会使用这些值。

- [`crypto/sha256`](https://go.dev/pkg/crypto/sha256/)

> 当`GOARCH=amd64`时，SHA-224和SHA-256操作现在使用本地指令，性能提高了3-4倍。

- [`ctypto/tls`](https://go.dev/pkg/crypto/tls/)

> 现在，服务器除了检查过期时间外，还会跳过验证恢复连接的客户端证书（包括不运行[`Config.VerifyPeerCertificate`](https://go.dev/pkg/crypto/tls/#Config.VerifyPeerCertificate)）。这样，在使用客户端证书时，会话票据(session tickets)就会变大。客户端已在恢复时跳过验证，但现在即使设置了[`Config.InsecureSkipVerify`](https://go.dev/pkg/crypto/tls/#Config.InsecureSkipVerify)也会检查过期时间。

> 应用现在可以控制会话票据的内容。
> - 新的[`SessionState`](https://go.dev/pkg/crypto/tls/#SessionState)类型描述了可恢复的会话
> - [`SessionState.Bytes`](https://go.dev/pkg/crypto/tls/#SessionState.Bytes)方法和[`ParseSessionState`](https://go.dev/pkg/crypto/tls/#ParseSessionState)函数对`SessionState`进行序列化和反序列化
> - [`Config.WrapSession`](https://go.dev/pkg/crypto/tls/#Config.WrapSession)和[`Config.UnwrapSession`](https://go.dev/pkg/crypto/tls/#Config.UnwrapSession)钩子可在服务器端将`SessionState`与票据进行转换
> - [`Config.EncryptTicket`](https://go.dev/pkg/crypto/tls/#Config.EncryptTicket)和[`Config.DecryptTicket`](https://go.dev/pkg/crypto/tls/#Config.DecryptTicket)方法提供了`WrapSession`和`UnwrapSession`的默认实现
> - `ClientSessionCache`实现可使用[`ClientSessionState.ResumptionState`](https://go.dev/pkg/crypto/tls/#ClientSessionState.ResumptionState)方法和[`NewResumptionState`](https://go.dev/pkg/crypto/tls/#NewResumptionState)函数在客户端存储和恢复会话


为了减少会话票据被用作跨连接跟踪机制的可能性，服务器现在会在每次恢复时（如果支持且未禁用）签发新票据，而且票据不再带有加密密钥的标识符。如果向[`Conn.SetSessionTicketKeys`](https://go.dev/pkg/crypto/tls/#Conn.SetSessionTicketKeys)传递大量密钥，可能会导致明显的性能损失。

客户端和服务器现在都实现了扩展主密码扩展（RFC 7627）(Extended Master Secret extension)。已恢复对[`ConnectionState.TLSUnique`](https://go.dev/pkg/crypto/tls/#ConnectionState.TLSUnique)的弃用，现在将其设置为支持扩展主密钥的恢复连接。

新的[`QUICConn`](https://go.dev/pkg/crypto/tls/#QUICConn)类型支持QUIC实现，包括`0-RTT`支持。请注意，这本身并不是QUIC实现，而且TLS仍不支持`0-RTT`。

新的[`VersionName`](https://go.dev/pkg/crypto/tls/#VersionName)函数返回TLS版本号的名称。

改进了服务器发送的客户端身份验证失败的TLS警告代码。之前，这种故障总是导致“坏证书(bad certificate)”警告。现在，根据RFC 5246和RFC 8446的规定，某些故障会导致更恰当的警告代码：

- 对于TSL 1.3连接，如果服务器配置为使用[`RequireAnyClientCert`](https://go.dev/pkg/crypto/tls/#RequireAnyClientCert)或[`RequireAndVerifyClientCert`](https://go.dev/pkg/crypto/tls/#RequireAndVerifyClientCert)要求客户端验证，而客户端未提供任何证书，服务器现在会返回“需要证书(certificate required)”警告。
- 如果客户端提供的证书不是由服务器上配置的受信任证书颁发机构签署的，服务器将返回“未知证书颁发机构(unknown certificate authority)”警告
- 如果客户端提供的证书过期或无效，服务器将返回“证书过期(expired certificate)”警告
- 在与客户端身份验证失败有关的所有其他情况下，服务器仍会返回“坏证书(bad certificate)”警告

- [`crypto/x509`](https://go.dev/pkg/crypto/x509/)

> [`RevocationList.RevokedCertificates`](https://go.dev/pkg/crypto/x509/#RevocationList.RevokedCertificates)已被启用，取而代之的是新的[`RevokedCertificateEntries`](https://go.dev/pkg/crypto/x509/#RevocationList.RevokedCertificateEntries)字段，它是[`RevocationListEntry`](https://go.dev/pkg/crypto/x509/#RevocationListEntry)的一个片段。[`RevocationListEntry`](https://go.dev/pkg/crypto/x509/#RevocationListEntry)包含[`pkix.RevokedCertificate`](https://go.dev/pkg/crypto/x509/pkix#RevokedCertificate)中的所有字段以及撤销原因代码
> 名称限制现在可在非叶证书上正确执行，而不是在表达名称限制的证书上执行

- [`debug/elf`](https://go.dev/pkg/debug/elf/)

> 新的[`File.DynValue`](https://go.dev/pkg/debug/elf/#File.DynValue)方法可用于检索指定动态标记列出的数值
> `DT_FLAGS_1`动态标记中允许使用的常量标记现在用[`DynFlag1`](https://go.dev/pkg/debug/elf/#DynFlag1)类型定义。这些标记的名称以`DF_1`开头
> 该包现在定义了常量[`COMPRESS_ZSTD`](https://go.dev/pkg/debug/elf/#COMPRESS_ZSTD)
> 该包现在定义了常量[`R_PPC64_REL24_P9NOTOC`](https://go.dev/pkg/debug/elf/#R_PPC64_REL24_P9NOTOC)

- [`debug/pe`](https://go.dev/pkg/debug/pe/)

> 尝试使用[`Section.Data`](https://go.dev/pkg/debug/pe/#Section.Data)或[`Section.Open`](https://go.dev/pkg/debug/pe/#Section.Open)返回的读取器读取包含未初始化数据的部分时，现在会返回错误信息

- [`embed`](https://go.dev/pkg/embed/)

> [`FS.Open`](https://go.dev/pkg/embed/#FS.Open)返回的[`io/fs.File`](https://go.dev/pkg/io/fs/#File)现在有了一个实现[`io.ReaderAt`](https://go.dev/pkg/io/#ReaderAt)的`ReadAt`方法
> 调用[`FS.Open`](https://go.dev/pkg/embed/FS.Open)[`.Stat`](https://go.dev/pkg/io/fs/#File.Stat)将返回一个现在实现了`String`方法的类型，该方法调用[`io/fs.FormatFileInfo](https://go.dev/pkg/io/fs/#FormatFileInfo)

- [`encoding/binary`](https://go.dev/pkg/encoding/binary/)

> 新的[`NativeEndian`](https://go.dev/pkg/encoding/binary/#NativeEndian)变量可用于使用当前机器的本地子节序进行字节片和整数之间的转换

- [`errors`](https://go.dev/pkg/errors/)

> 新的[`ErrUnsupported`](https://go.dev/pkg/errors/#ErrUnsupported)错误提供了一种标准化的方式，用于表示所请求的操作因不支持而无法执行。例如，在使用不支持硬链接的文件系统时调用[`os.Link`](https://go.dev/pkg/os/#Link)

- [`flag`](https://go.dev/pkg/flag/)

> 新的[`BoolFunc`](https://go.dev/pkg/flag/#BoolFunc)函数和[`FlagSet.BoolFunc`](https://go.dev/pkg/flag/#FlagSet.BoolFunc)方法定义了一个不需要参数的标志，并在使用标志时调用一个函数。这与[`Func`](https://go.dev/pkg/flag/#Func)类似，但针对的是布尔标志
> 如果已经在同名标志上调用了[`Set`](https://go.dev/pkg/flag/#Set)，则标志定义（通过[`Bool`](https://go.dev/pkg/flag/#Bool)， [`BoolVar`](https://go.dev/pkg/flag/#BoolVar)，[`Int`](https://go.dev/pkg/flag/#Int)，[`IntVar`](https://go.dev/pkg/flag/#IntVar)等）会出现问题。这一修改旨在检测[初始化顺序的变化](https://go.dev/doc/go1.21#language)导致标记操作发生的顺序与预期不同的情况。在许多情况下，解决这个问题的方法是引入一个显式的包依赖关系，以便在进行任何[`Set`](https://go.dev/pkg/flag/#Set)操作之前正确的对定义进行排序

- [`go/ast`](https://go.dev/pkg/go/ast/)

> 新的[`IsGenerated`](https://go.dev/pkg/go/ast/#IsGenerated)谓词可报告文件语法树是否包含[特殊注释](https://go.dev/s/generatedcode)，这种注释通常表示文件是由工具生成的。
> 新的[`File.GoVersion`](https://go.dev/pkg/go/ast/#File.GoVersion)字段记录了任何`//go:build`或`//+build`指令所需的最小GO版本。

- [`go/build`](https://go.dev/pkg/go/build/)

> 包现在可以解析文件头（在包声明前）中的构建指令（以`//go:`开头的注释）。这些指令可在新的包字段[`Directives`](https://go.dev/pkg/go/build#Package.Directives)、[`TestDirectives`](https://go.dev/pkg/go/build#Package.TestDirectives)和[`XTestDirectives`](https://go.dev/pkg/go/build#Package.XTestDirectives)中使用。

- [`go/build/constraint`](https://go.dev/pkg/go/build/constraint/)

> 新的[`GoVersion`](https://go.dev/pkg/go/build/constraint/#GoVersion)函数返回构建表达式所隐含的最小GO版本

- [`go/token`](https://go.dev/pkg/go/token/)

> 新的[`File.Lines`](https://go.dev/pkg/go/token/#File.Lines)方法以`File.SetLines`接受的相同形式返回文件的行号表

- [`go/types`](https://go.dev/pkg/go/types/)

> 新的[`Package.GoVersion`](https://go.dev/pkg/go/types/#Package.GoVersion)方法返回用于检查软件包的GO语言版本

- [`hash/maphash`](https://go.dev/pkg/hash/maphash/)

> `hash/maphash`包现在有了纯GO实现，可使用`purego`构建标记进行选择

- [`html/template`](https://go.dev/pkg/html/template/)

> 当一个操作出现在JavsScript模板文字中时，将返回新错误[`ErrJSTemplate`](https://go.dev/pkg/html/template/#ErrJSTemplate)。此前返回的是未导出错误(unexported error)。

- [`io/fs`](https://go.dev/pkg/io/fs/)

> 新的[`FormatFileInfo`](https://go.dev/pkg/io/fs/#FormatFileInfo)函数返回[`FileInfo`](https://go.dev/pkg/io/fs/#FileInfo)的格式化版本。新的[`FormatDirEntry`](https://go.dev/pkg/io/fs/#FormatDirEntry)函数返回[`DirEntry`](https://go.dev/pkg/io/fs/#FileInfo)的格式化版本。[`ReadDir`](https://go.dev/pkg/io/fs/#ReadDir)返回的[`DirEntry`](https://go.dev/pkg/io/fs/#DirEntry)实现了一个调用[`FormatDirEntry`](https://go.dev/pkg/io/fs/#FormatDirEntry)的`String`方法，传递给[`WalkDirFunc`](https://go.dev/pkg/io/fs/#WalkDirFunc)的[`DirEntry`](https://go.dev/pkg/io/fs/#DirEntry)值也是如此。

- [`math/big`](https://go.dev/pkg/math/big/)

> 新的[`Int.Float64`](https://go.dev/pkg/math/big/#Int.Float64)方法会返回与多精度整数最接近的浮点数值，并显示四舍五入的结果。

- [`net`](https://go.dev/pkg/net/)

> 在Linux上，当内核支持多路径TCP时，[`net`](https://go.dev/pkg/net/)包现在可以使用多路径TCP。默认情况下不会使用。要在客户端可用时使用多路径TCP，请在调用[`Dialer.Dial`](https://go.dev/pkg/net/#Dialer.Dial)或[`Dialer.DialContext`](https://go.dev/pkg/net/#Dialer.DialContext)方法前调用[`Dialer.SetMultipathTCP`](https://go.dev/pkg/net/#Dialer.SetMultipathTCP)方法。要在服务器上可用时使用多路径TCP，请在调用[`ListenConfig.Listen`](https://go.dev/pkg/net/#ListenConfig.Listen)方法前调用[`ListenConfig.SetMultipathTCP`](https://go.dev/pkg/net/#ListenConfig.SetMultipathTCP)方法。像往常一样将网络指定为"tcp", "tcp4"或"tcp6"。如果内核或远程主机不支持多路径TCP，连接将无声的退回到TCP。要测试特定连接是否使用多路径TCP，请使用[`TCPConn.MultipathTCP](https://go.dev/pkg/net/#TCPConn.MultipathTCP)方法。
> 在未来的GO版本中，我们可能会在支持多路径TCP的系统上默认启用多路径TCP。

- [`net/http`](https://go.dev/pkg/net/http/)

> 新的[`ResponseController.EnableFullDuplex`](https://go.dev/pkg/net/http#ResponseController.EnableFullDuplex)方法允许服务器处理程序在写入响应的同时读取`HTTP/1`请求正文。通常情况下，`HTTP/1`服务器在开始写入响应前，会自动消耗掉所有剩余的请求体，以避免客户端在读取响应前试图写入一个完成的请求而造成死锁。`EnableFullDuplex`方法会禁止这种行为。
> 当服务器以`HTTP`响应回应`HTTPS`请求时，[`Client`](https://go.dev/pkg/net/http/#Client)和[`Transport`](https://go.dev/pkg/net/http/#Transport)会返回新的[`ErrSchemeMismatch`](https://go.dev/pkg/net/http/#ErrSchemeMismatch)错误。
> [`net/http`](https://go.dev/pkg/net/http/)包现在支持[`errors.ErrUnsupported`](https://go.dev/pkg/errors/#ErrUnsupported)，表达式`errors.Is(http.ErrNotSupported, errors.ErrUnsupported)`将返回`true`。

- [`os`](https://go.dev/pkg/os/)

> 程序现在可以向[`Chtimes`](https://go.dev/pkg/os/#Chtimes)函数传递一个空的`time.Time`值，以保持访问时间或修改时间不变。
> 在Windows中，[`File.Chdir`](https://go.dev/pkg/os#File.Chdir)方法现在可将当前文件夹更改为文件，而不是总是返回错误信息。
> 在Unix系统上，如果向[`NewFile`](https://go.dev/pkg/os/#NewFile)传递了一个非阻塞描述符，调用[`File.Fd`](https://go.dev/pkg/os/#File.Fd)方法现在将返回一个非阻塞描述符。在此之前，描述符会转换为阻塞模式。
> 在Windows系统中，在不存在的文件上调用[`Truncate`](https://go.dev/pkg/os/#Truncate)会创建一个空文件。现在它会返回一个错误，表明文件不存在。
> 在Windows中调用[`TempDir`](https://go.dev/pkg/os/#TempDir)时，现在使用`GetTempPath2W`，而不是`GetTempPathW`。新行为是一种安全加固措施，可防止以`SYSTEM`身份运行的进程创建的临时文件被非`SYSTEM`进程访问。
> 在Windows系统中，`os`包现在支持处理无法以有效UTF-8表示而以UTF-16保存的文件名表示的文件。
> 在Windows中，[`Lstat`](https://go.dev/pkg/os/#Lstat)现在可为以路径分隔符结尾的路径解析符号链接，这与POSIX平台上的行为一致。
> [`ReadDir`](https://go.dev/pkg/os/#ReadDir)函数和[`File.ReadDir`](https://go.dev/pkg/os/#File.ReadDir)方法返回的[`io/fs.DirEntry`](https://go.dev/pkg/io/fs/#DirEntry)接口的实现现在实现了一个调用[`io/fs.FormatDirEntry`](https://go.dev/pkg/io/fs/#FormatDirEntry)的`String`方法。
> [`DirFS`](https://go.dev/pkg/os/#DirFS)函数返回的[`io/fs.FS`](https://go.dev/pkg/io/fs/#FS)接口的实现现在实现了[`io/fs.ReadFileFS`](https://go.dev/pkg/io/fs/#ReadFileFS)和[`io/fs.ReadDirFS`](https://go.dev/pkg/io/fs/#ReadDirFS)接口。

- [`path/filepath`](https://go.dev/pkg/path/filepath/)

> 传递给[`WalkDir`](https://go.dev/pkg/path/filepath/#WalkDir)函数参数的[`io/fs.DirEntry`](https://go.dev/pkg/io/fs/#DirEntry)接口的实现现在实现了一个调用[`io/fs.FormatDirEntry`](https://go.dev/pkg/io/fs/#FormatDirEntry)的`String`方法。

- [`reflect`](https://go.dev/pkg/reflect/)

> 在GO 1.21中，[`ValueOf`](https://go.dev/pkg/reflect/#ValueOf)不再强制在堆上分配其参数，而是允许在栈上分配`Value`的内容。对`Value`的大多数操作也允许在堆栈中分配底层值。
> 新的[`Value`](https://go.dev/pkg/reflect/#Value)方法[`Value.Clear`](https://go.dev/pkg/reflect/#Value.Clear)可清除映射的内容或将片段的内容清零。这与语言中新增的`clear`内置方法相对应。
> [`SliceHeader`](https://go.dev/pkg/reflect/#SliceHeader)和[`StringHeader`](https://go.dev/pkg/reflect/#StringHeader)类型现已弃用。新代码中，请首选[`unsafe.Slice`](https://go.dev/pkg/unsafe/#Slice)、[`unsafe.SliceData`](https://go.dev/pkg/unsafe/#SliceData)、[`unsafe.String`](https://go.dev/pkg/unsafe/#String)或[`unsafe.StringData`](https://go.dev/pkg/unsafe/#StringData)。

- [`regexp`](https://go.dev/pkg/regexp/)

> [`Regexp`](https://go.dev/pkg/regexp#Regexp)现在定义了[`MarshalText`](https://go.dev/pkg/regexp#Regexp.MarshalText)和[`UnmarshalText`](https://go.dev/pkg/regexp#Regexp.UnmarshalText)方法。这些方法实现了[`encoding.TextMarshaler`](https://go.dev/pkg/encoding#TextMarshaler)和[`encoding.TextUnmarshaler`](https://go.dev/pkg/encoding#TextUnmarshaler)，将被[`encoding/json`](https://go.dev/pkg/encoding/json)等包使用。

- [`runtime`](https://go.dev/pkg/runtime/)

> GO程序产生的文本堆栈跟踪（例如在崩溃，调用`runtime.Stack`或使用`debug=2`收集`goroutine profile`时产生的堆栈跟踪）现在包含了在堆栈跟踪中创建每个`goroutine`的`goroutine`的`ID`。
> 崩溃的GO应用现在可以通过设置环境变量`GOTRACEBACK=wer`或在崩溃前调用[`debug.SetTraceback("wer")`](https://go.dev/pkg/runtime/debug/#SetTraceback)来选择加入Windows错误报告(`WER`)。除启用`WER`外，运行时的行为与`GOTRACEBACK=crash`相同。在非Windows系统上，`GOTRACEBACK=wer`将被忽略。
> `GODEBUG=cgocheck=2`是对`cgo`指针传递规则的全面检查，不再作为[调试选项](https://go.dev/pkg/runtime#hdr-Environment_Variables)提供。取而代之的是使用`GOEXPERIMENT=cgocheck2`作为实验。这意味着必须在构建时而不是启动时选择该模式。
> `GODEBUG=cgocheck=1`仍然可用（且仍然是默认值）。
> `runtime`包中新增了`Pinner`类型。`Pinner`可以用来“钉住”GO内存，这样非GO代码就可以更自由的使用它。例如，现在允许向C代码传递引用固定GO内存的GO值。此前，[`cgo`指针传递规则](https://pkg.go.dev/cmd/cgo#hdr-Passing_pointers)不允许传递任何此类嵌套引用。参阅[文档](https://go.dev/pkg/runtime#Pinner)了解更多。

- [`runtime/metrics`](https://go.dev/pkg/runtime/metrics/)

> 一些以前内部的GC指标，如实时堆大小，现在也可用了。`GOGC`和`GOMEMLIMIT`现在也可作为指标使用。

- [`runtime/trace`](https://go.dev/pkg/runtime/trace/)

> 现在，在`amd64`和`arm64`上收集痕迹所需的CPU成本大大降低：与上一版本相比最多可提高10倍。
> 跟踪现在包含显式停止事件，可用于GO运行时可能停止的各种原因，而不仅仅是垃圾回收。

- [`sync`](https://go.dev/pkg/sync/)

> 新的[`OnceFunc`](https://go.dev/pkg/sync/#OnceFunc)、[`OnceValue`](https://go.dev/pkg/sync/#OnceValue)和[`OnceValues`](https://go.dev/pkg/sync/#OnceValues)函数捕捉了[`Once`](https://go.dev/pkg/sync/#Once)的一种常用用法，即在首次使用时懒初始化一个值。

- [`syscall`](https://go.dev/pkg/syscall/)

> 在Windows中，[`Fchdir`](https://go.dev/pkg/syscall#Fchdir)函数现在可以将当前文件夹更改为其参数，而不是总是返回错误。
> 在FreeBSD上，[`SysProcAttr`](https://go.dev/pkg/syscall#SysProcAttr)有一个新字段`Jail`，可用于将新创建的进程置于监禁环境中。
> 在Windows系统中，`syscall`包现在支持处理文件名无法以有效UTF-8表示而以UTF-16保存的文件。现在，[`UTF16ToString`](https://go.dev/pkg/syscall#UTF16ToString)和[`UTF16FromString`](https://go.dev/pkg/syscall#UTF16FromString)函数可在UTF-16数据和[WTF-8](https://simonsapin.github.io/wtf-8/)字符串之间进行转换。这是向下兼容的，因为WTF-8是早期版本中使用的UTF-8格式的超集。
> 有几个错误值与新的[`errors.ErrUnsupported`](https://go.dev/pkg/errors/#ErrUnsupported)匹配，因此`errors.Is(err, errors.ErrUnsupported)`返回`true`。

  - `ENOSYS`
  - `ENOTSUP`
  - `EOPNOTSUPP`
  - `EPLAN9` (Plan 9 only)
  - `ERROR_CALL_NOT_IMPLEMENTED` (Windows only)
  - `ERROR_NOT_SUPPORTED` (Windows only)
  - `EWINDOWS` (Windows only)
  
- [`testing`](https://go.dev/pkg/testing/)

> 新的`-test.fullpath`选项将在测试日志信息中打印完整的路径名，而不仅仅是基名。
> 新的[`Testing`](https://go.dev/pkg/testing/#Testing)函数会报告程序是否是由`go test`创建的测试。

- [`testing/fstest`](https://go.dev/pkg/testing/fstest/)

> 调用[`Open`](https://go.dev/pkg/testing/fstest/MapFS.Open)[`.Stat`](https://go.dev/pkg/io/fs/#File.Stat)将返回一个类型，该类型现在实现了调用[`io/fs.FormatFileInfo`](https://go.dev/pkg/io/fs/#FormatFileInfo)的`String`方法。

- [`unicode`](https://go.dev/pkg/unicode/)

> 整个系统的[`unicode`](https://go.dev/pkg/unicode/)包和相关支持已升级到[`Unicode 15.0.0`](https://www.unicode.org/versions/Unicode15.0.0/)。


## Ports

### WebAssembly

现在可以在GO程序中使用新的`go:wasmimport`指令从`WebAssembly`主机导入函数。

现在，GO调度器与JS事件循环的交互效率大大提高，尤其是在频繁阻塞异步事件的应用程序中。


### WebAssembly System Interface

GO 1.21为[WebAssembly系统接口（WASI）](https://wasi.dev/)预览版1增加了一个实验端口（`GOOS=wasip1, GOARCH=wasm`）。

由于增加了新的`GOOS`值"wasip1"，除非使用该`GOOS`值，否则名为`*_wasip1.go`的GO文件将[被GO工具忽略](https://go.dev/pkg/go/build/#hdr-Build_Constraints)。如果现有文件名和该模式匹配，需要重新命名。


### ppc64/ppc64le

在Linux上，`GOPPC64=power10`现在可以生成PC相关指令，前缀指令和其他新的`Power10`指令。在AIX上，`GOPPC64=power10`会生成`Power10`指令，但不会生成PC相关指令。

在为`GOPPC64=power10 GOOS=linux GOARCH=ppc64le`构建与位置无关的二进制文件时，用户可以期待在大多数情况下缩小二进制文件的大小，在某些情况下缩小`3.5%`。与位置无关的二进制文件是为`ppc64le`构建的，其`-buildmode`值如下：`c-archive, c-shared, shared, pie, plugin`。

### loong64

`linux/loong64`端口现在支持`-buildmode=c-archive`, `-buildmode=c-shared`和`-buildmode=pie`。


---

> 上周准备写的，一下给忘了，拖了一周。
