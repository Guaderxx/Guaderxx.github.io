---
title: runc-usage
date: 2024-04-23 17:36:31
tags:
- runc
categories:
- doc
- runc
keywords: runc
copyright: Guader
copyright_author_href:
copyright_info:
---

[runc](https://github.com/opencontainers/runc) 算是 [runtime caller](https://github.com/Guaderxx/runtime-spec_ZH/blob/chinese-translation/glossary_zh.md#runtime-caller) 的 [参考实现](https://github.com/Guaderxx/runtime-spec_ZH/blob/chinese-translation/implementations_zh.md#runtime-container) ，不过用的确实很多。

当然了，这是个底层工具，设计时也没考虑到终端用户，所以以下内容仅作为参考（乐子）看。

---

1. 下载一个喜欢的镜像，比如 busybox

> 用 docker 是因为方便，其他的也行

```bash
docker pull busybox
```

2. 创建 bundle

```bash
mkdir ~/mybusybox
cd ~/mybusybox

# create the rootfs directory
mkdir rootfs

# export busybox via Docker into the rootfs directory
docker export $(docker create busybox) | tar -C rootfs -xvf -

# create `config.json` by runc
runc spec
```

3. 运行该容器 run

`run` 命令会处理容器的创建，启动和删除。

```bash
# At ~/mybusybox
runc run [container-id] # like 123456789, whatever
```

4. 运行该容器 manual

先修改一下 `config.json`

```diff
...
"process": {
-    "terminal": true
+    "terminal": false
    "user": {
        "uid": 0,
        "gid": 0
    },
    "args": [
-        "sh"
+        "sleep", "5"
    ],
...
}
```

现在可以完成生命周期操作了

4.1. 创建容器

```bash
# ~/mybusybox
sudo runc create [container-id]
```

查看容器状态：

```bash
sudo runc list
```

```
ID          PID         STATUS      BUNDLE                                            CREATED                          OWNER
12345699    2201396     created     /home/user/github.com/Guaderxx/runc/mybusybox   2024-04-23T09:33:42.994710154Z   root
```

4.2. 运行容器

运行容器，容器会在 5S 内结束，sleep 期间可以查看容器状态

```bash
sudo runc start [container-id] # 比如我这里是 12345699
```

运行时容器状态：

```
ID          PID         STATUS      BUNDLE                                            CREATED                          OWNER
12345699    2201396     running     /home/user/github.com/Guaderxx/runc/mybusybox   2024-04-23T09:33:42.994710154Z   root
```

结束后容器状态：

```
ID          PID         STATUS      BUNDLE                                            CREATED                          OWNER
12345699    0           stopped     /home/user/github.com/Guaderxx/runc/mybusybox   2024-04-23T09:33:42.994710154Z   root
```

4.3. 删除容器

```bash
sudo runc delete [container-id] 
```

此时再查看容器状态，可以发现已经看不到了。
