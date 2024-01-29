---
title: sd_notify
date: 2024-01-29 16:27:20
tags:
categories:
- man-pages
keywords:
- sd_notify
copyright: Guader
copyright_author_href:
copyright_info:
---

[new-source][sd-notify]


## 名称

`sd_notify`, `sd_notifyf`, `sd_pid_notify`, `sd_pid_notifyf`, `sd_pid_notify_with_fds`, `sd_pid_notifyf_with_fds`, `sd_notify_barrier`, `sd_pid_notify_barrier` -- 通知服务管理器关于启动完成以及其他服务状态变化。


## 概述

```c
int sd_notify(int unset_environment, const char *state);

int sd_notifyf( int unset_environment, 
    const char *format, 
    …);

int sd_pid_notify( pid_t pid, 
    int unset_environment, 
    const char *state);

int sd_pid_notifyf( pid_t pid, 
    int unset_environment, 
    const char *format, 
    …);

int sd_pid_notify_with_fds(pid_t pid,
  int unset_environment,
  const char *state,
  const int *fds,
  unsigned n_fds);

int sd_pid_notifyf_with_fds( pid_t pid,
  int unset_environment,
  const int *fds,
  size_t n_fds,
  const char *format,
  …);

int sd_notify_barrier(int unset_environment,
  uint64_t timeout);

int sd_pid_notify_barrier( pid_t pid,
  int unset_environment,
  uint64_t timeout);
```


## 描述

服务可以调用 `sd_notify()` 来通知服务管理器有关状态更改的消息。
它可用于发送以类似环境块的字符串编码的任意信息。
最重要的是，它可用于启动或重新加载完成通知。

如果参数 `unset_environment` 非零，则 `sd_notify()` 将在返回之前取消设置 `$NOTIFY_SOCKET` 环境变量（无论函数调用本身是否成功）。
进一步调用 `sd_notify()` 将失败，且该变量不再由子进程继承。

`state` 参数应包含以换行符分隔的变量分配列表，其风格与环境块类似。
如果未指定，则隐含尾随换行符。
该字符串可以包含任何类型的变量赋值，但请参阅下一节以获取服务管理器理解的赋值列表。

请注意，只有在服务定义文件中正确设置了 `NofityAccess=` 选项时，`systemd` 才会接收从服务发送的状态数据。
查看 [systemd.service(5)][systemd-service-5] 了解更多。

请注意，只有当发送进程在 PID 1 处理消息时仍然存在，或者发送进程被服务管理器明确运行时追踪时，`sd_notify()` 通知才可以正确的归属于单元。
如果服务管理器最初分叉了该进程，则属于后一种情况，即，在与 `NotifyAccess=main` 或 `NotifyAccess=exec` 匹配的所有进程上。
相反，如果单元的辅助进程发送了 `sd_notify()` 消息并立即退出，服务管理器可能无法正确的将消息归属于单元，因此即时为其设置了 `NotifyAccess=all` ，也会忽略该消息。

因此，为了消除涉及客户端单元查找和正确将通知归属于单元的所有竞争条件，可以使用 `sd_notify_barrier()` 。
此调用充当同步点，并确保服务管理器在成功返回时已拾取此调用之前发送的所有通知。
对于不是由服务管理器调用的客户端，需要使用 `sd_notify_barrier()`，否则就没有必要使用这种同步机制将通知归属到单元。

`sd_notifyf()` 类似于 `sd_notify()`，但接收类似 `printf()` 的格式字符串和参数。

`sd_pid_notify()` 和 `sd_pid_notifyf()` 类似于 `sd_notify()` 和 `sd_notifyf()`，但以进程 ID（PID）作为第一个参数，将其用作消息的源 PID。
这对于代表其他进程发送通知消息非常有用，前提是要有相应的权限。
如果 PID 参数指定为 0，则使用调用进程的进程 ID，在这种情况下，调用与 `sd_notify()` 和 `sd_notifyf()` 完全等价。

