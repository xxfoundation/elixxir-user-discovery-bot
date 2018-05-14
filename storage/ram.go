////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Ram-only implementation of storage type
package storage

import (
	"fmt"
	"gitlab.com/privategrity/user-discovery-bot/fingerprint"
)

type RamStorage struct {
	Keys map[string][]byte // keyId -> publicKey
	Users map[uint64]string // cMix UID -> keyId
	KeyVal map[ValueType]map[string][]string // ValType -> search string -> keyIds
}

// Create a blank ram storage object
func NewRamStorage() *RamStorage {
	RS := RamStorage{
		Keys: make(map[string][]byte),
		Users: make(map[uint64]string),
		KeyVal: make(map[ValueType]map[string][]string),
	}
	// NOTE: We could init all the KeyVal maps here, but I
	// decided to leave that to the AddValue function in favor of
	// not needing to touch this file when we add types.
	return &RS
}

// Addkey - Add a key stream, return the fingerprint
func (RamStore RamStorage) AddKey(value []byte) (string, error) {
	keyFingerprint := fingerprint.Fingerprint(value)

	// Error out if the key exists already
	_, ok := RamStore.Keys[keyFingerprint]
	if ok {
		return "", fmt.Errorf("Fingerprint already exists: %s", keyFingerprint)
	}

	RamStore.Keys[keyFingerprint] = value
	return keyFingerprint, nil
}

// GetKey - Get a key based on the key id (retval of AddKey)
func (RamStore RamStorage) GetKey(keyId string) ([]byte, bool) {
	publicKey, ok := RamStore.Keys[keyId]
	return publicKey, ok
}

// AddUserKey - Add a user id to keyId (not used in high security)
func (RamStore RamStorage) AddUserKey(userId uint64, keyId string) error {
	_, ok := RamStore.Users[userId]
	if ok {
		return fmt.Errorf("UserId already exists: %d", userId)
	}
	RamStore.Users[userId] = keyId
	return nil
}

// GetUserKey - Get a user's keyId (not used in high security)
func (RamStore RamStorage) GetUserKey(userId uint64) (string, bool) {
	keyId, ok := RamStore.Users[userId]
	return keyId, ok
}

// AddValue - Add a searchable value (e-mail, nickname, etc)
func (RamStore RamStorage) AddValue(value string, valType ValueType,
	keyId string) error {
	_, ok := RamStore.KeyVal[valType]
	if ! ok {
		RamStore.KeyVal[valType] = make(map[string][]string)
	}
	_, ok = RamStore.KeyVal[valType][value]
	if ! ok {
		RamStore.KeyVal[valType][value] = make([]string, 1)
	}
	keyIds, _ := RamStore.KeyVal[valType][value]
	keyIds = append(keyIds, keyId)
	RamStore.KeyVal[valType][value] = keyIds
	return nil
}

// GetKeys - Returns all values that match the search criteria
func (RamStore RamStorage) GetKeys(value string, valType ValueType) (
	[]string, bool) {
	_, ok := RamStore.KeyVal[valType]
	if ok {
		keyIds, ok := RamStore.KeyVal[valType][value]
		if ok {
			return keyIds, ok
		}
	}
	return nil, false
}
