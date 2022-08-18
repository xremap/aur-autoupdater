package packageinfo

type PackageInfo struct {
	Name       string
	GitHubInfo GitHubInfo
}

type GitHubInfo struct {
	Owner string
	Repo  string
}

var PackageInfos = map[string]PackageInfo{
	"xremap-x11-bin": {
		Name: "xremap-x11-bin",
		GitHubInfo: GitHubInfo{
			Owner: "k0kubun",
			Repo:  "xremap",
		},
	},
}
