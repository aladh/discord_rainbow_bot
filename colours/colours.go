package colours

import (
	"fmt"
	"github.com/ali-l/discord_rainbow_bot/guildroles"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"time"
)

const maxColour = 16777216
const interval = 5

var timer = time.NewTicker(interval * time.Second)

func Change(session *discordgo.Session, guildRoles guildroles.GuildRoles) {
	for {
		<-timer.C

		changeColours(session, guildRoles)
	}
}

func changeColours(s *discordgo.Session, guildRoles guildroles.GuildRoles) {
	for _, guildRole := range guildRoles {
		err := changeColour(s, guildRole)
		if err != nil {
			fmt.Println(err)
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
