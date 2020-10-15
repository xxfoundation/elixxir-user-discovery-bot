////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package udb

import (
	"fmt"
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/client/globals"
	"gitlab.com/elixxir/client/user"
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	// Initialize client with ram storage
	client, err := api.NewClient(&globals.RamStorage{}, "", "", def)
	if err != nil {
		t.Fatalf("Failed to initialize UDB client: %s", err.Error())
	}

	err = client.InitNetwork()

	if err != nil {
		t.Errorf("Conneting to remotes failed: %+v", err)
	}

	err = client.GenerateKeys(nil, "")
	if err != nil {
		t.Errorf("GenerateKeys failed: %s", err.Error())
	}

	// Register with UDB registration code
	_, err = client.RegisterWithPermissioning(true, user.RegistrationCode(&id.UDB))
	if err != nil {
		t.Errorf("Register failed: %s", err.Error())
	}

	// Login to gateway
	_, err = client.Login("")

	if err != nil {
		t.Errorf("Login failed: %s", err.Error())
	}

	//store, _, _ := storage.NewDatabase("", "", "", "", "")

	fmt.Printf("ndf: %v\n",  client.GetNDF())
	//privKeyPem, _ := utils.ReadFile(testkeys.GetNodeKeyPath())
	//permPrivKey, _ := rsa.LoadPrivateKeyFromPem(privKeyPem)
	//h, _ := hash.NewCMixHash()
	//h.Write([]byte("testData"))
	//rsa.Sign(rand.Reader, permPrivKey, hash.CMixHash, )
	//
	//registerMsg := &pb.UDBUserRegistration{
	//	PermissioningSignature: nil,
	//	RSAPublicPem:           "",
	//	IdentityRegistration:   nil,
	//	IdentitySignature:      nil,
	//	Frs:                    nil,
	//	UID:                    nil,
	//	XXX_NoUnkeyedLiteral:   struct{}{},
	//	XXX_unrecognized:       nil,
	//	XXX_sizecache:          0,
	//}
	//
	//RegisterUser(registerMsg, client, store)

	

}