package main

import (
	"os"

	"github.com/urfave/cli"

	"go.polydawn.net/reppl/actions"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name: "put",
			Subcommands: []cli.Command{
				{
					Name:   "hash",
					Action: actions.PutHash,
				},
				{
					Name:   "file",
					Action: actions.PutFile,
				},
			},
		},
		{
			Name:   "eval",
			Action: actions.Eval,
		},
		{
			Name:   "init",
			Action: actions.Init,
		},
		{
			Name:   "show",
			Action: actions.Show,
		},

		{
			Name:   "rm",
			Action: actions.Remove,
		},
	}

	app.Run(os.Args)
}
