package discord

import "github.com/bwmarrin/discordgo"

func UpdateRoles(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hey there! Congratulations, you just executed the `update-my-roles` command",
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Emoji: &discordgo.ComponentEmoji{
								Name: "🔑",
							},
							Label:    "Mythic+ Keys",
							Style:    discordgo.PrimaryButton,
							CustomID: cmdToggleMplusRole,
						},
						discordgo.Button{
							Emoji: &discordgo.ComponentEmoji{
								Name: "🏆",
							},
							Label:    "Achievements",
							Style:    discordgo.PrimaryButton,
							CustomID: cmdToggleAcheivementRole,
						},
						discordgo.Button{
							Emoji: &discordgo.ComponentEmoji{
								Name: "⚔️",
							},
							Label:    "PVP",
							Style:    discordgo.PrimaryButton,
							CustomID: cmdTogglePVPRole,
						},
					},
				},
			},
		},
	})
}

func ToggleMPlusRole(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Mythic plus role has been added to your account",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
