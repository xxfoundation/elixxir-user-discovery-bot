////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Ram-only implementation of storage type
package storage

import (
	"fmt"
	"gitlab.com/elixxir/user-discovery-bot/fingerprint"
	"gitlab.com/elixxir/primitives/userid"
)

type RamStorage struct {
	Keys   map[string][]byte                 // keyId -> publicKey
	Users  map[id.UserID]string                // cMix UID -> keyId
	KeyVal map[ValueType]map[string][]string // ValType -> search string -> keyIds
}

// Create a blank ram storage object
func NewRamStorage() *RamStorage {
	RS := RamStorage{
		Keys:   make(map[string][]byte),
		Users:  make(map[id.UserID]string),
		KeyVal: make(map[ValueType]map[string][]string),
	}
	// NOTE: We could init all the KeyVal maps here, but I
	// decided to leave that to the AddValue function in favor of
	// not needing to touch this file when we add types.
	return &RS
}

// Addkey - Add a key stream, return the fingerprint
func (rs RamStorage) AddKey(value []byte) (string, error) {
	keyFingerprint := fingerprint.Fingerprint(value)

	// Error out if the key exists already
	_, ok := rs.Keys[keyFingerprint]
	if ok {
		return "", fmt.Errorf("Fingerprint already exists: %s", keyFingerprint)
	}

	rs.Keys[keyFingerprint] = value
	return keyFingerprint, nil
}

// GetKey - Get a key based on the key id (retval of AddKey)
func (rs RamStorage) GetKey(keyId string) ([]byte, bool) {
	publicKey, ok := rs.Keys[keyId]
	return publicKey, ok
}

// AddUserKey - Add a user id to keyId (not used in high security)
func (rs RamStorage) AddUserKey(userId *id.UserID, keyId string) error {
	_, ok := rs.Users[*userId]
	if ok {
		return fmt.Errorf("UserId already exists: %d", userId)
	}
	rs.Users[*userId] = keyId
	return nil
}

// GetUserKey - Get a user's keyId (not used in high security)
func (rs RamStorage) GetUserKey(userId *id.UserID) (string, bool) {
	keyId, ok := rs.Users[*userId]
	return keyId, ok
}

// AddValue - Add a searchable value (e-mail, nickname, etc)
func (rs RamStorage) AddValue(value string, valType ValueType,
	keyId string) error {
	_, ok := rs.KeyVal[valType]
	if !ok {
		rs.KeyVal[valType] = make(map[string][]string)
	}
	_, ok = rs.KeyVal[valType][value]
	if !ok {
		rs.KeyVal[valType][value] = make([]string, 0)
	}
	keyIds, _ := rs.KeyVal[valType][value]
	keyIds = append(keyIds, keyId)
	rs.KeyVal[valType][value] = keyIds
	return nil
}

// GetKeys - Returns all values that match the search criteria
func (rs RamStorage) GetKeys(value string, valType ValueType) (
	[]string, bool) {
	_, ok := rs.KeyVal[valType]
	if ok {
		keyIds, ok := rs.KeyVal[valType][value]
		if ok {
			return keyIds, ok
		}
	}
	return nil, false
}
