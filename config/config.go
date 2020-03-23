package config

import (
	"log"
	"os"
	"strconv"
)

// Config holds configuration values for the application
type Config struct {
	DiscordToken string
	InviteURL    string
	DelayMs      int
}

// New creates and returns a Config using values from environment variables
func New() *Config {
	return &Config{
		DiscordToken: getRequiredEnvString("DISCORD_TOKEN"),
		InviteURL:    getRequiredEnvString("INVITE_URL"),
		DelayMs:      getRequiredEnvInt("DELAY_MS"),
	}
}

func getRequiredEnvString(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	log.Fatalf("Environment variable %s must be set", key)
	return ""
}

func getRequiredEnvInt(key string) int {
	valueStr := getRequiredEnvString(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	log.Fatalf("Environment variable %s must be an integer", key)
	return 0
}
