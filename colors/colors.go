package colors

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/ali-l/discord_rainbow_bot/guildroles"
)

const maxColor = 16777216

// Rotate continuously rotates the color of every active GuildRole
func Rotate(session *discordgo.Session, delayMs int) {
	for {
		guildroles.ForEach(func(guildRole *guildroles.GuildRole) {
			err := changeColor(session, guildRole)
			if err != nil {
				log.Println(err)
			}

			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		})
	}
}

func changeColor(session *discordgo.Session, guildRole *guildroles.GuildRole) error {
	color := rand.Intn(maxColor)

	_, err := session.GuildRoleEdit(guildRole.GuildID, guildRole.ID, guildRole.Name, color, guildRole.Hoist, guildRole.Permissions, guildRole.Mentionable)
	if err != nil {
		return fmt.Errorf("error updating role color for role ID %s, guild ID %s: %w", guildRole.ID, guildRole.GuildID, err)
	}

	return nil
}
