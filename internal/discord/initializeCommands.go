package discord

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

const guildIDEnvKey = "GUILD_ID"

var (
	adminPermission int64 = discordgo.PermissionAdministrator
	adminAllowed          = true
)

var (
	cmdUpdateRole         = "update-my-roles"
	cmdProcessUpdateRoles = "process-updated-roles"
	cmdAddKarma           = "plus-one"
	cmdRemoveKarma        = "minus-one"
	cmdMyKarma            = "show-my-karma"
	cmdKarmaLeaderboard   = "show-karma-leaderboard"
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
	{
		Name:        cmdKarmaLeaderboard,
		Description: "Show the top & bottom 5 in karma points",
	},
	{
		Name:        cmdMyKarma,
		Description: "Show your current karma value",
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	cmdUpdateRole:       UpdateRoles,
	cmdAddKarma:         AddKarma,
	cmdRemoveKarma:      RemoveKarma,
	cmdMplusAffixes:     GetMPlusAffixes,
	cmdKarmaLeaderboard: ShowKarmaLeaderboard,
	cmdMyKarma:          ShowMyKarma,
}

var interactionComponentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	cmdProcessUpdateRoles: ProcessUpdateRoles,
}

var GuildID string

func InitializeCommands(discordBot *discordgo.Session) ([]*discordgo.ApplicationCommand, error) {
	envGuildID, guildIDExists := os.LookupEnv(guildIDEnvKey)
	if !guildIDExists {
		log.Fatal("No Guild ID in .env file")
	}
	GuildID = envGuildID

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
		successfulSetCmd, err := discordBot.ApplicationCommandCreate(discordBot.State.User.ID, GuildID, dCmd)
		if err != nil {
			return nil, fmt.Errorf("cannot create '%v' command: %v", discordCommands[0].Name, err)
		}

		registeredCommands[dCmdIdx] = successfulSetCmd
	}

	return registeredCommands, nil
}

func TeardownAllCommands(discordBot *discordgo.Session, registeredCmds []*discordgo.ApplicationCommand) {
	for _, rCmd := range registeredCmds {
		err := discordBot.ApplicationCommandDelete(discordBot.State.User.ID, GuildID, rCmd.ID)
		if err != nil {
			fmt.Printf("Cannot delete '%v' command: %v\n", rCmd.Name, err)
		}
	}
}
