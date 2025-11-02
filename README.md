<div align="center">

# winproxy

Change Windows system proxy settings
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

## Installation

### Using Cargo

```shell
cargo install winproxy
```

### Using Go

```shell
go install github.com/fhluo/winproxy/go/cmd/winproxy@latest
```

### Download Pre-built Binaries

Visit the [Releases](https://github.com/fhluo/winproxy/releases) page to download pre-built binaries.

## Usage

### View Current Proxy Settings

```shell
winproxy
```

This will display your current proxy configuration in a formatted table.

### Enable/Disable Proxy

```shell
# Enable proxy
winproxy -p true

# Disable proxy
winproxy -p false
```

### Set Proxy Server

```shell
# Set HTTP proxy
winproxy -p true --proxy-address "127.0.0.1:8080"

# Set proxy with bypass list (semicolon-separated)
winproxy -p true --proxy-address "127.0.0.1:8080" --bypass-list "localhost;127.*;<local>"
```

## Library Usage

Add `winproxy` to your `Cargo.toml`:

```toml
[dependencies]
winproxy = "0.5"
```

Example usage:

```rust
use winproxy::DefaultConnectionSettings;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Read current proxy settings
    let mut settings = DefaultConnectionSettings::from_registry()?;
    println!("Current settings: {:?}", settings);

    // Enable proxy and set address/bypass list
    settings.enable_proxy();
    settings.proxy_address = "127.0.0.1:8080".to_string();
    settings.set_bypass_list_from_str("localhost;127.*;<local>");

    // Apply settings
    settings.version += 1;
    settings.write_registry()?;
    println!("Proxy enabled!");

    Ok(())
}
```
