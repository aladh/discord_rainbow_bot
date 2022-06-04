package config

import (
	"log"
	"os"
)

// Config holds configuration values for the application
type Config struct {
	DiscordToken string
	InviteURL    string
}

// New creates and returns a Config using values from environment variables
func New() *Config {
	return &Config{
		DiscordToken: getRequiredEnvString("DISCORD_TOKEN"),
		InviteURL:    getRequiredEnvString("INVITE_URL"),
	}
}

func getRequiredEnvString(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	log.Fatalf("Environment variable %s must be set", key)
	return ""
}
