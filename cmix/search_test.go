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

func TestManager_SearchProcess(t *testing.T) {
	m := &Manager{db: storage.NewTestDB(t),
		searchListener: &mockSingleLookup{}}

	uid := id.NewIdFromString("zezima", id.User, t)
	grp := cyclic.NewGroup(large.NewInt(107), large.NewInt(2))

	rng := fastRNG.NewStreamGenerator(12, 1024, csprng.NewSystemRNG).GetStream()
	defer rng.Close()

	dhPriv := diffieHellman.GeneratePrivateKey(128, grp, rng)
	dhPub := diffieHellman.GeneratePublicKey(dhPriv, grp)

	// Insert mock user into DB
	err := m.db.InsertUser(&storage.User{Id: uid.Marshal(), DhPub: dhPub.Bytes()})
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}

	// Insert mock fact into DB
	fid := []byte("factHash")
	err = m.db.InsertFact(&storage.Fact{
		Hash:     fid,
		UserId:   uid.Marshal(),
		Fact:     "water is wet",
		Type:     0,
		Verified: true,
	})
	if err != nil {
		t.Errorf("Failed to insert dummy fact: %+v", err)
	}

	mockCmix := newMockCmix(uid, newMockCmixHandler(), t)

	sm := &searchManager{m: m}
	single.Listen(ud.SearchTag, uid, dhPriv, mockCmix, grp, sm)

	// Build search request
	searchMsg := &ud.SearchSend{Fact: []*ud.HashFact{{Hash: fid, Type: 0}}}
	payload, err := proto.Marshal(searchMsg)
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
	_, _, err = single.TransmitRequest(userContact, ud.SearchTag, payload,
		mockResponse, single.GetDefaultRequestParams(), mockCmix, rng, grp)
	if err != nil {
		t.Fatalf("Failed to transmit mock message: %v", err)
	}

	// Build expected response
	expectedResponse := &ud.SearchResponse{
		Contacts: []*ud.Contact{
			{
				UserID: uid.Bytes(),
				PubKey: dhPub.Bytes(),
				TrigFacts: []*ud.HashFact{
					{
						Hash: fid,
					},
				},
			},
		},
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

//func TestManager_handleSearch(t *testing.T) {
//	m := &Manager{db: storage.NewTestDB(t)}
//	uid := id.NewIdFromString("zezima", id.User, t)
//	c := single.NewContact(uid, &cyclic.Int{}, &cyclic.Int{}, singleUse.TagFP{}, 8)
//
//	expectedDhPub := []byte("DhPub")
//	fid := []byte("factHash")
//	f1 := &storage.Fact{
//		Hash:     fid,
//		UserId:   uid.Marshal(),
//		Fact:     "water is wet",
//		Type:     1,
//		Verified: true,
//	}
//	expectedContact := &ud.Contact{
//		UserID:    uid.Marshal(),
//		PubKey:    expectedDhPub,
//		TrigFacts: []*ud.HashFact{{Hash: f1.Hash, Type: int32(f1.Type)}},
//	}
//	err := m.db.InsertUser(&storage.User{Id: uid.Marshal(), DhPub: expectedDhPub})
//	if err != nil {
//		t.Errorf("Failed to insert dummy user: %+v", err)
//	}
//
//	err = m.db.InsertFact(f1)
//	if err != nil {
//		t.Errorf("Failed to insert dummy fact: %+v", err)
//	}
//
//	resp := m.handleSearch(&ud.SearchSend{Fact: []*ud.HashFact{{Hash: fid, Type: 0}}}, c)
//	if resp.Error != "" {
//		t.Errorf("Failed to handle search: %+v", resp.Error)
//	}
//	if len(resp.Contacts) != 1 {
//		t.Errorf("Did not receive expected number of contacts")
//	}
//	if !reflect.DeepEqual(expectedContact, resp.Contacts[0]) {
//		t.Errorf("Did not received expected contact."+
//			"\nexpected: %+v\nreceived: %+v", expectedContact, resp.Contacts[0])
//	}
//
//	resp = m.handleSearch(&ud.SearchSend{Fact: []*ud.HashFact{{Hash: fid, Type: int32(fact.Nickname)}}}, c)
//	if resp.Error == "" {
//		t.Errorf("Search should have returned error")
//	}
//	if len(resp.Contacts) != 0 {
//		t.Errorf("Should not be able to search with nickname")
//	}
//}

//// mockSingleSearch is used to test the search function, which uses the single-
//// use manager. It adheres to the SingleInterface interface.
//type mockSingleSearch struct {
//	callback func(payload []byte, c single.Contact)
//}
//
//func (s *mockSingleSearch) RegisterCallback(_ string, callback single.ReceiveComm) {
//	s.callback = callback
//}
//
//func (s *mockSingleSearch) RespondSingleUse(partner single.Contact, payload []byte, _ time.Duration) error {
//	go s.callback(payload, partner)
//	return nil
//}
//
//func (s *mockSingleSearch) StartProcesses() (stoppable.Stoppable, error) {
//	return stoppable.NewSingle(""), nil
//}
