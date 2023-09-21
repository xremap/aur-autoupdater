package main

import "github.com/xremap/aur-autoupdater/internal/processor"

func main() {
	packages := []string{
		"xremap-x11-bin",
		"xremap-hypr-bin",
	}
	for _, name := range packages {
		err := processor.Process(name)
		if err != nil {
			panic(err)
		}
	}
}
