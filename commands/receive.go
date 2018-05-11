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

func ReceiveMessage(message format.MessageInterface) {
	payload := message.GetPayload()
	sender := cyclic.NewIntFromBytes(message.GetSender()).Uint64()
	jww.INFO.Printf("Sender: %d, Receiver: %v", sender, payload)
}
