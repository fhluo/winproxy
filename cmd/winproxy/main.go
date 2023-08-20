package main

import (
	"github.com/fhluo/winproxy/cmd"
	"log/slog"
	"os"
)

func init() {
	slog.SetDefault(slog.New(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					a.Key = ""
				}
				return a
			},
		}),
	))
}

func main() {
	cmd.Execute()
}
