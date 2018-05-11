////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package commands

import (
	"testing"
	"gitlab.com/privategrity/crypto/format" // <-- FIXME: this is annoying, WHY?
)


func TestReceiveMessage(t *testing.T) {
	msg, err := format.NewMessage(1, 2, "Hello, World!")
	ReceiveMessage(msg[0])
	if err != nil {
		t.Errorf("Could not smoke test ReceiveMessage")
	}
}
