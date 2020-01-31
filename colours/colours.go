package colours

import (
	"fmt"
	"github.com/ali-l/discord_rainbow_bot/guildroles"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"time"
)

const maxColour = 16777216

func Change(session *discordgo.Session, intervalMs int) {
	var changeColoursWithSession = changeColours(session, intervalMs)

	for {
		guildroles.Run(changeColoursWithSession)
	}
}

func changeColours(session *discordgo.Session, intervalMs int) func(guildroles.GuildRoles) {
	return func(guildRoles guildroles.GuildRoles) {
		for _, guildRole := range guildRoles {
			err := changeColour(session, guildRole)
			if err != nil {
				fmt.Println(err)
			}

			time.Sleep(time.Duration(intervalMs) * time.Millisecond)
		}
	}
}

func changeColour(s *discordgo.Session, guildRole *guildroles.GuildRole) error {
	colour := rand.Intn(maxColour)

	_, err := s.GuildRoleEdit(guildRole.GuildId, guildRole.ID, guildRole.Name, colour, guildRole.Hoist, guildRole.Permissions, guildRole.Mentionable)
	if err != nil {
		return fmt.Errorf("error updating role colour for role ID %s, guild ID %s: %w", guildRole.ID, guildRole.GuildId, err)
	}

	return nil
}
