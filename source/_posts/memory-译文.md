---
title: memory 译文
date: 2023-12-26 22:03:03
tags:
categories:
- kernel doc
- memory
keywords:
- memory
copyright: Guader
copyright_author_href:
copyright_info:
---


[原文](https://www.kernel.org/doc/Documentation/cgroup-v1/memory.txt)


内存资源控制器

注意：这份文件已经过时了，它要求提供完整的信息改写。它仍然包含有用的信息，因此我们在这里保留它，但如果您需要更深入的了解，请务必检查当前最新代码。

注意：本文档中内存资源控制器通常被称为内存控制器。不要混淆这里的内存控制器和硬件中使用的内存控制器。

（对于编辑）
在本文档中： 当我们提到带有内存控制器的cgroup（cgroupfs的文件夹）时，我们称之为“memory cgroup”。当你看到 git-log 和源代码时，你会看到补丁的标题和函数名称倾向于使用“memcg”。在本文档中，我们避免使用它。



## 内存控制器的优点和用途

内存控制器隔离一组来自系统的其余部分任务的内存行为。LWN 上的文章提到了一些可能的内存控制器的用法。内存控制器可用于:

1. 隔离一个应用程序或一组应用程序。内存消耗大的应用程序可以被隔离并限制在较小的内存使用量范围内.
2. 创建一个内存有限的cgroup；这个可以用作启动是调用参数 `mem=XXXX` 的一个很好的替代方案。
3. 虚拟化解决方案可以控制他们想要分配给虚拟机实例的内存量。
4. CD/DVD 刻录机可以控制系统的其余部分所使用的内存量，以确保刻录不会因缺少可用内存而失败。
5. 还有其他几个用例；找到一个或为了好玩使用控制器（学习和破解 VM 子系统）。


当前状态：linux-2.6.34-mmotm（2010年4月开发版本）

*确实有点早了*

功能：
- 统计匿名页面、文件缓存、交换缓存的使用并限制它们。
- 页面专门链接到per-memcg LRU，并且没有全局LRU。
- 可选地，可以计算和限制内存+交换的使用。
- 分层会计
- 软限制
- 可以选择在移动任务时移动（充值）帐户。
- 使用阈值通知器
- 内存压力通知器
- oom-killer 禁用旋钮和 oom-notifier
- 根cgroup没有限制控制。

内核内存支持正在进行中，当前版本提供基本功能。（参见第 2.7 节）



控制文件的简要摘要。

文件 | 概要
--- | ---
 tasks | 附加任务（线程）并显示线程列表
 cgroup.procs | 显示进程列表
 cgroup.event_control | event_fd() 的接口
 memory.usage_in_bytes | 显示当前内存使用情况（详见5.5）
 memory.memsw.usage_in_bytes | 显示内存+交换的当前使用情况（详见5.5）
 memory.limit_in_bytes | 设置/显示内存使用限制
 memory.memsw.limit_in_bytes | 设置/显示内存+交换使用的限制
 memory.failcnt | 显示内存使用达到限制的次数
 memory.memsw.failcnt | 显示内存+Swap 达到限制的数量
 memory.max_usage_in_bytes | 显示记录的最大内存使用量
 memory.memsw.max_usage_in_bytes | 显示记录的最大内存+交换空间使用情况
 memory.soft_limit_in_bytes | 设置/显示内存使用的软限制
 memory.stat | 显示各种统计信息
 memory.use_hierarchy | 设置/显示启用的分层帐户
 memory.force_empty | 触发强制页面回收
 memory.Pressure_level | 设置内存压力通知
 memory.swappiness | 设置/显示 vmscan 的 swappiness 参数 （参见sysctl的vm.swappiness）
 memory.move_charge_at_immigrate | 设置/显示移动电荷的控制
 memory.oom_control | 设置/显示 oom 控件。
 memory.numa_stat | 显示每个numa节点的内存使用数量
 memory.kmem.limit_in_bytes | 设置/显示内核内存的硬限制
 memory.kmem.usage_in_bytes | 显示当前内核内存分配
 memory.kmem.failcnt | 显示内核内存使用达到限制的次数
 memory.kmem.max_usage_in_bytes | 显示记录的最大内核内存使用量
 memory.kmem.tcp.limit_in_bytes | 设置/显示 tcp buf 内存的硬限制
 memory.kmem.tcp.usage_in_bytes | 显示当前 tcp buf 内存分配
 memory.kmem.tcp.failcnt | 显示tcp buf内存使用达到限制的数量
 memory.kmem.tcp.max_usage_in_bytes | 显示记录的最大 tcp buf 内存使用量

## 1. 历史

内存控制器有着悠久的历史。内存评论请求控制器由 Balbir Singh 发布。RFC 发布时内存控制有多种实现方式。RFC 的目标旨在用于内存控制所需的最低功能建立共识和协议。第一个 RSS 控制器由 Balbir Singh 发布 于2007 年 2 月。Pavel Emelianov 此后发布了三个版本的 RSS 控制器。在OLS，在资源管理 BoF，每个人都建议我们同时处理页缓存和 RSS。还有一个要求是允许用户空间处理 OOM。当前的内存控制器是版本 6 ；它结合了映射（RSS）和未映射页缓存控制。



## 2. 内存控制

内存在环境中是一种独特的资源，因为它存在于有限的数量。如果某个任务需要大量 CPU 处理，则该任务可以将其处理时间拉长为数小时、数天、数月或数年，但内存，需要重复使用相同的物理内存来完成任务。

内存控制器的实现已分为几个阶段。这些是：

1. 内存控制器
2. `mlock(2)` 控制器
3. 内核用户内存核算和 slab 控制
4. 用户映射长度控制器

内存控制器是第一个开发的控制器。



### 2.1. 设计

设计的核心是一个称为 `page_counter` 的计数器。这个 `page_counter` 跟踪该组的与控制器相关的进程的当前内存使用情况和限制。每个 cgroup 都有一个与其关联的内存控制器特定数据结构 (`mem_cgroup`)。



## 2.2. 内存核算

```
TODO:
		+-------------------------------+
		| mem_cgroup                    |
		| （page_counter）              |
		+-------------------------------+
		 /            ^    \
		/             |     \
           +----------------+ |       +-------------+
           | mm_struct      | |....   | mm_struct   |
           |                | |       |             |
           +----------------+ |       +-------------+
                              |
                              + --------------+
                                              |
           +----------------+          +------+--------+
           | page           +----------> page_cgroup   |
           |                |          |               |
           +----------------+          +---------------+

             （图1：Hierarchy of Accounting）
```


图 1 显示了控制器的重要方面:

1. 按每个 cgroup 进行记账
2. 每个 `mm_struct` 都知道它属于哪个 `cgroup`
3. 每个页都有一个指向 `page_cgroup` 的指针，而 `page_cgroup` 又知道它所属的cgroup

内存核算过程如下：调用 `mem_cgroup_charge_common()`设置必要的数据结构并检查正在使用的 cgroup 已计算内存是否超过其限制。如果是，则在 cgroup 上调用回收。 更多详细信息可以在本文档的回收部分找到。
如果一切顺利，则会出现一个名为 `page_cgroup` 的 `page meta-data-structure`更新。`page_cgroup` 在 cgroup 上有自己的 LRU。

`page_cgroup` 结构在引导/内存热插拔时分配。


#### 2.2.1 内存核算细节

所有映射的匿名页面（RSS）和缓存页（Page Cache）都被计算在内。
一些永远不可回收且不会出现在 LRU 上的页不会被计算在内。 我们只在通常的虚拟机管理下记入页。

RSS 页面在 `page_fault` 处进行计算，除非它们之前已被计算过。文件页面插入 `inode`（基数树）时将被视为页缓存。当它被映射到进程的页表时，会小心地避免重复计算。

当 RSS 页面完全未映射时，该页就无法被核算。当 PageCache 页从基数树(radix-tree)中删除时，该页就不再被核算。即使 RSS 页面完全取消映射（通过 kswapd），它们也可能作为 SwapCache 存在于系统中，直到它们真正被释放。这样的 SwapCache 也会被核算。换入的页(A swapped-in page)在被映射之前不会被计算在内。

注意：内核会进行交换预读(swapin-readahead)并一次读取多个交换。这意味着换入的页可能包含其他任务的页面，而不只是导致页错误的任务。因此，我们避免在 `swap-in I/O` 时进行内存核算。

在页面迁移时，内存核算信息会保留。

注意：我们只核算 LRU 上的页，因为我们的目的是控制使用页的总量；从 VM 角度来看，非 LRU 页(not-on-LRU pages)往往会失控。



### 2.3 共享页内存核算

共享页按照首次接触方式进行核算。最先接触页的 cgroup 将被计入该页。这种方法背后的原理是，积极使用共享页的 cgroup 最终将为此进行分配内存（一旦它从引入它的 cgroup 中释放出来 - 这将在内存压力下发生）。

但请参阅第 8.2 节：当将任务移动到另一个 cgroup 时，如果已选择 `move_charge_at_immigrate`，则其页可能会重新分配到新的 cgroup。

例外：如果未使用 `CONFIG_MEMCG_SWAP`。当您执行 `swapoff` 并将 `shmem(tmpfs)` 的换出页强制备份到内存中时，页分配将由 `swapoff` 的调用者而不是 `shmem` 的用户负责。



### 2.4 交换扩展（`CONFIG_MEMCG_SWAP`）

交换扩展允许您记录交换的内存分配。如果可能的话，换入的页将被归还给原始页分配器。

当计算交换的内存时，会添加以下文件:

- `memory.memsw.usage_in_bytes`
- `memory.memsw.limit_in_bytes`

`memsw` 的意思是内存+交换区(memory + swap)。内存+交换的使用受到 `memsw.limit_in_bytes` 的限制。

示例：假设系统具有 4G 交换内存空间。在 2G 内存限制下（错误地）分配 6G 内存的任务将使用所有交换区。在这种情况下，设置 `memsw.limit_in_bytes=3G` 将防止交换区的不当使用。通过使用 memsw 限制，您可以避免因交换空间不足而导致的系统 OOM。

*为什么是“内存+交换”而不是交换。*

*全局 LRU(kswapd) 可以交换任意页。`Swap-out` 是指将核算从内存转移到 swap..., `memory+swap` 的使用没有变化。换句话说，当我们想要限制 swap 的使用而不影响全局LRU时，从操作系统的角度来看，memory+swap 限制比仅仅限制 swap 更好。*


*当 cgroup 达到 `memory.memsw.limit_in_bytes` 时会发生什么*

*当cgroup达到 `memory.memsw.limit_in_bytes` 时，在这个cgroup中进行 `swap-out` 是没有用的。然后，cgroup 例程将不会完成换出，并且文件缓存将被删除。但如上所述，全局 LRU 可以从中交换内存，以保证系统内存管理状态的完整性。 你不能通过 cgroup 来禁止它。*



### 2.5 回收

每个cgroup维护一个 `per cgroup LRU`，其结构与全局 VM 相同。当 cgroup 超出其限制时，我们首先尝试从 cgroup 回收内存，以便为 cgroup 触及的新页腾出空间。如果回收不成功，则会调用 OOM 例程来选择并终止 cgroup 中最庞大的任务。（参见下面的 10. OOM 控制。）

cgroup 的回收算法没有修改，只是选择回收的页来自每个 cgroup LRU 列表。

注意：回收对根 cgroup 不起作用，因为我们无法设置任何对根 cgroup 的限制。

注2：当 `panic_on_oom` 设置为 “2” 时，整个系统会 panic。

当 oom 事件通知程序被注册时，事件将被传递。（参见 `oom_control` 部分）



### 2.6 锁定

`lock_page_cgroup()/unlock_page_cgroup()` 不应在 `i_pages` 锁下调用。

其他锁顺序如下：
`PG_locked`.
`mm->page_table_lock`
    `pgdat->lru_lock`
        `lock_page_cgroup`.

在许多情况下，只调用 `lock_page_cgroup()`。
`per-zone-per-cgroup LRU`（cgroup 的私有 LRU ）仅由 `pgdat->lru_lock` 保护，它没有自己的锁。



### 2.7 内核内存扩展（`CONFIG_MEMCG_KMEM`）

通过内核内存扩展，内存控制器能够限制系统使用的内核内存量。内核内存从根本上来说是与用户内存不同，因为它不能被换出，这使得它通过消耗过多的这一宝贵资源，可能会对系统造成 DoS 攻击。

默认情况下，所有内存 cgroup 都会启用内核内存统计。但可以在启动时通过将 `cgroup.memory=nokmem` 传递给内核来在系统范围内禁用它。在这种情况下，根本不会核算内核内存。

根 cgroup 没有受到内核内存限制。根 cgroup 的内存用量可能会也可能不会被计算在内。使用的内存被累积到 `memory.kmem.usage_in_bytes`，或者在有意义的情况下在单独的计数器中。（目前仅适用于 TCP）。

主 “kmem” 计数器被馈送到主计数器中，因此 kmem 分配将也可以从用户计数器中看到。

目前内核内存没有实现软限制。当达到这些限制时触发 slab 回收是未来的工作。



#### 2.7.1 当前被计算的内核内存资源

*堆栈页*： 每个进程都会消耗一些堆栈页。通过计算内核内存，我们可以防止在内核内存使用率过高时创建新进程。

*slab页*： 跟踪由 SLAB 或 SLUB 分配器分配的页面。每次从 memcg 内部第一次触及缓存时，都会创建每个 `kmem_cache` 的副本。创建是延迟完成的，因此在创建缓存时仍然可以跳过某些对象。一个 `slab` 页面中的所有对象应该属于同一个 `memcg`。仅当缓存在页面分配期间将任务迁移到不同的 `memcg` 时，此情况才成立。

*套接字内存压力*： 一些套接字协议有内存压力阈值。内存控制器允许按 cgroup 单独控制它们，而不是全局控制。 

*tcp 内存压力*： tcp 协议的套接字内存压力。



#### 2.7.2 常见用例

因为“kmem”计数器被馈送到主用户计数器，内核内存永远不能完全独立于用户内存而受到限制。假设“U”是用户限制，“K”是内核限制。 可以通过三种可能的方式设置限制：

```
U != 0，K = 无限制：
这是在 kmem 核算之前就已经存在的标准 memcg 限制机制。内核内存被完全忽略。

U != 0，K < U：
内核内存是用户内存的子集。此设置在每个 cgroup 内存总量过量使用的部署中非常有用。绝对不建议过度使用内核内存限制，因为盒子仍然可能耗尽不可回收的内存。
在这种情况下，管理员可以设置 K，使所有组的总和永远不会大于总内存，并以牺牲 QoS 为代价自由设置 U。

警告：在当前的实现中，当 cgroup 达到 K 而保持在 U 以下时，不会触发内存回收，这使得此设置不切实际。

U != 0，K >= U：
由于 kmem 费用也将被馈送到用户计数器，并且这两种内存的 cgroup 的回收将被触发。此设置为管理员提供了统一的内存视图，对于只想跟踪内核内存使用情况的人也很有用。
```



## 3. 用户界面

### 3.0.配置

1. 启用 `CONFIG_CGROUPS`
2. 启用 `CONFIG_MEMCG`
3. 启用 `CONFIG_MEMCG_SWAP`（使用 swap 扩展）
4. 启用 `CONFIG_MEMCG_KMEM`（使用 kmem 扩展）



### 3.1. 准备 cgroup（请参阅 cgroups.txt，为什么需要 cgroup？）

```bash
mount -t tmpfs none /sys/fs/cgroup
mkdir /sys/fs/cgroup/memory
mount -t cgroup none /sys/fs/cgroup/memory -o memory
```



### 3.2. 创建新组并将 bash 移入其中

```bash
mkdir /sys/fs/cgroup/memory/0
echo $$ > /sys/fs/cgroup/memory/0/tasks
```

由于现在我们位于 0 cgroup，我们可以更改内存限制：

`echo 4M > /sys/fs/cgroup/memory/0/memory.limit_in_bytes`

注意：我们可以使用后缀（k、K、m、M、g 或 G）来表示以千为单位的值，
MB 或 GB。（这里，Kilo、Mega、Giga 分别是 Kibibytes、Mebibytes、Gibibytes。）

注意：我们可以写 “-1” 来重置 `*.limit_in_bytes`（无限制）。
注意：我们不能再对根 cgroup 设置限制。

`cat /sys/fs/cgroup/memory/0/memory.limit_in_bytes`
4194304

我们可以检查一下使用情况：

`cat /sys/fs/cgroup/memory/0/memory.usage_in_bytes`
1216512

成功写入此文件并不能保证成功将此限制设置为写入文件的值。这可能是由于多种因素造成的，例如向上舍入到页边界或系统上内存的总可用性。用户在写入后需要重新读取该文件，以保证内核提交的值。

```bash
echo 1 > memory.limit_in_bytes
cat memory.limit_in_bytes

# 4096
```

`memory.failcnt` 字段给出了超出 cgroup 限制的次数。

`memory.stat` 文件提供核算信息。现在，显示缓存数量、RSS 和活动页面/非活动页面。



## 4. 测试

有关测试功能和实现，请参阅 `memcg_test.txt`。

性能测试也很重要。要查看纯内存控制器的开销，对 tmpfs 进行测试将为您带来大量的小开销。
例如： 在 tmpfs 上执行内核 make。

页面错误可伸缩性也很重要。在测量并行页面错误测试时，多进程测试可能比多线程测试更好，因为它具有共享对象/状态的噪声。

但以上两个都是测试极端情况。
在内存控制器下尝试常规测试总是有帮助的。



### 4.1 故障排除

有时，用户可能会发现某个cgroup下的应用程序被 OOM Killer 终止。造成这种情况的原因有以下几个：

1. cgroup限制太低（太低了，啥也做不了）
2. 用户正在使用匿名内存并且 swap 被关闭或太低

同步之后 `echo 1 > /proc/sys/vm/drop_caches` 将帮助删除 cgroup 中缓存的一些页面（页面缓存页面）。

要了解发生了什么，请按照 “10. OOM Control”（如下）禁用 `OOM_Kill` 并看看发生了什么会有所帮助。



### 4.2 任务迁移

当一个任务从一个 cgroup 迁移到另一个 cgroup 时，它的分配默认不会转过去。从原来的 cgroup 分配的页仍然保持分配状态，当页面被释放或回收分配就会被移除。

您可以随着任务迁移而移动任务分配。
请参阅 8.“在任务迁移时移动分配”



### 4.3 删除cgroup

cgroup 可以通过 rmdir 删除，但如 4.1 和 4.2 节中讨论的，即使所有任务都已从 cgroup 迁移出去，cgroup 也可能有一些与之相关的分配。（因为我们按页分配，而不是按任务收费。）

我们将统计数据移动到根（如果 `use_hierarchy==0`）或父（如果
`use_hierarchy==1`)，除了从子释放外，分配没有变化。

