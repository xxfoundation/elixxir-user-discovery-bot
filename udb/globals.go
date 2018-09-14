////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package udb

import (
	"gitlab.com/privategrity/user-discovery-bot/storage"
	"gitlab.com/privategrity/crypto/id"
)

// The User Discovery Bot's userid & registrationn code
// (this is global in cMix systems)
var UDB_USERID = new(id.UserID).SetUints(&[4]uint64{0, 0, 0, 13})

var DataStore storage.Storage
