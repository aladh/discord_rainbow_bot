package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DiscordToken string
	InviteUrl    string
	IntervalMs   int
}

func New() *Config {
	return &Config{
		DiscordToken: getRequiredEnvString("DISCORD_TOKEN"),
		InviteUrl:    getRequiredEnvString("INVITE_URL"),
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
