///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package io

import (
	"crypto/rand"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/partnerships/crust"
	"gitlab.com/elixxir/user-discovery-bot/banned"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/primitives/ndf"
	"testing"
	"time"
)

// Happy path.
func TestValidateUsername(t *testing.T) {
	// Initialize fields needed for testing
	clientId, rsaPrivKey := initClientFields(t)
	store := storage.NewTestDB(t)
	ndfObj, _ := ndf.Unmarshal(getNDF())
	cert, err := loadPermissioningPubKey(ndfObj.Registration.TlsCertificate)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Register user first -------------------------------------------------------------------------
	testTime, err := time.Parse(time.RFC3339,
		"2012-12-21T22:08:41+00:00")
	if err != nil {
		t.Fatalf("Could not parse precanned time: %v", err.Error())
	}
	registerMsg, err := buildUserRegistrationMessage(clientId, rsaPrivKey, testTime, t)
	if err != nil {
		t.Fatalf("Failed to build registration message: %+v", err)
	}

	bannedManager, err := banned.NewManager("", "")
	if err != nil {
		t.Fatalf("Failed to construct ban manager: %v", err)
	}

	_, err = registerUser(registerMsg, cert, store, bannedManager)
	if err != nil {
		t.Errorf("Failed happy path: %v", err)
	}

	// Test Validate username ----------------------------------------------------------------------
	username := registerMsg.Frs.Fact.Fact
	pubKeyPem := []byte(registerMsg.RSAPublicPem)
	validationRequest := &pb.UsernameValidationRequest{
		UserId: registerMsg.UID,
	}

	validationResponse, err := validateUsername(validationRequest, store, rsaPrivKey, rand.Reader)
	if err != nil {
		t.Fatalf("Failed to validate username: %+v", err)
	}

	err = crust.VerifyVerificationSignature(rsaPrivKey.GetPublic(),
		username, pubKeyPem, validationResponse.Signature)
	if err != nil {
		t.Fatalf("validateUsername did not return a valid signature!")
	}

}

// Error path: Try to validate a username that does not belong to the user.
func TestValidateUsername_UsernameMismatch(t *testing.T) {
	// Initialize fields needed for testing
	clientId, rsaPrivKey := initClientFields(t)
	store := storage.NewTestDB(t)
	ndfObj, _ := ndf.Unmarshal(getNDF())
	cert, err := loadPermissioningPubKey(ndfObj.Registration.TlsCertificate)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Register user first -------------------------------------------------------------------------
	testTime, err := time.Parse(time.RFC3339,
		"2012-12-21T22:08:41+00:00")
	if err != nil {
		t.Fatalf("Could not parse precanned time: %v", err.Error())
	}
	registerMsg, err := buildUserRegistrationMessage(clientId, rsaPrivKey, testTime, t)
	if err != nil {
		t.Fatalf("Failed to build registration message: %+v", err)
	}

	bannedManager, err := banned.NewManager("", "")
	if err != nil {
		t.Fatalf("Failed to construct ban manager: %v", err)
	}

	_, err = registerUser(registerMsg, cert, store, bannedManager)
	if err != nil {
		t.Errorf("Failed happy path: %v", err)
	}

}
