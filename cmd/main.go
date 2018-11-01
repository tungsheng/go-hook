package main

import (
	"os"
	"time"

	"github.com/go-hook/config"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	cli "gopkg.in/urfave/cli.v2"
)

// Version set at compile-time
var Version = "v1.0.0-dev"

func setupLogging() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if config.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if config.Logs.Pretty {
		zlog.Logger = zlog.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: !config.Logs.Color,
			},
		)
	}
}

func main() {
	if env := os.Getenv("GOHOOK_ENV_FILE"); env != "" {
		godotenv.Load(env)
	}

	app := &cli.App{
		Name:      "gohook",
		Usage:     "Gohook Server",
		Copyright: "Copyright (c) 2018 Tony Lee",
		Version:   Version,
		Compiled:  time.Now(),

		Authors: []*cli.Author{
			{
				Name:  "Tony Lee",
				Email: "tungsheng@gmail.com",
			},
		},

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Value:       true,
				Usage:       "Activate debug information",
				EnvVars:     []string{"INN_DEBUG"},
				Destination: &config.Debug,
			},
		},

		Before: func(c *cli.Context) error {
			setupLogging()
			return nil
		},

		Commands: []*cli.Command{
			Server(),
		},
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "Show the help, so what you see now",
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print the current version of that tool",
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
