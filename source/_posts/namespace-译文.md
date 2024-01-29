---
title: namespace 译文
date: 2023-12-26 14:13:26
tags:
categories:
- man-pages
keywords:
- namespaces
copyright: Guader
copyright_author_href:
copyright_info:
---

[原文](https://man7.org/linux/man-pages/man7/namespaces.7.html)

## 名称

namespaces - Linux命名空间概述


## 描述

命名空间将全局系统资源封装在一个抽象概念中。在命名空间内的进程来看，它们拥有各自独立的全局资源实例。全局资源的更改对该命名空间内的其他进程可见，但对其他的进程不可见。命名空间的一个用途就是实现容器。

本页提供了有关各种命名空间类型信息的指针，描述了相关的 `/proc` 文件，并总结了用于处理命名空间的API。


### 命名空间类型

下表列出了Linux上可用的命名空间类型。表中第二列展示了在各种API中用于指定命名空间类型的标志值。第三列标识了提供命名空间类型详细信息的手册页面。最后一列是按命名空间类型隔离的资源摘要。

命名空间 | 标志 | man page | 隔离 
---|---|---|---
CGroup | `CLONE_NEWCGROUP` | `cgroup_namespaces(7)` | Cgroup root directory
IPC | `CLONE_NEWIPC` | `ipc_namespaces(7)` | System V IPC, POSIX message queues
Network | `CLONE_NEWNET` | `network_namespaces(7)` | Network devices, stacks, ports, etc.
Mount | `CLONE_NEWNS` | `mount_namespaces(7)` | Mount points
PID | `CLONE_NEWPID` | `pid_namespaces(7)` | Process IDs
Time | `CLONE_NEWTIME` | `time_namespaces(7)` | Boot and monotonic clocks
User | `CLONE_NEWUSER` | `user_namespaces(7)` | User and group IDs
UTS | `CLONE_NEWUTS` | `uts_namespaces(7)` | Hostname and NIS domain name



### 命名空间API

除了下面描述的各种 `/proc` 文件外，命名空间API包括以下系统调用：


- **`clone(2)`**  `clone(2)` 系统调用用于创建一个新的进程。如果调用的参数`flags`指定了上述的一个或多个`CLONE_NEW*`标志，那么将为每个标志创新新的命名空间，且子进程将成为这些命名空间的成员。（这个系统调用同时实现了一些与命名空间无关的功能。）
- **`setns(2)`**  `setns(2)` 系统调用允许调用的进程加入到现有的命名空间。要加入的命名空间通过下面要说的 `/proc/pid/ns` 文件的文件描述符指定。
- **`unshare(2)`**  `unshare(2)`系统调用将调用的进程移动到一个新的命名空间。如果调用的参数`flags`指定了上出的`CLONE_NEW*`的标志中的一个或多个，那么会为每个标志创建新的命名空间，同时调用的进程也会是这些命名空间的成员。（这个系统调用同时实现了一些与命名空间无关的功能。）
- **`ioctl(2)`**  各种`ioctl(2)`操作可用于发现命名空间的信息。这些操作在`ioctl_ns(2)`中有描述。


大多数情况下，使用`clone(2)`和`unshare(2)`创建新的命名空间需要 **`CAP_SYS_ADMIN`** 能力，因为在新命名空间中，创建者有权更改对随后创建或加入到该命名空间的进程可见的全局资源。用户命名空间是个例外：自 Linux 3.8 起，创建用户命名空间不需要权限。



### `/proc/pid/ns/` 文件夹

每个进程都有一个 `/proc/pid/ns/` 子文件夹，其中包含一个由`setns(2)`操作的命名空间的条目。

```
ls -l /proc/$$/ns | awk '{print $1, $9, $10, $11}'

total   
lrwxrwxrwx cgroup -> cgroup:[4026531835]
lrwxrwxrwx ipc -> ipc:[4026531839]
lrwxrwxrwx mnt -> mnt:[4026531841]
lrwxrwxrwx net -> net:[4026531840]
lrwxrwxrwx pid -> pid:[4026531836]
lrwxrwxrwx pid_for_children -> pid:[4026531836]
lrwxrwxrwx time -> time:[4026531834]
lrwxrwxrwx time_for_children -> time:[4026531834]
lrwxrwxrwx user -> user:[4026531837]
lrwxrwxrwx uts -> uts:[4026531838]
```

将该目录中的一个文件绑定挂载(见 `mount(2)`) 到文件系统中的其他地方，即使当前在该命名空间中的所有进程都已终止，也能保持pid 指定的进程的相应命名空间继续存在。

打开该文件夹中的一个文件（或绑定挂载到这些文件中的一个文件），将返回 pid 指定进程的相应命名空间的文件句柄。只要该文件描述符保持打开状态，即使命名空间中的所有进程都已终止，命名空间也将保持存活。 文件描述符可以传递给 `setns(2)`。

在 Linux 3.7 及更早版本中，这些文件以硬链接的形式显示。自 Linux 3.8 起，它们显示为符号链接。如果两个进程处于同一命名空间，那么它们的`/proc/pid/ns/xxx` 符号链接的设备 `ID` 和 `inode` 编号将是相同的；应用程序可以使用 `stat(2)` 返回的 `stat.st_dev` 和 `stat.st_ino` 字段检查这一点。 该符号链接的内容是一个包含命名空间类型和 inode 编号的字符串，如下例所示的字符串：

```
readlink /proc/$$/ns/nts

uts:[4026531838]
```

该子文件夹中的符号链接如下：

- `/proc/pid/ns/cgroup` 自 Linux 4.6 起。 该文件是cgroup命名空间进程的句柄。
- `/proc/pid/ns/ipc` 自 Linux 3.0 起。 该文件是IPC命名空间进程的句柄。
- `/proc/pid/ns/mnt` 自 Linux 3.8 起。 该文件是mount命名空间进程的句柄。
- `/proc/pid/ns/net` 自 Linux 3.0 起。 该文件是network命名空间进程的句柄。
- `/proc/pid/ns/pid` 自 Linux 3.8 起。 该文件是PID命名空间的句柄。该句柄在进程的生命周期内永久有效（即进程的PID命名空间成员资格永不改变）
- `/proc/pid/ns/pid_for_children` 自 Linux 4.12 起。该文件是这个进程创建的子进程的PID命名空间的句柄。这可能会因调用 `unshare(2)` 和 `setns(2)` 而发生变化（参见 `pid_namespaces(7)`）,因此文件可能与 `/proc/pid/ns/pid` 不同。只有在命名空间中创建了第一个子进程后，符号链接才会获得值。(在此之前，符号链接的 `readlink(2)` 将返回一个空缓冲区）。
- `/proc/pid/ns/time` 自 Linux 5.6 起。 该文件是time命名空间进程的句柄。
- `/proc/pid/ns/time_for_children` 自 Linux 5.6 起。 该文件是该进程创建的子进程的time命名空间的句柄。 调用 `unshare(2)` 和 `setns(2)` 后，time命名空间可能会发生变化（参见 `time_namespaces(7)`），因此该文件可能与 `/proc/pid/ns/time` 不同。
- `/proc/pid/ns/user` 自 Linux 3.8 起。 该文件是user命名空间进程的句柄。
- `/proc/pid/ns/nts` 自 Linux 3.0 起。该文件是UTS命名空间进程的句柄。

取消引用或读取（`readlink(2)`） 这些符号链接的权限受 `ptrace` 访问模式 **`PTRACE_MODE_READ_FSCREDS`** 检查的制约；请参阅 `ptrace(2)`。



### `/proc/sys/user` 文件夹

`/proc/sys/user` 文件夹中的文件（自 Linux 4.9 起就存在）暴露了对可创建的各种命名空间数量的限制。 这些文件如下

- `max_cgroup_namespaces`  该文件中的值定义了每个用户在用户命名空间中创建的 cgroup 命名空间的数量限制。
-  `max_ipc_namespaces`  该文件中的值定义了每个用户在用户命名空间中创建 ipc 命名空间的数量限制。
- `max_mnt_namespaces`  该文件中的值定义了每个用户在用户命名空间中创建的mount命名空间的数量限制。
- `max_net_namespaces`  该文件中的值定义了每个用户在用户命名空间中创建network命名空间的数量限制。
- `max_pid_namespaces`  该文件中的值定义了每个用户在用户命名空间中创建 PID 命名空间的数量限制。
- `max_time_namespaces` （自 Linux 5.7 起） 该文件中的值定义了每个用户在用户命名空间中创建time命名空间的数量限制。
- `max_user_namespaces`  该文件中的值定义了在用户命名空间中创建的每个user命名空间的数量限制。
- `max_uts_namespaces`  该文件中的值定义了每个用户在用户命名空间中创建 uts 命名空间的数量限制。


请注意这些文件的以下详细信息:

- 权限进程可以修改这些文件中的值。
- 这些文件显示的值是打开进程所在用户命名空间的限制。
- 这些限制以用户为单位。 同一用户命名空间中的每个用户可创建的命名空间不超过所定义的限制。
- 这些限制适用于所有用户，包括 UID 0。
- 这些限制除了适用于任何其他每个命名空间的的限制（如 PID 和用户命名空间的限制）之外，这些限制还可能被强制执行。
- 遇到这些限制时，`clone(2)` 和 `unshare(2)` 将失败，并显示错误 **`ENOSPC`** 。
- 对于初始用户命名空间，每个文件中的默认值为可创建线程数限制（`/proc/sys/kernel/threads-max`）的一半。 在所有后代用户命名空间中，每个文件的默认值都是 `MAXINT`。
- 在创建命名空间时，该对象也会与祖先命名空间进行比较。更准确地说;
  - 每个user命名空间都有创建者UID。
  - 当创建一个命名空间时，将根据每个祖先用户命名空间中的创建者 UID 进行核算, 内核会确保不超过祖先命名空间中创建者 UID 的相应命名空间限制。
  - 上述要点确保了创建新用户命名空间不能被用作逃避当前用户命名空间中有效限制的手段



### 命名空间寿命

如果没有其他因素，当命名空间中的最后一个进程终止或离开命名空间时，命名空间就会被自动删除。不过，还有一些其他因素可能会使命名空间在没有成员进程的情况下仍然存在。 这些因素包括：

- 相应的 `/proc/pid/ns/*` 文件存在打开的文件描述符或绑定挂载。
- 命名空间是分层的（即 PID 或用户命名空间），并有一个子命名空间。
- 用户命名空间拥有一个或多个非用户命名空间。
- 这是一个 PID 命名空间，有一个进程通过 `/proc/pid/ns/pid_for_children` 符号链接指向该命名空间。
- 这是一个时间命名空间，有一个进程通过 `/proc/pid/ns/time_for_children` 符号链接引用该命名空间。
- 这是一个 IPC 命名空间，mqueue 文件系统（参见 `mq_overview(7)`）的相应挂载指向该命名空间。
- 这是一个 PID 命名空间，`proc(5)` 文件系统的相应挂载会指向该命名空间。



## 例子

参见 `clone(2)`和 `user_namespaces(7)`



## 其他

`nsenter(1)`, `readlink(1)`, `unshare(1)`, `clone(2)`, `ioctl_ns(2)`, `setns(2)`, `unshare(2)`, `proc(5)`, `capabilities(7)`, `cgroup_namespaces(7)`, `cgroups(7)`, `credentials(7)`, `ipc_namespaces(7)`, `network_namespaces(7)`, `pid_namespaces(7)`,`user_namespaces(7)`, `uts_namespaces(7)`, `lsns(8)`, `switch_root(8)`

