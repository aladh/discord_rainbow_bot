package main

import (
	"fmt"
	"github.com/ali-l/discord_rainbow_bot/colours"
	"github.com/ali-l/discord_rainbow_bot/commands"
	"github.com/ali-l/discord_rainbow_bot/guildroles"
	"github.com/ali-l/discord_rainbow_bot/utils"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var session *discordgo.Session

func init() {
	var err error

	session, err = discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("DISCORD_TOKEN")))
	if err != nil {
		panic(fmt.Errorf("error creating Discord session: %w", err))
	}

	err = session.Open()
	if err != nil {
		panic(fmt.Errorf("error opening connection: %w", err))
	}

	rand.Seed(time.Now().Unix())

	err = guildroles.Initialize(session)
	if err != nil {
		panic(fmt.Errorf("error initializing guild roles: %w", err))
	}

	commands.Initialize(session)

	utils.AssignRoleToSelf(session)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
}

func main() {
	defer session.Close()

	go colours.Change(session)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<-sc
}
