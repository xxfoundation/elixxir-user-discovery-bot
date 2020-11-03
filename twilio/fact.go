////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Main methods for registering & confirming facts with twilio

package twilio

import (
	"github.com/pkg/errors"
	"gitlab.com/xx_network/crypto/hasher"
	"gitlab.com/xx_network/primitives/id"
)

// RegisterFact submits a fact for verification
func (m *Manager) RegisterFact(uid *id.ID, fact string, factType uint8, signature []byte) (string, error) {
	verifyId, err := m.verifier.Verification(fact, Channel(factType).String())
	if err != nil {
		return "", errors.WithStack(err)
	}
	h := hasher.SHA3_256.New()
	h.Write([]byte(fact))

	// Adds entry to facts and verifications tables
	err = m.storage.InsertFactTwilio(uid.Marshal(), h.Sum(nil), signature, uint(factType), fact, verifyId)
	// Makes call to Verification endpoint in twilio
	// Return the confirmation ID from db entry
	return verifyId, err
}

// ConfirmFact confirms a code and completes fact verification
func (m *Manager) ConfirmFact(confirmationID string, code int) (bool, error) {
	// Make call to verification check endpoint with code
	valid, err := m.verifier.VerificationCheck(code, confirmationID)
	if err != nil {
		return false, err
	}
	// If good, verify associated fact
	if valid {
		err = m.storage.MarkTwilioFactVerified(confirmationID)
		return valid, err
	} else {
		return valid, nil
	}
}
