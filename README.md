## WOL

> 这是一个使用 go 编写的简单的 WOL(Wake on LAN) 工具，主要用于远程开机(需要主板支持)，工具参考了 [sabhiram/go-wol](https://github.com/sabhiram/go-wol) 的部分代码

### 使用

```sh
➜ ~ wol --help
NAME:
   wol - Wake-on-LAN TOOL

USAGE:
   wol [global options] command [command options] [arguments...]

VERSION:
   v1.0.0 2020-12-07 15:47:57 858f64b34d7cfcb041e07c482b121a7c0d761d0a

AUTHOR:
   mritd <mritd@linux.com>

COMMANDS:
   add      add machine
   del      del machine
   wake     wake machine
   print    print machines
   example  print example config
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value  wol config (default: "/Users/bleem/.wol.yaml") [$WOL_CONFIG]
   --name value, -n value    machine name
   --mac value, -m value     machine mac address
   --help, -h                show help (default: false)
   --version, -v             print the version (default: false)

COPYRIGHT:
   Copyright (c) 2020 mritd, All rights reserved.
```

### 配置

可以通过 `wol example > ~/.wol.yaml` 生成样例配置文件，也可以通过 `--config` 选项指定配置文件位置；默认情况下 wol 只读取 `$HOME/.wol.yaml` 配置文件。

### Bash Completion

Bash Completion 以及 ZSH 支持请参考 [urfave/cli](https://github.com/urfave/cli/blob/master/docs/v2/manual.md#bash-completion) 文档了解如何使用。