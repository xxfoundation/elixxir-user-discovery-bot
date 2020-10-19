package udb

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

func TestRegisterFact(t *testing.T) {
	db, _, err := storage.NewDatabase("", "", "", "", "11")
	if err != nil {
		t.Fatalf("Failed to initialize mock database: %+v", err)
	}

	uid := id.NewIdFromString("zezima", id.User, t)

	request := &pb.FactRegisterRequest{
		UID: uid.Bytes(),
		Fact: &pb.Fact{
			Fact:     "Hair is blue",
			FactType: 5,
		},
		FactSig: []byte("test"),
	}

	response, err := RegisterFact(request, db, nil)
	if err != nil {
		t.Errorf("RegisterFact() produced an error: %+v", err)
	}

	t.Logf("response: %+v", response)
}
