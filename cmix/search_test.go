package cmix

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"git.xx.network/elixxir/client/single"
	"git.xx.network/elixxir/client/stoppable"
	"git.xx.network/elixxir/client/ud"
	"git.xx.network/elixxir/crypto/cyclic"
	"git.xx.network/elixxir/crypto/e2e/singleUse"
	"git.xx.network/elixxir/primitives/fact"
	"git.xx.network/elixxir/user-discovery-bot/storage"
	"git.xx.network/xx_network/crypto/large"
	"git.xx.network/xx_network/primitives/id"
	"reflect"
	"testing"
	"time"
)

func TestManager_SearchProcess(t *testing.T) {
	m := &Manager{storage.NewTestDB(t), &mockSingleSearch{}}
	uid := id.NewIdFromString("zezima", id.User, t)
	grp := cyclic.NewGroup(large.NewInt(107), large.NewInt(2))
	ct := single.NewContact(uid, grp.NewInt(42), grp.NewInt(43), singleUse.TagFP{}, 8)
	err := m.db.InsertUser(&storage.User{Id: uid.Marshal(), DhPub: []byte("DhPub")})
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}

	fid := []byte("factHash")
	searchMsg := &ud.SearchSend{Fact: []*ud.HashFact{{Hash: fid, Type: 0}}}
	payload, err := proto.Marshal(searchMsg)
	if err != nil {
		t.Errorf("Failed to marshal payload: %+v", err)
	}

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
	expectedPayload, err := proto.Marshal(m.handleSearch(searchMsg, ct))
	if err != nil {
		t.Errorf("Failed to marshal message: %+v", err)
	}

	callbackChan := make(chan struct {
		payload []byte
		c       single.Contact
	})
	callback := func(payload []byte, c single.Contact) {
		callbackChan <- struct {
			payload []byte
			c       single.Contact
		}{payload: payload, c: c}
	}
	m.singleUse.RegisterCallback("", callback)

	m.searchCallback(payload, ct)
	results := <-callbackChan
	if !results.c.Equal(ct) {
		t.Errorf("Callback did not return the expected contact."+
			"\nexpected: %s\nreceived: %s", ct, results.c)
	}
	if !bytes.Equal(expectedPayload, results.payload) {
		t.Errorf("Callback did not return the expected payload."+
			"\nexpected: %+v\nreceived: %+v", expectedPayload, results.payload)
	}
}

func TestManager_handleSearch(t *testing.T) {
	m := &Manager{db: storage.NewTestDB(t)}
	uid := id.NewIdFromString("zezima", id.User, t)
	c := single.NewContact(uid, &cyclic.Int{}, &cyclic.Int{}, singleUse.TagFP{}, 8)

	expectedDhPub := []byte("DhPub")
	fid := []byte("factHash")
	f1 := &storage.Fact{
		Hash:     fid,
		UserId:   uid.Marshal(),
		Fact:     "water is wet",
		Type:     1,
		Verified: true,
	}
	expectedContact := &ud.Contact{
		UserID:    uid.Marshal(),
		PubKey:    expectedDhPub,
		TrigFacts: []*ud.HashFact{{Hash: f1.Hash, Type: int32(f1.Type)}},
	}
	err := m.db.InsertUser(&storage.User{Id: uid.Marshal(), DhPub: expectedDhPub})
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}

	err = m.db.InsertFact(f1)
	if err != nil {
		t.Errorf("Failed to insert dummy fact: %+v", err)
	}

	resp := m.handleSearch(&ud.SearchSend{Fact: []*ud.HashFact{{Hash: fid, Type: 0}}}, c)
	if resp.Error != "" {
		t.Errorf("Failed to handle search: %+v", resp.Error)
	}
	if len(resp.Contacts) != 1 {
		t.Errorf("Did not receive expected number of contacts")
	}
	if !reflect.DeepEqual(expectedContact, resp.Contacts[0]) {
		t.Errorf("Did not received expected contact."+
			"\nexpected: %+v\nreceived: %+v", expectedContact, resp.Contacts[0])
	}

	resp = m.handleSearch(&ud.SearchSend{Fact: []*ud.HashFact{{Hash: fid, Type: int32(fact.Nickname)}}}, c)
	if resp.Error == "" {
		t.Errorf("Search should have returned error")
	}
	if len(resp.Contacts) != 0 {
		t.Errorf("Should not be able to search with nickname")
	}
}

// mockSingleSearch is used to test the search function, which uses the single-
// use manager. It adheres to the SingleInterface interface.
type mockSingleSearch struct {
	callback func(payload []byte, c single.Contact)
}

func (s *mockSingleSearch) RegisterCallback(_ string, callback single.ReceiveComm) {
	s.callback = callback
}

func (s *mockSingleSearch) RespondSingleUse(partner single.Contact, payload []byte, _ time.Duration) error {
	go s.callback(payload, partner)
	return nil
}

func (s *mockSingleSearch) StartProcesses() stoppable.Stoppable {
	return stoppable.NewSingle("")
}
