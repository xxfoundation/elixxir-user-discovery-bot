////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Receive and parse user discovery bot messages, and run the appropriate
// command
package udb

import (
	"github.com/mattn/go-shellwords"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/client/parse"
	"gitlab.com/elixxir/client/switchboard"
	"gitlab.com/elixxir/crypto/id"
	"gitlab.com/elixxir/client/cmixproto"
)

type SearchListener struct{}
type RegisterListener struct{}
type PushKeyListener struct{}
type GetKeyListener struct{}

// Register the UDB listeners
func init() {
	switchboard.Listeners.Register(id.ZeroID, cmixproto.Type_UDB_SEARCH,
		SearchListener{})
	switchboard.Listeners.Register(id.ZeroID, cmixproto.Type_UDB_REGISTER,
		RegisterListener{})
	switchboard.Listeners.Register(id.ZeroID, cmixproto.Type_UDB_PUSH_KEY,
		PushKeyListener{})
	switchboard.Listeners.Register(id.ZeroID, cmixproto.Type_UDB_GET_KEY,
		GetKeyListener{})
}

// Listen for Search Messages
func (s SearchListener) Hear(message *parse.Message, isHeardElsewhere bool) {
	sender := message.GetSender()
	if sender != nil {
		args, err := shellwords.Parse(string(message.GetPayload()))
		if err != nil {
			jww.ERROR.Printf("Error parsing message: %s", err)
		}
		Search(sender, args[1:])
	}
}

// Listen for Register Messages
func (s RegisterListener) Hear(message *parse.Message, isHeardElsewhere bool) {
	sender := message.GetSender()
	if sender != nil {
		args, err := shellwords.Parse(string(message.GetPayload()))
		if err != nil {
			jww.ERROR.Printf("Error parsing message: %s", err)
		}
		Register(sender, args[1:])
	}
}

// Listen for PushKey Messages
func (s PushKeyListener) Hear(message *parse.Message, isHeardElsewhere bool) {
	sender := message.GetSender()
	if sender != nil {
		args, err := shellwords.Parse(string(message.GetPayload()))
		if err != nil {
			jww.ERROR.Printf("Error parsing message: %s", err)
		}
		PushKey(sender, args[1:])
	}
}

// Listen for GetKey Messages
func (s GetKeyListener) Hear(message *parse.Message, isHeardElsewhere bool) {
	sender := message.GetSender()
	if sender != nil {
		args, err := shellwords.Parse(string(message.GetPayload()))
		if err != nil {
			jww.ERROR.Printf("Error parsing message: %s", err)
		}
		GetKey(sender, args[1:])
	}
}
