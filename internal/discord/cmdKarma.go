package discord

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/ammuench/rerolled-bot/internal/db"

	"github.com/bwmarrin/discordgo"
)

func AddKarma(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionUser := options[0].UserValue(s)

	if doesUserHaveKarmaCmdTimeout(i.Member.User.ID) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You're doing that too much, wait a few more seconds",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			LogCmdError(err, cmdRemoveKarma, s, i)
		}
	} else if optionUser.ID == i.Member.User.ID {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "📛 You can't give yourself points 📛",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			LogCmdError(err, cmdRemoveKarma, s, i)
		}
	} else {
		parsedUserID, err := strconv.Atoi(optionUser.ID)
		if err != nil {
			LogCmdError(err, cmdRemoveKarma, s, i)
			return
		}

		newScore, err := updateUserKarma(parsedUserID, 1)
		if err != nil {
			LogCmdError(err, cmdRemoveKarma, s, i)
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error updating karma ⛔",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				LogCmdError(err, cmdRemoveKarma, s, i)
			}
		} else {
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("%v gained a point! (Now at %v)", FmtMentionString(optionUser.ID), newScore),
				},
			})
			if err != nil {
				LogCmdError(err, cmdRemoveKarma, s, i)
			}
		}
	}
}

func RemoveKarma(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionUser := options[0].UserValue(s)

	if doesUserHaveKarmaCmdTimeout(i.Member.User.ID) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You're doing that too much, wait a few more seconds",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			LogCmdError(err, cmdRemoveKarma, s, i)
		}
	} else {
		parsedUserID, err := strconv.Atoi(optionUser.ID)
		if err != nil {
			LogCmdError(err, cmdRemoveKarma, s, i)
			return
		}

		newScore, err := updateUserKarma(parsedUserID, -1)
		if err != nil {
			LogCmdError(err, cmdRemoveKarma, s, i)
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error updating karma ⛔",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				LogCmdError(err, cmdRemoveKarma, s, i)
			}
		} else {
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("%v lost a point (Now at %v)", FmtMentionString(optionUser.ID), newScore),
				},
			})
			if err != nil {
				LogCmdError(err, cmdRemoveKarma, s, i)
			}
		}
	}
}

func ShowKarmaLeaderboard(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if doesUserHaveKarmaCmdTimeout(i.Member.User.ID) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You're doing that too much, wait a few more seconds",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			LogCmdError(err, cmdKarmaLeaderboard, s, i)
		}
	} else {
		topFiveList, err := getKarmaSublist(top)
		if err != nil {
			LogCmdError(err, cmdKarmaLeaderboard, s, i)
			return
		}

		bottomFiveList, err := getKarmaSublist(bottom)
		if err != nil {
			LogCmdError(err, cmdKarmaLeaderboard, s, i)
			return
		}

		var topListFormattedString string
		var botListFormattedString string

		for topFiveEntryIdx, topFiveEntry := range topFiveList {
			topListFormattedString = topListFormattedString + fmt.Sprintf(
				"%v. %v -- (%v points)\n",
				topFiveEntryIdx+1,
				FmtMentionString(strconv.Itoa(topFiveEntry.userID)),
				strconv.Itoa(topFiveEntry.score),
			)
		}

		for botFiveIdx, botFiveEntry := range bottomFiveList {
			botListFormattedString = botListFormattedString + fmt.Sprintf(
				"%v. %v -- (%v points)\n",
				botFiveIdx+1,
				FmtMentionString(strconv.Itoa(botFiveEntry.userID)),
				strconv.Itoa(botFiveEntry.score),
			)
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf(
					`### Karma Top Five:
					 %v
					
					 ### Karma Bottom Five:
					 %v
					`,
					topListFormattedString,
					botListFormattedString,
				),
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			LogCmdError(err, cmdKarmaLeaderboard, s, i)
		}

	}
}

type UserKarma struct {
	userID int
	score  int
}

type KarmaListDirection int8

const (
	top    KarmaListDirection = iota
	bottom KarmaListDirection = iota
)

func getKarmaSublist(sublist KarmaListDirection) ([]UserKarma, error) {
	database := db.GetDB()

	sublistOrder := "ASC"
	if sublist == top {
		sublistOrder = "DESC"
	}

	var userKarmaList []UserKarma
	qText := fmt.Sprintf("SELECT userID, score FROM karma_points ORDER BY score %v LIMIT 5", sublistOrder)
	usersSublistLookup, err := database.Query(qText)
	if err != nil {
		return nil, err
	}
	defer usersSublistLookup.Close()
	for usersSublistLookup.Next() {
		var userKarmaRecord UserKarma
		err := usersSublistLookup.Scan(&userKarmaRecord.userID, &userKarmaRecord.score)
		if err == nil {
			userKarmaList = append(userKarmaList, userKarmaRecord)
		}
	}

	return userKarmaList, nil
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

var karmaTimeoutMap = make(map[string]int64)

func doesUserHaveKarmaCmdTimeout(userID string) bool {
	lastCmdTime := karmaTimeoutMap[userID]

	if lastCmdTime == 0 {
		karmaTimeoutMap[userID] = time.Now().Unix()
		return false
	}

	timeDiff := time.Now().Unix() - lastCmdTime

	if timeDiff > 10 {
		karmaTimeoutMap[userID] = time.Now().Unix()
		return false
	}

	return true
}
