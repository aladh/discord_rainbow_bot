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

// Change rotates the colour of every active GuildRole
func Change(session *discordgo.Session, delayMs int) {
	for {
		changeColours(session, *guildroles.Get(), delayMs)
	}
}

func changeColours(session *discordgo.Session, guildRoles guildroles.GuildRoles, delayMs int) {
	for _, guildRole := range guildRoles {
		err := changeColour(session, guildRole)
		if err != nil {
			log.Println(err)
		}

		time.Sleep(time.Duration(delayMs) * time.Millisecond)
	}
}

func changeColour(session *discordgo.Session, guildRole *guildroles.GuildRole) error {
	colour := rand.Intn(maxColour)

	_, err := session.GuildRoleEdit(guildRole.GuildID, guildRole.ID, guildRole.Name, colour, guildRole.Hoist, guildRole.Permissions, guildRole.Mentionable)
	if err != nil {
		return fmt.Errorf("error updating role colour for role ID %s, guild ID %s: %w", guildRole.ID, guildRole.GuildID, err)
	}

	return nil
}
