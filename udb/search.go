////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Search Command
package udb

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/privategrity/user-discovery-bot/storage"
	"fmt"
)

const SEARCH_USAGE = ("Usage: 'SEARCH [EMAIL] [email-address]'")

// Search for an entry in the database
// The search command takes the form "SEARCH TYPE VALUE"
// WHERE:
// - TYPE = EMAIL
// - VALUE = "rick@privategrity.com"
// It returns a list of fingerprints if found (1 per message), or NOTFOUND
func Search(userId uint64, args []string) {
	jww.INFO.Printf("Search %d: %v", userId, args)
	SearchErr := func(msg string) {
		Send(userId, msg)
		Send(userId, SEARCH_USAGE)
		jww.INFO.Printf("User %d, error: %s", userId, msg)
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
	regTypeEnum := storage.Email

	keyFingerprints, ok := DataStore.GetKeys(regVal, regTypeEnum)
	if !ok {
		msg := fmt.Sprintf("SEARCH %s NOTFOUND", regVal)
		jww.INFO.Printf("User %d: %s", msg)
		Send(userId, msg)
		return
	}

	for i := range keyFingerprints {
		msg := fmt.Sprintf("SEARCH %s FOUND %s", keyFingerprints[i])
		jww.INFO.Printf("User %s: %s", msg)
		Send(userId, msg)
	}
}
