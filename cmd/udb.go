////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// User Discovery Bot main functions (Start Bot and register)
// This file covers all of the glue code necessary to run the bot. All of the
// interesting code is in the udb module.

package cmd

import (
	jww "github.com/spf13/jwalterweatherman"
	client "gitlab.com/privategrity/client/api"
	clientGlobals "gitlab.com/privategrity/client/globals"
	"gitlab.com/privategrity/user-discovery-bot/storage"
	"gitlab.com/privategrity/user-discovery-bot/udb"
	"os"
	"gitlab.com/privategrity/crypto/id"
)

// FIXME: Remove
var NUM_NODES uint

// Regular globals
var GATEWAY_ADDRESS string

// Message rate limit in ms (100 = 10 msg per second)
const RATE_LIMIT = 100

// The Session file used by UDB (hard coded)
const UDB_SESSIONFILE = ".udb-cMix-session"

// Startup the user discovery bot:
//  - Set up global variables
//  - Log into the server
//  - Start the main loop
func StartBot(gatewayAddr string, numNodes uint) {
	jww.DEBUG.Printf("Starting User Discovery Bot...")

	// Use RAM storage for now
	udb.DataStore = storage.NewRamStorage()

	// Globals we need to set
	NUM_NODES = numNodes
	GATEWAY_ADDRESS = gatewayAddr

	// API Settings (hard coded)
	client.DisableBlockingTransmission() // Deprecated
	// Up to 10 messages per second
	client.SetRateLimiting(uint32(RATE_LIMIT))

	// Initialize the client
	regCode := udb.UDB_USERID.RegistrationCode()
	userId := Init(UDB_SESSIONFILE, regCode)

	// Log into the server
	Login(userId)

	// TODO Set up message listeners to handle commands as they come in

	// Block forever as a keepalive
	quit := make(chan bool)
	<-quit
}

// Initialize a session using the given session file and other info
func Init(sessionFile string, regCode string) *id.UserID {
	userId := udb.UDB_USERID

	// We only register when the session file does not exist
	// FIXME: this is super weird -- why have to check for a file,
	// then init that file, then register optionally based on that check?
	_, err := os.Stat(sessionFile)
	// Init regardless, wow this is broken...
	initErr := client.InitClient(&clientGlobals.DefaultStorage{}, sessionFile)
	if initErr != nil {
		jww.FATAL.Panicf("Could not initialize: %v", initErr)
	}
	if os.IsNotExist(err) {
		userId, err = client.Register(regCode, GATEWAY_ADDRESS, NUM_NODES, false)
		if err != nil {
			jww.FATAL.Panicf("Could not register: %v", err)
		}
	}

	return userId
}

// Log into the server using the user id generated from Init
func Login(userId *id.UserID) {
	client.Login(userId, GATEWAY_ADDRESS)
}
