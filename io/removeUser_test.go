package io

import (
	"crypto/rand"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/factID"
	"gitlab.com/elixxir/crypto/hash"
	"gitlab.com/elixxir/primitives/fact"
	"gitlab.com/elixxir/user-discovery-bot/interfaces/params"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"
	"testing"
	"time"
)

// Check that our nil check works
func TestDeleteUser_NilCheck(t *testing.T) {
	_, err := removeUser(nil, nil)
	if err == nil {
		t.Error("removeUser receiving a nil msg didn't error")
	}

	badmsg := pb.FactRemovalRequest{
		UID:         nil,
		RemovalData: nil,
	}
	_, err = removeUser(&badmsg, nil)
	if err == nil {
		t.Error("removeUser receiving a nil msg didn't error")
	}
}

// Checks test for no user having registered the fact
func TestDeleteUser_UsersCheck(t *testing.T) {
	// Make a FactRemovalRequest to put into the Delete function
	badmsg := pb.FactRemovalRequest{
		UID: id.DummyUser.Marshal(),
		RemovalData: &pb.Fact{
			Fact:     "Testing",
			FactType: uint32(fact.Username),
		},
	}

	// Setup a Storage object
	store, _, err := storage.NewStorage(params.Database{})
	if err != nil {
		t.Fatal("connect.NewHost returned an error: ", err)
	}

	_, err = removeUser(&badmsg, store)
	if err == nil {
		t.Error("removeUser receiving a nil msg didn't error")
	}
}

// Test that the function doesn't work when a different user to the Fact owner
// tries to delete the Fact
func TestDeleteUser_WrongOwner(t *testing.T) {
	// Create an input message
	input_msg := pb.FactRemovalRequest{
		UID: []byte{0, 1, 2, 3},
		RemovalData: &pb.Fact{
			Fact:     "Testing",
			FactType: uint32(fact.Username),
		},
	}

	// Setup a Storage object
	store, _, err := storage.NewStorage(params.Database{})
	if err != nil {
		t.Fatal("connect.NewHost returned an error: ", err)
	}

	// Insert a user to assign the Fact to
	suser_factstorage := new([]storage.Fact)
	suser := storage.User{
		Id:        id.NewIdFromUInt(0, id.User, t).Marshal(),
		RsaPub:    "",
		DhPub:     nil,
		Salt:      nil,
		Signature: nil,
		Facts:     *suser_factstorage,
	}
	err = store.InsertUser(&suser)
	if err != nil {
		t.Fatal(err)
	}

	// Create a Fact object to put into our Storage object
	// Generate the hash function and hash the fact
	f, err := fact.NewFact(fact.FactType(input_msg.RemovalData.FactType),
		input_msg.RemovalData.Fact)
	if err != nil {
		t.Fatal(err)
	}
	sfact := storage.Fact{
		Hash:         factID.Fingerprint(f),
		UserId:       id.NewIdFromUInt(0, id.User, t).Marshal(),
		Fact:         "Testing",
		Type:         uint8(fact.Username),
		Signature:    nil,
		Verified:     false,
		Timestamp:    time.Time{},
		Verification: storage.TwilioVerification{},
	}
	err = store.InsertFact(&sfact)
	if err != nil {
		t.Fatal(err)
	}

	// Attempt to delete our User
	_, err = removeUser(&input_msg, store)
	if err == nil {
		t.Error("removeUser did not return an error deleting " +
			"a fact the input user doesn't own")
	}
}

