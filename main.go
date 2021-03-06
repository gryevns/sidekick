package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config-source",
			EnvVar: "SIDEKICK_CONFIG_SOURCE",
			Value:  "dynamodb",
		},
		cli.StringFlag{
			Name:   "config-table",
			EnvVar: "SIDEKICK_CONFIG_TABLE",
		},
		cli.StringSliceFlag{
			Name:  "env, e",
			Usage: "set environment variables",
		},
	}
	app.Before = func(c *cli.Context) error {
		// TODO: support multiple config sources
		var configSource ConfigSource
		switch c.String("config-source") {
		case "dynamodb":
			configSource = &DynamoDBConfigSource{
				Table: c.String("config-table"),
				Key:   "common", // TODO: parameterise
			}
		default:
			return cli.NewExitError("couldn't find that config source type", 2)
		}

		appendConfigSource(c, configSource)

		envs := c.StringSlice("env")
		if len(envs) > 0 {
			appendConfigSource(c, NewStringSliceConfigSource(envs))
		}

		return nil
	}
	app.Commands = []cli.Command{
		commandEnv,
		commandRun,
		commandStart,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
