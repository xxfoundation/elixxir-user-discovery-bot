////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package udb

import (
	"gitlab.com/privategrity/client/user"
	"gitlab.com/privategrity/user-discovery-bot/storage"
)

// The User Discovery Bot's userid & registrationn code
// (this is global in cMix systems)
const UDB_USERID = user.ID(13)
const UDB_NICK = "UDB"

var DataStore storage.Storage
