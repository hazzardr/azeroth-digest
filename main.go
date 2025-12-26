package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/log"
	"github.com/hazzardr/azeroth-digest/cmd/cli"
)

func main() {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})

	slog.SetDefault(slog.New(logger))

	azdg := cli.App{}
	ctx := kong.Parse(&azdg,
		kong.Name("azeroth-digest"),
		kong.Description("Azeroth Digest CLI."),
		kong.UsageOnError(),
	)

	err := ctx.Run()
	if err != nil {
		slog.Error("failed to start azeroth-digest", "error", err)
	}
}
