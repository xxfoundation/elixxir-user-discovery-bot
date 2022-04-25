package cmix

import (
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/client/single"
	"gitlab.com/elixxir/client/ud"
	"gitlab.com/elixxir/user-discovery-bot/storage"
)

// Manager struct for user discovery single use
type Manager struct {
	db             *storage.Storage
	cl             *api.Client
	lookupListener single.Listener
	searchListener single.Listener
}

// NewManager creates a CMIX Manager
func NewManager(client *api.Client, db *storage.Storage) *Manager {
	return &Manager{
		db: db,
		cl: client,
	}
}

// Start user discovery CMIX handler for single use messages
func (m *Manager) Start() {
	// Register the lookup listener
	m.lookupListener = single.Listen(ud.LookupTag, m.cl.GetE2EHandler().GetReceptionID(), m.cl.GetUser().E2eDhPublicKey, m.cl.GetNetworkInterface(), m.cl.GetE2EHandler().GetGroup(), &lookupManager{m: m})

	// Register the search listener
	m.searchListener = single.Listen(ud.LookupTag, m.cl.GetE2EHandler().GetReceptionID(), m.cl.GetUser().E2eDhPublicKey, m.cl.GetNetworkInterface(), m.cl.GetE2EHandler().GetGroup(), &searchManager{m: m})
}

// Stop the user discovery cmix handler
func (m *Manager) Stop() {
	m.searchListener.Stop()
	m.lookupListener.Stop()
}
