---
title: cgroups-v1 译文
date: 2024-01-03 15:00:10
tags:
categories:
- kernel doc
- cgroup-v1
keywords:
- cgroup-v1
copyright: Guader
copyright_author_href:
copyright_info:
---

[原文](https://www.kernel.org/doc/Documentation/cgroup-v1/cgroups.txt)


# 1. Control Groups


## 1.1 What are cgroups ?

控制组提供了一种聚合/分区任务及其所有后续子任务集的机制，按专门的行为以等级分组。

定义：

*cgroup* 将一组任务与一个或多个子系统的一组参数相关联。

*子系统* 是一个模块，它利用 cgroup 提供的任务分组功能以特定方式处理任务组。
子系统通常是一个“资源控制器”，用于调度资源或应用每个 cgroup 的限制，但它也可以是任何想要对一组进程执行操作的东西，例如 虚拟化子系统。

*hierarchy* 是一组排列在树中的 cgroup，这样系统中的每个任务都恰好位于层次结构中的一个 cgroup 中，并且是一组子系统； 每个子系统都具有附加到层次结构中每个 cgroup 的系统特定状态。
每个层次结构都有一个与其关联的 cgroup 虚拟文件系统实例。

在任一时刻，可能存在多个活动的任务 cgroup 层次结构。
每个层次结构是系统中所有任务的一个分区。

用户级代码可以在cgroup虚拟文件系统的实例中按名称创建和销毁cgroup，指定和查询任务被分配到哪个cgroup，并列出分配给cgroup的任务PID。
这些创建和分配仅影响与该 cgroup 文件系统实例关联的层次结构。

就其本身而言，cgroup 的唯一用途是进行简单的工作跟踪。
目的是让其他子系统接入通用 cgroup 支持，为 cgroup 提供新属性，例如记账/限制 cgroup 中的进程可以访问的资源。
例如，cpuset（请参阅 Documentation/cgroup-v1/cpusets.txt）允许您将一组 CPU 和一组内存节点与每个 cgroup 中的任务关联起来。


## 1.2 Why are cgroups needed ?

为了在 Linux 内核中提供进程聚合，人们做出了多种努力，主要是为了资源跟踪的目的。
此类工作包括 `cpuset`、`CKRM/ResGroups`、`UserBeanCounters` 和虚拟服务器命名空间。
这些都需要进程分组/分区的基本概念，新分叉的进程最终位于与其父进程相同的组 (cgroup) 中。

内核 cgroup 补丁提供了有效实现此类组所需的最低限度的基本内核机制。
它对系统快速路径的影响最小，并为特定子系统（例如 cpuset）提供钩子，以根据需要提供其他行为。

提供多层次结构支持，以允许不同子系统将任务划分为 cgroup 明显不同的情况 - 
具有并行层次结构允许每个层次结构成为任务的自然划分，而不必处理在以下情况下出现的复杂任务组合： 
需要将几个不相关的子系统强制放入同一 cgroup 树中。

在一种极端情况下，每个资源控制器或子系统可以位于单独的层次结构中； 
在另一个极端，所有子系统都将附属于同一层次结构。

作为可以从多个层次结构中受益的场景（最初由 `vatsa@in.ibm.com` 提出）的示例，请考虑具有各种用户（学生、教授、系统任务等）的大型大学服务器。
该服务器的资源规划可以遵循以下原则：

```
    CPU             "Top cpuset"
                    /          \
                CPUSet1       CPUSet2
                   |             |
                (Professors)  (Students)

                此外（系统任务）附加到topcpuset（以便它们可以在任何地方运行），限制为20%

    Memory:  Professors (50%), Students (30%), system (20%)
    Disk:  Professors (50%), Students (30%), system (20%)
    Network: WWW browsing (20%), Network File System (60%), others (20%)
                           / \
           Professors (15%)   Students (5%)
```

像 `Firefox/Lynx` 这样的浏览器进入 `WWW` 网络类，而 `(k)nfsd` 进入 `NFS` 网络类。

同时，`Firefox/Lynx` 将根据启动者（教授/学生）共享适当的 CPU/内存类别。

由于能够针对不同的资源对任务进行不同的分类（通过将这些资源子系统放在不同的层次结构中），管理员可以轻松设置一个接收执行通知的脚本，并且根据谁启动浏览器，他可以:

`echo browser_pid > /sys/fs/cgroup/<restype>/<userclass>/tasks`

由于只有一个层次结构，他现在可能必须为每个启动的浏览器创建一个单独的 cgroup，并将其与适当的网络和其他资源类相关联。
这可能会导致此类 cgroup 的扩散。

另外，假设管理员希望临时为学生的浏览器提供增强的网络访问权限
（因为现在是晚上，并且用户想要进行在线游戏:) ），
或者为学生的其中一个模拟应用程序提供增强的 CPU 能力。

