---
title: cgroup-v2 译文
date: 2024-01-03 20:38:15
tags:
categories:
- kernel doc
- cgroup-v2
keywords:
- cgroup-v2
copyright: Guader
copyright_author_href:
copyright_info:
---

[原文](https://www.kernel.org/doc/Documentation/cgroup-v2.txt)

这是关于cgroup v2的设计、接口和约定的权威文档。
它描述了 cgroup 的所有用户可见的方面，包括核心和特定控制器行为。
未来的所有变更都必须反映在本文件中。
v1 的文档可在 Documentation/cgroup-v1/ 下找到。


# 1. Introduction


## 1.1 Terminology

`cgroup` 代表 `control group` 并且从来不大写。
单数形式用于指定整个功能，也用作 `cgroup controller` 中的限定符。
当明确指代多个单独的对照组时，使用复数形式 `cgroups`。


## 1.2 What is cgroup ?

cgroup 是一种按层次结构组织进程并以受控和可配置的方式沿层次结构分配系统资源的机制。

cgroup 大体上由两部分组成 - 核心和控制器。
cgroup 核心主要负责分层组织进程。
cgroup 控制器通常负责沿层次结构分配特定类型的系统资源，尽管还有一些实用程序控制器可用于资源分配以外的目的。

cgroup 形成一种树形结构，系统中的每个进程都属于一个且唯一的 cgroup。
一个进程的所有线程都属于同一个 cgroup。
创建时，所有进程都被放入父进程当时所属的 cgroup 中。
一个进程可以迁移到另一个 cgroup。 
进程的迁移不会影响已经存在的后代进程。

遵循某些结构约束，可以在 cgroup 上选择性地启用或禁用控制器。
所有控制器行为都是分层的 - 如果在 cgroup 上启用控制器，它会影响属于包含该 cgroup 子层次结构的 cgroup 的所有进程。
当在嵌套 cgroup 上启用控制器时，它始终会进一步限制资源分配。
靠近层次结构中的根设置的限制不能从更远的地方覆盖。


# 2. Basic Operations


## 2.1 Mounting

和 v1 不同，cgroup v2 只有单一层次结构。
可以使用以下挂载命令挂载 cgroup v2 层次结构：

`mount -t cgroup2 none $MOUNT_POINT`

cgroup2 文件系统具有幻数 `0x63677270`（“cgrp”）。
所有支持 v2 并且未绑定到 v1 层次结构的控制器都会自动绑定到 v2 层次结构并显示在根目录中。
v2 层次结构中未主动使用的控制器可以绑定到其他层次结构。
这允许以完全向后兼容的方式将 v2 层次结构与旧版 v1 多个层次结构混合。

仅当控制器的当前层次结构中不再引用该控制器时，才可以跨层次结构移动该控制器。
由于每个 cgroup 控制器状态被异步销毁，并且控制器可能具有延迟引用，因此在先前层次结构最终卸载后，控制器可能不会立即显示在 v2 层次结构上。
同样，控制器应完全禁用才能移出统一层次结构，并且禁用的控制器可能需要一些时间才能可用于其他层次结构； 此外，由于控制器间的依赖性，其他控制器可能也需要禁用。

虽然对于开发和手动配置很有用，但强烈建议不要在生产使用中在 v2 和其他层次结构之间动态移动控制器。
建议在系统启动后开始使用控制器之前确定层次结构和控制器关联。

在过渡到 v2 期间，系统管理软件可能仍会自动挂载 v1 cgroup 文件系统，因此在启动期间劫持所有控制器，然后才能进行手动干预。
为了使测试和实验更容易，内核参数 `cgroup_no_v1=` 允许禁用 v1 中的控制器并使其在 v2 中始终可用。

cgroup v2 当前支持以下挂载选项。

- `nsdelegate`
    将 cgroup 命名空间视为委托边界。
    此选项是系统范围的，只能在挂载时设置或通过从 init 命名空间重新挂载进行修改。
    在非 `init` 命名空间挂载上，挂载选项将被忽略。
    详情请参阅 `Delegation` 部分。


## 2.2 Organizing Processes and Threads


### 2.2.1 Processes

最初，仅存在所有进程所属的根 cgroup。
可以通过创建子目录来创建子 cgroup::

`mkdir $CGROUP_NAME`

给定的 cgroup 可能有多个子 cgroup，形成树结构。
每个 cgroup 都有一个可读写的接口文件 `cgroup.procs`。
读取时，它会逐行列出属于该 cgroup 的所有进程的 PID。
PID 没有排序，如果进程移动到另一个 cgroup 然后又返回，或者 PID 在读取时被回收，则相同的 PID 可能会多次出现。

通过将进程的 PID 写入目标 cgroup 的 `cgroup.procs` 文件，可以将进程迁移到 cgroup 中。
一次 `write(2)` 调用只能迁移一个进程。
如果一个进程由多个线程组成，写入任意线程的PID就会迁移该进程的所有线程。

当一个进程 `fork` 出一个子进程时，新进程就会诞生到 `fork` 进程运行时所属的 cgroup 中。
退出后，进程将与其退出时所属的 cgroup 保持关联，直到被回收；
但是，僵尸进程不会出现在 `cgroup.procs` 中，因此无法移动到另一个 cgroup。

没有任何子进程或活动进程的 cgroup 可以通过删除目录来销毁。
请注意，没有任何子进程且仅与僵尸进程关联的 cgroup 被视为空的，可以删除::

`rmdir $CGROUP_NAME`

`/proc/$PID/cgroup` 列出进程的 cgroup 成员资格。
如果系统中正在使用旧版 cgroup，则此文件可能包含多行，每个层次结构一行。
cgroup v2 的条目始终采用 `0::$PATH` 格式

```bash
cat /proc/842/cgroup

...
0::test-cgroup/test-cgroup-nested
```

如果该进程成为僵尸进程，并且随后删除了与其关联的 cgroup，则 `(deleted)` 将附加到路径::

```bash
cat /proc/842/cgroup

...
0::/test/cgroup/test-cgroup-nested (deleted)
```


### 2.2.2 Threads

cgroup v2 支持控制器子集的线程粒度，以支持需要跨一组进程的线程进行分层资源分配的用例。
默认情况下，进程的所有线程都属于同一个 cgroup，该 cgroup 也充当资源域，用于托管非特定于进程或线程的资源消耗。
线程模式允许线程分布在子树上，同时仍然维护它们的公共资源域。

支持线程模式的控制器称为线程控制器 `threaded controllers`。
没有的称为域控制器 `domain controllers`。

将 cgroup 标记为线程化会使其作为线程化 cgroup 加入其父级的资源域。
父级可能是另一个线程 cgroup，其资源域在层次结构中更靠上。
线程子树的根，即非线程化的最近祖先，可互换地称为线程域或线程根，并充当整个子树的资源域。

在线程子树内部，进程的线程可以放入不同的 cgroup 中，并且不受无内部进程约束 - 可以在非叶 cgroup 上启用线程控制器，无论它们是否有线程。

由于 `threaded domain cgroup` 托管子树的所有域资源消耗，因此无论其中是否有进程，它都被认为具有内部资源消耗，并且不能填充非线程化的子 cgroup。
由于根 cgroup 不受内部进程约束，因此它既可以充当线程域，也可以充当域 cgroup 的父级。

cgroup当前的操作模式或类型显示在 `cgroup.type` 文件中，该文件指示 cgroup 是`normal domain`、`a domain which is serving as the domain of a threaded subtree` 还是 `threaded cgroup`。

创建时，cgroup 始终是 `domain cgroup`，并且可以通过将 `threaded` 写入`cgroup.type` 文件来使其线程化。
操作是单向的::

`echo threaded > cgroup.type`

一旦线程化，cgroup 就无法再次成为 `domain`。
要启用线程模式，必须满足以下条件。

- 因为 cgroup 将加入父级的资源域。父级必须是有效（线程）域或线程 cgroup。
- 当父域是非线程域时，它不得启用任何域控制器或填充域子域。根不受此要求的约束。

从拓扑角度来看，cgroup 可能处于无效状态。
请考虑以下拓扑:

`A (threaded domain) - B (threaded) - C (domain, just created)`

C 作为域创建，但未连接到可以托管子域的父域。
C 必须变成 `threaded cgroup` 才能使用。
在这些情况下，`cgroup.type` 文件将报告 `domain (invalid)`。
由于无效拓扑而失败的操作使用 `EOPNOTSUPP` 作为 `errno`。

当 `domain cgroup` 的一个子 cgroup 成为线程化或在 `cgroup.subtree_control` 文件中启用线程控制器且 cgroup 中有进程时，该 `threaded cgroup` 将转变为 `threaded domain` 。
当条件清除时，线程域将恢复为正常域。

读取时，`cgroup.threads` 包含 cgroup 中所有线程的线程ID列表。
除了操作是针对每个线程而不是针对每个进程之外，`cgroup.threads` 与 `cgroup.procs` 具有相同的格式和行为方式。
虽然 `cgroup.threads` 可以写入任何 cgroup 中，但由于它只能在同一线程域内移动线程，因此其操作仅限于每个线程子树内。

线程域 cgroup 充当整个子树的资源域，并且虽然线程可以分散在子树中，但所有进程都被视为位于线程域 cgroup 中。
线程域 cgroup 中的 `cgroup.procs` 包含子树中所有进程的 PID，并且在子树中不可读。
但是，可以从子树中的任何位置写入 `cgroup.procs` ，以将匹配进程的所有线程迁移到 cgroup。

在线程子树中只能启用线程控制器。
当线程控制器在线程子树内启用时，它仅考虑和控制与 cgroup 及其后代中的线程相关的资源消耗。
所有不依赖于特定线程的消耗都属于线程域 cgroup。

由于线程子树不受任何内部进程约束，因此线程控制器必须能够处理非叶 cgroup 及其子 cgroup 中的线程之间的竞争。
每个线程控制器定义了如何处理此类竞争。


## 2.3 [Un]populated Notification

每个非根 cgroup 都有一个 `cgroup.events` 文件，其中包含 `populated` 字段，指示 cgroup 的子层次结构中是否有实时进程。
如果 cgroup 及其后代中没有活动进程，则其值为 0；否则，1。
当值发生变化时，会触发 `poll` 和 `[id]notify` 事件。
例如，这可以用于在给定子层次结构的所有进程退出后启动清理操作。
填充的状态更新和通知是递归的。
考虑以下子层次结构，其中括号中的数字表示每个 cgroup 中的进程数：

```
A(4) - B(0) - C(1)
            \ D(0)
```

A、B 和 C 的 `populated` 字段将为 1，而 D 的字段将为 0。
C 中的一个进程退出后，B 和 C 的 `populated` 字段将翻转为 “0”，并且文件修改事件将在两个 cgroup 的 `cgroup.events` 文件上生成。


## 2.4 Controlling Controllers


### 2.4.1 Enabling and Disabling

每个 cgroup 都有一个 `cgroup.controllers` 文件，其中列出了可供该 cgroup 启用的所有控制器：

```bash
cat cgroup.controllers

# cpu io memory
```

默认情况下不启用任何控制器。
可以通过写入 `cgroup.subtree_control` 文件来启用和禁用控制器：

`echo "+cpu +memory -io" > cgroup.subtree_control`

只能启用 `cgroup.controllers` 中列出的控制器。
当如上所述指定多个操作时，它们要么全部成功，要么全部失败。
如果对同一控制器指定了多次操作，则最后一次有效。

启用 cgroup 中的控制器表示目标资源在其直接子级之间的分配将受到控制。
考虑以下子层次结构。
启用的控制器列在括号中:

```
A(cpu,memory) - B(memory) - C()
                          \ D()
```

由于 A 启用了 `cpu` 和 `memory`，因此 A 将控制向其子级（在本例中为 B）分配 CPU 周期和内存。
由于 B 启用了 `memory` 但没有启用 `CPU`，C 和 D 将在 CPU 周期上自由竞争，但 B 可用的内存分配将受到控制。

当控制器调节目标资源到 cgroup 子级的分配时，使其能够在子 cgroup 中创建控制器的接口文件。
在上面的示例中，在 B 上启用 `cpu` 将在 C 和 D 中创建带有 `cpu.` 前缀的控制器接口文件。
同样，从 B 禁用 `memory` 将从 C 和 D 中删除以 `memory` 为前缀的控制器接口文件。
这意味着控制器接口文件 - 任何不以 `cgroup.` 开头的文件都由父级而不是 cgroup 本身拥有。


### 2.4.2 Top-down Constraint

资源是自上而下分配的，只有当资源已从父级分配给 cgroup 时，cgroup 才能进一步分配资源。
这意味着所有非根 `cgroup.subtree_control` 文件只能包含在父级 `cgroup.subtree_control` 文件中启用的控制器。
仅当父级启用了控制器时才能启用该控制器，并且如果一个或多个子级启用了控制器则无法禁用该控制器。


### 2.4.3 No Internal Process Constraint

非根 cgroup 仅当它们没有自己的任何进程时才可以将域资源分配给其子级。
换句话说，只有不包含任何进程的域 cgroup 才能在其 `cgroup.subtree_control` 文件中启用域控制器。

这保证了当域控制器查看启用它的层次结构部分时，进程始终仅位于叶子上。
这排除了子 cgroup 与父 cgroup 的内部进程竞争的情况。

根 cgroup 不受此限制。
Root 包含进程和匿名资源消耗，它们不能与任何其他 cgroup 关联，并且需要大多数控制器的特殊处理。
如何管理根 cgroup 中的资源消耗取决于每个控制器（有关此主题的更多信息，请参阅控制器章节中的非规范信息部分）。

请注意，如果 cgroup 的 `cgroup.subtree_control` 中没有启用的控制器，则该限制不会妨碍。
这很重要，否则将无法创建已填充 cgroup 的子级。
为了控制 cgroup 的资源分配，cgroup 必须创建子级并将其所有进程转移到子级，然后才能在其 `cgroup.subtree_control` 文件中启用控制器。


## 2.5 Delegation


### 2.5.1 Model of Delegation

cgroup 可以通过两种方式进行委托。
首先，通过向用户授予目录及其 `cgroup.procs`、`cgroup.threads` 和`cgroup.subtree_control` 文件的写访问权限来授予特权较低的用户。
其次，如果设置了 `nsdelegate` 挂载选项，则在创建命名空间时自动到 cgroup 命名空间。

由于给定目录中的资源控制接口文件控制父级资源的分配，因此不应允许受委托者对它们进行写入。
对于第一种方法，这是通过不授予对这些文件的访问权限来实现的。
对于第二个，内核拒绝从命名空间内部写入命名空间根上除 `cgroup.procs` 和 `cgroup.subtree_control` 之外的所有文件。

两种委托类型的最终结果是相同的。
一旦委派，用户就可以在目录下构建子层次结构，根据需要组织其中的流程，并进一步分配从父级接收的资源。
所有资源控制器的限制和其他设置都是分层的，无论委托的子层次结构中发生什么，没有任何东西可以逃脱父级施加的资源限制。

目前，cgroup 对委派子层次结构中的 cgroup 数量或嵌套深度没有任何限制； 但是，将来这可能会受到明确限制。


### 2.5.2 Delegation Containment

委派的子层次结构包含在进程不能由受委派者移入或移出子层次结构的意义上。

对于权限较低的用户的委派，这是通过要求具有非 root euid 的进程满足以下条件来将目标进程迁移到 cgroup（通过将其 PID 写入 `cgroup.procs` 文件）来实现的。

- 写入者必须具有对 `cgroup.procs` 文件的写入权限
- 写入者必须对源 cgroup 和目标 cgroup 的共同祖先的 `cgroup.procs` 文件具有写入权限。

上述两个约束确保虽然委托者可以在委托的子层次结构中自由地迁移进程，但它不能从子层次结构外部拉入或推出到子层次结构之外。

例如，假设 cgroup `C0` 和 `C1` 已委托给用户 `U0`，用户 `U0` 在 `C0` 下创建了 `C00`、`C01`，在 `C1 下创建了 `C10`，如下所示，并且 `C0` 和 `C1` 下的所有进程都属于 `U0`：

```
~~~~~~~~~~~~~ - C0 - C00
~ cgroup    ~      \ C01
~ hierarchy ~       
~~~~~~~~~~~~~ - C1 - C10
```

还假设 `U0` 想要将当前位于 `C10` 中的进程的 PID 写入 `C00/cgroup.procs`。
`U0` 对文件有写权限； 然而，源 cgroup `C10` 和目标 cgroup `C00` 的共同祖先位于委派点之上，并且 `U0` 对其 `cgroup.procs` 文件没有写入访问权限，因此写入将被 `EACCES` 拒绝。

对于命名空间的委派，通过要求源 cgroup 和目标 cgroup 都可从尝试迁移的进程的命名空间访问来实现遏制。
如果其中一个无法访问，则迁移会被 `ENOENT` 拒绝。


## 2.6 Guidelines


### 2.6.1 Organize Once and Control

跨 cgroup 迁移进程是一项相对昂贵的操作，并且内存等有状态资源不会随进程一起移动。
这是一个明确的设计决策，因为迁移和各种热路径之间在同步成本方面通常存在固有的权衡。

因此，不鼓励频繁地跨 cgroup 迁移进程作为应用不同资源限制的手段。
启动时，应根据系统的逻辑和资源结构将工作负载分配给 cgroup。
可以通过接口文件更改控制器配置来动态调整资源分配。


### 2.6.2 Avoid Name Collisions

cgroup 及其子 cgroup 的接口文件占用相同的目录，并且可以创建与接口文件冲突的子 cgroup。

所有 cgroup 核心接口文件都以 “cgroup.” 为前缀，每个控制器的接口文件都以控制器名称和点为前缀。
控制器的名称由小写字母和 `_` 组成，但绝不以 `_` 开头，因此可以用作避免冲突的前缀字符。
此外，接口文件名不会以工作负载分类中常用的术语开头或结尾，例如 `job`、`service`、`slice`、`unit` 或 `workload`。

cgroup 不会采取任何措施来防止名称冲突，用户有责任避免名称冲突。


# 3. Resource Distribution Models

cgroup 控制器根据资源类型和预期用例实现多种资源分配方案。
本节描述了正在使用的主要方案及其预期行为。


## 3.1 Weights

通过将所有活动子级的权重相加并给予每个子级与其权重与总和的比率相匹配的分数来分配父级的资源。
由于只有当前可以使用资源的子级才参与分配，因此这是节省工作的。
由于动态特性，该模型通常用于无状态资源。

所有权重均在 [1, 10000] 范围内，默认值为 100。
这允许在两个方向上以足够细的粒度实现对称乘法偏差，同时保持在直观范围内。

只要权重在范围内，所有配置组合都是有效的，没有理由拒绝配置更改或流程迁移。

`cpu.weight` 按比例将 CPU 周期分配给活动的子进程，是这种类型的一个示例。


## 3.2 Limits

子级最多只能消耗配置的资源量。
限制可能会被过度使用 - 子级的限制总和可能超过父级可用的资源量。

限制范围为 [0, max]，默认为 `max`，即 noop。

由于限制可能会被过度使用，因此所有配置组合都是有效的，没有理由拒绝配置更改或流程迁移。

`io.max` 限制 cgroup 可在 IO 设备上消耗的最大 BPS 和/或 IOPS，是此类的一个示例。


## 3.3 Protections

如果 cgroup 的所有祖先的使用量都在其受保护级别之下，则该 cgroup 将受到保护，最多可分配配置的资源量。
保护可以是硬保证，也可以是尽力而为的软边界。
保护也可能会过度承诺，在这种情况下，子级只能受到父级可用的保护。

保护范围为 [0, max]，默认为 0，即 noop。

由于保护可能会被过度使用，因此所有配置组合都是有效的，没有理由拒绝配置更改或流程迁移。

`memory.low` 实现尽力而为的内存保护，并且是这种类型的一个示例。


## 3.4 Allocations

一个cgroup 被专门分配一定数量的有限资源。
分配不能过度承诺 - 子级分配的总和不能超过父级可用的资源量。

分配范围为 [0, max]，默认为 0，即没有资源。

由于分配不能过度分配，因此某些配置组合无效，应被拒绝。
此外，如果该资源对于进程的执行是必需的，则进程迁移可能会被拒绝。

`cpu.rt.max` 硬分配实时切片，是这种类型的一个示例。


# 4. Interface Files


## 4.1 Format

所有接口文件应尽可能采用以下格式之一：

1. 换行分隔值
    （当一次只能写入一个值时）
    ```
    VAL0\n
    VAL1\n
    ...
    ```
2. 空格分隔值
    （当只读或是一次可以写入多个值时）
    ```
    VAL0 VAL1 ...\n
    ```
3. 平铺键
    ```
    KEY0 VAL0\n
    KEY1 VAL1\n
    ...
    ```
4. 嵌套键
    ```
    KEY0 SUB_KEY0=VAL00 SUB_KEY1=VAL01...
    KEY1 SUB_KEY0=VAL10 SUB_KEY1=VAL11...
    ...
    ```

对于可写文件，写入的格式一般应与读取的格式一致； 然而，控制器可能允许省略后面的字段或为最常见的用例实现受限的快捷方式。

对于平面和嵌套键控文件，一次只能写入单个键的值。
对于嵌套键控文件，可以按任何顺序指定子键值对，并且不必指定所有对。


## 4.2 Conventions

- 单个功能的设置应包含在单个文件中。
- 根 cgroup 应不受资源控制，因此不应具有资源控制接口文件。
  此外，根 cgroup 上最终显示其他地方可用的全局信息的信息文件不应该存在。
- 如果控制器实现基于权重的资源分配，则其接口文件应命名为 `weight`，范围为 [1,10000]，默认为100。
  选择这些值是为了在两个方向上允许足够的对称偏差，同时保持直观（默认值为 100%）。
- 如果控制器实现绝对资源保证和/或限制，则接口文件应分别命名为 `min` 和 `max`。
  如果控制器实现尽力而为资源保证和/或限制，则接口文件应分别命名为 `low` 和 `high`。

  在上面的四个控制文件中，应该使用特殊标记 `max` 来表示读和写的向上无穷大。
- 如果设置具有可配置的默认值和键入的特定覆盖，则默认条目应键入 `default` 并显示为文件中的第一个条目。
    
  可以通过写入 `default $VAL` 或 `$VAL` 来更新默认值。

  当写入更新特定覆盖时，`default` 可以用作指示删除覆盖的值。
  使用 `default` 覆盖条目，因为读取时不得出现该值。

  例如，由具有整数值的 主设备号：次设备号 作为键控的设置可能如下所示：

  ```bash
  cat cgroup-example-interface-file

  default 150
  8:0 300
  ```

  默认值可以这样更新：

  `echo 125 > cgroup-example-interface-file`

  或者：

  `echo "default 125" > cgroup-example-interface-file`

  可以这样设置覆盖：

  `echo "8:16 170" > cgroup-example-interface-file`

  以及这样清除：

  ```bash
  echo "8:0 default" > cgroup-example-interface-file
  cat cgroup-example-interface-file

  default 125
  8:16 170
  ```

- 对于频率不是很高的事件，应创建一个接口文件 `events`，其中列出事件键值对。
  每当发生可通知事件时，应在文件上生成文件修改事件。


## 4.3 Core Interface Files
所有 cgroup 核心文件都以 `cgroup` 为前缀。

- **`cgroup.type`**
    存在于非根 cgroup 上的读写单值文件。

    读取时，它指示 cgroup 的当前类型，可以是以下值之一。

    - **`domain`**  正常的有效域 cgroup。
    - **`domain threaded`**  线程域 cgroup，用作线程子树的根。
    - **`domain invalid`**  处于无效状态的 cgroup。
        它无法填充或启用控制器。
        它可能被允许成为线程化的cgroup。
    - **`threaded`**  线程 cgroup，它是线程子树的成员。

    通过向此文件写入 `threaded`，可以将 cgroup 转变为 `threaded cgroup`。

- **`cgroup.procs`** 
    存在于所有 cgroup 上的读写换行分隔值文件。

    读取时，它会逐行列出属于该 cgroup 的所有进程的 PID。
    PID 没有排序，如果进程移动到另一个 cgroup 然后又返回，或者 PID 在读取时被回收，则相同的 PID 可能会多次出现。

    可以写入PID，将与该PID关联的进程迁移到cgroup中。
    写入者应满足以下所有条件。

    - 它必须具有对 `cgroup.procs` 文件的写访问权限。
    - 它必须对源 cgroup 和目标 cgroup 的共同祖先的 `cgroup.procs` 文件具有写入权限。

    委派子层次结构时，应授予对此文件及其包含目录的写访问权限。

    在线程 cgroup 中，读取此文件会失败并显示 EOPNOTSUPP，因为所有进程都属于线程根。
    支持写入，并将进程的每个线程移动到 cgroup。

- **`cgroup.threads`**
    存在于所有 cgroup 上的读写换行分隔值文件。

    读取时，它会逐行列出属于该 cgroup 的所有线程的 TID。
    TID 没有排序，如果线程移动到另一个 cgroup 然后又返回，或者 TID 在读取时被回收，则相同的 TID 可能会出现多次。

    可以写入TID，将与该TID关联的线程迁移到 cgroup 中。
    写入者应满足以下所有条件。

    - 它必须具有对 `cgroup.threads` 文件的写访问权限。
    - 线程当前所在的 cgroup 必须与目标 cgroup 位于同一资源域中。
    - 它必须对源 cgroup 和目标 cgroup 的共同祖先的 `cgroup.procs` 文件具有写入权限。

    委派子层次结构时，应授予对此文件及其包含目录的写访问权限。

- **`cgroup.controllers`**
    存在于所有 cgroup 上的只读空格分隔值文件。

    它显示了 cgroup 可用的所有控制器的空格分隔列表。 
    控制器未排序。

- **`cgroup.subtree_control`**
    存在于所有 cgroup 上的读写空格分隔值文件。 
    开始是空的。

    读取时，它显示以空格分隔的控制器列表，这些控制器用于控制从 cgroup 到其子级的资源分配。

    可以写入以“+”或“-”为前缀的空格分隔的控制器列表来启用或禁用控制器。
    控制器名称以 `+` 为前缀可启用控制器，`-` 则可禁用。
    如果某个控制器在列表中出现多次，则最后一个有效。
    当指定多个启用和禁用操作时，要么全部成功，要么全部失败。

- **`cgroup.events`**
    存在于非根 cgroup 上的只读平键文件。
    定义了以下条目。
    除非另有指定，否则此文件中的值更改会生成文件修改事件。

    - `populated`
        如果 cgroup 或其后代包含任何活动进程，为 1；否则，0。

- **`cgroup.max.descendants`**
    读写单值文件。 默认值为 `max`。

    允许的最大下降 cgroup 数量。
    如果后代的实际数量等于或大于，则尝试在层次结构中创建新的 cgroup 将失败。

- **`cgroup.max.depth`**
    读写单值文件。默认值为 `max`。

    当前 cgroup 下方允许的最大下降深度。
    如果实际下降深度等于或更大，则尝试创建新的子 cgroup 将失败。

- **`cgroup.stat`**
    具有以下条目的只读平键文件：

    - `nr_descendants`  可见后代 cgroups 的总数。
    - `nr_dying_descendants`  垂死的后代 cgroup 总数。
        cgroup 在被用户删除后就会死亡。
        在完全销毁之前，cgroup 将在一段不确定的时间内保持死亡状态（这可能取决于系统负载）。
        进程在任何情况下都无法进入垂死的 cgroup，垂死的 cgroup 无法复活。

        垂死的 cgroup 可以消耗不超过限制的系统资源，这些资源在 cgroup 删除时处于活动状态。


# 5. Controllers


## 5.1 CPU

`cpu` 控制器调节 CPU 周期的分配。
该控制器实现了正常调度策略的权重和绝对带宽限制模型以及实时调度策略的绝对带宽分配模型。

警告： cgroup2 尚不支持实时进程的控制，并且只有当所有 RT 进程都位于根 cgroup 中时才能启用 cpu 控制器。
请注意，系统管理软件可能已在系统引导过程中将 RT 进程放入非 root cgroup，并且可能需要将这些进程移至 root cgroup，然后才能启用 cpu 控制器。


### 5.1.1 CPU Interface Files

所有持续时间均以微秒为单位。

- **`cpu.stat`**
    存在于非根 cgroup 上的只读平键文件。
    无论控制器是否启用，该文件都存在。

    它始终报告以下三个统计数据：

    - `usage_usec`
    - `user_usec`
    - `system_usec`

    当控制器启用时，还有以下三个：

    - `nr_periods`
    - `nr_throttled`
    - `throttled_usec`

- **`cpu.weight`**
    存在于非根 cgroup 上的读写单值文件。 
    默认值为 `100`。

    权重的范围时 [1, 10000]。

- **`cpu.weight.nice`**
    存在于非根 cgroup 上的读写单值文件。 
    默认值为 `0`。

    `nice` 值的范围是 [-20, 19]。

    该接口文件是 `cpu.weight` 的替代接口，允许使用与 `nice(2)` 相同的值读取和设置权重。
    由于良好值的范围较小且粒度较粗，因此读取的值是当前权重的最接近的近似值。

- **`cpu.max`**
    存在于非根 cgroup 上的读写二值文件。
    默认值为 `max 100000`。

    最大带宽限制。它的格式如下:

    `$MAX $PERIOD`

    这表明该组在每个 `$PERIOD` 持续时间内最多可以消耗 `$MAX`。
    `$MAX` 的 `max` 表示没有限制。
    如果只写入一个数字，则更新 `$MAX`。


## 5.2 Memory

`memory` 控制器调节内存的分配。
内存是有状态的，并实现限制和保护模型。
由于内存使用和回收压力之间的相互交织以及内存的有状态特性，分配模型相对复杂。

虽然不是完全无懈可击，但会跟踪给定 cgroup 的所有主要内存使用情况，以便可以计算总内存消耗并将其控制在合理的范围内。
目前，跟踪以下类型的内存使用情况。

- 用户态内存 - 页缓存和匿名内存
- 内核数据结构，例如 `dentry` 和 `inode`。
- TCP 套接字缓冲区

上述列表将来可能会扩大，以获得更好的覆盖范围。


### 5.2.1 Memory Interface Files

所有内存量均以字节为单位。
如果写入的值未与 `PAGE_SIZE` 对齐，则读回时该值可能会向上舍入为最接近的 `PAGE_SIZE` 倍数。

- **`memory.current`**
    存在于非根 cgroup 上的只读单值文件。

    cgroup 及其后代当前使用的内存总量。

- **`memory.low`**
    存在于非根 cgroup 上的读写单值文件。 
    默认值为 `0`。

    尽最大努力的内存保护。 
    如果某个 cgroup 及其所有祖先的内存使用量低于其下限，则该 cgroup 的内存将不会被回收，除非可以从未受保护的 cgroup 回收内存。

    不鼓励在此保护下放置比一般可用内存更多的内存。

- **`memory.high`**
    存在于非根 cgroup 上的读写单值文件。    
    默认值为 `max`。

    内存使用限制。
    这是控制 cgroup 内存使用的主要机制。
    如果 cgroup 的使用量超过上限，则该 cgroup 的进程将受到限制并承受沉重的回收压力。

    超过上限永远不会调用 OOM 杀手，并且在极端条件下可能会突破该限制。

- **`memory.max`**
    存在于非根 cgroup 上的读写单值文件。 
    默认值为 `max`。

    内存使用硬限制。
    这是最终的保护机制。
    如果某个 cgroup 的内存使用量达到此限制并且无法减少，则会在该 cgroup 中调用 OOM Killer。
    在某些情况下，使用量可能会暂时超出限制。

    这是最终的保护机制。 
    只要正确使用和监控上限，该限制的效用就仅限于提供最终的安全网。

- **`memory.events`**
    存在于非根 cgroup 上的只读平键文件。
    定义了以下条目。 
    除非另有指定，否则此文件中的值更改会生成文件修改事件。

    - **`low`** 尽管 cgroup 的使用率低于低边界，但由于内存压力较高而回收 cgroup 的次数。这通常表明低边界被过度使用。
    - **`high`**    由于超出高内存边界而对 cgroup 的进程进行限制和路由以执行直接内存回收的次数。 对于内存使用量受到上限而不是全局内存压力限制的 cgroup，此事件的发生是预料之中的。
    - **`max`** cgroup 的内存使用量即将超过最大边界的次数。 如果直接回收无法将其关闭，则 cgroup 将进入 OOM 状态。
    - **`oom`** cgroup内存使用达到限制并且分配即将失败的次数。
        根据上下文结果，可能会调用 OOM Killer 并重试分配或分配失败。
        失败的分配又可以作为 -ENOMEM 返回到用户空间，或者在磁盘预读等情况下默默地忽略。 目前，内存中的 OOM 如果在页面错误内发生短缺，则 cgroup 会终止任务。

    - **`oom_kill`**    被任何类型的 OOM 杀手杀死的属于此 cgroup 的进程数。

- **`memory.stat`**
    存在于非根 cgroup 上的只读平键文件。

    这将 cgroup 的内存占用量分解为不同类型的内存、特定于类型的详细信息以及有关内存管理系统的状态和过去事件的其他信息。

    所有内存量均以字节为单位。

    这些条目被排序为人类可读的，并且新条目可以显示在中间。 不要依赖保持在固定位置的物品； 使用按键查找特定值！

    - **`anon`**    匿名映射（例如 `brk()`、`sbrk()` 和 `mmap(MAP_ANONYMOUS)`）中使用的内存量
    - **`file`**    用于缓存文件系统数据的内存量，包括 tmpfs 和共享内存。
    - **`kernel_stack`**    分配给内核堆栈的内存量。
    - **`slab`**    用于存储内核数据结构的内存量。
    - **`sock`**    网络传输缓冲区使用的内存量
    - **`shmem`**   支持交换的缓存文件系统数据量，例如 `tmpfs`、`shm` 段、共享匿名 `mmap()`
    - **`file_mapped`** 使用 `mmap()` 映射的缓存文件系统数据量
    - **`file_dirty`**  已修改但尚未写回磁盘的缓存文件系统数据量
    - **`file_writeback`**  已修改且当前正在写回磁盘的缓存文件系统数据量
    - **`inactive_anon, active_anon, inactive_file, active_file, unevictable`**
        页面回收算法使用的内部内存管理列表上的内存量（支持交换和支持文件系统）
    - **`slab_reclaimable`**    可能被回收的 `slab` 的一部分，例如 `dentry` 和 `inode`。
    - **`slab_unreclaimable`**  由于内存压力而无法回收的 `slab` 的一部分。
    - **`pgfault`** 发生的页面错误总数
    - **`pgmajfault`**  发生重大页面故障的次数
    - **`workingset_refault`**  先前被驱逐页面的拒绝次数
    - **`workingset_activate`** 立即激活的拒绝页面数
    - **`workingset_nodereclaim`**  影子节点被回收的次数
    - **`pgrefill`**    扫描页面的数量（在活动 LRU 列表中）
    - **`pgscan`**  扫描页面的数量（在非活动 LRU 列表中）
    - **`pgsteal`** 回收的页面数量
    - **`pgactivate`**  移动到活动 LRU 列表的页面数量
    - **`pgdeactivate`**    移动到非活动 LRU 列表的页面数量
    - **`pglazyfree`**  在内存压力下推迟释放的页面数量
    - **`pglazyfreed`** 回收的lazyfree页面数量
- **`memory.swap.current`** 存在于非根 cgroup 上的只读单值文件。
    cgroup 及其后代当前使用的交换总量。

- **`memory.swap.max`** 存在于非根 cgroup 上的读写单值文件。 默认值为 `max`。
    交换使用硬限制。 如果cgroup的交换使用量达到此限制，则该cgroup的匿名内存将不会被换出。


### 5.2.2 Usage Guidelines

`memory.high` 是控制内存使用的主要机制。
过度承诺上限（上限总和>可用内存）并让全局内存压力根据使用情况分配内存是一个可行的策略。

由于违反上限不会触发 OOM Killer，而是会限制违规 cgroup，因此管理代理有充足的机会来监视并采取适当的操作，例如授予更多内存或终止工作负载。

确定 cgroup 是否有足够的内存并非易事，因为内存使用情况并不表明工作负载是否可以从更多内存中受益。
例如，将从网络接收的数据写入文件的工作负载可以使用所有可用内存，但也可以使用少量内存进行高性能操作。
衡量内存压力（由于内存不足而影响工作负载的程度）对于确定工作负载是否需要更多内存是必要的； 不幸的是，内存压力监控机制尚未实现。


### 5.2.3 Memory Ownership

内存区域被实例化它的 cgroup 占用，并保持被 cgroup 占用，直到该区域被释放。
将进程迁移到不同的 cgroup 不会将其在前一个 cgroup 中实例化的内存使用量移动到新的 cgroup。

内存区域可以由属于不同 cgroup 的进程使用。
该区域将被计入哪个 cgroup 是不确定的； 然而，随着时间的推移，内存区域很可能最终出现在一个有足够内存空间以避免高回收压力的 cgroup 中。

如果一个 cgroup 扫描了大量内存，并且预计会被其他 cgroup 重复访问，则使用 `POSIX_FADV_DONTNEED` 放弃属于受影响文件的内存区域的所有权以确保正确的内存所有权可能是有意义的。


## 5.3 IO

`io` 控制器调节IO资源的分配。
该控制器实现基于权重和绝对带宽或 IOPS 限制分配； 
但是，仅当使用 `cfq-iosched` 时，基于权重的分配才可用，并且这两种方案均不适用于 `blk-mq` 设备。


### 5.3.1 IO Interface Files

- **`io.stat`** 存在于非 root cgroup 上的只读嵌套键控文件。
    线路由 `$MAJ:$MIN` 设备编号键入，并且不排序。
    定义了以下嵌套键。
    ```
    rbytes      Byted read
    wbytes      Bytes Written
    rios        Number of read IOs
    wios        Number of write IOs
    ```

    一个读取输出示例如下：

    ```
    8:16 rbytes=1459200 wbytes=314773504 rios=192 wios=353
    8:0 rbytes=90430464 wbytes=299008000 rios=8950 wios=1252
    ```

- **`io.weight`**   存在于非根 cgroup 上的读写平键文件。
    默认值为 `default 100`。

    第一行是应用于没有特定覆盖的设备的默认权重。 其余的都是由 `$MAJ:$MIN` 设备编号键入的覆盖，并且未排序。 权重范围为 [1, 10000]，指定 cgroup 与其同级组相比可以使用的相对 IO 时间量。

    可以通过写入 `default $WEIGHT` 或简单地 `$WEIGHT` 来更新默认权重。
    可以通过写入`$MAJ:$MIN $WEIGHT`来设置覆盖，并通过写入`$MAJ:$MIN default` 来取消设置。

    一个读取输出的示例如下：

    ```
    default 100
    8:16 200
    8:0 50
    ```

- **`io.max`**  存在于非 root cgroup 上的读写嵌套键控文件
    基于 BPS 和 IOPS 的 IO 限制。 
    线路由 `$MAJ:$MIN` 设备编号键入，并且不排序。 
    定义了以下嵌套键。

    ```
    rbps        Max read bytes per second
    wbps        Max write bytes per second
    riops       Max read IO operations per second
    wiops       Max write IO operations per second
    ```

    写入时，可以按任意顺序指定任意数量的嵌套键值对。
    可以将 `max` 指定为删除特定限制的值。 
    如果多次指定相同的键，则结果不确定。

    BPS 和 IOPS 在每个 IO 方向上进行测量，如果达到限制，IO 就会延迟。
    允许临时突发。

    将读取限制设置为 2M BPS，写入限制为 120 IOPS，持续 8:16:

    `echo 8:16 rbps=2097152 wiops=120" > io.max`

    读取返回下值：

    `8:16 rbps=2097152 wbps=max riops=max wiops=120`

    写入 IOPS 限制可以通过写入以下内容来删除：

    `echo "8:16 wiops=max" > io.max`

    现在读取会返回如下：

    `8:16 rbps=2097152 wbps=max riops=max wiops=max`


### 5.3.2 Writeback

页面缓存通过缓冲写入和共享 mmap 被弄脏，并通过写回机制异步写入到支持文件系统。
Writeback 位于内存和IO域之间，通过平衡脏数据和写IO来调节脏内存的比例。

io控制器与内存控制器配合，实现对页缓存写回IO的控制。
内存控制器定义计算和维护脏内存比率的内存域，io控制器定义为内存域写出脏页的io域。
系统范围和每个 cgroup 的脏内存状态都会被检查，并强制执行两者中更严格的状态。

cgroup 写回需要底层文件系统的显式支持。
目前，cgroup writeback 在 `ext2`、`ext4` 和 `btrfs` 上实现。
在其他文件系统上，所有写回 IO 都归属于根 cgroup。

内存和写回管理存在固有的差异，这会影响 cgroup 所有权的跟踪方式。
内存按页进行跟踪，而写回按索引节点进行。 出于写回的目的，一个 inode 被分配给一个 cgroup，所有从该 inode 写入脏页的 IO 请求都归属于该 cgroup。

由于内存的 cgroup 所有权是按页跟踪的，因此可能存在与与 inode 关联的 cgroup 不同的 cgroup 关联的页面。
这些被称为外部页面。 写回会不断跟踪外部页面，如果特定的外部 cgroup 在一段时间内成为多数，则将 inode 的所有权切换到该 cgroup。

虽然此模型足以满足大多数用例，其中给定 inode 大部分被单个 cgroup 弄脏，即使主要写入 cgroup 随着时间的推移而变化，但不能很好地支持多个 cgroup 同时写入单个 inode 的用例。
在这种情况下，很大一部分 IO 可能会被错误归因。
由于内存控制器在第一次使用时分配页面所有权，并且在释放页面之前不会更新它，即使回写严格遵循页面所有权，多个 cgroup 弄脏重叠区域也无法按预期工作。
建议避免此类使用模式。

影响写回行为的 sysctl 旋钮应用于 cgroup 写回，如下所示。

- **`vm.dirty_background_ratio, vm.dirty_ratio`**
    这些比率同样适用于 cgroup 写回，可用内存量受内存控制器和系统范围的干净内存施加的限制。

- **`vm.dirty_background_bytes, vm.dirty_bytes`**
    对于 cgroup 写回，这被计算为与总可用内存的比率，并以与 `vm.dirty[_background]_ratio` 相同的方式应用。


## 5.4 PID

进程号控制器用于允许 cgroup 在达到指定限制后停止任何新任务的 `fork()` 或 `clone()` 操作。

cgroup 中的任务数量可能会以其他控制器无法阻止的方式耗尽，因此需要有自己的控制器。
例如，分叉炸弹可能会在达到内存限制之前耗尽任务数量。

请注意，此控制器中使用的 PID 指的是内核使用的 TID、进程 ID。


### 5.4.1 PID Interface Files

- **`pids.max`**    存在于非根 cgroup 上的读写单值文件。 
    默认值为 `max`。

    进程数量的硬限制。

- **`pids.current`**    存在于所有 cgroup 中只读的单一值文件。
    cgroup 及其后代中当前进程的数量。

组织操作不会被 cgroup 策略阻止，因此 `pids.current > pids.max` 是可能的。
这可以通过将限制设置为小于 `pids.current` 或将足够的进程附加到 cgroup 以使 `pids.current` 大于 `pids.max` 来完成。
但是，不可能通过 `fork()` 或 `clone()` 违反 cgroup PID 策略。
如果创建新进程会导致违反 cgroup 策略，这些将返回 -EAGAIN。


## 5.5 Device controller

设备控制器管理对设备文件的访问。 
它包括创建新设备文件（使用 `mknod`）以及访问现有设备文件。

Cgroup v2 设备控制器没有接口文件，并且在 cgroup BPF 之上实现。
为了控制对设备文件的访问，用户可以创建 `BPF_CGROUP_DEVICE` 类型的 `bpf` 程序并将它们附加到 cgroup。
尝试访问设备文件时，将执行相应的 BPF 程序，并且根据返回值，尝试将成功或失败并显示 -EPERM。

`BPF_CGROUP_DEVICE` 程序采用指向 `bpf_cgroup_dev_ctx` 结构的指针，该结构描述了设备访问尝试：访问类型（`mknod`/`read`/`write`）和设备（类型、主设备号和次设备号）。
如果程序返回 0，则尝试失败并返回 -EPERM，否则成功。

`BPF_CGROUP_DEVICE` 程序的示例可以在 `tools/testing/selftests/bpf/dev_cgroup.c` 文件的内核源代码树中找到。


## 5.6 RDMA

`rdma` 控制器调节 RDMA 资源的分配和记账。


### 5.6.1 RDMA Interface Files

- **`rdma.max`**    除 root 之外的所有 cgroup 都存在的读写嵌套键控文件，用于描述 RDMA/IB 设备当前配置的资源限制。

    行按设备名称键入，并且不排序。
    每行包含空格分隔的资源名称及其可分发的配置限制。

    定义了以下嵌套键。

    ```
    hca_handle      Maximum number of HCA Handles
    hca_object      Maximum number of HCA Objects
    ```

    mlx4 和 ocrdma 设备的示例如下：

    ```
    mlx4_0 hca_handle=2 hca_object=2000
    ocrdma1 hca_handle=3 hca_object=max
    ```

- **`rdma.current`**    描述当前资源使用情况的只读文件。
    存在于除了 root 之外的所有 cgroup 中。

    mlx4 和 ocrdma 设备的示例如下：

    ```
    mlx4_0 hca_handle=1 hca_object=20
    ocrdma1 hca_handle=1 hca_object=23
    ```


## 5.7 Misc


### 5.7.1 `perf_event`

`perf_event` 控制器如果未安装在旧层次结构上，则会在 v2 层次结构上自动启用，以便始终可以通过 cgroup v2 路径过滤 perf 事件。
填充 v2 层次结构后，控制器仍然可以移动到旧层次结构。


## 5.N Non-normative information

本节包含不被视为稳定内核 API 一部分的信息，因此可能会发生更改。


### 5.N.1 CPU controller root cgroup process behaviour

在根 cgroup 中分配 CPU 周期时，该 cgroup 中的每个线程都被视为托管在根 cgroup 的单独子 cgroup 中。
该子 cgroup 的权重取决于其线程的良好级别。

有关此映射的详细信息，请参阅 `kernel/sched/core.c` 文件中的 `sched_prio_to_weight` 数组（该数组中的值应适当缩放，以便中性 - nice 0 - 值为 100 而不是 1024）。


### 5.N.2 IO controller root cgroup process behaviour

根 cgroup 进程托管在隐式叶子节点中。
分配 IO 资源时，会考虑此隐式子节点，就好像它是权重值为 200 的根 cgroup 的普通子 cgroup 一样。


# 6. Namespace


## 6.1 Basics

cgroup 命名空间提供了一种虚拟化 `/proc/$PID/cgroup` 文件和 cgroup 挂载视图的机制。
`CLONE_NEWCGROUP` 克隆标志可以与 `clone(2)` 和 `unshare(2)` 一起使用来创建新的 cgroup 命名空间。
在 cgroup 命名空间内运行的进程将其 `/proc/$PID/cgroup` 输出限制为 cgroupns 根目录。
cgroupns 根是创建 cgroup 命名空间时进程的 cgroup。

如果没有 cgroup 命名空间，`/proc/$PID/cgroup` 文件显示进程 cgroup 的完整路径。
在一组 cgroup 和命名空间旨在隔离进程的容器设置中，`/proc/$PID/cgroup` 文件可能会将潜在的系统级信息泄漏给隔离的进程。 例如：

```bash
cat /proc/self/cgroup

0::/batchjobs/container_id1
```

路径 `/batchjobs/container_id1` 可以被视为系统数据，并且不希望暴露给隔离的进程。
cgroup 命名空间可用于限制此路径的可见性。
例如，在创建 cgroup 命名空间之前，我们会看到：

```bash
ls -l /proc/self/ns/cgroup
lrwxrwxrwx 1 root root 0 2014-07-15 10:37 /proc/self/ns/cgroup -> cgroup:[4026531835]

cat /proc/self/cgroup
0::/batchjobs/container_id1
```

取消共享新命名空间后，视图会发生变化：

```bash
ls -l /proc/self/ns/cgroup
lrwxrwxrwx 1 root root 0 2014-07-15 10:35 /proc/self/ns/cgroup -> cgroup:[4026532183]

cat /proc/self/cgroup
0::/
```

当多线程进程中的某个线程取消共享其 cgroup 命名空间时，新的 cgroupns 将应用于整个进程（所有线程）。
这对于 v2 层次结构来说是很自然的； 然而，对于遗留层次结构来说，这可能是意想不到的。

只要内部有进程或挂载固定它，cgroup 命名空间就处于活动状态。
当最后一次使用消失时，cgroup 命名空间将被销毁。
cgroupns 根和实际的 cgroup 仍然存在。


## 6.2 The Root and Views

cgroup 命名空间的 `cgroupns root` 是调用 `unshare(2)` 的进程正在其中运行的 cgroup。
例如，如果 `/batchjobs/container_id1` cgroup 中的进程调用 `unshare`，则 cgroup `/batchjobs/container_id1` 将成为 cgroupns 根。
对于 `init_cgroup_ns`，这是真正的根（'/'）cgroup。

即使名称空间创建者进程稍后移动到不同的 cgroup，cgroupns 根 cgroup 也不会更改：

```bash
~/unshare -c  # unshare cgroupns in some cgroup
cat /proc/self/cgroup
0::/

mkdir sub_cgrp_1
echo 0 > sub_cgrp_1/cgroup.procs

cat /proc/self/cgroup
0::/sub_cgrp_1
```
每个进程都会获取其名称空间特定的 `/proc/$PID/cgroup` 视图

在 cgroup 命名空间内运行的进程将只能在其根 cgroup 内看到 cgroup 路径（在 `/proc/self/cgroup` 中）。
来自非共享 cgroupns:

```bash
sleep 100000 &
[1] 7353

echo 7353 > sub_cgrp_1/cgroup.procs
cat /proc/7353/cgroup
0::/sub_cgrp_1
```

从初始的 cgroup 命名空间中，真正的 cgroup 路径将是可见的:

```bash
cat /proc/7353/cgroup
0::/batchjobs/container_id1/sub_cgrp_1
```

从同级 cgroup 命名空间（即以不同 cgroup 为根的命名空间），将显示相对于其自己的 cgroup 命名空间根的 cgroup 路径。
例如，如果 PID 7353 的 cgroup 命名空间根位于 `/batchjobs/container_id2`，那么它将看到：

```bash
cat /proc/7353/cgroup
0::/../container_id2/sub_cgrp_1
```

请注意，相对路径始终以 `/` 开头，表示它相对于调用者的 cgroup 命名空间根。


## 6.3 Migration and setns(2)

如果 cgroup 命名空间内的进程具有对外部 cgroup 的适当访问权限，则它们可以移入和移出命名空间根。
例如，从 cgroupns 根位于 `/batchjobs/container_id1` 的命名空间内部，并假设全局层次结构仍然可以在 cgroupns 内访问：

```bash
cat /proc/7353/cgroup
0::/sub_cgrp_1

echo 7353 > batchjobs/container_id2/cgroup.procs

cat /proc/7353/cgroup
0::/../container_id2
```

请注意，不鼓励这种设置。 
cgroup 命名空间内的任务只能暴露给它自己的 cgroupns 层次结构。

在以下情况下允许 `setns(2)` 到另一个 cgroup 命名空间：

1. 该进程对其当前用户命名空间具有 `CAP_SYS_ADMIN`
2. 该进程具有针对目标 cgroup 命名空间的用户名的 `CAP_SYS_ADMIN`

连接到另一个 cgroup 命名空间时不会发生隐式 cgroup 更改。
预计某人会将附加进程移动到目标 cgroup 命名空间根下。


## 6.4 Interaction with Other Namespaces

命名空间特定的 cgroup 层次结构可以由在非 init cgroup 命名空间内运行的进程挂载：

`mount -t cgroup2 none $MOUNT_POINT`

这将挂载统一的 cgroup 层次结构，并将 cgroupns 根作为文件系统根。
该进程需要针对其用户和安装命名空间的 `CAP_SYS_ADMIN`。

`/proc/self/cgroup` 文件的虚拟化与通过命名空间私有 cgroupfs 挂载限制 cgroup 层次结构视图相结合，在容器内提供了正确隔离的 cgroup 视图。


# P. Information on Kernel Programming

本节包含需要与 cgroup 交互的区域的内核编程信息。
cgroup 核心和控制器不包括在内。


## P.1 Filesystem Support for Writeback

文件系统可以通过更新 `address_space_operations->writepage[s]()` 来支持 cgroup 写回，以使用以下两个函数注释 bio。

- **`wbc_init_bio(@wbc, @bio)`**
    应该为每个携带写回数据的 Bio 调用，并将 Bio 与 inode 的所有者 cgroup 相关联。
    可以在 bio 分配和提交之间随时调用。

- **`wbc_account_io(@wbc, @page, @bytes)`**
    应该为每个被写出的数据段调用。
    虽然此函数并不关心在写回会话期间何时调用它，但在将数据段添加到 Bio 时调用它是最简单、最自然的。

通过写回 Bio 的注释，可以通过在 `->s_iflags` 中设置 `SB_I_CGROUPWB` 来启用每个 `super_block` 的 cgroup 支持。
这允许有选择地禁用 cgroup 写回支持，这在某些文件系统功能（例如 日志数据模式，不兼容。

`wbc_init_bio()` 将指定的 bio 绑定到它的 cgroup。 
根据配置，bio 可能会以较低的优先级执行，并且如果写回会话持有共享资源（例如日志条目），可能会导致优先级反转。
对于这个问题，没有一种简单的解决方案。
文件系统可以尝试通过跳过 `wbc_init_bio()` 或直接使用 `bio_associate_blkcg()` 来解决特定问题。


# D. Deprecated v1 Core Features

- 不支持多个层次结构（包括命名层次结构）。
- 不支持所有 v1 挂载选项。
- `tasks` 文件被删除，`cgroup.procs` 未排序。
- `cgroup.clone_children` 已删除。
- `/proc/cgroups` 对于 v2 来说没有意义。 
    请改用根目录下的 `cgroup.controllers` 文件。


# R. Issues with v1 and Rationales for v2


## R.1 Multiple Hierarchies

cgroup v1 允许任意数量的层次结构，每个层次结构可以托管任意数量的控制器。
虽然这似乎提供了高度的灵活性，但在实践中并没有什么用处。

例如，由于每个控制器只有一个实例，因此可在所有层次结构中使用的诸如 `freezer` 之类的实用型控制器只能在一个层次结构中使用。
一旦填充了层次结构，控制器就无法移动到另一个层次结构，这一事实加剧了这个问题。
另一个问题是绑定到层次结构的所有控制器都被迫具有完全相同的层次结构视图。
无法根据特定控制器来改变粒度。

实际上，这些问题严重限制了哪些控制器可以放置在同一层次结构中，并且大多数配置都诉诸于将每个控制器放置在其自己的层次结构中。
只有紧密相关的控制器（例如 `cpu` 和 `cpuacct` 控制器）才有意义放在同一层次结构中。
这通常意味着用户空间最终会管理多个相似的层次结构，每当需要层次结构管理操作时，就会在每个层次结构上重复相同的步骤。

此外，对多个层次结构的支持需要付出高昂的代价。
它极大地复杂了 cgroup 核心实现，但更重要的是，对多个层次结构的支持限制了 cgroup 的一般使用方式以及控制器能够执行的操作。

可能存在的层次结构数量没有限制，这意味着线程的 cgroup 成员身份无法以有限长度进行描述。
密钥可能包含任意数量的条目并且长度不受限制，这使得操作非常困难，并导致添加仅用于识别成员身份的控制器，这反过来又加剧了层次结构数量激增的原始问题。

此外，由于控制器不能对其他控制器可能所在的层次结构的拓扑有任何期望，因此每个控制器必须假设所有其他控制器都附加到完全正交的层次结构。
这使得控制器之间的协作变得不可能，或者至少非常麻烦。

在大多数用例中，没有必要将控制器放在彼此完全正交的层次结构上。
通常需要的是能够根据特定的控制器具有不同级别的粒度。
换句话说，当从特定控制器查看时，层次结构可能会从叶向根折叠。
例如，给定的配置可能不关心超出特定级别的内存如何分配，但仍希望控制 CPU 周期的分配方式。


## R.2 Thread Granularity

cgroup v1 允许进程的线程属于不同的 cgroup。
这对于某些控制器来说没有意义，并且这些控制器最终实现了不同的方式来忽略此类情况，但更重要的是，它模糊了暴露给单个应用程序的 API 和系统管理接口之间的界限。

一般来说，进程内知识仅适用于进程本身； 因此，与进程的服务级组织不同，对进程的线程进行分类需要拥有目标进程的应用程序的积极参与。

cgroup v1 有一个定义不明确的委托模型，该模型与线程粒度结合起来被滥用。
cgroup 被委托给各个应用程序，以便它们可以创建和管理自己的子层次结构并控制它们的资源分配。
这有效地将 cgroup 提升为向非专业程序公开的类似系统调用的 API 的地位。

首先，cgroup 的接口根本不足以以这种方式公开。
对于要访问自己的旋钮的进程，它必须从 `/proc/self/cgroup` 中提取目标层次结构上的路径，通过将旋钮的名称附加到路径来构造路径，打开然后读取和/或写入它。
这不仅极其笨重和不寻常，而且本质上很活泼。
没有传统的方法来定义跨所需步骤的事务，并且无法保证流程实际上在其自己的子层次结构上运行。

cgroup 控制器实现了许多永远不会被接受为公共 API 的旋钮，因为它们只是向系统管理伪文件系统添加控制旋钮。
cgroup 最终得到了没有正确抽象或细化的接口旋钮，并且直接揭示了内核内部细节。
这些旋钮通过定义不明确的委托机制暴露给各个应用程序，有效地滥用 cgroup 作为实现公共 API 的捷径，而无需经过所需的审查。

这对于用户态和内核来说都是痛苦的。
用户态最终会出现行为不当和抽象不良的接口，并且内核会无意中暴露并锁定到构造中。


## R.3 Competition Between Inner Nodes and Threads

cgroup v1 允许线程位于任何 cgroup 中，这产生了一个有趣的问题，即属于父 cgroup 及其子 cgroup 的线程竞争资源。
这是令人讨厌的，因为两种不同类型的实体之间存在竞争，并且没有明显的方法来解决它。
不同的控制器做了不同的事情。

CPU 控制器将线程和 cgroup 视为等效项，并将良好级别映射到 cgroup 权重。
这在某些情况下有效，但当子级们想要分配特定比率的 CPU 周期和内部线程数量波动时，这种方法就会失败——这些比率随着竞争实体数量的波动而不断变化。
还有其他问题。
从良好级别到权重的映射并不明显或通用，并且还有各种其他旋钮根本不适用于线程。

io 控制器隐式地为每个 cgroup 创建一个隐藏的叶节点来托管线程。
隐藏的叶子有自己的所有旋钮的副本，并带有 `leaf_` 前缀。
虽然这允许对内部线程进行同等的控制，但它有严重的缺点。
它总是添加一个额外的嵌套层，否则就没有必要，使界面变得混乱，并使实现变得非常复杂。

内存控制器无法控制内部任务和子 cgroup 之间发生的情况，并且行为也没有明确定义。
有人尝试添加临时行为和旋钮来根据特定工作负载定制行为，但从长远来看，这会导致问题极难解决。

多个控制者都在努力处理内部任务，并想出了不同的方法来处理它； 不幸的是，所有方法都存在严重缺陷，而且，广泛不同的行为使得 cgroup 作为一个整体高度不一致。

这显然是一个需要从 cgroup 核心以统一方式解决的问题。


## R.4 Other Interface Issues

cgroup v1 在没有监督的情况下成长并产生了大量的特性和不一致。
cgroup 核心方面的一个问题是如何通知空 cgroup——为每个事件分叉并执行一个用户态辅助二进制文件。
事件传递不是递归的或可委托的。
该机制的局限性还导致内核内的事件传递过滤机制使接口进一步复杂化。

控制器接口也存在问题。 
一个极端的例子是控制器完全忽略层次结构并将所有 cgroup 视为直接位于根 cgroup 下。
一些控制器向用户空间暴露了大量不一致的实现细节。

控制器之间也没有一致性。
创建新的 cgroup 时，某些控制器默认不施加额外限制，而其他控制器则在明确配置之前不允许使用任何资源。
相同类型的控件的配置旋钮使用了截然不同的命名方案和格式。
统计和信息旋钮被任意命名，即使在同一控制器中也使用不同的格式和单位。

cgroup v2 在适当的情况下建立通用约定并更新控制器，以便它们公开最少且一致的接口。


## R.5 Controller Issues and Remedies


### R.5.1 Memory

原始下限（软限制）被定义为默认未设置的限制。
因此，全局回收首选的 cgroup 集是选择加入，而不是选择退出。
优化这些大多是负面查找的成本是如此之高，以至于尽管其规模巨大，但其实现甚至无法提供基本的理想行为。
首先，软限制没有等级意义。
所有配置的组都组织在全局 rbtree 中，并被视为平等的对等体，无论它们位于层次结构中的哪个位置。
这使得子树委托不可能。 其次，软限制回收过程非常激进，不仅会给系统带来较高的分配延迟，还会因过度回收而影响系统性能，甚至导致该功能弄巧成拙。

另一方面，`memory.low` 边界是自上而下分配的保留。
当 cgroup 及其所有祖先都低于其低边界时，它享有回收保护，这使得子树的委派成为可能。
其次，新的 cgroup 没有默认储备，并且在常见情况下，大多数 cgroup 都有资格获得首选回收通行证。
这使得新的低边界只需对通用回收代码进行少量添加即可有效实现，而不需要带外数据结构和回收通道。
由于通用回收代码会考虑除首选第一回收过程中运行速度较低的 cgroup 之外的所有 cgroup，因此也会消除各个组的过度回收，从而实现更好的整体工作负载性能。

最初的 `high` 边界，即硬限制，被定义为一个不能移动的严格限制，即使必须调用 OOM Killer。
但这通常违背了充分利用可用内存的目标。
工作负载的内存消耗在运行时会发生变化，这需要用户过量使用。
但要在严格的上限下做到这一点，需要对工作集大小进行相当准确的预测，或者在极限上增加松弛度。
由于工作集大小估计很困难且容易出错，并且错误会导致 OOM 终止，因此大多数用户倾向于选择更宽松的限制，最终浪费宝贵的资源。

另一方面，`memory.high` 边界可以设置得更加保守。
当受到攻击时，它会通过强制直接回收来消除多余的分配来限制分配，但它永远不会调用 OOM 杀手。
因此，过于激进地选择高边界不会终止进程，反而会导致性能逐渐下降。
用户可以对此进行监控并进行纠正，直到找到仍能提供可接受性能的最小内存占用量。

在极端情况下，由于组内有许多并发分配和回收进度完全崩溃，因此可能会超出上限。
但即便如此，从其他组或系统其余部分中可用的闲置资源中满足分配也比杀死该组要好。
否则，`memory.max` 会限制这种类型的溢出，并最终包含有错误甚至恶意的应用程序。

将原始 `memory.limit_in_bytes` 设置为低于当前使用量会受到竞争条件的影响，其中并发开销可能会导致限制设置失败。
另一方面，`memory.max` 将首先设置限制以防止新的费用，然后回收并 OOM 终止，直到满足新的限制 - 或者写入 `memory.max` 的任务被终止。

组合的 `memory+swap` 计算和限制被对交换空间的实际控制所取代。

原始 cgroup 设计中组合 `memory+swap` 设施的主要论点是，全局或父级压力始终能够交换子组的所有匿名内存，无论子组自己的（可能不受信任的）配置如何。
但是，不受信任的组可以通过其他方式破坏交换 - 例如在紧密循环中引用其匿名内存 - 并且管理员在过度提交不受信任的作业时无法假设完全可交换性。

另一方面，对于受信任的作业，组合计数器不是直观的用户空间界面，并且它违背了 cgroup 控制器应该考虑和限制特定物理资源的想法。
交换空间是一种与系统中所有其他资源一样的资源，这就是统一层次结构允许单独分配它的原因。

