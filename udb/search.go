////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Search Command
package udb

import (
	"encoding/base64"
	"fmt"
	"gitlab.com/elixxir/client/cmixproto"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/user-discovery-bot/storage"
)

const SEARCH_USAGE = "Usage: 'SEARCH [EMAIL] [email-address]'"

// Search for an entry in the database
// The search command takes the form "SEARCH TYPE VALUE"
// WHERE:
// - TYPE = EMAIL
// - VALUE = "rick@elixxir.io"
// It returns a list of fingerprints if found (1 per message), or NOTFOUND
func Search(userId *id.User, args []string) {
	Log.INFO.Printf("Search %d: %v", userId, args)
	SearchErr := func(msg string) {
		Send(userId, msg, cmixproto.Type_UDB_SEARCH_RESPONSE)
		Send(userId, SEARCH_USAGE, cmixproto.Type_UDB_SEARCH_RESPONSE)
		Log.INFO.Printf("User %d, error: %s", userId, msg)
	}
	if len(args) != 2 {
		SearchErr("Invalid command syntax!")
		return
	}

	regType := args[0]
	regVal := args[1]

	// Verify that regType == EMAIL
	// TODO: Functionalize this. Leaving it be for now.
	if regType != "EMAIL" {
		SearchErr("EMAIL is the only acceptable registration type")
		return
	}
	// TODO: Add parse func to storage class, embed into function and
	// pass it a string instead
	usr := storage.NewUser()
	usr.SetValue(regVal)

	// Get the userID associated to email
	searchedUser, err := storage.UserDiscoveryDb.GetUser(usr)
	//DataStore.GetUserID(regVal)
	if err != nil {
		msg := fmt.Sprintf("SEARCH %s NOTFOUND", regVal)
		Log.INFO.Printf("User %d: %s", userId, msg)
		Send(userId, msg, cmixproto.Type_UDB_SEARCH_RESPONSE)
		return
	}

	searchedUserID := searchedUser.Id
	searchedUserKeyID := searchedUser.KeyId

	// Correctly send the messages with actual userID based on email, followed by key fingerprint
	//TODO: Make sure this thing works
	for i := range searchedUserKeyID {
		msg := fmt.Sprintf("SEARCH %s FOUND %+v %+v", regVal,
			base64.StdEncoding.EncodeToString(searchedUserID[:]), searchedUserKeyID[i])
		Log.INFO.Printf("User %d: %s", userId, msg)
		Send(userId, msg, cmixproto.Type_UDB_SEARCH_RESPONSE)
	}
}
