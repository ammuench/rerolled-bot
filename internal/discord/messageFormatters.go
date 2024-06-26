package discord

import "fmt"

func FmtMentionString(userID string) string {
	return fmt.Sprintf("<@%v>", userID)
}
