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
	client "gitlab.com/elixxir/client/api"
	clientGlobals "gitlab.com/elixxir/client/globals"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/elixxir/user-discovery-bot/udb"
	"os"
	"gitlab.com/elixxir/crypto/id"
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
	udb.Log.DEBUG.Printf("Starting User Discovery Bot...")

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

	// Register the listeners with the user discovery bot
	udb.RegisterListeners()

	// Log into the server
	Login(userId)

	// TEMPORARILY try starting the reception thread here instead-it seems to
	// not be starting?
	//go io.Messaging.MessageReceiver(time.Second)

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
		udb.Log.FATAL.Panicf("Could not initialize: %v", initErr)
	}
	// SB: Trying to always register.
	// I think it's needed for some things to work correctly.
	// Need a more accurate descriptor of what the method actually does than
	// Register, or to remove the things that aren't actually used for
	// registration.
	userId, err = client.Register(regCode, GATEWAY_ADDRESS, NUM_NODES, false)
	if err != nil {
		udb.Log.FATAL.Panicf("Could not register: %v", err)
	}

	return userId
}

// Log into the server using the user id generated from Init
func Login(userId *id.UserID) {
	client.Login(userId, GATEWAY_ADDRESS)
}
