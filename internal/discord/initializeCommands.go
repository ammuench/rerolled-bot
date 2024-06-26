package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	adminPermission int64 = discordgo.PermissionAdministrator
	adminAllowed          = true
)

var (
	cmdUpdateRole         = "update-my-roles"
	cmdProcessUpdateRoles = "process-updated-roles"
	cmdAddKarma           = "plus-one"
	cmdRemoveKarma        = "minus-one"
	cmdMplusAffixes       = "current-mplus-affixes"
)

var discordCommands = []*discordgo.ApplicationCommand{
	{
		Name:        cmdUpdateRole,
		Description: "Command to update your opt-in roles in the server",
	},
	{
		Name:        cmdMplusAffixes,
		Description: "Get the current mplus affixes and print them in the channel",
	},
	{
		Name:        cmdAddKarma,
		Description: "Give someone +1",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "User to +1",
				Required:    true,
			},
		},
	},
	{
		Name:        cmdRemoveKarma,
		Description: "Give someone -1",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "User to -1",
				Required:    true,
			},
		},
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	cmdUpdateRole:   UpdateRoles,
	cmdAddKarma:     AddKarma,
	cmdRemoveKarma:  RemoveKarma,
	cmdMplusAffixes: GetMPlusAffixes,
}

var interactionComponentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	cmdProcessUpdateRoles: ProcessUpdateRoles,
}

func InitializeCommands(discordBot *discordgo.Session) ([]*discordgo.ApplicationCommand, error) {
	// Register command handlers
	discordBot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := interactionComponentHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	registeredCommands := make([]*discordgo.ApplicationCommand, len(discordCommands))
	for dCmdIdx, dCmd := range discordCommands {
		successfulSetCmd, err := discordBot.ApplicationCommandCreate(discordBot.State.User.ID, "1246302013860483142", dCmd)
		if err != nil {
			return nil, fmt.Errorf("cannot create '%v' command: %v", discordCommands[0].Name, err)
		}

		registeredCommands[dCmdIdx] = successfulSetCmd
	}

	return registeredCommands, nil
}

func TeardownAllCommands(discordBot *discordgo.Session, registeredCmds []*discordgo.ApplicationCommand) {
	for _, rCmd := range registeredCmds {
		err := discordBot.ApplicationCommandDelete(discordBot.State.User.ID, "1246302013860483142", rCmd.ID)
		if err != nil {
			fmt.Printf("Cannot delete '%v' command: %v\n", rCmd.Name, err)
		}
	}
}
