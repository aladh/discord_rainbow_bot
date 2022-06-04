package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/aladh/discord_rainbow_bot/colors"
	"github.com/aladh/discord_rainbow_bot/commands"
	"github.com/aladh/discord_rainbow_bot/config"
	"github.com/aladh/discord_rainbow_bot/guildroles"
)

var session *discordgo.Session
var conf = config.New()

func init() {
	rand.Seed(time.Now().UnixNano())

	var err error

	session, err = initSession()
	if err != nil {
		log.Fatal(err)
	}

	err = guildroles.Initialize(session)
	if err != nil {
		log.Fatalf("error initializing guild roles: %s", err)
	}

	commands.Initialize(session, conf.InviteURL)

	log.Println("Bot is now running. Press CTRL-C to exit.")
}

func main() {
	defer closeSession()

	go colors.Rotate(session)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-sc
}

func initSession() (*discordgo.Session, error) {
	session, err := discordgo.New(fmt.Sprintf("Bot %s", conf.DiscordToken))
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}

	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages
	session.Debug = true

	if err = session.Open(); err != nil {
		return nil, fmt.Errorf("error opening connection: %w", err)
	}

	return session, nil
}

func closeSession() {
	err := session.Close()
	if err != nil {
		log.Printf("error closing session: %s", err)
	}
}