删除 cgroup 时，交换信息中记录的分配不会更新。记录的信息将被丢弃，并且使用交换（swapcache）的 cgroup 将作为它的新所有者被分配。

关于 `use_hierarchy`，请参见第6节。



## 5. 其他。接口。

### 5.1 `force_empty`

提供了 `memory.force_empty` 接口来清空 cgroup 的内存使用量。
当写任何东西到这个

`echo 0 > memory.force_empty`

cgroup 将被回收，并回收尽可能多的页。

该接口的典型用例是在调用 `rmdir()` 之前。
虽然 `rmdir()` 使 `memcg` 离线(offline)，但由于文件缓存收费，memcg 可能仍留在那里。一些未使用的页缓存可能会一直分配，直到发生内存压力。如果你想避免这种情况，`force_empty` 会很有用。

另请注意，当设置了 `memory.kmem.limit_in_bytes` 时，仍会看到由于内核页面而产生的分配。这不被视为失败，写入仍会返回成功。在这种情况下，预计 `memory.kmem.usage_in_bytes == memory.usage_in_bytes`。

关于 `use_hierarchy`，请参见第6节。



### 5.2 统计文件

`memory.stat` 文件包含以下统计信息

```
每个内存 cgroup 本地状态

cache               - 页缓存的字节数。
rss                 - 匿名和交换缓存内存的字节数（包括透明大页）。
rss_huge            - 匿名透明大页的字节数。
mapped_file         - 映射文件的字节数（包括 tmpfs/shmem）
pgpgin              - 内存 cgroup 的分配事件数。每当一个页面被视为映射到 cgroup 的匿名页面 (RSS) 或缓存页面 (Page Cache) 时，就会发生分配事件。
pgpgout             - 内存 cgroup 的取消分配事件数。每当一个页面未从 cgroup 中记入时，就会发生取消分配事件。
swap                - 交换使用的字节数
dirty               - 等待写回磁盘的字节数。
writeback           - 排队等待同步到磁盘的文件/匿名缓存的字节数。
inactive_anon       - 非活动 LRU 列表上的匿名和交换高速缓存内存的字节数。
active_anon         - 活动 LUR 列表上匿名和交换高速缓存的字节数
inactive_file       - 非活动 LRU 列表上文件支持内存的字节数。
active_file         - 活动 LRU 列表上文件支持内存的字节数。
unevictable         - 无法回收的内存字节数（mlocked 等）。

# 考虑层次结构的状态（请参阅 memory.use_hierarchy 设置）

hierarchy_memory_limit  - 与内存 cgroup 所在层次结构相关的内存限制字节数
hierarchical_memsw_limit - 关于内存 cgroup 所在层次结构的内存+交换字节数限制。

total_<counter>_    - # <counter> 的分层版本，除了 cgroup 自己的值之外，还包括 <counter> 的所有分层子级值的总和，即 total_cache

# 以下附加统计信息取决于 CONFIG_DEBUG_VM。

centre_rotated_anon - VM 内部参数。（参见mm/vmscan.c）
centre_rotated_file - VM 内部参数。（参见mm/vmscan.c）
centre_scanned_anon - VM 内部参数。（参见mm/vmscan.c）
centre_scanned_file - VM 内部参数。（参见mm/vmscan.c）
```

