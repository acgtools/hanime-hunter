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

[toc]

## Installation

### Using `go`

```sh
$ go install -ldflags "-s -w" github.com/acgtools/hanime-hunter
```

### Download from releases

[release page](https://github.com/acgtools/hanime-hunter/releases)

## Quick Start

### Prerequisites

#### Ensure that your terminal charset is UTF-8

**Windows**

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

**Linux**

```sh
$ echo $LANG
en_US.UTF-8
```

### Command Help

```sh
$ hani.exe -h
HAnime downloader. Repo: https://github.com/acgtools/hanime-hunter

Usage:
  hani [command]

Available Commands:
  dl          download
  help        Help about any command
  version     Print version info

Flags:
  -h, --help               help for hani
      --log-level string   log level, options: debug, info, warn, error, fatal (default "info")

Use "hani [command] --help" for more information about a command.
```

#### Download

```sh
$ hani help dl
download

Usage:
  hani dl [flags]

Flags:
  -h, --help                help for dl
  -i, --info                get anime info only
      --low-quality         download the lowest quality video
  -o, --output-dir string   output directory
  -q, --quality string      specify video quality. e.g. 1080p, 720p, 480p ...
  -s, --series              download full series

Global Flags:
      --log-level string   log level, options: debug, info, warn, error, fatal (default "info")
```

## Download

### Only one episode

```sh
# Download from the watch page
# The anime will be saved in ./anime_series_title/
$ hani dl https://hanime1.me/watch?v=xxxx
```

### Full series based on the specified episode

```sh
# Download the full series
# E.g. If you give the link of the Anime_Foo_02
# then the full series of Anime_Foo will be downloaded (Anime_Foo_01, Anime_Foo_02, ...)
$ hani dl -s https://hanime1.me/watch?v=xxxx
```

### Specify the output directory

```sh
# The anime will be saved in output_dir/anime_series_title/
$ hani dl -o <output_dir>
```

### Specify the quality

```sh
# You can specify the quality of video
# if it is not exists, the default (highest quality) will be downloaded
$ hani dl -q "720p" https://hanime1.me/watch?v=xxxx
```

### Get info only

```sh
# Get only the downloadable video info
# title, quality, file extension
$ hani dl -i https://hanime1.me/watch?v=xxxx
```

## Supported Site

> **NSFW** Warning, the following site may contain sensitive content.

| Site       | Language | Episode    | Series   | Playlist   | Status        |
| ---------- | -------- | ---------- | -------- | ---------- | ------------- |
| hanime1.me | Chinese  | ✓          | ✓        | Developing | Available     |
| hanime.tv  | English  | Developing | Planning | Planning   | Not Available |

## Issue

Feel free to create issues to report bugs or request new features.
