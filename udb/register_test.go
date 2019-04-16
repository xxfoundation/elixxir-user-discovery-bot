////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package udb

import (
	"encoding/base64"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/client/cmixproto"
	"gitlab.com/elixxir/client/crypto"
	"gitlab.com/elixxir/client/globals"
	"gitlab.com/elixxir/client/parse"
	"gitlab.com/elixxir/comms/gateway"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"os"
	"testing"
)

type DummySender struct{}

var rl = RegisterListener{}
var sl = SearchListener{}
var pl = PushKeyListener{}
var gl = GetKeyListener{}

var GWAddress = "localhost:12345"
var GWHandler = api.TestInterface{LastReceivedMessage: pb.CmixMessage{}}

func (d DummySender) Send(recipientID *id.User, msg string) error {
	// do nothing
	jww.INFO.Printf("DummySender!")
	return nil
}

// Hack around the interface for client to do what we need for testing.
func NewMessage(msg string, msgType cmixproto.Type) *parse.Message {
	// Create the message body and assign its type
	tmp := parse.TypedBody{
		MessageType: int32(msgType),
		Body:      []byte(msg),
	}
	return &parse.Message{
		TypedBody: tmp,
		Sender:    new(id.User).SetUints(&[4]uint64{0, 0, 0, 4}),
		Receiver:  new(id.User).SetUints(&[4]uint64{0, 0, 0, 3}),
	}
}

func TestMain(m *testing.M) {
	UdbSender = DummySender{}
	defer gateway.StartGateway(
		GWAddress, &GWHandler, "", "",
	)()
	jww.SetStdoutThreshold(jww.LevelDebug)
	os.Exit(m.Run())
}

// Push the key then register
// NOTE: The send function defaults to a no-op when client is not set up. I am
//       not sure how I feel about it.
func TestRegisterHappyPath(t *testing.T) {
	DataStore = storage.NewRamStorage()
	pubKeyBits := "S8KXBczy0jins9uS4LgBPt0bkFl8t00MnZmExQ6GcOcu8O7DKgAsNzLU7a" +
		"+gMTbIsS995IL/kuFF8wcBaQJBY23095PMSQ/nMuetzhk9HdXxrGIiKBo3C/n4SClpq4H+PoF9XziEVKua8JxGM2o83KiCK3tNUpaZbAAElkjueY7wuD96h4oaA+WV5Nh87cnIZ+fAG0uLve2LSHZ0FBZb3glOpNAOv7PFWkvN2BO37ztOQCXTJe72Y5ReoYn7nWVNxGUh0ilal+BRuJt1GZ7whOGDRE0IXfURIoK2yjyAnyZJWWMhfGsL5S6iL4aXUs03mc8BHKRq3HRjvTE10l3YFA=="

	pubKey := make([]byte, 256)
	pubKey, _ = base64.StdEncoding.DecodeString(pubKeyBits)

	fingerprint := "8oKh7TYG4KxQcBAymoXPBHSD/uga9pX3Mn/jKhvcD8M="
	msgs := []string{
		"myKeyId " + pubKeyBits,
		"EMAIL rick@elixxir.io " + fingerprint,
		fingerprint,
	}

	msg := NewMessage(msgs[0], cmixproto.Type_UDB_PUSH_KEY)
	pl.Hear(msg, false)
	msg = NewMessage(msgs[1], cmixproto.Type_UDB_REGISTER)
	rl.Hear(msg, false)
	msg = NewMessage(msgs[2], cmixproto.Type_UDB_GET_KEY)
	gl.Hear(msg, false)

	// Assert expected state
	k, ok := DataStore.GetKey(fingerprint)
	if !ok {
		t.Errorf("Could not retrieve key %s", fingerprint)
	}
	for i := range k {
		if k[i] != pubKey[i] {
			t.Errorf("pubKey byte mismatch at %d: %d v %d", i, k[i], pubKey[i])
		}
	}

	u, ok2 := DataStore.GetUserKey(id.NewUserFromUint(4, t))
	if !ok2 {
		t.Errorf("Could not retrieve user key 1!")
	}
	if u != fingerprint {
		t.Errorf("GetUserKey fingerprint mismatch: %s v %s", u, fingerprint)
	}

	ks, ok3 := DataStore.GetKeys("rick@elixxir.io", storage.Email)
	if !ok3 {
		t.Errorf("Could not retrieve by e-mail address!")
	}
	if ks[0] != fingerprint {
		t.Errorf("GetKeys fingerprint mismatch: %v v %s", ks[0], fingerprint)
	}
}

func TestInvalidRegistrationCommands(t *testing.T) {
	DataStore = storage.NewRamStorage()
	msgs := []string{
		"PUSHKEY garbage doiandga daoinaosf adsoifn dsaoifa",
		"REGISTER NOTEMAIL something something",
		"REGISTER EMAIL garbage this is a garbage",
		"REGISTER EMAIL rick@elixxir 8oKh7TYG4KxQcBAymoXPBHSD/uga9pX3Mn/jKh" +
			"vcD8M=",
	}

	msg := NewMessage(msgs[0], cmixproto.Type_UDB_PUSH_KEY)
	pl.Hear(msg, false)

	for i := 1; i < len(msgs); i++ {
		msg = NewMessage(msgs[i], cmixproto.Type_UDB_REGISTER)
		rl.Hear(msg, false)
		_, ok := DataStore.GetKey("8oKh7TYG4KxQcBAymoXPBHSD/uga9pX3Mn/jKh")
		if ok {
			t.Errorf("Data store key 8oKh7TYG4KxQcBAymoXPBHSD/uga9pX3Mn/jKh should" +
				" not exist!")
		}
		_, ok2 := DataStore.GetUserKey(id.NewUserFromUint(1, t))
		if ok2 {
			t.Errorf("Data store user 1 should not exist!")
		}
		_, ok3 := DataStore.GetKeys("rick@elixxir.io", storage.Email)
		if ok3 {
			t.Errorf("Data store value rick@elixxir.io should not exist!")
		}
	}
}

func TestRegisterListeners(t *testing.T) {
	grp := crypto.InitCrypto()
	// Initialize client with ram storage
	client, err := api.NewClient(&globals.RamStorage{}, "")
	if err != nil {
		t.Errorf("Failed to initialize UDB client: %s", err.Error())
	}

	udbID := new(id.User).SetUints(&[4]uint64{0, 0, 0, 3})
	// Register with UDB registration code
	userID, err := client.Register(true,
		udbID.RegistrationCode(), "",
		"", []string{GWAddress}, false, grp)

	if err != nil {
		t.Errorf("Register failed: %s", err.Error())
	}

	// Login to gateway
	_, err = client.Login(userID, "", GWAddress, "")

	if err != nil {
		t.Errorf("Login failed: %s", err.Error())
	}

	// Register Listeners
	RegisterListeners(client)

	err = client.Logout()

	if err != nil {
		t.Errorf("Logout failed: %v", err)
	}
}
