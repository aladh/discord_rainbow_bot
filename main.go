package main

import (
	"fmt"
	"github.com/ali-l/discord_rainbow_bot/colours"
	"github.com/ali-l/discord_rainbow_bot/commands"
	"github.com/ali-l/discord_rainbow_bot/guildroles"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const discordToken = "***REMOVED***"

var session *discordgo.Session
var guildRoles guildroles.GuildRoles

func init() {
	var err error

	session, err = discordgo.New(fmt.Sprintf("Bot %s", discordToken))
	if err != nil {
		panic(fmt.Errorf("error creating Discord session: %w", err))
	}

	err = session.Open()
	if err != nil {
		panic(fmt.Errorf("error opening connection: %w", err))
	}

	rand.Seed(time.Now().Unix())

	guildRoles, err = guildroles.New(session)
	if err != nil {
		panic(fmt.Errorf("error finding/creating rainbow roles: %w", err))
	}

	commands.Setup(session, guildRoles)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
}

func main() {
	defer session.Close()

	go colours.Change(session, guildRoles)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<-sc
}
