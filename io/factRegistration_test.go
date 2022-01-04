package io

import (
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/elixxir/user-discovery-bot/twilio"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	jww.SetStdoutThreshold(jww.LevelTrace)
	connect.TestingOnlyDisableTLS = true
	os.Exit(m.Run())
}

// Happy path.
func TestRegisterFact(t *testing.T) {
	// Initialize client and storage
	clientId, clientKey := initClientFields(t)
	store := storage.NewTestDB(t)

	err := store.InsertUser(&storage.User{
		Id:     clientId.Bytes(),
		RsaPub: string(rsa.CreatePublicKeyPem(clientKey.GetPublic())),
	})
	if err != nil {
		t.Fatalf("Failed to insert user: %+v", err)
	}

	request, err := buildFactMessage("newUser123", clientId, clientKey)
	if err != nil {
		t.FailNow()
	}

	// test bad request error path
	old := request.Fact
	request.Fact = nil
	_, err = registerFact(request, twilio.NewMockManager(store), store)
	if err == nil || !strings.Contains(err.Error(), invalidFactRegisterRequestError) {
		t.Errorf("registerFact() did not fail with bad request, instead got: %+v", err)
	}
	request.Fact = old

	response, err := registerFact(request, twilio.NewMockManager(store), store)
	if err != nil {
		t.Errorf("registerFact() produced an error: %+v", err)
	}

	if response.ConfirmationID != "0" {
		t.Errorf("registerFact() produced incorrect ConfirmationID: %s", response.ConfirmationID)
	}

	expectedResponse := &pb.FactRegisterResponse{
		ConfirmationID: "0",
	}

	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("registerFact() produced incorrect FactRegisterRequest."+
			"\n\texpected: %+v\n\treceived: %+v", *expectedResponse, *response)
	}
}

// Error path: Invalid fact signature.
func TestRegisterFact_BadSigError(t *testing.T) {
	// Initialize client and storage
	clientId, clientKey := initClientFields(t)
	store := storage.NewTestDB(t)

	err := store.InsertUser(&storage.User{
		Id:     clientId.Bytes(),
		RsaPub: string(rsa.CreatePublicKeyPem(clientKey.GetPublic())),
	})
	if err != nil {
		t.Fatalf("Failed to insert user: %+v", err)
	}

	request, err := buildFactMessage("newUser123", clientId, clientKey)
	if err != nil {
		t.FailNow()
	}
	request.FactSig = []byte("Bad signature")

	response, err := registerFact(request, twilio.NewMockManager(store), store)
	if err == nil || !strings.Contains(err.Error(), invalidFactSigError) {
		t.Errorf("registerFact() did not produce an error for invalid signature."+
			"\n\texpected: %v\n\treceived: %v", invalidFactSigError, err)
	}

	if !reflect.DeepEqual(response, &pb.FactRegisterResponse{}) {
		t.Errorf("registerFact() produced incorrect FactRegisterRequest."+
			"\n\texpected: %+v\n\treceived: %+v", pb.FactRegisterResponse{}, *response)
	}
}

// Happy path.
func TestConfirmFact(t *testing.T) {
	jww.SetLogThreshold(jww.LevelTrace)
	jww.SetStdoutThreshold(jww.LevelTrace)
	// Initialize client and storage
	clientId, clientKey := initClientFields(t)
	store := storage.NewTestDB(t)

	err := store.InsertUser(&storage.User{
		Id:     clientId.Bytes(),
		RsaPub: string(rsa.CreatePublicKeyPem(clientKey.GetPublic())),
	})
	if err != nil {
		t.Fatalf("Failed to insert user: %+v", err)
	}
	manager := twilio.NewMockManager(store)
	req, err := buildFactMessage("newUser123", clientId, clientKey)
	if err != nil {
		t.FailNow()
	}
	response, err := registerFact(req, manager, store)
	if err != nil {
		t.Fatalf("registerFact() produced an error: %+v", err)
	}

	request := &pb.FactConfirmRequest{
		ConfirmationID: response.ConfirmationID,
		Code:           "420",
	}

	_, err = confirmFact(request, manager)
	if err != nil {
		t.Errorf("confirmFact() produced an error: %+v", err)
	}

}

// Error path: Twilio fails to verify fact because of invalid confirmation ID
// and code.
func TestConfirmFact_FailedVerification(t *testing.T) {
	// Initialize client and storage
	clientId, clientKey := initClientFields(t)
	store := storage.NewTestDB(t)

	err := store.InsertUser(&storage.User{
		Id:     clientId.Bytes(),
		RsaPub: string(rsa.CreatePublicKeyPem(clientKey.GetPublic())),
	})
	if err != nil {
		t.Fatalf("Failed to insert user: %+v", err)
	}
	req, err := buildFactMessage("newUser123", clientId, clientKey)
	if err != nil {
		t.FailNow()
	}
	manager := twilio.NewMockManager(store)
	_, err = registerFact(req, manager, store)
	if err != nil {
		t.Fatalf("registerFact() produced an error: %+v", err)
	}

	request := &pb.FactConfirmRequest{
		ConfirmationID: "5",
		Code:           "3343",
	}

	_, err = confirmFact(request, manager)
	if err == nil || err.Error() != twilioVerificationFailureError {
		t.Errorf("confirmFact() did not produce an error for invalid ConfirmationID and Code."+
			"\n\texpected: %v\n\treceived: %v", twilioVerificationFailureError, err)
	}

}
