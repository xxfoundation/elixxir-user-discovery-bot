package storage

import (
	"github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/primitives/id"
	"strings"
	"testing"
)

// Hidden function for one-time unit testing database implementation
func TestDatabaseImpl(t *testing.T) {

	jwalterweatherman.SetLogThreshold(jwalterweatherman.LevelTrace)
	jwalterweatherman.SetStdoutThreshold(jwalterweatherman.LevelTrace)

	db, _, err := NewDatabase("jonahhusson", "", "cmix_udb", "0.0.0.0", "5432")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	uid := id.NewIdFromString("testid", id.User, t)
	rsapub := "testrsa"

	err = db.InsertUser(&User{
		Id:        uid.Marshal(),
		RsaPub:    rsapub,
		DhPub:     []byte("testdh"),
		Salt:      []byte("testsalt"),
		Signature: []byte("testsig"),
	})
	if err != nil {
		t.Errorf("Failed to insert user: %+v", err)
	}

	u, err := db.GetUser(uid.Marshal())
	if err != nil {
		t.Errorf("Failed to get user: %+v", err)
	}
	if strings.Compare(u.RsaPub, rsapub) != 0 {
		t.Error("Retrieved user did not preserve data properly")
	}

	factid := []byte("facthash")
	err = db.InsertFact(&Fact{
		Hash:      factid,
		UserId:    uid.Marshal(),
		Fact:      "zezima",
		Type:      uint8(Username),
		Signature: []byte("factsig"),
		Verified:  false,
	})
	if err != nil {
		t.Errorf("Failed to insert fact: %+v", err)
	}

	err = db.CheckUser("zezima", uid, rsapub)
	if err == nil {
		t.Error("Should have returned error")
	}

	err = db.CheckUser("tim", id.NewIdFromString("tim", id.Node, t), "tim")
	if err != nil {
		t.Errorf("Failed to check user: %+v", err)
	}

	err = db.VerifyFact(factid)
	if err != nil {
		t.Errorf("Failed to verify fact: %+v", err)
	}

	factid2 := []byte("facthashtwo")
	err = db.InsertFactTwilio(uid.Marshal(), factid2, []byte("factsig2"), "twilio", 1, "conf")
	if err != nil {
		t.Errorf("Failed to insert twilio-verified fact: %+v", err)
	}

	err = db.VerifyFactTwilio("conf")
	if err != nil {
		t.Errorf("Failed to verify twilio fact: %+v", err)
	}

	err = db.DeleteFact(factid2)
	if err != nil {
		t.Errorf("Failed to delete fact2: %+v", err)
	}

	err = db.DeleteFact(factid)
	if err != nil {
		t.Errorf("Failed to delete fact: %+v", err)
	}

	err = db.DeleteUser(uid.Marshal())
	if err != nil {
		t.Errorf("Failed to delete user: %+v", err)
	}
}

// Unit test for mapimpl insert fact
func TestMapImpl_InsertFact(t *testing.T) {
	mapImpl, _, err := NewDatabase("", "", "", "", "")
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

}

// Unit test for mapimpl insert user
func TestMapImpl_InsertUser(t *testing.T) {
	mapImpl, _, err := NewDatabase("", "", "", "", "")
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
}

// Unit test for mapimpl get user
func TestMapImpl_GetUser(t *testing.T) {
	mapImpl, _, err := NewDatabase("", "", "", "", "")
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

// Unit test for mapimpl confirm fact
func TestMapImpl_ConfirmFact(t *testing.T) {
	mapImpl, _, err := NewDatabase("", "", "", "", "")
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
		Hash:      factHash,
		UserId:    uid.Marshal(),
		Fact:      "water is wet",
		Type:      0,
		Signature: []byte("John Hancock"),
		Verified:  true,
	}
	err = mapImpl.InsertFact(fact)
	if err != nil {
		t.Errorf("Failed to insert fact: %+v", err)
	}
}

// Unit test for mapimpl delete fact
func TestMapImpl_DeleteFact(t *testing.T) {
	mapImpl, _, err := NewDatabase("", "", "", "", "")
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
	mapImpl, _, err := NewDatabase("", "", "", "", "")
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
