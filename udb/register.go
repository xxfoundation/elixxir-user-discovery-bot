////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Registration Commands (Register, PushKey, GetKey)
package udb

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"gitlab.com/elixxir/client/cmixproto"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/user-discovery-bot/fingerprint"
	"gitlab.com/elixxir/user-discovery-bot/storage"
)

const REGISTER_USAGE = "Usage: 'REGISTER [EMAIL] [email-address] " +
	"[key-fingerprint]'"

// Add a user to the registry
// The register command takes the form "REGISTER TYPE VALUE KEYID",
// WHERE:
//  - TYPE = EMAIL (and later others, maybe)
//  - VALUE = "rick@elixxir.io"
//  - KEYFP = the key fingerprint
//
// The user ID is taken from the sender at this time, this will need to change
// when a registrar comes online.
// Registration fails if the KEYID is not already pushed and confirmed.
func Register(userId *id.User, args []string) {
	Log.DEBUG.Printf("Register %d: %v", userId, args)
	RegErr := func(msg string) {
		Send(userId, msg, cmixproto.Type_UDB_REGISTER_RESPONSE)
		Send(userId, REGISTER_USAGE, cmixproto.Type_UDB_REGISTER_RESPONSE)
		Log.INFO.Printf("Register user %d error: %s", userId, msg)
	}
	if len(args) != 3 {
		RegErr("Invalid command syntax!")
		return
	}

	regType := args[0]
	regVal := args[1]
	keyFp := args[2]

	// Verify that regType == EMAIL
	if regType != "EMAIL" {
		RegErr("EMAIL is the only acceptable registration type")
		return
	}
	// TODO: Add parse func to storage class, embed into function and
	//  pass it a string instead

	// Verify the key is accounted for
	retrievedUser, err := storage.UserDiscoveryDb.GetUserByKeyId(keyFp)
	if err != nil {
		msg := fmt.Sprintf("Could not find keyFp: %s", keyFp)
		RegErr(msg)
		return
	}

	//Check that the email has not been registered before
	_, err = storage.UserDiscoveryDb.GetUserByValue(regVal)
	if err==nil{
		msg := fmt.Sprintf("Can not register with existing email: %s", regVal)
		RegErr(msg)
	}

	//FIXME: Do you want to do these checks? My guess is no, but that's how it was done previously
	// W/o checks, you could get someone trying to overwrite someone else's account (? maybe?)
	//Check that the retrieved user's attributes have been set
	if bytes.Compare(retrievedUser.Id, make([]byte, 0)) != 0 {
		RegErr(fmt.Sprintf("UserId already exists: %d", retrievedUser.Id))
	} else {
		retrievedUser.SetID(userId.Bytes())
	}
	if retrievedUser.Value != "" {
		RegErr(fmt.Sprintf("email already exists: %s", retrievedUser.Value))
	} else {
		retrievedUser.SetValue(regVal)
	}

	//FIXME: Hardcoded to email value, change later
	retrievedUser.SetValueType(0)
	err = storage.UserDiscoveryDb.UpsertUser(retrievedUser)

	if err != nil {
		RegErr(err.Error())
	}

	Log.INFO.Printf("User %v registered successfully with %s, %s",
		*userId, regVal, keyFp)
	Send(userId, "REGISTRATION COMPLETE", cmixproto.Type_UDB_REGISTER_RESPONSE)
}

const PUSHKEY_USAGE = "Usage: 'PUSHKEY [temp-key-id] " +
	"[base64-encoded-bytestream]'"

// PushKey adds a key to the registration database and links it by fingerprint
// The PushKey command has the form PUSHKEY KEYID KEYMAT
// WHERE:
//  - KEYID = The Key ID -- not necessarily the same as the fingerprint
//  - KEYMAT = The part of the key corresponding to that index, in BASE64
// PushKey returns an ACK that it received the command OR a success/failure
// once it receives all pieces of the key.
func PushKey(userId *id.User, args []string) {
	Log.DEBUG.Printf("PushKey %d, %v", userId, args)
	PushErr := func(msg string) {
		Send(userId, msg, cmixproto.Type_UDB_PUSH_KEY_RESPONSE)
		Send(userId, PUSHKEY_USAGE, cmixproto.Type_UDB_PUSH_KEY_RESPONSE)
		Log.INFO.Printf("PushKey user %d error: %s", userId, msg)
	}
	if len(args) != 2 {
		PushErr("Invalid command syntax!")
		return
	}

	// keyId := args[0] Note: Legacy, key id is not needed anymore as it is
	//                        sent as a single message
	keyMat := args[1]
	// Decode keyMat
	// FIXME: Not sure I like having to base64 stuff here, but it's this or hex
	//  Maybe add support to client for these pubkey conversions?
	newKeyBytes, decErr := base64.StdEncoding.DecodeString(keyMat)
	if decErr != nil {
		PushErr(fmt.Sprintf("Could not decode new key bytes, "+
			"it must be in base64! %s", decErr))
		return
	}
	usr := storage.NewUser()
	usr.SetKey(newKeyBytes)
	keyFP := fingerprint.Fingerprint(newKeyBytes)
	fmt.Println()
	usr.SetKeyID(keyFP)
	_ = storage.UserDiscoveryDb.UpsertUser(usr)
	msg := fmt.Sprintf("PUSHKEY COMPLETE %s", keyFP)
	Log.DEBUG.Printf("User %d: %s", userId, msg)
	Send(userId, msg, cmixproto.Type_UDB_PUSH_KEY_RESPONSE)
}

const GETKEY_USAGE = "GETKEY [KEYFP]"

// GetKey retrieves a key based on its fingerprint
// The GetKey command has the form GETKEY KEYFP
// WHERE:
//  - KEYFP - The Key Fingerprint
// GetKey returns KEYFP IDX KEYMAT, where:
//  - KEYFP - The Key Fingerprint
//  - KEYMAT - Key material in BASE64 encoding
// It sends these messages until the entire key is transmitted.
func GetKey(userId *id.User, args []string) {
	Log.DEBUG.Printf("GetKey %d:, %v", userId, args)
	GetErr := func(msg string) {
		Send(userId, msg, cmixproto.Type_UDB_GET_KEY_RESPONSE)
		Send(userId, GETKEY_USAGE, cmixproto.Type_UDB_GET_KEY_RESPONSE)
		Log.INFO.Printf("User %d error: %s", userId, msg)
	}
	if len(args) != 1 {
		GetErr("Invalid command syntax!")
		return
	}

	keyFp := args[0]
	retrievedUser, err := storage.UserDiscoveryDb.GetUserByKeyId(keyFp)
	if err != nil {
		msg := fmt.Sprintf("GETKEY %s NOTFOUND", keyFp)
		Log.INFO.Printf("UserId %d: %s", userId, msg)
		Send(userId, msg, cmixproto.Type_UDB_GET_KEY_RESPONSE)
		return
	}

	keymat := base64.StdEncoding.EncodeToString(retrievedUser.Key)
	msg := fmt.Sprintf("GETKEY %s %s", keyFp, keymat)
	Log.DEBUG.Printf("UserId %d: %s", userId, msg)
	Send(userId, msg, cmixproto.Type_UDB_GET_KEY_RESPONSE)
}
