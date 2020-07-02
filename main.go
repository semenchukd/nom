package main

import (
	"fmt"
	"github.com/semenchukd/nom/commands"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "newenv",
				Usage: "TODO: add usage",
				Action:  commands.Newenv,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "up",
						Usage:   "Containers for docker-compose up",
					},
					&cli.BoolFlag{
						Name:    "noup",
						Usage: "Skip docker-compose up",
					},
					&cli.BoolFlag{
						Name:    "nobs",
						Usage: "Skip bootstrap-api",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("[error]:", err)
	}

}