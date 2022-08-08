package storage

import (
	"bytes"
	"gitlab.com/xx_network/primitives/id"
	"testing"
	"time"
)

// Hidden function for one-time unit testing database implementation
//func TestDatabaseImpl(t *testing.T) {
//
//	jwalterweatherman.SetLogThreshold(jwalterweatherman.LevelTrace)
//	jwalterweatherman.SetStdoutThreshold(jwalterweatherman.LevelTrace)
//
//	db, err := newDatabase("jonahhusson", "", "cmix_server", "0.0.0.0", "5432")
//	if err != nil {
//		t.Errorf(err.Error())
//		return
//	}
//
//	uid := id.NewIdFromString("testid", id.User, t)
//	rsapub := "testrsa"
//
//	err = db.InsertUser(&User{
//		Id:        uid.Marshal(),
//		RsaPub:    rsapub,
//		DhPub:     []byte("testdh"),
//		Salt:      []byte("testsalt"),
//		Signature: []byte("testsig"),
//	})
//	if err != nil {
//		t.Errorf("Failed to insert user: %+v", err)
//	}
//
//	u, err := db.GetUser(uid.Marshal())
//	if err != nil {
//		t.Errorf("Failed to get user: %+v", err)
//	}
//	if strings.Compare(u.RsaPub, rsapub) != 0 {
//		t.Error("Retrieved user did not preserve data properly")
//	}
//
//	factid := []byte("facthash")
//	err = db.InsertFact(&Fact{
//		Hash:      factid,
//		UserId:    uid.Marshal(),
//		Fact:      "zezima",
//		Type:      uint8(Username),
//		Signature: []byte("factsig"),
//		Verified:  false,
//	})
//	if err != nil {
//		t.Errorf("Failed to insert fact: %+v", err)
//	}
//
//	err = db.CheckUser("ZeZiMa", uid)
//	if err == nil {
//		t.Error("Should have returned error")
//	}
//
//	err = db.CheckUser("tim", id.NewIdFromString("tim", id.Node, t))
//	if err != nil {
//		t.Errorf("Failed to check user: %+v", err)
//	}
//
//	err = db.MarkFactVerified(factid)
//	if err != nil {
//		t.Errorf("Failed to verify fact: %+v", err)
//	}
//
//	factid2 := []byte("facthashtwo")
//	err = db.InsertFactTwilio(uid.Marshal(), factid2, []byte("factsig2"), 1, "conf")
//	if err != nil {
//		t.Errorf("Failed to insert twilio-verified fact: %+v", err)
//	}
//
//	err = db.MarkTwilioFactVerified("conf")
//	if err != nil {
//		t.Errorf("Failed to verify twilio fact: %+v", err)
//	}
//
//	users, err := db.Search([][]byte{
//		factid, factid2,
//	})
//	if err != nil {
//		t.Errorf("Failed to search for users: %+v", err)
//	}
//	if len(users) != 1 || len(users[0].Facts) != 2 {
//		t.Errorf("Search did not return expected results: %+v", users[0].Facts)
//	}
//
//	err = db.DeleteFact(factid2)
//	if err != nil {
//		t.Errorf("Failed to delete fact2: %+v", err)
//	}
//
//	err = db.DeleteFact(factid)
//	if err != nil {
//		t.Errorf("Failed to delete fact: %+v", err)
//	}
//
//	err = db.DeleteUser(uid.Marshal())
//	if err != nil {
//		t.Errorf("Failed to delete user: %+v", err)
//	}
//}
//
//// Unit test for mapimpl insert fact
//func TestMapImpl_InsertFact(t *testing.T) {
//	mapImpl, err := newDatabase("", "", "", "", "")
//	if err != nil {
//		t.Errorf("Failed to create map impl")
//		t.FailNow()
//	}
//	uid := id.NewIdFromString("testuserid", id.User, t)
//	user := &User{
//		Id:        uid.Marshal(),
//		RsaPub:    "testrsa",
//		DhPub:     []byte("testdh"),
//		Salt:      []byte("testsalt"),
//		Signature: []byte("testsig"),
//	}
//
//	err = mapImpl.InsertUser(user)
//	if err != nil {
//		t.Errorf("Failed to insert dummy user: %+v", err)
//	}
//
//	factHash := []byte("testconfid")
//	fact := &Fact{
//		UserId:    uid.Marshal(),
//		Fact:      "water is wet",
//		Type:      0,
//		Hash:      factHash,
//		Signature: []byte("John Hancock"),
//		Verified:  true,
//	}
//	err = mapImpl.InsertFact(fact)
//	if err != nil {
//		t.Errorf("Failed to insert fact: %+v", err)
//	}
//
//}
//
//// Unit test for mapimpl insert user
//func TestMapImpl_InsertUser(t *testing.T) {
//	mapImpl, err := newDatabase("", "", "", "", "")
//	if err != nil {
//		t.Errorf("Failed to create map impl")
//		t.FailNow()
//	}
//	uid := id.NewIdFromString("testuserid", id.User, t)
//	user := &User{
//		Id:        uid.Marshal(),
//		RsaPub:    "testrsa",
//		DhPub:     []byte("testdh"),
//		Salt:      []byte("testsalt"),
//		Signature: []byte("testsig"),
//	}
//
//	err = mapImpl.InsertUser(user)
//	if err != nil {
//		t.Errorf("Failed to insert dummy user: %+v", err)
//	}
//}

