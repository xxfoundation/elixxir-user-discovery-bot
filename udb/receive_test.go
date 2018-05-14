////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package udb

import (
	"gitlab.com/privategrity/crypto/format" // <-- FIXME: this is annoying, WHY?
	"testing"
)

func TestReceiveMessage(t *testing.T) {
	msg, err := format.NewMessage(1, 2, "Hello, World!")
	if err != nil {
		t.Errorf("Could not smoke test ReceiveMessage: %v", err)
	}
	ReceiveMessage(msg[0])
}
