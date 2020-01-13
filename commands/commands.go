package commands

import (
	"fmt"
	"github.com/ali-l/discord_rainbow_bot/guildroles"
	"github.com/bwmarrin/discordgo"
	"time"
)

const prefix = "+rainbow"
const addCommand = prefix + "add"
const removeCommand = prefix + "remove"
const pingCommand = prefix + "ping"

func Setup(s *discordgo.Session, guildRoles guildroles.GuildRoles) {
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
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

func addCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, guildRoles guildroles.GuildRoles) error {
	guildRole, err := guildRoles.FindGuildId(m.GuildID)
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

func removeCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate, guildRoles guildroles.GuildRoles) error {
	guildRole, err := guildRoles.FindGuildId(m.GuildID)
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
