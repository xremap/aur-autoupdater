package main

import "github.com/xremap/aur-autoupdater/internal/processor"

func main() {
	err := processor.Process("xremap-x11-bin")
	if err != nil {
		panic(err)
	}
}
