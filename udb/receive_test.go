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

// Hack around the interface for client to do what we need for testing.
func NewMessage(msg string) (format.MessageInterface, error) {
	msgs, err := format.NewMessage(1, 2, msg)
	return msgs[0], err
}

// Test with an unknown function
func TestReceiveMessage(t *testing.T) {
	msg, err := NewMessage("Hello, World!")
	if err != nil {
		t.Errorf("Could not smoke test ReceiveMessage: %v", err)
	}
	ReceiveMessage(msg)
}

func TestBrokenMessage(t *testing.T) {
	brokenMsg := "foo '" // From shellwords test cases
	msg, _ := NewMessage(brokenMsg)
	ReceiveMessage(msg)
	// We are only making sure this doesn't crash the program.
}
