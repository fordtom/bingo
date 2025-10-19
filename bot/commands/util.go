package commands

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fordtom/bingo/db"
)

var mentionRegex = regexp.MustCompile(`<@!?(\d+)>`)

// parseMentionsToIDs extracts Discord user IDs from mention strings
func parseMentionsToIDs(s string) []int64 {
	matches := mentionRegex.FindAllStringSubmatch(s, -1)
	ids := make([]int64, 0, len(matches))
	for _, match := range matches {
		if id, err := strconv.ParseInt(match[1], 10, 64); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

// parseUserID converts Discord snowflake ID to int64
func parseUserID(snowflake string) int64 {
	id, _ := strconv.ParseInt(snowflake, 10, 64)
	return id
}

// userDisplayName returns the member's guild display name (nickname) if set; otherwise the username.
// Falls back to the userID string if neither can be retrieved.
func userDisplayName(s *discordgo.Session, guildID, userID string) string {
	if guildID == "" || userID == "" {
		return userID
	}
	if m, err := s.State.Member(guildID, userID); err == nil && m != nil {
		if m.Nick != "" {
			return m.Nick
		}
		if m.User != nil && m.User.Username != "" {
			return m.User.Username
		}
	}
	if m, err := s.GuildMember(guildID, userID); err == nil && m != nil {
		if m.Nick != "" {
			return m.Nick
		}
		if m.User != nil && m.User.Username != "" {
			return m.User.Username
		}
	}
	return userID
}

// getGameIDOrActive returns specified game_id or active game
func getGameIDOrActive(ctx context.Context, database *db.DB, options []*discordgo.ApplicationCommandInteractionDataOption, optionName string) (int64, error) {
	for _, opt := range options {
		if opt.Name == optionName {
			return opt.IntValue(), nil
		}
	}

	game, err := database.GetActiveGame(ctx)
	if err != nil {
		return 0, fmt.Errorf("error fetching active game: %w", err)
	}
	if game == nil {
		return 0, fmt.Errorf("no active game found. Please specify a game_id or set an active game")
	}
	return game.ID, nil
}

// Embed color constants
const (
	colorInfo    = 0x3498db
	colorSuccess = 0x2ecc71
	colorError   = 0xe74c3c
	colorWin     = 0xf1c40f
)

// respondEmbed sends a message as an embed with optional ephemeral flag
func respondEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, title, desc string, color int, ephemeral bool) {
	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: desc,
		Color:       color,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	flags := discordgo.MessageFlags(0)
	if ephemeral {
		flags = discordgo.MessageFlagsEphemeral
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  flags,
		},
	})
}

// respondError sends an ephemeral error message using an embed
func respondError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	data := i.ApplicationCommandData()
	sub := ""
	if len(data.Options) > 0 {
		sub = data.Options[0].Name
	}
	actor := ""
	if i.Member != nil && i.Member.User != nil {
		actor = i.Member.User.ID
	}
	log.Printf("err %s/%s actor=%s msg=%q", data.Name, sub, actor, message)
	respondEmbed(s, i, "Error", message, colorError, true)
}

// respondSuccess sends a non-ephemeral info message using an embed
func respondSuccess(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	data := i.ApplicationCommandData()
	sub := ""
	if len(data.Options) > 0 {
		sub = data.Options[0].Name
	}
	actor := ""
	if i.Member != nil && i.Member.User != nil {
		actor = i.Member.User.ID
	}
	log.Printf("ok %s/%s actor=%s msg=%q", data.Name, sub, actor, message)
	respondEmbed(s, i, "", message, colorInfo, false)
}
