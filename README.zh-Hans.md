<div align="center">

# winproxy

更改 Windows 系统代理配置
<br><br>
<a href="https://github.com/fhluo/winproxy/actions/workflows/build.yaml">
<img src="https://github.com/fhluo/winproxy/actions/workflows/build.yaml/badge.svg" alt="build workflow"></a>
<a href="https://crates.io/crates/winproxy">
<img src="https://img.shields.io/crates/v/winproxy" alt="version"></a>
<a href="https://pkg.go.dev/github.com/fhluo/winproxy/go">
<img src="https://img.shields.io/github/v/tag/fhluo/winproxy?filter=go%2F*&label=pkg"></a>

<samp>

**[English](README.md)** ┃ **[简体中文](README.zh-Hans.md)**

</samp>
</div>

## 安装

### 使用 Cargo

```shell
cargo install winproxy
```

### 使用 Go

```shell
go install github.com/fhluo/winproxy/go/cmd/winproxy@latest
```

### 手动下载二进制文件

在 [Releases](https://github.com/fhluo/winproxy/releases) 页面获取预编译版本。

## 使用

### 查看当前代理配置

```shell
winproxy
```

以表格形式显示当前的代理配置。

### 启用或禁用代理

```shell
# 开启代理
winproxy -p true

# 关闭代理
winproxy -p false
```

### 设置代理服务器

```shell
# 设置 HTTP 代理
winproxy -p true --proxy-address "127.0.0.1:8080"

# 带忽略规则的代理配置
winproxy -p true --proxy-address "127.0.0.1:8080" --bypass-list "localhost,127.*,<local>"
```
