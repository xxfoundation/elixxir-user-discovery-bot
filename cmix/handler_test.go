package cmix

import (
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"testing"
)

// Test the start function on cmix manager
func TestManager_Start(t *testing.T) {
	m := &Manager{
		db: storage.NewTestDB(t),
	}
	t.Log(m)
	// err := m.Start()
	// if err != nil {
	// 	t.Errorf("Failed to start manager: %+v", err)
	// }
}
