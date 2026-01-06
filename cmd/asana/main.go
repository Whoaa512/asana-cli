package main

import (
	"os"

	"github.com/whoaa512/asana-cli/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
