package cmix

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"gitlab.com/elixxir/client/single"
	"gitlab.com/elixxir/client/stoppable"
	"gitlab.com/elixxir/client/ud"
	"gitlab.com/elixxir/crypto/cyclic"
	"gitlab.com/elixxir/crypto/e2e/singleUse"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/crypto/large"
	"gitlab.com/xx_network/primitives/id"
	"strings"
	"testing"
	"time"
)

func TestManager_lookupCallback(t *testing.T) {
	m := &Manager{storage.NewTestDB(t), &mockSingleLookup{}}
	uid := id.NewIdFromString("zezima", id.User, t)
	grp := cyclic.NewGroup(large.NewInt(107), large.NewInt(2))
	ct := single.NewContact(uid, grp.NewInt(42), grp.NewInt(43), singleUse.TagFP{}, 8)
	err := m.db.InsertUser(&storage.User{Id: uid.Marshal(), DhPub: []byte("DhPub")})
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}

	lookupMsg := &ud.LookupSend{UserID: uid.Marshal()}
	payload, err := proto.Marshal(lookupMsg)
	if err != nil {
		t.Errorf("Failed to marshal payload: %+v", err)
	}
	expectedPayload, err := proto.Marshal(m.handleLookup(lookupMsg, ct))
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

	m.lookupCallback(payload, ct)

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

// Happy path.
func TestManager_handleLookup(t *testing.T) {
	m := &Manager{db: storage.NewTestDB(t)}
	uid := id.NewIdFromString("zezima", id.User, t)
	c := single.NewContact(uid, &cyclic.Int{}, &cyclic.Int{}, singleUse.TagFP{}, 8)

	expectedDhPub := "DhPub"
	err := m.db.InsertUser(&storage.User{Id: uid.Marshal(), DhPub: []byte(expectedDhPub)})
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}

	resp := m.handleLookup(&ud.LookupSend{UserID: uid.Marshal()}, c)
	if resp.Error != "" {
		t.Errorf("handleLookup() returned a response with an error: %s", resp.Error)
	}

	if string(resp.PubKey) != expectedDhPub {
		t.Errorf("handleLookup() returned a response with inccorect PubKey."+
			"\nexpected: %s\nreceived: %s", expectedDhPub, resp.Error)
	}
}

// Error path: Id is malformed and fails to unmarshal.
func TestManager_handleLookup_IdUnmarshalError(t *testing.T) {
	m := &Manager{db: storage.NewTestDB(t)}
	uid := id.NewIdFromString("zezima", id.User, t)
	c := single.NewContact(uid, &cyclic.Int{}, &cyclic.Int{}, singleUse.TagFP{}, 8)

	err := m.db.InsertUser(&storage.User{Id: uid.Marshal(), DhPub: []byte("DhPub")})
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}

	resp := m.handleLookup(&ud.LookupSend{UserID: []byte{1}}, c)
	if !strings.Contains(resp.Error, "failed to unmarshal lookup ID in request") {
		t.Errorf("handleLookup() returned a response with an error: %s", resp.Error)
	}
}

// mockSingleLookup is used to test the lookup function, which uses the single-
// use manager. It adheres to the SingleInterface interface.
type mockSingleLookup struct {
	callback func(payload []byte, c single.Contact)
}

func (s *mockSingleLookup) RegisterCallback(_ string, callback single.ReceiveComm) {
	s.callback = callback
}

func (s *mockSingleLookup) RespondSingleUse(partner single.Contact, payload []byte, _ time.Duration) error {
	go s.callback(payload, partner)
	return nil
}

func (s *mockSingleLookup) StartProcesses() stoppable.Stoppable {
	return stoppable.NewSingle("")
}