// Test that the function does work given the right inputs and DB entries
func TestDeleteUser_Happy(t *testing.T) {
	clientId, clientKey := initClientFields(t)

	// Create an input message
	input_msg := pb.FactRemovalRequest{
		UID: clientId.Bytes(),
		RemovalData: &pb.Fact{
			Fact:     "Testing",
			FactType: uint32(fact.Username),
		},
	}

	// Setup a Storage object
	store, _, err := storage.NewStorage(params.Database{})
	if err != nil {
		t.Fatal("connect.NewHost returned an error: ", err)
	}

	// Insert a user to assign the Fact to
	suser_factstorage := new([]storage.Fact)
	suser := storage.User{
		Id: clientId.Bytes(),
		RsaPub: string(rsa.CreatePublicKeyPem(
			clientKey.GetPublic())),
		DhPub:     nil,
		Salt:      nil,
		Signature: nil,
		Facts:     *suser_factstorage,
	}
	err = store.InsertUser(&suser)
	if err != nil {
		t.Fatal(err)
	}

	// Create a Fact object to put into our Storage object
	// Generate the hash function and hash the fact
	f, err := fact.NewFact(fact.FactType(input_msg.RemovalData.FactType),
		input_msg.RemovalData.Fact)
	if err != nil {
		t.Fatal(err)
	}
	sfact := storage.Fact{
		Hash:         factID.Fingerprint(f),
		UserId:       clientId.Bytes(),
		Fact:         "Testing",
		Type:         uint8(fact.Username),
		Signature:    nil,
		Verified:     false,
		Timestamp:    time.Time{},
		Verification: storage.TwilioVerification{},
	}
	err = store.InsertFact(&sfact)
	if err != nil {
		t.Fatal(err)
	}

	// Sign the fact
	factSig, _ := rsa.Sign(rand.Reader, clientKey, hash.CMixHash,
		factID.Fingerprint(f), nil)
	input_msg.FactSig = factSig

	// Attempt to delete our Fact
	_, err = removeUser(&input_msg, store)
	if err != nil {
		t.Error("removeUser returned an error: ", err)
	}
}

// Test that the function does work given the right inputs and DB entries
func TestDeleteUser_DeletedUser(t *testing.T) {
	clientId, clientKey := initClientFields(t)

	// Create an input message
	input_msg := pb.FactRemovalRequest{
		UID: clientId.Bytes(),
		RemovalData: &pb.Fact{
			Fact:     "Testing",
			FactType: uint32(fact.Username),
		},
	}

	// Setup a Storage object
	store, _, err := storage.NewStorage(params.Database{})
	if err != nil {
		t.Fatal("connect.NewHost returned an error: ", err)
	}

	// Insert a user to assign the Fact to
	suser_factstorage := new([]storage.Fact)
	suser := storage.User{
		Id: clientId.Bytes(),
		RsaPub: string(rsa.CreatePublicKeyPem(
			clientKey.GetPublic())),
		DhPub:     nil,
		Salt:      nil,
		Signature: nil,
		Facts:     *suser_factstorage,
	}
	err = store.InsertUser(&suser)
	if err != nil {
		t.Fatal(err)
	}

	// Create a Fact object to put into our Storage object
	// Generate the hash function and hash the fact
	f, err := fact.NewFact(fact.FactType(input_msg.RemovalData.FactType),
		input_msg.RemovalData.Fact)
	if err != nil {
		t.Fatal(err)
	}
	sfact := storage.Fact{
		Hash:         factID.Fingerprint(f),
		UserId:       clientId.Bytes(),
		Fact:         "Testing",
		Type:         uint8(fact.Username),
		Signature:    nil,
		Verified:     false,
		Timestamp:    time.Time{},
		Verification: storage.TwilioVerification{},
	}
	err = store.InsertFact(&sfact)
	if err != nil {
		t.Fatal(err)
	}

	// Sign the fact
	factSig, _ := rsa.Sign(rand.Reader, clientKey, hash.CMixHash,
		factID.Fingerprint(f), nil)
	input_msg.FactSig = factSig

	// Attempt to delete our Fact
	_, err = removeUser(&input_msg, store)
	if err != nil {
		t.Error("removeUser returned an error: ", err)
	}

	// Insert a user to assign the Fact to
	suser_factstorage = new([]storage.Fact)
	suser = storage.User{
		Id: clientId.Bytes(),
		RsaPub: string(rsa.CreatePublicKeyPem(
			clientKey.GetPublic())),
		DhPub:     nil,
		Salt:      nil,
		Signature: nil,
		Facts:     *suser_factstorage,
	}
	err = store.InsertUser(&suser)
	if err != nil {
		t.Fatal(err)
	}

	// Create a Fact object to put into our Storage object
	// Generate the hash function and hash the fact
	ret, err := store.Search([][]byte{factID.Fingerprint(f)})
	if err != nil {
		t.Fatal(err)
	}
	if len(ret) == 0 {
		t.Fatal("Expected to find fact already present!")
	}

}
