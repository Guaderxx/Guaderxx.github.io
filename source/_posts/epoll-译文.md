---
title: epoll-译文
date: 2024-02-25 17:09:55
tags:
categories:
- Doc
keywords:
- epoll
copyright: Guader
copyright_author_href:
copyright_info:
---

epoll(7) -- Linux Manual Page

[Source][epoll-7]


## Name

epoll  -- I/O 事件通知工具


## Synopsis

```c
#include <sys/epoll.h>
```


## Description

**epoll** API 执行的任务与 [poll(2)][poll-2] 类似：监控多个文件描述符以查看是否有任何可以进行I/O操作的。
**epoll** API 可以用作边缘触发或水平触发接口，并且可以很好地扩展到大量被监控的文件描述符。

**epoll** API 的核心概念是 **epoll** *实例* ，这是一个在内核中的数据结构，从用户空间的角度来看，可以将其视为包含两个列表的容器：

- *兴趣* 列表（有时也称为 **epoll** 集合）： 进程已注册感兴趣监控的文件描述符集合。
- *准备* 列表： 一组 “准备就绪” 用于 I/O 的文件描述符。准备列表是兴趣列表中文件描述符的一个子集（或更准确地说，是一组指向这些文件描述符的引用）。
  准备列表是由内核动态填充的，这是由于这些文件描述符上的 I/O 活动所致。

以下系统调用用于创建和管理 **epoll** 实例：

- [`epoll-create(2)`][epoll-create-2]: 创建一个新的epoll实例，并返回一个引用该实例的文件描述符。
  （新的 [`epoll_create1(2)`][epoll-create1-2] 函数扩展了 [`epoll_create(2)`][epoll-create-2] 的功能。）
- 然后通过 [`epoll_ctl(2)`][epoll-ctl-2] 注册对特定文件描述符的兴趣，将项目添加到 epoll 实例的兴趣列表中。
- [`epoll_wait(2)`][epoll-wait-2] 等待 I/O 事件，如果当前没有事件可用，它会阻塞调用线程。
  （这个系统调用可以被认为是从 **epoll** 实例的就绪列表中获取项。）


### Level-triggered and edge-triggered

**epoll** 事件分发接口既可以作为边缘触发（ET），也可以作为水平触发（LT）。
两种机制之间的差异可以描述如下。
假设发生这种情况：

1. 表示管道读取端的文件描述符 (`rfd`) 已注册到 **epoll** 实例上。
2. 管道写入器在管道写入端写入 2 kB 数据。
3. 对 [`epoll_wait(2)`][epoll-wait-2] 的调用完成，将返回 `rfd` 作为 准备好的文件描述符。
4. 管道读取器从 `rfd` 读取 1 kB 数据。
5. 调用 [`epoll_wait(2)`][epoll-wait-2] 。

如果 `rfd` 文件描述符已使用 **EPOLLET**（边缘触发）标志添加到 **epoll** 接口，那么在步骤5中执行的 `epoll_wait(2)` 调用可能会挂起，尽管文件输入缓冲区中仍有可用的数据；同时，远程对等方可能正在期待基于它已发送的数据的响应。
这是因为边缘触发模式仅在监控的文件描述符发生变化时才传递事件。
因此，在步骤5中，调用者可能会最终等待已经在输入缓冲区中的数据。
在上面的例子中，由于步骤2中的写入操作，会在 `rfd` 上生成一个事件，该事件在步骤3中被消耗。
由于步骤4中的读取操作没有消耗完整个缓冲区数据，步骤5中对 [`epoll_wait(2)`][epoll-wait-2] 的调用可能会无限期地阻塞。

一个使用 **EPOLLET** 标志的应用程序应使用非阻塞文件描述符，以避免在处理多个文件描述符时，一个阻塞的读或写操作使任务饥饿。
建议以下面的方式使用 **epoll** 作为边缘触发（**EPOLLET**）接口：

1. 使用非阻塞文件描述符；以及
2. 通过在 [`read(2)`][read-2] 或 [`write(2)`][write-2] 返回 **EAGAIN** 后等待事件。

相比之下，当作为水平触发接口使用时（默认情况，未指定 **EPOLLET** 时），**epoll** 仅仅是一个更快的 [`poll(2)`][poll-2] ，并且可以在任何使用 [`poll(2)`][poll-2] 的地方使用，因为它们共享相同的语义。

即使在边缘触发的 **epoll** 中，由于接收到多个数据块可能会产生多个事件，调用者可以选择指定 **EPOLLONESHOT** 标志，告诉 **epoll** 在通过 [`epoll_wait(2)`][epoll-wait-2] 接收到一个事件后禁用相关的文件描述符。
当指定了 **EPOLLONESHOT** 标志时，调用者有责任使用带有 **`EPOLL_CTL_MOD`** 的 [`epoll_ctl(2)`][epoll-ctl-2] 重新启用该文件描述符。

