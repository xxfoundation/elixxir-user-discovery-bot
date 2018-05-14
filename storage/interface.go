////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Interface and enums for UDB storage systems
package storage

// The ValueType constant stores the allowable types we search on
// (e-mail, group, nickname, etc).
type ValueType int

// Note: because DB backends vary, and this list could have
//       items added and removed, we are not using an iota on purpose.
const (
	Email ValueType = 0 // An e-mail address
	Nick  ValueType = 1 // The user's nickname
	// TODO: Add more as necessary
)

type Storage interface {
	// AddKey - Add a key stream, return the fingerprint
	AddKey(value []byte) (string, error)
	// GetKey - Get a key based on the key id (retval of AddKey)
	GetKey(keyId string) ([]byte, bool)

	// NOTE: See doc on user discovery for these 2 functions.
	//       At this time we can't do a high-security version because
	//       we lack the anonymous return receipts.
	// AddUserKey - Add a user id to keyId (not used in high security)
	AddUserKey(userId uint64, keyId string) error
	// GetUserKey - Get a user's keyId (not used in high security)
	GetUserKey(userId uint64) (string, bool)

	// AddValue - Add a searchable value (e-mail, nickname, etc)
	AddValue(value string, valType ValueType, keyId string) error
	// GetKeys - Returns all values that match the search criteria
	GetKeys(value string, valType ValueType) ([]string, bool)
}

// Print strings for ValueType
func (v ValueType) String() string {
	names := [...]string{
		"Email",
		"Nick",
	}
	if v < Email || v > Nick {
		return "Unknown"
	}
	return names[v]
}
