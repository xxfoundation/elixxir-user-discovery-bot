package udb

import (
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/elixxir/user-discovery-bot/twilio"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"testing"
)

// Happy path.
func TestRegisterFact(t *testing.T) {
	// Initialize client and storage
	client := initTestClient(t)
	store, _, _ := storage.NewDatabase("", "", "", "", "")

	// Create a mock host
	params := connect.GetDefaultHostParams()
	params.MaxRetries = 0
	fakeHost, err := connect.NewHost(client.GetCurrentUser(), "", nil, params)
	if err != nil {
		t.Errorf("Failed to create fakeHost, %s", err)
	}

	// Construct mock auth object
	auth := &connect.Auth{
		IsAuthenticated: true,
		Sender:          fakeHost,
	}

	err = store.InsertUser(&storage.User{
		Id:     client.GetCurrentUser().Bytes(),
		RsaPub: string(rsa.CreatePublicKeyPem(client.GetCommManager().Comms.GetPrivateKey().GetPublic())),
	})
	if err != nil {
		t.Fatalf("Failed to insert user: %+v", err)
	}

	request := buildFactMessage("newUser123", client)

	response, err := RegisterFact(request, twilio.MV, store, auth)
	if err != nil {
		t.Errorf("RegisterFact() produced an error: %+v", err)
	}

	if response.ConfirmationID != "0" {
		t.Errorf("RegisterFact() produced incorrect ConfirmationID: %s", response.ConfirmationID)
	}
}
