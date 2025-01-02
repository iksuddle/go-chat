package messages

import (
	"fmt"
)

// get a message with sender name
func GetMessageFrom(bytes []byte, sourceName string) string {
	return fmt.Sprintf("<%s> %s", sourceName, string(bytes))
}

// get a message saying a client joined
func GetJoinMessage(joinerName string) string {
	return fmt.Sprintf("<%s joined>", joinerName)
}

// get a message saying a client left
func GetLeaveMessage(leaverName string) string {
	return fmt.Sprintf("<%s left>", leaverName)
}
