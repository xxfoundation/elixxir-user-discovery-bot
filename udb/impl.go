////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package udb

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/udb"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/comms/messages"
)

func NewImplementation() *udb.Implementation  {
	impl := udb.NewImplementation()

	impl.Functions.RegisterUser = func(registration *pb.UDBUserRegistration) (*messages.Ack, error) {
		return RegisterUser(registration, clientObj, storage.UserDiscoveryDB)
	}

	return impl
}