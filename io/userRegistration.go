////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package io

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "git.xx.network/elixxir/comms/mixmessages"
	"git.xx.network/elixxir/crypto/factID"
	"git.xx.network/elixxir/crypto/hash"
	"git.xx.network/elixxir/crypto/registration"
	"git.xx.network/elixxir/primitives/fact"
	"git.xx.network/elixxir/user-discovery-bot/storage"
	"git.xx.network/xx_network/comms/connect"
	"git.xx.network/xx_network/comms/messages"
	"git.xx.network/xx_network/crypto/signature/rsa"
	"git.xx.network/xx_network/primitives/id"
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
	if !auth.IsAuthenticated {
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

	// Verify the Permissioning signature provided
	err = registration.VerifyWithTimestamp(permPublicKey, msg.Timestamp, msg.RSAPublicPem, msg.PermissioningSignature)
	if err != nil {
		return &messages.Ack{}, errors.Errorf(
			"Could not verify permissioning signature. "+
				"Data: %s, Signature: %s, %+v",
			msg.RSAPublicPem, msg.PermissioningSignature, err)
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
		Id:                    msg.UID,
		RsaPub:                msg.RSAPublicPem,
		DhPub:                 msg.IdentityRegistration.DhPubKey,
		Salt:                  msg.IdentityRegistration.Salt,
		Signature:             msg.PermissioningSignature,
		RegistrationTimestamp: time.Unix(0, msg.Timestamp),
		Facts:                 []storage.Fact{f},
	}

	// Insert the user into the database
	err = store.InsertUser(u)
	if err != nil {
		return &messages.Ack{}, errors.New("Could not register username due " +
			"to internal error. Please try again later")

	}

	jww.INFO.Printf("User Registered: %s", uid)

	return &messages.Ack{}, nil
}
