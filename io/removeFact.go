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
)

// Takes in a FactRemovalRequest from a client and deletes the Fact if the client owns it
func removeFact(msg *pb.FactRemovalRequest, store *storage.Storage) (*messages.Ack, error) {
	// Generic copy of the internal error message
	e := errors.New("Removal could not be " +
		"completed do to internal error, please try again later")

	// Nil checks
	// Can we have a blank fact?
	if msg == nil || msg.RemovalData == nil || msg.UID == nil {
		return &messages.Ack{}, errors.New("Unable to parse required " +
			"fields in registration message")
	}

	// Generate the hash function and hash the fact
	f, err := fact.NewFact(fact.FactType(msg.RemovalData.FactType), msg.RemovalData.Fact)
	if err != nil {
		return &messages.Ack{}, errors.WithMessage(err, "Failed to create fact object")
	}
	hashFact := factID.Fingerprint(f)

	// Get the user who owns the fact
	users, err := store.Search([][]byte{hashFact})
	if err != nil {
		return &messages.Ack{}, err
	}
	if len(users) != 1 {
		jww.ERROR.Print("removeFact internal error users != 1")
		return &messages.Ack{}, e
	}
	// Unmarshal the owner ID
	uid, err := id.Unmarshal(users[0].Id)
	if err != nil {
		jww.ERROR.Print("removeFact internal error unmarshalling found user id", err)
		return &messages.Ack{}, e
	}

	hashedFact := factID.Fingerprint(f)

	// Parse the client's public key
	clientPubKey, err := rsa.LoadPublicKeyFromPem([]byte(users[0].RsaPub))
	if err != nil {
		return &messages.Ack{}, errors.WithMessage(err, invalidUserKeyError)
	}

	// Return an error if the fact signature cannot be verified
	err = rsa.Verify(clientPubKey, hash.CMixHash, hashedFact, msg.FactSig, nil)
	if err != nil {
		return &messages.Ack{}, errors.WithMessage(err, invalidFactSigError)
	}

	senderID, err := id.Unmarshal(msg.UID)
	if err != nil {
		jww.ERROR.Print("removeFact internal error unmarshalling sender uid", err)
		return &messages.Ack{}, e
	}
	// Check the owner ID matches the sender ID
	if !senderID.Cmp(uid) {
		jww.ERROR.Print("removeFact internal error Auth Sender mismatch")
		return &messages.Ack{}, errors.New("Removal could not be " +
			"completed because you do not own this fact.")
	}

	// Delete the fact
	err = store.DeleteFact(hashFact)
	if err != nil {
		jww.ERROR.Print("removeFact internal error store.DeleteHash", err)
		return &messages.Ack{}, e
	}

	return &messages.Ack{}, nil
}
