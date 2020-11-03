////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package io

import (
	"bytes"
	"crypto/rand"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/client/globals"
	"gitlab.com/elixxir/client/user"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/elixxir/crypto/hash"
	"gitlab.com/elixxir/primitives/utils"
	"gitlab.com/elixxir/user-discovery-bot/interfaces/params"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/crypto/tls"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/ndf"
	"testing"
	"time"
)

// Loads permissioning public key from the certificate
func loadPermissioningPubKey(cert string) (*rsa.PublicKey, error) {
	permCert, err := tls.LoadCertificate(cert)
	if err != nil {
		return nil, errors.Errorf("Could not decode permissioning tls cert file "+
			"into a tls cert: %v", err)
	}

	return tls.ExtractPublicKey(permCert)
}

func getNDF() *ndf.NetworkDefinition {
	return &ndf.NetworkDefinition{
		E2E: ndf.Group{
			Prime: "E2EE983D031DC1DB6F1A7A67DF0E9A8E5561DB8E8D49413394C049B" +
				"7A8ACCEDC298708F121951D9CF920EC5D146727AA4AE535B0922C688B55B3DD2AE" +
				"DF6C01C94764DAB937935AA83BE36E67760713AB44A6337C20E7861575E745D31F" +
				"8B9E9AD8412118C62A3E2E29DF46B0864D0C951C394A5CBBDC6ADC718DD2A3E041" +
				"023DBB5AB23EBB4742DE9C1687B5B34FA48C3521632C4A530E8FFB1BC51DADDF45" +
				"3B0B2717C2BC6669ED76B4BDD5C9FF558E88F26E5785302BEDBCA23EAC5ACE9209" +
				"6EE8A60642FB61E8F3D24990B8CB12EE448EEF78E184C7242DD161C7738F32BF29" +
				"A841698978825B4111B4BC3E1E198455095958333D776D8B2BEEED3A1A1A221A6E" +
				"37E664A64B83981C46FFDDC1A45E3D5211AAF8BFBC072768C4F50D7D7803D2D4F2" +
				"78DE8014A47323631D7E064DE81C0C6BFA43EF0E6998860F1390B5D3FEACAF1696" +
				"015CB79C3F9C2D93D961120CD0E5F12CBB687EAB045241F96789C38E89D796138E" +
				"6319BE62E35D87B1048CA28BE389B575E994DCA755471584A09EC723742DC35873" +
				"847AEF49F66E43873",
			SmallPrime: "2",
			Generator:  "2",
		},
		CMIX: ndf.Group{
			Prime: "9DB6FB5951B66BB6FE1E140F1D2CE5502374161FD6538DF1648218642F0B5C48" +
				"C8F7A41AADFA187324B87674FA1822B00F1ECF8136943D7C55757264E5A1A44F" +
				"FE012E9936E00C1D3E9310B01C7D179805D3058B2A9F4BB6F9716BFE6117C6B5" +
				"B3CC4D9BE341104AD4A80AD6C94E005F4B993E14F091EB51743BF33050C38DE2" +
				"35567E1B34C3D6A5C0CEAA1A0F368213C3D19843D0B4B09DCB9FC72D39C8DE41" +
				"F1BF14D4BB4563CA28371621CAD3324B6A2D392145BEBFAC748805236F5CA2FE" +
				"92B871CD8F9C36D3292B5509CA8CAA77A2ADFC7BFD77DDA6F71125A7456FEA15" +
				"3E433256A2261C6A06ED3693797E7995FAD5AABBCFBE3EDA2741E375404AE25B",
			SmallPrime: "F2C3119374CE76C9356990B465374A17F23F9ED35089BD969F61C6DDE9998C1F",
			Generator: "5C7FF6B06F8F143FE8288433493E4769C4D988ACE5BE25A0E24809670716C613" +
				"D7B0CEE6932F8FAA7C44D2CB24523DA53FBE4F6EC3595892D1AA58C4328A06C4" +
				"6A15662E7EAA703A1DECF8BBB2D05DBE2EB956C142A338661D10461C0D135472" +
				"085057F3494309FFA73C611F78B32ADBB5740C361C9F35BE90997DB2014E2EF5" +
				"AA61782F52ABEB8BD6432C4DD097BC5423B285DAFB60DC364E8161F4A2A35ACA" +
				"3A10B1C4D203CC76A470A33AFDCBDD92959859ABD8B56E1725252D78EAC66E71" +
				"BA9AE3F1DD2487199874393CD4D832186800654760E1E34C09E4D155179F9EC0" +
				"DC4473F996BDCE6EED1CABED8B6F116F7AD9CF505DF0F998E34AB27514B0FFE7",
		},
	}
}

