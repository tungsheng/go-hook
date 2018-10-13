package config

import (
	"github.com/joho/godotenv"
	"github.com/segmentio/go-env"
)

type server struct {
	Host string
	Addr string
	Key  string
}

var (
	// Server represents the informations about the server bindings.
	Server = &server{}
)

// Init load .env file
func Init() error {
	if err := godotenv.Load("../.env"); err != nil {
		return err
	}

	Server.Key, _ = env.Get("BIT_SECRET")

	return nil
}