能够将 PID 直接写入资源类，只需执行以下操作：

```bash
echo pid > /sys/fs/cgroup/network/<new_class>/tasks
# (after some time)
echo pid > sys/fs/cgroup/network/<orig_class>/tasks
```

如果没有这种能力，管理员就必须将 cgroup 拆分为多个单独的 cgroup，然后将新的 cgroup 与新的资源类相关联。


## 1.3 How are cgroups implemented ?

控制组对内核进行了如下扩展：

- 系统中的每个任务都有一个指向 `css_set` 的引用计数指针。
- `css_set` 包含一组指向 `cgroup_subsys_state` 对象的引用计数指针，每个 cgroup 子系统对应一个在系统中注册的子系统。
    从任务到它在每个层次结构中所属的 cgroup 之间没有直接链接，但这可以通过跟踪 `cgroup_subsys_state` 对象的指针来确定。
    这是因为访问子系统状态是在性能关键型代码中经常发生的事情，而需要任务的实际 cgroup 分配（特别是在 cgroup 之间移动）的操作不太常见。
    链表使用 `css_set` 贯穿每个 `task_struct` 的 `cg_list` 字段，锚定在 `css_set->tasks` 处。
- 可以安装 cgroup 层次结构文件系统，以便从用户空间进行浏览和操作。
- 您可以列出附加到任何 cgroup 的所有任务（按 PID ）。

cgroups 的实现需要一些简单的钩子到内核的其余部分，在性能关键路径中没有：

- 在 `init/main.c` 中，在系统启动时初始化根 cgroup 和初始 `css_set`。
- 在 fork 和 exit 中，从其 `css_set` 附加和分离任务。

此外，可以挂载 “cgroup” 类型的新文件系统，以允许浏览和修改内核当前已知的 cgroup。
挂载 cgroup 层次结构时，您可以指定要挂载的子系统的逗号分隔列表作为文件系统挂载选项。
默认情况下，挂载 cgroup 文件系统会尝试挂载包含所有已注册子系统的层次结构。

如果已存在具有完全相同的子系统集的活动层次结构，则它将被重新用于新挂载。
如果现有层次结构不匹配，并且任何请求的子系统正在现有层次结构中使用，则挂载将失败并显示 - `EBUSY`。 
否则，将激活与所请求的子系统相关联的新层次结构。

目前无法将新子系统绑定到活动 cgroup 层次结构，或从活动 cgroup 层次结构取消子系统的绑定。 
这在未来可能是可能的，但充满了令人讨厌的错误恢复问题。

当 cgroup 文件系统被卸载时，如果在顶级 cgroup 下创建了任何子 cgroup，则即使已卸载，该层次结构仍将保持活动状态； 如果没有子 cgroup，则层次结构将被停用。

没有为 cgroup 添加新的系统调用 - 对查询和修改 cgroup 的所有支持都是通过此 cgroup 文件系统。

`/proc` 下的每个任务都有一个名为 “cgroup” 的添加文件，为每个活动层次结构显示子系统名称和 cgroup 名称作为相对于 cgroup 文件系统根的路径。

每个 cgroup 由 cgroup 文件系统中的一个目录表示，其中包含描述该 cgroup 的以下文件：

- `tasks`: 附加到该 cgroup 的任务列表（按 PID）。 不保证此列表已排序。
    将线程 ID 写入此文件会将线程移动到此 cgroup 中。
- `cgroup.procs`: cgroup 中的线程组 ID 列表。
    不保证此列表已排序或没有重复的 `TGID`，并且如果需要此属性，用户空间应对该列表进行排序/唯一化。
    将线程组 ID 写入此文件会将该组中的所有线程移动到此 cgroup 中。
- `notify_on_release`: 退出时运行 `release agent` ?
- `release_agent`： 用于发布通知的路径（此文件仅存在于顶级 cgroup 中）

其他子系统（例如 cpuset）可能会在每个 cgroup 目录中添加其他文件。

