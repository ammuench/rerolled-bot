package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var adminPermission int64 = discordgo.PermissionAdministrator
var adminAllowed = true

var discordCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "update-my-roles",
		Description: "Command to update your opt-in roles in the server",
	},
	{
		Name:                     "add-assignable-role",
		Description:              "Adds a role to the public assignable roles list",
		DefaultMemberPermissions: &adminPermission,
	},
}
var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"update-my-roles": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hey there! Congratulations, you just executed the `update-my-roles` command",
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
	},
	"add-assignable-role": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hey there! Congratulations, you just executed the `add-assignable-role`  command",
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
	},
}

func InitializeCommands(discordBot *discordgo.Session) ([]*discordgo.ApplicationCommand, error) {
	discordBot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	registeredCommands := make([]*discordgo.ApplicationCommand, 2)
	for dCmdIdx, dCmd := range discordCommands {
		successfulSetCmd, err := discordBot.ApplicationCommandCreate(discordBot.State.User.ID, "1246302013860483142", dCmd)
		if err != nil {
			return nil, fmt.Errorf("cannot create '%v' command: %v", discordCommands[0].Name, err)
		}

		registeredCommands[dCmdIdx] = successfulSetCmd
	}
	return registeredCommands, nil
}

func ShutdownCommands(discordBot *discordgo.Session, registeredCmds []*discordgo.ApplicationCommand) {
	for _, rCmd := range registeredCmds {
		err := discordBot.ApplicationCommandDelete(discordBot.State.User.ID, "1246302013860483142", rCmd.ID)
		if err != nil {
			fmt.Printf("Cannot delete '%v' command: %v\n", rCmd.Name, err)
		}
	}
}
