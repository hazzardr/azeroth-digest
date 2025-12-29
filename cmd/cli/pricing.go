package cli

type PricingCommand struct {
	Scrape ScrapeCommand `cmd:"" help:"Scrape pricing data from TSM API."`
}

type ScrapeCommand struct {
	Token    string `required:"true" help:"TSM client secret / API token"`
	Out      string `required:"true" help:"Output format" default:"file"`
	URL      string `required:"true" help:"URL to TSM API" default:"https://pricing-api.tradeskillmaster.com"`
	ClientID string `required:"true" help:"TSM client ID" default:"c260f00d-1071-409a-992f-dda2e5498536"`
	AH       int    `required:"true" help:"Auction house ID to scrape" default:"554"`
}

func (s *ScrapeCommand) Run() error {

	return nil
}
