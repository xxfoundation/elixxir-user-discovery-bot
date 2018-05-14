////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Registration Commands (Register, PushKey, GetKey)
package udb

import (
	jww "github.com/spf13/jwalterweatherman"
)

// Add a user to the registry
// The register command takes the form "REGISTER TYPE VALUE KEYID",
// WHERE:
//  - TYPE = EMAIL (and later others, maybe)
//  - VALUE = "rick@privategrity.com"
//  - KEYID = the key fingerprint
//
// The user ID is taken from the sender at this time, this will need to change
// when a registrar comes online.
// Registration fails if the KEYID is not already pushed and confirmed.
func Register(userId uint64, args []string) {
	jww.INFO.Printf("Register: %d, %v", userId, args)
}

// PushKey adds a key to the registration database and links it by fingerprint
// The PushKey command has the form PUSHKEY KEYID IDX KEYMAT
// WHERE:
//  - KEYID = The Key ID
//  - IDX = byte index of the key
//  - KEYMAT = The part of the key corresponding to that index, in BASE64
// PushKey returns and ACK that it received the command OR a success/failure
// once it receives all pieces of the key.
func PushKey(userId uint64, args []string) {
	jww.INFO.Printf("PushKey: %d, %v", userId, args)
}
