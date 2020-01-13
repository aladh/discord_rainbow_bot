package colours

import (
	"fmt"
	"github.com/ali-l/discord_rainbow_bot/guildroles"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"time"
)

const maxColour = 16777216
const intervalMs = 8000

var numGuilds int

func Change(session *discordgo.Session, guildRoles guildroles.GuildRoles) {
	numGuilds = len(guildRoles)

	for {
		changeColours(session, guildRoles)
	}
}

func changeColours(session *discordgo.Session, guildRoles guildroles.GuildRoles) {
	for _, guildRole := range guildRoles {
		err := changeColour(session, guildRole)
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(time.Duration(intervalMs/numGuilds) * time.Millisecond)
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
