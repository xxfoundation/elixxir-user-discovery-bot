package twilio

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/primitives/id"
)

// RegisterFact submits a fact for verification
func RegisterFact(uid *id.ID, fact string, factType uint8, verifier VerificationService) (string, error) {
	verifyId, err := verifier.Verification(fact, Channel(factType).String())
	if err != nil {
		return "", errors.WithStack(err)
	}

	// Adds entry to facts and verifications tables
	err = storage.UserDiscoveryDB.InsertFact(&storage.Fact{
		ConfirmationId:     verifyId,
		UserId:             uid.Marshal(),
		Fact:               fact,
		FactType:           factType,
		FactHash:           []byte("temphash"),
		Signature:          []byte("tempsig"),
		VerificationStatus: 0,
		Manual:             false,
		Code:               0,
	})
	// Makes call to Verification endpoint in twilio
	// Return the confirmation ID from db entry
	return verifyId, nil
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

	storage.UserDiscoveryDB.ConfirmFact()
	return valid, nil
}