备忘录：
`recent_rotated` 表示最近的 LRU 轮转频率。
`recent_scanned` 表示最近扫描到 LRU 的次数。
为了更好的调试，请参阅代码以了解含义。

注意：
仅匿名和交换高速缓存内存被列为 “rss” 统计信息的一部分。
不应将其与真正的“驻留集大小(resident set size)”或 cgroup 使用的物理内存量相混淆。
`“rss + mapped_file”`将为您提供 cgroup 的驻留集大小。
（注意：文件和 shmem 可能在其他 cgroup 之间共享。在这种情况下，仅当内存 cgroup 是页面缓存的所有者时，`mapped_file` 才会被计算在内。）



### 5.3 交换性 swappiness

覆盖特定组的 `/proc/sys/vm/swappiness`。根 cgroup 中的可调参数对应于全局交换设置。

请注意，与全局回收期间不同，限制回收会强制执行 0 交换，即使有可用的交换存储，也确实会阻止任何交换。 如果没有要回收的文件页面，这可能会导致 `memcg OOM Killer`。



### 5.4 失败计数 failcnt

内存 cgroup 提供了 `memory.failcnt` 和 `memory.memsw.failcnt` 文件。
这个 failcnt(==failure count)显示了使用计数器达到其限制的次数。当内存 cgroup 达到限制时，failcnt 会增加，并且其下的内存将被回收。 

