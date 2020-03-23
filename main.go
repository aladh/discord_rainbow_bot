package main

import (
	"fmt"
	"github.com/ali-l/discord_rainbow_bot/colours"
	"github.com/ali-l/discord_rainbow_bot/commands"
	"github.com/ali-l/discord_rainbow_bot/config"
	"github.com/ali-l/discord_rainbow_bot/guildroles"
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var session *discordgo.Session
var conf = config.New()

func init() {
	rand.Seed(time.Now().Unix())

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

	go colours.Change(session, conf.DelayMs)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-sc
}

func initSession() (*discordgo.Session, error) {
	session, err := discordgo.New(fmt.Sprintf("Bot %s", conf.DiscordToken))
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}

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
