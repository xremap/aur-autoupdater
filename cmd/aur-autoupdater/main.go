package main

import "github.com/xremap/aur-autoupdater/internal/processor"

func main() {
	packages := []string{
		"xremap-gnome-bin",
		"xremap-hypr-bin",
		"xremap-x11-bin",
	}
	for _, name := range packages {
		err := processor.Process(name)
		if err != nil {
			panic(err)
		}
	}
}
