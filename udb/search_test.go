////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package udb

import (
	"gitlab.com/elixxir/client/cmixproto"
	"testing"
)

func TestSearchHappyPath(t *testing.T) {
	// Load a user
	TestRegisterHappyPath(t)
	// NOTE: This is kind of hard, since we can't see the response and search
	//       does not modify data we can check
	// TODO: Monkeypatch send so we can verify? -- this is tested in integration,
	//       so.. low priority.
	msgs := []string{
		"EMAIL rick@elixxir.io",
	}

	msg := NewMessage(msgs[0], cmixproto.Type_UDB_SEARCH)
	sl.Hear(msg, false)
}


func TestSearch_InvalidArgs(t *testing.T) {
	// Load a user
	TestRegisterHappyPath(t)
	// NOTE: This is kind of hard, since we can't see the response and search
	//       does not modify data we can check
	// TODO: Monkeypatch send so we can verify? -- this is tested in integration,
	//       so.. low priority.
	msgs := []string{
		"EMAIL rick@elixxir.io",
		"GETKEY MORETHAN 1ARG",
	}

	msg := NewMessage(msgs[0], cmixproto.Type_UDB_SEARCH)
	sl.Hear(msg, false)

	msg = NewMessage(msgs[1], cmixproto.Type_UDB_SEARCH)
	sl.Hear(msg,false)
}

// Test invalid search type
func TestSearch_Invalid_Type(t *testing.T) {
	fingerprint := "8oKh7TYG4KxQcBAymoXPBHSD/uga9pX3Mn/jKhvcD8M="
	msgs := []string{
		"SEARCH INVALID test",
		"GETKEY " + fingerprint,
	}

	msg := NewMessage(msgs[0], cmixproto.Type_UDB_SEARCH)
	sl.Hear(msg, false)
}

// Test invalid user
func TestSearch_Invalid_User(t *testing.T) {

	msgs := []string{
		"SEARCH EMAIL cat@elixxir.io",
	}

	msg := NewMessage(msgs[0], cmixproto.Type_UDB_SEARCH)
	sl.Hear(msg, false)
}
