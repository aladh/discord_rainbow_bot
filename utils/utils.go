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

	role, err := findRoleByName(roles)
	if err != nil {
		// TODO: Create role if it does not already exist
		return nil, fmt.Errorf("error finding rainbow role for guild %s: %w", guildId, err)
	}

	return role, nil
}

func findRoleByName(roles []*discordgo.Role) (*discordgo.Role, error) {
	for _, role := range roles {
		if roleNameRegex.MatchString(role.Name) {
			return role, nil
		}
	}

	return nil, fmt.Errorf("could not find matching role")
}
