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
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/crypto/certs"
	"gitlab.com/elixxir/crypto/cyclic"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/elixxir/user-discovery-bot/udb"
	"os"
)

// FIXME: Remove
var NUM_NODES uint

// Regular globals
var GATEWAY_ADDRESS string

// Message rate limit in ms (100 = 10 msg per second)
const RATE_LIMIT = 100

// The Session file used by UDB (hard coded)
const UDB_SESSIONFILE = ".udb-cMix-session"

var clientObj *api.Client

// Startup the user discovery bot:
//  - Set up global variables
//  - Log into the server
//  - Start the main loop
func StartBot(gatewayAddr string, numNodes uint, grpConf string) {
	udb.Log.DEBUG.Printf("Starting User Discovery Bot...")

	// Use RAM storage for now
	udb.DataStore = storage.NewRamStorage()

	// Globals we need to set
	NUM_NODES = numNodes
	GATEWAY_ADDRESS = gatewayAddr

	// Initialize the client
	regCode := udb.UDB_USERID.RegistrationCode()
	userId := Init(UDB_SESSIONFILE, regCode, grpConf)

	// API Settings (hard coded)
	clientObj.DisableBlockingTransmission() // Deprecated
	// Up to 10 messages per second
	clientObj.SetRateLimiting(uint32(RATE_LIMIT))

	// Log into the server
	Login(userId)

	// Register the listeners with the user discovery bot
	udb.RegisterListeners(clientObj)

	// TEMPORARILY try starting the reception thread here instead-it seems to
	// not be starting?
	//go io.Messaging.MessageReceiver(time.Second)

	// Block forever as a keepalive
	quit := make(chan bool)
	<-quit
}

// Initialize a session using the given session file and other info
func Init(sessionFile string, regCode string, grpConf string) *id.User {
	userId := udb.UDB_USERID

	// We only register when the session file does not exist
	// FIXME: this is super weird -- why have to check for a file,
	// then init that file, then register optionally based on that check?
	_, err := os.Stat(sessionFile)
	// Get new client. Setting storage to nil internally creates a
	// default storage
	var initErr error
	clientObj, initErr = api.NewClient(nil, sessionFile)
	if initErr != nil {
		udb.Log.FATAL.Panicf("Could not initialize: %v", initErr)
	}
	// SB: Trying to always register.
	// I think it's needed for some things to work correctly.
	// Need a more accurate descriptor of what the method actually does than
	// Register, or to remove the things that aren't actually used for
	// registration.
	grp := cyclic.Group{}
	err = grp.UnmarshalJSON([]byte(grpConf))
	if err != nil {
		udb.Log.FATAL.Panicf("Could Not Decode group from JSON: %s\n", err.Error())
	}
	userId, err = clientObj.Register(regCode, GATEWAY_ADDRESS, NUM_NODES, false, &grp)
	if err != nil {
		udb.Log.FATAL.Panicf("Could not register: %v", err)
	}

	return userId
}

// Log into the server using the user id generated from Init
func Login(userId *id.User) {
	clientObj.Login(userId, GATEWAY_ADDRESS, certs.GatewayTLS)
}
