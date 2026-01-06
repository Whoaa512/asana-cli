package main

import (
	"os"

	"github.com/whoaa512/asana-cli/internal/cli"
)

var version = "dev"

func main() {
	cli.SetVersion(version)
	os.Exit(cli.Execute())
}
