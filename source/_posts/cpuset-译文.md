---
title: cpuset 译文
date: 2023-12-26 16:58:00
tags:
- doc
categories:
- kernel doc
- cpusets
keywords:
- cpusets
copyright: Guader
copyright_author_href:
copyright_info:
---

[原文](https://www.kernel.org/doc/Documentation/cgroup-v1/cpusets.txt)

## Cpusets

### 什么是 cpusets

Cpusets 提供了一种将一组 CPU 和内存节点分配给一组任务的机制。在本文中，"内存节点 "是指包含内存的联机节点。

Cpuset 将任务的 CPU 和内存位置限制在任务当前 cpuset 的资源范围内。 它们形成了虚拟文件系统中可见的嵌套层次结构。这些都是在大型系统上管理动态任务分配所需的基本钩子，超出了已有的钩子范围。

Cpusets 使用文档[cgroups.txt](https://www.kernel.org/doc/Documentation/cgroup-v1/cgroups.txt) 中描述的通用 cgroup 子系统。

使用 `sched_setaffinity(2)` 系统调用将 CPU 纳入 CPU 亲和性掩码(CPU affinity mask)，以及使用 `mbind(2)` 和 `set_mempolicy(2)` 系统调用将内存节点纳入内存策略的任务请求，都会通过该任务的 cpuset 过滤，过滤掉不在该 cpuset 中的 CPU 或内存节点。
调度程序不会在 CPU 的 `cpus_allowed` 向量中不允许的 CPU 上调度任务，内核页分配器也不会在请求任务的 `mems_allowed` 向量中不允许的节点上分配页面。

用户级代码可以在 cgroup 虚拟文件系统中按名称创建和销毁 cpuset，管理这些 cpuset 的属性和权限以及为每个 cpuset 分配哪些 CPU 和内存节点，指定和查询任务分配到哪个 cpuset，并列出分配给 cpuset 的任务 pid。


### 为什么需要 cpusets

具有许多处理器 (CPUs)、复杂的内存缓存层次结构和具有非均匀访问时间 (NUMA, non-uniform access times) 的多个内存节点的大型计算机系统的管理给进程的高效调度和内存布局带来了额外的挑战。

通常，只需让操作系统在请求任务间自动共享可用的 CPU 和内存资源，就可以以足够的效率运行较小规模的系统。

但是，较大的系统可以从细致的处理器和内存布局中获益更多，以减少内存访问时间和争用，并且通常代表客户的更大投资，可以通过将工作明确地放置在系统的适当大小的子集上而受益。

这对于以下方面尤其有价值：

- 运行同一 Web 应用程序多个实例的 Web 服务器，
- 运行不同应用的服务器（例如，Web 服务器和数据库），或者
- 运行具有要求性能特征的大型 HPC 应用的 NUMA 系统。

这些子集或 “软分区(soft partitions)” 必须能够随着工作组合的变化而动态调整，而不影响其他并发执行的工作。当内存位置改变时，正在运行的工作页面的位置也可能被移动。

内核 cpuset 补丁提供了有效实现此类子集所需的最低限度的基本内核机制。它利用 Linux 内核中现有的 CPU 和内存布局设施，以避免对关键调度程序或内存分配器代码产生任何额外影响。


### cpusets 是如何实现的

Cpuset 提供了一种 Linux 内核机制来限制一个进程或一组进程使用哪些 CPU 和内存节点。

Linux 内核已经有一对机制来指定任务可以在哪些 CPU 上调度（`sched_setaffinity`）以及任务可以在哪些内存节点上获取内存（`mbind、set_mempolicy`）。

Cpuset 对这两种机制进行了如下扩展：

- Cpuset 是内核已知的，允许的 CPU 和内存节点的集合。
- 系统中的每个任务都通过任务结构中指向引用计数 cgroup 结构的指针附加到一个 cpuset。
- 对 `sched_setaffinity` 的调用将被过滤到该任务的 cpuset 中允许的那些 CPU。
- 对 `mbind` 和 `set_mempolicy` 的调用将被过滤到该任务的 cpuset 中允许的那些内存节点。
- 根 cpuset 包含所有系统 CPU 和内存节点。
- 对于任何 cpuset，都可以定义包含父 CPU 和内存节点资源的子集的子 cpuset。
- cpuset 的层次结构可以挂载在 `/dev/cpuset` 中，以便从用户空间进行浏览和操作。
- cpuset可以被标记为独占的，这确保没有其他cpuset（除了直接祖先和后代）可以包含任何重叠的CPU或内存节点。
- 您可以列出附加到任何 cpuset 上的所有任务（按 pid ）。


Cpuset 的实现需要一些简单的钩子到内核的其余部分，而不是性能关键路径：

- 在 `init/main.c` 中，在系统启动时初始化根 cpuset。
- 在 `fork` 和 `exit` 中，从其 cpuset 附加和分离任务。
- 在 `sched_setaffinity` 中，根据该任务的 cpuset 中允许的内容来屏蔽请求的 CPU。
- 在 `sched.c migrate_live_tasks()` 中，如果可能的话，在其 cpuset 允许的 CPU 内保留迁移任务。
- 在 `mbind` 和 `set_mempolicy` 系统调用中，根据该任务的 cpuset 中允许的内容来屏蔽所请求的内存节点。
- 在 `page_alloc.c` 中，将内存限制到允许的节点。
- 在 `vmscan.c` 中，将页面恢复限制为当前 cpuset。


您应该挂载“cgroup”文件系统类型，以便能够浏览和修改内核当前已知的 cpuset。 没有为 cpuset 添加新的系统调用 - 对查询和修改 cpuset 的所有支持都是通过此 cpuset 文件系统实现的。

每个任务的 `/proc/<pid>/status` 文件添加了四行，显示任务的 `cpus_allowed`（可以在哪些 CPU 上调度）和 `mems_allowed`（可以在哪些内存节点上获取内存），两种格式见下面的例子：

```
Cpus_allowed:   ffffffff,ffffffff,ffffffff,ffffffff
Cpus_allowed_list:      0-127
Mems_allowed:   ffffffff,ffffffff
Mems_allowed_list:      0-63
```

每个 cpuset 由 cgroup 文件系统中的一个文件夹表示，其中包含（在标准 cgroup 文件之上）描述该 cpuset 的以下文件：

- `cpuset.cpus`：该 cpuset 中的 CPU 列表
- `cpuset.mems`：该 cpuset 中的内存节点列表
- `cpuset.memory_migrate` 标志：如果设置，则将页面移动到 cpusets 节点
- `cpuset.cpu_exclusive` 标志：CPU 布局是否独占？
- `cpuset.mem_exclusive` 标志：内存布局是否独占？
- `cpuset.mem_hardwall` 标志：内存分配是否是硬墙的 hardwalled
- `cpuset.memory_pressure`：衡量cpuset中分页压力的大小
- `cpuset.memory_spread_page` 标志：如果设置，则在允许的节点上均匀分布页面缓存
- `cpuset.memory_spread_slab` 标志：如果设置，则将slab缓存均匀地分布在允许的节点上
- `cpuset.sched_load_balance` 标志：如果设置，则该 cpuset 上的 CPU 内的负载平衡
- `cpuset.sched_relax_domain_level`：迁移任务时的搜索范围


另外，只有根 cpuset 有以下文件：
- `cpuset.memory_pressure_enabled` 标志：计算内存压力

使用 mkdir 系统调用或 shell 命令创建新的 cpuset。 cpuset 的属性（例如其标志、允许的 CPU 和内存节点以及附加任务）可通过写入该 cpusets 文件夹中的相应文件来修改，如上所列。

嵌套 cpuset 的命名分层结构允许将大型系统划分为嵌套的、可动态更改的，“软分区”。

每个任务（由该任务的任何子任务在 fork 时自动继承）附加到 cpuset，允许将系统上的工作负载组织到相关的任务集中，这样每个任务集都被限制为使用特定 cpuset 的 CPU 和内存节点。如果必要的 cpuset 文件系统文件夹的权限允许，任务可以重新附加到任何其他 cpuset。

“大型”系统的这种管理与使用 `sched_setaffinity`、`mbind` 和 `set_mempolicy` 系统调用在各个任务和内存区域上完成的详细放置顺利集成。

以下规则适用于每个 cpuset：

- 它的 CPU 和内存节点必须是其父节点的子集。
- 除非它的父级是独占的，否则它不能被标记为独占。
- 如果它的 cpu 或内存是独占的，它们可能不会与任何同级重叠。


这些规则以及 cpuset 的自然层次结构可以有效执行独占保证，而无需在每次 cpuset 发生任何更改时扫描所有 cpuset，以确保没有任何内容与独占 cpuset 重叠。 此外，使用 Linux 虚拟文件系统 (vfs) 来表示 cpuset 层次结构为 cpuset 提供了熟悉的权限和名称空间，并且需要最少的附加内核代码。

根 (`top_cpuset`) cpuset 中的 cpus 和 mems 文件是只读的。cpus 文件使用 CPU 热插拔通知程序自动跟踪 `cpu_online_mask` 的值，而 mems 文件使用 `cpuset_track_online_nodes()` 钩子自动跟踪 `node_states[N_MEMORY]`（即具有内存的节点）的值。


### 什么是独占 cpusets

如果某个 cpuset 是 cpu 或 mem 独占的，则只有直接祖先或后代，可以共享任何相同的CPU或内存节点。

`cpuset.mem_exclusive` *或* `cpuset.mem_hardwall` 的 cpuset 是 "hardwalled"(硬墙？)。
即它限制内核对页面、缓冲区和内核在多个用户之间共享的其他数据的分配。
所有 cpuset，无论是否有硬墙，都会限制用户空间的内存分配。
这样可以配置系统，以便多个独立工作中可以共享公共内核数据，例如文件系统页，同时隔离每个工作在其自己的 cpuset 中的用户分配。
为此，请构造一个大型 `mem_exclusive cpuset` 来容纳所有工作，并为每个独立工作构造子，非 `mem_exclusive` 的 cpuset。
即使是 `mem_exclusive` cpuset，也只允许将少量典型内核内存，例如来自中断处理程序的请求，取出。



### 什么是内存压力 `memory_pressure` 

内存压力为每个 cpuset 提供了一个简单指标，用于衡量 cpuset 中的任务尝试释放 cpuset 节点上正在使用的内存以满足额外内存请求的速率。

这使得批处理管理器能够监控在专用 CPU 组中运行的工作，从而有效地检测该工作造成的内存压力级别。

这对于运行各种已提交作业的严格管理的系统很有用，这些系统可能会选择终止或重新确定尝试使用比分配给它们的节点上允许的内存更多的内存的作业，并且对于紧密耦合、长时间运行的情况，这非常有用。 大规模并行科学计算作业如果开始使用超过允许的内存，将严重无法满足所需的性能目标。

此机制为批处理管理器提供了一种非常经济的方式来监视 cpuset 的内存压力迹象。 由批次管理器或其他用户代码来决定如何处理并采取行动。

**除非通过向特殊文件 `/dev/cpuset/memory_pressure_enabled` 写入“1”来启用此功能，否则此指标的 `__alloc_pages()` 重新平衡代码中的钩子会简化为仅注意到 `cpuset_memory_pressure_enabled` 标志为零。 因此，只有启用此功能的系统才会计算该指标。**


为什么要计算每个 cpuset 的运行平均值：

由于此计量表是针对每个 cpuset 的，而不是针对每个任务或 mm，因此在大型系统上监视此指标的批处理调度程序所施加的系统负载会急剧减少，因为可以避免在每组查询上扫描任务列表。

由于该计量表是运行平均值，而不是累积计数器，因此批处理调度程序可以通过单次读取来检测内存压力，而不必读取并累积一段时间的结果。

由于此计量表是针对每个 cpuset 而不是针对每个任务或 mm，因此批处理调度程序可以通过单次读取来获取关键信息（cpuset 中的内存压力），而不必查询并累积所有（动态变化的）结果 cpuset 中的任务集。


每个 cpuset 的简单数字过滤器(simple digital filter)（需要一个自旋锁和每个 cpuset 3 个字的数据）将被保留，并由附加到该 cpuset 的任何任务更新，如果它进入同步（直接）页回收代码。

每个 cpuset 文件提供一个整数，表示最近由 cpuset 中的任务引起的直接页回收率（半衰期为 10 秒），以每秒尝试回收次数为单位，乘以 1000。



### 什么是内存分布 `memory spread`

每个 cpuset 有两个布尔标志文件，用于控制内核为文件系统缓冲区分配页的位置以及内核数据结构中的相关内容。它们被称为 `cpuset.memory_spread_page` 和 `cpuset.memory_spread_slab`。

如果设置了每个 cpuset 布尔标志文件 `cpuset.memory_spread_page`，则内核会将文件系统缓冲区（页缓存）均匀地分布在允许故障任务使用的所有节点上，而不是优先放置这些页到任务正在运行的节点上。

如果设置了每个 cpuset 布尔标志文件 `cpuset.memory_spread_slab`，那么内核将在允许故障任务使用的所有节点上均匀地分布一些与文件系统相关的 slab 缓存，例如索引节点和目录项，而不是更喜欢将这些页放在正在运行任务的节点上。

这些标志的设置不会影响任务的匿名数据段或堆栈段页。

默认情况下，这两种内存扩展都是关闭的，并且内存页分配在任务正在运行的本地节点上，除非任务的 NUMA mempolicy 或 cpuset 配置进行了修改，只要有足够的空闲内存页面可用。

创建新 cpuset 时，它们会继承其父代的内存分布设置。

设置内存扩展会导致受影响页或 slab 缓存的分配忽略任务的 NUMA 内存策略并进行扩展。使用 `mbind()` 或 `set_mempolicy()` 调用来设置 NUMA mempolicies 的任务不会注意到这些调用中的任何变化，因为它们包含任务的内存分布设置。如果关闭内存扩展，则当前指定的 NUMA mempolicy 再次应用于内存页面分配。

`cpuset.memory_spread_page` 和 `cpuset.memory_spread_slab` 都是布尔标志文件。 默认情况下，它们包含“0”，这意味着该 cpuset 的功能已关闭。 如果将“1”写入该文件，则会打开指定的功能。

实现很简单。

设置标志`cpuset.memory_spread_page`会为该 cpuset 中或随后加入该 cpuset 的每个任务打开每进程标志 `PFA_SPREAD_PAGE`。 页缓存的页分配调用被修改为对此 `PFA_SPREAD_PAGE` 任务标志执行内联检查，如果设置，对新例程 `cpuset_mem_spread_node()` 的调用将返回首选分配的节点。

类似地，设置`cpuset.memory_spread_slab`会打开标志 `PFA_SPREAD_SLAB`，并且适当标记的 slab 缓存将从`cpuset_mem_spread_node()`返回的节点分配页。

`cpuset_mem_spread_node()` 例程也很简单。 它使用每个任务转子 `cpuset_mem_spread_rotor` 的值来选择当前任务的 `mems_allowed` 中的下一个节点以进行分配。

这种内存放置策略也称为（在其他上下文中）循环(round-robin)或交错(interleave)。

此策略可以为需要将线程本地数据放置在相应节点上的工作提供实质性改进，但需要访问需要分布在工作 cpuset 中的多个节点上的大型文件系统数据集才能适应。 如果没有这一策略，特别是对于可能有一个线程读取数据集的工作，工作 cpuset 中节点之间的内存分配可能会变得非常不均匀。



### 什么是 `sched_load_balance`

内核调度程序（`kernel/sched/core.c`）自动对任务进行负载平衡。 如果一个 CPU 未得到充分利用，则在该 CPU 上运行的内核代码将在其他更过载的 CPU 上查找任务，并将这些任务在 cpusets 和 `sched_setaffinity` 等放置机制的限制内移动到自身。

负载平衡的算法成本及其对关键共享内核数据结构（例如任务列表）的影响随着平衡的 CPU 数量呈线性以上增长。因此，调度程序支持将系统 CPU 划分为多个调度域，以便仅在每个调度域内进行负载平衡。每个调度域覆盖系统中 CPU 的一些子集；两个调度域不重叠；某些 CPU 可能不在任何调度域中，因此不会进行负载平衡。

简而言之，在两个较小的调度域之间进行平衡的成本比在一个大的调度域之间进行平衡的成本要低，但这样做意味着两个域之一中的过载将不会负载平衡到另一个域。

默认情况下，有一个调度域覆盖所有 CPU，包括那些使用内核启动时间`isolcpus=`参数标记为隔离的 CPU。但是，隔离的 CPU 将不会参与负载平衡，并且除非明确分配，否则不会有任务在其上运行。

所有 CPU 之间的默认负载平衡不太适合以下两种情况：

1. 在大型系统上，跨多个 CPU 进行负载平衡的成本很高。 如果使用 cpuset 管理系统，将独立工作放置在不同的 CPU 组上，则不需要完全负载平衡。
2. 在某些 CPU 上支持实时的系统需要最大限度地减少这些 CPU 上的系统开销，包括避免任务负载平衡（如果不需要）。

当启用每个 cpuset 标志`cpuset.sched_load_balance`（默认设置）时，它会请求该 cpuset 中允许`cpuset.cpus`的所有 CPU 包含在单个调度域中，以确保负载平衡可以移动任务 （不以其他方式固定，如通过 `sched_setaffinity`）从该 cpuset 中的任何 CPU 到任何其他 CPU。

当每个 cpuset 标志`cpuset.sched_load_balance`被禁用时，调度程序将避免在该 cpuset 中的 CPU 之间进行负载平衡，除非有必要，因为某些重叠的 cpuset 启用了`sched_load_balance`。

因此，例如，如果顶级 cpuset 启用了`cpuset.sched_load_balance`标志，那么调度程序将有一个覆盖所有 CPU 的调度域，并且任何其他 cpuset 中的`cpuset.sched_load_balance`标志的设置并不重要 ，因为我们已经完全负载平衡了。

因此，在上述两种情况下，应该禁用顶级 cpuset 标志`cpuset.sched_load_balance`，并且只有一些较小的子 cpuset 启用此标志。

执行此操作时，您通常不希望在可能使用大量 CPU 的顶级 cpuset 中留下任何未固定的任务，因为此类任务可能会人为地限制到某些 CPU 子集，具体取决于此标志设置的详细信息 在后代 cpu 组中。 即使此类任务可以使用某些其他 CPU 中的空闲 CPU 周期，内核调度程序也可能不会考虑将该任务负载平衡到未充分利用的 CPU 的可能性。

当然，固定到特定 CPU 的任务可以保留在禁用`cpuset.sched_load_balance`的 cpuset 中，因为这些任务无论如何都不会去其他地方。

这里 cpuset 和调度域之间存在阻抗不匹配。Cpuset 是分层和嵌套的。调度域是平坦的；它们不重叠，并且每个 CPU 最多位于一个调度域中。

调度域必须是平坦的，因为部分重叠的 CPU 组之间的负载平衡会带来不稳定动态的风险，这超出了我们的理解。因此，如果两个部分重叠的 cpuset 中的每一个都启用标志`cpuset.sched_load_balance`，那么我们将形成一个单一的调度域，它是两者的超集。我们不会将任务移动到其 cpuset 之外的 CPU，但考虑到这种可能性，调度程序负载平衡代码可能会浪费一些计算周期。

这种不匹配就是为什么启用了`cpuset.sched_load_balance`标志的 cpuset 与调度域配置之间不存在简单的一对一关系的原因。如果某个 cpuset 启用该标志，它将在其所有 CPU 之间实现平衡，但如果禁用该标志，则只有在没有其他重叠 cpuset 启用该标志的情况下，才能确保没有负载平衡。

如果两个 cpuset 允许部分重叠的`cpuset.cpus`，并且只有其中一个启用了此标志，则另一个可能会发现其任务仅在重叠的 CPU 上部分负载平衡。 这只是上面几段给出的 top_cpuset 示例的一般情况。在一般情况下，如在 top cpuset 情况下，不要将可能使用大量 CPU 的任务留在此类部分负载平衡的 cpuset 中，因为它们可能被人为地限制为允许它们使用的某些 CPU 子集，例如 缺乏对其他 CPU 的负载平衡。

`cpuset.isolcpus`中的 CPU 被 `isolcpus=` 内核引导选项排除在负载平衡之外，并且无论任何 cpuset 中`cpuset.sched_load_balance` 的值如何，都永远不会进行负载平衡。



### `sched_setaffinity` 实现细节

每个 cpuset 标志`cpuset.sched_load_balance`默认为启用（与大多数 cpuset 标志相反）。当为某个 cpuset 启用时，内核将确保它可以在该 cpuset 中的所有 CPU 之间实现负载平衡（确保该 cpuset 的 `cpus_allowed` 的 CPU 位于同一调度域中。）

如果两个重叠的 cpuset 都启用了`cpuset.sched_load_balance`，那么它们将（必须）位于同一个调度域中。

如果按照默认情况，顶级 cpuset 启用了`cpuset.sched_load_balance`，那么根据上面的内容，这意味着有一个覆盖整个系统的调度域，而不管任何其他 cpuset 设置如何。

内核向用户空间承诺，它将尽可能避免负载平衡。 它将尽可能精细地选择调度域的粒度分区，同时仍然为启用了`cpuset.sched_load_balance`的 cpuset 允许的任何 CPU 集提供负载平衡。

内部内核 cpuset 到调度程序接口将系统中负载平衡 CPU 的分区从 cpuset 代码传递到调度程序代码。该分区是一组成对不相交的 CPU 子集（表示为 struct cpumask 数组），涵盖了必须进行负载平衡的所有 CPU。

cpuset 代码会构建一个新的此类分区，并将其传递给调度程序调度域设置代码，以便在以下情况下根据需要重建调度域：
- 具有非空 CPU 的 cpuset 的`cpuset.sched_load_balance`标志发生变化，
- 或者 CPU 来自启用此标志的 cpuset，
- 或具有非空 CPU 且启用此标志的 cpuset 的 `cpuset.sched_relax_domain_level` 值发生变化，
- 或者删除具有非空 CPU 并启用此标志的 cpuset，
- 或CPU 离线/在线。

该分区准确地定义了调度程序应设置的调度域 - 分区中的每个元素（struct cpumask）一个调度域。

调度程序会记住当前活动的调度域分区。 当从 cpuset 代码调用调度程序例程 `partition_sched_domains()` 来更新这些调度域时，它会将请求的新分区与当前分区进行比较，并更新其调度域，对于每次更改，删除旧的并添加新的。



### 什么是 `sched_relax_domain_level`

在调度域中，调度程序以两种方式迁移任务；周期性负载均衡和某些计划事件发生时定期进行负载平衡。

当任务被唤醒时，调度程序会尝试将任务移至空闲的 CPU 上。例如，如果在 CPU X 上运行的任务 A 激活了同一 CPU X 上的另一个任务 B，并且如果 CPU Y 是 X 的同级且处于空闲状态，则调度程序会将任务 B 迁移到 CPU Y，以便任务 B 可以在 CPU Y 上启动而无需 在 CPU X 上等待任务 A。

如果某个 CPU 用完其运行队列中的任务，该 CPU 会尝试从其他繁忙的 CPU 中拉出额外的任务来帮助它们，然后再进入空闲状态。

当然，寻找可移动任务和/或空闲 CPU 需要一些搜索成本，调度程序可能不会每次都搜索域中的所有 CPU。事实上，在某些架构中，事件的搜索范围仅限于 CPU 所在的同一个 socket 或节点，而周期性的负载均衡则搜索全部。

例如，假设 CPU Z 距离 CPU X 相对较远。即使 CPU Z 空闲而 CPU X 和同级 CPU 忙碌，调度程序也无法将唤醒的任务 B 从 X 迁移到 Z，因为它超出了其搜索范围。结果，CPU X 上的任务 B 需要等待任务 A 或等待下一个时钟周期的负载平衡。 对于某些特殊情况的应用程序，等待 1 个周期可能会太长。

`cpuset.sched_relax_domain_level`文件允许您根据需要请求更改此搜索范围。该文件采用 int 值，该值表示理想情况下级别搜索范围的大小，如下所示，否则初始值 -1 表示 cpuset 没有请求。

- -1：无请求。 使用系统默认或按照别人的要求。
-  0：不搜索。
-  1：搜索同级（核心中的超线程）。
-  2：搜索包中的核心。
-  3: 在节点中搜索 cpu [= 非 NUMA 系统上的系统范围]
-  4：在节点块中搜索节点[在NUMA系统上]
-  5：搜索系统范围[在 NUMA 系统上]

系统默认值取决于体系结构。 可以使用`relax_domain_level=`引导参数更改系统默认值。

该文件是针对每个 cpuset 的，影响 cpuset 所属的调度域。因此，如果 cpuset 的标志`cpuset.sched_load_balance`被禁用，则`cpuset.sched_relax_domain_level`无效，因为没有属于该 cpuset 的调度域。

如果多个 cpuset 重叠并因此形成单个调度域，则使用其中的最大值。请注意，如果一个请求 0，而其他请求 -1，则使用 0。

注意，修改这个文件会有好有坏的影响，是否可以接受取决于你的情况。如果您不确定，请勿修改此文件。

如果您的情况是：

- 由于您的特殊应用程序的行为或对 CPU 缓存的特殊硬件支持等，每个 cpu 之间的迁移成本可以假设相当小（对您来说）。
- 搜索成本不会对您产生影响，或者您可以通过管理 cpuset 进行压缩等来使搜索成本足够小。
- 即使牺牲缓存命中率等，延迟也是必需的，然后增加`sched_relax_domain_level`将使您受益。



### 如何使用 cpusets

为了最大限度地减少 cpusets 对关键内核代码（例如调度程序）的影响，并且由于内核不支持一个任务直接更新另一个任务的内存布局，更改其 cpuset CPU 或内存节点放置，或者更改任务附加到哪个 cpuset 对任务的影响是微妙的。

如果一个 cpuset 的内存节点被修改，那么对于附加到该 cpuset 的每个任务，下次内核尝试为该任务分配内存页时，内核将注意到该任务 cpuset 的变化，并更新其每个任务的内存节点。如果任务使用内存策略 `MPOL_BIND`，并且其绑定的节点与其新 cpuset 重叠，则该任务将继续使用新 cpuset 中仍允许的 `MPOL_BIND` 节点子集。如果任务使用 `MPOL_BIND` 并且现在新 cpuset 中不允许其任何 `MPOL_BIND` 节点，则该任务本质上将被视为绑定到新 cpuset 的 `MPOL_BIND`（即使其 NUMA 位置，如 `get_mempolicy()` 查询），不变）。如果一个任务从一个 cpuset 移动到另一个 cpuset，那么内核将在下次尝试为该任务分配内存页时调整该任务的内存布局，如上所述。

如果某个 cpuset 的`cpuset.cpus`已修改，则该 cpuset 中的每个任务将立即更改其允许的 CPU 布局。类似地，如果一个任务的 pid 被写入另一个 cpuset 的`tasks`文件，那么其允许的 CPU 布局会立即更改。如果此类任务已使用 `sched_setaffinity()` 调用绑定到其 cpuset 的某个子集，则该任务将被允许在其新 cpuset 中允许的任何 CPU 上运行，从而消除先前 `sched_setaffinity()` 调用的效果。

总之，cpuset 更改的任务的内存布局由内核在下次为该任务分配页时更新，并且处理器布局会立即更新。

通常，一旦分配了一个页（给定主内存的物理页），那么该页就会保留在它分配的任何节点上，只要它保持分配状态，即使 cpusets 内存放置策略`cpuset.mems`随后发生变化。
如果 cpuset 标志文件`cpuset.memory_migrate`设置为 true，则当任务附加到该 cpuset 时，该任务在其先前 cpuset 中的节点上分配给它的任何页都将迁移到该任务的新 cpuset。如果可能的话，在这些迁移操作期间会保留页在 cpuset 中的相对位置。
例如，如果该页位于先前 cpuset 的第二个有效节点上，则该页面将被放置在新 cpuset 的第二个有效节点上。

此外，如果`cpuset.memory_migrate`设置为 true，则如果该 cpuset 的`cpuset.mems`文件被修改，则分配给该 cpuset 中任务（位于之前`cpuset.mems`设置中的节点上）的页，将被移动到新设置`mems`中的节点。不在任务的先前 cpuset 中或不在 cpuset 的先前`cpuset.mems`设置中的页将不会被移动。

上述情况有一个例外。如果使用热插拔功能删除当前分配给某个 cpuset 的所有 CPU，则该 cpuset 中的所有任务都将移动到具有非空 cpu 的最近祖先。但是，如果 cpuset 与另一个对任务附加有一些限制的 cgroup 子系统绑定，则某些（或全部）任务的移动可能会失败。在这种失败的情况下，这些任务将保留在原始 cpuset 中，并且内核将自动更新其 `cpus_allowed` 以允许所有在线 CPU。当用于删除内存节点的内存热插拔功能可用时，预计也会出现类似的异常。 一般来说，内核更喜欢违反 cpuset 放置，而不是让所有允许的 CPU 或内存节点脱机的任务挨饿。

上述情况还有第二个例外。`GFP_ATOMIC` 请求是必须立即满足的内核内部分配。如果 `GFP_ATOMIC` 分配失败，内核可能会丢弃一些请求，在极少数情况下甚至会出现恐慌。如果当前任务的 cpuset 无法满足请求，那么我们会放松 cpuset，并在任何可以找到的地方寻找内存。侵犯 cpuset 比给内核施加压力要好。

要启动包含在 cpuset 中的新工作，步骤如下:

1. `mkdir /sys/fs/cgroup/cpuset`
2. `mount -t cgroup -ocpuset cpuset /sys/fs/cgroup/cpuset`
3. create the new cpuset by doing mkdir's and write's (or echo's) in the `/sys/fs/cgroup/cpuset` virtual file system.
4. start a task that will be the "founding father" of the new job
5. attach that task to the new cpuset by writing its pid to the `/sys/fs/cgroup/cpuset` tasks file for that cpuset
6. fork, exec or clone the job tasks from this founding father task.


例如，以下命令序列将设置一个名为“Charlie”的 cpuset，仅包含 CPU 2 和 3 以及内存节点 1，然后在该 cpuset 中启动一个子 shell 'sh'：

```bash
mount -t cgroup -ocpuset cpuset /sys/fs/cgroup/cpuset
cd /sys/fs/cgroup/cpuset
mkdir Charlie
cd Charlie
/bin/echo 2-3 > cpuset.cpus
/bin/echo 1 > cpuset.mems
/bin/echo $$ > tasks
sh
# The subshell 'sh' is now running in cpuset Charlie
# The next line should display '/Charlie'
cat /proc/self/cpuset
```

有多种方法可以查询或修改 cpusets：
- 直接通过 cpuset 文件系统，使用 shell 中的各种 cd、mkdir、echo、cat、rmdir 命令或 C 中的等效命令。
- 通过 C 库 libcpuset。
- 通过 C 库 libcgroup。（http://sourceforge.net/projects/libcg/）
- 通过 python 应用程序 cset。（http://code.google.com/p/cpuset/）


`sched_setaffinity` 调用也可以在 shell 提示符下使用 SGI 的 runon 或 Robert Love 的任务集来完成。 `mbind` 和 `set_mempolicy` 调用可以在 shell 提示符下使用 numactl 命令（Andi Kleen 的 numa 包的一部分）完成。



## 用法示例和语法


### 基本使用

创建、修改、使用 cpuset 可以通过 cpuset 虚拟文件系统来完成。

要挂载它，输入：

`mount -t cgroup -o cpuset cpuset /sys/fs/cgroup/cpuset`

然后在 /sys/fs/cgroup/cpuset 下，您可以找到与系统中 cpuset 树相对应的树。 例如，/sys/fs/cgroup/cpuset 是保存整个系统的 cpuset。

如果你想在`/sys/fs/cgroup/cpuset` 下创建新的 cpuset：

```bash
cd /sys/fs/cgroup/cpuset
mkdir my_cpuset
```

现在你想在这个 cpuset 上做点什么：

```bash
cd my_cpuset

# 在该文件夹下你会发现这些文件
ls
cgroup.clone_children   cpuset.memory_pressure
cgroup.event_control    cpuset.memory_spread_page
cgroup.procs            cpuset.memory_spread_slab
cpuset.cpu_exclusive    cpuset.mems
cpuset.cpus             cpuset.sched_load_balance
cpuset.mem_exclusive    cpuse.sched_relax_domain_level
cpuset.mem_hardwall     notify_on_release
cpuset.memory_migrate   tasks
```

读取它们将为您提供有关此 cpuset 状态的信息：它可以使用的 CPU 和内存节点、正在使用它的进程、它的属性。通过写入这些文件，您可以操纵 cpuset。

```bash
# set some flags
/bin/echo 1 > cpuset.cpu_exclusive

# add some cpus
/bin/echo 0-7 > cpuset.cpus

# add some mems
/bin/echo 0-7 > cpuset.mems

# now attach your shell to this cpuset
/bin/echo $$ > tasks

# you can also create cpusets inside your cpuset by using mkdir in this directory
mkdir my_sub_cs

# to remove a cpuset, just use rmdir
rmdir my_sub_cs
# this will fail if the cpuset is in use (has cpusets inside, or has processes attached).

# 请注意，由于遗留原因，“cpuset”文件系统作为 cgroup 文件系统的包装器而存在。

mount -t cpuset X /sys/fs/cgroup/cpuset
# 等价于
mount -t cgroup -ocpuset,noprefix X /sys/fs/cgroup/cpuset
echo "/sbin/cpuset_release_agent" > /sys/fs/cgroup/cpuset/release_agent
```


### 增加/移除 CPU

这是在 cpuset 目录中写入 cpus 或 mems 文件时使用的语法：

```bash
/bin/echo 1-4 > cpuset.cpus       # set cpus list to cpus 1,2,3,4
/bin/echo 1,2,3,4 -> cpuset.cpus  # set cpus list to cpus 1,2,3,4
```

要将 CPU 添加到 cpuset，请写入新的 CPU 列表，包括要添加的 CPU。 要将 6 添加到上述 cpuset：

`/bin/echo 1-4,6 > cpuset.cpus  # set cpus list to cpus 1,2,3,4,6`

类似地，要从 cpuset 中删除 CPU，请写入新的 CPU 列表，其中不包含要删除的 CPU。

移除所有的CPU：

`/bin/echo "" > cpuset.cpus   # clear cpus list`



### 设置标志

语法非常简单：

```bash
/bin/echo 1 > cpuset.cpu_exclusive  # set flag 'cpuset.cpu_exclusive'
/bin/echo 0 > cpuset.cpu_exclusive  # unset flag 'cpuset.cpu_exclusive'
```



### 附加进程

`/bin/echo PID > tasks`

注意，是PID，而不是PIDs。 您一次只能附加一项任务。
如果您有多个任务要附加，则必须一个接一个地执行：

```bash
/bin/echo PID1 > tasks
/bin/echo PID2 > tasks
# ...
/bin/echo PIDn > tasks
```



## QA

Q： 为什么是 `/bin/echo`
A： bash 的内置“echo”命令不会检查对 write() 的调用是否有错误。 如果您在 cpuset 文件系统中使用它，您将无法判断命令是成功还是失败。

Q： 当我附加多个进程时，只有第一行真正附加了
A： 每次调用 write() 只能返回一个错误代码。 所以你也应该只输入一个pid


