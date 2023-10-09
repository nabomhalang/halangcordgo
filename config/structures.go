package config

type Config struct {
	Token        string `cig:"token" validate:"required"`
	Owner        string `cig:"owner" validate:"required"`
	ClientID     string `cig:"client_id"`
	ClientSecret string `cig:"client_secret"`
	YoutubeKey   string `cig:"youtube_key" validate:"required"`
	Prefix       string `cig:"prefix" validate:"required"`
	CachePath    string `cig:"cache_path" validate:"required"`
}
