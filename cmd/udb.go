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
	"crypto/rand"
	"encoding/json"
	"github.com/mitchellh/go-homedir"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/crypto/certs"
	"gitlab.com/elixxir/crypto/cyclic"
	"gitlab.com/elixxir/crypto/signature"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/elixxir/user-discovery-bot/udb"
	"io/ioutil"
	"os"
)

// Regular globals
var GATEWAY_ADDRESSES []string
var REG_ADDRESS string

// Message rate limit in ms (100 = 10 msg per second)
const RATE_LIMIT = 100

// The Session file used by UDB (hard coded)
const UDB_SESSIONFILE = ".udb-cMix-session"

var clientObj *api.Client

// Startup the user discovery bot:
//  - Set up global variables
//  - Log into the server
//  - Start the main loop
func StartBot(gatewayAddr []string, registrationAddr, regCode, grpConf string) {
	udb.Log.DEBUG.Printf("Starting User Discovery Bot...")

	// Use RAM storage for now
	udb.DataStore = storage.NewRamStorage()

	GATEWAY_ADDRESSES = gatewayAddr
	REG_ADDRESS = registrationAddr

	// Initialize the client
	udb.UDB_USERID = Init(UDB_SESSIONFILE, regCode, grpConf)

	// Save DSA public key and user ID to JSON file
	outputDsaPubKeyToJson(udb.UDB_USERID, ".elixxir", "udb_info.json")

	// API Settings (hard coded)
	clientObj.DisableBlockingTransmission() // Deprecated
	// Up to 10 messages per second
	clientObj.SetRateLimiting(uint32(RATE_LIMIT))

	// Log into the server
	Login(udb.UDB_USERID)

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

	userId, err := clientObj.Register(false, regCode, "UDB",
		REG_ADDRESS, GATEWAY_ADDRESSES, false, &grp)

	if err != nil {
		udb.Log.FATAL.Panicf("Could not register: %v", err)
	} else {
		udb.Log.DEBUG.Printf("UDB registered as user %v", *userId)
	}

	return userId
}

// Log into the server using the user id generated from Init
func Login(userId *id.User) {
	_, err := clientObj.Login(userId, "", GATEWAY_ADDRESSES[0],
		certs.GatewayTLS)

	if err != nil {
		udb.Log.FATAL.Panicf("Could not log into the server: %s", err)
	}
}

// outputDsaPubKeyToJson encodes the DSA public key and user ID to JSON and
// outputs it to the specified directory with the specified file name.
func outputDsaPubKeyToJson(userID *id.User, dir, fileName string) {

	// Get the default parameters and generate a public key from it
	dsaParams := signature.GetDefaultDSAParams()
	publicKey := dsaParams.PrivateKeyGen(rand.Reader).PublicKeyGen()

	// Setup struct that will dictate the JSON structure
	jsonStruct := struct {
		Id             *id.User
		Dsa_public_key *signature.DSAPublicKey
	}{
		Id:             userID,
		Dsa_public_key: publicKey,
	}

	// Generate JSON from structure
	data, err := json.MarshalIndent(jsonStruct, "", "\t")
	if err != nil {
		jww.ERROR.Printf("Error encoding structure to JSON: %s", err)
	}

	// Get the user's home directory
	homeDir, err := homedir.Dir()
	if err != nil {
		jww.ERROR.Printf("Unable to retrieve user's home directory: %s", err)
	}

	// Write JSON to file
	err = ioutil.WriteFile(homeDir+"/"+dir+"/"+fileName, data, 0644)
	if err != nil {
		jww.ERROR.Printf("Error writing JSON file: %s", err)
	}
}
