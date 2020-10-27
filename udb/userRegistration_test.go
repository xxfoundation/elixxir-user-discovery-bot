////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package udb

import (
	"bytes"
	"crypto/rand"
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/client/globals"
	"gitlab.com/elixxir/client/user"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/elixxir/crypto/hash"
	"gitlab.com/elixxir/primitives/utils"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"
	"testing"
	"time"
)

// Happy path test
func TestRegisterUser(t *testing.T) {
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

	// Build user registration message
	registerMsg := buildUserRegistrationMessage(client, t)
	_, err = RegisterUser(registerMsg, client, store, auth)
	if err != nil {
		t.Errorf("Error in happy path: %v", err)
	}

	// Grab the inserted user from database
	retrievedUser, err := store.GetUser(client.GetCurrentUser().Bytes())
	if err != nil {
		t.Errorf("Failed to get user from storage: %v", err)
	}

	hashedFact := registerMsg.Frs.Fact.Digest()

	// Create the expected fact object
	f := storage.Fact{
		Hash:      hashedFact,
		UserId:    registerMsg.UID,
		Fact:      registerMsg.Frs.Fact.Fact,
		Type:      uint8(registerMsg.Frs.Fact.FactType),
		Signature: registerMsg.Frs.FactSig,
		Verified:  true,
		Timestamp: time.Now(),
	}

	// Create the expected user
	expectedUser := &storage.User{
		Id:        registerMsg.UID,
		RsaPub:    registerMsg.RSAPublicPem,
		DhPub:     registerMsg.IdentityRegistration.DhPubKey,
		Salt:      registerMsg.IdentityRegistration.Salt,
		Signature: registerMsg.PermissioningSignature,
		Facts:     []storage.Fact{f},
	}

	// Compare the retrieved user and the expected user
	// Cannot use reflect.DeepEquals because the Timestamp field will always differ
	if !bytes.Equal(expectedUser.Id, retrievedUser.Id) || expectedUser.RsaPub != retrievedUser.RsaPub ||
		!bytes.Equal(expectedUser.Salt, retrievedUser.Salt) {
		t.Errorf("Retrieved user did not match expected value!"+
			"\n\tExpected: %v"+
			"\n\tRecieved: %v", expectedUser, retrievedUser)
	}

}

// Error path: Pass in invalid signatures for every
//  signature in registration message
func TestRegisterUser_InvalidSignatures(t *testing.T) {
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

	// Set an invalid identity signature, check that error occurred
	registerMsg := buildUserRegistrationMessage(client, t)
	registerMsg.IdentitySignature = []byte("invalid")
	_, err = RegisterUser(registerMsg, client, store, auth)
	if err == nil {
		t.Errorf("Should not be able to verify identity signature: %v", err)
	}

	// Set invalid fact registration signature, check that error occurred
	registerMsg = buildUserRegistrationMessage(client, t)
	registerMsg.Frs.FactSig = []byte("invalid")
	_, err = RegisterUser(registerMsg, client, store, auth)
	if err == nil {
		t.Errorf("Should not be able to verify fact signature: %v", err)
	}

	// Set invalid permissioning signature, check that error occurred
	registerMsg = buildUserRegistrationMessage(client, t)
	registerMsg.PermissioningSignature = []byte("invalid")
	_, err = RegisterUser(registerMsg, client, store, auth)
	if err == nil {
		t.Errorf("Should not be able to verify permissioning signature: %v", err)
	}

}

// Error path: Pass in invalid messages
func TestRegisterUser_InvalidMessage(t *testing.T) {
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

	// Set an invalid message, check that error occurred
	registerMsg := buildUserRegistrationMessage(client, t)
	registerMsg = nil
	_, err = RegisterUser(registerMsg, client, store, auth)
	if err == nil {
		t.Errorf("Should not be able to handle nil message: %v", err)
	}

	// Set invalid fact registration, check that error occurred
	registerMsg = buildUserRegistrationMessage(client, t)
	registerMsg.Frs = nil
	_, err = RegisterUser(registerMsg, client, store, auth)
	if err == nil {
		t.Errorf("Should not be able to handle nil FactRegistration message: %v", err)
	}

	// Set invalid fact, check that error occurred
	registerMsg = buildUserRegistrationMessage(client, t)
	registerMsg.Frs.Fact = nil
	_, err = RegisterUser(registerMsg, client, store, auth)
	if err == nil {
		t.Errorf("Should not be able to handle nil Fact message: %v", err)
	}

	// Set invalid identity registration, check that error occurred
	registerMsg = buildUserRegistrationMessage(client, t)
	registerMsg.IdentityRegistration = nil
	_, err = RegisterUser(registerMsg, client, store, auth)
	if err == nil {
		t.Errorf("Should not be able to handle nil IdentityRegistration message: %v", err)
	}

}

