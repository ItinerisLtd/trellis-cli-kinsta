package cmd

import (
	"flag"
	"strings"

	"github.com/ItinerisLtd/trellis-cli-kinsta/kinsta"
	"github.com/mitchellh/cli"
	"github.com/roots/trellis-cli/trellis"
)

func NewSiteListCommand(ui cli.Ui, trellis *trellis.Trellis) *SiteListCommand {
	c := &SiteListCommand{UI: ui, Trellis: trellis}
	c.init()
	return c
}

type SiteListCommand struct {
	UI      cli.Ui
	Trellis *trellis.Trellis
	flags   *flag.FlagSet
	company string
}

func (c *SiteListCommand) init() {
	c.flags = flag.NewFlagSet("", flag.ContinueOnError)
	c.flags.Usage = func() { c.UI.Info(c.Help()) }
	c.flags.StringVar(&c.company, "company", "", "The company ID to query.")
}

func (c *SiteListCommand) Run(args []string) int {
	if err := c.Trellis.LoadProject(); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.Trellis.CheckVirtualenv(c.UI)

	if err := c.flags.Parse(args); err != nil {
		return 1
	}

	args = c.flags.Args()

	accessToken, err := kinsta.GetAccessToken(c.UI)
	if err != nil {
		c.UI.Error("Error: DigitalOcean access token is required.")
		return 1
	}
	kinsta.ListSites(c.UI, accessToken, c.company)

	return 0
}

func (c *SiteListCommand) Synopsis() string {
	return "help me"
}

func (c *SiteListCommand) Help() string {
	helpText := `
Usage: trellis kinsta site list [options]

List all sites associated to the given company ID:

  $ trellis kinsta site list --company=123

Arguments:
  ENVIRONMENT Name of environment (ie: production)

Options:
      --company The company ID to query
  -h, --help    show this help
`

	return strings.TrimSpace(helpText)
}
