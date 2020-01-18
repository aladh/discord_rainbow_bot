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

var guildRoles GuildRoles

func Initialize(session *discordgo.Session) error {
	err := syncGuilds(session)
	if err != nil {
		return err
	}

	session.AddHandler(onGuildCreate)
	session.AddHandler(onGuildDelete)

	return nil
}

func Run(f func(GuildRoles)) {
	f(guildRoles)
}

func FindByGuildId(guildId string) (*GuildRole, error) {
	for _, guildRole := range guildRoles {
		if guildRole.GuildId == guildId {
			return guildRole, nil
		}
	}

	return nil, fmt.Errorf("could not find role for guild %s", guildId)
}

func syncGuilds(session *discordgo.Session) error {
	guildRoles = nil

	guilds, err := session.UserGuilds(0, "", "")
	if err != nil {
		return fmt.Errorf("error getting guilds: %w", err)
	}

	for _, guild := range guilds {
		role, err := utils.FindOrCreateRole(session, guild.ID)
		if err != nil {
			return err
		}

		guildRoles = append(guildRoles, &GuildRole{GuildId: guild.ID, Role: role})
	}

	return nil
}

func onGuildCreate(session *discordgo.Session, guildCreate *discordgo.GuildCreate) {
	err := syncGuilds(session)
	if err != nil {
		panic(fmt.Sprintf("error finding/creating role for guildCreate ID %s: %s", guildCreate.ID, err))
	}
}

func onGuildDelete(session *discordgo.Session, guildDelete *discordgo.GuildDelete) {
	err := syncGuilds(session)
	if err != nil {
		panic(fmt.Sprintf("error handling guildDelete ID %s: %s", guildDelete.ID, err))
	}
}
