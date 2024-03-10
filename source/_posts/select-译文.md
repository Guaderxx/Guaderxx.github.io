---
title: select-译文
date: 2024-02-29 17:05:22
tags:
categories:
- Doc
keywords:
- select
copyright: Guader
copyright_author_href:
copyright_info:
---

select(2) -- System Calls Manual

[Source][select-2]


## Name

`select`, `pselect`, `FD_CLR`, `FD_ISSET`, `FD_SET`, `FD_ZERO`, `fd_set`  --  同步 I/O 复用


## Library

C标准库 （ `libc`, `-lc` ）


## Synopsis

```c
#include <sys/select.h>

typedef /* ... */ fd_set;

int select(int nfds, fd_set *_Nullable restrict readfds, 
    fd_set * _Nullable restrict writefds, 
    fd_set * _Nullable restrict exceptfds,
    struct timeval * _Nullable restrict timeout);

void FD_CLR(int fd, fd_set *set);
int FD_ISSET(int fd, fd_set *set);
void FD_SET(int fd, fd_set *set);
void FD_ZERO(fd_set *set);

int pselect(int nfds, fd_set *_Nullable restrict readfds,
    fd_set * _Nullable restrict writefds,
    fd_set * _Nullable restrict exceptfds,
    const struct timespec * _Nullable restrict timeout,
    const sigset_t * _Nullable restrict sigmask);
```

glibc 的功能测试宏要求 （见 [`feature_test_macros(7)`][feature-test-macro-7] ）：

```
pselect():
    _POSIX_C_SOURCE >= 200112L
```


## Description

**警告** ： **`select()`** 只能监视小于 **`FD_SETSIZE`** （1024）的文件描述符数量 —— 这对许多现代应用程序来说是一个低到不合理的限制，并且这个限制不会改变。
所有现代应用程序应该改用 [`poll(2)`][poll-2] 或 [`epoll(7)`][epoll-7] ，它们不受此限制。

**`select()`** 允许程序监视多个文件描述符，直到一个或多个文件描述符变为可以进行某些 I/O 操作类别（例如，可能进行输入）的 “就绪” 状态。如果可以不阻塞地执行相应的 I/O 操作（例如，[`read(2)`][read-2] ，或者足够小的 [`write(2)`][write-2] ），则文件描述符被认为是就绪的。


### `fd_set`

一种可以表示文件描述符集合的结构类型。
根据 POSIX 标准，`fd_set` 结构中文件描述符的最大数量是由宏 **`FD_SETSIZE`** 的值决定的。


### File descriptors sets

**`select()`** 函数的主要参数是三个 “集合” 的文件描述符（使用类型 `fd_set` 声明），它们允许调用者等待在指定的一组文件描述符上发生的三类事件。每个 `fd_set` 参数如果不需要监视相应类别的事件，可以指定为 NULL。

**请注意** ：返回后，每个文件描述符集合都会被就地修改，以指明哪些文件描述符当前是 “就绪” 的。因此，如果在循环中使用 **`select()`** ，必须在每次调用之前 **重新初始化** 集合。

文件描述符集的内容可以使用以下宏进行操作：

- **`FD_ZERO()`**
  这个宏用于清除（从集合中移除所有文件描述符）。
  它应该作为初始化文件描述符集的第一步使用。
- **`FD_SET()`**
  此宏将文件描述符 *fd* 添加到集合中。向集合中添加已经存在的文件描述符是一个空操作，不会产生错误。
- **`FD_CLR()`**
  这个宏从集合中移除文件描述符 *fd* 。移除一个不在集合中的文件描述符是空操作，不会产生错误。
- **`FD_ISSET()`**
  **`select()`** 根据下面描述的规则修改集合的内容。调用 **`select()`** 之后，可以使用 **`FD_ISSET()`** 宏来测试文件描述符是否仍存在于集合中。如果文件描述符 `fd` 在集合 `set` 中， **`FD_ISSET()`** 返回非零值；如果不在，则返回零。


### Arguments

**`select()`** 的参数如下：

- `readfds`
  在此集合中的文件描述符被监控以判断它们是否准备好读取。一个文件描述符如果进行读取操作不会阻塞，则认为是准备好读取的；特别是，在文件结束时文件描述符也会准备好。
  在 **`select()`** 函数返回之后，`readfds` 中将清除所有除了准备好读取之外的文件描述符。
- `writefds`
  这个集合中的文件描述符被监控以查看它们是否准备好写入。如果写入操作不会阻塞，则文件描述符就准备好写入。然而，即使文件描述符显示为可写，较大的写入操作仍可能会阻塞。
  在 **`select()`** 返回后，`writefds` 将除了准备好写入的文件描述符之外的所有文件描述符清除。
