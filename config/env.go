package config

import (
	"errors"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

type Environment struct {
	TOKEN  string
	Prefix string
}

var env *Environment

func initializeEnvironment() error {
	env = &Environment{
		TOKEN: strings.TrimSpace(os.Getenv("TOKEN")),
	}

	if len(env.TOKEN) == 0 {
		return errors.New("TOKEN or INVITE_URL is not set")
	}

	prefix := strings.TrimSpace(os.Getenv("PREFIX"))
	if len(prefix) == 0 {
		env.Prefix = "!"
	} else {
		env.Prefix = prefix
	}

	return nil
}

func GetEnv() *Environment {
	return env
}