`sd_pid_notify_with_fds()` 类似于 `sd_pid_notify()`，但多了一个文件描述符数组。
这些文件描述符会随通知信息一起发送给服务管理器。
如上所述，这对发送 `"FDSTORE=1"` 信息特别有用。
附加参数是指向文件描述符数组的指针，以及数组中文件描述符的个数。
如果传入的文件描述符数量为 0，则调用与 sd_pid_notify() 完全等价，即，没有文件描述符被传递。
请注意，通过不含 `"FDSTORE=1"` 的报文发送到服务管理器的文件描述符会在收到后立即关闭。

`sd_pid_notifyf_with_fds()` 是 `sd_pid_notify_with_fds()` 和 `sd_notifyf()` 的组合，即，它接受 PID 和一组文件描述符作为输入，并处理格式字符串以生成状态字符串。

`sd_notify_barrier()` 允许调用者同步接收先前发送的通知信息，并使用 `BARRIER=1` 命令。
它获取一个以微秒为单位的相对超时值，并将其传递给 [ppoll(2)][ppoll-2]。
`UINT64_MAX` 值被解释为无限超时。

`sd_pid_notify_barrier()` 与 `sd_notify_barrier()`类似，但允许为通知信息指定源 PID。


## 知名赋值

以下赋值具有定义的含义：

- `READY=1` 
    告诉服务管理器服务启动已完成，或者服务已完成重新加载其配置。只有当服务定义文件设置了 `Type=notify` 或 `Type=notify-reload` 时，`systemd` 才会使用该选项。由于发出 `non-readiness` 信号的价值不大，服务应发送的唯一值就是 `READY=1`（即, 未定义 `READY=0` ）。
- `RELOADING=1`
    告诉服务管理器该服务正在开始重新加载其配置。这有助于服务管理器跟踪服务的内部状态，并将其呈现给用户。请注意，发送此通知的服务在完成重新加载配置后，还必须发送 `READY=1` 通知。服务管理器通过该机制收到的重载通知的传播方式与最初通过服务管理器启动时的传播方式相同。此消息与 `Type=notify-reload` 服务特别相关，用于通知服务管理器已收到重新加载服务的请求并正在处理。

    在版本 217 添加。
- `STOPPING=1`
    通知服务管理器服务开始关闭。这对于允许服务管理器跟踪服务的内部状态并将其呈现给用户非常有用。

    在版本 217 添加。
- `MONOTONIC_USEC=...`
    一个字段，包含客户端生成通知信息时的单调时间戳（根据 `CLOCK_MONOTONIC`），格式为十进制，单位为 μs。通常与 `RELOADING=1` 结合使用，以便服务管理器正确同步重载周期。查看 [systemd.service(5)][systemd-service-5]了解更多细节，尤其是 `Type=notify-reload` 。

    在版本 253 添加。
- `STATUS=...`
    将描述服务状态的单行 UTF-8 状态字符串传回服务管理器。这是自由格式，可用于多种目的：一般状态反馈、类似 `fsck` 的程序可传递完成百分比，失败的程序可传递人类可读的错误信息。
    例如："STATUS=Completed 66% of file system check..."

    在版本 233 添加。
- `NOTIFYACCESS=...`
    在运行时重置服务状态通知socket的访问权限，覆盖服务单元文件中的 `NotifyAccess=` 设置。
    查看 [systemd.service(5)][systemd-service-5]了解更多细节，尤其是 `NotifyAccess=`，以获得可接受值的列表 。

    在版本 254 添加。
- `ERRNO=...`
    如果服务失败，则以字符串形式显示 errno 样式的错误代码。
    例如: "ERRNO=2"，表示 ENOENT。

    在版本 233 添加。
- `BUSERROR=...`
    如果服务失败，则会出现 `D-Bus` 错误类型的错误代码。
    例如： "BUSERROR=org.freedesktop.DBus.Error.TimedOut"。
    请注意，`systemd` 目前不使用此分配。

    在版本 233 添加。
- `EXIT_STATUS=...`
    服务或管理器本身的退出状态。
    请注意，`systemd` 目前不会在服务发送时使用该值，因此该赋值仅供参考。
    管理器会将此通知发送到它的通知套接字，该套接字可用于收集系统（容器或虚拟机）关闭时的退出状态。
    例如，[mkosi(1)][mkosi-1] 就使用了这一功能。
    返回值可以通过 [systemctl(1)][systemctl-1] **exit** verb 设置。

    在版本 254 添加。
