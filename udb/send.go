////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Wrapper for Send command
package udb

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/privategrity/client/io"
	"gitlab.com/privategrity/client/parse"
	"gitlab.com/privategrity/client/user"
)

// Sender interface -- the api is broken here (does not return the error), so
// we reimplement a new interface...
type Sender interface {
	Send(userID user.ID, msg string) error
}

// ApiSender calls the api send function
type APISender struct{}

// Send calls the api send function
func (a APISender) Send(userID user.ID, msg string) error {
	return io.Messaging.SendMessage(userID, msg)
}

// UdbSender is the sender interface to use
var UdbSender Sender = APISender{}

// Wrap the API Send function (useful for mock tests)
func Send(userID uint64, msg string, msgType parse.Type) {
	// Create the message body and assign its type
	message := string(parse.Pack(&parse.TypedBody{
		Type: msgType,
		Body: []byte(msg),
	}))
	// Send the message
	sendErr := UdbSender.Send(user.ID(userID), message)
	if sendErr != nil {
		jww.ERROR.Printf("Error responding to %d: %s", userID, sendErr)
	}
}
