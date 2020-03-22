package guildroles

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

const roleName = "Rainbow"

// GuildRole holds a guild and its corresponding rainbow role
type GuildRole struct {
	GuildID string
	*discordgo.Role
}

// GuildRoles is a slice of GuildRole
type GuildRoles []*GuildRole

var guildRoles GuildRoles

// Initialize loads active guilds and registers event handlers to keep them in sync
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

// Run invokes the given function with the current GuildRoles
func Run(f func(GuildRoles)) {
	f(guildRoles)
}

// FindByGuildID returns the GuildRole for the given guildID
func FindByGuildID(guildID string) (*GuildRole, error) {
	for _, guildRole := range guildRoles {
		if guildRole.GuildID == guildID {
			return guildRole, nil
		}
	}

	return nil, fmt.Errorf("could not find role for guild %s", guildID)
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

		guildRoles = append(guildRoles, &GuildRole{GuildID: guild.ID, Role: role})
	}

	return nil
}

func onGuildCreate(session *discordgo.Session, guildCreate *discordgo.GuildCreate) {
	err := syncGuilds(session)
	if err != nil {
		log.Panicf("error finding/creating role for guildCreate ID %s: %s", guildCreate.ID, err)
	}

	assignRoleToSelf(session)
}

func onGuildDelete(session *discordgo.Session, guildDelete *discordgo.GuildDelete) {
	err := syncGuilds(session)
	if err != nil {
		log.Panicf("error handling guildDelete ID %s: %s", guildDelete.ID, err)
	}
}

func findOrCreateRole(s *discordgo.Session, guildID string) (*discordgo.Role, error) {
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return nil, fmt.Errorf("error getting roles for guild %s: %w", guildID, err)
	}

	role := findRoleByName(roles)

	if role == nil {
		role, err = createRole(s, guildID)
		if err != nil {
			return nil, fmt.Errorf("error creating role for guild %s: %w", guildID, err)
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

func createRole(session *discordgo.Session, guildID string) (*discordgo.Role, error) {
	role, err := session.GuildRoleCreate(guildID)
	if err != nil {
		return nil, fmt.Errorf("error creating role for guild %s: %w", guildID, err)
	}

	role, err = session.GuildRoleEdit(guildID, role.ID, roleName, role.Color, role.Hoist, role.Permissions, role.Mentionable)
	if err != nil {
		return role, fmt.Errorf("error updating name for guild ID %s: %w", guildID, err)
	}

	return role, nil
}

func assignRoleToSelf(session *discordgo.Session) {
	userID := session.State.User.ID

	for _, guildRole := range guildRoles {
		err := session.GuildMemberRoleAdd(guildRole.GuildID, userID, guildRole.ID)
		if err != nil {
			log.Printf("error adding role %s to self in guild %s: %s", guildRole.ID, guildRole.GuildID, err)
		}
	}
}
