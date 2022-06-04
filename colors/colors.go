package colors

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/aladh/discord_rainbow_bot/guildroles"
)

const maxColor = 16777216

// Rotate continuously rotates the color of every active GuildRole
func Rotate(session *discordgo.Session, delayMs int) {
	for {
		guildroles.ForEach(func(guildRole *guildroles.GuildRole) {
			err := changeColor(session, guildRole)
			if err != nil {
				log.Println(err)
			}

			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		})
	}
}

func changeColor(session *discordgo.Session, guildRole *guildroles.GuildRole) error {
	color := rand.Intn(maxColor)

	_, err := guildRoleEdit(session, guildRole.GuildID, guildRole.ID, guildRole.Name, color, guildRole.Hoist, guildRole.Permissions, guildRole.Mentionable)
	if err != nil {
		return fmt.Errorf("error updating role color for role ID %s, guild ID %s: %w", guildRole.ID, guildRole.GuildID, err)
	}

	log.Printf("changed color for role ID %s, guild ID %s to %d\n", guildRole.ID, guildRole.GuildID, color)

	return nil
}

func guildRoleEdit(s *discordgo.Session, guildID, roleID, name string, color int, hoist bool, perm int64, mention bool) (st *discordgo.Role, err error) {
	// copied from discordgo/restapi.go
	// modified to return error when request fails

	// Prevent sending a color int that is too big.
	if color > 0xFFFFFF {
		err = fmt.Errorf("color value cannot be larger than 0xFFFFFF")
		return nil, err
	}

	data := struct {
		Name        string `json:"name"`               // The role's name (overwrites existing)
		Color       int    `json:"color"`              // The color the role should have (as a decimal, not hex)
		Hoist       bool   `json:"hoist"`              // Whether to display the role's users separately
		Permissions int64  `json:"permissions,string"` // The overall permissions number of the role (overwrites existing)
		Mentionable bool   `json:"mentionable"`        // Whether this role is mentionable
	}{name, color, hoist, perm, mention}

	body, err := s.RequestWithBucketID("PATCH", discordgo.EndpointGuildRole(guildID, roleID), data, discordgo.EndpointGuildRole(guildID, ""))
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	err = unmarshal(body, &st)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return
}

func unmarshal(data []byte, v interface{}) error {
	// copied from discordgo/restapi.go
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}