使用 `mkdir` 系统调用或 `shell` 命令创建新的 cgroup。
cgroup 的属性（例如其标志）可通过写入该 cgroups 目录中的相应文件来修改，如上面所列。

嵌套 cgroup 的命名分层结构允许将大型系统划分为嵌套的、动态可更改的 “软分区”。

每个任务（由该任务的任何子任务在 fork 时自动继承）附加到 cgroup 允许将系统上的工作负载组织成相关的任务集。
如果必要的 cgroup 文件系统目录的权限允许，任务可以重新附加到任何其他 cgroup。

当任务从一个 cgroup 移动到另一个 cgroup 时，它会获得一个新的 `css_set` 指针 - 如果已经存在具有所需 cgroup 集合的 `css_set`，则该组将被重用，否则会分配一个新的 `css_set`。
通过查看哈希表来找到适当的现有 `css_set`。

为了允许从 cgroup 访问组成它的 `css_sets`（以及任务），一组 `cg_cgroup_link` 对象形成一个网格； 
每个 `cg_cgroup_link` 都链接到其 `cgrp_link_list` 字段上的单个 cgroup 的 `cg_cgroup_links` 列表，以及其 `cg_link_list` 上的单个 `css_set` 的 `cg_cgroup_links` 列表。

因此，可以通过迭代引用该 cgroup 的每个 `css_set` 以及对每个 `css_set` 的任务集进行子迭代来列出 cgroup 中的任务集。

使用 Linux 虚拟文件系统 (vfs) 来表示 cgroup 层次结构为 cgroup 提供了熟悉的权限和名称空间，并且需要最少的附加内核代码。


## 1.4 What does `notify_on_release` do ?

如果在cgroup中启用了 `notify_on_release` 标志(`1`)，则每当cgroup中的最后一个任务离开（退出或附加到某个其他 cgroup ）并且该cgroup的最后一个子 cgroup 被删除时，
然后内核运行该层次结构根目录中 `release_agent` 文件内容指定的命令，提供废弃 cgroup 的路径名（相对于 cgroup 文件系统的挂载点）。
这可以自动删除废弃的 cgroup。
系统启动时 root cgroup 中的 `notification_on_release` 的默认值是禁用的 (`0`)。
其他 cgroup 在创建时的默认值是其父级的 `notify_on_release` 设置的当前值。
cgroup 层次结构的 `release_agent` 路径的默认值为空。


## 1.5 What does `clone_children` do ?

该标志仅影响 cpuset 控制器。
如果在 cgroup 中启用了 `clone_children` 标志 (`1`)，则新的 cpuset cgroup 将在初始化期间从父级复制其配置。


## 1.6 How do I use cgroups ?

要使用 “cpuset” cgroup 子系统启动包含在 cgroup 中的新工作，步骤如下：

```bash
mount -t tmpfs cgroup_root /sys/fs/cgroup
mkdir /sys/fs/cgroup/cpuset
mount -t cgroup -ocpuset cpuset /sys/fs/cgroup/cpuset
# 通过在 /sys/fs/cgroup/cpuset 虚拟文件系统中使用 mkdir 和 write (或 echo) 创建新的 cgroup
# 启动一项任务，该任务将成为新工作的“奠基人”
# 通过将该任务的 PID 写入该 cgroup 的 /sys/fs/cgroup/cpuset 任务文件，将该任务附加到新的 cgroup。
# fork, 执行或克隆该创始任务的工作任务
```

例如，以下命令序列将设置一个名为“Charlie”的 cgroup，仅包含 CPU 2 和 3 以及内存节点 1，然后在该 cgroup 中启动一个子 shell 'sh'：

```bash
mount -t tmpfs cgroup_root /sys/fs/cgroup
mkdir /sys/fs/cgroup/cpuset
mount -t cgroup cpuset -ocpuset /sys/fs/cgroup/cpuset

cd /sys/fs/cgroup/cpuset
mkdir Charlie
/bin/echo 2-3 > cpuset.cpus
/bin/echo 1 > cpuset.mems
/bin/echo $$ > tasks
sh
# 子shell `sh` 现在运行在 Charlie 控制组中
cat /proc/self/cgroup
# output: '/Charlie'
```


# 2. Usage Examples and Syntax


## 2.1 Basic Usage

创建、修改、使用 cgroup 可以通过 cgroup 虚拟文件系统来完成。

