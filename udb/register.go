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
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/privategrity/user-discovery-bot/storage"
	"strconv"
)

const REGISTER_USAGE = "Usage: 'REGISTER [EMAIL] [email-address] " +
	"[key-fingerprint]'"

// Add a user to the registry
// The register command takes the form "REGISTER TYPE VALUE KEYID",
// WHERE:
//  - TYPE = EMAIL (and later others, maybe)
//  - VALUE = "rick@privategrity.com"
//  - KEYFP = the key fingerprint
//
// The user ID is taken from the sender at this time, this will need to change
// when a registrar comes online.
// Registration fails if the KEYID is not already pushed and confirmed.
func Register(userId uint64, args []string) {
	jww.INFO.Printf("Register %d: %v", userId, args)
	RegErr := func(msg string) {
		Send(userId, msg)
		Send(userId, REGISTER_USAGE)
		jww.INFO.Printf("User %d error: %s", userId, msg)
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

	jww.INFO.Printf("User %d registered successfully with %s, %s",
		userId, regVal, keyFp)
	Send(userId, "REGISTRATION COMPLETE")
}

const PUSHKEY_USAGE = "Usage: 'PUSHKEY [temp-key-id] [starting-byte-index] " +
				"[base64-encoded-bytestream]'"
const PUSHKEY_SIZE = 256 // 2048 bits
var tempKeys = make(map[string][]byte)
var tempKeysState = make(map[string][]bool)

// PushKey adds a key to the registration database and links it by fingerprint
// The PushKey command has the form PUSHKEY KEYID IDX KEYMAT
// WHERE:
//  - KEYID = The Key ID -- not the same as the fingerprint
//  - IDX = byte index of the key
//  - KEYMAT = The part of the key corresponding to that index, in BASE64
// PushKey returns an ACK that it received the command OR a success/failure
// once it receives all pieces of the key.
func PushKey(userId uint64, args []string) {
	jww.INFO.Printf("PushKey %d:, %v", userId, args)
	PushErr := func(msg string) {
		Send(userId, msg)
		Send(userId, PUSHKEY_USAGE)
		jww.INFO.Printf("User %d error: %s", userId, msg)
	}
	if len(args) != 3 {
		PushErr("Invalid command syntax!")
		return
	}

	keyId := args[0]
	keyIdxStr := args[1]
	keyMat := args[2]

	// Decode keyMat
	// FIXME: Not sure I like having to base64 stuff here, but it's this or hex
	// Maybe add suppor to client for these pubkey conversions?
	newKeyBytes, decErr := base64.StdEncoding.DecodeString(keyMat)
	if decErr != nil {
		PushErr("Could not decode new key bytes, it must be in base64!")
		return
	}

	// Parse index
	keyIdx, pErr := strconv.Atoi(keyIdxStr)
	if pErr != nil || keyIdx < 0 || (keyIdx+len(newKeyBytes)) > PUSHKEY_SIZE {
		PushErr("Invalid key index!")
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
		j := keyIdx + i
		key[j] = newKeyBytes[i]
		keyState[j] = true
	}

	// Calculate how many more bytes are needed
	missingCnt := 0
	for i := 0; i < PUSHKEY_SIZE; i++ {
		if !keyState[i] {
			missingCnt += 1
		}
	}

	if missingCnt != 0 {
		tempKeys[keyId] = key
		tempKeysState[keyId] = keyState
		msg := fmt.Sprintf("PUSHKEY ACK NEED %d", missingCnt)
		jww.INFO.Printf("User %d: %s", userId, msg)
		Send(userId, msg)
		return
	}

	// Add key and remove from temporary storage
	delete(tempKeys, keyId)
	delete(tempKeysState, keyId)
	fingerprint, err := DataStore.AddKey(key)
	if err != nil {
		PushErr(err.Error())
	}
	msg := fmt.Sprintf("PUSHKEY COMPLETE %s", fingerprint)
	jww.INFO.Printf("User %d: %s", userId, msg)
	Send(userId, msg)
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
func GetKey(userId uint64, args []string) {
	jww.INFO.Printf("GetKey %d:, %v", userId, args)
	GetErr := func(msg string) {
		Send(userId, msg)
		Send(userId, GETKEY_USAGE)
		jww.INFO.Printf("User %d error: %s", userId, msg)
	}
	if len(args) != 1 {
		GetErr("Invalid command syntax!")
		return
	}

	keyFp := args[0]

	key, ok := DataStore.GetKey(keyFp)
	if !ok {
		msg := fmt.Sprintf("GETKEY %s NOTFOUND", keyFp)
		jww.INFO.Printf("UserId %d: %s", userId, msg)
		Send(userId, msg)
		return
	}

	keymat := base64.StdEncoding.EncodeToString(key)
	msg := fmt.Sprintf("GETKEY %s %s", keyFp, keymat)
	jww.INFO.Printf("UserId %d: %s", userId, msg)
	Send(userId, msg)
}
