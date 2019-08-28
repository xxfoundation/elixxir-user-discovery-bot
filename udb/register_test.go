////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package udb

import (
	"bytes"
	"encoding/base64"
	"fmt"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/client/cmixproto"
	"gitlab.com/elixxir/client/globals"
	"gitlab.com/elixxir/client/parse"
	"gitlab.com/elixxir/comms/gateway"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/primitives/ndf"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"math/rand"
	"os"
	"testing"
	"time"
)

type DummySender struct{}

var rl = RegisterListener{}
var sl = SearchListener{}
var pl = PushKeyListener{}
var gl = GetKeyListener{}

const NumNodes = 1
const NumGWs = NumNodes
const GWsStartPort = 10000

var GWComms [NumGWs]*gateway.GatewayComms

var def *ndf.NetworkDefinition

func TestMain(m *testing.M) {

	UdbSender = DummySender{}
	jww.SetStdoutThreshold(jww.LevelDebug)

	os.Exit(testMainWrapper(m))
}

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
		Body:        []byte(msg),
	}
	return &parse.Message{
		TypedBody: tmp,
		Sender:    id.NewUserFromUints(&[4]uint64{0, 0, 0, 4}),
		Receiver:  id.NewUserFromUints(&[4]uint64{0, 0, 0, 3}),
	}
}

// Push the key then register
// NOTE: The send function defaults to a no-op when client is not set up. I am
//       not sure how I feel about it.
func TestRegisterHappyPath(t *testing.T) {
	//DataStore = storage.NewRamStorage()
	m := &storage.MapImpl{
		Users: make(map[*id.User]*storage.User),
	}
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
	usr := storage.NewUser()
	usr.SetKeyID(fingerprint)
	retrievedUsr, err := m.GetUser(usr)
	if err != nil  {
		t.Errorf("Could not retrieve key %s", fingerprint)
	}

	if bytes.Compare(retrievedUsr.Key, pubKey) != 0 {
		t.Errorf("pubKey byte mismatch: %+v v %+v", retrievedUsr.Key, pubKey)

	}
	usr2ID := id.NewUserFromUint(4, t)
	usr2 := storage.NewUser()
	usr2.SetID(usr2ID.Bytes())
	err = m.UpsertUser(usr2)
	if err != nil {
		t.Errorf("Could not retrieve user key 1!")
	}
	if u != fingerprint {
		t.Errorf("GetUserKey fingerprint mismatch: %s v %s", u, fingerprint)
	}
	usr3 := storage.NewUser()
	usr3.SetValue("rick@elixxir.io")

	retrievedUser, err := storage.UserDiscoveryDb.GetUser(usr3)
	if err != nil {
		t.Errorf("Could not retrieve by e-mail address!")
	}
	if retrievedUser.KeyId != fingerprint {
		t.Errorf("GetKeys fingerprint mismatch: %v v %s", ks[0], fingerprint)
	}
}

func TestInvalidRegistrationCommands(t *testing.T) {
	//DataStore = storage.NewRamStorage()
	m := &storage.MapImpl{
		Users: make(map[*id.User]*storage.User),
	}
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
		usr := storage.NewUser()
		usr.SetKeyID("8oKh7TYG4KxQcBAymoXPBHSD/uga9pX3Mn/jKh")
		_, err := m.GetUser(usr)
		if err==nil {
			t.Errorf("Data store key 8oKh7TYG4KxQcBAymoXPBHSD/uga9pX3Mn/jKh should" +
				" not exist!")
		}
		usr2 := storage.NewUser()
		usr2.SetID(id.NewUserFromUint(1, t).Bytes())

		_, err = m.GetUser(usr2)
		if err == nil  {
			t.Errorf("Data store user 1 should not exist!")
		}
		usr3 := storage.NewUser()
		usr3.SetValue("rick@elixxir.io")
		_, err = m.GetUser(usr3)
			//DataStore.GetKeys("rick@elixxir.io", storage.Email)
		if err == nil {
			t.Errorf("Data store value rick@elixxir.io should not exist!")
		}
	}
}

func TestRegisterListeners(t *testing.T) {

	// Initialize client with ram storage
	client, err := api.NewClient(&globals.RamStorage{}, "", def)
	if err != nil {
		t.Fatalf("Failed to initialize UDB client: %s", err.Error())
	}

	udbID := id.NewUserFromUints(&[4]uint64{0, 0, 0, 3})
	// Register with UDB registration code
	_, err = client.Register(true, udbID.RegistrationCode(),
		"", "", "", nil)

	if err != nil {
		t.Errorf("Register failed: %s", err.Error())
	}

	err = client.Connect()

	if err != nil {
		t.Errorf("Conneting to remotes failed: %+v", err)
	}

	// Login to gateway
	_, err = client.Login("")

	if err != nil {
		t.Errorf("Login failed: %s", err.Error())
	}

	// Register Listeners
	RegisterListeners(client)

	err = client.StartMessageReceiver()

	if err != nil {
		t.Errorf("Could not start message reciever: %v", err)
	}

	err = client.Logout()

	if err != nil {
		t.Errorf("Logout failed: %v", err)
	}
}

// Handles initialization of mock registration server,
// gateways used for registration and gateway used for session
func testMainWrapper(m *testing.M) int {

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	rndPort := int(rng.Uint64() % 10000)

	def = getNDF()

	// Start mock gateways used by registration and defer their shutdown (may not be needed)
	for i := 0; i < NumGWs; i++ {

		gw := ndf.Gateway{
			Address: fmtAddress(GWsStartPort + i + rndPort),
		}

		def.Gateways = append(def.Gateways, gw)

		GWComms[i] = gateway.StartGateway(gw.Address,
			gateway.NewImplementation(), nil, nil)
	}

	for i := 0; i < NumNodes; i++ {
		nIdBytes := make([]byte, id.NodeIdLen)
		nIdBytes[0] = byte(i)
		n := ndf.Node{
			ID: nIdBytes,
		}
		def.Nodes = append(def.Nodes, n)
	}

	//defer testWrapperShutdown()
	return m.Run()
}

/*func testWrapperShutdown() {
	for _, gw := range GWComms {
		gw.Shutdown()
	}
}*/

func fmtAddress(port int) string { return fmt.Sprintf("localhost:%d", port) }

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
