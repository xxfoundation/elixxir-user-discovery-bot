////////////////////////////////////////////////////////////////////////////////
// Copyright © 2019 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// User Discovery Bot main functions (Start Bot and register)
// This file covers all of the glue code necessary to run the bot. All of the
// interesting code is in the udb module.

package cmd

import (
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/client/globals"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/primitives/ndf"
	"gitlab.com/elixxir/user-discovery-bot/udb"
	"os"
)

// RateLimit for messages in ms (100 = 10 msg per second)
const RateLimit = 100

// UDBSessionFileName used by UDB
var UDBSessionFileName string

var clientObj *api.Client

// StartBot starts the user discovery bot:
//  - Set up global variables
//  - Log into the server
//  - Start the main loop
func StartBot(sess string, def *ndf.NetworkDefinition) {
	udb.Log.DEBUG.Printf("Starting User Discovery Bot...")

	UDBSessionFileName = sess

	// Initialize the client
	regCode := udb.UDB_USERID.RegistrationCode()
	Init(UDBSessionFileName, regCode, def)

	udb.Log.INFO.Printf("Logging in")

	// Log into the server with a blank password
	_, err := clientObj.Login("")

	if err != nil {
		udb.Log.FATAL.Panicf("Could not login: %s", err)
	}

	// Register the listeners with the user discovery bot
	udb.RegisterListeners(clientObj)

	udb.Log.INFO.Printf("Starting UDB")

	// starting the reception thread
	err = clientObj.StartMessageReceiver()
	if err != nil {
		udb.Log.FATAL.Panicf("Could not start message recievers:  %+v", err)
	}

	// Block forever as a keepalive
	select {}
}

// Initialize a session using the given session file and other info
func Init(sessionFile string, regCode string, def *ndf.NetworkDefinition) *id.User {
	userID := udb.UDB_USERID

	// We only register when the session file does not exist
	// FIXME: this is super weird -- why have to check for a file,
	// then init that file, then register optionally based on that check?
	_, err := os.Stat(sessionFile)
	// Get new client. Setting storage to nil internally creates a
	// default storage
	var initErr error

	dummyConnectionStatusHandler := func(status uint32, timeout int) {
		globals.Log.INFO.Printf("Network status: %+v, %+v", status, timeout)
	}

	clientObj, initErr = api.NewClient(nil, sessionFile, def, dummyConnectionStatusHandler)
	if initErr != nil {
		udb.Log.FATAL.Panicf("Could not initialize: %v", initErr)
	}

	if noTLS {
		clientObj.DisableTLS()
	}

	// API Settings (hard coded)
	clientObj.DisableBlockingTransmission() // Deprecated
	// Up to 10 messages per second
	clientObj.SetRateLimiting(uint32(RateLimit))

	//connect udb to gateways
	err = clientObj.Connect()
	for err != nil {
		udb.Log.ERROR.Printf("UDB could not connect to gateways, "+
			"reconnecting: %+v", err)
		err = clientObj.Connect()
	}

	// SB: Trying to always register.
	// I think it's needed for some things to work correctly.
	// Need a more accurate descriptor of what the method actually does than
	// Register, or to remove the things that aren't actually used for
	// registration.

	userID, err = clientObj.Register(true, regCode, "",
		"", "", nil)
	if err != nil {
		udb.Log.FATAL.Panicf("Could not register: %v", err)
	}

	return userID
}
