////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Main methods for registering & confirming facts with twilio

package twilio

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/ttacon/libphonenumber"
	"gitlab.com/elixxir/crypto/factID"
	fact2 "gitlab.com/elixxir/primitives/fact"
	"gitlab.com/xx_network/primitives/id"
)

// RegisterFact submits a fact for verification
func (m *Manager) RegisterFact(uid *id.ID, fact string, factType uint8, signature []byte) (string, error) {
	var to string
	var channel string
	if fact2.FactType(factType) == fact2.Phone {
		l := len(fact)
		number := fact[:l-2]
		countryCode := fact[l-2:]
		num, err := libphonenumber.Parse(number, countryCode)
		if err != nil {
			return "", err
		}
		// Phone numbers sent to twilio must be in e.164 format
		to = libphonenumber.Format(num, libphonenumber.E164)
		channel = SMS.String()
	} else {
		to = fact
		channel = Email.String()
	}

	verifyId, err := m.verifier.Verification(to, channel)
	jww.INFO.Printf("Sent verification & received %s", verifyId)
	if err != nil {
		err = errors.WithMessage(err, "Twilio verification init failed")
		jww.ERROR.Println(err)
		return "", err
	}

	f, err := fact2.NewFact(fact2.FactType(factType), fact)
	if err != nil {
		return "", errors.WithMessage(err, "Failed to hash fact object")
	}
	factId := factID.Fingerprint(f)

	// Adds entry to facts and verifications tables
	err = m.storage.InsertFactTwilio(uid.Marshal(), factId, signature, uint(factType), verifyId)
	// Makes call to Verification endpoint in twilio
	// Return the confirmation ID from db entry
	return verifyId, err
}

// ConfirmFact confirms a code and completes fact verification
func (m *Manager) ConfirmFact(confirmationID string, code string) (bool, error) {
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
