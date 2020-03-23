package colours

import (
	"fmt"
	"github.com/ali-l/discord_rainbow_bot/guildroles"
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"time"
)

const maxColour = 16777216

// Change rotates the colour of every active GuildRole at the specified interval
func Change(session *discordgo.Session, intervalMs int) {
	for {
		changeColours(session, *guildroles.Get(), intervalMs)
	}
}

func changeColours(session *discordgo.Session, guildRoles guildroles.GuildRoles, intervalMs int) {
	for _, guildRole := range guildRoles {
		err := changeColour(session, guildRole)
		if err != nil {
			log.Println(err)
		}

		time.Sleep(time.Duration(intervalMs) * time.Millisecond)
	}
}

func changeColour(s *discordgo.Session, guildRole *guildroles.GuildRole) error {
	colour := rand.Intn(maxColour)

	_, err := s.GuildRoleEdit(guildRole.GuildID, guildRole.ID, guildRole.Name, colour, guildRole.Hoist, guildRole.Permissions, guildRole.Mentionable)
	if err != nil {
		return fmt.Errorf("error updating role colour for role ID %s, guild ID %s: %w", guildRole.ID, guildRole.GuildID, err)
	}

	return nil
}
