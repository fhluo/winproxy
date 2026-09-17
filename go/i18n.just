set shell := ["nu", "-c"]
set script-interpreter := ["nu"]
set indentation := "  "
set default-list := true

setup:
  go install -v github.com/nicksnyder/go-i18n/v2/goi18n@latest

[working-directory: 'cmd']
update:
  goi18n extract -outdir locales
  goi18n merge -outdir locales locales/active.*.toml locales/translate.*.toml
  cd locales; cp active.en.toml en.toml; cp active.zh-Hans.toml zh-Hans.toml
