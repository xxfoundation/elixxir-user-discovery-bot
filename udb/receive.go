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
	"gitlab.com/elixxir/primitives/format"
	"gitlab.com/elixxir/primitives/switchboard"
	"gitlab.com/elixxir/primitives/id"
)

type SearchListener struct{}
type RegisterListener struct{}
type PushKeyListener struct{}
type GetKeyListener struct{}

// Register the UDB listeners
func RegisterListeners(cl *api.Client) {
	Log.DEBUG.Println("Registering UDB listeners")
	cl.Listen(id.ZeroID, format.None, int32(cmixproto.Type_UDB_SEARCH),
		SearchListener{})
	cl.Listen(id.ZeroID, format.None, int32(cmixproto.Type_UDB_REGISTER),
		RegisterListener{})
	cl.Listen(id.ZeroID, format.None, int32(cmixproto.Type_UDB_PUSH_KEY),
		PushKeyListener{})
	cl.Listen(id.ZeroID, format.None, int32(cmixproto.Type_UDB_GET_KEY),
		GetKeyListener{})
	clientObj = cl
}

// Listen for Search Messages
func (s SearchListener) Hear(item switchboard.Item, isHeardElsewhere bool) {
	message := item.(*parse.Message)
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
func (s RegisterListener) Hear(item switchboard.Item, isHeardElsewhere bool) {
	message := item.(*parse.Message)
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
func (s PushKeyListener) Hear(item switchboard.Item, isHeardElsewhere bool) {
	message := item.(*parse.Message)
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
func (s GetKeyListener) Hear(item switchboard.Item, isHeardElsewhere bool) {
	message := item.(*parse.Message)
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