如果多个线程（或进程，如果子进程通过 [`fork(2)`][fork-2] 继承了 **epoll** 文件描述符）在 [`epoll_wait(2)`][epoll-wait-2] 中因等待同一个 **epoll** 文件描述符中的兴趣列表里的文件描述符就绪而阻塞，而这个文件描述符被标记为边缘触发（**EPOLLET**）通知，那么只有一个线程（或进程）会被 [`epoll_wait(2)`][epoll-wait-2] 唤醒。
这在某些情况下提供了避免“惊群效应”唤醒的有用优化。


### Interaction with autosleep

如果系统通过 `/sys/power/autosleep` 进入 **自动睡眠** 模式，并且有一个事件发生将设备从睡眠状态唤醒，设备驱动程序将只保持设备唤醒直到该事件被排队。
为了在事件处理完毕之前保持设备唤醒状态，需要使用 [`epoll_ctl(2)`][epoll-ctl-2] 的 **EPOLLWAKEUP** 标志。

当在 `struct epoll_event` 的结构中的 **events** 字段设置 **EPOLLWAKEUP** 标志时，系统将从事件被排队的那一刻起保持清醒状态，直到通过 [`epoll_wait(2)`][epoll-wait-2] 调用返回该事件，并持续到随后的 [`epoll_wait(2)`][epoll-wait-2] 调用。
如果事件需要在那个时间点之后继续保持系统清醒，那么应该在第二个 [`epoll_wait(2)`][epoll-wait-2] 调用之前获取一个单独的 `wake_lock` 。


### /proc interfaces

以下接口可用于限制epoll消耗的内核内存量：

- `/proc/sys/fs/epoll/max_user_watches`  自 Linux 2.6.28 起
  这指定了用户可以在系统上所有epoll实例中注册的文件描述符总数的限制。
  该限制是基于每个真实用户ID。
  在32位内核上，每个注册的文件描述符大约需要90字节，在64位内核上则大约需要160字节。
  目前，`max_user_watches` 的默认值是可用低内存的 1/25（4%），再除以注册成本（以字节为单位）。
  

### Example for suggested usage

尽管当 **epoll** 作为水平触发接口使用时与 [`poll(2)`][poll-2] 具有相同的语义，但边缘触发使用需要更多澄清以避免应用程序事件循环中的停滞。
在这个例子中，`listener` 是一个已经调用了 [`listen(2)`][listen-2] 的非阻塞套接字。
函数 `do_use_fd()` 使用新的就绪文件描述符，直到 [`read(2)`][read-2] 或 [`write(2)`][write-2] 返回 **EAGAIN** 。
事件驱动的状态机应用程序在收到 **EAGAIN** 后，应该记录其当前状态，以便在下一次调用 `do_use_fd()` 时，它可以从之前停止的地方继续 [`read(2)`][read-2] 或 [`write(2)`][write-2] 。

```c
#define MAX_EVENTS 10
struct epoll_event ev, events[MAX_EVENTS];
int listen_sock, conn_sock, nfds, epollfd;

/* Code to set up listening socket, 'listen_sock',
    (socket(), bind(), listen()) omitted. */

epollfd = epoll_create1(0);
if (epollfd == -1) {
    perror("epoll_create1");
    exit(EXIT_FAILURE);
}

ev.events = EPOLLIN;
ev.data.fd = listen_sock;
if (epoll_ctl(epollfd, EPOLL_CTL_ADD, listen_sock, &ev) == -1) {
    perror("epoll_ctl: listen_sock");
    exit(EXIT_FAILURE);
}

for (;;) {
    nfds = epoll_wait(epollfd, events, MAX_EVENTS, -1);
    if (nfds == -1) {
        perror("epoll_wait");
        exit(EXIT_FAILURE);
    }
    
    for (n = 0; n < nfds; ++n) {
        if (events[n].data.fd == listen_sock) {
            conn_sock = accept(listen_sock,
                (struct sockaddr *)&addr, &addrlen);
            if (conn_sock == -1) {
                perror("accept");
                exit(EXIT_FAILURE);
            }
            
            setnonblocking(conn_sock);
            ev.events = EPOLLIN | EPOLLET;
            ev.data.fd = conn_sock;
            if (epoll_ctl(epollfd, EPOLL_CTL_ADD, conn_sock, &ev) == -1) {
                perror("epoll_ctl: conn_sock");
                exit(EXIT_FAILURE);
            }
        } else {
            do_use_fd(events[n].data.fd);
        }
    }
}
```