// Helper function which generates a client for testing
func initTestClient(t *testing.T) *api.Client {
	// Initialize client with ram storage
	client, err := api.NewClient(&globals.RamStorage{}, "", "", def)
	if err != nil {
		t.Fatalf("Failed to initialize UDB client: %s", err.Error())
	}

	// Initialize the client's network
	err = client.InitNetwork()
	if err != nil {
		t.Errorf("Conneting to remotes failed: %+v", err)
	}

	// Generate private keys for client
	err = client.GenerateKeys(nil, "")
	if err != nil {
		t.Errorf("GenerateKeys failed: %s", err.Error())
	}

	// Register with UDB registration code
	_, err = client.RegisterWithPermissioning(true, user.RegistrationCode(&id.UDB))
	if err != nil {
		t.Errorf("Register failed: %s", err.Error())
	}

	// Login
	_, err = client.Login("")
	if err != nil {
		t.Errorf("Login failed: %s", err.Error())
	}

	return client
}

// Helper function which generates a user registration message
func buildUserRegistrationMessage(client *api.Client, t *testing.T) *pb.UDBUserRegistration {
	// Pull keys and ID out of client
	clientKey := client.GetCommManager().Comms.GetPrivateKey()
	clientPubKeyPem := rsa.CreatePublicKeyPem(clientKey.GetPublic())
	clientId := client.GetCurrentUser().Bytes()

	// Generate permissioning signature, identity and fact messages
	permSig := generatePermissioningSignature(clientPubKeyPem, t)
	requestedUsername := "newUser123"
	identity, identitySig := buildIdentityMsg(requestedUsername, client)
	frs := buildFactMessage(requestedUsername, client)

	// Construct the user registration message and return
	registerMsg := &pb.UDBUserRegistration{
		PermissioningSignature: permSig,
		RSAPublicPem:           string(clientPubKeyPem),
		IdentityRegistration:   identity,
		IdentitySignature:      identitySig,
		Frs:                    frs,
		UID:                    clientId,
	}

	return registerMsg
}

// Helper function which generates the identity message
func buildIdentityMsg(username string, client *api.Client) (*pb.Identity, []byte) {
	// Pull keys out of client
	dhPubKey := client.GetSession().GetCMIXDHPublicKey().Bytes()
	clientKey := client.GetCommManager().Comms.GetPrivateKey()

	// Construct the identity message
	identity := &pb.Identity{
		Username: username,
		DhPubKey: dhPubKey,
		Salt:     []byte("testSalt"),
	}

	// Sign the identity signature
	identitySig, _ := rsa.Sign(rand.Reader, clientKey, hash.CMixHash, identity.Digest(), nil)

	return identity, identitySig
}

// Helper function which builds the fact messsage
func buildFactMessage(username string, client *api.Client) *pb.FactRegisterRequest {
	// Pull keys and ID out of client
	clientKey := client.GetCommManager().Comms.GetPrivateKey()
	clientId := client.GetCurrentUser().Bytes()

	// Build the fact
	f := &pb.Fact{
		Fact:     username,
		FactType: uint32(storage.Username),
	}

	// Sign the fact
	factSig, _ := rsa.Sign(rand.Reader, clientKey, hash.CMixHash, f.Digest(), nil)

	// Build the fact registration request and return
	frs := &pb.FactRegisterRequest{
		UID:     clientId,
		Fact:    f,
		FactSig: factSig,
	}

	return frs
}

// Helper function which creates a permissioning signature for a client
func generatePermissioningSignature(clientPubKey []byte, t *testing.T) []byte {
	// Pull the key which matches cert in global ndf (see where def is initialized)
	privKeyPem, err := utils.ReadFile(testkeys.GetGatewayKeyPath())
	if err != nil {
		t.Errorf("Could not get test key: %v", err)
	}

	// Load key into object
	permPrivKey, err := rsa.LoadPrivateKeyFromPem(privKeyPem)
	if err != nil {
		t.Errorf("Could not get test key: %v", err)
	}

	// Construct hash
	h, err := hash.NewCMixHash()
	if err != nil {
		t.Errorf("Could make hash: %v", err)
	}

	// Hash public key
	h.Write(clientPubKey)
	hashed := h.Sum(nil)

	// Construct a permissioning signature
	permSig, err := rsa.Sign(rand.Reader, permPrivKey, hash.CMixHash, hashed, nil)
	if err != nil {
		t.Errorf("Could not get sign test data: %v", err)
	}

	return permSig

}
