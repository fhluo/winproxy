<div align="center">

# winproxy

Change Windows system proxy settings
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

## Usage

winproxy provides two ways to change Windows system proxy settings. One is through the command line and the other is through programming.

### Through the command line

Use`go install` to install or download and install it manually.

```shell
go install github.com/fhluo/winproxy/cmd/winproxy@latest
```

- Use `winproxy` to show the current proxy settings.
- Use `winproxy help` to view help.

### Through programming

```go
package main

import (
	"github.com/fhluo/winproxy"
	"log"
)

func main() {
	// Read
	settings, err := winproxy.ReadSettings()
	if err != nil {
		log.Fatalln(err)
	}

	// Change
	settings.Proxy = true
	settings.ProxyAddress = "127.0.0.1:8080"
	settings.Script = false
	settings.ScriptAddress = ""
	settings.AutoDetect = false
	settings.BypassList = []string{
		"<local>",
	}

	// Apply
	if err = settings.Apply(); err != nil {
		log.Fatalln(err)
	}
}

```
