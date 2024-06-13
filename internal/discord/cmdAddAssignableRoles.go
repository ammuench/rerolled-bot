package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func AddAssignableRole(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guildRoles, err := s.GuildRoles(i.GuildID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Uh oh, there was an error running this command",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	} else {
		minRoles := 0
		totalRoles := len(guildRoles) - 1 // Always remove @everyone from this count
		generatedRoleOptions := []discordgo.SelectMenuOption{}
		for _, role := range guildRoles {
			if role.Name != "@everyone" {
				generatedRoleOptions = append(generatedRoleOptions, discordgo.SelectMenuOption{
					Label: role.Name,
					Value: role.ID,
				})
			}
		}

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hey there! Congratulations, you just executed the `add-assignable-role`  command",
				Flags:   discordgo.MessageFlagsEphemeral,
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.SelectMenu{
								MenuType:    discordgo.StringSelectMenu,
								CustomID:    cmdSelectAssignableRoles,
								Placeholder: "Select roles that you want to be publically assignable",
								Options:     generatedRoleOptions,
								MinValues:   &minRoles,
								MaxValues:   totalRoles,
							},
						},
					},
				},
			},
		})

		if err != nil {
			log.Printf("Error running %v command ==> %v\n", cmdAddAssignableRole, err)
		}
	}

}

func HandleAddAssignableRoleSelect(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Thank you for selecting a role",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	data := i.MessageComponentData()

	for selectOptIdx, selectOpt := range data.Values {
		fmt.Printf("Selection #%v: \n", selectOptIdx)
		fmt.Printf(":::> %v\n", selectOpt)
	}
}
