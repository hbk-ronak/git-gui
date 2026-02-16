package main

import (
	"embed"
	"flag"

	"git-gui/backend"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	var projectPath string
	flag.StringVar(&projectPath, "project", "", "path to git project")
	flag.StringVar(&projectPath, "p", "", "path to git project (shorthand)")
	flag.Parse()

	backend.Run(assets, projectPath)
}