您可以通过向 failcnt 文件写入0 来重置 failcnt。

`echo 0 > .../memory.failcnt`



### 5.5 `usage_in_bytes`

为了提高效率，与其他内核组件一样，内存 cgroup 使用了一些优化来避免不必要的缓存行错误共享。`use_in_bytes` 受该方法的影响，不显示内存（和交换）使用量的 “准确” 值，它是有效访问的模糊值。（当然，必要时，它是同步的。）如果你想知道更准确的内存使用情况，你应该使用 `memory.stat` 中的 RSS+CACHE(+SWAP) 值（参见5.2）。



### 5.6 `numa_stat`

这与 `numa_maps` 类似，但基于每个 memcg 进行操作。 这对于提供对 memcg 内的 numa 位置信息的可见性非常有用，因为允许从任何物理节点分配页面。其中一个用例是通过将此信息与应用程序的 CPU 分配相结合来评估应用程序性能。

每个 memcg 的 `numa_stat` 文件都包含 “total”、“file”、“anon” 和 “unevictable” 每节点页计数，其中包括 `“hierarchical_<counter>”`，它除了 memcg 自身的值之外还总结了所有分层子级的值。

`memory.numa_stat` 的输出格式为：

```
total=<total pages> N0=<node 0 pages> N1=<node 1 pages> ...
file=<total file pages> N0=<node 0 pages> N1=<node 1 pages> ...
anon=<total anon pages> N0=<node 0 pages> N1=<node 1 pages> ...
unevictable=<total anon pages> N0=<node 0 pages> N1=<node 1 pages> ...
hierarchical_<counter>=<counter pages> N0=<node 0 pages> N1=<node 1 pages> ...
```

