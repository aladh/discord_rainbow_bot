package utils

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"regexp"
)

var roleNameRegex = regexp.MustCompile("Rainbow|Random")

func FindOrCreateRole(s *discordgo.Session, guildId string) (*discordgo.Role, error) {
	roles, err := s.GuildRoles(guildId)
	if err != nil {
		return nil, fmt.Errorf("error getting roles for guild %s: %w", guildId, err)
	}

	role := findRoleByName(roles)
	if role == nil {
		// TODO: Create role if it does not already exist
		return nil, fmt.Errorf("error finding rainbow role for guild %s", guildId)
	}

	return role, nil
}

func findRoleByName(roles []*discordgo.Role) *discordgo.Role {
	for _, role := range roles {
		if roleNameRegex.MatchString(role.Name) {
			return role
		}
	}

	return nil
}
