package config

import (
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

var (
	env Config
	log *Logger = NewLogger("config")
)

func init() {
	log.Info("Initializing config...")

	token := strings.TrimSpace(os.Getenv("TOKEN"))
	if len(token) == 0 || token == "" {
		log.Fatal("TOKEN is required")
	} else {
		env.Token = token
	}

	prefix := strings.TrimSpace(os.Getenv("PREFIX"))
	if len(prefix) == 0 {
		env.Prefix = "!"
	} else {
		env.Prefix = prefix
	}

	owner := strings.TrimSpace(os.Getenv("OWNER"))
	if len(owner) == 0 {
		log.Fatal("OWNER is required")
	} else {
		env.Owner = owner
	}

	cachePath := strings.TrimSpace(os.Getenv("CACHE_PATH"))
	if len(cachePath) == 0 {
		env.CachePath = "./audio_cache"
	} else {
		env.CachePath = cachePath
	}

	youtubeKey := strings.TrimSpace(os.Getenv("YOUTUBE_KEY"))
	if len(youtubeKey) == 0 {
		log.Fatal("YOUTUBE_KEY is required")
	} else {
		env.YoutubeKey = youtubeKey
	}
}

func Get() *Config {
	return &env
}
