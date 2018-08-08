////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package storage

import (
	"gitlab.com/privategrity/client/user"
	"testing"
)

func TestRamAddAndGetKey(t *testing.T) {
	RS := NewRamStorage()
	testKey := []byte{'a', 'b', 'c', 'd'}

	// Add K to the Store
	fingerprint, err := RS.AddKey(testKey)
	if err != nil {
		t.Errorf("Ram storage error on AddKey: %v", err)
	}
	// Verify it made it
	retKey, ok := RS.GetKey(fingerprint)
	if !ok {
		t.Errorf("Ram storage error GetKey: %v", err)
	}
	for i := range testKey {
		if retKey[i] != testKey[i] {
			t.Errorf("Ram Storage cannot store and load keys at index %d - "+
				"Expected: %v, Got: %v", i, testKey[i], retKey[i])
		}
	}

	// Now try to add it again and verify it fails
	_, err3 := RS.AddKey(testKey)
	if err3 == nil {
		t.Errorf("Ram storage AddKey allows duplicates!")
	}

	// Now check missing key
	_, ok2 := RS.GetKey("BlahThisIsABadKey")
	if ok2 {
		t.Errorf("Ram storage GetKey returns results on bad keys!")
	}
}

func TestRamAddAndGetUserKey(t *testing.T) {
	RS := NewRamStorage()
	keyId := "This is my keyId"
	userId := user.ID(1337)
	// Add key
	err := RS.AddUserKey(userId, keyId)
	if err != nil {
		t.Errorf("Ram storage AddUserKey failed to add a user: %v", err)
	}
	// Add duplicate
	err2 := RS.AddUserKey(userId, keyId)
	if err2 == nil {
		t.Errorf("Ram storage AddUserKey permits duplicates!")
	}

	// Get Key
	retrievedKeyId, ok := RS.GetUserKey(userId)
	if !ok {
		t.Errorf("Ram storage GetUserKey could not retrieve key!")
	}
	if retrievedKeyId != keyId {
		t.Errorf("Ram storage GetUserKey failed - Got: %s, Expected: %s",
			retrievedKeyId, keyId)
	}
}

func TestValueAndKeyStore(t *testing.T) {
	RS := NewRamStorage()
	value := "Hello, World!"
	KeyId := "This is a key id"
	err := RS.AddValue(value, Email, KeyId)
	if err != nil {
		t.Errorf("Ram storage could not AddValue!")
	}

	retKeys, ok := RS.GetKeys(value, Email)
	if !ok {
		t.Errorf("Ram storage could not GetKeys!")
	}
	if retKeys[0] != KeyId {
		t.Errorf("Ram storage GetKeys returned bad result - Got: %s, Expected: %s",
			retKeys, KeyId)
	}

	// check for empty value
	_, ok2 := RS.GetKeys("junk value", Email)
	if ok2 {
		t.Errorf("Ram storage GetKeys returned on junk input!")
	}
}