要挂载包含所有可用子系统的 cgroup 层次结构，输入：

`mount -t cgroup xxx /sys/fs/cgroup`

“xxx” 不由 cgroup 代码解释，但会出现在 `/proc/mounts` 中，因此可能是您喜欢的任何有用的标识字符串。

注意：如果没有一些用户输入，某些子系统将无法工作。
例如，如果启用了 cpuset，则用户必须先为创建的每个新 cgroup 填充 `cpus` 和 `mems` 文件，然后才能使用该组。

正如 “1.2 为什么需要 cgroup？” 一节中所解释的,您应该为要控制的每个资源或资源组创建不同的 cgroup 层次结构。
因此，您应该在 `/sys/fs/cgroup` 上挂载 tmpfs 并为每个 cgroup 资源或资源组创建目录。

```bash
mount -t tmpfs cgroup_root /sys/fs/cgroup
mkdir /sys/fs/cgroup/rg1
```

要安装仅包含 cpuset 和内存子系统的 cgroup 层次结构，输入：

`mount -t cgroup -o cpuset,memory hier1 /sys/fs/cgroup/rg1`

虽然目前支持重新挂载 cgroup，但不建议使用它。
重新挂载允许更改绑定的子系统和 `release_agent`。
重新绑定几乎没有用，因为它仅在层次结构为空且 `release_agent` 本身应替换为传统的 `fsnotify` 时才有效。
未来将取消对重新安装的支持。

指定层次结构的 `release_agent`：

```bash
mount -t cgroup -o cpuset,release_agent="/sbin/cpuset_release_agent" \
    xxx /sys/fs/cgroup/rg1
```

请注意，多次指定 `release_agent` 将返回失败。

请注意，当前仅当层次结构由单个（根）cgroup 组成时才支持更改子系统集。
支持从现有 cgroup 层次结构中任意绑定/取消绑定子系统的能力预计将在未来实现。

然后在 `/sys/fs/cgroup/rg1` 下你可以找到一棵与系统中cgroup的树相对应的树。
例如，`/sys/fs/cgroup/rg1` 是保存整个系统的 cgroup。

如果要更改 `release_agent` 的值：

`echo "/sbin/new_release_agent" > /sys/fs/cgroup/rg1/release_agent`

也可以通过重新挂载来更改。

如果要在 `/sys/fs/cgroup/rg1` 下创建一个新的 cgroup：

```bash
cd /sys/fs/cgroup/rg1
mkdir my_cgroup
```

现在可以操作这个 `cgroup`。

`cd my_cgroup`

在该目录下可以发现以下几个文件：

```bash
ls

cgroup.procs notify_on_release tasks
# (加上附加子系统添加的任何文件）。
```

现在将 shell 连接到 cgroup：

`/bin/echo $$ > tasks`

在此目录中还可以使用 mkdir 在你的 cgroup 中创建 cgroup。

`mkdir my_sub_cs`

要移除一个 cgroup, 用 `rmdir` 就可以：

`rmdir my_sub_cs`

如果 cgroup 正在使用（内部有 cgroup，或附加了进程，或由其他子系统特定的引用保持活动状态），则此操作将会失败。


## 2.2 Attaching processes

`/bin/echo PID > tasks`

注意，是PID，而不是PIDs。 您一次只能附加一项任务。
如果您有多个任务要附加，则必须一个接一个地执行：

```bash
/bin/echo PID1 > tasks
/bin/echo PID2 > tasks
# ...
/bin/echo PIDn > tasks
```

您可以通过`echo 0` 来附加当前的 shell 任务：

`echo 0 > tasks`

您可以使用 `cgroup.procs` 文件而不是任务文件来一次移动线程组中的所有线程。
将线程组中任何任务的 PID 回显到 `cgroup.procs` 会导致该线程组中的所有任务都附加到 cgroup。
将 0 写入 `cgroup.procs` 会移动写入任务的线程组中的所有任务。

注意：由于每个任务始终是每个已安装层次结构中一个 cgroup 的成员，因此要从当前 cgroup 中删除任务，您必须通过写入新 cgroup 的任务文件将其移至新 cgroup（可能是根 cgroup）。

注意：由于某些 cgroup 子系统强制执行的一些限制，将进程移动到另一个 cgroup 可能会失败。


## 2.3 Mounting hierarchies by name

