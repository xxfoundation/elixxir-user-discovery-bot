package twilio

import (
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

// Unit test for registerfact using mock twilio and db interfaces
func TestRegisterFact(t *testing.T) {
	mockDb, _, err := storage.newDatabase("", "", "", "", "11")
	if err != nil {
		t.Errorf("Failed to init mock db: %+v", err)
	}
	uid := id.NewIdFromString("zezima", id.User, t)
	_, err = RegisterFact(uid, "water is wet", 0, []byte("hancock"), MV, mockDb)
	if err != nil {
		t.Errorf("Failed to register fact: %+v", err)
	}
}

// unit test for confirmfact using mock twilio and db interfaces
func TestConfirmFact(t *testing.T) {
	mockDb, _, err := storage.newDatabase("", "", "", "", "11")
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
