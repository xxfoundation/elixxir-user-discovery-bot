////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package io

import (
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/factID"
	"gitlab.com/elixxir/crypto/hash"
	"gitlab.com/elixxir/primitives/fact"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"
	"time"
)

// Endpoint which handles a users attempt to register
func registerUser(msg *pb.UDBUserRegistration, permPublicKey *rsa.PublicKey,
	store *storage.Storage, auth *connect.Auth) (*messages.Ack, error) {

	// Nil checks
	if msg == nil || msg.Frs == nil || msg.Frs.Fact == nil ||
		msg.IdentityRegistration == nil {
		return &messages.Ack{}, errors.New("Unable to parse required " +
			"fields in registration message")
	}

	// Ensure client is properly authenticated
	if !auth.IsAuthenticated || auth.Sender.IsDynamicHost() {
		return &messages.Ack{}, connect.AuthError(auth.Sender.GetId())
	}

	// Parse the username and UserID
	username := msg.IdentityRegistration.Username
	uid, err := id.Unmarshal(msg.UID)
	if err != nil {
		return &messages.Ack{}, errors.New("Could not parse UID sent over. " +
			"Please try again")
	}

	// Check if username is taken
	err = store.CheckUser(username, uid, msg.RSAPublicPem)
	if err != nil {
		return &messages.Ack{}, errors.Errorf("Username %s is already taken. "+
			"Please try again", username)
	}

	// Hash the rsa public key, what permissioning signature was signed off of
	h, err := hash.NewCMixHash()
	if err != nil {
		return &messages.Ack{}, errors.New("Could not verify signature due " +
			"to internal error. Please try again later")
	}
	h.Write([]byte(msg.RSAPublicPem))
	hashedRsaKey := h.Sum(nil)

	// Verify the Permissioning signature provided
	err = rsa.Verify(permPublicKey, hash.CMixHash, hashedRsaKey, msg.PermissioningSignature, nil)
	if err != nil {
		return &messages.Ack{}, errors.Errorf("Could not verify permissioning signature")
	}

	// Parse the client's public key
	clientPubKey, err := rsa.LoadPublicKeyFromPem([]byte(msg.RSAPublicPem))
	if err != nil {
		return &messages.Ack{}, errors.New("Could not parse key passed in")
	}

	// Verify the signed identity data
	hashedIdentity := msg.IdentityRegistration.Digest()
	err = rsa.Verify(clientPubKey, hash.CMixHash, hashedIdentity, msg.IdentitySignature, nil)
	if err != nil {
		return &messages.Ack{}, errors.New("Could not verify identity signature")
	}

	// Verify the signed fact
	tf, err := fact.NewFact(fact.FactType(msg.Frs.Fact.FactType), msg.Frs.Fact.Fact)
	if err != nil {
		return &messages.Ack{}, errors.WithMessage(err, "Failed to hash fact")
	}
	hashedFact := factID.Fingerprint(tf)
	err = rsa.Verify(clientPubKey, hash.CMixHash, hashedFact, msg.Frs.FactSig, nil)
	if err != nil {
		return &messages.Ack{}, errors.New("Could not verify fact signature")
	}

	// Create fact off of username
	f := storage.Fact{
		Hash:      hashedFact,
		UserId:    msg.UID,
		Fact:      msg.Frs.Fact.Fact,
		Type:      uint8(msg.Frs.Fact.FactType),
		Signature: msg.Frs.FactSig,
		Verified:  true,
		Timestamp: time.Now(),
	}

	// Create the user to insert into the database
	u := &storage.User{
		Id:        msg.UID,
		RsaPub:    msg.RSAPublicPem,
		DhPub:     msg.IdentityRegistration.DhPubKey,
		Salt:      msg.IdentityRegistration.Salt,
		Signature: msg.PermissioningSignature,
		Facts:     []storage.Fact{f},
	}

	// Insert the user into the database
	err = store.InsertUser(u)
	if err != nil {
		return &messages.Ack{}, errors.New("Could not register username due " +
			"to internal error. Please try again later")

	}

	return &messages.Ack{}, nil
}