- `MAINPID=...`
    服务的主进程 ID (PID)，以防服务管理器没有自行分叉进程。
    例如: "MAINPID=4711"。

    在版本 233 添加。
- `WATCHDOG=1`
    告诉服务管理器更新看门狗时间戳。
    如果启用了` WatchdogSec=`，则服务需要定期发出 keep-alive ping。
    有关如何启用此功能的信息，请参见 [systemd.service(5)][systemd-service-5]；
    有关服务如何检查看门狗是否启用的详细信息，请参见 [sd_watchdog_enabled(3)][sd-watchdog-enabled-3]。
- `WATCHDOG=trigger`
    告诉服务管理器，服务检测到内部错误，应由配置的看门狗选项处理。
    这将触发与启用 `WatchdogSec=` 后服务未及时发送 "WATCHDOG=1 "相同的行为。
    请注意，无需启用 `WatchdogSec=`，`"WATCHDOG=trigger"` 即可触发看门狗操作。
    有关看门狗行为的信息，请参阅 [systemd.service(5)][systemd-service-5]。

    在版本 243 添加。
- `WATCHDOG_USEC=…`
    在运行期间重置 `watchdog_usec` 值。
    请注意，使用 `sd_event_set_watchdog()` 或 `sd_watchdog_enabled()` 时，此功能不可用。
    示例："WATCHDOG_USEC=20000000"。

    在版本 233 添加。
- `EXTEND_TIMEOUT_USEC=...`
    告诉服务管理器延长与当前状态相对应的启动、运行或关闭服务超时时间。
    指定的值是服务必须发送新信息的时间（以微秒为单位）。
    如果未收到信息，服务就会超时，但前提是当前状态的运行时间超过 `TimeoutStartSec=`、`RuntimeMaxSec=` 和 `TimeoutStopSec=` 的原始最大时间。
    请参阅 [systemd.service(5)][systemd-service-5] 了解服务超时的影响。
    在版本 236 添加。
- `FDSTORE=1`   
    在服务管理器中存储文件描述符。
    以这种方式发送的文件描述符将由服务管理器代为保管，并在下次启动或重启服务时使用通常的文件描述符传递逻辑交还，请参见 [sd_listen_fds(3)][sd-listen-fds-3]。
    重启时不应关闭的任何打开的套接字和其他文件描述符都可以这样存储。
    当服务停止时，其文件描述符存储空间将被丢弃，其中的所有文件描述符也将关闭，除非使用 `FileDescriptorStorePreserve=` 改写，参见 [systemd.service(5)][systemd-service-5]。

    只有当服务的 `FileDescriptorStoreMax=` 设置为非零（默认为 零，请参阅 [systemd.service(5)][systemd-service-5] ）时，服务管理器才会接受该服务的消息。
    服务管理器将为已启用文件描述符存储的服务设置 `$FDSTORE` 环境变量，请参阅 [systemd.exec(5)][systemd-exec-5]。

    如果未设置 `FDPOLL=0`，且文件描述符是可轮询的（参见 [epoll_ctl(2)][epoll-ctl-2] ），那么在这些文件描述符上出现的任何 `EPOLLHUP` 或 `EPOLLERR` 事件都将导致这些文件描述符自动从存储中删除。

    多组文件描述符可以在不同的报文中发送，在这种情况下，这些文件描述符会被合并。
    服务管理器会删除重复的文件描述符（指向相同对象的文件描述符），然后再将其传递给服务。
    
    该功能应用于实现可在明确请求或崩溃后重新启动而不会丢失状态的服务。
    应用程序状态既可以序列化到 `/run/` 中的文件，也可以存储在 [memfd_create(2)][memfd-create-2] 内存文件描述符中。
    使用 `sd_pid_notify_with_fds()` 发送 "FDSTORE=1 "的信息。
    建议将 `FDSTORE=` 与 `FDNAME=` 结合使用，以方便管理存储的文件描述符。

    有关文件描述符存储的更多信息，请参阅 [文件描述符存储][FILE-DESCRIPTOR-STORE] 概述。

    在版本 219 添加。