"total" 的计数是 `file + anon + unevictable` 的总和。



## 6. 层次结构支持

内存控制器支持深层次结构和分层记账。通过在 cgroup 文件系统中创建适当的 cgroup 来创建层次结构。 例如，考虑以下 cgroup 文件系统层次结构:

```
        root
       / | \
      /  |  \
     a   b   c
             | \
             |  \
             d   e
```

在上图中，启用分层记帐后，e 的所有内存使用量都会记入其祖先，直到启用了 `memory.use_hierarchy` 的根（即 c 和 root）。 如果祖先之一超出其限制，回收算法将从祖先及其子代中的任务中回收。



### 6.1 启用分级核算和回收

默认情况下，内存 cgroup 会禁用层次结构功能。可以通过将 1 写入根 cgroup 的 `memory.use_hierarchy` 文件来启用支持：

`echo 1 > memory.use_hierarchy`

该功能也可以这样禁用：

`echo 0 > memory.use_hierarchy`

注意1： 如果该 cgroup 已在其下方创建了其他 cgroup，或者父 cgroup 启用了 `use_hierarchy`，则启用/禁用将会失败。

注意2： 当 `panic_on_oom` 设置为“2”时，如果任何 cgroup 中发生 OOM 事件，整个系统都会panic。



## 7. 软限制

