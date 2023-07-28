////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package io

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/factID"
	"gitlab.com/elixxir/crypto/hash"
	"gitlab.com/elixxir/crypto/registration"
	"gitlab.com/elixxir/primitives/fact"
	"gitlab.com/elixxir/user-discovery-bot/banned"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/elixxir/user-discovery-bot/validation"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"
	"time"
)

// Endpoint which handles a users attempt to register
func registerUser(msg *pb.UDBUserRegistration, permPublicKey *rsa.PublicKey,
	store *storage.Storage, bannedManager *banned.Manager,
	skipVerification bool) (*messages.Ack, error) {

	// Nil checks
	if msg == nil || msg.Frs == nil || msg.Frs.Fact == nil ||
		msg.IdentityRegistration == nil {
		return &messages.Ack{}, errors.New("Unable to parse required " +
			"fields in registration message")
	}

	// Parse the username and UserID
	username := msg.Frs.Fact.Fact // TODO: this & msg.IdentityRegistration.Username seems redundant
	uid, err := id.Unmarshal(msg.UID)
	if err != nil {
		return &messages.Ack{}, errors.New("Could not parse UID sent over. " +
			"Please try again")
	}

	jww.INFO.Printf("Registering user %s, %s", uid, username)

	canonicalUsername := validation.Canonicalize(username)

	// Check if username is valid
	if err := validation.IsValidUsername(canonicalUsername); err != nil {
		return nil, errors.Errorf("Username %q is invalid: %v", username, err)
	}

	// Check if the username is banned
	if bannedManager.IsBanned(canonicalUsername) {
		// Return same error message as if the user was already taken
		return &messages.Ack{}, errors.Errorf("Username %s is already taken. "+
			"Please try again", username)
	}

	// Check if username is taken
	err = store.CheckUser(canonicalUsername, uid)
	if err != nil {
		return &messages.Ack{}, errors.Errorf("Username %q is already taken. "+
			"Please try again", username)
	}

	if skipVerification {
		jww.WARN.Printf("Skipping scheduling signature verification for "+
			"user %s...", uid)
	} else {
		// Verify the Permissioning signature provided
		err = registration.VerifyWithTimestamp(permPublicKey, msg.Timestamp, msg.RSAPublicPem,
			msg.PermissioningSignature)
		if err != nil {
			return &messages.Ack{}, errors.Errorf(
				"Could not verify permissioning signature. "+
					"Data: %s, Signature: %s, %+v",
				msg.RSAPublicPem, msg.PermissioningSignature, err)
		}
	}

	// Parse the client's public key
	clientPubKey, err := rsa.LoadPublicKeyFromPem([]byte(msg.RSAPublicPem))
	if err != nil {
		return &messages.Ack{}, errors.New("Could not parse key passed in")
	}

	// Verify the signed fact
	tf, err := fact.NewFact(fact.FactType(msg.Frs.Fact.FactType), msg.Frs.Fact.Fact)
	if err != nil {
		return &messages.Ack{}, errors.WithMessage(err, "Failed to hash fact")
	}
	hashedFact := factID.Fingerprint(tf) // TODO: does fingerprint still need to uppercase the fact?
	err = rsa.Verify(clientPubKey, hash.CMixHash, hashedFact, msg.Frs.FactSig, nil)
	if err != nil {
		return &messages.Ack{}, errors.New("Could not verify fact signature")
	}

	canonicalFact, err := fact.NewFact(fact.FactType(msg.Frs.Fact.FactType), canonicalUsername)
	if err != nil {
		return &messages.Ack{}, errors.WithMessage(err, "Failed to hash canonicalUsername fact")
	}

	// Verify the signed identity data
	hashedIdentity := msg.IdentityRegistration.Digest()
	err = rsa.Verify(clientPubKey, hash.CMixHash, hashedIdentity, msg.IdentitySignature, nil)
	if err != nil {
		return &messages.Ack{}, errors.New("Could not verify identity signature")
	}

	// Create fact off of username
	f := storage.Fact{
		Hash:      factID.Fingerprint(canonicalFact),
		UserId:    msg.UID,
		Fact:      canonicalUsername,
		Type:      uint8(msg.Frs.Fact.FactType),
		Signature: msg.Frs.FactSig,
		Verified:  true,
		Timestamp: time.Now(),
	}

	// Create the user to insert into the database
	u := &storage.User{
		Id:                    msg.UID,
		Username:              username,
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

	jww.INFO.Printf("User Registered: %s, %s", uid, f.Fact)

	return &messages.Ack{}, nil
}
