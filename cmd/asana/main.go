package main

import (
	"os"
	"runtime/debug"

	"github.com/whoaa512/asana-cli/internal/cli"
)

var version = "dev"

func main() {
	cli.SetVersion(resolveVersion())
	os.Exit(cli.Execute())
}

func resolveVersion() string {
	if version != "dev" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return version
}
