package kinsta

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mitchellh/cli"
	tCmd "github.com/roots/trellis-cli/cmd"
)

const accessTokenEnvVar = "KINSTA_API_ACCESS_TOKEN"
const apiUrl = "https://api.kinsta.com/v2"

func GetAccessToken(ui cli.Ui) (accessToken string, err error) {
	accessToken = os.Getenv(accessTokenEnvVar)

	if accessToken == "" {
		ui.Info(fmt.Sprintf("%s environment variable not set.", accessTokenEnvVar))
		accessToken, err = ui.Ask("Enter Access token:")

		if err != nil {
			return "", err
		}
	}

	return accessToken, nil
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
	Company Company `json:"company"`
}

func Request(ui cli.Ui, accessToken string, url string, target interface{}) {
	spinner := tCmd.NewSpinner(
		tCmd.SpinnerCfg{
			FailMessage: "Failed",
		},
	)
	spinner.Start()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", apiUrl, url), nil)
	if err != nil {
		spinner.StopFail()
		ui.Error(err.Error())
		return
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		spinner.StopFail()
		ui.Error(err.Error())
		return
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		spinner.StopFail()
		ui.Error(readErr.Error())
		return
	}

	jsonErr := json.Unmarshal(body, &target)
	if jsonErr != nil {
		spinner.StopFail()
		ui.Error(jsonErr.Error())
		return
	}

	spinner.Stop()
}

func ListSites(ui cli.Ui, accessToken string, company string) {
	var sl SitesList
	Request(ui, accessToken, fmt.Sprintf("sites/?company=%s", company), &sl)

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
}
