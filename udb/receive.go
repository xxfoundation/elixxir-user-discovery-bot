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
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/client/cmixproto"
	"gitlab.com/elixxir/client/parse"
	"gitlab.com/elixxir/client/switchboard"
	"gitlab.com/elixxir/primitives/id"
)

type SearchListener struct{}
type RegisterListener struct{}
type PushKeyListener struct{}
type GetKeyListener struct{}

// Register the UDB listeners
func RegisterListeners() {
	Log.DEBUG.Println("Registering UDB listeners")
	api.Listen(id.ZeroID, cmixproto.Type_UDB_SEARCH, SearchListener{}, switchboard.Listeners)
	api.Listen(id.ZeroID, cmixproto.Type_UDB_REGISTER, RegisterListener{}, switchboard.Listeners)
	api.Listen(id.ZeroID, cmixproto.Type_UDB_PUSH_KEY, PushKeyListener{}, switchboard.Listeners)
	api.Listen(id.ZeroID, cmixproto.Type_UDB_GET_KEY, GetKeyListener{}, switchboard.Listeners)
}

// Listen for Search Messages
func (s SearchListener) Hear(message *parse.Message, isHeardElsewhere bool) {
	Log.DEBUG.Printf("SearchListener heard message from %q to %q: %q",
		*message.GetSender(), *message.GetRecipient(), message.GetPayload())
	sender := message.GetSender()
	if sender != nil {
		args, err := shellwords.Parse(string(message.GetPayload()))
		if err != nil {
			Log.ERROR.Printf("Error parsing message: %s", err)
		}
		Search(sender, args)
	}
}

// Listen for Register Messages
func (s RegisterListener) Hear(message *parse.Message, isHeardElsewhere bool) {
	Log.DEBUG.Printf("RegisterListener heard message from %q to %q: %q",
		*message.GetSender(), *message.GetRecipient(), message.GetPayload())
	sender := message.GetSender()
	if sender != nil {
		args, err := shellwords.Parse(string(message.GetPayload()))
		if err != nil {
			Log.ERROR.Printf("Error parsing message: %s", err)
		}
		Register(sender, args)
	}
}

// Listen for PushKey Messages
func (s PushKeyListener) Hear(message *parse.Message, isHeardElsewhere bool) {
	Log.DEBUG.Printf("PushKeyListener heard message from %q to %q: %q",
		*message.GetSender(), *message.GetRecipient(), message.GetPayload())
	sender := message.GetSender()
	if sender != nil {
		args, err := shellwords.Parse(string(message.GetPayload()))
		if err != nil {
			Log.ERROR.Printf("Error parsing message: %s", err)
		}
		PushKey(sender, args)
	}
}

// Listen for GetKey Messages
func (s GetKeyListener) Hear(message *parse.Message, isHeardElsewhere bool) {
	Log.DEBUG.Printf("GetKeyListener heard message from %q to %q: %q",
		*message.GetSender(), *message.GetRecipient(), message.GetPayload())
	sender := message.GetSender()
	if sender != nil {
		args, err := shellwords.Parse(string(message.GetPayload()))
		if err != nil {
			Log.ERROR.Printf("Error parsing message: %s", err)
		}
		GetKey(sender, args)
	}
}
