package main

import (
	"github.com/codegangsta/cli"
	"github.com/nautsio/presenter/command"
)

// Commands is an array containing the available commands.
var Commands = []cli.Command{
	commandInit,
	commandServe,
}

var commandInit = cli.Command{
	Name:      "init",
	ShortName: "i",
	Usage:     "Initialize an empty presentation directory",
	ArgsUsage: "<destination path>",
	Action:    command.Init,
}

var commandServe = cli.Command{
	Name:      "serve",
	ShortName: "s",
	Usage:     "Serve a presentation directory",
	ArgsUsage: "<presentation directory>",
	Action:    command.Serve,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "master, m",
			Usage: "Start presenter in master mode",
		},
		cli.StringFlag{
			Name:  "theme, t",
			Usage: "Use one of the built in themes",
		},
	},
}
