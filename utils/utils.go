package utils

import (
	"fmt"
	"github.com/ali-l/discord_rainbow_bot/guildroles"
	"github.com/bwmarrin/discordgo"
)

const roleName = "Rainbow"

func FindOrCreateRole(s *discordgo.Session, guildId string) (*discordgo.Role, error) {
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
