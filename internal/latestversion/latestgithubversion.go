package latestversion

type LatestGitHubVersion struct {
	TagName string `json:"tag_name"`
}

func (version *LatestGitHubVersion) Version() string {
	return version.TagName
}
