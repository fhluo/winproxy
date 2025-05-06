<div align="center">

# winproxy

更改 Windows 系统代理配置
<br><br>
<a href="https://github.com/fhluo/winproxy/blob/main/LICENSE">
<img src="https://img.shields.io/github/license/fhluo/winproxy" alt="license"></a>
<a href="https://github.com/fhluo/winproxy/actions/workflows/build.yaml">
<img src="https://github.com/fhluo/winproxy/actions/workflows/build.yaml/badge.svg" alt="build workflow"></a>
<a href="https://goreportcard.com/report/github.com/fhluo/winproxy">
<img src="https://goreportcard.com/badge/github.com/fhluo/winproxy" alt="go report card"></a>
<a href="https://pkg.go.dev/github.com/fhluo/winproxy/go">
<img src="https://img.shields.io/github/v/tag/fhluo/winproxy?filter=go%2F*&label=pkg"></a>

<samp>

**[English](README.md)** ┃ **[简体中文](README.zh-Hans.md)**

</samp>
</div>

## 介绍

`winproxy` 是一个用于更改 Windows 系统代理设置的工具。它提供了通过命令行和代码两种方式来管理代理配置。

## 安装

你可以使用 `go install` 命令安装 `winproxy`，也可以手动下载和安装它。

```shell
go install github.com/fhluo/winproxy/cmd/winproxy@latest
```

## 使用方法

### 命令行

- 使用 `winproxy` 命令显示当前的代理配置。
- 使用 `winproxy help` 命令查看帮助。

#### 示例

- 启用代理：`winproxy on`
- 禁用代理：`winproxy off`

### 代码

```go
package main

import (
	"github.com/fhluo/winproxy"
	"log"
)

func main() {
	// 读取当前的代理配置
	settings, err := winproxy.ReadSettings()
	if err != nil {
		log.Fatalln(err)
	}

	// 更改代理配置
	settings.Proxy = true
	settings.ProxyAddress = "127.0.0.1:8080"
	settings.Script = false
	settings.ScriptAddress = ""
	settings.AutoDetect = false
	settings.BypassList = []string{
		"<local>",
	}

	// 应用代理配置
	if err = settings.Apply(); err != nil {
		log.Fatalln(err)
	}
}

```
