package cmix

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/client/interfaces"
	"gitlab.com/elixxir/client/interfaces/contact"
	"gitlab.com/elixxir/client/interfaces/message"
	"gitlab.com/elixxir/client/interfaces/params"
	"gitlab.com/elixxir/client/switchboard"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/primitives/id"
)

// CMIX Handler struct for user discovery
type Manager struct {
	client     ClientInterface
	lookupChan chan message.Receive
	searchChan chan message.Receive
	db         *storage.Storage
}

type ClientInterface interface {
	GetSwitchboard() interfaces.Switchboard
	StartNetworkFollower() error
	GetAuthRegistrar() interfaces.Auth
	ConfirmAuthenticatedChannel(contact2 contact.Contact) error
	SendUnsafe(m message.Send, param params.Unsafe) ([]id.Round,
		error)
}

// Create a CMIX Manager
func NewManager(storagedir string, password []byte, db *storage.Storage) (*Manager, error) {
	m := &Manager{
		client:     nil,
		lookupChan: make(chan message.Receive, 1000),
		searchChan: make(chan message.Receive, 1000),
		db:         db,
	}
	var err error
	m.client, err = api.Login(storagedir, password, params.GetDefaultNetwork())

	return m, err
}

// Start user discovery CMIX handler with a general callback that confirms all authenticated channel requests
func (m *Manager) Start() error {
	var err error

	//register the lookup listener
	m.client.GetSwitchboard().RegisterChannel("UDLookup",
		switchboard.AnyUser(), message.UdLookup, m.lookupChan)

	//register the search listener
	m.client.GetSwitchboard().RegisterChannel("UDSearch",
		switchboard.AnyUser(), message.UdSearch, m.searchChan)

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
