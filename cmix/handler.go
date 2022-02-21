package cmix

import (
	"gitlab.com/elixxir/client/single"
	"gitlab.com/elixxir/client/stoppable"
	"gitlab.com/elixxir/client/ud"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"time"
)

type SingleInterface interface {
	RegisterCallback(string, single.ReceiveComm)
	RespondSingleUse(single.Contact, []byte, time.Duration) error
	StartProcesses() (stoppable.Stoppable, error)
}

// CMIX Handler struct for user discovery
type Manager struct {
	db        *storage.Storage
	singleUse SingleInterface
}

// Create a CMIX Manager
func NewManager(singleUse *single.Manager, db *storage.Storage) *Manager {
	return &Manager{
		db:        db,
		singleUse: singleUse,
	}
}

// Start user discovery CMIX handler with a general callback that confirms all authenticated channel requests
func (m *Manager) Start() (stoppable.Stoppable, error) {
	// Register the lookup listener
	m.singleUse.RegisterCallback(ud.LookupTag, m.lookupCallback)

	// Register the search listener
	m.singleUse.RegisterCallback(ud.SearchTag, m.searchCallback)

	return m.singleUse.StartProcesses()
}
