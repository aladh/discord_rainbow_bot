package colours

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ali-l/discord_rainbow_bot/guildroles"
	"github.com/bwmarrin/discordgo"
)

const maxColour = 16777216

// Rotate continuously rotates the colour of every active GuildRole
func Rotate(session *discordgo.Session, delayMs int) {
	for {
		guildroles.ForEach(func(guildRole *guildroles.GuildRole) {
			err := changeColour(session, guildRole)
			if err != nil {
				log.Println(err)
			}

			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		})
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