- `exceptfds`
  在此集合中的文件描述符被监控以检测 “异常条件”。关于某些异常条件的示例，请参见 [`poll(2)`][poll-2] 中对 **POLLPRI** 的讨论。
  在 **`select()`** 返回后，`exceptfds` 将清除所有文件描述符，除了那些发生了异常条件的描述符。
- `nfds`
  这个参数应设置为三个集合中编号最高的文件描述符加1。将检查每个集合中指定的文件描述符，直至这个限制（但请参见 BUGS 部分）。
- `timeout`
  超时参数是一个 `timeval` 结构体（如下所示），它指定了 **`select()`** 应该阻塞等待文件描述符准备好的时间间隔。该调用将一直阻塞，直到以下两种情况之一发生：

  - 一个文件描述符准备就绪
  - 调用被信号处理程序中断；或是 `timeout` 到期

  请注意，超时时间间隔将被向上取整到系统时钟的粒度，而内核调度延迟意味着阻塞间隔可能会稍微超出。
  如果 `timeval` 结构的两个字段都是零，那么 **`select()`** 会立即返回。（这对于轮询很有用。）
  如果将超时设置为 NULL， **`select()`** 将无限期地阻塞，等待文件描述符准备就绪。


### pselect

**`pselect()`** 系统调用允许应用程序安全地等待，直到文件描述符就绪或者捕获到一个信号。

除了以下三个差异之外， **`select()`** 和 **`pselect()`** 的操作是相同的：

- **`select()`** 使用一个以秒和微秒为单位的 `struct timeval` 结构作为超时，而 **`pselect()`** 使用一个以秒和纳秒为单位的 `struct timespec` 结构。
- **`select()`** 可能会更新超时参数以指示剩余的时间。 **`pselect()`** 不会改变这个参数。
- **`select()`** 没有 `sigmask` 参数，其行为与使用 `NULL sigmask` 调用 **`pselect()`** 相同。

`sigmask` 是指向信号掩码的指针（参见 [`sigprocmask(2)`][sigprocmask-2] ）；如果它不为 `NULL`，则 **`pselect()`** 首先用 `sigmask` 指向的信号掩码替换当前的信号掩码，然后执行 `select` 函数，最后恢复原始的信号掩码。（如果 `sigmask` 为 `NULL` ，在 **`pselect()`** 调用期间不会修改信号掩码。）

除了超时参数的精度差异之外，以下 **`pselect()`** 调用的方式：

```c
ready = pselect(nfds, &readfds, &writefds, &exceptfds, timeout, &sigmask);
```

相当于 *原子性* 的执行以下调用：

```c
sigset_t origmask;

pthread_sigmask(SIG_SETMASK, &sigmask, &origmask);
ready = select(nfds, &readfds, &writefds, &exceptfds, timeout);
pthread_sigmask(SIG_SETMASK, &origmask, NULL);
```

**`pselect()`** 被需要的原因是，如果某人想要等待一个信号或者一个文件描述符变为就绪，那么需要一个原子性测试来防止竞态条件。（假设信号处理程序设置了一个全局标志并返回。那么在信号到达测试之后但在调用 **`select()`** 之前，对这个全局标志的测试可能导致无限期挂起。相比之下， **`pselect()`** 允许首先阻塞信号，处理已经到达的信号，然后使用所需的 *信号掩码* 调用 **`pselect()`** ，从而避免了竞态条件。）


### The timeout

**`select()`** 函数的超时参数是一个以下类型的结构体：

```c
struct timeval {
    time_t      tv_sec;   /* seconds */
    suseconds_t tv_usec;  /* microseconds */ 
};
```

**`pselect()`** 对应的参数是一个 **`timespec(3)`** 结构体。

在 Linux 上， **`select()`** 会修改 `timeout` 以反映未休眠的时间量；而大多数其他实现不会这样做。（POSIX.1 允许这两种行为。）这既会导致将读取 `timeout` 的 Linux 代码移植到其他操作系统时出现问题，也会导致将代码移植到 Linux 时出现问题，尤其是当在循环中使用相同的 `struct timeval` 结构体多次调用 **`select()`** 而没有重新初始化它时。认为在 **`select()`** 返回后 `timeout` 是未定义的。


## Return Value

成功时， **`select()`** 和 **`pselect()`** 函数会返回三个返回描述符集合中包含的文件描述符数量（即，`readfds` 、`writefds` 、`exceptfds` 中设置的总位数）。如果超时在任何一个文件描述符准备就绪之前发生，返回值可能为零。

出错时，返回 `-1` ，并设置 [`errno`][errno-3] 来指明错误；文件描述符集合保持不变，而 `timeout` 变量变得未定义。


## Errors

