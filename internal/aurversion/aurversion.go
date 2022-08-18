package aurversion

import (
	"fmt"
	"io"
	"io/ioutil"

	pkgbuild "github.com/mikkeloscar/gopkgbuild"
	"github.com/njkevlani/aur-autoupdater/internal/version"
	"github.com/sirupsen/logrus"
)

type AURVersion struct {
	version string
}

func (version *AURVersion) Version() string {
	return version.version
}

func GetAURVersion(srcinfoFile io.Reader) (version.Version, error) {
	srcinfoFileContent, err := ioutil.ReadAll(srcinfoFile)

	if err != nil {
		return nil, err
	}

	parsedScrinfo, err := pkgbuild.ParseSRCINFOContent(srcinfoFileContent)

	if err != nil {
		return nil, err
	}

	version := AURVersion{version: string(parsedScrinfo.Pkgver)}

	logrus.WithField("AURVersion", fmt.Sprintf("%#v", version)).Info("AURVersion")

	return &version, nil
}
