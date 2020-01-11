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

const interval = 5 * time.Second
const guildId = "***REMOVED***"
const roleId = "***REMOVED***"
const maxColour = 16777216

func main() {
	dg, err := discordgo.New("Bot ***REMOVED***")
	if err != nil {
		panic(fmt.Errorf("error creating Discord session: %w", err))
	}

	err = dg.Open()
	if err != nil {
		panic(fmt.Errorf("error opening connection: %w", err))
	}
	defer dg.Close()

	rand.Seed(time.Now().Unix())

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	timer := time.NewTicker(interval)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	role, err := getRole(dg, guildId, roleId)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-timer.C:
			err := changeRoleColour(dg, guildId, role)
			if err != nil {
				fmt.Println("error updating role colour: ", err)
			}
		case <-sc:
			fmt.Println("Shutting down")
			return
		}
	}
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
	var role *discordgo.Role

	for _, r := range roles {
		if r.ID == roleId {
			role = r
			break
		}

		return nil, fmt.Errorf("could not find role with ID: %s", roleId)
	}

	return role, nil
}