当作为边缘触发的接口使用时，出于性能考虑，可以通过指定（ **EPOLLIN** | **EPOLLOUT** ）一次性将文件描述符添加到 **epoll** 接口（ **`EPOLL_CTL_ADD`** ）中。
这允许你避免不断地在 **EPOLLIN** 和 **EPOLLOUT** 之间切换，通过使用 **`EPOLL_CTL_MOD`** 调用 [`epoll_ctl(2)`][epoll-ctl-2] 。



## Questions and Answers

- <details>
    <summary>用于区分兴趣列表中注册的文件描述符的关键是什么？ </summary>
    关键是文件描述符编号和打开的文件描述（也称为“打开的文件句柄”，是内核内部对打开文件的表示）的组合。
  </details>
- <details>
    <summary>如果在 <strong>epoll</strong> 实例上两次注册同一个文件描述符会发生什么？ </summary>
    
    你可能会得到 **EEXIST** 错误。
    然而，可以将一个重复的文件描述符（通过 [`dup(2)`][dup-2] 、[`dup2(2)`][dup2-2] 、[`fcntl(2)`][fcntl-2] `F_DUPFD` 创建）添加到同一个 **epoll** 实例中。
    如果这些重复的文件描述符使用不同的事件掩码注册，这可以成为过滤事件的有用技巧。
  </details>
- <details>
    <summary>两个 <strong>epoll</strong> 实例可以等待同一个文件描述符吗？如果是的话，事件会被报告给两个 <strong>epoll</strong> 文件描述符吗？ </summary>
    是的，事件会被报告给两个。但是，可能需要谨慎编程才能正确地做到这一点。
  </details>
- <details>
    <summary> <strong>epoll</strong> 文件描述符本身是可轮询/epoll/选择的？ </summary>
    
    是的，如果 **epoll** 文件描述符有等待的事件，那么它将表现为可读。
  </details>
- <details>
    <summary>如果尝试将一个 <strong>epoll</strong> 文件描述符放入它自己的文件描述符集合中会发生什么？ </summary>
    
    [`epoll_ctl(2)`][epoll-ctl-2] 调用会失败（**EINVAL**）。
    但是，你可以在另一个 **epoll** 文件描述符集合中添加一个 **epoll** 文件描述符。
  </details>
- <details>
    <summary>我可以将一个 <strong>epoll</strong> 文件描述符通过 UNIX 域套接字发送到另一个进程吗？ </summary>
    是的，你可以这样做，但这并没有意义，因为接收进程不会拥有兴趣列表中的文件描述符的副本。
  </details>
- <details>
    <summary>关闭文件描述符会导致它从所有的 <strong>epoll</strong> 兴趣列表中移除吗？</summary>
   
   是的，但请注意以下要点。
    文件描述符是对打开文件描述的引用（参见 [`open(2)`][open-2] ) 。
    每当通过 [`dup(2)`][dup-2] 、[`dup2(2)`][dup2-2] 、[`fcntl(2)`][fcntl-2] `F_DUPFD` 或 [`fork(2)`][fork-2] 复制文件描述符时，都会创建一个新的引用相同打开文件描述的文件描述符。
    打开文件描述会一直存在，直到所有引用它的文件描述符都被关闭。

    只有当所有引用底层打开文件描述的文件描述符都被关闭后，文件描述符才会从兴趣列表中移除。这意味着即使兴趣列表中的文件描述符已经被关闭，如果其他引用相同底层文件描述的文件描述符仍然打开，仍然可能会报告该文件描述符的事件。为了防止这种情况发生，必须在复制之前，使用 [`epoll_ctl(2)`][epoll-ctl-2] `EPOLL_CTL_DEL` 明确地从兴趣列表中移除文件描述符。或者，应用程序必须确保所有文件描述符都已关闭（如果文件描述符在幕后被使用 [`dup(2)`][dup-2] 或 [`fork(2)`][fork-2] 的库函数复制，这可能比较困难）。
  </details>
- <details>
    <summary>如果在 epoll_wait(2) 调用之间发生多个事件，它们会被合并还是分开报告？ </summary>
    它们会被合并。
  </details>
- <details>
    <summary>对文件描述符的操作是否会影响已经存在的文件描述符已收集但尚未报告的事件？</summary>
    您可以对现有文件描述符执行两个操作。
    对于这种情况，删除是没有意义的。修改后会重读可用的 I/O。
  </details>
