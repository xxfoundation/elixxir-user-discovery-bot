////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package io

import (
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/factID"
	"gitlab.com/elixxir/crypto/hash"
	"gitlab.com/elixxir/primitives/fact"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/elixxir/user-discovery-bot/twilio"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"
)

const (
	invalidFactRegisterRequestError = "Unable to parse required fields in FactRegisterRequest."
	factExistsError                 = "Cannot register fact that already exists."
	noUserError                     = "User associated with fact not registered: %s"
	invalidUserKeyError             = "Could not parse user's key."
	invalidFactSigError             = "Failed to verify fact signature."
	getUserFailureError             = "Failed to find user"
	invalidUserIdError              = "Failed to parse user ID."
	twilioRegFailureError           = "Failed to register fact with Twilio."

	invalidFactConfirmRequestError = "Unable to parse required fields in FactConfirmRequest."
	invalidFactCodeError           = "Failed to parse the FactConfirmRequest code."
	twilioConfirmFailureError      = "Failed to confirm fact with Twilio"
	twilioVerificationFailureError = "Twilio verification failed."
	nicknameFactError              = "Cannot register nickname type facts"
)

// registerFact is an endpoint that attempts to register a user's fact.
func registerFact(request *pb.FactRegisterRequest,
	verifier *twilio.Manager, store *storage.Storage) (*pb.FactRegisterResponse, error) {

	// Return an error if the request is invalid
	if request == nil || request.Fact == nil {
		return &pb.FactRegisterResponse{}, errors.New(invalidFactRegisterRequestError)
	}

	if fact.FactType(request.Fact.FactType) == fact.Nickname {
		return &pb.FactRegisterResponse{}, errors.New(nicknameFactError)
	}

	f, err := fact.NewFact(fact.FactType(request.Fact.FactType), request.Fact.Fact)
	if err != nil {
		return &pb.FactRegisterResponse{}, err
	}

	// Return an error if the fact is already registered
	hashedFact := factID.Fingerprint(f)
	ret, err := store.Search([][]byte{hashedFact})
	if err != nil {
		return &pb.FactRegisterResponse{}, err
	}
	if len(ret) != 0 {
		return &pb.FactRegisterResponse{}, errors.New(factExistsError)
	}

	// Marshal user ID
	userID, err := id.Unmarshal(request.UID)
	if err != nil {
		return &pb.FactRegisterResponse{}, errors.WithMessage(err, invalidUserIdError)
	}

	// Return an error if the fact's user is not registered
	user, err := store.GetUser(request.UID)
	if err != nil {
		return &pb.FactRegisterResponse{}, errors.WithMessage(err, getUserFailureError)
	} else if user == nil {
		return &pb.FactRegisterResponse{}, errors.Errorf(noUserError,
			userID)
	}

	// Parse the client's public key
	clientPubKey, err := rsa.LoadPublicKeyFromPem([]byte(user.RsaPub))
	if err != nil {
		return &pb.FactRegisterResponse{}, errors.WithMessage(err, invalidUserKeyError)
	}

	// Return an error if the fact signature cannot be verified
	err = rsa.Verify(clientPubKey, hash.CMixHash, hashedFact, request.FactSig, nil)
	if err != nil {
		return &pb.FactRegisterResponse{}, errors.WithMessage(err, invalidFactSigError)
	}

	// Register fact with Twilio to get confirmation ID
	confirmationID, err := verifier.RegisterFact(userID, request.Fact.Fact,
		uint8(request.Fact.FactType), request.FactSig)
	if err != nil {
		return &pb.FactRegisterResponse{}, errors.WithMessage(err, twilioRegFailureError)
	}

	// Create response
	response := &pb.FactRegisterResponse{
		ConfirmationID: confirmationID,
	}

	return response, nil
}

// confirmFact verifies the fact via Twilio and sets the fact in the database as
// confirmed.
func confirmFact(request *pb.FactConfirmRequest, verifier *twilio.Manager) (*messages.Ack, error) {

	// Return an error if the request is nil
	if request == nil || request.ConfirmationID == "" {
		return &messages.Ack{}, errors.New(invalidFactConfirmRequestError)
	}

	valid, err := verifier.ConfirmFact(request.ConfirmationID, request.Code)
	if err != nil {
		return &messages.Ack{}, errors.Errorf(twilioConfirmFailureError+": %+v", err)
	} else if !valid {
		return &messages.Ack{}, errors.New(twilioVerificationFailureError)
	}

	return &messages.Ack{}, nil
}
