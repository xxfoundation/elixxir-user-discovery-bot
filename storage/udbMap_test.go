package storage

import (
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

//Hidden function for one-time unit testing database implementation
//func TestDatabaseImpl(t *testing.T) {
//
//	jwalterweatherman.SetLogThreshold(jwalterweatherman.LevelTrace)
//	jwalterweatherman.SetStdoutThreshold(jwalterweatherman.LevelTrace)
//
//	db, _, err := NewDatabase("jonahhusson", "", "cmix_udb", "0.0.0.0", "5432")
//	if err != nil {
//		t.Errorf(err.Error())
//		return
//	}
//
//	uid := []byte("testuid")
//	rsapub := []byte("testrsa")
//
//	err = db.InsertUser(&User{
//		Id:        uid,
//		RsaPub:    rsapub,
//		DhPub:     []byte("testdh"),
//		Salt:      []byte("testsalt"),
//		Signature: []byte("testsig"),
//	})
//	if err != nil {
//		t.Errorf("Failed to insert user: %+v", err)
//	}
//
//	u, err := db.GetUser(uid)
//	if err != nil {
//		t.Errorf("Failed to get user: %+v", err)
//	}
//	if bytes.Compare(u.RsaPub, rsapub) != 0 {
//		t.Error("Retrieved user did not preserve data properly")
//	}
//
//	factid := []byte("factid")
//	err = db.InsertFact(&Fact{
//		ConfirmationId:     factid,
//		UserId:             uid,
//		Fact:               "water is wet",
//		FactType:           0,
//		FactHash:           []byte("facthash"),
//		Signature:          []byte("factsig"),
//		VerificationStatus: 0,
//		Manual:             false,
//		Code:               0,
//	})
//	if err != nil {
//		t.Errorf("Failed to insert fact: %+v", err)
//	}
//
//	_, err = db.GetFact(factid)
//	if err != nil {
//		t.Errorf("Failed to get fact: %+v", err)
//	}
//
//	err = db.ConfirmFact(factid)
//	if err != nil {
//		t.Errorf("Failed to confirm fact: %+v", err)
//	}
//
//	updated, err := db.GetFact(factid)
//	if err != nil {
//		t.Errorf("Failed to get fact after confirm: %+v", err)
//	}
//
//	if updated.VerificationStatus != 1 {
//		t.Error("Fact verification status did not update")
//	}
//
//	err = db.DeleteFact(factid)
//	if err != nil {
//		t.Errorf("Failed to delete fact: %+v", err)
//	}
//
//	err = db.DeleteUser(uid)
//	if err != nil {
//		t.Errorf("Failed to delete user")
//	}
//}

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
		RsaPub:    []byte("testrsa"),
		DhPub:     []byte("testdh"),
		Salt:      []byte("testsalt"),
		Signature: []byte("testsig"),
	}

	err = mapImpl.InsertUser(user)
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}

	factId := []byte("testconfid")
	fact := &Fact{
		ConfirmationId:     factId,
		UserId:             uid.Marshal(),
		Fact:               "water is wet",
		FactType:           0,
		FactHash:           []byte("Definitely a hash"),
		Signature:          []byte("John Hancock"),
		VerificationStatus: 0,
		Manual:             false,
		Code:               0,
	}
	err = mapImpl.InsertFact(fact)
	if err != nil {
		t.Errorf("Failed to insert fact: %+v", err)
	}

	retreived, err := mapImpl.GetFact(factId)
	if err != nil {
		t.Errorf("Failed to get fact after insert: %+v", err)
	}
	if retreived == nil {
		t.Error("Did not retrieve fact properly")
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
		RsaPub:    []byte("testrsa"),
		DhPub:     []byte("testdh"),
		Salt:      []byte("testsalt"),
		Signature: []byte("testsig"),
	}

	err = mapImpl.InsertUser(user)
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}
}

// Unit test for mapimpl get fact
func TestMapImpl_GetFact(t *testing.T) {
	mapImpl, _, err := NewDatabase("", "", "", "", "")
	if err != nil {
		t.Errorf("Failed to create map impl")
		t.FailNow()
	}
	uid := id.NewIdFromString("testuserid", id.User, t)
	user := &User{
		Id:        uid.Marshal(),
		RsaPub:    []byte("testrsa"),
		DhPub:     []byte("testdh"),
		Salt:      []byte("testsalt"),
		Signature: []byte("testsig"),
	}

	err = mapImpl.InsertUser(user)
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}

	factId := []byte("testconfid")
	fact := &Fact{
		ConfirmationId:     factId,
		UserId:             uid.Marshal(),
		Fact:               "water is wet",
		FactType:           0,
		FactHash:           []byte("Definitely a hash"),
		Signature:          []byte("John Hancock"),
		VerificationStatus: 0,
		Manual:             false,
		Code:               0,
	}
	err = mapImpl.InsertFact(fact)
	if err != nil {
		t.Errorf("Failed to insert fact: %+v", err)
	}

	_, err = mapImpl.GetFact(factId)
	if err != nil {
		t.Errorf("Failed to get fact after insert: %+v", err)
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
		RsaPub:    []byte("testrsa"),
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
		RsaPub:    []byte("testrsa"),
		DhPub:     []byte("testdh"),
		Salt:      []byte("testsalt"),
		Signature: []byte("testsig"),
	}

	err = mapImpl.InsertUser(user)
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}

	factId := []byte("testconfid")
	fact := &Fact{
		ConfirmationId:     factId,
		UserId:             uid.Marshal(),
		Fact:               "water is wet",
		FactType:           0,
		FactHash:           []byte("Definitely a hash"),
		Signature:          []byte("John Hancock"),
		VerificationStatus: 0,
		Manual:             false,
		Code:               0,
	}
	err = mapImpl.InsertFact(fact)
	if err != nil {
		t.Errorf("Failed to insert fact: %+v", err)
	}

	err = mapImpl.ConfirmFact(factId)
	if err != nil {
		t.Errorf("Failed to confirm fact: %+v", err)
	}

	newfact, err := mapImpl.GetFact(factId)
	if err != nil {
		t.Errorf("Failed to retreive fact after confirming: %+v", err)
	}

	if newfact.VerificationStatus != 1 {
		t.Error("Failed to verify fact")
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
		RsaPub:    []byte("testrsa"),
		DhPub:     []byte("testdh"),
		Salt:      []byte("testsalt"),
		Signature: []byte("testsig"),
	}

	err = mapImpl.InsertUser(user)
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}

	factId := []byte("testconfid")
	fact := &Fact{
		ConfirmationId:     factId,
		UserId:             uid.Marshal(),
		Fact:               "water is wet",
		FactType:           0,
		FactHash:           []byte("Definitely a hash"),
		Signature:          []byte("John Hancock"),
		VerificationStatus: 0,
		Manual:             false,
		Code:               0,
	}
	err = mapImpl.InsertFact(fact)
	if err != nil {
		t.Errorf("Failed to insert fact: %+v", err)
	}

	_, err = mapImpl.GetFact(factId)
	if err != nil {
		t.Errorf("Failed to get fact after insert: %+v", err)
	}

	err = mapImpl.DeleteFact(factId)
	if err != nil {
		t.Errorf("Failed to delete fact: %+v", err)
	}

	deleted, err := mapImpl.GetFact(factId)
	if err != nil {
		t.Errorf("GetFact should not error if fact is deleted: %+v", err)
	}
	if deleted != nil {
		t.Errorf("Fact %+v was not deleted", deleted)
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
		RsaPub:    []byte("testrsa"),
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
