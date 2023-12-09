# hanime-hunter

English | [简体中文](./README_ZH_CN.md)

A CLI app to download HAnime.

If this repo is helpful to you, please consider giving it a star (o゜▽゜)o☆ . Thank you OwO.

> Random Wink OvO

<img src="https://waifu-getter.vercel.app/sfw?eps=wink" />

<br />

<!--
  If you want to use your own Moe-Counter
  please refer to the tutorial
  in its original repo: https://github.com/journey-ad/Moe-Counter
  and deploy it to the Replit or Glitch
-->
![](https://political-capable-roll.glitch.me/get/@acg_tools_hanime_hunter?theme=rule34)

## Installation

### Using `go`

```sh
$ go install -ldflags "-s -w" github.com/acgtools/hanime-hunter
```

### Download from releases

[release page](https://github.com/acgtools/hanime-hunter/releases)

## Quick Start

```sh
$ moe-go -h
A TUI app for finding anime scene by image, using trace.moe api

Usage:
  moe-go [command]

Available Commands:
  file        search image by file
  help        Help about any command

Flags:
  -h, --help      help for moe-go
  -v, --version   version for moe-go

Use "moe-go [command] --help" for more information about a command.
```

### Ensure your terminal charset is UTF-8

#### Windows

```cmd
> chcp
Active code page: 65001

# if code page is not 65001(utf-8), change it temporarily
> chcp 65001
```

If you want to set the default charset, follow the steps:

1. Start -> Run -> regedit
2. Go to `[HKEY_LOCAL_MACHINE\Software\Microsoft\Command Processor\Autorun]`
3. Change the value to `@chcp 65001>nul`

If `Autorun` is not present, you can add a `New String`.

This approach will auto-execute `@chcp 65001>nul` when `cmd` starts.

#### Linux

```sh
$ echo $LANG
en_US.UTF-8
```

### Find by image file

```sh
$ moe-go file <image file path>
```

Keys:

- `up`, `down`: move the cursor
- `space` ,`enter`: select one result
- `q`: quit program

#### Example

![gochiusa_rize](https://raw.githubusercontent.com/dreamjz/pics/main/pics/2023/202312042054552.jpg)

![1](https://raw.githubusercontent.com/dreamjz/pics/main/pics/2023/202312042051978.gif)

## Issue

Feel free to create issues to report bugs or request new features.
