////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package cmix

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"gitlab.com/elixxir/client/v4/single"
	"gitlab.com/elixxir/client/v4/ud"
	"gitlab.com/elixxir/crypto/contact"
	"gitlab.com/elixxir/crypto/cyclic"
	"gitlab.com/elixxir/crypto/diffieHellman"
	"gitlab.com/elixxir/crypto/fastRNG"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/crypto/csprng"
	"gitlab.com/xx_network/crypto/large"
	"gitlab.com/xx_network/primitives/id"
	"testing"
	"time"
)

func TestManager_lookupCallback(t *testing.T) {
	// Initialize manager
	m := &Manager{db: storage.NewTestDB(t),
		lookupListener: &mockSingleLookup{}}
	uid := id.NewIdFromString("zezima", id.User, t)
	grp := cyclic.NewGroup(large.NewInt(107), large.NewInt(2))

	rng := fastRNG.NewStreamGenerator(12, 1024, csprng.NewSystemRNG).GetStream()
	defer rng.Close()

	dhPriv := diffieHellman.GeneratePrivateKey(128, grp, rng)
	dhPub := diffieHellman.GeneratePublicKey(dhPriv, grp)

	lm := &lookupManager{m: m}

	// Insert mock user into DB
	err := m.db.InsertUser(&storage.User{Id: uid.Marshal(), DhPub: dhPub.Bytes()})
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}

	mockCmix := newMockCmix(uid, newMockCmixHandler(), t)

	// Set up the listener
	single.Listen(ud.LookupTag, uid, dhPriv, mockCmix, grp, lm)

	// Build the request
	lookupMsg := &ud.LookupSend{UserID: uid.Marshal()}
	payload, err := proto.Marshal(lookupMsg)
	if err != nil {
		t.Errorf("Failed to marshal payload: %+v", err)
	}

	// Build the contact
	userContact := contact.Contact{
		ID:       uid,
		DhPubKey: dhPub,
	}

	// Set up the response handling
	callbackChan := make(chan []byte)
	mockResponse := newMockResponse(callbackChan, t)

	// Trigger the request
	_, _, err = single.TransmitRequest(userContact, ud.LookupTag, payload,
		mockResponse, single.GetDefaultRequestParams(), mockCmix, rng, grp)
	if err != nil {
		t.Fatalf("Failed to transmit mock message: %v", err)
	}

	// Build expected response
	expectedResponse := &ud.LookupResponse{
		PubKey: dhPub.Bytes(),
	}

	// Marshal expected response
	expected, err := proto.Marshal(expectedResponse)
	if err != nil {
		t.Fatalf("Marshal response: %v", err)
	}

	// Wait for response or timeout
	timeout := time.NewTicker(5 * time.Second)
	select {
	case response := <-callbackChan:
		if !bytes.Equal(expected, response) {
			t.Errorf("Did not receive expected response."+
				"\nExpected: %s"+
				"\nReceived: %s", expected, response)
		}
	case <-timeout.C:
		t.Fatalf("Failed to get response")
	}

}

// mockSingleLookup is used to test the lookup function, which uses the single-
// use manager. It adheres to the SingleInterface interface.
type mockSingleLookup struct {
}

func (s *mockSingleLookup) Stop() {

}
