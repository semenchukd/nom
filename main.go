package main

import (
	"fmt"
	"os"

	"github.com/semenchukd/nom/commands"
	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:     "Nsone Macros",
		Usage:    "Set of helpers for local dev",
		HelpName: "nom",
		Commands: []*cli.Command{
			{
				Name:   "newenv",
				Usage:  "Setup new evironment",
				Action: commands.Newenv,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "up",
						Usage: "Containers for docker-compose up",
					},
					&cli.BoolFlag{
						Name:  "noup",
						Usage: "Skip docker-compose up",
					},
					&cli.BoolFlag{
						Name:  "nobs",
						Usage: "Skip bootstrap-api",
					},
				},
			},
			{
				Name:   "dc",
				Usage:  "Docker-compose ops",
				Action: commands.DC,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "c",
						Usage: "List of containers",
					},
				},
			},
			{
				Name:   "build",
				Usage:  "Build and restart target container",
				Action: commands.Build,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("[error]:", err)
	}

}
