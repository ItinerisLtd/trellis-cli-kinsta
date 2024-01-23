package kinsta

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mitchellh/cli"
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

func Request(ui cli.Ui, accessToken string, url string) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", apiUrl, url), nil)
	if err != nil {
		ui.Error(err.Error())
	}

	var bearer = "Bearer " + accessToken

	req.Header.Add("Authorization", bearer)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		ui.Error(err.Error())
	}

	ui.Info(fmt.Sprintf("status code %d\n", res.StatusCode))

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		ui.Error(err.Error())
	}

	ui.Info(string(resBody))
}

func ListSites(ui cli.Ui, accessToken string, company string) {
	Request(ui, accessToken, fmt.Sprintf("sites/?company=%s", company))
}
