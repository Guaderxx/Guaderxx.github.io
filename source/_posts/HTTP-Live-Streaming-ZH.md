---
title: HTTP Live Streaming -- ZH
date: 2024-07-29 20:49:40
tags:
categories:
- Doc
keywords:
- HTTP Live Streaming
- m38
- m3u8
copyright: Guader
copyright_author_href:
copyright_info:
---

From: [rfc8216](https://datatracker.ietf.org/doc/html/rfc8216)

> 说起来我之前一直把 m3u8 当成直播流的，类似于 flv 这种。    
> 突然发现这其实是个文本协议。

> 翻译对照

 english | chinese
  ---    | ---
unbounded streams | 无界流
Memo | 备忘录
HTTP Live Streaming | HTTP 实时流传输
Media Playlist | 媒体播放列表
Master Playlist | 主播放列表
Variant Streams | 变体流
resolution | 分辨率
Program Association Table (PAT) | 节目关联表
Program Map Table (PMT) | 程序映射表
Service Description Table (SDT) | 节目业务描述表
Media Initialization Section | 媒体初始化部分
Movie Box (moov) | 影片盒
Media Data Box (mdat) | 媒体数据盒
Movie Fragment Box (moof) | 影片片段盒
File Type Box (ftyp) | 文件类型盒
Track Box (trak) | 轨道盒
Track Fragment Box (traf) | 轨道片段盒
Movie Header Box (mvhd) | 影片头部盒
Track Header Box (tkhd) | 轨道头部盒
Movie Extends Box (mvex) | 影片扩展盒
Common Media Application Format (CMAF) | 通用媒体应用格式
Track Fragment Decode Time Box (tfdt) | 轨道片段解码时间盒
Packed Audio | 压缩音频
Advanced Audio Coding (AAC) | 高级音频编码
Audio Data Transport Stream (ADTS) | 音频数据传输流
Private frame (PRIV) | 私有帧
Program Elementary Stream | 程序基本流
Byte Order Mark (BOM) | 字节顺序标记



> 摘要

本文档描述了传输多媒体数据的无界流的协议。
它指定了文件的数据格式，以及流的服务器（发送方）和客户端（接收方）要采取的操作。
它描述了这个协议的 7 号版本。


> 这个备忘录的状态

本文档不是互联网标准轨迹规范；其发布仅供参考。

这是对 RFC 系列的贡献，独立于任何其他 RFC 流。
RFC 编辑者自行决定发布此文档，并且不对其实施或部署的价值作出任何声明。
经 RFC 编辑者批准发布的文档不是任何级别的互联网标准的候选文档；请参阅 [RFC 7841 的第 2 节](https://datatracker.ietf.org/doc/html/rfc7841#section-2)。

有关本文档的当前状态、任何勘误表以及如何提供反馈的信息可在 [http://www.rfc-editor.org/info/rfc8216](http://www.rfc-editor.org/info/rfc8216) 上获得。

> 版权声明 - 省略


## 1. Introduction to HTTP Live Streaming

HTTP 实时流传输是一种可靠且经济高效的通过互联网传输连续长视频的方法。
它允许接收方根据当前网络条件调整媒体的比特率，以保持最佳质量的不间断播放。
它支持间隙内容边界。它为媒体加密提供了灵活的框架。
它可以高效地提供相同内容的多种版本，例如音频翻译。
它与大规模 HTTP 缓存基础设施兼容，以支持向大量受众交付。

自 2009 年互联网草案首次发布以来，HTTP 实时流式传输已被众多内容制作者、工具供应商、分销商和设备制造商实施和部署。
在随后的八年中，该协议通过与各种媒体流式传输实施者的广泛审查和讨论而得到完善。

本文档旨在通过描述媒体传输协议来促进 HTTP Live Streaming 实现之间的互操作性。
使用此协议，客户端可以从服务器接收连续的媒体流以进行并发演示。

该文档描述了这个协议的 7 号版本。

> 最新的已经到 23 了 - 2024/08/13


## 2. Overview

多媒体演示由统一资源标识符 (URI) [[RFC3986][RFC3986]] 指定给播放列表。

播放列表可以是媒体播放列表或主播放列表。
两者都是包含 URI 和描述性标签的 UTF-8 文本文件。

一个媒体播放列表包含媒体片段列表，按顺序播放时将播放多媒体演示。

这是一个媒体播放列表的示例:

```m3u
#EXTM3U
#EXT-X-TARGETDURATION:10

#EXTINF:9.009,
http://media.example.com/first.ts
#EXTINF:9.009,
http://media.example.com/second.ts
#EXTINF:3.003,
http://media.example.com/third.ts
```

第一行是格式标识符标签 `#EXTM3U` 。
包含 `#EXT-X-TARGETDURATION` 的行表示所有媒体片段的长度不超过 10 秒。
然后声明了三个媒体片段。
第一个和第二个媒体片段的长度为 9.009 秒；第三个媒体片段的长度为 3.003 秒。

要播放此播放列表，客户端首先下载它，然后下载并播放其中声明的每个媒体片段。
客户端按照本文档中的说明重新加载播放列表以发现任何添加的片段。
数据 *应该* 通过 HTTP [[RFC7230][RFC7230]] 传输，不过一般来说，URI 可以指定任何能够按需可靠传输指定资源的协议。

更复杂的演示可以通过主播放列表来描述。
主播放列表提供一组变体流，每个变体流描述同一内容的不同版本。

一个变体流包含一个媒体播放列表，该列表指定以特定比特率、特定格式和特定分辨率编码的包含视频的媒体。

变体流还可以指定一组版本。
版本是内容的替代版本，例如以不同语言制作的音频或是用摄像机从不同角度录制的视频。

客户端应根据网络状况在不同的变体流之间切换。
客户端应根据用户偏好选择版本。

本文档中关键字 "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY" 以及 "OPTIONAL" 都和 [[BCP 14](https://datatracker.ietf.org/doc/html/bcp14)], [[RFC2119](https://datatracker.ietf.org/doc/html/rfc2119)] [[RFC8174](https://datatracker.ietf.org/doc/html/rfc8174)] 中描述的一致。


## 3. Media Segments

媒体播放列表包含一系列构成整体演示的媒体片段。
一个媒体片段由一个 URI 和可选的字节范围指定。

每个媒体片段的持续时间由媒体播放列表中的 `EXTINF` 标签指示（第 [4.3.2.1](#4321-extinf) 节）。

媒体播放列表中的每个片段都有一个唯一的整数作为媒体序列号。
媒体播放列表中第一个片段的媒体序列号为 0 或在播放列表中声明（第 [4.3.3.2](#4332-ext-x-media-sequence) 节）。
其他每个片段的媒体序列号等于其前一个片段的媒体序列号加一。

每个媒体片段必须携带编码比特流的延续，从具有前一个媒体序列号的片段末尾开始，其中时间戳和连续性计数器等系列中的值必须不间断地连续。
唯一的例外是媒体播放列表中出现的第一个媒体片段和明确标记为不连续的媒体片段（第 [4.3.2.3](#4323-ext-x-discontinuity) 节）。
未标记的媒体不连续性可能会触发播放错误。

任何包含视频的媒体片段都应包含足够的信息来初始化视频解码器并解码包含片段中最后一帧的连续帧集；
如果片段中有足够的信息来解码片段中的所有帧，则网络效率会得到优化。
例如，任何包含 H.264 视频的媒体片段都应包含瞬时解码刷新 (IDR)；
第一个 IDR 之前的帧将被下载但可能会被丢弃。


### 3.1. Supported Media Segment Formats

所有媒体段 *必须* 采用本节中描述的格式。 
其他格式媒体文件的传输未定义。

某些媒体格式需要使用通用的字节序列来初始化解析器，然后才能解析媒体段。
此格式特定的序列称为媒体初始化部分。
媒体初始化部分可以通过 `EXT-X-MAP` 标签指定（第 [4.3.2.5](#4325-ext-x-map) 节）。
媒体初始化部分 *不得* 包含样本数据。


### 3.2. MPEG-2 Transport Streams

MPEG-2 传输流由 [[ISO 13818][ISO-13818]] 指定。

MPEG-2 传输流段的媒体初始化部分是一个节目关联表 (PAT)，后面跟着一个节目映射表 (PMT)。

传输流片段 *必须* 包含单个 MPEG-2 节目；
多个节目传输流的播放尚未定义。
每个传输流片段 *必须* 包含一个 PAT 和一个 PMT，或者应用一个 `EXT-X-MAP` 标签（第 [4.3.2.5](#4325-ext-x-map) 节）。
片段的传输流的前两个数据包中应该包含 `EXT-X-MAP` 标签，否则 *应该* 是 PAT 和 PMT。


### 3.3. Fragmented MPEG-4

MPEG-4 片段由 ISO 基础媒体文件格式 [[ISOBMFF][ISOBMFF]] 指定。
与常规 MPEG-4 文件不同，常规 MPEG-4 文件具有包含样本表的影片盒（ `moov` ）和包含相应样本的媒体数据盒（ `mdat` ），
而 MPEG-4 片段由包含样本表子集的影片片段盒（ `moof` ）和包含这些样本的媒体数据盒组成。
使用 MPEG-4 片段确实需要影片盒进行初始化，
但该影片盒仅包含非样本特定信息，例如轨道和样本描述。

片段化 MPEG-4 (fMP4) 片段是 [[ISOBMFF][ISOBMFF]] 第 3 节定义的 “片段”，
包括 [[ISOBMFF][ISOBMFF]] 第 8.16 节中对媒体数据盒的限制。

fMP4 段的媒体初始化部分是一个 ISO 基础媒体文件，可以为该段初始化解析器。

广义上讲，fMP4 段和媒体初始化部分是 [[ISOBMFF][ISOBMFF]] 文件，也满足本节所述的约束。

fMP4 片段的媒体初始化部分 *必须* 包含一个文件类型盒（ `ftyp` ），
其中包含与 `iso6` 或更高版本兼容的牌子。
文件类型盒后面 *必须* 跟一个影片盒。
影片盒 *必须* 包含 fMP4 片段中每个轨道片段盒（ `traf` ）的轨道盒（ `trak` ），
并具有匹配的 `track_ID`。
每个轨道盒都 *应该* 包含一个样本表，但其样本计数 *必须* 为零。
影片头部盒（ `mvhd` ）和轨道头部盒（ `tkhd` ）的持续时间 *必须* 为零。
影片扩展盒（ `mvex` ）*必须* 跟在最后一个轨道盒后面。
请注意，通用媒体应用格式 (CMAF) 头部 [[CMAF][CMAF]] 满足所有这些要求。

在 fMP4 片段中，
每个轨道片段盒 *必须* 包含一个轨道片段解码时间盒（ `tfdt` ）。
fMP4 片段 *必须* 使用电影片段相对寻址。
fMP4 片段 *不得* 使用外部数据引用。
请注意，CMAF 片段满足这些要求。

播放列表中包含 `EXT-X-I-FRAMES-ONLY` 标签（第 [4.3.3.6](#4336-ext-x-i-frames-only) 节）的 fMP4 片段 *可以* 省略帧内编码帧（I 帧）样本数据之后的媒体数据盒部分。

媒体播放列表中的每个 fMP4 片段都 *必须* 应用 `EXT-X-MAP` 标签。


### 3.4. Packed Audio

压缩音频片段包含编码的音频样本和 ID3 标签，
它们以最少的帧和不包含每个样本时间戳的方式简单打包在一起。
支持的压缩音频格式包括采用音频数据传输流 (ADTS) 帧的高级音频编码 (AAC) [[ISO_13818_7][ISO-13818-7]]、MP3 [[ISO_13818_3][ISO-13818-3]]、AC-3 [AC_3][AC-3] 和增强型 AC-3 [[AC_3][AC-3]]。

压缩音频段没有媒体初始化部分.

每个压缩音频片段 *必须* 在片段开头使用 ID3 私有帧 (PRIV) 标签 [[ID3][ID3]] 来标示其第一个样本的时间戳。
ID3 PRIV 所有者标识符 *必须* 是 `"com.apple.streaming.transportStreamTimestamp"` 。
ID3 有效负载 *必须* 是 33 位 MPEG-2 程序基本流时间戳，
以大端八位字节数表示，其中高 31 位设置为零。
客户端 *不应该* 播放没有此 ID3 标签的压缩音频片段。


### 3.5. WebVTT

WebVTT 片段是 WebVTT [[WebVTT][WebVTT]] 文件的一部分。
WebVTT 片段带有字幕。

WebVTT 段的媒体初始化部分是 WebVTT 标头。

每个 WebVTT 片段 *必须* 包含要在片段 `EXTINF` 持续时间指示的时间段内显示的所有字幕提示。
每个提示的开始时间偏移和结束时间偏移必须指示该提示的总显示时间，即使提示时间范围的一部分超出了片段时长。
WebVTT 片段 *可以* 不包含任何提示；这表示在该时间段内不会显示任何字幕。

每个 WebVTT 段 *必须* 以 WebVTT 标头开头或应用 `EXT-X-MAP` 标签。

为了同步音频/视频和字幕之间的时间戳， *应该* 将 `X-TIMESTAMP-MAP` 元数据标头添加到每个 WebVTT 标头。
此标头将 WebVTT 提示时间戳映射到变体流的其他版本中的 MPEG-2 (PES) 时间戳。
其格式为：

```plain
X-TIMESTAMP-MAP=LOCAL:<cue time>,MPEGTS:<MPEG-2 time>
e.g.
X-TIMESTAMP-MAP=LOCAL:00:00:00.000,MPEGTS:900000
```

`LOCAL` 属性中的提示时间戳 *可以* 超出该片段涵盖的时间范围。

如果 WebVTT 段没有 `X-TIMESTAMP-MAP`，则客户端 *必须* 假定 WebVTT 提示时间 0 映射到 MPEG-2 时间戳 0。

当将 WebVTT 与 PES 时间戳同步时，客户端 *应该* 考虑 33 位 PES 时间戳已换行而 WebVTT 提示时间尚未换行的情况。


## 4. Playlists 

本节介绍 HTTP 实时流式传输使用的播放列表文件。
本节中的 *必须* 和 *不得* 规定了合法播放列表文件的语法和结构规则。
违反这些规则的播放列表无效；客户端 *必须* 无法解析它们。
请参阅[第 6.3.2 节](#632-loading-the-media-playlist-file)。

播放列表文件的格式源自 M3U [[M3U][M3U]] 播放列表文件格式，并从早期文件格式继承了两个标签：
`EXTM3U`（[第 4.3.1.1 节](#4311-extm3u)）和 `EXTINF`（[第 4.3.2.1 节](#4321-extinf)）。

在标签语法规范中，用 `<>` 括起来的字符串标识一个标签参数，其具体格式在标签定义中有描述。
如果参数进一步用 `[]` 括起来，则为可选参数，否则为必选参数。

每个播放列表文件 *必须* 可以通过其 URI 的路径部分或 HTTP 内容类型来识别。
在第一种情况下，路径必须以 `.m3u8` 或 `.m3u` 结尾。
在第二种情况下，HTTP 内容类型必须是 `application/vnd.apple.mpegurl` 或 `audio/mpegurl` 。
客户端 *应该* 拒绝解析未如此识别的播放列表。


### 4.1. Definition of a Playlist

播放列表文件必须采用 UTF-8 [[RFC3629][RFC3629]] 编码。
它们 *不得* 包含任何字节顺序标记 (BOM)；客户端 *应该* 无法解析包含 BOM 或未解析为 UTF-8 的播放列表。
播放列表文件 *不得* 包含 UTF-8 控制字符（`U+0000` 到 `U+001F` 和 `U+007F` 到 `U+009F`），但 CR (`U+000D`) 和 LF (`U+000A`) 除外。
所有字符序列 *必须* 根据 Unicode 规范化形式 `NFC` [[UNICODE][UNICODE]] 进行规范化。
请注意，US-ASCII [[US_ASCII][US-ASCII]] 符合这些规则。

播放列表文件中的行以单个换行符 `\n`或回车符加换行符 `\r\n` 结尾。
每行都是一个 URI、空白或以字符 `#` 开头。
空白行将被忽略。
*不得* 出现空格，除非元素中明确指定了空格。

以字符 `#` 开头的行是注释或标签。
标签以 `#EXT` 开头。
它们区分大小写。
以 `#` 开头的所有其他行都是注释， *应该* 被忽略。

URI 行标识媒体片段或播放列表文件（参见[第 4.3.4.2 节](#4342-ext-x-stream-inf)）。
每个媒体片段由 URI 和适用于它的标签指定。

如果播放列表中的所有 URI 行都标识媒体段，则播放列表为媒体播放列表。
如果播放列表中的所有 URI 行都标识媒体播放列表，则播放列表为主播放列表。
播放列表 *必须* 是媒体播放列表或主播放列表；所有其他播放列表均无效。

播放列表中的 URI（无论是 URI 行还是标签的一部分）都 *可以* 是相对的。
任何相对 URI 都被视为相对于包含它的播放列表的 URI。

媒体播放列表的持续时间是其中媒体片段的持续时间的总和。

媒体片段的片段比特率等于媒体片段的大小除以其 `EXTINF` 持续时间（[第 4.3.2.1 节](#4321-extinf)）。
请注意，这包括容器开销，但不包括交付系统施加的开销，例如 HTTP、TCP 或 IP 标头。

媒体播放列表的峰值片段比特率是任何连续片段集的最大比特率，其总时长介于目标时长的 0.5 到 1.5 倍之间。
一组片段的比特率是通过将片段大小的总和除以片段时长的总和来计算的。

媒体播放列表的平均片段比特率是媒体播放列表中每个媒体片段的大小（以位为单位）的总和除以媒体播放列表的持续时间。
请注意，这包括容器开销，但不包括 HTTP 或交付系统施加的其他开销。


### 4.2. Attribute Lists

某些标签的值是属性列表。
属性列表是以逗号分隔的属性/值对列表，其中没有空格。

属性/值对具有以下语法：

`AttributeName=AttributeValue`

`AttributeName` 是一个不带引号的字符串，包含字符 `[A..Z]` 、 `[0..9]` 以及 `-` 。
因此， `AttributeName` 仅包含大写字母，不包含小写字母。
`AttributeName` 和 `=` 字符之间以及 `=` 字符和 `AttributeValue` 之间 *不得* 有任何空格。

`AttributeValue` 是以下之一：

- 十进制整数：来自集合 `[0..9]` 的不带引号的字符串，表示以 10 为基数的整数，范围从 0 到 2^64-1 (18446744073709551615)。十进制整数的长度可以是 1 到 20 个字符。
- 十六进制序列：来自集合 `[0..9]` 和 `[A..F]` 的不带引号的字符串，以 `0x` 或 `0X` 为前缀。十六进制序列的最大长度取决于其属性名称。
- 十进制浮点数：来自集合 `[0..9]` 和 `.` 的不带引号的字符串，以十进制位置表示法表示非负浮点数。
- 有符号十进制浮点数：由集合 `[0..9]` 、 `-` 和 `.` 组成的不带引号的字符串，以十进制位置表示法表示有符号浮点数。
- 引号字符串：一对双引号 (0x22) 内的字符串。以下字符 *不得* 出现在引号字符串中：换行符 (0xA)、回车符 (0xD) 或双引号 (0x22)。引号字符串属性值 *应该* 这样构造：逐字节比较足以测试两个引号字符串属性值是否相等。请注意，这意味着区分大小写的比较。
- 枚举字符串：来自由 `AttributeName` 明确定义的集合的未加引号的字符串。枚举字符串绝不会包含双引号 `"` 、逗号 `,` 或空格。
- 十进制分辨率：两个十进制整数，以 `x` 字符分隔。第一个整数是水平像素尺寸（宽度）；第二个是垂直像素尺寸（高度）。

给定 `AttributeName` 的 `AttributeValue` 类型由属性定义指定。

给定的 `AttributeName` *不得* 在给定的属性列表中出现多次。
客户端 *应该* 拒绝解析此类播放列表。


### 4.3. Playlist Tags 

播放列表标签指定播放列表的全局参数或有关其后出现的媒体片段或媒体播放列表的信息。


#### 4.3.1. Basic Tags

媒体播放列表和主播放列表中都允许使用这些标签。


##### 4.3.1.1. EXTM3U

`EXTM3U` 标签表示该文件是扩展 M3U [[M3U][M3U]] 播放列表文件。
它 *必须* 是每个媒体播放列表和每个主播放列表的第一行。其格式为：

`#EXTM3U` 


##### 4.3.1.2. EXT-X-VERSION

`EXT-X-VERSION` 标签表示播放列表文件、其关联媒体及其服务器的兼容版本。
`EXT-X-VERSION` 标签适用于整个播放列表文件。 
它的格式为：

`EXT-X-VERSION:<n>`

其中 `n` 是一个整数，表示协议兼容版本号。

它 *必须* 出现在包含与协议版本 1 不兼容的标签或属性的所有播放列表中，以支持与旧客户端的互操作性。
[第 7 节](#7-protocol-version-compatibility) 指定了任何给定播放列表文件的兼容性版本号的最小值。

播放列表文件 *不得* 包含多个 `EXT-X-VERSION` 标签。
如果客户端遇到包含多个 `EXT-X-VERSION` 标签的播放列表，则 *必须* 无法解析它。


#### 4.3.2. Media Segment Tags

每个媒体片段由一系列媒体片段标签和 URI 指定。
一些媒体片段标签仅适用于下一个片段；其他则适用于所有后续片段，直到出现相同标签的另一个实例。

媒体片段标签 *不得* 出现在主播放列表中。
客户端 *必须* 无法解析同时包含媒体片段标签和主播放列表标签的播放列表（[第 4.3.4 节](#434-master-playlist-tags)）。


##### 4.3.2.1. EXTINF

`EXTINF` 标签指定媒体片段的持续时间。
它仅适用于下一个媒体片段。此标签对于每个媒体片段都是 *必需* 的。
其格式为：

`#EXTINF:<duration>,[<title>]`

其中， `duration` 是十进制浮点数或十进制整数（如 [第 4.2 节](#42-attribute-lists) 所述），用于指定媒体片段的持续时间（以秒为单位）。
持续时间 *应该* 为十进制浮点数，且精度应足够高，以避免在片段持续时间累积时出现可察觉的误差。
但是，如果兼容版本号小于 3，则持续时间 *必须* 为整数。
以整数形式报告的持续时间 *应该* 四舍五入为最接近的整数。
逗号后面的行的其余部分是媒体片段的可选人类可读信息标题，以 UTF-8 文本表示。


##### 4.3.2.2. EXT-X-BYTERANGE

`EXT-X-BYTERANGE` 标签表示媒体片段是其 URI 所标识的资源的子范围。
它仅适用于播放列表中紧随其后的下一个 URI 行。
其格式为：

`#EXT-X-BYTERANGE:<n>[@<o>]`

其中 `n` 是一个十进制整数，表示子范围的长度（以字节为单位）。
如果存在，`o` 是一个十进制整数，表示子范围的开始，以相对于资源开头的字节偏移量表示。
如果 `o` 不存在，则子范围从上一个媒体段的子范围后面的下一个字节开始。

如果 `o` 不存在，则前一个媒体段 *必须* 出现在播放列表文件中，并且 *必须* 是同一媒体资源的子范围，或者媒体段未定义，并且客户端 *必须* 无法解析播放列表。

没有 `EXT-X-BYTERANGE` 标签的媒体片段由其 URI 标识的整个资源组成。

使用 `EXT-X-BYTERANGE` 标签 *需要* 兼容版本号为 4 或更高。


##### 4.3.2.3. EXT-X-DISCONTINUITY

`EXT-X-DISCONTINUITY` 标签表示其后的媒体片段与之前的媒体片段之间的不连续性。

它的格式为：

`#EXT-X-DISCONTINUITY`

如果以下任何特征发生变化，则 *必须* 存在 `EXT-X-DISCONTINUITY` 标签：

- 文件格式
- 轨道的数量、类型和标识符
- 时间戳序列

如果以下任何特征发生变化，则 *应该* 存在 `EXT-X-DISCONTINUITY` 标签：

- 编码参数
- 编码序列

有关 `EXT-X-DISCONTINUITY` 标签的更多信息，请参阅第 [3](#3-media-segments) 、[6.2.1](#621-general-server-responsibilities) 和 [6.3.3](#633-playing-the-media-playlist-file) 节。


##### 4.3.2.4. EXT-X-KEY

媒体片段 *可以* 加密。`EXT-X-KEY` 标签指定如何解密它们。
它适用于每个媒体片段和每个由 `EXT-X-MAP` 标签声明的媒体初始化部分，该标签出现在播放列表文件中具有相同 `KEYFORMAT` 属性的下一个 `EXT-X-KEY` 标签之间（或播放列表文件的末尾）。
如果两个或多个具有不同 `KEYFORMAT` 属性的 `EXT-X-KEY` 标签最终产生相同的解密密钥，则它们可以应用于同一个媒体片段。
其格式为：

`#EXT-X-KEY:<attribute-list>`

定义了以下属性：

- 方法 ( `METHOD` )
    该值是一个枚举字符串，用于指定加密方法。
    此属性是 *必需* 的。
    定义的方法有：`NONE` 、 `AES-128` 和 `SAMPLE-AES` 。
    加密方法为 `NONE` 表示媒体片段未加密。
  - 如果加密方法为 `NONE` ，则其他属性 *不得* 存在。
  - `AES-128` 加密方法表示媒体段已使用高级加密标准 (AES) [[AES_128](#52-iv-for-aes-128)] 和 128 位密钥、密码块链接 (CBC) 和公钥加密标准 #7 (PKCS7) 填充 [[RFC5652][RFC5652]] 完全加密。
    CBC 在每个段边界上重新启动，使用初始化向量 (IV) 属性值或媒体序列号作为 IV；请参阅 [第 5.2 节](#52-iv-for-aes-128) 。
  - `SAMPLE-AES` 加密方法表示媒体片段包含使用高级加密标准 [[AES_128](#52-iv-for-aes-128)] 加密的媒体样本，例如音频或视频。
    这些媒体流如何加密并封装在片段中取决于片段的媒体编码和媒体格式。
    `fMP4` 媒体片段使用通用加密 [[COMMON_ENC][COMMON-ENC]] 的 `cbcs` 方案进行加密。
    HTTP 实时流 (HLS) 示例加密规范 [[SampleEnc][Sample-Enc]] 中描述了包含 H.264 [[H_264][H-264]]、AAC [[ISO_14496][ISO-14496]]、AC-3 [[AC_3][AC-3]] 和增强型 AC-3 [[AC_3][AC-3]] 媒体流的其他媒体片段格式的加密。
    IV 属性 *可以* 存在；请参阅 [第 5.2 节](#52-iv-for-aes-128) 。
- URI
    该值是一个带引号的字符串，其中包含指定如何获取密钥的 URI。
    除非 `METHOD` 为 `NONE` ，否则此属性是 *必需* 的。
- IV
    该值是一个十六进制序列，指定与密钥一起使用的 128 位无符号整数初始化向量。
    使用 IV 属性 *需要* 兼容版本号为 2 或更高。
    有关何时使用 IV 属性，请参阅 [第 5.2 节](#52-iv-for-aes-128) 。
- KEYFORMAT
    该值是一个带引号的字符串，用于指定密钥在 URI 标识的资源中的表示方式；
    有关详细信息，请参阅 [第 5 节](#5-key-files) 。
    此属性是可选的；如果不存在，则表示隐式值为 `identity` 。
    使用 `KEYFORMAT` 属性要求兼容版本号为 5 或更高。
- KEYFORMATVERSIONS
    该值是一个带引号的字符串，其中包含一个或多个由 `/` 字符分隔的正整数（例如，`1` 、 `1/2` 或 `1/2/5` ）。
    如果定义了特定 `KEYFORMAT` 的多个版本，则可以使用此属性来指示此实例符合哪个版本。
    此属性是可选的；如果不存在，则其值被视为 `1` 。
    使用 `KEYFORMATVERSIONS` 属性 *需要* 兼容版本号为 5 或更高。

如果媒体播放列表文件不包含 `EXT-X-KEY` 标签，则媒体片段未加密。

有关密钥文件的格式请参见 [第 5 节](#5-key-files) ，有关媒体片段加密的更多信息请参见第 [5.2](#52-iv-for-aes-128) 、 [6.2.3](#623-encrypting-media-segments) 和 6.3.6 节。


##### 4.3.2.5. EXT-X-MAP

`EXT-X-MAP` 标签指定如何获取解析适用媒体片段所需的媒体初始化部分（ [第 3 节](#3-media-segments) ）。
它适用于播放列表中出现在其后面的每个媒体片段，直到下一个 `EXT-X-MAP` 标签或播放列表结束。

它的格式为：

`#EXT-X-MAP:<attribute-list>` 

定义了以下属性：

- URI
    该值是一个带引号的字符串，其中包含一个 URI，用于标识包含媒体初始化部分的资源。
    此属性是 *必需* 的。
- BYTERANGE
    该值是一个带引号的字符串，指定 URI 属性所标识的资源中的字节范围。
    此范围 *应该* 仅包含媒体初始化部分。
    字节范围的格式在 [第 4.3.2.2 节](#4322-ext-x-byterange) 中描述。​
    ​此属性是可选的；如果不存在，则字节范围是 URI 所指示的整个资源。

当播放列表中的第一个媒体片段（即 I 帧）（或 `EXT-X-DISCONTINUITY` 标签后面的第一个片段）没有紧跟在其资源开头的媒体初始化部分之后时， *应该* 为带有 `EXT-X-I-FRAMES-ONLY` 标签的播放列表中的媒体片段提供 `EXT-X-MAP` 标签。

在包含 `EXT-X-I-FRAMES-ONLY` 标签的媒体播放列表中使用 `EXT-X-MAP` 标签 *需要* 兼容版本号为 5 或更高。
在不包含 `EXT-X-I-FRAMES-ONLY` 标签的媒体播放列表中使用 `EXT-X-MAP` 标签 *需要* 兼容版本号为 6 或更高。

如果 `EXT-X-MAP` 标签声明的媒体初始化部分使用 `AES-128` 方法加密，则适用于 `EXT-X-MAP` 的 `EXT-X-KEY` 标签的 IV 属性是 *必需* 的。


##### 4.3.2.6. EXT-X-PROGRAM-DATE-TIME

`EXT-X-PROGRAM-DATE-TIME` 标签将媒体片段的第一个样本与绝对日期 和/或 时间关联起来。
它仅适用于下一个媒体片段。
其格式为：

`#EXT-X-PROGRAM-DATE-TIME:<date-time-msec>`

其中 `date-time-msec` 是 ISO/IEC 8601:2004 [[ISO_8601][ISO-8601]] 日期/时间表示形式，例如 `YYYY-MM-DDThh:mm:ss.SSSZ`。
它 *应该* 指示时区和秒的小数部分，精确到毫秒。

例如：

`#EXT-X-PROGRAM-DATE-TIME:2010-02-19T14:54:23.031+08:00` 

有关 `EXT-X-PROGRAM-DATE-TIME` 标签的更多信息，请参阅第 [6.2.1](#621-general-server-responsibilities) 和 [6.3.3](#633-playing-the-media-playlist-file) 节。


##### 4.3.2.7. EXT-X-DATERANGE

`EXT-X-DATERANGE` 标签将日期范围（即由开始日期和结束日期定义的时间范围）与一组属性/值对关联起来。
其格式为：

`#EXT-X-DATERANGE:<attribute-list>`

其中定义的属性是：

- ID
    带引号的字符串，用于唯一标识播放列表中的日期范围。此属性是 *必需* 的。
- CLASS
    客户端定义的带引号的字符串，用于指定一组属性及其关联的值语义。
    具有相同 CLASS 属性值的所有日期范围都 *必须* 遵循这些语义。此属性是可选的。
- START-DATE
    包含日期范围开始的 ISO-8601 日期的带引号的字符串。此属性是 *必需* 的。
- END-DATE
    包含日期范围结束的 ISO-8601 日期的带引号的字符串。
    它 *必须* 等于或晚于 START-DATE 属性的值。此属性是可选的。
- DURATION
    日期范围的持续时间以十进制浮点秒数表示。它 *不得* 为负数。
    单个时刻（例如越过终点线）的持续时间 *应该* 为 0。此属性是可选的。
- PLANNED-DURATION
    日期范围的预期持续时间，以十进制浮点数秒表示。它 *不得* 为负数。
    此属性应用于指示实际持续时间尚不清楚的日期范围的预期持续时间。
    它是可选的。
- `X-<client-attribute>`
    `X-` 前缀定义了为客户端定义属性保留的命名空间。客户端属性必须是合法的 `AttributeName` 。
    客户端在定义自己的属性名称时应使用反向 DNS 语法以避免冲突。
    属性值必须是带引号的字符串、十六进制序列或十进制浮点数。
    客户端定义属性的一个示例是 `X-COM-EXAMPLE-AD-ID="XYZ123"`。
    这些属性是可选的。
- `SCTE35-CMD, SCTE35-OUT, SCTE35-IN`
    用于承载 `SCTE-35` 数据；有关详细信息，请参阅 [第 4.3.2.7.1 节](#43271-mapping-scte-35-into-ext-x-daterange) 。这些属性是可选的。
- END-ON-NEXT
    枚举字符串，其值 *必须* 为 YES。此属性表示包含它的范围的结束等于其后续范围的 `START-DATE` 。
    后续范围是同一 `CLASS` 的日期范围，其最早的 `START-DATE` 晚于相关范围的 `START-DATE` 。
    此属性是可选的。

具有 `END-ON-NEXT=YES` 属性的 `EXT-X-DATERANGE` 标签 *必须* 具有 `CLASS` 属性。
具有相同 `CLASS` 属性的其他 `EXT-X-DATERANGE` 标签 *不得* 指定重叠的日期范围。

具有 `END-ON-NEXT=YES` 属性的 `EXT-X-DATERANGE` 标签 *不得* 包含 `DURATION` 或 `END-DATE` 属性。

如果日期范围没有 `DURATION` 、 `END-DATE` 或 `END-ON-NEXT=YES` 属性，则即使具有 `PLANNED-DURATION` ，其持续时间也是未知的。

如果播放列表包含 `EXT-X-DATERANGE` 标签，它还必须包含至少一个 `EXT-X-PROGRAM-DATE-TIME` 标签。

如果播放列表包含两个具有相同 ID 属性值的 `EXT-X-DATERANGE` 标签，则两个标签中出现的任何 `AttributeName` 都必须具有相同的 `AttributeValue` 。

如果日期范围同时包含 `DURATION` 属性和 `END-DATE` 属性，则 `END-DATE` 属性的值 *必须* 等于 `START-DATE` 属性的值加上 `DURATION` 属性的值。

客户端 *应该* 忽略具有非法语法的 `EXT-X-DATERANGE` 标签。


###### 4.3.2.7.1. Mapping SCTE-35 into EXT-X-DATERANGE

根据 SCTE-35 规范 [[SCTE35][SCTE35]]，源媒体中携带的拼接信息 *可以* 使用 `EXT-X-DATERANGE` 标签在媒体播放列表中表示。

每个包含 `splice_null()` 、 `splice_schedule()` 、 `bandwidth_reservation()` 或 `private_cmd()` 的 SCTE-35 `splice_info_section()` 都 *应该* 由带有 `SCTE35-CMD` 属性的 `EXT-X-DATERANGE` 标签表示，该属性的值是 `splice_info_section()` 的大端二进制表示形式，以十六进制序列表示。

由一对 `splice_insert()` 命令发出信号的 SCTE-35 拼接输出/输入对 *应该* 由一个或多个带有相同 ID 属性的 `EXT-X-DATERANGE` 标签表示，这些标签对于该拼接输出/输入对 *必须* 是唯一的。
"输出" `splice_info_section()`（将 `out_of_network_indicator` 设置为 1） *必须* 放置在 `SCTE35-OUT` 属性中，格式与 `SCTE35-CMD` 相同。
"输入" `splice_info_section()`（将 `out_of_network_indicator` 设置为 0） *必须* 放置在 `SCTE35-IN` 属性中，格式与 `SCTE35-CMD` 相同。

由一对 `time_signal()` 命令发出信号的 SCTE-35 拼接输出/输入对，每个命令都携带一个分段描述符 ( `segmentation_descriptor()` )， *应该* 由一个或多个携带相同 ID 属性的 `EXT-X-DATERANGE` 标签表示，这些标签对于该拼接输出/输入对 *必须* 是唯一的。
"输出" `splice_info_section()` *必须* 放置在 `SCTE35-OUT` 属性中；
"输入" `splice_info_section()` *必须* 放置在 `SCTE35-IN` 属性中。

即使两个或多个 `segmentation_descriptor()` 到达同一个 `splice_info_section()` ，不同类型的分段（如 `segmentation_descriptor()` 中的 `segmentation_type_id` 所示）也 *应该* 由单独的 `EXT-X-DATERANGE` 标签表示。
在这种情况下，每个 `EXT-X-DATERANGE` 标签将具有一个 `SCTE35-OUT` 、 `SCTE35-IN` 或 `SCTE35-CMD` 属性，其值是整个  `splice_info_section()` 。

未发出拼接输出点或输入点信号的 SCTE-35 `time_signal()` 命令 *应该* 由具有 `SCTE35-CMD` 属性的 `EXT-X-DATERANGE` 标签表示。

包含 `SCTE35-OUT` 属性的 `EXT-X-DATERANGE` 标签的 `START-DATE` *必须* 是与该拼接的程序时间相对应的日期和时间。

包含 `SCTE35-CMD` 的 `EXT-X-DATERANGE` 标签的 `START-DATE` 必须是命令中 `splice_time()` 指定的日期和时间，或者如果命令未指定 `splice_time()` ，则 *必须* 是该命令出现在源流中的程序时间。

包含 `SCTE35-OUT` 属性的 `EXT-X-DATERANGE` 标签 *可以* 包含 `PLANNED-DURATION` 属性。其值 *必须* 是拼接的计划持续时间。

包含 `SCTE35-IN` 属性的 `EXT-X-DATERANGE` 标签的 `DURATION` *必须* 是相应出点和入点之间的实际（非计划）程序持续时间。

包含 `SCTE35-IN` 属性的 `EXT-X-DATERANGE` 标签的 `END-DATE` *必须* 是该入点的实际（非计划）程序日期和时间。

如果在将 `SCTE35-OUT` 属性添加到播放列表时不知道实际的结束日期和时间，则 *不得* 存在 `DURATION` 属性和 `END-TIME` 属性；一旦建立，就应该通过另一个 `EXT-X-DATERANGE` 标签来发出拼接的实际结束日期的信号。

已取消的拼接 *不应该* 作为 `EXT-X-DATERANGE` 标签出现在播放列表中。

宣布拼接的 `EXT-X-DATERANGE` 标签 *应该* 与最后一个预拼接媒体片段同时添加到播放列表中，或者如果可能的话更早地添加到播放列表中。

`EXT-X-DATERANGE` 标签的 ID 属性 *可以* 包含 `splice_event_id` 和/或 `fragmentation_event_id` ，但它在播放列表中 *必须* 是唯一的。
如果有可能重复使用 SCTE-35 ID，则 ID 属性值必须包含歧义消除信息，例如日期或序列号。


#### 4.3.3. Media Playlist Tags

媒体播放列表标签描述媒体播放列表的全局参数。
任何媒体播放列表中每种类型的媒体播放列表标签 *不得* 超过一个。

媒体播放列表标签 *不得* 出现在主播放列表中。


##### 4.3.3.1. EXT-X-TARGETDURATION

`EXT-X-TARGETDURATION` 标签指定媒体片段的最大时长。
播放列表文件中每个媒体片段的 `EXTINF` 时长（四舍五入为最接近的整数） *必须* 小于或等于目标时长；较长的片段可能会触发播放停顿或其他错误。
它适用于整个播放列表文件。其格式为：

`#EXT-X-TARGETDURATION:<s>`

其中 `s` 是一个十进制整数，表示目标持续时间（以秒为单位）。
`EXT-X-TARGETDURATION` 标签是必需的。


##### 4.3.3.2. EXT-X-MEDIA-SEQUENCE

`EXT-X-MEDIA-SEQUENCE` 标签表示播放列表文件中出现的第一个媒体片段的媒体序列号。
其格式为：

`#EXT-X-MEDIA-SEQUENCE:<number>`

其中 `number` 是十进制整数。

如果媒体播放列表文件不包含 `EXT-X-MEDIA-SEQUENCE` 标签，则媒体播放列表中第一个媒体段的媒体序列号应被视为 0。
客户端不得假设不同媒体播放列表中具有相同媒体序列号的段包含匹配的内容（参见 [第 6.3.2 节](#632-loading-the-media-playlist-file) ）。

媒体段的 URI 不需要包含其媒体序列号。

有关设置 `EXT-X-MEDIA-SEQUENCE` 标签的更多信息，请参阅第 [6.2.1](#621-general-server-responsibilities) 和 [6.3.5](#635-determining-the-next-segment-to-load) 节。

`EXT-X-MEDIA-SEQUENCE` 标签 *必须* 出现在播放列表中的第一个媒体片段之前。


##### 4.3.3.3. EXT-X-DISCONTINUITY-SEQUENCE

`EXT-X-DISCONTINUITY-SEQUENCE` 标签允许同一变体流的不同版本或媒体播放列表中具有 `EXT-X-DISCONTINUITY` 标签的不同变体流之间的同步。

它的格式为：

`#EXT-X-DISCONTINUITY-SEQUENCE:<number>`

其中 `number` 是十进制整数。

如果媒体播放列表不包含 `EXT-X-DISCONTINUITY-SEQUENCE` 标签，则播放列表中第一个媒体段的不连续序列号 *应该* 被视为 0。

`EXT-X-DISCONTINUITY-SEQUENCE` 标签 *必须* 出现在播放列表中的第一个媒体片段之前。

`EXT-X-DISCONTINUITY-SEQUENCE` 标签 *必须* 出现在任何 `EXT-X-DISCONTINUITY` 标签之前。

有关设置 `EXT-X-DISCONTINUITY-SEQUENCE` 标签值的更多信息，请参阅第 [6.2.1](#621-general-server-responsibilities) 和 [6.2.2](#622-live-playlists) 节。


##### 4.3.3.4. EXT-X-ENDLIST

`EXT-X-ENDLIST` 标签表示不再向媒体播放列表文件添加媒体片段。
它 *可以* 出现在媒体播放列表文件的任何位置。
其格式为：

`#EXT-X-ENDLIST`


##### 4.3.3.5. EXT-X-PLAYLIST-TYPE

`EXT-X-PLAYLIST-TYPE` 标签提供有关媒体播放列表文件的可变性信息。
它适用于整个媒体播放列表文件。
它是可选的。其格式为：

`#EXT-X-PLAYLIST-TYPE:<type-enum>`

其中 `type-enum` 是 `EVENT` 或 `VOD` 。

[第 6.2.1 节](#621-general-server-responsibilities) 定义了 `EXT-X-PLAYLIST-TYPE` 标签的含义。

如果 `EXT-X-PLAYLIST-TYPE` 值为 `EVENT` ，则只能将媒体片段添加到媒体播放列表的末尾。
如果 `EXT-X-PLAYLIST-TYPE` 值为视频点播 ( `VOD` )，则媒体播放列表不能更改。

如果媒体播放列表中省略了 `EXT-X-PLAYLIST-TYPE` 标签，则可以根据 [第 6.2.1 节](#621-general-server-responsibilities) 中的规则更新播放列表，而无需其他限制。
例如， *可以* 更新实时播放列表（ [第 6.2.2 节](#622-live-playlists) ）以按媒体片段出现的顺序删除它们。


##### 4.3.3.6. EXT-X-I-FRAMES-ONLY

`EXT-X-I-FRAMES-ONLY` 标签表示播放列表中的每个媒体片段都描述单个 I 帧。
I 帧是经过编码的视频帧，其编码不依赖于任何其他帧。
I 帧播放列表可用于特技播放，例如快进、快速倒退和滚动播放。

`EXT-X-I-FRAMES-ONLY` 标签适用于整个播放列表。其格式为：

`#EXT-X-I-FRAMES-ONLY`

在带有 `EXT-X-I-FRAMES-ONLY` 标签的播放列表中，媒体片段持续时间（ `EXTINF` 标签值）是媒体片段中 I 帧的显示时间与播放列表中下一个 I 帧的显示时间之间的时间，或者如果它是播放列表中的最后一个 I 帧，则为显示结束的时间。

包含 I 帧段的媒体资源 *必须* 以媒体初始化段（ [第 3 节](#3-media-segments) ）开头，或附带指示媒体初始化段的 `EXT-X-MAP` 标记，以便客户端可以按任意顺序加载和解码 I 帧段。
应用了 `EXT-X-BYTERANGE` 标记的 I 帧段的字节范围（ [第 4.3.2.2 节](#4322-ext-x-byterange) ） *不得* 包括其媒体初始化段；客户端可以假设媒体初始化段由 `EXT-X-MAP` 标记定义，或位于从资源开头到该资源中第一个 I 帧段的偏移量之间。

使用 `EXT-X-I-FRAMES-ONLY` 需要兼容版本号为 4 或更高。


#### 4.3.4. Master Playlist Tags

主播放列表标签定义演示的变体流、版本和其他全局参数。

主播放列表标签 *不得* 出现在媒体播放列表中；
客户端 *不得* 解析任何同时包含主播放列表标签和媒体播放列表标签或媒体片段标签的播放列表。


##### 4.3.4.1. EXT-X-MEDIA

`EXT-X-MEDIA` 标签用于关联包含相同内容的备选版本（ [第 4.3.4.2.1 节](#43421-alternative-renditions) ）的媒体播放列表。
例如，三个 `EXT-X-MEDIA` 标签可用于识别包含相同演示的英语、法语和西班牙语版本的纯音频媒体播放列表。
或者，两个 `EXT-X-MEDIA` 标签可用于识别显示两个不同摄像机角度的纯视频媒体播放列表。

它的格式为：

`#EXT-X-MEDIA:<attribute-list>` 

定义了以下属性：

- TYPE
    该值是一个枚举字符串；有效字符串为 AUDIO 、 VIDEO 、 SUBTITLES 和 CLOSED-CAPTIONS。
    此属性是必需的。
    通常，隐藏字幕 [[CEA608][CEA608]] 媒体在视频流中传输。
    因此，类型为 `CLOSED-CAPTIONS` 的 `EXT-X-MEDIA` 标签不指定 Rendition；隐藏字幕媒体存在于每个视频 Rendition 的媒体片段中。
- URI
    该值是一个带引号的字符串，其中包含标识媒体播放列表文件的 URI。此属性是可选的；
    请参阅 [第 4.3.4.2.1 节](#43421-alternative-renditions) 。如果类型为 `CLOSED-CAPTIONS` ，则 URI 属性 *不得* 存在。
- GROUP-ID
    该值是一个带引号的字符串，用于指定 Rendition 所属的组。
    请参阅 [第 4.3.4.1.1 节](#43411-rendition-groups) 。此属性是必需的。
- LANGUAGE
    该值是一个带引号的字符串，其中包含用于识别语言的标准标记之一 [[RFC5646][RFC5646]]，用于标识 Rendition 中使用的主要语言。此属性是可选的。
- ASSOC-LANGUAGE
    该值是一个带引号的字符串，其中包含一个语言标记 [[RFC5646][RFC5646]] ，用于标识与 Rendition 关联的语言。
    关联语言通常用于与 `LANGUAGE` 属性指定的语言不同的角色（例如，书面与口语或后备方言）。此属性是可选的。
    例如， `LANGUAGE` 和 `ASSOC-LANGUAGE` 属性可用于链接使用不同口语和书面语言的挪威语版本。
- NAME
    该值是一个带引号的字符串，其中包含 Rendition 的可读描述。
    如果存在 `LANGUAGE` 属性，则该描述应使用该语言。此属性是必需的。
- DEFAULT
    该值是一个枚举字符串；有效字符串为 `YES` 和 `NO` 。
    如果该值为 `YES` ，则客户端应在用户未提供表明其他选择的信息的情况下播放此内容版本。
    此属性是可选的。如果未提供该属性，则隐含值为 `NO` 。
- AUTOSELECT
    该值是一个枚举字符串；有效字符串为 `YES` 和 `NO` 。
    此属性是可选的。如果不存在该属性，则表示隐式值为 `NO` 。
    如果值为 `YES` ，则客户端可以选择在没有明确用户偏好的情况下播放此版本，因为它与当前播放环境（例如所选的系统语言）相匹配。
    如果存在 `AUTOSELECT` 属性，则如果 `DEFAULT` 属性的值为 `YES` ，则其值必须为 `YES` 。
- FORCED
    该值是一个枚举字符串；有效字符串为 `YES` 和 `NO` 。
    此属性是可选的。若不存在则表示隐式值为 `NO` 。
    除非 `TYPE` 为 `SUBTITLES` ，否则不得存在 `FORCED` 属性。
    值为 `YES` 表示版本包含被视为播放必不可少的内容。
    选择强制版本时，客户端应选择与当前播放环境（例如语言）最匹配的版本。
    值为 `NO` 表示 Rendition 包含旨在响应明确的用户请求而播放的内容。
- INSTREAM-ID
    该值是一个带引号的字符串，用于指定媒体播放列表中片段内的 Rendition。
    如果 `TYPE` 属性为 `CLOSED-CAPTIONS` ，则此属性是必需的，在这种情况下，它必须具有以下值之一： `CC1` 、 `CC2` 、 `CC3` 、 `CC4` 或 `SERVICEn` ，其中 n 必须是 1 到 63 之间的整数（例如 `SERVICE3` 或 `SERVICE42` ）。
    值 `CC1` 、 `CC2` 、 `CC3` 和 `CC4` 标识 Line 21 数据服务信道 [CEA608]。“SERVICE”值标识数字电视隐藏式字幕 [CEA708] 服务块编号。
    对于所有其他 `TYPE` 值，不得指定 `INSTREAM-ID` 。
- CHARACTERISTICS
    该值是一个带引号的字符串，其中包含一个或多个由逗号 `,` 字符分隔的统一类型标识符 [[UTI][UTI]]。
    此属性是可选的。每个 UTI 表示 Rendition 的单独特征。

    字幕版本可能包含以下特征：
    
  - `"public.accessibility.transcribes-spoken-dialog"` 、
  - `"public.accessibility.describes-music-and-sound"` 和
  - `"public.easy-to-read"` （表示字幕已经过编辑以便于阅读）。

    音频演绎可能包含以下特征： `"public.accessibility.describes-video"` 。
    `CHARACTERISTICS` 属性可以包括私有 UTI
- CHANNELS
    该值是一个带引号的字符串，用于指定有序的、以反斜杠分隔的参数列表（ `/` ）。
    如果 `TYPE` 属性为 `AUDIO` ，则第一个参数是以十进制整数表示的音频通道数，表示 Rendition 中任何媒体段中存在的独立、同时音频通道的最大数量。
    例如，AC-3 5.1 Rendition 将具有 `CHANNELS="6"` 属性。目前未定义其他 `CHANNELS` 参数。
    所有音频 `EXT-X-MEDIA` 标签都应具有 `CHANNELS` 属性。
    如果主播放列表包含两个使用相同编解码器但不同声道数编码的版本，则 `CHANNELS` 属性是必需的；否则，它是可选的。


###### 4.3.4.1.1. Rendition Groups

一组具有相同 `GROUP-ID` 值和相同 `TYPE` 值的一个或多个 `EXT-X-MEDIA` 标签定义了一个 Rendition 组。
组中的每个成员都必须是相同内容的替代 Rendition；否则，可能会出现播放错误。

播放列表中的所有 `EXT-X-MEDIA` 标签必须满足以下限制：

- 同一组内的所有 `EXT-X-MEDIA` 标签必须具有不同的 `NAME` 属性。
- 一个组不得有多个 `DEFAULT` 属性为 `YES` 的成员。
- 每个具有 `AUTOSELECT=YES` 属性的 `EXT-X-MEDIA` 标签都应具有 `LANGUAGE` [[RFC5646][RFC5646]] 、 `ASSOC-LANGUAGE` 、 `FORCED` 和 `CHARACTERISTICS` 属性的组合，这些属性与其组中其他 `AUTOSELECT=YES` 成员的属性不同。

播放列表可以包含多个相同类型的组，以便提供该媒体类型的多种编码。
如果这样做，则相同类型的每个组必须具有相同的成员集，并且每个相应的成员必须具有相同的属性（ `URI` 和 `CHANNELS` 属性除外）。

一个版本组中的每个成员可以有不同的样本格式。
例如，英语版本可以用 `AC-3 5.1` 编码，而西班牙语版本则用 `AAC` 立体声编码。
但是，引用此类组的任何 `EXT-X-STREAM-INF` 标签（ [第 4.3.4.2 节](#4342-ext-x-stream-inf) ）或 `EXT-X-I-FRAME-STREAM-INF` 标签（ [第 4.3.4.3 节](#4343-ext-x-i-frame-stream-inf) ）必须具有 `CODECS` 属性，该属性列出了组中任何版本中存在的每种样本格式，否则可能会发生客户端播放失败。
在上面的示例中， `CODECS` 属性将包括 `"ac-3,mp4a.40.2"` 。


##### 4.3.4.2. EXT-X-STREAM-INF

`EXT-X-STREAM-INF` 标签指定变体流，它是一组可以组合起来播放演示文稿的 Rendition。
该标签的属性提供了有关变体流的信息。

`EXT-X-STREAM-INF` 标签后面的 `URI` 行指定了包含变体流版本 (Rendition) 的媒体播放列表。
此 `URI` 行是必需的。不支持多个视频版本 (Rendition) 的客户端应该播放此版本。

它的格式为：

```plain
#EXT-X-STREAM-INF:<attribute-list>
<URI>
```

定义了以下属性:

- BANDWIDTH
    该值是 比特/秒 的十进制整数。它表示变体流的峰值片段比特率。
    如果变体流中的所有媒体片段都已创建，则 `BANDWIDTH` 值必须是任何可播放的版本组合产生的峰值片段比特率的最大总和。
    （对于具有单个媒体播放列表的变体流，这只是该媒体播放列表的峰值片段比特率。）
    不准确的值可能会导致播放停顿或阻止客户端播放变体。

    如果要在演示中的所有媒体片段都完成编码之前提供主播放列表，则 `BANDWIDTH` 值应该是使用相同设置进行编码的类似内容的代表时段的 `BANDWIDTH` 值。

    每个 `EXT-X-STREAM-INF` 标签 *必须* 包含 `BANDWIDTH` 属性。
- AVERAGE-BANDWIDTH
    该值是比特/秒的十进制整数。它表示变体流的平均分段比特率。

    如果变体流中的所有媒体片段都已创建，则 `AVERAGE-BANDWIDTH` 值必须是任何可播放的 Rendition 组合产生的平均片段比特率的最大总和。
    （对于具有单个媒体播放列表的变体流，这只是该媒体播放列表的平均片段比特率。）
    不准确的值可能会导致播放停顿或阻止客户端播放变体。

    如果要在演示文稿中的所有媒体片段都完成编码之前提供主播放列表，则平均带宽值应该是使用相同设置进行编码的类似内容代表时段的平均带宽值。

    `AVERAGE-BANDWIDTH` 属性为可选属性。
- CODECS
    该值是一个带引号的字符串，其中包含以逗号分隔的格式列表，其中每种格式指定变体流指定的一个或多个呈现形式中存在的媒体样本类型。
    有效格式标识符是 ISO 基本媒体文件格式名称空间中的标识符，由" 'Bucket' 媒体类型的 'Codecs' 和 'Profiles' 参数" [[RFC6381][RFC6381]] 定义。

    例如，包含 `AAC` 低复杂度 ( `AAC-LC` ) 音频和 `H.264 Main Profile Level 3.0` 视频的流的 `CODECS` 值为 `"mp4a.40.2,avc1.4d401e"` 。
    每个 `EXT-X-STREAM-INF` 标签应该包含一个 `CODECS` 属性。
- RESOLUTION
    该值是一个十进制分辨率，描述了显示变体流中所有视频的最佳像素分辨率。
    RESOLUTION 属性是可选的，但如果变体流包含视频，则建议使用。
- FRAME-RATE
    该值是一个十进制浮点数，描述变体流中所有视频的最大帧速率，四舍五入到小数点后三位。
    `FRAME-RATE` 属性是可选的，但如果变体流包含视频，则建议使用。
    如果变体流中的任何视频超过每秒 30 帧，则应包含 `FRAME-RATE` 属性。
- HDCP-LEVEL
    该值是一个枚举字符串；有效字符串为 `TYPE-0` 和 `NONE` 。
    此属性是建议性的；`TYPE-0` 值表示变体流可能无法播放，除非输出受高带宽数字内容保护 (HDCP) Type 0 [[HDCP][HDCP]] 或同等保护。
    `NONE` 值表示内容不需要输出版权保护。

    具有不同 `HDCP` 级别的加密变体流应使用不同的媒体加密密钥。

    `HDCP-LEVEL` 属性是可选的。
    如果变体流中的任何内容在没有 `HDCP` 的情况下无法播放，则应存在该属性。
    没有输出版权保护的客户端不应加载具有 `HDCP-LEVEL` 属性的变体流，除非其值为 NONE。
- AUDIO
    该值是一个带引号的字符串。它必须与主播放列表中其他位置的 `EXT-X-MEDIA` 标签的 `GROUP-ID` 属性值匹配，该标签的 `TYPE` 属性为 `AUDIO` 。
    它表示播放演示时应使用的音频呈现集。请参阅 [第 4.3.4.2.1 节](#43421-alternative-renditions) 。
    `AUDIO` 属性是可选的。
- VIDEO
    该值是一个带引号的字符串。
    它 *必须* 与主播放列表中其他位置的 `EXT-X-MEDIA` 标签的 `GROUP-ID` 属性值匹配，该标签的 `TYPE` 属性为 `VIDEO` 。
    它表示播放演示时应使用的视频呈现集。请参阅 [第 4.3.4.2.1 节](#43421-alternative-renditions) 。
    `VIDEO` 属性是可选的。
- SUBTITLES
    该值是一个带引号的字符串。
    它必须与主播放列表中其他位置的 `EXT-X-MEDIA` 标签的 `GROUP-ID` 属性值匹配，该标签的 `TYPE` 属性为 `SUBTITLES` 。
    它表示播放演示时可以使用的字幕版本集。请参阅 [第 4.3.4.2.1 节](#43421-alternative-renditions) 。
    `SUBTITLES` 属性是可选的。
- CLOSED-CAPTIONS
    该值可以是带引号的字符串，也可以是值为 `NONE` 的枚举字符串。
    如果该值是带引号的字符串，则它必须与播放列表中其他位置的 `EXT-X-MEDIA` 标签的 `GROUP-ID` 属性值匹配，该标签的 `TYPE` 属性为 `CLOSED-CAPTIONS` ，并且它指示播放演示文稿时可以使用的隐藏字幕版本集。请参阅 [第 4.3.4.2.1 节](#43421-alternative-renditions) 。
    如果该值为枚举字符串值 `NONE` ，则所有 `EXT-X-STREAM-INF` 标签都必须具有该属性，其值为 `NONE` ，表示主播放列表中的任何变体流中均无隐藏字幕。如果一个变体流中有隐藏字幕，而另一个没有，则可能会引发播放不一致。
    `CLOSED-CAPTIONS` 属性是可选的。


###### 4.3.4.2.1. Alternative Renditions

当 `EXT-X-STREAM-INF` 标签包含 `AUDIO` 、 `VIDEO` 、 `SUBTITLES` 或 `CLOSED-CAPTIONS` 属性时，表示可以使用内容的替代版本来播放该变体流。

定义替代版本时， *必须* 满足以下约束以防止客户端播放错误：

- 与 `EXT-X-STREAM-INF` 标签关联的所有可播放的 Rendition 组合的总带宽必须小于或等于 `EXT-X-STREAM-INF` 标签的 `BANDWIDTH` 属性。
- 如果 `EXT-X-STREAM-INF` 标签包含 `RESOLUTION` 属性和 `VIDEO` 属性，则每个备选视频再现都必须具有与 `RESOLUTION` 属性的值匹配的最佳显示分辨率。
- 与 `EXT-X-STREAM-INF` 标签相关的每个替代渲染都必须满足 [第 6.2.4 节](#624-providing-variant-streams) 中描述的变体流的约束。

如果媒体类型为 `SUBTITLES` ，则 `EXT-X-MEDIA` 标签的 `URI` 属性是必需的；如果媒体类型为 `VIDEO` 或 `AUDIO` ，则 `URI` 属性是可选的。
如果媒体类型为 `VIDEO` 或 `AUDIO` ，则缺少 `URI` 属性表示此 Rendition 的媒体数据包含在引用此 `EXT-X-MEDIA` 标签的任何 `EXT-X-STREAM-INF` 标签的媒体播放列表中。
如果媒体类型为 `AUDIO` ，且缺少 `URI` 属性，则客户端必须假定此 Rendition 的音频数据存在于 `EXT-X-STREAM-INF` 标签指定的每个视频 Rendition 中。

如果媒体类型为 `CLOSED-CAPTIONS` ，则不得包含 `EXT-X-MEDIA` 标签的 `URI` 属性。


##### 4.3.4.3. EXT-X-I-FRAME-STREAM-INF

`EXT-X-I-FRAME-STREAM-INF` 标签标识包含多媒体演示的 I 帧的媒体播放列表文件。
它是独立的，因为它不适用于主播放列表中的特定 URI。其格式为：

`#EXT-X-I-FRAME-STREAM-INF:<attribute-list>`

`EXT-X-STREAM-INF` 标签定义的所有属性（ [第 4.3.4.2 节](#4342-ext-x-stream-inf) ）也为 `EXT-X-I-FRAME-STREAM-INF` 标签定义，但 `FRAME-RATE` 、 `AUDIO` 、 `SUBTITLES` 和 `CLOSED-CAPTIONS` 属性除外。
此外，还定义了以下属性：

- URI
    该值是一个带引号的字符串，其中包含标识 I 帧媒体播放列表文件的 URI。
    该播放列表文件必须包含 `EXT-X-I-FRAMES-ONLY` 标签。

每个 `EXT-X-I-FRAME-STREAM-INF` 标签必须包含一个 `BANDWIDTH` 属性和一个 URI 属性。

[4.3.4.2.1节](#43421-alternative-renditions) 中的规定也适用于具有 `VIDEO` 属性的 `EXT-X-I-FRAME-STREAM-INF` 标签。

指定替代视频演绎和 I 帧播放列表的主播放列表应为每个常规视频演绎包含一个替代 I 帧视频演绎，并具有相同的 `NAME` 和 `LANGUAGE` 属性。


##### 4.3.4.4. EXT-X-SESSION-DATA

`EXT-X-SESSION-DATA` 标签允许在主播放列表中携带任意会话数据。

它的格式为：

`#EXT-X-SESSION-DATA:<attribute-list>`

定义了以下属性：

- DATA-ID
    `DATA-ID` 的值是一个带引号的字符串，用于标识特定的数据值。
    `DATA-ID` 应符合反向 DNS 命名约定，例如 `"com.example.movi​​e.title"` ； 
    但是，没有中央注册机构，因此播放列表作者 *应该* 小心选择一个不太可能与其他值冲突的值。
    此属性是必需的。
- VALUE
    `VALUE` 是一个带引号的字符串。它包含由 `DATA-ID` 标识的数据。
    如果指定了 `LANGUAGE` ， `VALUE` 应该包含以指定语言编写的人类可读字符串。
- URI
    该值是包含 `URI` 的带引号的字符串。
    `URI` 标识的资源必须采用 `JSON` 格式 [[RFC7159][RFC7159]] ；否则，客户端可能无法解释该资源。
- LANGUAGE
    该值是一个带引号的字符串，其中包含用于标识 `VALUE` 的语言的语言标记 [[RFC5646][RFC5646]] 。
    此属性是可选的。

每个 `EXT-X-SESSION-DATA` 标签必须包含 `VALUE` 或 `URI` 属性，但不能同时包含两者。

播放列表 *可以* 包含多个具有相同 `DATA-ID` 属性的 `EXT-X-SESSION-DATA` 标签。
播放列表 *不得* 包含多个具有相同 `DATA-ID` 属性和 `LANGUAGE` 属性的 `EXT-X-SESSION-DATA` 标签。


##### 4.3.4.5. EXT-X-SESSION-KEY

`EXT-X-SESSION-KEY` 标签允许在主播放列表中指定媒体播放列表中的加密密钥。
这样客户端就可以预加载这些密钥，而无需先读取媒体播放列表。

它的格式为：

`#EXT-X-SESSION-KEY:<attribute-list>` 

`EXT-X-KEY` 标签定义的所有属性（ [第 4.3.2.4 节](#4324-ext-x-key) ）也为 `EXT-X-SESSION-KEY` 定义，但 `METHOD` 属性的值 *不得* 为 `NONE` 。
如果使用 `EXT-X-SESSION-KEY` ，则 `METHOD` 、 `KEYFORMAT` 和 `KEYFORMATVERSIONS` 属性的值必须与任何具有相同 URI 值的 `EXT-X-KEY` 匹配。

如果多个变体流或版本使用相同的加密密钥和格式，则 *应该* 添加 `EXT-X-SESSION-KEY` 标签。
`EXT-X-SESSION-KEY` 标签不与任何特定媒体播放列表相关联。

主播放列表 *不得* 包含多个具有相同 `METHOD` 、 `URI` 、 `IV` 、 `KEYFORMAT` 和 `KEYFORMATVERSIONS` 属性值的 `EXT-X-SESSION-KEY` 标签。

`EXT-X-SESSION-KEY` 标签是可选的。


#### 4.3.5. Media or Master Playlist Tags

本节中的标签可以出现在主播放列表或媒体播放列表中。
如果其中一个标签出现在主播放列表中，则它 *不应该* 出现在该主播放列表引用的任何媒体播放列表中。
同时出现在两者中的标签 *必须* 具有相同的值；否则，客户端 *应该* 忽略媒体播放列表中的值。

这些标签 *不得* 在播放列表中出现多次。
如果标签出现多次，客户端 *必须* 无法解析播放列表。


##### 4.3.5.1. EXT-X-INDEPENDENT-SEGMENTS

`EXT-X-INDEPENDENT-SEGMENTS` 标签表示媒体片段中的所有媒体样本都可以解码，而无需来自其他片段的信息。
它适用于播放列表中的每个媒体片段。

其格式为：

`#EXT-X-INDEPENDENT-SEGMENTS`

如果 `EXT-X-INDEPENDENT-SEGMENTS` 标签出现在主播放列表中，则它将适用于主播放列表中每个媒体播放列表的每个媒体片段。


##### 4.3.5.2. EXT-X-START

`EXT-X-START` 标签表示开始播放播放列表的首选点。
默认情况下，客户端在开始播放会话时 *应该* 在此点开始播放。
此标签是可选的。

其格式为：

`#EXT-X-START:<attribute-list>`

定义了以下属性：

- TIME-OFFSET
    `TIME-OFFSET` 的值是一个有符号的十进制浮点数，单位为秒。
    正数表示与播放列表开头之间的时间偏移。
    负数表示与播放列表中最后一个媒体片段结尾之间的负时间偏移。
    此属性是必需的。

    `TIME-OFFSET` 的绝对值不应大于播放列表的持续时间。
    如果 `TIME-OFFSET` 的绝对值超过播放列表的持续时间，则表示播放列表的结束（如果为正数）或播放列表的开始（如果为负数）。

    如果播放列表不包含 `EXT-X-ENDLIST` 标签，则 `TIME-OFFSET` *不应该* 位于播放列表文件结尾的三个目标持续时间内。
- PRECISE
    该值是一个枚举字符串；有效字符串为 `YES` 和 `NO` 。
    如果该值为 `YES` ，客户端应从包含 `TIME-OFFSET` 的媒体段开始播放，但不应渲染该段中呈现时间早于 `TIME-OFFSET` 的媒体样本。
    如果该值为 `NO` ，客户端应尝试渲染该段中的每个媒体样本。
    此属性是可选的。
    如果缺失，其值应被视为 `NO` 。


## 5. Key Files

### 5.1. Structure of Key Files

具有 `URI` 属性的 `EXT-X-KEY` 标签标识密钥文件。
密钥文件包含一个可以解密播放列表中的媒体片段的密钥。

[[AES_128](#52-iv-for-aes-128)] 加密使用 16 个八位字节的密钥。
如果 `EXT-X-KEY` 标签的 `KEYFORMAT` 为 `"identity"` ，则密钥文件是二进制格式的单个 16 个八位字节的打包数组。


### 5.2. IV for AES-128

[AES_128] *需要* 在加密和解密时提供相同的 16 个八位字节 IV 。
改变此 IV 可增加密码的强度。

`EXT-X-KEY` 标签上的 IV 属性（ `KEYFORMAT` 为 `"identity"` ）指定了在解密使用该密钥文件加密的媒体片段时可以使用的 IV 。
AES-128 的 IV 值为 128 位数字。

`KEYFORMAT` 为 `identity` 且没有 IV 属性的 `EXT-X-KEY` 标签表示在解密媒体段时将使用媒体序列号作为 IV，通过将其大端二进制表示放入 16 八位字节（128 位）缓冲区并用零填充（在左侧）。


## 6. Client/Server Responsibilities

### 6.1. Introduction

本节介绍服务器如何生成播放列表和媒体片段以及客户端如何下载它们进行播放。


### 6.2. Server Responsibilities

#### 6.2.1. General Server Responsibilities

源媒体的制作超出了本文档的范围，本文档仅假设包含演示的连续编码媒体源。

服务器 *必须* 将源媒体划分为单个媒体片段，这些片段的持续时间小于或等于固定目标持续时间。
超过计划目标持续时间的片段可能会触发播放停顿和其他错误。

服务器 *应该* 尝试在支持有效解码各个媒体段的点（例如，在数据包和关键帧边界）上划分源媒体。

服务器 *必须* 为每个媒体片段创建一个 URI ，以便其客户端获取片段数据。
如果服务器支持部分加载资源（例如，通过 HTTP 范围请求），则可以使用 `EXT-X-BYTERANGE` 标签将片段指定为较大资源的子范围。

客户端加载的播放列表中指定的任何媒体片段都 *必须* 可立即下载，否则可能会出现播放错误。
下载开始后，其传输速率 *不应该* 受到片段制作过程的限制。

如果客户端表示准备接受，HTTP 服务器 *应该* 使用 `gzip` 内容编码传输文本文件（例如播放列表和 `WebVTT` 片段）。

服务器必须创建一个媒体播放列表文件（ [第 4 节](#4-playlists) ），其中包含服务器希望提供的每个媒体片段的 URI ，并按照播放顺序排列。

`EXT-X-VERSION` 标签的值（ [第 4.3.1.2 节](#4312-ext-x-version) ）不应大于播放列表中标签和属性所需的值（参见 [第 7 节](#7-protocol-version-compatibility) ）。

从客户端的角度来看，对播放列表文件的更改 *必须* 以原子方式进行，否则可能会出现播放错误。

服务器 *不得* 更改媒体播放列表文件，除非：

- 向其添加行（ [第 6.2.1 节](#621-general-server-responsibilities) ）。
- 按照出现的顺序从播放列表中删除媒体片段 URI ，以及仅适用于这些片段的任何标签（ [第 6.2.2 节](#622-live-playlists) ）。
- 增加 `EXT-X-MEDIA-SEQUENCE` 或 `EXT-X-DISCONTINUITY-SEQUENCE` 标签的值（ [第 6.2.2 节](#622-live-playlists) ）。
- 将 `EXT-X-ENDLIST` 标签添加到播放列表（ [第 6.2.1 节](#621-general-server-responsibilities) ）。

如果媒体播放列表包含 `EXT-X-PLAYLIST-TYPE` 标签，则其更新会受到进一步限制。
值为 `VOD` 的 `EXT-X-PLAYLIST-TYPE` 标签表示播放列表文件不得更改。
值为 `EVENT` 的 `EXT-X-PLAYLIST-TYPE` 标签表示服务器 *不得* 更改或删除播放列表文件的任何部分；它可以向其中添加行。

媒体播放列表中的 `EXT-X-TARGETDURATION` 标签的值 *不得* 更改。
典型的目标持续时间为 10 秒。

除此处允许的更改之外的播放列表更改可能会触发播放错误和不一致的客户端行为。

媒体播放列表中的每个媒体片段都有一个整数 **不连续序列号** 。
除了媒体内的时间戳之外，还可以使用 **不连续序列号** 来同步不同版本之间的媒体片段。

片段的不连续性序列号是 `EXT-X-DISCONTINUITY-SEQUENCE` 标签的值（如果没有则为零）加上该片段 URI 行之前的播放列表中 `EXT-X-DISCONTINUITY` 标签的数量。

服务器 *可以* 通过将 `EXT-X-PROGRAM-DATE-TIME` 标签应用于媒体片段，从而将绝对日期和时间与媒体片段关联起来。
这定义了标签指定的（挂钟）日期和时间与片段中的第一个媒体时间戳之间的信息映射，可用作查找、显示或其他用途的基础。
如果服务器提供此映射，则 *应该* 将 `EXT-X-PROGRAM-DATE-TIME` 标签应用于每个已应用 `EXT-X-DISCONTINUITY` 标签的片段。

服务器不得向播放列表添加任何 `EXT-X-PROGRAM-DATE-TIME` 标签，这会导致节目日期和媒体片段之间的映射变得模糊。

如果范围内的任何日期映射到播放列表中的媒体段，则服务器 *不得* 从播放列表中删除 `EXT-X-DATERANGE` 标签。

服务器 *不得* 在同一播放列表中的任何新日期范围重复使用 `EXT-X-DATERANGE` 标签的 ID 属性值。

一旦将具有 `END-ON-NEXT=YES` 属性的日期范围的以下范围添加到播放列表，服务器 *不得* 随后添加具有相同 CLASS 属性的日期范围，其 `START-DATE` 位于 `END-ON-NEXT=YES` 范围和其以下范围之间。

对于具有 `PLANNED-DURATION` 属性的日期范围，服务器 *应该* 在范围建立后发出实际结束信号。
它可以通过添加另一个具有相同 `ID` 属性值和 `DURATION` 或 `END-DATE` 属性的 `EXT-X-DATERANGE` 标签来实现这一点，或者，如果日期范围具有 `END-ON-NEXT=YES` 属性，则通过添加以下范围来实现。

如果媒体播放列表包含演示的最终媒体片段，则播放列表文件 *必须* 包含 `EXT-X-ENDLIST` 标签；这允许客户端最大限度地减少无效的播放列表重新加载。

如果媒体播放列表不包含 `EXT-X-ENDLIST` 标签，则服务器 *必须* 提供包含至少一个新媒体片段的新版播放列表文件。
新版播放列表文件的可用时间 *必须* 与上一版播放列表文件的可用时间相对应：不早于该时间之后目标持续时间的一半，不晚于该时间之后目标持续时间的 1.5 倍。这样，客户端就可以高效利用网络。

如果服务器希望删除整个演示文稿，它 *应该* 向客户端明确指示播放列表文件不再可用（例如，使用 HTTP 404 或 410 响应）。
它 *必须* 确保播放列表文件中的所有媒体片段至少在删除时播放列表文件的持续时间内对客户端可用，以防止中断正在进行的播放。


#### 6.2.2. Live Playlists

服务器 *可以* 通过从播放列表文件中删除媒体片段来限制媒体片段的可用性（ [第 6.2.1 节](#621-general-server-responsibilities) ）。
如果要删除媒体片段，播放列表文件 *必须* 包含 `EXT-X-MEDIA-SEQUENCE` 标签。
每从播放列表文件中删除一个媒体片段，其值 *必须* 增加 1； *不得* 减少或换行。
如果每个媒体片段没有一致、唯一的媒体序列号，客户端可能会出现故障。

*必须* 按照媒体片段在播放列表中出现的顺序将其从播放列表文件中删除；否则，客户端播放可能会出现故障。

如果删除媒体片段会导致播放列表的时长少于目标时长的三倍，则服务器 *不得* 从不带 `EXT-X-ENDLIST` 标签的播放列表文件中删除媒体片段。
否则，可能会触发播放停顿。

当服务器从播放列表中移除媒体片段 URI 时，相应的媒体片段 *必须* 保持可供客户端使用的一段时间，该时间等于片段的持续时间加上由包含该片段的服务器分发的最长播放列表文件的持续时间。
在此之前移除媒体片段可能会中断正在进行的播放。

如果服务器希望从包含 `EXT-X-DISCONTINUITY` 标签的媒体播放列表中删除片段，则媒体播放列表 *必须* 包含 `EXT-X-DISCONTINUITY-SEQUENCE` 标签。
如果没有 `EXT-X-DISCONTINUITY-SEQUENCE` 标签，客户端就无法在 Rendition 之间找到相应的片段。

如果服务器从媒体播放列表中删除了 `EXT-X-DISCONTINUITY` 标签，则 *必须* 增加 `EXT-X-DISCONTINUITY-SEQUENCE` 标签的值，以使仍在媒体播放列表中的片段的不连续序列号保持不变。
`EXT-X-DISCONTINUITY-SEQUENCE` 标签的值 *不得* 减少或换行。
如果每个媒体片段没有一致的 **不连续序列号** ，客户端可能会出现故障。

如果服务器计划在通过 HTTP 将媒体片段传送给客户端后删除它，则它 *应该* 确保 HTTP 响应包含反映计划生存时间的 `Expires` 标头。

直播播放列表 *不得* 包含 `EXT-X-PLAYLIST-TYPE` 标签，因为该标签的任何值都不允许删除媒体片段。


#### 6.2.3. Encrypting Media Segments

媒体片段 *可以* 加密。
每个加密的媒体片段必须应用一个 `EXT-X-KEY` 标签（ [第 4.3.2.4 节](#4324-ext-x-key) ），并带有一个 URI，客户端可以使用该 URI 获取包含解密密钥的密钥文件（ [第 5 节](#5-key-files) ）。

媒体片段只能使用一种加密 `METHOD` 、一个加密密钥和 IV 进行加密。
但是，服务器 *可以* 通过提供多个 `EXT-X-KEY` 标签（每个标签具有不同的 `KEYFORMAT` 属性值）来提供多种检索该密钥的方法。

服务器 *可以* 在密钥响应中设置 HTTP Expires 标头，以指示密钥可以缓存的持续时间。

播放列表中任何未加密的媒体片段（如果前面有一个加密的媒体片段）都 *必须* 应用 `EXT-X-KEY` 标签，且 `METHOD` 属性为 `NONE` 。
否则，客户端会误认为这些片段已加密。

如果加密 `METHOD` 是 `AES-128` ，且播放列表不包含 `EXT-X-I-FRAMES-ONLY` 标签，则 *应* 将 [第 4.3.2.4 节](#4324-ext-x-key) 中所述的 AES 加密应用于各个媒体片段。

如果加密 `METHOD` 是 `AES-128` ，且播放列表包含 `EXT-X-I-FRAMES-ONLY` 标签，则必须使用带 PKCS7 填充 [RFC5652][RFC5652] 的 AES-128 CBC 对整个资源进行加密。
除非第一个块包含 I 帧，否则 *可以* 在 16 字节块边界上重新开始加密。
用于加密的 IV *必须* 是媒体段的 **媒体序列号** 或 `EXT-X-KEY` 标签的 IV 属性值，如 [第 5.2 节](#52-iv-for-aes-128) 所述。
这些限制允许客户端加载和解密指定为常规加密媒体段子范围及其媒体初始化部分的各个 I 帧。

如果加密 `METHOD` 是 `SAMPLE-AES` ，则媒体样本 *可以* 在封装到媒体段之前进行加密。

如果 `EXT-X-KEY` 标签适用于播放列表文件中的任何媒体段，则服务器 *不得* 从播放列表文件中删除该标签，否则随后加载该播放列表的客户端将无法解密这些媒体段。


#### 6.2.4. Providing Variant Streams

服务器 *可以* 提供多个媒体播放列表文件，以提供同一演示文稿的不同编码。
如果这样做，它 *应该* 提供一个列出每个变体流的主播放列表文件，以允许客户端动态切换编码。

主播放列表使用 `EXT-X-STREAM-INF` 标签描述常规变体流，使用 `EXT-X-I-FRAME-STREAM-INF` 标签描述 I 帧变体流。

如果 `EXT-X-STREAM-INF` 标签或 `EXT-X-I-FRAME-STREAM-INF` 标签包含 `CODECS` 属性，则属性值 *必须* 包含变体流指定的任何版本中任何媒体段中存在的每种媒体格式 [RFC6381][RFC6381] 。

为了允许客户端在它们之间无缝切换，服务器在生成变体流时 *必须* 满足以下约束：

- 每个变体流 *必须* 呈现相同的内容。
- 变体流中的匹配内容 *必须* 具有匹配的时间戳。这允许客户端同步媒体。
- 变体流中的匹配内容 *必须* 具有匹配的 **不连续序列号** （参见 [第 4.3.3.3 节](#4333-ext-x-discontinuity-sequence) ）。
- 每个变体流中的每个媒体播放列表都 *必须* 具有相同的目标时长。
    唯一的例外是字幕版本和包含 `EXT-X-I-FRAMES-ONLY` 标签的媒体播放列表，如果它们的 `EXT-X-PLAYLIST-TYPE` 为 `VOD` ，则它们 *可以* 具有不同的目标时长。
- 在一个变体流的媒体播放列表中出现但未在另一个变体流中出现的内容 *必须* 出现在媒体播放列表文件的开头或结尾，并且不得长于目标持续时间。
- 如果任何媒体播放列表具有 `EXT-X-PLAYLIST-TYPE` 标签，则所有媒体播放列表都 *必须* 具有相同值的 `EXT-X-PLAYLIST-TYPE` 标签。
- 如果播放列表包含值为 `VOD` 的 `EXT-X-PLAYLIST-TYPE` 标签，则每个变体流中每个媒体播放列表的第一个片段 *必须* 从相同的媒体时间戳开始。
- 如果主播放列表中的任何媒体播放列表包含 `EXT-X-PROGRAM-DATE-TIME` 标签，则该主播放列表中的所有媒体播放列表都必须包含 `EXT-X-PROGRAM-DATE-TIME` 标签，并且日期和时间与媒体时间戳具有一致的映射。
- 每个变体流 *必须* 包含相同的一组日期范围，每个日期范围由具有相同 ID 属性值的 `EXT-X-DATERANGE` 标签标识，并包含相同的一组属性/值对。

此外，为了实现最广泛的兼容性，变体流 *应该* 包含相同的编码音频比特流。
这样，客户端就可以在变体流之间切换，而不会出现声音故障。

变体流的规则也适用于替代演绎版 （参见 [第 4.3.4.2.1 节](#43421-alternative-renditions) ）。


### 6.3. Client Responsibilities

#### 6.3.1. General Client Responsibilities 

客户端如何获取播放列表文件的 URI 超出了本文档的范围；假定已经这样做了。

客户端从 URI 中获取播放列表文件，如果获取到的播放列表文件是主播放列表，则客户端可以从主播放列表中选择一个变体流进行加载。

客户端 *必须* 确保加载的播放列表符合 [第 4 节](#4-playlists) ，并且 `EXT-X-VERSION` 标签（如果存在）指定客户端支持的协议版本；
如果任一检查失败，客户端 *不得* 尝试使用播放列表，否则可能会发生意外行为。

如果播放列表中的任何 URI 元素包含客户端无法处理的 URI 方案，则客户端 *必须* 停止播放。
所有客户端都 *必须* 支持 HTTP 方案。

为了支持前向兼容性，在解析播放列表时，客户端必须：

- 忽略任何无法识别的标签。
- 忽略任何具有无法识别的 `AttributeName` 的属性/值对。
- 忽略任何包含枚举字符串类型属性/值对的标签，该标签的 `AttributeName` 可以被识别，但 `AttributeValue` 无法被识别，除非该属性的定义另有说明。

客户端用于在变体流之间切换的算法超出了本文档的范围。


#### 6.3.2. Loading the Media Playlist File

每次从播放列表 URI 加载或重新加载媒体播放列表时，如果客户端打算正常播放演示（即按播放列表顺序以标称播放速率播放），则 *必须* 确定下一个要加载的媒体片段，如 [第 6.3.5 节](#635-determining-the-next-segment-to-load) 所述。

如果媒体播放列表包含 `EXT-X-MEDIA-SEQUENCE` 标签，则客户端 *应该* 假定其中的每个媒体片段在播放列表文件加载时以及播放列表文件的持续时间内将变得不可用。

重新加载播放列表时，客户端 *可以* 使用片段 **媒体序列号** 来跟踪播放列表中媒体片段的位置。

客户端 *不得* 假设不同变体流或演绎版中具有相同 **媒体序列号** 的片段在演示中具有相同的位置；播放列表 *可以* 具有独立的媒体序列号。
相反，客户端 *必须* 使用播放列表时间线上每个片段的相对位置及其 **不连续序列号** 来定位相应的片段。

客户端 *必须* 加载每个选定要播放的版本媒体播放列表文件，以便找到特定于该版本媒体的位置。
但为了防止服务器产生不必要的负载，客户端 *不应该* 加载任何其他版本的播放列表文件。

对于某些变体流，可以选择不包含 `EXT-X-STREAM-INF` 标签指定的 Rendition 的 Rendition。
如上所述，在这种情况下，客户端 *不应该* 加载该 Rendition。


#### 6.3.3. Playing the Media Playlist File

播放开始时，客户端 *应* 从媒体播放列表中选择首先播放哪个媒体片段。
如果不存在 `EXT-X-ENDLIST` 标签，并且客户端打算正常播放媒体，则客户端 *不应该* 选择从播放列表文件末尾开始的片段，该片段的起始位置少于三个目标持续时间。
这样做可能会触发播放停顿。

按照播放列表中出现的顺序播放媒体片段即可实现正常播放。
客户端 *可以* 按照自己希望的任何方式呈现可用媒体，包括正常播放、随机访问和特技模式。

媒体片段中以及媒体播放列表中多个媒体片段中的样本的编码参数 *应该* 保持一致。
但是，客户端应在遇到编码变化时进行处理，例如，通过缩放视频内容来适应分辨率变化。
如果变体流包含 RESOLUTION 属性，则客户端 *应该* 将所有视频显示在与该分辨率具有相同比例的矩形内。

客户端 *应该* 准备好处理特定类型（例如音频或视频）的多个轨道。
没有其他偏好的客户端 *应该* 选择其可以播放的具有最小数字轨道标识符的轨道。

客户端 *应该* 忽略传输流中它们无法识别的私有流。
私有流可用于支持使用相同流的不同设备，但流作者 *应该* 注意由此带来的额外网络负载。

在播放应用了 `EXT-X-DISCONTINUITY` 标签的媒体片段之前，客户端 *必须* 准备重置其解析器和解码器；否则，可能会发生播放错误。

客户端 *应该* 尝试在需要不间断播放之前提前加载媒体片段，以补偿延迟和吞吐量的暂时变化。

客户端 *可以* 使用 `EXT-X-PROGRAM-DATE-TIME` 标签的值向用户显示节目的起始时间。
如果该值包含时区信息，客户端 *应* 将其考虑在内；如果不包含，客户端 *可以* 假设时间为本地时间。

请注意，播放列表中的日期可以指内容的制作时间（或其他时间），与播放时间无关。

如果播放列表中的第一个 `EXT-X-PROGRAM-DATE-TIME` 标签出现在一个或多个媒体片段 URI 之后，客户端 *应该* 从该标签向后推断（使用 `EXTINF` 持续时间和/或媒体时间戳）以将日期与这些片段关联。
要将日期与任何其他未直接应用 `EXT-X-PROGRAM-DATE-TIME` 标签的媒体片段关联，客户端 *应该* 从播放列表中该片段之前的最后一个 `EXT-X-PROGRAM-DATE-TIME` 标签向前推断。


#### 6.3.4. Reloading the Media Playlist File

客户端 *必须* 定期重新加载媒体播放列表文件以了解当前可用的媒体，除非它包含值为 `VOD` 的 `EXT-X-PLAYLIST-TYPE` 标签，或者值为 `EVENT` 并且 `EXT-X-ENDLIST` 标签也存在。

但是，为了限制服务器上的集体负载，客户端 *不得* 尝试以比本节指定的更频繁的频率重新加载播放列表文件。

当客户端首次加载播放列表文件或重新加载播放列表文件并发现自上次加载以来已发生变化时，客户端 *必须* 至少等待目标时长，才能尝试再次重新加载播放列表文件，该时间从客户端上次开始加载播放列表文件的时间开始计算。

如果客户端重新加载播放列表文件并发现它没有改变，那么它 *必须* 等待目标持续时间的一半才能重试。

重新加载媒体播放列表后，客户端 *应该* 验证其中的每个媒体片段是否与上一个媒体播放列表中具有相同媒体序列号的媒​​体​​片段具有相同的 URI（和字节范围，如果指定）。
如果没有，它 *应该* 停止播放，因为这通常表示服务器错误。

为了减少服务器负载，客户端 *不应该* 重新加载当前未播放的变体流或替代版本播放列表文件。
如果客户端决定切换到其他变体流播放，则 *应该* 停止重新加载旧变体流播放列表，并开始加载新变体流播放列表。
客户端可以使用 `EXTINF` 持续时间和 [第 6.2.4 节](#624-providing-variant-streams) 中的约束来确定相应媒体的大致位置。
加载新变体流中的媒体后，媒体片段中的时间戳可用于精确同步新旧时间线。

客户端 *不得* 尝试使用媒体序列号在流之间进行同步（参见 [第 6.3.2 节](#632-loading-the-media-playlist-file) ）。


#### 6.3.5. Determining the Next Segment to Load

客户端 *必须* 在每次加载或重新加载时检查媒体播放列表文件以确定下一个要加载的媒体段，因为可用媒体集 *可能* 会发生变化。

第一个加载的片段通常是客户端选择首先播放的片段（参见 [第 6.3.3 节](#633-playing-the-media-playlist-file) ）。

为了正常播放演示文稿，下一个要加载的媒体段是具有最低媒体序列号的媒​​体​​段，该媒体序列号大于最后一个加载的媒体段的媒体序列号。


#### 6.3.6. Decrypting Encrypted Media Segments

如果媒体播放列表文件包含指定密钥文件 URI 的 `EXT-X-KEY` 标签，则客户端可以获取该密钥文件并使用其中的密钥解密该 `EXT-X-KEY` 标签适用的所有媒体段。

客户端 *必须* 忽略任何具有不受支持或无法识别的 `KEYFORMAT` 属性的 `EXT-X-KEY` 标签，以实现跨设备寻址。
如果播放列表包含仅应用了具有无法识别或无法支持的 `KEYFORMAT` 属性的 `EXT-X-KEY` 标签的媒体片段，则播放应会失败。

客户端 *不得* 尝试解密任何 `EXT-X-KEY` 标签具有其无法识别的 `METHOD` 属性的段。

如果加密 `METHOD` 是 `AES-128` ，则 *应该* 将 AES-128 CBC 解密应用于各个媒体段，其加密格式如 [第 4.3.2.4 节](#4324-ext-x-key) 所述。

如果加密方法是 `AES-128` ，且媒体片段是 I 帧播放列表的一部分（ [第 4.3.3.6 节](#4336-ext-x-i-frames-only) ），并且它应用了 `EXT-X-BYTERANGE` 标签，则在加载和解密该片段时需要特别小心，因为 URI 标识的资源从资源开始就以 16 字节块的形式加密。

可以通过首先扩大其字节范围（由 `EXT-X-BYTERANGE` 标签指定）来恢复解密的 I 帧，以便它从资源的开头开始和结束在 16 字节边界上。

接下来，字节范围进一步扩大，在范围开头包含一个 16 字节块。这个 16 字节块允许计算下一个块的正确 IV。

然后，可以使用任意 IV 加载扩展的字节范围并使用 AES-128 CBC 解密。解密的字节中丢弃添加到原始字节范围开头和结尾的字节数；剩下的就是解密的 I 帧。

如果加密 `METHOD` 是 `SAMPLE-AES` ，则 *应* 将 AES-128 解密应用于媒体段内的加密媒体样本。

`EXT-X-KEY` 标签的 `METHOD` 为 `NONE` ，表示其适用的媒体片段未加密。


## 7. Protocol Version Compatibility

协议兼容性由 `EXT-X-VERSION` 标签指定。
包含与协议版本 1 不兼容的标签或属性的播放列表必须包含 `EXT-X-VERSION` 标签。

如果客户端不支持 `EXT-X-VERSION` 标签指定的协议版本，则客户端 *不得* 尝试播放，否则可能会发生意外行为。

如果媒体播放列表包含以下内容，则必须指示 `EXT-X-VERSION` 为 2 或更高版本：

- `EXT-X-KEY` 标签的 IV 属性。
    
如果媒体播放列表包含以下内容，则必须指示 `EXT-X-VERSION` 为 3 或更高版本：

- 浮点 `EXTINF` 持续时间值。

如果媒体播放列表包含以下内容，则必须指示 `EXT-X-VERSION` 为 4 或更高版本：

- `EXT-X-BYTERANGE` 标签。
- `EXT-X-I-FRAMES-ONLY` 标签。

如果媒体播放列表包含以下内容，则必须指示 `EXT-X-VERSION` 为 5 或更高版本：

- `EXT-X-KEY` 标签的 `KEYFORMAT` 和 `KEYFORMATVERSIONS` 属性。
- `EXT-X-MAP` 标签

如果媒体播放列表包含以下内容，则必须指示 `EXT-X-VERSION` 为 6 或更高版本：

- 媒体播放列表中不包含 `EXT-X-I-FRAMES-ONLY` 的 `EXT-X-MAP` 标签。

如果媒体播放列表包含以下内容，则必须指示 `EXT-X-VERSION` 为 7 或更高版本：

- `EXT-X-MEDIA` 标签的 `INSTREAM-ID` 属性的 `"SERVICE"` 值。

`EXT-X-MEDIA` 标签以及 `EXT-X-STREAM-INF` 标签的 `AUDIO` 、 `VIDEO` 和 `SUBTITLES` 属性向后兼容协议版本 1，但在较旧的客户端上播放可能并不理想。
服务器可以考虑在主播放列表中指示 `EXT-X-VERSION` 为 4 或更高版本，但这不是必须的。

协议版本 6 中删除了 `EXT-X-STREAM-INF` 和 `EXT-X-I-FRAME-STREAM-INF` 标签的 `PROGRAM-ID` 属性。

`EXT-X-ALLOW-CACHE` 标签在协议版本 7 中已被删除。


## 8. Playlist Examples

### 8.1. Simple Media Playlist 

```m3u
#EXTM3U
#EXT-X-TARGETDURATION:10
#EXT-X-VERSION:3
#EXTINF:9.009,
http://media.example.com/first.ts
#EXTINF:9.009,
http://media.example.com/second.ts
#EXTINF:3.003,
http://media.example.com/third.ts
#EXT-X-ENDLIST
```


### 8.2. Live Media Playlist Using HTTPS 

```m3u
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:8
#EXT-X-MEDIA-SEQUENCE:2680

#EXTINF:7.975,
https://priv.example.com/fileSequence2680.ts
#EXTINF:7.941,
https://priv.example.com/fileSequence2681.ts
#EXTINF:7.975,
https://priv.example.com/fileSequence2682.ts
```


### 8.3. Playlist With Encrypted Media Segments

```m3u
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-MEDIA-SEQUENCE:7794
#EXT-X-TARGETDURATION:15

#EXT-X-KEY:METHOD=AES-128,URI="https://priv.example.com/key.php?r=52"

#EXTINF:2.833,
http://media.example.com/fileSequence52-A.ts
#EXTINF:15.0,
http://media.example.com/fileSequence52-B.ts
#EXTINF:13.333,
http://media.example.com/fileSequence52-C.ts

#EXT-X-KEY:METHOD=AES-128,URI="https://priv.example.com/key.php?r=53"

#EXTINF:15.0,
http://media.example.com/fileSequence53-A.ts
```


### 8.4. Master Playlist

```m3u
#EXTM3U
#EXT-X-STREAM-INF:BANDWIDTH=1280000,AVERAGE-BANDWIDTH=1000000
http://example.com/low.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=2560000,AVERAGE-BANDWIDTH=2000000
http://example.com/mid.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=7680000,AVERAGE-BANDWIDTH=6000000
http://example.com/hi.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=65000,CODECS="mp4a.40.5"
http://example.com/audio-only.m3u8
```


### 8.5. Master Playlist with I-Frames

```m3u
#EXTM3U
#EXT-X-STREAM-INF:BANDWIDTH=1280000
low/audio-video.m3u8
#EXT-X-I-FRAME-STREAM-INF:BANDWIDTH=86000,URI="low/iframe.m3u8"
#EXT-X-STREAM-INF:BANDWIDTH=2560000
mid/audio-video.m3u8
#EXT-X-I-FRAME-STREAM-INF:BANDWIDTH=150000,URI="mid/iframe.m3u8"
#EXT-X-STREAM-INF:BANDWIDTH=7680000
hi/audio-video.m3u8
#EXT-X-I-FRAME-STREAM-INF:BANDWIDTH=550000,URI="hi/iframe.m3u8"
#EXT-X-STREAM-INF:BANDWIDTH=65000,CODECS="mp4a.40.5"
audio-only.m3u8
```


### 8.6. Master Playlist with Alternative Audio

在此示例中， `CODECS` 属性已压缩以节省空间。
使用 `\` 表示标签在删除空格后继续到下一行：

```m3u
#EXTM3U
#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="aac",NAME="English", \
    DEFAULT=YES,AUTOSELECT=YES,LANGUAGE="en", \
    URI="main/english-audio.m3u8"
#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="aac",NAME="Deutsch", \
    DEFAULT=NO,AUTOSELECT=YES,LANGUAGE="de", \
    URI="main/german-audio.m3u8"
#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="aac",NAME="Commentary", \
    DEFAULT=NO,AUTOSELECT=NO,LANGUAGE="en", \
    URI="commentary/audio-only.m3u8"
#EXT-X-STREAM-INF:BANDWIDTH=1280000,CODECS="...",AUDIO="aac"
low/video-only.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=2560000,CODECS="...",AUDIO="aac"
mid/video-only.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=7680000,CODECS="...",AUDIO="aac"
hi/video-only.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=65000,CODECS="mp4a.40.5",AUDIO="aac"
main/english-audio.m3u8
```


### 8.7. Master Playlist with Alternative Video

此示例显示了三种不同的视频版本（ `Main` 、 `Centerfield` 和 `Dugout` ）和三种不同的变体流（低、中和高）。
在此示例中，不支持 `EXT-X-MEDIA` 标签和 `EXT-X-STREAM-INF` 标签的 `VIDEO` 属性的客户端只能播放视频版本 `Main` 。

由于 `EXT-X-STREAM-INF` 标签没有 `AUDIO` 属性，因此所有视频版本都需要包含音频。

在此示例中， `CODECS` 属性已压缩以节省空间。
使用 `\` 表示标签在删除空格后继续到下一行：

```m3u
#EXTM3U
#EXT-X-MEDIA:TYPE=VIDEO,GROUP-ID="low",NAME="Main", \
    DEFAULT=YES,URI="low/main/audio-video.m3u8"
#EXT-X-MEDIA:TYPE=VIDEO,GROUP-ID="low",NAME="Centerfield", \
    DEFAULT=NO,URI="low/centerfield/audio-video.m3u8"
#EXT-X-MEDIA:TYPE=VIDEO,GROUP-ID="low",NAME="Dugout", \
    DEFAULT=NO,URI="low/dugout/audio-video.m3u8"

#EXT-X-STREAM-INF:BANDWIDTH=1280000,CODECS="...",VIDEO="low"
low/main/audio-video.m3u8

#EXT-X-MEDIA:TYPE=VIDEO,GROUP-ID="mid",NAME="Main", \
    DEFAULT=YES,URI="mid/main/audio-video.m3u8"
#EXT-X-MEDIA:TYPE=VIDEO,GROUP-ID="mid",NAME="Centerfield", \
    DEFAULT=NO,URI="mid/centerfield/audio-video.m3u8"
#EXT-X-MEDIA:TYPE=VIDEO,GROUP-ID="mid",NAME="Dugout", \
    DEFAULT=NO,URI="mid/dugout/audio-video.m3u8"

#EXT-X-STREAM-INF:BANDWIDTH=2560000,CODECS="...",VIDEO="mid"
mid/main/audio-video.m3u8

#EXT-X-MEDIA:TYPE=VIDEO,GROUP-ID="hi",NAME="Main", \
    DEFAULT=YES,URI="hi/main/audio-video.m3u8"
#EXT-X-MEDIA:TYPE=VIDEO,GROUP-ID="hi",NAME="Centerfield", \
    DEFAULT=NO,URI="hi/centerfield/audio-video.m3u8"
#EXT-X-MEDIA:TYPE=VIDEO,GROUP-ID="hi",NAME="Dugout", \
    DEFAULT=NO,URI="hi/dugout/audio-video.m3u8"

#EXT-X-STREAM-INF:BANDWIDTH=7680000,CODECS="...",VIDEO="hi"
hi/main/audio-video.m3u8
```


### 8.8. Session Data in a Master Playlist

在此示例中，仅显示 `EXT-X-SESSION-DATA` ：

```m3u
#EXT-X-SESSION-DATA:DATA-ID="com.example.lyrics",URI="lyrics.json"

#EXT-X-SESSION-DATA:DATA-ID="com.example.title",LANGUAGE="en", \
    VALUE="This is an example"
#EXT-X-SESSION-DATA:DATA-ID="com.example.title",LANGUAGE="es", \
    VALUE="Este es un ejemplo"
```


### 8.9. CHARACTERISTICS Attribute Containing Multiple Characteristics

某些特征组合起来是有效的，例如：

```m3u
CHARACTERISTICS=
"public.accessibility.transcribes-spoken-dialog,public.easy-to-read"
```


### 8.10. EXT-X-DATERANGE Carrying SCTE-35 Tags

此示例显示了两个描述单个日期范围的 `EXT-X-DATERANGE` 标签，其中 `SCTE-35 "out" splice_insert()` 命令随后使用 `SCTE-35 "in" splice_insert()` 命令进行更新。

```m3u
#EXTM3U
...
#EXT-X-DATERANGE:ID="splice-6FFFFFF0",START-DATE="2014-03-05T11:
15:00Z",PLANNED-DURATION=59.993,SCTE35-OUT=0xFC002F0000000000FF0
00014056FFFFFF000E011622DCAFF000052636200000000000A0008029896F50
000008700000000

... Media Segment declarations for 60s worth of media

#EXT-X-DATERANGE:ID="splice-6FFFFFF0",DURATION=59.993,SCTE35-IN=
0xFC002A0000000000FF00000F056FFFFFF000401162802E6100000000000A00
08029896F50000008700000000
...
```


## 9. IANA Considerations

IANA 已注册以下媒体类型 [RFC2046][RFC2046]：

Type name: application  
Subtype name: vnd.apple.mpegurl  
Required parameters: none  
Optional parameters: none  
Encoding considerations：编码为 UTF-8，即 8 位文本。此媒体类型可能需要在无法处理 8 位文本的传输上进行编码。有关更多信息，请参阅[第 4 节](#4-playlists) 。  
Security considerations: See [Section 10](#10-security-considerations) .  
Compression: 此媒体类型不采用压缩。

互操作性注意事项：由于文件是 8 位文本，因此不存在字节顺序问题。应用程序可能会遇到无法识别的标签，这些标签 *应该* 被忽略。

已发布的规范：参见 [第 4 节](#4-playlists) 。

使用此媒体类型的应用程序：多媒体应用程序，例如 iOS 3.0 及更高版本中的 iPhone 媒体播放器和 Mac OS X 版本 10.6 及更高版本中的 QuickTime Player。

片段标识符注意事项：此媒体类型未定义片段标识符。

附加信息：

> 此类型的弃用别名：无  
> 幻数(magic numbers)：#EXTM3U  
> 文件扩展名：.m3u8、.m3u（参见 [第 4 节](#4-playlists)）  
> Macintosh 文件类型代码：无  


如需更多信息，请联系以下人员和电子邮件地址：David Singer，<singer@apple.com>。

Intended usage: LIMITED USE

Restrictions on usage: none

Author: Roger Pantos

Change Controller: David Singer


## 10. Security Considerations

由于该协议通常使用 HTTP 传输数据，因此大多数相同的安全注意事项都适用。请参阅 HTTP [[RFC7230][RFC7230]] 的 [第 15 节](https://datatracker.ietf.org/doc/html/rfc8216#section-15)。

媒体文件解析器通常会受到 *模糊测试* 攻击。
实施者 *应该* 特别注意解析从服务器接收的数据的代码，并确保正确处理所有可能的输入。

播放列表文件包含 URI，客户端将使用该 URI 发出任意实体的网络请求。
客户端 *应该* 检查响应的范围以防止缓冲区溢出。
另请参阅 “统一资源标识符 (URI)：通用语法” [[RFC3986][RFC3986]] 的安全注意事项部分。

除了 URL 解析之外，此格式不采用任何形式的主动内容。

客户端 *应该* 将每个播放会话限制为合理的并发下载数量（例如四个），以避免导致拒绝服务（dos）攻击。

HTTP 请求通常包含会话状态（ `cookies` ），其中可能包含私人用户数据。
实现 *必须* 遵循 “HTTP 状态管理机制” [[RFC6265][RFC6265]] 指定的 cookie 限制和到期规则，以保护自己免受攻击。
另请参阅该文档的安全注意事项部分和 “HTTP 状态管理的使用” [[RFC2964][RFC2964]]。

加密密钥由 URI 指定。
这些密钥的传送 *应该* 通过 HTTP Over TLS [[RFC2818][RFC2818]]（以前称为 SSL）等机制与安全领域或会话令牌相结合来保护。



[RFC2046]: https://datatracker.ietf.org/doc/html/rfc2046
[RFC2818]: https://datatracker.ietf.org/doc/html/rfc2818
[RFC2964]: https://datatracker.ietf.org/doc/html/rfc2964
[RFC3629]: https://datatracker.ietf.org/doc/html/rfc3629
[RFC3986]: https://datatracker.ietf.org/doc/html/rfc3986
[RFC5646]: https://datatracker.ietf.org/doc/html/rfc5646
[RFC5652]: https://datatracker.ietf.org/doc/html/rfc5652
[RFC6265]: https://datatracker.ietf.org/doc/html/rfc6265
[RFC6381]: https://datatracker.ietf.org/doc/html/rfc6381
[RFC7159]: https://datatracker.ietf.org/doc/html/rfc7159
[RFC7230]: https://datatracker.ietf.org/doc/html/rfc7230
<!-- [COMMON-ENC]: https://datatracker.ietf.org/doc/html/rfc8216#ref-COMMON_ENC -->
[COMMON-ENC]: https://www.iso.org/standard/68042.html
<!-- [H-264]: https://datatracker.ietf.org/doc/html/rfc8216#ref-H_264 -->
[H-264]: https://www.itu.int/rec/T-REC-H.264
<!-- [ISO-14496]: https://datatracker.ietf.org/doc/html/rfc8216#ref-ISO_14496 -->
[ISO-14496]: https://www.iso.org/standard/53943.html
<!-- [AC-3]: https://datatracker.ietf.org/doc/html/rfc8216#ref-AC_3 -->
[AC-3]: https://www.atsc.org/wp-content/uploads/2015/03/A52-201212-17.pdf
<!-- [Sample-Enc]: https://datatracker.ietf.org/doc/html/rfc8216#ref-SampleEnc -->
[Sample-Enc]: https://developer.apple.com/library/archive/documentation/AudioVideo/Conceptual/HLS_Sample_Encryption/Intro/Intro.html
<!-- SCTE35 Page Not Found -->
[SCTE35]: https://datatracker.ietf.org/doc/html/rfc8216#ref-SCTE35
<!-- [ISO-13818]: https://datatracker.ietf.org/doc/html/rfc8216#ref-ISO_13818 -->
[ISO-13818]: https://www.iso.org/standard/44169.html
<!-- [ISO-13818-3]: https://datatracker.ietf.org/doc/html/rfc8216#ref-ISO_13818_3 -->
[ISO-13818-3]: https://www.iso.org/standard/26797.html
<!-- [ISO-13818-7]: https://datatracker.ietf.org/doc/html/rfc8216#ref-ISO_13818_7 -->
[ISO-13818-7]: https://www.iso.org/standard/43345.html
<!-- [ISOBMFF]: https://datatracker.ietf.org/doc/html/rfc8216#ref-ISOBMFF -->
[ISOBMFF]: https://www.iso.org/standard/68960.html
<!-- [CMAF]: https://datatracker.ietf.org/doc/html/rfc8216#ref-CMAF -->
[CMAF]: https://www.iso.org/standard/71975.html
<!-- ID3 Not Found -->
[ID3]: https://datatracker.ietf.org/doc/html/rfc8216#ref-ID3
<!-- WebVTT Not Found -->
[WebVTT]: https://datatracker.ietf.org/doc/html/rfc8216#ref-WebVTT
<!-- [M3U]: https://datatracker.ietf.org/doc/html/rfc8216#ref-M3U -->
<!-- M3U Not found, I think this one -->
[M3U]: https://en.wikipedia.org/wiki/M3U
<!-- [UNICODE]: https://datatracker.ietf.org/doc/html/rfc8216#ref-UNICODE -->
[UNICODE]: https://www.unicode.org/versions/Unicode15.1.0/
[US-ASCII]: https://datatracker.ietf.org/doc/html/rfc8216#ref-US_ASCII
<!-- [ISO-8601]: https://datatracker.ietf.org/doc/html/rfc8216#ref-ISO_8601 -->
[ISO-8601]: https://www.iso.org/standard/40874.html
[CEA608]: https://datatracker.ietf.org/doc/html/rfc8216#ref-CEA608
<!-- [UTI]: https://datatracker.ietf.org/doc/html/rfc8216#ref-UTI -->
[UTI]: https://developer.apple.com/library/archive/documentation/General/Conceptual/DevPedia-CocoaCore/UniformTypeIdentifier.html
<!-- [HDCP]: https://datatracker.ietf.org/doc/html/rfc8216#ref-HDCP -->
[HDCP]: https://www.digital-cp.com/sites/default/files/specifications/HDCP%20on%20HDMI%20Specification%20Rev2_2_Final1.pdf

