////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package cmix

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/client/globals"
	"io/ioutil"
	"log"
	"os"
)

var Log = jww.NewNotepad(jww.LevelDebug, jww.LevelDebug, os.Stdout,
	ioutil.Discard, "CLIENT", log.Ldate|log.Ltime)
var BannedUsernameList BlackList

func init() {
	globals.Log = Log
}