- `FDSTOREREMOVE=1`
    从文件描述符存储中删除文件描述符。
    该字段需要与 `FDNAME=` 结合使用，以指定要删除的文件描述符的名称。

    在版本 236 添加。
- `FDNAME=...`
    与 `FDSTORE=1` 结合使用时，指定提交的文件描述符的名称。
    与 `FDSTOREREMOVE=1` 结合使用时，指定要删除的文件描述符的名称。
    该名称在激活过程中传递给服务，可使用 [`sd_listen_fds_with_names(3)`][sd-listen-fds-with-names-3] 进行查询。
    如果提交的文件描述符中没有该字段，则会被称为 "stored"。

    名称可由任意 ASCII 字符组成，控制字符或 `:` 除外。
    长度不得超过 255 个字符。
    如果提交的名称不符合这些限制，将被忽略。
    请注意，如果在一条信息中提交了多个文件描述符，则所有文件描述符都将使用指定的名称。
    如果要为已提交的文件描述符指定不同的名称，请在不同的信息中提交。

    在版本 233 添加。
- `FDPOLL=0`
    与 `FDSTORE=1` 结合使用时，将禁用对已存储文件描述符的轮询，无论这些文件描述符是否可轮询。
    由于该选项禁止在 `EPOLLERR` 和 `EPOLLHUP` 上自动清理已存储的文件描述符，因此必须注意确保正确的手动清理。
    一般情况下不建议使用该选项，除非自动清理产生了不必要的行为，如过早地从存储中丢弃文件描述符。

    在版本 246 添加。
- `BARRIER=1`
    告诉服务管理器，客户端通过关闭随此命令发送的文件描述符明确请求同步。
    服务管理器保证，只有在处理完 `BARRIER=1` 命令之前发送的所有通知信息后，才会处理 `BARRIER=1` 命令。
    因此，该命令与单个文件描述符一起使用时，可以同步接收以前的所有状态信息。
    请注意，该命令不能与其他通知混用，必须以单独消息的形式发送给服务管理器，否则所有分配都将被忽略。
    请注意，使用此命令发送 0 个或超过 1 个文件描述符是违反协议的。

    在版本 246 添加。


服务发送的通知信息由服务管理器解释。
未知任务可能会被记录，但在其他情况下会被忽略。
因此，发送不在此列表中的任务是没有用的。
服务管理器还会向其通知套接字发送一些消息，然后由机器或容器管理器使用这些消息。


## 返回值

如果调用失败，这些调用会返回一个负的 errno 样式的错误代码。
如果未设置 `$NOTIFY_SOCKET`，因此无法发送状态信息，则返回 0。
如果状态已发送，这些函数将返回一个正值。
为了同时支持实施和不实施该方案的服务管理器，一般建议忽略该调用的返回值。
请注意，返回值仅表示通知信息是否被正确排队，并不反映信息是否能被成功处理。
具体来说，当使用 `FDSTORE=1` 试图存储文件描述符，但服务实际上未配置为允许存储文件描述符时，不会返回错误（见上文）。


## 说明

此处描述的功能以共享库的形式提供，可通过 libsystemd [pkg-config(1)][pkg-config-1] 文件对其进行编译和链接。

这里描述的代码使用了 [getenv(3)][getenv-3]，它被声明为非多线程安全的。
这意味着调用此处所述函数的代码不得从并行线程中调用 [setenv(3)][setenv-3]。
建议只在程序的早期阶段，即没有其他线程启动时调用 `setenv()`。

