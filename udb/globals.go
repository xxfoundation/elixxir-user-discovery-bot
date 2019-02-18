////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package udb

import (
	"gitlab.com/elixxir/user-discovery-bot/storage"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/client/globals"
	"io/ioutil"
	"log"
	"os"
	"gitlab.com/elixxir/primitives/userid"
)

// The User Discovery Bot's userid & registrationn code
// (this is global in cMix systems)
var UDB_USERID *userid.UserID = new(userid.UserID).SetUints(&[4]uint64{0, 0, 0, 3})

var DataStore storage.Storage

var Log = jww.NewNotepad(jww.LevelDebug, jww.LevelDebug, os.Stdout,
	ioutil.Discard, "CLIENT", log.Ldate|log.Ltime)

func init() {
	globals.Log = Log
}
