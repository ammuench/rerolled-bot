package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func LogCmdError(err error, cmd *string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Error running %v command ==> %v\n", cmd, err)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Uh oh, there was an error running this command",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
