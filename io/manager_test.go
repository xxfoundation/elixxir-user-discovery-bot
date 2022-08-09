package io

import (
	"crypto/ed25519"
	"gitlab.com/elixxir/user-discovery-bot/banned"
	"gitlab.com/elixxir/user-discovery-bot/interfaces/params"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/elixxir/user-discovery-bot/twilio"
	"gitlab.com/xx_network/crypto/csprng"
	"gitlab.com/xx_network/primitives/id"
	"reflect"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	p := params.IO{
		Cert: nil,
		Key:  nil,
		Port: "",
	}
	store := storage.NewTestDB(t)
	tm := twilio.NewMockManager(store)
	bannedManager, err := banned.NewManager("", "")
	if err != nil {
		t.Fatalf("Failed to construct ban manager: %v", err)
	}

	_, edPriv, err := ed25519.GenerateKey(csprng.NewSystemRNG())

	m := NewManager(p, id.NewIdFromString("zezima", id.User, t), nil, params.Channels{
		Enabled:          true,
		LeaseTime:        500 * time.Hour,
		LeaseGracePeriod: 24 * time.Hour,
		Ed25519Key:       edPriv,
	}, tm, bannedManager, store, false)
	if m == nil || reflect.TypeOf(m) != reflect.TypeOf(&Manager{}) {
		t.Errorf("Did not receive a manager")
	}
}