// Unit test for mapimpl get user
func TestMapImpl_GetUser(t *testing.T) {
	mapImpl, err := newDatabase("", "", "", "", "")
	if err != nil {
		t.Errorf("Failed to create map impl")
		t.FailNow()
	}
	uid := id.NewIdFromString("testuserid", id.User, t)
	user := &User{
		Id:        uid.Marshal(),
		RsaPub:    "testrsa",
		DhPub:     []byte("testdh"),
		Salt:      []byte("testsalt"),
		Signature: []byte("testsig"),
	}

	err = mapImpl.InsertUser(user)
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}

	u, err := mapImpl.GetUser(uid.Marshal())
	if err != nil {
		t.Errorf("Failed to retrieve user: %+v", err)
	}
	if u == nil {
		t.Errorf("User was not properly inserted")
	}
}

// Unit test for mapimpl delete fact
func TestMapImpl_DeleteFact(t *testing.T) {
	mapImpl, err := newDatabase("", "", "", "", "")
	if err != nil {
		t.Errorf("Failed to create map impl")
		t.FailNow()
	}
	uid := id.NewIdFromString("testuserid", id.User, t)
	user := &User{
		Id:        uid.Marshal(),
		RsaPub:    "testrsa",
		DhPub:     []byte("testdh"),
		Salt:      []byte("testsalt"),
		Signature: []byte("testsig"),
	}

	err = mapImpl.InsertUser(user)
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}

	factHash := []byte("testconfid")
	fact := &Fact{
		UserId:    uid.Marshal(),
		Fact:      "water is wet",
		Type:      0,
		Hash:      factHash,
		Signature: []byte("John Hancock"),
		Verified:  true,
	}
	err = mapImpl.InsertFact(fact)
	if err != nil {
		t.Errorf("Failed to insert fact: %+v", err)
	}

	err = mapImpl.DeleteFact(factHash)
	if err != nil {
		t.Errorf("Failed to delete fact: %+v", err)
	}

}

// Unit test for mapimpl delete user
func TestMapImpl_DeleteUser(t *testing.T) {
	mapImpl, err := newDatabase("", "", "", "", "")
	if err != nil {
		t.Errorf("Failed to create map impl")
		t.FailNow()
	}
	uid := id.NewIdFromString("testuserid", id.User, t)
	user := &User{
		Id:        uid.Marshal(),
		RsaPub:    "testrsa",
		DhPub:     []byte("testdh"),
		Salt:      []byte("testsalt"),
		Signature: []byte("testsig"),
	}

	err = mapImpl.InsertUser(user)
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}

	u, err := mapImpl.GetUser(uid.Marshal())
	if err != nil {
		t.Errorf("Failed to retrieve user: %+v", err)
	}
	if u == nil {
		t.Errorf("User was not properly inserted")
	}

	err = mapImpl.DeleteUser(uid.Marshal())
	if err != nil {
		t.Errorf("Failed to delete user: %+v", err)
	}

	u, err = mapImpl.GetUser(uid.Marshal())
	if err != nil {
		t.Errorf("Failed to retrieve user after delete: %+v", err)
	}
	if u != nil {
		t.Errorf("User was not properly deleted")
	}
}

// Unit test for mapimpl verify fact
func TestMapImpl_VerifyFact(t *testing.T) {
	mapImpl, err := newDatabase("", "", "", "", "")
	if err != nil {
		t.Errorf("Failed to create map impl")
		t.FailNow()
	}
	uid := id.NewIdFromString("testuserid", id.User, t)
	user := &User{
		Id:        uid.Marshal(),
		RsaPub:    "testrsa",
		DhPub:     []byte("testdh"),
		Salt:      []byte("testsalt"),
		Signature: []byte("testsig"),
	}

	err = mapImpl.InsertUser(user)
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}

	factHash := []byte("testconfid")
	fact := &Fact{
		UserId:    uid.Marshal(),
		Fact:      "water is wet",
		Type:      0,
		Hash:      factHash,
		Signature: []byte("John Hancock"),
		Verified:  false,
	}
	err = mapImpl.InsertFact(fact)
	if err != nil {
		t.Errorf("Failed to insert fact: %+v", err)
	}

	err = mapImpl.MarkFactVerified(factHash)
	if err != nil {
		t.Errorf("Failed to verify fact: %+v", err)
	}
}

// unit test for mapimpl check user
func TestMapImpl_CheckUser(t *testing.T) {
	mapImpl, err := newDatabase("", "", "", "", "")
	if err != nil {
		t.Errorf("Failed to create map impl")
		t.FailNow()
	}
	err = mapImpl.CheckUser("", id.NewIdFromString("test", id.User, t))
	if err != nil {
		t.Errorf("This should always return nil from map impl: %+v", err)
	}
}

