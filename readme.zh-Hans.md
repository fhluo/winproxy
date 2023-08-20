<div align="center">

# winproxy

更改 Windows 系统代理配置
<br><br>
<a href="https://github.com/fhluo/winproxy/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/fhluo/winproxy" alt="license">
</a>
<a href="https://github.com/fhluo/winproxy/actions/workflows/build.yaml">
    <img src="https://github.com/fhluo/winproxy/actions/workflows/build.yaml/badge.svg" alt="build workflow">
</a>
<a href="https://goreportcard.com/report/github.com/fhluo/winproxy">
    <img src="https://goreportcard.com/badge/github.com/fhluo/winproxy" alt="go report card">
</a>

<samp>

**[English](readme.md)** ┃ **[简体中文](readme.zh-Hans.md)**

</samp>
</div>

## 用法

winproxy 提供了两种方式来更改 Windows 系统的代理配置：通过命令行或通过代码。

### 命令行

你可以使用 `go install` 命令安装 winproxy，也可以手动下载和安装它。

```shell
go install github.com/fhluo/winproxy/cmd/winproxy@latest
```

- 使用 `winproxy` 命令显示当前的代理配置。
- 使用 `winproxy help` 命令查看帮助。

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
