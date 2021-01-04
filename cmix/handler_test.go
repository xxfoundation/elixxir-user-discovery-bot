package cmix

import (
	"gitlab.com/elixxir/client/auth"
	"gitlab.com/elixxir/client/interfaces"
	"gitlab.com/elixxir/client/interfaces/contact"
	"gitlab.com/elixxir/client/interfaces/message"
	"gitlab.com/elixxir/client/interfaces/params"
	"gitlab.com/elixxir/client/switchboard"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

/* Mock api.client so we can test the logic without one */
type MockClient struct {
	sw interfaces.Switchboard
	ar interfaces.Auth
}

func (mc *MockClient) GetSwitchboard() interfaces.Switchboard {
	return mc.sw
}
func (mc *MockClient) StartNetworkFollower() error {
	return nil
}
func (mc *MockClient) GetAuthRegistrar() interfaces.Auth {
	return mc.ar
}
func (mc *MockClient) ConfirmAuthenticatedChannel(contact2 contact.Contact) error {
	return nil
}
func (mc *MockClient) SendUnsafe(m message.Send, param params.Unsafe) ([]id.Round,
	error) {
	return nil, nil
}

// Test the start function on cmix manager
func TestManager_Start(t *testing.T) {
	c := &MockClient{
		sw: switchboard.New(),
	}
	c.ar = auth.NewManager(c.sw, nil, nil)
	m := &Manager{
		client:     c,
		lookupChan: make(chan message.Receive, 1000),
		searchChan: make(chan message.Receive, 1000),
		db:         storage.NewTestDB(t),
	}
	err := m.Start()
	if err != nil {
		t.Errorf("Failed to start manager: %+v", err)
	}
}
