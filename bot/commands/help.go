package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fordtom/bingo/db"
)

// Help returns the help subcommand definition
func Help() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "help",
		Description: "Display help information about all bot commands",
	}
}

// HandleHelp processes the help command
func HandleHelp(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption, database *db.DB) {
	prefix := "/" + Prefix
	helpText := "**Game Management**\n" +
		"• `" + prefix + " new_game` - Create a game with events and player boards (requires CSV)\n" +
		"• `" + prefix + " delete_game <game_id>` - Delete a game and all data\n" +
		"• `" + prefix + " set_active_game <game_id>` - Set the active game\n\n" +
		"**Game Information**\n" +
		"• `" + prefix + " list_games` - List all games with stats\n" +
		"• `" + prefix + " list_events [game_id]` - List events with vote counts\n" +
		"• `" + prefix + " view_board <user> [game_id]` - View a player's board\n\n" +
		"**Gameplay**\n" +
		"• `" + prefix + " vote <event_id> [game_id]` - Vote that an event occurred\n" +
		"• `" + prefix + " help` - Show this help\n\n" +
		"**CSV Format**\n" +
		"One event per line, optional header row:\n" +
		"```\ndescription\nFirst event\nSecond event\n```\n\n" +
		"**Voting**\n" +
		"• Consensus: 100% for ≤3 players, 60% for larger games\n" +
		"• When consensus reached, event closes and winners are checked"

	respondEmbed(s, i, "BingoBot Commands", helpText, colorInfo, false)
}