挂载 cgroups 层次结构时传递 `name=<x>` 选项会将给定名称与层次结构关联起来。
这可以在安装预先存在的层次结构时使用，以便通过名称而不是通过其活动子系统集来引用它。
每个层次结构要么是无名的，要么具有唯一的名称。

名称应匹配 `[\w.-]+`

当为新层次结构传递 `name=<x>` 选项时，您需要手动指定子系统；
当您为子系统命名时，不支持在未显式指定任何子系统时安装所有子系统的旧行为。

子系统的名称作为层次结构描述的一部分出现在 `/proc/mounts` 和 `/proc/<pid>/cgroups` 中。


# 3. Kernel API


## 3.1 Overview

每个想要挂接到通用 cgroup 系统的内核子系统都需要创建一个 `cgroup_subsys` 对象。
其中包含各种方法，这些方法是来自 cgroup 系统的回调，以及将由 cgroup 系统分配的子系统 ID。

`cgroup_subsys` 对象中的其他字段包括：

- `subsys_id` 子系统的唯一数组索引，指示该子系统应管理 `cgroup->subsys[]` 中的哪个条目。
- `name`  应初始化为唯一的子系统名称。
    长度不得超过 `MAX_CGROUP_TYPE_NAMELEN`
- `early_init`  指示子系统是否需要在系统启动时提前初始化。

系统创建的每个cgroup对象都有一个指针数组，以子系统ID为索引； 该指针完全由子系统管理； 通用 cgroup 代码永远不会触及该指针。


## 3.2 Synchronization

有一个全局互斥体 `cgroup_mutex`，由 cgroup 系统使用。
任何想要修改 cgroup 的人都应该采取此操作。
也可以采取措施防止 cgroup 被修改，但在这种情况下更具体的锁定可能更合适。

查阅 [kernel/cgroup.c][] 了解更多。

子系统可以通过函数 `cgroup_lock()`/`cgroup_unlock()` 获取/释放 `cgroup_mutex`。

访问任务的 cgroup 指针可以通过以下方式完成：

- 持有 `cgroup_mutex` 时
- 同时持有任务的 `alloc_lock`（通过 `task_lock()` ）。
- 通过 `rcu_dereference()` 在 `rcu_read_lock()` 部分内


## 3.3 Subsystem API

每个子系统应该：

- 在 `linux/cgroup_subsys.h` 中添加条目
- 定义名为 `<name>_cgrp_subsys` 的 `cgroup_subsys` 对象

每个子系统可以导出以下方法。
唯一的强制方法是 `css_alloc/free`。
任何其他为空的操作都被认为是成功的空操作。

`struct cgroup_subsys_state *css_alloc(struct cgroup *cgrp)`
（`cgroup_mutex` 由调用者持有）

调用为 cgroup 分配子系统状态对象。
子系统应该为传递的 cgroup 分配其子系统状态对象，成功时返回指向新对象的指针或 `ERR_PTR()` 值。
成功后，子系统指针应指向 `cgroup_subsys_state` 类型的结构（通常嵌入到较大的子系统特定对象中），该结构将由 cgroup 系统初始化。
请注意，这将在初始化时被调用，以为此子系统创建根子系统状态； 
这种情况可以通过传递的具有 NULL 父级的 cgroup 对象来识别（因为它是层次结构的根），并且可能是初始化代码的适当位置。

`int css_online(struct cgroup *cgrp)`
（`cgroup_mutex` 由调用者持有）

在 @cgrp 成功完成所有分配并对 `cgroup_for_each_child/descendant_*()` 迭代器可见后调用。
子系统可以通过返回 -errno 选择创建失败。
此回调可用于实现沿层次结构的可靠状态共享和传播。
有关详细信息，请参阅 `cgroup_for_each_descendant_pre()` 上的注释。

`void css_offline(struct cgroup *cgrp);`
（`cgroup_mutex` 由调用者持有）

这是 `css_online()` 的对应部分，当且仅当 `css_online()` 在 @cgrp 上成功时调用。
这标志着@cgrp结束的开始。
@cgrp 正在被删除，子系统应该开始删除它在 @cgrp 上持有的所有引用。
当所有引用都被删除后，cgroup 删除将继续进行下一步 - `css_free()`。
在此回调之后，@cgrp 对于子系统来说应该被视为死亡。

`void css_free(struct cgroup *cgrp)`
（`cgroup_mutex` 由调用者持有）

