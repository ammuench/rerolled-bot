package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-zoox/fetch"
)

var (
	affixesLastFetch int64
	affixesCache     string
)

func GetMPlusAffixes(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if doesUserHaveAffixesCmdTimeout(i.Member.User.ID) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You're doing that too much, wait a few more seconds",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			LogCmdError(err, cmdMplusAffixes, s, i)
		}
	} else {
		affixes, err := fetchAffixes()
		if err != nil {
			LogCmdError(err, cmdMplusAffixes, s, i)
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error fetching affixes, try again in a bit",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				LogCmdError(err, cmdMplusAffixes, s, i)
			}

		}
		affixSlice := strings.Split(affixes, ", ")
		formattedAffixString := "* `" + strings.Join(affixSlice, "`\n* `") + "`"
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Current Mythic+ Affixes are:\n%v", formattedAffixString),
			},
		})
		if err != nil {
			LogCmdError(err, cmdMplusAffixes, s, i)
		}
	}
}

type Affixes struct {
	Affixes string `json:"title"`
}

func fetchAffixes() (string, error) {
	if (time.Now().Unix() - affixesLastFetch) < 3600 {
		return affixesCache, nil
	}

	response, err := fetch.Get("https://mythicpl.us/affix-us")
	if err != nil {
		return "", err
	}

	var affixes Affixes

	err = response.UnmarshalJSON(&affixes)
	if err != nil {
		return "", err
	}

	affixesCache = affixes.Affixes
	affixesLastFetch = time.Now().Unix()

	return affixes.Affixes, nil
}

var affixesTimeoutMap = make(map[string]int64)

func doesUserHaveAffixesCmdTimeout(userID string) bool {
	lastCmdTime := affixesTimeoutMap[userID]

	if lastCmdTime == 0 {
		affixesTimeoutMap[userID] = time.Now().Unix()
		return false
	}

	timeDiff := time.Now().Unix() - lastCmdTime

	if timeDiff > 60 {
		affixesTimeoutMap[userID] = time.Now().Unix()
		return false
	}

	return true
}
