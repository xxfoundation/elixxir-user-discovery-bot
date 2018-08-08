////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package udb

import (
	"gitlab.com/privategrity/client/parse"
	"gitlab.com/privategrity/user-discovery-bot/storage"
	"testing"
)

func TestSearchHappyPath(t *testing.T) {
	DataStore = storage.NewRamStorage()
	// Load a user
	TestRegisterHappyPath(t)
	// NOTE: This is kind of hard, since we can't see the response and search
	//       does not modify data we can check
	// TODO: Monkeypatch send so we can verify? -- this is tested in integration,
	//       so.. low priority.
	msgs := []string{
		"SEARCH EMAIL rick@privategrity.com",
	}

	msg := NewMessage(msgs[0], parse.Type_UDB_SEARCH)
	sl.Hear(msg, false)
}

// Test invalid search type
func TestSearch_Invalid_Type(t *testing.T) {
	fingerprint := "8oKh7TYG4KxQcBAymoXPBHSD/uga9pX3Mn/jKhvcD8M="
	msgs := []string{
		"SEARCH INVALID test",
		"GETKEY " + fingerprint,
	}

	msg := NewMessage(msgs[0], parse.Type_UDB_SEARCH)
	sl.Hear(msg, false)
}

// Test invalid user
func TestSearch_Invalid_User(t *testing.T) {
	DataStore = storage.NewRamStorage()
	msgs := []string{
		"SEARCH EMAIL cat@privategrity.com",
	}

	msg := NewMessage(msgs[0], parse.Type_UDB_SEARCH)
	sl.Hear(msg, false)
}
