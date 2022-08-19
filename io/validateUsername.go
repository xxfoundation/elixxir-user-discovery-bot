///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package io

import (
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/partnerships/crust"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"
	"io"
)

const (
	usernameNotAssociatedWithUser = "username %s is not associated with user ID %s"
)

// ValidateUsername validates that a user owns a username by signing the contents of the
// mixmessages.UsernameValidationRequest.
func validateUsername(request *pb.UsernameValidationRequest,
	store *storage.Storage, privKey *rsa.PrivateKey, rng io.Reader) (*pb.UsernameValidation, error) {
	// Return an error if the request is invalid
	if request == nil || request.Username == "" || request.ReceptionPublicKeyPem == nil {
		return &pb.UsernameValidation{}, errors.New("Unable to parse required " +
			"fields in registration message")
	}

	// Marshal user ID
	userID, err := id.Unmarshal(request.UserId)
	if err != nil {
		return &pb.UsernameValidation{}, errors.WithMessage(err, invalidUserIdError)
	}

	// Return an error if the user is not registered
	user, err := store.GetUser(request.UserId)
	if err != nil {
		return &pb.UsernameValidation{}, errors.WithMessage(err, getUserFailureError)
	} else if user == nil {
		return &pb.UsernameValidation{}, errors.Errorf(noUserError,
			userID)
	}

	// Check that the username in storage matches the username in the request message
	// before validating
	if user.Username != request.Username {
		return &pb.UsernameValidation{}, errors.Errorf(usernameNotAssociatedWithUser,
			request.Username, userID)
	}

	// Create a signature verifying the user owns their username
	verificationSignature, err := crust.SignVerification(rng, privKey,
		request.Username, request.ReceptionPublicKeyPem)
	if err != nil {
		return nil, errors.Errorf("Failed to create verification signature: %v", err)
	}

	// Return signature to user
	return &pb.UsernameValidation{Signature: verificationSignature}, nil
}
