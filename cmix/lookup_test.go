package cmix

//func TestManager_lookupCallback(t *testing.T) {
//	m := &Manager{db: storage.NewTestDB(t),
//		lookupListener: &mockSingleLookup{}}
//	uid := id.NewIdFromString("zezima", id.User, t)
//	grp := cyclic.NewGroup(large.NewInt(107), large.NewInt(2))
//
//	dhPub := grp.NewInt(25)
//
//	c := contact.Contact{
//		ID:             uid,
//		DhPubKey:       dhPub,
//	}
//
//	lm := &lookupManager{m: m}
//
//	err := m.db.InsertUser(&storage.User{Id: uid.Marshal(), DhPub: dhPub.Bytes())}
//	if err != nil {
//		t.Errorf("Failed to insert dummy user: %+v", err)
//	}
//
//
//	lookupMsg := &ud.LookupSend{UserID: uid.Marshal()}
//	payload, err := proto.Marshal(lookupMsg)
//	if err != nil {
//		t.Errorf("Failed to marshal payload: %+v", err)
//	}
//
//	single.Listen(ud.LookupTag, uid, payload, lm, single.GetDefaultRequestParams(), mockCMix, rng, grp)
//
//	expectedPayload, err := proto.Marshal(lookupMsg)
//	if err != nil {
//		t.Errorf("Failed to marshal message: %+v", err)
//	}
//
//	callbackChan := make(chan struct {
//		payload []byte
//		c       single.Contact
//	})
//	callback := func(payload []byte, c single.Contact) {
//		callbackChan <- struct {
//			payload []byte
//			c       single.Contact
//		}{payload: payload, c: c}
//	}
//	m.singleUse.RegisterCallback("", callback)
//
//	m.lookupCallback(payload, ct)
//
//	results := <-callbackChan
//	if !results.c.Equal(ct) {
//		t.Errorf("Callback did not return the expected contact."+
//			"\nexpected: %s\nreceived: %s", ct, results.c)
//	}
//	if !bytes.Equal(expectedPayload, results.payload) {
//		t.Errorf("Callback did not return the expected payload."+
//			"\nexpected: %+v\nreceived: %+v", expectedPayload, results.payload)
//	}
//}

//
//// Happy path.
//func TestManager_handleLookup(t *testing.T) {
//	m := &Manager{db: storage.NewTestDB(t)}
//	uid := id.NewIdFromString("zezima", id.User, t)
//
//	username := "ZeZima"
//	expectedDhPub := "DhPub"
//	err := m.db.InsertUser(
//		&storage.User{Id: uid.Marshal(),
//			DhPub:    []byte(expectedDhPub),
//			Username: username,
//			Facts: []storage.Fact{
//				{
//					Hash:      []byte("hash"),
//					UserId:    uid.Marshal(),
//					Fact:      strings.ToLower(username),
//					Type:      0,
//					Signature: []byte("Signature"),
//					Verified:  true,
//				},
//			}})
//	if err != nil {
//		t.Errorf("Failed to insert dummy user: %+v", err)
//	}
//
//	lookupManager := lookupManager{m: m}
//
//	single.TransmitRequest()
//
//	resp := lookupManager.handleLo  okup(&ud.LookupSend{UserID: uid.Marshal()}, c)
//	if resp.Error != "" {
//		t.Errorf("handleLookup() returned a response with an error: %s", resp.Error)
//	}
//
//	if string(resp.PubKey) != expectedDhPub {
//		t.Errorf("handleLookup() returned a response with inccorect PubKey."+
//			"\nexpected: %s\nreceived: %s", expectedDhPub, resp.Error)
//	}
//
//	if resp.Username != username {
//		t.Errorf("Should have gotten username %s, instead got: %s", username, resp.Username)
//	}
//}
//
//// Error path: Id is malformed and fails to unmarshal.
//func TestManager_handleLookup_IdUnmarshalError(t *testing.T) {
//	m := &Manager{db: storage.NewTestDB(t)}
//	uid := id.NewIdFromString("zezima", id.User, t)
//	c := single.NewContact(uid, &cyclic.Int{}, &cyclic.Int{}, singleUse.TagFP{}, 8)
//
//	err := m.db.InsertUser(&storage.User{Id: uid.Marshal(), DhPub: []byte("DhPub")})
//	if err != nil {
//		t.Errorf("Failed to insert dummy user: %+v", err)
//	}
//
//	resp := m.handleLookup(&ud.LookupSend{UserID: []byte{1}}, c)
//	if !strings.Contains(resp.Error, "failed to unmarshal lookup ID in request") {
//		t.Errorf("handleLookup() returned a response with an error: %s", resp.Error)
//	}
//}

// mockSingleLookup is used to test the lookup function, which uses the single-
// use manager. It adheres to the SingleInterface interface.
type mockSingleLookup struct {
}

func (s *mockSingleLookup) Stop() {

}
