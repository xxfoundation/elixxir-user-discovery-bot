///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Handles Manager interface for the IO layer

package io

import (
	"fmt"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/udb"
	"gitlab.com/elixxir/crypto/fastRNG"
	"gitlab.com/elixxir/user-discovery-bot/banned"
	"gitlab.com/elixxir/user-discovery-bot/interfaces/params"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/elixxir/user-discovery-bot/twilio"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"
)

// Manager is the main UserDiscovery instance object.
type Manager struct {
	Comms                  *udb.Comms
	rng                    *fastRNG.StreamGenerator
	PermissioningPublicKey *rsa.PublicKey
	Storage                *storage.Storage
	Twilio                 *twilio.Manager
	Banned                 *banned.Manager
	rsaPrivateKey          *rsa.PrivateKey
	skipVerification       bool
}

// NewManager creates a new UserDiscovery Manager given a set of Params.
func NewManager(p params.IO, id *id.ID, rsaPrivateKey *rsa.PrivateKey,
	permissioningCert *rsa.PublicKey,
	twilio *twilio.Manager, banned *banned.Manager,
	storage *storage.Storage, rng *fastRNG.StreamGenerator,
	skipVerification bool) *Manager {

	m := &Manager{
		Storage:                storage,
		PermissioningPublicKey: permissioningCert,
		Twilio:                 twilio,
		Banned:                 banned,
		skipVerification:       skipVerification,
		rsaPrivateKey:          rsaPrivateKey,
		rng:                    rng,
	}

	m.Comms = udb.StartServer(id, fmt.Sprintf("0.0.0.0:%s", p.Port),
		newImplementation(m), p.Cert, p.Key)
	return m
}

// Create a new Comms implementation for UserDiscovery
func newImplementation(m *Manager) *udb.Implementation {
	impl := udb.NewImplementation()

	impl.Functions.RegisterUser =
		func(registration *pb.UDBUserRegistration) (*messages.Ack, error) {
			return registerUser(registration, m.PermissioningPublicKey, m.Storage,
				m.Banned, m.skipVerification)
		}

	impl.Functions.RemoveUser = func(msg *pb.FactRemovalRequest) (*messages.Ack, error) {
		return removeUser(msg, m.Storage)
	}

	impl.Functions.RegisterFact = func(request *pb.FactRegisterRequest) (*pb.FactRegisterResponse, error) {
		return registerFact(request, m.Twilio, m.Storage)
	}

	impl.Functions.ConfirmFact = func(request *pb.FactConfirmRequest) (*messages.Ack, error) {
		return confirmFact(request, m.Twilio)
	}

	impl.Functions.RemoveFact = func(msg *pb.FactRemovalRequest) (*messages.Ack, error) {
		return removeFact(msg, m.Storage)
	}

	impl.Functions.ValidateUsername = func(request *pb.UsernameValidationRequest) (*pb.UsernameValidation, error) {
		stream := m.rng.GetStream()
		defer stream.Close()
		return validateUsername(request, m.Storage, m.rsaPrivateKey, stream)
	}

	return impl
}
