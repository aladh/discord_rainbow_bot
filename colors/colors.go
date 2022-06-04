package colors

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/aladh/discord_rainbow_bot/guildroles"
)

const maxColor = 16777216
const delay = 2

// Rotate continuously rotates the color of every active GuildRole
func Rotate(session *discordgo.Session) {
	for {
		guildroles.ForEach(func(guildRole *guildroles.GuildRole) {
			err := changeColor(session, guildRole)
			if err != nil {
				log.Println(err)
			}

			time.Sleep(delay * time.Minute)
		})
	}
}

func changeColor(session *discordgo.Session, guildRole *guildroles.GuildRole) error {
	log.Printf("changing color for role ID %s, guild ID %s\n", guildRole.ID, guildRole.GuildID)

	color := rand.Intn(maxColor)

	_, err := session.GuildRoleEdit(guildRole.GuildID, guildRole.ID, guildRole.Name, color, guildRole.Hoist, guildRole.Permissions, guildRole.Mentionable)
	if err != nil {
		return fmt.Errorf("error updating role color for role ID %s, guild ID %s: %w", guildRole.ID, guildRole.GuildID, err)
	}

	log.Printf("changed color for role ID %s, guild ID %s to %d\n", guildRole.ID, guildRole.GuildID, color)

	return nil
}
