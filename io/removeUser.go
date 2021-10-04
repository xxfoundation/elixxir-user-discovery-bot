///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package io

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/factID"
	"gitlab.com/elixxir/crypto/hash"
	"gitlab.com/elixxir/primitives/fact"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"
	"time"
)

// Takes in a FactRemovalRequest from a client and deletes the Fact if the client owns it
func removeUser(msg *pb.FactRemovalRequest, store *storage.Storage) (
	*messages.Ack, error) {
	// Generic copy of the internal error message
	e := errors.New("Removal could not be " +
		"completed do to internal error, please try again later")

	if msg == nil || msg.RemovalData == nil || msg.UID == nil {
		return &messages.Ack{}, errors.New("Unable to parse required " +
			"fields in registration message")
	}

	if fact.FactType(msg.RemovalData.FactType) != fact.Username {
		return &messages.Ack{}, errors.New(
			"RemoveUser requires a username")
	}

	// Get the user ID from the username

	// Generate the hash function and hash the fact
	// NOTE: We don't call NewFact here because there's a trap door that
	// causes it to return an empty fact for usernames.
	f := fact.Fact{
		T:    fact.FactType(msg.RemovalData.FactType),
		Fact: msg.RemovalData.Fact,
	}
	hashFact := factID.Fingerprint(f)

	// Get the user who owns the fact
	users, err := store.Search([][]byte{hashFact})
	if err != nil {
		return &messages.Ack{}, err
	}
	if len(users) != 1 {
		jww.ERROR.Print("removeUser internal error users != 1")
		return &messages.Ack{}, e
	}
	// Unmarshal the owner ID
	uid, err := id.Unmarshal(users[0].Id)
	if err != nil {
		jww.ERROR.Print("removeUser internal error unmarshalling "+
			"found user id", err)
		return &messages.Ack{}, e
	}

	// Parse the client's public key
	clientPubKey, err := rsa.LoadPublicKeyFromPem([]byte(users[0].RsaPub))
	if err != nil {
		return &messages.Ack{}, errors.WithMessage(err,
			invalidUserKeyError)
	}

	// Return an error if the fact signature cannot be verified
	err = rsa.Verify(clientPubKey, hash.CMixHash, hashFact,
		msg.FactSig, nil)
	if err != nil {
		return &messages.Ack{}, errors.WithMessage(err,
			invalidFactSigError)
	}

	senderID, err := id.Unmarshal(msg.UID)
	if err != nil {
		jww.ERROR.Print("removeUser internal error unmarshalling "+
			"sender uid", err)
		return &messages.Ack{}, e
	}
	// Check the owner ID matches the sender ID
	if !senderID.Cmp(uid) {
		jww.ERROR.Print("removeUser internal error Auth Sender mismatch")
		return &messages.Ack{}, errors.New("Removal could not be " +
			"completed because you do not own this fact.")
	}

	// Delete the user, which deletes the facts for them as well.
	err = store.DeleteUser(msg.UID)
	if err != nil {
		jww.ERROR.Printf("Could not delete user: %+v", err)
		return &messages.Ack{}, err
	}

	// Ensure the dummy user is present in the system
	suser := storage.User{
		Id:        id.DummyUser.Bytes(),
		RsaPub:    "DUMMY KEY",
		DhPub:     nil,
		Salt:      nil,
		Signature: nil,
		Facts:     *new([]storage.Fact),
	}
	_ = store.InsertUser(&suser)

	// insert a dummy fact for the same username which prevents
	// reregistration.
	sfact := &storage.Fact{
		Hash:      hashFact,
		UserId:    id.DummyUser[:],
		Fact:      f.Fact,
		Type:      uint8(f.T),
		Signature: nil,
		Verified:  true,
		Timestamp: time.Time{},
	}
	err = store.InsertFact(sfact)
	if err != nil {
		jww.ERROR.Printf("Deleted user, but couldn't preserve "+
			"username: %+v", err)
		return &messages.Ack{}, err
	}

	jww.DEBUG.Printf("Deleted user with username %s", f.Fact)

	return &messages.Ack{}, nil
}
