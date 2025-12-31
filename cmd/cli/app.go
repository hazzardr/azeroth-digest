package cli

type App struct {
	Bot     BotCommand     `cmd:"" help:"Perform discord actions."`
	Pricing PricingCommand `cmd:"" help:"Perform pricing actions."`
}