- <details>
    <summary>当使用 <strong>EPOLLET</strong> 标志时 （边缘触发行为），我是否需要持续 读/写 文件描述符直到 <strong>EAGAIN</strong>？ </summary>
    
    从 [`epoll_wait(2)`][epoll-wait-2] 接收到事件应该会让你知道该文件描述符已准备好进行请求的 I/O 操作。你必须认为它已经准备好，直到下一次（非阻塞）读写返回 **EAGAIN** 。何时以及如何使用文件描述符完全取决于你。

   对于面向数据包/令牌的文件（例如，数据报套接字，处于规范模式的终端），检测读写 I/O 空间结束的唯一方法是不停地进行读写直到遇到 **EAGAIN** 。

   对于面向流的文件（例如，管道，`FIFO` ，流套接字），通过检查从/写入目标文件描述符的数据量，也可以检测到读写I/O空间已耗尽的条件。例如，如果你调用 [`read(2)`][read-2] 请求读取一定量的数据，而 [`read(2)`][read-2] 返回的字节数较少，你可以确信已经耗尽了该文件描述符的读取I/O空间。在使用 [`write(2)`][write-2] 写入时也是如此。（如果你不能保证监控的文件描述符总是指向面向流的文件，请避免使用后者技术。）
  </details>
  

## Possible pitfalls and ways to avoid them

- Starvation (edge-triggered)
    如果存在大量的I/O空间，试图清空它可能会导致其他文件无法得到处理，从而造成资源饥饿。（这个问题并不仅限于 **epoll** 。）

    解决方案是维护一个就绪列表，并在其关联的数据结构中将文件描述符标记为就绪，这样应用程序就能记住哪些文件需要处理，同时仍然在所有就绪文件之间进行轮询。这也支持忽略那些已经就绪的文件描述符后续收到的事件。
    
- If using an event cache...
    如果您使用事件缓存或存储从 [`epoll_wait(2)`][epoll-wait-2] 返回的所有文件描述符，请确保提供一种动态标记其关闭的方法（即由先前事件的处理引起的）。假设您从 [`epoll_wait(2)`][epoll-wait-2] 接收到 100 个事件，在事件 #47 中某个条件导致事件 #13 被关闭。如果您移除事件 #13 的结构并使用 [`close(2)`][close-2] 关闭文件描述符，那么您的事件缓存可能仍然会说该文件描述符有待处理的事件，从而造成混淆。
    解决此问题的方法之一是在处理事件 47 期间调用 **`epoll_ctl(EPOLL_CTL_DEL)`** 删除文件描述符 13 并执行 [`close(2)`][close-2] ，然后将关联的数据结构标记为已移除并链接到清理列表中。如果在您的批量处理中找到另一个针对文件描述符 13 的事件，您会发现该文件描述符已被先前移除，这样就不再会有混淆。




## Versions

其他一些系统也提供类似的机制；例如,
FreeBSD 有 `kqueue` ，Solaris 有 `/dev/poll` 。


## Notes

可以通过查看进程中 `/proc/pid/fdinfo` 目录下对应 epoll 文件描述符的条目来查看通过 epoll 文件描述符被监控的文件描述符集合。有关更多详细信息，请参见 [`proc(5)`][proc-5] 。

可以使用 [`kcmp(2)`][kcmp-2] 的 **`KCMP_EPOLL_TFD`** 操作来测试一个文件描述符是否存在于 epoll 实例中。



[epoll-7]: https://man7.org/linux/man-pages/man7/epoll.7.html 
[poll-2]: https://man7.org/linux/man-pages/man2/poll.2.html
[epoll-create-2]: https://man7.org/linux/man-pages/man2/epoll_create.2.html
[epoll-create1-2]: https://man7.org/linux/man-pages/man2/epoll_create1.2.html
[epoll-ctl-2]: https://man7.org/linux/man-pages/man2/epoll_ctl.2.html
[epoll-wait-2]: https://man7.org/linux/man-pages/man2/epoll_wait.2.html
[read-2]: https://man7.org/linux/man-pages/man2/read.2.html
[write-2]: https://man7.org/linux/man-pages/man2/write.2.html
[fork-2]: https://man7.org/linux/man-pages/man2/fork.2.html
[listen-2]: https://man7.org/linux/man-pages/man2/listen.2.html
[dup-2]: https://man7.org/linux/man-pages/man2/dup.2.html
[dup2-2]: https://man7.org/linux/man-pages/man2/dup2.2.html
[fcntl-2]: https://man7.org/linux/man-pages/man2/fcntl.2.html
[close-2]: https://man7.org/linux/man-pages/man2/close.2.html
[proc-5]: https://man7.org/linux/man-pages/man5/proc.5.html
[kcmp-2]: https://man7.org/linux/man-pages/man2/kcmp.2.html
[open-2]: https://man7.org/linux/man-pages/man2/open.2.html