软限制允许更大程度的内存共享。 软限制背后的想法是允许控制组根据需要使用尽可能多的内存，前提是：

1. 不存在内存争用问题
2. 不超过硬限制

当系统检测到内存争用或内存不足时，控制组将被推回到其软限制。如果每个对照组的软限制非常高，则它们会被尽可能地推迟，以确保一个对照组不会耗尽其他对照组的内存。

请注意，软限制是尽力而为的功能；它没有任何保证，但它会尽力确保当内存严重竞争时，根据软限制提示/设置来分配内存。当前设置了基于软限制的回收，以便从 `balance_pgdat` (kswapd) 调用它。



### 7.1 接口

可以使用以下命令设置软限制（在本示例中，我们假设软限制为 256 MiB）

`echo 256M > memory.soft_limit_in_bytes`

如果我们想将其改为 1G， 我们可以随时：

`echo 1G > memory.soft_limit_in_bytes`

注意1： 软限制会在很长一段时间内生效，因为它们涉及回收内存以在内存 cgroup 之间进行平衡。

注意2： 建议将软限制设置为始终低于硬限制，否则硬限制将优先。



## 8. 在任务迁移时移动分配

用户可以在任务迁移的同时移动与任务相关的分配，即从旧 cgroup 中取消任务页面的分配，并将其重新在新 cgroup 中分配。由于缺少页表，`!CONFIG_MMU` 环境中不支持此功能。



