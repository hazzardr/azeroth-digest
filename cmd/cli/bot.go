package cli

import (
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/hazzardr/azeroth-digest/internal/discord"
)

type ServeCommand struct {
	Token string `required:"true" help:"Discord bot token"`
	DSN   string `required:"true" default:"data/azdg.duckdb" help:"database path"`
}
type SyncCommand struct {
	Token string `required:"true" help:"Discord bot token"`
}

type BotCommand struct {
	Serve ServeCommand `cmd:"" help:"Start the Discord Bot"`
	Sync  SyncCommand  `cmd:"" help:"Sync the Discord Application Commands. Necessary when adding a new slash command."`
}

func (s *ServeCommand) Run() error {
	b, err := discord.NewBot(s.Token)
	if err != nil {
		slog.Error("failed to create discord bot", slog.Any("err", err))
		os.Exit(1)
	}
	err = b.Start()
	slog.Info("bot is running")
	if err != nil {
		return errors.Join(errors.New("failed to start the discord bot"), err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sig
	return nil
}