- **EBADF**  在某个集合中给出了一个无效的文件描述符。（可能是已经关闭的文件描述符，或者是在其上发生错误的文件描述符。）但是，请参见 BUGS 。
- **EINTR**  捕获了一个信号；请参阅 [`signal(7)`][signal-7] 。
- **EINVAL** `nfds` 为负数或超过了 **`RLIMIT_NOFILE`** 资源限制（请参见 [`getrlimit(2)`][getrlimit-2] ）。
- **EINVAL** `timeout` 中包含的值无效。
- **ENOMEM** 无法为内部表分配内存。


## Versions

在某些其他的 UNIX 系统中，如果系统无法分配内核级的内部资源， **`select()`** 可能会返回 **EAGAIN** 错误，而不是像 Linux 那样返回 **ENOMEM** 错误。POSIX 为 [`poll(2)`][poll-2] 函数定义了此错误，但并未对 **`select()`** 这样做。可移植的程序可能需要检查 **EAGAIN** 错误，并进行循环处理，这与处理 **EINTR** 的方式类似。


## Standards

POSIX.1-2008


## Notes

以下头文件也提供了 `fd_set` 类型： `<sys/time.h>` 。

`fd_set` 是一个固定大小的缓冲区。当以负数或等于或大于 **`FD_SETSIZE`** 的 `fd` 值执行 **`FD_CLR()`** 或 **`FD_SET()`** 时，会导致未定义的行为。此外，POSIX 要求 `fd` 必须是一个有效的文件描述符。

**`select()`** 和 **`pselect()`** 的操作不受 **`O_NONBLOCK`** 标志的影响。


### The self-pipe trick

在缺少 **`pselect()`** 的系统上，可以使用自管道（self-pipe）技巧来实现可靠（且更可移植）的信号捕获。在这种技术中，信号处理程序向一个管道写入一个字节，而该管道的另一端由主程序中的 **`select()`** 来监控。（为了避免在写入可能已满的管道或从可能为空的管道读取时可能发生的阻塞，在从管道读取和写入时使用非阻塞 I/O。）


### Emulating usleep(3)

在 [`usleep(3)`][usleep-3] 出现之前，某些代码使用 **`select()`** 调用，其中所有三个集合为空，`nfds` 为零，且超时参数不为 NULL，作为一种相当可移植的方式来以亚秒级精度休眠。


### Correspondence between select() and poll() notifications

在Linux内核源码中，我们可以找到以下定义，它们展示了 **`select()`** 的可读、可写和异常条件通知与 [`poll(2)`][poll-2] 和 [`epoll(7)`][epoll-7] 提供的事件通知之间的对应关系：

```c
#define POLLIN_SET (EPOLLRDNORM | EPOLLRDBAND | EPOLLIN | EPOLLHUP | EPOLLERR)
    /* Ready for reading */
#define POLLOUT_SET (EPOLLWRBAND | EPOLLWRNORM | EPOLLOUT | EPOLLERR)
    /* Ready for writing */
#define POLLEX_SET (EPOLLPRI)
    /* Exceptional condition */
```


### Multithreaded applications

如果在一个线程中关闭了由 **`select()`** 监控的文件描述符，结果是不确定的。在一些 UNIX 系统中， **`select()`** 会解除阻塞并返回，表明文件描述符已准备好（随后的 I/O 操作很可能会失败并出现错误，除非另一个进程在 **`select()`** 返回和执行 I/O 操作之间的时间内重新打开文件描述符）。在 Linux（以及一些其他系统）上，在另一个线程中关闭文件描述符对 **`select()`** 没有效果。总之，任何依赖在这种情况下特定行为的应用程序都被认为是错误的。


### C library/kernel differences

Linux 内核允许任意大小的文件描述符集合，通过 `nfds` 的值来确定要检查的集合长度。然而，在 glibc 实现中，`fd_set` 类型的大小是固定的。另见 BUGS。

本页描述的 **`pselect()`** 接口由 glibc 实现。底层 Linux 系统调用名为 **`pselect6()`** 。这个系统调用的行为与 glibc 包装函数有些不同。

Linux **`pselect6()`** 系统调用会修改其超时参数。然而，glibc 包装函数通过使用一个局部变量作为传递给系统调用的超时参数，隐藏了这种行为。因此，glibc 的 **`pselect()`** 函数不会修改其超时参数；这是 POSIX.1-2001 要求的行为。

**`pselect6()`** 系统调用的最后一个参数不是一个 `sigset_t *` 指针，而是一个以下形式的结构体：

```c
struct {
    const kernel_sigset_t *ss;  /* Pointer to signal set */
    size_t ss_len;              /* Size (In bytes) of object
                                   pointed to by 'ss' */
}
```