// unit test for insert twilio fact on map backend
func TestMapImpl_InsertFactTwilio(t *testing.T) {
	mapImpl, err := newDatabase("", "", "", "", "")
	if err != nil {
		t.Errorf("Failed to create map impl")
		t.FailNow()
	}
	uid := id.NewIdFromString("testuserid", id.User, t)
	user := &User{
		Id:        uid.Marshal(),
		RsaPub:    "testrsa",
		DhPub:     []byte("testdh"),
		Salt:      []byte("testsalt"),
		Signature: []byte("testsig"),
	}

	err = mapImpl.InsertUser(user)
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}

	factHash := []byte("testconfid")
	conf := "twilio"
	err = mapImpl.InsertFactTwilio(uid.Marshal(), factHash, []byte("John Hancock"), 0, conf)
	if err != nil {
		t.Errorf("Failed to insert fact: %+v", err)
	}
}

// unit test for verifying a twilio fact in the map backend
func TestMapImpl_VerifyFactTwilio(t *testing.T) {
	mapImpl, err := newDatabase("", "", "", "", "")
	if err != nil {
		t.Errorf("Failed to create map impl")
		t.FailNow()
	}
	uid := id.NewIdFromString("testuserid", id.User, t)
	user := &User{
		Id:        uid.Marshal(),
		RsaPub:    "testrsa",
		DhPub:     []byte("testdh"),
		Salt:      []byte("testsalt"),
		Signature: []byte("testsig"),
	}

	err = mapImpl.InsertUser(user)
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}

	factHash := []byte("testconfid")
	conf := "twilio"
	err = mapImpl.InsertFactTwilio(uid.Marshal(), factHash, []byte("John Hancock"), 0, conf)
	if err != nil {
		t.Errorf("Failed to insert fact: %+v", err)
	}

	err = mapImpl.MarkTwilioFactVerified(conf)
	if err != nil {
		t.Errorf("Failed to verify twilio fact: %+v", err)
	}
}

// Search unit test for map backend
func TestMapImpl_Search(t *testing.T) {
	mapImpl, err := newDatabase("", "", "", "", "")
	if err != nil {
		t.Errorf("Failed to create map impl")
		t.FailNow()
	}
	uid := id.NewIdFromString("testuserid", id.User, t)
	user := &User{
		Id:        uid.Marshal(),
		RsaPub:    "testrsa",
		DhPub:     []byte("testdh"),
		Salt:      []byte("testsalt"),
		Signature: []byte("testsig"),
	}

	err = mapImpl.InsertUser(user)
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}

	factHash := []byte("testconfid")
	conf := "twilio"
	err = mapImpl.InsertFactTwilio(uid.Marshal(), factHash, []byte("John Hancock"), 0, conf)
	if err != nil {
		t.Errorf("Failed to insert fact: %+v", err)
	}

	err = mapImpl.MarkTwilioFactVerified(conf)
	if err != nil {
		t.Errorf("Failed to verify twilio fact: %+v", err)
	}

	ulist, err := mapImpl.Search([][]byte{factHash})
	if err != nil {
		t.Errorf("Failed to search: %+v", err)
	}
	if len(ulist) != 1 {
		t.Errorf("Did not receive expected num users.  Received: %d, expected: %d", len(ulist), 1)
	}
}

// Unit test for inserting channel identity
func TestMapImpl_InsertChannelIdentity(t *testing.T) {
	mapImpl, err := newDatabase("", "", "", "", "")
	if err != nil {
		t.Errorf("Failed to create map impl")
		t.FailNow()
	}
	uid := id.NewIdFromString("testuserid", id.User, t)
	ci := &ChannelIdentity{
		UserId:     uid.Marshal(),
		Ed25519Pub: []byte("eidpub"),
		Lease:      time.Now().UnixNano(),
		Banned:     false,
	}
	err = mapImpl.InsertChannelIdentity(ci)
	if err != nil {
		t.Fatalf("Failed to insert channel Identity: %+v", err)
	}
}

// Unit test for getting channel identity
func TestMapImpl_GetChannelIdentity(t *testing.T) {
	mapImpl, err := newDatabase("", "", "", "", "")
	if err != nil {
		t.Errorf("Failed to create map impl")
		t.FailNow()
	}
	uid := id.NewIdFromString("testuserid", id.User, t)
	ci := &ChannelIdentity{
		UserId:     uid.Marshal(),
		Ed25519Pub: []byte("eidpub"),
		Lease:      time.Now().UnixNano(),
		Banned:     false,
	}
	err = mapImpl.InsertChannelIdentity(ci)
	if err != nil {
		t.Fatalf("Failed to insert channel Identity: %+v", err)
	}
	ciLoaded, err := mapImpl.GetChannelIdentity(uid.Marshal())
	if err != nil {
		t.Fatalf("Failed to load channel identity: %+v", err)
	}
	if !bytes.Equal(ci.Ed25519Pub, ciLoaded.Ed25519Pub) {
		t.Errorf("Did not receive expected data from map impl\n\tExpected: %+v\n\tReceived: %+v\n", ci, ciLoaded)
	}
}
