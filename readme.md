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

winproxy provides two ways to change the proxy settings of a Windows system: through the command line or through code.

### Command line

You can install winproxy using the `go install` command or by downloading and installing it manually.

```shell
go install github.com/fhluo/winproxy/cmd/winproxy@latest
```

- Use the `winproxy` command to display the current proxy settings.
- Use the `winproxy help` command to view the help.

### Code

```go
package main

import (
	"github.com/fhluo/winproxy"
	"log"
)

func main() {
	// Read the current proxy settings
	settings, err := winproxy.ReadSettings()
	if err != nil {
		log.Fatalln(err)
	}

	// Change the proxy settings
	settings.Proxy = true
	settings.ProxyAddress = "127.0.0.1:8080"
	settings.Script = false
	settings.ScriptAddress = ""
	settings.AutoDetect = false
	settings.BypassList = []string{
		"<local>",
	}

	// Apply the proxy settings
	if err = settings.Apply(); err != nil {
		log.Fatalln(err)
	}
}

```
