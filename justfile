mod go

set shell := ["nu", "-c"]

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