这允许系统调用同时获取信号集的指针和其大小，同时考虑到大多数架构支持系统调用最多6个参数的事实。有关内核和 libc 中对信号集概念差别的讨论，请参见 [`sigprocmask(2)`][sigprocmask-2] 。


### Historical glibc details

glibc 2.0 提供了一个不正确的 **`pselect()`** 版本，该版本没有带 `sigmask` 参数。

从 glibc 2.1 到 glibc 2.2.1，我们必须定义 **`_GNU_SOURCE`** 以便从 `<sys/select.h>` 中获取 **`pselect()`** 的声明。


## Bugs

POSIX 允许实现定义一个上限，通过常数 **`FD_SETSIZE`** 来指定可以在文件描述符集中指定的文件描述符的范围。Linux 内核没有设定固定的限制，但是 glibc 实现将`fd_set` 设置为固定大小的类型， **`FD_SETSIZE`** 定义为 1024，并且 **`FD_*()`** 宏根据这个限制操作。要监控大于 1023 的文件描述符，请使用 [`poll(2)`][poll-2] 或 [`epoll(7)`][epoll-7] 。

将 `fd_set` 参数实现为 值-结果 参数是一个设计错误，这在 [`poll(2)`][poll-2] 和 [`epoll(7)`][epoll-7] 中得到了避免。

根据 POSIX， **`select()`** 应该检查三个文件描述符集中的所有指定的文件描述符，直到 `nfds-1` 的限制。然而，当前的实现忽略这些集中任何大于当前进程打开的最大文件描述符编号的文件描述符。根据POSIX，任何在这些集合中指定的此类文件描述符都应该导致 **EBADF** 错误。

从 glibc 2. 1开始，glibc 提供了一个使用 [`sigprocmask(2)`][sigprocmask-2] 和 **`select()`** 实现的 **`pselect()`** 的仿真，这个实现仍然容易受到 **`pselect()`** 旨在防止的竞态条件。现代版本的 glibc 在提供 **`pselect()`** 系统调用的内核上使用（无竞态条件的） **`pselect()`** 。

在 Linux 上， **`select()`** 可能会将一个套接字文件描述符报告为 “准备好读取” ，但随后的读取操作却会阻塞。例如，这可能发生在数据到达但经过检查后发现校验和错误并被丢弃的情况下。可能还有其他情况，文件描述符会被错误地报告为就绪。因此，对于不应该阻塞的套接字，使用 **O_NONBLOCK** 可能会更安全。

在 Linux 上，如果 **`select()`** 被信号处理程序中断（即返回 **EINTR** 错误），它还会修改超时时间。这是 POSIX.1 所不允许的。Linux 的 **`pselect()`** 系统调用具有相同的行为，但 glibc 包装器通过将超时时间内部复制到局部变量，并将该变量传递给系统调用来隐藏这种行为。


## Examples

```c
#include <stdio.h>
#include <stdlib.h>
#include <sys/select.h>

int
main(void)
{
    int             retval;
    fd_set          rfds;
    struct timeval  tv;

    /* Watch stdin (fd 0) to see when it has input. */

    FD_ZERO(&rfds);
    FD_SET(0, &rfds);

    /* Wait up to five seconds. */

    tv.tv_sec = 5;
    tv.tv_usec = 0;

    retval = select(1, &rfds, NULL, NULL, &tv);
    /* Don't rely on the value of tv now! */

    if (retval == -1)
        perror("select()");
    else if (retval)
        printf("Data is available now.\n");
        /* FD_ISSET(0, &rfds) will be true. */
    else
        printf("No data within five seconds.\n");

    exit(EXIT_SUCCESS);
}
```


## See Also

包含讨论和教程的示例，参见 [`select_tut(2)`][select-tut-2]












[select-2]: https://man7.org/linux/man-pages/man2/select.2.html
[feature-test-macro-7]: https://man7.org/linux/man-pages/man7/feature_test_macros.7.html
[poll-2]: https://man7.org/linux/man-pages/man2/poll.2.html
[epoll-7]: https://man7.org/linux/man-pages/man7/epoll.7.html
[read-2]: https://man7.org/linux/man-pages/man2/read.2.html
[write-2]: https://man7.org/linux/man-pages/man2/write.2.html
[sigprocmask-2]: https://man7.org/linux/man-pages/man2/sigprocmask.2.html
[errno-3]: https://man7.org/linux/man-pages/man3/errno.3.html
[signal-7]: https://man7.org/linux/man-pages/man7/signal.7.html
[getrlimit-2]: https://man7.org/linux/man-pages/man2/getrlimit.2.html
[usleep-3]: https://man7.org/linux/man-pages/man3/usleep.3.html
[sigprocmask-2]: https://man7.org/linux/man-pages/man2/sigprocmask.2.html
[select-tut-2]: https://man7.org/linux/man-pages/man2/select_tut.2.html

