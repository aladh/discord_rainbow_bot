package guildroles

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

const roleName = "Rainbow"

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

	assignRoleToSelf(session)

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
		role, err := findOrCreateRole(session, guild.ID)
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

	assignRoleToSelf(session)
}

func onGuildDelete(session *discordgo.Session, guildDelete *discordgo.GuildDelete) {
	err := syncGuilds(session)
	if err != nil {
		panic(fmt.Sprintf("error handling guildDelete ID %s: %s", guildDelete.ID, err))
	}
}

func findOrCreateRole(s *discordgo.Session, guildId string) (*discordgo.Role, error) {
	roles, err := s.GuildRoles(guildId)
	if err != nil {
		return nil, fmt.Errorf("error getting roles for guild %s: %w", guildId, err)
	}

	role := findRoleByName(roles)

	if role == nil {
		role, err = createRole(s, guildId)
		if err != nil {
			return nil, fmt.Errorf("error creating role for guild %s: %w", guildId, err)
		}
	}

	return role, nil
}

func findRoleByName(roles []*discordgo.Role) *discordgo.Role {
	for _, role := range roles {
		if role.Name == roleName {
			return role
		}
	}

	return nil
}

func createRole(session *discordgo.Session, guildId string) (*discordgo.Role, error) {
	role, err := session.GuildRoleCreate(guildId)
	if err != nil {
		return nil, fmt.Errorf("error creating role for guild %s: %w", guildId, err)
	}

	role, err = session.GuildRoleEdit(guildId, role.ID, roleName, role.Color, role.Hoist, role.Permissions, role.Mentionable)
	if err != nil {
		return role, fmt.Errorf("error updating name for guild ID %s: %w", guildId, err)
	}

	return role, nil
}

func assignRoleToSelf(session *discordgo.Session) {
	userID := session.State.User.ID

	for _, guildRole := range guildRoles {
		err := session.GuildMemberRoleAdd(guildRole.GuildId, userID, guildRole.ID)
		if err != nil {
			fmt.Println("error adding role to user ", userID, ": ", err)
		}
	}
}
