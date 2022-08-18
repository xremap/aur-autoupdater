package version

type Version interface {
	Version() string
}

func stripV(v string) string {
	if len(v) > 0 && v[0] == 'v' {
		return v[1:]
	}

	return v
}

func Equal(v1, v2 Version) bool {
	return stripV(v1.Version()) == stripV(v2.Version())
}
