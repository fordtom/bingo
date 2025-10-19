package commands

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"

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

// respondError sends an ephemeral error message
func respondError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	data := i.ApplicationCommandData()
	sub := ""
	if len(data.Options) > 0 {
		sub = data.Options[0].Name
	}
	log.Printf("err %s/%s actor=%s msg=%q", data.Name, sub, i.Member.User.ID, message)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "❌ " + message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

// respondSuccess sends a success message
func respondSuccess(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	data := i.ApplicationCommandData()
	sub := ""
	if len(data.Options) > 0 {
		sub = data.Options[0].Name
	}
	log.Printf("ok %s/%s actor=%s msg=%q", data.Name, sub, i.Member.User.ID, message)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}
