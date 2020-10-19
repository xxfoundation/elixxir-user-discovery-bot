package twilio

import (
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

func TestRegisterFact(t *testing.T) {
	mockDb, _, err := storage.NewDatabase("", "", "", "", "11")
	if err != nil {
		t.Errorf("Failed to init mock db: %+v", err)
	}
	uid := id.NewIdFromString("zezima", id.User, t)
	_, err = RegisterFact(uid, "water is wet", 0, []byte("hancock"), MV, mockDb)
	if err != nil {
		t.Errorf("Failed to register fact: %+v", err)
	}
}

func TestConfirmFact(t *testing.T) {
	mockDb, _, err := storage.NewDatabase("", "", "", "", "11")
	if err != nil {
		t.Errorf("Failed to init mock db: %+v", err)
	}
	uid := id.NewIdFromString("zezima", id.User, t)
	confId, err := RegisterFact(uid, "water is wet", 0, []byte("hancock"), MV, mockDb)
	if err != nil {
		t.Errorf("Failed to register fact: %+v", err)
	}

	_, err = ConfirmFact(confId, 1234, MV, mockDb)
	if err != nil {
		t.Errorf("Failed to confirm fact: %+v", err)
	}
}
