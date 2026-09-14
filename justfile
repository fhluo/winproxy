mod go

set shell := ["nu", "-c"]
set script-interpreter := ["nu"]
set indentation := "  "

default:
  @just --list

run *args:
  cargo run -p winproxy-cli -- {{args}}

build *args:
  cargo build -p winproxy-cli {{args}}

release: (build "--release")

install:
  cargo install --path winproxy-cli

test: go::test
  cargo test --workspace

[group: 'release']
[script]
check-version version:
  if "{{version}}" not-like '^v?\d+\.\d+\.\d+$' {
    print $"(ansi red)invalid version: {{version}}, expected a version like 0.5.0 or v0.5.0(ansi reset)"
    exit 1
  }

[group: 'release']
[script]
bump-version version: (check-version version)
  let status = (git status --porcelain --untracked-files=no)
  if ($status | is-not-empty) {
    print $"(ansi red)working tree is not clean, please commit or stash changes first(ansi reset)"
    exit 1
  }

  let ver = ("{{version}}" | str replace -r '^v' '')
  print $"(ansi light_gray)Syncing version to ($ver)...(ansi reset)"
  just sync-version $ver

  print $"(ansi light_gray)Committing...(ansi reset)"
  git add winproxy/Cargo.toml winproxy-cli/Cargo.toml Cargo.lock winproxy/src/lib.rs README.md README.zh-Hans.md
  git commit -m $"chore: bump version to ($ver)"

  print $"(ansi light_gray)Tagging v($ver)...(ansi reset)"
  git tag $"v($ver)" -m $"v($ver)"

  print $"(ansi light_gray)Publishing winproxy ($ver) to crates.io...(ansi reset)"
  cargo publish -p winproxy

  print $"(ansi light_gray)Publishing winproxy-cli ($ver) to crates.io...(ansi reset)"
  cargo publish -p winproxy-cli

  print $"(ansi green)✓ Bumped to ($ver), published winproxy and winproxy-cli(ansi reset)"

[group: 'release']
[script]
sync-version version:
  let version = "{{version}}" | into semver --loose
  let major_minor = ($version | into record | $"($in.major).($in.minor)")
  let package_version = ($version | into record | reject prefix | into semver)

  open --raw winproxy/Cargo.toml
  | str replace -r '(?m)^version = "\d+\.\d+\.\d+"' $"version = \"($package_version)\""
  | save -f winproxy/Cargo.toml

  open --raw winproxy-cli/Cargo.toml
  | str replace -r '(?m)^version = "\d+\.\d+\.\d+"' $"version = \"($package_version)\""
  | str replace -r 'winproxy = \{ version = "\d+\.\d+"' $"winproxy = { version = \"($major_minor)\""
  | save -f winproxy-cli/Cargo.toml

  # let cargo refresh winproxy/winproxy-cli entries in Cargo.lock
  cargo metadata --format-version 1 | ignore

  open --raw winproxy/src/lib.rs
  | str replace -r 'winproxy = "\d+\.\d+"' $"winproxy = \"($major_minor)\""
  | save -f winproxy/src/lib.rs

  open --raw README.md
  | str replace -r 'winproxy = "\d+\.\d+"' $"winproxy = \"($major_minor)\""
  | save -f README.md

  open --raw README.zh-Hans.md
  | str replace -r 'winproxy = "\d+\.\d+"' $"winproxy = \"($major_minor)\""
  | save -f README.zh-Hans.md
