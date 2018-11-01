package config

import (
	"github.com/joho/godotenv"
	"github.com/segmentio/go-env"
)

type server struct {
	Host string
	Addr string
	Key  string
	Root string
}

type logs struct {
	Color  bool
	Debug  bool
	Pretty bool
}

// ContextKey for context package
type ContextKey string

func (c ContextKey) String() string {
	return "backend context key " + string(c)
}

var (
	// Debug represents the flag to enable or disable debug logging.
	Debug bool

	// Server represents the informations about the server bindings.
	Server = &server{}

	// ContextKeyUser for user
	ContextKeyUser = ContextKey("user")

	// Logs for zerolog
	Logs = &logs{}
)

// Init load .env file
func Init() error {
	if err := godotenv.Load("../.env"); err != nil {
		return err
	}

	Server.Key, _ = env.Get("BIT_SECRET")

	return nil
}
