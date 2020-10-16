package twilio

import (
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

func TestRegisterFact(t *testing.T) {
	uid := id.NewIdFromString("zezima", id.User, t)
	_, err := RegisterFact(uid, "water is wet", 0, []byte("hancock"), MV)
	if err != nil {
		t.Errorf("Failed to register fact: %+v", err)
	}
}

func TestConfirmFact(t *testing.T) {
	uid := id.NewIdFromString("zezima", id.User, t)
	confId, err := RegisterFact(uid, "water is wet", 0, []byte("hancock"), MV)
	if err != nil {
		t.Errorf("Failed to register fact: %+v", err)
	}

	_, err = ConfirmFact(confId, 1234, MV)
	if err != nil {
		t.Errorf("Failed to confirm fact: %+v", err)
	}
}
