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
	"gitlab.com/elixxir/user-discovery-bot/interfaces/params"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/elixxir/user-discovery-bot/twilio"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"
)

// The main UserDiscovery instance object
type Manager struct {
	Comms                  *udb.Comms
	PermissioningPublicKey *rsa.PublicKey
	Storage                *storage.Storage
	Twilio                 *twilio.Manager
}

// Create a new UserDiscovery Manager given a set of Params
func NewManager(p params.IO, id *id.ID, permissioningCert *rsa.PublicKey, twilio *twilio.Manager, storage *storage.Storage) *Manager {
	m := &Manager{
		Storage:                storage,
		PermissioningPublicKey: permissioningCert,
		Twilio:                 twilio,
	}
	m.Comms = udb.StartServer(id, fmt.Sprintf("0.0.0.0:%s", p.Port),
		newImplementation(m), p.Cert, p.Key)
	return m
}

// Create a new Comms implementation for UserDiscovery
func newImplementation(m *Manager) *udb.Implementation {
	impl := udb.NewImplementation()

	impl.Functions.RegisterUser = func(registration *pb.UDBUserRegistration, auth *connect.Auth) (*messages.Ack, error) {
		return registerUser(registration, m.PermissioningPublicKey, m.Storage, auth)
	}

	impl.Functions.RegisterFact = func(request *pb.FactRegisterRequest, auth *connect.Auth) (*pb.FactRegisterResponse, error) {
		return registerFact(request, m.Twilio, m.Storage, auth)
	}

	impl.Functions.ConfirmFact = func(request *pb.FactConfirmRequest, auth *connect.Auth) (*messages.Ack, error) {
		return confirmFact(request, m.Twilio, m.Storage, auth)
	}

	impl.Functions.RemoveFact = func(msg *pb.FactRemovalRequest, auth *connect.Auth) (*messages.Ack, error) {
		return removeFact(msg, m.Storage, auth)
	}

	return impl
}
