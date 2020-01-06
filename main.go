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

func main() {
	dg, err := discordgo.New("Bot ***REMOVED***")

	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection: ", err)
		return
	}
	defer dg.Close()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	timer := time.NewTicker(interval)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case <-timer.C:
			err := changeRoleColour(dg, "***REMOVED***", "***REMOVED***")
			if err != nil {
				fmt.Println("error updating role colour: ", err)
			}
		case <-sc:
			return
		}
	}
}

func changeRoleColour(s *discordgo.Session, guildId string, roleId string) error {
	roles, err := s.GuildRoles(guildId)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("error getting roles for guild %s :", guildId), err)
	}

	var role *discordgo.Role

	for _, r := range roles {
		if r.ID == roleId {
			role = r
		}
	}

	rand.Seed(time.Now().Unix())
	colour := rand.Intn(16777216)

	_, err = s.GuildRoleEdit(guildId, role.ID, role.Name, colour, role.Hoist, role.Permissions, role.Mentionable)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("error updating role colour for guild %s :", guildId), err)
	}

	return nil
}