这些函数向 `$NOTIFY_SOCKET` 环境变量中引用的套接字发送以状态字符串为有效载荷的单个数据报。
如果 `$NOTIFY_SOCKET` 的第一个字符是 `/` 或 `@` ，则字符串会被理解为 `AF_UNIX` 或 `Linux` 抽象命名空间套接字（分别），在这两种情况下，数据报都会使用 `SCM_CREDENTIALS` 附带发送服务的进程凭据。
如果字符串以 `vsock:` 开头，那么该字符串就会被理解为 `AF_VSOCK` 地址，这对主机上的管理程序/VMM 或其他进程在虚拟机完成启动时接收通知非常有用。
请注意，如果管理程序不支持 `SOCK_DGRAM over AF_VSOCK`，则将使用 `SOCK_SEQPACKET` 代替。
地址格式应为 `vsock:CID:PORT` 。
请注意，与 `vsock` 的其他用途不同，`CID` 是强制性的，不能为 `VMADDR_CID_ANY` 。
请注意，PID1 将从特权端口（即低于 1024 端口）发送 `VSOCK` 数据包，以解决客户机中的无特权进程可能试图向主机发送恶意通知，从而导致主机据此做出破坏性决定的问题。


## 环境

- `$NOTIFY_SOCKET`
    由服务管理器设置，用于受监控进程的状态和启动完成通知。
    该环境变量指定 `sd_notify()` 与之对话的套接字。详见上文。


## 例子

### 例1. Start-up Notification

服务启动完成后，可能会发出以下调用来通知服务管理器：

```c
sd_notify(0, "READY=1");
```


### 例2. Extended Start-up Notification

服务在完成初始化后可发送以下信息：

```c
sd_notify(0, "READY=1\n"
        "STATUS=Processing requests...\n"
        "MAINPID=%lu",
        (unsigned long) getpid());
```


### 例3. Error Cause Notification

如果服务失败，可在退出前不久发送以下信息：

```c
sd_notifyf(0, "STATUS=Failed to start up: %s\n"
            "ERRNO=%i",
            strerror_r(errnum, (char[1024]){}, 1024),
            errnum);
```


### 例4. Store a File Descriptor in the Service Manager

要在服务管理器中存储打开的文件描述符，以便在服务重启后继续运行而不丢失状态，请使用 "FDSTORE=1"：

```c
sd_pid_notify_with_fds(0, 0, "FDSTORE=1\nFDNAME=foobar", &fd, 1);
```


### 例5. Eliminating race conditions

当发送通知的客户端不是由服务管理器生成时，它可能会过快退出，而服务管理器可能无法将这些通知正确地归属于该单元。
为防止此类竞赛，请使用 `sd_notify_barrier()`，同步接收在此调用之前发送的所有通知。

```c
sd_notify(0, "READY=1");
/* set timeout to 5 seconds */
sd_notify_barrier(0, 5 * 1000000);
```



[sd-notify]: https://www.freedesktop.org/software/systemd/man/latest/sd_notify.html
[systemd-service-5]: https://www.freedesktop.org/software/systemd/man/latest/systemd.service.html#
[ppoll-2]: https://man7.org/linux/man-pages/man2/ppoll.2.html
[mkosi-1]: https://manpages.debian.org/unstable/mkosi/mkosi.1.en.html
[systemctl-1]: https://www.freedesktop.org/software/systemd/man/latest/systemctl.html#
[sd-watchdog-enabled-3]: https://www.freedesktop.org/software/systemd/man/latest/sd_watchdog_enabled.html#
[sd-listen-fds-3]: https://www.freedesktop.org/software/systemd/man/latest/sd_listen_fds.html#
[systemd-exec-5]: https://www.freedesktop.org/software/systemd/man/latest/systemd.exec.html#
[epoll-ctl-2]: https://man7.org/linux/man-pages/man2/epoll_ctl.2.html
[memfd-create-2]: https://man7.org/linux/man-pages/man2/memfd_create.2.html
[FILE-DESCRIPTOR-STORE]: https://systemd.io/FILE_DESCRIPTOR_STORE
[sd-listen-fds-with-names-3]: https://www.freedesktop.org/software/systemd/man/latest/sd_listen_fds_with_names.html#
[pkg-config-1]: http://linux.die.net/man/1/pkg-config
[getenv-3]: https://man7.org/linux/man-pages/man3/getenv.3.html
[setenv-3]: https://man7.org/linux/man-pages/man3/setenv.3.html
