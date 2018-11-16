////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Registration Commands (Register, PushKey, GetKey)
package udb

import (
	"encoding/base64"
	"fmt"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/elixxir/crypto/id"
	"gitlab.com/elixxir/client/cmixproto"
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
func Register(userId *id.UserID, args []string) {
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
	// pass it a string instead
	regTypeEnum := storage.Email

	// Verify the key is accounted for
	_, ok := DataStore.GetKey(keyFp)
	if !ok {
		msg := fmt.Sprintf("Could not find keyFp: %s", keyFp)
		RegErr(msg)
		return
	}

	err := DataStore.AddUserKey(userId, keyFp)
	if err != nil {
		RegErr(err.Error())
	}
	err2 := DataStore.AddValue(regVal, regTypeEnum, keyFp)
	if err2 != nil {
		RegErr(err.Error())
	}

	Log.INFO.Printf("User %d registered successfully with %s, %s",
		userId, regVal, keyFp)
	Send(userId, "REGISTRATION COMPLETE", cmixproto.Type_UDB_REGISTER_RESPONSE)
}

const PUSHKEY_USAGE = "Usage: 'PUSHKEY [temp-key-id] " +
	"[base64-encoded-bytestream]'"
const PUSHKEY_SIZE = 256 // 2048 bits
var tempKeys = make(map[string][]byte)
var tempKeysState = make(map[string][]bool)

// PushKey adds a key to the registration database and links it by fingerprint
// The PushKey command has the form PUSHKEY KEYID KEYMAT
// WHERE:
//  - KEYID = The Key ID -- not necessarily the same as the fingerprint
//  - KEYMAT = The part of the key corresponding to that index, in BASE64
// PushKey returns an ACK that it received the command OR a success/failure
// once it receives all pieces of the key.
func PushKey(userId *id.UserID, args []string) {
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

	keyId := args[0]
	keyMat := args[1]

	// Decode keyMat
	// FIXME: Not sure I like having to base64 stuff here, but it's this or hex
	// Maybe add suppor to client for these pubkey conversions?
	newKeyBytes, decErr := base64.StdEncoding.DecodeString(keyMat)
	if decErr != nil {
		PushErr(fmt.Sprintf("Could not decode new key bytes, "+
			"it must be in base64! %s", decErr))
		return
	}

	// Does it exist yet?
	key, ok := tempKeys[keyId]
	keyState, _ := tempKeysState[keyId]
	if !ok {
		key = make([]byte, PUSHKEY_SIZE)
		keyState = make([]bool, PUSHKEY_SIZE)
	}

	// Update temporary storage
	for i := range newKeyBytes {
		key[i] = newKeyBytes[i]
		keyState[i] = true
	}

	// Add key and remove from temporary storage
	delete(tempKeys, keyId)
	delete(tempKeysState, keyId)
	fingerprint, err := DataStore.AddKey(key)
	if err != nil {
		PushErr(err.Error())
	}
	msg := fmt.Sprintf("PUSHKEY COMPLETE %s", fingerprint)
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
func GetKey(userId *id.UserID, args []string) {
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

	key, ok := DataStore.GetKey(keyFp)
	if !ok {
		msg := fmt.Sprintf("GETKEY %s NOTFOUND", keyFp)
		Log.INFO.Printf("UserId %d: %s", userId, msg)
		Send(userId, msg, cmixproto.Type_UDB_GET_KEY_RESPONSE)
		return
	}

	keymat := base64.StdEncoding.EncodeToString(key)
	msg := fmt.Sprintf("GETKEY %s %s", keyFp, keymat)
	Log.DEBUG.Printf("UserId %d: %s", userId, msg)
	Send(userId, msg, cmixproto.Type_UDB_GET_KEY_RESPONSE)
}
