////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package cmix

import (
	"fmt"
	"gitlab.com/elixxir/client/cmixproto"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/primitives/id"
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
	fmt.Println(storage.UserDiscoveryDb)

	sender := id.NewIdFromUInt(89, id.User, t)

	msg := NewMessage(msgs[0], cmixproto.Type_UDB_SEARCH, sender, t)
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
	}
	sender := id.NewIdFromUInt(122, id.User, t)
	msg := NewMessage(msgs[0], cmixproto.Type_UDB_SEARCH, sender, t)
	sl.Hear(msg, false)

}

func TestSearch_InvalidArgs_Email(t *testing.T) {
	// Load a user
	TestRegisterHappyPath(t)
	// NOTE: This is kind of hard, since we can't see the response and search
	//       does not modify data we can check
	// TODO: Monkeypatch send so we can verify? -- this is tested in integration,
	//       so.. low priority.
	msgs := []string{
		"NotEMAIL rick@elixxir.io",
	}

	sender := id.NewIdFromUInt(43, id.User, t)

	msg := NewMessage(msgs[0], cmixproto.Type_UDB_SEARCH, sender, t)
	sl.Hear(msg, false)

}

// Test invalid search type
func TestSearch_Invalid_Type(t *testing.T) {
	defer func() {}()
	fingerprint := "8oKh7TYG4KxQcBAymoXPBHSD/uga9pX3Mn/jKhvcD8M="
	msgs := []string{
		"SEARCH INVALID test",
		"GETKEY " + fingerprint,
	}
	sender := id.NewIdFromUInt(222, id.User, t)
	msg := NewMessage(msgs[0], cmixproto.Type_UDB_SEARCH, sender, t)
	sl.Hear(msg, false)
}

// Test invalid user
func TestSearch_Invalid_User(t *testing.T) {

	msgs := []string{
		"SEARCH EMAIL cat@elixxir.io",
	}
	sender := id.NewIdFromUInt(9000, id.User, t)
	msg := NewMessage(msgs[0], cmixproto.Type_UDB_SEARCH, sender, t)
	sl.Hear(msg, false)
}
