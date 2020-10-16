package twilio

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/crypto/hasher"
	"gitlab.com/xx_network/primitives/id"
)

// RegisterFact submits a fact for verification
func RegisterFact(uid *id.ID, fact string, factType uint8, signature []byte, verifier VerificationService) (string, error) {
	verifyId, err := verifier.Verification(fact, Channel(factType).String())
	if err != nil {
		return "", errors.WithStack(err)
	}
	h := hasher.SHA3_256.New()
	h.Write([]byte(fact))

	// Adds entry to facts and verifications tables
	err = storage.UserDiscoveryDB.InsertFactTwilio(uid.Marshal(), h.Sum(nil), signature, fact, uint(factType), verifyId)
	// Makes call to Verification endpoint in twilio
	// Return the confirmation ID from db entry
	return verifyId, err
}

// ConfirmFact confirms a code and completes fact verification
func ConfirmFact(confirmationID string, code int, verifier VerificationService) (bool, error) {
	// Get verifications entry by confirmation id
	// Make call to verification check endpoint with code
	// If good, update verification to status 2 (complete)
	// If not, update verification to status 1, return error
	valid, err := verifier.VerificationCheck(code, confirmationID)
	if err != nil {
		return false, err
	}
	if valid {
		err = storage.UserDiscoveryDB.VerifyFactTwilio(confirmationID)
		return valid, err
	} else {
		return valid, nil
	}
}
