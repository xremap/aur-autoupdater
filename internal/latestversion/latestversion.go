package latestversion

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"

	"github.com/xremap/aur-autoupdater/internal/version"
	"github.com/sirupsen/logrus"
)

// Currently, only fetches from GitHub
func GetLatestVersion(owner, repo string) (version.Version, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

	versionFromGitHub := LatestGitHubVersion{}
	resp, err := resty.New().R().
		SetAuthToken(os.Getenv("GH_TOKEN_FOR_AUR_AUTO_UPDATE")).
		SetResult(&versionFromGitHub).
		Get(url)

	if err != nil {
		logrus.WithError(err).
			WithField("resp", resp).
			WithField("url", url).
			Error("got error while making request to GitHub")

		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		err = errors.New("got non-200 response from GitHub")
		logrus.WithError(err).
			WithField("resp", resp).
			WithField("url", url).
			Error()

		return nil, err
	}

	// TODO: This should not be required after better error handling above.
	// If that is the cases, this will be removed in upcoming commits.
	if len(versionFromGitHub.Version()) == 0 {
		logrus.Error("latest version length is 0")
		return nil, errors.New("latest version length is 0")
	}

	logrus.WithField("LatestGitHubVersion", fmt.Sprintf("%#v", versionFromGitHub)).Info("LatestGitHubVersion")

	return &versionFromGitHub, nil
}
