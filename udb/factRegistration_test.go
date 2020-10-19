package udb

import (
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/primitives/id"
	"math/rand"
	"testing"
)

func TestRegisterFact(t *testing.T) {
	store, _, err := storage.NewDatabase("", "", "", "", "")
	if err != nil {
		t.Fatalf("Failed to initialize new database with map backend: %+v", err)
	}

	auth := connect.Auth{
		IsAuthenticated: true,
		Sender:          nil,
		Reason:          "",
	}

	user := &storage.User{
		Id:        id.NewIdFromUInt(rand.Uint64(), id.User, t).Bytes(),
		RsaPub:    "",
		DhPub:     nil,
		Salt:      nil,
		Signature: nil,
		Facts:     nil,
	}

	store.InsertUser(user)

}
