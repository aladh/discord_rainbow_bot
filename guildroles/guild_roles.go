package guildroles

import (
	"fmt"
	"github.com/ali-l/discord_rainbow_bot/utils"
	"github.com/bwmarrin/discordgo"
)

type GuildRole struct {
	GuildId string
	*discordgo.Role
}

type GuildRoles []*GuildRole

func New(s *discordgo.Session, guilds []*discordgo.UserGuild) (GuildRoles, error) {
	var guildRoles []*GuildRole

	for _, guild := range guilds {
		guildRole, err := utils.FindOrCreateRole(s, guild.ID)
		if err != nil {
			return nil, err
		}

		guildRoles = append(guildRoles, &GuildRole{GuildId: guild.ID, Role: guildRole})
	}

	return guildRoles, nil
}

func (guildRoles GuildRoles) FindGuildId(guildId string) (*GuildRole, error) {
	for _, guildRole := range guildRoles {
		if guildRole.GuildId == guildId {
			return guildRole, nil
		}
	}

	return nil, fmt.Errorf("could not find role for guild %s", guildId)
}
