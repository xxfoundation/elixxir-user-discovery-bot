package udb

import (
	"encoding/base64"
	"gitlab.com/elixxir/client/cmixproto"
	"gitlab.com/elixxir/primitives/id"
	fingerprint2 "gitlab.com/elixxir/user-discovery-bot/fingerprint"
	"gitlab.com/elixxir/user-discovery-bot/testutil"
	"testing"
	"time"
)

func TestSearchListener_Hear(t *testing.T) {
	mockSender := testutil.NewMockSender("FOUND", "")
	mockDB := testutil.GetMockDatabase("testid", "testval", "testkeyid", "testkey", false, false)
	var listener = SearchListener{
		Sender: mockSender,
		db:     mockDB,
	}
	var senderID = id.NewUserFromUint(5, t)

	msg := NewMessage("EMAIL rick@elixxir.io", cmixproto.Type_UDB_SEARCH, senderID)

	listener.Hear(msg, false)

	time.Sleep(3 * time.Second)

	if !mockSender.Val {
		t.Error("Did not receive search found message")
	}
}

func TestGetKeyListener_Hear(t *testing.T) {
	mockSender := testutil.NewMockSender("GETKEY", "NOTFOUND")
	mockDB := testutil.GetMockDatabase("testid", "testval", "testkeyid", "testkey", false, false)
	pubKeyBits := "S8KXBczy0jins9uS4LgBPt0bkFl8t00MnZmExQ6GcOcu8O7DKgAsNzLU7a" +
		"+gMTbIsS995IL/kuFF8wcBaQJBY23095PMSQ/nMuetzhk9HdXxrGIiKBo3C/n4SClpq4H+PoF9XziEVKua8JxGM2o83KiCK3tNUpaZbAAElkjueY7wuD96h4oaA+WV5Nh87cnIZ+fAG0uLve2LSHZ0FBZb3glOpNAOv7PFWkvN2BO37ztOQCXTJe72Y5ReoYn7nWVNxGUh0ilal+BRuJt1GZ7whOGDRE0IXfURIoK2yjyAnyZJWWMhfGsL5S6iL4aXUs03mc8BHKRq3HRjvTE10l3YFA=="

	pubKey := make([]byte, 256)
	pubKey, _ = base64.StdEncoding.DecodeString(pubKeyBits)

	fingerprint := fingerprint2.Fingerprint(pubKey)
	var listener = GetKeyListener{
		Sender: mockSender,
		db:     mockDB,
	}
	var senderID = id.NewUserFromUint(5, t)

	msg := NewMessage(fingerprint, cmixproto.Type_UDB_GET_KEY, senderID)
	listener.Hear(msg, false)

	time.Sleep(3 * time.Second)

	if !mockSender.Val {
		t.Error("Did not receive getkey found message")
	}
}

func TestPushKeyListener_Hear(t *testing.T) {
	mockSender := testutil.NewMockSender("PUSHKEY COMPLETE", "")
	mockDB := testutil.GetMockDatabase("testid", "testval", "testkeyid", "testkey", false, true)
	pubKeyBits := "S8KXBczy0jins9uS4LgBPt0bkFl8t00MnZmExQ6GcOcu8O7DKgAsNzLU7a" +
		"+gMTbIsS995IL/kuFF8wcBaQJBY23095PMSQ/nMuetzhk9HdXxrGIiKBo3C/n4SClpq4H+PoF9XziEVKua8JxGM2o83KiCK3tNUpaZbAAElkjueY7wuD96h4oaA+WV5Nh87cnIZ+fAG0uLve2LSHZ0FBZb3glOpNAOv7PFWkvN2BO37ztOQCXTJe72Y5ReoYn7nWVNxGUh0ilal+BRuJt1GZ7whOGDRE0IXfURIoK2yjyAnyZJWWMhfGsL5S6iL4aXUs03mc8BHKRq3HRjvTE10l3YFA=="

	var listener = PushKeyListener{
		Sender: mockSender,
		db:     mockDB,
	}
	var senderID = id.NewUserFromUint(5, t)

	msg := NewMessage("myKeyId "+pubKeyBits, cmixproto.Type_UDB_PUSH_KEY, senderID)
	listener.Hear(msg, false)

	time.Sleep(3 * time.Second)

	if !mockSender.Val {
		t.Error("Did not receive pushkey complete message")
	}
}

func TestRegisterListener_Hear(t *testing.T) {
	mockDB := testutil.GetMockDatabase("testid", "testval", "testkeyid", "testkey", true, false)
	mockSender := testutil.NewMockSender("REGISTRATION COMPLETE", "")
	pubKeyBits := "S8KXBczy0jins9uS4LgBPt0bkFl8t00MnZmExQ6GcOcu8O7DKgAsNzLU7a" +
		"+gMTbIsS995IL/kuFF8wcBaQJBY23095PMSQ/nMuetzhk9HdXxrGIiKBo3C/n4SClpq4H+PoF9XziEVKua8JxGM2o83KiCK3tNUpaZbAAElkjueY7wuD96h4oaA+WV5Nh87cnIZ+fAG0uLve2LSHZ0FBZb3glOpNAOv7PFWkvN2BO37ztOQCXTJe72Y5ReoYn7nWVNxGUh0ilal+BRuJt1GZ7whOGDRE0IXfURIoK2yjyAnyZJWWMhfGsL5S6iL4aXUs03mc8BHKRq3HRjvTE10l3YFA=="

	pubKey := make([]byte, 256)
	pubKey, _ = base64.StdEncoding.DecodeString(pubKeyBits)

	fingerprint := fingerprint2.Fingerprint(pubKey)

	var listener = RegisterListener{
		Sender:    mockSender,
		db:        mockDB,
		blacklist: *InitBlackList("./blacklists/bannedNames.txt"),
	}
	var senderID = id.NewUserFromUint(5, t)

	msg := NewMessage("EMAIL rick@elixxir.io "+fingerprint, cmixproto.Type_UDB_REGISTER, senderID)
	listener.Hear(msg, false)

	time.Sleep(3 * time.Second)

	if !mockSender.Val {
		t.Error("Did not receive register message")
	}

}
