package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	dg, err := discordgo.New("Bot ***REMOVED***")

	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == "***REMOVED***" {
		fmt.Println(fmt.Sprintf("reacting to %s", m.Content))

		err := s.MessageReactionAdd(m.ChannelID, m.ID, "***REMOVED***")
		if err != nil {
			fmt.Println("error reacting to message", err)
		}
	}
}