### 8.1 接口

默认情况下禁用此功能。可以通过写入目标 cgroup 的 `memory.move_charge_at_immigrate` 来启用（并再次禁用）它。

如果你要启用它：

`echo (some positive value) > memory.move_charge_at_immigrate`

注意： `move_charge_at_immigrate` 的每一位对于应移动哪种类型的分配都有其自己的含义。详情参见8.2。

注意： 仅当您移动 `mm->owner`（换句话说，线程组的领导者）时，分配才会移动。

注意： 如果我们在目标 cgroup 中找不到足够的空间来容纳该任务，我们会尝试通过回收内存来腾出空间。如果无法腾出足够的空间，任务迁移可能会失败。

注意： 如果您大量移动分配，可能需要几秒钟的时间。


如果你想再次禁用它：

`echo 0 > memory.move_charge_at_immigrate`



### 8.2 可以移动的分配类型

`move_charge_at_immigrate` 中的每一位对于应移动哪种类型的分配都有其自己的含义。但无论如何，必须注意的是，页面或交换的分配只有在记入任务当前（旧）内存 cgroup 时才能移动。

位 bit | 可以移动的分配类型
--- | ---
0   | 目标任务使用的匿名页（或其交换）的分配。必须启用交换扩展(见 2.4)才能移动交换分配
1   | 目标任务映射的文件页面（普通文件、tmpfs 文件（例如 ipc 共享内存）和 tmpfs 文件的交换）的费用。与匿名页面的情况不同，即使任务没有发生页面错误，任务映射范围内的文件页面（和交换）也会被移动，即它们可能不是该任务的“RSS”，而是其他任务的“RSS” 映射相同的文件。并且页面的 mapcount 被忽略（即使 `page_mapcount(page) > 1`，页也可以被移动）。您必须启用交换扩展（请参阅 2.4）才能转移交换分配。



### 8.3 TODO

- 所有的移动分配操作都是在 `cgroup_mutex` 下完成的。保持互斥体时间过长并不是一个好行为，因此我们可能需要一些技巧。



## 9. 内存阈值

内存 cgroup 使用 cgroups 通知 API 实现内存阈值（请参阅 cgroups.txt）。 它允许注册多个内存和 memsw 阈值，并在超过阈值时收到通知。

要注册阈值，应用必须：

- 通过 `eventfd(2)` 创建一个 `eventfd`
- 打开 `memory.usage_in_bytes` 或是 `memory.memsw.usage_in_bytes`
- 将类似 `<event_fd> <fd of memory.usage_in_bytes> <threshold>` 的字符串写入 `cgroup.event_control` 

当内存使用量在任何方向上超过阈值时，应用程序都会通过 eventfd 收到通知。

它适用于根 cgroup 和非根 cgroup。



## 10. 内存溢出控制 OOM Control

`memory.oom_control` 文件用于 OOM 通过和其他控件。

内存 cgroup 使用 cgroup 通知 API 实现 OOM 通知程序（请参阅 cgroups.txt）。 它允许注册多个 OOM 通知传递并在 OOM 发生时获取通知。

要注册通知器，应用必须：

- 通过 `eventfd(2)` 创建一个 `eventfd`
- 打开 `memory.oom_control` 文件
- 将类似 `<event_fd> <fd of memory.oom_control>` 的字符串写入 `cgroup.event_control`

当 OOM 发生时，应用程序将通过 `eventfd` 得到通知。
OOM 通知对根 cgroup 不起作用。

你可以通过向 `memory.oom_control` 文件写入 “1” 来禁用 `OOM-killer`，如下：

`echo 1 > memory.oom_control`

如果禁用 `OOM-killer`，则 cgroup 下的任务在请求可用内存时将挂起/睡眠在内存 cgroup 的 `OOM-waitqueue` 中。

