package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "presenter"
	app.Version = Version
	app.Usage = "HTTP server for serving your slides."
	app.Author = "nauts.io"
	app.Commands = Commands

	app.Run(os.Args)
}
