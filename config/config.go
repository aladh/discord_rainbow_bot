package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds configuration values for the application
type Config struct {
	DiscordToken string
	InviteURL    string
	IntervalMs   int
}

// New creates and returns a Config using values from environment variables
func New() *Config {
	return &Config{
		DiscordToken: getRequiredEnvString("DISCORD_TOKEN"),
		InviteURL:    getRequiredEnvString("INVITE_URL"),
		IntervalMs:   getRequiredEnvInt("INTERVAL_MS"),
	}
}

func getRequiredEnvString(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	panic(fmt.Sprintf("Environment variable %s must be set", key))
}

func getRequiredEnvInt(key string) int {
	valueStr := getRequiredEnvString(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	panic(fmt.Sprintf("Environment variable %s must be an integer", key))
}
