////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Handles Manager interface for the Twilio layer

package twilio

import (
	"gitlab.com/elixxir/user-discovery-bot/interfaces/params"
	"gitlab.com/elixxir/user-discovery-bot/storage"
)

type Manager struct {
	storage  *storage.Storage
	verifier VerificationService
}

// Create a new Manager given Twilio params and the Storage interface
func NewManager(twilio params.Twilio, storage *storage.Storage) *Manager {
	return &Manager{
		storage:  storage,
		verifier: &verifier{p: twilio},
	}
}

// Create a new Manager given the Storage interface
func NewMockManager(storage *storage.Storage) *Manager {
	return &Manager{
		storage:  storage,
		verifier: newMockVerifier(),
	}
}
