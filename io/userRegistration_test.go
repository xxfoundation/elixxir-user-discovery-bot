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
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/elixxir/crypto/factID"
	"gitlab.com/elixxir/crypto/hash"
	"gitlab.com/elixxir/crypto/registration"
	"gitlab.com/elixxir/primitives/fact"
	"gitlab.com/elixxir/user-discovery-bot/interfaces/params"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/crypto/tls"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/ndf"
	"gitlab.com/xx_network/primitives/utils"
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

func getNDF() []byte {
	cert, _ := utils.ReadFile(testkeys.GetGatewayCertPath())
	addr := "0.0.0.0:4321"

	ndfObj := &ndf.NetworkDefinition{
		Registration: ndf.Registration{
			Address:        addr,
			TlsCertificate: string(cert),
		},
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

	ndfBytes, _ := ndfObj.Marshal()
	return ndfBytes
}

// Happy path test
func TestRegisterUser(t *testing.T) {
	// Initialize client and storage
	clientId, clientKey := initClientFields(t)
	store, _, _ := storage.NewStorage(params.Database{})
	ndfObj, _ := ndf.Unmarshal(getNDF())

	// Create a mock host
	p := connect.GetDefaultHostParams()
	p.MaxRetries = 0
	fakeHost, err := connect.NewHost(clientId, "", nil, p)
	if err != nil {
		t.Errorf("Failed to create fakeHost, %s", err)
	}

	// Construct mock auth object
	auth := &connect.Auth{
		IsAuthenticated: true,
		Sender:          fakeHost,
	}

	cert, err := loadPermissioningPubKey(ndfObj.Registration.TlsCertificate)
	if err != nil {
		t.Errorf(err.Error())
	}

	testTime, err := time.Parse(time.RFC3339,
		"2012-12-21T22:08:41+00:00")
	if err != nil {
		t.Fatalf("Could not parse precanned time: %v", err.Error())
	}

	// Set an invalid identity signature, check that error occurred
	registerMsg, err := buildUserRegistrationMessage(clientId, clientKey, testTime, t)
	if err != nil {
		t.FailNow()
	}

	_, err = registerUser(registerMsg, cert, store, auth)
	if err != nil {
		t.Errorf("Failed happy path: %v", err)
	}

	// Grab the inserted user from database
	retrievedUser, err := store.GetUser(clientId.Bytes())
	if err != nil {
		t.Errorf("Failed to get user from storage: %v", err)
		t.FailNow()
	}

	tf, err := fact.NewFact(fact.FactType(registerMsg.Frs.Fact.FactType), registerMsg.Frs.Fact.Fact)
	if err != nil {
		t.Errorf(" failed to hash fact: %+v", err)
		t.FailNow()
	}
	// Create the expected fact object
	f := storage.Fact{
		Hash:      factID.Fingerprint(tf),
		UserId:    registerMsg.UID,
		Fact:      registerMsg.Frs.Fact.Fact,
		Type:      uint8(registerMsg.Frs.Fact.FactType),
		Signature: registerMsg.Frs.FactSig,
		Verified:  true,
		Timestamp: time.Now(),
	}

	// Create the expected user
	expectedUser := &storage.User{
		Id:                    registerMsg.UID,
		RsaPub:                registerMsg.RSAPublicPem,
		DhPub:                 registerMsg.IdentityRegistration.DhPubKey,
		Salt:                  registerMsg.IdentityRegistration.Salt,
		Signature:             registerMsg.PermissioningSignature,
		Facts:                 []storage.Fact{f},
		RegistrationTimestamp: testTime,
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
	clientId, clientKey := initClientFields(t)
	store, _, _ := storage.NewStorage(params.Database{})
	ndfObj, _ := ndf.Unmarshal(getNDF())

	// Create a mock host
	p := connect.GetDefaultHostParams()
	p.MaxRetries = 0
	fakeHost, err := connect.NewHost(clientId, "", nil, p)
	if err != nil {
		t.Errorf("Failed to create fakeHost, %s", err)
	}

	// Construct mock auth object
	auth := &connect.Auth{
		IsAuthenticated: true,
		Sender:          fakeHost,
	}
	cert, err := loadPermissioningPubKey(ndfObj.Registration.TlsCertificate)
	if err != nil {
		t.Errorf(err.Error())
	}

	testTime, err := time.Parse(time.RFC3339,
		"2012-12-21T22:08:41+00:00")
	if err != nil {
		t.Fatalf("Could not parse precanned time: %v", err.Error())
	}

	// Set an invalid identity signature, check that error occurred
	registerMsg, err := buildUserRegistrationMessage(clientId, clientKey, testTime, t)
	if err != nil {
		t.FailNow()
	}
	registerMsg.IdentitySignature = []byte("invalid")
	_, err = registerUser(registerMsg, cert, store, auth)
	if err == nil {
		t.Errorf("Should not be able to verify identity signature: %v", err)
	}

	// Set invalid fact registration signature, check that error occurred
	registerMsg, err = buildUserRegistrationMessage(clientId, clientKey, testTime, t)
	if err != nil {
		t.FailNow()
	}
	registerMsg.Frs.FactSig = []byte("invalid")
	_, err = registerUser(registerMsg, cert, store, auth)
	if err == nil {
		t.Errorf("Should not be able to verify fact signature: %v", err)
	}

	// Set invalid permissioning signature, check that error occurred
	registerMsg, err = buildUserRegistrationMessage(clientId, clientKey, testTime, t)
	if err != nil {
		t.FailNow()
	}
	registerMsg.PermissioningSignature = []byte("invalid")
	_, err = registerUser(registerMsg, cert, store, auth)
	if err == nil {
		t.Errorf("Should not be able to verify permissioning signature: %v", err)
	}

}

// Error path: Pass in invalid messages
func TestRegisterUser_InvalidMessage(t *testing.T) {
	// Initialize client and storage
	clientId, clientKey := initClientFields(t)
	store, _, _ := storage.NewStorage(params.Database{})
	ndfObj, _ := ndf.Unmarshal(getNDF())

	// Create a mock host
	p := connect.GetDefaultHostParams()
	p.MaxRetries = 0
	fakeHost, err := connect.NewHost(clientId, "", nil, p)
	if err != nil {
		t.Errorf("Failed to create fakeHost, %s", err)
	}

	// Construct mock auth object
	auth := &connect.Auth{
		IsAuthenticated: true,
		Sender:          fakeHost,
	}

	cert, err := loadPermissioningPubKey(ndfObj.Registration.TlsCertificate)
	if err != nil {
		t.Errorf(err.Error())
	}

	testTime, err := time.Parse(time.RFC3339,
		"2012-12-21T22:08:41+00:00")
	if err != nil {
		t.Fatalf("Could not parse precanned time: %v", err.Error())
	}
	// Set an invalid message, check that error occurred
	registerMsg, err := buildUserRegistrationMessage(clientId, clientKey, testTime, t)
	if err != nil {
		t.FailNow()
	}
	registerMsg = nil
	_, err = registerUser(registerMsg, cert, store, auth)
	if err == nil {
		t.Errorf("Should not be able to handle nil message: %v", err)
	}

	// Set invalid fact registration, check that error occurred
	registerMsg, err = buildUserRegistrationMessage(clientId, clientKey, testTime, t)
	if err != nil {
		t.FailNow()
	}
	registerMsg.Frs = nil
	_, err = registerUser(registerMsg, cert, store, auth)
	if err == nil {
		t.Errorf("Should not be able to handle nil FactRegistration message: %v", err)
	}

	// Set invalid fact, check that error occurred
	registerMsg, err = buildUserRegistrationMessage(clientId, clientKey, testTime, t)
	if err != nil {
		t.FailNow()
	}
	registerMsg.Frs.Fact = nil
	_, err = registerUser(registerMsg, cert, store, auth)
	if err == nil {
		t.Errorf("Should not be able to handle nil Fact message: %v", err)
	}

	// Set invalid identity registration, check that error occurred
	registerMsg, err = buildUserRegistrationMessage(clientId, clientKey, testTime, t)
	if err != nil {
		t.FailNow()
	}
	registerMsg.IdentityRegistration = nil
	_, err = registerUser(registerMsg, cert, store, auth)
	if err == nil {
		t.Errorf("Should not be able to handle nil IdentityRegistration message: %v", err)
	}

}

// Helper function which generates needed client fields for a simulated client for testing
func initClientFields(t *testing.T) (*id.ID, *rsa.PrivateKey) {

	clientID := id.NewIdFromString("zezima", id.User, t)

	// RSA Keygen (4096 bit defaults)
	rsaKey, err := rsa.GenerateKey(rand.Reader, rsa.DefaultRSABitLen)
	if err != nil {
		t.Errorf(err.Error())
	}

	return clientID, rsaKey
}

// Helper function which generates a user registration message
func buildUserRegistrationMessage(clientId *id.ID, clientKey *rsa.PrivateKey,
	registrationTimestamp time.Time, t *testing.T) (*pb.UDBUserRegistration, error) {
	// Pull keys and ID out of client
	clientPubKeyPem := rsa.CreatePublicKeyPem(clientKey.GetPublic())

	// Generate permissioning signature, identity and fact messages
	permSig := generatePermissioningSignature(clientPubKeyPem, registrationTimestamp, t)
	requestedUsername := "newUser123"
	identity, identitySig := buildIdentityMsg(requestedUsername, clientId, clientKey)
	frs, err := buildFactMessage(requestedUsername, clientId, clientKey)
	if err != nil {
		return nil, err
	}

	// Construct the user registration message and return
	registerMsg := &pb.UDBUserRegistration{
		PermissioningSignature: permSig,
		RSAPublicPem:           string(clientPubKeyPem),
		IdentityRegistration:   identity,
		IdentitySignature:      identitySig,
		Frs:                    frs,
		UID:                    clientId.Bytes(),
		Timestamp:              registrationTimestamp.UnixNano(),
	}

	return registerMsg, nil
}

// Helper function which generates the identity message
func buildIdentityMsg(username string, clientID *id.ID, clientKey *rsa.PrivateKey) (*pb.Identity, []byte) {
	// We don't need the key to be actually DH for testing purposes
	dhPubKey := rsa.CreatePublicKeyPem(clientKey.GetPublic())
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
func buildFactMessage(username string, clientId *id.ID, clientKey *rsa.PrivateKey) (*pb.FactRegisterRequest, error) {

	// Build the fact
	f, err := fact.NewFact(fact.Username, username)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to hash fact")
	}

	// Sign the fact
	factSig, _ := rsa.Sign(rand.Reader, clientKey, hash.CMixHash, factID.Fingerprint(f), nil)

	// Build the fact registration request and return
	frs := &pb.FactRegisterRequest{
		UID: clientId.Bytes(),
		Fact: &pb.Fact{
			Fact:     username,
			FactType: 0,
		},
		FactSig: factSig,
	}

	return frs, nil
}

// Helper function which creates a permissioning signature for a client
func generatePermissioningSignature(clientPubKey []byte, regTimestamp time.Time, t *testing.T) []byte {
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

	// Construct a permissioning signature
	permSig, err := registration.SignWithTimestamp(rand.Reader, permPrivKey, regTimestamp.UnixNano(), string(clientPubKey))
	if err != nil {
		t.Errorf("Could not get sign test data: %v", err)
	}

	return permSig

}
