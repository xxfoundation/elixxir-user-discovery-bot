////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package udb

import (
	"gitlab.com/elixxir/user-discovery-bot/storage"
	jww "github.com/spf13/jwalterweatherman"
	"io/ioutil"
	"log"
	"os"
	"gitlab.com/elixxir/client/globals"
	"gitlab.com/elixxir/primitives/userid"
)

// The User Discovery Bot's userid & registrationn code
// (this is global in cMix systems)
var UDB_USERID *id.UserID = new(id.UserID).SetUints(&[4]uint64{0, 0, 0, 13})

var DataStore storage.Storage

var Log = jww.NewNotepad(jww.LevelDebug, jww.LevelDebug, os.Stdout,
	ioutil.Discard, "CLIENT", log.Ldate|log.Ltime)

func init() {
	globals.Log = Log
}