package colours

import (
	"fmt"
	"github.com/ali-l/discord_rainbow_bot/guildroles"
	"github.com/bwmarrin/discordgo"
	"math/rand"
)

const maxColour = 16777216

func Change(s *discordgo.Session, guildRoles guildroles.GuildRoles) error {
	for _, guildRole := range guildRoles {
		colour := rand.Intn(maxColour)

		_, err := s.GuildRoleEdit(guildRole.GuildId, guildRole.ID, guildRole.Name, colour, guildRole.Hoist, guildRole.Permissions, guildRole.Mentionable)
		if err != nil {
			return fmt.Errorf("error updating role colour for role ID %s, guild ID %s: %w", guildRole.ID, guildRole.GuildId, err)
		}
	}

	return nil
}
