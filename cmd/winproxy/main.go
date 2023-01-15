package main

import (
	"github.com/fhluo/winproxy/cmd"
	"golang.org/x/exp/slog"
	"os"
)

func init() {
	slog.SetDefault(slog.New(
		slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					a.Key = ""
				}
				return a
			},
		}.NewTextHandler(os.Stderr),
	))
}

func main() {
	cmd.Execute()
}
