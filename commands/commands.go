package commands

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/ali-l/discord_rainbow_bot/guildroles"
)

const commandPrefix = "+rainbow"
const addCommand = "add"
const removeCommand = "remove"
const pingCommand = "ping"
const inviteCommand = "invite"

// Initialize registers bot commands and handlers
func Initialize(session *discordgo.Session, inviteURL string) {
	bindCommands(session, inviteURL)
	setStatus(session)
}

func setStatus(session *discordgo.Session) {
	err := session.UpdateListeningStatus(commandPrefix)
	if err != nil {
		log.Printf("error stetting status: %s", err)
	}
}

func bindCommands(session *discordgo.Session, inviteURL string) {
	session.AddHandler(func(session *discordgo.Session, messageCreate *discordgo.MessageCreate) {
		if !strings.HasPrefix(messageCreate.Content, commandPrefix) {
			return
		}

		switch extractCommand(messageCreate.Content) {
		case addCommand:
			err := addCommandHandler(session, messageCreate)
			if err != nil {
				log.Println(err)
			}
		case removeCommand:
			err := removeCommandHandler(session, messageCreate)
			if err != nil {
				log.Println(err)
			}
		case pingCommand:
			err := pingCommandHandler(session, messageCreate)
			if err != nil {
				log.Println(err)
			}
		case inviteCommand:
			err := inviteCommandHandler(session, messageCreate, inviteURL)
			if err != nil {
				log.Println(err)
			}
		default:
			err := defaultCommandHandler(session, messageCreate)
			if err != nil {
				log.Println(err)
			}
		}
	})
}

func extractCommand(message string) string {
	return strings.TrimSpace(
		strings.TrimPrefix(message, commandPrefix),
	)
}

func inviteCommandHandler(session *discordgo.Session, messageCreate *discordgo.MessageCreate, inviteURL string) error {
	_, err := session.ChannelMessageSend(messageCreate.ChannelID, "Invite me: "+inviteURL)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}

func addCommandHandler(session *discordgo.Session, messageCreate *discordgo.MessageCreate) error {
	guildRole, err := guildroles.FindByGuildID(messageCreate.GuildID)
	if err != nil {
		return err
	}

	err = session.GuildMemberRoleAdd(messageCreate.GuildID, messageCreate.Author.ID, guildRole.ID)
	if err != nil {
		return fmt.Errorf("error adding role to user %s: %w", messageCreate.Author.ID, err)
	}

	err = addCheckMarkReaction(session, messageCreate)
	if err != nil {
		return err
	}

	return nil
}

func removeCommandHandler(session *discordgo.Session, messageCreate *discordgo.MessageCreate) error {
	guildRole, err := guildroles.FindByGuildID(messageCreate.GuildID)
	if err != nil {
		return err
	}

	err = session.GuildMemberRoleRemove(messageCreate.GuildID, messageCreate.Author.ID, guildRole.ID)
	if err != nil {
		return fmt.Errorf("error removing role from user %s: %w", messageCreate.Author.ID, err)
	}

	err = addCheckMarkReaction(session, messageCreate)
	if err != nil {
		return err
	}

	return nil
}

func pingCommandHandler(session *discordgo.Session, messageCreate *discordgo.MessageCreate) error {
	message, err := session.ChannelMessageSend(messageCreate.ChannelID, "Pong!")
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	timestamp, err := message.Timestamp.Parse()
	if err != nil {
		return fmt.Errorf("error parsing timestamp: %w", err)
	}

	latency := (time.Now().UnixNano() - timestamp.UnixNano()) / 1000000

	_, err = session.ChannelMessageEdit(message.ChannelID, message.ID, fmt.Sprintf("Pong! (%dms)", latency))
	if err != nil {
		return fmt.Errorf("error editing message: %w", err)
	}

	return nil
}

func defaultCommandHandler(session *discordgo.Session, messageCreate *discordgo.MessageCreate) error {
	embed := discordgo.MessageEmbed{
		Title:  "Commands",
		Color:  0,
		Author: &discordgo.MessageEmbedAuthor{Name: "Rainbow Bot"},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Add Rainbow role to yourself", Value: fmt.Sprintf("%s add", commandPrefix), Inline: false},
			{Name: "Remove Rainbow role from yourself", Value: fmt.Sprintf("%s remove", commandPrefix), Inline: false},
			{Name: "Show Rainbow Bot invite link", Value: fmt.Sprintf("%s invite", commandPrefix), Inline: false},
		},
	}

	_, err := session.ChannelMessageSendComplex(messageCreate.ChannelID, &discordgo.MessageSend{Embed: &embed})
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}

func addCheckMarkReaction(session *discordgo.Session, messageCreate *discordgo.MessageCreate) error {
	err := session.MessageReactionAdd(messageCreate.ChannelID, messageCreate.ID, "✅")
	if err != nil {
		return fmt.Errorf("error adding check mark rection: %w", err)
	}

	return nil
}
