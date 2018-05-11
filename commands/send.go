////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Wrapper for Send command
package commands

import (
	jww "github.com/spf13/jwalterweatherman"
	client "gitlab.com/privategrity/client/api"
	clientGlobals "gitlab.com/privategrity/client/globals"
	"gitlab.com/privategrity/crypto/format" // <-- FIXME: this is annoying, WHY?
)

// Wrap the API Send function (useful for mock tests)
func Send(userId uint64, msg string) {
	myId := clientGlobals.Session.GetCurrentUser().UserID
	messages, err := format.NewMessage(myId, userId, msg)
	if err != nil {
		jww.FATAL.Panicf("Error creating message: %d, %d, %s",
			myId, userId, msg)
	}

	for i := range messages {
		sendErr := client.Send(messages[i])
		if sendErr != nil {
			jww.ERROR.Println("Error responding to %d", userId)
		}
	}
}
