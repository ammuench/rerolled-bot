package discord

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/ammuench/rerolled-bot/internal/db"

	"github.com/bwmarrin/discordgo"
)

func AddKarma(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionUser := options[0].UserValue(s)

	parsedUserID, err := strconv.Atoi(optionUser.ID)
	if err != nil {
		LogCmdError(err, &cmdRemoveKarma, s, i)
		return
	}

	newScore, err := updateUserKarma(parsedUserID, 1)
	if err != nil {
		LogCmdError(err, &cmdRemoveKarma, s, i)
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error updating karma ⛔",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			LogCmdError(err, &cmdRemoveKarma, s, i)
		}
	} else {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("%v gained a point! (Now at %v)", FmtMentionString(optionUser.ID), newScore),
			},
		})
		if err != nil {
			LogCmdError(err, &cmdRemoveKarma, s, i)
		}
	}
}

func RemoveKarma(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionUser := options[0].UserValue(s)
	parsedUserID, err := strconv.Atoi(optionUser.ID)
	if err != nil {
		LogCmdError(err, &cmdRemoveKarma, s, i)
		return
	}

	newScore, err := updateUserKarma(parsedUserID, -1)
	if err != nil {
		LogCmdError(err, &cmdRemoveKarma, s, i)
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error updating karma ⛔",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			LogCmdError(err, &cmdRemoveKarma, s, i)
		}
	} else {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("%v lost a point (Now at %v)", FmtMentionString(optionUser.ID), newScore),
			},
		})
		if err != nil {
			LogCmdError(err, &cmdRemoveKarma, s, i)
		}
	}
}

type UserKarma struct {
	userID int
	score  int
}

func updateUserKarma(userID int, adjustmentAmt int) (int, error) {
	database := db.GetDB()

	var userKarmaState UserKarma

	userLookup := database.QueryRow("SELECT * FROM karma_points WHERE userID = ?", userID)
	err := userLookup.Scan(&userKarmaState.userID, &userKarmaState.score)
	if err != nil {
		if err == sql.ErrNoRows {
			_, updateErr := database.Exec("INSERT INTO karma_points (userID, score) VALUES (?, ?);", userID, adjustmentAmt)
			if updateErr != nil {
				return 0, updateErr
			}

			return adjustmentAmt, nil
		}

		return 0, err
	} else {

		newUserScore := userKarmaState.score + adjustmentAmt

		_, updateErr := database.Exec("UPDATE karma_points SET score = ? WHERE userID = ?", &newUserScore, &userKarmaState.userID)
		if updateErr != nil {
			return 0, updateErr
		}

		return newUserScore, nil
	}
}
