package main

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tungsheng/go-hook/config"
	"github.com/tungsheng/go-hook/router"
	cli "gopkg.in/urfave/cli.v2"
)

var (
	defaultAddr = ":3003"
)

// Server provides the sub-command to start the API server.
func Server() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start the Gohook API",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Value:       "http://localhost:3003",
				Usage:       "External access to server",
				EnvVars:     []string{"GOHOOK_SERVER_HOST"},
				Destination: &config.Server.Host,
			},
			&cli.StringFlag{
				Name:        "addr",
				Value:       defaultAddr,
				Usage:       "Address to bind the server",
				EnvVars:     []string{"GOHOOK_SERVER_ADDR"},
				Destination: &config.Server.Addr,
			},
			&cli.StringFlag{
				Name:        "root",
				Value:       "/",
				Usage:       "Root folder of the app",
				EnvVars:     []string{"GOHOOK_SERVER_ROOT"},
				Destination: &config.Server.Root,
			},
		},
		Action: func(c *cli.Context) error {
			// load global script
			log.Info().Msg("Initial module engine.")
			//router.GlobalInit()

			server := &http.Server{
				Addr:         config.Server.Addr,
				Handler:      router.Load(),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
			}

			if err := startServer(server); err != nil {
				log.Fatal().Err(err)
			}

			return nil
		},
	}
}
