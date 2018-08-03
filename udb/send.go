////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Wrapper for Send command
package udb

import (
	jww "github.com/spf13/jwalterweatherman"
	client "gitlab.com/privategrity/client/api"
	"gitlab.com/privategrity/client/parse"
	"gitlab.com/privategrity/crypto/format" // <-- FIXME: this is annoying, WHY?
)

// Sender interface -- the api is broken here (does not return the error), so
// we reimplement a new interface...
type Sender interface {
	Send(messageInterface format.MessageInterface) error
}

// ApiSender calls the api send function
type APISender struct{}

// Send calls the api send function
func (a APISender) Send(message format.MessageInterface) error {
	return client.Send(message)
}

// UdbSender is the sender interface to use
var UdbSender Sender = APISender{}

// Wrap the API Send function (useful for mock tests)
func Send(userID uint64, msg string, msgType parse.Type) {
	// Create the message body and assign its type
	msgBody := &parse.TypedBody{
		Type: msgType,
		Body: []byte(msg),
	}
	myID := uint64(UDB_USERID)
	messages, err := format.NewMessage(myID, userID, msg)
	if err != nil {
		jww.FATAL.Panicf("Error creating message: %d, %d, %s",
			myID, userID, msg)
	}

	for i := range messages {
		sendErr := UdbSender.Send(messages[i])
		if sendErr != nil {
			jww.ERROR.Printf("Error responding to %d", userID)
		}
	}
}
