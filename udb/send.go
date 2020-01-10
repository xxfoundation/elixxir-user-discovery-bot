////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Wrapper for Send command
package udb

import (
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/client/cmixproto"
	"gitlab.com/elixxir/client/parse"
	"gitlab.com/elixxir/primitives/id"
)

// Sender interface -- the api is broken here (does not return the error), so
// we reimplement a new interface...
type Sender interface {
	Send(recipientID *id.User, msg string, msgType cmixproto.Type)
}

// ApiSender calls the api send function
type APISender struct {
	ClientObj *api.Client
}

// Send calls the api send function
func (a APISender) Send(recipientID *id.User, msg string, msgType cmixproto.Type) {
	message := string(parse.Pack(&parse.TypedBody{
		MessageType: int32(msgType),
		Body:        []byte(msg),
	}))
	sendErr := a.ClientObj.Send(api.APIMessage{
		Payload:     []byte(message),
		SenderID:    UDB_USERID,
		RecipientID: recipientID,
	})
	if sendErr != nil {
		Log.ERROR.Printf("Error responding to %d: %s", recipientID, sendErr)
	}
}
