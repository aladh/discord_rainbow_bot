package utils

import (
	"fmt"
	"github.com/ali-l/discord_rainbow_bot/guildroles"
	"github.com/bwmarrin/discordgo"
)

func AssignRoleToSelf(session *discordgo.Session) {
	guildroles.Run(func(guildRoles guildroles.GuildRoles) {
		userID := session.State.User.ID

		for _, guildRole := range guildRoles {
			err := session.GuildMemberRoleAdd(guildRole.GuildId, userID, guildRole.ID)
			if err != nil {
				fmt.Println("error adding role to user ", userID, ": ", err)
			}
		}
	})
}
