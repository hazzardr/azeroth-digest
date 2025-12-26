package discord

import (
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/handler/middleware"
)

type Bot struct {
	client *bot.Client
}

func NewBot(discordToken string) (*Bot, error) {
	b := &Bot{}
	r := initializeEventListeners()

	client, err := disgo.New(discordToken,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
				gateway.IntentGuilds,
				gateway.IntentDirectMessages,
			),
			gateway.WithPresenceOpts(
				gateway.WithListeningActivity("Crunching the numbers..."),
				gateway.WithOnlineStatus(discord.OnlineStatusOnline),
			),
		),
		bot.WithEventListeners(r),
	)
	if err != nil {
		return nil, err
	}
	b.client = client
	return b, nil
}

// Start starts the bot, initializing the discord client against the gateway.
func (b *Bot) Start() error {
	return nil
}

func (b *Bot) GracefulShutdown() {}

func initializeEventListeners() *handler.Mux {
	r := handler.New()
	r.Use(middleware.Logger)
	r.SlashCommand("/ping", HandlePing)
	return r
}
