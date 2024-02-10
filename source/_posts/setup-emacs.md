---
title: setup emacs
date: 2024-02-10 14:58:27
tags:
- emacs
categories:
- doc
keywords:
- emacs
copyright: Guader
copyright_author_href:
copyright_info:
---


> 前一阵子 Emacs 抽风了打不开，我用着 vscode 也凑合就没管。
> 今天正得闲，重新配置一遍。


1. 首先，下载

进入 [GNU Emacs download][gnu-emacs-download] 页。
根据自己的系统选择合适的下载方式即可。

比如我是用的 ubuntu，就是：

```bash
sudo apt-get install emacs
```


2. 配置

下载完正常就可以打开了，不过这时候都是默认配置。
方便起见，我 fork 了 [purcell][purcell-profile] 的 [emacs.d][purcell-emacs-d]。
然后进入根目录下载即可：

```bash
cd 

git clone https://github.com/Guaderxx/emacs.d.git .emacs.d
```

然后打开 Emacs GUI，运行 `M-x package-refresh-contents`

*这个配置是没有包含 golang 的配置的*

下载需要的包通过 [melpa][melpa]，也就是在 Emacs 中运行： `M-x package-install [package-name]`

所以重启一下再下载： 

- [`go-mode`][go-mode]
  - gopls
    - `go install golang.org/x/tools/gopls@latest`
- [`flycheck-golangci-lint`][flycheck-golangci-lint]
- [`go-projectile`][go-projectile]
- [`treemacs`][treemacs]
- `treemacs-projectile`
  - 运行 `M-x treemacs-projectile` 将项目加入 treemacs


然后根据需要再进一步配置即可。



[gnu-emacs-download]: https://www.gnu.org/software/emacs/download.html
[purcell-profile]: https://github.com/purcell
[purcell-emacs-d]: https://github.com/purcell/emacs.d
[melpa]: https://melpa.org/#/
[go-mode]: https://github.com/dominikh/go-mode.el
[flycheck-golangci-lint]: https://github.com/weijiangan/flycheck-golangci-lint
[go-projectile]: https://github.com/dougm/go-projectile
[treemacs]: https://github.com/Alexander-Miller/treemacs
