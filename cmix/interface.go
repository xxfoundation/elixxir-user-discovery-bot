package cmix

import (
	"gitlab.com/elixxir/client/cmix"
	"gitlab.com/elixxir/client/storage"
	"gitlab.com/elixxir/client/storage/user"
)

// client is a sub-interface of xxdk.Cmix containing methods relevant to
// this package.
type client interface {
	GetUser() user.Info
	GetCmix() cmix.Client
	GetStorage() storage.Session
}
