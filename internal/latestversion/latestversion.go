package latestversion

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/njkevlani/aur-autoupdater/internal/version"
	"github.com/sirupsen/logrus"
)

// Currently, only fetches from GitHub
func GetLatestVersion(owner, repo string) (version.Version, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	version := LatestGitHubVersion{}
	err = json.NewDecoder(resp.Body).Decode(&version)

	if err != nil {
		return nil, err
	}

	logrus.WithField("LatestGitHubVersion", fmt.Sprintf("%#v", version)).Info("LatestGitHubVersion")

	return &version, nil
}
