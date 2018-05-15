////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package udb

import (
	"encoding/base64"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/privategrity/user-discovery-bot/storage"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	jww.SetStdoutThreshold(jww.LevelDebug)
	os.Exit(m.Run())
}

// Push the key then register
// NOTE: The send function defaults to a no-op when client is not set up. I am
//       not sure how I feel about it.
func TestRegisterHappyPath(t *testing.T) {
	DataStore = storage.NewRamStorage()
	pubKeyBits := []string{
		"S8KXBczy0jins9uS4LgBPt0bkFl8t00MnZmExQ6GcOcu8O7DKgAsNz" +
			"LU7a+gMTbIsS995IL/kuFF8wcBaQJBY23095PMSQ/nMuetzhk9HdXxrGIiKBo3C/n4SClp" +
			"q4H+PoF9XziEVKua8JxGM2o83KiCK3tNUpaZbAAElkjueY4=",
		"8Lg/eoeKGgPlleTYfO3JyGfnwBtLi73ti0h2dBQWW94JTqTQDr+z" +
			"xVpLzdgTt+87TkAl0yXu9mOUXqGJ+51lTcRlIdIpWpfgUbibdRme8IThg0RNCF31ESKCts" +
			"o8gJ8mSVljIXxrC+Uuoi+Gl1LNN5nPARykatx0Y70xNdJd2BQ=",
	}
	pubKey := make([]byte, 256)
	for i := range pubKeyBits {
		bytes, _ := base64.StdEncoding.DecodeString(pubKeyBits[i])
		for j := range bytes {
			pubKey[j+i*128] = bytes[j]
		}
	}

	fingerprint := "8oKh7TYG4KxQcBAymoXPBHSD/uga9pX3Mn/jKhvcD8M="
	msgs := []string{
		"PUSHKEY myKeyId 0 " + pubKeyBits[0],
		"PUSHKEY myKeyId 128 " + pubKeyBits[1],
		"REGISTER EMAIL rick@privategrity.com " + fingerprint,
		"GETKEY " + fingerprint,
	}

	for i := range msgs {
		msg, err := NewMessage(msgs[i])
		if err != nil {
			t.Errorf("Error generating message: %v", err)
		}
		ReceiveMessage(msg)
	}

	// Assert expected state
	k, ok := DataStore.GetKey(fingerprint)
	if !ok {
		t.Errorf("Could not retrieve key %s", fingerprint)
	}
	for i := range k {
		if k[i] != pubKey[i] {
			t.Errorf("pubKey byte mismatch at %d: %d v %d", i, k[i], pubKey[i])
		}
	}

	u, ok2 := DataStore.GetUserKey(uint64(1))
	if !ok2 {
		t.Errorf("Could not retriever user key 1!")
	}
	if u != fingerprint {
		t.Errorf("GetUserKey fingerprint mismatch: %s v %s", u, fingerprint)
	}

	ks, ok3 := DataStore.GetKeys("rick@privategrity.com", storage.Email)
	if !ok3 {
		t.Errorf("Could not retrieve by e-mail address!")
	}
	if ks[0] != fingerprint {
		t.Errorf("GetKeys fingerprint mismatch: %v v %s", ks, fingerprint)
	}
}

func TestInvalidRegistrationCommands(t *testing.T) {
	DataStore = storage.NewRamStorage()
	msgs := []string{
		"PUSHKEY garbage doiandga daoinaosf adsoifn dsaoifa",
		"REGISTER NOTEMAIL something something",
		"REGISTER EMAIL garbage this is a garbage",
		"REGISTER EMAIL rick@privategrity 8oKh7TYG4KxQcBAymoXPBHSD/uga9pX3Mn/jKh" +
			"vcD8M=",
	}

	for i := range msgs {
		msg, err := NewMessage(msgs[i])
		if err != nil {
			t.Errorf("Error generating message: %v", err)
		}
		ReceiveMessage(msg)
		_, ok := DataStore.GetKey("8oKh7TYG4KxQcBAymoXPBHSD/uga9pX3Mn/jKh")
		if ok {
			t.Errorf("Data store key 8oKh7TYG4KxQcBAymoXPBHSD/uga9pX3Mn/jKh should" +
				" not exist!")
		}
		_, ok2 := DataStore.GetUserKey(uint64(1))
		if ok2 {
			t.Errorf("Data store user 1 should not exist!")
		}
		_, ok3 := DataStore.GetKeys("rick@privategrity.com", storage.Email)
		if ok3 {
			t.Errorf("Data store value rick@privategrity.com should not exist!")
		}
	}
}
