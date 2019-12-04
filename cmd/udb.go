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
	"crypto/rand"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/client/globals"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/csprng"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/primitives/ndf"
	"gitlab.com/elixxir/user-discovery-bot/udb"
	"math/big"
	"os"
	"strings"
	"time"
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
func StartBot(sess string, def *ndf.NetworkDefinition) error {
	udb.Log.DEBUG.Printf("Starting User Discovery Bot...")

	UDBSessionFileName = sess

	// Initialize the client
	regCode := udb.UDB_USERID.RegistrationCode()
	_, err := Init(UDBSessionFileName, regCode, def)
	if err != nil {
		return err
	}

	udb.Log.INFO.Printf("Logging in")

	// Log into the server with a blank password
	_, err = clientObj.Login("")
	if err != nil {
		return err
	}

	// get the newest message ID on the reception gateway to stop the UDB from
	// replaying old messages in the event of a redeploy where the session file
	// is lost
	lastMessageID, err := getLatestMessageID()
	if err != nil {
		return err
	}

	clientObj.GetSession().SetLastMessageID(lastMessageID)

	// Register the listeners with the user discovery bot
	udb.RegisterListeners(clientObj)

	udb.Log.INFO.Printf("Starting UDB")

	// starting the reception thread
	receiverCallback := func(err error) {
		if err != nil {
			udb.Log.ERROR.Printf("Start Message Reciever Callback Error: %v", err)
			backoff(clientObj, 0)
		}
	}

	err = clientObj.StartMessageReceiver(receiverCallback)
	if err != nil {
		return err
	}

	return nil
}

// Initialize a session using the given session file and other info
func Init(sessionFile string, regCode string, def *ndf.NetworkDefinition) (*id.User, error) {

	// We only register when the session file does not exist
	// FIXME: this is super weird -- why have to check for a file,
	// then init that file, then register optionally based on that check?
	_, err := os.Stat(sessionFile)
	// Get new client. Setting storage to nil internally creates a
	// default storage
	var initErr error

	if noTLS {
		//Set all tls certificates as empty effectively disabling tls
		for i := range def.Gateways {
			def.Gateways[i].TlsCertificate = ""
		}
		def.Registration.TlsCertificate = ""

		udb.Log.INFO.Printf("TURNING OFF TLS NOW, THESE ARE THE GATEWAYS %v, and this is Registration %v", def.Gateways, def.Registration)
	}

	secondarySessionFile := sessionFile + "-2"
	clientObj, initErr = api.NewClient(nil, sessionFile, secondarySessionFile, def)
	if initErr != nil {
		return nil, initErr
	}

	// API Settings (hard coded)
	clientObj.DisableBlockingTransmission() // Deprecated
	// Up to 10 messages per second
	clientObj.SetRateLimiting(uint32(RateLimit))

	// connect udb to gateways
	for {
		err = clientObj.InitNetwork()
		if err == nil {
			break
		}
		time.Sleep(10 * time.Second)
		udb.Log.ERROR.Printf("UDB could not connect to gateways, "+
			"reconnecting: %+v", err)
	}

	// SB: Trying to always register.
	// I think it's needed for some things to work correctly.
	// Need a more accurate descriptor of what the method actually does than
	// Register, or to remove the things that aren't actually used for
	// registration.
	//  RegisterWithPermissioning(preCan bool, registrationCode, nick, email,
	//	password string, privateKeyRSA *rsa.PrivateKey) (*id.User, error)
	userID, err := clientObj.RegisterWithPermissioning(true, regCode, "",
		"", "", nil)
	if err != nil {
		return nil, err
	}

	return userID, nil
}

// getLatestMessageID gets the newest message ID on the reception gateway, used
// to stop the UDB from replaying old messages in the event of a redeploy where
// the session file is lost
func getLatestMessageID() (string, error) {
	//get the newest message id to
	clientComms := clientObj.GetCommManager().Comms

	msg := &mixmessages.ClientRequest{
		UserID:        udb.UDB_USERID.Bytes(),
		LastMessageID: "",
	}

	receiveGateway := id.NewNodeFromBytes(clientObj.GetNDF().Nodes[len(clientObj.GetNDF().Gateways)-1].ID).NewGateway()

	var idList *mixmessages.IDList

	for {
		var err error
		host, ok := clientComms.GetHost(receiveGateway.String())
		if !ok {
			//ERROR getting host log it here
			//Needs to be part of a larger discussion for error handling
			return "", errors.Errorf("Failed to find the host with ID %v", receiveGateway.String())
		}

		idList, err = clientComms.SendCheckMessages(host, msg)
		if err != nil {
			globals.Log.WARN.Printf("Failed to get the latest message "+
				"IDs from the reception gateway: %s", err.Error())
			if strings.Contains(err.Error(),
				"Could not find any message IDs for this user") {
				break
			}
		} else {
			break
		}

		time.Sleep(2 * time.Second)
	}

	lastMessage := ""

	if idList != nil && idList.IDs != nil && len(idList.IDs) != 0 {
		lastMessage = idList.IDs[len(idList.IDs)-1]
	}

	globals.Log.INFO.Printf("Discarding messages before ID `%s`", lastMessage)

	return lastMessage, nil
}

//This is a recursive function used to restart startMessageReciever whenever it fails.
func backoff(cl *api.Client, backoffCount int) {
	receiverCallback := func(err error) {
		backoff(cl, backoffCount+1)
	}
	// Compute backoff time
	var delay time.Duration
	var block = false
	if backoffCount > 15 {
		delay = time.Hour
		block = true
	}
	wait := 2 ^ backoffCount
	if wait > 180 {
		wait = 180
	}
	jitter, _ := rand.Int(csprng.NewSystemRNG(), big.NewInt(1000))
	delay = time.Second*time.Duration(wait) + time.Millisecond*time.Duration(jitter.Int64())

	timer := time.NewTimer(delay)
	if block {
		timer.Stop()
	}
	select {
	case <-timer.C:
		backoffCount = 0
	}
	// attempt to start the message receiver
	err := cl.StartMessageReceiver(receiverCallback)
	if err != nil {
		udb.Log.ERROR.Printf("Start Message receiver failed %v", err)
		return
	}
}
