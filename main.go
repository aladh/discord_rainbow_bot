package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const discordToken = "***REMOVED***"

const interval = 5 * time.Second
const maxColour = 16777216

const addCommand = "+rainbow add"
const removeCommand = "+rainbow remove"
const pingCommand = "+rainbow ping"
const rainbowRoleName = "Rainbow"

type GuildRole struct {
	GuildId string
	*discordgo.Role
}

func main() {
	dg, err := discordgo.New(fmt.Sprintf("Bot %s", discordToken))
	if err != nil {
		panic(fmt.Errorf("error creating Discord session: %w", err))
	}

	err = dg.Open()
	if err != nil {
		panic(fmt.Errorf("error opening connection: %w", err))
	}
	defer dg.Close()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	guilds, err := dg.UserGuilds(0, "", "")
	if err != nil {
		panic(fmt.Errorf("error getting user guilds: %w", err))
	}

	guildRoles, err := findOrCreateRainbowRoles(dg, guilds)
	if err != nil {
		panic(fmt.Errorf("error finding/creating rainbow roles: %w", err))
	}

	rand.Seed(time.Now().Unix())

	setupCommands(dg, guildRoles)

	timer := time.NewTicker(interval)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case <-timer.C:
			err := changeRoleColour(dg, guildRoles[0])
			if err != nil {
				fmt.Println(err)
			}
		case <-sc:
			fmt.Println("Shutting down")
			return
		}
	}
}

func findOrCreateRainbowRoles(s *discordgo.Session, guilds []*discordgo.UserGuild) ([]*GuildRole, error) {
	var roles []*GuildRole

	for _, guild := range guilds {
		role, err := findOrCreateRainbowRole(s, guild.ID)
		if err != nil {
			return nil, err
		}

		roles = append(roles, &GuildRole{GuildId: guild.ID, Role: role})
	}

	return roles, nil
}

func setupCommands(dg *discordgo.Session, guildRoles []*GuildRole) {
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		switch m.Content {
		case addCommand:
			err := addCommandHandler(s, m, guildRoles)
			if err != nil {
				fmt.Println(err)
			}
		case removeCommand:
			err := removeCommandHandler(s, m, guildRoles)
			if err != nil {
				fmt.Println(err)
			}
		case pingCommand:
			err := pingCommandHandler(s, m)
			if err != nil {
				fmt.Println(err)
			}
		}
	})
}

func changeRoleColour(s *discordgo.Session, guildRole *GuildRole) error {
	colour := rand.Intn(maxColour)

	_, err := s.GuildRoleEdit(guildRole.GuildId, guildRole.ID, guildRole.Name, colour, guildRole.Hoist, guildRole.Permissions, guildRole.Mentionable)
	if err != nil {
		return fmt.Errorf("error updating role colour for role ID %s, guild ID %s: %w", guildRole.ID, guildRole.GuildId, err)
	}

	return nil
}

func findRoleByName(roles []*discordgo.Role, roleName string) (*discordgo.Role, error) {
	for _, role := range roles {
		if role.Name == roleName {
			return role, nil
		}
	}

	return nil, fmt.Errorf("could not find role with name: %s", roleName)
}

func findOrCreateRainbowRole(s *discordgo.Session, guildId string) (*discordgo.Role, error) {
	roles, err := s.GuildRoles(guildId)
	if err != nil {
		return nil, fmt.Errorf("error getting roles for guild %s: %w", guildId, err)
	}

	role, err := findRoleByName(roles, rainbowRoleName)
	if err != nil {
		return nil, fmt.Errorf("error finding rainbow role for guild %s: %w", guildId, err)
	}

	return role, nil
}

func findGuildRoleByGuildId(guildRoles []*GuildRole, guildId string) (*GuildRole, error) {
	for _, guildRole := range guildRoles {
		if guildRole.GuildId == guildId {
			return guildRole, nil
		}
	}

	return nil, fmt.Errorf("could not find role for guild %s", guildId)
}

func addCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, guildRoles []*GuildRole) error {
	guildRole, err := findGuildRoleByGuildId(guildRoles, m.GuildID)
	if err != nil {
		return err
	}

	err = s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, guildRole.ID)
	if err != nil {
		return fmt.Errorf("error adding role to user %s: %w", m.Author.ID, err)
	}

	err = addCheckMarkReaction(s, m)
	if err != nil {
		return err
	}

	return nil
}

func removeCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, guildRoles []*GuildRole) error {
	guildRole, err := findGuildRoleByGuildId(guildRoles, m.GuildID)
	if err != nil {
		return err
	}

	err = s.GuildMemberRoleRemove(m.GuildID, m.Author.ID, guildRole.ID)
	if err != nil {
		return fmt.Errorf("error removing role from user %s: %w", m.Author.ID, err)
	}

	err = addCheckMarkReaction(s, m)
	if err != nil {
		return err
	}

	return nil
}

func pingCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) error {
	message, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	timestamp, err := message.Timestamp.Parse()
	if err != nil {
		return fmt.Errorf("error parsing timestamp: %w", err)
	}

	latency := (time.Now().UnixNano() - timestamp.UnixNano()) / 1000000

	_, err = s.ChannelMessageEdit(message.ChannelID, message.ID, fmt.Sprintf("Pong! (%dms)", latency))
	if err != nil {
		return fmt.Errorf("error editing message: %w", err)
	}

	return nil
}

func addCheckMarkReaction(s *discordgo.Session, m *discordgo.MessageCreate) error {
	err := s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	if err != nil {
		return fmt.Errorf("error adding check mark rection: %w", err)
	}

	return nil
}
