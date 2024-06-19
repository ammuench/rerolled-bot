package discord

import (
	"log"
	"slices"
	"strconv"

	"github.com/ammuench/rerolled-bot/internal/db"
	"github.com/bwmarrin/discordgo"
)

type SelectableRole struct {
	roleID int
	name   string
	emoji  string
}

func cmdError(err error, s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("Error running %v command ==> %v\n", cmdUpdateRole, err)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Uh oh, there was an error running this command",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func UpdateRoles(s *discordgo.Session, i *discordgo.InteractionCreate) {
	database := db.GetDB()
	sqlAssignableRoles, err := database.Query("SELECT roleId as roleID, name, emoji from assignable_roles WHERE enabled = true AND guildId = ?", i.GuildID)
	if err != nil {
		cmdError(err, s, i)
	}
	defer sqlAssignableRoles.Close()

	var assignableRoles []SelectableRole

	for sqlAssignableRoles.Next() {
		var role SelectableRole
		err = sqlAssignableRoles.Scan(&role.roleID, &role.name, &role.emoji)
		if err != nil {
			cmdError(err, s, i)
		}
		assignableRoles = append(assignableRoles, role)
	}

	minRoles := 0
	totalRoles := len(assignableRoles)
	generatedRoleOptions := []discordgo.SelectMenuOption{}
	for _, role := range assignableRoles {
		roleIDString := strconv.Itoa(role.roleID)
		generatedRoleOptions = append(generatedRoleOptions, discordgo.SelectMenuOption{
			Label:   role.name,
			Value:   roleIDString,
			Default: slices.Contains(i.Member.Roles, roleIDString),
			Emoji: &discordgo.ComponentEmoji{
				Name: role.emoji,
			},
		})
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hey there! Select the opt-in roles you want from the list below",
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							MenuType:    discordgo.StringSelectMenu,
							CustomID:    cmdProcessUpdateRoles,
							Placeholder: "Select roles that you want to assign to yourself",
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
		cmdError(err, s, i)
	}
}

func ProcessUpdateRoles(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()

	var selectedRoleIds []string

	for _, selectOpt := range data.Values {
		err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, selectOpt)
		selectedRoleIds = append(selectedRoleIds, selectOpt)
		if err != nil {
			cmdError(err, s, i)
			break
		}
	}

	database := db.GetDB()
	sqlAssignableRoles, err := database.Query("SELECT roleId, name, emoji from assignable_roles WHERE enabled = true AND guildId = ?", i.GuildID)
	if err != nil {
		cmdError(err, s, i)
	}
	defer sqlAssignableRoles.Close()

	for sqlAssignableRoles.Next() {
		var role SelectableRole
		err = sqlAssignableRoles.Scan(&role.roleID, &role.name, &role.emoji)
		if err != nil {
			cmdError(err, s, i)
		}
		roleString := strconv.Itoa(role.roleID)
		if !slices.Contains(selectedRoleIds, roleString) && slices.Contains(i.Member.Roles, roleString) {
			err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, roleString)
			if err != nil {
				cmdError(err, s, i)
				break
			}
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Roles updated ✅",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