为了运行它们，你必须通过以下方式放松内存 cgroup 的 OOM 状态：

- 扩大限制或减少使用

要减少使用量：

- 终止一些任务
- 通过核算迁移将一些任务移动到其他组
- 移除一些文件 （在 tmpfs 上？）

然后，停止的任务将再次运行。

读取时，会显示 OOM 的当前状态。

- `oom_kill_disable`    0 或 1  （如果是 1， oom-killer 被禁用状态）
- `under_oom`           0 或 1  （如果是 1， 内存cgroup OOM，任务可能会停止。）



## 11. 内存压力

压力级别通知可用于监控内存分配成本； 根据压力，应用程序可以实施不同的策略来管理其内存资源。 压力水平定义如下：

- `low`  意味着系统正在回收内存以进行新的分配。 监视此回收活动可能有助于维护缓存级别。收到通知后，程序（通常是“活动管理器”）可能会分析 vmstat 并提前采取行动（即提前关闭不重要的服务）。
- `medium`  意味着系统正在经历中等内存压力，系统可能正在进行交换、分页活动文件缓存等。在此事件上，应用程序可能决定进一步分析 `vmstat/zoneinfo/memcg` 或内部内存使用统计信息，并释放任何可以轻松重建或从磁盘重新读取的资源。
- `critical`    意味着系统正在主动抖动，即将内存不足（OOM），甚至内核中的 OOM 杀手正在触发。 应用程序应该尽其所能来帮助系统。 现在咨询 vmstat 或任何其他统计数据可能为时已晚，因此建议立即采取行动。

默认情况下，事件向上传播直到事件被处理，即事件不传递。例如，您有三个 cgroup： `A->B->C`。现在，您在 cgroup A、B 和 C 上设置了一个事件侦听器，并假设 C 组遇到了一些压力。在这种情况下，只有C组会收到通知，即 A 组和 B 组不会收到通知。这样做是为了避免过度 “广播” 消息，这会扰乱系统，如果内存不足或出现抖动，情况尤其糟糕。B 组，仅当 C 组没有事件列表时才会收到通知。

共有三种可选模式指定不同的传播行为：

- `default` 这是上面指定的默认行为。 此模式与省略可选模式参数相同，通过向后兼容性保留。
- `hierarchy`   事件总是向上传播到根，与默认行为类似，不同之处在于，无论每个级别是否有事件侦听器，传播都会继续，采用 “层次结构” 模式。在上面的示例中，A、B、C 组将收到内存压力通知。
- `local`   事件是传递的，即，它们仅在注册通知的 memcg 中遇到内存压力时接收通知。在上面的示例中，如果注册了 “本地” 通知并且该组遇到内存压力，则组 C 将收到通知。但是，如果 B 组注册了本地通知，则无论是否有 C 组的事件侦听器，B 组都永远不会收到通知。

级别和事件通知模式（`hierarchy` 或 `local`，如果需要）由逗号分隔的字符串指定，即 `low,hierarchy` 指定所有祖先 memcgs 的分层、直通、通知。通知是默认的非传递行为，不指定模式。`medium,local` 指定中等级别的直通通知。

文件 `memory.pressure_level` 仅用于设置 `eventfd`。要注册通知，应用程序必须：

- 通过 `eventfd(2)` 创建一个 `eventfd`
- 打开 `memory.pressure_level`
- 将类似于 `<event_fd> <fd of memory.pressure_level> <level[,mode]>` 格式的字符串写入 `cgroup.event_control`

当内存压力达到特定级别（或更高）时，应用程序将通过 eventfd 收到通知。未实现对 `memory.pressure_level` 的读/写操作。


测试： 下面是一个小脚本示例，它创建一个新的 cgroup，设置内存限制，在 cgroup 中设置通知，然后使子 cgroup 经历临界压力：

```bash
cd /sys/fs/cgroup/memory
mkdir foo
cd foo
cgroup_event_listener memory.pressure_level low,hierarchy & 
echo 8000000 > memory.limit_in_bytes
echo 8000000 > memory.memsw.limit_in_bytes
echo && > tasks
dd if=/dev/zero | read x

# （预计会出现一堆通知，最终，oom-killer 将会触发。）
```



## 12. TODO

1. 让每个 cgroup 扫描程序首先回收非共享面
2. 教导控制器核算共享页
3. 当尚未达到限制但使用量越来越接近时在后台启动回收



## 概述

总的来说，内存控制器是一个稳定的控制器，并且在社区中得到了广泛的评论和讨论。

