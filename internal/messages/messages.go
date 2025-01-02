package messages

import (
	"fmt"
)

func GetMessageFrom(bytes []byte, sourceName string) string {
	return fmt.Sprintf("<%s> %s", sourceName, string(bytes))
}

func GetJoinMessage(joinerName string) string {
	return fmt.Sprintf("<%s joined>", joinerName)
}

func GetLeaveMessage(leaverName string) string {
	return fmt.Sprintf("<%s left>", leaverName)
}
