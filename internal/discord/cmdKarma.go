package discord

import (
	// "github.com/ammuench/rerolled-bot/internal/db"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func AddKarma(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// database := db.GetDB()
	options := i.ApplicationCommandData().Options
	optionUser := options[0].UserValue(s)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%v gave %v +1", FmtMentionString(i.Member.User.ID), FmtMentionString(optionUser.ID)),
		},
	})
	if err != nil {
		LogCmdError(err, &cmdAddKarma, s, i)
	}
}

func RemoveKarma(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// database := db.GetDB()
	options := i.ApplicationCommandData().Options
	optionUser := options[0].UserValue(s)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%v gave %v -1", FmtMentionString(i.Member.User.ID), FmtMentionString(optionUser.ID)),
		},
	})
	if err != nil {
		LogCmdError(err, &cmdRemoveKarma, s, i)
	}
}
