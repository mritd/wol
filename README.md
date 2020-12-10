## WOL

> 这是一个使用 go 编写的简单的 WOL(Wake on LAN) 工具，主要用于远程开机(需要主板支持)，工具参考了 [sabhiram/go-wol](https://github.com/sabhiram/go-wol) 的部分代码

### 安装

wol 使用 Go 编写，运行时只需要一个单独的二进制文件即可执行；用户可以从 [RELEASE](https://github.com/mritd/wol/releases) 页下载预编译的二进制文件到本地执行；以下为 Linux 安装示例:

``` sh
# 后续新版本发布请自行修改版本号
export WOL_VERSION='v1.0.1'

# 此命令将下载 Linux 系统 x64 架构的预编译文件，其他架构、平台请自行替换
wget https://github.com/mritd/wol/releases/download/${WOL_VERSION}/wol_linux_amd64

# 增加可执行权限
chmod +x wol_linux_amd64

# 移动到 PATH 目录
mv wol_linux_amd64 /usr/local/bin/wol

# 运行测试
wol --help
```

### 使用

> wol 工具默认读取家目录(`$HOME`)下的 `.wol.yaml` 配置文件，并从该配置文件中获取当前存在的主机列表以及得知如何发送数据包；
> **默认情况下该配置文件不存在，且 wol 工具也不会自动创建，所以直接运行 `print、wake` 等命令可能会出现相关错误提示，属于正常情况。**

在使用之前请执行 `wol example > ~/.wol.yaml` 创建样例配置文件，以下为 `wol` 命令列表:

- `wol example`: 向控制台输出 WOL 的样例配置文件
- `wol wake 主机名/Mac地址`: 唤醒某台主机
- `wol add --name 主机名 --mac Mac地址`: 向配置文件中增加一台主机
- `wol del 主机名/Mac地址`: 从配置文件中删除一台主机
- `wol print`: 打印配置文件中所有主机列表

使用样例:

```sh
# 唤醒名字为 iMac 的主机(真正的 mac 电脑似乎不支持远程开机)
wol wake iMac

# 使用 Mac 地址进行唤醒
wol wake E0:D5:5E:6E:30:C9

# 添加一台叫做 nas 的主机
wol add --name nas --mac E1:D2:3E:6E:20:C5

# 删除一台叫做 nas 的主机
wol del nas

# 通过 Mac 地址来删除
wol del E1:D2:3E:6E:20:C5

# 显示当前配置文件中的所有主机
wol print
```

**其他更详细的使用请使用 `--help` 选项查看帮助文档(支持子命令的 `--help` 选项):**

```sh
➜ wol --help
NAME:
   wol - Wake-on-LAN TOOL

USAGE:
   wol [global options] command [command options] [arguments...]

VERSION:
   v1.0.0 2020-12-07 16:26:16 b36316e772b9ca3abecb6b34fd05797ccbd98044

AUTHOR:
   mritd <mritd@linux.com>

COMMANDS:
   add      add device
   del      del device
   wake     wake device
   print    print devices
   example  print example config
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value  wol config (default: "$HOME/.wol.yaml") [$WOL_CONFIG]
   --name value, -n value    device name
   --mac value, -m value     device mac address
   --help, -h                show help (default: false)
   --version, -v             print the version (default: false)

COPYRIGHT:
   Copyright (c) 2020 mritd, All rights reserved.
```

### 配置

默认配置文件格式可以通过 `wol example` 查看:

```sh
➜  ~ wol example

devices:
- name: iMac
  mac: e0:d5:5e:6e:30:c9
```

主机配置默认存储在 `$HOME/.wol.yaml` 中，这是一个标准的 `yaml` 配置文件，除了 `name`、`mac` 字段外完整的配置文件支持列表如下:

```yaml
# 被隐藏的配置属于高级配置，除非你明确了解其含义和作用，否则不建议修改
# 错误修改只会导致发送的唤醒命令失败
devices:
- name: iMac
  mac: e0:d5:5e:6e:30:c9
  # 数据包发送端口
  port: 7
  # 数据包发送接口
  broadcast_interface: eth0
  # 数据包发送 IP
  broadcast_ip: 192.168.2.1
```

### 自动提示

针对于习惯命令行操作的高级用户，大部分人可能更习惯使用 `tab` 键来实现自动完成；`wol` 命令提供了自动完成支持包括 zsh，以下为自动完成的配置说明:

#### bash 用户

- 下载 `autocomplete/bash_autocomplete` 文件到任意位置
- 确保已经安装好了 `wol` 工具(在 `PATH` 中可以找到)
- 在 `~/.bashrc` 中添加 `PROG=wol source path/to/cli/autocomplete/bash_autocomplete`(路径请自行替换)
- 退出 bash 重新登录，或执行 `source ~/.bashrc` 命令
- 最后执行 `wol` + `tab` 进行测试是否有智能提示

#### zsh 用户

注意: 以下仅在 ohmyzsh 测试成功，标准 zsh 理论上也兼容

- 下载 `autocomplete/zsh_autocomplete` 文件到任意位置
- 确保已经安装好了 `wol` 工具(在 `PATH` 中可以找到)
- 在 `~/.zshrc` 中添加 `PROG=wol _CLI_ZSH_AUTOCOMPLETE_HACK=1 source path/to/cli/autocomplete/zsh_autocomplete`(路径请自行替换)
- 退出 bash 重新登录，或执行 `source ~/.zshrc` 命令
- 最后执行 `wol` + `tab` 进行测试是否有智能提示

#### powershell 用户

...powershell 用户请自行参考 [urfave/cli 文档](https://github.com/urfave/cli/blob/master/docs/v2/manual.md#powershell-support)
