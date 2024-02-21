package main

import (
	"github.com/cduggn/ccexplorer/internal/cli"
	"github.com/sagikazarmark/slog-shim"
	"os"
)

func main() {
	root := cli.RootCommand()

	if err := root.Execute(); err != nil {
		slog.Error("error", ErrAttr(err))
		os.Exit(126)
	}

}

func ErrAttr(err error) slog.Attr {
	return slog.Any("error", err)
}
