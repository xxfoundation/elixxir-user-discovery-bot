package io

import (
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
	store, _, err := storage.NewStorage(params.Database{})
	if err != nil {
		t.Errorf("Failed to create storage")
	}
	tm := twilio.NewMockManager(store)

	m := NewManager(p, id.NewIdFromString("zezima", id.User, t), nil, tm, store)
	if m == nil || reflect.TypeOf(m) != reflect.TypeOf(&Manager{}) {
		t.Errorf("Did not receive a manager")
	}
}
