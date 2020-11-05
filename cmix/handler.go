package cmix

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/client/interfaces/contact"
)

// CMIX Handler struct for user discovery
type UD struct {
	client *api.Client
}

// Start user discovery CMIX handler with a general callback that confirms all authenticated channel requests
func (ud *UD) Start(storagedir string, password []byte) error {
	c, err := api.Login(storagedir, password)
	if err != nil {
		return err
	}

	err = c.StartNetworkFollower()
	if err != nil {
		return err
	}

	// Create and register authenticated channel receiver
	registrar := ud.client.GetAuthRegistrar()
	rcb := func(requestor contact.Contact, message string) {
		err := c.ConfirmAuthenticatedChannel(requestor)
		if err != nil {
			jww.ERROR.Println(err)
		}
	}
	registrar.AddGeneralRequestCallback(rcb)

	ud.client = c
	return nil
}
