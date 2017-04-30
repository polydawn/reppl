package main

import (
	"os"

	"github.com/urfave/cli"

	"go.polydawn.net/reppl/actions"
	"go.polydawn.net/reppl/lib/efmt"
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
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "warehouse",
							Usage: "A URL giving coordinates to a warehouse where we should be able to find this ware.",
						},
						cli.StringFlag{
							Name:  "kind",
							Usage: "The kind of transit format this ware is.  Defaults to tar if unspecified.",
						},
					},
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
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "Optional -- Force this formula to *always* be run, even if we have a memoized result for its current setup.",
				},
				cli.StringSliceFlag{
					Name:  "env, e",
					Usage: "Apply additional environment vars to formula before launch.  Format like '-e KEY=val'",
				},
			},
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
			Name:   "unpack",
			Action: actions.Unpack,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "warehouse",
					Usage: "Optional -- A URL giving coordinates to a warehouse where we should be able to find this ware.",
				},
			},
		},
		{
			Name:   "rm",
			Action: actions.Remove,
		},
	}

	cli.ErrWriter = efmt.LineFlankingWriter(os.Stderr,
		efmt.Ansi(efmt.Ansi_textBrightRed), efmt.Ansi(efmt.Ansi_reset))

	app.Run(os.Args)
}
