package config

import (
	"fmt"
	"os"
)

type Config struct {
	DiscordToken string
	InviteUrl    string
}

func New() *Config {
	return &Config{
		DiscordToken: getRequiredEnv("DISCORD_TOKEN"),
		InviteUrl:    getRequiredEnv("INVITE_URL"),
	}
}

func getRequiredEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	panic(fmt.Sprintf("Environment variable %s must be set", key))
}
