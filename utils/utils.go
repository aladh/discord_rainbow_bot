package utils

import (
	"fmt"
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

func findRoleByName(roles []*discordgo.Role) *discordgo.Role {
	for _, role := range roles {
		// TODO: Remove random color exception
		if role.Name == roleName || role.Name == "Random Color" {
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