// Happy path test
func TestRegisterUser(t *testing.T) {
	// Initialize client and storage
	client := initTestClient(t)
	store, _, _ := storage.NewStorage(params.Database{})

	// Create a mock host
	p := connect.GetDefaultHostParams()
	p.MaxRetries = 0
	fakeHost, err := connect.NewHost(client.GetCurrentUser(), "", nil, p)
	if err != nil {
		t.Errorf("Failed to create fakeHost, %s", err)
	}

	// Construct mock auth object
	auth := &connect.Auth{
		IsAuthenticated: true,
		Sender:          fakeHost,
	}

	cert, err := loadPermissioningPubKey(client.GetNDF().Registration.TlsCertificate)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Build user registration message
	registerMsg := buildUserRegistrationMessage(client, t)
	_, err = registerUser(registerMsg, cert, store, auth)
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
	store, _, _ := storage.NewStorage(params.Database{})

	// Create a mock host
	p := connect.GetDefaultHostParams()
	p.MaxRetries = 0
	fakeHost, err := connect.NewHost(client.GetCurrentUser(), "", nil, p)
	if err != nil {
		t.Errorf("Failed to create fakeHost, %s", err)
	}

	// Construct mock auth object
	auth := &connect.Auth{
		IsAuthenticated: true,
		Sender:          fakeHost,
	}
	cert, err := loadPermissioningPubKey(client.GetNDF().Registration.TlsCertificate)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Set an invalid identity signature, check that error occurred
	registerMsg := buildUserRegistrationMessage(client, t)
	registerMsg.IdentitySignature = []byte("invalid")
	_, err = registerUser(registerMsg, cert, store, auth)
	if err == nil {
		t.Errorf("Should not be able to verify identity signature: %v", err)
	}

	// Set invalid fact registration signature, check that error occurred
	registerMsg = buildUserRegistrationMessage(client, t)
	registerMsg.Frs.FactSig = []byte("invalid")
	_, err = registerUser(registerMsg, cert, store, auth)
	if err == nil {
		t.Errorf("Should not be able to verify fact signature: %v", err)
	}

	// Set invalid permissioning signature, check that error occurred
	registerMsg = buildUserRegistrationMessage(client, t)
	registerMsg.PermissioningSignature = []byte("invalid")
	_, err = registerUser(registerMsg, cert, store, auth)
	if err == nil {
		t.Errorf("Should not be able to verify permissioning signature: %v", err)
	}

}

// Error path: Pass in invalid messages
func TestRegisterUser_InvalidMessage(t *testing.T) {
	// Initialize client and storage
	client := initTestClient(t)
	store, _, _ := storage.NewStorage(params.Database{})

	// Create a mock host
	p := connect.GetDefaultHostParams()
	p.MaxRetries = 0
	fakeHost, err := connect.NewHost(client.GetCurrentUser(), "", nil, p)
	if err != nil {
		t.Errorf("Failed to create fakeHost, %s", err)
	}

	// Construct mock auth object
	auth := &connect.Auth{
		IsAuthenticated: true,
		Sender:          fakeHost,
	}

	cert, err := loadPermissioningPubKey(client.GetNDF().Registration.TlsCertificate)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Set an invalid message, check that error occurred
	registerMsg := buildUserRegistrationMessage(client, t)
	registerMsg = nil
	_, err = registerUser(registerMsg, cert, store, auth)
	if err == nil {
		t.Errorf("Should not be able to handle nil message: %v", err)
	}

	// Set invalid fact registration, check that error occurred
	registerMsg = buildUserRegistrationMessage(client, t)
	registerMsg.Frs = nil
	_, err = registerUser(registerMsg, cert, store, auth)
	if err == nil {
		t.Errorf("Should not be able to handle nil FactRegistration message: %v", err)
	}

	// Set invalid fact, check that error occurred
	registerMsg = buildUserRegistrationMessage(client, t)
	registerMsg.Frs.Fact = nil
	_, err = registerUser(registerMsg, cert, store, auth)
	if err == nil {
		t.Errorf("Should not be able to handle nil Fact message: %v", err)
	}

	// Set invalid identity registration, check that error occurred
	registerMsg = buildUserRegistrationMessage(client, t)
	registerMsg.IdentityRegistration = nil
	_, err = registerUser(registerMsg, cert, store, auth)
	if err == nil {
		t.Errorf("Should not be able to handle nil IdentityRegistration message: %v", err)
	}

}

// Helper function which generates a client for testing
func initTestClient(t *testing.T) *api.Client {
	// Initialize client with ram storage
	client, err := api.NewClient(&globals.RamStorage{}, "", "", getNDF())
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
