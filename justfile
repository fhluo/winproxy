mod go
mod version

set shell := ["nu", "-c"]
set script-interpreter := ["nu"]
set indentation := "  "
set default-list := true

run *args:
  cargo run -p winproxy-cli -- {{args}}

build *args:
  cargo build -p winproxy-cli {{args}}

release: (build "--release")

install:
  cargo install --path winproxy-cli

test:
  cargo test --workspace