cgroup系统即将释放@cgrp； 子系统应该释放其子系统状态对象。
当这个方法被调用的时候，@cgrp已经完全没有被使用了； `@cgrp->parent` 仍然有效。
（注意 - 如果在为新 cgroup 调用该子系统的 `create()` 方法后发生错误，也可以为新创建的 cgroup 调用该方法）。

`int can_attach(struct cgroup *cgrp, struct cgroup_taskset *tset)`
（`cgroup_mutex` 由调用者持有）

在将一项或多项任务移入 cgroup 之前调用； 如果子系统返回错误，这将中止附加操作。
@tset 包含要附加的任务，并保证其中至少有一个任务。

如果任务集中有多个任务，则：

- 保证所有线程都来自同一个线程组
- @tset 包含线程组中的所有任务，无论它们是否正在切换 cgroup
- 第一个任务是领导者

每个 @tset 条目还包含任务的旧 cgroup，并且可以使用 `cgroup_taskset_for_each()` 迭代器轻松跳过不切换 cgroup 的任务。
请注意，这不是在 fork 上调用的。
如果此方法返回 0（成功），那么当调用者持有 `cgroup_mutex` 时，该方法应该保持有效，并确保将来调用 `Attach()` 或 `cancel_attach()` 。

`void css_reset(struct cgroup_subsys_state *css)`
（`cgroup_mutex` 由调用者持有）

一个可选操作，应将 `@css` 的配置恢复到初始状态。
当前仅在通过 `cgroup.subtree_control` 在 cgroup 上禁用子系统时在统一层次结构上使用，但应保持启用状态，因为其他子系统依赖于它。
cgroup core 通过删除关联的界面文件来使此类 css 不可见，并调用此回调，以便隐藏子系统可以返回到初始中性状态。
这可以防止来自隐藏 CSS 的意外资源控制，并确保配置在稍后再次可见时处于初始状态。

`void cancel_attach(struct cgroup *cgrp, struct cgroup_taskset *tset)`
（`cgroup_mutex` 由调用者持有）

当 `can_attach()` 成功后任务附加操作失败时调用。
一个 `can_attach()` 有副作用的子系统应该提供这个函数，以便子系统可以实现回滚。 如果没有，则没有必要。
仅当子系统的 `can_attach()` 操作成功时才会调用此方法。
参数与 `can_attach()` 相同。

`void attach(struct cgroup *cgrp, struct cgroup_taskset *tset)`
（`cgroup_mutex` 由调用者持有）

在任务附加到 cgroup 后调用，以允许任何需要内存分配或阻塞的附加后活动。
参数与 `can_attach()` 相同。

`void fork(struct task_struct *task)`

当任务被分叉到 cgroup 时调用。

`void exit(struct task_struct *task)`

在任务退出期间调用。

`void free(struct task_struct *task)`

当`task_struct` 被释放时调用。

`void bind(struct cgroup *root)`
（`cgroup_mutex` 由调用者持有）

当 cgroup 子系统重新绑定到不同的层次结构和根 cgroup 时调用。
目前，这只涉及默认层次结构（从不具有子 cgroup）和正在创建/销毁的层次结构（因此没有子 cgroup）之间的移动。


# 4. Extended attributes usage

cgroup 文件系统在其目录和文件中支持某些类型的扩展属性。
目前支持的类型有：

- `Trusted (XATTR_TRUSTED)`
- `Security (XATTR_SECURITY)`

两者都需要 `CAP_SYS_ADMIN` 功能来设置。

与 tmpfs 一样，cgroup 文件系统中的扩展属性是使用内核内存存储的，建议将使用量保持在最低限度。
这就是不支持用户定义的扩展属性的原因，因为任何用户都可以执行此操作，并且值大小没有限制。

当前使用此功能的已知用户是 SELinux，用于限制容器中 cgroup 的使用，以及 systemd 用于存储各种元数据，例如 cgroup 中的主 PID（systemd 为每个服务创建一个 cgroup）。

# 5. Questions

Q: 为什么是用 `/bin/echo`?
A: bash 的内置“echo”命令不会检查对 write() 的调用是否有错误。
如果您在 cgroup 文件系统中使用它，您将无法判断命令是成功还是失败。

Q: 当我附加进程时，只有第一行真正附加！
A: 每次调用 write() 只能返回一个错误代码。所以你也应该只输入一个 PID。

