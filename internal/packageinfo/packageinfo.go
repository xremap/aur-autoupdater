package packageinfo

import "fmt"

type PackageInfo struct {
	Name              string
	GitHubInfo        GitHubInfo
	GitHubInfoAarch64 GitHubInfo
	PkgbuildInfo      PkgbuildInfo
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
	"xremap-gnome-bin": {
		Name: "xremap-gnome-bin",
		GitHubInfo: GitHubInfo{
			Owner: "k0kubun",
			Repo:  "xremap",
			ReleaseAssetURL: func(version string) string {
				return fmt.Sprintf("https://github.com/k0kubun/xremap/releases/download/v%s/xremap-linux-x86_64-gnome.zip", version)
			},
		},
		GitHubInfoAarch64: GitHubInfo{
			Owner: "k0kubun",
			Repo:  "xremap",
			ReleaseAssetURL: func(version string) string {
				return fmt.Sprintf("https://github.com/k0kubun/xremap/releases/download/v%s/xremap-linux-aarch64-gnome.zip", version)
			},
		},
		PkgbuildInfo: PkgbuildInfo{
			PkgbuildTemplateFilepath: "assets/xremap-gnome-bin/PKGBUILD.tmpl",
			SrcinfoTemplateFilepath:  "assets/xremap-gnome-bin/.SRCINFO.tmpl",
		},
	},
	"xremap-wlroots-bin": {
		Name: "xremap-wlroots-bin",
		GitHubInfo: GitHubInfo{
			Owner: "k0kubun",
			Repo:  "xremap",
			ReleaseAssetURL: func(version string) string {
				return fmt.Sprintf("https://github.com/k0kubun/xremap/releases/download/v%s/xremap-linux-x86_64-wlroots.zip", version)
			},
		},
		PkgbuildInfo: PkgbuildInfo{
			PkgbuildTemplateFilepath: "assets/xremap-wlroots-bin/PKGBUILD.tmpl",
			SrcinfoTemplateFilepath:  "assets/xremap-wlroots-bin/.SRCINFO.tmpl",
		},
	},
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
}
