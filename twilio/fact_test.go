package twilio

import (
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

// Unit test for registerfact using mock twilio and db interfaces
func TestRegisterFact(t *testing.T) {
	mockDb := storage.NewTestDB(t)
	m := Manager{
		storage:  mockDb,
		verifier: newMockVerifier(),
	}
	uid := id.NewIdFromString("zezima", id.User, t)
	_, err := m.RegisterFact(uid, "7813151633US", 2, []byte("hancock"))
	if err != nil {
		t.Errorf("Failed to register fact: %+v", err)
	}
}

// unit test for confirmfact using mock twilio and db interfaces
func TestConfirmFact(t *testing.T) {
	mockDb := storage.NewTestDB(t)
	m := Manager{
		storage:  mockDb,
		verifier: newMockVerifier(),
	}
	uid := id.NewIdFromString("zezima", id.User, t)
	confId, err := m.RegisterFact(uid, "water is wet", 0, []byte("hancock"))
	if err != nil {
		t.Errorf("Failed to register fact: %+v", err)
	}

	_, err = m.ConfirmFact(confId, "01234")
	if err != nil {
		t.Errorf("Failed to confirm fact: %+v", err)
	}
}
