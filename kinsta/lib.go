package kinsta

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mitchellh/cli"
	tCmd "github.com/roots/trellis-cli/cmd"
)

const apiUrl = "https://api.kinsta.com/v2"

type Command interface {
	UI() cli.Ui
	Flags() flag.FlagSet
}

func GetFlagValue(c Command, varName string) (value string, err error) {
	ui := c.UI()
	flags := c.Flags()

	value = flags.Lookup(varName).Value.String()
	if value != "" {
		return value, nil
	}

	envVarName := "KINSTA_" + strings.ToUpper(strings.ReplaceAll(varName, "-", "_"))
	value = os.Getenv(envVarName)
	if value != "" {
		return value, nil
	}

	value, err = ui.Ask(fmt.Sprintf("Enter %s:", varName))
	if value == "" {
		return "", fmt.Errorf("Error: %s is required.", varName)
	}

	if err != nil {
		return "", err
	}

	return value, nil
}

type SiteLabel struct {
	Id   json.Number `json:"id"`
	Name string      `json:"name"`
}

type Site struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	DisplayName string      `json:"display_name"`
	Status      string      `json:"status"`
	SiteLabels  []SiteLabel `json:"site_labels"`
}

type Company struct {
	Sites []Site `json:"sites"`
}

type SitesList struct {
	Error   string  `json:"error"`
	Company Company `json:"company"`
}

func Request(ui cli.Ui, apiKey string, url string, target interface{}) (errCode int) {
	spinner := tCmd.NewSpinner(
		tCmd.SpinnerCfg{
			FailMessage: "Failed",
		},
	)

	spinner.Start()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", apiUrl, url), nil)
	if err != nil {
		spinner.StopFailMessage(err.Error())
		spinner.StopFail()
		return 1
	}

	req.Header.Add("Authorization", "Bearer "+apiKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		spinner.StopFailMessage(err.Error())
		spinner.StopFail()
		return 1
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		spinner.StopFailMessage(readErr.Error())
		spinner.StopFail()
		return 1
	}

	// Unmarshal into a map to check for error property before unmarshalling to target
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		spinner.StopFailMessage(err.Error())
		spinner.StopFail()
		return 1
	}

	if errMsg, ok := raw["message"].(string); res.StatusCode != http.StatusOK && ok && errMsg != "" {
		spinner.StopFailMessage(errMsg)
		spinner.StopFail()
		return 1
	}

	if errMsg, ok := raw["error"].(string); ok && errMsg != "" {
		spinner.StopFailMessage(errMsg)
		spinner.StopFail()
		return 1
	}

	jsonErr := json.Unmarshal(body, &target)
	if jsonErr != nil {
		spinner.StopFailMessage(jsonErr.Error())
		spinner.StopFail()
		return 1
	}

	spinner.Stop()

	return 0
}

func ListSites(ui cli.Ui, apiKey string, company string) int {
	var sl SitesList
	if err := Request(ui, apiKey, fmt.Sprintf("sites/?company=%s", company), &sl); err != 0 {
		return 1
	}

	var trs []table.Row
	for _, v := range sl.Company.Sites {
		var lbls []string
		for _, lbl := range v.SiteLabels {
			lbls = append(lbls, lbl.Name)
		}
		tr := table.Row{v.Id, v.DisplayName, v.Name, v.Status, strings.Join(lbls, ", ")}
		trs = append(trs, tr)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "NAME", "DISPLAY NAME", "STATUS", "SITE LABELS"})
	t.AppendRows(trs)
	t.Render()

	return 0
}
