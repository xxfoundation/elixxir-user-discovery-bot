package cmix

import (
	"gitlab.com/elixxir/client/single"
	"gitlab.com/elixxir/client/ud"
	"gitlab.com/elixxir/client/xxdk"
	"gitlab.com/elixxir/user-discovery-bot/storage"
)

// Manager struct for user discovery single use
type Manager struct {
	db             *storage.Storage
	e2eClient      *xxdk.E2e
	lookupListener single.Listener
	searchListener single.Listener
}

// NewManager creates a CMIX Manager
func NewManager(e2eClient *xxdk.E2e,
	db *storage.Storage) *Manager {
	return &Manager{
		db:        db,
		e2eClient: e2eClient,
	}
}

// Start user discovery CMIX handler for single use messages
func (m *Manager) Start() {
	// Register the lookup listener
	m.lookupListener = single.Listen(ud.LookupTag,
		m.e2eClient.GetUser().ReceptionID,
		m.e2eClient.GetUser().E2eDhPrivateKey,
		m.e2eClient.GetCmix(),
		m.e2eClient.GetStorage().GetE2EGroup(),
		&lookupManager{m: m})

	// Register the search listener
	m.searchListener = single.Listen(ud.SearchTag,
		m.e2eClient.GetUser().ReceptionID,
		m.e2eClient.GetUser().E2eDhPrivateKey,
		m.e2eClient.GetCmix(),
		m.e2eClient.GetStorage().GetE2EGroup(),
		&searchManager{m: m})
}

// Stop the user discovery cmix handler
func (m *Manager) Stop() {
	m.searchListener.Stop()
	m.lookupListener.Stop()
}
