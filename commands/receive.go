////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Receive and parse user discovery bot messages, and run the appropriate
// command
package commands

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/privategrity/crypto/cyclic" // <-- FIXME: this is annoying, WHY?
	"gitlab.com/privategrity/crypto/format" // <-- FIXME: this is annoying, WHY?
)

// Parse the command and run the corresponding function
func ReceiveMessage(message format.MessageInterface) {
	payload := message.GetPayload()
	sender := cyclic.NewIntFromBytes(message.GetSender()).Uint64()

	// Parse the command and run the returned function
	cmdFn, args := ParseCommand(payload)

	cmdFn(sender, args)
}

// Respond to the sender that the command does not exist
func UnknownCommand(userId uint64, args []string) {
	// 1 argument, the exact command string send to the function
	jww.WARN.Printf("Received Unknown Command from %d: %s", userId, args[0])
	msg := "Unknown Command: " + args[0]
	Send(userId, msg)
}

// ParseCommand parses the message payload and return the function with it's
// arguments
func ParseCommand(cmdMsg string) (func(uint64, []string), []string) {
	args := []string{}
	args = append(args, cmdMsg)
	return UnknownCommand, args
}
