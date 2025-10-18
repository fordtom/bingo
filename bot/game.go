package bot

import (
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handleCreateCommand(s *discordgo.Session, i *discordgo.InteractionCreate, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	gameName := opts[0].StringValue()
	// Later: b.db.CreateGame(gameName)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Creating game: " + gameName,
		},
	})
}
