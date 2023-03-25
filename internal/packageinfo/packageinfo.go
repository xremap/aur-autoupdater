package packageinfo

import "fmt"

type PackageInfo struct {
	Name         string
	GitHubInfo   GitHubInfo
	PkgbuildInfo PkgbuildInfo
}

type GitHubInfo struct {
	Owner           string
	Repo            string
	ReleaseAssetURL func(string) string
}

type PkgbuildInfo struct {
	PkgbuildTemplateFilepath string
	SrcinfoTemplateFilepath  string
}

var PackageInfos = map[string]PackageInfo{
	"xremap-x11-bin": {
		Name: "xremap-x11-bin",
		GitHubInfo: GitHubInfo{
			Owner: "k0kubun",
			Repo:  "xremap",
			ReleaseAssetURL: func(version string) string {
				return fmt.Sprintf("https://github.com/k0kubun/xremap/releases/download/v%s/xremap-linux-x86_64-x11.zip", version)
			},
		},
		PkgbuildInfo: PkgbuildInfo{
			PkgbuildTemplateFilepath: "assets/xremap-x11-bin/PKGBUILD.tmpl",
			SrcinfoTemplateFilepath:  "assets/xremap-x11-bin/.SRCINFO.tmpl",
		},
	},
	"xremap-hypr-bin": {
		Name: "xremap-hypr-bin",
		GitHubInfo: GitHubInfo{
			Owner: "k0kubun",
			Repo:  "xremap",
			ReleaseAssetURL: func(version string) string {
				return fmt.Sprintf("https://github.com/k0kubun/xremap/releases/download/v%s/xremap-linux-x86_64-hypr.zip", version)
			},
		},
		PkgbuildInfo: PkgbuildInfo{
			PkgbuildTemplateFilepath: "assets/xremap-hypr-bin/PKGBUILD.tmpl",
			SrcinfoTemplateFilepath:  "assets/xremap-hypr-bin/.SRCINFO.tmpl",
		},
	},
}
