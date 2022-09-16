package io

import (
	"gitlab.com/elixxir/user-discovery-bot/banned"
	"gitlab.com/elixxir/user-discovery-bot/interfaces/params"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/elixxir/user-discovery-bot/twilio"
	"gitlab.com/xx_network/primitives/id"
	"reflect"
	"testing"
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

	m := NewManager(p, id.NewIdFromString("zezima", id.User, t), nil, nil, tm, bannedManager, store, nil)
	if m == nil || reflect.TypeOf(m) != reflect.TypeOf(&Manager{}) {
		t.Errorf("Did not receive a manager")
	}
}
