package cmd

import (
	"flag"
	"strings"

	"github.com/ItinerisLtd/trellis-cli-kinsta/kinsta"
	"github.com/hashicorp/cli"
	"github.com/roots/trellis-cli/trellis"
)

type SiteListCommand struct {
	ui      cli.Ui
	Trellis *trellis.Trellis
	flags   *flag.FlagSet
	apiKey  string
	company string
}

func (s SiteListCommand) UI() cli.Ui          { return s.ui }
func (s SiteListCommand) Flags() flag.FlagSet { return *s.flags }

func (c *SiteListCommand) init() {
	c.flags = flag.NewFlagSet("", flag.ContinueOnError)
	c.flags.Usage = func() { c.ui.Info(c.Help()) }
	c.flags.StringVar(&c.apiKey, "api-key", "", "The API key used to query Kinsta APIs.")
	c.flags.StringVar(&c.company, "company", "", "The company ID to query.")
}

func NewSiteListCommand(ui cli.Ui, trellis *trellis.Trellis) *SiteListCommand {
	c := &SiteListCommand{ui: ui, Trellis: trellis}
	c.init()
	return c
}

func (c *SiteListCommand) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 1
	}

	args = c.flags.Args()

	apiKey, err := kinsta.GetFlagValue(c, "api-key")
	if err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	company, err := kinsta.GetFlagValue(c, "company")
	if err != nil {
		c.ui.Error(err.Error())
		return 1
	}

	return kinsta.ListSites(c.ui, apiKey, company)
}

func (c *SiteListCommand) Synopsis() string {
	return "List all sites associated to the given company ID"
}

func (c *SiteListCommand) Help() string {
	helpText := `
Usage: trellis kinsta site list [options]

List all sites associated to the given company ID:

  $ trellis kinsta site list --company=123 --api-key=123

Options:
      --company The company ID to query
      --api-key The API key used to query Kinsta APIs
  -h, --help    show this help
`

	return strings.TrimSpace(helpText)
}
