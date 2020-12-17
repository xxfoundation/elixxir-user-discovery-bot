package cmix

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/client/interfaces/contact"
	"gitlab.com/elixxir/client/interfaces/message"
	"gitlab.com/elixxir/client/switchboard"
	"gitlab.com/elixxir/user-discovery-bot/storage"
)

// CMIX Handler struct for user discovery
type Manager struct {
	client     *api.Client
	lookupChan chan message.Receive
	searchChan chan message.Receive
	db         *storage.Storage
}

// Start user discovery CMIX handler with a general callback that confirms all authenticated channel requests
func NewManager(storagedir string, password []byte) error {

	m := &Manager{
		client:     nil,
		lookupChan: make(chan message.Receive, 1000),
		searchChan: make(chan message.Receive, 1000),
	}

	var err error

	m.client, err = api.Login(storagedir, password)
	if err != nil {
		return err
	}

	//register the lookup listener
	m.client.GetSwitchboard().RegisterChannel("UDLookup",
		switchboard.AnyUser(), message.UdLookup, m.lookupChan)

	err = m.client.StartNetworkFollower()
	if err != nil {
		return err
	}

	// Create and register authenticated channel receiver
	registrar := m.client.GetAuthRegistrar()
	rcb := func(requestor contact.Contact, message string) {
		err := m.client.ConfirmAuthenticatedChannel(requestor)
		if err != nil {
			jww.ERROR.Println(err)
		}
	}
	registrar.AddGeneralRequestCallback(rcb)

	go m.LookupProcess()
	go m.SearchProcess()

	return nil
}
