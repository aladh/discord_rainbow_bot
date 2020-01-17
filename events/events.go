package events

import (
	"fmt"
	"github.com/ali-l/discord_rainbow_bot/guildroles"
	"github.com/ali-l/discord_rainbow_bot/utils"
	"github.com/bwmarrin/discordgo"
)

func Setup(session *discordgo.Session, createGuildRole chan<- guildroles.GuildRole, deleteGuildRole chan<- string) {
	session.AddHandler(func(session *discordgo.Session, guildCreate *discordgo.GuildCreate) {
		role, err := utils.FindOrCreateRole(session, guildCreate.ID)
		if err != nil {
			fmt.Println("error finding/creating role for guildCreate ID: ", guildCreate.ID)
			return
		}

		guildRole := guildroles.GuildRole{GuildId: guildCreate.ID, Role: role}
		createGuildRole <- guildRole
	})

	session.AddHandler(func(session *discordgo.Session, guildDelete *discordgo.GuildDelete) {
		deleteGuildRole <- guildDelete.ID
	})
}
