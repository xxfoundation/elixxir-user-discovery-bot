package cmix

import (
	"git.xx.network/elixxir/client/single"
	"git.xx.network/elixxir/client/stoppable"
	"git.xx.network/elixxir/client/ud"
	"git.xx.network/elixxir/user-discovery-bot/storage"
	"time"
)

type SingleInterface interface {
	RegisterCallback(string, single.ReceiveComm)
	RespondSingleUse(single.Contact, []byte, time.Duration) error
	StartProcesses() stoppable.Stoppable
}

// CMIX Handler struct for user discovery
type Manager struct {
	db        *storage.Storage
	singleUse SingleInterface
}

type SingleUseInterface interface {
}

// Create a CMIX Manager
func NewManager(singleUse *single.Manager, db *storage.Storage) *Manager {
	return &Manager{
		db:        db,
		singleUse: singleUse,
	}
}

// Start user discovery CMIX handler with a general callback that confirms all authenticated channel requests
func (m *Manager) Start() stoppable.Stoppable {
	// Register the lookup listener
	m.singleUse.RegisterCallback(ud.LookupTag, m.lookupCallback)

	// Register the search listener
	m.singleUse.RegisterCallback(ud.SearchTag, m.searchCallback)

	return m.singleUse.StartProcesses()
}
