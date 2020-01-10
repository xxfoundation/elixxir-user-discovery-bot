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
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/primitives/switchboard"
)

type SearchListener struct {
	Sender
}
type RegisterListener struct {
	Sender
	blacklist BlackList
}
type PushKeyListener struct {
	Sender
}
type GetKeyListener struct {
	Sender
}

// Register the UDB listeners
func RegisterListeners(cl *api.Client, blacklist BlackList) {
	Log.DEBUG.Println("Registering UDB listeners")
	sender := APISender{cl}
	cl.Listen(id.ZeroID, int32(cmixproto.Type_UDB_SEARCH), SearchListener{sender})
	cl.Listen(id.ZeroID, int32(cmixproto.Type_UDB_REGISTER), RegisterListener{sender, blacklist})
	cl.Listen(id.ZeroID, int32(cmixproto.Type_UDB_PUSH_KEY), PushKeyListener{sender})
	cl.Listen(id.ZeroID, int32(cmixproto.Type_UDB_GET_KEY), GetKeyListener{sender})
}

// Listen for Search Messages
func (s SearchListener) Hear(item switchboard.Item, isHeardElsewhere bool, i ...interface{}) {
	message := item.(*parse.Message)
	Log.DEBUG.Printf("SearchListener heard message from %q to %q: %q",
		*message.GetSender(), *message.GetRecipient(), message.GetPayload())
	senderID := message.GetSender()
	if senderID != nil {
		args, err := shellwords.Parse(string(message.GetPayload()))
		if err != nil {
			Log.ERROR.Printf("Error parsing message: %s", err)
		}
		Search(senderID, args, s.Sender)
	}
}

// Listen for Register Messages
func (s RegisterListener) Hear(item switchboard.Item, isHeardElsewhere bool, i ...interface{}) {
	message := item.(*parse.Message)
	Log.DEBUG.Printf("RegisterListener heard message from %q to %q: %q",
		*message.GetSender(), *message.GetRecipient(), message.GetPayload())
	senderID := message.GetSender()
	if senderID != nil {
		args, err := shellwords.Parse(string(message.GetPayload()))
		if err != nil {
			Log.ERROR.Printf("Error parsing message: %s", err)
		}
		Register(senderID, args, s.blacklist, s.Sender)
	}
}

// Listen for PushKey Messages
func (s PushKeyListener) Hear(item switchboard.Item, isHeardElsewhere bool, i ...interface{}) {
	message := item.(*parse.Message)
	Log.DEBUG.Printf("PushKeyListener heard message from %q to %q: %q",
		*message.GetSender(), *message.GetRecipient(), message.GetPayload())
	senderID := message.GetSender()
	if senderID != nil {
		args, err := shellwords.Parse(string(message.GetPayload()))
		if err != nil {
			Log.ERROR.Printf("Error parsing message: %s", err)
		}
		PushKey(senderID, args, s.Sender)
	}
}

// Listen for GetKey Messages
func (s GetKeyListener) Hear(item switchboard.Item, isHeardElsewhere bool, i ...interface{}) {
	message := item.(*parse.Message)
	Log.DEBUG.Printf("GetKeyListener heard message from %q to %q: %q",
		*message.GetSender(), *message.GetRecipient(), message.GetPayload())
	senderID := message.GetSender()
	if senderID != nil {
		args, err := shellwords.Parse(string(message.GetPayload()))
		if err != nil {
			Log.ERROR.Printf("Error parsing message: %s", err)
		}
		GetKey(senderID, args, s.Sender)
	}
}
