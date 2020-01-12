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
const guildId = "***REMOVED***"
const roleId = "***REMOVED***"
const maxColour = 16777216
const addCommand = "+rainbow add"
const removeCommand = "+rainbow remove"

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

	rand.Seed(time.Now().Unix())

	role, err := getRole(dg, guildId, roleId)
	if err != nil {
		panic(err)
	}

	setupCommands(dg, role)

	timer := time.NewTicker(interval)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case <-timer.C:
			err := changeRoleColour(dg, guildId, role)
			if err != nil {
				fmt.Println(err)
			}
		case <-sc:
			fmt.Println("Shutting down")
			return
		}
	}
}

func setupCommands(dg *discordgo.Session, role *discordgo.Role) {
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		switch m.Content {
		case addCommand:
			err := addCommandHandler(s, m, role)
			if err != nil {
				fmt.Println(err)
			}
		case removeCommand:
			err := removeCommandHandler(s, m, role)
			if err != nil {
				fmt.Println(err)
			}
		}
	})
}

func changeRoleColour(s *discordgo.Session, guildId string, role *discordgo.Role) error {
	colour := rand.Intn(maxColour)

	_, err := s.GuildRoleEdit(guildId, role.ID, role.Name, colour, role.Hoist, role.Permissions, role.Mentionable)
	if err != nil {
		return fmt.Errorf("error updating role colour for role ID %s, guild ID %s: %w", role.ID, guildId, err)
	}

	return nil
}

func getRole(s *discordgo.Session, guildId string, roleId string) (*discordgo.Role, error) {
	roles, err := s.GuildRoles(guildId)
	if err != nil {
		return nil, fmt.Errorf("error getting roles for guild %s: %w", guildId, err)
	}

	return findRoleById(roles, roleId)
}

func findRoleById(roles []*discordgo.Role, roleId string) (*discordgo.Role, error) {
	for _, role := range roles {
		if role.ID == roleId {
			return role, nil
		}
	}

	return nil, fmt.Errorf("could not find role with ID: %s", roleId)
}

func addCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, role *discordgo.Role) error {
	err := s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, role.ID)
	if err != nil {
		return fmt.Errorf("error adding role to user %s: %w", m.Author.ID, err)
	}

	return nil
}

func removeCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, role *discordgo.Role) error {
	err := s.GuildMemberRoleRemove(m.GuildID, m.Author.ID, role.ID)
	if err != nil {
		return fmt.Errorf("error removing role from user %s: %w", m.Author.ID, err)
	}

	return nil
}
