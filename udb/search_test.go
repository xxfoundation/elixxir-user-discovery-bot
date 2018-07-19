////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package udb

import (
	"gitlab.com/privategrity/user-discovery-bot/storage"
	"testing"
)

func TestSearchHappyPath(t *testing.T) {
	DataStore = storage.NewRamStorage()
	// Load a user
	TestRegisterHappyPath(t)
	fingerprint := "8oKh7TYG4KxQcBAymoXPBHSD/uga9pX3Mn/jKhvcD8M="
	// NOTE: This is kind of hard, since we can't see the response and search
	//       does not modify data we can check
	// TODO: Monkeypatch send so we can verify? -- this is tested in integration,
	//       so.. low priority.
	msgs := []string{
		"SEARCH EMAIL rick@privategrity.com",
		"GETKEY " + fingerprint,
	}

	for i := range msgs {
		msg, err := NewMessage(msgs[i])
		if err != nil {
			t.Errorf("Error generating message: %v", err)
		}
		ReceiveMessage(msg)
	}
}

// Test invalid search type
func TestSearch_Invalid_Type(t *testing.T) {
	fingerprint := "8oKh7TYG4KxQcBAymoXPBHSD/uga9pX3Mn/jKhvcD8M="
	msgs := []string{
		"SEARCH INVALID test",
		"GETKEY " + fingerprint,
	}

	for i := range msgs {
		msg, err := NewMessage(msgs[i])
		if err != nil {
			t.Errorf("Error generating message: %v", err)
		}
		ReceiveMessage(msg)
	}
}

// Test invalid user
func TestSearch_Invalid_User(t *testing.T) {
	msgs := []string{
		"SEARCH EMAIL cat@privategrity.com",
	}

	for i := range msgs {
		msg, err := NewMessage(msgs[i])
		if err != nil {
			t.Errorf("Error generating message: %v", err)
		}
		ReceiveMessage(msg)
	}
}
