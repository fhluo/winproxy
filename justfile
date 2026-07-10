mod go

set shell := ["nu", "-c"]

default:
  @just --list

run *args:
  cargo run -p winproxy-cli -- {{args}}

build:
  cargo build -p winproxy-cli

build-release:
  cargo build -p winproxy-cli --release

install:
  cargo install --path winproxy-cli

test: go::test
  cargo test --workspace
